package roundrobin

import (
	"fmt"
	. "github/NeichS/simu/internal/structs"
)

func updateAllCounters(tiempo int) {

	unidadesDeTiempo = unidadesDeTiempo + tiempo

	for _, proceso := range listaProcesosListos {
		proceso.PCB.TiempoEnListo += tiempo
	}

	for _, proceso := range listaProcesosBloqueados {
		proceso.PCB.TiempoRafagaIOEmitido += tiempo
	}

	quantumUsage += tiempo
}

var procesoEjecutando *Process

var colaProcesosListos Queue
var listaProcesosListos []*Process

var listaProcesosBloqueados []*Process

var listaProcesosTerminados []*Process
var unidadesDeTiempo int

var quantumUsage int

func StartRoundRobin(procesosNuevos []*Process, procesosTotales, tip, tfp, tcp, quantum int) error {

	cantidadProcesosTerminados := 0

	colaProcesosListos = *NewQueue()
	unidadesDeTiempo = 0

	red := "\033[31m"
	reset := "\033[0m"

	quantumUsage = 0

	for cantidadProcesosTerminados < procesosTotales{

		if procesoEjecutando != nil {
			//corriendo a terminado
			if procesoEjecutando.PCB.RafagasCompletadas == procesoEjecutando.BurstNeeded {

				updateAllCounters(tfp)
				cantidadProcesosTerminados++
				procesoEjecutando.State = "finished"
				fmt.Printf("Tiempo %d: El proceso %s finalizo su ejecucion\n", unidadesDeTiempo, procesoEjecutando.PID)
				listaProcesosTerminados = append(listaProcesosTerminados, procesoEjecutando)
				procesoEjecutando = nil
				quantumUsage = 0
			}
			//corriendo a bloqueado unicamente por I/O
			if procesoEjecutando != nil {
				if procesoEjecutando.PCB.TiempoRafagaEmitido == procesoEjecutando.BurstDuration {

					fmt.Printf("Tiempo %d: Se atendio una interrupcion de I/O del proceso %s \n", unidadesDeTiempo, procesoEjecutando.PID)
					procesoEjecutando.State = "blocked"
					listaProcesosBloqueados = append(listaProcesosBloqueados, procesoEjecutando)
					procesoEjecutando = nil
					quantumUsage = 0
				}
			}
			//corriendo a listo 
			if quantumUsage == quantum {
				listaProcesosListos = append(listaProcesosListos, procesoEjecutando)
				colaProcesosListos.Enqueue(procesoEjecutando)
				fmt.Printf("Tiempo %d: El proceso %s uso todo el quantum\n", unidadesDeTiempo, procesoEjecutando.PID)
				procesoEjecutando = nil
				quantumUsage = 0
			}
		}

		//bloqueado a listo sucede instantaneamente
		for _, element := range listaProcesosBloqueados {
			if element.IOBurstDuration <= element.PCB.TiempoRafagaIOEmitido {
				listaProcesosBloqueados = remove(listaProcesosBloqueados, *element)
				element.PCB.RafagasCompletadas++
				element.PCB.TiempoRafagaEmitido = 0
				element.PCB.TiempoRafagaIOEmitido = 0
				fmt.Printf("Tiempo %d: El proceso %s finalizo su operacion de I/O\n", unidadesDeTiempo, element.PID)
				colaProcesosListos.Enqueue(element)
				listaProcesosListos = append(listaProcesosListos, element)
			}
		}

		//nuevo a listo
		atleastone := false
		for i := len(procesosNuevos) - 1; i >= 0; i-- {
			if unidadesDeTiempo >= procesosNuevos[i].ArrivalTime {
				atleastone = true
				fmt.Printf("Tiempo %d: El proceso %s llega al sistema\n", unidadesDeTiempo, procesosNuevos[i].PID)
				procesosNuevos[i].State = "ready"
				listaProcesosListos = append(listaProcesosListos, procesosNuevos[i]) //falta considerar tip
				colaProcesosListos.Enqueue(procesosNuevos[i])
				procesosNuevos = remove(procesosNuevos, *procesosNuevos[i])
			}
		}
		if atleastone {
			updateAllCounters(tip)
		}

		//listo a corriendo
		if procesoEjecutando == nil && !colaProcesosListos.IsEmpty() {
			quantumUsage = 0
			procesoEjecutando = colaProcesosListos.Dequeue()
			listaProcesosListos = remove(listaProcesosListos, *procesoEjecutando)
			procesoEjecutando.State = "running"
			updateAllCounters(tcp)
			fmt.Printf("Tiempo %d: El proceso %s fue despachado\n", unidadesDeTiempo, procesoEjecutando.PID)
		}

		if procesoEjecutando != nil {
			procesoEjecutando.PCB.TiempoRafagaEmitido++ //recibe su cuota de cpu
		}
		updateAllCounters(1)

		fmt.Printf("%sTiempo %d: Procesos finalizados %d %s \n", string(red), unidadesDeTiempo, cantidadProcesosTerminados, string(reset))
	}

	for _, element := range listaProcesosTerminados {

		fmt.Printf("Descripcion del PID: %s\n", element.PID)
		fmt.Printf("Tiempo en estado listo: %d\n", element.PCB.TiempoEnListo)
	}
	
	return nil
}

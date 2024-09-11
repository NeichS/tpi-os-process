package roundrobin

import (
	"fmt"
	. "github/NeichS/simu/internal/structs"
	s "github/NeichS/simu/cmd/scheduling"
)

func updateAllCounters(tiempo int, so ...string) {

	unidadesDeTiempo = unidadesDeTiempo + tiempo

	for _, proceso := range listaProcesosListos {
		proceso.PCB.TiempoEnListo += tiempo
	}

	for _, proceso := range listaProcesosBloqueados {
		proceso.PCB.TiempoRafagaIOEmitido += tiempo
	}

	quantumUsage += tiempo

	if len(so) == 0 {
		tiempoSO++
	}
}

var procesoEjecutando *Process

var colaProcesosListos Queue
var listaProcesosListos []*Process

var listaProcesosBloqueados []*Process

var listaProcesosTerminados []*Process
var unidadesDeTiempo int
var tiempoSO int
var quantumUsage int

func StartRoundRobin(procesosNuevos []*Process, procesosTotales, tip, tfp, tcp, quantum int) []string {

	cantidadProcesosTerminados := 0

	colaProcesosListos = *NewQueue()
	unidadesDeTiempo = 0

	tiempoPrimerProceso := -1
	tiempoSO = 0

	var logs []string
	quantumUsage = 0

	for cantidadProcesosTerminados < procesosTotales{

		if procesoEjecutando != nil {
			//corriendo a terminado
			if procesoEjecutando.PCB.RafagasCompletadas == procesoEjecutando.BurstNeeded {

				updateAllCounters(tfp)
				cantidadProcesosTerminados++
				procesoEjecutando.State = "finished"
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s finalizo su ejecucion\n", unidadesDeTiempo, procesoEjecutando.PID) )
				procesoEjecutando.TiempoRetorno = unidadesDeTiempo
				listaProcesosTerminados = append(listaProcesosTerminados, procesoEjecutando)
				procesoEjecutando = nil
				quantumUsage = 0
			}
			//corriendo a bloqueado unicamente por I/O
			if procesoEjecutando != nil {
				if procesoEjecutando.PCB.TiempoRafagaEmitido == procesoEjecutando.BurstDuration {

					logs = append(logs, fmt.Sprintf("Tiempo %d: Se atendio una interrupcion de I/O del proceso %s \n", unidadesDeTiempo, procesoEjecutando.PID))
					procesoEjecutando.State = "blocked"
					listaProcesosBloqueados = append(listaProcesosBloqueados, procesoEjecutando)
					procesoEjecutando = nil
					quantumUsage = 0
				}
			}
			//corriendo a listo 
			if quantumUsage == quantum && procesoEjecutando != nil {
				listaProcesosListos = append(listaProcesosListos, procesoEjecutando)
				colaProcesosListos.Enqueue(procesoEjecutando)
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s uso todo el quantum\n", unidadesDeTiempo, procesoEjecutando.PID) )
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
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s finalizo su operacion de I/O\n", unidadesDeTiempo, element.PID) )
				colaProcesosListos.Enqueue(element)
				listaProcesosListos = append(listaProcesosListos, element)
			}
		}

		//nuevo a listo
		atleastone := false
		for i := len(procesosNuevos) - 1; i >= 0; i-- {
			if unidadesDeTiempo >= procesosNuevos[i].ArrivalTime {
				if tiempoPrimerProceso == -1 {
					tiempoPrimerProceso = procesosNuevos[i].ArrivalTime 
				}
				atleastone = true
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s llega al sistema\n", unidadesDeTiempo, procesosNuevos[i].PID))
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
			logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s fue despachado\n", unidadesDeTiempo, procesoEjecutando.PID))
		}

		if procesoEjecutando != nil {
			procesoEjecutando.PCB.TiempoRafagaEmitido++ //recibe su cuota de cpu
			updateAllCounters(1, "tiempo que no usa el SO")
		} else {
			updateAllCounters(1, "nadie usa el cpu") 
		}

	}

	s.ImprimirResultados(listaProcesosTerminados, unidadesDeTiempo, tiempoPrimerProceso, procesosTotales, tiempoSO)
	return logs
}

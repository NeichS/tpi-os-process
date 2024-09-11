package spn

import (
	. "github/NeichS/simu/internal/structs"
	"fmt"
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

func StartSPN(procesosNuevos []*Process, procesosTotales, tip, tfp, tcp int) error {

	cantidadProcesosTerminados := 0

	colaProcesosListos = *NewQueue()
	unidadesDeTiempo = 0

	tiempoPrimerProceso := -1
	tiempoSO = 0

	for cantidadProcesosTerminados < procesosTotales{

		if procesoEjecutando != nil {
			//corriendo a terminado
			if procesoEjecutando.PCB.RafagasCompletadas == procesoEjecutando.BurstNeeded {

				updateAllCounters(tfp)
				cantidadProcesosTerminados++
				procesoEjecutando.State = "finished"
				fmt.Printf("Tiempo %d: El proceso %s finalizo su ejecucion\n", unidadesDeTiempo, procesoEjecutando.PID)
				procesoEjecutando.TiempoRetorno = unidadesDeTiempo
				listaProcesosTerminados = append(listaProcesosTerminados, procesoEjecutando)
				procesoEjecutando = nil
			}
			//corriendo a bloqueado unicamente por I/O
			if procesoEjecutando != nil {
				if procesoEjecutando.PCB.TiempoRafagaEmitido == procesoEjecutando.BurstDuration {

					fmt.Printf("Tiempo %d: Se atendio una interrupcion de I/O del proceso %s \n", unidadesDeTiempo, procesoEjecutando.PID)
					procesoEjecutando.State = "blocked"
					listaProcesosBloqueados = append(listaProcesosBloqueados, procesoEjecutando)
					procesoEjecutando = nil
				}
			}
			//corriendo a listo no hay interrupciones debido a que es no preemptive
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
				colaProcesosListos.Sort("burstDuration")
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
				fmt.Printf("Tiempo %d: El proceso %s llega al sistema\n", unidadesDeTiempo, procesosNuevos[i].PID)
				procesosNuevos[i].State = "ready"
				listaProcesosListos = append(listaProcesosListos, procesosNuevos[i]) //falta considerar tip
				colaProcesosListos.Enqueue(procesosNuevos[i])
				procesosNuevos = remove(procesosNuevos, *procesosNuevos[i])
			}
		}
		if atleastone {
			colaProcesosListos.Sort("burstDuration")
			updateAllCounters(tip)
		}

		//listo a corriendo
		if procesoEjecutando == nil && !colaProcesosListos.IsEmpty() {
			procesoEjecutando = colaProcesosListos.Dequeue()
			listaProcesosListos = remove(listaProcesosListos, *procesoEjecutando)
			procesoEjecutando.State = "running"
			updateAllCounters(tcp)
			fmt.Printf("Tiempo %d: El proceso %s fue despachado\n", unidadesDeTiempo, procesoEjecutando.PID)
		}

		if procesoEjecutando != nil {
			procesoEjecutando.PCB.TiempoRafagaEmitido++ //recibe su cuota de cpu
			updateAllCounters(1, "tiempo que no usa el SO")
		} else {
			updateAllCounters(1, "nadie usa el cpu") 
		}

	}

	
	s.ImprimirResultados(listaProcesosTerminados, unidadesDeTiempo, tiempoPrimerProceso, procesosTotales, tiempoSO)

	return nil
}

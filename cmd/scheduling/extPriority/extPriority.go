package extpriority

import (
	"fmt"
	s "github/NeichS/simu/cmd/scheduling"
	. "github/NeichS/simu/internal/structs"
)

func updateAllCounters(tiempo int, so ...string) {

	unidadesDeTiempo = unidadesDeTiempo + tiempo

	for _, proceso := range listaProcesosListos {
		proceso.PCB.TiempoEnListo += tiempo
	}

	for _, proceso := range listaProcesosBloqueados {
		proceso.PCB.TiempoRafagaIOEmitido += tiempo
	}

	//siempre le sumo tiempo de uso al SO excepto cuando lo usa un proceso o se desperdicia la rafaga
	if len(so) == 0 {
		tiempoSO++
	} else if so[0] == "desperdicio" {
		desperdicio++
	}
}

var procesoEjecutando *Process

var colaProcesosListos Queue
var listaProcesosListos []*Process
var desperdicio int
var listaProcesosBloqueados []*Process
var listaProcesosTerminados []*Process
var unidadesDeTiempo int
var tiempoSO int

func StartExternalPriority(procesosNuevos []*Process, procesosTotales, tip, tfp, tcp int) []string {

	var logs []string
	cantidadProcesosTerminados := 0
	colaProcesosListos = *NewQueue()
	unidadesDeTiempo = 0
	desperdicio = 0
	// red := "\033[31m"
	// reset := "\033[0m"

	tiempoPrimerProceso := -1
	tiempoSO = 0
	for cantidadProcesosTerminados < procesosTotales {

		if procesoEjecutando != nil {
			//corriendo a terminado
			if procesoEjecutando.PCB.RafagasCompletadas == procesoEjecutando.BurstNeeded {
				updateAllCounters(tfp)
				cantidadProcesosTerminados++
				procesoEjecutando.State = "finished"
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s finalizo su ejecucion\n", unidadesDeTiempo, procesoEjecutando.PID))
				procesoEjecutando.TiempoRetorno = unidadesDeTiempo
				listaProcesosTerminados = append(listaProcesosTerminados, procesoEjecutando)
				procesoEjecutando = nil
				continue
			}
			//corriendo a bloqueado unicamente por I/O
			if procesoEjecutando.PCB.TiempoRafagaEmitido == procesoEjecutando.BurstDuration {

				logs = append(logs, fmt.Sprintf("Tiempo %d: Se atendio una interrupcion de I/O del proceso %s \n", unidadesDeTiempo, procesoEjecutando.PID))
				procesoEjecutando.State = "blocked"
				listaProcesosBloqueados = append(listaProcesosBloqueados, procesoEjecutando)
				procesoEjecutando = nil
				continue
			}

			//corriendo a listo
			if !colaProcesosListos.IsEmpty() {
				if (colaProcesosListos.Peek().ExternalPriority) < (procesoEjecutando.ExternalPriority) {
					procesoEjecutando.State = "ready"
					colaProcesosListos.Enqueue(procesoEjecutando)
					colaProcesosListos.Sort("remaining")
					logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s es interrumpido debido a que hay uno de mayor prioridad en estado listo \n", unidadesDeTiempo, procesoEjecutando.PID))
					listaProcesosListos = append(listaProcesosListos, procesoEjecutando)
					procesoEjecutando = nil
					continue
				}
			}
		}

		//bloqueado a listo sucede instantaneamente
		for _, element := range listaProcesosBloqueados {
			if element.IOBurstDuration <= element.PCB.TiempoRafagaIOEmitido {
				listaProcesosBloqueados = s.Remove(listaProcesosBloqueados, *element)
				element.PCB.RafagasCompletadas++
				if element.PCB.RafagasCompletadas == element.BurstNeeded {
					element.PCB.TiempoRafagaEmitido = element.BurstDuration
					element.PCB.TiempoRafagaIOEmitido = 0
				} else {
					element.PCB.TiempoRafagaEmitido = 0
					element.PCB.TiempoRafagaIOEmitido = 0
				}
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s finalizo su operacion de I/O\n", unidadesDeTiempo, element.PID))
				colaProcesosListos.Enqueue(element)
				listaProcesosListos = append(listaProcesosListos, element)
				colaProcesosListos.Sort()
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
				listaProcesosListos = append(listaProcesosListos, procesosNuevos[i])
				colaProcesosListos.Enqueue(procesosNuevos[i])
				procesosNuevos = s.Remove(procesosNuevos, *procesosNuevos[i])
			}
		}
		if atleastone {
			colaProcesosListos.Sort("remaining")
			updateAllCounters(tip)
		}

		//listo a corriendo
		if procesoEjecutando == nil && !colaProcesosListos.IsEmpty() {
			procesoEjecutando = colaProcesosListos.Dequeue()
			listaProcesosListos = s.Remove(listaProcesosListos, *procesoEjecutando)
			procesoEjecutando.State = "running"
			updateAllCounters(tcp)
			logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s fue despachado\n", unidadesDeTiempo, procesoEjecutando.PID))
		}

		if procesoEjecutando != nil {
			procesoEjecutando.PCB.TiempoRafagaEmitido++ //recibe su cuota de cpu
			updateAllCounters(1, "tiempo que no usa el SO")
		} else {
			updateAllCounters(1, "desperdicio")
		}

	}

	s.ImprimirResultados(listaProcesosTerminados, unidadesDeTiempo, tiempoPrimerProceso, procesosTotales, tiempoSO, desperdicio)

	return logs
}

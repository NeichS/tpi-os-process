package srt

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

	if len(so) == 0 {
		tiempoSO++
	} else if so[0] == "desperdicio" {
		desperdicio++
	}
}

var procesoEjecutando *Process

var colaProcesosListos Queue
var listaProcesosListos []*Process
var listaProcesosBloqueados []*Process

var listaProcesosTerminados []*Process
var unidadesDeTiempo int
var tiempoSO int
var desperdicio int

func StartSRT(procesosNuevos []*Process, procesosTotales, tip, tfp, tcp int) []string {

	cantidadProcesosTerminados := 0

	colaProcesosListos = *NewQueue()
	unidadesDeTiempo = 0
	var logs []string
	tiempoSO = 0
	tiempoPrimerProceso := -1
	desperdicio = 0
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
				if (colaProcesosListos.Peek().GetRemaining()) < (procesoEjecutando.GetRemaining()) {
					logs = append(logs, fmt.Sprintf("Tiempo %d: Se interrumpio al proceso %s debido a que hay un proceso con menor tiempo restante de rafaga en la cola de listo\n", unidadesDeTiempo, procesoEjecutando.PID))
					procesoEjecutando.State = "ready"
					colaProcesosListos.Enqueue(procesoEjecutando)
					colaProcesosListos.Sort("remaining")
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
				colaProcesosListos.Sort("remaining")
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

package srt

import (
	"fmt"
	s "github/NeichS/simu/cmd/scheduling"
	. "github/NeichS/simu/internal/structs"
)

// siempre tiempo = 1
func updateAllCounters(tiempo int, so ...string) {

	unidadesDeTiempo = unidadesDeTiempo + tiempo

	for _, proceso := range listaProcesosListos {
		proceso.PCB.TiempoEnListo += tiempo
	}

	for _, proceso := range listaProcesosBloqueados {
		proceso.PCB.TiempoRafagaIOEmitido += tiempo
		logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s recibio rafaga IO %d/%d\n", unidadesDeTiempo, proceso.PID, proceso.PCB.TiempoRafagaIOEmitido, proceso.IOBurstDuration))
	}

	if len(so) == 0 {
		tiempoSO++
	} else if so[0] == "desperdicio" {
		desperdicio++
	}
}

var (
	procesoEjecutando       *Process
	procesoEjecutandoSO     *Process
	desperdicio             int
	colaProcesosListos      Queue
	listaProcesosListos     []*Process
	listaProcesosBloqueados []*Process
	listaProcesosTerminados []*Process
	listaProcesosSO         []*Process
	colaProcesosT           Queue
	unidadesDeTiempo        int
	tiempoSO                int
	logs                    []string
)

func StartSRT(procesosNuevos []*Process, procesosTotales, TIP, TFP, TCP int) []string {

	cantidadProcesosTerminados := 0
	procesoEjecutando = nil
	procesoEjecutandoSO = nil
	colaProcesosListos = *NewQueue()
	unidadesDeTiempo = 0
	colaProcesosT = *NewQueue() //cola de procesos que ejecutan una operacion de SO (TIP, TCP o TFP)
	tiempoPrimerProceso := -1
	tiempoSO = 0
	for cantidadProcesosTerminados < procesosTotales {
		//listo a corriendo
		if procesoEjecutando == nil && !colaProcesosListos.IsEmpty() {
			procesoEjecutando = colaProcesosListos.Dequeue()
			listaProcesosListos = s.Remove(listaProcesosListos, *procesoEjecutando)
			logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s fue despachado\n", unidadesDeTiempo, procesoEjecutando.PID))
		}

		if procesoEjecutando != nil {
			//corriendo a terminado
			if procesoEjecutando.PCB.RafagasCompletadas == procesoEjecutando.BurstNeeded {
				if procesoEjecutando.PCB.TiempoTFP == TFP {
					cantidadProcesosTerminados++
					listaProcesosTerminados = append(listaProcesosTerminados, procesoEjecutando)
					logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s pasa a estado terminado\n", unidadesDeTiempo, procesoEjecutando.PID))
					procesoEjecutando.TiempoRetorno = unidadesDeTiempo
					procesoEjecutando = nil
				} else if procesoEjecutando.PCB.OperacionSOActual != "TFP" {
					procesoEjecutando.PCB.OperacionSOActual = "TFP"
					//logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s comienza a ejecutar su TFP\n", unidadesDeTiempo, procesoEjecutando.PID))
					colaProcesosT.Enqueue(procesoEjecutando)
					listaProcesosSO = append(listaProcesosSO, procesoEjecutando)
				}
			} else {
				//corriendo a bloqueado unicamente por I/O
				if procesoEjecutando != nil {
					// Verificamos si ha completado su ráfaga de CPU
					if procesoEjecutando.PCB.TiempoRafagaEmitido == procesoEjecutando.BurstDuration {
						// Si ya es tiempo de hacer una operación del SO
						if procesoEjecutando.PCB.TiempoTCP == TCP {
							// Pasa a bloqueado
							listaProcesosBloqueados = append(listaProcesosBloqueados, procesoEjecutando)
							procesoEjecutando.PCB.OperacionSOActual = ""
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s pasa a estado bloqueado\n", unidadesDeTiempo, procesoEjecutando.PID))
							procesoEjecutando.PCB.TiempoTCP = 0
							procesoEjecutando = nil
						} else {
							// Verificar si no está en listaProcesosSO
							if procesoEjecutando.PCB.OperacionSOActual != "TCP" {
								procesoEjecutando.PCB.OperacionSOActual = "TCP"
								colaProcesosT.Enqueue(procesoEjecutando)
								listaProcesosSO = append(listaProcesosSO, procesoEjecutando)
								logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s comienza a ejecutar su TCP por operacion I/O\n", unidadesDeTiempo, procesoEjecutando.PID))
							}
						}

					}
				}
				//corriendo a listo
				if procesoEjecutando != nil && !colaProcesosListos.IsEmpty() {
					if (colaProcesosListos.Peek().GetRemaining()) < (procesoEjecutando.GetRemaining()) {

						if procesoEjecutando.PCB.TiempoTCP == TCP {
							// Pasa a bloqueado
							listaProcesosListos = append(listaProcesosListos, procesoEjecutando)
							colaProcesosListos.Enqueue(procesoEjecutando)
							colaProcesosListos.Sort("remaining")
							procesoEjecutando.PCB.OperacionSOActual = ""
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s pasa a estado listo\n", unidadesDeTiempo, procesoEjecutando.PID))
							procesoEjecutando.PCB.TiempoTCP = 0
							procesoEjecutando = nil
						} else {
							// Verificar si no está en listaProcesosSO
							if procesoEjecutando.PCB.OperacionSOActual != "TCP" {
								procesoEjecutando.PCB.OperacionSOActual = "TCP"
								colaProcesosT.Enqueue(procesoEjecutando)
								listaProcesosSO = append(listaProcesosSO, procesoEjecutando)
								logs = append(logs, fmt.Sprintf("Tiempo %d: Se interrumpio al proceso %s debido a que hay un proceso con menor tiempo restante de rafaga en la cola de listo. Tiempo restante del proceso listo = %d | Tiempo restante proceso en ejecucion = %d. Comienza su TCP\n", unidadesDeTiempo, procesoEjecutando.PID, colaProcesosListos.Peek().GetRemaining(), procesoEjecutando.GetRemaining()))
							}
						}
					}
				}
			}
		}

		//bloqueado a listo sucede instantaneamente
		var procesosParaEliminar []*Process

		// Primero, recorres la lista y manejas los procesos que cambian de estado
		for _, element := range listaProcesosBloqueados {
			if element.IOBurstDuration == element.PCB.TiempoRafagaIOEmitido {
				element.PCB.RafagasCompletadas++
				logs = append(logs, fmt.Sprintf("Tiempo %d: Proceso %s rafagas completadas %d/%d \n", unidadesDeTiempo, element.PID, element.PCB.RafagasCompletadas, element.BurstNeeded))
				if element.PCB.RafagasCompletadas == element.BurstNeeded {
					element.PCB.TiempoRafagaEmitido = element.BurstDuration //si era su ultima rafaga el remaining va a ser 0
				} else {
					element.PCB.TiempoRafagaEmitido = 0
				}
				element.PCB.TiempoRafagaIOEmitido = 0
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s finalizo su operacion de I/O\n", unidadesDeTiempo, element.PID))
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s pasa a estado listo I/O\n", unidadesDeTiempo, element.PID))
				colaProcesosListos.Enqueue(element)
				listaProcesosListos = append(listaProcesosListos, element)
				procesosParaEliminar = append(procesosParaEliminar, element) // Marcar para eliminar
			}
		}
		for _, proceso := range procesosParaEliminar {
			listaProcesosBloqueados = s.Remove(listaProcesosBloqueados, *proceso)
		}

		//nuevo a listo
		var procesosParaEliminarNuevos []*Process // Variable para almacenar los procesos a eliminar

		// Primero recorres la lista y manejas los procesos que cumplen la condición
		for _, element := range procesosNuevos {
			if element.ArrivalTime == unidadesDeTiempo {
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s ingresa al sistema\n", unidadesDeTiempo, element.PID))
				element.PCB.OperacionSOActual = "TIP"
				listaProcesosSO = append(listaProcesosSO, element)
				colaProcesosT.Enqueue(element)                                           // lo mando a ejecutar su TIP
				procesosParaEliminarNuevos = append(procesosParaEliminarNuevos, element) // Marcar para eliminar
			}
		}
		// Luego eliminas los procesos marcados de la lista original
		for _, proceso := range procesosParaEliminarNuevos {
			procesosNuevos = s.Remove(procesosNuevos, *proceso)
		}

		//Pregunto donde uso la rafaga del cpu
		if cantidadProcesosTerminados != procesosTotales {
			if procesoEjecutandoSO != nil {
				switch procesoEjecutandoSO.PCB.OperacionSOActual {
				case "TIP":
					if TIP == 0 {
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TIP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						colaProcesosListos.Enqueue(procesoEjecutandoSO)
						listaProcesosListos = append(listaProcesosListos, procesoEjecutandoSO)
						colaProcesosListos.Sort("remaining")
						procesoEjecutandoSO = nil
					} else {
						procesoEjecutandoSO.PCB.TiempoTIP++
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s ejecutó su TIP %d/%d \n", unidadesDeTiempo, procesoEjecutandoSO.PID, procesoEjecutandoSO.PCB.TiempoTIP, TIP))
						if procesoEjecutandoSO.PCB.TiempoTIP >= TIP {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TIP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
							colaProcesosListos.Enqueue(procesoEjecutandoSO)
							listaProcesosListos = append(listaProcesosListos, procesoEjecutandoSO)
							colaProcesosListos.Sort("remaining")
							procesoEjecutandoSO = nil
						}
						updateAllCounters(1)
					}

				case "TCP":
					if TCP == 0 {
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TCP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						procesoEjecutandoSO = nil
					} else {
						procesoEjecutandoSO.PCB.TiempoTCP++
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s ejecutó su TCP %d/%d \n", unidadesDeTiempo, procesoEjecutandoSO.PID, procesoEjecutandoSO.PCB.TiempoTCP, TCP))
						if procesoEjecutandoSO.PCB.TiempoTCP >= TCP {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TCP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
							procesoEjecutandoSO = nil
						}
						updateAllCounters(1)
					}

				case "TFP":
					if TFP == 0 {
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TFP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						procesoEjecutandoSO = nil
					} else {
						procesoEjecutandoSO.PCB.TiempoTFP++
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s ejecutó su TFP %d/%d \n", unidadesDeTiempo, procesoEjecutandoSO.PID, procesoEjecutandoSO.PCB.TiempoTFP, TFP))
						if procesoEjecutandoSO.PCB.TiempoTFP >= TFP {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TFP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
							procesoEjecutandoSO = nil
						}
						updateAllCounters(1)
					}

				}

			} else if !colaProcesosT.IsEmpty() {
				procesoEjecutandoSO = colaProcesosT.Dequeue()
				s.Remove(listaProcesosSO, *procesoEjecutandoSO)
				switch procesoEjecutandoSO.PCB.OperacionSOActual {
				case "TIP":
					if TIP == 0 {
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TIP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						colaProcesosListos.Enqueue(procesoEjecutandoSO)
						listaProcesosListos = append(listaProcesosListos, procesoEjecutandoSO)
						procesoEjecutandoSO = nil
						colaProcesosListos.Sort("remaining")
					} else {
						if procesoEjecutandoSO.PCB.TiempoTIP == 0 {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s comienza a ejecutar su TIP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						}
						procesoEjecutandoSO.PCB.TiempoTIP++
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s ejecutó su TIP %d/%d \n", unidadesDeTiempo, procesoEjecutandoSO.PID, procesoEjecutandoSO.PCB.TiempoTIP, TIP))
						if procesoEjecutandoSO.PCB.TiempoTIP >= TIP {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TIP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
							colaProcesosListos.Enqueue(procesoEjecutandoSO)
							listaProcesosListos = append(listaProcesosListos, procesoEjecutandoSO)
							procesoEjecutandoSO = nil
							colaProcesosListos.Sort("remaining")
						}
						updateAllCounters(1)
					}

				case "TCP":
					if TCP == 0 {
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TCP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						procesoEjecutandoSO = nil
					} else {
						if procesoEjecutandoSO.PCB.TiempoTCP == 0 {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s comienza a ejecutar su TCP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						}
						procesoEjecutandoSO.PCB.TiempoTCP++
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s ejecutó su TCP %d/%d \n", unidadesDeTiempo, procesoEjecutandoSO.PID, procesoEjecutandoSO.PCB.TiempoTCP, TCP))
						if procesoEjecutandoSO.PCB.TiempoTCP >= TCP {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TCP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
							procesoEjecutandoSO = nil
						}
						updateAllCounters(1)
					}

				case "TFP":
					if TFP == 0 {
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TFP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						procesoEjecutandoSO = nil
					} else {
						if procesoEjecutandoSO.PCB.TiempoTFP == 0 {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s comienza a ejecutar su TFP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
						}
						procesoEjecutandoSO.PCB.TiempoTFP++
						logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s ejecutó su TFP %d/%d \n", unidadesDeTiempo, procesoEjecutandoSO.PID, procesoEjecutandoSO.PCB.TiempoTFP, TFP))
						if procesoEjecutandoSO.PCB.TiempoTFP >= TFP {
							logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s terminó su TFP \n", unidadesDeTiempo, procesoEjecutandoSO.PID))
							procesoEjecutandoSO = nil
						}
						updateAllCounters(1)
					}

				}
			} else if procesoEjecutando != nil && procesoEjecutando.BurstNeeded > procesoEjecutando.PCB.RafagasCompletadas {
				procesoEjecutando.PCB.TiempoRafagaEmitido++
				logs = append(logs, fmt.Sprintf("Tiempo %d: El proceso %s ejecutó rafaga de CPU %d/%d \n", unidadesDeTiempo, procesoEjecutando.PID, procesoEjecutando.PCB.TiempoRafagaEmitido, procesoEjecutando.BurstDuration))

				updateAllCounters(1, "proceso usa cpu")
			} else if !colaProcesosListos.IsEmpty() {
				continue
			} else {
				logs = append(logs, fmt.Sprintf("Tiempo %d: Se desperdicio una rafaga de cpu \n", unidadesDeTiempo))
				updateAllCounters(1, "desperdicio")
			}
		}
	}

	s.ImprimirResultados(listaProcesosTerminados, unidadesDeTiempo, tiempoPrimerProceso, procesosTotales, tiempoSO, desperdicio)
	return logs
}

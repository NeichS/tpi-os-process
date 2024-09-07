package extpriority

import (
	"fmt"
	"github/NeichS/simu/internal/structs"
)

// aumenta todos los contadores excepto los relacionados al proceso que esta usando el cpu
func updateAllCounters(tiempo int) {

	unidadesDeTiempo = unidadesDeTiempo + tiempo

	for _, proceso := range procesosListos {
        proceso.PCB.TiempoEnListo += tiempo
    }
   
    for _, proceso := range procesosBloqueados {
        proceso.PCB.TiempoRafagaIOEmitido += tiempo
    }

	tipLeft = tipLeft + tiempo
	tiempoUsoSO = tiempoUsoSO + tiempo
}

var tiempoUsoSO int
var unidadesDeTiempo int
var procesosListos []*structs.Process
var procesosBloqueados []*structs.Process
var tipLeft int
var procesoEjecutando *structs.Process

func StartExternalPriority(procesos *[]structs.Process, procesosTotales, tip, tfp, tcp int) error {

	procesosTerminados := []structs.Process{}
	procesosBloqueados = []*structs.Process{}
	procesosListos = []*structs.Process{}

	listosPrioridadO := structs.NewQueue()
	listosPrioridad1 := structs.NewQueue()
	listosPrioridad2 := structs.NewQueue()
	listosPrioridad3 := structs.NewQueue()
	listosPrioridad4 := structs.NewQueue()
	listosPrioridad5 := structs.NewQueue()

	cantidadProcesosTerminados := 0
	unidadesDeTiempo = 0

	tipLeft = tip //lo igualo para que inicie rapidamente por primera vez

	red := "\033[31m"
	reset := "\033[0m"
	for cantidadProcesosTerminados < procesosTotales && unidadesDeTiempo < 1000 {

		fmt.Println(string(red), "Cantidad de procesos terminados: ", cantidadProcesosTerminados, string(reset))
		if procesoEjecutando != nil {

			//Corriendo a terminado
			if procesoEjecutando.BurstNeeded == procesoEjecutando.PCB.RafagasCompletadas {
				procesosTerminados = append(procesosTerminados, *procesoEjecutando)
				fmt.Printf("Tiempo %d: Finaliza la ejecucion del proceso %s\n", unidadesDeTiempo, procesoEjecutando.PID)
				cantidadProcesosTerminados++
				procesoEjecutando.State = "Terminado"
				procesoEjecutando = nil
				updateAllCounters(tfp) //sumo el tiempo en el que tarda en finalizar el proceso
			}

			if procesoEjecutando != nil {
				//Corriendo a bloqueado
				if procesoEjecutando.BurstDuration == procesoEjecutando.PCB.TiempoRafagaEmitido && procesoEjecutando.IOBurstDuration > procesoEjecutando.PCB.TiempoRafagaIOEmitido {
					procesosBloqueados = append(procesosBloqueados, procesoEjecutando)
					fmt.Printf("Tiempo %d: Se atiende una interrupcion I/O del proceso %s\n", unidadesDeTiempo, procesoEjecutando.PID)
					procesoEjecutando = nil
				}
			}

			if procesoEjecutando != nil {
				//Corriendo a listo (sucede si aparece un proceso con mayor prioridad)
				if !listosPrioridadO.IsEmpty() && procesoEjecutando.ExternalPriority > 0 {
					fmt.Printf("Tiempo %d: El proceso %s fue interrumpido\n", unidadesDeTiempo, procesoEjecutando.PID)
					switch procesoEjecutando.ExternalPriority {
					case 0:
						listosPrioridadO.Enqueue(procesoEjecutando)
					case 1:
						listosPrioridad1.Enqueue(procesoEjecutando)
					case 2:
						listosPrioridad2.Enqueue(procesoEjecutando)
					case 3:
						listosPrioridad3.Enqueue(procesoEjecutando)
					case 4:
						listosPrioridad4.Enqueue(procesoEjecutando)
					default:
						listosPrioridad5.Enqueue(procesoEjecutando)
					}
					procesoEjecutando.State = "ready"
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad1.IsEmpty() && procesoEjecutando.ExternalPriority > 1 {
					fmt.Printf("Tiempo %d: El proceso %s fue interrumpido\n", unidadesDeTiempo, procesoEjecutando.PID)
					switch procesoEjecutando.ExternalPriority {
					case 0:
						listosPrioridadO.Enqueue(procesoEjecutando)
					case 1:
						listosPrioridad1.Enqueue(procesoEjecutando)
					case 2:
						listosPrioridad2.Enqueue(procesoEjecutando)
					case 3:
						listosPrioridad3.Enqueue(procesoEjecutando)
					case 4:
						listosPrioridad4.Enqueue(procesoEjecutando)
					default:
						listosPrioridad5.Enqueue(procesoEjecutando)
					}
					procesoEjecutando.State = "ready"
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad2.IsEmpty() && procesoEjecutando.ExternalPriority > 2 {
					fmt.Printf("Tiempo %d: El proceso %s fue interrumpido\n", unidadesDeTiempo, procesoEjecutando.PID)
					switch procesoEjecutando.ExternalPriority {
					case 0:
						listosPrioridadO.Enqueue(procesoEjecutando)
					case 1:
						listosPrioridad1.Enqueue(procesoEjecutando)
					case 2:
						listosPrioridad2.Enqueue(procesoEjecutando)
					case 3:
						listosPrioridad3.Enqueue(procesoEjecutando)
					case 4:
						listosPrioridad4.Enqueue(procesoEjecutando)
					default:
						listosPrioridad5.Enqueue(procesoEjecutando)
					}
					procesoEjecutando.State = "ready"
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad3.IsEmpty() && procesoEjecutando.ExternalPriority > 3 {
					fmt.Printf("Tiempo %d: El proceso %s fue interrumpido\n", unidadesDeTiempo, procesoEjecutando.PID)
					switch procesoEjecutando.ExternalPriority {
					case 0:
						listosPrioridadO.Enqueue(procesoEjecutando)
					case 1:
						listosPrioridad1.Enqueue(procesoEjecutando)
					case 2:
						listosPrioridad2.Enqueue(procesoEjecutando)
					case 3:
						listosPrioridad3.Enqueue(procesoEjecutando)
					case 4:
						listosPrioridad4.Enqueue(procesoEjecutando)
					default:
						listosPrioridad5.Enqueue(procesoEjecutando)
					}
					procesoEjecutando.State = "ready"
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad4.IsEmpty() && procesoEjecutando.ExternalPriority > 4 {
					fmt.Printf("Tiempo %d: El proceso %s fue interrumpido\n", unidadesDeTiempo, procesoEjecutando.PID)
					switch procesoEjecutando.ExternalPriority {
					case 0:
						listosPrioridadO.Enqueue(procesoEjecutando)
					case 1:
						listosPrioridad1.Enqueue(procesoEjecutando)
					case 2:
						listosPrioridad2.Enqueue(procesoEjecutando)
					case 3:
						listosPrioridad3.Enqueue(procesoEjecutando)
					case 4:
						listosPrioridad4.Enqueue(procesoEjecutando)
					default:
						listosPrioridad5.Enqueue(procesoEjecutando)
					}
					procesoEjecutando.State = "ready"
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad5.IsEmpty() {
					fmt.Printf("Tiempo %d: El proceso %s fue interrumpido\n", unidadesDeTiempo, procesoEjecutando.PID)
					switch procesoEjecutando.ExternalPriority {
					case 0:
						listosPrioridadO.Enqueue(procesoEjecutando)
					case 1:
						listosPrioridad1.Enqueue(procesoEjecutando)
					case 2:
						listosPrioridad2.Enqueue(procesoEjecutando)
					case 3:
						listosPrioridad3.Enqueue(procesoEjecutando)
					case 4:
						listosPrioridad4.Enqueue(procesoEjecutando)
					default:
						listosPrioridad5.Enqueue(procesoEjecutando)
					}
					procesoEjecutando.State = "ready"
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else {
					fmt.Printf("Tiempo %d: El proceso %s continua con su ejecucion\n", unidadesDeTiempo, procesoEjecutando.PID)
				}
			}
		}

		//bloqueado a listo (asum que no suma tiempo)

		for i := len(procesosBloqueados) - 1; i >= 0; i-- {
			element := procesosBloqueados[i]
			if element.IOBurstDuration <= element.PCB.TiempoRafagaIOEmitido {
				// Remove the element from procesosBloqueados
				procesosBloqueados = append(procesosBloqueados[:i], procesosBloqueados[i+1:]...)

				element.PCB.TiempoRafagaIOEmitido = 0
				element.PCB.RafagasCompletadas++
				element.PCB.TiempoRafagaEmitido = 0

				switch element.ExternalPriority {
				case 0:
					listosPrioridadO.Enqueue(element)
				case 1:
					listosPrioridad1.Enqueue(element)
				case 2:
					listosPrioridad2.Enqueue(element)
				case 3:
					listosPrioridad3.Enqueue(element)
				case 4:
					listosPrioridad4.Enqueue(element)
				default:
					listosPrioridad5.Enqueue(element)
				}
				element.State = "ready"
				procesosListos = append(procesosListos, element)

				fmt.Printf("Tiempo %d: Se completa una rafaga del proceso %s vuelve a la cola de listo %d\n", unidadesDeTiempo, element.PID, element.ExternalPriority)
			} else {
				procesosBloqueados[i].PCB.TiempoRafagaIOEmitido++
			}
		}

		//Nuevo a listo
		if tip <= tipLeft && cantidadProcesosTerminados < procesosTotales {
			fmt.Printf("Tiempo %d: El sistema operativo va a aceptar nuevos procesos \n", unidadesDeTiempo)
			for i := 0; i < len(*procesos); i++ {

				if (*procesos)[i].ArrivalTime <= unidadesDeTiempo && !Contains(procesosTerminados, (*procesos)[i]) && (*procesos)[i].State == "New" {
					(*procesos)[i].State = "ready"
					switch (*procesos)[i].ExternalPriority {
					case 0:
						listosPrioridadO.Enqueue(&(*procesos)[i])
					case 1:
						listosPrioridad1.Enqueue(&(*procesos)[i])
					case 2:
						listosPrioridad2.Enqueue(&(*procesos)[i])
					case 3:
						listosPrioridad3.Enqueue(&(*procesos)[i])
					case 4:
						listosPrioridad4.Enqueue(&(*procesos)[i])
					default:
						listosPrioridad5.Enqueue(&(*procesos)[i])
					}
					procesosListos = append(procesosListos, &(*procesos)[i])
					fmt.Printf("Tiempo %d: El proceso %s ingreso a esperar a la cola %d \n", unidadesDeTiempo, (*procesos)[i].PID, (*procesos)[i].ExternalPriority)
				}
			}
			tipLeft = 0
		} else {
			updateAllCounters(1)
		}

		//Despacho de listo a corriendo aca se aplica la politica.
		if procesoEjecutando == nil {
			if !listosPrioridadO.IsEmpty() {
				procesoEjecutando = listosPrioridadO.Dequeue()
				
				fmt.Printf("Tiempo %d: Nuevo proceso despachado %s\n", unidadesDeTiempo, procesoEjecutando.PID)
			} else if !listosPrioridad1.IsEmpty() {
				procesoEjecutando = listosPrioridad1.Dequeue()
				fmt.Printf("Tiempo %d: Nuevo proceso despachado %s\n", unidadesDeTiempo, procesoEjecutando.PID)
			} else if !listosPrioridad2.IsEmpty() {
				procesoEjecutando = listosPrioridad2.Dequeue()
				fmt.Printf("Tiempo %d: Nuevo proceso despachado %s\n", unidadesDeTiempo, procesoEjecutando.PID)
			} else if !listosPrioridad3.IsEmpty() {
				procesoEjecutando = listosPrioridad3.Dequeue()
				fmt.Printf("Tiempo %d: Nuevo proceso despachado %s\n", unidadesDeTiempo, procesoEjecutando.PID)
			} else if !listosPrioridad4.IsEmpty() {
				procesoEjecutando = listosPrioridad4.Dequeue()
				fmt.Printf("Tiempo %d: Nuevo proceso despachado %s\n", unidadesDeTiempo, procesoEjecutando.PID)
			} else if !listosPrioridad5.IsEmpty() {
				procesoEjecutando = listosPrioridad5.Dequeue()
				fmt.Printf("Tiempo %d: Nuevo proceso despachado %s\n", unidadesDeTiempo, procesoEjecutando.PID)
			} else {
				fmt.Printf("Tiempo %d: No existen prcesos listos\n", unidadesDeTiempo)
			}
			if procesoEjecutando != nil {
				Remove(procesosListos, *procesoEjecutando)
				procesoEjecutando.State = "running"
				updateAllCounters(tcp)
				fmt.Printf("Tiempo %d: Se completo cambio de proceso, ahora %s tiene el control de la CPU \n", unidadesDeTiempo, procesoEjecutando.PID)
			}
		}

		//fin de la vuelta

		if procesoEjecutando != nil {
			fmt.Printf("Tiempo %d: El proceso %s tiene el control de la CPU \n", unidadesDeTiempo, procesoEjecutando.PID)
		}

		//el proceso recibe su dosis de cpu
		if procesoEjecutando != nil {
			fmt.Printf("Tiempo %d: El proceso %s recibe una rafaga de cpu\n", unidadesDeTiempo, procesoEjecutando.PID)
			procesoEjecutando.PCB.TiempoRafagaEmitido++
			updateAllCounters(1)
		}

	}

	fmt.Printf("Procesos terminados\n")
	for _, element := range procesosTerminados {

		fmt.Printf("Descripcion del PID: %s\n", element.PID)
		fmt.Printf("Tiempo en estado listo: %d\n", element.PCB.TiempoEnListo)

	}

	return nil
}


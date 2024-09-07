package scheduling

import (
	"fmt"
	"github/NeichS/simu/internal/structs"
)

// aumenta todos los contadores excepto los relacionados al proceso que esta usando el cpu
func updateAllCounters(tiempo int) {

	unidadesDeTiempo = unidadesDeTiempo + tiempo
	
	for _, element := range procesosListos {
		element.PCB.TiempoEnListo = element.PCB.TiempoEnListo + tiempo 
	}

	for _, element := range procesosBloqueados {
		element.PCB.TiempoRafagaIOEmitido = element.PCB.TiempoRafagaIOEmitido + tiempo 
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
	for cantidadProcesosTerminados < procesosTotales {

		fmt.Println(string(red), "Cantidad de procesos terminados: ", cantidadProcesosTerminados, string(reset))
		if procesoEjecutando != nil {

			//Corriendo a terminado
			if procesoEjecutando.BurstNeeded == procesoEjecutando.PCB.RafagasCompletadas {
				procesosTerminados = append(procesosTerminados, *procesoEjecutando)
				fmt.Printf("Tiempo %d: Finaliza la ejecucion del proceso %s\n", unidadesDeTiempo, procesoEjecutando.PID)
				cantidadProcesosTerminados++
				procesoEjecutando.State = "Terminado"
				procesoEjecutando = nil
				updateAllCounters(tfp)
			}
			
			if procesoEjecutando != nil {
				//Corriendo a bloqueado
				if procesoEjecutando.BurstDuration == procesoEjecutando.PCB.TiempoRafagaEmitido && procesoEjecutando.IOBurstDuration > procesoEjecutando.PCB.TiempoRafagaIOEmitido {
					procesoEjecutando.PCB.TiempoRafagaEmitido = 0

					procesosBloqueados = append(procesosBloqueados, procesoEjecutando)
					procesoEjecutando = nil
					updateAllCounters(1)
					fmt.Printf("Tiempo %d: Se atiende una interrupcion I/O\n", unidadesDeTiempo)
				}
			}

			if procesoEjecutando != nil {
				//Corriendo a listo (sucede si aparece un proceso con mayor prioridad)
				if !listosPrioridadO.IsEmpty() && procesoEjecutando.ExternalPriority > 0 {
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad1.IsEmpty() && procesoEjecutando.ExternalPriority > 1 {
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad2.IsEmpty() && procesoEjecutando.ExternalPriority > 2 {
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad3.IsEmpty() && procesoEjecutando.ExternalPriority > 3 {
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad4.IsEmpty() && procesoEjecutando.ExternalPriority > 4 {
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else if !listosPrioridad5.IsEmpty() {
					procesosListos = append(procesosListos, procesoEjecutando)
					procesoEjecutando = nil
				} else {
					updateAllCounters(1)
					fmt.Printf("Tiempo %d: El proceso %s continua con su ejecucion\n", unidadesDeTiempo, procesoEjecutando.PID)
				}
			}
		}



		for _, element := range procesosBloqueados {
			if element.IOBurstDuration == element.PCB.TiempoRafagaIOEmitido {
				procesosBloqueados = remove(procesosBloqueados, *element)
				element.PCB.TiempoRafagaIOEmitido = 0
				element.PCB.RafagasCompletadas++
				element.PCB.TiempoRafagaEmitido = 0
				fmt.Printf("Tiempo %d: Se completa una rafaga del proceso %s\n", unidadesDeTiempo, procesoEjecutando.PID)
			} else {
				element.PCB.TiempoRafagaIOEmitido++
			}
		}
	
		updateAllCounters(1)

		//Nuevo a listo
		if tip <= tipLeft {
			fmt.Printf("Tiempo %d: El sistema operativo va a aceptar nuevos procesos \n", unidadesDeTiempo)
			for i := 0; i < len(*procesos); i++ {
				if (*procesos)[i].ArrivalTime >= unidadesDeTiempo && !contains(procesosTerminados,(*procesos)[i]) {
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
				fmt.Printf("Nuevo proceso despachado %s\n", procesoEjecutando.PID)
			} else if !listosPrioridad1.IsEmpty() {
				procesoEjecutando = listosPrioridad1.Dequeue()
				fmt.Printf("Nuevo proceso despachado %s\n", procesoEjecutando.PID)
			} else if !listosPrioridad2.IsEmpty() {
				procesoEjecutando = listosPrioridad2.Dequeue()
				fmt.Printf("Nuevo proceso despachado %s\n", procesoEjecutando.PID)
			} else if !listosPrioridad3.IsEmpty() {
				procesoEjecutando = listosPrioridad3.Dequeue()
				fmt.Printf("Nuevo proceso despachado %s\n", procesoEjecutando.PID)
			} else if !listosPrioridad4.IsEmpty() {
				procesoEjecutando = listosPrioridad4.Dequeue()
				fmt.Printf("Nuevo proceso despachado %s\n", procesoEjecutando.PID)
			} else if !listosPrioridad5.IsEmpty() {
				procesoEjecutando = listosPrioridad5.Dequeue()
				fmt.Printf("Nuevo proceso despachado %s\n", procesoEjecutando.PID)
			} else {
				fmt.Printf("Tiempo %d: No existen prcesos listos\n", unidadesDeTiempo)
			}
			if procesoEjecutando != nil {
				unidadesDeTiempo = unidadesDeTiempo + tcp
				fmt.Printf("Tiempo %d: Se completo cambio de proceso, ahora %s tiene el control de la CPU \n", unidadesDeTiempo, procesoEjecutando.PID)
			}
		}

		//fin de la vuelta

		if procesoEjecutando != nil {
			fmt.Printf("Tiempo %d: El proceso %s tiene el control de la CPU \n", unidadesDeTiempo, procesoEjecutando.PID)
		}

		//el proceso recibe su dosis de cpu
		
		procesoEjecutando.PCB.TiempoRafagaEmitido++
	}

	fmt.Printf("Procesos terminados\n");
	for _, element := range procesosTerminados {
		fmt.Printf("%s\n", element.PID)
	}

	return nil
}

func remove(procesos []*structs.Process, element structs.Process) []*structs.Process {
	index := -1
    for i, v := range procesos {
        if v.PID == element.PID {
            index = i
            break
        }
    }

    if index != -1 {
        procesos = append(procesos[:index], procesos[index+1:]...)
    }

	return procesos
}

func contains(slice []structs.Process, item structs.Process) bool {
    for _, element := range slice {
        if element.PID == item.PID {
            return true
        }
    }
    return false
}
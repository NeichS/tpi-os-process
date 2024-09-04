package scheduling

import (
	"fmt"
	"github/NeichS/simu/internal/structs"
)

// Cola de procesos para cada prioridad

func StartExternalPriority(procesos *[]structs.Process, procesosTotales, tip int) error {

	procesosTerminados := structs.NewLinkedList()
	procesosBloqueados := structs.NewLinkedList()

	listosPrioridadO := structs.NewQueue() 
	listosPrioridad1 := structs.NewQueue()
	listosPrioridad2 := structs.NewQueue()
	listosPrioridad3 := structs.NewQueue()
	listosPrioridad4 := structs.NewQueue()
	listosPrioridad5 := structs.NewQueue()

	cantidadProcesosTerminados := 0
	var procesoEjecutando *structs.Process
	unidadesDeTiempo := 0
	tipLeft := tip;

	for cantidadProcesosTerminados < procesosTotales {

		//Corriendo a terminado
		if procesoEjecutando != nil {

			if procesoEjecutando.BurstNeeded == procesoEjecutando.PCB.RafagasCompletadas && procesoEjecutando.IOBurstDuration == procesoEjecutando.PCB.TiempoRafagaIOEmitido {
				procesosTerminados.Append(procesoEjecutando)
				cantidadProcesosTerminados++
				procesoEjecutando = nil
			}

			//Corriendo a bloqueado
			if procesoEjecutando.BurstDuration == procesoEjecutando.PCB.TiempoRafagaEmitido && procesoEjecutando.IOBurstDuration < procesoEjecutando.PCB.TiempoRafagaIOEmitido {
				procesosBloqueados.Append(procesoEjecutando)
			}

			//Corriendo a listo (en esta politica no deberia suceder esto ya que es no preemptive
		}

		//Bloqueado a listo
		var recorrer *structs.Node = procesosBloqueados.Head

		for recorrer != nil {
			if recorrer.Process.IOBurstDuration == recorrer.Process.PCB.TiempoRafagaIOEmitido {
				procesosBloqueados.Remove(recorrer.Process.Name)
			} else {
				recorrer.Process.PCB.TiempoRafagaIOEmitido++
			}

			recorrer = recorrer.Next
		} 

		//Nuevo a listo
		if tip == tipLeft {
			fmt.Printf("Tiempo %d: El sistema operativo va a aceptar nuevos procesos \n", unidadesDeTiempo)
			for i := 0; i < len(*procesos); i++{
				if (*procesos)[i].ArrivalTime == unidadesDeTiempo  {
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

					fmt.Printf("Tiempo %d: El proceso %s ingreso a esperar a la cola %d \n", unidadesDeTiempo, (*procesos)[i].Name, (*procesos)[i].ExternalPriority )
				}
			}
			tipLeft = 0
		} else {
			tipLeft++
		}

		//Despacho de listo a corriendo aca se aplica la politica
		if procesoEjecutando != nil {
			if !listosPrioridadO.IsEmpty() {
				procesoEjecutando = listosPrioridadO.Dequeue()
				
			} else if !listosPrioridad1.IsEmpty() {
				procesoEjecutando = listosPrioridad1.Dequeue()
			} else if !listosPrioridad2.IsEmpty() {
				procesoEjecutando = listosPrioridad2.Dequeue()
			} else if !listosPrioridad3.IsEmpty() {
				procesoEjecutando = listosPrioridad3.Dequeue()
			} else if !listosPrioridad4.IsEmpty() {
				procesoEjecutando = listosPrioridad4.Dequeue()
			} else if !listosPrioridad5.IsEmpty() {
				procesoEjecutando = listosPrioridad5.Dequeue()
			} else {
				fmt.Printf("Tiempo %d: No existen prcesos listos\n", unidadesDeTiempo)
			}
			if procesoEjecutando != nil {
				fmt.Printf("Tiempo %d: El proceso %s tiene el control de la CPU \n", unidadesDeTiempo, procesoEjecutando.Name)
			}
		}

		procesoEjecutando.PCB.TiempoRafagaEmitido++
		if procesoEjecutando.PCB.TiempoRafagaEmitido == procesoEjecutando.BurstDuration {
			procesoEjecutando.PCB.RafagasCompletadas++
			procesoEjecutando.PCB.TiempoRafagaEmitido = 0
		}

		unidadesDeTiempo++
		
	}
	return nil
}

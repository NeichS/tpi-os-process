package scheduling

import (
	"fmt"
	"github/NeichS/simu/internal/structs"
	"sync"
	"time"
)


func startProcessRace(process *structs.Process, colaEspera chan structs.Process,wg *sync.WaitGroup, start chan struct{}) {

	defer wg.Done()
	<-start

	//lo atraso lo que tarde en arrivar 
	time.Sleep(time.Duration(process.ArrivalTime) * time.Millisecond)
	for {
		select {
		case <-colaEspera:
			//obtengo rafaga de cpu 
			fmt.Printf("El proceso %s esta usando la CPU", process.Name)
			time.Sleep(time.Duration(process.BurstDuration))
			return
		default:
			
		}

	}
}

func StartExternalPriority(procesos *[]structs.Process) error {

	//Todavia no se si va a ser un canal sincronico o asincronico
	prioridadCero := make(chan structs.Process) 
	prioridadUno := make(chan structs.Process)
	prioridadDos := make(chan structs.Process)
	prioridadTres := make(chan structs.Process)
	prioridadCuatro := make(chan structs.Process)
	prioridadCinco := make(chan structs.Process)

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0 ; i < len(*procesos) ; i++ {
		wg.Add(1)
		
		switch (*procesos)[i].ExternalPriority {
		case 0:
			go startProcessRace(&(*procesos)[i], prioridadCero ,&wg, start)
		case 1:
			go startProcessRace(&(*procesos)[i], prioridadUno ,&wg, start)
		case 2:
			go startProcessRace(&(*procesos)[i], prioridadDos ,&wg, start)
		case 3: 
			go startProcessRace(&(*procesos)[i], prioridadTres ,&wg, start)
		case 4: 
			go startProcessRace(&(*procesos)[i], prioridadCuatro ,&wg, start)
		default:
			//si llega a ser un valor que no tiene sentido se lo manda a la cola con menor prioridad
			go startProcessRace(&(*procesos)[i], prioridadCinco ,&wg, start)
		}				
	}

	close(start)

	wg.Wait()
	return nil
}
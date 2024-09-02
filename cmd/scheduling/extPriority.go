package scheduling

import (
	"fmt"
	"github/NeichS/simu/internal/structs"
	"sync"
	"time"
)


func startProcessRace(process *structs.Process, colaEspera chan struct{},wg *sync.WaitGroup, start chan struct{}) {

	defer wg.Done()
	<-start

	//lo atraso lo que tarde en arrivar 
	time.Sleep(time.Duration(process.ArrivalTime) * time.Millisecond)
	
	for (process.BurstNeeded > 0){
		select {
		case <-colaEspera:
			//obtengo rafaga de cpu 
			fmt.Printf("El proceso %s esta usando la CPU", process.Name)
			time.Sleep(time.Duration(process.BurstDuration))
			process.BurstNeeded--;
		default:
			
		}

	}
}

func cpu(colasPrioridad [5]chan struct{} ,wg *sync.WaitGroup, start chan struct{}) {
	
}

var colasPrioridad [5]chan struct{}

func StartExternalPriority(procesos *[]structs.Process) error {

	for i:= 0 ; i < len(colasPrioridad); i++ {
		colasPrioridad[i] = make(chan struct{})
	} 

	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0 ; i < len(*procesos) ; i++ {
		wg.Add(1)
		go startProcessRace(&(*procesos)[i], colasPrioridad[(*procesos)[i].ExternalPriority], &wg, start)		
	}

	close(start)

	wg.Wait()
	return nil
}
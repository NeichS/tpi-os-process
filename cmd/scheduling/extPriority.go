package scheduling

import (
	"fmt"
	"github/NeichS/simu/internal/structs"
	"sync"
	"time"
)

// Cola de procesos para cada prioridad

var finalizarUsoCpu structs.Semaphore

func startProcessRace(process *structs.Process, cola *structs.Semaphore, wg *sync.WaitGroup, start chan struct{}) {

	defer wg.Done()
	<-start

	//lo atraso lo que tarde en arrivar
	time.Sleep(time.Duration(process.ArrivalTime) * time.Millisecond)

	//espero por CPU
	cola.Wait()
	//Como es no preemptive una vez que acceda al cpu va a usar todas las rafagas que necesite
	for i := 0; i < process.BurstNeeded; i++ {
		fmt.Printf("El proceso %s esta usando la CPU", process.Name)
		time.Sleep(time.Duration(process.BurstDuration))   //rafaga de cpu
		time.Sleep(time.Duration(process.IOBurstDuration)) //rafaga entrada salida
		process.BurstNeeded--
	}

	finalizarUsoCpu.Signal()
	
}

func cpuDispatcher(colasDeEspera [5]*structs.Semaphore, procesosTotales int, wg *sync.WaitGroup, start chan struct{}) {

	procesosTerminados := 0
	defer wg.Done()
	<-start

	for procesosTerminados < procesosTotales {
		for i := 0 ; i < len(colasDeEspera); {
			if (!colasDeEspera[i].IsEmpty()) {
				colasDeEspera[i].Signal()
				
				finalizarUsoCpu.Wait() 		
			}
		}

	}	
}


func StartExternalPriority(procesos *[]structs.Process, procesosTotales int) error {

	finalizarUsoCpu = *structs.NewSemaphore(0)

	var wg sync.WaitGroup
	start := make(chan struct{})

	var colas [5]*structs.Semaphore
	for i := 0 ; i < 5; i++ {
		colas[i] = structs.NewSemaphore(0)
	}

	wg.Add(1)
	go cpuDispatcher(colas, procesosTotales, &wg, start)

	for i := 0; i < procesosTotales; i++{
		wg.Add(1)
		proceso := &(*procesos)[i]
		go startProcessRace(proceso, colas[proceso.ExternalPriority], &wg, start)
	}

	close(start)
	wg.Wait()
	return nil
}

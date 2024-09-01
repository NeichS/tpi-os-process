package scheduling

import "github/NeichS/simu/structs"


func startProcessRace(process *structs.Process) {

	select {

	}
}

func StartExternalPriority(procesos *[]structs.Process) error {

	//Todavia no se si va a ser un canal sincronico o asincronico
	// prioridadCero := make(chan structs.Process, 1000) 
	// prioridadUno := make(chan structs.Process, 1000)
	// prioridadTres := make(chan structs.Process, 1000)
	// prioridadCuatro := make(chan structs.Process, 1000)
	// prioridadCinco := make(chan structs.Process, 1000)

	for i := 0 ; i < len(*procesos) ; i++ {
		go startProcessRace(&(*procesos)[i])			
	}

	return nil
}
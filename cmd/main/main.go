package main

import (
	"encoding/csv"
	"fmt"

	//"github/NeichS/simu/scheduling"
	"github/NeichS/simu/cmd/scheduling"
	"github/NeichS/simu/internal/structs"
	"log"
	"strconv"

	"os"

	"github.com/nexidian/gocliselect"
)

func strToInt(text string) int {
	num, err := strconv.ParseInt(text, 10, 0)
	if err != nil {
		fmt.Println("Error al convertir el string a entero:", err)
		log.Fatal(err)
	}

	return int(num)
}

var procesosTotales int

func extraerProcesos(file *os.File) (*[]structs.Process, error) {
	procesos := make([]structs.Process, 0)

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(records) > 0 {
		records = records[1:] // Ignorar el encabezado
	}

	procesosTotales = 0
	for _, record := range records {
		procesosTotales++
		tiempoArribo := strToInt(record[1])
		rafagasNecesarias := strToInt(record[2])
		duracionRafaga := strToInt(record[3])
		duracionRafagaIO := strToInt(record[4])
		prioridadExterna := strToInt(record[5])

		proceso := structs.NewProcess(record[0], tiempoArribo, rafagasNecesarias, duracionRafaga, duracionRafagaIO, prioridadExterna)
		procesos = append(procesos, *proceso)
	}

	return &procesos, nil
}

func main() {

	var param string
	if len(os.Args) > 1 {
		// Acceder al primer argumento después del nombre del programa
		param = os.Args[1]
		fmt.Printf("El parámetro es: %s\n", param)
	} else {
		fmt.Println("No se proporcionó ningún parámetro.")
	}

	fileDir := "csv-files/" + param

	file, err := os.Open(fileDir)

	if err != nil {
		log.Fatal(err)
	}

	procesos, err := extraerProcesos(file)

	if err != nil {
		log.Fatal(err)
	}

	menu := gocliselect.NewMenu("Elige una politica de scheduling")

	menu.AddItem("Round robin", "rr")
	menu.AddItem("First come first server", "fcfs")
	menu.AddItem("External priority", "exPriority")
	menu.AddItem("Short process next", "spn")
	menu.AddItem("Shortest remaining time next", "srtn")

	choice := menu.Display()
	
	var tipInput, tfpInput, tcpInput string
	var tip, tfp, tcp int

	// TIP
	for {
		fmt.Print("Ingrese el tiempo que utiliza el sistema operativo para aceptar nuevos procesos (TIP): ")
		_, err = fmt.Scanln(&tipInput) // Corregido: escanea la referencia
		if err == nil {
			tip, err = strconv.Atoi(tipInput)
			if err == nil {
				break
			}
			fmt.Println("Error: ingrese un número válido para TIP.")
		}
	}

	// TFP
	for {
		fmt.Print("Ingrese el tiempo que utiliza el sistema operativo para terminar los procesos (TFP): ")
		_, err = fmt.Scanln(&tfpInput) // Corregido: escanea la referencia
		if err == nil {
			tfp, err = strconv.Atoi(tfpInput)
			if err == nil {
				break
			}
			fmt.Println("Error: ingrese un número válido para TFP.")
		}
	}

	// TCP
	for {
		fmt.Print("Ingrese el tiempo de conmutación de proceso (TCP): ")
		_, err = fmt.Scanln(&tcpInput) // Corregido: escanea la referencia
		if err == nil {
			tcp, err = strconv.Atoi(tcpInput)
			if err == nil {
				break
			}
			fmt.Println("Error: ingrese un número válido para TCP.")
		}
	}


	var quantum int
	if choice == "rr" {
		for {
			var quantumInput string
			fmt.Print("Ingrese quantum: ")
			_, err = fmt.Scanln(&quantumInput) // Corregido: escanea la referencia
			if err == nil {
				quantum, err = strconv.Atoi(quantumInput)
				if err == nil {
					break
				}
				fmt.Println("Error: ingrese un número válido para el quantum.")
			}
		}
		scheduling.StartRoundRobin(quantum)
	}

	switch choice {
	case "rr":
		scheduling.StartRoundRobin(quantum)
	case "fcfs":
	case "exPriority":
		scheduling.StartExternalPriority(procesos, procesosTotales, tip, tfp, tcp)
		//scheduling.StartExternalPriority(procesos)
	case "spn":
	case "srtn":
	}
	fmt.Printf("Choice: %s\n", choice)

}

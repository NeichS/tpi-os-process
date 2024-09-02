package main

import (
	"encoding/csv"
	"fmt"
	//"github/NeichS/simu/scheduling"
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

	for _, record := range records {
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

	fmt.Println(procesos)

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

	switch choice {
	case "rr":
		break
	case "fcfs":
		break
	case "exPriority":

		//scheduling.StartExternalPriority(procesos)
		break
	case "spn":
		break
	case "srtn":
		break
	}
	fmt.Printf("Choice: %s\n", choice)

}

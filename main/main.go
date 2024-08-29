package main

import (
	"fmt"
	"github/NeichS/simu/structs"

	"github.com/nexidian/gocliselect"
	"os"
)

func main() {
	fmt.Print("Creamos un proceso\n")

	if len(os.Args) > 1 {
		// Acceder al primer argumento después del nombre del programa
		param := os.Args[1]
		fmt.Printf("El parámetro es: %s\n", param)
	} else {
		fmt.Println("No se proporcionó ningún parámetro.")
	}

	proc := structs.NewProcess("vscode", 1, 3, 4, 5, 6)
	fmt.Printf("%s\n", proc.Name)

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
		break
	case "spn":
		break
	case "srtn":
		break
	}
	fmt.Printf("Choice: %s\n", choice)

}

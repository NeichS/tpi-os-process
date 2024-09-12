package scheduling

import (
	"fmt"
	s "github/NeichS/simu/internal/structs"
)



func ImprimirResultados(listaProcesosTerminados []*s.Process, unidadesDeTiempo, tiempoPrimerProceso, procesosTotales, tiempoSO, desperdicio int) {
	sumaTiempos := 0
	fmt.Println()
	for _, element := range listaProcesosTerminados {
		fmt.Println()
		fmt.Printf("PID: %s\n", element.PID)
		fmt.Printf("Tiempo en estado listo: %d\n", element.PCB.TiempoEnListo)
		fmt.Printf("Tiempo retorno: %d\n", element.TiempoRetorno - element.ArrivalTime)
		fmt.Printf("Tiempo de retorno normalizado: %d\n", element.TiempoRetorno / (element.BurstDuration * element.BurstNeeded))
		sumaTiempos += element.TiempoRetorno
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("Tiempo retorno de la tanda: %d\n", unidadesDeTiempo - tiempoPrimerProceso)
	fmt.Printf("Tiempo medio de retorno de la tanda: %d\n",sumaTiempos / procesosTotales)
	fmt.Printf("Tiempo desperdiciado: %d\n", desperdicio)
	procesosUso := (unidadesDeTiempo - tiempoSO - desperdicio) * 100 / unidadesDeTiempo
	fmt.Printf("Tiempo utilizado por los procesos: %d (%d%%)\n", unidadesDeTiempo - tiempoSO - desperdicio, procesosUso)
	soUso := tiempoSO * 100 / unidadesDeTiempo
	fmt.Printf("Tiempo de CPU utilizado por el SO: %d (%d%%)\n", tiempoSO, soUso)
}
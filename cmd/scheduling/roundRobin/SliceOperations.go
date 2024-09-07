package roundrobin

import "github/NeichS/simu/internal/structs"

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

package spn

import s "github/NeichS/simu/internal/structs"

func remove(procesos []*s.Process, element s.Process) []*s.Process {
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

func contains(slice []s.Process, item s.Process) bool {
	for _, element := range slice {
		if element.PID == item.PID {
			return true
		}
	}
	return false
}



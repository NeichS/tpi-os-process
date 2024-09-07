package structs

type Process struct {
	PID              string
	ArrivalTime      int //asumo que se refiere a unidad de tiempo
	BurstNeeded      int //rafagas de cpu
	BurstDuration    int
	IOBurstDuration  int
	ExternalPriority int
	State            string
	PCB              *PCB
}

func NewProcess(Name string, ArrivalTime, BurstNeeded, BurstDuration, IOBurstDuration, ExternalPriority int) *Process {
	return &Process{
		PID:              Name,
		ArrivalTime:      ArrivalTime,
		BurstNeeded:      BurstNeeded,
		BurstDuration:    BurstDuration,
		IOBurstDuration:  IOBurstDuration,
		ExternalPriority: ExternalPriority,
		State:            "New",
		PCB:              NewPCB(),
	}
}

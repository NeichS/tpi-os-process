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
	TiempoRetorno    int
	//tiempo retorno medio va a ser calculado tiempoRetorno / BurstDuration * burstNeeded
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
		TiempoRetorno:    -1,
	}
}

func (s Process) GetRemaining() int {
	return s.BurstDuration - s.PCB.TiempoRafagaEmitido
}

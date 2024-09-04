package structs

type PCB struct {
	RafagasCompletadas      int
	TiempoRafagaEmitido   int
	TiempoRafagaIOEmitido int
}

func NewPCB() *PCB {
	return &PCB{
		RafagasCompletadas:      0,
		TiempoRafagaEmitido:   0,
		TiempoRafagaIOEmitido: 0,
	}
}

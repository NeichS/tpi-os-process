package structs

type PCB struct {
	RafagasCompletadas       int
	TiempoRafagaEmitido      int
	TiempoRafagaIOEmitido    int
	TiempoEnListo            int
	TiempoRetorno            int
	TiempoRetornoNormalizado int
}

func NewPCB() *PCB {
	return &PCB{
		RafagasCompletadas:    0,
		TiempoRafagaEmitido:   0,
		TiempoRafagaIOEmitido: 0,
		TiempoEnListo:         0,
	}
}

package structs

type PCB struct {
	RafagasCompletadas       int
	TiempoRafagaEmitido      int
	TiempoRafagaIOEmitido    int
	TiempoEnListo            int
	TiempoRetorno            int
	TiempoRetornoNormalizado int
	TiempoTCP                int
	TiempoTIP                int
	TiempoTFP                int
	OperacionSOActual        string //TFP, TCP TIP o vacio
}

func NewPCB() *PCB {
	return &PCB{
		RafagasCompletadas:    0,
		TiempoRafagaEmitido:   0,
		TiempoRafagaIOEmitido: 0,
		TiempoEnListo:         0,
		TiempoTIP:             0,
		TiempoTCP:             0,
		TiempoTFP:             0,
		OperacionSOActual:     "",
	}
}

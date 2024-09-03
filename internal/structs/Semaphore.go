package structs

type Semaphore struct {
	channel chan struct{}
}

// NewSemaphore crea un nuevo semáforo con un número dado de permisos.
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		channel: make(chan struct{}, n),
	}
}

// Acquire intenta adquirir un permiso del semáforo.
func (s *Semaphore) Wait() {
	s.channel <- struct{}{}
}

// Release libera un permiso al semáforo.
func (s *Semaphore) Signal() {
	<-s.channel
}

func (s *Semaphore) IsEmpty() bool {
	return len(s.channel) == 0
}
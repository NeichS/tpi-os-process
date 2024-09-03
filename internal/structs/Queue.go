package structs

import (
	"sync"

)

// Queue es una estructura que representa una cola para manejar procesos
type Queue struct {
	items    []*Process // Slice para almacenar los procesos en la cola
	lock     sync.Mutex         // Mutex para asegurar la exclusión mutua
	notEmpty *sync.Cond         // Condición para manejar la disponibilidad de elementos
}

// NewQueue crea una nueva cola vacía para procesos
func NewQueue() *Queue {
	q := &Queue{
		items: make([]*Process, 0),
	}
	q.notEmpty = sync.NewCond(&q.lock)
	return q
}

// Enqueue agrega un proceso al final de la cola
func (q *Queue) Enqueue(process *Process) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.items = append(q.items, process) // Agrega el proceso a la cola
	q.notEmpty.Signal()                // Notifica que la cola ya no está vacía
}


func (q *Queue) Dequeue() *Process {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Espera hasta que la cola no esté vacía
	for len(q.items) == 0 {
		q.notEmpty.Wait()
	}

	process := q.items[0]    // Obtiene el primer proceso
	q.items = q.items[1:]    // Remueve el primer proceso de la cola

	return process
}

// IsEmpty devuelve verdadero si la cola está vacía
func (q *Queue) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.items) == 0
}

// Size devuelve el número de procesos en la cola
func (q *Queue) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.items)
}

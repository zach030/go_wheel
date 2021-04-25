package runtime

import (
	"fmt"
	"sync"
)

type waitQueue struct {
	lock  sync.Mutex
	queue []*GoRoutine
}

func NewWaitQueue() *waitQueue {
	return &waitQueue{
		queue: make([]*GoRoutine, 0),
	}
}

func (q *waitQueue) add(g *GoRoutine) {
	q.lock.Lock()
	defer q.lock.Unlock()
	g.Block()
	fmt.Printf("[Block Queue] Goroutine ID :%v is add to block queue now\n", g.ID)
	q.queue = append(q.queue, g)
}

func (q *waitQueue) get() (*GoRoutine, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.isEmpty() {
		return nil, EmptyQueue
	}
	g := q.queue[0]
	q.queue = q.queue[1:len(q.queue)]
	fmt.Printf("[Block Queue] Get Goroutine ID :%v from block queue now\n", g.ID)
	return g, nil
}

func (q *waitQueue) isEmpty() bool {
	return q.size() == 0
}

func (q *waitQueue) size() int {
	return len(q.queue)
}

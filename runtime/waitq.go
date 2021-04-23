package runtime

import "sync"

type waitQueue struct {
	lock  sync.Mutex
	queue []*GoRoutine
}

func NewWaitQueue() *waitQueue {
	return &waitQueue{
		queue: make([]*GoRoutine, 0),
	}
}

func (q *waitQueue) isEmpty() bool {
	return q.size() == 0
}

func (q *waitQueue) size() int {
	return len(q.queue)
}

func (q *waitQueue) add(g *GoRoutine) {
	q.lock.Lock()
	defer q.lock.Unlock()
	g.Block()
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
	return g, nil
}

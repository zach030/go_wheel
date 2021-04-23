package runtime

import (
	"fmt"
	"time"
)

const (
	Run   = "run"
	Block = "Block"
)

type GoRoutine struct {
	ID     string
	status string
	data   map[*GoChan]interface{}
}

func NewGoroutine(id string) *GoRoutine {
	return &GoRoutine{
		ID:     id,
		status: Run,
	}
}

func (g *GoRoutine) SendChannel(data interface{}, ch *GoChan) {
	ch.Send(g, data)
}

func (g *GoRoutine) RecvChannel(ch *GoChan) interface{} {
	ch.Recv(g)
	for {
		if g.data[ch] == nil {
			g.Block()
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	return g.data
}

func (g *GoRoutine) Recv(ch *GoChan, data interface{}) {
	fmt.Printf("[Goroutine ID:%v] is Wake by channel:%+v\n", g.ID, ch)
	g.Wake()
	g.data[ch] = data
}

func (g *GoRoutine) Block() {
	g.status = Block
	fmt.Printf("[Goroutine ID:%v] is Block now\n", g.ID)
}

func (g *GoRoutine) Wake() {
	g.status = Run
	fmt.Printf("[Goroutine ID:%v] is Wake now\n", g.ID)
}

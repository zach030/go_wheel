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
	ID     string                  // goroutine的唯一ID
	status string                  // 当前状态，运行和阻塞
	data   map[*GoChan]interface{} // 数据map
}

func NewGoroutine(id string) *GoRoutine {
	fmt.Printf("[Goroutine] ID=%v is created now!\n", id)
	return &GoRoutine{
		ID:     id,
		status: Run,
		data:   make(map[*GoChan]interface{}, 0),
	}
}

// 根据channel，选择数据
func (g *GoRoutine) DirectSend(ch *GoChan) interface{} {
	if v, ok := g.data[ch]; ok {
		delete(g.data, ch)
		fmt.Printf("[Goroutine] Get data from GoRoutine ID:%v , data :%v\n", g.ID, v)
		return v
	}
	return nil
}

// 将数据发送到ch
func (g *GoRoutine) SendChannel(data interface{}, ch *GoChan) {
	g.data[ch] = data
	ch.Send(g, data)
}

// 从ch接收数据
func (g *GoRoutine) RecvChannel(ch *GoChan) interface{} {
	ch.Recv(g)
	for {
		if g.data[ch] == nil {
			g.Block()
			time.Sleep(time.Second)
		} else {
			fmt.Printf("[Goroutine ID:%v] recv data:%v now!\n", g.ID, g.data[ch])
			break
		}
	}
	return g.data[ch]
}

// 从阻塞状态，被唤醒，接收数据
func (g *GoRoutine) BlockRecv(ch *GoChan, data interface{}) {
	fmt.Printf("[Goroutine ID:%v] is Wake by channel, recv data:%v\n", g.ID, data)
	g.Wake()
	g.data[ch] = data
}

// 运行状态接收数据
func (g *GoRoutine) WakeRecv(ch *GoChan, data interface{}) {
	fmt.Printf("[GoRoutine] ID:%v Wake Recv data from channel now!\n", g.ID)
	g.data[ch] = data
}

// 阻塞
func (g *GoRoutine) Block() {
	g.status = Block
	fmt.Printf("[Goroutine ID:%v] is Block now, time:%v\n", g.ID, time.Now().String())
}

// 唤醒
func (g *GoRoutine) Wake() {
	g.status = Run
	fmt.Printf("[Goroutine ID:%v] is Wake now, time:%v\n", g.ID, time.Now().String())
}

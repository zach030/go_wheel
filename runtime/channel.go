package runtime

import (
	"errors"
	"fmt"
	"sync"
)

const (
	String = "string"
	Int    = "int"
)

var (
	EmptyQueue = errors.New("empty goroutine queue")
	FullChan   = errors.New("channel is full now")
	EmptyChan  = errors.New("channel is empty now")
)

type GoChan struct {
	queueSize  uint // 可存放的元素个数,最大
	remainSize uint // 队列中剩余的元素个数(空的元素个数)
	// 单向队列实现
	// todo implement by ringBuffer
	buf        []int // 数据队列
	bufPointer int   // 队列指针

	elementType string // 存放的元素类型
	elemSize    uint   // 元素类型大小

	sendPos uint // 写入通道时的位置
	recvPos uint // 下一个被读取的元素在数组中的位置

	closed bool // 关闭状态

	lock  sync.Mutex // 互斥锁
	sendQ *waitQueue // 发送阻塞队列
	recvQ *waitQueue // 接收阻塞队列
}

func NewGoChan(typ string, size uint) *GoChan {
	return &GoChan{
		queueSize:   size,
		remainSize:  size,
		buf:         make([]int, size),
		bufPointer:  0,
		elementType: typ,
		elemSize:    getSize(typ),
		sendPos:     0,
		recvPos:     0,
		closed:      false,
		lock:        sync.Mutex{},
		sendQ:       NewWaitQueue(),
		recvQ:       NewWaitQueue(),
	}
}

func getSize(typ string) uint {
	switch typ {
	case "int":
		return 8
	default:
		return 0
	}
}

// goRoutine send data to GoChan
func (g *GoChan) Send(goRoutine *GoRoutine, data interface{}) {
	v := data.(int)
	// 如果接收阻塞队列非空,唤醒头一个阻塞g,赋值
	if !g.recvQ.isEmpty() {
		headG, err := g.recvQ.get()
		if err != nil {
			fmt.Println("get recv block queue failed,err:", err)
			return
		}
		// 被阻塞的goroutine接收此data
		headG.Recv(g, data)
		return
	} else {
		// 如果接收阻塞队列为空
		if g.isBufFull() {
			// 如果缓冲区已满,则加入sendQ
			g.sendQ.add(goRoutine)
			return
		}
		// 缓冲区有空余位置
		g.buf = append(g.buf, v)
		fmt.Printf("GoRoutine ID:%v success send data:%v to channel:%+v\n", goRoutine.ID, v, g)
	}
}

// 从GoChan里读数据
func (g *GoChan) Recv(routine *GoRoutine) interface{} {
	if !g.sendQ.isEmpty() {
		// 发送阻塞队列不为空
		if len(g.buf) == 0 {
			// 无缓冲区
			headG, err := g.sendQ.get()
			if err != nil {
				fmt.Println("get send block queue failed,err:", err)
				return nil
			}

		}
	}
	// 发送阻塞队列空
	// 缓冲区有数据
	if g.remainSize > 0 {
		return g.buf[g.bufPointer]
	}
	return nil
}

// 判断buf是否有空闲位
func (g *GoChan) isBufFull() bool {
	return g.remainSize == 0
}

// 从buf读数据
func (g *GoChan) getDataFromBuf() interface{} {
	if g.remainSize == g.queueSize {
		// 缓冲区为空
		fmt.Println(EmptyChan)
		return nil
	}
	// 读第一个元素
	data := g.buf[0]
	// 读指针前移
	g.recvPos--

	return data
}

func (g *GoChan) addDataToBuf(data interface{}) {
	v := data.(int)
	if g.remainSize == 0 {
		// 缓冲队列已满
		fmt.Println(FullChan)
		return
	}
	// 加入缓冲队列
	g.buf[g.bufPointer] = v
	// 指针后移
	g.bufPointer++
	// 空闲数减少
	g.remainSize--
}

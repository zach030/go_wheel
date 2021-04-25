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
	fmt.Printf("[Channel] new channel with type:%v size:%v is created now!\n", typ, size)
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
		headG.BlockRecv(g, data)
		return
	} else {
		// 如果接收阻塞队列为空
		if g.isBufFull() {
			// 如果缓冲区已满,则加入sendQ
			g.sendQ.add(goRoutine)
			return
		}
		// 缓冲区有空余位置
		g.addDataToBuf(data)
		fmt.Printf("GoRoutine ID:%v success send data:%v to channel\n", goRoutine.ID, v)
	}
}

// goroutine 从 GoChan里读数据
func (g *GoChan) Recv(routine *GoRoutine) {
	// 1 判断发送消息阻塞队列是否为空
	if !g.sendQ.isEmpty() {
		// 1.1 发送阻塞队列不为空
		headG, err := g.sendQ.get()
		if err != nil {
			fmt.Println("get send block queue failed,err:", err)
			return
		}
		if len(g.buf) == 0 {
			// 1.2 无缓冲区，从发送消息阻塞队列中取一个go
			// 1.3 取出这个go要发的数据，直接给routine接收
			routine.WakeRecv(g, headG.DirectSend(g))
			return
		}
		// 1.3 有缓冲区，则缓冲区已满，先从buf中取出一个队头数据给当前go
		routine.WakeRecv(g, g.getDataFromBuf())
		// 1.4 再从发送阻塞队列中取出队头，把队头的数据放入buf中
		g.addDataToBuf(headG.DirectSend(g))
		return
	}
	// 2 发送阻塞队列为空
	if g.isBufHasData() {
		// 2.1 缓冲区有数据，直接读取buf中的第一个数据
		routine.WakeRecv(g, g.getDataFromBuf())
		return
	}
	// 2.2 缓冲区无数据，此goroutine被加入读取阻塞队列
	g.recvQ.add(routine)
	return
}

func (g *GoChan) isBufHasData() bool {
	return g.remainSize != g.queueSize
}

// 判断buf是否有空闲位
func (g *GoChan) isBufFull() bool {
	return g.remainSize == 0
}

// todo 从buf读数据
func (g *GoChan) getDataFromBuf() interface{} {
	if g.remainSize == g.queueSize {
		// 缓冲区为空
		fmt.Println(EmptyChan)
		return nil
	}
	// 读第一个元素
	data := g.buf[0]
	// 数组左移一位
	newBuf := make([]int, g.queueSize)
	copy(newBuf, g.buf[1:])
	g.buf = newBuf
	// 指针前移
	g.bufPointer--
	// 空闲数增加
	g.remainSize++
	fmt.Printf("[Channel] get data:%v from channel\n", data)
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
	fmt.Printf("[Channel] add data:%v to channel\n", data)
	return
}

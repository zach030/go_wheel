package ds

import (
	"sync/atomic"
	"unsafe"
)

type WaitGroup struct {
	state1 [3]uint32 // 还未执行结束的：counter 等待者数量：waiter  信号量：semaphore
}

func (wg *WaitGroup) Add(delta int) {
	// 取地址：counter和 waiter
	statep, semap := wg.state()
	// 取出state数组，也就是state1的前两位，给counter加上delta
	state := atomic.AddUint64(statep, uint64(delta)<<32)
	c := int32(state >> 32) // counter
	w := uint32(state)      // waiter
	// 如果counter<0，等待者小于0，触发panic
	if c < 0 {
		panic("negative counter")
	}
	// 如果counter>0 或者 waiter==0，说明无等待者，不需要释放信号量
	if c > 0 || w == 0 {
		return
	}
	// counter==0 waiter >0
	*statep = 0
	for ; w != 0; w-- {
		// 释放信号量,唤醒等待者
		semeRelease(semap, false)
	}
}

func (wg *WaitGroup) Wait() {
	// waiter++
}

func (wg *WaitGroup) Done() {
	// counter--
	wg.Add(-1)
}

// state returns pointers to the state and sema fields stored within wg.state1.
func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
	if uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
		return (*uint64)(unsafe.Pointer(&wg.state1)), &wg.state1[2]
	} else {
		return (*uint64)(unsafe.Pointer(&wg.state1[1])), &wg.state1[0]
	}
}

func semeRelease(p *uint32, state bool) {

}

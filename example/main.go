package main

import (
	"fmt"
	"go_implement/runtime"
)

func main() {
	ch := runtime.NewGoChan(runtime.Int, 1)

	g1 := runtime.NewGoroutine("1")
	g2 := runtime.NewGoroutine("2")
	//g3 := goroutine.NewGoroutine("3")

	g1.SendChannel(1, ch)
	data := g2.RecvChannel(ch)
	//g3.SendChannel(3, ch)

	fmt.Println(data)
}

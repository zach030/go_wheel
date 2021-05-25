package main

import "fmt"

func main() {
	//ch := runtime.NewGoChan(runtime.Int, 1)
	//
	//g1 := runtime.NewGoroutine("1")
	//g2 := runtime.NewGoroutine("2")
	//g3 := runtime.NewGoroutine("3")
	//
	//g1.SendChannel(1, ch)
	//g3.SendChannel(3, ch)
	//data := g2.RecvChannel(ch)
	//fmt.Println(data)
	var slice []int
	fmt.Println(len(slice),cap(slice))
	s1 := append(slice,1,2,3)
	fmt.Println(len(s1),cap(s1))
	s2 := append(s1,4)
	fmt.Println(len(s2),cap(s2))
	fmt.Println(&s1[0]==&s2[0])
}

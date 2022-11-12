package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Mutex
	m.Lock()
	m.Unlock()
	m.Unlock()
}

func putVal(c chan int) {
	for i := 0; i < 100; i++ {
		c <- i
	}
}

func getVal(name string, c chan int) {
	for i := 0; i < 20; i++ {
		val := <-c
		fmt.Printf("from %s, vali %d \n", name, val)
	}
}

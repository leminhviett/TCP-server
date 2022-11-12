package main

import (
	"fmt"
)

func main() {
	// c1 := make(chan int, 100)
	// go putVal(c1)

	// go getVal("reader1", c1)
	// go getVal("reader2", c1)

	// time.Sleep(5*time.Second)

	a := make([]int, 0, 5)
	fmt.Println(len(a), cap(a))
	a = append(a, 1,2,3,4,5,6)
	fmt.Println(cap(a))

	fmt.Println(a)
}

func putVal(c chan int) {
	for i := 0; i < 100; i++{
		c <- i
	}
}

func getVal(name string, c chan int) {
	for i := 0; i < 20; i++{
		val := <- c
		fmt.Printf("from %s, vali %d \n", name, val)
	}
}

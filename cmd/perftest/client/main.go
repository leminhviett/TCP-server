package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main(){
	performanceTester()
}

func performanceTester() {
	var wg sync.WaitGroup
	simpleCaller := func() {
		resp, err := http.Get("http://localhost:8001/")
		if err != nil {
			fmt.Println(err.Error())
		}
	
		fmt.Println(resp)
		wg.Done()
	}

	timeStart := time.Now()

	times := 100
	for j:= 0; j < times; j ++ {
		wg.Add(1)
		go simpleCaller()
	}

	wg.Wait()
	timeEnd := time.Now()

	fmt.Println("Elapsed: ...")
	fmt.Println(timeEnd.Sub(timeStart).Seconds())
}
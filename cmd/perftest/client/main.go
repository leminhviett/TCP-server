package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main(){
	performanceTester(1000, 8)
}

func performanceTester(times, threads int) {
	var wg sync.WaitGroup
	var m sync.Mutex

	failedTimes := 0
	
	simpleCaller := func() {
		for j:= 0; j < times; j ++ {
			_, err := http.Get("http://localhost:8001/")
			if err != nil {
				m.Lock()
				failedTimes += 1 
				defer m.Unlock()
				fmt.Println(err.Error())
			}
		}
		wg.Done()
	}

	timeStart := time.Now()
	for i := 0; i < threads; i ++ {
		fmt.Println("thread created")
		wg.Add(1)
		go simpleCaller()
	}


	wg.Wait()
	timeEnd := time.Now()

	fmt.Printf("Elapsed: %f \n", timeEnd.Sub(timeStart).Seconds())

	var failedRate float64 = float64(failedTimes)/float64(times)
	fmt.Printf("Failed rate: %f \n", failedRate)
}
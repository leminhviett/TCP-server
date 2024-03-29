package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/leminhviett/TCP-server/config"
)

func main() {
	performanceTester(50, 20)
}

func performanceTester(times, threads int) {
	var wg sync.WaitGroup
	var m sync.Mutex
	var m2 sync.Mutex

	failedTimes := 0
	count := 0

	simpleCaller := func() {
		for j := 0; j < times; j++ {
			resp, err := http.Get(fmt.Sprintf("http://%s:%s", config.BFF_SERVER_CONN_HOST, config.BFF_SERVER_CONN_PORT))
			if err != nil || (resp != nil && resp.StatusCode == 500) {
				m.Lock()
				failedTimes += 1
				m.Unlock()
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			m2.Lock()
			count += 1
			m2.Unlock()
		}
		wg.Done()
	}

	timeStart := time.Now()
	for i := 0; i < threads; i++ {
		fmt.Println("thread created")
		wg.Add(1)
		go simpleCaller()
	}

	wg.Wait()
	timeEnd := time.Now()

	fmt.Printf("Elapsed: %f \n", timeEnd.Sub(timeStart).Seconds())

	var failedRate float64 = float64(failedTimes) / float64(times)
	fmt.Printf("Failed rate: %f \n", failedRate)
	fmt.Printf("failed times: %d; total request: %d", failedTimes, count)
}

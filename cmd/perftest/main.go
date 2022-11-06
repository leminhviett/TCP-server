package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/utils"
)

func main(){
	go startBFFWithoutConnPool()
	clientPerfTester()
}


func startBFFWithoutConnPool() {
	r := mux.NewRouter()

	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := net.Dial(config.CONN_TYPE, 
			fmt.Sprintf("%s:%s", config.CONN_HOST, config.CONN_PORT))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer conn.Close()

		messsage := &utils.Message{
			ApplicationRoute: "hellodummy",
			ApplicationData: []byte{1,2,3,4,5,6,7,8,9,10},
		}
		
		utils.WriteTo(conn, messsage)
		fmt.Fprintf(w, "Message sent")
	}

    r.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8001", r)
}

func clientPerfTester() {
	times := 100
	batchN := 5

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
	for i:= 0; i < batchN; i ++ {
		for j:= 0; j < times; j ++ {
			wg.Add(1)
			go simpleCaller()
		}
		// time.Sleep(100*time.Millisecond)
	}

	wg.Wait()
	timeEnd := time.Now()

	fmt.Println("Elapsed: ...")
	fmt.Println(timeEnd.Sub(timeStart).Seconds())
}

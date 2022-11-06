package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/utils"
)

func main() {
	startBFFWithoutConnPool()
}

func startBFFWithoutConnPool() {
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

	r := mux.NewRouter()
    r.HandleFunc("/", handler)
	http.ListenAndServe("localhost:8001", r)
}

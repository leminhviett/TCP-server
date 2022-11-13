package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/common"
	connpool "github.com/leminhviett/TCP-server/domain/connPool"
)

func BeForFe() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := net.Dial(config.TCP_CONN_TYPE,
			fmt.Sprintf("%s:%s", config.TCP_SERVER_CONN_HOST, config.TCP_SERVER_CONN_PORT))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer conn.Close()
		fmt.Printf("%s created \n", conn.LocalAddr().String())

		common.WriteToConn(conn, dummyMessage)
		fmt.Fprintf(w, "Message sent")
	}

	startHTTPServer(handler)
}

func BeForFeWConnPool() {
	ctx := context.Background()
	pool := connpool.NewConnPool(ctx, 12, 12)

	handler := func(w http.ResponseWriter, r *http.Request) {
		ctxRequest, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		conn, err := pool.GetConn(ctxRequest)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(500)
			return
		}
		defer pool.PutConn(ctx, conn)

		_, err = common.WriteToConn(conn, dummyMessage)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		message, _ := common.ReadFromConn(conn)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(message)
		return
	}

	startHTTPServer(handler)
}

func startHTTPServer(handler func(w http.ResponseWriter, r *http.Request)) {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	fmt.Println("Listening on " + config.BFF_SERVER_CONN_PORT)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.BFF_SERVER_CONN_HOST, config.BFF_SERVER_CONN_PORT), r)
	if err != nil {
		fmt.Println(err)
	}
}

var dummyMessage = &common.Message{
	ApplicationRoute: "hellodummy",
	ApplicationData:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
}

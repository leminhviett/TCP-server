package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/leminhviett/TCP-server/config"
	connpool "github.com/leminhviett/TCP-server/domain/connPool"
	"github.com/leminhviett/TCP-server/domain/utils"
)

var (
	messsage = &utils.Message{
		ApplicationRoute: "hellodummy",
		ApplicationData:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
)

func main() {
	startBFFWithConnPool()
}

func startBFFWithoutConnPool() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := net.Dial(config.TCP_CONN_TYPE,
			fmt.Sprintf("%s:%s", config.TCP_SERVER_CONN_HOST, config.TCP_SERVER_CONN_PORT))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer conn.Close()
		fmt.Printf("%s created \n", conn.LocalAddr().String())

		utils.WriteTo(conn, messsage)
		fmt.Fprintf(w, "Message sent")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	fmt.Println("Listening on " + config.BFF_SERVER_CONN_PORT)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.BFF_SERVER_CONN_HOST, config.BFF_SERVER_CONN_PORT), r)
	if err != nil {
		fmt.Println(err)
	}
}

func startBFFWithConnPool() {
	ctx := context.Background()
	pool := connpool.NewConnPool(ctx, 15, 15)

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

		_, err = utils.WriteTo(conn, messsage)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		message, _ := utils.ReadFrom(conn)
		fmt.Println(message)
		w.WriteHeader(200)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	fmt.Println("Listening on " + config.BFF_SERVER_CONN_PORT)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.BFF_SERVER_CONN_HOST, config.BFF_SERVER_CONN_PORT), r)
	if err != nil {
		fmt.Println(err)
	}
}

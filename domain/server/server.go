package server

import (
	"fmt"
	"net"
	"os"

	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/utils"
)

func StartServer() {
	l, err := net.Listen(config.CONN_TYPE, config.CONN_HOST+":"+config.CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + config.CONN_HOST + ":" + config.CONN_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	fmt.Println("client at: " + conn.RemoteAddr().String())

	for {
		message, err := utils.ReadFrom(conn)
		if err != nil {
			conn.Write([]byte("Error: " + err.Error()))
			return
		}
		fmt.Println(message)

		conn.Write([]byte("Message received."))
	}

}

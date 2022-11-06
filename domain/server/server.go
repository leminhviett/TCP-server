package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/utils"
	"github.com/leminhviett/TCP-server/domain/utils/customerror"
)

func StartBackend() {
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
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	fmt.Println("client at: " + conn.RemoteAddr().String())

	for {
		message, err := utils.ReadFrom(conn)
		switch err {
		case nil:
			conn.Write([]byte("Message received."))
			fmt.Println(message)
		case io.EOF:
			conn.Write([]byte(customerror.ErrorConnClosed.Error()))
		default:
			conn.Write([]byte("Error: " + err.Error()))
			return
		}
	}

}

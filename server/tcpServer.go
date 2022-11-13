package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/common"
	"github.com/leminhviett/TCP-server/domain/customError"
)

func StartTCPServer() {
	l, err := net.Listen(config.TCP_CONN_TYPE, config.TCP_SERVER_CONN_HOST+":"+config.TCP_SERVER_CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + config.TCP_SERVER_CONN_HOST + ":" + config.TCP_SERVER_CONN_PORT)

	connTracker := make(map[string]bool)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		connAddress := conn.RemoteAddr().String()
		if _, ok := connTracker[connAddress]; !ok {
			fmt.Println("client at: " + connAddress)
			connTracker[connAddress] = true
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		message, err := common.ReadFromConn(conn)
		switch err {
		case nil:
			common.WriteToConn(conn, &common.Message{
				ApplicationRoute: message.ApplicationRoute,
				ApplicationData:  []byte("Data received"),
			})
		case io.EOF:
			common.WriteToConn(conn, &common.Message{
				ApplicationData: []byte(customError.ErrorConnClosed.Error()),
			})
		default:
			common.WriteToConn(conn, &common.Message{
				ApplicationData: []byte("Error: " + err.Error()),
			})
			return
		}
	}

}

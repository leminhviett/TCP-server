package server

import (
	"fmt"
	"net"
	"os"

	"github.com/leminhviett/TCP-server/config"
)

func StartServer() {
	// Listen for incoming connections.
	l, err := net.Listen(config.CONN_TYPE, config.CONN_HOST+":"+config.CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + config.CONN_HOST + ":" + config.CONN_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("client at: " + conn.RemoteAddr().String())
	// Make a buffer to hold incoming data.

	for {
		buf := make([]byte, 1024)

		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		fmt.Println("Receiving (bytes): ", reqLen)
		fmt.Println(string(buf))

		// Send a response back to person contacting us.
		conn.Write([]byte("Message received."))
		// Close the connection when you're done with it.
	}

}

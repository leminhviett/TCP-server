package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/leminhviett/TCP-server/config"
)

func StartClient() {
	conn, err := net.Dial(config.CONN_TYPE, 
		fmt.Sprintf("%s:%s", config.CONN_HOST, config.CONN_PORT))
	if err != nil {
		panic(err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		returnB := make([]byte, 1024)
		conn.Read(returnB)
		fmt.Println("->: " + string(returnB))

	}

}

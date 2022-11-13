package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/common"
)

func main() {
	startTCPClient()
}

func startTCPClient() {
	conn, err := net.Dial(config.TCP_CONN_TYPE,
		fmt.Sprintf("%s:%s", config.TCP_SERVER_CONN_HOST, config.TCP_SERVER_CONN_PORT))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> Enter route ")
		text, _ := reader.ReadString('\n')

		messsage := &common.Message{
			ApplicationRoute: text,
		}
		fmt.Println(messsage)

		_, err := common.WriteToConn(conn, messsage)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		returnB := make([]byte, 1024)
		conn.Read(returnB)
		fmt.Println("->: " + string(returnB))
	}
}

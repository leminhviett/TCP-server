package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/utils"
)

func StartClientCmd() {
	conn, err := net.Dial(config.CONN_TYPE, 
		fmt.Sprintf("%s:%s", config.CONN_HOST, config.CONN_PORT))
	if err != nil {
		panic(err)
	}
	defer conn.Close()


	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> Enter route ")
		text, _ := reader.ReadString('\n')

		messsage := &utils.Message{
			ApplicationRoute: text,
		}
		fmt.Println(messsage)
		
		_, err := utils.WriteTo(conn, messsage)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		returnB := make([]byte, 1024)
		conn.Read(returnB)
		fmt.Println("->: " + string(returnB))
	}
}

package main

import "github.com/leminhviett/TCP-server/server"

func main() {
	go server.StartTCPServer()
	server.BeForFeWConnPool()
}

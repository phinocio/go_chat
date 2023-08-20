package main

import (
	"go_chat/src/client"
	"go_chat/src/server"
)

func main() {
	client.Run()
	server.Run()
}

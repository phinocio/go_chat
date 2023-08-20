package main

import (
	"fmt"
	"go_chat/src/client"
	"go_chat/src/server"
	"os"
	"strings"
)

func main() {
	var arg = strings.ToLower(os.Args[1])

	if arg == "client" {
		client.Run()
	} else if arg == "server" {
		server.Run()
	} else {
		fmt.Println("Please pass one of 'client' or 'server' as an arg")
	}
}

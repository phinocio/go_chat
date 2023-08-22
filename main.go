package main

import (
	"fmt"
	"os"
	"strings"

	"go_chat/src/client"
	"go_chat/src/server"
)

func main() {
	var instance string
	var host string
	var port string

	// Get client or server
	if len(os.Args) >= 2 {
		instance = strings.ToLower(os.Args[1])
	} else {
		instance = "client" // Default to client
	}


	// Get the host
	if len(os.Args) >= 3 {
		host = strings.ToLower(os.Args[2])
	} else {
		host = "127.0.0.1"
	}

	// Get the port
	if len(os.Args) >= 4 {
		port = strings.ToLower(os.Args[3])
	} else {
		port = "4444"
	}

	// Run client or server
	if instance == "client" {
		client.Run(host, port)
	} else if instance == "server" {
		server.Run(host, port)
	} else {
		fmt.Println("Something went wrong with args passed. Found ", len(os.Args))
	}
}

package server

import (
	"fmt"
	"net"
)

func Run(host string, port string) {
	fmt.Println(host, port)
	fmt.Println("Hello from Server!")
	net.Listen("tcp4", "4444")
}


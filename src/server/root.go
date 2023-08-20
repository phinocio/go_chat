package server

import (
	"io"
	"fmt"
	"log"
	"net"
	// "os"
)

func Run(host string, port string) {
	fmt.Println(host, port)
	ln, err := net.Listen("tcp4", host+":"+port)

	if err != nil {
		fmt.Println("An error happened, ", err)
	}

	fmt.Println("[INFO] Listening on " + host + ":" + port)

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("An error happened in the for loop, ", err)
		}

		log.Println("[INFO] Connection received from: ", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	streamMessages(conn)
	defer closeConnection(conn)
}

func streamMessages(conn net.Conn) {
    tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
    for {
        n, err := conn.Read(tmp)
        if err != nil {
            if err != io.EOF {
                fmt.Println("read error:", err)
            }
            break
        }
        // fmt.Println("got", n, "bytes.")
		fmt.Print("[INFO] Msg from ", conn.RemoteAddr(),  ": ", string(tmp[:n]))
    }
}

func closeConnection(conn net.Conn) {
	log.Println("[INFO] Connection closed from:", conn.RemoteAddr())
	conn.Close()
}

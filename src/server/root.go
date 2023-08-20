package server

import (
	"fmt"
	"net"
	"io"
	"os"
)

func Run(host string, port string) {
	fmt.Println(host, port)
	ln, err := net.Listen("tcp4", host+":"+port)

	if err != nil {
		fmt.Println("An error happened, ", err)
	}

	fmt.Println("Listening on " + host + ":" + port)

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("An error happened in the for loop, ", err)
		}

		// go func(c net.Conn) {
  //           defer c.Close()
  //           io.Copy(os.Stdout, c)
  //       }(conn)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// defer conn.Close()
	// buf, read_err := io.ReadAll(conn)
	// fmt.Println("Connection received from", conn.RemoteAddr())
	// if read_err != nil {
	// 	fmt.Println("failed:", read_err)
	// 	return
	// }
	// fmt.Println("Got: ", string(buf))
	//
	// _, write_err := conn.Write([]byte("Message received.\n"))
	// if write_err != nil {
	// 	fmt.Println("failed:", write_err)
	// 	return
	// }
	defer conn.Close()
    io.Copy(os.Stdout, conn) }


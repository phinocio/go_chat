package server

import (
	"io"
	"net"

	// "os"

	"go_chat/src/utils/log_msgs"
)

func Run(host string, port string) {
	log_msgs.InfoLog( host + " " + port)
	ln, err := net.Listen("tcp4", host+":"+port)

	if err != nil {
		log_msgs.ErrorTimeLog( err.Error() )
	}

	log_msgs.InfoLog("Listening on " + host + ":" + port)

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log_msgs.ErrorLog( "An error happened in the for loop" + err.Error())
		}

		log_msgs.InfoTimeLog( "Connection received from: " + conn.RemoteAddr().String()  )

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
				log_msgs.ErrorLog("read error:" + err.Error())
            }
            break
        }
        // fmt.Println("got", n, "bytes.")
		log_msgs.InfoLog("Msg from " + conn.RemoteAddr().String() + ":" + string(tmp[:n]))
    }
}

func closeConnection(conn net.Conn) {
	log_msgs.InfoTimeLog("Connection close from: " + conn.RemoteAddr().String())
	conn.Close()
}

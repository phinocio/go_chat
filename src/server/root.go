package server

import (
	"io"
	"net"
	"strings"

	// "os"

	"go_chat/src/utils/log_msgs"
)


type Connection struct {
	name string // The name of this client instance
	peer string	// The person we are trying to talk with
}

var clients []net.Conn

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

		log_msgs.InfoTimeLog( "Connection received from: " + conn.RemoteAddr().String())
		tmp := make([]byte, 256)
		n, err := conn.Read(tmp)
		var metaData = strings.Trim(string(tmp[:n]), "\n")
		nameTarget := strings.Split(metaData, ":")
		log_msgs.InfoLog("Metadata received: " + string(tmp[:n]))
		log_msgs.InfoLog("Name: " + nameTarget[0] + ". Tareget: " + nameTarget[1])


		clients = append(clients, conn)
		log_msgs.InfoLog("Last connection: " + clients[len(clients)-1].RemoteAddr().String())
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
		var msg = strings.Trim(string(tmp[:n]), "\n")
		// log_msgs.InfoLog("Msg from " + conn.RemoteAddr().String() + ": " + msg)

		if conn == clients[0] {
			log_msgs.InfoLog("Msg from " + conn.RemoteAddr().String() + " sent to " + clients[1].RemoteAddr().String() + ". Msg: " + msg)
			clients[1].Write([]byte(msg))
		}
		if conn == clients[1] {
			log_msgs.InfoLog("Msg from " + conn.RemoteAddr().String() + " sent to " + clients[1].RemoteAddr().String() + ". Msg: " + msg)
			clients[0].Write([]byte(msg))
		}
    }
}

func closeConnection(conn net.Conn) {
	log_msgs.InfoTimeLog("Connection close from: " + conn.RemoteAddr().String())
	conn.Close()
}

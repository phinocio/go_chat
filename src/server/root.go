package server

import (
	"io"
	"net"
	"strings"

	// "os"

	"go_chat/src/utils/log_msgs"
	"go_chat/src/utils/network"
)

// type Connection struct {
// 	name string // The name of this client instance
// 	peer string	// The person we are trying to talk with
// 	conn net.Conn // The client's connection
// }

var clients []network.Connection

func Run(host string, port string) {
	log_msgs.InfoLog(host + " " + port)
	ln, err := net.Listen("tcp4", host+":"+port)

	if err != nil {
		log_msgs.ErrorTimeLog(err.Error())
	}

	log_msgs.InfoLog("Listening on " + host + ":" + port)

	defer ln.Close()
	
	for {
		conn, err := ln.Accept()
		
		if err != nil {
			log_msgs.ErrorLog("An error happened in the for loop" + err.Error())
		}
<<<<<<< HEAD
		
		log_msgs.InfoTimeLog( "Connection received from: " + conn.RemoteAddr().String())
||||||| parent of 9ae7276 ([networking] work on sendMsg chunking)

		log_msgs.InfoTimeLog( "Connection received from: " + conn.RemoteAddr().String())
=======

		log_msgs.InfoTimeLog("Connection received from: " + conn.RemoteAddr().String())
>>>>>>> 9ae7276 ([networking] work on sendMsg chunking)

		tmp := make([]byte, 256)
		n, err := conn.Read(tmp)
		var metaData = strings.Trim(string(tmp[:n]), "\n")
		nameTarget := strings.Split(metaData, ":")
		log_msgs.InfoLog("Metadata received: " + metaData)
		log_msgs.InfoLog("Name: " + nameTarget[0] + ". Target: " + nameTarget[1])
		clients = append(clients, network.Connection{Name: nameTarget[0], Peer: nameTarget[1], Conn: conn})
		log_msgs.InfoLog("Last connection: " + clients[len(clients)-1].Conn.RemoteAddr().String())
		go handleConnection(clients[len(clients)-1])
	}
}

func handleConnection(client network.Connection) {
	streamMessages(client)
	defer closeConnection(client.Conn)
}

func streamMessages(client network.Connection) {
	tmp := make([]byte, network.ChunkSize * network.MaxChunks + 1) // using small tmo buffer for demonstrating
	for {
		n, err := client.Conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log_msgs.ErrorLog("read error:" + err.Error())
			}
			break
		}
		// fmt.Println("got", n, "bytes.")
		var msg = strings.Trim(string(tmp[:n]), "\n")
		// log_msgs.InfoLog("Msg from " + conn.RemoteAddr().String() + ": " + msg)

		for _, c := range clients {
			if client.Peer == c.Name {
				// c.Conn.Write([]byte(msg))
				network.SendMsg(c, msg)
			}
		}
	}
}

func closeConnection(conn net.Conn) {
	log_msgs.InfoTimeLog("Connection close from: " + conn.RemoteAddr().String())
	conn.Close()
}

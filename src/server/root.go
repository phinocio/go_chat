package server

import (
	// "io"
	"encoding/base64"
	"net"
	"strings"

	// "os"

	"go_chat/src/utils/log_msgs"
	"go_chat/src/utils/network"
)

var clients []network.Connection

func Run(host string, port string) {
	log_msgs.InfoLog(host + " " + port)
	ln, err := net.Listen("tcp4", host+":"+port)

	if err != nil {
		// log_msgs.ErrorLog(err.Error())
		log_msgs.ErrorTimeLog(err.Error())
	}

	log_msgs.InfoLog("Listening on " + host + ":" + port)

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log_msgs.ErrorLog("An error happened in the for loop" + err.Error())
		}
		// log_msgs.InfoTimeLog("Connection received from: " + conn.RemoteAddr().String())
		log_msgs.InfoLog("Connection received from: " + conn.RemoteAddr().String())
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
	for {
		var msg = network.RecvMsg(client.Conn)
		for _, c := range clients {
			if client.Peer == c.Name {
				// c.Conn.Write([]byte(msg))
				log_msgs.InfoLog("Sending msg from " + client.Name + " to " + c.Name )
				msg = append([]byte(client.Name + ":"), msg[0:]...)
				network.SendMsg(c.Conn, msg)
			}
		}
		log_msgs.InfoLog("Base64 of the payload: " + base64.StdEncoding.EncodeToString(msg))
	}
}

func closeConnection(conn net.Conn) {
	// log_msgs.InfoTimeLog("Connection close from: " + conn.RemoteAddr().String())
	log_msgs.InfoLog("Connection close from: " + conn.RemoteAddr().String())
	conn.Close()
}

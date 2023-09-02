package network

import (
	"go_chat/src/utils/log_msgs"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	ChunkSize = 8
	MaxChunks = 5
)

func getNumChunks(msg string) int {
	return len(msg) / ChunkSize
}

func RecvMsg(source net.Conn) {
	// pseude code
	// 1. recieve num chunks
	// 2. for loop of reading chunks until done
	// return msg
	
	// 1. recieve num chunks
	recv_buf := make([]byte, 256)
	n, err := source.Read(recv_buf)
	if err != nil {
		if err != io.EOF {
			log_msgs.ErrorLog("read error:" + err.Error())
		}
	}
	var x = string(recv_buf[:n])
	x = strings.Trim(x, "\x00\r\n")
	num_chunks, err := strconv.Atoi(x)
	if err != nil {
		log_msgs.ErrorLog("failed to convert chunk string to int")
	}
	
	// 2. for loop of reading chunks until done
	var msg []byte
	recv_buf = make([]byte, ChunkSize)
	for i := 0; i < num_chunks; i++ {
		n, err := source.Read(recv_buf)
		if err != nil {
			if err != io.EOF {
				log_msgs.ErrorLog("read error:" + err.Error())
			}
		}
		recv_buf := recv_buf[:n]
		msg = append(msg, recv_buf...)	
	}
	println(strings.TrimSpace(string(msg)))
}

func SendMsg(conn net.Conn, msg string) {
	// send a msg
}

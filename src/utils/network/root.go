package network

import (
	"fmt"
	"go_chat/src/utils/log_msgs"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	NumChunksByteSize = 8
	ChunkSize         = 1500
	MaxChunks         = 50
)

type Connection struct {
	Name string   // The name of this client instance
	Peer string   // The person we are trying to talk with
	Conn net.Conn // The client's connection
}

func getNumChunks(msg []byte) int {
	return len(msg)/ChunkSize + 1
}

func RecvMsg(source net.Conn) []byte {
	// pseude code
	// 1. recieve num chunks
	// 2. for loop of reading chunks until done
	// return msg

	// 1. recieve num chunks
	recv_buf := make([]byte, NumChunksByteSize)
	n, err := source.Read(recv_buf)
	if err != nil {
		if err != io.EOF {
			log_msgs.ErrorLog("read error:" + err.Error())
		}
	}

	var x = strings.Trim(string(recv_buf[:n]), "\x00\n")
	num_chunks, err := strconv.Atoi(x)
	if err != nil {
		log_msgs.ErrorLog("failed to convert chunk \"" + x + "\" string to int for " + source.RemoteAddr().String())
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
	// println(strings.TrimSpace(string(msg)))
	return msg
}

func SendMsg(target net.Conn, msg []byte) {
	// send a msg
	// log_msgs.InfoLog(msg)
	var chunks = getNumChunks(msg)
	if chunks > MaxChunks {
		log_msgs.ErrorLog("Msg too long bro. Got " + fmt.Sprint(chunks) + " chunks, expected max of " + fmt.Sprint(MaxChunks))
		return
	}
	//else {
	// log_msgs.InfoLog("Got " + fmt.Sprint(chunks) + " chunks")
	//}

	// Send number of chunks
	var num_chunks = []byte(strconv.Itoa(chunks))
	target.Write(num_chunks[:NumChunksByteSize])

	for i := 0; i < chunks; i++ {
		var start = i * ChunkSize
		var end = (i + 1) * ChunkSize
		var remainingLen = len(msg[start:])
		if remainingLen > ChunkSize {
			end = (i + 1) * ChunkSize
		} else {
			end = start + remainingLen
		}

		var tmp = make([]byte, ChunkSize)
		tmp = []byte(msg[start:end])
		target.Write(tmp)
	}
	// println()
}

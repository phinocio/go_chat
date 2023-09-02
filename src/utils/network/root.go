package network

import (
	"net"
)

const (
	ChunkSize = 8
	MaxChunks = 5
)

func getNumChunks(msg string) int {
	return len(msg) / ChunkSize
}

func RecvMsg(conn net.Conn, msg string) {
	// receive a msg
}

func SendMsg(conn net.Conn, msg string) {
	// send a msg
}

package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"go_chat/src/utils/colors"
	"go_chat/src/utils/encryption"
	"go_chat/src/utils/log_msgs"
	"go_chat/src/utils/network"
)

// Global Constants Avaiable to All go-routines
var global_prompt = colors.ColorWrap(colors.Purple, "[go_chat]> ")

var client_config encryption.H8go

func Run(host string, port string, nameTarget string) {
	log_msgs.InfoLog("client entry called")
	log_msgs.InfoTimeLog("client entry called")
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log_msgs.ErrorLog("failed to connect")
		os.Exit(1)
	}
	// define name:target
	log_msgs.InfoLog(nameTarget)
	log_msgs.InfoLog("myname: " + strings.Split(nameTarget, ":")[0])
	client_config = encryption.Load_Keys(strings.Split(nameTarget, ":")[0])
	log_msgs.InfoLog("targetname: " + strings.Split(nameTarget, ":")[1])
	writeToConn(conn, nameTarget)

	go readFromServer(conn)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	l, err := readline.NewEx(&readline.Config{
		// Prompt:          "\033[31m»\033[0m ",
		Prompt:          global_prompt,
		HistoryFile:     dir + "/example.history",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				continue // disable ctrl+c
				// break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case line == "exit":
			goto exit
		case line == "help":
			usage(l.Stderr())
		case line == "status":
			get_status(conn)
		case line == "":
		default:
			var msg []byte
			msg = encryption.Encryptor([]byte(line), client_config.Priv_key, client_config.Peers[0].Publ_key)
			log_msgs.InfoLog(base64.StdEncoding.EncodeToString(msg))
			network.SendMsg(conn, msg)
		}
	}
exit:
	conn.Close()
}

func get_status(conn net.Conn) {
	var remote_infos = strings.Split(conn.RemoteAddr().String(), ":")
	var remote_addr = remote_infos[0]
	var remote_port = remote_infos[1]
	fmt.Println()
	fmt.Println(colors.ColorWrap(colors.LightBlue, "\tAddr: ") + remote_addr)
	fmt.Println(colors.ColorWrap(colors.LightBlue, "\tPort: ") + remote_port)
	fmt.Println()
}

func writeToConn(conn net.Conn, line string) {
	var buffer = []byte(line + "\n")

	conn.Write(buffer)
}

func readFromServer(conn net.Conn) {
	log_msgs.InfoLog("Reading from server!")
	for {
		var msg = network.RecvMsg(conn)
		var decrypted []byte
		var encMsg = bytes.Split(msg, []byte(": "))
		decrypted = encryption.Decryptor(encMsg[1], client_config.Priv_key, client_config.Peers[0].Publ_key)
		log_msgs.InfoLog(base64.StdEncoding.EncodeToString(encMsg[1]))
		fmt.Println("")
		log_msgs.InfoLog("Msg from " + conn.RemoteAddr().String() + ": ")
		os.Stderr.WriteString("\n" + string(decrypted) + "\n\n")
		os.Stderr.WriteString(global_prompt)
	}

}

func usage(w io.Writer) {
	io.WriteString(w, "\ndefault behavior is to take input and write it to the server\n\n")
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
	io.WriteString(w, "\n")
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("exit"),
	readline.PcItem("help"),
	readline.PcItem("status"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

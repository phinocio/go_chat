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
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log_msgs.ErrorLog("failed to connect")
		os.Exit(1)
	}
	// define name:target
	log_msgs.InfoLog(nameTarget)
	var src_name = strings.Split(nameTarget, ":")[0]
	var dst_name = strings.Split(nameTarget, ":")[1]
	log_msgs.InfoLog("source is: " + src_name)
	log_msgs.InfoLog("destination is: " + dst_name)
	log_msgs.InfoLog("myname: " + src_name)
	client_config = encryption.Load_Keys(src_name)
	log_msgs.InfoLog("targetname: " + dst_name)
	writeToConn(conn, nameTarget)

	go readFromServer(conn, dst_name, src_name)		// can we dynamically change dst_name?

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	l, err := readline.NewEx(&readline.Config{
		// Prompt:          "\033[31m»\033[0m ",
		// Prompt:          global_prompt,				// TODO, go back to this
		Prompt:          "["+ src_name +"]> ",
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
			for _,v := range client_config.Peers {
				if v.Name == dst_name {
					msg = encryption.Encryptor([]byte(line), client_config.Priv_key, v.Publ_key)
					network.SendMsg(conn, msg)
					break
				}
			}
			// log_msgs.InfoLog(base64.StdEncoding.EncodeToString(msg))
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

func readFromServer(conn net.Conn, dst_name string, src_name string) {
	log_msgs.InfoLog("Reading from server!")
	for {
		var msg = network.RecvMsg(conn)
		var decrypted []byte
		var encMsg = bytes.Split(msg, []byte(":"))
		for _,v := range client_config.Peers {
			if v.Name == string(encMsg[0]) {
				decrypted = encryption.Decryptor(encMsg[1], client_config.Priv_key, v.Publ_key)
				break
			}
		}
		log_msgs.InfoLog(base64.StdEncoding.EncodeToString(encMsg[1]))
		fmt.Println("")
		log_msgs.InfoLog("Msg from " + conn.RemoteAddr().String() + ": ")
		os.Stderr.WriteString("\n" + string(decrypted) + "\n\n")
		// os.Stderr.WriteString(global_prompt)
		os.Stderr.WriteString("["+ src_name +"]> ")

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

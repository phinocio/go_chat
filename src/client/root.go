package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"go_chat/src/utils/log_msgs"
)

func Run(host string, port string) {
	log_msgs.InfoLog("client entry called")
	log_msgs.InfoTimeLog("client entry called")
    conn, err := net.Dial("tcp", host+":"+port)
    if err != nil {
			log_msgs.ErrorLog("failed to connect")
            os.Exit(1)
    }
    
    dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	l, err := readline.NewEx(&readline.Config{
		// Prompt:          "\033[31m»\033[0m ",
		Prompt:          "[go_chat]> ",
		HistoryFile:     dir  + "/example.history",
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
				break
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
                writeToConn(conn, line)
		}
	}
exit:
    // writeToConn(conn)
    conn.Close()
}

func get_status(conn net.Conn) {
	var remote_infos = strings.Split(conn.RemoteAddr().String(), ":")
	var remote_addr = remote_infos[0]
	var remote_port = remote_infos[1]
	fmt.Println()
	// fmt.Println("You are connected to:")
	fmt.Println("\tAddr: " + remote_addr)
	fmt.Println("\tPort: " + remote_port)
	fmt.Println()
}

func writeToConn(conn net.Conn, line string) {
    var buffer = []byte(line + "\n")
    conn.Write(buffer)
    // buffer := make([]byte, 1024)
    // _, err := conn.Read(buffer)
    // if err != nil {
    //         fmt.Println("failed to read the client connection")
    // }
    // fmt.Print(string(buffer))
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

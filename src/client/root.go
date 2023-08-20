package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

const (
    HOST = "localhost"
    PORT = "4444"
    TYPE = "tcp"
)

func Run(host string, port string) {
    fmt.Println("[INFO]  client entry called")

    conn, err := net.Dial("tcp", HOST+":"+PORT)
    if err != nil {
            fmt.Println("[ERROR] failed to connect")
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
		case line == "help":
			usage(l.Stderr())
		case strings.HasPrefix(line, "say"):
			line := strings.TrimSpace(line[3:])
			if len(line) == 0 {
				log.Println("say what?")
				break
			}
			go func() {
				for range time.Tick(time.Second) {
					log.Println(line)
				}
			}()
		case line == "exit":
			goto exit
		case line == "":
		default:
			// log.Println("you said:", strconv.Quote(line))
            writeToConn(conn)
		}
	}
exit:
    // writeToConn(conn)
    conn.Close()
}

func writeToConn(conn net.Conn) {
    conn.Write([]byte("hello server\n"))

    // buffer := make([]byte, 1024)
    // _, err := conn.Read(buffer)
    // if err != nil {
    //         fmt.Println("failed to read the client connection")
    // }
    // fmt.Print(string(buffer))
}

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("say",
		readline.PcItemDynamic(listFiles("./"),
			readline.PcItem("with",
				readline.PcItem("following"),
				readline.PcItem("items"),
			),
		),
		readline.PcItem("hello"),
		readline.PcItem("exit"),
	),
	readline.PcItem("exit"),
	readline.PcItem("help"),
	// readline.PcItem("go",
	// 	readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
	// 	readline.PcItem("install",
	// 		readline.PcItem("-v"),
	// 		readline.PcItem("-vv"),
	// 		readline.PcItem("-vvv"),
	// 	),
	// 	readline.PcItem("test"),
	// ),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

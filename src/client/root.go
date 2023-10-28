package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"go_chat/src/client/command"
	"go_chat/src/utils/colors"
	"go_chat/src/utils/encryption"
	"go_chat/src/utils/log_msgs"
	"go_chat/src/utils/network"
)

// Global Constants Avaiable to All go-routines
var global_prompt = colors.ColorWrap(colors.Purple, "[go_chat]> ")

var client_config encryption.H8go

// TODO: Extract bootstrap and code creation out of client into other modules/functions
func bootstrap(src_name string) {
	// create $XDG_CONFIG_HOME/go_chat if doesn't exist
	// create and store json file skeleton
	log_msgs.InfoLog("Bootstrapping client by creating directory...")

	confDir, err := os.UserConfigDir()
	if err != nil {
		log_msgs.ErrorLog("Error fetching user config directory")
		log.Fatal(err)
	}
	confDir += "/go_chat"

	if _, err := os.Stat(confDir); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(confDir, 0700)
			if err != nil {
				log_msgs.ErrorLog("Error creating go_chat directory")
				log.Fatal(err)
			}
		} else {
			log_msgs.ErrorLog("An error occurred checking if go_chat directory exists.")
			log.Fatal(err)
		}
	}

	// json skeleton
	// TODO: priv/pub key need to be generated on fly if file doesn't already exist.
	var userConfigFile = src_name + ".json"

	var fullConfigFilePath = confDir + "/" + userConfigFile

	if _, err := os.Stat(fullConfigFilePath); err != nil {
		if os.IsNotExist(err) {
			data := map[string]interface{}{
				"name":     src_name,
				"priv_key": "UkVDMgAAAC2WXSbNAMNzZBCJCD7EjJhEnKeAPASMDKTBOySyXqOrAL4VbXVc",
				"publ_key": "VUVDMgAAAC0VZx8oAzXCDUmNAD5oQAEqkxvxjpajjozZ+++FZzfxMeHDbvzm",
				"peers":    map[string]interface{}{},
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("could not marshal json: %s\n", err)
				return
			}
			fmt.Printf("json data: %s\n", jsonData)

			err2 := os.WriteFile(fullConfigFilePath, []byte(jsonData), 0600)
			if err2 != nil {
				log.Fatal(err)
			}
		} else {
			log_msgs.ErrorLog("An error occurred checking if go_chat directory exists.")
			log.Fatal(err)
		}
	}

	log_msgs.InfoLog("Bootstrapping compelted...")
}

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
	bootstrap(src_name)
	writeToConn(conn, nameTarget)

	go readFromServer(conn, dst_name, src_name) // can we dynamically change dst_name?

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	l, err := readline.NewEx(&readline.Config{
		// Prompt:          "\033[31m»\033[0m ",
		// Prompt:          global_prompt,				// TODO, go back to this
		Prompt:      "[" + src_name + "]> ",
		HistoryFile: dir + "/example.history",
		// AutoComplete:    completer,
		AutoComplete:    buildCompleter(client_config),
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
		cmd, args, _ := strings.Cut(line, " ")
		switch {
		case cmd == "exit":
			goto exit
		case cmd == "help":
			usage(l.Stderr(), client_config)
		case cmd == "status":
			command.Status(conn, dst_name)
			// command.status_command(conn, dst_name)
		case cmd == "peer":
			var dst_ptr *string
			dst_ptr = &dst_name
			command.Peer(args, client_config, dst_ptr)
		case cmd == "":
		default:
			var msg []byte
			for _, v := range client_config.Peers {
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
		for _, v := range client_config.Peers {
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
		os.Stderr.WriteString("[" + src_name + "]> ")

	}

}

func usage(w io.Writer, client_config encryption.H8go) {
	io.WriteString(w, "\ndefault behavior is to take input and write it to the server\n\n")
	io.WriteString(w, "commands:\n")
	io.WriteString(w, buildCompleter(client_config).Tree("    "))
	io.WriteString(w, "\n")
}

func buildCompleter(client_config encryption.H8go) *readline.PrefixCompleter {
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("exit"),
		readline.PcItem("help"),
		readline.PcItem("status"),
		readline.PcItem("peer",
			readline.PcItem("list"),
			readline.PcItem(
				"swap",
				readline.PcItemDynamic(client_config.GetPeerArrayGetter()),
			),
		),
	)
	return completer
}
func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

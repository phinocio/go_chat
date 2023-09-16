package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"encoding/base64"

	"github.com/chzyer/readline"

	"go_chat/src/utils/colors"
	"go_chat/src/utils/encryption"
	"go_chat/src/utils/log_msgs"
	"go_chat/src/utils/network"

	"github.com/cossacklabs/themis/gothemis/keys"
	"github.com/cossacklabs/themis/gothemis/message"
)

// Global Constants Avaiable to All go-routines
var global_prompt = colors.ColorWrap(colors.Purple, "[go_chat]> ")

var	aliceKeys = encryption.Gen_Keys()
var	bobKeys = encryption.Gen_Keys()

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
				continue 					// disable ctrl+c
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
				if (strings.Split(nameTarget, ":")[1] == "bob") {
					msg = encryption.Encryptor([]byte(line), aliceKeys.Private, bobKeys.Public)
				}
				if (strings.Split(nameTarget, ":")[1] == "alice") {
					msg = encryption.Encryptor([]byte(line), bobKeys.Private, aliceKeys.Public)
				}
				print(base64.StdEncoding.EncodeToString(msg))
				// encryption.Encryptor(line)
				network.SendMsg(conn, msg)
                // writeToConn(conn, line)
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
	fmt.Println(colors.ColorWrap(colors.LightBlue, "\tAddr: ") + remote_addr)
	fmt.Println(colors.ColorWrap(colors.LightBlue, "\tPort: ") + remote_port)
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

func gen_keys() ( *keys.Keypair) {
	alice_keyPair, err := keys.New(keys.TypeEC)
	if nil != err {
		fmt.Println("Keypair generating error")
		os.Exit(1)
	}
	return alice_keyPair
}

func readFromServer(conn net.Conn) {
	println("Reading from server!")
	// for {
		// tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
		for {
			// n, err := conn.Read(tmp)
			// if err != nil {
			// 	if err != io.EOF {
			// 		log_msgs.ErrorLog("read error:" + err.Error())
			// 	}
			// 	break
			// }
			// // fmt.Println("got", n, "bytes.")
			// var msg = strings.Trim(string(tmp[:n]), "\n")
			var msg = network.RecvMsg(conn)
			var whoIsThePrependedTag = bytes.Split(msg, []byte(": "))
			print(base64.StdEncoding.EncodeToString(whoIsThePrependedTag[1]))
			if (string(whoIsThePrependedTag[0]) == "[bob]") {
				msg = encryption.Decryptor(msg, bobKeys.Private, aliceKeys.Public)
			}
			if (string(whoIsThePrependedTag[0]) == "[alice]") {
				msg = encryption.Decryptor(whoIsThePrependedTag[1], aliceKeys.Private, bobKeys.Public)
			}

			fmt.Println("")
			log_msgs.InfoLog("Msg from " + conn.RemoteAddr().String() + ": ")
			os.Stderr.WriteString("\n" + string(msg) + "\n\n")
			os.Stderr.WriteString(global_prompt)
			// println(msg)
		}
	// }
}

func encryptor(clear_text []byte, my_priv_key *keys.PrivateKey, peer_publ_key *keys.PublicKey) ([]byte) {
	aliceToBob := message.New(my_priv_key, peer_publ_key)
	cipher_text, err := aliceToBob.Wrap(clear_text)
	if err != nil {
		fmt.Println("message cannot be empty")
	}
	return cipher_text
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

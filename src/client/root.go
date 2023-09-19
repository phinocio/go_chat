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
	"github.com/cossacklabs/themis/gothemis/keys"

	"go_chat/src/utils/colors"
	"go_chat/src/utils/encryption"
	"go_chat/src/utils/log_msgs"
	"go_chat/src/utils/network"
)

// Global Constants Avaiable to All go-routines
var global_prompt = colors.ColorWrap(colors.Purple, "[go_chat]> ")

var gosux, _ = base64.StdEncoding.DecodeString("VUVDMgAAAC0VZx8oAzXCDUmNAD5oQAEqkxvxjpajjozZ+++FZzfxMeHDbvzm")
var alicePublicKey = &keys.PublicKey{
	Value: gosux,
}
var gosux2, _ = base64.StdEncoding.DecodeString("UkVDMgAAAC2WXSbNAMNzZBCJCD7EjJhEnKeAPASMDKTBOySyXqOrAL4VbXVc")
var alicePrivateKey = &keys.PrivateKey{
	Value: gosux2,
}

var gosux3, _ = base64.StdEncoding.DecodeString("VUVDMgAAAC3+kngIAkAdHfaub4y5+VVHZglc/8+oJ7nBpwUpxvH4EOzuQCbS")
var bobPublicKey = &keys.PublicKey{
	Value: gosux3,
}
var gosux4, _ = base64.StdEncoding.DecodeString("UkVDMgAAAC2JIg6cANXF+T4cntocTUjWvO9Z8VPUgE1N3eouBrb9gGOupKyF")
var bobPrivateKey = &keys.PrivateKey{
	Value: gosux4,
}

type self_config struct {
	Name string `json:"name"`
	Priv_key string `json:"priv_key"`
	Publ_key string `json:"publ_key"`
}
type peer_config struct {
	Name string `json:"name"`
	Publ_key string `json:"publ_key"`
}
type config_pack struct {
	Self_config self_config `json:"self"`
	Peer_config peer_config `json:"peer"`
}
type CONFIG_PACK interface {
	debug_print()
}
func (self config_pack) debug_print() {
	fmt.Println("SELF")
	fmt.Println(self.Self_config.Name)
	fmt.Println(self.Self_config.Priv_key)
	fmt.Println(self.Self_config.Publ_key)
	fmt.Println("")
	fmt.Println("PEER")
	fmt.Println(self.Peer_config.Name)
	fmt.Println(self.Peer_config.Publ_key)
}
func load_config_file(filename string) (config_pack) {
	var result config_pack

	b, err := os.ReadFile(filename) // just pass the file name
    if err != nil {
        fmt.Print(err)
    }

	json.Unmarshal(b, &result)		// unmarshal means convert to struct

	return result
}
func we_developing(){
	var config_file_name_1 = "bob.json"
	var config_file_name_2 = "alice.json"
	
	var conn_pack_1 = load_config_file(config_file_name_1)
	conn_pack_1.debug_print()

	fmt.Print("\n\n")

	var conn_pack_2 =  load_config_file(config_file_name_2)
	conn_pack_2.debug_print()
	
	os.Exit(1)
}




func Run(host string, port string, nameTarget string) {
	//
	we_developing()	// Gadget function for developing
	//
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
					// msg = encryption.Encryptor([]byte(line), aliceKeys.Private, bobKeys.Public) // ORIGINAL
					msg = encryption.Encryptor([]byte(line), alicePrivateKey, bobPublicKey) 
				}
				if (strings.Split(nameTarget, ":")[1] == "alice") {
					// msg = encryption.Encryptor([]byte(line), bobKeys.Private, aliceKeys.Public)	// ORIGINAL
					msg = encryption.Encryptor([]byte(line), bobPrivateKey, alicePublicKey) // ORIGINAL
				}
				log_msgs.InfoLog(base64.StdEncoding.EncodeToString(msg))
				var decrypted []byte
				if (strings.Split(nameTarget, ":")[1] == "bob") {
					// decrypted = encryption.Decryptor(msg, bobKeys.Private, aliceKeys.Public)	// ORIGINAL
					decrypted = encryption.Decryptor(msg, bobPrivateKey, alicePublicKey)	
				}
				if (strings.Split(nameTarget, ":")[1] == "alice") {
					// decrypted = encryption.Decryptor(msg, aliceKeys.Private, bobKeys.Public)	// ORIGINAL
					decrypted = encryption.Decryptor(msg, alicePrivateKey, bobPublicKey)	
				}
				// encryption.Encryptor(line)
				log_msgs.InfoLog(base64.StdEncoding.EncodeToString(decrypted))
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
}


func readFromServer(conn net.Conn) {
	log_msgs.InfoLog("Reading from server!")
	for {
		var msg = network.RecvMsg(conn)
		var decrypted []byte
		var whoIsThePrependedTag = bytes.Split(msg, []byte(": "))
		log_msgs.InfoLog(base64.StdEncoding.EncodeToString(whoIsThePrependedTag[1]))
		if (string(whoIsThePrependedTag[0]) == "[bob]") {
			// decrypted = encryption.Decryptor(whoIsThePrependedTag[1], bobKeys.Private, aliceKeys.Public)	// ORIGINAL
			decrypted = encryption.Decryptor(whoIsThePrependedTag[1], bobPrivateKey, alicePublicKey)	
		}
		if (string(whoIsThePrependedTag[0]) == "[alice]") {
			// decrypted = encryption.Decryptor(whoIsThePrependedTag[1], aliceKeys.Private, bobKeys.Public)	// ORIGINAL
			decrypted = encryption.Decryptor(whoIsThePrependedTag[1], alicePrivateKey, bobPublicKey)	
		}

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

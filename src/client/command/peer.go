package command

import (
	"fmt"
	"go_chat/src/utils/colors"
	"go_chat/src/utils/encryption"
	"go_chat/src/utils/log_msgs"
	"strings"
)

func Peer(args string, client_config encryption.H8go, dst_ptr *string) {

	subcmd, args, _ := strings.Cut(args, " ")
	switch {
		case subcmd == "":
			peer_usage()
		case subcmd == "list":
			peer_list(client_config)
		case subcmd == "swap":
			if !client_config.Peer_exists(args) {
				log_msgs.ErrorLog("peer not in config")
				return
			}
			peer_swap(args, dst_ptr)
	}
}
func peer_usage() {
	usage:= `
	peer list 		// list available peers
	peer swap bob	// swap to chatting with bob
	`
	fmt.Println(colors.ColorWrap(colors.Yellow, usage))
}
func peer_list(client_config encryption.H8go) {
	fmt.Println()
	fmt.Println(colors.ColorWrap(colors.LightBlue, "Avalaible peers:"))
	
	for _,v := range client_config.Peers {
		fmt.Println("\t" + v.Name)
	}
	
	fmt.Println()
}
func peer_swap(args string, dst_ptr *string) {
	log_msgs.WarnLog("POINTER ARE BEING USED, WOAH DUDE")
	*dst_ptr = args
	fmt.Println()
	fmt.Println(colors.ColorWrap(colors.LightBlue, "\tPeer: ") + *dst_ptr)
	fmt.Println()
}
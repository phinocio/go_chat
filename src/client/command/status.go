package command

import (
	"fmt"
	"go_chat/src/utils/colors"
	"net"
	"strings"
)



func Status(conn net.Conn, dst_name string) {
	var remote_infos = strings.Split(conn.RemoteAddr().String(), ":")
	var remote_addr = remote_infos[0]
	var remote_port = remote_infos[1]
	fmt.Println()
	fmt.Println(colors.ColorWrap(colors.LightBlue, "\tAddr: ") + remote_addr)
	fmt.Println(colors.ColorWrap(colors.LightBlue, "\tPort: ") + remote_port)
	fmt.Println(colors.ColorWrap(colors.LightBlue, "\tPeer: ") + dst_name)
	fmt.Println()
}
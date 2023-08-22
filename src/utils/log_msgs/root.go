package log_msgs

import (
	"fmt"
	"log"

	"go_chat/src/utils/colors"
)

const (
	info_prefix  = colors.LightBlue + "[INFO]  " + colors.ResetColor
	warn_prefix  = colors.Yellow + "[WARN]  " + colors.ResetColor
	error_prefix = colors.Red + "[ERROR] " + colors.ResetColor
)

func InfoLog(info_details string) {
	fmt.Println( info_prefix + info_details )
}
func InfoTimeLog(info_details string) {
	log.Println( info_prefix + info_details )
}

func WarnLog(_details string) {
	fmt.Println( warn_prefix + _details )
}
func WarnTimeLog(_details string) {
	log.Println( warn_prefix + _details )
}

func ErrorLog(_details string) {
	fmt.Println( error_prefix + _details )
}
func ErrorTimeLog(_details string) {
	log.Println( error_prefix + _details )
}

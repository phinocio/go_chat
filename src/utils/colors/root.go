package colors

// https://gist.github.com/ik5/d8ecde700972d4378d87
// https://gist.github.com/JBlond/2fea43a3049b38287e5e9cefc87b2124
//
// WARNING: \e in bash is \033 in ansi
//
const (
	LightBlue 	= "\033[1;36m"
	Purple  	= "\033[0;35m"
	Red 		= "\033[1;31m"
	ResetColor 	= "\033[0m"
	Yellow 		= "\033[1;33m"
)

func ColorWrap( color string, message string ) string {
	// var result = color + message + ResetColor
	// return result
	return color + message + ResetColor
}

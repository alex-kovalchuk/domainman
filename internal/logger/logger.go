package logger

import "fmt"

var (
	ColorInfo    = "\033[1;34m%s\033[0m"
	ColorNotice  = "\033[1;36m%s\033[0m"
	ColorWarning = "\033[1;33m%s\033[0m"
	ColorError   = "\033[1;31m%s\033[0m"
	ColorDebug   = "\033[0;36m%s\033[0m"
)

func Print(text string, color string) {
	fmt.Println(fmt.Sprintf(color, text))
}

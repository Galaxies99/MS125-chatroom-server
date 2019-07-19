package chatroom

import (
	"fmt"
	"time"
)

//something about color, no need modifying
const (
	CHAT_LOG    = "[LOG %s] %s\n"
	CHAT_ERR    = "[ERR %s] %s\n"
	CHAT_MSG    = "[MSG %s] %s\n"
	COLOR_RED   = "\033[31m"
	COLOR_GREEN = "\033[32m"
	COLOR_RESET = "\033[0m"
)

//get current time and return a string
func GetCurrentTimeString() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006/01/02 15:04:05")
}

//print a log info
func PrintLog(msg string) {
	fmt.Printf(CHAT_LOG, GetCurrentTimeString(), msg)
}

//print an error info
func PrintErr(msg string) {
	fmt.Printf(COLOR_RED+CHAT_ERR+COLOR_RESET, GetCurrentTimeString(), msg)
}

//print message receive from some client
func PrintClientMsg(msg string) {
	fmt.Printf(COLOR_GREEN+CHAT_MSG+COLOR_RESET, GetCurrentTimeString(), msg)
}

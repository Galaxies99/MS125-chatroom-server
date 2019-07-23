package main

import (
	"./src/chatroom"
)

func main() {
	server := chatroom.NewChatServer("0.0.0.0", 6666)
	server.StartListen()
}

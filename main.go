package main

import (
	"./src/chatroom"
)

func main() {
	server := chatroom.NewChatServer("127.0.0.1", 6666)
	server.StartListen()
}

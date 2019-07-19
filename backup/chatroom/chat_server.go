package chatroom

import (
	"net"
	"strconv"
)

const (
	LISTEN_TCP = "tcp"
	PING_MSG   = "receive connection from "
)

var users []net.Conn
var userNum int
var userName []string

//data structure of server
type ChatServer struct {
	listenAddr string
	status     bool
	listener   net.Listener
}

//create a new server, you should explain why we do this
func NewChatServer(addr string, port int) *ChatServer {
	server := new(ChatServer)
	server.listenAddr = addr + ":" + strconv.Itoa(port)
	return server
}

//main server function
func (server ChatServer) StartListen() {
	//start listen on the address given
	listener, err := net.Listen(LISTEN_TCP, server.listenAddr)
	server.listener = listener

	users = make([]net.Conn, 15)
	userName = make([]string, 15)
	userNum = 0

	//exit when server listen fail
	if err != nil {
		PrintErr(err.Error())
	} else {
		PrintLog("Start Listen " + server.listenAddr)
	}

	//main server loop, you should explain how this server loop works
	for {
		client, acceptError := server.listener.Accept() //when a user comes in
		if acceptError != nil {
			PrintErr(acceptError.Error()) //if accept go wrong, then the server will exit
			break
		} else {
			go server.userHandler(client) //then create a coroutine go process the user (What is coroutine? What's the function of keyword 'go'?)
		}
	}
}

func (server ChatServer) userHandler(client net.Conn) {
	buffer := make([]byte, 1024)      //create a buffer
	clientAddr := client.RemoteAddr() //get user's address
	var msg string
	var userID int
	for userID = 0; userID < userNum; userID++ {
		if client == users[userID] {
			break
		}
	}
	if userID == userNum {
		users[userNum] = client
		userName[userNum] = client.RemoteAddr().String()
		userID = userNum
		userNum++
	}
	PrintClientMsg(PING_MSG + clientAddr.String() + " (" + userName[userID] + ").") //print a log to show that a new client comes in
	for i := 0; i < userNum; i++ {
		if i != userID {
			users[i].Write([]byte(userName[userID] + " log in."))
		} else {
			users[i].Write([]byte("Successfully log in, your name is " + userName[userID] + ", you can use \"SetName@yourname\" to set your name."))
		}
	}
	for {                                          //main serve loop
		readSize, readError := client.Read(buffer) //why we can put a read in for loop?
		if readError != nil {
			PrintErr(clientAddr.String() + " fail") //if read error occurs, close the connection with user
			client.Close()
			break
		} else {
			msg = string(buffer[0:readSize])

			var id int
			for i := 0; i < userNum; i++ {
				if users[i] == client {
					id = i
					break
				}
			}

			PrintClientMsg(clientAddr.String() + " (" + userName[id] + "): " + msg)

			if readSize >= 8 && msg[0:8] == "SetName@" {
				nameEnd := 8
				for ; nameEnd < readSize; nameEnd ++ {
					if msg[nameEnd] == ' ' || msg[nameEnd] == '\t' || msg[nameEnd] == '\n' {
						break
					}
				}

				blocked := false
				for i := 0; i < userNum; i++ {
					if i != id && msg[8:nameEnd] == userName[i] {
						blocked = true
						break
					}
				}
				if blocked == false {
					userName[id] = msg[8:nameEnd]
					client.Write([]byte("Server: Change your name to " + msg[8:nameEnd] + " successfully."))
				} else {
					client.Write([]byte("Server: Invalid Name"))
				}
				continue
			}

			// spread message to all
			toPrint := []byte(userName[id] + ": " + msg)

			for i := 0; i < userNum; i++ {
				if users[i] != client {
					users[i].Write(toPrint)
				}
			}
		}
	}
}

package chatroom

import (
	"net"
	"strconv"
)

const (
	LISTEN_TCP = "tcp"
	PING_MSG   = "receive connection from "
)

var dataLib map[string]string
var activePort []net.Conn
var activePortNum int

type ChatServer struct {
	listenAddr string
	status     bool
	listener   net.Listener
}

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

	dataLib = make(map[string]string)
	activePort = make([]net.Conn, 30000)
	activePortNum = 0

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

func checkifLogin(msg string) (bool) {
	if len(msg) < 15 {
		return false
	}
	if msg[0:15] == "~@LoginRequest#" {
		return true
	} else {
		return false
	}
}

func checkifTourist(msg string) (bool) {
	if len(msg) < 15 {
		return false
	}
	if(msg[15:] == "Tourist") {
		return true
	} else {
		return false
	}
}

func checkifRegister(msg string) (bool) {
	if len(msg) < 18 {
		return false
	}
	if msg[0:18] == "~@RegisterRequest#" {
		return true
	} else {
		return false
	}
}

func checkifChange(msg string) (bool) {
	if len(msg) < 16 {
		return false
	}
	if msg[0:16] == "~@ChangeRequest#" {
		return true
	} else {
		return false
	}
}

func (server ChatServer) userHandler(client net.Conn) {
	buffer := make([]byte, 1024)
	clientAddr := client.RemoteAddr()
	clientType := -1
	PrintClientMsg(PING_MSG + clientAddr.String())

	found := false
	for i := 0; i < activePortNum; i++ {
		if activePort[i] == client {
			found = true
			break
		}
	}
	if found == false {
		activePort[activePortNum] = client
		activePortNum ++
	}

	var msg string
	for {
		readSize, readError := client.Read(buffer)
		if readError != nil {
			PrintErr(clientAddr.String() + " fail")
			for i := 0; i < activePortNum; i++ {
				if activePort[i] == client {
					activePort[i] = activePort[activePortNum - 1]
					activePortNum --
					break
				}
			}
			client.Close()
			break
		} else {
			msg = string(buffer[0:readSize])

			PrintClientMsg(clientAddr.String() + ": " + msg)

			if clientType == -1 {
				// check type
				if checkifLogin(msg) {
					clientType = 2
				} else {
					if checkifRegister(msg) {
						clientType = 3
					} else {
						if checkifChange(msg) {
							clientType = 4;
						} else {
							clientType = 1
						}
					}
				}
			}

			if clientType == 1 {
				// normal user, send message.
				toPrint := []byte(msg)
				for i := 0; i < activePortNum; i++ {
					if activePort[i] != client {
						activePort[i].Write(toPrint)
					}
				}
			}

			if clientType == 2 {
				// login request.
				if checkifTourist(msg) {
					client.Write([]byte("Accept"))
				} else {
					msg = msg[15:]
					k := 0
					for ; k < len(msg); k++ {
						if msg[k] == '#' {
							break
						}
					}
					usrName := msg[0:k]
					pwd := msg[k+1:]
					usrpwd, ok := dataLib[usrName]
					if ok {
						if pwd == usrpwd {
							client.Write([]byte("Accept"))
						} else {
							client.Write([]byte("WrongPwd"))
						}
					} else {
						client.Write([]byte("NoUser"))
					}
				}
			}

			if clientType == 3 {
				// register request.
				msg = msg[18:]
				k := 0
				for ; k < len(msg); k++ {
					if msg[k] == '#' {
						break
					}
				}
				usrName := msg[0:k]
				pwd := msg[k+1:]
				_, ok := dataLib[usrName]
				if ok {
					client.Write([]byte("UserExists"))
				} else {
					dataLib[usrName] = pwd;
					client.Write([]byte("Accept"))
				}
			}

			if clientType == 4 {
				// change
				msg = msg[16:]
				k := 0
				for ; k < len(msg); k++ {
					if msg[k] == '#' {
						break
					}
				}
				oldusrName := msg[0:k]
				msg = msg[k+1:]
				k = 0
				for ; k < len(msg); k++ {
					if msg[k] == '#' {
						break
					}
				}
				newusrName := msg[0:k]
				msg = msg[k+1:]
				k = 0
				for ; k < len(msg); k++ {
					if msg[k] == '#' {
						break
					}
				}
				oldpwd := msg[0:k]
				newpwd := msg[k+1:]
				pwd, ok := dataLib[oldusrName]
				if ok == false {
					client.Write([]byte("NoUser"))
				} else {
					if pwd != oldpwd {
						client.Write([]byte("WrongPwd"))
					} else {
						_, ok2 := dataLib[newusrName]
						if ok2 == true {
							client.Write([]byte("UserExists"))
						} else {
							delete(dataLib, oldusrName)
							dataLib[newusrName] = newpwd;
							for i := 0; i < activePortNum; i++ {
								if activePort[i] != client {
									activePort[i].Write([]byte("<" + oldusrName + "> change his/her name to <" + newusrName + ">."))
								}
							}
							client.Write([]byte("Accept"))
						}
					}
				}
			}
		}
	}
}

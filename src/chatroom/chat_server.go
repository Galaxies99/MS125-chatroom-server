package chatroom

import (
	"net"
	"strconv"
)

const (
	LISTEN_TCP = "tcp"
	PING_MSG   = "receive connection from "
)

type userLib struct {
	pwd	string
	pts int
}

var dataLib map[string]userLib
var activePort []net.Conn
var activePortNum int
var TouristID int

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
	TouristID = 0

	dataLib = make(map[string]userLib)
	activePort = make([]net.Conn, 30000)
	activePortNum = 0

	// for test
	dataLib["test1"] = userLib{"test", 0}
	dataLib["test2"] = userLib{"test", 100}
	dataLib["test3"] = userLib{"test", 1000}


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

func checkifLoginRequest(msg string) (bool) {
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

func checkifUnlog(msg string) (bool) {
	if len(msg) < 8 {
		return false
	}
	if msg[0:8] == "~@Unlog#" {
		return true
	} else {
		return false
	}
}

func checkifLogin(msg string) (bool) {
	if len(msg) < 8 {
		return false
	}
	if msg[0:8] == "~@Login#" {
		return true
	} else {
		return false
	}
}

func getLevel(pts int) (string) {
	if pts < 5 {
		return "Level 1"
	}
	if pts < 20 {
		return "Level 2"
	}
	if pts < 50 {
		return "Level 3"
	}
	if pts < 100 {
		return "Level 4"
	}
	if pts < 200 {
		return "Level 5"
	}
	if pts < 500 {
		return "Level 6"
	}
	if pts < 1000 {
		return "Level 7"
	}
	if pts < 2000 {
		return "Level 8"
	}
	if pts < 5000 {
		return "Level 9"
	}
	return "Level 10"
}

func (server ChatServer) userHandler(client net.Conn) {
	buffer := make([]byte, 1024)
	clientAddr := client.RemoteAddr()
	clientType := -1
	PrintClientMsg(PING_MSG + clientAddr.String())

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
				if checkifLoginRequest(msg) {
					clientType = 2
				} else {
					if checkifRegister(msg) {
						clientType = 3
					} else {
						if checkifChange(msg) {
							clientType = 4
						} else {
							clientType = 1

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
						}
					}
				}
			}

			if clientType == 1 {
				if checkifUnlog(msg) {
					msg = msg[8:]
					toPrint := "<" + msg + "> log out."
					for i := 0; i < activePortNum; i++ {
						if activePort[i] != client {
							activePort[i].Write([]byte(toPrint))
						}
					}
				} else if checkifLogin(msg) {
					msg = msg[8:]
					if len(msg) >= 7 && msg[0:7] == "Tourist" {
						toPrint := "<" + msg + "> log in."
						for i := 0; i < activePortNum; i++ {
							activePort[i].Write([]byte(toPrint))
						}
					} else {
						toPrint := "<" + msg + "> log in."
						for i := 0; i < activePortNum; i++ {
							if activePort[i] != client {
								activePort[i].Write([]byte(toPrint))
							} else {
								usrData := dataLib[msg]
								activePort[i].Write([]byte("~@" + getLevel(usrData.pts) + ", Coins: " + strconv.Itoa(usrData.pts) + "#" + toPrint))
							}
						}
					}
				} else {
					i := 0
					for ; i < len(msg); i++ {
						if msg[i] == '#' {
							break
						}
					}
					usrName := msg[0:i]
					if len(usrName) >= 7 && usrName[0:7] == "Tourist" {
						toPrint := usrName + "  " + GetCurrentTimeString() + "#" + msg[i+1:]
						for i := 0; i < activePortNum; i++ {
							activePort[i].Write([]byte(toPrint))
						}
					} else {
						usrData := dataLib[usrName]
						usrData.pts++
						toPrint := usrName + " (" + getLevel(usrData.pts) + ") " + GetCurrentTimeString() + "#" + msg[i+1:]
						dataLib[usrName] = usrData
						for i := 0; i < activePortNum; i++ {
							if activePort[i] != client {
								activePort[i].Write([]byte(toPrint))
							} else {
								tmp := "~@" + getLevel(usrData.pts) + ", Coins: " + strconv.Itoa(usrData.pts)
								activePort[i].Write([]byte(tmp + "#" + toPrint))
							}
						}
					}
				}
			}

			if clientType == 2 {
				// login request.
				if checkifTourist(msg) {
					TouristID ++
					client.Write([]byte("Accept#" + "Tourist" + strconv.Itoa(TouristID)))
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
					usrData, ok := dataLib[usrName]
					if ok {
						if pwd == usrData.pwd {
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
				usrPwd := msg[k+1:]
				_, ok := dataLib[usrName]
				if ok {
					client.Write([]byte("UserExists"))
				} else {
					dataLib[usrName] = userLib{pwd: usrPwd, pts: 0}
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
				usrData, ok := dataLib[oldusrName]
				if ok == false {
					client.Write([]byte("NoUser"))
				} else {
					if usrData.pwd != oldpwd {
						client.Write([]byte("WrongPwd"))
					} else {
						_, ok2 := dataLib[newusrName]
						if ok2 == true {
							client.Write([]byte("UserExists"))
						} else {
							delete(dataLib, oldusrName)
							dataLib[newusrName] = userLib{pwd: newpwd, pts: usrData.pts}
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

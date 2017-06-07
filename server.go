package main


import (
	//"strings"
	//"encoding/json"
	"net"
	"fmt"
	"os"
	"encoding/json"
	_ "math/rand"
)

type MessageType uint8


type Client struct {
	clientID int
}

var clientList []Client

const (
	CONNECTION MessageType = iota + 1
	SEND_MESSAGE
	QUIT
	CHECK
)

type Message struct {
	ClientID   	int    		`json:"clientId"`
	TypeOfMessage   MessageType 	`json:"typeOfMessage"`
	Message   	string 		`json:"message"`
}


/* A Simple function to verify error */
func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}
}

func main() {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr,err := net.ResolveUDPAddr("udp",":10001")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n,addr,err := ServerConn.ReadFromUDP(buf)
		mesJson := string(buf[0:n])
		//fmt.Println("Received ",mesJson, " from ",addr)

		var mes Message

		b := []byte(mesJson)
		err = json.Unmarshal(b, &mes)

		if err != nil {
			fmt.Println("Json Error: ",err)
			os.Exit(1)
		}

		if CONNECTION == mes.TypeOfMessage{
			if !clientIsExists(mes.ClientID) {
				var client Client
				client.clientID = mes.ClientID
				clientList = append (clientList, client)
				fmt.Println("New Client: ", mes.ClientID, " messageType ",mes.TypeOfMessage, " Message ", mes.Message)
			} else {
				fmt.Println("Client with this ID already exists")
			}
		}

	}
}

func clientIsExists( clientID int) (bool){
	for _, v := range clientList {
		if v.clientID == clientID {
			return true
		}
	}
	return false
}
package main


import (
	//"strings"
	//"encoding/json"
	"net"
	"fmt"
	"os"
	"encoding/json"
	_ "math/rand"
	"bufio"
	"strings"
	"golang.org/x/crypto/openpgp/errors"
)

type MessageType uint8


type Client struct {
	ipAddr *net.UDPAddr
	id string
}

var clientList []Client

var ServerAddr *net.UDPAddr
var ServerConn *net.UDPConn

const (
	CONNECTION MessageType = iota + 1
	SEND_MESSAGE
	QUIT
	CHECK
)

type Message struct {
	ClientIP   	string    	`json:"clientIP"`
	ClientID	string		`json:"clientID"`
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
	var err error

	ServerAddr,err = net.ResolveUDPAddr("udp",":10001")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	go readConsole ()
	for {
		n,addr,err := ServerConn.ReadFromUDP(buf)
		mesJson := string(buf[0:n])
		fmt.Println("Received ",mesJson, " from ",addr)

		var mes Message

		b := []byte(mesJson)
		err = json.Unmarshal(b, &mes)

		if err != nil {
			fmt.Println("Json Error: ",err)
			os.Exit(1)
		}

		switch mes.TypeOfMessage {
		case CONNECTION:
			{
				if !clientIsExists(mes.ClientIP) {
					var client Client
					client.ipAddr = addr
					client.id = mes.ClientID
					clientList = append(clientList, client)
					fmt.Println("New Client: ", mes.ClientIP, " messageType ", mes.TypeOfMessage, " Message ", mes.Message)
				} else {
					fmt.Println("Client with this ID already exists")
					b := []byte("Client already exists")
					err = json.Unmarshal(b, &mes)
				}
			}
		case CHECK:{
			sendResponse(addr)
			}
		case SEND_MESSAGE: {
			fmt.Println("Client " , addr.IP.String(), " Snet message: n", mes.Message)
		}
		}
	}
}

func sendResponse(addr *net.UDPAddr) {
	_,err := ServerConn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
		cli, _, _:= getClientByIp(addr.IP.String())
		deleteClientFromList(cli.id)
	}
}

func clientIsExists( clientIP string) (bool){
	for _, v := range clientList {
		if v.ipAddr.IP.String() == clientIP {
			return true
		}
	}
	return false
}

func printClients(){
	for i, element := range clientList {
		fmt.Println("Client ID ", i, " Client IP ", element.ipAddr)
	}
}

func deleteClientFromList(clientID string){
	_, i, err := getClientByID(clientID)
	if err == nil {
		clientList = append(clientList[:], clientList[i+1:]...)
	}
}

func getClientByID(clientID string) (Client, int, error) {
	var cli Client
	var i int
	for i, cli = range clientList {
		if cli.id == clientID {
			return cli, i, nil
		}
	}
	return cli, -1, errors.ErrUnknownIssuer
}

func manualDisconnectClient() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter ID  of client")
	clientID, _ := reader.ReadString('\n')
	client, _, _ := getClientByID(clientID)

	deleteClientFromList(client.id)
	var msg Message
	msg.TypeOfMessage = QUIT
	msg.ClientID = client.id
	msg.ClientIP  = client.ipAddr.IP.String()

	msgJson, _ := json.Marshal(msg)
	buf := []byte(msgJson)
	ServerConn.WriteToUDP(buf, client.ipAddr)
}

func getClientByIp(ipaddr string) (Client, int, error) {
	var cli Client
	var i int
	for i, cli = range clientList {
		if cli.ipAddr.IP.String() == ipaddr {
			return cli, i, nil
		}
	}
	return cli, -1, errors.ErrUnknownIssuer
}

func sendMessageToAllClient( message string){

	for _, cli := range clientList {
		sendMessageToClient(cli, message)
	}
}


func sendMessageToClient(cli Client, msgS string) {

	var msg Message
	msg.TypeOfMessage = SEND_MESSAGE
	msg.ClientID = cli.id
	msg.ClientIP  = cli.ipAddr.IP.String()
	msg.Message = msgS
	msgJson, _ := json.Marshal(msg)
	buf := []byte(msgJson)
	ServerConn.WriteToUDP(buf, cli.ipAddr)

}

func readConsole() ( exit bool,  msg Message) {

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message text: ")
		text, _ := reader.ReadString('\n')

		service := strings.Index(text, "/")
		if (service == 0) {
			message := strings.ToUpper(text)

			switch  message {
			case "/QUIT\n":
				{
					msg.TypeOfMessage = QUIT
					fmt.Println(exit)
					exit = true
					break
				}
			case "/GETCLIENTS\n":
				{
					printClients()
					break
				}
			case "/DISC\n":  //disconnect
				{
					manualDisconnectClient()
				}

			}

		}
	}
}
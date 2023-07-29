package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	mod        int
}

func NewClient(serverip string, serverport int) *Client {
	//create the object
	client := &Client{
		ServerIp:   serverip,
		ServerPort: serverport,
		mod:        999,
	}
	// connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverip, serverport))
	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return nil
	}
	client.conn = conn
	return client
}

// display the content of the conn on terminal
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

// displag the content of the menu
func (client *Client) menu() bool {
	var chose int
	fmt.Println(">>>1.public chat")
	fmt.Println(">>>2.private chat")
	fmt.Println(">>>3.rename")
	fmt.Println(">>>0.exit")

	fmt.Scanln(&chose)
	if chose >= 0 && chose <= 3 {
		client.mod = chose
		return true
	} else {
		fmt.Println("chose err!")
		return false
	}
}

// choose to user to chat privately
func (client *Client) SelectUsers() {
	// show users table
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return
	}
}

// do private chat
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>please input the name who you want to chat with (input exit is to quit)")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>input your content to send (input exit is to quit)")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err: ", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>input your content to send (input exit is to quit)")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println(">>>please input the name who you want to chat with (input exit is to quit)")
		fmt.Scanln(&remoteName)
	}
}

// do public chat
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>please input your content to chat (input exit is to quit)")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		//send to the server
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err: ", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>please input your content to chat, exit mean quit")
		fmt.Scanln(&chatMsg)
	}
}

// update the new name to server
func (client *Client) UpdateName() bool {
	fmt.Println(">>>please input the new name:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}
	return true
}

// run the client
func (client *Client) Run() {
	for client.mod != 0 {
		for client.menu() != true {
		}

		switch client.mod {
		case 1:
			fmt.Println(">>>public chat mod")
			client.PublicChat()
		case 2:
			fmt.Println(">>>private chat mod")
			client.PrivateChat()
		case 3:
			fmt.Println(">>>rename mod:")
			client.UpdateName()
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set the server's ip (127.0.0.1 by default)")
	flag.IntVar(&serverPort, "port", 8888, "set the server's port (8888 by default)")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>connection fail!<<<<")
		return
	}
	fmt.Println(">>>>connection to server success!<<<<")

	go client.DealResponse()
	client.Run()

}

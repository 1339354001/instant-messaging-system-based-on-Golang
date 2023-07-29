package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// hash table to save yser
	OnlineMap map[string]*User
	MapLock   sync.RWMutex
	// channel for message to send
	Message chan string
}

// create a new server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// listen the server's message channel and receive the message
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.MapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.MapLock.Unlock()
	}
}

// send the message to server
func (s *Server) BroadCast(user *User, msg string) {
	if msg == "who" {
		s.MapLock.Lock()
		for _, u := range s.OnlineMap {
			OnlineMsg := "[" + user.Addr + "]" + u.Name + ":" + " online...\n"
			user.conn.Write([]byte(OnlineMsg))
		}
		s.MapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// get new name
		newName := strings.Split(msg, "|")[1]

		// check if the name has existed
		_, ok := s.OnlineMap[newName]
		if ok {
			user.conn.Write([]byte("this name has existed\n"))
		} else {
			s.MapLock.Lock()
			delete(s.OnlineMap, user.Name)
			s.OnlineMap[newName] = user
			user.Name = newName
			s.MapLock.Unlock()
			user.conn.Write([]byte("your name has update: " + newName + "\n"))
		}
	} else {
		sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
		s.Message <- sendMsg
	}

}

// what to do after connection built
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("connection success!")

	// creat a new user
	user := NewUser(conn)
	s.MapLock.Lock()
	// save the user to hash table
	s.OnlineMap[user.Name] = user
	s.MapLock.Unlock()

	// broadcast the message that this user is online
	s.BroadCast(user, "is online.")

	// receive message from the user with a goroutine
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "is offline")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err: ", err)
				return
			}
			// without the last "\n"
			msg := string(buf[:n-1])
			// broadcast this message
			s.BroadCast(user, msg)
		}
	}()

	select {}
}

// run the server
func (s *Server) Start() {
	// build and listen the connection
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen err: ", err)
		return
	}

	defer listener.Close()
	// listen the message channel with a goroutine
	go s.ListenMessage()

	for {
		// check if new user has been online
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept() err: ", err)
			continue
		}
		// run handler with a goroutine
		go s.Handler(conn)
	}
}

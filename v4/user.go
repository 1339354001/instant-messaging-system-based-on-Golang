package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// create a new user
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// check if the user has received a new message with a goroutine
	go user.ListenMessage()

	return user
}

// check if the user has received a new message
func (s *User) ListenMessage() {
	for {
		msg := <-s.C
		// write the message in terminal
		s.conn.Write([]byte(msg + "\n"))
	}
}

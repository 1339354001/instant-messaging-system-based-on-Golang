package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// create a new server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// what to do after connection built
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("connection success!")
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

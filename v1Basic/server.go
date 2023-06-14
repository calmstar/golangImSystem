package main

import (
	"fmt"
	"net"
)

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

type Server struct {
	Ip   string
	Port int
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	defer listener.Close()
	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 处理每个链接
		go this.Handler(conn)
	}
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接成功")
}

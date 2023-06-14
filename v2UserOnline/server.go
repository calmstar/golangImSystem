package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	Lock          *sync.RWMutex
	OnlineMap     map[string]*User // name->User
	broadcastChan chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:            ip,
		Port:          port,
		Lock:          new(sync.RWMutex),
		OnlineMap:     make(map[string]*User),
		broadcastChan: make(chan string),
	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%v", this.Ip, this.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	// server处理监听广播chan
	go this.ListenBroadcastChan()

	// 阻塞监听链接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go this.Handler(conn)
	}
}

func (this *Server) Handler(conn net.Conn) {
	name := conn.RemoteAddr().String()
	user := NewUser(conn)
	fmt.Println(name, ", 链接成功")

	this.Lock.Lock()
	this.OnlineMap[name] = user
	this.Lock.Unlock()

	//发送消息
	this.SendBroadcastChan(name + "上线嘞。。")

	// 阻塞保持对客户端的链接
	select {}
}

func (this *Server) SendBroadcastChan(msg string) {
	this.broadcastChan <- msg
}

func (this *Server) ListenBroadcastChan() {
	for {
		msg := <-this.broadcastChan
		this.Lock.RLock()
		for _, user := range this.OnlineMap {
			user.MsgChan <- msg
		}
		this.Lock.RUnlock()
	}
}

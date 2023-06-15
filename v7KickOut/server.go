package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

// Handler 处理来自客户端的链接
func (this *Server) Handler(conn net.Conn) {
	name := conn.RemoteAddr().String()
	user := NewUser(conn, this)
	fmt.Println(name, ", 链接成功")

	//发送消息
	user.Online()

	// 接受客户端发送的消息 this.conn.read(读取客户端消息)
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 { // 断线触发
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}

			msg := string(buf[:n-2])
			user.DoMsg(msg)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	timeOutChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if 5.0 < time.Now().Sub(user.LastAliveTime).Seconds() {
					fmt.Printf("%v, 被踢下线了\n", user.Name)
					user.SendMsg("你被踢下线了")
					timeOutChan <- struct{}{}
					return
				}
			}
		}
	}()

	// 阻塞保持对客户端的链接
	select {
	case <-timeOutChan:
		conn.Close()
	}
}

func (this *Server) SendBroadcastChan(user *User, msg string) {
	msg = user.Name + ": " + msg
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

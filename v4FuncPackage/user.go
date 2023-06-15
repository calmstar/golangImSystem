package main

import "net"

type User struct {
	Name    string
	Addr    string
	MsgChan chan string
	Conn    net.Conn
	Server  *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	u := &User{
		Name:    conn.RemoteAddr().String(),
		Addr:    conn.RemoteAddr().String(),
		MsgChan: make(chan string),
		Conn:    conn,
		Server:  server,
	}
	go u.ListenMsgChan()
	return u
}

func (this *User) ListenMsgChan() {
	for {
		msg := <-this.MsgChan
		this.Conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	this.Server.Lock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.Lock.Unlock()

	this.Server.SendBroadcastChan(this, "上线")
}

func (this *User) Offline() {
	this.Server.Lock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.Lock.Unlock()

	this.Server.SendBroadcastChan(this, "下线")
}

func (this *User) DoMsg(msg string) {
	this.Server.SendBroadcastChan(this, msg)
}

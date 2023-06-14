package main

import "net"

type User struct {
	Name    string
	Addr    string
	MsgChan chan string
	Conn    net.Conn
}

func NewUser(conn net.Conn) *User {
	u := &User{
		Name:    conn.RemoteAddr().String(),
		Addr:    conn.RemoteAddr().String(),
		MsgChan: make(chan string),
		Conn:    conn,
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

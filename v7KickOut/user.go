package main

import (
	"fmt"
	"net"
	"time"
)

type User struct {
	Name          string
	Addr          string
	MsgChan       chan string
	Conn          net.Conn
	Server        *Server
	LastAliveTime time.Time
}

func NewUser(conn net.Conn, server *Server) *User {
	u := &User{
		Name:          conn.RemoteAddr().String(),
		Addr:          conn.RemoteAddr().String(),
		MsgChan:       make(chan string),
		Conn:          conn,
		Server:        server,
		LastAliveTime: time.Now(),
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
	if msg == "who" {
		this.Who()
	} else if len(msg) > 7 && msg[:7] == "rename|" { // "rename|zhangsan" 修改名字
		this.Rename(msg[7:])
		this.SendMsg("rename success \n")
	} else {
		this.Server.SendBroadcastChan(this, msg)
	}
	this.UpdateAliveTime()
}

func (this *User) SendMsg(msg string) {
	this.Conn.Write([]byte(msg + "\n"))
}

func (this *User) Who() {
	msg := "\n"
	i := 1
	this.Server.Lock.RLock()
	defer this.Server.Lock.RUnlock()
	for _, v := range this.Server.OnlineMap {
		msg += fmt.Sprintf("%v. %v \n", i, v.Name)
		i++
	}
	this.SendMsg(msg)
}

func (this *User) Rename(newName string) {
	this.Server.Lock.Lock()
	defer this.Server.Lock.Unlock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.OnlineMap[newName] = this
	this.Name = newName
}

func (this *User) UpdateAliveTime() {
	this.LastAliveTime = time.Now()
}

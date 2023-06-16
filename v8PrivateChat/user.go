package main

import (
	"fmt"
	"net"
	"strings"
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
	} else if msg == "me" {
		this.SendMsg(this.Name)
	} else if len(msg) > 7 && msg[:7] == "rename|" { // "rename|zhangsan" 修改名字
		this.Rename(msg[7:])
		this.SendMsg("rename success \n")
	} else if len(msg) > 6 && msg[:3] == "to|" { // 私聊模式 "to|zhangsan|我是xx"
		msgArr := strings.Split(msg, "|")
		if len(msgArr) != 3 {
			this.SendMsg("msg err")
			return
		}
		name := msgArr[1]
		content := msgArr[2]
		this.Server.Lock.Lock()
		toUser, ok := this.Server.OnlineMap[name]
		this.Server.Lock.Unlock()

		if !ok {
			this.SendMsg("user not exist")
			return
		}
		if len(content) == 0 {
			this.SendMsg("content empty")
			return
		}
		toUser.SendMsg(this.Name + "，对您说：" + content)

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

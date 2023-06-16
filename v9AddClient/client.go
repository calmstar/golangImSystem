package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	Ip   string
	Port int
	Name string
	Conn net.Conn
	Flag int // 当前client的模式
}

func NewClient(ip string, port int) *Client {
	c := &Client{
		Ip:   ip,
		Port: port,
		Flag: 999, // 默认模式
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", ip, port))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	c.Conn = conn

	return c
}

func (this *Client) DealResponse() {
	io.Copy(os.Stdout, this.Conn)
}

func (this *Client) Menu() bool {
	fmt.Println("请选择你要使用的功能")
	fmt.Println("0 退出")
	fmt.Println("1 私聊")
	fmt.Println("2 公聊")
	fmt.Println("3 改名字")

	var choice int
	fmt.Scanln(&choice)
	if 0 <= choice && choice <= 3 {
		this.Flag = choice
		return true
	} else {
		fmt.Println("输入错误，请重试")
		return false
	}

	return false
}

func (this *Client) Run() {
	for this.Flag != 0 {
		for !this.Menu() {
		}

		switch this.Flag {
		case 1:
			this.PrivateChat()
		case 2:
			this.PublicChat()
		case 3:
			this.Rename()
		case 0:
			fmt.Println("退出")
		default:
			fmt.Println("错误选择")
		}
	}
}

func (this *Client) PublicChat() {
	var msg string
	for msg != "exit" {
		fmt.Println("请输入你要公聊的内容，exit退出")
		fmt.Scanln(&msg)
		if len(msg) > 0 {
			this.Send(msg)
		}
	}
}

func (this *Client) PrivateChat() {
	this.SelectUser()
	var remoteUser, msg string
	fmt.Println("请输入你要私聊的人的名字, exit退出")
	fmt.Scanln(&remoteUser)

	for remoteUser != "exit" {
		fmt.Println("请输入你要私聊的内容")
		fmt.Scanln(&msg)

		if len(msg) != 0 {
			content := fmt.Sprintf("to|%v|%v", remoteUser, msg)
			this.Send(content)
		}

		fmt.Println("请输入你要私聊的人的名字, exit退出")
		fmt.Scanln(&remoteUser)
	}
}

func (this *Client) Send(msg string) {
	_, err := this.Conn.Write([]byte(msg + "\r\n"))
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("发送成功")
	}
}

func (this *Client) SelectUser() {
	msg := "who\r\n"
	_, err := this.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (this *Client) Rename() {
	var newName string
	fmt.Println("请输入你要修改的名字, exit退出")
	fmt.Scanln(&newName)
	for newName != "exit" {
		if len(newName) != 0 {
			msg := fmt.Sprintf("rename|%v", newName)
			this.Send(msg)
		}
		fmt.Println("请输入你要修改的名字, exit退出")
		fmt.Scanln(&newName)
	}

}

var (
	serverIp string
	port     int
)

func init() {
	flag.StringVar(&serverIp, "server_ip", "127.0.0.1", "服务端ip")
	flag.IntVar(&port, "port", 9999, "服务端端口")
}

func main() {
	flag.Parse()
	c := NewClient(serverIp, port)
	if c == nil {
		return
	}

	go c.DealResponse()

	// 开始运行
	c.Run()
}

package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) (*Client, error) {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net dial error: ", err)
		return nil, err
	}

	client.conn = conn

	return client, nil
}

func (c *Client) menu() bool {
	fmt.Println("菜单选择")
	fmt.Println("1. 群聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	var flag int
	fmt.Scanln(&flag)
	if flag >= 0 && flag < 4 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的命令")
		return false
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}
		switch c.flag {
		case 1:
			fmt.Println("你选择群聊模式")
		case 2:
			fmt.Println("你选择的是私聊模式")
		case 3:
			fmt.Println("你选择更新用户名模式")
		}
	}
	fmt.Println("程序即将退出")
	c.conn.Close()
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器地址,默认127.0.0.1")
	flag.IntVar(&serverPort, "port", 12345, "设置服务器端口,默认12345")
}

func main() {
	// 命令行解析
	flag.Parse()
	c, err := NewClient(serverIp, serverPort)
	if err != nil {
		fmt.Print("连接服务器失败,请检查连接参数")
		return
	}
	fmt.Println(c.ServerIp + " 连接成功\n")
	c.Run()
}

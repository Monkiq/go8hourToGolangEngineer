package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

func (c *Client) UpdateName() {
	fmt.Println("请输入你新的用户名")
	fmt.Scanln(&c.Name)
	c.conn.Write([]byte("rename|" + c.Name + "\n"))
}

func (c *Client) PublicChat() {
	fmt.Println("请输入消息")
	msg := ""
	fmt.Scanln(&msg)

	for msg != "exit" {
		if len(msg) != 0 {
			_, err := c.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		msg = ""
		fmt.Scanln(&msg)
	}

}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}
		go func() {
			// 一旦conn有数据,就拷贝到标准输出上,会一直阻塞
			io.Copy(os.Stdout, c.conn)
		}()
		switch c.flag {
		case 1:
			c.PublicChat()
		case 2:
			fmt.Println("你选择的是私聊模式")
		case 3:
			c.UpdateName()
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

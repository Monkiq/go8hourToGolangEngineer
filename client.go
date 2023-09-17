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
}

func NewClient(serverIp string, serverPort int) (*Client, error) {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net dial error: ", err)
		return nil, err
	}

	client.conn = conn

	return client, nil
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
	} else {
		fmt.Println(c.ServerIp + " 连接成功\n")
	}
}

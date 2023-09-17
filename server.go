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

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel的goroutine,一旦有消息就发送给全部的在线user
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message

		// 将msg发给全部的在线user
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 广播消息方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	//  当前连接的业务
	//fmt.Println("连接建立成功")

	user := NewUser(conn, s)
	user.Online()

	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户消息，去除'\n'
			msg := string(buf[:n-1])

			//将得到的消息进行广播
			user.DoMessage(msg)

			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 如果用户活跃就什么也不做
			// 不做任何事情是为了重置下面的定时器
		case <-time.After(time.Second * 35):
			// 已经超时 ,将当前user关闭
			user.SendMsg("你被踢了")
			s.mapLock.Lock()
			delete(s.OnlineMap, user.Name)
			s.mapLock.Unlock()
			// 释放资源
			close(user.C)
			// 关闭连接
			conn.Close()
			// 返回Handler
			return
		}
	}
}

// 启动服务器的接口
func (s *Server) Start() {
	// socker listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
	}
	defer listener.Close()

	//启动监听Message的goroutine
	go s.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err: ", err)
			continue
		}

		// do handler
		go s.Handler(conn)
	}

	// close listensocket
}

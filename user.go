package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 创建一个用户的API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	// 启动监听当前user channel消息的goroutine
	go user.ListenMessager()
	return user
}

// 监听当前User channel的方法，一旦有消息，就发送给对端客户端
func (u *User) ListenMessager() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}

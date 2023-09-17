package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// 启动监听当前user channel消息的goroutine
	go user.ListenMessager()
	return user
}

// 用户上线业务
func (u *User) Online() {
	// 用户上线,将用户加入到onlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	// 广播当前用户上线消息
	u.server.BroadCast(u, "已上线")
}

// 用户下线业务
func (u *User) OffLine() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	u.server.BroadCast(u, "下线了")
}

// 给当前User对应的客户端发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户都有哪些
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式 rename|张三
		uName := msg[7:]
		// 判断名字是否存在
		if _, ok := u.server.OnlineMap[uName]; uName == u.Name || !ok {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.Name = uName
			u.server.OnlineMap[u.Name] = u
			u.server.mapLock.Unlock()
			u.SendMsg("你已经更新用户名为" + u.Name + "\n")
		} else {
			u.SendMsg("该用户名已经存在")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式: to|张三|消息内容
		// 1 获取对方的用户名
		name := strings.Split(msg, "|")[1]
		if name == "" {
			u.SendMsg("消息格式不正确,请使用\"to|张三|消息内容\"")
		}
		// 2 根据用户名获取对方User对象
		if t, ok := u.server.OnlineMap[name]; ok {
			// 3 获取消息内容,通过对方的User对象将消息内容发送过去
			c := strings.Split(msg, "|")[2]
			t.SendMsg(u.Name + "对您说: " + c + "\n")
		} else {
			u.SendMsg("你私聊的用户不存在\n")
		}
	} else {
		u.server.BroadCast(u, msg)
	}

}

// 监听当前User channel的方法，一旦有消息，就发送给对端客户端
func (u *User) ListenMessager() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}

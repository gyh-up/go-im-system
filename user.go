package main

import "net"

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}

	// 启动监听当前User channel消息的goroutine
	go user.ListenMessage()

	return user
}

func (this *User) Online() {
	// 用户上线，将用户加入到onlinemap 中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户上线消息
	this.server.BroadCast(this, "i am online")
}

func (this *User) Offline() {
	// 用户下线，将用户从onlinemap 中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户下线消息
	this.server.BroadCast(this, "i am offline")
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else {
		this.server.BroadCast(this, msg)
	}
}

func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

// 监听当前User channel 的方法，一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for  {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\r\n"))
	}
}
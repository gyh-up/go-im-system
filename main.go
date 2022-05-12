package main

func main () {
	// 实例化一个server对象
	server := NewServer("127.0.0.1", 8888)
	// 启动
	server.Start()
}
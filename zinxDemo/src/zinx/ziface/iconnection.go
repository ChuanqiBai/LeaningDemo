package ziface

import "net"

//链接的抽象层
type IConnection interface {
	//启动链接 让当前链接准备开始工作
	Start()
	//停止链接 结束当前链接的工作
	Stop()
	//获取当前链接的绑定socket conn
	GetTCPConnetcion() *net.TCPConn
	//获取当前链接模块的链接ID
	GetConnID() uint32
	//获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr
	//发送数据，将数据发送给远程客户端
	SendMsg(uint32, []byte) error
	//设置链接属性
	SetProperty(string, interface{})
	//获取链接属性
	GetProperty(string) (interface{}, error)
	//移除链接属性
	RemoveProperty(string)
}

//定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
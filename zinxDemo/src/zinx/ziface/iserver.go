package ziface

//定义服务器接口

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//路由功能，给当前的服务组成一个路由方法，供客户端连接处理使用
	AddRouter(uint32, IRouter)
	//获取当前的ConnManager
	GetConnManager() IConnManager
	//注册钩子函数OnConnStart
	SetOnConnStart(func(IConnection))
	//注册钩子函数OnConnStop
	SetOnConnStop(func(IConnection))
	//调用注册的钩子函数OnConnStart
	CallOnConnStart(IConnection)
	//调用注册的钩子函数OnConnStop
	CallOnConnStop(IConnection)
}

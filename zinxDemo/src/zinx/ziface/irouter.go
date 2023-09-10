package ziface

//路由的抽象接口,路由里的数据都是IRequest

type IRouter interface {
	//处理业务之前的钩子方法Hook
	PreHandle(IRequest)
	//处理业务的主方法
	Handle(IRequest)
	//处理业务之后的方法Hook
	PostHandle(IRequest)
}

package znet

import "github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"

//实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

//这里之所以所有的BaseRouter的方法都为空，是因为有的router不希望有PreHandle、PostHandle

//处理业务之前的钩子方法Hook
func (router *BaseRouter) PreHandle(req ziface.IRequest) {}

//处理业务的主方法
func (router *BaseRouter) Handle(req ziface.IRequest) {}

//处理业务之后的方法Hook
func (router *BaseRouter) PostHandle(req ziface.IRequest) {}

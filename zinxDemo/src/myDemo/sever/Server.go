package main

import (
	"fmt"

	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/znet"
)

//基于Zinx框架开发的 服务器应用程序

// ping test自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// test handle
func (this *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call router handle...")
	//先读取客户端的数据，再会写ping...ping...
	fmt.Println("recv from client msgID: ", req.GetMsgID(), " data= ", string(req.GetData()))

	err := req.GetConnection().SendMsg(200, []byte("ping...ping..."))
	if err != nil {
		fmt.Println("send rsp err ", err)
	}
}

// Hello zinx自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

// test handle
func (this *HelloRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call router handle...")
	//先读取客户端的数据，再会写ping...ping...
	fmt.Println("recv from client msgID: ", req.GetMsgID(), " data= ", string(req.GetData()))

	err := req.GetConnection().SendMsg(201, []byte("hello, welcome to zinx"))
	if err != nil {
		fmt.Println("send rsp err ", err)
	}
}

// 创建链接之后执行的钩子函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("====> DoconnBegin is called....")
	if err := conn.SendMsg(202, []byte("DoConnBegin")); err != nil {
		fmt.Println(err)
	}

	//给当前的链接设置一些属性
	fmt.Println("[Set conn property]")
	conn.SetProperty("Name", "six six six")
	conn.SetProperty("Home", "github.com/ChuanqiBai")
}

// 关闭链接之后执行的钩子函数
func DoConnEnd(conn ziface.IConnection) {
	fmt.Println("====> DoconnEnd is called....")
	fmt.Println("connID: ", conn.GetConnID(), " offline")

	//获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Home = ", home)
	}
}

func main() {
	//创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.1]")

	//注册链接的hook钩子函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnEnd)

	//给当前zinx框架添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//启动sever
	s.Serve()
}

package znet

import (
	"errors"
	"fmt"
	"net"

	"github.com/ChuanqiBai/zinxDemo/src/utils"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
)

// IServer的接口实现
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前server的消息管理模块，用来绑定MsgID与对应的处理handler
	MsgHandler ziface.IMsgHandle
	//该server的链接管理器
	ConnManager ziface.IConnManager
	//创建链接之后的钩子函数OnConnStart
	OnConnStart func(conn ziface.IConnection)
	//销毁链接之后的钩子函数OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

// 定义当前客户端连接所绑定handle api(目前是写死的)
func handler(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] callback to client...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back to client failed, err", err)
		return errors.New("CallbackToClient error")
	}
	return nil
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s ,listener Ip:%s, port:%d is starting\n", utils.Globalobject.Name, utils.Globalobject.Host, utils.Globalobject.TcpPort)
	fmt.Printf("[Zinx] Version %s, Maxconnect: %d, MaxPackageSize:%d\n", utils.Globalobject.Version, utils.Globalobject.MaxConn, utils.Globalobject.MaxPackageSize)
	fmt.Printf("[Start] Server Listenner at IP: %s, Port: %d\n", s.IP, s.Port)

	go func() {
		//开启消息队列并初始化workerPool
		s.MsgHandler.StartWorkerPool()

		//获取一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error: ", err)
			return
		}
		//监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Create Listener failed with error: ", err)
		}

		fmt.Println("Start Zinx server success, Listening")
		var connID uint32
		connID = 0

		//阻塞地等待客户端连接，处理客户端连接业务
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大链接个数的判断，如果当前已经达到上限，直接关闭连接
			if s.ConnManager.Len() >= utils.Globalobject.MaxConn {
				//TODO给客户端响应一个超出最大链接的错误包
				fmt.Println("=====>conn num reach max limit = ", utils.Globalobject.MaxConn)
				conn.Close()
				continue
			}

			connID++
			//将处理链接的业务方法与conn绑定
			dealConn := NewConnection(s, conn, connID, s.MsgHandler)

			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	//将一些服务器的资源进行回收
	fmt.Println("[zinx sever stop]")
	s.ConnManager.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外功能

	//阻塞状态
	select {}
}

func (s *Server) AddRouter(id uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(id, router)
	fmt.Println("Add router succ!")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

// 注册钩子函数OnConnStart
func (s *Server) SetOnConnStart(f func(ziface.IConnection)) {
	s.OnConnStart = f
}

// 注册钩子函数OnConnStop
func (s *Server) SetOnConnStop(f func(ziface.IConnection)) {
	s.OnConnStop = f
}

// 调用注册的钩子函数OnConnStart
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("====> call onConnStart")
		s.OnConnStart(conn)
	}
}

// 调用注册的钩子函数OnConnStop
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("====> call onConnStart")
		s.OnConnStop(conn)
	}
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:        utils.Globalobject.Name,
		IPVersion:   "tcp4",
		IP:          utils.Globalobject.Host,
		Port:        utils.Globalobject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
}

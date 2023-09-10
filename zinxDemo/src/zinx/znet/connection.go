package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/ChuanqiBai/zinxDemo/src/utils"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
)

// 链接模块
type Connection struct {
	//当前conn创建时，隶属于的server
	TcpServer ziface.IServer

	//当前链接的socketTcp 套接字
	Conn *net.TCPConn
	//链接ID
	ConnID uint32

	//当前链接的状态
	IsClosed bool

	//告知当前链接退出的channel
	StopChan chan bool

	//消息管理的模块，管理MsgID和对应的处理模块
	MsgHandler ziface.IMsgHandle

	//无缓冲的管道，用于读写gorotinue之间的消息通信
	MsgChan chan []byte

	//链接属性集合
	Property map[string]interface{}

	//保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgH ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		StopChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
		MsgHandler: msgH,
		Property:   make(map[string]interface{}),
	}

	c.TcpServer.GetConnManager().Add(c)
	//将conn加入到ConnManager中
	return c
}

// 链接的读业务方法
func (conn *Connection) StartReader() {
	fmt.Println("[Reader gorotinue is running]")
	defer fmt.Println("connId = ", conn.ConnID, "[reader is exit], remote addr is", conn.RemoteAddr().String())
	defer conn.Stop()

	for {
		//读取客户端的数据到buf中
		// buf := make([]byte, utils.Globalobject.MaxPackageSize)
		// _, err := conn.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err", err)
		// 	continue
		// }

		//创建一个拆包 解包的对象
		dp := NewDataPack()
		//读取客户端的Msg Head
		headData := make([]byte, dp.GetHeadLen())
		n, err := io.ReadFull(conn.GetTCPConnetcion(), headData)
		if n < int(dp.GetHeadLen()) || err != nil {
			fmt.Println("read head error", err)
			break
		}

		//拆包 得到MgsID和MsgDataLen放入Msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err", err)
			break
		}

		//根据dataLen，再次读取Data数据，放入Msg。Data
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if n, err := io.ReadFull(conn.GetTCPConnetcion(), data); n < int(msg.GetMsgLen()) || err != nil {
				fmt.Println("read msg data failed")
				break
			}
			msg.SetMsgData(data)
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: conn,
			msg:  msg,
		}

		if utils.Globalobject.WorkerPoolSize > 0 {
			conn.MsgHandler.SendRequestToTaskQueue(&req)
		} else {
			go func(req ziface.IRequest) {
				conn.MsgHandler.DoMsgHandler(req)
			}(&req)
		}

	}
}

// 写消息的gorotinue，专门发送给客户端消息的模块
func (conn *Connection) StartWriter() {
	fmt.Println("[Writer gorotinue is running]")
	defer fmt.Println("[ conn writer exit!]", conn.RemoteAddr().String())

	// 不断的阻塞等待channel消息，写给客户端
	for {
		select {
		case data := <-conn.MsgChan:
			//有数据要写给客户端
			if _, err := conn.Conn.Write(data); err != nil {
				fmt.Println("Send data err ", err)
				return
			}
		case <-conn.StopChan:
			//代表reader已经退出，此时Writer也要退出
			return
		}
	}

}

// 提供一个sendMsg方法，将发送给客户端的数据先进行封包再进学校发送
func (conn *Connection) SendMsg(msgID uint32, data []byte) error {
	if conn.IsClosed == true {
		return errors.New("Conn is closed")
	}

	//将data进行封包 msgDataLen/MsgID/MsgData
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPack(msgID, data))
	if err != nil {
		fmt.Println("pack msg failed with msgID: ", msgID)
		return errors.New("Pack msg failed")
	}

	//将数据发送到客户端
	// if _, err := conn.Conn.Write(binaryMsg); err != nil {
	// 	fmt.Println("write msg id= ", msgID, " error ", err)
	// 	return errors.New("conn write failed")
	// }
	conn.MsgChan <- binaryMsg

	return nil
}

// 启动链接 让当前链接准备开始工作
func (conn *Connection) Start() {
	fmt.Println("Conn start... connID = ", conn.ConnID)

	//启动从当前链接读数据的业务
	go conn.StartReader()
	//启动从当前链接写数据的业务
	go conn.StartWriter()

	//按照开发者传递进来的，创建链接之后的hook函数，调用hook
	conn.TcpServer.CallOnConnStart(conn)
}

// 停止链接 结束当前链接的工作
func (conn *Connection) Stop() {
	fmt.Println("Conn stop... connID = ", conn.ConnID)

	if conn.IsClosed == true {
		return
	}

	conn.IsClosed = true

	//调用开发者注册的，消耗链接之前的hook函数
	conn.TcpServer.CallOnConnStop(conn)

	conn.Conn.Close()

	conn.StopChan <- true
	//告知writer关闭

	//将当前链接从ConnMgr中移除
	conn.TcpServer.GetConnManager().Remove(conn)

	//回收资源
	close(conn.StopChan)
	close(conn.MsgChan)
}

// 获取当前链接的绑定socket conn
func (conn *Connection) GetTCPConnetcion() *net.TCPConn {
	return conn.Conn
}

// 获取当前链接模块的链接ID
func (conn *Connection) GetConnID() uint32 {
	return conn.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (conn *Connection) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

// 设置链接属性
func (conn *Connection) SetProperty(key string, val interface{}) {
	conn.propertyLock.Lock()
	defer conn.propertyLock.Unlock()
	//添加一个链接属性
	conn.Property[key] = val
}

// 获取链接属性
func (conn *Connection) GetProperty(key string) (interface{}, error) {
	conn.propertyLock.RLock()
	defer conn.propertyLock.RUnlock()

	//读取属性
	if val, ok := conn.Property[key]; !ok {
		return nil, errors.New("property not found")
	} else {
		return val, nil
	}
}

// 移除链接属性
func (conn *Connection) RemoveProperty(key string) {
	conn.propertyLock.Lock()
	defer conn.propertyLock.Unlock()
	//删除一个链接属性

	if _, ok := conn.Property[key]; ok {
		delete(conn.Property, key)
	}
}

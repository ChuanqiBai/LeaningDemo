package znet

import "github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"

type Request struct {
	//已经和客户端建立好的链接
	conn ziface.IConnection
	//客户端请求的数据
	msg ziface.IMessage
}

//得到当前链接
func (req *Request) GetConnection() ziface.IConnection {
	return req.conn
}

//得到数据ID
func (req *Request) GetMsgID() uint32 {
	return req.msg.GetMsgId()
}

//得到请求的数据
func (req *Request) GetData() []byte {
	return req.msg.GetData()
}

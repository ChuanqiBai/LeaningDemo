package ziface

//IRequest接口，将客户端请求的链接和请求的数据封装到一个Request中

type IRequest interface {
	//得到当前链接
	GetConnection() IConnection
	//得到请求的数据
	GetData() []byte
	//得到请求的ID
	GetMsgID() uint32
}

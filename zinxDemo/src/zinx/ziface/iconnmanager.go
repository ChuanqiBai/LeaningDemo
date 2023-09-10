package ziface

//链接管理模块抽象层

type IConnManager interface {
	//添加链接
	Add(IConnection)
	//删除链接
	Remove(IConnection)
	//根据connID获取链接
	Get(uint32) (IConnection, error)
	//得到当前链接总数
	Len() int
	//清除并终止所有链接
	ClearConn()
}

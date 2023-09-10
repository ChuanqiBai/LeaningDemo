package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
)

//链接管理模块

type ConnManager struct {
	connectionMap map[uint32]ziface.IConnection //管理链接的集合
	connLock      sync.RWMutex                  //保护链接读写的锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connectionMap: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn加入到ConnManager中
	connMgr.connectionMap[conn.GetConnID()] = conn
	fmt.Println("connection add to connMgr successs: connId = ", conn.GetConnID(), ",cur conn num = ", connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connectionMap, conn.GetConnID())
	fmt.Println("connection remove from connMgr successs: connId = ", conn.GetConnID(), ",cur conn num = ", connMgr.Len())
}

// 根据connID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	conn, ok := connMgr.connectionMap[connID]
	if !ok {
		fmt.Println("cannot find conn with id:", connID)
		return nil, errors.New(fmt.Sprintf("not such conn with id:%d", connID))
	}
	return conn, nil
}

// 得到当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connectionMap)
}

// 清除并终止所有链接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connectionMap {
		//停止链接，并从map中清除
		conn.Stop()
		delete(connMgr.connectionMap, connID)
	}

	fmt.Println("Clear all conn success, cur num = ", connMgr.Len())
}

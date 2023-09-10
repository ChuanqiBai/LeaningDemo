package core

import (
	"fmt"
	"sync"
)

//一个AOI地图中的格子类型

type Grid struct {
	//格子ID
	GID int
	//格子左边边界坐标
	MinX int
	//格子右边边界坐标
	MaxX int
	//格子上边边界坐标
	MinY int
	//格子下边边界坐标
	MaxY int
	//当前格子内玩家或物体成员的ID集合
	playerIDs map[int]bool
	//保护当前集合的锁
	pIDlocks sync.RWMutex
}

// 初始化当前的格子的方法
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 给格子添加一个玩家
func (g *Grid) AddPlayer(playerID int) {
	g.pIDlocks.Lock()
	defer g.pIDlocks.Unlock()

	g.playerIDs[playerID] = true
}

// 从格子中删除一个玩家
func (g *Grid) RemovePlayer(playerID int) {
	g.pIDlocks.Lock()
	defer g.pIDlocks.Unlock()
	if _, ok := g.playerIDs[playerID]; ok {
		delete(g.playerIDs, playerID)
	}
}

// 得到当前格子中所有的玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDlocks.RLock()
	defer g.pIDlocks.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

// 调式使用-打印格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerID:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}

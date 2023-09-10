package core

import "sync"

//定义一些AOI的边界值

const (
	AOI_MIN_X  int = 0
	AOI_MAX_X  int = 150
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 0
	AOI_MAX_Y  int = 200
	AOI_CNTS_Y int = 20
)

//当前世界的总管理模块
type WorldManager struct {
	//AOIManager 当前世界AOI的管理模块
	AoiMgr *AOIManager
	//当前在线玩家的集合
	Players map[int32]*Player
	//保护Player集合的锁
	pLock sync.RWMutex
}

//提供一个对外的世界管理模块的句柄
var WorldMgrObj *WorldManager

//初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		//创建世界AOI地图规划
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}

//添加一个对象
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	//将player添加到AOIManager中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

//删除玩家
func (wm *WorldManager) RemovePlayerByGid(pid int32) {
	wm.pLock.Lock()
	player := wm.Players[pid]
	delete(wm.Players, pid)
	wm.pLock.Unlock()

	wm.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Z)
}

//通过玩家ID查询player玩家
func (wm *WorldManager) GetPlayerByPID(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pid]
}

//获取全部的在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, p := range wm.Players {
		players = append(players, p)
	}
	return players
}

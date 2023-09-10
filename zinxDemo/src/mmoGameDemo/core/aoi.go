package core

import (
	"fmt"
)

// AOI区域管理模块

type AOIManager struct {
	//区域的左边界坐标
	MinX int
	//区域的右边界坐标
	MaxX int
	//X方向格子的数量
	CntsX int
	//区域的上边界坐标
	MinY int
	//区域的下边界坐标
	MaxY int
	//Y方向格子的数量
	CntsY int
	//当前区域中有哪些格子map  key=格子的ID，value=格子对象
	Grids map[int]*Grid
}

// 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		Grids: make(map[int]*Grid),
	}
	//给AOI初始化区域的格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//计算格子ID，根据x，y编号
			//格子编号：id = idy*cntX+idx
			gid := y*cntsX + x
			//初始化gid格子
			aoiMgr.Grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}
	return aoiMgr
}

// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在Y轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印格子信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, CntsX:%d, MinY:%d, MaxY:%d, CntsY:%d\n Grids in AoiManager\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	return s
}

// 根据格子GID得到周边九宫格子的ID集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//判断gID是否在AOIManager中
	if _, ok := m.Grids[gID]; !ok {
		return nil
	}

	//初始化grids返回值切片,将当前gid本身加入返回值
	// grids = make([]*Grid, 4)	//这里是绝对不对的，因为这里就会初始化4个nil的数组，然后后面append的时候，不会覆盖之前的nil
	grids = append(grids, m.Grids[gID])

	//需要通过gID得到当前格子x轴的编号idx=gID%nx
	idx := gID % m.CntsX

	//判断idx左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.Grids[gID-1])
	}

	//判断idx右边是否还有格子
	if idx < m.CntsX-1 {
		grids = append(grids, m.Grids[gID+1])
	}

	//遍历gridX集合中每个格子的gid
	gridX := make([]int, 0) //必须使用0，因为不使用0，其append不会覆盖之前的默认值
	for _, v := range grids {
		gridX = append(gridX, v.GID)
	}

	for _, v := range gridX {
		idy := v / m.CntsY

		//gid上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.Grids[v-m.CntsX])
		}
		//gid下边是否还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.Grids[v+m.CntsX])

		}
	}

	return
}

//通过横纵坐标得到周边九宫格内全部的playerIDs

func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	//得到玩家的GID格子ID
	gID := m.GetGIDByPos(x, y)
	//得到GID周边的GID
	grids := m.GetSurroundGridsByGid(gID)
	//将九宫格的信息里的全部playerID累加到playerID中
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}
	return
}

func (m *AOIManager) GetGIDByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntsX + idx
}

// 添加一个playerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.Grids[gID].AddPlayer(pID)
}

// 移除一个格子中的playerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.Grids[gID].RemovePlayer(pID)
}

// 通过GID获取全部的playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.Grids[gID].GetPlayerIDs()
	return
}

// 通过坐标将player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	// gID = int(math.Max(10, float64(gID)))
	m.AddPidToGrid(pID, gID)
}

// 通过坐标把一个player从一个格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	m.RemovePidFromGrid(pID, gID)
}

package core

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/ChuanqiBai/zinxDemo/src/mmoGameDemo/pb"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
	"google.golang.org/protobuf/proto"
)

type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //当前玩家的链接
	X    float32            //平面的x坐标
	Y    float32            //高度
	Z    float32            //平面y坐标
	V    float32            //旋转0~360度
}

//playerID生成器

var PidGen int32 = 1  //用来生产玩家ID的计数器
var IdLock sync.Mutex //保护PidGen的Mutex

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(120 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}
	return p
}

//提供一个发送给客户端消息的方法
//主要是酱pb的protobuf数据序列化后，再调用zinx的sendMsg方法

func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	//将proto Message结构体序列化转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg error: ", err)
		return
	}

	//将二进制文件 通过zinx框架的sendMsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("Player sendMsg error: ", err)
		return
	}
	return
}

// 告知客户端玩家PID，同步已经生成的玩家ID给客户端
func (p *Player) SyncPID() {
	//组建MsgID:0的proto数据
	data := &pb.SyncPid{
		Pid: p.Pid,
	}
	// 将消息发送给客户端
	p.SendMsg(1, data)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	//组建MsgID:200的proto数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //TP=2代表广播的位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, msg)
}

// 广播消息
func (p *Player) Talk(content string) {
	//组建MsgID 200的数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, //tp=1代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	//得到当前世界所有的玩家
	players := WorldMgrObj.GetAllPlayers()
	//向所有玩家发送数据 MsgID200
	for _, player := range players {
		//player分别给对应的客户端发送消息
		player.SendMsg(200, msg)
	}
}

// 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	//获取当前玩家周围的玩家有哪些
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))

	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pid)))
	}

	//将当前玩家的位置信息通过MsgID:200发送给周围玩家
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //tp=2 广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	for _, player := range players {
		player.SendMsg(200, msg)
	}

	//将周围的全部玩家的位置信息发送给当前的玩家MsgID202客户端
	player_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		//制作一个message Player
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.Z,
			},
		}
		player_msg = append(player_msg, p)
	}

	syncPlayer_msg := &pb.SyncPlayers{
		Ps: player_msg[:],
	}

	p.SendMsg(202, syncPlayer_msg)
}

// 广播当前玩家的位置移动信息
func (p *Player) UpdatePos(x, y, z, v float32) {
	//更新当前玩家player对象的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	//组建广播proto协议MsgID:200 tp-4
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4, //tp4移动之后的坐标位置
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//获取当前玩家的周边玩家AOI九宫格之内的玩家
	players := p.GetSurroundingPlayers()

	//一次性给每个玩家对应的客户端发送当前玩家位置更新的消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

// 获取当前玩家的周边玩家AOI九宫格内的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	//得到当前AOI九宫格内的所有玩家PID
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	//将所有的pid对应的player放入切片
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pid)))
	}
	return players
}
func (p *Player) Offline() {
	//得到当前玩家周边的九宫格内有哪些玩家
	players := p.GetSurroundingPlayers()
	//给周围的玩家广播MsgID:201消息
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}

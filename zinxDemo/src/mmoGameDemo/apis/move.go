package apis

import (
	"fmt"

	"github.com/ChuanqiBai/zinxDemo/src/mmoGameDemo/core"
	"github.com/ChuanqiBai/zinxDemo/src/mmoGameDemo/pb"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/znet"
	"google.golang.org/protobuf/proto"
)

//玩家移动

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(req ziface.IRequest) {
	//解析客户端传进来的协议
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(req.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move: Postion Unmarshal error ", err)
		return
	}
	//得到当前发送位置的是哪个玩家
	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("Get pid property error ", err)
		return
	}

	fmt.Printf("Player pid=%dm move(%f,%f,%f, %f)\n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
	//给其他玩家进行当前玩家的位置信息广播
	player := core.WorldMgrObj.GetPlayerByPID(pid.(int32))
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}

package apis

import (
	"fmt"

	"github.com/ChuanqiBai/zinxDemo/src/mmoGameDemo/core"
	"github.com/ChuanqiBai/zinxDemo/src/mmoGameDemo/pb"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/znet"
	"google.golang.org/protobuf/proto"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(req ziface.IRequest) {
	//解析客户端传进来的proto协议
	msg := &pb.Talk{}
	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("talk Unmarshal err: ", err)
	}
	//当前的聊天数据是属于哪个玩家发送的
	pid, err := req.GetConnection().GetProperty("pid")

	//根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPID(pid.(int32))

	//将这个消息广播给其他全部在线的玩家
	player.Talk(msg.Content)

}

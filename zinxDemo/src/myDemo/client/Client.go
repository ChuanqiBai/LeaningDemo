package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/ChuanqiBai/zinxDemo/src/zinx/znet"
)

func main() {
	fmt.Println("client start....")
	time.Sleep(time.Second)

	//链接服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("client start err, exit", err)
		return
	}

	//向链接写入数据
	for {
		//发送封包的msg消息 mgsOD:0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPack(0, []byte("zinx V0.5 client test msg")))
		if err != nil {
			fmt.Println("Pack err: ", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write err: ", err)
			return
		}

		//服务器回复一个msg数据

		//先读取流中的head，得到ID和dataLen
		head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, head); err != nil {
			fmt.Println("read head err ", err)
			break
		}

		//将二进制的head拆包到msg结构体中
		msg, err := dp.Unpack(head)
		if err != nil {
			fmt.Println("client Unpack failed err:", err)
			break
		}

		//再根据DataLen进行第二次读取，将data读取
		if msg.GetMsgLen() > 0 {
			msg := msg.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data failed with err: ", err)
				break
			}
			fmt.Println("Recv sever msg: ID = ", msg.ID, "msg len: ", msg.DataLen, " data = ", string(msg.Data))
		}

		//阻塞 避免cpu空转
		time.Sleep(time.Second)
	}
}

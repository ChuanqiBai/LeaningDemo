package znet_test

import (
	"io"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/ChuanqiBai/zinxDemo/src/zinx/znet"
)

// 只是负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	//模拟的服务器
	t.Parallel()
	//创建socketTcp
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	require.NoError(t, err)

	//创建一个gorotinue 负责从客户端处理业务
	go func() {
		//从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			require.NoError(t, err)

			go func(conn net.Conn) {
				//处理客户端的请求
				//—————————拆包过程———————————
				dp := NewDataPack()
				for {
					//第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					n, err := io.ReadFull(conn, headData)
					if err == io.EOF {
						log.Println("read file eof")
						break
					}
					require.NoError(t, err)
					require.EqualValues(t, n, dp.GetHeadLen())

					msgHead, err := dp.Unpack(headData)
					require.NoError(t, err)
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要再进行一次读取
						//第二次从conn读，根据head的dataLen再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.DataLen)

						//根据dataLen的长度，再次从IO流中读取
						n, err := io.ReadFull(conn, msg.Data)
						require.NoError(t, err)
						require.EqualValues(t, n, msg.DataLen)

						log.Println("Recv MsgId: ", msg.ID, " dataLen = ", msg.DataLen, " data= ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	//模拟客户端

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	require.NoError(t, err)
	//创建一个封包对象
	dp := NewDataPack()

	//模拟粘包过程，封装2个msg一起发送
	//封装第一个msg包
	msg1 := &Message{
		ID:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	require.NoError(t, err)

	//封装第二个msg包
	msg2 := &Message{
		ID:      2,
		DataLen: 9,
		Data:    []byte{'h', 'e', 'l', 'l', 'w', 'o', 'r', 'l', 'd'},
	}
	sendData2, err := dp.Pack(msg2)
	require.NoError(t, err)

	//将2个包拼接在一起
	sendData1 = append(sendData1, sendData2...)
	n, err := conn.Write(sendData1)
	require.EqualValues(t, n, len(sendData1))
	require.NoError(t, err)

}

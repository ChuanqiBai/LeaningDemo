package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/ChuanqiBai/zinxDemo/src/utils"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
)

//封包 拆包的具体模块

type DataPack struct{}

// 拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	//dataLen uint32(4字节) + Id uint32(4字节)
	return 8
}

// 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//存放一个byte字节流的缓存
	dataBuff := bytes.NewBuffer([]byte{})
	//将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgID写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将Msg data写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// 拆包方法 (将包的head信息读出来) 之后根据head信息里的data长度再进行一次读
func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	//创建一个从输入的二进制数据读取的ioReader
	dataBuff := bytes.NewReader(data)

	//只解压headxinxi买得到datalen和msgID
	msg := &Message{}

	//读取dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读取MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	if utils.Globalobject.MaxPackageSize > 0 && msg.DataLen > utils.Globalobject.MaxPackageSize {
		return nil, errors.New("msg data is too large")
	}

	return msg, nil
}

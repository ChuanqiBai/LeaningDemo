package znet

type Message struct {
	ID      uint32 //消息的ID
	DataLen uint32 //消息的长度
	Data    []byte //消息的数据
}

// 创建一个Msg
func NewMsgPack(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息的ID
func (msg *Message) GetMsgId() uint32 {
	return msg.ID
}

// 获取消息的长度
func (msg *Message) GetMsgLen() uint32 {
	return msg.DataLen
}

// 获取消息的内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

// 设置消息的ID
func (msg *Message) SetMsgId(id uint32) {
	msg.ID = id
}

// 设置消息的内容
func (msg *Message) SetMsgLen(length uint32) {
	msg.DataLen = length
}

// 设置消息的长度
func (msg *Message) SetMsgData(data []byte) {
	msg.Data = data
}

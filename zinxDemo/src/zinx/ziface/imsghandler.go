package ziface

//消息管理抽象层

type IMsgHandle interface {
	//调度/执行对应的Router消息处理方法
	DoMsgHandler(IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)
	//启动Worker工作池
	StartWorkerPool()
	//将消息交给MsgHandler，由MsgHandler来负责分配到worker的TaskQueue中
	SendRequestToTaskQueue(IRequest)
}

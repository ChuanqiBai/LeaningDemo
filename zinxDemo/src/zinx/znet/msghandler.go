package znet

import (
	"fmt"
	"sync"

	"github.com/ChuanqiBai/zinxDemo/src/utils"
	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
)

// 消息处理模块的实现
type MsgHandler struct {
	//存放每个MsgID对应的处理方法
	ApiMgr map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkPoolSize uint32
	//用来保证对worker pool的初始化只会执行一次
	once sync.Once
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		ApiMgr:       make(map[uint32]ziface.IRouter),
		TaskQueue:    make([]chan ziface.IRequest, utils.Globalobject.WorkerPoolSize),
		WorkPoolSize: utils.Globalobject.WorkerPoolSize, //从全局配置中获取
	}
}

// 调度/执行对应的Router消息处理方法
func (msgH *MsgHandler) DoMsgHandler(req ziface.IRequest) {
	//从req中找到msgID
	handler, ok := msgH.ApiMgr[req.GetMsgID()]
	if !ok {
		fmt.Println("api msgID= ", req.GetMsgID(), " is not registed!")
		return
	}
	//根据MsgID调度对应的router业务
	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

// 为消息添加具体的处理逻辑
func (msgH *MsgHandler) AddRouter(id uint32, router ziface.IRouter) {
	//判断当前Msg绑定API处理方法是否存在
	if _, ok := msgH.ApiMgr[id]; ok {
		//id已经注册
		fmt.Println("id 已经注册了")
		return
	}

	//添加msg与API的绑定关系
	msgH.ApiMgr[id] = router
	fmt.Println("Add api MsgID=", id, " success")
}

// 启动一个worker工作池(只能发生一次，一个zinx框架只有一个工作池)
func (msgH *MsgHandler) StartWorkerPool() {
	msgH.once.Do(func() {
		//根据workerPoolSize，分别开启worker，每个worker用一个gorotinue来承载
		for i := 0; i < int(msgH.WorkPoolSize); i++ {
			//启动一个worker

			//当前的worker对应的channel消息队列开辟空间
			msgH.TaskQueue[i] = make(chan ziface.IRequest, utils.Globalobject.MaxWorkerTaskLen)
			//启动当前的worker，阻塞等待消息从channel中传递进来
			go msgH.StartOneWorker(i, msgH.TaskQueue[i])
		}
	})
}

// 启动一个worker工作流程
func (msgH *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started...")

	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		case req := <-taskQueue:
			msgH.DoMsgHandler(req)
		}
	}
}

// 将消息交给MsgHandler，由MsgHandler来负责分配到worker的TaskQueue中
func (msgH *MsgHandler) SendRequestToTaskQueue(req ziface.IRequest) {
	//将消息瓶颈分配给不同的worker
	//根据客户端的ConnID来进行分配,采用瓶颈分配的轮询算法
	workerID := req.GetConnection().GetConnID() % msgH.WorkPoolSize
	fmt.Println("Add ConnID = ", req.GetConnection().GetConnID(),
		" request MsgID = ", req.GetMsgID(),
		" to Worker: workerID = ", workerID)
	//将消息推送到worker对应的taskQueue
	msgH.TaskQueue[workerID] <- req
}

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ChuanqiBai/zinxDemo/src/zinx/ziface"
)

//存储一切关于Zinx框架的全局参数，供其他模块使用，一些参数是可以通过zinx.json由用户进行配置

type Globalobj struct {
	//server
	TcpServer ziface.IServer //当前Zinx全局的server对象
	Host      string         //当前服务器主机监听的Ip
	TcpPort   int            //当前服务器主机监听端口
	Name      string         //当前服务器名称

	//zinx
	Version          string //当前ZINX的版本号
	MaxConn          int    //主机允许的最大连接数
	MaxPackageSize   uint32 //数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作池的gorotinue数量
	MaxWorkerTaskLen uint32 //zinx框架允许用户worker任务队列的最大长度(限定条件)
}

//定义一个全局的对外Globalobj

var Globalobject *Globalobj

//提供一个init方法，初始化当前的Globalobj

func init() {
	Globalobject = &Globalobj{
		Name:             "ZinxSeverApp",
		Version:          "V0.4",
		TcpPort:          8080,
		Host:             "0.0.0.0",
		MaxConn:          1024,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//尝试从json文件加载用户定义的参数
	Globalobject.Reload()
}

func (g *Globalobj) Reload() {
	fileStr, _ := os.Getwd()
	fmt.Println("当前的工作路径", fileStr)
	filePath := "conf/zinx.json"
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("文件不存在:", filePath)
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("read data from json file failed", err)
		panic(err)
	}

	err = json.Unmarshal(data, &Globalobject)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"appconfig"
	"gamelog"
	"logsvr/mainlogic"
	"strconv"
	"tcpserver"
	"utility"
)

func main() {
	//加载配制文件
	appconfig.LoadConfig()

	//初始化日志系统
	gamelog.InitLogger("logsvr", appconfig.LogSvrLogLevel)

	//开启控制台窗口，可以接受一些调试命令
	utility.StartConsole()

	//注册控制台命令处理方法
	RegConsoleCmdHandler()

	//注册所有TCP消息处理方法
	RegTcpMsgHandler()

	//逻辑初始化
	mainlogic.Init()

	//启动TCP服务器
	gamelog.Error("----Log Server Start-----")
	tcpserver.ServerRun(":"+strconv.Itoa(appconfig.LogSvrPort), 5000)
}

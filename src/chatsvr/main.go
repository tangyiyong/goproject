package main

import (
	"appconfig"
	"chatsvr/msgprocess"
	"gamelog"
	"strconv"
	"tcpserver"
	"utility"
)

func main() {
	//加载配制文件
	appconfig.LoadConfig()

	//初始化日志系统
	gamelog.InitLogger("chat")
	gamelog.SetLevel(appconfig.ChatLogLevel)

	//开启控制台窗口，可以接受一些调试命令
	utility.StartConsole()

	//注册控制台命令处理方法
	RegConsoleCmdHandler()

	//注册所有TCP消息处理方法
	RegTcpMsgHandler()

	//消息处理逻辑初始化
	msgprocess.Init()

	//启动TCP服务器
	gamelog.Warn("----Chat Server Start-----")
	tcpserver.ServerRun(":"+strconv.Itoa(appconfig.ChatSvrPort), 5000)
}

package main

import (
	"appconfig"
	"battlesvr/mainlogic"
	"flag"
	"gamelog"
	"strconv"
	"tcpserver"
	"utility"
)

func main() {
	//加载配制文件
	appconfig.LoadConfig()

	//需要先从命令行参数中取当前的端口号
	port := flag.Int("port", appconfig.BattleSvrPort, "Need a listen port")
	flag.Parse()
	appconfig.BattleSvrPort = *port

	//初始化日志系统
	gamelog.InitLogger("battle"+strconv.Itoa(appconfig.BattleSvrPort), appconfig.BattleLogLevel)

	//初始化工具系统
	utility.Init()

	//开启控制台窗口，可以接受一些调试命令
	utility.StartConsole()

	//注册所有HTTP消息处理方法
	RegHttpMsgHandler()

	//注册控制台命令处理方法
	RegConsoleCmdHandler()

	//注册所有TCP消息处理方法
	RegTcpMsgHandler()

	//消息处理逻辑初始化
	mainlogic.Init()

	//注册到游戏服
	mainlogic.RegisterToGameSvr()
	//启动TCP服务器
	gamelog.Error("----Battle Server Start--Port:%d---", appconfig.BattleSvrPort)
	tcpserver.MsgDispatcher = mainlogic.BatSvrMsgDispatcher
	tcpserver.ServerRun(":"+strconv.Itoa(appconfig.BattleSvrPort), 5000)
}

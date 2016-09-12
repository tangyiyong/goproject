package main

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"gamesvr/mainlogic"
	"gamesvr/reggamesvr"
	"mongodb"
	"strconv"
	"utility"
)

func main() {
	//加载配制文件
	appconfig.LoadConfig()

	//初始化日志系统
	gamelog.InitLogger("game")
	gamelog.SetLevel(appconfig.GameLogLevel)

	//设置mongodb的服务器地址
	mongodb.Init(appconfig.GameDbAddr)

	//初始化工具系统
	utility.Init()

	//加载所有游戏配制数据
	gamedata.LoadGameData()

	//开启输入控制台程序
	utility.StartConsole()

	//注册控制台命令处理方法
	RegConsoleCmdHandler()

	//注册所有HTTP消息处理方法
	RegHttpMsgHandler()

	//注册所有的TCP消息处理方法
	RegTcpMsgHandler()

	//初始化主逻辑模块
	mainlogic.Init()

	//注册到账号服务器
	reggamesvr.RegisterToSvr()

	//连接到其它服务器
	mainlogic.ConnectToOtherSvr()

	//启动http监听服务
	gamelog.Error("----Game Server Start-----")
	//http.ListenAndServe(/*appconfig.GameSvr+*/ ":"+strconvt.Itoa(appconfig.GameSvrPort), nil)
	err := utility.HttpLimitListen(":"+strconv.Itoa(appconfig.GameSvrPort), appconfig.GameMaxCon)
	if err != nil {
		gamelog.Error("----Http Listen Error :%s-----", err.Error())
	}
}

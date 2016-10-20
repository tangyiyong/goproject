package main

import (
	"accountsvr/mainlogic"
	"appconfig"
	"gamelog"
	"mongodb"
	"strconv"
	"utility"
)

func main() {
	//加载配制文件
	appconfig.LoadConfig()

	//初始化日志系统
	gamelog.InitLogger("account", appconfig.AccountLogLevel)

	//设置mongodb的服务器地址
	mongodb.Init(appconfig.AccountDbAddr)

	//初始化工具系统
	utility.Init()

	//初始化游戏服务器管理对象
	mainlogic.Init()

	//开启控制台窗口，接受用户输入命令
	utility.StartConsole()

	//注册控制台命令处理方法
	RegConsoleCmdHandler()

	//注册所有HTTP消息处理方法
	RegHttpMsgHandler()

	//启动http监听服务器
	gamelog.Error("----Account Server Start-----")
	utility.HttpLimitListenTimeOut(":"+strconv.Itoa(appconfig.AccountSvrPort), appconfig.AccountMaxCon)

	//http.ListenAndServe( /*appconfig.AccountSvr+*/ ":"+strconv.Itoa(appconfig.AccountSvrPort), nil)
}

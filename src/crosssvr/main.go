package main

import (
	"appconfig"
	"crosssvr/mainlogic"
	"gamelog"
	"net/http"
	"strconv"
	"utility"
)

func main() {
	//加载配制文件
	appconfig.LoadConfig()

	//初始化日志系统
	gamelog.InitLogger("cross", true)
	gamelog.SetLevel(appconfig.CrossLogLevel)

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

	//启动TCP服务器
	gamelog.Warn("----Cross Server Start-----")
	http.ListenAndServe(":"+strconv.Itoa(appconfig.CrossSvrHttpPort), nil)
}

package main

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"sdk/sdklogic"
	"strconv"
	"utility"
)

// 1 开一个http server
// 2 读取gamesvr list配置表，取得各游戏服的地址 —— 能够根据svrId往各游戏服推送数据

// SDK向gamesvr post http

func main() {
	//加载配制文件
	appconfig.LoadConfig()

	//初始化日志系统
	gamelog.InitLogger("sdk", true)
	gamelog.SetLevel(appconfig.SdkLogLevel)

	//设置mongodb的服务器地址
	mongodb.Init(appconfig.GameDbAddr)

	//开启控制台窗口，可以接受一些调试命令
	utility.StartConsole()

	//注册控制台命令处理方法
	utility.HandleFunc("setloglevel", HandCmd_SetLogLevel)

	sdklogic.LoadSvrAddrList()

	//注册所有http消息处理方法
	RegSdkHttpMsgHandler()

	err := utility.HttpLimitListen(":"+strconv.Itoa(appconfig.SdkSvrPort), 0)
	if err != nil {
		gamelog.Error("----Http Listen Error :%s-----", err.Error())
	}
}

func HandCmd_SetLogLevel(args []string) bool {
	level, err := strconv.Atoi(args[1])
	if err != nil {
		gamelog.Error("HandCmd_SetLogLevel Error : Invalid param :%s", args[1])
		return true
	}
	gamelog.SetLevel(level)
	return true
}

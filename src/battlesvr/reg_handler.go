package main

import (
	"battlesvr/mainlogic"
	"utility"
)

func RegHttpMsgHandler() {

}

//注册TCP消息处理方法
func RegTcpMsgHandler() {
}

//注册控制台消息处理方法
func RegConsoleCmdHandler() {
	utility.HandleFunc("setloglevel", mainlogic.HandCmd_SetLogLevel) //例如 setloglevel [1]
	//utility.HandleFunc()
	//utility.HandleFunc()
	//utility.HandleFunc()
}

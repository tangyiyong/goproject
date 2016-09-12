package main

import (
	"logsvr/mainlogic"
	"tcpserver"
	//"utility"
	"msg"
)

//注册TCP消息处理方法
func RegTcpMsgHandler() {
	tcpserver.HandleFunc(msg.MSG_DISCONNECT, mainlogic.Hand_DisConnect)
	tcpserver.HandleFunc(msg.MSG_CHECK_IN_REQ, mainlogic.Hand_CheckInReq)
	tcpserver.HandleFunc(msg.MSG_SVR_LOGDATA, mainlogic.Hand_OnLogData)
}

//注册控制台消息处理方法
func RegConsoleCmdHandler() {

	//utility.HandleFunc()
	//utility.HandleFunc()
	//utility.HandleFunc()
}

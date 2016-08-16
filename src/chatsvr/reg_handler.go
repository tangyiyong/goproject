package main

import (
	"chatsvr/msgprocess"
	"tcpserver"
	//"utility"
	"msg"
)

//注册TCP消息处理方法
func RegTcpMsgHandler() {
	tcpserver.HandleFunc(msg.MSG_CHECK_IN_REQ, msgprocess.Hand_CheckInReq)
	tcpserver.HandleFunc(msg.MSG_CHATMSG_REQ, msgprocess.Hand_ChatMsgReq)
	tcpserver.HandleFunc(msg.MSG_DISCONNECT, msgprocess.Hand_DisConnect)
	tcpserver.HandleFunc(msg.MSG_GAME_TO_CLIENT, msgprocess.Hand_Game_To_Client)
	tcpserver.HandleFunc(msg.MSG_GUILD_NOTIFY, msgprocess.Hand_GuildChange_Notify)
	tcpserver.HandleFunc(msg.MSG_HEART_BEAT, msgprocess.Hand_HeartBeat)

}

//注册控制台消息处理方法
func RegConsoleCmdHandler() {

	//utility.HandleFunc()
	//utility.HandleFunc()
	//utility.HandleFunc()
}

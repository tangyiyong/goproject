package main

import (
	"accountsvr/mainlogic"
	"net/http"
)

//注册http处理消息
func RegHttpMsgHandler() {

	//玩家登录
	http.HandleFunc("/login", mainlogic.Handle_Login)

	//玩家注册
	http.HandleFunc("/register", mainlogic.Handle_Register)

	//游客玩家注册
	http.HandleFunc("/tourist_register", mainlogic.Handle_TouristRegister)
	http.HandleFunc("/bind_tourist", mainlogic.Handle_BindTourist)

	//玩家请求服务器列表
	http.HandleFunc("/serverlist", mainlogic.Handle_ServerList)

	//游戏服务器查询用户是否己登录
	http.HandleFunc("/verifyuserlogin", mainlogic.Handle_VerifyUserLogin)

	//游戏服务器注册
	http.HandleFunc("/reggameserver", mainlogic.Handle_RegisterGameSvr)

	//以下是GM后台的指令
	http.HandleFunc("/set_gamesvr_flag", mainlogic.Handle_SetGamesvrFlag)
	http.HandleFunc("/get_server_list", mainlogic.Handle_GetServerList)
	http.HandleFunc("/gm_login", mainlogic.Handle_GmLogin)
}

//注册TCP处理消息
func RegTcpMsgHandler() {

}

//注册命令行处理消息
func RegConsoleCmdHandler() {

	//utility.HandleFunc(msgprocess.MSG_CHECK_IN_REQ, msgprocess.Hand_CheckInReq)
	//utility.HandleFunc(msgprocess.MSG_CHAT_REQ, msgprocess.Hand_ChatMsgReq)
	//utility.HandleFunc(msgprocess.MSG_DISCONNECT, msgprocess.Hand_DisConnect)
}

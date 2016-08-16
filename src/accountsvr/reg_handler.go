package main

import (
	"accountsvr/gamesvrmgr"
	"accountsvr/login"
	"net/http"
)

//注册http处理消息
func RegHttpMsgHandler() {

	//玩家登录
	http.HandleFunc("/login", login.Handle_Login)

	//玩家注册
	http.HandleFunc("/register", login.Handle_Register)

	//游客玩家注册
	http.HandleFunc("/tourist_register", login.Handle_TouristRegister)
	http.HandleFunc("/bind_tourist", login.Handle_BindTourist)

	//玩家请求服务器列表
	http.HandleFunc("/serverlist", login.Handle_ServerList)

	//游戏服务器查询用户是否己登录
	http.HandleFunc("/verifyuserlogin", login.Handle_VerifyUserLogin)

	//游戏服务器注册
	http.HandleFunc("/reggameserver", gamesvrmgr.Handle_RegisterGameSvr)
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

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
	http.HandleFunc("/get_serverlist", mainlogic.Handle_GetServerList)

	//游戏服务器查询用户是否己登录
	http.HandleFunc("/verifyuserlogin", mainlogic.Handle_VerifyUserLogin)

	//游戏服务器注册
	http.HandleFunc("/reggameserver", mainlogic.Handle_RegisterGameSvr)

	//获取游戏公告
	http.HandleFunc("/get_game_public", mainlogic.Handle_GetGamePublic)

	//游戏服查询激活码状态并领取激活码
	http.HandleFunc("/gamesvr_giftcode", mainlogic.Handle_GameSvrGiftCode)

	//以下是GM后台的指令
	http.HandleFunc("/gm_set_svrstate", mainlogic.Handle_SetGameSvrState)
	http.HandleFunc("/gm_server_list", mainlogic.Handle_GmServerList)
	http.HandleFunc("/gm_login", mainlogic.Handle_GmLogin)
	http.HandleFunc("/gm_enable_account", mainlogic.Handle_GmEnableAccount)
	http.HandleFunc("/gm_add_giftaward", mainlogic.Handle_AddGiftAward)
	http.HandleFunc("/gm_make_giftcode", mainlogic.Handle_MakeGiftCode)
	http.HandleFunc("/get_account_info", mainlogic.Handle_GetPlayerInfo)
	http.HandleFunc("/get_net_list", mainlogic.Handle_GetNetList)
	http.HandleFunc("/add_net_list", mainlogic.Handle_AddNetList)
	http.HandleFunc("/del_net_list", mainlogic.Handle_DelNetList)
	http.HandleFunc("/gm_query_svrip", mainlogic.Handle_QuerySvrIp)
	http.HandleFunc("/gm_get_giftaward", mainlogic.Handle_GetGiftAward)
	http.HandleFunc("/gm_del_giftaward", mainlogic.Handle_DelGiftAward)
	http.HandleFunc("/gm_get_enablelst", mainlogic.Handle_GmGetEnableLst)
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

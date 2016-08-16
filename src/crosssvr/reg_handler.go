package main

import (
	"crosssvr/mainlogic"
	//"tcpserver"
	//"utility"
	//"msg"
	"net/http"
)

func RegHttpMsgHandler() {
	http.HandleFunc("/reggameserver", mainlogic.Handle_RegisterGameSvr)
	http.HandleFunc("/cross_query_score_rank", mainlogic.Hand_GetScoreRank)     //! 玩家请求积分排行榜
	http.HandleFunc("/cross_query_score_target", mainlogic.Hand_GetScoreTarget) //! 玩家请求积分战目标
	http.HandleFunc("/cross_get_fight_target", mainlogic.Handle_GetFightTarget) //! 玩家请求战斗目标数据
}

//注册TCP消息处理方法
func RegTcpMsgHandler() {

}

//注册控制台消息处理方法
func RegConsoleCmdHandler() {

	//utility.HandleFunc()
	//utility.HandleFunc()
	//utility.HandleFunc()
}

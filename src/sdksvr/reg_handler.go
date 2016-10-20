package main

import (
	"net/http"
	"sdksvr/mainlogic"
)

//注册http消息处理方法
func RegSdkHttpMsgHandler() {
	//! From Gamesvr
	http.HandleFunc("/reggameserver", mainlogic.Hand_RegGamesvrAddr)

	//! From 第三方
	http.HandleFunc("/sdk_recharge_info", mainlogic.HandSdk_RechargeSuccess)
}

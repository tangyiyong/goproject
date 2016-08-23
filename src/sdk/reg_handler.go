package main

import (
	"net/http"
	"sdk/sdklogic"
)

//注册http消息处理方法
func RegSdkHttpMsgHandler() {
	//! From Gamesvr
	http.HandleFunc("/create_recharge_order", sdklogic.HandSvr_CreateRechargeOrder)

	//! From 第三方
	http.HandleFunc("/sdk_recharge_info", sdklogic.HandSdk_RechargeSuccess)
}

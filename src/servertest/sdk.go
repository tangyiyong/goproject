package main

import (
	// "bytes"
	// "encoding/json"
	"fmt"
	"msg"
)

func (self *TPlayer) Create_Recharge_Order_2Gamesvr() {
	reqUrl := "http://127.0.0.1:8082/create_recharge_order"
	var req msg.Msg_create_recharge_order_Req
	req.SessionKey = self.SessoinKey
	req.PlayerID = self.PlayerID
	req.OrderID = "abcdefg233"
	req.Channel = "360"
	req.PlatformEnum = 0
	req.ChargeCsvID = 2

	backBuf := PostMsg(reqUrl, &req)
	fmt.Println(backBuf)
}

func (self *TPlayer) Recharge_Syccess_2SDK() {
	reqUrl := "http://127.0.0.1:8110/sdk_recharge_info"
	var req msg.SDKMsg_recharge_result
	req.OrderID = "abcdefg233"
	req.ThirdOrderID = "zzzzzzz"
	req.Channel = "360"
	req.PlayerID = 10023
	req.RMB = 233

	backBuf := PostMsg(reqUrl, &req)
	fmt.Println(backBuf)
}

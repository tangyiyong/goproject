/***********************************************************************
* @ 游戏服通知SDK进程
* @ brief
    1、gamesvr先通知SDK进程，建立新充值订单

    2、第三方充值信息到达后，验证是否为有效订单

* @ author zhoumf
* @ date 2016-8-18
***********************************************************************/
package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

func Handle_Recharge_Notify(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.Msg_Recharge_Notify
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Handle_Recharge_Notify unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	defer func() {
		w.Write([]byte("ok"))
	}()

	// 充值到账，增加钻石数量
	var player *TPlayer = GetPlayerByID(req.PlayerID)
	if player == nil {
		gamelog.Error("Handle_Recharge_Notify Error: Invalid PlayerID:%d, chargeid:%d, chargemoney:%d", req.PlayerID, req.ChargeID, req.RMB)
		return
	}
	player.HandChargeRenMinBi(req.RMB, req.ChargeID)

	return
}

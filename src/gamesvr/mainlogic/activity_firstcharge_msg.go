package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家查询首冲活动信息
func Hand_QueryFirstChargeData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryActivity_FirstRecharge_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryActivityFirstRechargeInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryActivity_FirstRecharge_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.ActivityModule.CheckReset()

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.FirstCharge.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		gamelog.Error("Hand_QueryActivityFirstRechargeInfo Error: Not open")
		return
	}

	response.RetCode = msg.RE_SUCCESS

	//! 首充奖励状态
	response.FirstRechargeStatus = player.ActivityModule.FirstCharge.FirstAward
	response.NextRechargeStatus = player.ActivityModule.FirstCharge.NextAward
}

//! 请求领取首充奖励
func Hand_GetFirstRechargeAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetActivity_FirstRecharge_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFirstRechargeAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetActivity_FirstRecharge_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.ActivityModule.CheckReset()

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.FirstCharge.ActivityID) == false {
		gamelog.Error("Hand_GetFirstRechargeAward Error: Activity not open.")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 判断当前是否能够领取首充奖励
	if (player.ActivityModule.FirstCharge.FirstAward != 1 && req.GetAwardType == 1) ||
		(player.ActivityModule.FirstCharge.NextAward != 1 && req.GetAwardType == 2) {
		gamelog.Error("Hand_GetFirstRechargeAward Error: Can't recvice first recharge award")
		response.RetCode = msg.RE_NOT_RECHARGE
		return
	}

	//! 发放奖励
	awardLst := []gamedata.ST_ItemData{}

	if req.GetAwardType == 1 {
		awardLst = gamedata.GetItemsFromAwardID(gamedata.FirstRechargeAwardID)
	} else {
		awardLst = gamedata.GetItemsFromAwardID(gamedata.NextRechargeAwardID)
	}

	player.BagMoudle.AddAwardItems(awardLst)

	for _, v := range awardLst {
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum})
	}

	//! 修改标记
	if req.GetAwardType == 1 {
		player.ActivityModule.FirstCharge.FirstAward = 2
	} else if req.GetAwardType == 2 {
		player.ActivityModule.FirstCharge.NextAward = 2
	}

	player.ActivityModule.FirstCharge.DB_SetFirstRechargeMark()

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

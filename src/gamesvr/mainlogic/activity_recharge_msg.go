package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 查询累计充值活动信息
func Hand_QueryActivityTotalRechargeInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryActivity_TotalRecharge_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryActivityTotalRechargeInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryActivity_TotalRecharge_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.ActivityModule.CheckReset()

	var totalRechargeInfo *TActivityRecharge
	for i, v := range player.ActivityModule.Recharge {
		if v.ActivityID == req.ActivityID {
			totalRechargeInfo = &player.ActivityModule.Recharge[i]
			break
		}
	}

	if totalRechargeInfo == nil {
		gamelog.Error("Hand_GetRechargeAward Error: Activity not exist %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	response.RetCode = msg.RE_SUCCESS

	//! 充值回馈
	response.RechargeNum = totalRechargeInfo.RechargeValue
	response.AwardMark = int(totalRechargeInfo.AwardMark)
	response.ActivityID = req.ActivityID
	response.AwardType = G_GlobalVariables.GetActivityAwardType(req.ActivityID)
}

//! 玩家请求领取充值反馈
func Hand_GetRechargeAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetRechargeAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetRechargeAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetRechargeAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测当前是否有此活动
	var totalRechargeInfo *TActivityRecharge
	var activityIndex int
	for i, v := range player.ActivityModule.Recharge {
		if v.ActivityID == req.ActivityID {
			totalRechargeInfo = &player.ActivityModule.Recharge[i]
			activityIndex = i
			break
		}
	}

	if totalRechargeInfo == nil {
		gamelog.Error("Hand_GetRechargeAward Error: Activity not exist %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	player.ActivityModule.CheckReset()

	//! 获取当前活动奖励返利
	awardType := G_GlobalVariables.GetActivityAwardType(req.ActivityID)
	awardLst := gamedata.GetRechargeInfo(awardType)
	if req.Index > len(awardLst) {
		gamelog.Error("Hand_GetRechargeAward Error: Invalid index %d", req.Index)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断是否有过领取
	if totalRechargeInfo.AwardMark.Get(uint(req.Index)) == true {
		gamelog.Error("Hand_GetRechargeAward Error: Aleady reveice")
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	awardInfo := awardLst[req.Index-1]

	//! 判断充值额度是否满足
	if totalRechargeInfo.RechargeValue < awardInfo.Recharge {
		gamelog.Error("Hand_GetRechargeAward Error: Rechare value not enough")
		response.RetCode = msg.RE_NOT_RECHARGE
		return
	}

	//! 发放奖励
	award := gamedata.GetItemsFromAwardID(awardInfo.Award)
	player.BagMoudle.AddAwardItems(award)

	//! 修改标记
	totalRechargeInfo.AwardMark.Set(uint(req.Index))
	go totalRechargeInfo.DB_UpdateRechargeMark(activityIndex, int(totalRechargeInfo.AwardMark))

	response.RetCode = msg.RE_SUCCESS

	for _, v := range award {
		var awardData msg.MSG_ItemData
		awardData.ID = v.ItemID
		awardData.Num = v.ItemNum
		response.AwardItem = append(response.AwardItem, awardData)
	}

	response.AwardMark = int(totalRechargeInfo.AwardMark)
}

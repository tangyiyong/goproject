package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家请求周周盈状态
func Hand_GetWeekAwardStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetWeekAwardStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetWeekAwardStatus Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetWeekAwardStatus_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.WeekAward.ActivityID) == false {
		gamelog.Error("Hand_GetWeekAwardStatus Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	response.RechargeNum = player.ActivityModule.WeekAward.RechargeNum
	response.LoginDay = player.ActivityModule.WeekAward.LoginDay
	response.AwardMark = int(player.ActivityModule.WeekAward.AwardMark)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求周周盈奖励
func Hand_GetWeekAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetWeekAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetWeekAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetWeekAward_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.WeekAward.ActivityID) == false {
		gamelog.Error("Hand_GetWeekAward Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 获取奖励信息
	awardType := G_GlobalVariables.GetActivityAwardType(player.ActivityModule.WeekAward.ActivityID)
	awardInfo := gamedata.GetWeekAwardInfo(awardType, req.Index)
	if awardInfo == nil {
		gamelog.Error("GetWeekAwardInfo Error: Invalid Param %d", req.Index)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查是否领取
	if player.ActivityModule.WeekAward.AwardMark.Get(uint32(req.Index)) != false {
		gamelog.Error("Hand_GetWeekAward Error: Aleady received")
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 检查登录天数与充值是否满足
	if player.ActivityModule.WeekAward.LoginDay < awardInfo.LoginDay {
		gamelog.Error("Hand_GetWeekAward Error: Not enough LoginDay")
		response.RetCode = msg.RE_NOT_ENOUGH_LOGIN_DAY
		return
	}

	if player.ActivityModule.WeekAward.RechargeNum < awardInfo.RechargeNum {
		gamelog.Error("Hand_GetWeekAward Error: Not enough Recharge")
		response.RetCode = msg.RE_NOT_RECHARGE
		return
	}

	awardLst := gamedata.GetItemsFromAwardID(awardInfo.AwardID)

	//! 判断是否为多选一
	if awardInfo.IsSelect != 0 {
		if req.Select > len(awardLst) || req.Select <= 0 {
			gamelog.Error("Hand_GetWeekAward Error: Invalid select %d", req.Select)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}

		award := awardLst[req.Select-1]
		player.BagMoudle.AddAwardItem(award.ItemID, award.ItemNum)

		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{award.ItemID, award.ItemNum})
	} else {
		player.BagMoudle.AddAwardItems(awardLst)

		for i := 0; i < len(awardLst); i++ {
			response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{awardLst[i].ItemID, awardLst[i].ItemNum})
		}
	}

	player.ActivityModule.WeekAward.AwardMark.Set(uint32(req.Index))
	go player.ActivityModule.WeekAward.DB_UpdateAwardMark()

	response.RetCode = msg.RE_SUCCESS
	response.AwardMark = int(player.ActivityModule.WeekAward.AwardMark)

}

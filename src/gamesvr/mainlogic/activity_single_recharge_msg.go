package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 查询单笔充值活动信息
func Hand_QueryActivitySingleRechargeInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryActivity_SingleRecharge_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryActivitySingleRechargeInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryActivity_SingleRecharge_Ack
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

	var singleRecharge *TActivitySingleRecharge
	var activityIndex int
	for i, v := range player.ActivityModule.SingleRecharge {
		if v.ActivityID == req.ActivityID {
			singleRecharge = &player.ActivityModule.SingleRecharge[i]
			activityIndex = i
			break
		}
	}

	if singleRecharge == nil {
		gamelog.Error("Hand_GetSingleRechargeAward Error: Activity not exist %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	response.RetCode = msg.RE_SUCCESS
	response.ActivityID = req.ActivityID
	response.AwardType = G_GlobalVariables.GetActivityAwardType(req.ActivityID)
	activtiyInfo := gamedata.GetActivityInfo(req.ActivityID)
	AwardMark := gamedata.GetRechargeInfo(activtiyInfo.AwardType)
	if len(AwardMark) <= 0 {
		gamelog.Error("Hand_GetActivity Error: GetRechargeInfo nil")
		return
	}

	for i, v := range AwardMark {
		var info msg.MSG_SingleRecharge
		info.Index = i + 1

		//! 判断次数是否足够
		activityTimes, _ := singleRecharge.GetSingleRechargeAwardTimes(info.Index)
		if activityTimes == nil {
			activityTimes = &TActivityRechargeInfo{info.Index, 0}
			singleRecharge.SingleAwardLst = append(singleRecharge.SingleAwardLst, *activityTimes)
			singleRecharge.DB_AddSingleRecharge(activityIndex, *activityTimes)
		}

		info.Times = v.Times - activityTimes.Times
		if activityTimes.Times < v.Times {
			info.Status = 1
		} else {
			info.Status = 2
		}

		//! 检查是否有过单充记录
		isHave, _ := player.ActivityModule.CheckSingleRecharge(req.ActivityID, v.Recharge)

		if isHave == false && info.Status != 2 {
			info.Status = 0
		}

		response.SingleRechargeLst = append(response.SingleRechargeLst, info)
	}
}

//! 玩家请求领取单充奖励
func Hand_GetSingleRechargeAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSingleAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSingleRechargeAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetSingleAward_Ack
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

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	//! 检测当前是否有此活动
	var singleRecharge *TActivitySingleRecharge
	var activityIndex int
	for i, v := range player.ActivityModule.SingleRecharge {
		if v.ActivityID == req.ActivityID {
			singleRecharge = &player.ActivityModule.SingleRecharge[i]
			activityIndex = i
			break
		}
	}

	if singleRecharge == nil {
		gamelog.Error("Hand_GetSingleRechargeAward Error: Activity not exist %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	awardType := G_GlobalVariables.GetActivityAwardType(req.ActivityID)
	AwardMark := gamedata.GetRechargeInfo(awardType)

	if len(AwardMark) < req.Index {
		gamelog.Error("Hand_GetSingleRechargeAward Error: Inavlid index")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	award := AwardMark[req.Index-1]

	//! 判断次数是否足够
	activityTimes, infoIndex := singleRecharge.GetSingleRechargeAwardTimes(req.Index)
	if activityTimes == nil {
		gamelog.Error("Hand_GetSingleRechargeAward Error: GetSingleRechargeAwardTimes nil.")
		return
	}

	//! 检查是否有过单充记录
	isHave, index := player.ActivityModule.CheckSingleRecharge(req.ActivityID, award.Recharge)

	if isHave == false {
		gamelog.Error("Hand_GetSingleRechargeAward Error: Recharge not enough")
		response.RetCode = msg.RE_NOT_RECHARGE
		return
	}

	if activityTimes.Times >= award.Times {
		gamelog.Error("Hand_GetSingleRechargeAward Error: Receive times is use up.")
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 修改充值记录状态
	singleRecharge.RechargeRecord[index].Status = 1
	player.ActivityModule.DB_UpdateRechargeRecord(activityIndex, index, 1)

	//! 修改领取次数
	activityTimes.Times += 1
	player.ActivityModule.DB_UpdateSingelAward(activityIndex, infoIndex, activityTimes.Times)

	//! 给予奖励
	awardLst := gamedata.GetItemsFromAwardID(award.Award)
	player.BagMoudle.AddAwardItems(awardLst)

	response.RetCode = msg.RE_SUCCESS

	for _, v := range awardLst {
		var awardData msg.MSG_ItemData
		awardData.ID = v.ItemID
		awardData.Num = v.ItemNum
		response.AwardItem = append(response.AwardItem, awardData)
	}

}

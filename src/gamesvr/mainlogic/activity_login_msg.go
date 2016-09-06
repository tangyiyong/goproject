package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家查询累计登录活动信息
func Hand_QueryActivityLoginInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryActivity_Login_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryActivityLoginInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryActivity_Login_Ack
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

	var activity *TActivityLogin
	for i, v := range player.ActivityModule.Login {
		if v.ActivityID == req.ActivityID {
			activity = &player.ActivityModule.Login[i]
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_QueryActivityLoginInfo Error: Activity not exist %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	response.RetCode = msg.RE_SUCCESS

	//! 登录活动信息
	response.AwardType = G_GlobalVariables.GetActivityAwardType(req.ActivityID)
	response.AwardMark = int(activity.LoginAward)
	response.LoginDay = activity.LoginDay
	response.ActivityID = req.ActivityID
}

//! 玩家请求领取累计登录奖励
func Hand_GetActivityLoginAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetActivity_LoginAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetActivityLoginInfo : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetActivity_LoginAward_Ack
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
	var activity *TActivityLogin
	var activityIndex int
	for i, v := range player.ActivityModule.Login {
		if v.ActivityID == req.ActivityID {
			activity = &player.ActivityModule.Login[i]
			activityIndex = i
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_GetActivityLoginInfo Error: Activity not exist %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测当前奖励是否能够领取
	if activity.LoginDay < req.Index {
		gamelog.Error("Hand_GetActivityLoginInfo Error: Login day is not enough %d", req.Index)
		response.RetCode = msg.RE_NOT_ENOUGH_LOGIN_DAY
		return
	}

	//! 检测是否已经领取
	if activity.LoginAward.Get(uint32(req.Index)) == true {
		gamelog.Error("Hand_GetActivityLoginInfo Error: Aleady recvice award %d", req.Index)
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 获取奖励
	activityInfo := gamedata.GetActivityInfo(req.ActivityID)
	loginAwardLst := gamedata.GetActivityLoginInfo(activityInfo.AwardType)
	if len(loginAwardLst) < req.Index {
		gamelog.Error("Hand_GetActivityLoginInfo Error: Login award error len(loginAwardLst): %d req.Index: %d", len(loginAwardLst), req.Index)
		return
	}

	loginAward := loginAwardLst[req.Index-1]
	awardLst := gamedata.GetItemsFromAwardID(loginAward.Award)
	if loginAward.IsSelect == 1 {
		if len(awardLst) < req.Choice {
			gamelog.Error("Hand_GetActivityLoginInfo Error: Login award error len(awardLst): %d req.Index: %d", len(awardLst), req.Index)
			return
		}

		player.BagMoudle.AddAwardItem(awardLst[req.Choice-1].ItemID, awardLst[req.Choice-1].ItemNum)
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{awardLst[req.Choice-1].ItemID, awardLst[req.Choice-1].ItemNum})
	} else {
		player.BagMoudle.AddAwardItems(awardLst)

		for _, v := range awardLst {
			response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum})
		}
	}

	//! 改变玩家标记
	activity.LoginAward.Set(uint32(req.Index))
	activity.DB_UpdateLoginAward(activityIndex)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
	response.AwardMark = int(activity.LoginAward)
}

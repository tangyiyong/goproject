package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

//! 玩家请求查询奖励中心信息
func Hand_GetAwardCenterInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_AwardCenter_Query_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetAwardCenterInfo unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_AwardCenter_Query_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 先检查是否有全服奖励可领取
	SendSvrAwardToPlayer(req.PlayerID)

	//! 查询信息
	response.AwardLst = []msg.MSG_AwardCenter_Data{}
	for _, v := range player.AwardCenterModule.AwardLst {
		data := msg.MSG_AwardCenter_Data{}
		data.ID = v.ID
		data.TextType = v.TextType
		data.Value = v.Value
		data.Time = v.Time

		for _, i := range v.ItemLst {
			data.ItemLst = append(data.ItemLst, msg.MSG_ItemData{i.ItemID, i.ItemNum})
		}

		response.AwardLst = append(response.AwardLst, data)
	}

	//! 反馈结果
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求一键领取领奖中心奖励
func Hand_RecvAwardCenterAwardOneyKey(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_OneKey_AwardCenter_Get_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetAwardCenterInfo unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_OneKey_AwardCenter_Get_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	for _, v := range player.AwardCenterModule.AwardLst {
		player.BagMoudle.AddAwardItems(v.ItemLst)

		for _, n := range v.ItemLst {
			var award msg.MSG_ItemData
			award.ID = n.ItemID
			award.Num = n.ItemNum
			response.AwardLst = append(response.AwardLst, award)
		}
	}

	//! 清空奖励列表
	player.AwardCenterModule.AwardLst = []TAwardData{}
	player.AwardCenterModule.DB_UpdateDatabaseLst()
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求领取领奖中心奖励
func Hand_RecvAwardCenter(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_AwardCenter_Get_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetAwardCenterInfo unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_AwardCenter_Get_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 获取奖励内容
	awardInfo := player.AwardCenterModule.GetAwardData(req.AwardID)
	if awardInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	player.BagMoudle.AddAwardItems(awardInfo.ItemLst)

	player.AwardCenterModule.RemoveAward(req.AwardID)
	response.RetCode = msg.RE_SUCCESS
}

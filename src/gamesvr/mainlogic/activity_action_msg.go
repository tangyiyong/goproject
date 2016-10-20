package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"msg"
	"net/http"
	"time"
)

//! 查询领取体力活动信息
func Hand_QueryActivityActionInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryActivity_Action_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryActivityActionInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryActivity_Action_Ack
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

	//	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
	//		return
	//	}

	//	defer player.FinishMsgProcess()

	player.ActivityModule.CheckReset()

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.ReceiveAction.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		gamelog.Error("Hand_QueryActivityActionInfo Error: Not open")
		return
	}

	response.RetCode = msg.RE_SUCCESS

	//! 领取体力
	response.RecvAction = uint32(player.ActivityModule.ReceiveAction.RecvAction)
	response.RetroactiveCostMoneyID = gamedata.ActionActivityRetroactiveMoneyID
	response.RetroactiveCostMoneyNum = gamedata.ActionActivityRetroactiveMoneyNum
}

//! 玩家请求领取体力
func Hand_ReceiveActivityAction(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetActivity_Action_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ReceiveActivityAction : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetActivity_Action_Ack
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
	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.ReceiveAction.ActivityID) == false {
		gamelog.Error("Hand_ReceiveActivityAction Error: Activity not open.")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 判断玩家是否已经领取
	if player.ActivityModule.ReceiveAction.RecvAction.Get(req.Index) == true {
		gamelog.Error("Hand_ActionRetroactive Error: Aleady received:mark:%d, index:%d", uint32(player.ActivityModule.ReceiveAction.RecvAction), req.Index)
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 获取奖励信息
	awardType := G_GlobalVariables.GetActivityAwardType(player.ActivityModule.ReceiveAction.ActivityID)
	pAtyAction := gamedata.GetActivityAction(awardType)
	if pAtyAction == nil {
		gamelog.Error("Hand_ReceiveActivityAction Error: GetRecvAction nil")
		return
	}

	if false == pAtyAction.IsTimeOK(req.Index) {
		gamelog.Error("Hand_ReceiveActivityAction Error: Cur index is not time ok:%d.", req.Index)
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 修改领取标记
	player.ActivityModule.ReceiveAction.RecvAction.Set(req.Index)
	player.ActivityModule.ReceiveAction.DB_Refresh()

	response.Index = req.Index

	//! 增加玩家体力
	player.RoleMoudle.AddAction(pAtyAction.ActionID, pAtyAction.ActionNum)

	//! 获取用户体力以及下次恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(pAtyAction.ActionID)

	//! 随机奖励
	randValue := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := randValue.Intn(1000)
	if random < pAtyAction.AwardPro {
		player.RoleMoudle.AddMoney(pAtyAction.MoneyID, pAtyAction.MoneyNum)
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{pAtyAction.MoneyID, pAtyAction.MoneyNum})
	}

	response.RecvAction = uint32(player.ActivityModule.ReceiveAction.RecvAction)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求补签领取体力
func Hand_ActionRetroactive(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetAction_Retroactive_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ActionRetroactive : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetAction_Retroactive_Ack
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
	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.ReceiveAction.ActivityID) == false {
		gamelog.Error("Hand_ActionRetroactive Error: Activity not open.")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 获取领取奖励信息
	activityInfo := gamedata.GetActivityInfo(player.ActivityModule.ReceiveAction.ActivityID)
	pAtyAction := gamedata.GetActivityAction(activityInfo.AwardType)
	if pAtyAction == nil {
		gamelog.Error("Hand_ActionRetroactive Error: Invalid index %d", req.Index)
		return
	}

	//! 检测当前时间是否超过最后领取时间
	now := time.Now()
	sec := now.Hour()*3600 + now.Minute()*60 + now.Second()
	if sec < pAtyAction.Time_End[req.Index-1] {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ActionRetroactive Error: Time yet")
		return
	}

	//! 检测是否已领取该时间段奖励
	if player.ActivityModule.ReceiveAction.RecvAction.Get(req.Index) == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		gamelog.Error("Hand_ActionRetroactive Error: Aleady received:mark:%d, index:%d", uint32(player.ActivityModule.ReceiveAction.RecvAction), req.Index)
		return
	}

	//! 检测玩家货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.ActionActivityRetroactiveMoneyID, gamedata.ActionActivityRetroactiveMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_ActionRetroactive Error: Not enough money")
		return
	}

	response.Index = req.Index

	//! 扣除玩家货币
	player.RoleMoudle.CostMoney(gamedata.ActionActivityRetroactiveMoneyID, gamedata.ActionActivityRetroactiveMoneyNum)
	response.CostItem = append(response.CostItem, msg.MSG_ItemData{gamedata.ActionActivityRetroactiveMoneyID, gamedata.ActionActivityRetroactiveMoneyNum})

	//! 修改领取标记
	player.ActivityModule.ReceiveAction.RecvAction.Set(req.Index)
	player.ActivityModule.ReceiveAction.DB_Refresh()

	//! 增加玩家体力
	player.RoleMoudle.AddAction(pAtyAction.ActionID, pAtyAction.ActionNum)

	//! 获取用户体力以及下次恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(pAtyAction.ActionID)

	//! 随机奖励
	randValue := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := randValue.Intn(1000)
	if random < pAtyAction.AwardPro {
		player.RoleMoudle.AddMoney(pAtyAction.MoneyID, pAtyAction.MoneyNum)
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{pAtyAction.MoneyID, pAtyAction.MoneyNum})
	}

	response.RecvAction = uint32(player.ActivityModule.ReceiveAction.RecvAction)
	response.RetCode = msg.RE_SUCCESS
}

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

//! 查询迎财神活动信息
func Hand_QueryActivityMoneyGodInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryActivity_MoneyGod_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryActivityMoneyGodInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryActivity_MoneyGod_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.MoneyGod.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		gamelog.Error("Hand_QueryActivityMoneyGodInfo Error: Not open")
		return
	}

	response.RetCode = msg.RE_SUCCESS

	//! 迎财神
	response.CurrentTimes = player.ActivityModule.MoneyGod.CurrentTimes
	response.TotalMoney = player.ActivityModule.MoneyGod.TotalMoney

	response.NextTime = player.ActivityModule.MoneyGod.NextTime - time.Now().Unix()
	if response.NextTime < 0 {
		response.NextTime = 0
	}
	response.CumulativeTimes = player.ActivityModule.MoneyGod.CumulativeTimes
}

//! 玩家请求迎财神
func Hand_WelcomeMoneyGold(w http.ResponseWriter, r *http.Request) {

	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_WelcomeMoneyGod_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_WelcomeMoneyGold : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_WelcomeMoneyGod_Ack
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

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.MoneyGod.ActivityID) == false {
		gamelog.Error("Hand_WelcomeMoneyGold Error: Activity not open.")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.ActivityModule.CheckReset()

	activityInfo := gamedata.GetActivityInfo(player.ActivityModule.MoneyGod.ActivityID)

	//! 获取静态配置
	moneyInfo := gamedata.GetMoneyGoldInfo(activityInfo.AwardType)
	if moneyInfo == nil {
		gamelog.Error("Hand_WelcomeMoneyGold Error: GetMoneyGoldInfo nil")
		return
	}

	//! 检测当前领取时间
	if player.ActivityModule.MoneyGod.NextTime != 0 {
		gamelog.Error("Hand_WelcomeMoneyGold Error: Not to receive time")
		response.RetCode = msg.RE_NOT_REACH_TIME
		return
	}

	//! 检查是否满足领取聚宝盆条件
	if player.ActivityModule.MoneyGod.CumulativeTimes >= moneyInfo.AwardTimes {
		gamelog.Error("Hand_WelcomeMoneyGold Error: Please get Ju Bao Pen award first")
		response.RetCode = msg.RE_PLEASE_GET_JUBAOPENG
		return
	}

	//! 检查次数
	if player.ActivityModule.MoneyGod.CurrentTimes <= 0 {
		gamelog.Error("Hand_WelcomeMoneyGold Error: Times is zero")
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 计算获取银币
	response.MoneyID = moneyInfo.MoneyID
	response.MoneyNum = moneyInfo.MoneyNum * player.GetLevel()

	//! 设置时间
	player.ActivityModule.MoneyGod.NextTime = time.Now().Unix() + int64(moneyInfo.CDTime)
	player.ActivityModule.MoneyGod.CurrentTimes -= 1
	player.ActivityModule.MoneyGod.CumulativeTimes += 1
	player.ActivityModule.MoneyGod.TotalMoney += response.MoneyNum

	response.TotalMoney = player.ActivityModule.MoneyGod.TotalMoney
	response.NextTime = player.ActivityModule.MoneyGod.NextTime - time.Now().Unix()
	if response.NextTime < 0 {
		response.NextTime = 0
	}

	player.ActivityModule.MoneyGod.DB_Refresh()

	//! 给予玩家银币
	player.RoleMoudle.AddMoney(response.MoneyID, response.MoneyNum)

	//! 计算幸运道具奖励
	randValue := rand.New(rand.NewSource(time.Now().UnixNano()))

	random := randValue.Intn(1000)

	if random < moneyInfo.LuckPro {
		//! 额外奖励
		player.BagMoudle.AddAwardItem(moneyInfo.ItemID, moneyInfo.ItemNum)
		response.ExAwardID = moneyInfo.ItemID
		response.ExAwardNum = moneyInfo.ItemNum
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求聚宝盆奖励
func Hand_MoneyGoldAward(w http.ResponseWriter, r *http.Request) {

	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetMoneyGodTotalAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MoneyGoldAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetMoneyGodTotalAward_Ack
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

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.MoneyGod.ActivityID) == false {
		gamelog.Error("Hand_MoneyGoldAward Error: Activity not open.")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.ActivityModule.CheckReset()

	activityInfo := gamedata.GetActivityInfo(player.ActivityModule.MoneyGod.ActivityID)

	//! 获取静态配置
	moneyInfo := gamedata.GetMoneyGoldInfo(activityInfo.AwardType)

	if player.ActivityModule.MoneyGod.CumulativeTimes < moneyInfo.AwardTimes {
		gamelog.Error("Hand_MoneyGoldAward Error: Not enough times")
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 给予玩家聚宝盆累积奖励
	player.RoleMoudle.AddMoney(moneyInfo.MoneyID, player.ActivityModule.MoneyGod.TotalMoney)
	response.MoneyID = moneyInfo.MoneyID
	response.MoneyNum = player.ActivityModule.MoneyGod.TotalMoney

	//! 清空次数
	player.ActivityModule.MoneyGod.CumulativeTimes = 0
	player.ActivityModule.MoneyGod.TotalMoney = 0
	player.ActivityModule.MoneyGod.DB_UpdateCumulativeTimes()

	response.RetCode = msg.RE_SUCCESS
}

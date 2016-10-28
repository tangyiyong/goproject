package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 查询VIP日常福利领取状态
func Hand_GetDailyVipStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetVipDailyWelfareStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_DailyVipWelfare Unmarshal fail. Error: %v", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetVipDailyWelfareStatus_Ack
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

	if player.ActivityModule.VipGift.IsRecvWelfare == true {
		response.SignStatus = 2
	} else {
		response.SignStatus = 1
	}

	//! 周一零点刷新
	player.ActivityModule.VipGift.CheckWeekGiftRefresh()

	//! 获取礼包内容
	for _, v := range player.ActivityModule.VipGift.WeekGift {
		var gift msg.MSG_WeekGiftInfo
		gift.ID = v.ID
		gift.BuyTimes = v.BuyTimes
		response.GiftLst = append(response.GiftLst, gift)
	}

	response.RetCode = msg.RE_SUCCESS

}

//! 领取VIP日常福利
func Hand_DailyVipWelfare(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetVipDailyWelfare_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_DailyVipWelfare Unmarshal fail. Error: %v", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetVipDailyWelfare_Ack
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

	if player.ActivityModule.VipGift.IsRecvWelfare == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		gamelog.Error("Hand_DailyVipWelfare already recevied. playerID: %v ", player.playerid)
		return
	}

	player.ActivityModule.VipGift.IsRecvWelfare = true

	//! 发送福利
	//! 获取VIP信息
	info := gamedata.GetVipInfo(player.GetVipLevel())

	//! 获取日常奖励信息
	awardLst := gamedata.GetItemsFromAwardID(info.VipAward)

	//! 发送奖励
	player.BagMoudle.AddAwardItems(awardLst)

	//! 更新至数据库
	player.ActivityModule.VipGift.DB_SaveDailyResetTime()

	for _, v := range awardLst {
		award := msg.MSG_ItemData{v.ItemID, v.ItemNum}
		response.AwardItem = append(response.AwardItem, award)
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买VIP每周礼包
func Hand_BuyVipWeekGift(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyVipWeekGiftInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetVipWeekGift Unmarshal fail. Error: %v", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuyVipWeekGiftInfo_Ack
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

	//! 检查玩家是否能够购买
	money := 0
	for i, _ := range player.ActivityModule.VipGift.WeekGift {

		if player.ActivityModule.VipGift.WeekGift[i].ID == req.ID {
			data := gamedata.GetVipWeekItemFromID(player.ActivityModule.VipGift.WeekGift[i].ID)

			if player.ActivityModule.VipGift.WeekGift[i].BuyTimes+req.BuyTimes > data.BuyTimes {
				//! 购买次数不足
				response.RetCode = msg.RE_NOT_ENOUGH_TIMES
				return
			}

			money = data.MoneyNum
			if player.RoleMoudle.CheckMoneyEnough(data.MoneyID, data.MoneyNum*req.BuyTimes) == false {
				//! 钻石不足
				response.RetCode = msg.RE_NOT_ENOUGH_MONEY
				return
			}

			response.MoneyID = gamedata.VipWeeklyGiftMoneyID
			response.MoneyNum = money * req.BuyTimes

			//! 扣除金钱,增加次数
			player.RoleMoudle.CostMoney(gamedata.VipWeeklyGiftMoneyID, money*req.BuyTimes)
			player.ActivityModule.VipGift.WeekGift[i].BuyTimes += req.BuyTimes

			response.BuyTimes = player.ActivityModule.VipGift.WeekGift[i].BuyTimes

			//! 发货
			itemLst := gamedata.GetItemsFromAwardID(data.Award)
			for i, _ := range itemLst {
				itemLst[i].ItemNum *= req.BuyTimes

				response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{itemLst[i].ItemID, itemLst[i].ItemNum})
			}

			player.BagMoudle.AddAwardItems(itemLst)

			player.ActivityModule.VipGift.DB_UpdateBuyTimes(player.ActivityModule.VipGift.WeekGift[i].ID, player.ActivityModule.VipGift.WeekGift[i].BuyTimes)
			break
		}
	}

	response.ID = req.ID
	response.RetCode = msg.RE_SUCCESS
}

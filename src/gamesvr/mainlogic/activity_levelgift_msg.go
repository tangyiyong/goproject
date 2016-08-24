package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家请求等级礼包信息
func Hand_GetLevelGiftInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetLevelGiftInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetLevelGiftInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetLevelGiftInfo_Ack
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
	player.ActivityModule.LevelGift.CheckDeadLine()

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.LevelGift.ActivityID) == false {
		gamelog.Error("Hand_GetLevelGiftInfo Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	giftLst := player.ActivityModule.LevelGift.GiftLst
	length := len(player.ActivityModule.LevelGift.GiftLst)
	for i := 0; i < length; i++ {
		var gift msg.MSG_LevelGiftInfo
		gift.ID = giftLst[i].GiftID
		gift.BuyTimes = giftLst[i].BuyTimes
		gift.DeadLine = giftLst[i].DeadLine
		response.GiftLst = append(response.GiftLst, gift)
	}

	//! 因为查看了购买列表, 将新商品标记置位false
	player.ActivityModule.LevelGift.IsHaveNewItem = false
	go player.ActivityModule.LevelGift.DB_UpdateNewItemMark()

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买等级礼包
func Hand_BuyLevelGift(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyLevelGift_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyLevelGift Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_BuyLevelGift_Ack
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
	player.ActivityModule.LevelGift.CheckDeadLine()

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.LevelGift.ActivityID) == false {
		gamelog.Error("Hand_BuyLevelGift Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	awardType := G_GlobalVariables.GetActivityAwardType(player.ActivityModule.LevelGift.ActivityID)

	//! 获取等级礼包信息
	levelGift := player.ActivityModule.LevelGift.GetLevelGiftInfo(req.GiftID)
	if levelGift == nil {
		gamelog.Error("Hand_BuyLevelGift Error: Not find gift info. id: %d", req.GiftID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if levelGift.BuyTimes <= 0 {
		gamelog.Error("Hand_BuyLevelGift Error: BuyTimes not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	levelGiftInfo := gamedata.GetLevelGiftInfo(awardType, player.ActivityModule.LevelGift.ActivityID)
	if levelGiftInfo == nil {
		gamelog.Error("Hand_BuyLevelGift Error: GetLevelGiftInfo nil")
		response.RetCode = msg.RE_UNKNOWN_ERR
		return
	}

	if levelGiftInfo.MoneyID != 0 {
		//! 收费领取

		//! 检测货币是否足够
		if player.RoleMoudle.CheckMoneyEnough(levelGiftInfo.MoneyID, levelGiftInfo.MoneyNum) == false {
			gamelog.Error("Hand_BuyLevelGift Error: Not enough money ID: %d", req.GiftID)
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			return
		}

		//! 扣除货币
		player.RoleMoudle.CostMoney(levelGiftInfo.MoneyID, levelGiftInfo.MoneyNum)
		response.CostMoneyID = levelGiftInfo.MoneyID
		response.CostMoneyNum = levelGiftInfo.MoneyNum
	}

	levelGift.BuyTimes -= 1
	go player.ActivityModule.LevelGift.DB_UpdateBuyTimes(levelGift.GiftID, levelGift.BuyTimes)
	response.BuyTimes = levelGift.BuyTimes

	//! 给予商品
	awardLst := gamedata.GetItemsFromAwardID(levelGiftInfo.Award)
	player.BagMoudle.AddAwardItems(awardLst)

	for _, v := range awardLst {
		response.AwardItem = append(response.AwardItem,
			msg.MSG_ItemData{v.ItemID, v.ItemNum})
	}

	response.RetCode = msg.RE_SUCCESS
}

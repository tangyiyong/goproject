package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"msg"
	"net/http"
	"time"
	"utility"
)

//! 查询巡回探宝状态
func Hand_QueryHuntTreasure(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryHuntTreasure_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryHuntTreasure : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryHuntTreasure_Ack
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
	//! 获取活动
	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.HuntTreasure.ActivityID) == false {
		gamelog.Error("Hand_QueryHuntTreasure Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	_, response.EndCountDown = G_GlobalVariables.IsActivityTime(player.ActivityModule.HuntTreasure.ActivityID)

	response.CurrentPos = player.ActivityModule.HuntTreasure.CurrentPos
	response.HuntTurns = player.ActivityModule.HuntTreasure.HuntTurns
	response.Score = player.ActivityModule.HuntTreasure.TodayScore[utility.GetCurDayMod()]
	response.IsHaveStore = player.ActivityModule.HuntTreasure.IsHaveStore
	response.TotalRank = -1
	response.FreeTimes = player.ActivityModule.HuntTreasure.FreeTimes
	response.TodayRank = -1
	for i, v := range G_HuntTreasureTotalRanker.List {
		if v.RankID == player.playerid {
			response.TotalRank = i + 1
			break
		}
	}

	for i, v := range G_HuntTreasureTodayRanker.List {
		if v.RankID == player.playerid {
			response.TodayRank = i + 1
			break
		}
	}

	response.RetCode = msg.RE_SUCCESS
	response.AwardType = G_GlobalVariables.GetActivityAwardType(player.ActivityModule.HuntTreasure.ActivityID)
}

//! 玩家开始掷骰
func Hand_StartHuntTreasure(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.Msg_StartHuntTreasure_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_StartHuntTreasure : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.Msg_StartHuntTreasure_Ack
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

	//! 获取活动
	activityID := player.ActivityModule.HuntTreasure.ActivityID

	if G_GlobalVariables.IsActivityOpen(activityID) == false {
		gamelog.Error("Hand_StartHuntTreasure Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 判断是否为活动时间
	isEnd, _ := G_GlobalVariables.IsActivityTime(player.ActivityModule.HuntTreasure.ActivityID)
	if isEnd == false {
		gamelog.Error("Hand_StartHuntTreasure Error: Activity is over")
		response.RetCode = msg.RE_ACTIVITY_IS_OVER
		return
	}

	if req.IsUseLucklyDice == 0 {
		if req.IsStartTenTimes == 1 {
			needItemNum := 10 - player.ActivityModule.HuntTreasure.FreeTimes
			if player.BagMoudle.IsItemEnough(gamedata.HuntTicketItemID, needItemNum) == false {
				//! 判断钱
				if player.RoleMoudle.CheckMoneyEnough(gamedata.HuntCostMoneyID, gamedata.HuntCostMoneyNum*10) == false {
					gamelog.Error("Hand_StartHuntTreasure Error: %v Not enough game money", player.playerid)
					response.RetCode = msg.RE_ITEM_NOT_ENOUGH
					return
				} else {
					player.RoleMoudle.CostMoney(gamedata.HuntCostMoneyID, gamedata.HuntCostMoneyNum*10)
					response.CostItem = msg.MSG_ItemData{gamedata.HuntCostMoneyID, gamedata.HuntCostMoneyNum * 10}
				}
			} else {
				response.CostFreeTimes = player.ActivityModule.HuntTreasure.FreeTimes
				player.ActivityModule.HuntTreasure.FreeTimes = 0
				player.BagMoudle.RemoveNormalItem(gamedata.HuntTicketItemID, needItemNum)
				response.CostItem = msg.MSG_ItemData{gamedata.HuntTicketItemID, needItemNum}
				player.ActivityModule.HuntTreasure.DB_SaveFreeTiems()
			}
		} else {
			if player.ActivityModule.HuntTreasure.FreeTimes >= 1 {
				player.ActivityModule.HuntTreasure.FreeTimes -= 1
				response.CostFreeTimes = 1
				player.ActivityModule.HuntTreasure.DB_SaveFreeTiems()
			} else if player.BagMoudle.IsItemEnough(gamedata.HuntTicketItemID, 1) == true {
				player.BagMoudle.RemoveNormalItem(gamedata.HuntTicketItemID, 1)
				response.CostItem = msg.MSG_ItemData{gamedata.HuntTicketItemID, 1}
			} else if player.RoleMoudle.CheckMoneyEnough(gamedata.HuntCostMoneyID, gamedata.HuntCostMoneyNum) == true {
				player.RoleMoudle.CostMoney(gamedata.HuntCostMoneyID, gamedata.HuntCostMoneyNum)
				response.CostItem = msg.MSG_ItemData{gamedata.HuntCostMoneyID, gamedata.HuntCostMoneyNum}
			} else {
				gamelog.Error("Hand_StartHuntTreasure Error: %v Not enough game money", player.playerid)
				response.RetCode = msg.RE_ITEM_NOT_ENOUGH
				return
			}
		}

	}

	//! 随机步数
	steps := 0
	curPos := player.ActivityModule.HuntTreasure.CurrentPos
	moveStep := IntLst{}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	if req.IsUseLucklyDice == 1 {
		//! 使用幸运骰子
		if player.BagMoudle.IsItemEnough(gamedata.LuckyDiceItemID, 1) == false {
			if player.RoleMoudle.CheckMoneyEnough(gamedata.HuntCostMoneyID, 30) == false {
				gamelog.Error("Hand_StartHuntTreasure Error: %v Not enough luck dice", player.playerid)
				response.RetCode = msg.RE_ITEM_NOT_ENOUGH
				return
			} else {
				player.RoleMoudle.CostMoney(gamedata.HuntCostMoneyID, 30)
				response.CostItem = msg.MSG_ItemData{gamedata.HuntCostMoneyID, 30}
			}

		} else {
			player.BagMoudle.RemoveNormalItem(gamedata.LuckyDiceItemID, 1)
			response.CostItem = msg.MSG_ItemData{gamedata.LuckyDiceItemID, 1}
		}

		if req.Steps > 6 || req.Steps < 1 {
			gamelog.Error("Hand_StartHuntTreasure Error: Invlid steps %d", req.Steps)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}

		steps = req.Steps
	} else if req.IsStartTenTimes == 0 {
		//! 普通随机, 六分之一
		randValue := random.Intn(6000)
		curValue := 0
		for i := 1; i <= 6; i++ {
			if randValue >= curValue && randValue < curValue+1000 {
				steps = i
				break
			}
			curValue += 1000
		}
	} else if req.IsStartTenTimes == 1 {
		for j := 0; j < 10; j++ {
			randValue := random.Intn(6000)
			curValue := 0
			for i := 1; i <= 6; i++ {
				if randValue >= curValue && randValue < curValue+1000 {
					steps += i
					moveStep.Add(i)
					break
				}
				curValue += 1000
			}
		}
	}

	//! 消耗游戏券
	if req.IsStartTenTimes == 1 {
		player.BagMoudle.RemoveNormalItem(gamedata.HuntTicketItemID, 10)
	} else {
		player.BagMoudle.RemoveNormalItem(gamedata.HuntTicketItemID, 1)
	}

	response.RandomScore = steps

	//! 前进步数
	activityInfo := gamedata.GetActivityInfo(activityID)
	mapCount := gamedata.GetHuntTreasureMapCount(activityInfo.AwardType)
	player.ActivityModule.HuntTreasure.CurrentPos += steps

	for {
		if player.ActivityModule.HuntTreasure.CurrentPos > mapCount {
			player.ActivityModule.HuntTreasure.CurrentPos -= mapCount
			player.ActivityModule.HuntTreasure.HuntTurns += 1
			gamelog.Info("HuntTurns: %v  mapCount: %v   Pos: %v", player.ActivityModule.HuntTreasure.HuntTurns, mapCount, player.ActivityModule.HuntTreasure.CurrentPos)
		} else {
			break
		}
	}

	//! 获取积分
	player.ActivityModule.HuntTreasure.Score += steps

	indexToday := 0
	if utility.GetCurDayMod() == 1 {
		indexToday = 1
	}
	player.ActivityModule.HuntTreasure.TodayScore[indexToday] += steps
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_GET_HUNT_SCORE, steps)

	//! 判断当前所在格类型
	if req.IsStartTenTimes == 0 {
		for {
			mapInfo := gamedata.GetHuntTreasureMap(player.ActivityModule.HuntTreasure.CurrentPos, activityInfo.AwardType)
			response.MoveTypeLst = append(response.MoveTypeLst, mapInfo.Type)
			if mapInfo.Type == gamedata.HuntTreasureMap_Move {
				//! 额外随机1-3格前进
				exMove := random.Intn(3) + 1
				response.ExMove = append(response.ExMove, exMove)
				player.ActivityModule.HuntTreasure.CurrentPos += exMove
				if player.ActivityModule.HuntTreasure.CurrentPos > mapCount {
					player.ActivityModule.HuntTreasure.CurrentPos -= mapCount
					player.ActivityModule.HuntTreasure.HuntTurns += 1
				}
			} else if mapInfo.Type == gamedata.HuntTreasureMap_Store {
				itemLst := gamedata.RandHuntTreasureStoreItem(5, activityInfo.AwardType)
				player.ActivityModule.HuntTreasure.StoreItemLst = []THuntStoreItem{}

				player.ActivityModule.HuntTreasure.IsHaveStore = true
				player.ActivityModule.HuntTreasure.DB_SaveStoreMark()

				for _, v := range itemLst {
					player.ActivityModule.HuntTreasure.StoreItemLst = append(player.ActivityModule.HuntTreasure.StoreItemLst,
						THuntStoreItem{v, false})
				}

				player.ActivityModule.HuntTreasure.DB_UpdateHuntStore()
				break
			} else {
				if mapInfo.Award != 0 {
					awardLst := gamedata.GetItemsFromAwardID(mapInfo.Award)
					for _, v := range awardLst {
						response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum})
					}

					player.BagMoudle.AddAwardItems(awardLst)
				}

				break
			}
		}
	} else if req.IsStartTenTimes == 1 {
		for _, m := range moveStep {
			curPos += m
			if curPos > mapCount {
				curPos -= mapCount
				player.ActivityModule.HuntTreasure.HuntTurns += 1
			}
			for {
				mapInfo := gamedata.GetHuntTreasureMap(curPos, activityInfo.AwardType)
				response.MoveTypeLst = append(response.MoveTypeLst, mapInfo.Type)
				if mapInfo.Type == gamedata.HuntTreasureMap_Move {
					//! 额外随机1-3格前进
					exMove := random.Intn(3) + 1
					response.ExMove = append(response.ExMove, exMove)
					curPos += exMove
					if curPos > mapCount {
						curPos -= mapCount
						player.ActivityModule.HuntTreasure.HuntTurns += 1
					}
				} else if mapInfo.Type == gamedata.HuntTreasureMap_Store {
					itemLst := gamedata.RandHuntTreasureStoreItem(5, activityInfo.AwardType)

					player.ActivityModule.HuntTreasure.IsHaveStore = true
					player.ActivityModule.HuntTreasure.DB_SaveStoreMark()

					player.ActivityModule.HuntTreasure.StoreItemLst = []THuntStoreItem{}
					for _, v := range itemLst {
						player.ActivityModule.HuntTreasure.StoreItemLst = append(player.ActivityModule.HuntTreasure.StoreItemLst,
							THuntStoreItem{v, false})
					}

					player.ActivityModule.HuntTreasure.DB_UpdateHuntStore()
					break
				} else {
					if mapInfo.Award != 0 {
						awardLst := gamedata.GetItemsFromAwardID(mapInfo.Award)
						for _, v := range awardLst {
							response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum})
						}

						player.BagMoudle.AddAwardItems(awardLst)
					}

					break
				}
			}
		}

	}

	player.ActivityModule.HuntTreasure.DB_SaveHuntStatus()

	response.TodayRank = -1
	response.TotalRank = -1

	response.TodayRank = G_HuntTreasureTodayRanker.SetRankItem(player.playerid, player.ActivityModule.HuntTreasure.TodayScore[indexToday])
	response.TotalRank = G_HuntTreasureTotalRanker.SetRankItem(player.playerid, player.ActivityModule.HuntTreasure.Score)

	response.Score = player.ActivityModule.HuntTreasure.TodayScore[indexToday]
	response.CurrentPos = player.ActivityModule.HuntTreasure.CurrentPos
	response.HuntTurn = player.ActivityModule.HuntTreasure.HuntTurns
	response.RetCode = msg.RE_SUCCESS

}

//! 玩家查询巡回奖励领取情况
func Hand_QueryHuntTurnsAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryHuntTurn_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryHuntTurnsAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryHuntTurn_Ack
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

	response.AwardMask = int(player.ActivityModule.HuntTreasure.HuntAward)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家领取巡回奖励
func Hand_GetHuntTurnsAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetHuntTurnAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetHuntTurnsAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetHuntTurnAward_Ack
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

	//! 获取活动
	activityID := player.ActivityModule.HuntTreasure.ActivityID

	if G_GlobalVariables.IsActivityOpen(activityID) == false {
		gamelog.Error("Hand_GetHuntTurnsAward Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 检查参数合法性
	activityInfo := gamedata.GetActivityInfo(activityID)
	if req.ID > gamedata.GetHuntTreasureAwardCount(activityInfo.AwardType) || req.ID < 1 {
		gamelog.Error("Hand_GetHuntTurnsAward Error: %v Invalid param", player.playerid)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断该奖励是否已经领取
	if player.ActivityModule.HuntTreasure.HuntAward.Get(uint32(req.ID)) == true {
		gamelog.Error("Hand_GetHuntTurnsAward Error: %v Repeat get award %d", player.playerid, req.ID)
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 判断领奖条件
	award := gamedata.GetHuntTreasureAward(req.ID, activityInfo.AwardType)
	if player.ActivityModule.HuntTreasure.HuntTurns < award.NeedTurn {
		gamelog.Error("Hand_GetHuntTurnsAward Error: %v Turns not enough", player.playerid)
		response.RetCode = msg.RE_TURNS_NOT_ENOUGH
		return
	}

	//! 领取奖励
	awardLst := gamedata.GetItemsFromAwardID(award.Award)
	player.BagMoudle.AddAwardItems(awardLst)
	for _, v := range awardLst {
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum})
	}

	//! 改变标记
	player.ActivityModule.HuntTreasure.HuntAward.Set(uint32(req.ID))
	player.ActivityModule.HuntTreasure.DB_SaveHuntTurnsAwardMark()

	response.RetCode = msg.RE_SUCCESS
}

//! 查询巡回商店
func Hand_QueryHuntTreasureStore(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryHuntStore_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryHuntTreasureStore Error : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryHuntStore_Ack
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

	activityID := player.ActivityModule.HuntTreasure.ActivityID

	if G_GlobalVariables.IsActivityOpen(activityID) == false {
		gamelog.Error("Hand_QueryHuntTreasureStore Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	activityInfo := gamedata.GetActivityInfo(activityID)

	for _, v := range player.ActivityModule.HuntTreasure.StoreItemLst {
		itemInfo := gamedata.GetHuntTreasureStoreItem(v.ID, activityInfo.AwardType)

		var item msg.MSG_HuntStoreItem
		item.ID = v.ID
		item.ItemID = itemInfo.ItemID
		item.ItemNum = itemInfo.ItemNum
		item.MoneyID = itemInfo.MoneyID
		item.MoneyNum = itemInfo.MoneyNum
		item.Score = itemInfo.Score
		item.IsBuy = v.IsBuy
		response.GoodsLst = append(response.GoodsLst, item)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 购买巡回商店物品
func Hand_BuyHuntTreasureStroreItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyHuntStoreItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyHuntTreasureStroreItem Error : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_BuyHuntStoreItem_Ack
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

	activityID := player.ActivityModule.HuntTreasure.ActivityID

	if G_GlobalVariables.IsActivityOpen(activityID) == false {
		gamelog.Error("Hand_BuyHuntTreasureStroreItem Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 判断是否为活动时间
	isEnd, _ := G_GlobalVariables.IsActivityTime(player.ActivityModule.HuntTreasure.ActivityID)
	if isEnd == false {
		gamelog.Error("Hand_BuyHuntTreasureStroreItem Error: Activity is over")
		response.RetCode = msg.RE_ACTIVITY_IS_OVER
		return
	}

	activityInfo := gamedata.GetActivityInfo(activityID)

	//! 获取物品信息
	if req.ID <= 0 || req.ID > len(player.ActivityModule.HuntTreasure.StoreItemLst) {
		gamelog.Error("Hand_BuyHuntTreasureStroreItem Error: Invalid ID %d", req.ID)
	}

	item := player.ActivityModule.HuntTreasure.StoreItemLst[req.ID-1]
	itemInfo := gamedata.GetHuntTreasureStoreItem(item.ID, activityInfo.AwardType)
	if itemInfo == nil {
		gamelog.Error("Hand_BuyHuntTreasureStroreItem Error: Invalid ID %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	itemPos := req.ID - 1

	//! 检测货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(itemInfo.MoneyID, itemInfo.MoneyNum) == false {
		gamelog.Error("Hand_BuyHuntTreasureStroreItem Error: Player money not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 检测是否被购买
	if player.ActivityModule.HuntTreasure.StoreItemLst[itemPos].IsBuy == true {
		gamelog.Error("Hand_BuyHuntTreasureStoreItem Error: Item is aleady buy")
		response.RetCode = msg.RE_ALEADY_BUY
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(itemInfo.MoneyID, itemInfo.MoneyNum)

	//! 给予物品
	player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum)

	//! 改变标记
	player.ActivityModule.HuntTreasure.StoreItemLst[itemPos].IsBuy = true
	player.ActivityModule.HuntTreasure.DB_ChangeHuntStoreItemMark(itemPos)

	//! 增加积分
	indexToday := 0
	if utility.GetCurDayMod() == 1 {
		indexToday = 1
	}

	player.ActivityModule.HuntTreasure.Score += itemInfo.Score
	player.ActivityModule.HuntTreasure.TodayScore[indexToday] += itemInfo.Score
	player.ActivityModule.HuntTreasure.DB_SaveHuntScore()
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_GET_HUNT_SCORE, itemInfo.Score)

	player.ActivityModule.HuntTreasure.IsHaveStore = false
	player.ActivityModule.HuntTreasure.DB_SaveStoreMark()

	G_HuntTreasureTodayRanker.SetRankItem(player.playerid, player.ActivityModule.HuntTreasure.TodayScore[indexToday])

	G_HuntTreasureTotalRanker.SetRankItem(player.playerid, player.ActivityModule.HuntTreasure.Score)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

func Hand_CleanHuntStore(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_CleanHuntStore_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_CleanHuntStore Error : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_CleanHuntStore_Ack
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

	player.ActivityModule.HuntTreasure.StoreItemLst = []THuntStoreItem{}
	player.ActivityModule.HuntTreasure.DB_UpdateHuntStore()

	response.RetCode = msg.RE_SUCCESS
}

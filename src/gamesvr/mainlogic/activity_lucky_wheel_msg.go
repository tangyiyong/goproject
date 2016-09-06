package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"utility"
)

//! 玩家请求查询幸运轮盘
func Hand_QueryLuckyWheel(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryLuckyWheel_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryLuckyWheel Error : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryLuckyWheel_Ack
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
	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.LuckyWheel.ActivityID) == false {
		gamelog.Error("Hand_QueryLuckyWheel Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	_, response.EndCountDown = G_GlobalVariables.IsActivityTime(player.ActivityModule.LuckyWheel.ActivityID)

	indexToday := 0
	if utility.GetCurDayMod() == 1 {
		indexToday = 1
	}

	response.Score = player.ActivityModule.LuckyWheel.TodayScore[indexToday]
	response.NormalMoneyPoor = G_GlobalVariables.NormalMoneyPoor
	response.ExcitedMoneyPoor = G_GlobalVariables.ExcitedMoneyPoor

	for _, v := range player.ActivityModule.LuckyWheel.NormalAwardLst {
		itemInfo := gamedata.GetLuckyWheelItemFromID(v)
		if itemInfo == nil {
			gamelog.Error("Hand_QueryLuckyWheel Error: Can't find id: %d", v)
			return
		}

		var items msg.MSG_LuckyWheelAward
		items.ItemID = itemInfo.ItemID
		items.ItemNum = itemInfo.ItemNum
		items.IsSpecial = itemInfo.IsSpecial
		response.NormalAwardLst = append(response.NormalAwardLst, items)
	}

	for _, v := range player.ActivityModule.LuckyWheel.ExcitedAwardLst {
		itemInfo := gamedata.GetLuckyWheelItemFromID(v)
		if itemInfo == nil {
			gamelog.Error("Hand_QueryLuckyWheel Error: Can't find id: %d", v)
			return
		}

		var items msg.MSG_LuckyWheelAward
		items.ItemID = itemInfo.ItemID
		items.ItemNum = itemInfo.ItemNum
		items.IsSpecial = itemInfo.IsSpecial
		response.ExcitedAwardLst = append(response.ExcitedAwardLst, items)
	}

	response.CostMoneyID[0] = gamedata.NormalWheelMoneyID
	response.CostMoneyID[1] = gamedata.NormalWheelMoneyID
	response.CostMoneyID[2] = gamedata.ExcitedWheelMoneyID
	response.CostMoneyID[3] = gamedata.ExcitedWheelMoneyID

	response.CostMoneyNum[0] = gamedata.NormalWheelMoneyNum
	response.CostMoneyNum[1] = gamedata.NormalWheelMoneyNum * 10
	response.CostMoneyNum[2] = gamedata.ExcitedWheelMoneyNum
	response.CostMoneyNum[3] = gamedata.ExcitedWheelMoneyNum * 10

	response.TotalRank = -1
	response.TodayRank = -1
	for i, v := range G_LuckyWheelTotalRanker.List {
		if v.RankID == player.playerid {
			response.TotalRank = i + 1
			break
		}
	}

	for i, v := range G_LuckyWheelTodayRanker.List {
		if v.RankID == player.playerid {
			response.TodayRank = i + 1
			break
		}
	}

	response.NormalFreeTimes = player.ActivityModule.LuckyWheel.NormalFreeTimes

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家申请转动转盘
func Hand_RotatingWheel(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RotatingWheel_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RotatingWheel Error : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_RotatingWheel_Ack
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
	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.LuckyWheel.ActivityID) == false {
		gamelog.Error("Hand_RotatingWheel Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	indexToday := 0
	if utility.GetCurDayMod() == 1 {
		indexToday = 1
	}

	moneyNum := player.RoleMoudle.GetMoney(gamedata.ExcitedWheelMoneyID)

	//! 判断是否为豪华转盘
	if req.IsExcited == 1 {
		//! 普通转盘
		if req.IsStartTenTimes == 0 {
			//! 判断是否有免费次数
			if player.ActivityModule.LuckyWheel.NormalFreeTimes < 1 {
				//! 没有免费次数,判断道具
				if player.BagMoudle.IsItemEnough(gamedata.LuckyWheelCostItemID, 1) == false {
					//! 花费货币
					if player.RoleMoudle.CheckMoneyEnough(gamedata.NormalWheelMoneyID, gamedata.NormalWheelMoneyNum) == false {
						gamelog.Error("Hand_RotatingWheel Error: Money not enough Need: %d %d ", gamedata.NormalWheelMoneyID, gamedata.NormalWheelMoneyNum)
						response.RetCode = msg.RE_NOT_ENOUGH_MONEY
						return
					} else {
						response.CostItem = msg.MSG_ItemData{gamedata.NormalWheelMoneyID, gamedata.NormalWheelMoneyNum}

						//! 扣除金币
						player.RoleMoudle.CostMoney(gamedata.NormalWheelMoneyID, gamedata.NormalWheelMoneyNum)
					}
				} else {
					//! 扣除道具
					player.BagMoudle.RemoveNormalItem(gamedata.LuckyWheelCostItemID, 1)
					response.CostItem = msg.MSG_ItemData{gamedata.LuckyWheelCostItemID, 1}
				}

			} else {
				//! 直接扣除免费次数
				player.ActivityModule.LuckyWheel.NormalFreeTimes -= 1
				player.ActivityModule.LuckyWheel.DB_SaveLuckyWheelFreeTimes()
				response.CostFreeTimes = 1
			}
			//! 奖金池变化
			G_GlobalVariables.NormalMoneyPoor += 1
			G_GlobalVariables.DB_SaveMoneyPoor()

			player.ActivityModule.LuckyWheel.TodayScore[indexToday] += 10
			player.ActivityModule.LuckyWheel.TotalScore += 10
		} else {
			needItemNum := 10 - player.ActivityModule.LuckyWheel.NormalFreeTimes
			if player.BagMoudle.IsItemEnough(gamedata.LuckyWheelCostItemID, needItemNum) == false {
				//! 道具+免费次数不足十次, 则直接扣去钻石
				if player.RoleMoudle.CheckMoneyEnough(gamedata.NormalWheelMoneyID, gamedata.NormalWheelMoneyNum*10) == false {
					gamelog.Error("Hand_RotatingWheel Error: Money not enough Need: %d %d ", gamedata.NormalWheelMoneyID, gamedata.NormalWheelMoneyNum)
					response.RetCode = msg.RE_NOT_ENOUGH_MONEY
					return
				} else {
					//! 扣除金币
					response.CostItem = msg.MSG_ItemData{gamedata.NormalWheelMoneyID, gamedata.NormalWheelMoneyNum * 10}
					player.RoleMoudle.CostMoney(gamedata.NormalWheelMoneyID, gamedata.NormalWheelMoneyNum*10)
				}
			} else {
				//! 扣除免费次数
				response.CostFreeTimes = player.ActivityModule.LuckyWheel.NormalFreeTimes
				player.ActivityModule.LuckyWheel.NormalFreeTimes = 0
				player.ActivityModule.LuckyWheel.DB_SaveLuckyWheelFreeTimes()

				//! 扣除道具
				player.BagMoudle.RemoveNormalItem(gamedata.LuckyWheelCostItemID, needItemNum)
				response.CostItem = msg.MSG_ItemData{gamedata.LuckyWheelCostItemID, needItemNum}
			}

			//! 奖金池变化
			G_GlobalVariables.NormalMoneyPoor += 10
			G_GlobalVariables.DB_SaveMoneyPoor()

			player.ActivityModule.LuckyWheel.TodayScore[indexToday] += 100
			player.ActivityModule.LuckyWheel.TotalScore += 100
		}

		response.TodayRank = -1
		response.TotalRank = -1
		response.TodayRank = G_LuckyWheelTodayRanker.SetRankItem(player.playerid, player.ActivityModule.LuckyWheel.TodayScore[indexToday])

		response.TotalRank = G_LuckyWheelTotalRanker.SetRankItem(player.playerid, player.ActivityModule.LuckyWheel.TotalScore)

		player.ActivityModule.LuckyWheel.DB_SaveLuckyWheelScore()

	} else if req.IsExcited == 2 {
		//! 豪华转盘
		if req.IsStartTenTimes == 0 {
			//! 只允许使用钻石
			if player.RoleMoudle.CheckMoneyEnough(gamedata.ExcitedWheelMoneyID, gamedata.ExcitedWheelMoneyNum) == false {
				gamelog.Error("Hand_RotatingWheel Error: Money not enough  Need: %d %d", gamedata.ExcitedWheelMoneyID, gamedata.ExcitedWheelMoneyNum)
				response.RetCode = msg.RE_NOT_ENOUGH_MONEY
				return
			} else {
				//! 扣除金币
				player.RoleMoudle.CostMoney(gamedata.ExcitedWheelMoneyID, gamedata.ExcitedWheelMoneyNum)
				response.CostItem = msg.MSG_ItemData{gamedata.ExcitedWheelMoneyID, gamedata.ExcitedWheelMoneyNum}
			}

			//! 奖金池变化
			G_GlobalVariables.ExcitedMoneyPoor += 10
			G_GlobalVariables.DB_SaveMoneyPoor()

			//! 积分变化
			player.ActivityModule.LuckyWheel.TodayScore[indexToday] += 100
			player.ActivityModule.LuckyWheel.TotalScore += 100
		} else {
			//! 十连抽只允许使用钻石
			if player.RoleMoudle.CheckMoneyEnough(gamedata.ExcitedWheelMoneyID, gamedata.ExcitedWheelMoneyNum*10) == false {
				gamelog.Error("Hand_RotatingWheel Error: Money not enough Need: %d %d ", gamedata.ExcitedWheelMoneyNum*10, gamedata.ExcitedWheelMoneyNum)
				response.RetCode = msg.RE_NOT_ENOUGH_MONEY
				return
			} else {
				//! 扣除金币
				response.CostItem = msg.MSG_ItemData{gamedata.ExcitedWheelMoneyID, gamedata.ExcitedWheelMoneyNum * 10}

				player.RoleMoudle.CostMoney(gamedata.ExcitedWheelMoneyID, gamedata.ExcitedWheelMoneyNum*10)
			}

			//! 奖金池变化
			G_GlobalVariables.ExcitedMoneyPoor += 100
			G_GlobalVariables.DB_SaveMoneyPoor()

			player.ActivityModule.LuckyWheel.TodayScore[indexToday] += 1000
			player.ActivityModule.LuckyWheel.TotalScore += 1000
		}

		response.TodayRank = -1
		response.TotalRank = -1
		response.TodayRank = G_LuckyWheelTodayRanker.SetRankItem(player.playerid, player.ActivityModule.LuckyWheel.TodayScore[indexToday])

		response.TotalRank = G_LuckyWheelTotalRanker.SetRankItem(player.playerid, player.ActivityModule.LuckyWheel.TotalScore)

		player.ActivityModule.LuckyWheel.DB_SaveLuckyWheelScore()

	} else {
		gamelog.Error("Hand_RotatingWheel Error: Invalid param %d", req.IsExcited)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 随机奖励
	if req.IsStartTenTimes != 0 {
		for i := 0; i < 10; i++ {
			itemID, itemNum, isSpecial, index := player.ActivityModule.LuckyWheel.RandWheelAward(req.IsExcited)
			if isSpecial == 1 {
				//! 奖励奖金池的百分之50
				if req.IsExcited == 1 {
					awardMoney := G_GlobalVariables.NormalMoneyPoor * itemNum / 10000
					G_GlobalVariables.NormalMoneyPoor -= awardMoney
					if G_GlobalVariables.NormalMoneyPoor < 0 {
						G_GlobalVariables.NormalMoneyPoor = 0
					}
					player.RoleMoudle.AddMoney(itemID, awardMoney)
					response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{itemID, awardMoney})

				} else {
					awardMoney := G_GlobalVariables.ExcitedMoneyPoor * itemNum / 10000
					G_GlobalVariables.ExcitedMoneyPoor -= awardMoney
					if G_GlobalVariables.ExcitedMoneyPoor < 0 {
						G_GlobalVariables.ExcitedMoneyPoor = 0
					}

					player.RoleMoudle.AddMoney(itemID, awardMoney)
					response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{itemID, awardMoney})
				}

			} else {
				player.BagMoudle.AddAwardItem(itemID, itemNum)
				response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{itemID, itemNum})
			}
			if i == 9 {
				//! 十连转取最后一次转到的索引
				response.AwardIndex = index + 1
			}

		}

		G_GlobalVariables.DB_SaveMoneyPoor()
	} else {
		itemID, itemNum, isSpecial, index := player.ActivityModule.LuckyWheel.RandWheelAward(req.IsExcited)
		if isSpecial == 1 {
			//! 奖励奖金池的百分之50
			if req.IsExcited == 1 {
				awardMoney := G_GlobalVariables.NormalMoneyPoor * itemNum / 10000
				G_GlobalVariables.NormalMoneyPoor -= awardMoney
				if G_GlobalVariables.NormalMoneyPoor < 0 {
					G_GlobalVariables.NormalMoneyPoor = 0
				}
				player.RoleMoudle.AddMoney(itemID, awardMoney)
				G_GlobalVariables.DB_SaveMoneyPoor()
				response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{itemID, awardMoney})

			} else {
				awardMoney := G_GlobalVariables.ExcitedMoneyPoor * itemNum / 10000
				G_GlobalVariables.ExcitedMoneyPoor -= awardMoney
				if G_GlobalVariables.ExcitedMoneyPoor < 0 {
					G_GlobalVariables.ExcitedMoneyPoor = 0
				}
				player.RoleMoudle.AddMoney(itemID, awardMoney)
				G_GlobalVariables.DB_SaveMoneyPoor()
				response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{itemID, itemNum})
			}

		} else {
			player.BagMoudle.AddAwardItem(itemID, itemNum)
			response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{itemID, itemNum})
		}

		response.AwardIndex = index + 1
	}

	response.MoneyNum = moneyNum - player.RoleMoudle.GetMoney(gamedata.ExcitedWheelMoneyID)
	response.NormalMoneyPoor = G_GlobalVariables.NormalMoneyPoor
	response.ExcitedMoneyPoor = G_GlobalVariables.ExcitedMoneyPoor
	response.Score = player.ActivityModule.LuckyWheel.TodayScore[indexToday]
	response.RetCode = msg.RE_SUCCESS
}

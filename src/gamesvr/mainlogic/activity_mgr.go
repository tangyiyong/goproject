package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

//! 活动刷新
func ActivityTimerFunc(now int32) bool {
	gamelog.Info("Timer: ActivityTimerFunc")
	//! 检测新启活动
	for i := 0; i < len(G_GlobalVariables.ActivityLst); i++ {
		pActivity := gamedata.GetActivityInfo(G_GlobalVariables.ActivityLst[i].ActivityID)
		if pActivity == nil {
			gamelog.Error("ActivityTimerFunc Error: Invalid ActivityID %d", G_GlobalVariables.ActivityLst[i].ActivityID)
			return false
		}

		if now < G_GlobalVariables.ActivityLst[i].endTime || pActivity.CycleType == gamedata.CyCle_All {
			//! 没有到结束时间或者永久存在并正在运行的活动执行刷新
			RefreshActivity(G_GlobalVariables.ActivityLst[i].activityType, G_GlobalVariables.ActivityLst[i].ActivityID)
			G_GlobalVariables.ActivityLst[i].VersionCode += 1
		} else if now >= G_GlobalVariables.ActivityLst[i].endTime {
			//! 已经超过结束时间的非永久并正在运行的活动执行结束
			EndActivity(G_GlobalVariables.ActivityLst[i].activityType, G_GlobalVariables.ActivityLst[i].ActivityID)
			G_GlobalVariables.ActivityLst[i].ResetCode += 1
			G_GlobalVariables.ActivityLst[i].VersionCode = 0
			G_GlobalVariables.ActivityLst[i].beginTime, G_GlobalVariables.ActivityLst[i].actEndTime, G_GlobalVariables.ActivityLst[i].endTime = CalcActivityTime(G_GlobalVariables.ActivityLst[i].ActivityID, GetOpenServerDay())
		}
	}

	G_GlobalVariables.DB_UpdateActivityLst()

	return true
}

func EndActivity(activityType int, activityID int32) bool {
	if activityType == gamedata.Activity_Hunt_Treasure {
		G_HuntTreasureTodayRanker.Clear()
		G_HuntTreasureTotalRanker.Clear()
		G_HuntTreasureYesterdayRanker.Clear()
		mongodb.UpdateToDBAll("PlayerActivity", nil, &bson.M{"$set": bson.M{
			"hunttreasure.todayscore.0": 0,
			"hunttreasure.todayscore.1": 0,
			"hunttreasure.score":        0}})
	} else if activityType == gamedata.Activity_Card_Master {
		G_CardMasterTodayRanker.Clear()
		G_CardMasterTotalRanker.Clear()
		G_CardMasterYesterdayRanker.Clear()
		mongodb.UpdateToDBAll("PlayerActivity", nil, &bson.M{"$set": bson.M{
			"cardmaster.jifen.0":    0,
			"cardmaster.jifen.1":    0,
			"cardmaster.totaljifen": 0}})
	} else if activityType == gamedata.Activity_MoonlightShop {
	} else if activityType == gamedata.Activity_Luckly_Wheel {
		G_LuckyWheelTodayRanker.Clear()
		G_LuckyWheelTotalRanker.Clear()
		G_LuckyWheelYesterdayRanker.Clear()
		mongodb.UpdateToDBAll("PlayerActivity", nil, &bson.M{"$set": bson.M{
			"luckywheel.todayscore.0": 0,
			"luckywheel.todayscore.1": 0,
			"luckywheel.totalscore":   0}})
	} else if activityType == gamedata.Activity_Beach_Baby {
		G_BeachBabyTotalRanker.Clear()
		G_BeachBabyYesterdayRanker.Clear()
		mongodb.UpdateToDBAll("PlayerActivity", nil, &bson.M{"$set": bson.M{
			"beachbaby.score.0":    0,
			"beachbaby.score.1":    0,
			"beachbaby.totalscore": 0}})
	} else if activityType == gamedata.Activity_Seven { //! 七天
		G_GlobalVariables.DB_CleanSevenDayInfo(activityID)
	} else if activityType == gamedata.Activity_Competition { //! 战力排行
		//! 开始统计全服战力值排名
		activityInfo := gamedata.GetActivityInfo(activityID)

		//! 获取排名发放奖励
		for i, v := range G_FightRanker.List {
			if v.RankID == 0 {
				break
			}

			award := gamedata.GetCompetitionAward(activityInfo.AwardType, i+1)
			if award == 0 {
				break
			}

			var awardData TAwardData
			awardData.TextType = Text_CompetitionRankAward
			awardData.ItemLst = gamedata.GetItemsFromAwardID(award)
			awardData.Time = utility.GetCurTime()
			value := fmt.Sprintf("%d", i+1)
			awardData.Value = []string{value}

			SendAwardToPlayer(v.RankID, &awardData)
		}
	} else if activityType == gamedata.Activity_Group_Purchase { //! 团购
		activityLst := []TActivityModule{}
		s := mongodb.GetDBSession()
		defer s.Close()

		//! 获取所有参与过团购玩家信息
		err := s.DB(appconfig.GameDbName).C("PlayerActivity").Find(bson.M{"grouppurchase.score": bson.M{"$gt": 0}}).All(&activityLst)
		if err != nil {
			gamelog.Error("PlayerActivity Load Error :%s", err.Error())
			return false
		}

		for i := 0; i < len(activityLst); i++ {
			activity := activityLst[i]

			awardType := G_GlobalVariables.GetActivityAwardType(activity.GroupPurchase.ActivityID)

			//! 计算差价
			diffPrice := 0
			var awardData TAwardData
			for i := 0; i < len(activity.GroupPurchase.PurchaseCostLst); i++ {
				costItemID := activity.GroupPurchase.PurchaseCostLst[i].ItemID
				costMoney := 0
				costTimes := 0
				for _, v := range activity.GroupPurchase.PurchaseCostLst {
					if v.ItemID == costItemID {
						costMoney += v.MoneyNum
						costTimes += v.Times
					}
				}

				//! 获取现价
				saleInfo, _ := G_GlobalVariables.GetGroupPurchaseItemInfo(costItemID)
				salePriceInfo := gamedata.GetGroupPurchaseItemInfoFromSale(costItemID, awardType, saleInfo.SaleNum)

				//! 获取差价
				diffPrice = costMoney - costTimes*salePriceInfo.MoneyNum

				awardData.TextType = Text_Group_Purchase
				awardData.Time = utility.GetCurTime()
				awardData.ItemLst = append(awardData.ItemLst, gamedata.ST_ItemData{1, diffPrice})
				SendAwardToPlayer(activity.PlayerID, &awardData)
			}

		}

		//! 清除全局团购记录
		G_GlobalVariables.DB_CleanGroupPurchase()
	}

	return true
}

//! 刷新运营活动 每天整点刷新一次
func RefreshActivity(activityType int, activityID int32) bool {
	gamelog.Info("Timer: RefreshActivity")
	if activityType == gamedata.Activity_Moon_Card {
	} else if activityType == gamedata.Activity_Singel_Recharge {
	} else if activityType == gamedata.Activity_Limit_Daily_Task {
	} else if activityType == gamedata.Activity_Hunt_Treasure {
		if len(G_HuntTreasureYesterdayRanker.List) != 0 {
			//! 今日排行榜计入昨天
			G_HuntTreasureYesterdayRanker.Clear()
			G_HuntTreasureYesterdayRanker.CopyFrom(&G_HuntTreasureTodayRanker)
			G_HuntTreasureTodayRanker.Clear()

			//! 清除玩家今日积分
			FieldName := fmt.Sprintf("hunttreasure.todayscore.%d", utility.GetCurDayMod())
			mongodb.UpdateToDBAll("PlayerActivity", nil, &bson.M{"$set": bson.M{FieldName: 0}})
		}
	} else if activityType == gamedata.Activity_Card_Master {
		if len(G_CardMasterYesterdayRanker.List) != 0 {
			//! 今日排行榜计入昨天
			G_CardMasterYesterdayRanker.Clear()
			G_CardMasterYesterdayRanker.CopyFrom(&G_CardMasterTodayRanker)
			G_CardMasterTodayRanker.Clear()

			//! 清除玩家今日积分
			FieldName := fmt.Sprintf("cardmaster.jifen.%d", utility.GetCurDayMod())
			mongodb.UpdateToDBAll("PlayerActivity", nil, &bson.M{"$set": bson.M{FieldName: 0}})
		}
	} else if activityType == gamedata.Activity_MoonlightShop {
	} else if activityType == gamedata.Activity_Luckly_Wheel {
		//! 今日排行榜计入昨天
		if len(G_LuckyWheelYesterdayRanker.List) != 0 {
			G_LuckyWheelYesterdayRanker.Clear()
			G_LuckyWheelYesterdayRanker.CopyFrom(&G_LuckyWheelTodayRanker)
			G_LuckyWheelTodayRanker.Clear()

			//! 清除玩家今日积分
			FieldName := fmt.Sprintf("luckywheel.todayscore.%d", utility.GetCurDayMod())
			mongodb.UpdateToDBAll("PlayerActivity", nil, &bson.M{"$set": bson.M{FieldName: 0}})
		}
	} else if activityType == gamedata.Activity_Beach_Baby {
		//! 今日排行榜计入昨天
		G_BeachBabyYesterdayRanker.Clear()
		G_BeachBabyYesterdayRanker.CopyFrom(&G_BeachBabyTodayRanker)
		G_BeachBabyTodayRanker.Clear()

		//! 清除玩家今日积分
		FieldName := fmt.Sprintf("beachbaby.score.%d", utility.GetCurDayMod())
		mongodb.UpdateToDBAll("PlayerActivity", nil, &bson.M{"$set": bson.M{FieldName: 0}})
	} else if activityType == gamedata.Activity_Group_Purchase {
	} else if activityType == gamedata.Activity_Festival {
	}

	return true
}

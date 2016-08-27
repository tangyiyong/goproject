package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"time"
	"utility"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//! 活动刷新
func ActivityTimerFunc(now int64) bool {
	gamelog.Info("Timer: ActivityTimerFunc")
	//! 检测新启活动
	for i := 0; i < len(G_GlobalVariables.ActivityLst); i++ {
		if now >= G_GlobalVariables.ActivityLst[i].endTime && G_GlobalVariables.ActivityLst[i].endTime > 0 {
			//! 已经超过结束时间的非永久并正在运行的活动执行结束
			EndActivity(G_GlobalVariables.ActivityLst[i].activityType, G_GlobalVariables.ActivityLst[i].ActivityID)
			G_GlobalVariables.ActivityLst[i].ResetCode += 1
			G_GlobalVariables.ActivityLst[i].VersionCode = 0

			//! 计算下一次开始时间
			beginTime, endTime := gamedata.GetActivityNextBeginTime(G_GlobalVariables.ActivityLst[i].ActivityID, GetOpenServerDay())
			G_GlobalVariables.ActivityLst[i].beginTime, G_GlobalVariables.ActivityLst[i].endTime = beginTime, endTime
		} else if now < G_GlobalVariables.ActivityLst[i].endTime ||
			G_GlobalVariables.ActivityLst[i].endTime == 0 {
			//! 没有到结束时间或者永久存在并正在运行的活动执行刷新
			RefreshActivity(G_GlobalVariables.ActivityLst[i].activityType, G_GlobalVariables.ActivityLst[i].ActivityID)
			G_GlobalVariables.ActivityLst[i].VersionCode += 1
		}
	}

	go G_GlobalVariables.DB_UpdateActivityLst()

	return true
}

func CheckActivityAdd() {
	//! 获取今日开启活动
	openDay := GetOpenServerDay()
	for _, v := range gamedata.GT_ActivityLst {
		if v.ID == 0 {
			gamelog.Error("CheckActivityAdd Error: Invalid ActivityID:%d", v.ID)
			continue
		}

		isExist := false
		for _, n := range G_GlobalVariables.ActivityLst {
			if n.ActivityID == v.ID {
				isExist = true
				break
			}
		}

		if isExist == true {
			continue
		}

		if v.ActivityType == gamedata.Activity_Seven {
			seven := TSevenDayBuyInfo{}
			seven.ActivityID = v.ID
			G_GlobalVariables.SevenDayLimit = append(G_GlobalVariables.SevenDayLimit, seven)
			G_GlobalVariables.DB_AddSevenDayBuyInfo(seven)
		}

		var activity TActivityData
		activity.ActivityID = v.ID
		activity.activityType = v.ActivityType
		activity.award = v.AwardType
		activity.beginTime, activity.endTime = gamedata.GetActivityEndTime(v.ID, openDay)
		activity.VersionCode = 0
		activity.Status = v.Status
		activity.ResetCode = 0
		G_GlobalVariables.ActivityLst = append(G_GlobalVariables.ActivityLst, activity)
		G_GlobalVariables.DB_AddNewActivity(activity)
	}
}

func EndActivity(activityType int, activityID int) bool {
	if activityType == gamedata.Activity_Hunt_Treasure {
		G_HuntTreasureTodayRanker.Clear()
		G_HuntTreasureTotalRanker.Clear()
		G_HuntTreasureYesterdayRanker.Clear()
		mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerActivity", nil, bson.M{"$set": bson.M{
			"hunttreasure.todayscore.0": 0,
			"hunttreasure.todayscore.1": 0,
			"hunttreasure.score":        0}})
	} else if activityType == gamedata.Activity_Card_Master {
		G_CardMasterTodayRanker.Clear()
		G_CardMasterTotalRanker.Clear()
		G_CardMasterYesterdayRanker.Clear()
		mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerActivity", nil, bson.M{"$set": bson.M{
			"cardmaster.jifen.0":    0,
			"cardmaster.jifen.1":    0,
			"cardmaster.totaljifen": 0}})
	} else if activityType == gamedata.Activity_MoonlightShop {
	} else if activityType == gamedata.Activity_Luckly_Wheel {
		G_LuckyWheelTodayRanker.Clear()
		G_LuckyWheelTotalRanker.Clear()
		G_LuckyWheelYesterdayRanker.Clear()
		mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerActivity", nil, bson.M{"$set": bson.M{
			"luckywheel.todayscore.0": 0,
			"luckywheel.todayscore.1": 0,
			"luckywheel.totalscore":   0}})
	} else if activityType == gamedata.Activity_Beach_Baby {
		G_BeachBabyTotalRanker.Clear()
		G_BeachBabyYesterdayRanker.Clear()
		mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerActivity", nil, bson.M{"$set": bson.M{
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
			awardData.Time = time.Now().Unix()
			value := fmt.Sprintf("%d", i+1)
			awardData.Value = []string{value}

			SendAwardToPlayer(v.RankID, &awardData)
		}
	} else if activityType == gamedata.Activity_Group_Purchase { //! 团购
		//! 清除全局团购记录
		G_GlobalVariables.DB_CleanGroupPurchase()
	}

	return true
}

//! 刷新运营活动 每天整点刷新一次
func RefreshActivity(activityType int, activityID int) bool {
	//gamelog.Info("Timer: RefreshActivity")
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
			mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerActivity", nil, bson.M{"$set": bson.M{FieldName: 0}})
		}
	} else if activityType == gamedata.Activity_Card_Master {
		if len(G_CardMasterYesterdayRanker.List) != 0 {
			//! 今日排行榜计入昨天
			G_CardMasterYesterdayRanker.Clear()
			G_CardMasterYesterdayRanker.CopyFrom(&G_CardMasterTodayRanker)
			G_CardMasterTodayRanker.Clear()

			//! 清除玩家今日积分
			FieldName := fmt.Sprintf("cardmaster.jifen.%d", utility.GetCurDayMod())
			mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerActivity", nil, bson.M{"$set": bson.M{FieldName: 0}})
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
			mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerActivity", nil, bson.M{"$set": bson.M{FieldName: 0}})
		}
	} else if activityType == gamedata.Activity_Beach_Baby {
		//! 今日排行榜计入昨天
		G_BeachBabyYesterdayRanker.Clear()
		G_BeachBabyYesterdayRanker.CopyFrom(&G_BeachBabyTodayRanker)
		G_BeachBabyTodayRanker.Clear()

		//! 清除玩家今日积分
		FieldName := fmt.Sprintf("beachbaby.score.%d", utility.GetCurDayMod())
		mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerActivity", nil, bson.M{"$set": bson.M{FieldName: 0}})
	} else if activityType == gamedata.Activity_Group_Purchase {
	} else if activityType == gamedata.Activity_Festival {
	}

	return true
}

func ClearPlayerActData(dataName string, handler func(*TPlayer)) {
	//! 清除玩家积分
	rankLst := []TActivityModule{}
	db_session := mongodb.GetDBSession()
	defer db_session.Close()
	collection := db_session.DB(appconfig.GameDbName).C("PlayerActivity")
	err := collection.Find(bson.M{dataName: bson.M{"$gt": 0}}).All(&rankLst)
	if err != mgo.ErrNotFound && err != nil {
		gamelog.Error3("Find_Sort error: %v \r\ndbName: %s \r\ntable: %s \r\nfind: %s \r\n",
			err.Error(), appconfig.GameDbName, "PlayerActivity", dataName)
		return
	}

	for _, v := range rankLst {
		//! 刷新玩家数据
		player := GetPlayerByID(v.PlayerID)
		if player != nil {
			handler(player)
		}
	}
}

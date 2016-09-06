package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"utility"
)

//! 获取今日活动
func Hand_GetActivity(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetActivity_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetActivity Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetActivity_Ack
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

	response.RetCode = msg.RE_SUCCESS

	for _, v := range G_GlobalVariables.ActivityLst {
		var activity msg.MSG_ActivityInfo
		if G_GlobalVariables.IsActivityOpen(v.ActivityID) == false {
			continue
		}

		activity.ID = v.ActivityID
		activityInfo := gamedata.GetActivityInfo(v.ActivityID)
		openday := GetOpenServerDay()
		if activityInfo == nil {
			gamelog.Error("GetActivityType Fail: type: %d", v.activityType)
			continue
		}

		activity.Icon = activityInfo.Icon
		activity.Type = v.activityType
		activity.Name = activityInfo.Name
		activity.AwardType = v.award
		activity.BeginTime = v.beginTime
		activity.EndTime = v.endTime
		activity.AwardTime = int(v.endTime)
		activity.IsInside = activityInfo.Inside

		pActivity, ok := player.ActivityModule.activityPtrs[v.ActivityID]
		if ok && pActivity != nil {
			activity.RedTip = pActivity.RedTip()
		}

		if (player.ActivityModule.FirstRecharge.FirstRechargeAward == 2 &&
			activity.Type == gamedata.Activity_First_Recharge && activity.ID == 1) ||
			(player.ActivityModule.FirstRecharge.NextRechargeAward == 2 &&
				activity.Type == gamedata.Activity_First_Recharge && activity.ID != 1) {
			//! 已领取首充/次充则不放入列表
			continue
		}

		//! 次充必须在首充结束后出现
		if player.ActivityModule.FirstRecharge.NextRechargeAward != 2 && activity.Type == gamedata.Activity_First_Recharge && activity.ID != 1 {
			if player.ActivityModule.FirstRecharge.FirstRechargeAward != 2 {
				continue
			}
		}

		if activityInfo.Inside == 3 && openday > activityInfo.Days {
			response.RemoveActivityIcon = append(response.RemoveActivityIcon, v.ActivityID)
			response.ActivityLst = append(response.ActivityLst, activity)
			continue
		}

		response.ActivityLst = append(response.ActivityLst, activity)
	}

}

func HuntTreasureRankAward(awardType int, player *TPlayer, indexToday int) (awardItem []gamedata.ST_ItemData, retCode int) {
	activityID := player.ActivityModule.HuntTreasure.ActivityID
	activityInfo := gamedata.GetActivityInfo(activityID)
	info := &player.ActivityModule.HuntTreasure
	if awardType == 1 { //! 昨日榜
		//! 判断标记
		if player.ActivityModule.HuntTreasure.IsRecvTodayRankAward == true {
			gamelog.Error("Hand_GetHuntRankAward Error: TodayRankAward Aleady received")
			retCode = msg.RE_ALREADY_RECEIVED
			return awardItem, retCode
		}

		//! 获取名次
		yesterdayRank := G_HuntTreasureYesterdayRanker.GetRankIndex(player.playerid, info.TodayScore[1-indexToday])
		if yesterdayRank <= 0 {
			retCode = msg.RE_NOT_ENOUGH_RANK
			return awardItem, retCode
		}

		rankAward := gamedata.GetOperationalRankAwardFromRank(gamedata.Activity_Hunt_Treasure, activityInfo.AwardType, yesterdayRank)

		if rankAward.TodayNormalRankAward != 0 {
			awardItem = gamedata.GetItemsFromAwardID(rankAward.TodayNormalRankAward)
		}

		if info.TodayScore[1-indexToday] >= gamedata.EliteHuntRankNeedScore && rankAward.TodayEliteRankAward != 0 {
			awardItem = append(awardItem, gamedata.GetItemsFromAwardID(rankAward.TodayEliteRankAward)...)
		}

		player.ActivityModule.HuntTreasure.IsRecvTodayRankAward = true
		player.ActivityModule.HuntTreasure.DB_UpdateHuntTodayRankAward()

	} else if awardType == 2 { //! 总榜
		//! 判断时间
		isActivityTime, _ := G_GlobalVariables.IsActivityTime(activityID)
		if G_GlobalVariables.IsActivityOpen(activityID) == false || isActivityTime == true {
			retCode = msg.RE_ACTIVITY_NOT_OVER
			return awardItem, retCode
		}

		//! 判断标记
		if player.ActivityModule.HuntTreasure.IsRecvTotalRankAward == true {
			gamelog.Error("Hand_GetHuntRankAward Error: TotalRankAward Aleady received")
			retCode = msg.RE_ALREADY_RECEIVED
			return awardItem, retCode
		}

		//! 获取名次
		totalRank := G_HuntTreasureTotalRanker.GetRankIndex(player.playerid, info.Score)
		if totalRank <= 0 {
			retCode = msg.RE_NOT_ENOUGH_RANK
			gamelog.Error("Hand_GetHuntRankAward Error: Can't receive award rank: %d", totalRank)
			return awardItem, retCode
		}

		rankAward := gamedata.GetOperationalRankAwardFromRank(gamedata.Activity_Hunt_Treasure, activityInfo.AwardType, totalRank)
		if rankAward.TotalNormalRankAward != 0 {
			awardItem = gamedata.GetItemsFromAwardID(rankAward.TotalNormalRankAward)
		}

		if info.Score >= gamedata.EliteHuntRankNeedScore && rankAward.TotalEliteRankAward != 0 {
			awardItem = append(awardItem, gamedata.GetItemsFromAwardID(rankAward.TotalEliteRankAward)...)
		}

		player.ActivityModule.HuntTreasure.IsRecvTotalRankAward = true
		player.ActivityModule.HuntTreasure.DB_UpdateHuntTotalRankAward()
	} else {
		gamelog.Error("Hand_GetActivityRankAward Error: Invalid Param AwardType: %d", awardType)
		retCode = msg.RE_INVALID_PARAM
		return awardItem, retCode
	}

	retCode = msg.RE_SUCCESS
	return awardItem, retCode
}

func LuckyWheelRankAward(awardType int, player *TPlayer, indexToday int) (awardItem []gamedata.ST_ItemData, retCode int) {
	activityID := player.ActivityModule.LuckyWheel.ActivityID
	activityInfo := gamedata.GetActivityInfo(activityID)
	info := &player.ActivityModule.LuckyWheel
	if awardType == 1 { //! 昨日榜
		//! 判断标记
		if player.ActivityModule.LuckyWheel.IsRecvTodayRankAward == true {
			gamelog.Error("Hand_GetHuntRankAward Error: TodayRankAward Aleady received")
			retCode = msg.RE_ALREADY_RECEIVED
			return awardItem, retCode
		}

		//! 获取名次
		yesterdayRank := G_LuckyWheelYesterdayRanker.GetRankIndex(player.playerid, info.TodayScore[1-indexToday])
		if yesterdayRank <= 0 {
			retCode = msg.RE_NOT_ENOUGH_RANK
			return awardItem, retCode
		}

		rankAward := gamedata.GetOperationalRankAwardFromRank(gamedata.Activity_Hunt_Treasure, activityInfo.AwardType, yesterdayRank)

		if rankAward.TodayNormalRankAward != 0 {
			awardItem = gamedata.GetItemsFromAwardID(rankAward.TodayNormalRankAward)
		}

		if info.TodayScore[1-indexToday] >= gamedata.EliteHuntRankNeedScore && rankAward.TodayEliteRankAward != 0 {
			awardItem = append(awardItem, gamedata.GetItemsFromAwardID(rankAward.TodayEliteRankAward)...)
		}

		player.ActivityModule.LuckyWheel.IsRecvTodayRankAward = true
		player.ActivityModule.LuckyWheel.DB_UpdateWheelTodayRankAward()

	} else if awardType == 2 { //! 总榜
		//! 判断时间
		isActivityTime, _ := G_GlobalVariables.IsActivityTime(activityID)
		if G_GlobalVariables.IsActivityOpen(activityID) == false || isActivityTime == false {
			retCode = msg.RE_ACTIVITY_NOT_OVER
			return awardItem, retCode
		}

		//! 判断标记
		if player.ActivityModule.LuckyWheel.IsRecvTotalRankAward == true {
			gamelog.Error("Hand_GetHuntRankAward Error: TotalRankAward Aleady received")
			retCode = msg.RE_ALREADY_RECEIVED
			return awardItem, retCode
		}

		//! 获取名次
		totalRank := G_LuckyWheelTotalRanker.GetRankIndex(player.playerid, info.TotalScore)
		if totalRank <= 0 {
			retCode = msg.RE_NOT_ENOUGH_RANK
			return awardItem, retCode
		}

		rankAward := gamedata.GetOperationalRankAwardFromRank(gamedata.Activity_Hunt_Treasure, activityInfo.AwardType, totalRank)
		if rankAward.TotalNormalRankAward != 0 {
			awardItem = gamedata.GetItemsFromAwardID(rankAward.TotalNormalRankAward)
		}

		if info.TotalScore >= gamedata.EliteHuntRankNeedScore && rankAward.TotalEliteRankAward != 0 {
			awardItem = append(awardItem, gamedata.GetItemsFromAwardID(rankAward.TotalEliteRankAward)...)
		}

		player.ActivityModule.LuckyWheel.IsRecvTotalRankAward = true
		player.ActivityModule.LuckyWheel.DB_UpdateWheelTotalRankAward()
	} else {
		gamelog.Error("Hand_GetActivityRankAward Error: Invalid Param AwardType: %d", awardType)
		retCode = msg.RE_INVALID_PARAM
		return awardItem, retCode
	}

	retCode = msg.RE_SUCCESS
	return awardItem, retCode
}

//! 获取排行榜奖励
func Hand_GetActivityRankAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetActivityRankAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetActivityRankAward Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetActivityRankAward_Ack
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

	awardItem := []gamedata.ST_ItemData{}
	indexToday := 0
	if utility.GetCurDayMod() == 1 {
		indexToday = 1
	}

	switch req.ActivityType {
	case gamedata.Activity_Hunt_Treasure: //! 巡回探宝
		{
			awardItem, response.RetCode = HuntTreasureRankAward(req.AwardType, player, indexToday)
		}
	case gamedata.Activity_Luckly_Wheel: //! 幸运轮盘
		{
			awardItem, response.RetCode = LuckyWheelRankAward(req.AwardType, player, indexToday)
		}
	case gamedata.Activity_Card_Master: //! 卡牌大师
		{
			awardItem, response.RetCode = CardMasterRankAward(req.AwardType, player, indexToday)
		}
	case gamedata.Activity_Beach_Baby: //! 沙滩宝贝
		{
			awardItem, response.RetCode = BeachBabyRankAward(req.AwardType, player, indexToday)
		}
	}

	if response.RetCode != msg.RE_SUCCESS {
		gamelog.Error("Hand_GetActivityRankAward Error: ErrorCode %d", response.RetCode)
		return
	}

	//! 奖励物品
	player.BagMoudle.AddAwardItems(awardItem)

	for _, v := range awardItem {
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum})
	}

	response.RetCode = msg.RE_SUCCESS

}

//! 获取活动排行榜
func Hand_GetActivityRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetActivityRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetActivityRank Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetActivityRank_Ack
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

	switch req.Type {
	case gamedata.Activity_Hunt_Treasure: //! 巡回探宝
		{
			module := &player.ActivityModule.HuntTreasure
			response.EliteScore = gamedata.EliteHuntRankNeedScore
			response.IsRecvTodayRankAward = module.IsRecvTodayRankAward
			response.IsRecvTotalRankAward = module.IsRecvTotalRankAward

			MakeMsgRankInfo(req.PlayerID, &response, &G_HuntTreasureTodayRanker, &G_HuntTreasureYesterdayRanker, &G_HuntTreasureTotalRanker, module)
		}
	case gamedata.Activity_Luckly_Wheel: //! 幸运转盘
		{
			module := &player.ActivityModule.LuckyWheel
			response.EliteScore = gamedata.EliteHuntRankNeedScore
			response.IsRecvTodayRankAward = module.IsRecvTodayRankAward
			response.IsRecvTotalRankAward = module.IsRecvTotalRankAward

			MakeMsgRankInfo(req.PlayerID, &response, &G_LuckyWheelTodayRanker, &G_LuckyWheelYesterdayRanker, &G_LuckyWheelTotalRanker, module)
		}
	case gamedata.Activity_Card_Master: //! 卡牌大师
		{
			module := &player.ActivityModule.CardMaster
			response.EliteScore = gamedata.CardMaster_TotalRank_Limit
			response.IsRecvTodayRankAward = module.IsGetTodayRankAward
			response.IsRecvTotalRankAward = module.IsGetTotalRankAward

			MakeMsgRankInfo(req.PlayerID, &response, &G_CardMasterTodayRanker, &G_CardMasterYesterdayRanker, &G_CardMasterTotalRanker, module)
		}
	case gamedata.Activity_Beach_Baby: //! 沙滩宝贝
		{
			module := &player.ActivityModule.BeachBaby
			response.EliteScore = gamedata.BeachBaby_TotalRank_Limit
			response.IsRecvTodayRankAward = module.IsGetTodayRankAward
			response.IsRecvTotalRankAward = module.IsGetTotalRankAward

			MakeMsgRankInfo(req.PlayerID, &response, &G_BeachBabyTodayRanker, &G_BeachBabyYesterdayRanker, &G_BeachBabyTotalRanker, module)
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

type ActivityScoreFunc interface {
	GetTodayScore() int
	GetYesterdayScore() int
	GetTotalScore() int
}

func MakeMsgRankInfo(playerID int32, response *msg.MSG_GetActivityRank_Ack, todayRanker, yesterdayRanker, totalRanker *utility.TRanker, module ActivityScoreFunc) {
	response.TodayRankLst = GetActivityRankList(todayRanker)
	response.TotalRankLst = GetActivityRankList(totalRanker)

	response.ScoreLst[0] = module.GetYesterdayScore()
	response.ScoreLst[1] = module.GetTodayScore()
	response.ScoreLst[2] = module.GetTotalScore()

	response.RankLst[0] = yesterdayRanker.GetRankIndex(playerID, response.ScoreLst[0])
	response.RankLst[1] = todayRanker.GetRankIndex(playerID, response.ScoreLst[1])
	response.RankLst[2] = totalRanker.GetRankIndex(playerID, response.ScoreLst[2])

	response.YesterdayRankLst = GetActivityRankList(yesterdayRanker)
}
func GetActivityRankList(ranker *utility.TRanker) (ret []msg.MSG_OperationalActivityRank) {
	ranker.ForeachShow(
		func(rankID int32, rankVal int) {
			simpleInfo := G_SimpleMgr.GetSimpleInfoByID(rankID)
			if simpleInfo != nil {
				ret = append(ret, msg.MSG_OperationalActivityRank{
					PlayerID:   simpleInfo.PlayerID,
					PlayerName: simpleInfo.Name,
					HeroID:     simpleInfo.HeroID,
					Quality:    simpleInfo.Quality,
					Level:      simpleInfo.Level,
					Score:      rankVal})
			}
		})
	return ret
}

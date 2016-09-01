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

func CopyCheck(player *TPlayer, copyID int, chapter int, copyType int) (bool, int) {
	//! 判断称号加成是否过期
	player.TitleModule.CheckTitleDeadLine()

	//验证是否符合条件
	pCopyInfo := gamedata.GetCopyBaseInfo(copyID)
	if pCopyInfo == nil {
		gamelog.Error("CopyCheck error : Invalid copyid : %d", copyID)
		return false, msg.RE_INVALID_PARAM
	}

	isEnough := player.RoleMoudle.CheckActionEnough(pCopyInfo.ActionType, pCopyInfo.ActionValue)
	if isEnough == false { //! 体力不足挑战
		gamelog.Error("CopyCheck error : Not Enough Action, Type:%d, value :%d", pCopyInfo.ActionType, pCopyInfo.ActionValue)
		return false, msg.RE_STRENGTH_NOT_ENOUGH
	}

	//! 检测英雄背包是否超载
	if player.BagMoudle.IsHeroBagFull() == true {
		gamelog.Error("CopyCheck error : Hero bag is full")
		return false, msg.RE_HERO_BAG_OVERLOAD
	}

	if copyType == gamedata.COPY_TYPE_Main { //! 主线副本
		nextCopyID, _ := gamedata.GetNextCopy(player.CopyMoudle.Main.CurCopyID, player.CopyMoudle.Main.CurChapter, gamedata.COPY_TYPE_Main)
		if nextCopyID < copyID {
			gamelog.Error("CopyCheck error nextCopy: %d  copyID: %d", nextCopyID, copyID)
			return false, msg.RE_NEED_PASS_PRE_COPY
		}

		//! 检查挑战次数是否足够
		for _, v := range player.CopyMoudle.Main.CopyInfo {
			if v.CopyID == copyID {
				if v.BattleTimes >= pCopyInfo.MaxBattleTimes {
					return false, msg.RE_NOT_ENOUGH_TIMES
				}
			}
		}

	} else if copyType == gamedata.COPY_TYPE_Elite {
		//! 检测功能是否开启
		if gamedata.IsFuncOpen(gamedata.FUNC_ELITE_COPY, player.GetLevel(), player.GetVipLevel()) == false {
			gamelog.Error("CopyCheck error elite copy not open.")
			return false, msg.RE_FUNC_NOT_OPEN
		}

		nextCopyID, _ := gamedata.GetNextCopy(player.CopyMoudle.Elite.CurCopyID, player.CopyMoudle.Elite.CurChapter, gamedata.COPY_TYPE_Elite)
		if nextCopyID < copyID {
			gamelog.Error("CopyCheck error nextCopy: %d  copyID: %d", nextCopyID, copyID)
			return false, msg.RE_NEED_PASS_PRE_COPY
		}

		//! 检查挑战次数是否足够
		for _, v := range player.CopyMoudle.Elite.CopyInfo {
			if v.CopyID == copyID {
				if v.BattleTimes >= pCopyInfo.MaxBattleTimes {
					return false, msg.RE_NOT_ENOUGH_TIMES
				}
			}
		}

	} else if copyType == gamedata.COPY_TYPE_Daily { //! 日常副本
		//! 检测功能是否开启
		if gamedata.IsFuncOpen(gamedata.FUNC_DAILY_COPY, player.GetLevel(), player.GetVipLevel()) == false {
			gamelog.Error("CopyCheck error daily copy not open.")
			return false, msg.RE_FUNC_NOT_OPEN
		}

		//! 检查是否当前是否能够挑战
		dailyCopy := gamedata.GetDailyCopyData(pCopyInfo.CopyID)
		if dailyCopy == nil {
			gamelog.Error("GetDailyCopyData fail. copyID: %d", pCopyInfo.CopyID)
			return false, msg.RE_UNKNOWN_ERR
		}

		openResType := player.CopyMoudle.GetTodayDailyCopy()
		isCan := false
		for _, v := range openResType {
			if v == dailyCopy.ResType {
				isCan = true
				break
			}
		}

		if isCan == false {
			return false, msg.RE_COPY_IS_LOCK
		}

		//! 检查等级是否足够挑战
		level := player.GetLevel()
		if level < dailyCopy.Level {
			return false, msg.RE_COPY_IS_LOCK
		}

		//! 检查今日是否已经挑战
		for _, v := range player.CopyMoudle.Daily.CopyInfo {
			if v.ResID == dailyCopy.ResType && v.IsChallenge == true {
				gamelog.Error("CopyCheck Error: daily copy aleady challenge  res: %d  copyID: %d", dailyCopy.ResType, copyID)
				return false, msg.RE_CHALLENGE_ALEADY_END
			}
		}
	} else if copyType == gamedata.COPY_TYPE_Famous { //! 名将副本
		if gamedata.IsFuncOpen(gamedata.FUNC_FAMOUS_COPY, player.GetLevel(), player.GetVipLevel()) == false {
			gamelog.Error("CopyCheck error Famous copy not open.")
			return false, msg.RE_FUNC_NOT_OPEN
		}

		isSerialid := gamedata.IsSerialCopy(chapter, copyID)
		if isSerialid == false {
			//! 如果不是连环计,则检测挑战次数
			if player.CopyMoudle.Famous.BattleTimes > gamedata.FamousCopyChallengeTimes {
				return false, msg.RE_NOT_ENOUGH_TIMES
			}
		} else {
			//! 如果是连环计,则检测是否有通关记录
			if player.CopyMoudle.GetFamousCopyInfo(copyID, chapter) != nil {
				//! 已经通关后的连环计无法重复挑战
				gamelog.Error("Can not repeat the challenge series copy")
				return false, msg.RE_NOT_ENOUGH_TIMES
			}
		}

		//! 每个小关卡每天只能挑战一次
		famousCopy := player.CopyMoudle.GetFamousCopyInfo(copyID, chapter)
		if famousCopy != nil && famousCopy.BattleTimes >= 1 {
			return false, msg.RE_NOT_ENOUGH_TIMES
		}
	} else if copyType == gamedata.COPY_TYPE_Elite_Invade { //! 精英关卡入侵
		//! 检测该章节是否有入侵
		if player.CopyMoudle.IsHaveInvade(chapter) == false {
			gamelog.Error("CopyCheck error Not have Invade")
			return false, msg.RE_INVADE_ALEADY_ESCAPE
		}
	}
	return true, msg.RE_SUCCESS
}

//! 挑战条件检查
func Hand_BattleCheck(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_BattleCheck_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_BattleCheck : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_BattleCheck_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	ok, errcode := CopyCheck(player, req.CopyID, req.Chapter, req.CopyType)
	if ok != true {
		response.RetCode = errcode
		return
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 客户端汇报战斗结果
func Hand_BattleResult(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//以下为开启MD5消息验证后的代码
	/////		if false == utility.MsgDataCheck(buffer) {
	/////			//存在作弊的可能
	/////			gamelog.Error("Hand_BattleResult : Message Data Check Error!!!!")
	/////			return
	/////		}
	/////		var dLen = len(buffer) - 16
	/////		var req msg.MSG_BattleResult_Req
	/////		if json.Unmarshal(buffer[:dLen], &req) != nil {
	/////			gamelog.Error("Hand_BattleResult : Unmarshal error!!!!")
	/////			return
	/////		}
	///////////////////////////////////////
	var req msg.MSG_BattleResult_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_BattleResult : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_BattleResult_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.StarNum <= 0 || req.StarNum > 3 {
		gamelog.Error("Hand_BattleResult error : Invalid Star Num %d, playerid:%d", req.StarNum, req.PlayerID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	pCopyInfo := gamedata.GetCopyBaseInfo(req.CopyID)
	if pCopyInfo == nil {
		gamelog.Error("Hand_BattleResult error : Invalid copyid : %d", req.CopyID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	isEnough := player.RoleMoudle.CheckActionEnough(pCopyInfo.ActionType, pCopyInfo.ActionValue)
	if isEnough == false { //! 体力不足挑战
		gamelog.Error("Hand_BattleResult error : Not Enough Action")
		response.RetCode = msg.RE_STRENGTH_NOT_ENOUGH
		return
	}

	if req.CopyType == gamedata.COPY_TYPE_Main { //! 通关主线关卡
		ok, errcode := CopyCheck(player, req.CopyID, req.Chapter, gamedata.COPY_TYPE_Main)
		if ok != true {
			response.RetCode = errcode
			return
		}

		endCopyID := gamedata.GetChaperCopyEndID(player.CopyMoudle.Main.CurChapter, gamedata.COPY_TYPE_Main)
		if req.CopyID == endCopyID {
			player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PASS_MAIN_COPY_CHAPTER, req.Chapter)
		}

		if pCopyInfo.FirstAward != 0 && pCopyInfo.CopyID > player.CopyMoudle.Main.CurCopyID {
			awardItems := gamedata.GetItemsFromAwardID(pCopyInfo.FirstAward)
			for _, v := range awardItems {
				var item msg.MSG_ItemData
				item.ID = v.ItemID
				item.Num = v.ItemNum
				response.FirstItem = append(response.FirstItem, item)
			}
			player.BagMoudle.AddAwardItems(awardItems)
		}

		player.CopyMoudle.PlayerPassMainLevels(req.CopyID, req.Chapter, req.StarNum)

		//! 随机出现叛军围剿
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		//! 随机出现黑市
		isBlackMarket := false
		if player.BlackMarketModule.IsOpen == false && player.GetVipLevel() < int8(gamedata.EnterVipLevel) && player.GetLevel() >= 30 {
			randValue := random.Intn(1000)
			randValue = 1
			if randValue < gamedata.BlackMarketPro {
				//! 随机出现黑市
				player.BlackMarketModule.RefreshGoods(true)
				response.OpenEndTime = player.BlackMarketModule.OpenEndTime
				isBlackMarket = true
			}
		}

		isHadRebel := player.RebelModule.IsHaveRebel()
		if isHadRebel == false && isBlackMarket != true && player.GetLevel() >= 35 {
			randValue := random.Intn(100)

			//! 随机出现叛军
			if randValue < gamedata.FindRebelPro {
				//! 随机叛军属性
				player.RebelModule.RandRebel()
				response.IsFindRebel = true
			}
		}

	} else if req.CopyType == gamedata.COPY_TYPE_Elite {
		ok, errcode := CopyCheck(player, req.CopyID, req.Chapter, gamedata.COPY_TYPE_Elite)
		if ok != true {
			response.RetCode = errcode
			return
		}

		endCopyID := gamedata.GetChaperCopyEndID(player.CopyMoudle.Elite.CurChapter, gamedata.COPY_TYPE_Elite)
		if req.CopyID == endCopyID {
			player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PASS_ELITE_COPY_CHAPTER, req.Chapter)
		}

		player.CopyMoudle.PlayerPassEliteLevels(req.CopyID, req.Chapter, req.StarNum)

		//! 随机出现叛军围剿
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		//! 随机出现黑市
		isBlackMarket := false
		if player.BlackMarketModule.IsOpen == false && player.GetVipLevel() < int8(gamedata.EnterVipLevel) && player.GetLevel() >= 30 {
			randValue := random.Intn(1000)

			if randValue < gamedata.BlackMarketPro {
				//! 随机出现黑市
				player.BlackMarketModule.RefreshGoods(true)
				response.OpenEndTime = player.BlackMarketModule.OpenEndTime
				isBlackMarket = true
			}
		}

		isHadRebel := player.RebelModule.IsHaveRebel()
		if isHadRebel == false && isBlackMarket != true && player.GetLevel() >= 35 {
			randValue := random.Intn(100)

			//! 随机出现叛军
			if randValue < gamedata.FindRebelPro {
				//! 随机叛军属性
				player.RebelModule.RandRebel()
				response.IsFindRebel = true
			}
		}

	} else if req.CopyType == gamedata.COPY_TYPE_Daily { //! 通关日常副本
		ok, errcode := CopyCheck(player, req.CopyID, req.Chapter, gamedata.COPY_TYPE_Daily)
		if ok != true {
			response.RetCode = errcode
			return
		}
		player.CopyMoudle.PlayerPassDailyLevels(req.CopyID)
	} else if req.CopyType == gamedata.COPY_TYPE_Famous { //! 通关名将副本
		ok, errcode := CopyCheck(player, req.CopyID, req.Chapter, req.CopyType)
		if ok != true {
			response.RetCode = errcode
			return
		}

		//! 记录通关并判断是否领取首胜奖励
		isFirstVictory := player.CopyMoudle.PlayerPassFamousLevels(req.CopyID, req.Chapter)
		if isFirstVictory == true {
			firstVictoryAward := gamedata.GetItemsFromAwardID(pCopyInfo.FirstAward)
			for _, v := range firstVictoryAward {
				var item msg.MSG_ItemData
				item.ID = v.ItemID
				item.Num = v.ItemNum
				response.FirstItem = append(response.FirstItem, item)
			}
			player.BagMoudle.AddAwardItems(firstVictoryAward)
		}
	}

	//! 给予玩家经验
	response.Exp = pCopyInfo.Experience * player.GetLevel()

	//! 工会技能经验加成
	if player.HeroMoudle.GuildSkiLvl[8] > 0 {
		expInc := gamedata.GetGuildSkillExpValue(player.HeroMoudle.GuildSkiLvl[8])
		response.Exp += response.Exp * expInc / 1000
	}

	if response.Exp != 0 {
		player.HeroMoudle.AddMainHeroExp(response.Exp)
	}

	//! 给予玩家货币
	moneyNum := pCopyInfo.MoneyNum * player.GetLevel()
	if moneyNum != 0 {
		player.RoleMoudle.AddMoney(pCopyInfo.MoneyID, moneyNum)
	}

	//! 扣除体力
	player.RoleMoudle.CostAction(pCopyInfo.ActionType, pCopyInfo.ActionValue)

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(pCopyInfo.ActionType)

	//! 掉落物品
	response.ItemLst = []msg.MSG_ItemData{}
	dropItem := gamedata.GetItemsFromAwardID(pCopyInfo.AwardID)
	for _, v := range dropItem {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}

	if dropItem != nil {
		player.BagMoudle.AddAwardItems(dropItem)
	}

	response.RetCode = msg.RE_SUCCESS
}

func Hand_SweepCopy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_BattleResult_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_SweepCopy : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_BattleResult_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.StarNum <= 0 || req.StarNum > 3 {
		gamelog.Error("Hand_SweepCopy error : Invalid Star Num %d, playerid:%d", req.StarNum, req.PlayerID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	pCopyInfo := gamedata.GetCopyBaseInfo(req.CopyID)
	if pCopyInfo == nil {
		gamelog.Error("Hand_SweepCopy error : Invalid copyid : %d", req.CopyID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	isEnough := player.RoleMoudle.CheckActionEnough(pCopyInfo.ActionType, pCopyInfo.ActionValue)
	if isEnough == false { //! 体力不足挑战
		gamelog.Error("Hand_SweepCopy error : Not Enough Action")
		response.RetCode = msg.RE_STRENGTH_NOT_ENOUGH
		return
	}

	if req.CopyType == gamedata.COPY_TYPE_Main { //! 通关主线关卡
		ok, errcode := CopyCheck(player, req.CopyID, req.Chapter, gamedata.COPY_TYPE_Main)
		if ok != true {
			response.RetCode = errcode
			return
		}

		//! 检查副本是否通关
		if req.CopyID > player.CopyMoudle.Main.CurCopyID {
			gamelog.Error("Hand_SweepCopy error : Not Pass Level")
			response.RetCode = msg.RE_NEED_PASS_PRE_COPY
			return
		}

		//! 检测扫荡副本是否三星通关
		for _, v := range player.CopyMoudle.Main.CopyInfo {
			if v.CopyID == req.CopyID && v.StarNum != 3 {
				response.RetCode = msg.RE_STAR_NOT_ENOUGH
				return
			}
		}

		player.CopyMoudle.PlayerPassMainLevels(req.CopyID, req.Chapter, req.StarNum)

		//! 随机出现叛军围剿
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		//! 随机出现黑市
		isBlackMarket := false
		if player.BlackMarketModule.IsOpen == false && player.GetVipLevel() < int8(gamedata.EnterVipLevel) && player.GetLevel() >= 30 {
			randValue := random.Intn(1000)

			if randValue < gamedata.BlackMarketPro {
				//! 随机出现黑市
				player.BlackMarketModule.RefreshGoods(true)
				response.OpenEndTime = player.BlackMarketModule.OpenEndTime
				isBlackMarket = true

			}
		}

		isHadRebel := player.RebelModule.IsHaveRebel()
		if isHadRebel == false && isBlackMarket != true {
			randValue := random.Intn(100)

			//! 随机出现叛军
			if randValue < gamedata.FindRebelPro {
				//! 随机叛军属性
				player.RebelModule.RandRebel()
				response.IsFindRebel = true
			}
		}
	} else if req.CopyType == gamedata.COPY_TYPE_Elite {
		ok, errcode := CopyCheck(player, req.CopyID, req.Chapter, gamedata.COPY_TYPE_Elite)
		if ok != true {
			response.RetCode = errcode
			return
		}

		//! 检查副本是否通关
		if req.CopyID > player.CopyMoudle.Elite.CurCopyID {
			gamelog.Error("Hand_SweepCopy error : Not Pass Level")
			response.RetCode = msg.RE_NEED_PASS_PRE_COPY
			return
		}

		//! 检测扫荡副本是否三星通关
		for _, v := range player.CopyMoudle.Elite.CopyInfo {
			if v.CopyID == req.CopyID && v.StarNum != 3 {
				response.RetCode = msg.RE_STAR_NOT_ENOUGH
				return
			}
		}

		//! 随机出现叛军围剿
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		//! 随机出现黑市
		isBlackMarket := false
		if player.BlackMarketModule.IsOpen == false && player.GetVipLevel() < int8(gamedata.EnterVipLevel) && player.GetLevel() >= 30 {
			randValue := random.Intn(1000)

			if randValue < gamedata.BlackMarketPro {
				//! 随机出现黑市
				player.BlackMarketModule.RefreshGoods(true)
				response.OpenEndTime = player.BlackMarketModule.OpenEndTime
				isBlackMarket = true
			}
		}

		isHadRebel := player.RebelModule.IsHaveRebel()
		if isHadRebel == false && isBlackMarket != true {
			randValue := random.Intn(100)

			//! 随机出现叛军
			if randValue < gamedata.FindRebelPro {
				//! 随机叛军属性
				player.RebelModule.RandRebel()
				response.IsFindRebel = true
			}
		}

	}

	//! 给予玩家经验
	response.Exp = pCopyInfo.Experience * player.GetLevel()
	//! 工会技能经验加成
	if player.HeroMoudle.GuildSkiLvl[8] > 0 {
		expInc := gamedata.GetGuildSkillExpValue(player.HeroMoudle.GuildSkiLvl[8])
		response.Exp += response.Exp * expInc / 1000
	}
	player.HeroMoudle.AddMainHeroExp(response.Exp)

	//! 给予玩家货币
	moneyNum := pCopyInfo.MoneyNum * player.GetLevel()
	player.RoleMoudle.AddMoney(pCopyInfo.MoneyID, moneyNum)

	//! 给予玩家觉醒道具
	player.BagMoudle.AddWakeItem(pCopyInfo.MoneyID, pCopyInfo.MoneyNum*player.GetLevel())

	//! 扣除体力
	player.RoleMoudle.CostAction(pCopyInfo.ActionType, pCopyInfo.ActionValue)

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(pCopyInfo.ActionType)

	//! 掉落物品
	response.ItemLst = make([]msg.MSG_ItemData, 0, 5)
	dropItem := gamedata.GetItemsFromAwardID(pCopyInfo.AwardID)
	for _, v := range dropItem {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}
	player.BagMoudle.AddAwardItems(dropItem)

	response.RetCode = msg.RE_SUCCESS
}

//! 获取叛军简单信息
func Hand_GetRebelFindInfo(w http.ResponseWriter, r *http.Request) {
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetRebelFindInfo_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetRebelFindInfo : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetRebelFindInfo_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 获取叛军信息
	rebelInfo := gamedata.GetRebelInfo(player.RebelModule.RebelID)
	if rebelInfo != nil {
		response.RebelID = player.RebelModule.RebelID
		response.Level = player.RebelModule.GetRebelLevel()

		//! 发放发现叛军奖励
		awardItem := gamedata.GetRebelActionAward(gamedata.Find_Rebel, rebelInfo.Difficulty)

		var award TAwardData
		award.TextType = Text_Arean_Win
		award.ItemLst = []gamedata.ST_ItemData{gamedata.ST_ItemData{awardItem.ItemID, awardItem.ItemNum}}
		award.Time = time.Now().Unix()

		rebelcopy := gamedata.GetCopyBaseInfo(rebelInfo.CopyID)
		award.Value = []string{rebelcopy.Name}

		SendAwardToPlayer(player.playerid, &award)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家获取精英副本入侵信息
func Hand_GetEliteCopyInvadeInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetEliteInvadeStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetEliteInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetEliteInvadeStatus_Ack
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

	//! 副本重置
	player.CopyMoudle.CheckReset()

	//! 入侵检测
	player.CopyMoudle.CheckEliteInvade()

	response.InvadeID = []int{}
	response.InvadeID = append(response.InvadeID, player.CopyMoudle.Elite.InvadeChapter...)

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家获取精英副本信息
// func Hand_GetEliteInfo(w http.ResponseWriter, r *http.Request) {
// 	gamelog.Info("message: %s", r.URL.String())

// 	//! 接收消息
// 	buffer := make([]byte, r.ContentLength)
// 	r.Body.Read(buffer)

// 	var req msg.MSG_GetEliteChapterInfo_Req
// 	err := json.Unmarshal(buffer, &req)
// 	if err != nil {
// 		gamelog.Error("Hand_GetEliteInfo Unmarshal fail. Error: %s", err.Error())
// 		return
// 	}

// 	//! 创建回复
// 	var response msg.MSG_GetEliteChapterInfo_Ack
// 	response.RetCode = msg.RE_UNKNOWN_ERR

// 	defer func() {
// 		b, _ := json.Marshal(&response)
// 		w.Write(b)
// 	}()

// 	//! 常规检查
// 	var player *TPlayer = nil
// 	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
// 	if player == nil {
// 		return
// 	}

// 	//! 副本重置
// 	player.CopyMoudle.UpdateTimeReset()

// 	//! 入侵检测
// 	player.CopyMoudle.CheckEliteInvade()

// 	response.Info.CurChapter = player.CopyMoudle.Elite.CurChapter
// 	response.Info.CurCopyID = player.CopyMoudle.Elite.CurCopyID

// 	for _, v := range player.CopyMoudle.Elite.CopyInfo {
// 		var copyInfo msg.TEliteCopy
// 		copyInfo.CopyID = v.CopyID
// 		copyInfo.BattleTimes = v.BattleTimes
// 		copyInfo.ResetCount = v.ResetCount
// 		copyInfo.StarNum = v.StarNum
// 		response.Info.CopyInfo = append(response.Info.CopyInfo, copyInfo)
// 	}

// 	for _, v := range player.CopyMoudle.Elite.Chapter {
// 		var chapter msg.TEliteChapter
// 		chapter.Chapter = v.Chapter
// 		chapter.SceneAward = v.SceneAward
// 		chapter.StarAward = v.StarAward
// 		response.Info.Chapter = append(response.Info.Chapter, chapter)
// 	}

// 	response.RetCode = msg.RE_SUCCESS
// }

//! 玩家获取主线副本章节详细信息
// func Hand_GetMainDetailInfo(w http.ResponseWriter, r *http.Request) {
// 	gamelog.Info("message: %s", r.URL.String())

// 	//! 接收消息
// 	buffer := make([]byte, r.ContentLength)
// 	r.Body.Read(buffer)

// 	var req msg.MSG_GetMainDetailInfo_Req
// 	err := json.Unmarshal(buffer, &req)
// 	if err != nil {
// 		gamelog.Error("Hand_GetMainDetailInfo Unmarshal fail. Error: %s", err.Error())
// 		return
// 	}

// 	//! 创建回复
// 	var response msg.MSG_GetMainDetailInfo_Ack
// 	response.RetCode = msg.RE_UNKNOWN_ERR

// 	defer func() {
// 		b, _ := json.Marshal(&response)
// 		w.Write(b)
// 	}()

// 	//! 常规检查
// 	var player *TPlayer = nil
// 	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
// 	if player == nil {
// 		return
// 	}

// 	if req.Chapter >= len(player.CopyMoudle.Main.Chapter) {
// 		response.RetCode = msg.RE_UNKNOWN_ERR
// 		gamelog.Error("Hand_GetMainDetailInfo Invalid Chapter: %d", req.Chapter)
// 		return
// 	}

// 	player.CopyMoudle.UpdateTimeReset()

// 	//! 获取该章节关卡的详细信息
// 	for _, v := range player.CopyMoudle.Main.CopyInfo {
// 		var info msg.MainChapterDetailInfo
// 		info.CopyID = v.CopyID
// 		info.BattleTimes = v.BattleTimes
// 		info.StarNum = v.StarNum
// 		response.Info = append(response.Info, info)
// 	}

// 	for _, v := range player.CopyMoudle.Main.Chapter {
// 		if v.Chapter == req.Chapter {
// 			response.SceneAward = v.SceneAward
// 			response.StarAward = v.StarAward
// 		}
// 	}

// 	response.RetCode = msg.RE_SUCCESS
// }

//! 玩家获取精英副本章节详细信息
// func Hand_GetEliteDetailInfo(w http.ResponseWriter, r *http.Request) {
// 	gamelog.Info("message: %s", r.URL.String())

// 	//! 接收消息
// 	buffer := make([]byte, r.ContentLength)
// 	r.Body.Read(buffer)

// 	var req msg.MSG_GetEliteDetailInfo_Req
// 	err := json.Unmarshal(buffer, &req)
// 	if err != nil {
// 		gamelog.Error("Hand_GetEliteDetailInfo Unmarshal fail. Error: %s", err.Error())
// 		return
// 	}

// 	//! 创建回复
// 	var response msg.MSG_GetEliteDetailInfo_Ack
// 	response.RetCode = msg.RE_UNKNOWN_ERR

// 	defer func() {
// 		b, _ := json.Marshal(&response)
// 		w.Write(b)
// 	}()

// 	//! 常规检查
// 	var player *TPlayer = nil
// 	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
// 	if player == nil {
// 		return
// 	}

// 	player.CopyMoudle.UpdateTimeReset()
// 	if req.Chapter >= len(player.CopyMoudle.Elite.Chapter) {
// 		response.RetCode = msg.RE_UNKNOWN_ERR
// 		gamelog.Error("Hand_GetEliteInfo Invalid Chapter: %d", req.Chapter)
// 		return
// 	}

// 	//! 获取该章节关卡的详细信息
// 	for _, v := range player.CopyMoudle.Elite.CopyInfo {
// 		var info msg.MainChapterDetailInfo
// 		info.CopyID = v.CopyID
// 		info.BattleTimes = v.BattleTimes
// 		info.StarNum = v.StarNum
// 		response.Info = append(response.Info, info)
// 	}

// 	for _, v := range player.CopyMoudle.Elite.Chapter {
// 		if v.Chapter == req.Chapter {
// 			response.SceneAward = v.SceneAward
// 			response.StarAward = v.StarAward
// 		}
// 	}

// 	response.RetCode = msg.RE_SUCCESS
// }

//! 领取主线副本星级奖励
func Hand_GetMainStarAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetMainStarAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMainStarAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetMainStarAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	//! 检查是否已经领取
	if len(player.CopyMoudle.Main.Chapter) <= 0 {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	for _, v := range player.CopyMoudle.Main.Chapter {
		if v.Chapter == req.Chapter && v.StarAward[req.StarAward-1] != false {
			response.RetCode = msg.RE_ALREADY_RECEIVED
			return
		}
	}

	//! 检查是否够格领取
	chapterStarNumber := player.CopyMoudle.GetMainChapterStarNumber(req.Chapter)
	chapterData := gamedata.GetMainChapterInfo(req.Chapter)
	if chapterStarNumber < (chapterData.StarAwards[req.StarAward-1].StarNum) {
		response.RetCode = msg.RE_STAR_NOT_ENOUGH
		gamelog.Error("Hand_GetMainStarAward error: star not enough: %d", chapterStarNumber)
		return
	}

	//! 发放星级奖励
	player.CopyMoudle.PaymentMainAward(req.Chapter, req.StarAward-1, MAIN_AWARD_TYPE_STAR)

	response.RetCode = msg.RE_SUCCESS
}

//! 领取精英副本星级奖励
func Hand_GetEliteStarAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetEliteStarAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetEliteStarAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetEliteStarAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	//! 检查是否已经领取
	if len(player.CopyMoudle.Elite.Chapter) <= 0 {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	for _, v := range player.CopyMoudle.Elite.Chapter {
		if v.Chapter == req.Chapter && v.StarAward[req.StarAward-1] != false {
			response.RetCode = msg.RE_ALREADY_RECEIVED
			return
		}
	}

	//! 检查是否够格领取
	chapterStarNumber := player.CopyMoudle.GetEliteChapterStarNumber(req.Chapter)
	chapterData := gamedata.GetEliteChapterInfo(req.Chapter)
	if chapterStarNumber < (chapterData.StarAwards[req.StarAward-1].StarNum) {
		response.RetCode = msg.RE_STAR_NOT_ENOUGH
		gamelog.Error("Hand_GetMainStarAward error: star not enough: %d", chapterStarNumber)
		return
	}

	//! 发放星级奖励
	player.CopyMoudle.PaymentEliteAward(req.Chapter, req.StarAward-1, MAIN_AWARD_TYPE_STAR)

	response.RetCode = msg.RE_SUCCESS
}

//! 领取主线副本场景奖励
func Hand_GetMainSceneAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetMainSceneAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMainSceneAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetMainSceneAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	//! 检查是否已经领取
	isExist := false
	for _, v := range player.CopyMoudle.Main.Chapter {
		if v.Chapter == req.Chapter {
			if v.SceneAward[req.SceneAward-1] != false {
				response.RetCode = msg.RE_ALREADY_RECEIVED
				return
			}
			isExist = true
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 检查是否够格领取
	chapterData := gamedata.GetMainChapterInfo(req.Chapter)
	needCopyID := chapterData.SceneAwards[req.SceneAward-1].Levels
	if player.CopyMoudle.Main.CurCopyID < needCopyID {
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	//! 发放奖励
	player.CopyMoudle.PaymentMainAward(req.Chapter, req.SceneAward-1, MAIN_AWARD_TYPE_SCENE)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求精英副本场景奖励
func Hand_GetEliteSceneAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetEliteSceneAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetEliteSceneAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetEliteSceneAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	//! 检查是否已经领取
	isExist := false
	for _, v := range player.CopyMoudle.Elite.Chapter {
		if v.Chapter == req.Chapter {
			if v.SceneAward != false {
				response.RetCode = msg.RE_ALREADY_RECEIVED
				return
			}
			isExist = true
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 检查是否够格领取
	chapterData := gamedata.GetEliteChapterInfo(req.Chapter)
	needCopyID := chapterData.SceneAwards.Levels
	if player.CopyMoudle.Elite.CurCopyID < needCopyID {
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	//! 发放奖励
	player.CopyMoudle.PaymentEliteAward(req.Chapter, 0, MAIN_AWARD_TYPE_SCENE)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家查询主线副本重置次数
func Hand_GetMainResetTimes(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetMainRefreshTimes_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMainResetTimes Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetMainRefreshTimes_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	if req.Chapter > player.CopyMoudle.Main.CurChapter {
		gamelog.Error("Hand_GetMainResetTimes Error Invalid Chapter:%d", req.Chapter)
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	vipLevel := player.GetVipLevel()
	refreshLimit := gamedata.GetFuncVipValue(gamedata.FUNC_MAIN_COPY_RESET, vipLevel)

	isExist := false
	for _, v := range player.CopyMoudle.Main.CopyInfo {
		if v.CopyID == req.CopyID {
			response.RefreshTimes = refreshLimit - v.ResetCount
			isExist = true
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_COPY_NOT_PASS
		return
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家查询精英副本重置次数
func Hand_GetEliteResetTimes(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetEliteRefreshTimes_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetEliteResetTimes Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetEliteRefreshTimes_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	if req.Chapter > player.CopyMoudle.Elite.CurChapter {
		gamelog.Error("Hand_GetMainResetTimes Error Invalid Chapter:%d", req.Chapter)
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	vipLevel := player.GetVipLevel()
	refreshLimit := gamedata.GetFuncVipValue(gamedata.FUNC_MAIN_COPY_RESET, vipLevel)

	isExist := false
	for _, v := range player.CopyMoudle.Elite.CopyInfo {
		if v.CopyID == req.CopyID {
			response.RefreshTimes = refreshLimit - v.ResetCount
			isExist = true
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_COPY_NOT_PASS
		return
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求主线副本重置挑战
func Hand_ResetMainBattleTimes(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ResetMainBattleTimes_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ResetMainBattleTimes Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ResetMainBattleTimes_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	//! 获取重置次数额度
	vipLevel := player.GetVipLevel()
	refreshLimit := gamedata.GetFuncVipValue(gamedata.FUNC_MAIN_COPY_RESET, vipLevel)

	//! 检查当前刷新次数
	var copyInfo *TMainCopy
	curRefreshCounts := 0
	isExist := false
	copyIndex := 0
	for index, v := range player.CopyMoudle.Main.CopyInfo {
		if v.CopyID == req.CopyID {
			curRefreshCounts = v.ResetCount
			copyInfo = &player.CopyMoudle.Main.CopyInfo[index]
			copyIndex = index
			isExist = true
			break
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_COPY_NOT_PASS
		return
	}

	//! 可重置次数不足
	if curRefreshCounts >= refreshLimit {
		gamelog.Error("Hand_ResetMainBattleTimes Error: Refresh times not enough now: %d  limit: %d", curRefreshCounts, refreshLimit)
		response.RetCode = msg.RE_REFRESH_TIMES_NOT_ENOUGH
		return
	}

	cost := gamedata.GetFuncTimeCost(gamedata.FUNC_MAIN_COPY_RESET, curRefreshCounts+1)

	//! 判断玩家元宝
	if player.RoleMoudle.CheckMoneyEnough(gamedata.EliteCopyResetMoneyID, cost) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除元宝,完成重置
	player.RoleMoudle.CostMoney(gamedata.EliteCopyResetMoneyID, cost)

	copyInfo.ResetCount += 1
	copyInfo.BattleTimes = 0

	response.MoneyID = gamedata.EliteCopyResetMoneyID
	response.MoneyNum = cost

	response.RetCode = msg.RE_SUCCESS

	go player.CopyMoudle.DB_UpdateMainCopyAt(copyIndex)
}

//! 玩家请求精英副本重置挑战
func Hand_ResetEliteBattleTimes(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ResetEliteBattleTimes_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ResetEliteBattleTimes Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ResetEliteBattleTimes_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	//! 获取重置次数额度
	vipLevel := player.GetVipLevel()
	refreshLimit := gamedata.GetFuncVipValue(gamedata.FUNC_MAIN_COPY_RESET, vipLevel)

	//! 检查当前刷新次数
	var copyInfo *TEliteCopy
	curRefreshCounts := 0
	isExist := false
	copyIndex := 0
	for index, v := range player.CopyMoudle.Elite.CopyInfo {
		if v.CopyID == req.CopyID {
			curRefreshCounts = v.ResetCount
			copyInfo = &player.CopyMoudle.Elite.CopyInfo[index]
			copyIndex = index
			isExist = true
			break
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_COPY_NOT_PASS
		return
	}

	//! 可重置次数不足
	if curRefreshCounts >= refreshLimit {
		response.RetCode = msg.RE_REFRESH_TIMES_NOT_ENOUGH
		return
	}

	var cost int
	cost = gamedata.GetFuncTimeCost(gamedata.FUNC_MAIN_COPY_RESET, curRefreshCounts+1)

	//! 判断玩家元宝
	if player.RoleMoudle.CheckMoneyEnough(gamedata.EliteCopyResetMoneyID, cost) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除元宝,完成重置
	player.RoleMoudle.CostMoney(gamedata.EliteCopyResetMoneyID, cost)

	response.MoneyID = gamedata.EliteCopyResetMoneyID
	response.MoneyNum = cost

	copyInfo.ResetCount += 1
	copyInfo.BattleTimes = 0

	response.RetCode = msg.RE_SUCCESS

	go player.CopyMoudle.DB_UpdateEliteCopyAt(copyIndex)
}

//! 玩家请求获取日常副本信息
// func Hand_GetDailyCopyInfo(w http.ResponseWriter, r *http.Request) {
// 	gamelog.Info("message: %s", r.URL.String())

// 	//! 接收消息
// 	buffer := make([]byte, r.ContentLength)
// 	r.Body.Read(buffer)

// 	//! 解析消息
// 	var req msg.MSG_GetDailyCopyInfo_Req
// 	err := json.Unmarshal(buffer, &req)
// 	if err != nil {
// 		gamelog.Error("Hand_GetDailyCopyInfo Unmarshal fail. Error: %s", err.Error())
// 		return
// 	}

// 	//! 创建回复
// 	var response msg.MSG_GetDailyCopyInfo_Ack
// 	response.RetCode = msg.RE_UNKNOWN_ERR
// 	defer func() {
// 		b, _ := json.Marshal(&response)
// 		w.Write(b)

// 	}()

// 	//! 常规检测
// 	var player *TPlayer = nil
// 	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
// 	if player == nil {
// 		return
// 	}

// 	player.CopyMoudle.UpdateTimeReset()

// 	//! 获取今天开启的副本
// 	todayDailyCopy := player.CopyMoudle.GetTodayDailyCopy()

// 	response.RetCode = msg.RE_SUCCESS

// 	for _, b := range todayDailyCopy {
// 		if b == 0 {
// 			continue
// 		}

// 		var data msg.MSG_DailyCopy
// 		data.IsChallenge = false
// 		data.ResType = b

// 		for _, v := range player.CopyMoudle.Daily.CopyInfo {
// 			if b == v.ResID {
// 				//! 若有挑战记录,则返回挑战信息
// 				data.IsChallenge = v.IsChallenge
// 				data.ResType = v.ResID
// 			}
// 		}
// 		response.CopyInfo = append(response.CopyInfo, data)
// 	}

// }

//! 玩家请求获取名将副本章节信息
// func Hand_GetFamousCopyChapterInfo(w http.ResponseWriter, r *http.Request) {
// 	gamelog.Info("message: %s", r.URL.String())

// 	//! 接收消息
// 	buffer := make([]byte, r.ContentLength)
// 	r.Body.Read(buffer)

// 	//! 解析消息
// 	var req msg.MSG_GetFamousCopyChapterInfo_Req
// 	err := json.Unmarshal(buffer, &req)
// 	if err != nil {
// 		gamelog.Error("Hand_GetFamousCopyChapterInfo Unmarshal fail. Error: %s", err.Error())
// 		return
// 	}

// 	//! 创建回复
// 	var response msg.MSG_GetFamousCopyChapterInfo_Ack
// 	response.RetCode = msg.RE_UNKNOWN_ERR
// 	defer func() {
// 		b, _ := json.Marshal(&response)
// 		w.Write(b)

// 	}()

// 	//! 常规检测
// 	var player *TPlayer = nil
// 	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
// 	if player == nil {
// 		return
// 	}

// 	response.CurCopyID = player.CopyMoudle.Famous.CurCopyID

// 	for i, _ := range player.CopyMoudle.Famous.Chapter {
// 		for _, n := range player.CopyMoudle.Famous.Chapter[i].PassedCopy {
// 			if player.CopyMoudle.Famous.CurCopyID == n.CopyID {
// 				response.CurChapter = i
// 			}

// 		}
// 	}

// 	response.BattleTimes = player.CopyMoudle.Famous.BattleTimes
// 	response.RetCode = msg.RE_SUCCESS
// }

//! 玩家请求获取名将副本详细信息
func Hand_GetFamousCopyDetailInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetFamousCopyDetailInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFamousCopyDetailInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetFamousCopyDetailInfo_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 获取玩家指定关卡
	chapterInfo := player.CopyMoudle.Famous.Chapter[req.Chapter]

	for _, v := range chapterInfo.PassedCopy {
		var info msg.MSG_FamousCopyDetailInfo
		info.CopyID = v.CopyID
		info.BattleTimes = v.BattleTimes
		response.CopyLst = append(response.CopyLst, info)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求攻击精英副本入侵
func Hand_AttackInvade(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_AttackEliteInvade_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_AttackInvade Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_AttackEliteInvade_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测该章节是否有入侵
	if player.CopyMoudle.IsHaveInvade(req.Chapter) == false {
		response.RetCode = msg.RE_INVADE_ALEADY_ESCAPE
		return
	}

	//! 获得掉落奖励
	chapterInfo := gamedata.GetEliteChapterInfo(req.Chapter)
	if chapterInfo == nil {
		gamelog.Error("Hand_AttackInvade GetEliteChapterInfo fail. Chapter: %d", req.Chapter)
		return
	}

	copyInfo := gamedata.GetCopyBaseInfo(chapterInfo.InvadeID)
	if copyInfo == nil {
		gamelog.Error("Hand_AttackInvade GetCopyBaseInfo fail. InvadeID: %d", chapterInfo.InvadeID)
		return
	}

	if player.RoleMoudle.CheckActionEnough(copyInfo.ActionType, copyInfo.ActionValue) == false {
		gamelog.Error("Hand_AttackInvade CheckActionEnough fail.")
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	awardItems := gamedata.GetItemsFromAwardID(copyInfo.AwardID)
	player.BagMoudle.AddAwardItems(awardItems)

	response.Exp = copyInfo.Experience
	//! 工会技能经验加成
	if player.HeroMoudle.GuildSkiLvl[8] != 0 {
		expInc := gamedata.GetGuildSkillExpValue(player.HeroMoudle.GuildSkiLvl[8])
		response.Exp += response.Exp * expInc / 1000
	}

	player.HeroMoudle.AddMainHeroExp(response.Exp)

	for _, v := range awardItems {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.DropItems = append(response.DropItems, item)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求名将副本章节奖励
func Hand_GetFamousCopyAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetFamousCopyAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFamousCopyAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetFamousCopyAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断参数合法
	if req.Chapter > len(player.CopyMoudle.Famous.Chapter)-1 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断领取状态
	if player.CopyMoudle.Famous.Chapter[req.Chapter].ChapterAward == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 判断通关状态
	chapterInfo := gamedata.GetFamousChapterInfo(req.Chapter)
	if player.CopyMoudle.Famous.CurCopyID < chapterInfo.EndID {
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	//! 发放奖励
	itemLsts := gamedata.GetItemsFromAwardID(chapterInfo.Award)
	player.BagMoudle.AddAwardItems(itemLsts)

	//! 记录状态
	player.CopyMoudle.Famous.Chapter[req.Chapter].ChapterAward = true
	go player.CopyMoudle.DB_UpdateFamousAward(req.Chapter)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求副本数据
func Hand_GetCopyData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetCopyData_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetCopyData Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetCopyData_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.CopyMoudle.CheckReset()

	//! 主线
	response.CopyMainInfo.CurChapter = player.CopyMoudle.Main.CurChapter
	response.CopyMainInfo.CurCopyID = player.CopyMoudle.Main.CurCopyID

	refreshLimit := gamedata.GetFuncVipValue(gamedata.FUNC_MAIN_COPY_RESET, player.GetVipLevel())
	for _, v := range player.CopyMoudle.Main.CopyInfo {
		var copyInfo msg.TMainCopy
		copyInfo.CopyID = v.CopyID
		copyInfo.BattleTimes = v.BattleTimes
		copyInfo.ResetCount = refreshLimit - v.ResetCount
		copyInfo.StarNum = v.StarNum
		response.CopyMainInfo.CopyInfo = append(response.CopyMainInfo.CopyInfo, copyInfo)
	}

	for _, v := range player.CopyMoudle.Main.Chapter {
		var chapter msg.TMainChapter
		chapter.Chapter = v.Chapter
		chapter.SceneAward = v.SceneAward
		chapter.StarAward = v.StarAward
		response.CopyMainInfo.Chapter = append(response.CopyMainInfo.Chapter, chapter)
	}

	//! 精英
	//! 入侵检测
	player.CopyMoudle.CheckEliteInvade()

	response.CopyEliteInfo.CurChapter = player.CopyMoudle.Elite.CurChapter
	response.CopyEliteInfo.CurCopyID = player.CopyMoudle.Elite.CurCopyID

	for _, v := range player.CopyMoudle.Elite.CopyInfo {
		var copyInfo msg.TEliteCopy
		copyInfo.CopyID = v.CopyID
		copyInfo.BattleTimes = v.BattleTimes
		copyInfo.ResetCount = v.ResetCount
		copyInfo.StarNum = v.StarNum
		response.CopyEliteInfo.CopyInfo = append(response.CopyEliteInfo.CopyInfo, copyInfo)
	}

	for _, v := range player.CopyMoudle.Elite.Chapter {
		var chapter msg.TEliteChapter
		chapter.Chapter = v.Chapter
		chapter.SceneAward = v.SceneAward
		chapter.StarAward = v.StarAward
		response.CopyEliteInfo.Chapter = append(response.CopyEliteInfo.Chapter, chapter)
	}

	//! 名将
	response.CopyFamousInfo.CurCopyID = player.CopyMoudle.Famous.CurCopyID

	for i, _ := range player.CopyMoudle.Famous.Chapter {
		for _, n := range player.CopyMoudle.Famous.Chapter[i].PassedCopy {
			if player.CopyMoudle.Famous.CurCopyID == n.CopyID {
				response.CopyFamousInfo.CurChapter = i
			}

		}
	}
	response.CopyFamousInfo.BattleTimes = player.CopyMoudle.Famous.BattleTimes

	//! 日常
	//! 获取今天开启的副本
	todayDailyCopy := player.CopyMoudle.GetTodayDailyCopy()

	for _, b := range todayDailyCopy {
		if b == 0 {
			continue
		}

		var data msg.MSG_DailyCopy
		data.IsChallenge = false
		data.ResType = b

		for _, v := range player.CopyMoudle.Daily.CopyInfo {
			if b == v.ResID {
				//! 若有挑战记录,则返回挑战信息
				data.IsChallenge = v.IsChallenge
				data.ResType = v.ResID
			}
		}
		response.CopyDailyInfo = append(response.CopyDailyInfo, data)
	}

	response.RetCode = msg.RE_SUCCESS
}

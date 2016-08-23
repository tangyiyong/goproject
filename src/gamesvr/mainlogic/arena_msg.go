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

//! 玩家请求竞技场信息
func Hand_GetArenaInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetArenaInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetArenaInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetArenaInfo_Ack
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

	//! 检测功能是否开启
	isFuncOpen := gamedata.IsFuncOpen(gamedata.FUNC_ARENA, player.GetLevel(), player.GetVipLevel())
	if isFuncOpen == false {
		gamelog.Error("Hand_GetArenaInfo Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	response.SelfRank = player.ArenaModule.CurrentRank

	randValue := rand.New(rand.NewSource(time.Now().UnixNano()))

	challangeLst := player.ArenaModule.RefreshChallangeLst()
	for _, v := range challangeLst {
		var challangeInfo msg.MSG_ArenaPlayerInfo
		if v.IsRobot == false {
			//! 真实玩家逻辑
			playerInfo := G_SimpleMgr.GetSimpleInfoByID(v.PlayerID)
			if playerInfo == nil {
				gamelog.Error("GetPlayer Error: %d", v.PlayerID)
			} else {
				challangeInfo.PlayerID = v.PlayerID
				challangeInfo.Name = playerInfo.Name
				challangeInfo.Rank = v.Rank
				challangeInfo.Level = 10 + randValue.Intn(5) + 1
				challangeInfo.FightValue = playerInfo.FightValue
				challangeInfo.HeroID = playerInfo.HeroID
				challangeInfo.Quality = playerInfo.Quality
			}

			response.PlayerLst = append(response.PlayerLst, challangeInfo)
		} else {
			//! 机器人逻辑
			robotInfo := gamedata.GetRobot(v.PlayerID)
			if robotInfo == nil {
				gamelog.Error("GetRobot error: invalid robot id %d", v.PlayerID)
			} else {
				challangeInfo.PlayerID = v.PlayerID
				challangeInfo.Rank = v.Rank
				challangeInfo.Name = robotInfo.Name
				challangeInfo.Level = 10 + randValue.Intn(5) + 1
				challangeInfo.HeroID = robotInfo.Heros[0].HeroID
				challangeInfo.Quality = 2
				challangeInfo.FightValue = robotInfo.FightValue
			}

			response.PlayerLst = append(response.PlayerLst, challangeInfo)
		}

	}
	response.IDLst = []int{}
	response.IDLst = append(response.IDLst, player.ArenaModule.StoreAward...)
	response.HistoryRank = player.ArenaModule.HistoryRank
	response.RetCode = msg.RE_SUCCESS
}

func Hand_ArenaCheck(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ArenaCheck_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ArenaCheck Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ArenaCheck_Ack
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

	//! 检测功能是否开启
	isFuncOpen := gamedata.IsFuncOpen(gamedata.FUNC_ARENA, player.GetLevel(), player.GetVipLevel())
	if isFuncOpen == false {
		gamelog.Error("Hand_ArenaCheck Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 检测玩家行动力是否足够
	arenaConfig := gamedata.GetArenaConfig()
	if arenaConfig == nil {
		gamelog.Error("GetArenaConfig fail")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	copyInfo := gamedata.GetCopyBaseInfo(arenaConfig.CopyID)
	if copyInfo == nil {
		gamelog.Error("GetCopyBaseInfo fail. CopyID: %d", arenaConfig.CopyID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if player.RoleMoudle.CheckActionEnough(copyInfo.ActionType, copyInfo.ActionValue) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 检测该玩家是否在挑战队列中
	if req.Rank > 5000 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	challangeInfo := &G_Rank_List[req.Rank-1]
	if challangeInfo == nil {
		response.RetCode = msg.RE_NOT_IN_CHALLANGE_LIST
		return
	}

	response.RetCode = msg.RE_SUCCESS
	//! 可以挑战该玩家, 返回玩家详细武将信息
	if challangeInfo.IsRobot == true {
		response.TargetType = 2
		robot := gamedata.GetRobot(challangeInfo.PlayerID)
		response.PlayerData.FightValue = int32(robot.FightValue)
		response.PlayerData.Quality = int32(robot.Quality)
		response.PlayerData.PlayerID = challangeInfo.PlayerID
		response.Name = robot.Name
		for i := 0; i < BATTLE_NUM; i++ {
			response.PlayerData.Heros[i].HeroID = int32(robot.Heros[i].HeroID)
			for j := 0; j < 11; j++ {
				response.PlayerData.Heros[i].PropertyValue[j] = int32(robot.Heros[i].Propertys[j])
			}

		}
	} else {
		response.TargetType = 1
		response.PlayerData.PlayerID = int32(challangeInfo.PlayerID)
		var pHeroMoudle *THeroMoudle = nil
		target := GetPlayerByID(challangeInfo.PlayerID)
		if target != nil {
			pHeroMoudle = &player.HeroMoudle
		} else {
			var hm THeroMoudle
			if hm.OnPlayerLoad(challangeInfo.PlayerID, nil) == false {
				pHeroMoudle = &hm
			} else {
				pHeroMoudle = nil
			}
		}

		var HeroResults = make([]THeroResult, BATTLE_NUM)
		response.PlayerData.FightValue = int32(pHeroMoudle.CalcFightValue(HeroResults))
		response.PlayerData.Quality = int32(pHeroMoudle.CurHeros[0].Quality)
		for i := 0; i < BATTLE_NUM; i++ {
			response.PlayerData.Heros[i].HeroID = int32(HeroResults[i].HeroID)
			for j := 0; j < 11; j++ {
				response.PlayerData.Heros[i].PropertyValue[j] = int32(HeroResults[i].PropertyValues[j])
				response.PlayerData.Heros[i].PropertyPercent[j] = int32(HeroResults[i].PropertyPercents[j])
			}

			for j := 0; j < 5; j++ {
				response.PlayerData.Heros[i].CampDef[j] = int32(HeroResults[i].CampDef[j])
				response.PlayerData.Heros[i].CampKill[j] = int32(HeroResults[i].CampKill[j])
			}
		}
	}
}

//! 玩家反馈挑战竞技场结果
func Hand_ChallengeArenaResult(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ArenaResult_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ChallengeArenaResult Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ArenaResult_Ack
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

	//! 检测功能是否开启
	isFuncOpen := gamedata.IsFuncOpen(gamedata.FUNC_ARENA, player.GetLevel(), player.GetVipLevel())
	if isFuncOpen == false {
		gamelog.Error("Hand_ChallengeArenaResult Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 检测玩家行动力是否足够
	arenaConfig := gamedata.GetArenaConfig()
	if arenaConfig == nil {
		gamelog.Error("GetArenaConfig Fail")
	}

	copyInfo := gamedata.GetCopyBaseInfo(arenaConfig.CopyID)
	if copyInfo == nil {
		gamelog.Error("GetCopyBaseInfo fail. CopyID: %d", arenaConfig.CopyID)
		return
	}

	if player.RoleMoudle.CheckActionEnough(copyInfo.ActionType, copyInfo.ActionValue) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 检测该玩家是否在挑战队列中
	if req.Rank > 5000 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	challangeInfo := &G_Rank_List[req.Rank-1]
	if challangeInfo == nil {
		response.RetCode = msg.RE_NOT_IN_CHALLANGE_LIST
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(copyInfo.ActionType, copyInfo.ActionValue)

	//! 经验银币奖励
	player.RoleMoudle.AddMoney(copyInfo.MoneyID, copyInfo.MoneyNum*player.GetLevel())

	//! 增加玩家经验
	exp := copyInfo.Experience * player.GetLevel()

	//! 工会技能经验加成
	if player.RoleMoudle.ExpIncLvl != 0 {
		expInc := gamedata.GetGuildSkillExpValue(player.RoleMoudle.ExpIncLvl)
		exp += exp * expInc / 1000
	}

	player.HeroMoudle.AddMainHeroExp(exp)

	if req.IsVictory == 1 {
		//! 获胜的翻牌随机奖励
		awardLst := gamedata.GetItemsFromAwardIDEx(copyInfo.AwardID)
		if len(awardLst) != 3 {
			gamelog.Error("GetItemsFromAwardIDEx error: %v  awardID: %d", awardLst, copyInfo.AwardID)
			return
		}

		//! 发放奖励
		player.BagMoudle.AddAwardItem(awardLst[0].ItemID, awardLst[0].ItemNum)

		for _, v := range awardLst {
			var item msg.MSG_ItemData
			item.ID = v.ItemID
			item.Num = v.ItemNum

			response.DropItem = append(response.DropItem, item)
		}

		if len(response.DropItem) != 3 {
			gamelog.Error("Hand_ChallengeArenaResult error: rand drop item fail")
			return
		}

		//! 发放获胜奖励货币
		player.RoleMoudle.AddMoney(arenaConfig.VictoryMoneyID, arenaConfig.VictoryMoneyNum)
	} else {
		//! 发放失败奖励货币
		player.RoleMoudle.AddMoney(arenaConfig.FailedMoneyID, arenaConfig.FailedMoneyNum)

		if challangeInfo.IsRobot == false {
			challangePlayer := GetPlayerByID(challangeInfo.PlayerID)
			if challangePlayer != nil {
				challangePlayer.ArenaModule.CurrentRank = player.ArenaModule.CurrentRank
			}
			SendArenaMail(challangeInfo.PlayerID, player.RoleMoudle.Name, player.ArenaModule.CurrentRank, 0, false)
		}
	}

	response.IsVictory = req.IsVictory

	//! 记录玩家排名
	challengeRank := req.Rank
	if req.IsVictory == 1 && challengeRank < player.ArenaModule.CurrentRank {
		//! 败者排名修改
		loserID := G_Rank_List[challengeRank-1].PlayerID
		if challangeInfo.IsRobot == false {
			challangePlayer := GetPlayerByID(challangeInfo.PlayerID)
			if challangePlayer != nil {
				challangePlayer.ArenaModule.CurrentRank = player.ArenaModule.CurrentRank
			}
			player.ArenaModule.UpdateChallangeRank(challangeInfo.PlayerID, player.ArenaModule.CurrentRank)

			if player.ArenaModule.CurrentRank <= 5000 {
				G_Rank_List[player.ArenaModule.CurrentRank-1].PlayerID = loserID
				G_Rank_List[player.ArenaModule.CurrentRank-1].IsRobot = false

			}

			//! 败者邮件
			SendArenaMail(challangeInfo.PlayerID, player.RoleMoudle.Name, player.ArenaModule.CurrentRank, 1, true)
		} else {
			//! 败者若是机器人,则将玩家原本排名改为机器人信息
			if player.ArenaModule.CurrentRank <= 5000 {
				G_Rank_List[player.ArenaModule.CurrentRank-1].PlayerID = loserID
				G_Rank_List[player.ArenaModule.CurrentRank-1].IsRobot = true
			}
		}

		//! 如果为前5000,则对应修改内存数据
		if challengeRank <= 5000 {
			G_Rank_List[challengeRank-1].PlayerID = player.ArenaModule.PlayerID
			G_Rank_List[challengeRank-1].IsRobot = false
		}

		//! 胜者排名修改
		player.ArenaModule.CurrentRank = challengeRank
		player.ActivityModule.RankGift.CheckRankUp(challengeRank)

		//! 比较历史最高排名
		if challengeRank < player.ArenaModule.HistoryRank {
			//! 判断玩家当前名次是否拥有挑战元宝奖励
			moneyID, moneyNum := gamedata.GetArenaMoneyAward(player.ArenaModule.HistoryRank, challengeRank)
			if moneyID != 0 {
				response.ExtraAward.ID = moneyID
				response.ExtraAward.Num = moneyNum

				player.BagMoudle.AddAwardItem(response.ExtraAward.ID, response.ExtraAward.Num)
			}

			player.ArenaModule.HistoryRank = challengeRank
		}

		//! 存储数据
		go player.ArenaModule.UpdateRankToDatabase()

	}

	if req.IsVictory == 1 && challengeRank > player.ArenaModule.CurrentRank && challangeInfo.IsRobot == false {
		SendArenaMail(challangeInfo.PlayerID, player.RoleMoudle.Name, player.ArenaModule.CurrentRank, 1, false)
	}

	response.RetCode = msg.RE_SUCCESS
	response.HistoryRank = player.ArenaModule.HistoryRank
	response.SelfRank = player.ArenaModule.CurrentRank
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ARENA_CHALLENGE, 1)
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ARENA_RANK, player.ArenaModule.HistoryRank)
}

/*
type FightObject struct {
	ObjectID int
	Propertys [11]int
	buffers int
	size int
	HeroID int
	TargetObjectID int
	posx, posy int
	speed
}

type Move struct {
	time  int
	ObjectID int
	TargetPosx, posy int
}

type Attack struct {
	time  int
	ObjectID int
	ObjectID2 int
	Skill int
	Blood int
}

func Fight(f1 THeroMoudle, f2 THeroMoudle){
	var SideOne [6]FightObject

	for i:=0; i < 600; i++ {
		if SideOne[j].TargetObjectID == 0 {
			id = findTarget()
			if id == 0 {
				moveforward()
			}else{
				canattack()
				attack();
			}
		}

		if SideOne[j].TargetObjectID.is Alive() {
		 	Attack(i)
		}



	}

}

func FightBoss(f1 THeroMoudle, boos){
	var SideOne [6]FightObject

	for i:=0; i < 600; i++ {
		if SideOne[j].TargetObjectID == 0 {
			id = findTarget()
			if id == 0 {
				moveforward()
			}else{
				canattack()
				attack();
			}
		}

		if SideOne[j].TargetObjectID.is Alive() {
		 	Attack(i)
		}



	}
}
*/

//! 玩家请求购买声望商店已购买奖励ID
func Hand_GetArenaStoreAleadyBuyAwardLst(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryArenaStoreAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetArenaStoreAleadyBuyAwardLst Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_QueryArenaStoreAward_Ack
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

	//! 检测功能是否开启
	isFuncOpen := gamedata.IsFuncOpen(gamedata.FUNC_ARENA, player.GetLevel(), player.GetVipLevel())
	if isFuncOpen == false {
		gamelog.Error("Hand_GetArenaStoreAleadyBuyAwardLst Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	response.IDLst = []int{}
	response.IDLst = append(response.IDLst, player.ArenaModule.StoreAward...)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买声望商店物品
func Hand_BuyArenaStoreItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetArenaStoreItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyArenaStoreItem Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetArenaStoreItem_Ack
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

	//! 检测参数
	if req.Num <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_BuyArenaStoreItem invalid num. Num: %v  PlayerID: %v", req.Num, player.playerid)
		return
	}

	//! 检测功能是否开启
	isFuncOpen := gamedata.IsFuncOpen(gamedata.FUNC_ARENA, player.GetLevel(), player.GetVipLevel())
	if isFuncOpen == false {
		gamelog.Error("Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取要购买的商品
	item := gamedata.GetArenaStoreItem(req.ID)
	if item == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断玩家等级是否足够
	if player.GetLevel() < item.NeedLevel {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Not enough level")
		return
	}

	//! 如果购买物品属于奖励,则判断排名
	if item.Type == 2 {
		if player.ArenaModule.CurrentRank > item.NeedRank {
			response.RetCode = msg.RE_NOT_ENOUGH_RANK
			return
		}

		//! 判断是否已经购买
		if player.ArenaModule.StoreAward.IsExist(item.ID) >= 0 {
			response.RetCode = msg.RE_NOT_ENOUGH_TIMES
			return
		}
	}

	//! 检测金钱是否足够
	if player.RoleMoudle.CheckMoneyEnough(item.MoneyID, item.MoneyNum*req.Num) == false {
		gamelog.Error("Not enough money. NeedMoney: %d & NeedNum: %d", item.MoneyID, item.MoneyNum)
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 检测道具是否足够
	if item.CostItemID != 0 {
		if player.BagMoudle.IsItemEnough(item.CostItemID, item.CostItemNum*req.Num) == false {
			gamelog.Error("Not enough item. NeedItem: %d & NeedNum: %d", item.CostItemID, item.CostItemNum)
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			return
		}
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(item.MoneyID, item.MoneyNum*req.Num)

	//! 给予物品
	player.BagMoudle.AddAwardItem(item.ItemID, item.ItemNum*req.Num)

	//! 记录购买
	if item.Type == 2 {
		player.ArenaModule.StoreAward = append(player.ArenaModule.StoreAward, item.ID)
		go player.ArenaModule.UpdateStoreToDatabase()
	}

	//! 扣除道具
	if item.CostItemID != 0 {
		player.BagMoudle.RemoveNormalItem(item.CostItemID, item.CostItemNum*req.Num)
	}

	response.RetCode = msg.RE_SUCCESS
}

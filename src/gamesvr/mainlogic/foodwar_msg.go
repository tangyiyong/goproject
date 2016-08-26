package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

func SelectFoodTarget(player *TPlayer, value int) bool {
	if (player.FoodWarModule.TotalFood - player.FoodWarModule.FixedFood) <= value {
		return false
	}
	return true
}

//! 玩家请求挑战列表
func Hand_FoodWar_GetChallenger(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_GetChallenger_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_GetChallenger Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_GetChallenger_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动是否开启
	if player.FoodWarModule.IsActivityOpen() == false {
		gamelog.Error("Hand_FoodWar_GetChallenger Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()
	if len(g_SelectPlayers) < 5 {
		response.ChallengerLst = []msg.MSG_FoodWar_Challenger{}
	} else {
		for i := 0; i < 4; i++ {
			pTarget := GetSelectPlayer(SelectScoreTarget, 1000)
			if pTarget != nil && pTarget.playerid != 0 {
				var challenge msg.MSG_FoodWar_Challenger
				challenge.PlayerID = pTarget.playerid

				if challenge.PlayerID == player.playerid {
					i -= 1
					continue
				}

				simpleInfo := G_SimpleMgr.GetSimpleInfoByID(challenge.PlayerID)
				challenge.PlayerName = simpleInfo.Name
				challenge.Quality = simpleInfo.Quality
				challenge.HeroID = simpleInfo.HeroID
				challenge.Level = simpleInfo.Level
				challenge.FightValue = simpleInfo.FightValue
				challenge.CanRobFood = (pTarget.FoodWarModule.TotalFood - pTarget.FoodWarModule.FixedFood) * gamedata.FoodWarRobBili / 1000
				challenge.TotalFood = pTarget.FoodWarModule.TotalFood
				response.ChallengerLst = append(response.ChallengerLst, challenge)
			}
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求抢夺次数以及恢复时间信息
func Hand_FoodWar_GetTime(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_GetFoodWarTime_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_GetTime Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_GetFoodWarTime_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动是否开启
	if player.FoodWarModule.IsActivityOpen() == false {
		gamelog.Error("Hand_FoodWar_GetTime Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	response.AttackTimes = player.FoodWarModule.AttackTimes
	response.RecoverTime = player.FoodWarModule.NextTime
	response.RevengeTimes = player.FoodWarModule.RevengeTimes
	response.BuyAttackTimes = player.FoodWarModule.BuyAttackTimes
	response.BuyRevengeTimes = player.FoodWarModule.BuyRevengeTimes
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求夺粮
func Hand_RobFood(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_RobFood_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_RobFood Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_RobFood_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动开启
	if player.FoodWarModule.IsActivityOpen() == false {
		gamelog.Error("Hand_RobFood Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	//! 检测攻击次数
	if player.FoodWarModule.AttackTimes <= 0 {
		response.RetCode = msg.RE_NOT_ENOUGH_ATTACK_TIMES
		gamelog.Error("Hand_RobFood Error: Attack times not enough")
		return
	}

	//! 扣除掠夺次数
	player.FoodWarModule.AttackTimes -= 1
	go player.FoodWarModule.DB_SaveAttackTimes()

	//! 获取目标玩家粮草信息
	targetFood := player.FoodWarModule.GetPlayerFoodInfo(req.TargetPlayerID)

	if req.IsWin == 1 {
		//! 计算抢夺粮草
		response.RobFood = (targetFood.TotalFood - targetFood.FixedFood) * gamedata.FoodWarRobBili / 1000

		player.FoodWarModule.FixedFood += response.RobFood / 3
		player.FoodWarModule.TotalFood = player.FoodWarModule.TotalFood + response.RobFood/3 + response.RobFood
		response.TotalFood = player.FoodWarModule.TotalFood

		go player.FoodWarModule.DB_SaveFood()

		//! 扣除目标流动粮草
		targetFood.TotalFood -= response.RobFood
		go targetFood.DB_SaveFood()

		//! 排行榜变动
		G_FoodWarRanker.SetRankItem(req.TargetPlayerID, targetFood.TotalFood)

		//! 排行榜变动
		response.Rank = G_FoodWarRanker.SetRankItem(player.playerid, player.FoodWarModule.TotalFood) + 1

		//! 给予胜利货币奖励
		player.RoleMoudle.AddMoney(gamedata.FoodWarVictoryMoneyID, gamedata.FoodWarVictoryMoneyNum)
		response.MoneyID = gamedata.FoodWarVictoryMoneyID
		response.MoneyNum = gamedata.FoodWarVictoryMoneyNum

		//! 复仇名单增加
		targetFood.RevengeLst = append(targetFood.RevengeLst, TRevengeInfo{player.playerid, response.RobFood})
		go targetFood.DB_AddRevengeLst(TRevengeInfo{player.playerid, response.RobFood})

	} else {
		//! 给予失败货币奖励
		player.RoleMoudle.AddMoney(gamedata.FoodWarFailedMoneyID, gamedata.FoodWarFailedMoneyNum)
		response.MoneyID = gamedata.FoodWarFailedMoneyID
		response.MoneyNum = gamedata.FoodWarFailedMoneyNum
	}

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
	response.RobFood = (player.FoodWarModule.TotalFood - player.FoodWarModule.FixedFood) * gamedata.FoodWarRobBili / 1000

	if len(g_SelectPlayers) < 5 {
		response.ChallengerLst = []msg.MSG_FoodWar_Challenger{}
	} else {
		for i := 0; i < 4; i++ {
			pTarget := GetSelectPlayer(SelectScoreTarget, 1000)
			if pTarget != nil && pTarget.playerid != 0 {
				var challenge msg.MSG_FoodWar_Challenger
				challenge.PlayerID = pTarget.playerid

				if challenge.PlayerID == player.playerid {
					i -= 1
					continue
				}

				simpleInfo := G_SimpleMgr.GetSimpleInfoByID(challenge.PlayerID)
				challenge.PlayerName = simpleInfo.Name
				challenge.Quality = simpleInfo.Quality
				challenge.HeroID = simpleInfo.HeroID
				challenge.Level = simpleInfo.Level
				challenge.FightValue = simpleInfo.FightValue
				challenge.CanRobFood = (pTarget.FoodWarModule.TotalFood - pTarget.FoodWarModule.FixedFood) * gamedata.FoodWarRobBili / 1000
				challenge.TotalFood = pTarget.FoodWarModule.TotalFood
				response.ChallengerLst = append(response.ChallengerLst, challenge)
			}
		}
	}
}

//! 玩家请求复仇
func Hand_FoodWar_Revenge(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_RevengeRob_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_Revenge Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_RevengeRob_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动开启
	if player.FoodWarModule.IsActivityOpen() == false {
		gamelog.Error("Hand_FoodWar_Revenge Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	//! 检测复仇次数
	if player.FoodWarModule.RevengeTimes <= 0 {
		gamelog.Error("Hand_FoodWar_Revenge Error: Revenge times not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_REVENGE_TIMES
		return
	}

	//! 扣除复仇次数
	player.FoodWarModule.RevengeTimes -= 1
	go player.FoodWarModule.DB_SaveRevengeTimes()

	revengeInfo := player.FoodWarModule.GetRevengeInfo(req.TargetPlayerID)
	if revengeInfo == nil {
		gamelog.Error("Hand_FoodWar_Revenge Error: Not find player in revenge list")
		return
	}

	if req.IsWin == 1 {
		response.RobFood = revengeInfo.RobFood
		player.FoodWarModule.TotalFood += response.RobFood
		response.FixFood = player.FoodWarModule.FixedFood
		go player.FoodWarModule.DB_SaveFood()

		//! 排行榜变动
		response.Rank = G_FoodWarRanker.SetRankItem(player.playerid, player.FoodWarModule.TotalFood) + 1

		//! 给予胜利货币奖励
		player.RoleMoudle.AddMoney(gamedata.FoodWarVictoryMoneyID, gamedata.FoodWarVictoryMoneyNum)
		response.MoneyID = gamedata.FoodWarVictoryMoneyID
		response.MoneyNum = gamedata.FoodWarVictoryMoneyNum

		pos := 0
		for i, v := range player.FoodWarModule.RevengeLst {
			if v.PlayerID == req.TargetPlayerID {
				pos = i
				go player.FoodWarModule.DB_RemoveRevengeLst(v)
				break
			}
		}

		if pos == 0 {
			player.FoodWarModule.RevengeLst = player.FoodWarModule.RevengeLst[1:]
		} else if (pos + 1) == len(player.FoodWarModule.RevengeLst) {
			player.FoodWarModule.RevengeLst = player.FoodWarModule.RevengeLst[:pos]
		} else {
			player.FoodWarModule.RevengeLst = append(player.FoodWarModule.RevengeLst[:pos], player.FoodWarModule.RevengeLst[pos+1:]...)
		}
	} else {
		//! 给予失败货币奖励
		player.RoleMoudle.AddMoney(gamedata.FoodWarFailedMoneyID, gamedata.FoodWarFailedMoneyNum)
		response.MoneyID = gamedata.FoodWarFailedMoneyID
		response.MoneyNum = gamedata.FoodWarFailedMoneyNum
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求自身状态
func Hand_FoodWar_GetStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_GetStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_GetStatus Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_GetStatus_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动开启
	if player.FoodWarModule.IsActivityOpen() == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		gamelog.Error("Hand_FoodWar_GetStatus Error: Activity not open")
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	if len(g_SelectPlayers) < 5 {
		response.ChallengerLst = []msg.MSG_FoodWar_Challenger{}
	} else {
		for i := 0; i < 4; i++ {
			pTarget := GetSelectPlayer(SelectScoreTarget, 1000)
			if pTarget != nil && pTarget.playerid != 0 {
				var challenge msg.MSG_FoodWar_Challenger
				challenge.PlayerID = pTarget.playerid

				if challenge.PlayerID == player.playerid {
					i -= 1
					continue
				}

				simpleInfo := G_SimpleMgr.GetSimpleInfoByID(challenge.PlayerID)
				challenge.PlayerName = simpleInfo.Name
				challenge.Quality = simpleInfo.Quality
				challenge.HeroID = simpleInfo.HeroID
				challenge.Level = simpleInfo.Level
				challenge.FightValue = simpleInfo.FightValue
				challenge.CanRobFood = (pTarget.FoodWarModule.TotalFood - pTarget.FoodWarModule.FixedFood) * gamedata.FoodWarRobBili / 1000
				challenge.TotalFood = pTarget.FoodWarModule.TotalFood
				response.ChallengerLst = append(response.ChallengerLst, challenge)
			}
		}
	}

	response.AttackTimes = player.FoodWarModule.AttackTimes
	response.RecoverTime = player.FoodWarModule.NextTime
	response.BuyAttackTimes = player.FoodWarModule.BuyAttackTimes
	response.AttackTimesMax = gamedata.FoodWarAttackTimes
	response.FixFood = player.FoodWarModule.FixedFood
	response.RecoverFood = gamedata.FoodWarTimeAddFood

	response.Rank = G_FoodWarRanker.SetRankItem(player.playerid, player.FoodWarModule.TotalFood) + 1
	response.TotalFood = player.FoodWarModule.TotalFood
	response.RobFood = (player.FoodWarModule.TotalFood - player.FoodWarModule.FixedFood) * gamedata.FoodWarRobBili / 1000
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求复仇状态
func Hand_FoodWar_RevengeStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_RevengeStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_RevengeStatus Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_RevengeStatus_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动开启
	if player.FoodWarModule.IsActivityOpen() == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		gamelog.Error("Hand_FoodWar_RevengeStatus Error: Activity not open")
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	for _, v := range player.FoodWarModule.RevengeLst {
		var challenge msg.MSG_FoodWar_Challenger
		challenge.PlayerID = v.PlayerID

		simpleInfo := G_SimpleMgr.GetSimpleInfoByID(challenge.PlayerID)
		if simpleInfo == nil {
			gamelog.Error("GetSimpleInfoByID Error: invalid playerid %v", challenge.PlayerID)
			continue
		}
		challenge.PlayerName = simpleInfo.Name
		challenge.Quality = simpleInfo.Quality
		challenge.HeroID = simpleInfo.HeroID
		challenge.Level = simpleInfo.Level
		challenge.FightValue = simpleInfo.FightValue
		foodModule := player.FoodWarModule.GetPlayerFoodInfo(challenge.PlayerID)
		if foodModule == nil {
			gamelog.Error("GetPlayerFoodInfo Error: invalid playerid %v", challenge.PlayerID)
			continue
		}

		challenge.CanRobFood = v.RobFood
		challenge.TotalFood = foodModule.TotalFood
		response.RevengeLst = append(response.RevengeLst, challenge)
	}

	response.RecoverTime = player.FoodWarModule.NextTime
	response.RevengeTimes = player.FoodWarModule.RevengeTimes
	response.BuyRevengeTimes = player.FoodWarModule.BuyRevengeTimes

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求获取粮草排行
func Hand_FoodWar_GetRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_GetRank_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_GetRank Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_GetRank_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动开启
	if player.FoodWarModule.IsActivityOpen() == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		gamelog.Error("Hand_FoodWar_GetRank Error: Activity not open")
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	for i := 0; i < G_FoodWarRanker.List.Len(); i++ {
		var challenge msg.MSG_FoodWar_Challenger
		challenge.PlayerID = G_FoodWarRanker.List[i].RankID

		simpleInfo := G_SimpleMgr.GetSimpleInfoByID(challenge.PlayerID)
		challenge.PlayerName = simpleInfo.Name
		challenge.Quality = simpleInfo.Quality
		challenge.HeroID = simpleInfo.HeroID
		challenge.Level = simpleInfo.Level
		challenge.FightValue = simpleInfo.FightValue
		challenge.TotalFood = G_FoodWarRanker.List[i].RankValue
		response.RankLst = append(response.RankLst, challenge)

		if len(response.RankLst) == 5 {
			break
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 请求购买次数
func Hand_FoodWar_BuyTimes(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_BuyTimes_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_BuyTimes Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_BuyTimes_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动开启
	if player.FoodWarModule.IsActivityOpen() == false {
		gamelog.Error("Hand_FoodWar_BuyTimes Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	if req.TimesType == 1 {
		resetTimes := gamedata.GetFuncVipValue(gamedata.FUNC_FOODWAR_ATTACK_TIMES, player.GetVipLevel())

		if player.FoodWarModule.BuyAttackTimes+req.Times > resetTimes {
			gamelog.Error("Hand_FoodWar_BuyTimes Error: Buy times limit")
			response.RetCode = msg.RE_NOT_ENOUGH_ATTACK_TIMES
			return
		}

		response.CostMoneyID = gamedata.FoodWarBuyTimesNeedMoneyID
		for i := player.FoodWarModule.BuyAttackTimes + 1; i <= player.FoodWarModule.BuyAttackTimes+req.Times; i++ {
			cost := gamedata.GetFuncTimeCost(gamedata.FUNC_FOODWAR_ATTACK_TIMES, i)
			response.CostMoneyNum += cost
		}

		//! 检查金钱是否足够
		if player.RoleMoudle.CheckMoneyEnough(gamedata.FoodWarBuyTimesNeedMoneyID, response.CostMoneyNum) == false {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			gamelog.Error("Hand_FoodWar_BuyTimes Error: Not enough money")
			return
		}

		//! 扣除金钱
		player.RoleMoudle.CostMoney(gamedata.FoodWarBuyTimesNeedMoneyID, response.CostMoneyNum)

		//! 增加次数
		player.FoodWarModule.BuyAttackTimes += req.Times
		player.FoodWarModule.AttackTimes += req.Times
		go player.FoodWarModule.DB_SaveBuyAttackTimes()

		response.RetCode = msg.RE_SUCCESS
	} else {
		resetTimes := gamedata.GetFuncVipValue(gamedata.FUNC_FOODWAR_REVENGE_TIMES, player.GetVipLevel())

		if player.FoodWarModule.BuyRevengeTimes+req.Times > resetTimes {
			gamelog.Error("Hand_FoodWar_BuyTimes Error: Buy times limit")
			response.RetCode = msg.RE_NOT_ENOUGH_REVENGE_TIMES
			return
		}

		response.CostMoneyID = gamedata.FoodWarBuyTimesNeedMoneyID
		for i := player.FoodWarModule.BuyRevengeTimes + 1; i <= player.FoodWarModule.BuyRevengeTimes+req.Times; i++ {
			cost := gamedata.GetFuncTimeCost(gamedata.FUNC_FOODWAR_REVENGE_TIMES, i)
			response.CostMoneyNum += cost
		}

		//! 检查金钱是否足够
		if player.RoleMoudle.CheckMoneyEnough(gamedata.FoodWarBuyTimesNeedMoneyID, response.CostMoneyNum) == false {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			gamelog.Error("Hand_FoodWar_BuyTimes Error: Not enough money")
			return
		}

		//! 扣除金钱
		player.RoleMoudle.CostMoney(gamedata.FoodWarBuyTimesNeedMoneyID, response.CostMoneyNum)

		//! 增加次数
		player.FoodWarModule.BuyRevengeTimes += req.Times
		player.FoodWarModule.RevengeTimes += req.Times
		go player.FoodWarModule.DB_SaveBuyRevengeTimes()

		response.RetCode = msg.RE_SUCCESS
	}
}

//! 请求获取粮草奖励
func Hand_FoodWar_RecvAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_GetFoodAward_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_RecvAward Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_GetFoodAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动开启
	if player.FoodWarModule.IsActivityOpen() == false {
		gamelog.Error("Hand_FoodWar_RecvAward Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	foodAward := gamedata.GetFoodWarAward(req.ID)
	if foodAward == nil {
		gamelog.Error("Hand_FoodWar_RecvAward Error: FoodWarAward get nil Req.ID: %d", req.ID)
		return
	}

	if player.FoodWarModule.TotalFood < foodAward.Target {
		gamelog.Error("Hand_FoodWar_RecvAward Error: Food not enough.ID: %d", req.ID)
		response.RetCode = msg.RE_NOT_ENOUGH_FOOD
		return
	}

	if player.FoodWarModule.AwardRecvLst.IsExist(req.ID) >= 0 {
		gamelog.Error("Hand_FoodWar_RecvAward Error: Aleady recv this award.ID: %d", req.ID)
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	awardItems := gamedata.GetItemsFromAwardID(foodAward.Award)
	player.BagMoudle.AddAwardItems(awardItems)

	player.FoodWarModule.AwardRecvLst.Add(req.ID)
	go player.FoodWarModule.DB_AddAwardRecvRecord(req.ID)
	response.Award = foodAward.Award
	response.RetCode = msg.RE_SUCCESS
}

//! 请求查询粮草奖励
func Hand_FoodWar_QueryAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FoodWar_QueryFoodAward_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FoodWar_QueryAward Error: unmarshal fail.")
		return
	}

	//! 定义返回
	var response msg.MSG_FoodWar_QueryFoodAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Retrun: %s", b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测活动开启
	if player.FoodWarModule.IsActivityOpen() == false {
		gamelog.Error("Hand_FoodWar_QueryAward Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	player.FoodWarModule.CheckTime()
	player.FoodWarModule.CheckReset()

	response.AwardLst = append(response.AwardLst, player.FoodWarModule.AwardRecvLst...)

	response.RetCode = msg.RE_SUCCESS
}

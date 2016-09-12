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
func Hand_GetLevelRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetLevelRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetLevelRank : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetLevelRank_Ack
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

	response.PlayerLst = []msg.MSG_PlayerInfo{}

	response.MyRank = -1
	for i := 0; i < len(G_LevelRanker.List); i++ {
		if G_LevelRanker.List[i].RankID <= 0 {
			break
		}

		if len(response.PlayerLst) >= G_LevelRanker.ShowNum {
			break
		}

		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(G_LevelRanker.List[i].RankID)
		if pSimpleInfo != nil {
			var info msg.MSG_PlayerInfo
			info.FightValue = pSimpleInfo.FightValue
			info.Level = pSimpleInfo.Level
			info.Name = pSimpleInfo.Name
			info.Quality = pSimpleInfo.Quality
			info.HeroID = pSimpleInfo.HeroID
			response.PlayerLst = append(response.PlayerLst, info)
		}

		if G_LevelRanker.List[i].RankID == req.PlayerID {
			response.MyRank = i + 1
		}
	}

	if response.MyRank < 0 {
		response.MyRank = G_LevelRanker.GetRankIndex(player.playerid, player.HeroMoudle.CurHeros[0].Level)
	}

	response.RetCode = msg.RE_SUCCESS
}

func Hand_GetFightRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetFightRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFightRank : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetFightRank_Ack
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

	response.PlayerLst = []msg.MSG_PlayerInfo{}

	response.MyRank = -1
	for i := 0; i < len(G_FightRanker.List); i++ {
		if G_FightRanker.List[i].RankID <= 0 {
			break
		}

		if len(response.PlayerLst) >= G_FightRanker.ShowNum {
			break
		}
		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(G_FightRanker.List[i].RankID)
		if pSimpleInfo != nil {
			var info msg.MSG_PlayerInfo
			info.FightValue = pSimpleInfo.FightValue
			info.Level = pSimpleInfo.Level
			info.Name = pSimpleInfo.Name
			info.Quality = pSimpleInfo.Quality
			info.HeroID = pSimpleInfo.HeroID
			response.PlayerLst = append(response.PlayerLst, info)
		}

		if G_FightRanker.List[i].RankID == req.PlayerID {
			response.MyRank = i + 1
		}
	}

	if response.MyRank < 0 {
		response.MyRank = G_FightRanker.GetRankIndex(player.playerid, int(player.GetFightValue()))
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求获取全服星数排行榜
func Hand_GetSanguowsRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetSanguows_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSanguows_Rank Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSanguows_Ack
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

	response.PlayerLst = []msg.MSG_SanguowsInfo{}

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}
	response.MyStar = player.SangokuMusouModule.HistoryStar
	response.MyRank = -1
	for i := 0; i < len(G_SgwsStarRanker.List); i++ {
		if G_SgwsStarRanker.List[i].RankID <= 0 {
			break
		}
		if len(response.PlayerLst) >= G_SgwsStarRanker.ShowNum {
			break
		}
		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(G_SgwsStarRanker.List[i].RankID)
		if pSimpleInfo != nil {
			var info msg.MSG_SanguowsInfo
			info.Name = pSimpleInfo.Name
			info.Star = G_SgwsStarRanker.List[i].RankValue
			info.FightValue = pSimpleInfo.FightValue
			info.HeroID = pSimpleInfo.HeroID
			info.Quality = pSimpleInfo.Quality
			response.PlayerLst = append(response.PlayerLst, info)
		}
	}

	response.MyRank = G_SgwsStarRanker.GetRankIndex(player.playerid, player.SangokuMusouModule.HistoryStar)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求竞技场排行榜信息
func Hand_GetArenaRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetArenaRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetArenaRank Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetArenaRank_Ack
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

	response.PlayerLst = []msg.MSG_ArenaInfo{}

	//! 检测功能是否开启
	isFuncOpen := gamedata.IsFuncOpen(gamedata.FUNC_ARENA, player.GetLevel(), player.GetVipLevel())
	if isFuncOpen == false {
		gamelog.Error("Hand_GetArenaRank Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	if player.ArenaModule.CurrentRank > 5000 {
		response.MyRank = -1
	} else {
		response.MyRank = player.ArenaModule.CurrentRank
	}

	for _, v := range G_Rank_List {
		var info msg.MSG_ArenaInfo
		info.PlayerID = v.PlayerID
		if v.IsRobot == false { //! 真人
			simpleInfo := G_SimpleMgr.GetSimpleInfoByID(info.PlayerID)
			info.FightValue = simpleInfo.FightValue
			info.Name = simpleInfo.Name
			info.Level = simpleInfo.Level
			info.HeroID = simpleInfo.HeroID
			info.Quality = simpleInfo.Quality
			response.PlayerLst = append(response.PlayerLst, info)
		} else { //! 机器人
			pRobotInfo := gamedata.GetRobot(v.PlayerID)
			info.FightValue = pRobotInfo.FightValue
			info.Name = pRobotInfo.Name
			info.Level = pRobotInfo.Level
			info.HeroID = pRobotInfo.Heros[0].HeroID
			info.Quality = pRobotInfo.Quality
			response.PlayerLst = append(response.PlayerLst, info)
		}

		if len(response.PlayerLst) >= 30 {
			break
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询排行榜
func Hand_GetRebelRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetRebelRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetRebelRank Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetRebelRank_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.ExploitRankLst = []msg.MSG_RebelRankInfo{}
	response.DamageRankLst = []msg.MSG_RebelRankInfo{}

	//! 获取功勋排行
	for i, v := range G_RebelExploitRanker.List {
		if G_RebelExploitRanker.List[i].RankID <= 0 {
			break
		}

		if len(response.ExploitRankLst) < G_RebelExploitRanker.ShowNum {
			playerInfo := G_SimpleMgr.GetSimpleInfoByID(v.RankID)
			if playerInfo != nil {
				var info msg.MSG_RebelRankInfo
				info.PlayerID = v.RankID
				info.HeroID = playerInfo.HeroID
				info.Name = playerInfo.Name
				info.Level = playerInfo.Level
				info.FightValue = playerInfo.FightValue
				info.Value = v.RankValue
				info.Quality = playerInfo.Quality
				response.ExploitRankLst = append(response.ExploitRankLst, info)
			}
		}
		if req.PlayerID == v.RankID {
			response.MyExploitRank = i + 1
		}
	}

	//! 获取伤害排行
	for i, v := range G_RebelDamageRanker.List {
		if G_RebelDamageRanker.List[i].RankID <= 0 {
			break
		}
		if len(response.DamageRankLst) < G_RebelDamageRanker.ShowNum {
			playerInfo := G_SimpleMgr.GetSimpleInfoByID(v.RankID)
			if playerInfo != nil {
				var info msg.MSG_RebelRankInfo
				info.PlayerID = v.RankID
				info.HeroID = playerInfo.HeroID
				info.Name = playerInfo.Name
				info.Level = playerInfo.Level
				info.FightValue = playerInfo.FightValue
				info.Value = v.RankValue
				info.Quality = playerInfo.Quality
				response.DamageRankLst = append(response.DamageRankLst, info)
			}
		}

		if req.PlayerID == v.RankID {
			response.MyDamageRank = i + 1
		}
	}

	if response.MyDamageRank == 0 {
		response.MyDamageRank = 5000
	}

	if response.MyExploitRank == 0 {
		response.MyExploitRank = 5000
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求公会等级排行榜
func Hand_GetGuildLevelRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetGuildLevelRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetGuildLevelRank Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetGuildLevelRank_Ack
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

	response.GuildList = []msg.MSG_GuildRankInfo{}

	for i := 0; i < len(G_GuildLevelRanker.List); i++ {
		if G_GuildLevelRanker.List[i].RankID <= 0 {
			break
		}
		if len(response.GuildList) >= G_GuildLevelRanker.ShowNum {
			break
		}
		pGuildInfo := GetGuildByID(G_GuildLevelRanker.List[i].RankID)
		if pGuildInfo != nil {
			var info msg.MSG_GuildRankInfo
			info.GuildID = pGuildInfo.GuildID
			info.Icon = pGuildInfo.Icon
			info.CurNum = len(pGuildInfo.MemberList)
			info.Level = pGuildInfo.Level
			info.MaxNum = 30
			info.Name = pGuildInfo.bossName
			info.GuildName = pGuildInfo.Name
			info.CopyChapter = pGuildInfo.HistoryPassChapter
			response.GuildList = append(response.GuildList, info)
		}

	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求公会副本排行榜
func Hand_GetGuildCopyRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetGuildCopyRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetGuildCopyRank Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetGuildCopyRank_Ack
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

	response.GuildList = []msg.MSG_GuildRankInfo{}

	for i := 0; i < len(G_GuildCopyRanker.List); i++ {
		if G_GuildCopyRanker.List[i].RankID <= 0 {
			break
		}
		if len(response.GuildList) >= G_GuildCopyRanker.ShowNum {
			break
		}
		pGuildInfo := GetGuildByID(G_GuildLevelRanker.List[i].RankID)
		if pGuildInfo != nil {
			var info msg.MSG_GuildRankInfo
			info.GuildID = pGuildInfo.GuildID
			info.Icon = pGuildInfo.Icon
			info.CurNum = len(pGuildInfo.MemberList)
			info.Level = pGuildInfo.Level
			info.MaxNum = 30
			info.Name = pGuildInfo.bossName
			info.GuildName = pGuildInfo.Name
			info.CopyChapter = pGuildInfo.HistoryPassChapter
			response.GuildList = append(response.GuildList, info)
		}

	}

	response.RetCode = msg.RE_SUCCESS
}

func Hand_GetWanderRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetWanderRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetWanderRank : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetWanderRank_Ack
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

	response.PlayerLst = []msg.MSG_PlayerInfo{}

	response.MyRank = -1
	for i := 0; i < len(G_WanderRanker.List); i++ {
		if G_WanderRanker.List[i].RankID <= 0 {
			break
		}

		if len(response.PlayerLst) >= G_WanderRanker.ShowNum {
			break
		}
		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(G_WanderRanker.List[i].RankID)
		if pSimpleInfo != nil {
			var info msg.MSG_PlayerInfo
			info.FightValue = pSimpleInfo.FightValue
			info.Level = pSimpleInfo.Level
			info.Name = pSimpleInfo.Name
			info.Quality = pSimpleInfo.Quality
			info.HeroID = pSimpleInfo.HeroID
			response.PlayerLst = append(response.PlayerLst, info)
		}

		if G_WanderRanker.List[i].RankID == req.PlayerID {
			response.MyRank = i + 1
		}
	}

	if response.MyRank < 0 {
		response.MyRank = G_WanderRanker.GetRankIndex(player.playerid, int(player.GetFightValue()))
	}

	response.RetCode = msg.RE_SUCCESS
}

func GetCampBatRankList(ranker *utility.TRanker) (ret []msg.MSG_PlayerInfo) {
	ranker.ForeachShow(
		func(rankID int32, rankVal int) {
			simpleInfo := G_SimpleMgr.GetSimpleInfoByID(rankID)
			if simpleInfo != nil {
				ret = append(ret, msg.MSG_PlayerInfo{
					PlayerID: simpleInfo.PlayerID,
					Name:     simpleInfo.Name,
					HeroID:   simpleInfo.HeroID,
					Quality:  simpleInfo.Quality,
					Level:    simpleInfo.Level,
					Camp:     simpleInfo.BatCamp,
					Value:    rankVal})
			}
		})
	return ret
}

func Hand_GetCampBatRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetCampBatRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetCampBatRank : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetCampBatRank_Ack
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

	response.PlayerLst = []msg.MSG_PlayerInfo{}
	response.MyRank = -1
	response.RetCode = msg.RE_SUCCESS

	switch req.RankType {
	case RT_TodayKill:
		response.PlayerLst = GetCampBatRankList(&G_CampBat_TodayKill)
		response.MyRank = G_CampBat_TodayKill.GetRankIndex(req.PlayerID, player.CamBattleModule.Kill)
	case RT_TodayDestroy:
		response.PlayerLst = GetCampBatRankList(&G_CampBat_TodayDestroy)
		response.MyRank = G_CampBat_TodayDestroy.GetRankIndex(req.PlayerID, player.CamBattleModule.Destroy)
	case RT_KillSum:
		response.PlayerLst = GetCampBatRankList(&G_CampBat_KillSum)
		response.MyRank = G_CampBat_KillSum.GetRankIndex(req.PlayerID, player.CamBattleModule.KillSum)
	case RT_DestroySum:
		response.PlayerLst = GetCampBatRankList(&G_CampBat_DestroySum)
		response.MyRank = G_CampBat_DestroySum.GetRankIndex(req.PlayerID, player.CamBattleModule.DestroySum)
	case RT_CampDestroy:
		response.PlayerLst = GetCampBatRankList(&G_CampBat_CampDestroy[player.CamBattleModule.BattleCamp-1])
		response.MyRank = G_CampBat_CampDestroy[player.CamBattleModule.BattleCamp-1].GetRankIndex(req.PlayerID, player.CamBattleModule.Destroy)
	case RT_CampKill:
		response.PlayerLst = GetCampBatRankList(&G_CampBat_CampKill[player.CamBattleModule.BattleCamp-1])
		response.MyRank = G_CampBat_CampKill[player.CamBattleModule.BattleCamp-1].GetRankIndex(req.PlayerID, player.CamBattleModule.Kill)
	}
}

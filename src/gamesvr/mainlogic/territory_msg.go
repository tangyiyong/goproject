package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"utility"
)

//! 玩家请求当前领地状态
func Hand_GetTerritoryStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetTerritoryStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetTerritoryStatus Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetTerritoryStatus_Ack
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

	//! 获取领地消息
	totalTimes := gamedata.GetFuncVipValue(gamedata.FUNC_SUPPRESS_TERRITORY, player.GetVipLevel())
	response.RetCode = msg.RE_SUCCESS
	response.SuppressRiotTimes = totalTimes - player.TerritoryModule.SuppressRiotTimes
	response.TotalPatrolTime = player.TerritoryModule.TotalPatrolTime
	response.RiotTime = gamedata.RiotTime
	response.TerritoryLst = []msg.MSG_TerritoryInfo{}
	for _, v := range player.TerritoryModule.TerritoryLst {
		var territory msg.MSG_TerritoryInfo
		territory.ID = v.ID
		territory.PatrolBeginTime = v.PatrolEndTime - int32(v.PatrolTime)
		territory.PatrolEndTime = v.PatrolEndTime
		territory.SkillLevel = v.SkillLevel
		territory.HeroID = v.HeroID
		territory.PatrolType = v.AwardTime

		territory.RiotInfo = []msg.MSG_TerritoryRiotData{}
		for _, b := range v.RiotInfo {
			var riotInfo msg.MSG_TerritoryRiotData
			riotInfo.BeginTime = b.BeginTime
			riotInfo.DealTime = b.DealTime
			riotInfo.HelperName = b.HelperName
			territory.RiotInfo = append(territory.RiotInfo, riotInfo)
		}

		territory.AwardItem = []msg.MSG_ItemData{}
		for _, n := range v.AwardItem {
			var award msg.MSG_ItemData
			award.ID = n.ItemID
			award.Num = n.ItemNum
			territory.AwardItem = append(territory.AwardItem, award)
		}

		response.TerritoryLst = append(response.TerritoryLst, territory)
	}
}

//! 玩家回馈挑战领地结果
func Hand_ChallengeTerritory(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	if false == utility.MsgDataCheck(buffer, G_XorCode) {
		//存在作弊的可能
		gamelog.Error("Hand_ChallengeTerritory : Message Data Check Error!!!!")
		return
	}
	var req msg.MSG_ChallengeTerritory_Req
	if json.Unmarshal(buffer[:len(buffer)-16], &req) != nil {
		gamelog.Error("Hand_ChallengeTerritory : Unmarshal error!!!!")
		return
	}

	//! 创建回复
	var response msg.MSG_ChallengeTerritory_Ack
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

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	//检查英雄数据是否一致
	if !player.CheckHeroData(req.HeroCkD) {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_MiningEvent_Monster : CheckHeroData Error!!!!")
		return
	}

	if gamedata.IsFuncOpen(gamedata.FUNC_TERRITORY, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 检测是否已被挑战
	isChallenged := player.TerritoryModule.IsChallenged(req.TerritoryID)
	if isChallenged == true {
		gamelog.Error("Hand_ChallengeTerritory Error: terrtiory is challenged. TerritoryID: %d", req.TerritoryID)
		response.RetCode = msg.RE_CHALLENGE_ALEADY_END
		return
	}

	//! 获取领地信息
	territoryInfo := gamedata.GetTerritoryData(req.TerritoryID)
	if territoryInfo == nil {
		gamelog.Error("GetTerritoryData Fail. TerritoryID: %d", req.TerritoryID)
		return
	}

	//! 记录挑战结果
	player.TerritoryModule.ChallengeTerritory(req.TerritoryID)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家置放武将到领地巡逻
func Hand_PatrolTerritory(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_PatrolTerritory_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_PatrolTerritory Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_PatrolTerritory_Ack
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

	//! 判断领地是否已经攻略
	isChallenged := player.TerritoryModule.IsChallenged(req.TerritoryID)
	if isChallenged == false {
		response.RetCode = msg.RE_NOT_CHALLANGE
		return
	}

	//! 判断领地是否已经有武将巡逻
	territory, _ := player.TerritoryModule.GetTerritory(req.TerritoryID)
	if territory.PatrolEndTime > utility.GetCurTime() {
		response.RetCode = msg.RE_ALEADY_HAVE_HERO
		return
	}

	//! 判断VIP等级是否足够
	awardTime, funcID := gamedata.GetTerritoryAwardType(req.AwardType)

	if gamedata.IsFuncOpen(funcID, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 判断精力是否足够
	patrolInfo := gamedata.GetPatrolTypeInfo(req.PatrolType)
	needAction := patrolInfo.ActionNum
	if player.RoleMoudle.CheckActionEnough(patrolInfo.ActionType, needAction) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 判断玩家是否拥有该英雄
	if false == player.IsHasHero(req.HeroID) {
		response.RetCode = msg.RE_NOT_HAVE_HERO
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(patrolInfo.ActionType, needAction)

	//! 开始巡逻
	player.TerritoryModule.PatrolTerritory(req.TerritoryID, req.HeroID, patrolInfo, awardTime)
	response.RetCode = msg.RE_SUCCESS
	response.PatrolBeginTime = utility.GetCurTime()
	territory, _ = player.TerritoryModule.GetTerritory(req.TerritoryID)
	for _, v := range territory.AwardItem {
		var award msg.MSG_ItemData
		award.ID = v.ItemID
		award.Num = v.ItemNum
		response.AwardItem = append(response.AwardItem, award)
	}

	for _, b := range territory.RiotInfo {
		var riotInfo msg.MSG_TerritoryRiotData
		riotInfo.BeginTime = b.BeginTime
		riotInfo.DealTime = b.DealTime
		riotInfo.HelperName = b.HelperName
		response.RiotInfo = append(response.RiotInfo, riotInfo)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_TERRITORY_HUNT, 1)
}

//! 玩家请求查询领地暴动信息
func Hand_GetTerritoryRiotInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetTerritoryRiot_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetTerritoryRiotInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetTerritoryRiot_Ack
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

	if gamedata.IsFuncOpen(gamedata.FUNC_TERRITORY, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	territory, _ := player.TerritoryModule.GetTerritory(req.TerritoryID)
	for _, v := range territory.RiotInfo {
		var riotInfo msg.MSG_TerritoryRiotData
		riotInfo.BeginTime = v.BeginTime
		riotInfo.DealTime = v.DealTime
		riotInfo.HelperName = v.HelperName
		response.RiotInfo = append(response.RiotInfo, riotInfo)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求收获领地奖励
func Hand_GetTerritoryAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetTerritoryAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetTerritoryAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetTerritoryAward_Ack
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

	if gamedata.IsFuncOpen(gamedata.FUNC_TERRITORY, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 判断是否巡逻结束
	territory, _ := player.TerritoryModule.GetTerritory(req.TerritoryID)
	if territory.PatrolEndTime > utility.GetCurTime() {
		//! 尚未结束
		response.RetCode = msg.RE_PATROL_NOT_END
		return
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_TERRITORY_PATROLTIME, territory.PatrolTime/3600)

	//! 获取奖励,返回成功
	player.TerritoryModule.GetTerritoryAward(req.TerritoryID)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求升级领地技能
func Hand_TerritorySkillLevelUp(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_TerritorySkillUp_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_TerritorySkillLevelUp Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_TerritorySkillUp_Ack
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

	if gamedata.IsFuncOpen(gamedata.FUNC_TERRITORY, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取领地信息
	territory, _ := player.TerritoryModule.GetTerritory(req.TerritoryID)
	if territory == nil {
		response.RetCode = msg.RE_NOT_CHALLANGE
		return
	}

	//! 判断领地等级是否已满
	if territory.SkillLevel >= 5 {
		response.RetCode = msg.RE_MAX_TERRITORY_SKILL_LEVEL
		return
	}

	//! 获取开启技能需求累积时间
	skillInfo := gamedata.GetTerritorySkillData(req.TerritoryID, territory.SkillLevel+1)
	if skillInfo.SkillOpenTime > player.TerritoryModule.TotalPatrolTime {
		//! 累积时间不足以升级该技能
		response.RetCode = msg.RE_NOT_ENOUGH_PATROL_TIME
		return
	}

	//! 检查金钱是否足够
	if player.RoleMoudle.CheckMoneyEnough(skillInfo.CostMoneyID, skillInfo.CostMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(skillInfo.CostMoneyID, skillInfo.CostMoneyNum)

	//! 升级领地技能
	player.TerritoryModule.TerritorySkillLevelUp(req.TerritoryID)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 请求一键领取领地奖励
func Hand_OneKeyTerritoryAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_OneKeyTerritoryAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_OneKeyTerritoryAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_OneKeyTerritoryAward_Ack
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

	response.AwardInfo = []msg.MSG_ItemData{}
	for _, v := range player.TerritoryModule.TerritoryLst {

		if v.PatrolEndTime > utility.GetCurTime() {
			//! 尚未结束
			continue
		}

		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_TERRITORY_PATROLTIME, v.PatrolTime/3600)

		//! 获取奖励
		player.TerritoryModule.GetTerritoryAward(v.ID)

		for _, n := range v.AwardItem {
			var award msg.MSG_ItemData
			award.ID = n.ItemID
			award.Num = n.ItemNum
			response.AwardInfo = append(response.AwardInfo, award)
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 请求查询好友状态
func Hand_GetFriendTerritoryStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetFriendTerritoryStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFriendTerritoryStatus Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetFriendTerritoryStatus_Ack
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

	if gamedata.IsFuncOpen(gamedata.FUNC_TERRITORY, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取好友状态
	for _, v := range player.FriendMoudle.FriendList {

		friendTerritory := player.TerritoryModule.GetFriendTerritory(v.PlayerID)
		if friendTerritory == nil {
			gamelog.Error("Get territory info fail. playerID: %v  friendID: %v", player.playerid, v.PlayerID)
			return
		}

		playerInfo := G_SimpleMgr.GetSimpleInfoByID(friendTerritory.PlayerID)

		var status msg.MSG_FriendTerritoryStatus
		status.PlayerID = friendTerritory.PlayerID
		status.Level = playerInfo.Level
		status.Quality = playerInfo.Quality
		for _, n := range friendTerritory.TerritoryLst {
			var territory msg.MSG_TerritoryInfo
			territory.ID = n.ID
			territory.PatrolEndTime = n.PatrolEndTime
			territory.SkillLevel = n.SkillLevel
			territory.HeroID = n.HeroID
			for _, b := range n.RiotInfo {
				var riotInfo msg.MSG_TerritoryRiotData
				riotInfo.BeginTime = b.BeginTime
				riotInfo.DealTime = b.DealTime
				riotInfo.HelperName = b.HelperName
				territory.RiotInfo = append(territory.RiotInfo, riotInfo)
			}
			status.TerritoryLst = append(status.TerritoryLst, territory)
		}
		status.LastLoginTime = G_SimpleMgr.GetPlayerLogoffTime(friendTerritory.PlayerID)
		response.FriendInfo = append(response.FriendInfo, status)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查看好友领地详情
func Hand_GetFriendTerritoryDetail(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetFriendTerritoryInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFriendTerritoryDetail Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetFriendTerritoryInfo_Ack
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

	if gamedata.IsFuncOpen(gamedata.FUNC_TERRITORY, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取好友领地模块信息
	friendTerritory := player.TerritoryModule.GetFriendTerritory(req.FriendID)
	if friendTerritory == nil {
		gamelog.Error("Get territory info fail. playerID: %v  friendID: %v", player.playerid, req.FriendID)
		return
	}

	territoryInfo, _ := friendTerritory.GetTerritory(req.TerritoryID)
	if territoryInfo == nil {
		response.RetCode = msg.RE_NOT_CHALLANGE
		return
	}

	for _, v := range territoryInfo.AwardItem {
		var awardData msg.MSG_ItemData
		awardData.ID = v.ItemID
		awardData.Num = v.ItemNum
		response.AwardInfo = append(response.AwardInfo, awardData)
	}

	territory, _ := player.TerritoryModule.GetTerritory(req.TerritoryID)
	for _, v := range territory.RiotInfo {
		var riotInfo msg.MSG_TerritoryRiotData
		riotInfo.BeginTime = v.BeginTime
		riotInfo.DealTime = v.DealTime
		riotInfo.HelperName = v.HelperName
		response.RiotInfo = append(response.RiotInfo, riotInfo)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 请求帮助好友镇压暴动
func Hand_SuppressRiot(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_HelpRiot_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SuppressRiot Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_HelpRiot_Ack
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

	if gamedata.IsFuncOpen(gamedata.FUNC_TERRITORY, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取好友模块信息
	friendTerritory := player.TerritoryModule.GetFriendTerritory(req.TargetID)
	if friendTerritory == nil {
		gamelog.Error("Get territory info fail. playerID: %v  friendID: %v", player.playerid, req.TargetID)
		return
	}

	friendPlayerInfo := G_SimpleMgr.GetSimpleInfoByID(friendTerritory.PlayerID)

	territoryInfo, index := friendTerritory.GetTerritory(req.TargetTerritoryID)
	if territoryInfo == nil {
		//! 未攻略
		response.RetCode = msg.RE_NOT_CHALLANGE
		return
	}

	//! 判断是否暴动
	isRiot := friendTerritory.IsRiot(req.TargetTerritoryID)
	if isRiot == false {
		response.RetCode = msg.RE_NOT_RIOT
		return
	}

	//! 检查当前是否还有镇压暴动次数
	totalTimes := gamedata.GetFuncVipValue(gamedata.FUNC_SUPPRESS_TERRITORY, player.GetVipLevel())
	if player.TerritoryModule.SuppressRiotTimes >= totalTimes {
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 处理暴动奖励
	award := gamedata.ST_ItemData{gamedata.SuppressRiotFriendAwardItem, gamedata.SuppressRiotFriendAwardNum}
	territoryInfo.AwardItem = append(territoryInfo.AwardItem, award)
	friendTerritory.DB_AddTerritoryAward(territoryInfo.ID, award)

	//! 设置领地信息暴动信息
	for i, n := range territoryInfo.RiotInfo {
		//! 判断暴动
		if utility.GetCurTime() >= n.BeginTime &&
			utility.GetCurTime() < n.BeginTime+int32(gamedata.RiotTime) &&
			n.IsRoit == true {
			territoryInfo.RiotInfo[i].IsRoit = false
			territoryInfo.RiotInfo[i].DealTime = utility.GetCurTime()
			territoryInfo.RiotInfo[i].HelperName = friendPlayerInfo.Name
			friendTerritory.DB_UpdateRiotInfo(index, i, territoryInfo.RiotInfo[i])
		}
	}

	//! 暴动次数加一
	player.TerritoryModule.SuppressRiotTimes += 1
	player.TerritoryModule.DB_UpdateRiotTimes()

	//! 自己获取元宝奖励
	response.ItemID = gamedata.SuppressRiotAwardItem
	response.ItemNum = gamedata.SuppressRiotAwardItemNum

	//! 发放奖励
	player.BagMoudle.AddAwardItem(response.ItemID, response.ItemNum)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"msg"
	"net/http"
	"strings"
	"time"
)

//! 请求查询某个玩家详细信息
//! 消息: /get_player_info
type MSG_GetPlayerInfo_Req struct {
	PlayerID       int32
	SessionKey     string
	TargetPlayerID int32
}

type MSG_GetPlayerInfo_Ack struct {
	RetCode  int
	HeroInfo THeroMoudle
}

//! 玩家请求查询公会状态
func Hand_GetGuildStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildStatus Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildStatus_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_SUCCESS
		return
	}

	//! 检测行动力恢复
	player.GuildModule.RecoverAction()

	response.ActionTimes = player.GuildModule.ActionTimes
	response.NextRecoverTime = player.GuildModule.ActionRecoverTime

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	if guild == nil {
		gamelog.Error("Hand_GetGuildCopyStatus Error: invalid guild %d", player.pSimpleInfo.GuildID)
		return
	}

	guild.CheckReset()

	//! 副本
	for _, v := range guild.CampLife {
		var campLife msg.MSG_CampLife
		campLife.CopyID = v.CopyID
		campLife.Life = v.Life
		response.CampLife = append(response.CampLife, campLife)
	}

	response.IsBack = guild.IsBack
	response.PassChapter = guild.PassChapter
	response.HistoryPassChapter = guild.HistoryPassChapter

	for _, v := range guild.CopyTreasure {
		var treasure msg.MSG_GuildCopyTreasure
		treasure.CopyID = v.CopyID
		treasure.Index = v.Index
		treasure.AwardID = v.AwardID
		treasure.PlayerName = v.PlayerName
		response.CopyTreasure = append(response.CopyTreasure, treasure)
	}

	response.AwardChapter = []msg.MSG_PassAwardChapter{}
	for _, v := range guild.AwardChapterLst {
		var awardChapter msg.MSG_PassAwardChapter
		awardChapter.CopyID = v.CopyID
		awardChapter.PassChapter = v.PassChapter
		awardChapter.PassTime = v.PassTime
		awardChapter.PlayerName = v.PlayerName
		response.AwardChapter = append(response.AwardChapter, awardChapter)
	}

	for _, v := range guild.CopyTreasure {
		if v.PlayerID == req.PlayerID {
			award := gamedata.GetGuildCampAwardInfo(v.AwardID)
			var mark msg.MSG_RecvCopyMark
			mark.Chapter = award.Chapter
			mark.CopyID = award.CopyID
			response.IsRecvCopyAward = append(response.IsRecvCopyAward, mark)
		}

	}

	//! 商店
	for _, v := range player.GuildModule.ShoppingLst {
		var goods msg.MSG_GuildGoods
		goods.ID = v.ID
		goods.Times = v.BuyTimes
		response.BuyLst = append(response.BuyLst, goods)
	}

	//! 祭天

	response.SacrificeStatus = player.GuildModule.SacrificeStatus
	response.SacrificeNum = guild.Sacrifice
	response.SacrificeSchedule = guild.SacrificeSchedule

	//! 获取进度奖励
	awardLst := gamedata.GetGuildSacrificeAwardFromLevel(guild.Level)

	response.RecvLst = [4]int{0, 0, 0, 0}

	for i, v := range awardLst {
		if player.GuildModule.SacrificeAwardLst.IsExist(v) >= 0 {
			response.RecvLst[i] = 1
		}
	}

	response.SkillLst = []msg.TGuildSkill{}
	if player.pSimpleInfo.GuildID == 0 {
		response.SkillLst = []msg.TGuildSkill{}
		response.RetCode = msg.RE_SUCCESS
	}

	for _, v := range player.GuildModule.SkillLst {
		var skill msg.TGuildSkill
		skill.SkillID = v.SkillID
		skill.Level = v.Level
		response.SkillLst = append(response.SkillLst, skill)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求创建公会
func Hand_CreateGuild(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_CreateNewGuild_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuild Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_CreateNewGuild_Ack
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

	if player.pSimpleInfo.GuildID != 0 {
		response.RetCode = msg.RE_ALEADY_HAVE_GUILD
		return
	}

	//! 检测玩家货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.CreateGuildMoneyID, gamedata.CreateGuildMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 检查公会名是否重复
	if GetGuildByName(req.Name) != nil {
		response.RetCode = msg.RE_GUILD_NAME_REPEAT
		gamelog.Error("Hand_CreateGuild Error: Repeat guild name")
		return
	}

	//! 扣除玩家货币
	player.RoleMoudle.CostMoney(gamedata.CreateGuildMoneyID, gamedata.CreateGuildMoneyNum)

	//! 设置玩家公会ID

	guild := CreateNewGuild(req.PlayerID, req.Name, req.Icon)
	G_SimpleMgr.Set_GuildID(player.playerid, guild.GuildID)

	//! 移除其他帮派收到目标玩家申请记录
	for i, v := range G_Guild_List {
		for _, n := range v.ApplyList {
			if n == player.playerid {
				go G_Guild_List[i].RemoveApplyList(player.playerid)
			}
		}
	}

	player.GuildModule.ApplyGuildList = IntLst{}
	go player.GuildModule.DB_ResetApplyList()

	response.RetCode = msg.RE_SUCCESS
	response.NewGuild.BossName = player.RoleMoudle.Name
	response.NewGuild.BossID = player.playerid
	response.NewGuild.CurExp = guild.CurExp
	response.NewGuild.GuildID = guild.GuildID
	response.NewGuild.Icon = guild.Icon
	response.NewGuild.Level = guild.Level
	response.NewGuild.MemberNum = len(guild.MemberList)
	response.NewGuild.Name = guild.Name
	response.NewGuild.Notice = guild.Notice

}

//! 玩家请求公会状态
func Hand_GetGuild(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildInfo_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuild Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildInfo_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测帮会重置
	player.GuildModule.CheckReset()

	if player.pSimpleInfo.GuildID == 0 {
		response.IsHaveGuild = false

		//! 获取前五名公会
		guildLst := player.GuildModule.GetGuildLst(0)
		for _, v := range guildLst {
			var guild msg.MSG_GuildInfo
			guild.GuildID = v.GuildID
			guild.Name = v.Name
			guild.Notice = v.Notice
			guild.Level = v.Level
			guild.Icon = v.Icon
			guild.CurExp = v.CurExp

			//! 获取会长名字
			boss := v.GetGuildLeader()
			if boss != nil && boss.PlayerID != 0 {
				bossInfo := G_SimpleMgr.GetSimpleInfoByID(boss.PlayerID)
				guild.BossName = bossInfo.Name
				guild.BossID = bossInfo.PlayerID
			}

			guild.MemberNum = len(v.MemberList)

			response.GuildLst = append(response.GuildLst, guild)
		}

	} else {
		response.IsHaveGuild = true

		//! 获取公会信息
		guildInfo := GetGuildByID(player.pSimpleInfo.GuildID)

		//! 检测重置
		guildInfo.CheckReset()
		var guild msg.MSG_GuildInfo
		guild.GuildID = player.pSimpleInfo.GuildID
		guild.Icon = guildInfo.Icon
		guild.CurExp = guildInfo.CurExp
		guild.Level = guildInfo.Level
		guild.Name = guildInfo.Name
		guild.Notice = guildInfo.Notice

		//! 获取会长名字
		boss := guildInfo.GetGuildLeader()
		if boss != nil && boss.PlayerID != 0 {
			bossInfo := G_SimpleMgr.GetSimpleInfoByID(boss.PlayerID)
			guild.BossName = bossInfo.Name
			guild.BossID = bossInfo.PlayerID
		}
		guild.MemberNum = len(guildInfo.MemberList)

		response.GuildLst = append(response.GuildLst, guild)
	}

	response.CopyEndTime = GetTodayTime() + int64(gamedata.GuildCopyBattleTimeEnd)
	response.RetCode = msg.RE_SUCCESS
}

//! 请求更多公会列表
func Hand_GetMoreGuild(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildLst_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetMoreGuild Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildLst_Ack
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

	guildLst := player.GuildModule.GetGuildLst(req.Index - 1)

	for _, v := range guildLst {
		var guild msg.MSG_GuildInfo
		guild.GuildID = v.GuildID
		guild.Icon = v.Icon
		guild.CurExp = v.CurExp
		guild.Level = v.Level
		guild.Name = v.Name
		guild.Notice = v.Notice
		response.GuildLst = append(response.GuildLst, guild)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求加入公会
func Hand_EnterGuild(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_EnterGuild_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetMoreGuild Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_EnterGuild_Ack
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

	//! 检测帮会重置
	player.GuildModule.CheckReset()

	//! 判断玩家是否拥有公会
	if player.pSimpleInfo.GuildID != 0 {
		gamelog.Error("Hand_EnterGuild Error: Player don't have guild")
		response.RetCode = msg.RE_ALEADY_HAVE_GUILD
		return
	}

	//! 判断重复申请
	if player.GuildModule.ApplyGuildList.IsExist(req.GuildID) >= 0 {
		gamelog.Error("Hand_EnterGuild Error: Repeat apply guild")
		response.RetCode = msg.RE_ALEADY_APPLY
		return
	}

	//! 判断是否距离离开公会24小时
	if time.Now().Unix()-player.GuildModule.ExitGuildTime <= 24*60*60 {
		gamelog.Error("Hand_EnterGuild Error: Exit guild time not enough 24 hours")
		response.RetCode = msg.RE_EXIT_GUILD_TIME_NOT_ENOUGH
		return
	}

	//! 判断公会ID是否存在
	guildInfo := GetGuildByID(req.GuildID)
	if guildInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 加入申请列表
	player.GuildModule.ApplyGuildList.Add(req.GuildID)
	guildInfo.ApplyList = append(guildInfo.ApplyList, player.playerid)

	go player.GuildModule.DB_AddApplyGuildList(req.GuildID)
	go DB_AddApplyList(req.GuildID, player.playerid)

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询申请公会列表
func Hand_GetApplyGuildList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetApplyGuildList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetApplyGuildList Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetApplyGuildList_Ack
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

	for _, v := range player.GuildModule.ApplyGuildList {
		guildInfo := GetGuildByID(v)
		if guildInfo == nil {
			continue
		}

		//! 检测重置
		guildInfo.CheckReset()
		var guild msg.MSG_GuildInfo
		guild.GuildID = guildInfo.GuildID
		guild.Icon = guildInfo.Icon
		guild.CurExp = guildInfo.CurExp
		guild.Level = guildInfo.Level
		guild.Name = guildInfo.Name
		guild.Notice = guildInfo.Notice

		//! 获取会长名字
		boss := guildInfo.GetGuildLeader()
		if boss != nil && boss.PlayerID != 0 {
			bossInfo := G_SimpleMgr.GetSimpleInfoByID(boss.PlayerID)
			guild.BossName = bossInfo.Name
			guild.BossID = bossInfo.PlayerID
		}
		guild.MemberNum = len(guildInfo.MemberList)

		response.GuildLst = append(response.GuildLst, guild)

	}
	response.RetCode = msg.RE_SUCCESS
}

//! 请求撤回公会申请
func Hand_CancellationGuildApply(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_CancellationGuildApply_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_CancellationGuildApply Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_CancellationGuildApply_Ack
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

	//! 移除该工会申请
	index := -1
	for i, v := range player.GuildModule.ApplyGuildList {
		if v == req.GuildID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Hand_CancellationGuildApply Error: Apply list not exist guild id: %d", req.GuildID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	go player.GuildModule.DB_RemoveApplyGuildList(req.GuildID)
	if index == 0 {
		player.GuildModule.ApplyGuildList = player.GuildModule.ApplyGuildList[1:]
	} else if (index + 1) == len(player.GuildModule.ApplyGuildList) {
		player.GuildModule.ApplyGuildList = player.GuildModule.ApplyGuildList[:index]
	} else {
		player.GuildModule.ApplyGuildList = append(player.GuildModule.ApplyGuildList[:index], player.GuildModule.ApplyGuildList[index+1:]...)
	}

	//! 删除对应公会申请名单
	guild := GetGuildByID(req.GuildID)

	index = -1
	for i, v := range guild.ApplyList {
		if v == player.playerid {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Hand_CancellationGuildApply Error: Apply list not exist guild id: %d", req.GuildID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	go DB_RemoveApplyList(req.GuildID, player.playerid)
	if index == 0 {
		guild.ApplyList = guild.ApplyList[1:]
	} else if (index + 1) == len(guild.ApplyList) {
		guild.ApplyList = guild.ApplyList[:index]
	} else {
		guild.ApplyList = append(guild.ApplyList[:index], guild.ApplyList[index+1:]...)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 请求搜索公会
func Hand_SearchGuild(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_SearchGuild_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_SearchGuild Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_SearchGuild_Ack
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

	response.RetCode = msg.RE_SUCCESS

	for _, v := range G_Guild_List {
		if strings.Contains(v.Name, req.GuildName) == true {
			//! 检测重置
			v.CheckReset()
			var guild msg.MSG_GuildInfo
			guild.GuildID = v.GuildID
			guild.Icon = v.Icon
			guild.CurExp = v.CurExp
			guild.Level = v.Level
			guild.Name = v.Name
			guild.Notice = v.Notice

			//! 获取会长名字
			boss := v.GetGuildLeader()
			if boss != nil && boss.PlayerID != 0 {
				bossInfo := G_SimpleMgr.GetSimpleInfoByID(boss.PlayerID)
				guild.BossName = bossInfo.Name
				guild.BossID = bossInfo.PlayerID
			}
			guild.MemberNum = len(v.MemberList)

			response.GuildLst = append(response.GuildLst, guild)
		}
	}

}

//! 请求查询申请公会成员列表
func Hand_GetApplyGuildMemberList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetApplyGuildMemberList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetApplyGuildMemberList Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetApplyGuildMemberList_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guildInfo := GetGuildByID(player.pSimpleInfo.GuildID)
	guildMemberInfo := guildInfo.GetGuildMember(player.playerid)

	if gamedata.HasPermission(guildMemberInfo.Pose, gamedata.Permission_Income) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	for _, v := range guildInfo.ApplyList {
		simpleInfo := G_SimpleMgr.GetSimpleInfoByID(v)
		if simpleInfo == nil {
			continue
		}

		targetPlayer := GetPlayerByID(v)

		var member msg.MSG_MemberInfo
		member.PlayerID = simpleInfo.PlayerID
		member.Name = simpleInfo.Name
		member.OfflineTime = simpleInfo.LogoffTime
		member.Quality = simpleInfo.Quality
		member.Level = simpleInfo.Level
		member.Role = 0
		member.FightValue = simpleInfo.FightValue
		member.IsOnline = simpleInfo.isOnline
		member.Contribution = targetPlayer.GuildModule.HistoryContribution
		response.MemberInfoLst = append(response.MemberInfoLst, member)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 查询公会成员列表
func Hand_GetGuildMemberList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildMemberList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildMemberList Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildMemberList_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	for _, v := range guild.MemberList {
		simpleInfo := G_SimpleMgr.GetSimpleInfoByID(v.PlayerID)
		if simpleInfo == nil {
			continue
		}

		var member msg.MSG_MemberInfo
		targetPlayer := GetPlayerByID(v.PlayerID)
		if targetPlayer == nil {
			targetPlayer = LoadPlayerFromDB(v.PlayerID)
		}

		member.PlayerID = v.PlayerID
		member.Name = simpleInfo.Name
		member.OfflineTime = simpleInfo.LogoffTime
		member.Quality = simpleInfo.Quality
		member.Level = simpleInfo.Level
		member.Role = v.Pose
		member.FightValue = simpleInfo.FightValue
		member.IsOnline = simpleInfo.isOnline
		member.Contribution = targetPlayer.GuildModule.HistoryContribution
		response.MemberLst = append(response.MemberLst, member)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 查询玩家详细信息
func Hand_GetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req MSG_GetPlayerInfo_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetPlayerInfo Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response MSG_GetPlayerInfo_Ack
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

	targetPlayer := GetPlayerByID(req.TargetPlayerID)
	if targetPlayer == nil {
		targetPlayer = LoadPlayerFromDB(req.TargetPlayerID)
		if targetPlayer == nil {
			response.RetCode = msg.RE_INVALID_PLAYERID
			return
		}
	}
	response.HeroInfo = targetPlayer.HeroMoudle
	response.RetCode = msg.RE_SUCCESS
}

//! 接受玩家入帮
func Hand_ApplicationThrough(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ApplyThrough_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ApplicationThrough Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_ApplyThrough_Ack
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

	//! 判断玩家是否拥有帮派
	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	targetPlayer := GetPlayerByID(req.TargetPlayerID)
	if targetPlayer == nil {
		targetPlayer = LoadPlayerFromDB(req.TargetPlayerID)
	}

	if targetPlayer.pSimpleInfo.GuildID != 0 {
		response.RetCode = msg.RE_ALEADY_HAVE_GUILD
		return
	}

	//! 判断玩家权限
	guildInfo := GetGuildByID(player.pSimpleInfo.GuildID)
	guildMemberInfo := guildInfo.GetGuildMember(player.playerid)

	if gamedata.HasPermission(guildMemberInfo.Pose, gamedata.Permission_Income) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	//! 判断目标玩家是否在申请列表
	isExist := false
	for _, v := range guildInfo.ApplyList {
		if v == req.TargetPlayerID {
			isExist = true
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_NOT_HAVE_APPLY
		return
	}

	//! 判断公会成员是否上限
	guilddata := gamedata.GetGuildBaseInfo(guildInfo.Level)
	if len(guildInfo.MemberList) >= guilddata.MemberLimit {
		response.RetCode = msg.RE_GUILD_MEMBER_MAX
		return
	}

	//! 增加新成员
	guildInfo.AddGuildMember(req.TargetPlayerID)
	G_SimpleMgr.DB_SetGuildID(req.TargetPlayerID, player.pSimpleInfo.GuildID)
	SendGuildChangeMsg(req.TargetPlayerID, player.pSimpleInfo.GuildID)

	//! 移除目标玩家申请列表
	targetPlayer.GuildModule.ApplyGuildList = []int{}

	hour := gamedata.GuildCopyBattleTimeBegin / 60 * 60
	min := (gamedata.GuildCopyBattleTimeBegin - hour*3600) / 60
	sec := gamedata.GuildCopyBattleTimeBegin - hour*3600 - min*60

	beginTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, min, sec, 0, time.Now().Location()).Unix()

	for {
		beginTime += int64(gamedata.GuildActionRecoverTime)
		if beginTime > time.Now().Unix() {
			break
		}
	}

	targetPlayer.GuildModule.ActionRecoverTime = beginTime
	targetPlayer.GuildModule.RecoverAction()

	targetPlayer.GuildModule.DB_CleanApplyList()

	//! 移除帮派收到目标玩家申请记录
	for i, v := range G_Guild_List {
		for _, n := range v.ApplyList {
			if n == req.TargetPlayerID {
				G_Guild_List[i].RemoveApplyList(req.TargetPlayerID)
			}
		}
	}

	//! 增加玩家公会技能加成
	if len(player.GuildModule.SkillLst) > 0 {
		for _, v := range player.GuildModule.SkillLst {
			for i := 0; i < v.Level; i++ {
				if v.SkillID == 15 {
					player.RoleMoudle.AddGuildSkillExpIncLevel()
				} else {
					player.HeroMoudle.AddGuildSkillProLevel(v.SkillID)
				}
			}
		}
	}

	//! 计入公会时间
	guildInfo.AddGuildEvent(targetPlayer.playerid, GuildEvent_AddMember, 0, 0)

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求退出公会
func Hand_ExitGuild(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_LeaveGuild_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ExitGuild Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_LeaveGuild_Ack
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

	//! 判断玩家是否拥有帮派
	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 删除公会成员
	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	if guild == nil {
		gamelog.Error("Hand_ExitGuild Error: Invalid guild id %v", player.pSimpleInfo.GuildID)
		return
	}

	//! 判断玩家是否为公会会长
	if guild.GetGuildLeader().PlayerID == player.playerid {
		if len(guild.MemberList) != 1 {
			//! 公会长不允许解散公会
			response.RetCode = msg.RE_GUILD_LEADER_CAN_NOT_EXIT
			return
		} else {
			//! 解散公会
			go RemoveGuild(guild.GuildID)
		}
	} else {
		guild.RemoveGuildMember(player.playerid)
	}

	//! 清空公会技能加成
	player.RoleMoudle.ClearGuildSkillExpIncLevel()
	player.HeroMoudle.ClearGuildSkillProLevel()

	G_SimpleMgr.Set_GuildID(player.playerid, 0)
	player.GuildModule.ActionRecoverTime = 0
	player.GuildModule.ExitGuildTime = time.Now().Unix()
	go player.GuildModule.DB_ExitGuild()

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求祭天状态
func Hand_GetSacrificeStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSacrificeStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetApplyGuildList Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetSacrificeStatus_Ack
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

	//! 检测帮会重置
	player.GuildModule.CheckReset()

	//! 检测会长弃坑
	player.GuildModule.CheckGuildLeader()

	//! 判断玩家是否拥有帮派
	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	if guild == nil {
		gamelog.Error("Hand_GetSacrificeStatus Error: invalid GuildID %v", player.pSimpleInfo.GuildID)
		return
	}

	//! 检测公会重置
	guild.CheckReset()

	response.SacrificeStatus = player.GuildModule.SacrificeStatus
	response.SacrificeNum = guild.Sacrifice
	response.SacrificeSchedule = guild.SacrificeSchedule

	//! 获取进度奖励
	awardLst := gamedata.GetGuildSacrificeAwardFromLevel(guild.Level)

	response.RecvLst = [4]int{0, 0, 0, 0}

	for i, v := range awardLst {
		if player.GuildModule.SacrificeAwardLst.IsExist(v) >= 0 {
			response.RecvLst[i] = 1
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求开始祭天
func Hand_GuildSacrifice(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GuildSacrifice_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GuildSacrifice Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GuildSacrifice_Ack
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

	//! 检测帮会重置
	player.GuildModule.CheckReset()

	//! 判断玩家是否拥有帮派
	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 判断玩家是否祭天
	if player.GuildModule.SacrificeStatus != 0 {
		response.RetCode = msg.RE_ALEADY_SACRIFICE
		return
	}

	//! 检测祭天次数是否已满
	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	guild.CheckReset()
	guildData := gamedata.GetGuildBaseInfo(guild.Level)

	if guild.Sacrifice >= guildData.SacrificeTimes {
		response.RetCode = msg.RE_NOT_ENOUGH_SACRIFICE_TIMES
		return
	}

	//! 获取祭天方式信息
	sacrificeData := gamedata.GetGuildSacrificeInfo(req.SacrificeID)

	//! 检测金钱是否足够
	if player.RoleMoudle.CheckMoneyEnough(sacrificeData.CostMoneyID, sacrificeData.CostMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 开始祭天
	player.RoleMoudle.CostMoney(sacrificeData.CostMoneyID, sacrificeData.CostMoneyNum)

	//! 检查暴击
	randValue := rand.New(rand.NewSource(time.Now().UnixNano()))

	value := randValue.Intn(1000)
	isCril := 100
	if value < gamedata.GuildSacrificeCrit {
		isCril = 150
	}

	//! 增加军团贡献
	player.GuildModule.AddContribution(sacrificeData.MoneyNum * isCril / 100)
	player.RoleMoudle.AddMoney(sacrificeData.MoneyID, sacrificeData.MoneyNum*isCril/100)

	//! 增加祭天进度
	guild.AddSacrifice(sacrificeData.Schedule)

	//! 增加军团经验
	guild.AddExp(sacrificeData.Exp * isCril / 100)

	//! 增加事件
	if isCril == 150 {
		guild.AddGuildEvent(player.playerid, GuildEvent_Sacrifice_Crit, sacrificeData.Exp*150/100, sacrificeData.ID)
	} else {
		guild.AddGuildEvent(player.playerid, GuildEvent_Sacrifice, sacrificeData.Exp, sacrificeData.ID)
	}

	response.MoneyID = sacrificeData.MoneyID
	response.MoneyNum = sacrificeData.MoneyNum
	response.CurExp = guild.CurExp
	response.GuildLevel = guild.Level
	response.SacrificeNum = guild.Sacrifice
	response.SacrificeSchedule = guild.SacrificeSchedule
	response.RetCode = msg.RE_SUCCESS

	player.GuildModule.SacrificeStatus = req.SacrificeID
	go player.GuildModule.DB_UpdateSacrifice()

}

//! 玩家请求领取祭天奖励
func Hand_GetSacrificeAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSacrificeAward_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetSacrificeAward Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetSacrificeAward_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	if guild == nil {
		gamelog.Error("Hand_GetSacrificeAward Error: invalid guildid %v", player.pSimpleInfo.GuildID)
		return
	}

	guild.CheckReset()

	//! 获取奖励静态信息
	awardData := gamedata.GetGuildSacrificeAwardInfo(req.ID)

	if guild.SacrificeSchedule < awardData.NeedSchedule {
		response.RetCode = msg.RE_NOT_ENOUGH_POINT
		return
	}

	//! 判断奖励ID是否合法
	if awardData.Level != guild.Level {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断是否已领取
	if player.GuildModule.SacrificeAwardLst.IsExist(req.ID) >= 0 {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 领取物品
	awardLst := gamedata.GetItemsFromAwardID(awardData.Award)
	player.BagMoudle.AddAwardItems(awardLst)

	//! 记录领取
	player.GuildModule.SacrificeAwardLst.Add(req.ID)
	go player.GuildModule.DB_AddSacrificeMark(req.ID)

	response.RetCode = msg.RE_SUCCESS
}

//! 查询公会商店购买信息
func Hand_GetGuildStoreInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryGuildStoreStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildStoreInfo Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_QueryGuildStoreStatus_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	guild.CheckReset()

	for _, v := range player.GuildModule.ShoppingLst {
		var goods msg.MSG_GuildGoods
		goods.ID = v.ID
		goods.Times = v.BuyTimes
		response.BuyLst = append(response.BuyLst, goods)
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 公会商店购买
func Hand_BuyGuildItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyGuildStoreItem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_BuyGuildItem Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_BuyGuildStoreItem_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	guild.CheckReset()

	//! 获取商品信息
	itemData := gamedata.GetGuildItemInfo(req.ID)

	//! 检测公会等级是否足够
	if itemData.NeedLevel > guild.Level {
		response.RetCode = msg.RE_NOT_ENOUGH_GUILD_LEVEL
		return
	}

	//! 检测限购次数
	goodsInfo := player.GuildModule.ShoppingLst.Get(req.ID)
	if goodsInfo == nil && req.Num > itemData.Limit {
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	} else if goodsInfo != nil && goodsInfo.BuyTimes+req.Num > itemData.Limit {
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 检测金钱是否足够
	if itemData.CostMoneyID1 != 0 {
		if player.RoleMoudle.CheckMoneyEnough(itemData.CostMoneyID1, itemData.CostMoneyNum1*req.Num) == false {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			return
		}
	}

	if itemData.CostMoneyID2 != 0 {
		if player.RoleMoudle.CheckMoneyEnough(itemData.CostMoneyID2, itemData.CostMoneyNum2*req.Num) == false {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			return
		}
	}

	//! 扣除金钱
	if itemData.CostMoneyID1 != 0 {
		player.RoleMoudle.CostMoney(itemData.CostMoneyID1, itemData.CostMoneyNum1*req.Num)
	}

	if itemData.CostMoneyID2 != 0 {
		player.RoleMoudle.CostMoney(itemData.CostMoneyID2, itemData.CostMoneyNum2*req.Num)
	}

	//! 发放物品
	player.BagMoudle.AddAwardItem(itemData.ItemID, itemData.ItemNum*req.Num)

	//! 记录购买次数
	if goodsInfo == nil {
		index := player.GuildModule.ShoppingLst.Add(req.ID, itemData.Type, req.Num)
		go player.GuildModule.DB_AddShoppingInfo(player.GuildModule.ShoppingLst[index])
	} else {
		goodsInfo.BuyTimes += req.Num
		player.GuildModule.DB_AddShoppingTimes(req.ID, goodsInfo.BuyTimes)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求公会副本状态
func Hand_GetGuildCopyStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildCopyStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildCopyStatus Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildCopyStatus_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_SUCCESS
		return
	}

	//! 检测行动力恢复
	player.GuildModule.RecoverAction()

	response.ActionTimes = player.GuildModule.ActionTimes
	response.NextRecoverTime = player.GuildModule.ActionRecoverTime

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	if guild == nil {
		gamelog.Error("Hand_GetGuildCopyStatus Error: invalid guild %d", player.pSimpleInfo.GuildID)
		return
	}

	guild.CheckReset()

	for _, v := range guild.CampLife {
		var campLife msg.MSG_CampLife
		campLife.CopyID = v.CopyID
		campLife.Life = v.Life
		response.CampLife = append(response.CampLife, campLife)
	}

	response.IsBack = guild.IsBack
	response.PassChapter = guild.PassChapter
	response.HistoryPassChapter = guild.HistoryPassChapter

	for _, v := range guild.CopyTreasure {
		var treasure msg.MSG_GuildCopyTreasure
		treasure.CopyID = v.CopyID
		treasure.Index = v.Index
		treasure.AwardID = v.AwardID
		treasure.PlayerName = v.PlayerName
		response.CopyTreasure = append(response.CopyTreasure, treasure)
	}

	response.AwardChapter = []msg.MSG_PassAwardChapter{}
	for _, v := range guild.AwardChapterLst {
		var awardChapter msg.MSG_PassAwardChapter
		awardChapter.CopyID = v.CopyID
		awardChapter.PassChapter = v.PassChapter
		awardChapter.PassTime = v.PassTime
		awardChapter.PlayerName = v.PlayerName
		response.AwardChapter = append(response.AwardChapter, awardChapter)
	}

	for _, v := range guild.CopyTreasure {
		if v.PlayerID == req.PlayerID {
			award := gamedata.GetGuildCampAwardInfo(v.AwardID)
			var mark msg.MSG_RecvCopyMark
			mark.Chapter = award.Chapter
			mark.CopyID = award.CopyID
			response.IsRecvCopyAward = append(response.IsRecvCopyAward, mark)
		}

	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求攻击公会副本
func Hand_AttackGuildCopy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_AttackGuildCopy_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_AttackGuildCopy Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_AttackGuildCopy_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 检测行动力
	if player.GuildModule.ActionTimes <= 0 {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 判断副本是否关闭
	endTime := GetTodayTime() + int64(gamedata.GuildCopyBattleTimeEnd)
	if time.Now().Unix() > endTime {
		response.RetCode = msg.RE_COPY_IS_LOCK
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	if guild.GetCopyLifeInfo(req.CopyID) <= 0 {
		response.RetCode = msg.RE_CAMP_IS_KILLED
		return
	}

	if req.Chapter > guild.PassChapter {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 扣除行动力
	player.GuildModule.ActionTimes -= 1
	go player.GuildModule.DB_UpdateCopyAction()

	//! 记录伤害与攻打次数
	memberInfo := guild.GetGuildMember(player.playerid)
	memberInfo.BattleTimes += 1
	if memberInfo.BattleDamage < req.Damage {
		memberInfo.BattleDamage = req.Damage
	}

	go guild.DB_UpdateDamageAndTimes(memberInfo.PlayerID, memberInfo.BattleTimes, memberInfo.BattleDamage)

	//! 扣除阵营血量
	isVictory, isKilled := guild.SubCampLife(req.CopyID, req.Damage, player.RoleMoudle.Name)
	for _, v := range guild.CampLife {
		var campLife msg.MSG_CampLife
		campLife.CopyID = v.CopyID
		campLife.Life = v.Life
		response.CampLife = append(response.CampLife, campLife)
	}

	if isVictory == true {
		//! 进入下一章副本

		guild.NextChapter()
		response.IsPass = true
	} else {
		response.IsPass = false
	}

	chapter := gamedata.GetGuildChapterInfo(req.Chapter)
	if isKilled == true {
		//! 击杀奖励经验
		guild.AddExp(chapter.Exp)
	}

	//! 奖励军团贡献
	random := rand.New(rand.NewSource(time.Now().Unix()))

	//! 活动贡献
	contribution := chapter.Contribution_min + random.Intn(chapter.Contribution_max-chapter.Contribution_min)
	player.RoleMoudle.AddMoney(chapter.MoneyID, contribution)

	copybaseInfo := gamedata.GetCopyBaseInfo(req.CopyID)

	//! 增加经验与金钱
	player.HeroMoudle.AddMainHeroExp(copybaseInfo.Experience * player.GetLevel())
	response.RetCode = msg.RE_SUCCESS

	response.AwardChapter = []msg.MSG_PassAwardChapter{}
	for _, v := range guild.AwardChapterLst {
		var awardChapter msg.MSG_PassAwardChapter
		awardChapter.CopyID = v.CopyID
		awardChapter.PassChapter = v.PassChapter
		awardChapter.PassTime = v.PassTime
		awardChapter.PlayerName = v.PlayerName
		response.AwardChapter = append(response.AwardChapter, awardChapter)
	}

	response.GuildLevel = guild.Level
	response.CurExp = guild.CurExp
}

//! 玩家请求领取公会副本宝藏奖励
func Hand_GetGuildCopyTreasure(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildCopyAward_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildCopyTreasure Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildCopyAward_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	//! 检查参数
	if req.Chapter > guild.PassChapter {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if req.Chapter == guild.PassChapter {
		//! 领取正在攻打的章节奖励
		if guild.GetCopyLifeInfo(req.CopyID) > 0 {
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}
	} else {
		member := guild.GetGuildMember(player.playerid)

		//! 检查通关时间与入帮时间
		for _, v := range guild.AwardChapterLst {
			if v.PassChapter == req.Chapter && v.CopyID == req.CopyID {
				if v.PassTime < member.EnterTime {
					response.RetCode = msg.RE_CANNOT_BE_RECV
					return
				}
			}
		}
	}

	if guild.IsRecvCampAward(player.playerid, req.CopyID, req.Chapter) == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 获取已经领取奖励ID
	awardIDLst := guild.GetAleadyRecvAwardIDLst(req.Chapter, req.CopyID)

	//! 随机奖励
	award := gamedata.RandGuildCampAward(req.Chapter, req.CopyID, awardIDLst)
	if award == nil {
		response.RetCode = msg.RE_NOT_HAVE_TREASURE
		return
	}

	response.ItemID = award.ItemID
	response.ItemNum = award.ItemNum
	response.AwardID = award.ID

	//! 发放奖励
	player.BagMoudle.AddAwardItem(response.ItemID, response.ItemNum)

	//! 记录发放
	guild.PlayerRecvAward(player.playerid, player.RoleMoudle.Name, req.CopyID, req.ID, award.ID)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询章节奖励领取列表
func Hand_GetGuildChapterRecvLst(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildChapterAwardStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildChapterRecvLst Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildChapterAwardStatus_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RecvLst = []int{}
		response.RetCode = msg.RE_SUCCESS
		return
	}

	response.RecvLst = player.GuildModule.CopyAwardMark
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求一键领取所有章节奖励
func Hand_GetAllGuildChapterAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildChapterAwardAll_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetAllGuildChapterAward Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildChapterAwardAll_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	for i := 1; i < guild.PassChapter; i++ {
		if player.GuildModule.CopyAwardMark.IsExist(i) >= 0 {
			continue
		}

		chapter := gamedata.GetGuildChapterInfo(i)
		awardLst := gamedata.GetItemsFromAwardID(chapter.Award)
		player.BagMoudle.AddAwardItems(awardLst)

		for _, v := range awardLst {
			response.Award = append(response.Award, msg.MSG_ItemData{v.ItemID, v.ItemNum})
		}

		player.GuildModule.CopyAwardMark.Add(i)
		response.RecvChapter = append(response.RecvChapter, i)
		go player.GuildModule.DB_AddChapterAwardRecord(i)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求领取章节奖励
func Hand_GetGuildChapterAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildChapterAward_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildChapterAward Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildChapterAward_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	if player.GuildModule.CopyAwardMark.IsExist(req.Chapter) >= 0 {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	if req.Chapter > guild.PassChapter {
		response.RetCode = msg.RE_CANNOT_BE_RECV
		return
	}

	chapter := gamedata.GetGuildChapterInfo(req.Chapter)
	awardLst := gamedata.GetItemsFromAwardID(chapter.Award)
	player.BagMoudle.AddAwardItems(awardLst)

	player.GuildModule.CopyAwardMark.Add(req.Chapter)
	go player.GuildModule.DB_AddChapterAwardRecord(req.Chapter)

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求修改帮派信息
func Hand_UpdateGuildInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_UpdateGuildInfo_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpdateGuildInfo Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_UpdateGuildInfo_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	memberInfo := guild.GetGuildMember(player.playerid)

	//! 判断权限
	if gamedata.HasPermission(memberInfo.Pose, gamedata.Permission_UpdateNotice) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	guild.Notice = req.Notice
	guild.Declaration = req.Declaration
	guild.Icon = req.Icon
	go guild.DB_UpdateGuildInfo()

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求修改公会副本回退状态
func Hand_UpdateGuildChapterBackStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_UpdateGuildBackStatus_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpdateGuildChapterBackStatus Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_UpdateGuildBackStatus_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 判断权限
	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	memberInfo := guild.GetGuildMember(player.playerid)

	if gamedata.HasPermission(memberInfo.Pose, gamedata.Permission_ResetCopy) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	if req.IsBack == 0 {
		guild.IsBack = false
	} else {
		guild.IsBack = true
	}

	go guild.DB_UpdateGuildBackStatus()

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求修改公会名称
func Hand_UpdateGuildName(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_UpdateGuildName_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpdateGuildInfo Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_UpdateGuildName_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 判断权限
	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	memberInfo := guild.GetGuildMember(player.playerid)

	if gamedata.HasPermission(memberInfo.Pose, gamedata.Permission_UpdateNotice) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	//! 判断金钱是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.UpdateGuildNameMoneyID, gamedata.UpdateGuildNameMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(gamedata.UpdateGuildNameMoneyID, gamedata.UpdateGuildNameMoneyNum)

	//! 修改公会名
	guild.Name = req.Name
	go guild.DB_UpdateGuildName()

	response.RetCode = msg.RE_SUCCESS
}

//! 踢出公会成员
func Hand_KickGuildMember(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_KickGuildMember_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_KickGuildMember Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_KickGuildMember_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 检测操作权限
	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	memberInfo := guild.GetGuildMember(player.playerid)
	if gamedata.HasPermission(memberInfo.Pose, gamedata.Permission_Kick) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	guild.RemoveGuildMember(req.KickPlayerID)

	//! 修改身份参数
	targetPlayer := GetPlayerByID(req.KickPlayerID)
	if targetPlayer != nil {
		G_SimpleMgr.Set_GuildID(req.KickPlayerID, 0)
		targetPlayer.GuildModule.ActionRecoverTime = 0
	}
	player.GuildModule.DB_KickPlayer(req.KickPlayerID)

	guild.AddGuildEvent(targetPlayer.playerid, GuildEvent_ExpelMember, 0, 0)

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求使用公会留言板
func Hand_UseGuildMsgBoard(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %v", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_WriteGuildMsgBoard_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UseGuildMsgBoard Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_WriteGuildMsgBoard_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 检查字数
	if len(req.Message) > 256 {
		response.RetCode = msg.RE_MESSAGE_TOO_LONG
		return
	}

	//! 生成留言
	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	guild.AddMsgBoard(req.PlayerID, req.Message)

}

//! 管理删除留言板记录
func Hand_RemoveGuildMsgBoard(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %v", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RemoveGuildMsgBoard_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_RemoveGuildMsgBoard Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_RemoveGuildMsgBoard_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	//! 删除留言
	guild.RemoveMsgBoard(req.PlayerID, req.TargetTime)
}

//! 查询公会留言板信息
func Hand_QueryGuildMsgBoard(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryGuildMsgBoard_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryGuildMsgBoard Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_QueryGuildMsgBoard_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	for _, v := range guild.MsgBoard {
		var message msg.MSG_GuildBoard
		message.PlayerID = v.PlayerID

		targetPlayer := G_SimpleMgr.GetSimpleInfoByID(v.PlayerID)
		if targetPlayer == nil {
			targetPlayer := GetPlayerByID(v.PlayerID)
			message.PlayerName = targetPlayer.RoleMoudle.Name
		} else {
			message.PlayerName = targetPlayer.Name
		}

		message.Message = v.Message
		message.Time = v.Time
		response.MsgLst = append(response.MsgLst, message)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 查询公会副本排行榜
func Hand_QueryGuildCopyRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryGuildCopyRank_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryGuildCopyRank Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_QueryGuildCopyRank_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	guild.SortDamage()

	for _, v := range guild.MemberList {
		var rankInfo msg.MSG_GuildCopyRank
		rankInfo.PlayerID = v.PlayerID
		rankInfo.BattleTimes = v.BattleTimes
		rankInfo.Damage = v.BattleDamage

		playerInfo := G_SimpleMgr.GetSimpleInfoByID(v.PlayerID)
		if playerInfo == nil {
			playerInfo := GetPlayerByID(v.PlayerID)
			rankInfo.PlayerName = playerInfo.RoleMoudle.Name
		} else {
			rankInfo.PlayerName = playerInfo.Name
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 查询副本奖励领取情况
func Hand_QueryGuildCopyTreasure(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryGuildCopyTreasure_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryGuildCopyTreasure Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_QueryGuildCopyTreasure_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	for _, v := range guild.CopyTreasure {
		awardInfo := gamedata.GetGuildCampAwardInfo(v.AwardID)
		if awardInfo.Chapter == req.Chapter {
			var copyTreasure msg.MSG_GuildCopyTreasure
			copyTreasure.CopyID = v.CopyID
			copyTreasure.Index = v.Index
			copyTreasure.AwardID = v.AwardID
			copyTreasure.PlayerName = v.PlayerName
			response.CopyTreasure = append(response.CopyTreasure, copyTreasure)
		}

	}

	response.RetCode = msg.RE_SUCCESS
}

//! 研究公会技能
func Hnad_ResearchGuildSkill(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ResearchGuildSkill_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hnad_ResearchGuildSkill Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_ResearchGuildSkill_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 判断权限
	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	memberInfo := guild.GetGuildMember(player.playerid)

	if gamedata.HasPermission(memberInfo.Pose, gamedata.Permission_Research) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	//! 判断公会等级
	guildSkillLimit := 0
	if req.SkillID == 15 {
		guildSkillLimit = gamedata.GetGuildExpIncSKillLimit(guild.Level)
	} else {
		guildSkillLimit = gamedata.GetGuildSkillLimit(guild.Level, req.SkillID)
	}

	//! 获取技能等级
	skillLevel := guild.GetGuildSkillLevel(req.SkillID)
	if skillLevel+1 > guildSkillLimit {
		response.RetCode = msg.RE_GUILD_SKILL_LIMIT
		return
	}

	//! 获取需求经验
	needExp := 0
	if req.SkillID == 15 {
		needExp = gamedata.GetGuildSkillExpNeedExp(skillLevel)
	} else {
		needExp = gamedata.GetGuildSkillNeedExp(skillLevel+1, req.SkillID)
	}

	if guild.CurExp < needExp {
		response.RetCode = msg.RE_NOT_ENOUGH_GUILD_EXP
		return
	}

	guild.AddGuildSkillLevel(req.SkillID, needExp)

	response.RetCode = msg.RE_SUCCESS
}

//! 学习公会技能
func Hand_StudyGuildSkill(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_StudyGuildSkill_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hnad_ResearchGuildSkill Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_StudyGuildSkill_Ack
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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 获取公会技能等级上限
	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	playerSkillLevel := player.GuildModule.GetPlayerGuildSKillLevel(req.SkillID)
	if playerSkillLevel+1 > guild.GetGuildSkillLevel(req.SkillID) {
		response.RetCode = msg.RE_GUILD_SKILL_LIMIT
		return
	}

	//! 检查金钱是否足够
	moneyID, moneyNum := 0, 0
	if req.SkillID == 15 {
		moneyID, moneyNum = gamedata.GetGuildSkillExpNeedMoney(playerSkillLevel + 1)
	} else {
		moneyID, moneyNum = gamedata.GetGuildSkillNeedMoney(playerSkillLevel+1, req.SkillID)
	}

	if player.RoleMoudle.CheckMoneyEnough(moneyID, moneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	player.GuildModule.AddPlayerGuildSkillLevel(req.SkillID)

	//! 增加属性
	if req.SkillID != 15 {
		player.HeroMoudle.AddGuildSkillProLevel(req.SkillID)
	} else {
		player.RoleMoudle.AddGuildSkillExpIncLevel()
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 获取公会技能信息
func Hand_GetGuildSkillInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildSkillInfo_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildSkillInfo Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildSkillInfo_Ack

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

	response.SkillLst = []msg.TGuildSkill{}
	if player.pSimpleInfo.GuildID == 0 {
		response.SkillLst = []msg.TGuildSkill{}
		response.RetCode = msg.RE_SUCCESS
		return
	}

	for _, v := range player.GuildModule.SkillLst {
		var skill msg.TGuildSkill
		skill.SkillID = v.SkillID
		skill.Level = v.Level
		response.SkillLst = append(response.SkillLst, skill)
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 获取公会技能信息
func Hand_GetGuildSkillResearchInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildSkillResearch_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildSkillResearchInfo Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildSkillResearch_Ack

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

	if player.pSimpleInfo.GuildID == 0 {
		response.SkillLst = []msg.TGuildSkill{}
		response.RetCode = msg.RE_SUCCESS
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	for _, v := range guild.SkillLst {
		var skill msg.TGuildSkill
		skill.SkillID = v.SkillID
		skill.Level = v.Level
		response.SkillLst = append(response.SkillLst, skill)
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 获取公会事件
func Hand_GetGuildLog(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetGuildLog_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuildLog Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetGuildLog_Ack

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

	if player.pSimpleInfo.GuildID == 0 {
		response.LogLst = []msg.GuildEvent{}
		response.RetCode = msg.RE_SUCCESS
		return
	}

	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	for _, v := range guild.EventLst {
		var log msg.GuildEvent
		log.Action = v.Action
		log.Type = v.Type
		log.Value = v.Value
		log.Time = v.Time
		log.PlayerID = v.PlayerID

		playerInfo := G_SimpleMgr.GetSimpleInfoByID(v.PlayerID)
		log.PlayerName = playerInfo.Name

		response.LogLst = append(response.LogLst, log)
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 修改公会职位
func Hand_ChangeGuildMemberPose(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ChangeGuildRole_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ChangeGuildMemberPose Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_ChangeGuildRole_Ack

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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 检查权限
	guild := GetGuildByID(player.pSimpleInfo.GuildID)

	memberInfo := guild.GetGuildMember(player.playerid)

	if gamedata.HasPermission(memberInfo.Pose, gamedata.Permission_Change) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	//! 获取目标玩家

	targetmemeberInfo := guild.GetGuildMember(req.TargetPlayerID)
	if memberInfo == nil || targetmemeberInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if targetmemeberInfo.Pose < memberInfo.Pose {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if guild.GetGuildLeader().PlayerID == player.playerid {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测权限
	if memberInfo.Pose > targetmemeberInfo.Pose {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	//! 会长或者副会长禅让
	limitNum := gamedata.GetGuildRoleLimit(req.Role)

	curNum := guild.GetPoseNumber(req.Role)

	//! 角色相同且人数不足
	if memberInfo.Pose == req.Role && memberInfo.Pose == Pose_Boss {

		//! 赋予职位
		targetmemeberInfo.Pose = memberInfo.Pose
		guild.UpdateGuildMemeber(req.TargetPlayerID, targetmemeberInfo.Pose, targetmemeberInfo.Contribute)

		//! 解除自身职位
		memberInfo.Pose = Pose_Member
		guild.UpdateGuildMemeber(player.playerid, Pose_Member, memberInfo.Contribute)

	} else {

		//gamelog.Info("curNum: %d  limitNum: %d", curNum, limitNum)
		if curNum < limitNum {
			guild.UpdateGuildMemeber(req.TargetPlayerID, targetmemeberInfo.Pose, targetmemeberInfo.Contribute)
		} else {
			response.RetCode = msg.RE_GUILD_MEMBER_MAX
			return
		}

	}

	guild.AddGuildEvent(req.TargetPlayerID, GuildEvent_ChangePose, req.Role, 0)

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求升级工会(暂时不用)
func Hand_GuildLevelUp(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GuildLevelUp_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GuildLevelUp Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GuildLevelUp_Ack

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

	if player.pSimpleInfo.GuildID == 0 {
		response.RetCode = msg.RE_HAVE_NOT_GUILD
		return
	}

	//! 检查权限
	guild := GetGuildByID(player.pSimpleInfo.GuildID)
	memberInfo := guild.GetGuildMember(player.playerid)
	if gamedata.HasPermission(memberInfo.Pose, gamedata.Permission_UpdateGuild) == false {
		response.RetCode = msg.RE_NOT_HAVE_PERMISSION
		return
	}

	//! 检测经验是否足够
	guildData := gamedata.GetGuildBaseInfo(guild.Level + 1)
	if guildData == nil {
		response.RetCode = msg.RE_ALREADY_MAX_LEVEL
		return
	}

	if guildData.NeedExp > guild.CurExp {
		response.RetCode = msg.RE_NOT_ENOUGH_GUILD_EXP
		return
	}

	guild.LevelUp()

	response.RetCode = msg.RE_SUCCESS
}

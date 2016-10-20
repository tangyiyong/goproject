package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"utility"
)

//! 玩家请求获取叛军信息
func Hand_GetRebelInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetRebelInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetRebelInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetRebelInfo_Ack
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

	player.RebelModule.CheckEscapeTime()
	player.RebelModule.CheckReset()

	//! 获取自身发现叛军
	response.InfoLst = []msg.MSG_RebelInfo{}
	if player.RebelModule.RebelID != 0 {
		var rebelInfo msg.MSG_RebelInfo
		rebelInfo.PlayerID = player.playerid
		rebelInfo.RebelID = player.RebelModule.RebelID
		rebelInfo.Level = player.RebelModule.GetRebelLevel()
		rebelInfo.FindName = player.RoleMoudle.Name
		rebelInfo.CurLife = player.RebelModule.CurLife

		if player.RebelModule.EscapeTime-utility.GetCurTime() < 0 {
			rebelInfo.EscapeTime = 0
		} else {
			rebelInfo.EscapeTime = player.RebelModule.EscapeTime - utility.GetCurTime()
		}

		rebelInfo.IsShare = player.RebelModule.IsShare

		response.InfoLst = append(response.InfoLst, rebelInfo)
	}

	//! 获取好友发现的叛军
	for _, v := range player.FriendMoudle.FriendList {
		rebelModulePtr, playerName := player.RebelModule.GetPlayerRebelPtr(v.PlayerID)
		rebelModulePtr.CheckEscapeTime()
		if rebelModulePtr.RebelID != 0 && rebelModulePtr.IsShare == true {
			var rebelInfo msg.MSG_RebelInfo
			rebelInfo.PlayerID = v.PlayerID
			rebelInfo.RebelID = rebelModulePtr.RebelID
			rebelInfo.Level = rebelModulePtr.GetRebelLevel()
			rebelInfo.FindName = playerName
			rebelInfo.CurLife = rebelModulePtr.CurLife
			rebelInfo.EscapeTime = utility.GetCurTime() - rebelInfo.EscapeTime
			rebelInfo.IsShare = rebelModulePtr.IsShare

			response.InfoLst = append(response.InfoLst, rebelInfo)
		}
	}

	//! 获取功勋排行
	for i, v := range G_RebelExploitRanker.List {
		if v.RankID == player.playerid {
			response.ExploitRank = i + 1
			break
		}
	}

	//! 获取伤害排行
	for i, v := range G_RebelDamageRanker.List {
		if v.RankID == player.playerid {
			response.DamageRank = i + 1
			break
		}
	}

	response.Exploit = player.RebelModule.Exploit
	response.TopDamage = player.RebelModule.Damage
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求攻击叛军
func Hand_AttackRebel(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! MD5消息验证
	if false == utility.MsgDataCheck(buffer, G_XorCode) {
		//存在作弊的可能
		gamelog.Error("Hand_AttackRebel : Message Data Check Error!!!!")
		return
	}
	var req msg.MSG_Attack_Rebel_Req
	if json.Unmarshal(buffer[:len(buffer)-16], &req) != nil {
		gamelog.Error("Hand_AttackRebel : Unmarshal error!!!!")
		return
	}

	//! 创建回复
	var response msg.MSG_Attack_Rebel_Ack
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

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	player.RebelModule.CheckEscapeTime()
	player.RebelModule.CheckReset()

	//! 获取叛军信息
	rebelModulePtr, _ := player.RebelModule.GetPlayerRebelPtr(req.TargetPlayerID)
	if rebelModulePtr.RebelID == 0 {
		response.RetCode = msg.RE_NOT_FIND_REBEL //! 没有发现叛军
		return
	}

	rebelModulePtr.CheckEscapeTime()

	//! 检查叛军逃走时间
	if rebelModulePtr.EscapeTime <= utility.GetCurTime() {
		response.RetCode = msg.RE_REBEL_ALEADY_ESCAPE
		return
	}

	//! 获取叛军血量
	if rebelModulePtr.CurLife <= 0 {
		response.RetCode = msg.RE_REBEL_ALEADY_KILL //! 叛军已被击杀
		return
	}

	//! 检测检查行动力是否足够
	needActionNum := req.AttackType
	if req.AttackType == 1 {
		needActionNum = gamedata.NormalAttackRebelNeedActionNum
	} else if req.AttackType == 2 {
		needActionNum = gamedata.SeniorAttackRebelNeedActionNum
	} else {
		gamelog.Error("Hand_AttackRebel error: invalid attacktype: %d  PlayerID: %v", req.AttackType, player.playerid)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查是否处于征讨令减半活动期间
	if player.RebelModule.GetOpenActivity() == 1 {
		needActionNum /= 2
		if needActionNum == 0 {
			needActionNum = 1
		}
	}

	bEnough := player.RoleMoudle.CheckActionEnough(gamedata.AttackRebelActionID, needActionNum)
	if !bEnough {
		gamelog.Error("Hand_AttackRebel error: Action Not Enough  , AttackType: %d , PlayerID: %d", req.AttackType, player.playerid)
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		return
	}

	//! 对叛军造成伤害
	if req.Damage > rebelModulePtr.CurLife {
		rebelModulePtr.CurLife = 0
		req.Damage = 0
	} else {
		rebelModulePtr.CurLife -= req.Damage
	}

	if req.Damage > player.RebelModule.Damage {
		//! 记录单次最高伤害
		player.RebelModule.Damage = req.Damage

		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ATTACK_REBEL_DAMAGE, int(req.Damage))

		response.DamageRank = G_RebelDamageRanker.SetRankItem(req.PlayerID, req.Damage) + 1
	}

	//! 根据伤害累加功勋
	exploit := int(req.Damage / gamedata.RebelExploitPoint)

	if player.RebelModule.GetOpenActivity() == 2 {
		//! 活动期间功勋翻倍
		exploit *= 2
	}

	player.RebelModule.Exploit += exploit

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_REBEL_EXPLOIT, player.RebelModule.Exploit)

	response.ExploitRank = G_RebelExploitRanker.SetRankItem(req.PlayerID, player.RebelModule.Exploit) + 1

	player.BagMoudle.AddAwardItem(gamedata.RebelAchievements, req.AttackType*50)

	//! 存储伤害与功勋
	player.RebelModule.DB_UpdateExploit()

	//! 扣除道具
	player.RoleMoudle.CostAction(gamedata.AttackRebelActionID, needActionNum)
	//gamelog.Info("CostActionID: %d  CostActionNum: %d  CurAction: %d", gamedata.AttackRebelActionID, needActionNum, player.RoleMoudle.GetAction(gamedata.AttackRebelActionID))

	if rebelModulePtr.CurLife <= 0 {
		//! 限时日常相关
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_KILL_REBEL, 1)

		//! 已击杀叛军
		rebelData := gamedata.GetRebelInfo(rebelModulePtr.RebelID)
		rebelModulePtr.CurLife = 0
		rebelModulePtr.RebelID = 0
		rebelModulePtr.Level += 1 //! 叛军成长
		rebelModulePtr.IsShare = false
		rebelModulePtr.EscapeTime = 0
		response.IsKill = 1

		//! 给予发现奖励
		awardItem := gamedata.GetRebelActionAward(gamedata.Find_Rebel, rebelData.Difficulty)
		if awardItem == nil {
			gamelog.Error("GetRebelActionAward fail. type: %d difficulty: %d", gamedata.Find_Rebel, rebelData.Difficulty)
			return
		}
		var award TAwardData
		award.TextType = Text_Rebel_Find
		award.ItemLst = []gamedata.ST_ItemData{*awardItem}
		award.Time = utility.GetCurTime()

		rebelcopy := gamedata.GetCopyBaseInfo(rebelData.CopyID)

		award.Value = []string{rebelcopy.Name}

		SendAwardToPlayer(req.TargetPlayerID, &award)

		//! 给予击杀奖励
		awardItem = gamedata.GetRebelActionAward(gamedata.Kill_Rebel, rebelData.Difficulty)
		if awardItem == nil {
			gamelog.Error("GetRebelActionAward fail. type: %d difficulty: %d", gamedata.Kill_Rebel, rebelData.Difficulty)
			return
		}

		award.TextType = Text_Rebel_Killed
		award.ItemLst = []gamedata.ST_ItemData{*awardItem}
		award.Time = utility.GetCurTime()

		award.Value = []string{rebelcopy.Name}
		SendAwardToPlayer(player.playerid, &award)
	}

	rebelModulePtr.DB_UpdateRebelInfo()

	response.ActionID = gamedata.AttackRebelActionID
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(response.ActionID)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 请求分享叛军
func Hand_ShareRebel(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ShareRebel_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ShareRebel Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ShareRebel_Ack
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

	player.RebelModule.CheckEscapeTime()
	player.RebelModule.CheckReset()

	player.RebelModule.IsShare = true

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求领取功勋奖励
func Hand_GetExploitAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetExploitAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetExploitAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetExploitAward_Ack
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

	player.RebelModule.CheckEscapeTime()
	player.RebelModule.CheckReset()

	//! 获取功勋奖励信息
	awardInfo := gamedata.GetExploitAward(req.ExploitAwardID)
	if awardInfo == nil {
		gamelog.Error("GetExploitAward Fail. ExploitAwardID: %d", req.ExploitAwardID)
		return
	}

	//! 检测功勋是否足够
	if player.RebelModule.Exploit < awardInfo.NeedExploit {
		gamelog.Error("Hand_GetExploitAward error: Not Enough Exploit")
		response.RetCode = msg.RE_NOT_ENOUGH_EXPLOIT
		return
	}

	//! 检查今天是否已经领取过
	if player.RebelModule.ExploitAwardLst.IsExist(req.ExploitAwardID) >= 0 {
		gamelog.Error("Hand_GetExploitAward error: Player Aleady Receive")
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 检查等级是否满足
	if awardInfo.MinLevel > player.GetLevel() || awardInfo.MaxLevel < player.GetLevel() {
		gamelog.Error("Hand_GetExploitAward error: Player invalid level. Level: %d", player.GetLevel())
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 领取奖励
	player.BagMoudle.AddAwardItem(awardInfo.ItemID, awardInfo.ItemNum)

	//! 设置标记
	player.RebelModule.ExploitAwardLst.Add(req.ExploitAwardID)
	player.RebelModule.DB_UpdateExploitAward(req.ExploitAwardID)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询功勋奖励领取状态
func Hand_GetExploitAwardStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetExploitAwardStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetExploitAwardStatus Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetExploitAwardStatus_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	response.RecvLst = []int{}

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

	player.RebelModule.CheckEscapeTime()
	player.RebelModule.CheckReset()

	for _, v := range player.RebelModule.ExploitAwardLst {
		response.RecvLst = append(response.RecvLst, v)
	}

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买战功商店物品
func Hand_BuyRebelStore(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyRebelStore_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyRebelStore Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuyRebelStore_Ack
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

	//! 检测参数
	if req.Num <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_BuyRebelStore invalid item num. Num: %v  PlayerID: %v", req.Num, player.playerid)
		return
	}

	//! 获取商品信息
	itemInfo := gamedata.GetExploitStoreItemInfo(req.ID)
	if itemInfo == nil {
		gamelog.Error("GetExploitStoreItemInfo fail. ID: %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查金币是否足够
	if player.RoleMoudle.CheckMoneyEnough(itemInfo.NeedMoneyID, itemInfo.NeedMoneyNum*req.Num) == false {
		gamelog.Error("Hand_BuyRebelStore: Not Enough Money: MoneyID: %d  MoneyNum: %d", itemInfo.NeedMoneyID, itemInfo.NeedMoneyNum*req.Num)
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 检查道具是否足够
	if itemInfo.NeedItemID != 0 {
		bEnough := player.BagMoudle.IsItemEnough(itemInfo.NeedItemID, itemInfo.NeedItemNum*req.Num)
		if !bEnough {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			gamelog.Error("Hand_BuyRebelStore : Not Enough Item")
			return
		}
	}

	//! 扣除金币与道具
	player.RoleMoudle.CostMoney(itemInfo.NeedMoneyID, itemInfo.NeedMoneyNum*req.Num)

	if itemInfo.NeedItemID != 0 {
		player.BagMoudle.RemoveNormalItem(itemInfo.NeedItemID, itemInfo.NeedItemNum*req.Num)
	}

	//! 给予物品
	player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum*req.Num)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

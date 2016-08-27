package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家请求抢劫名单
func Hand_GetRobList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetRobList_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetRobList Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetRobList_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_ROB_GEM, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取名单
	exclude := Int32Lst{}
	robLst := player.RobModule.GetRobList(req.TreasureID, exclude)
	for _, v := range robLst {
		var info msg.MSG_RobPlayerInfo
		info.PlayerID = v.PlayerID
		info.Name = v.Name
		info.Level = v.Level
		info.HeroID = v.HeroID
		info.IsRobot = v.IsRobot
		response.Lst = append(response.Lst, info)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 请求刷新抢劫名单
func Hand_RefreshRobList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RefreshRobList_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RefreshRobList Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_RefreshRobList_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_ROB_GEM, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取名单
	exclude := Int32Lst{}
	for _, v := range req.CurRobLst {
		exclude = append(exclude, v)
	}
	robLst := player.RobModule.GetRobList(req.TreasureID, exclude)
	for _, v := range robLst {
		var info msg.MSG_RobPlayerInfo
		info.PlayerID = v.PlayerID
		info.Name = v.Name
		info.Level = v.Level
		info.HeroID = v.HeroID
		info.IsRobot = v.IsRobot
		response.Lst = append(response.Lst, info)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询免战时间
func Hand_GetFreeWarTime(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_FreeWarTime_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFreeWarTime Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_FreeWarTime_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_ROB_GEM, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 刷新免战时间
	player.RobModule.RefreshFreeWarTime()

	//! 获取免战时间
	response.FreeWarTime = player.RobModule.FreeWarTime
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求抢劫武将信息
func Hand_GetRobHeroInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetRobPlayerInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetRobHeroInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetRobPlayerInfo_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_ROB_GEM, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 检查玩家精力是否足够
	config := gamedata.GetRobConfig()
	if config == nil {
		gamelog.Error("GetRobConfig nil.")
		return
	}

	//! 获取副本基本信息
	copyInfo := gamedata.GetCopyBaseInfo(config.CopyID)
	if copyInfo == nil {
		gamelog.Error("GetCopyBaseInfo fail. CopyID: %d", config.CopyID)
		return
	}

	bRet := player.RoleMoudle.CheckActionEnough(copyInfo.ActionType, copyInfo.ActionValue)
	if bRet == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	if req.IsRobot == 0 {
		//! 获取玩家信息
		robPlayer := GetPlayerByID(req.RobPlayerID)
		if robPlayer == nil {
			//! 尝试从数据库中获取玩家数据
			robPlayer = LoadPlayerFromDB(req.RobPlayerID)
		}
		if !player.BagMoudle.IsItemEnough(req.TreasureID, 1) {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			return
		}

		response.PlayerData.PlayerID = int32(req.RobPlayerID)
		var HeroResults = make([]THeroResult, BATTLE_NUM)
		response.PlayerData.FightValue = int32(robPlayer.HeroMoudle.CalcFightValue(HeroResults))
		response.PlayerData.Quality = robPlayer.HeroMoudle.CurHeros[0].Quality
		for i := 0; i < BATTLE_NUM; i++ {
			response.PlayerData.Heros[i].HeroID = int32(HeroResults[i].HeroID)
			response.PlayerData.Heros[i].PropertyValue = HeroResults[i].PropertyValues
			response.PlayerData.Heros[i].PropertyPercent = HeroResults[i].PropertyPercents
			response.PlayerData.Heros[i].CampDef = HeroResults[i].CampDef
			response.PlayerData.Heros[i].CampKill = HeroResults[i].CampKill

		}
	} else {
		robot := gamedata.GetRobot(req.RobPlayerID)
		response.PlayerData.PlayerID = int32(req.RobPlayerID)
		for i := 0; i < BATTLE_NUM; i++ {
			response.PlayerData.Heros[i].HeroID = int32(robot.Heros[i].HeroID)
			for j := 0; j < 11; j++ {
				response.PlayerData.Heros[i].PropertyValue[j] = int32(robot.Heros[i].Propertys[j])
			}

		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求抢劫
func Hand_RobTreasure(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RobTreasure_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RobTreasure Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_RobTreasure_Ack
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

	//! 检测参数
	if req.RobPlayerID <= 0 || req.PlayerID <= 0 || req.TreasureID <= 0 {
		gamelog.Error("Hand_RobTreasure Error: invalid parma. Player:%v  req: %s", player.playerid, buffer)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_ROB_GEM, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 检查玩家精力是否足够
	config := gamedata.GetRobConfig()
	if config == nil {
		gamelog.Error("GetRobConfig nil.")
		return
	}

	if player.RobModule.FreeWarTime != 0 { //! 主动抢夺, 则免战时间归零
		player.RobModule.FreeWarTime = 0
	}

	//! 获取副本基本信息
	copyInfo := gamedata.GetCopyBaseInfo(config.CopyID)
	if copyInfo == nil {
		gamelog.Error("GetCopyBaseInfo fail. CopyID: %d", config.CopyID)
		return
	}

	bRet := player.RoleMoudle.CheckActionEnough(copyInfo.ActionType, copyInfo.ActionValue)
	if bRet == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 扣除精力
	player.RoleMoudle.CostAction(copyInfo.ActionType, copyInfo.ActionValue)

	//! 获取玩家信息
	level := 0

	if req.IsRobot == 1 {
		response.RobSuccess = player.RobModule.RobNPC(req.TreasureID)
		if response.RobSuccess == true { //! 抢夺成功则给予碎片
			player.BagMoudle.AddAwardItem(req.TreasureID, 1)
		}
	} else {

		robPlayer := GetPlayerByID(req.RobPlayerID)
		if robPlayer == nil {
			//! 尝试从数据库中获取玩家数据
			robPlayer = LoadPlayerFromDB(req.RobPlayerID)
		}

		//! 检测对方是否持有碎片
		if robPlayer.BagMoudle.GetGemPieceCount(req.TreasureID) <= 0 {
			gamelog.Error("Player have not this gem piece.")
			response.RetCode = msg.RE_NOT_ENOUGH_PIECE
			return
		}

		level = robPlayer.GetLevel()
		response.RobSuccess = player.RobModule.RobPlayer(level)
		if response.RobSuccess == true { //! 抢夺成功则给予碎片
			robPlayer.BagMoudle.RemoveGemPiece(req.TreasureID, 1)
			player.BagMoudle.AddAwardItem(req.TreasureID, 1)
		}
	}

	dropItem := gamedata.GetItemsFromAwardIDEx(copyInfo.AwardID)
	if len(dropItem) != 3 {
		gamelog.Error("Hand_RobTreasure GetItemsFromAwardIDEx fail. AwardID: %d", copyInfo.AwardID)
		return
	}

	player.BagMoudle.AddAwardItem(dropItem[0].ItemID, dropItem[0].ItemNum)
	for i, v := range dropItem {
		response.DropItem[i].ID = v.ItemID
		response.DropItem[i].Num = v.ItemNum
		break
	}

	//! 增加玩家经验
	response.Exp = copyInfo.Experience * player.GetLevel()
	//! 工会技能经验加成
	if player.HeroMoudle.GuildSkiLvl[8] > 0 {
		expInc := gamedata.GetGuildSkillExpValue(player.HeroMoudle.GuildSkiLvl[8])
		response.Exp += response.Exp * expInc / 1000
	}

	player.HeroMoudle.AddMainHeroExp(response.Exp)

	//! 给予货币
	player.RoleMoudle.AddMoney(copyInfo.MoneyID, copyInfo.MoneyNum*player.GetLevel())

	response.MoneyID = copyInfo.MoneyID
	response.MoneyNum = copyInfo.MoneyNum * player.GetLevel()
	response.RetCode = msg.RE_SUCCESS
	response.FreeWarTime = player.RobModule.FreeWarTime

	//! 限时日常相关
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ROB_TIMES, 1)

}

//! 玩家请求合成宝物
func Hand_TreasureComposed(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_TreasureComposed_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_TreasureComposed Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_TreasureComposed_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_ROB_GEM, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取物品信息
	pGemInfo := gamedata.GetGemInfo(req.GemID)
	if pGemInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查玩家是否全部拥有
	for _, v := range pGemInfo.PieceIDs {
		if player.BagMoudle.GetGemPieceCount(v) < req.Num {
			response.RetCode = msg.RE_NOT_ENOUGH_PIECE
			return
		}
	}

	//! 扣除碎片
	for _, v := range pGemInfo.PieceIDs {
		isSuccess := player.BagMoudle.RemoveGemPiece(v, req.Num)
		if isSuccess == false {
			return
		}
	}

	//! 给予宝物
	player.BagMoudle.AddGems(req.GemID, req.Num)
	response.RetCode = msg.RE_SUCCESS
	response.GemID = req.GemID
	response.Num = req.Num

	if pGemInfo.Quality == 4 { //! 紫色品质
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_COMPOSITION_PURPLE, req.Num)
	} else if pGemInfo.Quality == 5 { //! 橙色品质
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_COMPOSITION_ORANGE, req.Num)
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_COMPOSITION, req.Num)
}

//! 宝物熔炼
func Hand_TreasureMelting(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_TreasureMelting_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_TreasureMelting Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_TreasureMelting_Ack
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

	//! 获取背包宝物
	gem := player.BagMoudle.GetGemByPos(req.GemPos)
	if gem == nil {
		gamelog.Error("Hand_TreasureMelting Error: Can't get gem by pos")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	targetPiece := gamedata.GetItemInfo(req.TargetPieceID)
	targetGem := gamedata.GetGemInfo(targetPiece.Data1)
	gemInfo := gamedata.GetGemInfo(gem.ID)
	if gemInfo.Quality+1 != targetGem.Quality {
		response.RetCode = msg.RE_NOT_ENOUGH_QUALITY
		return
	}

	costMoneyID, costMoneyNum := gamedata.GetTreasureMeltingInfo(targetGem.GemID)
	if costMoneyID == 0 {
		gamelog.Error("Hand_TreasureMelting Error: Get Cost fail")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if player.RoleMoudle.CheckMoneyEnough(costMoneyID, costMoneyNum) == false {
		gamelog.Error("Hand_TreasureMelting Error: Money not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(costMoneyID, costMoneyNum)

	//! 扣除宝物
	player.BagMoudle.RemoveGemAt(req.GemPos)
	player.BagMoudle.DB_RemoveGemAt(req.GemPos)

	//! 给予碎片
	player.BagMoudle.AddGemPiece(req.TargetPieceID, 1)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

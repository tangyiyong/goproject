package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

func Hand_ChangeRoleName(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_ChangeRoleName_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ChangeRoleName : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ChangeRoleName_Ack
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

	response.RetCode = CheckPlayerName(req.NewName)
	if response.RetCode != msg.RE_SUCCESS {
		return
	}

	player.RoleMoudle.Name = req.NewName
	player.RoleMoudle.DB_SaveRoleName()
	response.RetCode = msg.RE_SUCCESS

	return
}

func Hand_GetRoleData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetRoleData_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetRoleData : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetRoleData_Ack
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

	//首先将当前的速率的收益结清
	player.RoleMoudle.UpdateAllAction()
	response.Moneys = player.RoleMoudle.Moneys
	response.BatCamp = player.CamBattleModule.BattleCamp
	response.Actions = make([]int, len(player.RoleMoudle.Actions))
	response.ActionTime = make([]int64, len(player.RoleMoudle.Actions))
	for i := 0; i < len(player.RoleMoudle.Actions); i++ {
		response.Actions[i] = player.RoleMoudle.Actions[i].Value
		response.ActionTime[i] = player.RoleMoudle.Actions[i].StartTime
	}
	response.VipLevel = player.GetVipLevel()
	response.VipExp = player.GetVipExp()
	response.NewWizard = player.RoleMoudle.NewWizard
	return
}

//读取新手向导
func Hand_GetNewWizard(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetNewWizard_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetNewWizard : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetNewWizard_Ack
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

	response.NewWizard = player.RoleMoudle.NewWizard
	response.RetCode = msg.RE_SUCCESS
}

func Hand_GetCollectHeros(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetCollectionHeros_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetCollectHeros : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetCollectionHeros_Ack
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

	response.Heros = player.BagMoudle.ColHeros
	response.RetCode = msg.RE_SUCCESS
}

//设置新手向导
func Hand_SetNewWizard(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_SetNewWizard_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_SetNewWizard : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_SetNewWizard_Ack
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

	player.RoleMoudle.NewWizard = req.NewWizard
	player.RoleMoudle.DB_SaveNewWizard()
	response.RetCode = msg.RE_SUCCESS
}

//! 请求主界面的红点提示
func Hand_GetMainUITip(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetMainUITip_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetMainUITip : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetMainUITip_Ack
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

	isRed, redLst := player.StoreModule.RedTip()
	if isRed == true { //!
		for i := 0; i < len(redLst); i++ {
			if redLst[i] == gamedata.StoreType_Hero {
				response.FuncID = append(response.FuncID, gamedata.FUNC_HERO_STORE)
			} else if redLst[i] == gamedata.StoreType_Awake {
				response.FuncID = append(response.FuncID, gamedata.FUNC_AWAKEN_STORE)
			} else if redLst[i] == gamedata.StoreType_Pet {
				response.FuncID = append(response.FuncID, gamedata.FUNC_PET_STORE)
			}
		}
	} else if player.AwardCenterModule.RedTip() == true { //! 奖励中心
		response.FuncID = append(response.FuncID, gamedata.FUNC_AWARDCENTER)
	} else if player.MailMoudle.RedTip() == true { //! 邮件
		response.FuncID = append(response.FuncID, gamedata.FUNC_MAIL)
	} else if player.GuildModule.RedTip() == true { //! 公会
		response.FuncID = append(response.FuncID, gamedata.FUNC_GUILD)
	} else if player.HeroSoulsModule.RedTip() == true { //! 英灵
		response.FuncID = append(response.FuncID, gamedata.FUNC_HEROSOULS_STORE)
	} else if player.FriendMoudle.RedTip() == true { //! 好友
		response.FuncID = append(response.FuncID, gamedata.FUNC_FRIEND)
	} else if player.FameHallModule.RedTip() == true { //! 名人堂
		response.FuncID = append(response.FuncID, gamedata.FUNC_FAMOUSHALL)
	} else if player.RoleMoudle.RedTip() == true { //! 三国志
		response.FuncID = append(response.FuncID, gamedata.FUNC_SANGUOZHI)
	} else if player.ArenaModule.RedTip() == true { //! 竞技场
		response.FuncID = append(response.FuncID, gamedata.FUNC_ARENA)
	} else if player.SangokuMusouModule.RedTip() == true { //! 三国无双
		response.FuncID = append(response.FuncID, gamedata.FUNC_SANGUOWUSHUANG)
	} else if player.MiningModule.RedTip() == true { //! 挖矿
		response.FuncID = append(response.FuncID, gamedata.FUNC_MINING)
	} else if player.RebelModule.RedTip() == true { //! 围剿叛军
		response.FuncID = append(response.FuncID, gamedata.FUNC_REBEL_SIEGE)
	} else if player.TaskMoudle.RedTip() == true { //! 日常任务
		response.FuncID = append(response.FuncID, gamedata.FUNC_DAILYTASK)
	} else if player.SummonModule.RedTip() == true { //! 商城召唤
		response.FuncID = append(response.FuncID, gamedata.FUNC_SUMMON)
	} else if player.TerritoryModule.RedTip() == true { //! 领地征讨
		response.FuncID = append(response.FuncID, gamedata.FUNC_TERRITORY)
	}
}

func Hand_SanGuoZhiInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSanGuoZhiInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SanGuoZhiInfo Unmarshal is fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSanGuoZhiInfo_Ack
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
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOZHI, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取已命星的星
	response.CurOpenID = player.RoleMoudle.CurStarID
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求命星
func Hand_SetSanGuoZhi(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_SetSanGuoZhi_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetSanGuoZhi Unmarshal is fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SetSanGuoZhi_Ack
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
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOZHI, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	if gamedata.IsStarEnd(player.RoleMoudle.CurStarID) == true {
		response.RetCode = msg.RE_MAX_STAR
		return
	}

	//! 检测命星材料是否足够
	info := gamedata.GetSanGuoZhiInfo(player.RoleMoudle.CurStarID + 1)
	if info == nil {
		//! 无法获取该星信息
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	bEnough := player.BagMoudle.IsItemEnough(info.CostType, info.CostNum)
	if !bEnough {
		response.RetCode = msg.RE_SANGUOZHI_ITEM_NOT_ENOUGH
		return
	}

	player.BagMoudle.RemoveNormalItem(info.CostType, info.CostNum)

	//! 开始升星
	player.RoleMoudle.CurStarID += 1
	player.RoleMoudle.DB_SaveSanGuoZhiStar()

	if info.Type == gamedata.Sanguo_Add_Attr {
		//! 全队增加指定属性
		player.HeroMoudle.AddExtraProperty(info.AttrID, int32(info.Value), false, 0)
		player.HeroMoudle.DB_SaveExtraProperty()
	} else if info.Type == gamedata.Sanguo_Give_Item {
		//! 给予道具
		player.BagMoudle.AddAwardItem(info.AttrID, int(info.Value))
		response.AwardItem = msg.MSG_ItemData{info.AttrID, int(info.Value)}
	} else if info.Type == gamedata.Sanguo_Main_Hero_Up {
		//! 提升主角品质
		player.HeroMoudle.ChangeMainQuality(info.Value)
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_QUALITY, int(info.Value))
	}

	response.FightValue = player.CalcFightValue()
	response.Quality = player.HeroMoudle.CurHeros[0].Quality
	response.RetCode = msg.RE_SUCCESS
}

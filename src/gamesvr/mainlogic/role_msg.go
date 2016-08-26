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
	} else if player.SanGuoZhiModule.RedTip() == true { //! 三国志
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

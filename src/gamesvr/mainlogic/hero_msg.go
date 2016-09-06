package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"time"
	"utility"
)

//玩家请求上阵英雄列表
func Hand_GetBattleData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBattleData_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBattleData : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBattleData_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		comdata := utility.CompressData(b)
		//gamelog.Error("Hand_GetBattleData : orginalLen:%d, compressLen:%d", len(b), len(comdata))
		w.Write(comdata)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.CurHeros = player.HeroMoudle.CurHeros
	response.BackHeros = player.HeroMoudle.BackHeros
	response.Equips = player.HeroMoudle.CurEquips
	response.Gems = player.HeroMoudle.CurGems
	response.Pets = player.HeroMoudle.CurPets
	response.Title = player.HeroMoudle.TitleID
	response.GuildSkiLvl = player.HeroMoudle.GuildSkiLvl
	response.FashionID = player.HeroMoudle.FashionID
	response.FashionLvl = player.HeroMoudle.FashionLvl
	response.ExtraProValue = player.HeroMoudle.ExtraProValue
	response.ExtraProPercent = player.HeroMoudle.ExtraProPercent
	response.ExtraCampDef = player.HeroMoudle.ExtraCampDef
	response.ExtraCampKill = player.HeroMoudle.ExtraCampKill
	response.CurStarID = player.RoleMoudle.CurStarID
	response.RetCode = msg.RE_SUCCESS
	return
}

//升级英雄
func Hand_UpgradeHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_UpgradeHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpgradeHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UpgradeHero_Ack
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

	if req.TargetHero.PosType == POSTYPE_BATTLE && req.TargetHero.HeroPos == 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_UpgradeHero error : Main Hero Can't Upgrade!")
		return
	}

	var pTargetHeroData *THeroData = nil
	if req.TargetHero.PosType == POSTYPE_BATTLE {
		pTargetHeroData = player.HeroMoudle.GetBattleHeroByPos(req.TargetHero.HeroPos)
	} else if req.TargetHero.PosType == POSTYPE_BAG {
		pTargetHeroData = player.BagMoudle.GetBagHeroByPos(req.TargetHero.HeroPos)
	} else if req.TargetHero.PosType == POSTYPE_BACK {
		pTargetHeroData = player.HeroMoudle.GetBackHeroByPos(req.TargetHero.HeroPos)
	}

	//检验目标英雄是不是正确
	if pTargetHeroData == nil || pTargetHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_UpgradeHero error : data.heroid:%d--req.heroid:%d", pTargetHeroData.ID, req.TargetHero.HeroID)
		return
	}

	//检验目标英雄的等级是不是己经不能进行升级了
	if pTargetHeroData.Level >= player.GetLevel() {
		gamelog.Error("Hand_UpgradeHero error : normal hero level can't greater than main hero")
		response.RetCode = msg.RE_CNT_OVER_MAIN_HERO_LEVEL
		return
	}

	var OldLevel int = pTargetHeroData.Level

	//验证消耗英雄的顺序
	//统计消耗英雄产生的经验
	var tempPos = 10000
	var ExpSum = 0
	for _, t := range req.CostHeros {
		pTempHeroData := player.BagMoudle.GetBagHeroByPos(t.HeroPos)
		if pTempHeroData == nil || pTempHeroData.ID != t.HeroID || t.HeroID == 0 {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpgradeHero error :  costid: %d, costpos :%d", t.HeroID, t.HeroPos)
			return
		}

		if t.HeroPos > tempPos {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpgradeHero error :  Wrong Squence: %d", t.HeroPos)
			return
		}

		tempPos = t.HeroPos
		pHeroInfo := gamedata.GetHeroInfo(t.HeroID)
		ExpSum += pTempHeroData.CurExp + pHeroInfo.HeroExp

		if req.TargetHero.PosType == POSTYPE_BAG {
			if t.HeroPos == req.TargetHero.HeroPos {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_UpgradeHero error :   TargetPos == costpos: %d", t.HeroPos)
				return
			}
		}
	}

	pHeroLevelInfo := gamedata.GetHeroLevelInfo(pTargetHeroData.Quality, pTargetHeroData.Level)
	if false == player.RoleMoudle.CheckMoneyEnough(pHeroLevelInfo.MoneyID, ExpSum*pHeroLevelInfo.MoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_UpgradeHero Error : Not Enough Money moneyid:%d, moneynum:%d", pHeroLevelInfo.MoneyID, ExpSum*pHeroLevelInfo.MoneyNum)
		return
	}

	player.HeroMoudle.AddHeroExp(req.TargetHero.PosType, req.TargetHero.HeroPos, ExpSum)
	player.RoleMoudle.CostMoney(pHeroLevelInfo.MoneyID, ExpSum*pHeroLevelInfo.MoneyNum)
	response.NewLevel = pTargetHeroData.Level
	response.NewExp = pTargetHeroData.CurExp

	//必须以不影响的索引的方式删除
	for t := 0; t < len(req.CostHeros); t++ {
		player.BagMoudle.RemoveHeroAt(req.CostHeros[t].HeroPos)
	}

	player.BagMoudle.DB_SaveHeroBag()
	response.RetCode = msg.RE_SUCCESS
	response.CostMoney = ExpSum * pHeroLevelInfo.MoneyNum

	if req.TargetHero.PosType == POSTYPE_BATTLE && OldLevel < response.NewLevel {
		response.FightValue = player.CalcFightValue()
	}

	return
}

func Hand_LevelUpNotify(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_LevelUpNotify_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_LevelUpNotify : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_LevelUpNotify_Ack
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

	pMainHero := &player.HeroMoudle.CurHeros[0]

	oldLevel := pMainHero.Level

	if oldLevel >= gamedata.G_HeroMaxLevel {
		pMainHero.Level = gamedata.G_HeroMaxLevel
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_LevelUpNotify : Main Has already reach the max level!!!")
		return
	}

	pHeroInfo := gamedata.GetHeroInfo(pMainHero.ID)
	if pHeroInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_LevelUpNotify : Invalid Main Hero ID:%d!!!", pMainHero.ID)
		return
	}

	pStHeroLevelInfo := gamedata.GetHeroLevelInfo(pHeroInfo.Quality, pMainHero.Level)
	if pStHeroLevelInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_LevelUpNotify error :  Invalid Quality: %d and Level:%d", pMainHero.Quality, pMainHero.Level)
		return
	}

	if pMainHero.CurExp >= pStHeroLevelInfo.MainNeedExp {
		pMainHero.CurExp -= pStHeroLevelInfo.MainNeedExp
		pMainHero.Level += 1
		player.DB_SaveHeroLevelExp(POSTYPE_BATTLE, 0)
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_LEVEL_UP, pMainHero.Level-oldLevel)
		G_LevelRanker.SetRankItem(req.PlayerID, pMainHero.Level)
		response.FightValue = player.CalcFightValue()
		player.ActivityModule.LevelGift.CheckLevelUp(pMainHero.Level)
	} else {
		gamelog.Error("Hand_LevelUpNotify Error : CurExp:%d, needExp:%d!!", pMainHero.CurExp, pStHeroLevelInfo.MainNeedExp)
	}

	response.Level = player.HeroMoudle.CurHeros[0].Level
	response.CurExp = player.HeroMoudle.CurHeros[0].CurExp
	response.CurSvrTime = time.Now().Unix()
	response.RetCode = msg.RE_SUCCESS
	return
}

func Hand_ChangeBackHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ChangeBackHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ChangeBackHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ChangeBackHero_Ack
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

	pSrcInfo := gamedata.GetHeroInfo(req.SourceID)
	if pSrcInfo == nil {
		gamelog.Error("Hand_ChangeHero Invalid req.SourceID:%d", req.SourceID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if pSrcInfo.Setup <= 0 {
		gamelog.Error("Hand_ChangeHero heor %d can set to battle", req.SourceID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if req.TargetID == 0 { //上阵
		if !gamedata.IsFuncOpen(gamedata.FUNC_BACK_POS_BEGIN+req.TargetPos-1, player.GetLevel(), 0) {
			gamelog.Error("Hand_ChangeBackHero battle pos is not open!")
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}
	}

	if req.TargetPos < 0 || req.TargetPos >= BACK_NUM {
		gamelog.Error("Hand_ChangeBackHero error Invalid TargetPos:%d", req.TargetPos)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if req.SourcePos < 0 || req.SourcePos >= len(player.BagMoudle.HeroBag.Heros) {
		gamelog.Error("Hand_ChangeBackHero error Invalid SourcePos:%d", req.SourcePos)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	tempTarget := player.HeroMoudle.BackHeros[req.TargetPos]
	if tempTarget.ID != req.TargetID {
		gamelog.Error("Hand_ChangeBackHero error req.TargetID:%d, tempTarget.HeroID:%d", req.TargetID, tempTarget.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	tempSource := player.BagMoudle.HeroBag.Heros[req.SourcePos]
	if tempSource.ID != req.SourceID {
		gamelog.Error("Hand_ChangeBackHero error req.SourceID:%d, req.SourcePos:%d, Source.HeroID:%d", req.SourceID, req.SourcePos, tempSource.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//将英雄装到上阵英雄中
	player.HeroMoudle.BackHeros[req.TargetPos] = tempSource
	player.DB_SaveHeroAt(POSTYPE_BACK, req.TargetPos)

	if req.TargetID == 0 { //上阵
		//删除掉背包中的英雄
		player.BagMoudle.RemoveHeroAt(req.SourcePos)
		player.BagMoudle.DB_SaveHeroBag()
	} else {
		player.BagMoudle.HeroBag.Heros[req.SourcePos] = tempTarget
		player.DB_SaveHeroAt(POSTYPE_BAG, req.SourcePos)
	}

	response.RetCode = msg.RE_SUCCESS
	response.FightValue = player.CalcFightValue()
}

func Hand_ChangeHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ChangeHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ChangeHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ChangeHero_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	pSrcInfo := gamedata.GetHeroInfo(req.SourceID)
	if pSrcInfo == nil {
		gamelog.Error("Hand_ChangeHero Invalid req.SourceID:%d", req.SourceID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if pSrcInfo.Setup <= 0 {
		gamelog.Error("Hand_ChangeHero hero %d can set to battle", req.SourceID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if req.TargetPos < 0 || req.TargetPos >= BATTLE_NUM {
		gamelog.Error("Hand_ChangeHero error :Invalid TargetPos:%d", req.TargetPos)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if req.SourcePos < 0 || req.SourcePos >= len(player.BagMoudle.HeroBag.Heros) {
		gamelog.Error("Hand_ChangeHero error :Invalid SourcePos:%d", req.SourcePos)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	tempTarget := player.HeroMoudle.CurHeros[req.TargetPos]
	if tempTarget.ID != req.TargetID {
		gamelog.Error("Hand_ChangeHero error : TargetHeroID :%d,  req.TargetID:%d", tempTarget.ID, req.TargetID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	tempSource := player.BagMoudle.HeroBag.Heros[req.SourcePos]
	if tempSource.ID != req.SourceID || req.SourceID == 0 {
		gamelog.Error("Hand_ChangeHero req.SourePos :%d : SourceID:%d, req.SourceID:%d", req.SourcePos, tempSource.ID, req.SourceID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if req.TargetID == 0 { //上阵
		if !gamedata.IsFuncOpen(gamedata.FUNC_POS_START+req.TargetPos-1, player.GetLevel(), 0) {
			gamelog.Error("Hand_ChangeHero battle pos is not open!, tPos:%d", req.TargetPos)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}
	}

	//将英雄装到上阵英雄中
	player.HeroMoudle.CurHeros[req.TargetPos] = tempSource
	player.DB_SaveHeroAt(POSTYPE_BATTLE, req.TargetPos)

	if req.TargetID == 0 { //上阵
		//删除掉背包中的英雄
		player.BagMoudle.RemoveHeroAt(req.SourcePos)
		player.BagMoudle.DB_RemoveHeroAt(req.SourcePos)
		//player.DB_SaveHeros(POSTYPE_BAG)
	} else {
		player.BagMoudle.HeroBag.Heros[req.SourcePos] = tempTarget
		player.DB_SaveHeroAt(POSTYPE_BAG, req.SourcePos)
	}

	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS

}

func Hand_SetWakeItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_SetWakeItem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_SetWakeItem : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_SetWakeItem_Ack
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

	if player.BagMoudle.GetWakeItemCount(req.SourceItemID) <= 0 {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_SetWakeItem : Not Enough Wake Item:%d", req.SourceItemID)
		return
	}

	pTargetHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if (pTargetHeroData == nil) || pTargetHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_SetWakeItem : req.posType:%d, req.Pos:%d, req.id:%d, targetID:%d", req.TargetHero.PosType,
			req.TargetHero.HeroPos, req.TargetHero.HeroID, pTargetHeroData.ID)
		return
	}

	pWakeLevelInfo := gamedata.GetWakeLevelItem(pTargetHeroData.WakeLevel)
	if pWakeLevelInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_SetWakeItem : Invalid Wake Level:%d", pTargetHeroData.WakeLevel)
		return
	}

	if pWakeLevelInfo.NeedItem[req.TargetItemPos] != req.SourceItemID || pTargetHeroData.WakeItem[req.TargetItemPos] != 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_SetWakeItem : Invalid Wake NeedID:%d, SourceID:%d, CurID:%d", pWakeLevelInfo.NeedItem[req.TargetItemPos], req.SourceItemID, pTargetHeroData.WakeItem[req.TargetItemPos])
		return
	}

	player.BagMoudle.RemoveWakeItem(req.SourceItemID, 1)
	pTargetHeroData.WakeItem[req.TargetItemPos] = req.SourceItemID
	player.DB_SaveHeroWakeItem(req.TargetHero.PosType, req.TargetHero.HeroPos)
	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS

	return
}

func Hand_ComposeWakeItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ComposeWakeItem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ComposeWakeItem : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ComposeWakeItem_Ack
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

	pWakeComposeInfo := gamedata.GetWakeComposeInfo(req.ItemID)
	if pWakeComposeInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ComposeWakeItem : Invalid Wake ItemID:%d", req.ItemID)
		return
	}

	//钱是否足够
	if false == player.RoleMoudle.CheckMoneyEnough(pWakeComposeInfo.MoneyID, pWakeComposeInfo.MoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_ComposeWakeItem : Not Enough Money")
		return
	}

	for i := 0; i < 4; i++ {
		if pWakeComposeInfo.Items[i].ItemID == 0 {
			break
		}

		if player.BagMoudle.GetWakeItemCount(pWakeComposeInfo.Items[i].ItemID) < pWakeComposeInfo.Items[i].ItemNum {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			gamelog.Error("Hand_ComposeWakeItem : Not Enough Item :%d", pWakeComposeInfo.Items[i].ItemID)
			return
		}
	}

	player.BagMoudle.AddWakeItem(req.ItemID, 1)
	player.RoleMoudle.CostMoney(pWakeComposeInfo.MoneyID, pWakeComposeInfo.MoneyNum)
	for i := 0; i < 4; i++ {
		if pWakeComposeInfo.Items[i].ItemID == 0 {
			break
		}

		player.BagMoudle.RemoveWakeItem(pWakeComposeInfo.Items[i].ItemID, pWakeComposeInfo.Items[i].ItemNum)
	}

	response.RetCode = msg.RE_SUCCESS
}

func Hand_UpWakeLevel(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_UpWakeLevel_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpWakeLevel : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UpWakeLevel_Ack
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

	pTargetHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if (pTargetHeroData == nil) || pTargetHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_UpWakeLevel : req.posType:%d, req.Pos:%d, req.id:%d, targetID:%d", req.TargetHero.PosType,
			req.TargetHero.HeroPos, req.TargetHero.HeroID, pTargetHeroData.ID)
		return
	}

	var bHost = false
	if req.TargetHero.PosType == POSTYPE_BATTLE && req.TargetHero.HeroPos == 0 {
		bHost = true
	}

	pWakeLevel := gamedata.GetWakeLevelItem(pTargetHeroData.WakeLevel)
	if pWakeLevel == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//道具是否齐全
	for i := 0; i < len(pWakeLevel.NeedItem); i++ {
		if pWakeLevel.NeedItem[i] != 0 {
			if pTargetHeroData.WakeItem[i] == 0 {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_UpWakeLevel : Not Enough Items!!!!")
				return
			}
		}
	}

	//是否达到需要的等级
	if pTargetHeroData.Level < pWakeLevel.NeedLevel {
		response.RetCode = msg.RE_NOT_ENOUGH_HERO_LEVEL
		gamelog.Error("Hand_UpWakeLevel : Not Enough hero Level!!!!")
		return
	}

	//是否有足够的觉醒丹
	needCount := 0
	if bHost {
		needCount = pWakeLevel.HostWakeNum
	} else {
		needCount = pWakeLevel.NeedWakeNum
	}
	if false == player.BagMoudle.IsItemEnough(pWakeLevel.NeedWakeID, needCount) {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_UpWakeLevel : Not Enough Wake Items!!!!")
		return
	}

	//是否有足够的货币
	if false == player.RoleMoudle.CheckMoneyEnough(pWakeLevel.NeedMoneyID, pWakeLevel.NeedMoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_UpWakeLevel : Not Enough Money!!!!")
		return
	}

	//是否有需要的同名英雄
	if pWakeLevel.NeedHeroNum != 0 && bHost == false {
		pSourceHero := player.BagMoudle.GetBagHeroByPos(req.SourcePos)
		if pSourceHero == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpWakeLevel : Need hero!!!!")
			return
		}

		if pSourceHero.ID != pTargetHeroData.ID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpWakeLevel : Need hero!!!!")
			return
		}
	}

	pTargetHeroData.WakeLevel += 1
	pTargetHeroData.WakeItem[0] = 0
	pTargetHeroData.WakeItem[1] = 0
	pTargetHeroData.WakeItem[2] = 0
	pTargetHeroData.WakeItem[3] = 0

	if pWakeLevel.NeedHeroNum != 0 && bHost == false {
		player.BagMoudle.RemoveHeroAt(req.SourcePos)
		player.BagMoudle.DB_RemoveHeroAt(req.SourcePos)
		//player.DB_SaveHeros(POSTYPE_BAG)
	}

	player.RoleMoudle.CostMoney(pWakeLevel.NeedMoneyID, pWakeLevel.NeedMoneyNum)

	player.BagMoudle.RemoveNormalItem(pWakeLevel.NeedWakeID, needCount)

	player.DB_SaveHeroWakeLevel(req.TargetHero.PosType, req.TargetHero.HeroPos)
	response.RetCode = msg.RE_SUCCESS
	response.FightValue = player.CalcFightValue()
	response.WakeLevel = pTargetHeroData.WakeLevel

	minlevel := pTargetHeroData.WakeLevel
	for i := 0; i < len(player.HeroMoudle.CurHeros); i++ {
		if player.HeroMoudle.CurHeros[i].WakeLevel < minlevel {
			minlevel = player.HeroMoudle.CurHeros[i].WakeLevel
		}
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_WAKE, minlevel)

	return
}

func Hand_Change_Career(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ChangeCareer_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ComposeHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ChangeCareer_Ack
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

	if 1 == player.HeroMoudle.CurHeros[0].ID%2 {
		player.HeroMoudle.CurHeros[0].ID += 1
	} else {
		player.HeroMoudle.CurHeros[0].ID -= 1
	}

	player.HeroMoudle.DB_SaveMainHeroID()

	response.NewHeroID = player.HeroMoudle.CurHeros[0].ID
	response.RetCode = msg.RE_SUCCESS

}

func Hand_UpgodHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_UpgodHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpgodHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UpgodHero_Ack
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

	if false == gamedata.IsFuncOpen(gamedata.FUNC_HEROGOD, player.GetLevel(), player.GetVipLevel()) {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Hand_UpgodHero : Hero God Not Open!!!")
		return
	}

	if req.PosType == POSTYPE_BATTLE && req.PosIndex == 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_UpgodHero : can't upgod main hero!!!")
		return
	}

	pTargetHeroData := player.GetHeroByPos(req.PosType, req.PosIndex)
	if (pTargetHeroData == nil) || pTargetHeroData.ID != req.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_UpgodHero : req.posType:%d, req.Pos:%d, req.id:%d, targetID:%d", req.PosType, req.PosIndex, req.HeroID, pTargetHeroData.ID)
		return
	}

	if pTargetHeroData.Quality < 5 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_UpgodHero : quality < 5 hero cant up god")
		return
	}

	var condgod int = 0
	var targetgod int = 0
	if pTargetHeroData.GodLevel <= 0 {
		if pTargetHeroData.Quality == 6 {
			condgod = 15
			targetgod = 16
		} else {
			condgod = 0
			targetgod = 1
		}
	} else {
		condgod = pTargetHeroData.GodLevel
		targetgod = pTargetHeroData.GodLevel + 1
	}
	pHeroGodInfo := gamedata.GetHeroGodInfo(condgod)
	if pHeroGodInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_UpgodHero : Invalid condgod :%d!!!!", condgod)
		return
	}

	//检测所需的道具是否足够
	if false == player.RoleMoudle.CheckMoneyEnough(pHeroGodInfo.NeedMoneyID, pHeroGodInfo.NeedMoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_UpgodHero : Not Enough Money!,id:%d, num;%d", pHeroGodInfo.NeedMoneyID, pHeroGodInfo.NeedMoneyNum)
		return
	}

	if false == player.BagMoudle.IsItemEnough(pHeroGodInfo.NeedItemID, pHeroGodInfo.NeedItemNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_UpgodHero : Not Enough Item:%d!!!", pHeroGodInfo.NeedItemID)
		return
	}

	if pHeroGodInfo.NeedType == 1 { //货币
		if false == player.RoleMoudle.CheckMoneyEnough(pHeroGodInfo.NeedID, pHeroGodInfo.NeedNum) {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			gamelog.Error("Hand_UpgodHero : Not Enough Money!,id:%d, num;%d", pHeroGodInfo.NeedMoneyID, pHeroGodInfo.NeedMoneyNum)
			return
		}

	} else if pHeroGodInfo.NeedType == 2 { //碎片
		pHeroInfo := gamedata.GetHeroInfo(req.HeroID)
		if pHeroInfo == nil {
			gamelog.Error("Hand_UpgodHero : Invalid Hero ID:%d!!!", req.HeroID)
			return
		}
		if player.BagMoudle.GetHeroPieceCount(pHeroInfo.PieceID) < pHeroGodInfo.NeedNum {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			gamelog.Error("Hand_UpgodHero : Not Enough Hero Piece Num!!!")
			return
		}

	} else if pHeroGodInfo.NeedType == 3 { //道具
		if false == player.BagMoudle.IsItemEnough(pHeroGodInfo.NeedID, pHeroGodInfo.NeedNum) {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			gamelog.Error("Hand_UpgodHero : Not Enough Hero Piece Num!!!")
			return
		}
	}

	pTargetHeroData.GodLevel = targetgod
	if pTargetHeroData.Quality == 5 && targetgod == 16 {
		pTargetHeroData.Quality += 1
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_QUALITY, int(pTargetHeroData.Quality))
	}

	player.DB_SaveHeroGodLevel(req.PosType, req.PosIndex)
	player.RoleMoudle.CostMoney(pHeroGodInfo.NeedMoneyID, pHeroGodInfo.NeedMoneyNum)
	player.BagMoudle.RemoveNormalItem(pHeroGodInfo.NeedItemID, pHeroGodInfo.NeedItemNum)
	response.GodLevel = pTargetHeroData.GodLevel
	response.Quality = pTargetHeroData.Quality
	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS

	if pHeroGodInfo.NeedType == 1 { //货币
		player.RoleMoudle.CostMoney(pHeroGodInfo.NeedID, pHeroGodInfo.NeedNum)
	} else if pHeroGodInfo.NeedType == 2 { //碎片
		pHeroInfo := gamedata.GetHeroInfo(req.HeroID)
		if pHeroInfo == nil {
			gamelog.Error("Hand_UpgodHero : Invalid Hero ID:%d!!!", req.HeroID)
			return
		}
		player.BagMoudle.RemoveHeroPiece(pHeroInfo.PieceID, pHeroGodInfo.NeedNum)
	} else if pHeroGodInfo.NeedType == 3 { //道具
		player.BagMoudle.RemoveNormalItem(pHeroGodInfo.NeedID, pHeroGodInfo.NeedNum)
	}
}

func Hand_ComposeHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ComposeHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ComposeHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ComposeHero_Ack
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

	pHeroPieceInfo := gamedata.GetItemInfo(req.HeroPieceID)
	if pHeroPieceInfo == nil {
		gamelog.Error("Hand_ComposeHero Error : Invalid PieceID :%d", req.HeroPieceID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pHeroInfo := gamedata.GetHeroInfo(pHeroPieceInfo.Data1)
	if pHeroInfo == nil {
		gamelog.Error("Hand_ComposeHero Error : Invalid HeroID :%d", pHeroPieceInfo.Data1)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pieceCount := player.BagMoudle.GetHeroPieceCount(req.HeroPieceID)
	if pieceCount < pHeroInfo.PieceNum {
		response.RetCode = msg.RE_NOT_ENOUGH_PIECE
		gamelog.Error("Hand_ComposeHero Error : Not Enough Hero Piece :%d", pieceCount)
		return
	}

	player.BagMoudle.AddHeroByID(pHeroInfo.HeroID, 1)
	player.BagMoudle.RemoveHeroPiece(req.HeroPieceID, pHeroInfo.PieceNum)

	response.HeroID = pHeroInfo.HeroID
	response.RetCode = msg.RE_SUCCESS
	return
}

func Hand_QueryHeroDestiny(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_QueryDestinyState_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryHeroDestiny : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_QueryDestinyState_Ack
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

	pCurHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if pCurHeroData == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_QueryHeroDestiny : Invalie TargetHero PosType:%d, HeroPos:%d, ID:%d", req.TargetHero.PosType, req.TargetHero.HeroPos, req.TargetHero.HeroID)
		return
	}

	if pCurHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_QueryHeroDestiny : Invalie TargetHero msg.HeroID:%d, data.HeroID:%d", req.TargetHero.HeroID, pCurHeroData.ID)
		return
	}
	response.NewDestinyState = pCurHeroData.DestinyState
	DestinyLevel := pCurHeroData.DestinyState >> 24 & 0x000F
	DestinyIndex := pCurHeroData.DestinyState >> 16 & 0x000F
	DestinyLight := pCurHeroData.DestinyState & 0x000F
	daychange := utility.GetCurDayByUnix() - pCurHeroData.DestinyTime
	if DestinyLight <= 0 || daychange <= 0 {
		response.RetCode = msg.RE_SUCCESS
		return
	}

	DestinyLight = DestinyLight - daychange
	if DestinyLight < 0 || DestinyLight > 4 {
		DestinyLight = 0
	}

	pCurHeroData.DestinyState = DestinyLevel
	pCurHeroData.DestinyState = pCurHeroData.DestinyState << 8
	pCurHeroData.DestinyState += DestinyIndex
	pCurHeroData.DestinyState = pCurHeroData.DestinyState << 16
	pCurHeroData.DestinyState += DestinyLight
	player.DB_SaveHeroDestiny(req.TargetHero.PosType, req.TargetHero.HeroPos)
	response.NewDestinyState = pCurHeroData.DestinyState
	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS
	return
}

func Hand_DestinyHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_DestinyHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_DestinyHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_DestinyHero_Ack
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

	pCurHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if (pCurHeroData == nil) || pCurHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_DestinyHero : Wrong Param curid:%d, tid:%d", pCurHeroData.ID, req.TargetHero.HeroID)
		return
	}

	DestinyLevel := pCurHeroData.DestinyState >> 24 & 0x000F
	DestinyIndex := pCurHeroData.DestinyState >> 16 & 0x000F
	DestinyLight := pCurHeroData.DestinyState & 0x000F

	pHeroDestinyInfo := gamedata.GetHeroDestinyInfo(int(DestinyLevel))
	if pHeroDestinyInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_DestinyHero : Invalid Destiny Level:%d, State:%d", DestinyLevel, pCurHeroData.DestinyState)
		return
	}

	bEnough := player.BagMoudle.IsItemEnough(pHeroDestinyInfo.CostItemID, pHeroDestinyInfo.OneTimeCost)
	if !bEnough {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_DestinyHero : Not enough item %d", pHeroDestinyInfo.CostItemID)
		return
	}

	if utility.Rand() < (pHeroDestinyInfo.UpgradeRatio * 10) {
		DestinyLight += 1
		if DestinyLight >= 4 {
			DestinyIndex += 1
			DestinyLight = 0
			response.FightValue = player.CalcFightValue()
			if DestinyIndex >= 5 {
				DestinyLevel += 1
				DestinyIndex = 0
				DestinyLight = 0
			}
		}

		pCurHeroData.DestinyState = DestinyLevel
		pCurHeroData.DestinyState = pCurHeroData.DestinyState << 8
		pCurHeroData.DestinyState += DestinyIndex
		pCurHeroData.DestinyState = pCurHeroData.DestinyState << 16
		pCurHeroData.DestinyState += DestinyLight
		pCurHeroData.DestinyTime = utility.GetCurDayByUnix()
		player.DB_SaveHeroDestiny(req.TargetHero.PosType, req.TargetHero.HeroPos)

		//! 点亮火炬次数+1
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_GOD_TIMES, 1)
	}

	player.BagMoudle.RemoveNormalItem(pHeroDestinyInfo.CostItemID, pHeroDestinyInfo.OneTimeCost)
	response.CostItemNum = pHeroDestinyInfo.OneTimeCost
	response.NewDestinyState = pCurHeroData.DestinyState
	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS

	//! 上阵武将天命等级
	isExist := true
	var minLevel uint32 = 0x7FFFFFFF
	var maxlevel uint32 = 0
	for i := 0; i < BATTLE_NUM; i++ {
		hero := &player.HeroMoudle.CurHeros[i]
		if hero.ID == 0 {
			isExist = false
			continue
		}

		dlevel := hero.DestinyState >> 24 & 0x000F
		if dlevel > maxlevel {
			maxlevel = dlevel
		}

		if minLevel > dlevel {
			minLevel = dlevel
		}
	}

	if isExist == true {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_DESTINY_LEVEL, int(minLevel))
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_DESTINY_LEVEL_MAX, int(maxlevel))

	return
}

func Hand_CultureHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_CultureHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_CultureHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_CultureHero_Ack
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

	pHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if (pHeroData == nil) || pHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_CultureHero : Invalid PosType:%d, Pos:%d, ID:%d!!!!", req.TargetHero.PosType, req.TargetHero.HeroPos, req.TargetHero.HeroID)
		return
	}

	pHeroInfo := gamedata.GetHeroInfo(pHeroData.ID)
	pCultureMaxInfo := gamedata.GetCultureMaxInfo(pHeroInfo.AttackType)
	if (pHeroInfo == nil) || pCultureMaxInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_CultureHero : cant get the static config data")
		return
	}

	//需求是否足够
	if false == player.BagMoudle.IsItemEnough(gamedata.CultureItemID, gamedata.CultureItemNum*req.Times) {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_CultureHero : Not Enough Items; itemid:%d", gamedata.CultureItemID)
		return
	}

	//培养的次数处理
	for j := 0; j < req.Times; j++ {
		pHeroData.Cultures[0] += utility.Rand() % 10 //生命
		pHeroData.Cultures[2] += utility.Rand() % 10 //物防
		pHeroData.Cultures[4] += utility.Rand() % 10 //法防
		rValue := utility.Rand() % 10                //攻击力
		pHeroData.Cultures[1] += rValue              //物攻
		pHeroData.Cultures[3] += rValue              //魔攻
	}

	//上限处理
	for i := 0; i < 5; i++ {
		if pHeroData.Cultures[i] > pHeroData.Level*pCultureMaxInfo.MaxRation[i] {
			pHeroData.Cultures[i] = pHeroData.Level * pCultureMaxInfo.MaxRation[i]
		}
	}

	player.BagMoudle.RemoveNormalItem(gamedata.CultureItemID, gamedata.CultureItemNum*req.Times)
	pHeroData.CulturesCost += gamedata.CultureItemNum * req.Times
	player.DB_SaveHeroCulture(req.TargetHero.PosType, req.TargetHero.HeroPos)

	response.CostItems = gamedata.CultureItemNum * req.Times
	response.Cultures = pHeroData.Cultures
	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS
	//! 增加日常任务进度
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_CULTURE, req.Times)
}

func Hand_BreakOutHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_BreakOut_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_BreakOutHero : Unmarshal error!!!!%s", string(buffer))
		return
	}

	var response msg.MSG_BreakOut_Ack
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

	pTargetHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if (pTargetHeroData == nil) || pTargetHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_BreakOutHero : req.posType:%d, req.Pos:%d, req.id:%d, targetID:%d", req.TargetHero.PosType,
			req.TargetHero.HeroPos, req.TargetHero.HeroID, pTargetHeroData.ID)
		return
	}

	var bHost = false
	if req.TargetHero.PosType == POSTYPE_BATTLE && req.TargetHero.HeroPos == 0 {
		bHost = true
	}

	pBreakLevelInfo := gamedata.GetHeroBreakInfo(pTargetHeroData.BreakLevel)
	if pBreakLevelInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_BreakOutHero : Invalid BreakLevel :%d!!!!", pTargetHeroData.BreakLevel)
		return
	}

	if pTargetHeroData.Level < pBreakLevelInfo.NeedLevel {
		response.RetCode = msg.RE_NOT_ENOUGH_HERO_LEVEL
		gamelog.Error("Hand_BreakOutHero : Not Enough Hero Level :%d!!!!", pTargetHeroData.Level)
		return
	}

	needHeroCount := pBreakLevelInfo.HeroNum
	needItemCount := pBreakLevelInfo.ItemNum

	//如果是英雄则需求数目需要调整
	if bHost {
		needHeroCount = 0
		needItemCount = pBreakLevelInfo.HostItemNum
	}

	bEnough := player.BagMoudle.IsItemEnough(pBreakLevelInfo.ItemID, needItemCount)
	if !bEnough {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_BreakOutHero : Invalid HeroBreakItemID :%d!!!!", pBreakLevelInfo.ItemID)
		return
	}

	if false == player.RoleMoudle.CheckMoneyEnough(pBreakLevelInfo.MoneyID, pBreakLevelInfo.MoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_BreakOutHero : Not Enough Money!")
		return
	}

	if needHeroCount > len(req.CostHeros) {
		response.RetCode = msg.RE_NOT_ENOUGH_HERO
		gamelog.Error("Hand_BreakOutHero : lack of same name heros! need:%d, has :%d", needHeroCount, len(req.CostHeros))
		return
	}

	var tempPos = 100000
	var pHeroData *THeroData = nil
	for _, t := range req.CostHeros {
		pHeroData = player.BagMoudle.GetBagHeroByPos(t.HeroPos)
		if pHeroData == nil || pHeroData.ID != t.HeroID || t.HeroID == 0 {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_BreakOutHero error :  Invalid SourcePos: %d", t.HeroPos)
			return
		}

		if t.HeroPos > tempPos {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_BreakOutHero error :  Wrong Squence: %d", t.HeroPos)
			return
		}

		tempPos = t.HeroPos

		if req.TargetHero.PosType == POSTYPE_BAG {
			if t.HeroPos == req.TargetHero.HeroPos {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_BreakOutHero error :  Invalid TargetPos: %d", t.HeroPos)
				return
			}
		}

		if t.HeroID != req.TargetHero.HeroID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_BreakOutHero error :  Invalid SourceID: %d not same as the target heroid:%d", t.HeroID, req.TargetHero.HeroID)
			return
		}
	}

	pTargetHeroData.BreakLevel += 1
	player.DB_SaveHeroBreakLevel(req.TargetHero.PosType, req.TargetHero.HeroPos)
	player.RoleMoudle.CostMoney(pBreakLevelInfo.MoneyID, pBreakLevelInfo.MoneyNum)
	player.BagMoudle.RemoveNormalItem(pBreakLevelInfo.ItemID, needItemCount)
	response.NewLevel = pTargetHeroData.BreakLevel

	//必须以不影响的索引的方式删除
	for t := len(req.CostHeros) - 1; t >= 0; t-- {
		player.BagMoudle.RemoveHeroAt(req.CostHeros[t].HeroPos)
	}
	player.BagMoudle.DB_SaveHeroBag()

	response.FightValue = player.CalcFightValue()
	response.CostItems = needItemCount
	response.CostMoney = pBreakLevelInfo.MoneyNum
	response.RetCode = msg.RE_SUCCESS

	minLevel := player.HeroMoudle.CurHeros[0].BreakLevel
	for _, v := range player.HeroMoudle.CurHeros {
		if v.BreakLevel < minLevel {
			minLevel = v.BreakLevel
		}
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_CUR_HERO_BREAK, int(minLevel))

	return
}

func Hand_ChangeEquip(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ChangeEquip_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ChangeEquip : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ChangeEquip_Ack
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

	if req.TargetPos < 0 || req.TargetPos >= EQUIP_NUM {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ChangeEquip : Invalid TargetPos :%d", req.TargetPos)
		return
	}

	targetEquipData := player.HeroMoudle.CurEquips[req.TargetPos]
	if targetEquipData.ID != req.TargetID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ChangeEquip : Invalid req.ID :%d, equipID:%d", req.TargetID, targetEquipData.ID)
		return
	}

	var sourceEquipData TEquipData
	if req.SourceID != 0 {
		if req.SourcePos < 0 || req.SourcePos >= len(player.BagMoudle.EquipBag.Equips) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeEquip : Invalid SourcePos :%d", req.SourcePos)
			return
		}

		sourceEquipData = player.BagMoudle.EquipBag.Equips[req.SourcePos]
		if sourceEquipData.ID != req.SourceID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeEquip : Invalid SourceID :%d, localid :%d, sourcepos:%d", req.SourceID, sourceEquipData.ID, req.SourcePos)
			return
		}
	}

	if req.TargetID == 0 { //上阵
		pEquipInfo := gamedata.GetEquipmentInfo(sourceEquipData.ID)
		if pEquipInfo == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeEquip : Invalid EquipID :%d", sourceEquipData.ID)
			return
		}

		if (pEquipInfo.Position - 1) != (req.TargetPos % 4) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeEquip : Change To The Wrong Position")
			return
		}

		player.HeroMoudle.CurEquips[req.TargetPos] = sourceEquipData
		player.HeroMoudle.DB_SaveBattleEquipAt(req.TargetPos)
		player.BagMoudle.RemoveEquipAt(req.SourcePos)
		player.BagMoudle.DB_RemoveEquipAt(req.SourcePos)
		//player.BagMoudle.DB_SaveBagEquips()
	} else if req.SourceID == 0 { //下阵
		player.BagMoudle.AddEqiupData(&targetEquipData)
		player.HeroMoudle.CurEquips[req.TargetPos].Clear()
		player.HeroMoudle.DB_SaveBattleEquipAt(req.TargetPos)
	} else {
		pEquipInfo := gamedata.GetEquipmentInfo(sourceEquipData.ID)
		if pEquipInfo == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeEquip : Invalid EquipID :%d", sourceEquipData.ID)
			return
		}

		if (pEquipInfo.Position - 1) != (req.TargetPos % 4) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeEquip : Change To The Wrong Position")
			return
		}

		player.HeroMoudle.CurEquips[req.TargetPos] = sourceEquipData
		player.HeroMoudle.DB_SaveBattleEquipAt(req.TargetPos)
		player.BagMoudle.EquipBag.Equips[req.SourcePos] = targetEquipData
		player.BagMoudle.DB_SaveBagEquipAt(req.SourcePos)
	}

	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS

	//! 判断装备强化品质
	eqiuQuality := 0x7FFFFFFF
	isExist := true
	for i := 0; i < EQUIP_NUM; i++ {
		equi := &player.HeroMoudle.CurEquips[i]
		if equi.ID == 0 {
			isExist = false
			break
		}

		equiData := gamedata.GetEquipmentInfo(equi.ID)

		//! 获取最小品质
		if equiData.Quality < eqiuQuality {
			eqiuQuality = equiData.Quality
		}
	}

	if isExist == true {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_EQUI_QUALITY, eqiuQuality)
	}

	return
}

func Hand_EquipStrengthen(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_EquipStrengthen_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_EquipStrengthen : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_EquipStrengthen_Ack
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

	var pEquipData *TEquipData = nil

	if req.PosType == POSTYPE_BATTLE {
		if req.PosIndex < 0 || req.PosIndex >= EQUIP_NUM {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_EquipStrengthen Error : Invalid posIndex")
			return
		}
		pEquipData = &player.HeroMoudle.CurEquips[req.PosIndex]
	} else if req.PosType == POSTYPE_BAG {
		if req.PosIndex >= len(player.BagMoudle.EquipBag.Equips) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_EquipStrengthen Error : Invalid posIndex")
			return
		}
		pEquipData = &player.BagMoudle.EquipBag.Equips[req.PosIndex]
	}

	if pEquipData.ID != req.EquipID {
		gamelog.Error("Hand_EquipStrengthen Error : Invalid posIndex")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pEquipInfo := gamedata.GetEquipmentInfo(pEquipData.ID)
	if pEquipInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_EquipStrengthen Error : GetEquipmentInfo return nil")
		return
	}

	if pEquipData.StrengLevel >= player.GetLevel()*2 {
		response.RetCode = msg.RE_ALREADY_MAX_LEVEL
		gamelog.Error("Hand_EquipStrengthen Error : Already reach the max level limit")
		return
	}

	costMoney := 0
	costmoneyId := 0
	oldlevel := pEquipData.StrengLevel

	doubleratio := gamedata.GetFuncVipValue(gamedata.FUNC_DOUBLE_POWER_UP, player.GetVipLevel())
	tripleratio := gamedata.GetFuncVipValue(gamedata.FUNC_TRIPLE_POWER_UP, player.GetVipLevel())

	for i := 0; i < req.Times; i++ {
		pEquipStrengCost := gamedata.GetEquipStrengthCostInfo(pEquipData.StrengLevel)
		if pEquipStrengCost == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_EquipStrengthen Error : Invalid pEquipData.StrengLevel :%d", pEquipData.StrengLevel)
			return
		}

		if pEquipData.StrengLevel >= player.GetLevel()*2 {
			break
		}

		costmoneyId = pEquipStrengCost.MoneyID
		tempCost := costMoney + pEquipStrengCost.MoneyNum[pEquipInfo.Quality-1]

		if false == player.RoleMoudle.CheckMoneyEnough(costmoneyId, tempCost) {
			break
		}

		randvalue := utility.Rand() % 1000
		if randvalue < tripleratio {
			pEquipData.StrengLevel += 3
			response.BaoJi = 1
		} else if randvalue < doubleratio {
			pEquipData.StrengLevel += 2
			response.BaoJi = 1
		} else {
			pEquipData.StrengLevel += 1
		}

		costMoney = tempCost
	}

	if pEquipData.StrengLevel > player.GetLevel()*2+1 {
		pEquipData.StrengLevel = player.GetLevel()*2 + 1
	}

	if oldlevel < pEquipData.StrengLevel {
		player.RoleMoudle.CostMoney(costmoneyId, costMoney)
		player.DB_SaveEquipStrength(req.PosType, req.PosIndex)
		response.CostMoney = costMoney
		response.NewLevel = pEquipData.StrengLevel
		response.RetCode = msg.RE_SUCCESS
		if req.PosType == POSTYPE_BATTLE {
			response.FightValue = player.CalcFightValue()
		}
	} else {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_EquipStrengthen Error : Not Enough Money")
		return
	}

	//! 增加日常任务进度
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_EQUI_STRENGTHEN, req.Times)

	//! 判断装备强化等级
	strengthenLevel := player.GetLevel() * 2
	isExist := true
	for i := 0; i < EQUIP_NUM; i++ {
		eqiu := &player.HeroMoudle.CurEquips[i]
		if eqiu.ID == 0 {
			isExist = false
			break
		}

		//! 获取最小等级
		if eqiu.StrengLevel < strengthenLevel {
			strengthenLevel = eqiu.StrengLevel
		}
	}

	if isExist == true {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_EQUI_STRENGTH, strengthenLevel)
	}

	return
}

func Hand_ComposeEquip(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ComposeEquip_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ComposeEquip : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ComposeEquip_Ack
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

	pEquipPieceInfo := gamedata.GetItemInfo(req.EquipPieceID)
	if pEquipPieceInfo == nil {
		gamelog.Error("Hand_ComposeEquip Error : Invalid PieceID :%d", req.EquipPieceID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pEquipInfo := gamedata.GetEquipmentInfo(pEquipPieceInfo.Data1)
	if pEquipInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ComposeEquip Error : Invalid EquipID :%d", pEquipPieceInfo.Data1)
		return
	}

	pieceCount := player.BagMoudle.GetEqiupPieceCount(req.EquipPieceID)
	if pieceCount < pEquipInfo.PieceNum {
		response.RetCode = msg.RE_NOT_ENOUGH_PIECE
		gamelog.Error("Hand_ComposeEquip Error : Not Enough Piece Num :%d", pieceCount)
		return
	}

	player.BagMoudle.AddEqiupByID(pEquipInfo.EquipID)
	player.BagMoudle.RemoveEquipPiece(req.EquipPieceID, pEquipInfo.PieceNum)

	response.EquipID = pEquipInfo.EquipID
	response.RetCode = msg.RE_SUCCESS

	return
}

func Hand_EquipRiseStar(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_EquipRiseStar_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_EquipRiseStar : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_EquipRiseStar_Ack
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

	var pEquipData *TEquipData = nil
	if req.PosType == POSTYPE_BATTLE {
		if req.PosIndex < 0 || req.PosIndex >= EQUIP_NUM {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_EquipRiseStar Error : Invalid posIndex")
			return
		}
		pEquipData = &player.HeroMoudle.CurEquips[req.PosIndex]
	} else if req.PosType == POSTYPE_BAG {
		if req.PosIndex >= len(player.BagMoudle.EquipBag.Equips) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_EquipRiseStar Error : Invalid posIndex")
			return
		}
		pEquipData = &player.BagMoudle.EquipBag.Equips[req.PosIndex]
	}

	if pEquipData == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_EquipRiseStar : pEquipData == nil!!!!")
		return
	}

	pEquipInfo := gamedata.GetEquipmentInfo(pEquipData.ID)
	if pEquipInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_EquipRiseStar Error : Invalid equipid :%d", pEquipData.ID)
		return
	}

	pEquipStarInfo := gamedata.GetEquipStarInfo(pEquipInfo.Quality, pEquipInfo.Position, pEquipData.Star)
	if req.CondIndex == 1 {
		if false == player.RoleMoudle.CheckMoneyEnough(pEquipStarInfo.MoneyID[0], pEquipStarInfo.MoneyNum[0]) || pEquipStarInfo.MoneyID[0] <= 0 {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			gamelog.Error("Hand_EquipRiseStar Error : Not Enough Money")
			return
		}

		player.RoleMoudle.CostMoney(pEquipStarInfo.MoneyID[0], pEquipStarInfo.MoneyNum[0])
		pEquipData.StarCost[0] += int32(pEquipStarInfo.MoneyNum[0])
	} else if req.CondIndex == 2 {
		if false == player.RoleMoudle.CheckMoneyEnough(pEquipStarInfo.MoneyID[1], pEquipStarInfo.MoneyNum[1]) {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			gamelog.Error("Hand_EquipRiseStar Error : Not Enough Money")
			return
		}

		player.RoleMoudle.CostMoney(pEquipStarInfo.MoneyID[1], pEquipStarInfo.MoneyNum[1])
		pEquipData.StarCost[1] += int32(pEquipStarInfo.MoneyNum[1])
	} else if req.CondIndex == 3 {
		if player.BagMoudle.GetEqiupPieceCount(pEquipInfo.PieceID) < pEquipStarInfo.PieceNum {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			gamelog.Error("Hand_EquipRiseStar Error : Not Enough Equipment Piece")
			return
		}
		player.BagMoudle.RemoveEquipPiece(pEquipInfo.PieceID, pEquipStarInfo.PieceNum)
		pEquipData.StarCost[2] += int32(pEquipStarInfo.PieceNum)
	} else {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_EquipRiseStar Error : Invalid CondIndex :%d", req.CondIndex)
		return
	}

	ratio := pEquipStarInfo.Ratio[0]
	if ratio < 1000 {
		for i := 0; i < 3; i++ {
			if pEquipData.StarLuck > pEquipStarInfo.Luck[i] {
				ratio = pEquipStarInfo.Ratio[i+1]
			}
		}

		randvalue := utility.Rand() % 1000
		if randvalue < ratio {
			pEquipData.StarExp += pEquipStarInfo.AddExp
			if pEquipData.StarExp >= pEquipStarInfo.NeedExp {
				pEquipData.StarExp = 0
				pEquipData.Star += 1
				pEquipData.StarLuck = 0
			}
		} else {
			pEquipData.StarLuck += pEquipStarInfo.AddLuck
		}
	} else {
		pEquipData.StarExp += pEquipStarInfo.AddExp
		if pEquipData.StarExp >= pEquipStarInfo.NeedExp {
			pEquipData.StarExp = 0
			pEquipData.Star += 1
		}
		response.FightValue = player.CalcFightValue()
	}

	player.DB_SaveEquipStar(req.PosType, req.PosIndex)

	if pEquipData.Star == 3 {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_EQUI_STAR_THREE, 1)
	} else if pEquipData.Star == 5 {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_EQUI_STAR_FIVE, 1)
	}

	response.Exp = pEquipData.StarExp
	response.Level = pEquipData.Star
	response.Luck = pEquipData.StarLuck
	response.RetCode = msg.RE_SUCCESS
}

func Hand_EquipRefine(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_EquipRefine_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_EquipRefine : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_EquipRefine_Ack
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

	bEnough := player.BagMoudle.IsItemEnough(req.ItemID, req.ItemNum)
	if !bEnough {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_EquipRefine : Not Enough Item ID:%d, Num:%d", req.ItemID, player.BagMoudle.GetNormalItemCount(req.ItemID))
		return
	}

	pItemInfo := gamedata.GetItemInfo(req.ItemID)
	if pItemInfo == nil || pItemInfo.SubType != gamedata.SUB_TYPE_EQUIP_REFINE {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_EquipRefine : Not Equip Strength Item :%d", req.ItemID)
		return
	}

	expCount := pItemInfo.Data1 * req.ItemNum

	var pEquipData *TEquipData = nil
	if req.PosType == POSTYPE_BATTLE {
		if req.PosIndex < 0 || req.PosIndex >= EQUIP_NUM {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_EquipRefine Error : Invalid posIndex")
			return
		}
		pEquipData = &player.HeroMoudle.CurEquips[req.PosIndex]
	} else if req.PosType == POSTYPE_BAG {
		if req.PosIndex >= len(player.BagMoudle.EquipBag.Equips) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_EquipRefine Error : Invalid posIndex")
			return
		}
		pEquipData = &player.BagMoudle.EquipBag.Equips[req.PosIndex]
	}

	if pEquipData.ID != req.EquipID {
		gamelog.Error("Hand_EquipRefine Error : Invalid posIndex")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pEquipInfo := gamedata.GetEquipmentInfo(pEquipData.ID)
	if pEquipInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_EquipRefine Error : Invalid EquipID :%d", pEquipData.ID)
		return
	}

	pEquipRefineCost := gamedata.GetEquipRefineCostInfo(pEquipData.RefineLevel)
	if pEquipRefineCost == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_EquipRefine Error : Invalid pEquipData.RefineLevel")
		return
	}

	pEquipData.RefineExp += expCount

	response.FightValue = 0
	if pEquipData.RefineExp >= pEquipRefineCost.NeedExp[pEquipInfo.Quality-1] {
		pEquipData.RefineExp -= pEquipRefineCost.NeedExp[pEquipInfo.Quality-1]
		pEquipData.RefineLevel += 1

		if req.PosType == POSTYPE_BATTLE {
			response.FightValue = player.CalcFightValue()
		}
	}

	player.DB_SaveEquipRefine(req.PosType, req.PosIndex)
	player.BagMoudle.RemoveNormalItem(req.ItemID, req.ItemNum)

	response.Exp = pEquipData.RefineExp
	response.Level = pEquipData.RefineLevel
	response.RetCode = msg.RE_SUCCESS

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_EQUI_REFINED, req.ItemNum)

	//! 上阵武将精炼等级
	isExist := true
	minLevel := 0x7FFFFFFF
	maxlevel := 0
	for i := 0; i < EQUIP_NUM; i++ {
		equi := &player.HeroMoudle.CurEquips[i]
		if equi.ID == 0 {
			isExist = false
			continue
		}

		if equi.RefineLevel > maxlevel {
			maxlevel = equi.RefineLevel
		}

		if minLevel > equi.RefineLevel {
			minLevel = equi.RefineLevel
		}
	}

	if isExist == true {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_EQUI_REFINED, minLevel)
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_EQUI_REFINED_MAX, maxlevel)

	return
}

func Hand_ChangeGem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ChangeGem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ChangeGem : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ChangeGem_Ack
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

	if req.TargetPos < 0 || req.TargetPos >= GEM_NUM {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ChangeGem : Invalid TargetPos :%d", req.TargetPos)
		return
	}

	targetGemData := player.HeroMoudle.CurGems[req.TargetPos]
	if targetGemData.ID != req.TargetID {
		gamelog.Error("Hand_ChangeGem : Invalid TargetID :%d", req.TargetID)
		return
	}

	var sourceGemData TGemData
	if req.SourceID != 0 {
		if req.SourcePos < 0 || req.SourcePos >= len(player.BagMoudle.GemBag.Gems) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeGem : Invalid SourcePos :%d", req.SourcePos)
			return
		}

		sourceGemData = player.BagMoudle.GemBag.Gems[req.SourcePos]
		if sourceGemData.ID != req.SourceID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeGem : Invalid SourceID :%d, SourcePos:%d, localid:%d", req.SourceID, req.SourcePos, sourceGemData.ID)
			return
		}
	}

	if req.TargetID == 0 { //上阵
		pGemInfo := gamedata.GetGemInfo(sourceGemData.ID)
		if pGemInfo == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeGem : Invalid GemID :%d", sourceGemData.ID)
			return
		}

		if (pGemInfo.Position - 5) != (req.TargetPos % 2) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeGem : Change To The Wrong Position")
			return
		}
		player.HeroMoudle.CurGems[req.TargetPos] = sourceGemData
		player.HeroMoudle.DB_SaveBattleGemAt(req.TargetPos)
		player.BagMoudle.RemoveGemAt(req.SourcePos)
		player.BagMoudle.DB_RemoveGemAt(req.SourcePos)
	} else if req.SourceID == 0 {
		player.BagMoudle.AddGemData(&targetGemData)
		player.HeroMoudle.CurGems[req.TargetPos].Clear()
		player.HeroMoudle.DB_SaveBattleGemAt(req.TargetPos)
	} else {
		pGemInfo := gamedata.GetGemInfo(sourceGemData.ID)
		if pGemInfo == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeGem : Invalid GemID :%d", sourceGemData.ID)
			return
		}

		if (pGemInfo.Position - 5) != (req.TargetPos % 2) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangeGem : Change To The Wrong Position")
			return
		}
		player.HeroMoudle.CurGems[req.TargetPos] = sourceGemData
		player.HeroMoudle.DB_SaveBattleGemAt(req.TargetPos)
		player.BagMoudle.GemBag.Gems[req.SourcePos] = targetGemData
		player.BagMoudle.DB_SaveBagGemAt(req.SourcePos)
	}

	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS

	return
}
func Hand_GemStrengthen(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GemStrengthen_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GemStrengthen : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GemStrengthen_Ack
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

	if len(req.CostGems) <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_GemStrengthen Error : Invalid CostGems Len:%d", len(req.CostGems))
	}

	var pGemData *TGemData = nil
	if req.GemPosType == POSTYPE_BATTLE {
		if req.GemIndex < 0 || req.GemIndex >= GEM_NUM {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemStrengthen Error : Invalid BATTLE posIndex:%d", req.GemIndex)
			return
		}
		pGemData = &player.HeroMoudle.CurGems[req.GemIndex]
	} else if req.GemPosType == POSTYPE_BAG {
		if req.GemIndex >= len(player.BagMoudle.GemBag.Gems) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemStrengthen Error : Invalid BAG posIndex:%d", req.GemIndex)
			return
		}
		pGemData = &player.BagMoudle.GemBag.Gems[req.GemIndex]
	}

	if pGemData.ID != req.GemID {
		gamelog.Error("Hand_GemStrengthen Error : Invalid gemid:%d, %d", pGemData.ID, req.GemID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pTargetGemInfo := gamedata.GetGemInfo(req.GemID)
	if pTargetGemInfo == nil {
		gamelog.Error("Hand_GemStrengthen Error : Invalid gemid:%d", req.GemID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	var tempPos = 10000
	var ExpSum = 0
	for _, t := range req.CostGems {
		pTemData := &player.BagMoudle.GemBag.Gems[t.GemPos]
		if pTemData == nil || pTemData.ID != t.GemID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemStrengthen error :  Invalid costGemID: %d", t.GemID)
			return
		}

		if t.GemPos > tempPos {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemStrengthen error :  Wrong Squence: %d", t.GemPos)
			return
		}

		tempPos = t.GemPos
		pGemInfo := gamedata.GetGemInfo(pTemData.ID)
		ExpSum += pTemData.StrengExp + pGemInfo.Experience

		if pGemInfo.Experience <= 0 {
			gamelog.Error("Hand_GemStrengthen error : gem experience is 0, gemid:%d", pTemData.ID)
		}

		if req.GemPosType == POSTYPE_BAG {
			if t.GemPos == req.GemIndex {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_GemStrengthen error :  Cannot cost gem itself pos:%d", t.GemPos)
				return
			}
		}
	}

	pGemStrengthCostInfo := gamedata.GetGemStrengthCostInfo(pGemData.StrengLevel)
	if pGemStrengthCostInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_GemStrengthen Error : Invalid Gem level :%d", pGemData.StrengLevel)
		return
	}

	if false == player.RoleMoudle.CheckMoneyEnough(pGemStrengthCostInfo.MoneyID, ExpSum*pGemStrengthCostInfo.MoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_GemStrengthen Error : Not Enough Money!! needmoney:%d, hasmoney:%d", ExpSum*pGemStrengthCostInfo.MoneyNum, player.RoleMoudle.GetMoney(pGemStrengthCostInfo.MoneyID))
		return
	}

	player.RoleMoudle.CostMoney(pGemStrengthCostInfo.MoneyID, ExpSum*pGemStrengthCostInfo.MoneyNum)
	response.CostMoneyID = pGemStrengthCostInfo.MoneyID
	response.CostMoneyNum = ExpSum * pGemStrengthCostInfo.MoneyNum
	var oldLevel = pGemData.StrengLevel
	pGemData.StrengExp += ExpSum
	for {
		pGemStrengthCostInfo = gamedata.GetGemStrengthCostInfo(pGemData.StrengLevel)
		if pGemStrengthCostInfo == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemStrengthen Error : Invalid Gem Strengthlevel :%d", pGemData.StrengLevel)
			return
		}
		if pGemData.StrengExp >= pGemStrengthCostInfo.NeedExp[pTargetGemInfo.Quality-1] {
			pGemData.StrengLevel += 1
			pGemData.StrengExp -= pGemStrengthCostInfo.NeedExp[pTargetGemInfo.Quality-1]
		} else {
			break
		}
	}

	response.Exp = pGemData.StrengExp
	response.Level = pGemData.StrengLevel
	response.NewPos = req.GemIndex
	//必须以不影响的索引的方式删除
	for t := 0; t < len(req.CostGems); t++ {
		player.BagMoudle.RemoveGemAt(req.CostGems[t].GemPos)
		if req.GemPosType == POSTYPE_BAG && req.CostGems[t].GemPos < req.GemIndex {
			response.NewPos -= 1
		}
	}
	player.BagMoudle.DB_SaveGemBag()
	if oldLevel < pGemData.StrengLevel && req.GemPosType == POSTYPE_BATTLE {
		response.FightValue = player.CalcFightValue()
	}

	if req.GemPosType == POSTYPE_BATTLE {
		player.HeroMoudle.DB_SaveBattleGemAt(req.GemIndex)
	}

	response.RetCode = msg.RE_SUCCESS
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_GEM_STRENGTHEN, len(req.CostGems))
}

func Hand_GemRefine(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GemRefine_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GemRefine : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GemRefine_Ack
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

	var pGemData *TGemData = nil
	if req.GemPosType == POSTYPE_BATTLE {
		if req.GemIndex < 0 || req.GemIndex >= GEM_NUM {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemRefine Error : Invalid Battle posIndex :%d", req.GemIndex)
			return
		}
		pGemData = &player.HeroMoudle.CurGems[req.GemIndex]
	} else if req.GemPosType == POSTYPE_BAG {
		if req.GemIndex >= len(player.BagMoudle.GemBag.Gems) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemRefine Error : Invalid Bag posIndex:%d", req.GemIndex)
			return
		}
		pGemData = &player.BagMoudle.GemBag.Gems[req.GemIndex]
	}

	if pGemData.ID != req.GemID {
		gamelog.Error("Hand_GemRefine Error : Invalid GemID:%d; %d", pGemData.ID, req.GemID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pGemRefineCostInfo := gamedata.GetGemRefineCostInfo(pGemData.RefineLevel)
	if pGemRefineCostInfo == nil {
		gamelog.Error("Hand_GemRefine Error : Invalid Gem level :%d", pGemData.RefineLevel)
		return
	}

	//判断钱够不够
	if false == player.RoleMoudle.CheckMoneyEnough(pGemRefineCostInfo.MoneyID, pGemRefineCostInfo.MoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_GemRefine Error : Not Enough Money, need:%d, has :%d", pGemRefineCostInfo.MoneyNum, player.RoleMoudle.GetMoney(pGemRefineCostInfo.MoneyID))
		return
	}

	//检查宝物精炼石是否足够
	bEnough := player.BagMoudle.IsItemEnough(pGemRefineCostInfo.ItemID, pGemRefineCostInfo.ItemNum)
	if !bEnough {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_GemRefine Error : Not Enough Items")
		return
	}

	if pGemRefineCostInfo.GemNum > len(req.CostGems) {
		response.RetCode = msg.RE_NOT_ENOUGH_GEM
		gamelog.Error("Hand_GemRefine Error : Not Enough Same Gems")
		return
	}

	//检查同名宝物是否足够
	var tempPos = 10000
	var pTemData *TGemData = nil
	for _, t := range req.CostGems {
		pTemData = &player.BagMoudle.GemBag.Gems[t.GemPos]
		if pTemData == nil || pTemData.ID != t.GemID || t.GemID == 0 {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemRefine error :  Invalid costGemID: %d", t.GemID)
			return
		}

		if t.GemPos > tempPos {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GemRefine error :  Wrong Squence: %d", t.GemPos)
			return
		}

		tempPos = t.GemPos

		if req.GemPosType == POSTYPE_BAG {
			if t.GemPos == req.GemIndex {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_GemRefine error :  Wrong Parameter: %d", t.GemPos)
				return
			}
		}
	}

	pGemData.RefineLevel += 1
	response.Level = pGemData.RefineLevel
	//必须以不影响的索引的方式删除
	for t := 0; t < len(req.CostGems); t++ {
		player.BagMoudle.RemoveGemAt(req.CostGems[t].GemPos)
	}
	player.BagMoudle.DB_SaveGemBag()
	player.RoleMoudle.CostMoney(pGemRefineCostInfo.MoneyID, pGemRefineCostInfo.MoneyNum)
	if req.GemPosType == POSTYPE_BATTLE {
		response.FightValue = player.CalcFightValue()
		player.DB_SaveGemAt(req.GemPosType, req.GemIndex)
	}

	response.CostMoneyID = pGemRefineCostInfo.MoneyID
	response.CostMoneyNum = pGemRefineCostInfo.MoneyNum
	response.RetCode = msg.RE_SUCCESS

	//! 上阵武将宝物精炼等级
	isExist := true
	minLevel := 0x7FFFFFFF
	maxlevel := 0
	for i := 0; i < GEM_NUM; i++ {
		gem := &player.HeroMoudle.CurGems[i]
		if gem.ID == 0 {
			isExist = false
			continue
		}

		if gem.RefineLevel > maxlevel {
			maxlevel = gem.RefineLevel
		}

		if minLevel > gem.RefineLevel {
			minLevel = gem.RefineLevel
		}
	}

	if isExist == true {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_GEM_REFINED, minLevel)
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_GEM_REFINED_MAX, maxlevel)
}

//! 查询分解英雄消耗
func Hand_QueryHeroDecomposeCost(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_QueryHeroDecomposeCost_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_DecomposeHero : Unmarshal error!!!!")
		return
	}

	gamelog.Info("Recv: %s", buffer)

	var response msg.MSG_QueryHeroDecomposeCost_Ack
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

	var resmap map[int]int
	resmap = make(map[int]int)

	for _, t := range req.CostHeros {
		pTempHeroData := player.BagMoudle.GetBagHeroByPos(t.HeroPos)
		if pTempHeroData == nil || pTempHeroData.ID != t.HeroID || t.HeroID == 0 {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_DecomposeHero error :  Invalid SourcePos: %d  HeroID: %d", t.HeroPos, pTempHeroData.ID)
			return
		}

		heroInfo := gamedata.GetHeroInfo(pTempHeroData.ID)

		//还原将魂
		resmap[heroInfo.DecomposeID] += heroInfo.DecomposePrice

		//还原等级材料
		levelInfo := gamedata.GetHeroLevelInfo(pTempHeroData.Quality, pTempHeroData.Level)
		totalExp := levelInfo.TotalNeedExp + pTempHeroData.CurExp

		itemInfo := gamedata.GetItemInfo(gamedata.HeroExpDecomposeItemID)
		itemNum := totalExp / itemInfo.SellPrice

		resmap[gamedata.HeroExpDecomposeItemID] = itemNum
		resmap[levelInfo.MoneyID] += (levelInfo.TotalMoney + pTempHeroData.CurExp*levelInfo.MoneyNum)

		//还原突破材料
		breakInfo := gamedata.GetHeroBreakInfo(pTempHeroData.BreakLevel)
		resmap[breakInfo.MoneyID] += breakInfo.TotalMoneyNum
		resmap[breakInfo.ItemID] += breakInfo.TotalItemNum
		resmap[gamedata.HeroGodDecomposeSoulsID] += breakInfo.TotalHeroNum * heroInfo.DecomposePrice

		//还原培养材料
		resmap[gamedata.CultureItemID] = pTempHeroData.CulturesCost

		//还原天命材料
		destinyLevel := pTempHeroData.DestinyState >> 24 & 0x000F
		if destinyLevel != 0 {
			pDestinyInfo := gamedata.GetHeroDestinyInfo(int(destinyLevel))
			resmap[pDestinyInfo.CostItemID] = pDestinyInfo.Return
		}

		//还原觉醒材料
		if pTempHeroData.WakeLevel != 0 {
			for i := 0; i < pTempHeroData.WakeLevel; i++ {
				wakeInfo := gamedata.GetWakeLevelItem(i)
				for _, v := range wakeInfo.NeedItem {
					if v != 0 {
						resmap[v] += 1
					}
				}

				resmap[wakeInfo.NeedMoneyID] += wakeInfo.NeedMoneyNum
				resmap[wakeInfo.NeedWakeID] += wakeInfo.NeedWakeNum
			}
		}

		//还原化神材料
		if pTempHeroData.GodLevel != 0 {
			godInfo := gamedata.GetHeroGodInfo(pTempHeroData.GodLevel)
			resmap[gamedata.HeroGodDecomposeSoulsID] += godInfo.TotalSouls
			resmap[gamedata.HeroGodDecomposeItemID] += godInfo.TotalItem
			resmap[gamedata.HeroGodDecomposeSoulsID] += ((godInfo.TotalPiece / heroInfo.PieceNum) * heroInfo.DecomposePrice)
			resmap[godInfo.NeedMoneyID] += godInfo.TotalMoney
		}
	}

	//! 发放奖励
	for i, v := range resmap {
		if v != 0 {
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//分解英雄
func Hand_DecomposeHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_DecomposeHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_DecomposeHero : Unmarshal error!!!!")
		return
	}

	gamelog.Info("Recv: %s", buffer)

	var response msg.MSG_DecomposeHero_Ack
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

	if len(req.CostHeros) > 5 || len(req.CostHeros) < 1 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_DecomposeHero Error: Invalid Hero Num :%d", len(req.CostHeros))
		return
	}

	var resmap map[int]int
	resmap = make(map[int]int)

	for _, t := range req.CostHeros {
		pTempHeroData := player.BagMoudle.GetBagHeroByPos(t.HeroPos)
		if pTempHeroData == nil || pTempHeroData.ID != t.HeroID || t.HeroID == 0 {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_DecomposeHero error :  Invalid SourcePos: %d  HeroID: %d", t.HeroPos, pTempHeroData.ID)
			return
		}

		heroInfo := gamedata.GetHeroInfo(pTempHeroData.ID)

		//还原将魂
		resmap[heroInfo.DecomposeID] += heroInfo.DecomposePrice

		//还原等级材料
		levelInfo := gamedata.GetHeroLevelInfo(pTempHeroData.Quality, pTempHeroData.Level)
		totalExp := levelInfo.TotalNeedExp + pTempHeroData.CurExp

		itemInfo := gamedata.GetItemInfo(gamedata.HeroExpDecomposeItemID)
		itemNum := totalExp / itemInfo.SellPrice

		resmap[gamedata.HeroExpDecomposeItemID] = itemNum
		resmap[levelInfo.MoneyID] += (levelInfo.TotalMoney + pTempHeroData.CurExp*levelInfo.MoneyNum)

		//还原突破材料
		breakInfo := gamedata.GetHeroBreakInfo(pTempHeroData.BreakLevel)
		resmap[breakInfo.MoneyID] += breakInfo.TotalMoneyNum
		resmap[breakInfo.ItemID] += breakInfo.TotalItemNum
		resmap[gamedata.HeroGodDecomposeSoulsID] += breakInfo.TotalHeroNum * heroInfo.DecomposePrice

		//还原培养材料
		resmap[gamedata.CultureItemID] = pTempHeroData.CulturesCost

		//还原天命材料
		destinyLevel := pTempHeroData.DestinyState >> 24 & 0x000F
		if destinyLevel != 0 {
			pDestinyInfo := gamedata.GetHeroDestinyInfo(int(destinyLevel))
			resmap[pDestinyInfo.CostItemID] = pDestinyInfo.Return
		}

		//还原觉醒材料
		if pTempHeroData.WakeLevel != 0 {
			for i := 0; i < pTempHeroData.WakeLevel; i++ {
				wakeInfo := gamedata.GetWakeLevelItem(i)
				for _, v := range wakeInfo.NeedItem {
					if v != 0 {
						resmap[v] += 1
					}
				}

				resmap[wakeInfo.NeedMoneyID] += wakeInfo.NeedMoneyNum
				resmap[wakeInfo.NeedWakeID] += wakeInfo.NeedWakeNum
			}
		}

		//还原化神材料
		if pTempHeroData.GodLevel != 0 {
			godInfo := gamedata.GetHeroGodInfo(pTempHeroData.GodLevel)
			resmap[gamedata.HeroGodDecomposeSoulsID] += godInfo.TotalSouls
			resmap[gamedata.HeroGodDecomposeItemID] += godInfo.TotalItem
			resmap[gamedata.HeroGodDecomposeSoulsID] += ((godInfo.TotalPiece / heroInfo.PieceNum) * heroInfo.DecomposePrice)
			resmap[godInfo.NeedMoneyID] += godInfo.TotalMoney
		}
	}

	pos := -1
	for t := 0; t < len(req.CostHeros); t++ {
		if pos >= 0 && req.CostHeros[t].HeroPos > pos {
			req.CostHeros[t].HeroPos -= 1
		}

		pos = req.CostHeros[t].HeroPos
		player.BagMoudle.RemoveHeroAt(req.CostHeros[t].HeroPos)
	}

	player.BagMoudle.DB_SaveHeroBag()

	//! 发放奖励
	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				player.BagMoudle.AddAwardItem(i, v)
				continue
			}
			player.BagMoudle.AddAwardItem(i, v*80/100)
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//查询分解装备消耗
func Hand_QueryEquipDecomposeCost(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_QueryEquipDecomposeCost_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryEquipDecomposeCost : Unmarshal error!!!!")
		return
	}

	gamelog.Info("Recv: %s", buffer)

	var response msg.MSG_QueryEquipDecomposeCost_Ack
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

	var resmap map[int]int
	resmap = make(map[int]int)
	for _, v := range req.CostEquips {
		//! 获取装备信息
		equiInfo := player.BagMoudle.GetEqiupByPos(v.EquipPos)
		if equiInfo == nil || v.EquipID == 0 || equiInfo.ID != v.EquipID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_DecomposeEquip Error: Invalid EquiPos %d", v.EquipPos)
			return
		}

		equiData := gamedata.GetEquipmentInfo(v.EquipID)

		//! 获取强化使用银币
		for i := 1; i < equiInfo.StrengLevel; i++ {
			costInfo := gamedata.GetEquipStrengthCostInfo(i)
			resmap[costInfo.MoneyID] += costInfo.MoneyNum[equiData.Quality-1]
		}

		//! 获取分解威名
		resmap[equiData.SellID[1]] += equiData.SellPrice[1]

		//! 获取升星素材
		resmap[2] += int(equiInfo.StarCost[0])
		resmap[1] += int(equiInfo.StarCost[1])
		resmap[equiData.PieceID] += int(equiInfo.StarCost[2])

		//! 获取精炼
		totalExp := gamedata.GetEquipRefineCostInfo(equiInfo.RefineLevel).TotalExp[equiData.Quality-1]
		totalExp += equiInfo.RefineExp
		itemInfo := gamedata.GetItemInfo(gamedata.EquipRefineDecomposeItemID)
		resmap[gamedata.EquipRefineDecomposeItemID] += totalExp / itemInfo.Data1
	}

	//! 奖励物品
	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, int(v)})
				continue
			}
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, int(v * 80 / 100)})
		}
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 查询分解宠物所得
func Hand_QueryDecomposePetCost(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_QueryPetDecomposeCost_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryDecomposePetCost : Unmarshal error!!!!")
		return
	}

	gamelog.Info("Recv: %s", buffer)

	var response msg.MSG_QueryPetDecomposeCost_Ack
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

	var resmap map[int]int
	resmap = make(map[int]int)

	//! 获取宠物
	petInfo := player.BagMoudle.GetPetByPos(req.PetPos)
	if petInfo.ID != req.PetID || petInfo == nil || req.PetID == 0 {
		gamelog.Error("GetPetByPos Errpr: Invalid pos %d", req.PetPos)
		response.RetCode = msg.RE_SUCCESS
		return
	}

	//! 分解宠物对应兽魂
	petData := gamedata.GetPetInfo(petInfo.ID)
	resmap[gamedata.PetDecomposeSoulsID] += petData.SellPrice

	//! 宠物升级花费
	levelInfo := gamedata.GetPetLevelInfo(petInfo.ID, petInfo.Level)
	totalExp := levelInfo.TotalExp + petInfo.Exp
	totalMoney := levelInfo.TotalMoney + (petInfo.Exp * levelInfo.MoneyNum)

	itemInfo := gamedata.GetItemInfo(gamedata.PetExpDecomposeItemID)
	itemNum := totalExp / itemInfo.Data1

	resmap[gamedata.PetExpDecomposeItemID] += itemNum
	resmap[levelInfo.MoneyID] += totalMoney

	//! 宠物升星花费
	starInfo := gamedata.GetPetStarInfo(petData.Quality, petInfo.Star)
	resmap[starInfo.MoneyID] += starInfo.TotalMoney
	resmap[starInfo.NeedItemID] += starInfo.TotalItemNum
	resmap[petData.PieceID] += starInfo.TotalPiece

	//! 宠物神练
	godInfo := gamedata.GetPetGodInfo(petInfo.ID, petInfo.God)
	totalExp = godInfo.TotalExp + petInfo.GodExp

	itemInfo = gamedata.GetItemInfo(gamedata.PetGodDecomposeItemID)
	itemNum = totalExp / itemInfo.Data1

	resmap[gamedata.PetGodDecomposeItemID] += itemNum

	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				continue
			}
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 分解宠物
func Hand_DecomposePet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_DecomposePet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_DecomposePet : Unmarshal error!!!!")
		return
	}

	gamelog.Info("Recv: %s", buffer)

	var response msg.MSG_DecomposePet_Ack
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

	var resmap map[int]int
	resmap = make(map[int]int)

	//! 获取宠物
	petInfo := player.BagMoudle.GetPetByPos(req.PetPos)
	if petInfo.ID != req.PetID || petInfo == nil || req.PetID == 0 {
		gamelog.Error("GetPetByPos Errpr: Invalid pos %d", req.PetPos)
		response.RetCode = msg.RE_SUCCESS
		return
	}

	//! 分解宠物对应兽魂
	petData := gamedata.GetPetInfo(petInfo.ID)
	resmap[gamedata.PetDecomposeSoulsID] += petData.SellPrice

	//! 宠物升级花费
	levelInfo := gamedata.GetPetLevelInfo(petInfo.ID, petInfo.Level)
	totalExp := levelInfo.TotalExp + petInfo.Exp
	totalMoney := levelInfo.TotalMoney + (petInfo.Exp * levelInfo.MoneyNum)

	itemInfo := gamedata.GetItemInfo(gamedata.PetExpDecomposeItemID)
	if itemInfo == nil {
		gamelog.Error("GetItemInfo Error: Invalid ItemID %d", gamedata.PetExpDecomposeItemID)
		return
	}
	itemNum := totalExp / itemInfo.Data1

	resmap[gamedata.PetExpDecomposeItemID] += itemNum
	resmap[levelInfo.MoneyID] += totalMoney

	//! 宠物升星花费
	starInfo := gamedata.GetPetStarInfo(petData.Quality, petInfo.Star)
	resmap[starInfo.MoneyID] += starInfo.TotalMoney
	resmap[starInfo.NeedItemID] += starInfo.TotalItemNum
	resmap[petData.PieceID] += starInfo.TotalPiece

	//! 宠物神练
	godInfo := gamedata.GetPetGodInfo(petInfo.ID, petInfo.God)
	totalExp = godInfo.TotalExp + petInfo.GodExp

	itemInfo = gamedata.GetItemInfo(gamedata.PetGodDecomposeItemID)
	itemNum = totalExp / itemInfo.Data1

	resmap[gamedata.PetGodDecomposeItemID] += itemNum

	//! 删除宠物
	player.BagMoudle.RemovePetAt(req.PetPos)
	player.BagMoudle.DB_RemovePetAt(req.PetPos)

	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				player.BagMoudle.AddAwardItem(i, v)
				continue
			}
			player.BagMoudle.AddAwardItem(i, v*80/100)
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 查询重生所得
func Hand_QueryHeroRelive(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ReliveHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ReliveHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ReliveHero_Ack
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

	pTargetHero := player.BagMoudle.GetBagHeroByPos(req.HeroPos)
	if pTargetHero.ID != req.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ReliveHero error :  Invalid SourcePos: %d", req.HeroPos)
		return
	}

	//! 判断重生货币是否足够
	var resmap map[int]int
	resmap = make(map[int]int)

	heroInfo := gamedata.GetHeroInfo(pTargetHero.ID)

	//还原将魂
	resmap[heroInfo.PieceID] += heroInfo.PieceNum

	//还原等级材料
	levelInfo := gamedata.GetHeroLevelInfo(pTargetHero.Quality, pTargetHero.Level)
	totalExp := levelInfo.TotalNeedExp + pTargetHero.CurExp

	itemInfo := gamedata.GetItemInfo(gamedata.HeroExpDecomposeItemID)
	itemNum := totalExp / itemInfo.SellPrice

	resmap[gamedata.HeroExpDecomposeItemID] = itemNum
	resmap[levelInfo.MoneyID] += (levelInfo.TotalMoney + pTargetHero.CurExp*levelInfo.MoneyNum)

	//还原突破材料
	breakInfo := gamedata.GetHeroBreakInfo(pTargetHero.BreakLevel)
	resmap[breakInfo.MoneyID] += breakInfo.TotalMoneyNum
	resmap[breakInfo.ItemID] += breakInfo.TotalItemNum
	resmap[heroInfo.PieceID] += heroInfo.PieceNum * breakInfo.TotalHeroNum

	//还原培养材料
	resmap[gamedata.CultureItemID] = pTargetHero.CulturesCost

	//还原天命材料
	destinyLevel := pTargetHero.DestinyState >> 24 & 0x000F
	if destinyLevel != 0 {
		pDestinyInfo := gamedata.GetHeroDestinyInfo(int(destinyLevel))
		resmap[pDestinyInfo.CostItemID] = pDestinyInfo.Return
	}

	//还原觉醒材料
	if pTargetHero.WakeLevel != 0 {
		for i := 0; i < pTargetHero.WakeLevel; i++ {
			wakeInfo := gamedata.GetWakeLevelItem(i)
			for _, v := range wakeInfo.NeedItem {
				if v != 0 {
					resmap[v] += 1
				}
			}

			resmap[wakeInfo.NeedMoneyID] += wakeInfo.NeedMoneyNum
			resmap[wakeInfo.NeedWakeID] += wakeInfo.NeedWakeNum
		}
	}

	//还原化神材料
	if pTargetHero.GodLevel != 0 {
		godInfo := gamedata.GetHeroGodInfo(pTargetHero.GodLevel)
		resmap[gamedata.HeroGodDecomposeSoulsID] += godInfo.TotalSouls
		resmap[gamedata.HeroGodDecomposeItemID] += godInfo.TotalItem
		resmap[heroInfo.PieceID] += godInfo.TotalPiece
		resmap[godInfo.NeedMoneyID] += godInfo.TotalMoney
	}

	//! 发放奖励
	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				continue
			}
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 查询重生所得
func Hand_QueryGemRelive(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ReliveGem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ReliveGem : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ReliveGem_Ack
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

	var resmap map[int]int
	resmap = make(map[int]int)

	//! 获取宝物信息
	gemInfo := player.BagMoudle.GetGemByPos(req.GemPos)
	if gemInfo == nil || req.GemID != gemInfo.ID || req.GemID == 0 {
		gamelog.Error("Hand_ReliveGem Error: Invalid Gempos %d", req.GemPos)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	gemData := gamedata.GetGemInfo(req.GemID)

	//! 给予宝物
	resmap[gemData.ItemID] += 1

	//! 宝物强化花费
	gemStrengInfo := gamedata.GetGemStrengthCostInfo(gemInfo.StrengLevel)
	totalCostExp := gemStrengInfo.TotalExp[gemData.Quality-1] + gemInfo.StrengExp
	itemInfo := gamedata.GetItemInfo(gamedata.GemStrengthDecomposeItemID)
	resmap[gamedata.GemStrengthDecomposeItemID] += (totalCostExp / itemInfo.SellPrice)

	//! 宝物精炼花费
	gemRefineInfo := gamedata.GetGemRefineCostInfo(gemInfo.RefineLevel)
	resmap[gemData.GemID] += gemRefineInfo.TotalGem
	resmap[gamedata.GemRefineDecomposeItemID] += gemRefineInfo.TotalItem
	resmap[gemRefineInfo.MoneyID] += gemRefineInfo.TotalMoney

	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				continue
			}

			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 查询重生所得
func Hand_QueryPetRelive(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_RelivePet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryGemRelive : Unmarshal error!!!!")
		return
	}

	gamelog.Info("Recv: %s", buffer)

	var response msg.MSG_RelivePet_Ack
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

	var resmap map[int]int
	resmap = make(map[int]int)

	//! 获取宠物
	petInfo := player.BagMoudle.GetPetByPos(req.PetPos)
	if petInfo.ID != req.PetID || petInfo == nil || req.PetID == 0 {
		gamelog.Error("GetPetByPos Errpr: Invalid pos %d", req.PetPos)
		response.RetCode = msg.RE_SUCCESS
		return
	}

	//! 重生宠物
	petData := gamedata.GetPetInfo(petInfo.ID)
	resmap[petData.PieceID] += petData.PieceNum

	//! 宠物升级花费
	levelInfo := gamedata.GetPetLevelInfo(petInfo.ID, petInfo.Level)
	totalExp := levelInfo.TotalExp + petInfo.Exp
	totalMoney := levelInfo.TotalMoney + (petInfo.Exp * levelInfo.MoneyNum)

	itemInfo := gamedata.GetItemInfo(gamedata.PetExpDecomposeItemID)
	itemNum := totalExp / itemInfo.Data1

	resmap[gamedata.PetExpDecomposeItemID] += itemNum
	resmap[levelInfo.MoneyID] += totalMoney

	//! 宠物升星花费
	starInfo := gamedata.GetPetStarInfo(petData.Quality, petInfo.Star)
	resmap[starInfo.MoneyID] += starInfo.TotalMoney
	resmap[starInfo.NeedItemID] += starInfo.TotalItemNum
	resmap[petData.PieceID] += starInfo.TotalPiece

	//! 宠物神练
	godInfo := gamedata.GetPetGodInfo(petInfo.ID, petInfo.God)
	totalExp = godInfo.TotalExp + petInfo.GodExp

	itemInfo = gamedata.GetItemInfo(gamedata.PetGodDecomposeItemID)
	itemNum = totalExp / itemInfo.Data1

	resmap[gamedata.PetGodDecomposeItemID] += itemNum

	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				continue
			}
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 查询重生所得
func Hand_QueryEquipRelive(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ReliveEquip_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ReliveEquip : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ReliveEquip_Ack
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

	var resmap map[int]int
	resmap = make(map[int]int)
	//! 获取装备信息
	equiInfo := player.BagMoudle.GetEqiupByPos(req.EquipPos)
	if equiInfo == nil || req.EquipID == 0 || equiInfo.ID != req.EquipID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_DecomposeEquip Error: Invalid EquiPos %d", req.EquipPos)
		return
	}

	equiData := gamedata.GetEquipmentInfo(req.EquipID)

	//! 获取强化使用银币
	for i := 1; i < equiInfo.StrengLevel; i++ {
		costInfo := gamedata.GetEquipStrengthCostInfo(i)
		resmap[costInfo.MoneyID] += costInfo.MoneyNum[equiData.Quality-1]
	}

	//! 获取分解碎片
	resmap[equiData.PieceID] += equiData.PieceNum

	//! 获取升星素材
	resmap[2] += int(equiInfo.StarCost[0])
	resmap[1] += int(equiInfo.StarCost[1])
	resmap[equiData.PieceID] += int(equiInfo.StarCost[2])

	//! 获取精炼
	totalExp := gamedata.GetEquipRefineCostInfo(equiInfo.RefineLevel).TotalExp[equiData.Quality-1]
	totalExp += equiInfo.RefineExp
	itemInfo := gamedata.GetItemInfo(gamedata.EquipRefineDecomposeItemID)
	resmap[gamedata.EquipRefineDecomposeItemID] += totalExp / itemInfo.Data1

	//! 奖励物品
	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				continue
			}
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 重生宠物
func Hand_RelivePet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_RelivePet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_DecomposePet : Unmarshal error!!!!")
		return
	}

	gamelog.Info("Recv: %s", buffer)

	var response msg.MSG_RelivePet_Ack
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

	//! 判断重生货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.RebornCostMoneyID, gamedata.RebornCostMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("CheckMoneyEnough Error: Not enough reborn money")
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(gamedata.RebornCostMoneyID, gamedata.RebornCostMoneyNum)

	var resmap map[int]int
	resmap = make(map[int]int)

	//! 获取宠物
	petInfo := player.BagMoudle.GetPetByPos(req.PetPos)
	if petInfo.ID != req.PetID || petInfo == nil || req.PetID == 0 {
		gamelog.Error("GetPetByPos Errpr: Invalid pos %d", req.PetPos)
		response.RetCode = msg.RE_SUCCESS
		return
	}

	//! 重生宠物
	petData := gamedata.GetPetInfo(petInfo.ID)
	resmap[petData.PieceID] += petData.PieceNum

	//! 宠物升级花费
	levelInfo := gamedata.GetPetLevelInfo(petInfo.ID, petInfo.Level)
	totalExp := levelInfo.TotalExp + petInfo.Exp
	totalMoney := levelInfo.TotalMoney + (petInfo.Exp * levelInfo.MoneyNum)

	itemInfo := gamedata.GetItemInfo(gamedata.PetExpDecomposeItemID)
	itemNum := totalExp / itemInfo.Data1

	resmap[gamedata.PetExpDecomposeItemID] += itemNum
	resmap[levelInfo.MoneyID] += totalMoney

	//! 宠物升星花费
	starInfo := gamedata.GetPetStarInfo(petData.Quality, petInfo.Star)
	resmap[starInfo.MoneyID] += starInfo.TotalMoney
	resmap[starInfo.NeedItemID] += starInfo.TotalItemNum
	resmap[petData.PieceID] += starInfo.TotalPiece

	//! 宠物神练
	godInfo := gamedata.GetPetGodInfo(petInfo.ID, petInfo.God)
	totalExp = godInfo.TotalExp + petInfo.GodExp

	itemInfo = gamedata.GetItemInfo(gamedata.PetGodDecomposeItemID)
	itemNum = totalExp / itemInfo.Data1

	resmap[gamedata.PetGodDecomposeItemID] += itemNum

	//! 删除宠物
	player.BagMoudle.RemovePetAt(req.PetPos)
	player.BagMoudle.DB_RemovePetAt(req.PetPos)

	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				player.BagMoudle.AddAwardItem(i, v)
				continue
			}
			player.BagMoudle.AddAwardItem(i, v*80/100)
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//分解装备
func Hand_DecomposeEquip(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_DecomposeEquip_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_DecomposeEquip : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_DecomposeEquip_Ack
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

	var resmap map[int]int
	resmap = make(map[int]int)

	for _, v := range req.CostEquips {
		//! 获取装备信息
		equiInfo := player.BagMoudle.GetEqiupByPos(v.EquipPos)
		if equiInfo == nil || v.EquipID == 0 || equiInfo.ID != v.EquipID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_DecomposeEquip Error: Invalid EquiPos %d", v.EquipPos)
			return
		}

		equiData := gamedata.GetEquipmentInfo(v.EquipID)

		//! 获取强化使用银币
		for i := 1; i < equiInfo.StrengLevel; i++ {
			costInfo := gamedata.GetEquipStrengthCostInfo(i)
			resmap[costInfo.MoneyID] += costInfo.MoneyNum[equiData.Quality-1]
		}

		//! 获取分解威名
		resmap[equiData.SellID[1]] += equiData.SellPrice[1]

		//! 获取升星素材
		resmap[2] += int(equiInfo.StarCost[0])
		resmap[1] += int(equiInfo.StarCost[1])
		resmap[equiData.PieceID] += int(equiInfo.StarCost[2])

		//! 获取精炼
		totalExp := gamedata.GetEquipRefineCostInfo(equiInfo.RefineLevel).TotalExp[equiData.Quality-1]
		totalExp += equiInfo.RefineExp
		itemInfo := gamedata.GetItemInfo(gamedata.EquipRefineDecomposeItemID)
		resmap[gamedata.EquipRefineDecomposeItemID] += totalExp / itemInfo.Data1
	}

	pos := -1
	for _, item := range req.CostEquips {
		if item.EquipPos > pos && pos >= 0 {
			item.EquipPos -= 1
		}

		pos = item.EquipPos
		player.BagMoudle.RemoveEquipAt(item.EquipPos)
	}
	player.BagMoudle.DB_SaveBagEquips()

	//! 奖励物品
	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				player.BagMoudle.AddAwardItem(i, v)
				continue
			}
			player.BagMoudle.AddAwardItem(i, v*80/100)
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}
	response.RetCode = msg.RE_SUCCESS
}

//重生英雄
func Hand_ReliveHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ReliveHero_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ReliveHero : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ReliveHero_Ack
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

	pTargetHero := player.BagMoudle.GetBagHeroByPos(req.HeroPos)
	if pTargetHero.ID != req.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ReliveHero error :  Invalid SourcePos: %d", req.HeroPos)
		return
	}

	//! 判断重生货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.RebornCostMoneyID, gamedata.RebornCostMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("CheckMoneyEnough Error: Not enough reborn money")
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(gamedata.RebornCostMoneyID, gamedata.RebornCostMoneyNum)

	var resmap map[int]int
	resmap = make(map[int]int)

	heroInfo := gamedata.GetHeroInfo(pTargetHero.ID)

	//还原将魂
	resmap[heroInfo.PieceID] += heroInfo.PieceNum

	//还原等级材料
	levelInfo := gamedata.GetHeroLevelInfo(pTargetHero.Quality, pTargetHero.Level)
	totalExp := levelInfo.TotalNeedExp + pTargetHero.CurExp

	itemInfo := gamedata.GetItemInfo(gamedata.HeroExpDecomposeItemID)
	itemNum := totalExp / itemInfo.SellPrice

	resmap[gamedata.HeroExpDecomposeItemID] = itemNum
	resmap[levelInfo.MoneyID] += (levelInfo.TotalMoney + pTargetHero.CurExp*levelInfo.MoneyNum)

	//还原突破材料
	breakInfo := gamedata.GetHeroBreakInfo(pTargetHero.BreakLevel)
	resmap[breakInfo.MoneyID] += breakInfo.TotalMoneyNum
	resmap[breakInfo.ItemID] += breakInfo.TotalItemNum
	resmap[heroInfo.PieceID] += heroInfo.PieceNum * breakInfo.TotalHeroNum

	//还原培养材料
	resmap[gamedata.CultureItemID] = pTargetHero.CulturesCost

	//还原天命材料
	destinyLevel := pTargetHero.DestinyState >> 24 & 0x000F
	if destinyLevel != 0 {
		pDestinyInfo := gamedata.GetHeroDestinyInfo(int(destinyLevel))
		resmap[pDestinyInfo.CostItemID] = pDestinyInfo.Return
	}

	//还原觉醒材料
	if pTargetHero.WakeLevel != 0 {
		for i := 0; i < pTargetHero.WakeLevel; i++ {
			wakeInfo := gamedata.GetWakeLevelItem(i)
			for _, v := range wakeInfo.NeedItem {
				if v != 0 {
					resmap[v] += 1
				}
			}

			resmap[wakeInfo.NeedMoneyID] += wakeInfo.NeedMoneyNum
			resmap[wakeInfo.NeedWakeID] += wakeInfo.NeedWakeNum
		}
	}

	//还原化神材料
	if pTargetHero.GodLevel != 0 {
		godInfo := gamedata.GetHeroGodInfo(pTargetHero.GodLevel)
		resmap[gamedata.HeroGodDecomposeSoulsID] += godInfo.TotalSouls
		resmap[gamedata.HeroGodDecomposeItemID] += godInfo.TotalItem
		resmap[heroInfo.PieceID] += godInfo.TotalPiece
		resmap[godInfo.NeedMoneyID] += godInfo.TotalMoney
	}

	//! 删除英雄
	player.BagMoudle.RemoveHeroAt(req.HeroPos)
	player.BagMoudle.DB_RemoveHeroAt(req.HeroPos)

	//! 发放奖励
	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				player.BagMoudle.AddAwardItem(i, v)
				continue
			}
			player.BagMoudle.AddAwardItem(i, v*80/100)
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}

	response.RetCode = msg.RE_SUCCESS

	return
}

//重生装备
func Hand_ReliveEquip(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ReliveEquip_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ReliveEquip : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ReliveEquip_Ack
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

	//! 判断重生货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.RebornCostMoneyID, gamedata.RebornCostMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("CheckMoneyEnough Error: Not enough reborn money")
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(gamedata.RebornCostMoneyID, gamedata.RebornCostMoneyNum)

	var resmap map[int]int
	resmap = make(map[int]int)
	//! 获取装备信息
	equiInfo := player.BagMoudle.GetEqiupByPos(req.EquipPos)
	if equiInfo == nil || req.EquipID == 0 || equiInfo.ID != req.EquipID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_DecomposeEquip Error: Invalid EquiPos %d", req.EquipPos)
		return
	}

	equiData := gamedata.GetEquipmentInfo(req.EquipID)

	//! 获取强化使用银币
	for i := 1; i < equiInfo.StrengLevel; i++ {
		costInfo := gamedata.GetEquipStrengthCostInfo(i)
		resmap[costInfo.MoneyID] += costInfo.MoneyNum[equiData.Quality-1]
	}

	//! 获取分解碎片
	resmap[equiData.PieceID] += equiData.PieceNum

	//! 获取升星素材
	resmap[2] += int(equiInfo.StarCost[0])
	resmap[1] += int(equiInfo.StarCost[1])
	resmap[equiData.PieceID] += int(equiInfo.StarCost[2])

	//! 获取精炼
	totalExp := gamedata.GetEquipRefineCostInfo(equiInfo.RefineLevel).TotalExp[equiData.Quality-1]
	totalExp += equiInfo.RefineExp
	itemInfo := gamedata.GetItemInfo(gamedata.EquipRefineDecomposeItemID)
	resmap[gamedata.EquipRefineDecomposeItemID] += totalExp / itemInfo.Data1

	//! 删除装备
	player.BagMoudle.RemoveEquipAt(req.EquipPos)
	player.BagMoudle.DB_RemoveEquipAt(req.EquipPos)

	//! 奖励物品
	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				player.BagMoudle.AddAwardItem(i, v)
				continue
			}
			player.BagMoudle.AddAwardItem(i, v*80/100)
			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
		}
	}
	response.RetCode = msg.RE_SUCCESS
	return
}

//重生宝物
func Hand_ReliveGem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ReliveGem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ReliveGem : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ReliveGem_Ack
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

	//! 判断重生货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.RebornCostMoneyID, gamedata.RebornCostMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("CheckMoneyEnough Error: Not enough reborn money")
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(gamedata.RebornCostMoneyID, gamedata.RebornCostMoneyNum)

	var resmap map[int]int
	resmap = make(map[int]int)

	//! 获取宝物信息
	gemInfo := player.BagMoudle.GetGemByPos(req.GemPos)
	if gemInfo == nil || req.GemID != gemInfo.ID || req.GemID == 0 {
		gamelog.Error("Hand_ReliveGem Error: Invalid Gempos %d", req.GemPos)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	gemData := gamedata.GetGemInfo(req.GemID)

	//! 给予宝物
	resmap[gemData.ItemID] += 1

	//! 宝物强化花费
	gemStrengInfo := gamedata.GetGemStrengthCostInfo(gemInfo.StrengLevel)
	totalCostExp := gemStrengInfo.TotalExp[gemData.Quality-1] + gemInfo.StrengExp
	itemInfo := gamedata.GetItemInfo(gamedata.GemStrengthDecomposeItemID)
	resmap[gamedata.GemStrengthDecomposeItemID] += (totalCostExp / itemInfo.SellPrice)

	//! 宝物精炼花费
	gemRefineInfo := gamedata.GetGemRefineCostInfo(gemInfo.RefineLevel)
	resmap[gemData.GemID] += gemRefineInfo.TotalGem
	resmap[gamedata.GemRefineDecomposeItemID] += gemRefineInfo.TotalItem
	resmap[gemRefineInfo.MoneyID] += gemRefineInfo.TotalMoney

	//! 删除宝物
	player.BagMoudle.RemoveGemAt(req.GemPos)
	player.BagMoudle.DB_RemoveGemAt(req.GemPos)

	for i, v := range resmap {
		if v != 0 {
			if v == 1 {
				response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v})
				player.BagMoudle.AddAwardItem(i, v)
				continue
			}

			response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{i, v * 80 / 100})
			player.BagMoudle.AddAwardItem(i, v*80/100)
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//玩家升品雕文
func Hand_UpgradeDiaoWen(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_UpgradeDiaoWen_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpgradeDiaoWen : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UpgradeDiaoWen_Ack
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

	if req.DiaoWenID <= 0 || req.DiaoWenID > 6 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pTargetHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if (pTargetHeroData == nil) || pTargetHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_UpgradeDiaoWen : req.posType:%d, req.Pos:%d, req.id:%d, targetID:%d", req.TargetHero.PosType,
			req.TargetHero.HeroPos, req.TargetHero.HeroID, pTargetHeroData.ID)
		return
	}

	if pTargetHeroData.DiaoWenQuality[req.DiaoWenID-1] < 2 {
		pTargetHeroData.DiaoWenQuality[req.DiaoWenID-1] = 2
	}

	pDiaoWenItem := gamedata.GetDiaoWenInfo(req.DiaoWenID, int(pTargetHeroData.DiaoWenQuality[req.DiaoWenID-1]))

	//首先雕文是否解锁
	for culture := 0; culture < 5; culture++ {
		if pDiaoWenItem.NeedCulture[culture] > pTargetHeroData.Cultures[culture] {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpgradeDiaoWen : need to culture :%d", culture)
			return
		}
	}

	//等级是否足够
	if pDiaoWenItem.NeedLevel > pTargetHeroData.Level {
		response.RetCode = msg.RE_NOT_ENOUGH_HERO_LEVEL
		gamelog.Error("Hand_UpgradeDiaoWen : Not Enough Level:Need:%d, Has:%d", pDiaoWenItem.NeedLevel, pTargetHeroData.Level)
		return
	}
	//需要人货币是否足够
	if false == player.RoleMoudle.CheckMoneyEnough(pDiaoWenItem.CostMoneyID, pDiaoWenItem.CostMoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_UpgradeDiaoWen : Not Enough money:Need:%d, Has:%d", pDiaoWenItem.CostMoneyNum, player.RoleMoudle.GetMoney(pDiaoWenItem.CostMoneyID))
		return
	}

	//雕文品质加1
	pTargetHeroData.DiaoWenQuality[req.DiaoWenID-1] += 1
	player.RoleMoudle.CostMoney(pDiaoWenItem.CostMoneyID, pDiaoWenItem.CostMoneyNum)
	response.RetCode = msg.RE_SUCCESS
	response.DiaoWenID = req.DiaoWenID
	response.DiaoWenQuality = pTargetHeroData.DiaoWenQuality[req.DiaoWenID-1]
	response.FightValue = player.CalcFightValue()
	player.DB_SaveHeroXiLian(req.TargetHero.PosType, req.TargetHero.HeroPos)

	//! 任务进度
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_DIAOWEN_QUALITY, int(response.DiaoWenQuality))

	return
}

//玩家洗炼雕文
func Hand_XiLianDiaoWen(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_XiLianDiaoWen_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_XiLianDiaoWen : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_XiLianDiaoWen_Ack
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

	if req.DiaoWenID <= 0 || req.DiaoWenID > 6 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pTargetHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if (pTargetHeroData == nil) || pTargetHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_XiLianDiaoWen : req.posType:%d, req.Pos:%d, req.id:%d, targetID:%d", req.TargetHero.PosType,
			req.TargetHero.HeroPos, req.TargetHero.HeroID, pTargetHeroData.ID)
		return
	}

	pDiaoWenItem := gamedata.GetDiaoWenInfo(req.DiaoWenID, int(pTargetHeroData.DiaoWenQuality[req.DiaoWenID-1]))
	pDiaoWenXiLian := gamedata.GetXiLianInfo(req.LockIndex[0] + req.LockIndex[1] + req.LockIndex[2] + req.LockIndex[3])
	if pDiaoWenItem == nil || pDiaoWenXiLian == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_XiLianDiaoWen : pDiaoWenItem == nil || pDiaoWenXiLian == nil")
		return
	}

	//首先雕文是否解锁
	for culture := 0; culture < 5; culture++ {
		if pDiaoWenItem.NeedCulture[culture] > pTargetHeroData.Cultures[culture] {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpgradeDiaoWen : need to culture :%d", culture)
			return
		}
	}

	//需要人货币是否足够
	response.CostMoneyID = pDiaoWenXiLian.FirstMoneyID
	response.CostMoneyNum = pDiaoWenXiLian.FirstMoneyNum
	bMoney := false
	if false == player.BagMoudle.IsItemEnough(pDiaoWenXiLian.FirstMoneyID, pDiaoWenXiLian.FirstMoneyNum) {
		if false == player.RoleMoudle.CheckMoneyEnough(pDiaoWenXiLian.SecondMoneyID, pDiaoWenXiLian.SecondMoneyNum) {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			gamelog.Error("Hand_UpgradeDiaoWen : Not Enough money:Need:%d, Has:%d", pDiaoWenXiLian.SecondMoneyNum, player.RoleMoudle.GetMoney(pDiaoWenXiLian.SecondMoneyID))
			return
		}

		response.CostMoneyID = pDiaoWenXiLian.SecondMoneyID
		response.CostMoneyNum = pDiaoWenXiLian.SecondMoneyNum
		bMoney = true
	}

	if req.LockIndex[0] == 0 {
		pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5] = int32(utility.Rand()%(pDiaoWenItem.Propertys[0].Value[1]-pDiaoWenItem.Propertys[0].Value[0]) + pDiaoWenItem.Propertys[0].Value[0])
		response.RandValue[0] = pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5]
	} else {
		response.RandValue[0] = pTargetHeroData.DiaoWenPtys[(req.DiaoWenID-1)*5]
	}
	if req.LockIndex[1] == 0 {
		pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+1] = int32(utility.Rand()%(pDiaoWenItem.Propertys[1].Value[1]-pDiaoWenItem.Propertys[1].Value[0]) + pDiaoWenItem.Propertys[1].Value[0])
		response.RandValue[1] = pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+1]
		pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+3] = pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+1]
		response.RandValue[3] = pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+3]
	} else {
		response.RandValue[1] = pTargetHeroData.DiaoWenPtys[(req.DiaoWenID-1)*5+1]
		response.RandValue[3] = pTargetHeroData.DiaoWenPtys[(req.DiaoWenID-1)*5+3]
	}
	if req.LockIndex[2] == 0 {
		pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+2] = int32(utility.Rand()%(pDiaoWenItem.Propertys[2].Value[1]-pDiaoWenItem.Propertys[2].Value[0]) + pDiaoWenItem.Propertys[2].Value[0])
		response.RandValue[2] = pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+2]
	} else {
		response.RandValue[2] = pTargetHeroData.DiaoWenPtys[(req.DiaoWenID-1)*5+2]
	}
	if req.LockIndex[3] == 0 {
		pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+4] = int32(utility.Rand()%(pDiaoWenItem.Propertys[4].Value[1]-pDiaoWenItem.Propertys[4].Value[0]) + pDiaoWenItem.Propertys[4].Value[0])
		response.RandValue[4] = pTargetHeroData.DiaoWenBack[(req.DiaoWenID-1)*5+4]
	} else {
		response.RandValue[4] = pTargetHeroData.DiaoWenPtys[(req.DiaoWenID-1)*5+4]
	}

	player.DB_SaveHeroXiLian(req.TargetHero.PosType, req.TargetHero.HeroPos)

	if bMoney {
		player.RoleMoudle.CostMoney(response.CostMoneyID, response.CostMoneyNum)
	} else {
		player.BagMoudle.RemoveNormalItem(response.CostMoneyID, response.CostMoneyNum)
	}

	response.RetCode = msg.RE_SUCCESS
	response.DiaoWenID = req.DiaoWenID

	return
}

//玩家洗炼雕文
func Hand_XiLianTiHuan(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_XiLianTiHuan_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GemRefine : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_XiLianTiHuan_Ack
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

	if req.DiaoWenID <= 0 || req.DiaoWenID > 6 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pTargetHeroData := player.GetHeroByPos(req.TargetHero.PosType, req.TargetHero.HeroPos)
	if (pTargetHeroData == nil) || pTargetHeroData.ID != req.TargetHero.HeroID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_XiLianDiaoWen : req.posType:%d, req.Pos:%d, req.id:%d, targetID:%d", req.TargetHero.PosType,
			req.TargetHero.HeroPos, req.TargetHero.HeroID, pTargetHeroData.ID)
		return
	}

	for i := 0; i < 5; i++ {
		if pTargetHeroData.DiaoWenBack[i+(req.DiaoWenID-1)*5] > 0 {
			pTargetHeroData.DiaoWenPtys[i+(req.DiaoWenID-1)*5] = pTargetHeroData.DiaoWenBack[i+(req.DiaoWenID-1)*5]
			pTargetHeroData.DiaoWenBack[i+(req.DiaoWenID-1)*5] = 0
		}

		response.PropertyValue[i] = pTargetHeroData.DiaoWenPtys[i+(req.DiaoWenID-1)*5]
	}

	response.RetCode = msg.RE_SUCCESS
	response.DiaoWenID = req.DiaoWenID
	response.FightValue = player.CalcFightValue()

	player.DB_SaveHeroXiLian(req.TargetHero.PosType, req.TargetHero.HeroPos)

	return
}

func Hand_UpgradePet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_UpgradePet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpgradePet : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UpgradePet_Ack
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

	if false == gamedata.IsFuncOpen(gamedata.FUNC_BATTLE_PET, player.GetLevel(), player.GetVipLevel()) {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Hand_UpgradePet : Func is not open!!!!")
		return
	}

	var pTargetPetData *TPetData = nil
	if req.PosType == POSTYPE_BATTLE {
		if req.PosIndex < 0 || req.PosIndex > BATTLE_NUM {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpgradePet error : Invalid PosIndex :%d", req.PosIndex)
		}
		pTargetPetData = &player.HeroMoudle.CurPets[req.PosIndex]
	} else if req.PosType == POSTYPE_BAG {
		if req.PosIndex < 0 || req.PosIndex >= len(player.BagMoudle.PetBag.Pets) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpgradePet error : Invalid PosIndex :%d", req.PosIndex)
		}
		pTargetPetData = &player.BagMoudle.PetBag.Pets[req.PosIndex]
	}
	//检验目标宠物是不是正确
	if pTargetPetData == nil || pTargetPetData.ID != req.PetID {
		gamelog.Error("Hand_UpgradePet error : RE_INVALID_PARAM")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//检验目标宠物的等级是不是己经不能进行升级了
	if pTargetPetData.Level >= player.GetLevel() {
		gamelog.Error("Hand_UpgradePet error : pet level can't greater than main hero level")
		response.RetCode = msg.RE_CNT_OVER_MAIN_HERO_LEVEL
		return
	}

	var OldLevel int = pTargetPetData.Level

	//统计消耗道具产生的经验
	var ExpSum = 0
	for i, t := range req.ItemID {
		pItemInfo := gamedata.GetItemInfo(t)
		if pItemInfo == nil {
			gamelog.Error("Hand_UpgradePet error : Invalid Item ID:%d", t)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}
		if player.BagMoudle.IsItemEnough(t, req.ItemNum[i]) == false {
			gamelog.Error("Hand_UpgradePet error : Invalid Item Num:%d", req.ItemNum[i])
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}
		ExpSum += pItemInfo.Data1 * req.ItemNum[i]
	}

	pPetLevelInfo := gamedata.GetPetLevelInfo(pTargetPetData.ID, pTargetPetData.Level)
	if false == player.RoleMoudle.CheckMoneyEnough(pPetLevelInfo.MoneyID, ExpSum*pPetLevelInfo.MoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_UpgradePet Error : Not Enough Money moneyid:%d, moneynum:%d", pPetLevelInfo.MoneyID, ExpSum*pPetLevelInfo.MoneyNum)
		return
	}

	pTargetPetData.Exp += ExpSum

	for {
		pPetLevelInfo = gamedata.GetPetLevelInfo(pTargetPetData.ID, pTargetPetData.Level)
		if pTargetPetData.Exp < pPetLevelInfo.NeedExp {
			break
		}

		pTargetPetData.Exp -= pPetLevelInfo.NeedExp
		pTargetPetData.Level += 1
	}

	player.DB_SavePetLevel(req.PosType, req.PosIndex)
	player.RoleMoudle.CostMoney(pPetLevelInfo.MoneyID, ExpSum*pPetLevelInfo.MoneyNum)

	for i, t := range req.ItemID {
		player.BagMoudle.RemoveNormalItem(t, req.ItemNum[i])
	}

	response.NewLevel = pTargetPetData.Level
	response.NewExp = pTargetPetData.Exp
	response.RetCode = msg.RE_SUCCESS
	response.CostMoney = ExpSum * pPetLevelInfo.MoneyNum

	if req.PosType == POSTYPE_BATTLE && OldLevel < response.NewLevel {
		response.FightValue = player.CalcFightValue()
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PET_LEVEL, response.NewLevel)

	return
}

func Hand_UpstarPet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_UpstarPet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpstarPet : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UpstarPet_Ack
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

	if false == gamedata.IsFuncOpen(gamedata.FUNC_BATTLE_PET, player.GetLevel(), player.GetVipLevel()) {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Hand_UpstarPet : Func is not open!!!!")
		return
	}

	var pTargetPetData *TPetData = nil
	if req.PosType == POSTYPE_BATTLE {
		if req.PosIndex < 0 || req.PosIndex > BATTLE_NUM {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpstarPet error : Invalid PosIndex :%d", req.PosIndex)
		}
		pTargetPetData = &player.HeroMoudle.CurPets[req.PosIndex]
	} else if req.PosType == POSTYPE_BAG {
		if req.PosIndex < 0 || req.PosIndex >= len(player.BagMoudle.PetBag.Pets) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpstarPet error : Invalid PosIndex :%d", req.PosIndex)
		}
		pTargetPetData = &player.BagMoudle.PetBag.Pets[req.PosIndex]
	}
	//检验目标宠物是不是正确
	if pTargetPetData == nil || pTargetPetData.ID != req.PetID {
		gamelog.Error("Hand_UpstarPet error : RE_INVALID_PARAM")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pPetInfo := gamedata.GetPetInfo(pTargetPetData.ID)
	if pPetInfo == nil {
		gamelog.Error("Hand_UpstarPet error : Invalid Pet ID:%d", pTargetPetData.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pPetStarInfo := gamedata.GetPetStarInfo(pPetInfo.Quality, pTargetPetData.Star)
	if pPetStarInfo == nil {
		gamelog.Error("Hand_UpstarPet error : Invalid Pet Quality:%d, Star :%d", pPetInfo.Quality, pTargetPetData.Star)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if pTargetPetData.Level < pPetStarInfo.NeedLevel {
		gamelog.Error("Hand_UpstarPet error : pet level can't greater than main hero level")
		response.RetCode = msg.RE_CNT_OVER_MAIN_HERO_LEVEL
		return
	}

	if false == player.BagMoudle.IsItemEnough(pPetStarInfo.NeedItemID, pPetStarInfo.NeedItemNum) {
		gamelog.Error("Hand_UpstarPet error : Not Enough Item :%d", pPetStarInfo.NeedItemID)
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		return
	}

	if player.BagMoudle.GetPetPieceCount(pPetInfo.PieceID) < pPetStarInfo.PieceNum {
		gamelog.Error("Hand_UpstarPet error : Not Enough Piece Num :%d", pPetInfo.PieceID)
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		return
	}

	if false == player.RoleMoudle.CheckMoneyEnough(pPetStarInfo.MoneyID, pPetStarInfo.MoneyNum) {
		gamelog.Error("Hand_UpstarPet error : Not Enough money :%d", pPetStarInfo.MoneyNum)
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		return
	}

	pTargetPetData.Star += 1
	player.DB_SavePetStar(req.PosType, req.PosIndex)

	player.RoleMoudle.CostMoney(pPetStarInfo.MoneyID, pPetStarInfo.MoneyNum)
	player.BagMoudle.RemoveNormalItem(pPetStarInfo.NeedItemID, pPetStarInfo.NeedItemNum)
	player.BagMoudle.RemovePetPiece(pPetInfo.PieceID, pPetStarInfo.PieceNum)

	response.NewStar = pTargetPetData.Star
	response.CostMoney = pPetStarInfo.MoneyNum
	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PET_STAR, pTargetPetData.Star)
}

func Hand_UpgodPet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_UpgodPet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UpgodPet : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UpgodPet_Ack
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

	if false == gamedata.IsFuncOpen(gamedata.FUNC_BATTLE_PET, player.GetLevel(), player.GetVipLevel()) {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Hand_UpgodPet : Func is not open!!!!")
		return
	}

	if false == gamedata.IsFuncOpen(gamedata.FUNC_BATTLE_PET, player.GetLevel(), player.GetVipLevel()) {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Hand_UpgodPet : Func is not open!!!!")
		return
	}

	var pTargetPetData *TPetData = nil
	if req.PosType == POSTYPE_BATTLE {
		if req.PosIndex < 0 || req.PosIndex > BATTLE_NUM {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpgodPet error : Invalid PosIndex :%d", req.PosIndex)
		}
		pTargetPetData = &player.HeroMoudle.CurPets[req.PosIndex]
	} else if req.PosType == POSTYPE_BAG {
		if req.PosIndex < 0 || req.PosIndex >= len(player.BagMoudle.PetBag.Pets) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_UpgodPet error : Invalid PosIndex :%d", req.PosIndex)
		}
		pTargetPetData = &player.BagMoudle.PetBag.Pets[req.PosIndex]
	}
	//检验目标宠物是不是正确
	if pTargetPetData == nil || pTargetPetData.ID != req.PetID {
		gamelog.Error("Hand_UpgodPet error : RE_INVALID_PARAM")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pItemInfo := gamedata.GetItemInfo(req.ItemID)
	if pItemInfo == nil || pItemInfo.SubType != gamedata.SUB_TYPE_PET_GOD {
		gamelog.Error("Hand_UpgodPet error : RE_INVALID_PARAM")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pTargetPetData.GodExp += pItemInfo.Data1
	oldGod := pTargetPetData.God
	var pPetGodInfo *gamedata.Pet_God
	for {
		pPetGodInfo = gamedata.GetPetGodInfo(pTargetPetData.ID, pTargetPetData.God)
		if pTargetPetData.GodExp < pPetGodInfo.NeedExp {
			break
		}

		pTargetPetData.GodExp -= pPetGodInfo.NeedExp
		pTargetPetData.God += 1
	}

	player.DB_SavePetGod(req.PosType, req.PosIndex)
	player.BagMoudle.RemoveNormalItem(req.ItemID, 1)
	response.Level = pTargetPetData.God
	response.Exp = pTargetPetData.GodExp
	response.RetCode = msg.RE_SUCCESS

	if req.PosType == POSTYPE_BATTLE && oldGod < response.Level {
		response.FightValue = player.CalcFightValue()
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PET_GOD, 1)

	return

}

func Hand_ChangePet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ChangePet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ChangePet : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ChangePet_Ack
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

	if false == gamedata.IsFuncOpen(gamedata.FUNC_BATTLE_PET, player.GetLevel(), player.GetVipLevel()) {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Hand_ChangePet : Func is not open!!!!")
		return
	}

	if req.TargetID == 0 {
		petcount := 0
		for _, itor := range player.HeroMoudle.CurPets {
			if itor.ID != 0 {
				petcount += 1
			}
		}

		if !gamedata.IsFuncOpen(gamedata.FUNC_PET_POS_BEGIN+petcount-1, player.GetLevel(), 0) {
			gamelog.Error("Hand_ChangePet battle pos is not open!")
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}
	}

	if req.TargetPos < 0 || req.TargetPos >= BATTLE_NUM {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ChangePet : Invalid TargetPos :%d", req.TargetPos)
		return
	}

	targetPetData := player.HeroMoudle.CurPets[req.TargetPos]
	if targetPetData.ID != req.TargetID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ChangePet : Invalid TargetID :%d", req.TargetID)
		return
	}

	var sourcePetData TPetData
	var pPetInfo *gamedata.ST_PetInfo = nil
	if req.SourceID != 0 {
		if req.SourcePos < 0 || req.SourcePos >= len(player.BagMoudle.PetBag.Pets) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangePet : Invalid SourcePos :%d", req.SourcePos)
			return
		}

		sourcePetData = player.BagMoudle.PetBag.Pets[req.SourcePos]
		if sourcePetData.ID != req.SourceID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangePet : Invalid SourceID :%d, localid :%d, sourcepos:%d", req.SourceID, sourcePetData.ID, req.SourcePos)
			return
		}

		pPetInfo = gamedata.GetPetInfo(sourcePetData.ID)
		if pPetInfo == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ChangePet : Invalid Source PetID :%d", sourcePetData.ID)
			return
		}
	}

	if req.TargetID == 0 { //上阵
		player.HeroMoudle.CurPets[req.TargetPos] = sourcePetData
		player.HeroMoudle.DB_SaveBattlePetAt(req.TargetPos)
		player.BagMoudle.RemovePetAt(req.SourcePos)
		player.BagMoudle.DB_RemovePetAt(req.SourcePos)
		//player.BagMoudle.DB_SavePetBag()
	} else if req.SourceID == 0 { //下阵
		player.BagMoudle.AddPetData(&targetPetData)
		player.HeroMoudle.CurPets[req.TargetPos].Clear()
		player.HeroMoudle.DB_SaveBattlePetAt(req.TargetPos)
	} else {
		player.HeroMoudle.CurPets[req.TargetPos] = sourcePetData
		player.HeroMoudle.DB_SaveBattlePetAt(req.TargetPos)
		player.BagMoudle.PetBag.Pets[req.SourcePos] = targetPetData
		player.BagMoudle.DB_SaveBagPetAt(req.SourcePos)
	}

	response.FightValue = player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS

}

func Hand_UnsetPet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_UnsetPet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UnsetPet : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UnsetPet_Ack
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

	if false == gamedata.IsFuncOpen(gamedata.FUNC_BATTLE_PET, player.GetLevel(), player.GetVipLevel()) {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Hand_UnsetPet : Func is not open!!!!")
		return
	}

	tempTarget := player.HeroMoudle.CurPets[req.TargetPos]
	if tempTarget.ID != req.TargetID {
		gamelog.Error("Hand_UnsetPet error : RE_INVALID_PARAM")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//将宠物装到背包中
	player.BagMoudle.AddPetData(&tempTarget)
	player.BagMoudle.DB_AddPetAtLast(false)

	player.HeroMoudle.CurPets[req.TargetPos].Clear()
	player.DB_SavePetAt(POSTYPE_BATTLE, req.TargetPos)

	response.RetCode = msg.RE_SUCCESS
	response.FightValue = player.CalcFightValue()
}

func Hand_ComposePet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ComposePet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ComposePet : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ComposePet_Ack
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

	pPetPieceInfo := gamedata.GetItemInfo(req.PetPieceID)
	if pPetPieceInfo == nil {
		gamelog.Error("Hand_ComposePet Error : Invalid PieceID :%d", req.PetPieceID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pPetInfo := gamedata.GetPetInfo(pPetPieceInfo.Data1)
	if pPetInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ComposePet Error : Invalid PetID :%d", pPetPieceInfo.Data1)
		return
	}

	pieceCount := player.BagMoudle.GetPetPieceCount(req.PetPieceID)
	if pieceCount < pPetInfo.PieceNum {
		response.RetCode = msg.RE_NOT_ENOUGH_PIECE
		gamelog.Error("Hand_ComposePet Error : Not Enough Piece Num :%d", pieceCount)
		return
	}

	//增加宠物图鉴功能
	bNew := true
	for i := 0; i < len(player.BagMoudle.ColPets); i++ {
		if player.BagMoudle.ColPets[i] == int16(pPetInfo.PetID) {
			bNew = false
			break
		}
	}

	player.BagMoudle.AddPetByID(pPetInfo.PetID)
	player.BagMoudle.RemovePetPiece(req.PetPieceID, pPetInfo.PieceNum)
	response.PetID = pPetInfo.PetID
	response.RetCode = msg.RE_SUCCESS

	if !bNew {
		return
	}

	bAdd := false
	for j := 0; j < len(gamedata.GT_PetMap_List); j++ {
		for k := 0; k < 3; k++ {
			if gamedata.GT_PetMap_List[j].PetIds[k] == pPetInfo.PetID && gamedata.GT_PetMap_List[j].IsMapOK(player.BagMoudle.ColPets) {
				bAdd = true
				for _, n := range gamedata.GT_PetMap_List[j].Buffs {
					if n.PropertyID != 0 {
						player.HeroMoudle.AddExtraProperty(n.PropertyID, int32(n.Value), n.IsPercent, 0)
					}
				}
			}
		}
	}

	if bAdd {
		player.HeroMoudle.DB_SaveExtraProperty()
	}

	return
}

//装备时装
func Hand_FashionSet(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_FashionSet_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FashionSet : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_FashionSet_Ack
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

	nIndex := -1
	//先检测背包是否存在时装
	for i := 0; i < len(player.BagMoudle.FashionBag.Fashions); i++ {
		if player.BagMoudle.FashionBag.Fashions[i].ID == req.FashionID {
			nIndex = i
			break
		}
	}

	//索引小于0， 背包里不存在时装，不能装备
	if nIndex < 0 {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_FashionSet Error: NO Fashion :%d", req.FashionID)
		return
	}

	player.HeroMoudle.FashionID = player.BagMoudle.FashionBag.Fashions[nIndex].ID
	player.HeroMoudle.FashionLvl = player.BagMoudle.FashionBag.Fashions[nIndex].Level
	player.HeroMoudle.DB_SaveFashionInfo()
	response.RetCode = msg.RE_SUCCESS
	response.FightValue = player.HeroMoudle.CalcFightValue(nil)
}

//时装强化
func Hand_FashionStrength(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_FashionStrength_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FashionStrength Error: Unmarshal error!!!!")
		return
	}

	var response msg.MSG_FashionStrength_Ack
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

	nIndex := -1
	//先检测是否存在时装
	for i := 0; i < len(player.BagMoudle.FashionBag.Fashions); i++ {
		if player.BagMoudle.FashionBag.Fashions[i].ID == req.FashionID {
			nIndex = i
			break
		}
	}

	if nIndex < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_FashionStrength Error: No this Fashion %d", req.FashionID)
		return
	}

	pFashionInfo := gamedata.GetFashionInfo(req.FashionID)
	if pFashionInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_FashionStrength Error: Invalid Fashion ID %d", req.FashionID)
		return
	}

	//检查需要道具

	pFashionLevel := gamedata.GetFashionLevelInfo(pFashionInfo.ID, player.BagMoudle.FashionBag.Fashions[nIndex].Level)
	if pFashionLevel == nil {
		gamelog.Error("Hand_FashionStrength Error: Invalid Fashion quality %d, level :%d", pFashionInfo.Quality, player.BagMoudle.FashionBag.Fashions[nIndex].Level)
		return
	}

	//扣除需要的道具
	if false == player.BagMoudle.IsItemEnough(pFashionLevel.CostItemID, pFashionLevel.CostItemNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_FashionStrength Error: Not Enough Item %d", pFashionLevel.CostItemID)
		return
	}

	if false == player.RoleMoudle.CheckMoneyEnough(pFashionLevel.CostMoneyID, pFashionLevel.CostMoneyNum) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_FashionStrength Error: Not Enough Money %d", pFashionLevel.CostMoneyID)
		return
	}

	player.BagMoudle.FashionBag.Fashions[nIndex].Level += 1
	player.BagMoudle.DB_SaveFashionAt(nIndex)
	player.BagMoudle.RemoveNormalItem(pFashionLevel.CostItemID, pFashionLevel.CostItemNum)
	player.RoleMoudle.CostMoney(pFashionLevel.CostMoneyID, pFashionLevel.CostMoneyNum)
	if player.HeroMoudle.FashionID == player.BagMoudle.FashionBag.Fashions[nIndex].ID {
		player.HeroMoudle.FashionLvl = player.BagMoudle.FashionBag.Fashions[nIndex].Level
		player.HeroMoudle.DB_SaveFashionInfo()
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_FASHION_STRENGTH, 1)

}

//时装重铸
func Hand_FashionRecast(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_FashionRecast_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FashionRecast : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_FashionRecast_Ack
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

	nIndex := -1
	//先检测是否存在时装
	for i := 0; i < len(player.BagMoudle.FashionBag.Fashions); i++ {
		if player.BagMoudle.FashionBag.Fashions[i].ID == req.FashionID {
			nIndex = i
			break
		}
	}

	if nIndex < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_FashionRecast Error: No this Fashion %d", req.FashionID)
		return
	}

	pFashionInfo := gamedata.GetFashionInfo(req.FashionID)
	if pFashionInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_FashionRecast Error: Invalid Fashion ID %d", req.FashionID)
		return
	}

	pFashionLevel := gamedata.GetFashionLevelInfo(pFashionInfo.ID, player.BagMoudle.FashionBag.Fashions[nIndex].Level)
	if pFashionLevel == nil {
		gamelog.Error("Hand_FashionRecast Error: Invalid Fashion quality %d, level :%d", pFashionInfo.Quality, player.BagMoudle.FashionBag.Fashions[nIndex].Level)
		return
	}
}

func IsFashionMapOK(pMap *gamedata.ST_FashionMapInfo, Fashions []TFashionData) bool {
	for i := 0; i < 3; i++ {
		bFind := false
		for j := 0; j < len(Fashions); j++ {
			if pMap.FashionIds[i] == int(Fashions[j].ID) || pMap.FashionIds[i] == 0 {
				bFind = true
				break
			}
		}

		if bFind == false {
			return false
		}
	}

	return true
}

//时装合成
func Hand_FashionCompose(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_FashionCompose_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FashionCompose : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_FashionCompose_Ack
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

	for i := 0; i < len(player.BagMoudle.FashionBag.Fashions); i++ {
		if player.BagMoudle.FashionBag.Fashions[i].ID == req.FashionID {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			gamelog.Error("Hand_FashionCompose Error: Fashion %d already exist", req.FashionID)
			return
		}
	}

	pFashionInfo := gamedata.GetFashionInfo(req.FashionID)
	if pFashionInfo == nil {
		gamelog.Error("Hand_FashionCompose Error: Invalid Fashion ID %d", req.FashionID)
		return
	}

	if player.BagMoudle.GetFashionPieceCount(pFashionInfo.PieceID) < pFashionInfo.PieceNum {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_FashionCompose Error: Not Enouth Piece Num %d", pFashionInfo.PieceNum)
		return
	}

	player.BagMoudle.AddFashionByID(req.FashionID)
	player.BagMoudle.RemoveFashionPiece(pFashionInfo.PieceID, pFashionInfo.PieceNum)

	response.RetCode = msg.RE_SUCCESS
	response.FashionID = req.FashionID

	bAdd := false
	for j := 0; j < len(gamedata.GT_FashionMap_List); j++ {
		for k := 0; k < 3; k++ {
			if gamedata.GT_FashionMap_List[j].FashionIds[k] == req.FashionID && IsFashionMapOK(&gamedata.GT_FashionMap_List[j], player.BagMoudle.FashionBag.Fashions) {
				bAdd = true
				for _, n := range gamedata.GT_FashionMap_List[j].Buffs {
					if n.PropertyID != 0 {
						player.HeroMoudle.AddExtraProperty(n.PropertyID, int32(n.Value), n.IsPercent, 0)
					}
				}
			}
		}
	}

	if bAdd {
		player.HeroMoudle.DB_SaveExtraProperty()
	}

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_FASHION_COMPOSE, pFashionInfo.Quality)
}

//时装熔炼
func Hand_FashionMelting(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_FashionMelting_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_FashionMelting : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_FashionMelting_Ack
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

	MeltingValue := 0
	//检查所有有碎片是不是足够
	for i := 0; i < len(req.PieceIDs); i++ {
		if player.BagMoudle.GetFashionPieceCount(req.PieceIDs[i]) < req.PieceNums[i] {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			gamelog.Error("Hand_FashionMelting : Not Enought Piece Item!!!!")
			return
		}

		pItemInfo := gamedata.GetItemInfo(req.PieceIDs[i])
		if pItemInfo == nil {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_FashionMelting : Invalid PieceID:%d!!!!", req.PieceIDs[i])
			return
		}

		MeltingValue += pItemInfo.Data2
	}

	if MeltingValue < gamedata.FashionMeltingSum {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_FashionMelting : Not Enought Piece Item!!!!")
		return
	}

	ret := gamedata.GetItemsFromAwardID(gamedata.FashionMeltingAwardID)
	if ret == nil || len(ret) <= 0 {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_FashionMelting : GetItemsFromAwardID Failed!!!!")
		return
	}
	bFind := false
	for i := 0; i < len(player.BagMoudle.FashionBag.Fashions); i++ {
		if ret[0].ItemID == player.BagMoudle.FashionBag.Fashions[i].ID {
			bFind = true
		}
	}

	if bFind == true {

	}

}

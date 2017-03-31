package mainlogic

import (
	"bytes"
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"time"
	"utility"
)

func Hand_TestGetAction(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetTestAction_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetTestMoney : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetTestAction_Ack
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

	for i := 1; i <= len(player.RoleMoudle.Actions); i++ {
		player.RoleMoudle.AddAction(i, 100)
	}

	response.RetCode = msg.RE_SUCCESS
	response.Actions = make([]int, len(player.RoleMoudle.Actions))
	for i := 0; i < len(player.RoleMoudle.Actions); i++ {
		response.Actions[i] = player.RoleMoudle.Actions[i].Value
	}

	return
}

func Hand_TestAddCharge(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_ChargeTestMoney_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetTestMoney : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ChargeTestMoney_Ack
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

	player.HandChargeRenMinBi(req.RMB, req.ChargeID)

	response.RetCode = msg.RE_SUCCESS
	response.VIPExp = player.GetVipExp()
	response.VIPLevel = player.GetVipLevel()
}

func Hand_TestGetMoney(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetTestMoney_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetTestMoney : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetTestMoney_Ack
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

	for i := 1; i < 14; i++ {
		player.RoleMoudle.AddMoney(i, 100000)
	}

	response.RetCode = msg.RE_SUCCESS
	response.Moneys = player.RoleMoudle.Moneys

	//! 玩家领取激活码
	//! 消息: /recv_gift_code
	type MSG_RecvGiftCode_Req struct {
		PlayerID   int32
		SessionKey string
		GiftCode   string //激活码
	}

	type MSG_RecvGiftCode_Ack struct {
		RetCode   int                //返回码
		AwardItem []msg.MSG_ItemData //! 奖励
	}

	var createGift_req MSG_RecvGiftCode_Req
	createGift_req.PlayerID = req.PlayerID
	createGift_req.SessionKey = req.SessionKey
	createGift_req.GiftCode = "2389b788d779a87dc5dec98da"

	b, _ := json.Marshal(createGift_req)
	http.DefaultClient.Timeout = time.Second * 3
	url := "http://192.168.0.167:8082/recv_gift_code"
	httpRet, err := http.Post(url, "text/HTML", bytes.NewReader(b))
	if err != nil {
		gamelog.Error("err : %s !!!!", err.Error())
		return
	}

	httpRetBuf := make([]byte, httpRet.ContentLength)
	httpRet.Body.Read(httpRetBuf)
	httpRet.Body.Close()

	var createGift_ack msg.MSG_RecvGiftCode_Ack
	err = json.Unmarshal(httpRetBuf, &createGift_ack)
	if err != nil {
		gamelog.Error("Hand_RecvGiftCode Unmarshal createGift_ack fail, Error: %s", err.Error())
		return
	}

	if createGift_ack.RetCode != msg.RE_SUCCESS {
		gamelog.Error("Hand_RecvGiftCode Error, createGift_ack.RetCode:%d", createGift_ack.RetCode)
		return
	}

	gamelog.Info("Get Award: %v", createGift_ack.AwardItem)

	return
}

func Hand_TestAddVip(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_TestAddVip_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetTestAward : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_TestAddVip_Ack
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

	//! 调用接口增加VIP经验
	player.RoleMoudle.AddVipExp(100)

	response.RetCode = msg.RE_SUCCESS
	response.VipExp = player.GetVipExp()
	response.VipLevel = player.GetVipLevel()
}

func Hand_TestAddGuildExp(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_TestAddGuild_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_TestAddGuildExp : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_TestAddGuild_Ack
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
	guild.AddExp(10000)
	response.RetCode = msg.RE_SUCCESS
	response.GuildExp = guild.CurExp
	response.GuildLevel = guild.Level
}

func Hand_TestCompress(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//defer func() {
	b := []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u'}
	cp := utility.CompressData(b)
	gamelog.Error("Hand_TestCompress : len:%d", len(cp))
	//var test []byte = []byte{31, 139, 8, 0, 0, 9, 110, 136, 0, 255, 98, 96, 100, 98, 102, 97, 101, 99, 231, 224, 4, 0, 0, 0, 255, 255}
	w.Write(cp)

	//w.Write(b)
	//}()

	return
}

func Hand_GetHerosProperty(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_TestHerosProperty_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetHerosProperty Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_TestHerosProperty_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = GetPlayerByID(req.PlayerID)
	if player == nil {
		response.RetCode = msg.RE_INVALID_PLAYERID
		gamelog.Error("Hand_QueryCampBatInfo Error: Invalid PlayerID", req.PlayerID)
		return
	}

	response.BattleCamp = player.CamBattleModule.BattleCamp
	response.Level = player.GetLevel()
	response.PlayerID = player.playerid

	var HeroResults = make([]THeroResult, BATTLE_NUM)
	player.HeroMoudle.CalcFightValue(HeroResults)

	for i := 0; i < BATTLE_NUM; i++ {
		response.Heros[i].HeroID = HeroResults[i].HeroID
		response.Heros[i].PropertyValue = HeroResults[i].PropertyValues
		response.Heros[i].PropertyPercent = HeroResults[i].PropertyPercents
		response.Heros[i].CampDef = HeroResults[i].CampDef
		response.Heros[i].CampKill = HeroResults[i].CampKill
	}

	response.RetCode = msg.RE_SUCCESS
}

func Hand_TestUplevel(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_TestUpLevel_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_TestUplevel : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_TestUpLevel_Ack
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

	player.HeroMoudle.CurHeros[0].Level += 1
	if player.HeroMoudle.CurHeros[0].Level > gamedata.G_HeroMaxLevel {
		player.HeroMoudle.CurHeros[0].Level = gamedata.G_HeroMaxLevel
	}
	G_LevelRanker.SetRankItem(req.PlayerID, player.HeroMoudle.CurHeros[0].Level)
	player.DB_SaveHeroLevelExp(POSTYPE_BATTLE, 0)
	response.RetCode = msg.RE_SUCCESS
	response.RetLevel = player.HeroMoudle.CurHeros[0].Level
	player.CalcFightValue()

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_LEVEL_UP, 1)
	player.ActivityModule.LevelGift.CheckLevelUp(response.RetLevel)
	return

}

func Hand_TestUplevelTen(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_TestUpLevelTen_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_TestUplevelTen : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_TestUpLevelTen_Ack
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

	player.HeroMoudle.CurHeros[0].Level += 10
	if player.HeroMoudle.CurHeros[0].Level > gamedata.G_HeroMaxLevel {
		player.HeroMoudle.CurHeros[0].Level = gamedata.G_HeroMaxLevel
	}
	G_LevelRanker.SetRankItem(req.PlayerID, player.HeroMoudle.CurHeros[0].Level)
	player.DB_SaveHeroLevelExp(POSTYPE_BATTLE, 0)
	player.CalcFightValue()
	response.RetCode = msg.RE_SUCCESS
	response.RetLevel = player.HeroMoudle.CurHeros[0].Level

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_LEVEL_UP, 10)
	player.ActivityModule.LevelGift.CheckLevelUp(response.RetLevel)
	return

}

func Hand_TestAddItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_TestAddItem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_TestAddItem : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_TestAddItem_Ack
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
	player.BagMoudle.AddAwardItem(req.ItemID, req.AddNum)

	response.Count = req.AddNum
	response.RetCode = msg.RE_SUCCESS
}

func Hand_TestPassCopy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetPassCopy_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_TestPassCopy : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetPassCopy_Ack
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

	if req.CopyType == gamedata.COPY_TYPE_Main {
		//! 获取章节最后关卡
		chapterInfo := gamedata.GetMainChapterInfo(req.Chapter)
		if chapterInfo == nil {
			gamelog.Error("GetMainChapterInfo Error: Not exist chapter %d", req.Chapter)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}

		player.CopyMoudle.Main.CurChapter = req.Chapter
		player.CopyMoudle.Main.CurID = chapterInfo.EndID
		player.CopyMoudle.DB_UpdateMainCopyInfo()

	} else if req.CopyType == gamedata.COPY_TYPE_Elite {
		//! 获取章节最后关卡
		chapterInfo := gamedata.GetEliteChapterInfo(req.Chapter)
		if chapterInfo == nil {
			gamelog.Error("GetEliteChapterInfo Error: Not exist chapter %d", req.Chapter)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}

		player.CopyMoudle.Elite.CurChapter = req.Chapter
		player.CopyMoudle.Elite.CurID = chapterInfo.EndID
		player.CopyMoudle.DB_UpdateEliteCopyInfo()

	}

	response.RetCode = msg.RE_SUCCESS
}

package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

//! 选择玩家简要信息
func Hand_SelectTargetPlayer(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GameSelectPlayer_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SelectTargetPlayer : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GameSelectPlayer_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	pTarget := GetSelectPlayer(SelectScoreTarget, 1000)
	if pTarget != nil && pTarget.playerid != 0 {
		response.RetCode = msg.RE_SUCCESS
		response.Target.FightValue = pTarget.GetFightValue()
		response.Target.HeroID = pTarget.HeroMoudle.CurHeros[0].HeroID
		response.Target.Level = pTarget.HeroMoudle.CurHeros[0].Level
		response.Target.Name = pTarget.RoleMoudle.Name
		response.Target.PlayerID = pTarget.playerid
		response.Target.SvrID = GetCurServerID()
		response.Target.SvrName = GetCurServerName()
	} else {
		response.RetCode = msg.RE_SELECT_PLAYRE_FAILED
		gamelog.Error("Hand_SelectTargetPlayer Error : Cant Select Player!!")
	}

	return
}

//! 提交战斗目标信息
func Hand_GetFightTarget(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetFightTarget_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFightTarget : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetFightTarget_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	pTarget := GetPlayerByID(req.PlayerID)
	if pTarget == nil {
		gamelog.Error("Hand_GetFightTarget Error  Cant get player %d data", req.PlayerID)
		return
	}

	var HeroResults = make([]THeroResult, BATTLE_NUM)
	response.PlayerData.FightValue = pTarget.HeroMoudle.CalcFightValue(HeroResults)
	response.PlayerData.Quality = pTarget.HeroMoudle.CurHeros[0].Quality
	for i := 0; i < BATTLE_NUM; i++ {
		response.PlayerData.Heros[i].HeroID = HeroResults[i].HeroID
		response.PlayerData.Heros[i].PropertyValue = HeroResults[i].PropertyValues
		response.PlayerData.Heros[i].PropertyPercent = HeroResults[i].PropertyPercents
		response.PlayerData.Heros[i].CampDef = HeroResults[i].CampDef
		response.PlayerData.Heros[i].CampKill = HeroResults[i].CampKill
	}

	response.RetCode = msg.RE_SUCCESS

	return
}

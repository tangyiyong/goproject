package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"gamesvr/tcpclient"
	"msg"
	"time"
)

func Hand_Connect(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: Hand_Connect")
	SendCheckInMsg(pTcpConn)

	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		return
	}

	if pClient.ConType == tcpclient.CON_TYPE_CHAT {

	} else if pClient.ConType == tcpclient.CON_TYPE_BATSVR {
		SetBattleSvrConnectOK(pClient.SvrID, true)
	}

	return
}

func Hand_DisConnect(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: Hand_DisConnect")

	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		return
	}

	if pClient.ConType == tcpclient.CON_TYPE_CHAT {

	} else if pClient.ConType == tcpclient.CON_TYPE_BATSVR {
		SetBattleSvrConnectOK(pClient.SvrID, false)
	}

	return
}

func Hand_KillEventReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_KILL_EVENT_REQ")

	var req msg.MSG_KillEvent_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_KillEventReq : Message Reader Error!!!!")
		return
	}

	pPlayer := GetPlayerByID(req.Killer)
	if pPlayer == nil {
		gamelog.Error("Hand_KillEventReq : Invalid PlayerID :%d!!!!", req.Killer)
		return
	}

	var response msg.MSG_KillEvent_Ack

	pPlayer.CamBattleModule.Kill += req.Kill
	response.CurRank = G_CampBat_TodayKill.SetRankItem(req.Killer, pPlayer.CamBattleModule.Kill)
	G_CampBat_KillSum.SetRankItem(req.Killer, pPlayer.CamBattleModule.Kill)
	G_CampBat_CampKill[pPlayer.CamBattleModule.BattleCamp-1].SetRankItem(req.Killer, pPlayer.CamBattleModule.Kill)
	if req.Destroy > 0 {
		pPlayer.CamBattleModule.Destroy += req.Destroy
		pPlayer.CamBattleModule.DestroySum += req.Destroy
		G_CampBat_TodayDestroy.SetRankItem(req.Killer, pPlayer.CamBattleModule.Destroy)
		G_CampBat_DestroySum.SetRankItem(req.Killer, pPlayer.CamBattleModule.DestroySum)
		G_CampBat_CampDestroy[pPlayer.CamBattleModule.BattleCamp-1].SetRankItem(req.Killer, pPlayer.CamBattleModule.Kill)
	}

	if req.SeriesKill == gamedata.CampBat_NtyKillNum {
		var nty msg.MSG_HorseLame_Notify
		nty.TextType = TextCampBatHorseLamp
		nty.Camps = append(nty.Camps, pPlayer.CamBattleModule.BattleCamp)
		nty.Params = append(nty.Params, pPlayer.RoleMoudle.Name)
		b, _ := json.Marshal(nty)
		SendMessageToClient(0, msg.MSG_HORSELAME_NOTIFY, b)
	}

	if pPlayer.CamBattleModule.KillHonor < gamedata.CampBat_KillHonorMax {
		AddHonor := req.Kill * gamedata.Campbat_KillHonorOne
		if (AddHonor + pPlayer.CamBattleModule.KillHonor) <= gamedata.CampBat_KillHonorMax {
			pPlayer.CamBattleModule.KillHonor += AddHonor
			pPlayer.RoleMoudle.AddMoney(4, AddHonor)
		} else {
			pPlayer.RoleMoudle.AddMoney(4, gamedata.CampBat_KillHonorMax-pPlayer.CamBattleModule.KillHonor)
			pPlayer.CamBattleModule.KillHonor = gamedata.CampBat_KillHonorMax
		}
	}

	pPlayer.CamBattleModule.DB_SaveKillData()
	response.KillHonor = pPlayer.CamBattleModule.KillHonor
	response.KillNum = pPlayer.CamBattleModule.Kill

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_KILL_EVENT_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())

	return
}

func Hand_PlayerQueryReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_QUERY_REQ")

	var req msg.MSG_PlayerQuery_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerQueryReq : Message Reader Error!!!!")
		return
	}

	pPlayer := GetPlayerByID(req.PlayerID)
	if pPlayer == nil {
		gamelog.Error("Hand_PlayerQueryReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	var response msg.MSG_PlayerQuery_Ack

	//如果己经开始搬运
	if pPlayer.CamBattleModule.EndTime > 0 {
		//如果己经超时，则搬运置为停止
		if time.Now().Unix() > pPlayer.CamBattleModule.EndTime {
			pPlayer.CamBattleModule.EndTime = 0
			pPlayer.CamBattleModule.CrystalID = 1
			pPlayer.CamBattleModule.DB_SaveMoveStautus()
		} else {
			gamelog.Error("Hand_PlayerQueryReq : Has already set the crystal quality!!!!")
		}
	} else { // 如果没有开始搬运
		pPlayer.CamBattleModule.CrystalID = 1
	}

	response.Quality = pPlayer.CamBattleModule.CrystalID
	response.PlayerID = req.PlayerID

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_PLAYER_QUERY_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())
	return
}

func Hand_PlayerReviveReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_REVIVE_REQ")

	var req msg.MSG_PlayerRevive_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerReviveReq : Message Reader Error!!!!")
		return
	}

	pPlayer := GetPlayerByID(req.PlayerID)
	if pPlayer == nil {
		gamelog.Error("Hand_PlayerReviveReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	pReviveInfo := gamedata.GetReviveInfo(req.ReviveOpt)
	if pReviveInfo == nil {
		gamelog.Error("Hand_PlayerReviveReq : Invalid ReviveOpt :%d!!!!", req.ReviveOpt)
		return
	}

	//查是否足够
	if pReviveInfo.CostMoneyID > 0 {
		if pPlayer.RoleMoudle.CostMoney(pReviveInfo.CostMoneyID, pReviveInfo.CostMoneyNum) == false {
			gamelog.Error("Hand_PlayerReviveReq : Not Enough Money id:%d, num:%d!!!!", pReviveInfo.CostMoneyID, pReviveInfo.CostMoneyNum)
			return
		}
	}

	var response msg.MSG_ServerRevive_Ack
	response.RetCode = msg.RE_SUCCESS
	response.PlayerID = req.PlayerID
	response.Stay = pReviveInfo.Stay
	response.ProInc = pReviveInfo.IncRatio
	response.BuffTime = pReviveInfo.BuffTime

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_PLAYER_REVIVE_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())
	return
}

func Hand_PlayerChangeReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_CHANGE_REQ")

	var req msg.MSG_PlayerChange_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerChangeReq : Message Reader Error!!!!")
		return
	}

	pPlayer := GetPlayerByID(req.PlayerID)
	if pPlayer == nil {
		gamelog.Error("Hand_PlayerChangeReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	//如果己经开始搬运
	if pPlayer.CamBattleModule.EndTime > 0 {
		gamelog.Error("Hand_PlayerChangeReq : now the moving is not finished!!!")
		return
	}

	//检查玩家是不是有足够的钱
	if req.HighQuality == 1 {
		pCrystalInfo := gamedata.GetCrystalInfo(4)
		if pCrystalInfo == nil {
			gamelog.Error("Hand_PlayerChangeReq : Invalid Crystal ID :%d!!!", 4)
			return
		}

		if false == pPlayer.RoleMoudle.CheckMoneyEnough(pCrystalInfo.CostMoneyID, pCrystalInfo.CostMoneyNum) {
			gamelog.Error("Hand_PlayerChangeReq : Not Enough Money", req.PlayerID)
			return
		}

	} else {
		if false == pPlayer.RoleMoudle.CheckMoneyEnough(gamedata.CampBat_Chg_MoneyID, gamedata.CampBat_Chg_MoneyNum) {
			gamelog.Error("Hand_PlayerChangeReq : Not Enough Money", req.PlayerID)
			return
		}
	}

	var response msg.MSG_PlayerChange_Ack
	if req.HighQuality == 1 {
		pPlayer.CamBattleModule.CrystalID = 4
		response.NewQuality = 4
	} else {
		pPlayer.CamBattleModule.CrystalID = 1
		response.NewQuality = 1
	}

	pPlayer.CamBattleModule.DB_SaveMoveStautus()

	response.PlayerID = req.PlayerID
	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_PLAYER_CHANGE_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())
	return
}

func Hand_PlayerCarryReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_CARRY_REQ")
	var req msg.MSG_PlayerCarry_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerCarryReq : Message Reader Error!!!!")
		return
	}

	pPlayer := GetPlayerByID(req.PlayerID)
	if pPlayer == nil {
		gamelog.Error("Hand_PlayerCarryReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	//开始搬运
	if req.CarryEvt == 1 {
		if pPlayer.CamBattleModule.LeftTimes <= 0 {
			gamelog.Error("Hand_PlayerCarryReq : Not Enough Carry Time!!!!")
			return
		}

		if time.Now().Unix() < pPlayer.CamBattleModule.EndTime {
			gamelog.Error("Hand_PlayerCarryReq : Still On Moving!!!!")
			return
		}

		pPlayer.CamBattleModule.LeftTimes = pPlayer.CamBattleModule.LeftTimes - 1
		pPlayer.CamBattleModule.EndTime = time.Now().Unix()

	} else { //完成搬运
		if pPlayer.CamBattleModule.EndTime <= 0 {
			gamelog.Error("Hand_PlayerCarryReq : Has not start move!!!!")
			return
		}

		if time.Now().Unix() > pPlayer.CamBattleModule.EndTime {
			gamelog.Error("Hand_PlayerCarryReq : Has already out of time!!!!")
			return
		}

		//完成了搬运
		pCrystal := gamedata.GetCrystalInfo(pPlayer.CamBattleModule.CrystalID)
		if pCrystal == nil {
			gamelog.Error("Hand_PlayerCarryReq : Invalid Crystal ID:%d!!!!", pPlayer.CamBattleModule.CrystalID)
			return
		}

		pPlayer.RoleMoudle.AddMoney(pCrystal.MoneyID[0], pCrystal.MoneyNum[0])
		pPlayer.RoleMoudle.AddMoney(pCrystal.MoneyID[1], pCrystal.MoneyNum[1])
		pPlayer.CamBattleModule.EndTime = 0
		pPlayer.CamBattleModule.CrystalID = 1
	}

	var response msg.MSG_PlayerCarry_Ack
	response.CarryEvt = req.CarryEvt
	response.EndTime = int(pPlayer.CamBattleModule.EndTime)
	response.RetCode = msg.RE_SUCCESS
	response.LeftTimes = pPlayer.CamBattleModule.LeftTimes

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_PLAYER_CARRY_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())

	return
}

//! 阵营战服务器加载阵营战数据
func Hand_LoadCampBatInfo(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_LOAD_CAMPBAT_REQ")
	var req msg.MSG_LoadCampBattle_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_LoadCampBatInfo : Message Reader Error!!!!")
		return
	}

	//! 常规检测
	var player *TPlayer = GetPlayerByID(req.PlayerID)
	if player == nil {
		gamelog.Error("Hand_LoadCampBatInfo Error: Invalid PlayerID:%d", req.PlayerID)
		return
	}

	//if req.EnterCode != player.CamBattleModule.enterCode {
	//	gamelog.Error("Hand_LoadCampBatInfo Error: Invalide enterCode, req.EnterCode:%d, player.EnterCode:%d", req.EnterCode, player.CamBattleModule.enterCode)
	//	return
	//}

	var response msg.MSG_LoadCampBattle_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	player.CamBattleModule.CheckReset()

	response.BattleCamp = player.CamBattleModule.BattleCamp
	response.Level = player.GetLevel()
	response.LeftTimes = player.CamBattleModule.LeftTimes
	response.CurRank = G_CampBat_TodayKill.SetRankItem(req.PlayerID, player.CamBattleModule.Kill)
	response.KillNum = player.CamBattleModule.Kill
	response.KillHonor = player.CamBattleModule.KillHonor
	response.PlayerID = player.GetPlayerID()
	//if response.Level <= gamedata.CampBat_RoomMatchLvl {
	response.RoomType = 1
	//} else {
	//	response.RoomType = 2
	//}

	var HeroResults = make([]THeroResult, BATTLE_NUM)
	player.HeroMoudle.CalcFightValue(HeroResults)
	for i := 0; i < BATTLE_NUM; i++ {
		response.Heros[i].HeroID = HeroResults[i].HeroID
		response.Heros[i].Camp = HeroResults[i].Camp
		response.Heros[i].PropertyValue = HeroResults[i].PropertyValues
		response.Heros[i].PropertyPercent = HeroResults[i].PropertyPercents
		response.Heros[i].CampDef = HeroResults[i].CampDef
		response.Heros[i].CampKill = HeroResults[i].CampKill
		if response.Heros[i].HeroID != 0 {
			pHeroInfo := gamedata.GetHeroInfo(response.Heros[i].HeroID)
			if pHeroInfo != nil {
				response.Heros[i].SkillID = pHeroInfo.Skills[0]
				if pHeroInfo.AttackType == 1 || pHeroInfo.AttackType == 3 {
					response.Heros[i].AttackID = gamedata.AttackPhysicID
				} else {
					response.Heros[i].AttackID = gamedata.AttackMagicID
				}

			} else {
				gamelog.Error("Hand_LoadCampBatInfo Error: Invalid HeroID:%d", response.Heros[i].HeroID)
			}
		}
	}
	response.LeftTimes = player.CamBattleModule.LeftTimes
	if time.Now().Unix() > player.CamBattleModule.EndTime {
		player.CamBattleModule.EndTime = 0
	}

	response.MoveEndTime = int(player.CamBattleModule.EndTime)
	response.RetCode = msg.RE_SUCCESS

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_LOAD_CAMPBAT_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())
}

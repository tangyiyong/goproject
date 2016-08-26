package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"gamesvr/tcpclient"
	"msg"
	"time"
	"utility"
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

	player := GetPlayerByID(req.Killer)
	if player == nil {
		gamelog.Error("Hand_KillEventReq : Invalid PlayerID :%d!!!!", req.Killer)
		return
	}

	var response msg.MSG_KillEvent_Ack

	player.CamBattleModule.Kill += int(req.Kill)
	response.CurRank = int32(G_CampBat_TodayKill.SetRankItem(req.Killer, player.CamBattleModule.Kill))
	G_CampBat_KillSum.SetRankItem(req.Killer, player.CamBattleModule.Kill)
	G_CampBat_CampKill[player.CamBattleModule.BattleCamp-1].SetRankItem(req.Killer, player.CamBattleModule.Kill)
	if req.Destroy > 0 {
		player.CamBattleModule.Destroy += int(req.Destroy)
		player.CamBattleModule.DestroySum += int(req.Destroy)
		G_CampBat_TodayDestroy.SetRankItem(req.Killer, player.CamBattleModule.Destroy)
		G_CampBat_DestroySum.SetRankItem(req.Killer, player.CamBattleModule.DestroySum)
		G_CampBat_CampDestroy[player.CamBattleModule.BattleCamp-1].SetRankItem(req.Killer, player.CamBattleModule.Kill)
	}

	if req.SeriesKill == int32(gamedata.CampBat_NtyKillNum) {
		var nty msg.MSG_HorseLame_Notify
		nty.TextType = TextCampBatHorseLamp
		nty.Camps = append(nty.Camps, player.CamBattleModule.BattleCamp)
		nty.Params = append(nty.Params, player.RoleMoudle.Name)
		b, _ := json.Marshal(nty)
		SendMessageToClient(0, msg.MSG_HORSELAME_NOTIFY, b)
	}

	if player.CamBattleModule.KillHonor < gamedata.CampBat_KillHonorMax {
		AddHonor := int(req.Kill) * gamedata.Campbat_KillHonorOne
		if (AddHonor + player.CamBattleModule.KillHonor) <= gamedata.CampBat_KillHonorMax {
			player.CamBattleModule.KillHonor += AddHonor
			player.RoleMoudle.AddMoney(4, AddHonor)
		} else {
			player.RoleMoudle.AddMoney(4, gamedata.CampBat_KillHonorMax-player.CamBattleModule.KillHonor)
			player.CamBattleModule.KillHonor = gamedata.CampBat_KillHonorMax
		}
	}

	player.CamBattleModule.DB_SaveKillData()
	response.KillHonor = int32(player.CamBattleModule.KillHonor)
	response.KillNum = int32(player.CamBattleModule.Kill)

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_KILL_EVENT_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_CAMP_BATTLE_KILL, player.CamBattleModule.KillSum)
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_CAMP_BATTLE_GROUP_KILL, player.CamBattleModule.DestroySum)

	return
}

func Hand_PlayerQueryReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_QUERY_REQ")

	var req msg.MSG_PlayerQuery_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerQueryReq : Message Reader Error!!!!")
		return
	}

	player := GetPlayerByID(req.PlayerID)
	if player == nil {
		gamelog.Error("Hand_PlayerQueryReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	var response msg.MSG_PlayerQuery_Ack

	//如果己经开始搬运
	if player.CamBattleModule.EndTime > 0 {
		//如果己经超时，则搬运置为停止
		if int(time.Now().Unix()) > player.CamBattleModule.EndTime {
			player.CamBattleModule.EndTime = 0
			player.CamBattleModule.CrystalID = utility.Rand()%2 + 1
		} else {
			gamelog.Error("Hand_PlayerQueryReq : Has already set the crystal quality!!!!")
			return
		}
	} else { // 如果没有开始搬运
		player.CamBattleModule.CrystalID = utility.Rand()%2 + 1
	}

	response.Quality = int32(player.CamBattleModule.CrystalID)
	response.PlayerID = req.PlayerID

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_PLAYER_QUERY_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())

	player.CamBattleModule.DB_SaveMoveStautus()
	return
}

func Hand_PlayerReviveReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_REVIVE_REQ")

	var req msg.MSG_PlayerRevive_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerReviveReq : Message Reader Error!!!!")
		return
	}

	player := GetPlayerByID(req.PlayerID)
	if player == nil {
		gamelog.Error("Hand_PlayerReviveReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	pReviveInfo := gamedata.GetReviveInfo(int(req.ReviveOpt))
	if pReviveInfo == nil {
		gamelog.Error("Hand_PlayerReviveReq : Invalid ReviveOpt :%d!!!!", req.ReviveOpt)
		return
	}

	//查是否足够
	if pReviveInfo.CostMoneyID > 0 {
		if player.RoleMoudle.CostMoney(pReviveInfo.CostMoneyID, pReviveInfo.CostMoneyNum) == false {
			gamelog.Error("Hand_PlayerReviveReq : Not Enough Money id:%d, num:%d!!!!", pReviveInfo.CostMoneyID, pReviveInfo.CostMoneyNum)
			return
		}
	}

	var response msg.MSG_ServerRevive_Ack
	response.RetCode = msg.RE_SUCCESS
	response.PlayerID = req.PlayerID
	response.Stay = int32(pReviveInfo.Stay)
	response.ProInc = int32(pReviveInfo.IncRatio)
	response.BuffTime = int32(pReviveInfo.BuffTime)

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

	player := GetPlayerByID(req.PlayerID)
	if player == nil {
		gamelog.Error("Hand_PlayerChangeReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	//如果己经开始搬运
	if player.CamBattleModule.EndTime > 0 {
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

		if false == player.RoleMoudle.CheckMoneyEnough(pCrystalInfo.CostMoneyID, pCrystalInfo.CostMoneyNum) {
			gamelog.Error("Hand_PlayerChangeReq : Not Enough Money:%d", req.PlayerID)
			return
		}

	} else {
		if false == player.RoleMoudle.CheckMoneyEnough(gamedata.CampBat_Chg_MoneyID, gamedata.CampBat_Chg_MoneyNum) {
			gamelog.Error("Hand_PlayerChangeReq : Not Enough Money:%d", req.PlayerID)
			return
		}
	}

	var response msg.MSG_PlayerChange_Ack
	if req.HighQuality == 1 {
		player.CamBattleModule.CrystalID = 4
		response.NewQuality = 4
	} else {
		player.CamBattleModule.CrystalID = utility.Rand()%2 + 1
		response.NewQuality = int32(player.CamBattleModule.CrystalID)
	}

	response.PlayerID = req.PlayerID
	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_PLAYER_CHANGE_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())

	player.CamBattleModule.DB_SaveMoveStautus()

	return
}

func Hand_StartCarryReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_START_CARRY_REQ")
	var req msg.MSG_StartCarry_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_StartCarryReq : Message Reader Error!!!!")
		return
	}

	player := GetPlayerByID(req.PlayerID)
	if player == nil {
		gamelog.Error("Hand_StartCarryReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	if player.CamBattleModule.LeftTimes <= 0 {
		gamelog.Error("Hand_StartCarryReq : Not Enough Carry Time!!!!")
		return
	}

	if int(time.Now().Unix()) < player.CamBattleModule.EndTime {
		gamelog.Error("Hand_StartCarryReq : Still On Moving!!!!")
		return
	}

	player.CamBattleModule.LeftTimes = player.CamBattleModule.LeftTimes - 1
	player.CamBattleModule.EndTime = int(time.Now().Unix()) + gamedata.Campbat_MaxMoveTime

	var response msg.MSG_StartCarry_Ack
	response.PlayerID = req.PlayerID
	response.EndTime = int32(player.CamBattleModule.EndTime)
	response.RetCode = msg.RE_SUCCESS
	response.LeftTimes = int32(player.CamBattleModule.LeftTimes)

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_START_CARRY_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())

	player.CamBattleModule.DB_SaveMoveStautus()

	return
}

func Hand_FinishCarryReq(pTcpConn *tcpclient.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_FINISH_CARRY_REQ")
	var req msg.MSG_FinishCarry_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_FinishCarryReq : Message Reader Error!!!!")
		return
	}

	player := GetPlayerByID(req.PlayerID)
	if player == nil {
		gamelog.Error("Hand_FinishCarryReq : Invalid PlayerID :%d!!!!", req.PlayerID)
		return
	}

	if player.CamBattleModule.EndTime <= 0 {
		gamelog.Error("Hand_FinishCarryReq : Has not start move!!!!")
		return
	}

	if int(time.Now().Unix()) > player.CamBattleModule.EndTime {
		gamelog.Error("Hand_FinishCarryReq : Has already out of time!!!!")
		return
	}

	//完成了搬运
	pCrystal := gamedata.GetCrystalInfo(player.CamBattleModule.CrystalID)
	if pCrystal == nil {
		gamelog.Error("Hand_FinishCarryReq : Invalid Crystal ID:%d!!!!", player.CamBattleModule.CrystalID)
		return
	}

	player.RoleMoudle.AddMoney(pCrystal.MoneyID[0], pCrystal.MoneyNum[0])
	player.RoleMoudle.AddMoney(pCrystal.MoneyID[1], pCrystal.MoneyNum[1])
	player.CamBattleModule.EndTime = 0
	player.CamBattleModule.CrystalID = 1

	var response msg.MSG_FinishCarry_Ack
	response.PlayerID = req.PlayerID
	response.RetCode = msg.RE_SUCCESS
	response.MoneyID[0] = int32(pCrystal.MoneyID[0])
	response.MoneyID[1] = int32(pCrystal.MoneyID[1])
	response.MoneyNum[0] = int32(pCrystal.MoneyNum[0])
	response.MoneyNum[1] = int32(pCrystal.MoneyNum[1])

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_FINISH_CARRY_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())

	player.CamBattleModule.DB_SaveMoveStautus()

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

	if req.EnterCode != player.CamBattleModule.enterCode {
		gamelog.Error("Hand_LoadCampBatInfo Error: Invalide enterCode, req.EnterCode:%d, player.EnterCode:%d", req.EnterCode, player.CamBattleModule.enterCode)
	}

	player.CamBattleModule.CheckReset()
	player.CamBattleModule.enterCode = 0
	var response msg.MSG_LoadCampBattle_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	response.BattleCamp = player.CamBattleModule.BattleCamp
	response.Level = int32(player.GetLevel())
	response.LeftTimes = int32(player.CamBattleModule.LeftTimes)
	response.CurRank = int32(G_CampBat_TodayKill.SetRankItem(player.playerid, player.CamBattleModule.Kill))
	response.KillNum = int32(player.CamBattleModule.Kill)
	response.KillHonor = int32(player.CamBattleModule.KillHonor)
	response.PlayerID = player.playerid
	if int(time.Now().Unix()) > player.CamBattleModule.EndTime {
		player.CamBattleModule.EndTime = 0
	}

	response.MoveEndTime = int32(player.CamBattleModule.EndTime)
	response.RetCode = msg.RE_SUCCESS

	if response.Level <= int32(gamedata.CampBat_RoomMatchLvl) {
		response.RoomType = 1
	} else {
		response.RoomType = 2
	}

	var HeroResults = make([]THeroResult, BATTLE_NUM)
	player.HeroMoudle.CalcFightValue(HeroResults)
	for i := 0; i < BATTLE_NUM; i++ {
		response.Heros[i].HeroID = int32(HeroResults[i].HeroID)
		response.Heros[i].Camp = HeroResults[i].Camp
		response.Heros[i].PropertyValue = HeroResults[i].PropertyValues
		response.Heros[i].PropertyPercent = HeroResults[i].PropertyPercents
		response.Heros[i].CampDef = HeroResults[i].CampDef
		response.Heros[i].CampKill = HeroResults[i].CampKill

		if response.Heros[i].HeroID != 0 {
			pHeroInfo := gamedata.GetHeroInfo(int(HeroResults[i].HeroID))
			if pHeroInfo != nil {
				response.Heros[i].SkillID = int32(pHeroInfo.Skills[0])
				if pHeroInfo.AttackType == 1 || pHeroInfo.AttackType == 3 {
					response.Heros[i].AttackID = int32(gamedata.AttackPhysicID)
				} else {
					response.Heros[i].AttackID = int32(gamedata.AttackMagicID)
				}

			} else {
				gamelog.Error("Hand_LoadCampBatInfo Error: Invalid HeroID:%d", response.Heros[i].HeroID)
			}
		}
	}

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_LOAD_CAMPBAT_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pTcpConn.WriteMsgData(writer.GetDataPtr())
}

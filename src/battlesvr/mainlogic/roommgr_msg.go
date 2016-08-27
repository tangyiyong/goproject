package mainlogic

import (
	"battlesvr/gamedata"
	"encoding/json"
	"gamelog"
	"msg"
	"tcpserver"
)

func (self *TRoomMgr) Hand_EnterRoom(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_ENTER_ROOM_REQ")
	var req msg.MSG_EnterRoom_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_EnterRoom : Message Reader Error!!!!")
		return
	}

	if req.PlayerID < 10000 {
		gamelog.Error("Hand_EnterRoom : req.PlayerID :%d", req.PlayerID)
		return
	}

	poldConn := GetConnByID(req.PlayerID)
	if poldConn != nil {
		DelConnByID(req.PlayerID)
		pBattleData := poldConn.Data.(*TBattleData)
		if pBattleData != nil && pBattleData.RoomID > 0 {
			pRoom := G_RoomMgr.GetRoomByID(pBattleData.RoomID)
			if pRoom == nil {
				gamelog.Error("Hand_EnterRoom : Invalid RoomID :%d", pBattleData.RoomID)
				return
			}

			var tmsg TMessage
			tmsg.MsgID = msg.MSG_LEAVE_BY_DISCONNT
			tmsg.MsgData = make([]byte, 4)
			tmsg.MsgData = append(tmsg.MsgData, byte(req.PlayerID), byte(req.PlayerID>>8), byte(req.PlayerID>>16), byte(req.PlayerID>>24))
			pRoom.MsgList <- tmsg
		} else {
			gamelog.Error("Hand_EnterRoom : old BattleData roomid <= 0  not in a room yet!!!")
		}

		poldConn.Cleaned = true
		poldConn.Close()
		gamelog.Info("CheckAndClean Error: Clean the unclosed Connection:%d", req.PlayerID)
	}

	AddTcpConn(req.PlayerID, 0, pTcpConn)
	var LoadReq msg.MSG_LoadCampBattle_Req
	LoadReq.PlayerID = req.PlayerID
	LoadReq.EnterCode = req.EnterCode
	SendMessageToGameSvr(msg.MSG_LOAD_CAMPBAT_REQ, 0, &LoadReq)
}

func CreateBattleObject(loadack *msg.MSG_LoadCampBattle_Ack) *TBattleObj {
	pBattleObj := new(TBattleObj)
	pBattleObj.PlayerID = loadack.PlayerID
	pBattleObj.Level = loadack.Level
	pBattleObj.BatCamp = loadack.BattleCamp
	pBattleObj.MoveEndTime = loadack.MoveEndTime
	for i := 0; i < 6; i++ {
		if loadack.Heros[i].HeroID == 0 {
			break
		}
		pBattleObj.HeroObj[i].HeroID = loadack.Heros[i].HeroID
		pBattleObj.HeroObj[i].PropertyValue = loadack.Heros[i].PropertyValue
		pBattleObj.HeroObj[i].PropertyPercent = loadack.Heros[i].PropertyPercent
		pBattleObj.HeroObj[i].CampDef = loadack.Heros[i].CampDef
		pBattleObj.HeroObj[i].CampKill = loadack.Heros[i].CampKill
		pBattleObj.HeroObj[i].Camp = loadack.Heros[i].Camp
		pBattleObj.HeroObj[i].Position = gamedata.GetCampHeroPos(loadack.BattleCamp)
		pBattleObj.HeroObj[i].SkillState.ID = loadack.Heros[i].SkillID
		pBattleObj.HeroObj[i].AttackPID = loadack.Heros[i].AttackID
		pBattleObj.HeroObj[i].CalcCurProperty(true)
	}

	pBattleObj.InitSkillState()

	return pBattleObj
}

func (self *TRoomMgr) Hand_LoadCampBatAck(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_LOAD_CAMPBAT_ACK")
	var req msg.MSG_LoadCampBattle_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_LoadCampBatAck : Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_LoadCampBatAck : Cannot find the connection!!!!", req.PlayerID)
		return
	}

	pBattleObj := CreateBattleObject(&req)
	RoomID := G_RoomMgr.AddPlayerToRoom(1 /*req.RoomType*/, pBattleObj.BatCamp, pBattleObj)
	if RoomID <= 0 {
		gamelog.Error("Hand_LoadCampBatAck Error AddPlayerToRoom Failed!!!")
		return
	}
	pConn.Data.(*TBattleData).RoomID = RoomID

	var resAck msg.MSG_EnterRoom_Ack
	resAck.MoveEndTime = req.MoveEndTime
	resAck.BatCamp = pBattleObj.BatCamp
	resAck.LeftTimes = req.LeftTimes
	resAck.CurRank = req.CurRank
	resAck.KillNum = req.KillNum
	resAck.KillHonor = req.KillHonor
	for j := 0; j < 4; j++ {
		resAck.SkillID[j] = pBattleObj.SkillState[j].ID
	}

	for i := 0; i < 6; i++ {
		resAck.Heros[i].HeroID = pBattleObj.HeroObj[i].HeroID
		resAck.Heros[i].ObjectID = pBattleObj.HeroObj[i].ObjectID
		resAck.Heros[i].Position = pBattleObj.HeroObj[i].Position
		resAck.Heros[i].CurHp = pBattleObj.HeroObj[i].CurProperty[0]
	}

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_ENTER_ROOM_ACK, 0)
	resAck.Write(&writer)
	writer.EndWrite()
	pConn.WriteMsgData(writer.GetDataPtr())

	//
	var resNotify msg.MSG_EnterRoom_Notify
	var obj msg.MSG_BattleObj
	obj.BatCamp = pBattleObj.BatCamp
	for i := 0; i < 6; i++ {
		obj.Heros[i].HeroID = pBattleObj.HeroObj[i].HeroID
		obj.Heros[i].ObjectID = pBattleObj.HeroObj[i].ObjectID
		obj.Heros[i].Position = pBattleObj.HeroObj[i].Position
		obj.Heros[i].CurHp = pBattleObj.HeroObj[i].CurProperty[0]
	}
	resNotify.BatObjs = append(resNotify.BatObjs, obj)
	resNotify.BatObjs_Cnt = 1
	SendMessageToRoom(req.PlayerID, RoomID, msg.MSG_ENTER_ROOM_NTY, &resNotify)

	OtherCount := 0
	var resNotify2 msg.MSG_EnterRoom_Notify
	pRoom := G_RoomMgr.GetRoomByID(RoomID)
	for i := 0; i < len(pRoom.Players); i++ {
		if pRoom.Players[i] != nil && pRoom.Players[i].PlayerID > 0 && pRoom.Players[i].PlayerID != pBattleObj.PlayerID {
			OtherCount += 1
			var obj msg.MSG_BattleObj
			obj.BatCamp = pRoom.Players[i].BatCamp
			for j := 0; j < 6; j++ {
				obj.Heros[j].HeroID = pRoom.Players[i].HeroObj[j].HeroID
				obj.Heros[j].ObjectID = pRoom.Players[i].HeroObj[j].ObjectID
				obj.Heros[j].Position = pRoom.Players[i].HeroObj[j].Position
				obj.Heros[j].CurHp = pRoom.Players[i].HeroObj[j].CurProperty[0]
			}
			resNotify2.BatObjs_Cnt += 1
			resNotify2.BatObjs = append(resNotify2.BatObjs, obj)
		}
	}

	if OtherCount > 0 {
		writer.BeginWrite(msg.MSG_ENTER_ROOM_NTY, 0)
		resNotify2.Write(&writer)
		writer.EndWrite()
		pConn.WriteMsgData(writer.GetDataPtr())
	}
	gamelog.Info("Hand_EnterRoom : player :%d  Enter Room Successed!!!", req.PlayerID)
	return
}

func (self *TRoomMgr) Hand_CheckInReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_CHECK_IN_REQ")
	var req msg.MSG_CheckIn_Req
	if json.Unmarshal(pdata, &req) != nil {
		gamelog.Error("Hand_CheckInReq : Unmarshal error!!!!:%s", pdata)
		return
	}

	if req.PlayerID == 0 || req.PlayerID >= 10000 {
		gamelog.Error("Hand_CheckInReq  Invalid PlayerID:%d", req.PlayerID)
		return
	}

	//收到的是服务器连接

	pData := new(TBattleData)
	pData.RoomID = 0
	pData.PlayerID = req.PlayerID
	pTcpConn.Data = pData
	pTcpConn.Cleaned = false
	G_GameSvrConns = pTcpConn
	gamelog.Info("message: Hand_CheckInReq id:%d, name:%s", req.PlayerID, req.PlayerName)
	return
}

func (self *TRoomMgr) Hand_Disconnect(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_DISCONNECT")
	if pTcpConn == nil || pTcpConn.Data == nil {
		gamelog.Error("Hand_Disconnect Error :pTcpConn == nil || pTcpConn.Data == nil")
		return
	}

	if pTcpConn.Cleaned == true {
		gamelog.Info("Hand_Disconnect :pTcpConn.Cleaned == true")
		return
	}

	pData := pTcpConn.Data.(*TBattleData)
	if pData == nil || pData.PlayerID <= 0 || pData.RoomID <= 0 {
		gamelog.Error("Hand_Disconnect Error :pData == nil || pData.PlayerID <= 0 || pData.RoomID <= 0")
		return
	}

	DelConnByID(pData.PlayerID)

	pRoom := G_RoomMgr.GetRoomByID(pData.RoomID)
	if pRoom == nil {
		gamelog.Error("Hand_Disconnect Error : Invalid RoomID:%d", pData.RoomID)
		return
	}

	var tmsg TMessage
	tmsg.MsgID = msg.MSG_LEAVE_BY_DISCONNT
	tmsg.MsgData = make([]byte, 4)
	tmsg.MsgData = append(tmsg.MsgData, byte(pData.PlayerID), byte(pData.PlayerID>>8), byte(pData.PlayerID>>16), byte(pData.PlayerID>>24))
	pRoom.MsgList <- tmsg
	return
}

func (self *TRoomMgr) Hand_HeartBeat(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_HEART_BEAT")
	if pTcpConn.Data == nil {
		gamelog.Info("Hand_PlayerChatReq Error: pTcpConn.Data == nil")
		return
	}
	var req msg.MSG_HeartBeat_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_HeartBeat : Message Reader Error!!!!")
		return
	}

	return
}

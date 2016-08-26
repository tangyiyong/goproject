package mainlogic

import (
	"battlesvr/gamedata"
	"encoding/json"
	"gamelog"
	"msg"
	"tcpserver"
	"time"
)

func Hand_CheckInReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_CHECK_IN_REQ")
	var req msg.MSG_CheckIn_Req
	if json.Unmarshal(pdata, &req) != nil {
		gamelog.Error("Hand_CheckInReq : Unmarshal error!!!!")
		return
	}

	if req.PlayerID == 0 || req.PlayerID >= 10000 {
		gamelog.Error("Hand_CheckInReq  Invalid PlayerID:%d", req.PlayerID)
		return
	}

	//收到的是服务器连接
	G_GameSvrConns = pTcpConn
	gamelog.Info("message: Hand_CheckInReq id:%d, name:%s", req.PlayerID, req.PlayerName)
	return
}

func Hand_DisConnect(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_DISCONNECT")
	if pTcpConn == nil || pTcpConn.Data == nil {
		gamelog.Error("pTcpConn == nil || pTcpConn.Data == nil")
		return
	}

	if pTcpConn.Cleaned == true {
		gamelog.Info("pTcpConn.Cleaned == true")
		return
	}

	pData := pTcpConn.Data.(*TBattleData)
	if pData == nil || pData.PlayerID <= 0 || pData.RoomID <= 0 {
		gamelog.Error("pData == nil || pData.PlayerID <= 0 || pData.RoomID <= 0")
		return
	}

	DelConnByID(pData.PlayerID)

	var response msg.MSG_LeaveRoom_Notify
	response.ObjectIDs = G_RoomMgr.GetPlayerHeroIDs(pData.RoomID, pData.PlayerID)
	SendMessageToRoom(pData.PlayerID, pData.RoomID, msg.MSG_LEAVE_ROOM_NTY, &response)
	G_RoomMgr.RemovePlayerFromRoom(pData.RoomID, pData.PlayerID)

	return
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

func Hand_EnterRoom(pTcpConn *tcpserver.TCPConn, pdata []byte) {
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

	/////////////////////////////////
	//如果消息编号不对，说明存在外挂嫌疑，需要产即断网
	//if req.MsgNo != pTcpConn.NeedNo {
	//	gamelog.Error("")
	//	pTcpConn.Close() //关闭网络连接
	//	return
	//}
	//////////////////////////////////
	poldConn := GetConnByID(req.PlayerID)
	if poldConn != nil {
		DelConnByID(req.PlayerID)
		pBattleData := poldConn.Data.(*TBattleData)
		if pBattleData != nil && pBattleData.RoomID > 0 {
			var leaventy msg.MSG_LeaveRoom_Notify
			leaventy.ObjectIDs = G_RoomMgr.GetPlayerHeroIDs(pBattleData.RoomID, req.PlayerID)
			SendMessageToRoom(req.PlayerID, pBattleData.RoomID, msg.MSG_LEAVE_ROOM_NTY, &leaventy)
			G_RoomMgr.RemovePlayerFromRoom(pBattleData.RoomID, req.PlayerID)
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
	SendMessageToGameSvr(msg.MSG_LOAD_CAMPBAT_REQ, &LoadReq)
}

func UpdateHeroStateToRoom(roomid int32) {
	if roomid <= 0 {
		gamelog.Error("UpdateHeroStateToRoom Error: Invalid RoomID %d", roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Error("UpdateHeroStateToRoom Error: Invalid RoomID2 %d", roomid)
		return
	}

	var ackHeroState msg.MSG_HeroState_Nty
	for i := 0; i < len(pRoom.Players); i++ {
		if pRoom.Players[i] != nil && pRoom.Players[i].PlayerID > 0 {
			for j := 0; j < 6; j++ {
				ackHeroState.Heros = append(ackHeroState.Heros, msg.MSG_HeroItem{ObjectID: pRoom.Players[i].HeroObj[j].ObjectID, CurHp: pRoom.Players[i].HeroObj[j].CurHp})
			}
		}
	}
	ackHeroState.Heros_Cnt = int32(len(ackHeroState.Heros))
	SendMessageToRoom(0, roomid, msg.MSG_HERO_STATE, &ackHeroState)

	return
}

func SendHeroStateToRoom(roomid int32, objectid int32) {
	if roomid <= 0 {
		gamelog.Error("SendHeroStateToRoom Error: Invalid RoomID %d", roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Error("SendHeroStateToRoom Error: Invalid RoomID2 %d", roomid)
		return
	}

	var ackHeroState msg.MSG_HeroState_Nty
	for i := 0; i < len(pRoom.Players); i++ {
		if pRoom.Players[i] != nil && pRoom.Players[i].PlayerID > 0 {
			for j := 0; j < 6; j++ {
				ackHeroState.Heros = append(ackHeroState.Heros, msg.MSG_HeroItem{ObjectID: pRoom.Players[i].HeroObj[j].ObjectID, CurHp: pRoom.Players[i].HeroObj[j].CurHp})
			}
		}
	}
	ackHeroState.Heros_Cnt = int32(len(ackHeroState.Heros))
	SendMessageToRoom(0, roomid, msg.MSG_HERO_STATE, &ackHeroState)
	return
}

func Hand_LeaveRoom(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_LEAVE_ROOM_REQ")
	var req msg.MSG_LeaveRoom_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_LeaveRoom : Message Reader Error!!!!")
		return
	}

	var roomID int32 = 0
	var playerid int32 = 0
	if pTcpConn != nil && pTcpConn.Data != nil {
		roomID = pTcpConn.Data.(*TBattleData).RoomID
		playerid = pTcpConn.Data.(*TBattleData).PlayerID
	}

	if req.PlayerID == 0 || req.PlayerID != playerid {
		gamelog.Error("Hand_LeaveRoom Error: req.PlayerID :%d, playerid:%d", req.PlayerID, playerid)
		return
	}

	G_RoomMgr.RemovePlayerFromRoom(roomID, req.PlayerID)
	DelConnByID(req.PlayerID)

	var response msg.MSG_LeaveRoom_Notify
	response.ObjectIDs = G_RoomMgr.GetPlayerHeroIDs(roomID, req.PlayerID)
	SendMessageToRoom(req.PlayerID, roomID, msg.MSG_LEAVE_ROOM_NTY, &response)
	pTcpConn.Cleaned = true
	pTcpConn.Close()
	return
}

func Hand_MoveState(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	//gamelog.Info("message: MSG_MOVE_STATE")
	if pTcpConn.Data == nil {
		gamelog.Info("Hand_MoveState Error: pTcpConn.Data == nil")
		return
	}

	playerid := pTcpConn.Data.(*TBattleData).PlayerID
	roomid := pTcpConn.Data.(*TBattleData).RoomID
	var req msg.MSG_Move_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_MoveState : Message Reader Error!!!!")
		return
	}

	if req.MoveEvents_Cnt <= 0 {
		gamelog.Error("Hand_MoveState Error: Invalid MoveEvents_Cnt:%d!!!!", req.MoveEvents_Cnt)
		return
	}

	SendMessageToRoom(playerid, roomid, msg.MSG_MOVE_STATE, &req)
	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Error("Hand_MoveState : Invalid RoomID:%d!!!!", roomid)
		return
	}

	for i := 0; i < len(req.MoveEvents); i++ {
		pHeroObject := pRoom.GetHeroObject(req.MoveEvents[i].S_ID)
		if pHeroObject == nil {
			gamelog.Error("Hand_MoveState Error: Invalid Hero ObjectID:%d", req.MoveEvents[i].S_ID)
			continue
		}

		pHeroObject.Position = req.MoveEvents[i].Position
	}

	return
}

func Hand_BuffState(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_BUFF_STATE")
}

func Hand_PlayerQueryReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_QUERY_REQ")
	if pTcpConn.Data == nil {
		gamelog.Info("Hand_PlayerQueryReq Error: pTcpConn.Data == nil")
		return
	}
	playerid := pTcpConn.Data.(*TBattleData).PlayerID
	roomid := pTcpConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_PlayerQueryReq Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_PlayerQueryReq Error: Invalid  roomid :%d", roomid)
		return
	}

	G_GameSvrConns.WriteMsg(msg.MSG_PLAYER_QUERY_REQ, pdata)
	return
}

func Hand_PlayerQueryAck(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_QUERY_ACK")

	var req msg.MSG_PlayerQuery_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerQueryAck : Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_PlayerQueryAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	pConn.WriteMsg(msg.MSG_PLAYER_QUERY_ACK, pdata)

	return
}

func Hand_LoadCampBatAck(pTcpConn *tcpserver.TCPConn, pdata []byte) {
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
	writer.BeginWrite(msg.MSG_ENTER_ROOM_ACK)
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
		writer.BeginWrite(msg.MSG_ENTER_ROOM_NTY)
		resNotify2.Write(&writer)
		writer.EndWrite()
		pConn.WriteMsgData(writer.GetDataPtr())
	}
	gamelog.Info("Hand_EnterRoom : player :%d  Enter Room Successed!!!", req.PlayerID)
	return

}

func Hand_StartCarryReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_START_CARRY_REQ")
	if pTcpConn.Data == nil {
		gamelog.Info("Hand_PlayerChatReq Error: pTcpConn.Data == nil")
		return
	}
	playerid := pTcpConn.Data.(*TBattleData).PlayerID
	roomid := pTcpConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_StartCarryReq Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_StartCarryReq Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pBattleObj := pRoom.GetBattleByPID(playerid)
	if pBattleObj == nil {
		gamelog.Error("Hand_StartCarryReq Error: Invalid playerid:%d", playerid)
		return
	}

	if pBattleObj.MoveEndTime > 0 {
		gamelog.Error("Hand_StartCarryReq Error: Has Already Carry a Ctystal:%d", playerid)
		return
	}

	pRect := &gamedata.GetSceneInfo().Camps[pBattleObj.BatCamp-1].MoveBegin
	if !pBattleObj.IsTeamIn(pRect) {
		gamelog.Error("Hand_StartCarryReq Error: Not In The Start Carry Rect:%d", playerid)
		return
	}

	G_GameSvrConns.WriteMsg(msg.MSG_START_CARRY_REQ, pdata)

	return
}

func Hand_FinishCarryReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_FINISH_CARRY_REQ")
	if pTcpConn.Data == nil {
		gamelog.Info("Hand_PlayerChatReq Error: pTcpConn.Data == nil")
		return
	}
	playerid := pTcpConn.Data.(*TBattleData).PlayerID
	roomid := pTcpConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_FinishCarryReq Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_FinishCarryReq Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pBattleObj := pRoom.GetBattleByPID(playerid)
	if pBattleObj == nil {
		gamelog.Error("Hand_FinishCarryReq Error: Invalid playerid:%d", playerid)
		return
	}

	var req msg.MSG_FinishCarry_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_FinishCarryReq Error: Message Reader Error!!!!")
		return
	}

	if pBattleObj.MoveEndTime <= 0 {
		gamelog.Error("Hand_FinishCarryReq Error: Has not start:%d", playerid)
		return
	}

	if int32(time.Now().Unix()) > pBattleObj.MoveEndTime {
		gamelog.Error("Hand_FinishCarryReq Error: Too late:%d", playerid)
		return
	}

	pRect := &gamedata.GetSceneInfo().Camps[pBattleObj.BatCamp-1].MoveEnd
	if !pBattleObj.IsTeamIn(pRect) {
		gamelog.Error("Hand_FinishCarryReq Error: Not In The Start Carry Rect:%d", playerid)
		return
	}

	G_GameSvrConns.WriteMsg(msg.MSG_FINISH_CARRY_REQ, pdata)

	return
}

func Hand_StartCarryAck(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_START_CARRY_ACK")
	var req msg.MSG_StartCarry_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_StartCarryAck Error: Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_StartCarryAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	playerid := pConn.Data.(*TBattleData).PlayerID
	roomid := pConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_StartCarryAck Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_StartCarryAck Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pBattleObj := pRoom.GetBattleByPID(playerid)
	if pBattleObj == nil {
		gamelog.Error("Hand_StartCarryAck Error: Invalid playerid:%d", playerid)
		return
	}

	if req.RetCode == msg.RE_SUCCESS {
		pBattleObj.MoveEndTime = req.EndTime
	}

	pConn.WriteMsg(msg.MSG_START_CARRY_ACK, pdata)

	return
}

func Hand_FinishCarryAck(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_FINISH_CARRY_ACK")
	var req msg.MSG_FinishCarry_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_FinishCarryAck Error: Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_FinishCarryAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	if req.RetCode == msg.RE_SUCCESS {
		playerid := pConn.Data.(*TBattleData).PlayerID
		roomid := pConn.Data.(*TBattleData).RoomID
		if playerid <= 0 || roomid <= 0 {
			gamelog.Info("Hand_FinishCarryAck Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
			return
		}

		pRoom := G_RoomMgr.GetRoomByID(roomid)
		if pRoom == nil {
			gamelog.Info("Hand_FinishCarryAck Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
			return
		}

		pBattleObj := pRoom.GetBattleByPID(playerid)
		if pBattleObj == nil {
			gamelog.Error("Hand_FinishCarryAck Error: Invalid playerid:%d", playerid)
			return
		}

		pBattleObj.MoveEndTime = 0
	}

	pConn.WriteMsg(msg.MSG_FINISH_CARRY_ACK, pdata)

	return
}

func Hand_PlayerChangeReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_CHANGE_REQ")
	if pTcpConn.Data == nil {
		gamelog.Info("Hand_PlayerChatReq Error: pTcpConn.Data == nil")
		return
	}
	playerid := pTcpConn.Data.(*TBattleData).PlayerID
	roomid := pTcpConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_PlayerChangeReq Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_PlayerChangeReq Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	G_GameSvrConns.WriteMsg(msg.MSG_PLAYER_CHANGE_REQ, pdata)

	return
}

func Hand_PlayerChangeAck(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_CHANGE_ACK")
	var req msg.MSG_PlayerChange_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerChangeAck : Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_PlayerChangeAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	pConn.WriteMsg(msg.MSG_PLAYER_CHANGE_ACK, pdata)

	return
}

func Hand_PlayerReviveReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_REVIVE_REQ")
	if pTcpConn.Data == nil {
		gamelog.Info("Hand_PlayerChatReq Error: pTcpConn.Data == nil")
		return
	}
	playerid := pTcpConn.Data.(*TBattleData).PlayerID
	roomid := pTcpConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_PlayerReviveReq Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_PlayerReviveReq Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	G_GameSvrConns.WriteMsg(msg.MSG_PLAYER_REVIVE_REQ, pdata)

	return
}

func Hand_PlayerReviveAck(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_REVIVE_ACK")
	var req msg.MSG_ServerRevive_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerReviveAck : Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_PlayerReviveAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	playerid := pConn.Data.(*TBattleData).PlayerID
	roomid := pConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_PlayerReviveReq Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_PlayerChangeReq Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pBattleObj := pRoom.GetBattleByPID(playerid)
	if pBattleObj == nil {
		gamelog.Info("Hand_PlayerChangeReq Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	if req.ReviveOpt == 2 {
		pBattleObj.ReviveTime[0] += 1
	} else if req.ReviveOpt == 3 {
		pBattleObj.ReviveTime[1] += 1
	} else if req.ReviveOpt == 4 {
	}

	if req.Stay != 1 { //表示要安全区复活

	}

	var response msg.MSG_PlayerRevive_Ack
	response.MoneyID = req.MoneyID
	response.MoneyNum = req.MoneyNum
	response.BatCamp = pBattleObj.BatCamp
	response.RetCode = req.RetCode

	for i := 0; i < 6; i++ {
		if pBattleObj.HeroObj[i].HeroID <= 0 {
			break
		}
		response.Heros = append(response.Heros, msg.MSG_HeroObj{pBattleObj.HeroObj[i].HeroID,
			pBattleObj.HeroObj[i].ObjectID,
			pBattleObj.HeroObj[i].CurHp,
			pBattleObj.HeroObj[i].Position})
	}
	response.Heros_Cnt = int32(len(response.Heros))

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_PLAYER_REVIVE_ACK)
	response.Write(&writer)
	writer.EndWrite()
	pConn.WriteMsgData(writer.GetDataPtr())

	return
}

func Hand_HeartBeat(pTcpConn *tcpserver.TCPConn, pdata []byte) {
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

func Hand_PlayerChatReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_CAMPBAT_CHAT_REQ")
	if pTcpConn.Data == nil {
		gamelog.Info("Hand_PlayerChatReq Error: pTcpConn.Data == nil")
		return
	}

	var req msg.MSG_CmapBatChat_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerChatReq : Message Reader Error!!!!")
		return
	}

	roomid := pTcpConn.Data.(*TBattleData).RoomID
	SendMessageToRoom(req.PlayerID, roomid, msg.MSG_CAMPBAT_CHAT_REQ, &req)
	return
}

package main

import (
	"battlesvr/mainlogic"
	"msg"
	"tcpserver"
)

func RegHttpMsgHandler() {

}

//注册TCP消息处理方法
func RegTcpMsgHandler() {
	tcpserver.HandleFunc(msg.MSG_CHECK_IN_REQ, mainlogic.Hand_CheckInReq)
	tcpserver.HandleFunc(msg.MSG_DISCONNECT, mainlogic.Hand_DisConnect)
	tcpserver.HandleFunc(msg.MSG_HEART_BEAT, mainlogic.Hand_HeartBeat)

	tcpserver.HandleFunc(msg.MSG_MOVE_STATE, mainlogic.Hand_MoveState)
	tcpserver.HandleFunc(msg.MSG_SKILL_STATE, mainlogic.Hand_SkillState)
	tcpserver.HandleFunc(msg.MSG_BUFF_STATE, mainlogic.Hand_BuffState)
	tcpserver.HandleFunc(msg.MSG_ENTER_ROOM_REQ, mainlogic.Hand_EnterRoom)
	tcpserver.HandleFunc(msg.MSG_LEAVE_ROOM_REQ, mainlogic.Hand_LeaveRoom)
	tcpserver.HandleFunc(msg.MSG_PLAYER_QUERY_REQ, mainlogic.Hand_PlayerQueryReq)
	tcpserver.HandleFunc(msg.MSG_PLAYER_QUERY_ACK, mainlogic.Hand_PlayerQueryAck)
	tcpserver.HandleFunc(msg.MSG_PLAYER_CHANGE_REQ, mainlogic.Hand_PlayerChangeReq)
	tcpserver.HandleFunc(msg.MSG_PLAYER_CHANGE_ACK, mainlogic.Hand_PlayerChangeAck)
	tcpserver.HandleFunc(msg.MSG_PLAYER_REVIVE_REQ, mainlogic.Hand_PlayerReviveReq)
	tcpserver.HandleFunc(msg.MSG_PLAYER_REVIVE_ACK, mainlogic.Hand_PlayerReviveAck)
	tcpserver.HandleFunc(msg.MSG_LOAD_CAMPBAT_ACK, mainlogic.Hand_LoadCampBatAck)
	tcpserver.HandleFunc(msg.MSG_CAMPBAT_CHAT_REQ, mainlogic.Hand_PlayerChatReq)
	tcpserver.HandleFunc(msg.MSG_START_CARRY_REQ, mainlogic.Hand_StartCarryReq)
	tcpserver.HandleFunc(msg.MSG_FINISH_CARRY_REQ, mainlogic.Hand_FinishCarryReq)
	tcpserver.HandleFunc(msg.MSG_START_CARRY_ACK, mainlogic.Hand_StartCarryAck)
	tcpserver.HandleFunc(msg.MSG_FINISH_CARRY_ACK, mainlogic.Hand_FinishCarryAck)

}

//注册控制台消息处理方法
func RegConsoleCmdHandler() {

	//utility.HandleFunc()
	//utility.HandleFunc()
	//utility.HandleFunc()
}

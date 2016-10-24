package mainlogic

import (
	"appconfig"
	"gamelog"
	"msg"
	"utility"
)

const (
	EVENT_LOGIN_GAME   = 1  //登录游戏
	EVENT_CHARGE_MONEY = 2  //充值事件
	EVENT_GET_DIAMOND  = 3  //获取钻石
	EVENT_LOSE_DIAMOND = 4  //消耗钻石
	EVENT_GET_GOLD     = 5  //获取金币
	EVENT_LOSE_GOLD    = 6  //消耗金币
	EVENT_GET_ACT      = 7  //体力获取
	EVENT_LOSE_ACT     = 8  //体力消耗
	EVENT_GET_HERO     = 9  //获取英雄
	EVENT_LOSE_HERO    = 10 //失去英雄
	EVENT_GET_EQUIP    = 11 //获取装备
	EVENT_LOSE_EQUIP   = 12 //失去装备
	EVENT_GET_ITEM     = 13 //获取道具
	EVENT_LOSE_ITEM    = 14 //失去道具
)

const (
	SOURCE_BUY_ITEM = 1
)

func EventCharge(player *TPlayer, realmoney int32, charid int32) {
	SendLogNotify(player.playerid, EVENT_CHARGE_MONEY, player.GetLevel(), player.GetVipLevel(), realmoney, charid)

}

func SendLogNotify(playerid int32, eventid int32, level int, viplvl int8, param1, param2 int32) bool {
	if G_LogClient.TcpConn == nil {
		gamelog.Error("SendLogNotify Error: G_LogClient.TcpConn is nullptr!!!")
		return false
	}

	var req msg.MSG_SvrLogData
	req.SvrID = int32(appconfig.GameSvrID)
	req.EventID = eventid
	req.PlayerID = playerid
	req.Level = int32(level)
	req.VipLvl = viplvl
	req.Time = utility.GetCurTime()
	req.Param[0] = param1
	req.Param[1] = param2

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_SVR_LOGDATA, int16(req.SvrID))
	req.Write(&writer)
	writer.EndWrite()

	return G_LogClient.TcpConn.WriteMsgData(writer.GetDataPtr())
}

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

const ()

//充值事件
func EventCharge(player *TPlayer, realmoney int32, chargeid int32) {
	SendLogNotify(player.playerid, EVENT_CHARGE_MONEY, 0, player.GetLevel(), player.GetVipLevel(), realmoney, chargeid)
}

//获取钻石事件
func EventGetDiamond(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_GET_DIAMOND, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//失去钻石事件
func EventLoseDiamond(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_LOSE_DIAMOND, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//获取金币事件
func EventGetGold(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_GET_GOLD, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//失去金币事件
func EventLoseGold(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_LOSE_GOLD, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//获取体力事件
func EventGetAct(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_GET_ACT, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//失去体力事件
func EventLoseAct(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_LOSE_ACT, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//获取英雄事件
func EventGetHero(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_GET_HERO, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//失去英雄事件
func EventLoseHero(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_LOSE_HERO, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//获取装备事件
func EventGetEquip(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_GET_EQUIP, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//失去装备事件
func EventLoseEquip(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_LOSE_EQUIP, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//获取道具事件
func EventGetItem(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_GET_ITEM, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

//失去道具事件
func EventLoseItem(player *TPlayer, srcid int32, num int32) {
	SendLogNotify(player.playerid, EVENT_LOSE_ITEM, srcid, player.GetLevel(), player.GetVipLevel(), num, 0)
}

func SendLogNotify(playerid int32, eventid int32, srcid int32, level int, viplvl int8, param1, param2 int32) bool {
	if G_LogClient.TcpConn == nil {
		gamelog.Error("SendLogNotify Error: G_LogClient.TcpConn is nullptr!!!")
		return false
	}

	var req msg.MSG_SvrLogData
	req.SvrID = int32(appconfig.GameSvrID)
	req.ChnlID = 0
	req.EventID = eventid
	req.SrcID = srcid
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

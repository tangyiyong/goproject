package mainlogic

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"msg"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//! 玩家请求充值结果
func Hand_GetChargeInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetChargeInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetChargeInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetChargeInfo_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.ChargeTimes = player.ChargeModule.ChargeTimes[1:] //0位是空的，不发给client了
	response.CardDays = player.ActivityModule.MonthCard.CardDays
	response.ActivityChargeID = player.ActivityModule.LimitSale.GetDiscountCharge()
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求充值结果
func Hand_GetChargeResult(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetChargeResult_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetChargeResult Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetChargeResult_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 返回领取状态
	response.RetCode = msg.RE_SUCCESS
	response.VipLevel = player.GetVipLevel()
	response.VipExp = player.GetVipExp()
	response.MoneyNum = player.RoleMoudle.GetMoney(gamedata.ChargeMoneyID)
}

//! 玩家请求领取激活码礼包
func Hand_RecvGiftCode(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_RecvGiftCode_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RecvGiftCode Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_RecvGiftCode_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	//查数据库，这个礼包是否己被领取。
	//如果己被领取,则不可领。
	//如果没有被领取, 则向礼包服务器请求领取。
	//如果领取成功， 则将礼包码存库，奖品领取。
	if len(req.GiftCode) <= 24 {
		gamelog.Error("Hand_RecvGiftCode Error: Invalid giftCode %s", req.GiftCode)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	giftCode := req.GiftCode[:24]
	s := mongodb.GetDBSession()
	defer s.Close()
	n, _ := s.DB(appconfig.GameDbName).C("PlayerGiftCode").Find(&bson.M{"pid": player.playerid, "gid": giftCode}).Count()
	if n > 0 {
		gamelog.Error("Hand_RecvGiftCode Error: Already Received GiftCode:%s", req.GiftCode)
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//领取
	var giftcodeReq msg.MSG_GameSvrGiftCode_Req
	giftcodeReq.AccountID = player.playerid
	giftcodeReq.SvrID = int32(appconfig.GameSvrID)
	giftcodeReq.ID = req.GiftCode
	b, _ := json.Marshal(giftcodeReq)
	http.DefaultClient.Timeout = 2 * time.Second
	httpRet, err := http.Post(appconfig.GiftCodeSvrUrl, "text/HTML", bytes.NewReader(b))
	if err != nil {
		gamelog.Error("Hand_RecvGiftCode Error: Query Gift Code Status err : %s !!!!", err.Error())
		return
	}

	httpRetBuf := make([]byte, httpRet.ContentLength)
	httpRet.Body.Read(httpRetBuf)
	httpRet.Body.Close()

	var giftcodeAck msg.MSG_GameSvrGiftCode_Ack
	err = json.Unmarshal(httpRetBuf, &giftcodeAck)
	if err != nil {
		gamelog.Error("Hand_RecvGiftCode Unmarshal giftcodeAck fail, Error: %s", err.Error())
		return
	}

	if giftcodeAck.RetCode != msg.RE_SUCCESS {
		response.RetCode = giftcodeAck.RetCode
		gamelog.Error("Hand_RecvGiftCode Error, giftcodeAck.RetCode:%d", giftcodeAck.RetCode)
		return
	}

	info, err := s.DB(appconfig.GameDbName).C("PlayerGiftCode").Upsert(bson.M{"pid": player.playerid, "gid": giftCode}, bson.M{"pid": player.playerid, "gid": giftCode})
	if info.Updated > 0 {
		gamelog.Error("Hand_RecvGiftCode  fail, Already Received :%d--%s ", player.playerid, req.GiftCode)
		//己经被领取过了，不让领
		return
	}

	for i := 0; i < len(giftcodeAck.ItemID); i++ {
		player.BagMoudle.AddAwardItem(giftcodeAck.ItemID[i], giftcodeAck.ItemNum[i])
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{giftcodeAck.ItemID[i], giftcodeAck.ItemNum[i]})
	}

	response.RetCode = msg.RE_SUCCESS
}

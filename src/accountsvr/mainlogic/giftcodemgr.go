package mainlogic

import (
	"appconfig"
	"encoding/json"
	"fmt"
	"gamelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"msg"
	"net/http"
	"utility"
)

//礼品码表结构
type TGiftCode struct {
	ID          string  `bson:"_id"` //礼包ID
	SvrID       []int32 //分区ID
	Platform    []int32 //平台ID
	GiftAwardID int32   //礼包码奖励ID
	EndTime     int32   //截止时间
	IsRecv      bool    //是否己领取
	IsAll       bool    //是否为全服可领
}

type TGiftAward struct {
	ID      int `bson:"_id"` //奖励ID
	Name    string
	ItemID  []int //物品ID
	ItemNum []int //物品数量
}

var G_GiftAwardID = 0
var G_GiftAwardLst = []TGiftAward{}

func InitGiftCodeMgr() {
	s := mongodb.GetDBSession()
	defer s.Close()

	G_GiftAwardLst = []TGiftAward{}
	err := s.DB(appconfig.AccountDbName).C("GiftAward").Find(nil).Sort("+_id").All(&G_GiftAwardLst)
	if err != nil {
		if err == mgo.ErrNotFound {
			G_GiftAwardID = 1
		} else {
			gamelog.Error("Init GiftCodeAward Failed Error : %s!!", err.Error())
			return
		}
	}

	if len(G_GiftAwardLst) <= 0 {
		G_GiftAwardID = 1
	} else {
		lastIndex := len(G_GiftAwardLst) - 1
		G_GiftAwardID = int(G_GiftAwardLst[lastIndex].ID) + 1
	}
}

func Handle_AddGiftAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_AddGiftAward_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_AddGiftAward : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_AddGiftAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var giftAward TGiftAward
	giftAward.ID = G_GiftAwardID
	giftAward.Name = req.Name
	giftAward.ItemID = req.ItemID
	giftAward.ItemNum = req.ItemNum
	mongodb.InsertToDB("GiftAward", giftAward)
	response.AwardID = giftAward.ID
	G_GiftAwardLst = append(G_GiftAwardLst, giftAward)
	G_GiftAwardID += 1
	response.RetCode = msg.RE_SUCCESS
}

func Handle_GetGiftAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetGiftAward_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GetGiftAward : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetGiftAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	for _, v := range G_GiftAwardLst {
		var award msg.MSG_GiftAward
		award.ID = v.ID
		award.Name = v.Name
		award.ItemID = v.ItemID
		award.ItemNum = v.ItemNum
		response.AwardLst = append(response.AwardLst, award)
	}

	response.RetCode = msg.RE_SUCCESS
}

func Handle_DelGiftAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_Del_giftaward_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_DelGiftAward : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_Del_giftaward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	for _, n := range req.AwardID {
		for i, v := range G_GiftAwardLst {
			if v.ID == n {
				G_GiftAwardLst = append(G_GiftAwardLst[:i], G_GiftAwardLst[i+1:]...)
				mongodb.RemoveFromDB(appconfig.AccountDbName, "GiftAward", &bson.M{"_id": n})

				response.RetCode = msg.RE_SUCCESS
				break
			}
		}
	}

}

func Handle_MakeGiftCode(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_MakeGiftCode_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_MakeGiftCode : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_MakeGiftCode_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var gift TGiftCode
	gift.SvrID = req.SvrID
	gift.Platform = req.Platform
	gift.EndTime = req.EndTime
	gift.IsAll = req.IsAll
	gift.GiftAwardID = req.GiftAwardID
	gift.IsRecv = false

	giftCode := utility.GetGuid()
	giftCode = giftCode[:24]
	for i := 0; i < req.GiftCodeNum; i++ {
		number := fmt.Sprintf("%x", i+1)
		gift.ID = giftCode + number

		response.GiftCodes = append(response.GiftCodes, gift.ID)
		mongodb.InsertToDB("GiftCode", gift)
	}
}

func Handle_GetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_QueryAccountInfo_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GetPlayerInfo : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_QueryAccountInfo_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	accountInfo, isExist := G_AccountMgr.GetAccountByID(req.AccountID)
	if isExist == false {
		gamelog.Error("Handle_GetPlayerInfo Error: AccountID %d not exist", req.AccountID)
		response.RetCode = msg.RE_ACCOUNT_NOT_EXIST
		return
	}

	response.AccountName = accountInfo.Name
	response.AccountPwd = accountInfo.Pwd
	response.CreateTime = accountInfo.CreateTime
	response.Enable = accountInfo.Enable
	response.LastLoginTime = accountInfo.LastTime
	response.Platform = accountInfo.Channel
	response.RetCode = msg.RE_SUCCESS
}

func Handle_GameSvrGiftCode(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GameSvrGiftCode_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GameSvrGiftCode : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GameSvrGiftCode_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//首先确认这个码是否可以领取
	//1.码是否有效(是否过期， 是否可以重复领， 是否己被领)
	//2.玩家所有的区是否符合
	//3.玩家所在的平台是否符合

	//如果玩家可以领
	//发给玩家礼包
	//不可以重复领的礼包标记己领取

	s := mongodb.GetDBSession()
	defer s.Close()
	var gift TGiftCode
	err := s.DB(appconfig.AccountDbName).C("GiftCode").Find(&bson.M{"_id": req.ID}).One(&gift)
	if err != nil {
		gamelog.Error("Handle_GameSvrGiftCode Error: %s", err.Error())
		response.RetCode = msg.RE_GIFTCODE_NOT_EXIST
		return
	}

	//! 检查领取
	if gift.IsRecv == true {
		gamelog.Error("Handle_GameSvrGiftCode Error: Aleady received gift code: %s", req.ID)
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 检查时间
	now := utility.GetCurTime()
	if gift.EndTime < now {
		gamelog.Error("Handle_GameSvrGiftCode Error: Gift code is outdated code: %s", req.ID)
		response.RetCode = msg.RE_GIFTCODE_OUTDATED
		return
	}

	//! 获取账户信息
	accountInfo, isExist := G_AccountMgr.GetAccountByID(req.AccountID)
	if isExist == false || accountInfo == nil {
		gamelog.Error("Handle_GameSvrGiftCode Error: Account not exist code: %s", req.ID)
		response.RetCode = msg.RE_ACCOUNT_NOT_EXIST
		return
	}

	//! 检查平台与服务器对应
	isExist = false
	for _, v := range gift.Platform {
		if v == 0 { //! 全服
			isExist = true
			break
		}
		if accountInfo.Channel == v {
			isExist = true
			break
		}
	}

	if isExist == false {
		gamelog.Error("Handle_GameSvrGiftCode Error: Non-matched Platform or SvrID code: %s", req.ID)
		response.RetCode = msg.RE_NON_MATCHED_PLATFORM_SVRID
	}

	isExist = false
	for _, v := range gift.SvrID {
		if v == 0 { //! 全服
			isExist = true
			break
		}
		if req.SvrID == v {
			isExist = true
			break
		}
	}

	if isExist == false {
		gamelog.Error("Handle_GameSvrGiftCode Error: Non-matched Platform or SvrID code: %s", req.ID)
		response.RetCode = msg.RE_NON_MATCHED_PLATFORM_SVRID
	}

	var giftAward TGiftAward
	err = s.DB(appconfig.AccountDbName).C("GiftAward").Find(&bson.M{"_id": gift.GiftAwardID}).One(&giftAward)
	if err != nil {
		gamelog.Error("Handle_GameSvrGiftCode Error: %s", err.Error())
		response.RetCode = msg.RE_GIFTCODE_NOT_EXIST
		return
	}

	//! 若非全服可领, 则设置标记
	if gift.IsAll == false {
		gift.IsRecv = true
		mongodb.UpdateToDB("GiftCode", &bson.M{"_id": req.ID}, &bson.M{"$set": bson.M{
			"isrecv": true}})
	}

	response.ItemID = giftAward.ItemID
	response.ItemNum = giftAward.ItemNum
	response.RetCode = msg.RE_SUCCESS
}

package mainlogic

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"gamelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"msg"
	"net/http"
	"time"
	"utility"
)

type DisableInfo struct {
	AccountID  int32 //! 账户ID
	BeginTime  int32 //! 起始时间
	DisableDay int32 //! 封禁天数 -1为永久封禁
}

type ServerDisableInfo struct {
	ServerID   int32         `bson:"_id"` //! 服务器ID
	DisableLst []DisableInfo //! 封禁信息
}

type TDisableAccountMgr struct {
	DisableLst map[int32]*ServerDisableInfo
}

var G_DisableMgr TDisableAccountMgr

//! 初始化
func InitDisableMgr() bool {
	var disableLst []ServerDisableInfo
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.AccountDbName).C("DisableAccount").Find(nil).All(&disableLst)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitDisableMgr DB Error!!!")
		return false
	}

	G_DisableMgr.DisableLst = make(map[int32]*ServerDisableInfo)
	for i := 0; i < len(disableLst); i++ {
		info := new(ServerDisableInfo)
		info = &disableLst[i]
		G_DisableMgr.DisableLst[info.ServerID] = info
	}

	return true
}

//! 获取服务器封禁信息
func (self *TDisableAccountMgr) GetServerDisableInfo(serverID int32) *ServerDisableInfo {
	disableInfo, isExist := self.DisableLst[serverID]
	if isExist == false {
		//! 不存在名单信息, 创建新名单
		newList := new(ServerDisableInfo)
		newList.ServerID = serverID
		self.DisableLst[serverID] = newList
		gamelog.Error("%v", self.DisableLst[serverID])
		mongodb.InsertToDB("DisableAccount", newList)

		return self.DisableLst[serverID]
	}

	return disableInfo
}

//! 获取封禁玩家信息
func (self *TDisableAccountMgr) GetDisablePlayerInifo(serverID int32, accountID int32) *DisableInfo {
	serverInfo := self.GetServerDisableInfo(serverID)

	for i := 0; i < len(serverInfo.DisableLst); i++ {
		if serverInfo.DisableLst[i].AccountID == accountID {
			return &serverInfo.DisableLst[i]
		}
	}
	return nil
}

func (self *TDisableAccountMgr) CheckExpire() {
	now := utility.GetCurTime()
	for j, v := range self.DisableLst {
		for i := 0; i < len(self.DisableLst[j].DisableLst); i++ {
			//! 非永久封禁且到期的玩家接触封禁
			if now >= v.DisableLst[i].DisableDay*(24*60*60)+v.DisableLst[i].BeginTime && v.DisableLst[i].DisableDay > 0 {
				//! 解除封禁
				pAccount, _ := G_AccountMgr.GetAccountByID(v.DisableLst[i].AccountID)
				if pAccount == nil {
					gamelog.Error("DisableAccountFunc : Invalid AccountID:%d!!!!", v.DisableLst[i].AccountID)
					continue
				}

				pAccount.Enable = 1
				mongodb.UpdateToDB("Account", &bson.M{"_id": v.DisableLst[i].AccountID},
					&bson.M{"$set": bson.M{"enable": pAccount.Enable}})

				newArray := v.DisableLst[:i]
				newArray = append(newArray, v.DisableLst[i+1:]...)
				self.DisableLst[j].DisableLst = newArray
				i -= 1
			}
		}
		self.DB_UpdateDisableLst(v.ServerID)
	}
}

func (self *TDisableAccountMgr) DisableAccount(svrID int32, accountID int32, disableDay int32) {
	accountInfo := self.GetDisablePlayerInifo(svrID, accountID)
	if accountInfo != nil {
		//! 已被封禁的时间叠加
		accountInfo.DisableDay += disableDay
		self.DB_AddDisableDay(svrID, accountID, accountInfo.DisableDay)
		return
	}
	serverLst := self.GetServerDisableInfo(svrID)

	//! 未找到信息则新加
	info := new(DisableInfo)
	info.AccountID = accountID
	info.BeginTime = utility.GetCurTime()
	info.DisableDay = disableDay
	self.DB_AddDisableInfo(svrID, info)
	serverLst.DisableLst = append(serverLst.DisableLst, *info)

	pAccount, _ := G_AccountMgr.GetAccountByID(accountID)
	if pAccount == nil {
		gamelog.Error("DisableAccountFunc : Invalid AccountID:%d!!!!", accountID)
		return
	}

	pAccount.Enable = 0
	mongodb.UpdateToDB("Account", &bson.M{"_id": accountID}, &bson.M{"$set": bson.M{"enable": pAccount.Enable}})
}

func (self *TDisableAccountMgr) EnableAccount(svrID int32, accountID int32) {
	accountInfo := self.GetDisablePlayerInifo(svrID, accountID)
	if accountInfo == nil {
		gamelog.Error("EnableAccount Error: not account disable info %d", accountID)
		return
	}

	self.DB_RemoveDisableInfo(svrID, accountInfo)

	disablelst := self.GetServerDisableInfo(svrID)
	for i := 0; i < len(disablelst.DisableLst); i++ {
		//! 非永久封禁且到期的玩家接触封禁
		if disablelst.DisableLst[i].AccountID == accountID {
			//! 解除封禁
			pAccount, _ := G_AccountMgr.GetAccountByID(accountID)
			if pAccount == nil {
				gamelog.Error("DisableAccountFunc : Invalid AccountID:%d!!!!", accountID)
				continue
			}

			pAccount.Enable = 1
			mongodb.UpdateToDB("Account", &bson.M{"_id": accountID}, &bson.M{"$set": bson.M{"enable": pAccount.Enable}})

			newArray := disablelst.DisableLst[:i]
			newArray = append(newArray, disablelst.DisableLst[i+1:]...)
			disablelst.DisableLst = newArray
			break
		}
	}

}

func (self *TDisableAccountMgr) DB_UpdateDisableLst(svrID int32) {
	mongodb.UpdateToDB("DisableAccount", &bson.M{"_id": svrID},
		&bson.M{"$set": bson.M{"disablelst": self.DisableLst[svrID].DisableLst}})
}

func (self *TDisableAccountMgr) DB_AddDisableInfo(svrID int32, info *DisableInfo) {
	mongodb.UpdateToDB("DisableAccount", &bson.M{"_id": svrID}, &bson.M{"$push": bson.M{"disablelst": *info}})
}

func (self *TDisableAccountMgr) DB_RemoveDisableInfo(svrID int32, info *DisableInfo) {
	mongodb.UpdateToDB("DisableAccount", &bson.M{"_id": svrID}, &bson.M{"$pull": bson.M{"disablelst": *info}})
}

func (self *TDisableAccountMgr) DB_AddDisableDay(svrID int32, accountID int32, disableDay int32) {
	mongodb.UpdateToDB("DisableAccount", &bson.M{"_id": svrID, "disablelst.accountid": accountID},
		&bson.M{"$set": bson.M{"disablelst.$.disableday": disableDay}})
}

func Handle_GmGetEnableLst(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	G_DisableMgr.CheckExpire()

	var req msg.MSG_GmGetEnableLst_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GmGetEnableLst : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GmGetEnableLst_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.PlayerID == 0 {
		//查询角色的账号ID
		req.PlayerID = QueryAccountID(req.PlayerName, req.SvrID)
	}

	now := utility.GetCurTime()

	for _, v := range G_DisableMgr.DisableLst {

		for _, n := range v.DisableLst {

			if (req.Type == 2 && req.PlayerID == n.AccountID) ||
				(req.Type == 1 && req.SvrID == v.ServerID) ||
				(req.Type == 0) {
				var _req msg.MSG_QueryPlayerInfo_Req
				_req.PlayerID = n.AccountID
				b, _ := json.Marshal(&_req)
				requrl := "http://" + GetGameSvrOutAddr(v.ServerID) + "/query_player_info"
				http.DefaultClient.Timeout = 2 * time.Second
				httpret, err := http.Post(requrl, "text/HTML", bytes.NewReader(b))
				if err != nil {
					gamelog.Error("query_player_info Error:  err : %s !!!!", err.Error())
					return
				}

				buffer = make([]byte, httpret.ContentLength)
				httpret.Body.Read(buffer)
				httpret.Body.Close()

				var ack msg.MSG_QueryPlayerInfo_Ack
				err = json.Unmarshal(buffer, &ack)
				if err != nil {
					gamelog.Error("query_player_info Error: Error: %s", err.Error())
					return
				} else if ack.RetCode != 0 {
					response.RetCode = ack.RetCode
					return
				}

				var info msg.MSG_GmEnablePlayerInfo
				info.SvrID = v.ServerID
				info.PlayerID = ack.PlayerID
				info.PlayerName = ack.PlayerName
				info.Level = ack.Level
				info.VipLevel = ack.VIPLevel

				info.Day = (now - n.BeginTime) / (24 * 60 * 60)
				info.DisableDay = n.DisableDay
				response.PlayerInfo = append(response.PlayerInfo, info)
			}
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

func Handle_GmEnableAccount(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GmEnableAccount_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GmEnableAccount : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GmEnableAccount_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.PlayerID == 0 {
		//查询角色的账号ID
		req.PlayerID = QueryAccountID(req.RoleName, req.SvrID)
	}

	if req.PlayerID <= 0 {
		gamelog.Error("Handle_GmEnableAccount : Cant find the player with name:%s, id:%d!!!!", req.RoleName, req.SvrID)
		return
	}

	if req.Enable == 0 {
		G_DisableMgr.DisableAccount(req.SvrID, req.PlayerID, req.DisableDay)
	} else {
		G_DisableMgr.EnableAccount(req.SvrID, req.PlayerID)
	}

	response.RetCode = msg.RE_SUCCESS
	return
}

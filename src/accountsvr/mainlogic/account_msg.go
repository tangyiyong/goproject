package mainlogic

import (
	"appconfig"
	"encoding/json"
	"gamelog"
	"math/rand"
	"mongodb"
	"msg"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

//处理登录请求
func Handle_Login(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	var buffer []byte
	buffer = make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_Login_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_Login : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_Login_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if !CheckAccountName(req.Name) {
		response.RetCode = msg.RE_INVALID_NAME
		return
	}

	if !CheckPassword(req.Password) {
		response.RetCode = msg.RE_INVALID_PASSWORD
		return
	}

	response.RetCode = msg.RE_SUCCESS
	result, bret := G_AccountMgr.GetAccountByName(req.Name)
	if !bret {
		response.RetCode = msg.RE_ACCOUNT_NOT_EXIST
	} else if result.Forbidden {
		response.RetCode = msg.RE_FORBIDDED_ACCOUNT
	} else if req.Password == result.Password {
		response.RetCode = msg.RE_SUCCESS
		response.AccountID = result.AccountID
		response.LoginKey = bson.NewObjectId().Hex()
		response.LastSvrID = result.LastSvrID
		var pGameInfo *TGameServerInfo = nil
		if result.LastSvrID <= 0 {
			pGameInfo = GetRecommendSvrID()
		} else {
			pGameInfo = GetGameSvrInfo(result.LastSvrID)
		}

		if pGameInfo != nil {
			response.LastSvrName = pGameInfo.SvrName
			response.LastSvrAddr = pGameInfo.svrOutAddr
		}
		G_AccountMgr.AddLoginKey(response.AccountID, response.LoginKey)
	} else {
		response.RetCode = msg.RE_INVALID_PASSWORD
	}

	if response.RetCode == msg.RE_SUCCESS {
		go mongodb.UpdateToDB(appconfig.AccountDbName, "Account", &bson.M{"_id": response.AccountID}, &bson.M{"$inc": bson.M{"logincount": 1}})
	}
}

//处理用户账户注册请求
func Handle_Register(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	var buffer []byte
	buffer = make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_RegAccount_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_Register : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_RegAccount_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	if !CheckAccountName(req.Name) {
		response.RetCode = msg.RE_INVALID_NAME
		return
	}

	if !CheckPassword(req.Password) {
		response.RetCode = msg.RE_INVALID_PASSWORD
		return
	}

	var pAccount *TAccount = nil
	pAccount, response.RetCode = G_AccountMgr.AddNewAccount(req.Name, req.Password)
	if response.RetCode == msg.RE_SUCCESS {
		mongodb.InsertToDB(appconfig.AccountDbName, "Account", pAccount)
	}
}

//游客玩家注册
func Handle_TouristRegister(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	var buffer []byte
	buffer = make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_TourRegAccount_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_TouristRegister : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_TourRegAccount_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var name string = "youke" + strconv.Itoa(rand.Intn(100000))
	var password string = strconv.Itoa(rand.Intn(100000) + 100000)

	var pAccount *TAccount = nil
	pAccount, response.RetCode = G_AccountMgr.AddNewAccount(name, password)
	if response.RetCode == msg.RE_SUCCESS {
		mongodb.InsertToDB(appconfig.AccountDbName, "Account", pAccount)
	}

	response.Name = name
	response.Password = password
}

//邦定游客账号
func Handle_BindTourist(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	var buffer []byte
	buffer = make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_BindTourist_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_BindTourist : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_BindTourist_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if !CheckAccountName(req.NewName) {
		response.RetCode = msg.RE_INVALID_NAME
		return
	}

	if !CheckPassword(req.NewPassword) {
		response.RetCode = msg.RE_INVALID_PASSWORD
		return
	}

	if G_AccountMgr.ResetAccount(req.Name, req.Password, req.NewName, req.NewPassword) == true {
		response.RetCode = msg.RE_SUCCESS
	} else {
		response.RetCode = msg.RE_FAILED
	}
	response.Name = req.NewName
	response.Password = req.NewPassword
}

func Handle_VerifyUserLogin(w http.ResponseWriter, r *http.Request) {
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	gamelog.Info("message: %s", r.URL.String())

	var req msg.MSG_VerifyUserLogin_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_VerifyUserLogin : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_VerifyUserLogin_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	if G_AccountMgr.CheckLoginKey(req.AccountID, req.LoginKey) {
		response.RetCode = msg.RE_SUCCESS
	}

	b, _ := json.Marshal(&response)
	w.Write(b)

	go ChangeLoginCountAndLast(req.AccountID, req.DomainID)
}

//处理登录请求
func Handle_ServerList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_ServerList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_ServerList : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_ServerList_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	nCount := len(G_ServerList)
	response.SvrList = make([]msg.ServerNode, 0, 10)
	for i := 0; i < nCount; i++ {
		if G_ServerList[i].SvrID != 0 && (G_ServerList[G_RecommendID].SvrFlag&SFG_VISIBLE > 0) {
			response.SvrList = append(response.SvrList, msg.ServerNode{G_ServerList[i].SvrID,
				G_ServerList[i].SvrName,
				G_ServerList[i].SvrFlag,
				G_ServerList[i].svrOutAddr})
		}
	}
}

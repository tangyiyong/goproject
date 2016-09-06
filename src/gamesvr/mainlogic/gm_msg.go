package mainlogic

import (
	"appconfig"
	"encoding/binary"
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"io/ioutil"
	"msg"
	"net/http"
	"os"
	"strings"
	"utility"
)

func Hand_SendAwardToPlayer(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_Send_Award_Player_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SendAwardToPlayer unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_Send_Award_Player_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_SendAwardToPlayer Error Invalid Gm request!!!")
		return
	}

	var data TAwardData
	data.TextType = Text_GM_Mail
	data.Value = append(data.Value, req.Value)
	for _, v := range req.ItemLst {
		data.ItemLst = append(data.ItemLst, gamedata.ST_ItemData{v.ID, v.Num})
	}

	SendAwardToPlayer(req.TargetID, &data)
	response.RetCode = msg.RE_SUCCESS
}
func Hand_AddSvrAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SvrAward_Add_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_AddSvrAward unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SvrAward_Add_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_AddSvrAward Error Invalid Gm request!!!")
		return
	}

	var data TAwardData
	data.TextType = Text_GM_Mail
	data.Value = req.Value
	for _, v := range req.ItemLst {
		data.ItemLst = append(data.ItemLst, gamedata.ST_ItemData{v.ID, v.Num})
	}

	G_GlobalVariables.AddSvrAward(&data)
	response.RetCode = msg.RE_SUCCESS
}
func Hand_DelSvrAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SvrAward_Del_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_DelSvrAward unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_SvrAward_Del_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_DelSvrAward Error Invalid Gm request!!!")
		return
	}

	G_GlobalVariables.DelSvrAward(req.ID)
	response.RetCode = msg.RE_SUCCESS
}

func Hand_UpdateGameData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	msglen := binary.LittleEndian.Uint32(buffer[:4])
	var req msg.MSG_UpdateGameData_Req
	if json.Unmarshal(buffer[4:4+msglen], &req) != nil {
		gamelog.Error("Hand_UpdateGameData : Unmarshal error!!!!")
		return
	}

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_UpdateGameData Error Invalid Gm request!!!")
		return
	}

	b, _ := utility.UnCompressData(buffer[4+msglen:])

	file := utility.GetCurrCsvPath() + req.TbName + ".csv"
	ioutil.WriteFile(file, b, 777)
	gamedata.ReloadOneFile(req.TbName)
	OnConfigChange(req.TbName)
	var response msg.MSG_UpdateGameData_Ack
	response.RetCode = msg.RE_SUCCESS
	ret, _ := json.Marshal(&response)
	w.Write(ret)
	return

}

func Hand_GetServerInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	var response msg.MSG_GetServerInfo_Ack
	response.SvrID = int32(appconfig.GameSvrID)
	response.SvrName = appconfig.GameSvrName
	response.OnlineCnt = G_OnlineCnt
	response.MaxOnlineCnt = G_MaxOnlineCnt
	response.RegisterCnt = G_RegisterCnt

	//	var ms runtime.MemStats
	//	runtime.ReadMemStats(&ms)
	//	response.MemAlloc = ms.HeapAlloc / 1024 / 1024
	//	response.MemInuse = ms.HeapSys / 1024 / 1024
	//	response.MenObjNum = ms.HeapObjects

	ret, _ := json.Marshal(&response)
	w.Write(ret)
	return

}

var clientlog *os.File = nil

func Hand_SaveClientInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var err error
	if clientlog == nil {
		clientlog, err = os.OpenFile(utility.GetCurrPath()+"log/client.log", os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			gamelog.Error("Hand_SaveClientInfo Error : %s", err.Error())
			return
		}
	}

	clientlog.Write(buffer)

	return

}

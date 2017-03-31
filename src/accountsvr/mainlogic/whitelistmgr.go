package mainlogic

import (
	"appconfig"
	"encoding/json"
	"gamelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"msg"
	"net/http"
	"strings"
)

var G_NetMgr TNetMgr

type TNetInfo struct {
	SvrID       int32    `bson:"_id"` //! 分区ID
	WhiteList   []string //! 白名单
	BlackList   []string //! 黑名单
	ChannelList []int    //! 渠道可见
}

//! 协议管理器
type TNetMgr struct {
	NetList map[int32]*TNetInfo
}

//! 初始化管理器
func InitNetMgr() bool {
	var infoLst []TNetInfo
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.AccountDbName).C("NetManagement").Find(nil).All(&infoLst)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitNetMgr DB Error!!!")
		return false
	}

	G_NetMgr.NetList = make(map[int32]*TNetInfo)
	for i := 0; i < len(infoLst); i++ {
		info := new(TNetInfo)
		info = &infoLst[i]
		G_NetMgr.NetList[info.SvrID] = info
	}
	return true
}

//! 获得服务器网络管理信息
func (self *TNetMgr) GetSvrNetInfo(svrID int32) *TNetInfo {
	netInfo, isExist := self.NetList[svrID]
	if isExist == false {
		//! 不存在名单信息, 创建新名单
		newList := new(TNetInfo)
		newList.SvrID = svrID
		self.NetList[svrID] = newList
		mongodb.InsertToDB("NetManagement", newList)

		return self.NetList[svrID]
	}

	return netInfo
}

//! 获得服务器白名单
func (self *TNetMgr) GetSvrNetWhiteList(svrID int32) []string {
	netInfo := self.GetSvrNetInfo(svrID)
	return netInfo.WhiteList
}

//! 获得服务器黑名单
func (self *TNetMgr) GetSvrNetBlackList(svrID int32) []string {
	netInfo := self.GetSvrNetInfo(svrID)
	return netInfo.BlackList
}

//! 获得服务器可见渠道
func (self *TNetMgr) GetSvrNetChannel(svrID int32) []int {
	netInfo := self.GetSvrNetInfo(svrID)
	return netInfo.ChannelList
}

//! 添加网络白名单
func (self *TNetMgr) AddSvrWhiteList(svrID int32, ip string) {
	netInfo := self.GetSvrNetInfo(svrID)
	netInfo.WhiteList = append(netInfo.WhiteList, ip)
	G_NetMgr.DB_AddSvrWhiteList(svrID, ip)
}

//! 添加可见渠道
func (self *TNetMgr) AddSvrChannelList(svrID int32, id int) {
	netInfo := self.GetSvrNetInfo(svrID)
	netInfo.ChannelList = append(netInfo.ChannelList, id)
	G_NetMgr.DB_AddSvrChannelList(svrID, id)
}

//! 删除可见渠道
func (self *TNetMgr) DelSvrChannelList(svrID int32, id int) bool {
	netInfo := self.GetSvrNetInfo(svrID)
	index := -1
	for i, v := range netInfo.ChannelList {
		if v == id {
			index = i
			break
		}
	}

	if index < 0 {
		return false
	}

	newArray := netInfo.ChannelList[:index]
	newArray = append(newArray, netInfo.ChannelList[index+1:]...)
	netInfo.ChannelList = newArray

	G_NetMgr.DB_DelSvrChannelList(svrID, id)
	return true
}

//! 删除网络白名单
func (self *TNetMgr) DelSvrWhiteList(svrID int32, ip string) bool {
	netInfo := self.GetSvrNetInfo(svrID)
	index := -1
	for i, v := range netInfo.WhiteList {
		if v == ip {
			index = i
			break
		}
	}

	if index < 0 {
		return false
	}

	newArray := netInfo.WhiteList[:index]
	newArray = append(newArray, netInfo.WhiteList[index+1:]...)
	netInfo.WhiteList = newArray

	G_NetMgr.DB_DelSvrWhiteList(svrID, ip)
	return true
}

//! 添加网络黑名单
func (self *TNetMgr) AddSvrBlackList(svrID int32, ip string) {
	netInfo := self.GetSvrNetInfo(svrID)
	netInfo.BlackList = append(netInfo.BlackList, ip)
	G_NetMgr.DB_AddSvrBlackList(svrID, ip)
}

//! 删除网络黑名单
func (self *TNetMgr) DelSvrBlackList(svrID int32, ip string) bool {
	netInfo := self.GetSvrNetInfo(svrID)

	index := -1
	for i, v := range netInfo.BlackList {
		if v == ip {
			index = i
			break
		}
	}

	if index < 0 {
		return false
	}

	newArray := netInfo.BlackList[:index]
	newArray = append(newArray, netInfo.BlackList[index+1:]...)
	netInfo.BlackList = newArray

	G_NetMgr.DB_DelSvrBlackList(svrID, ip)
	return true
}

//! 判断是否在白名单
func (self *TNetMgr) IsInWhiteList(svrID int32, ip string) bool {
	whiteList := self.GetSvrNetWhiteList(svrID)

	for _, v := range whiteList {
		if v == ip {
			return true
		}
	}

	return false
}

//! 判断是否在可见渠道中
func (self *TNetMgr) IsInChannelList(svrID int32, channelid int) bool {
	channelList := self.GetSvrNetChannel(svrID)

	for _, v := range channelList {
		if v == channelid {
			return true
		}
	}
	return false
}

//! 判断是否在黑名单
func (self *TNetMgr) IsInBlackList(svrID int32, ip string) bool {
	whiteList := self.GetSvrNetBlackList(svrID)
	for _, v := range whiteList {
		if v == ip {
			return true
		}
	}

	return false
}

func (self *TNetMgr) DB_AddSvrBlackList(svrID int32, ip string) {
	mongodb.UpdateToDB("NetManagement", &bson.M{"_id": svrID}, &bson.M{"$push": bson.M{"blacklist": ip}})
}

func (self *TNetMgr) DB_DelSvrBlackList(svrID int32, ip string) {
	mongodb.UpdateToDB("NetManagement", &bson.M{"_id": svrID}, &bson.M{"$pull": bson.M{"blacklist": ip}})
}

func (self *TNetMgr) DB_AddSvrWhiteList(svrID int32, ip string) {
	mongodb.UpdateToDB("NetManagement", &bson.M{"_id": svrID}, &bson.M{"$push": bson.M{"whitelist": ip}})
}

func (self *TNetMgr) DB_AddSvrChannelList(svrID int32, id int) {
	mongodb.UpdateToDB("NetManagement", &bson.M{"_id": svrID}, &bson.M{"$push": bson.M{"channellist": id}})
}

func (self *TNetMgr) DB_DelSvrChannelList(svrID int32, id int) {
	mongodb.UpdateToDB("NetManagement", &bson.M{"_id": svrID}, &bson.M{"$pull": bson.M{"channellist": id}})
}

func (self *TNetMgr) DB_DelSvrWhiteList(svrID int32, ip string) {
	mongodb.UpdateToDB("NetManagement", &bson.M{"_id": svrID}, &bson.M{"$pull": bson.M{"whitelist": ip}})
}

func Handle_GetNetList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetNetList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GetNetList : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetNetList_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Handle_GetNetList Error Invalid Gm request!!!")
		response.RetCode = msg.RE_INVALID_NAME
		return
	}

	netInfo := G_NetMgr.GetSvrNetInfo(req.SvrID)

	response.BlackList = netInfo.BlackList
	response.WhiteList = netInfo.WhiteList
	response.ChannelList = netInfo.ChannelList
	response.RetCode = msg.RE_SUCCESS
}

func Handle_AddNetList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_AddNetList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_AddNetList : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_AddNetList_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Handle_AddNetList Error Invalid Gm request!!!")
		response.RetCode = msg.RE_INVALID_NAME
		return
	}

	if req.ListType == 1 {
		G_NetMgr.AddSvrWhiteList(req.SvrID, req.IP)
	} else if req.ListType == 2 {
		G_NetMgr.AddSvrBlackList(req.SvrID, req.IP)
	} else if req.ListType == 3 {
		for _, v := range req.ChannelID {
			G_NetMgr.AddSvrChannelList(req.SvrID, v)
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

func Handle_DelNetList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_DelNetList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_DelNetList : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_DelNetList_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Handle_DelNetList Error Invalid Gm request!!!")
		response.RetCode = msg.RE_INVALID_NAME
		return
	}

	if req.ListType == 1 {
		G_NetMgr.DelSvrWhiteList(req.SvrID, req.IP)
	} else if req.ListType == 2 {
		G_NetMgr.DelSvrBlackList(req.SvrID, req.IP)
	} else if req.ListType == 3 {
		for _, v := range req.ChannelID {
			G_NetMgr.DelSvrChannelList(req.SvrID, v)
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

func Handle_QuerySvrIp(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_QuerySvrIp_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_QuerySvrIp : Unmarshal error!!!!")
		return
	}

	//var response msg.MSG_QuerySvrIp_Ack
	//response.RetCode = msg.RE_UNKNOWN_ERR
	//defer func() {
	//b, _ := json.Marshal(&response)

	//}()

	//response.SvrIp = GetGameSvrOutAddr(req.SvrID)
	//response.RetCode = msg.RE_SUCCESS
	var strIp string = GetGameSvrOutAddr(req.SvrID)
	w.Write([]byte(strIp))
}

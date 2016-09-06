package mainlogic

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

//! 增加申请名单
func DB_AddFriendAppList(hostid int32, appid int32) {
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": hostid}, &bson.M{"$addToSet": bson.M{"applylist": appid}})
}

//! 删除申请名单
func DB_RemoveFriendAppList(hostid int32, appid int32) {
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": hostid}, &bson.M{"$pull": bson.M{"applylist": appid}})
}

func (self *TFriendMoudle) DB_ClearAppList() {
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"applylist": self.ApplyList}})
}

//! 增加好友
func DB_AddFriend(hostid int32, friend *TFriendInfo) {
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": hostid}, &bson.M{"$push": bson.M{"friendlist": friend}})
}

//! 删除好友
func DB_RemoveFriend(hostid int32, appid int32) {
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": hostid}, &bson.M{"$pull": bson.M{"friendlist": appid}})
}

//! 更新好友里邻取状态
func DB_UpdateIsGive(hostid int32, nIndex int, IsGive bool) {
	FieldName := fmt.Sprintf("friendlist.%d.isgive", nIndex)
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": hostid}, &bson.M{"$set": bson.M{FieldName: IsGive}})
}

func DB_UpdateHasAct(hostid int32, nIndex int, HasAct bool) {
	FieldName := fmt.Sprintf("friendlist.%d.hasact", nIndex)
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": hostid}, &bson.M{"$set": bson.M{FieldName: HasAct}})
}

func DB_UpdateRcvNum(hostid int32, RcvNum int) {
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": hostid}, &bson.M{"$set": bson.M{"rcvnum": RcvNum}})
}

func (self *TFriendMoudle) DB_UpdateFriend() {
	GameSvrUpdateToDB("PlayerFriend", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"rcvnum": self.RcvNum,
		"friendlist": self.FriendList,
		"resetday":   self.ResetDay}})
}

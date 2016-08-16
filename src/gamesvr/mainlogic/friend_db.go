package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 增加申请名单
func DB_AddFriendAppList(hostid int, appid int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": hostid}, bson.M{"$addToSet": bson.M{"applylist": appid}})
}

//! 删除申请名单
func DB_RemoveFriendAppList(hostid int, appid int) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": hostid}, "applylist", appid)
}

func (self *TFriendMoudle) DB_ClearAppList() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"applylist": self.ApplyList}})
}

//! 增加好友
func DB_AddFriend(hostid int, friend *TFriendInfo) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": hostid}, bson.M{"$push": bson.M{"friendlist": friend}})
}

//! 删除好友
func DB_RemoveFriend(hostid int, appid int) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": hostid}, "friendlist", appid)
}

//! 更新好友里邻取状态
func DB_UpdateIsGive(hostid int, nIndex int, IsGive bool) {
	FieldName := fmt.Sprintf("friendlist.%d.isgive", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": hostid}, bson.M{"$set": bson.M{FieldName: IsGive}})
}

func DB_UpdateHasAct(hostid int, nIndex int, HasAct bool) {
	FieldName := fmt.Sprintf("friendlist.%d.hasact", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": hostid}, bson.M{"$set": bson.M{FieldName: HasAct}})
}

func DB_UpdateRcvNum(hostid int, RcvNum int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": hostid}, bson.M{"$set": bson.M{"rcvnum": RcvNum}})
}

func (self *TFriendMoudle) DB_UpdateFriend() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFriend", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"rcvnum": self.RcvNum,
		"friendlist": self.FriendList,
		"resetday":   self.ResetDay}})
}

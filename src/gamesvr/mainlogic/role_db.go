package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//保存体力值
func (self *TRoleMoudle) DB_SaveActions() {
	mongodb.UpdateToDB("PlayerRole", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"actions": self.Actions}})
}

func (self *TRoleMoudle) DB_SaveActionsAt(actionid int) {
	var FieldName = []byte("actions.$")
	FieldName[8] = byte(actionid - 1 + '0')
	mongodb.UpdateToDB("PlayerRole", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{string(FieldName): self.Actions[actionid-1]}})
}

//保存全部货币
func (self *TRoleMoudle) DB_SaveMoneys() {
	mongodb.UpdateToDB("PlayerRole", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"moneys": self.Moneys}})
}

//保存指定货币ID的货币
func (self *TRoleMoudle) DB_SaveMoneysAt(moneyid int) {
	FieldName := fmt.Sprintf("moneys.%d", moneyid-1)
	mongodb.UpdateToDB("PlayerRole", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.Moneys[moneyid-1]}})
}

//保存玩家的角色名
func (self *TRoleMoudle) DB_SaveRoleName() {
	mongodb.UpdateToDB("PlayerRole", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"name": self.Name}})
}

//保存玩家的角色名
func (self *TRoleMoudle) DB_SaveNewWizard() {
	mongodb.UpdateToDB("PlayerRole", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"newwizard": self.NewWizard}})
}

//保存全部货币
func (self *TRoleMoudle) DB_UpdateChargeMoney() {
	mongodb.UpdateToDB("PlayerRole", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"todaycharge": self.TodayCharge, "totalcharge": self.TotalCharge}})
}

//! 更新VIP等级
func (self *TRoleMoudle) DB_SaveVipLevel() {
	mongodb.UpdateToDB("PlayerRole", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"viplevel": self.VipLevel}})

	self.ownplayer.ActivityModule.VipGift.DB_SaveDailyResetTime()
}

//! 将升星信息存储到数据库
func (self *TRoleMoudle) DB_SaveSanGuoZhiStar() {
	mongodb.UpdateToDB("PlayerSanGuoZhi", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"curstarid": self.CurStarID}})
}

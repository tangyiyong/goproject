package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//保存体力值
func (role *TRoleMoudle) DB_SaveActions() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{"actions": role.Actions}})
}

func (role *TRoleMoudle) DB_SaveActionsAt(actionid int) {
	var FieldName = []byte("actions.$")
	FieldName[8] = byte(actionid - 1 + '0')
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{string(FieldName): role.Actions[actionid-1]}})
}

//保存全部货币
func (role *TRoleMoudle) DB_SaveMoneys() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{"moneys": role.Moneys}})
}

//保存指定货币ID的货币
func (role *TRoleMoudle) DB_SaveMoneysAt(moneyid int) {
	FieldName := fmt.Sprintf("moneys.%d", moneyid-1)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{FieldName: role.Moneys[moneyid-1]}})
}

//保存玩家的角色名
func (role *TRoleMoudle) DB_SaveRoleName() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{"name": role.Name}})
}

//保存玩家的角色名
func (role *TRoleMoudle) DB_SaveNewWizard() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{"newwizard": role.NewWizard}})
}

//保存经验加成等级信息
func (role *TRoleMoudle) DB_SaveExpIncLevel() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{"expinclvl": role.ExpIncLvl}})
}

//保存经验加成等级信息
func (role *TRoleMoudle) DB_AddColHero(heroid int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$push": bson.M{"colheros": heroid}})
}

//保存经验加成等级信息
func (role *TRoleMoudle) DB_AddColPet(petid int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$push": bson.M{"colpets": petid}})
}

//保存全部货币
func (role *TRoleMoudle) DB_UpdateChargeMoney() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{"todaycharge": role.TodayCharge, "totalcharge": role.TotalCharge}})
}

//! 更新VIP等级
func (role *TRoleMoudle) DB_SaveVipLevel() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRole", bson.M{"_id": role.PlayerID}, bson.M{"$set": bson.M{
		"viplevel": role.VipLevel}})

	role.ownplayer.ActivityModule.VipGift.DB_SaveDailyResetTime()
}

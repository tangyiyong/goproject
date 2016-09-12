package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//保存上阵装备信息
func (self *THeroMoudle) DB_SaveBattleEquipAt(nIndex int) {
	FieldName := fmt.Sprintf("curequips.%d", nIndex)
	mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.CurEquips[nIndex]}})
}

//保存上阵宝物信息
func (self *THeroMoudle) DB_SaveBattleGemAt(nIndex int) {
	FieldName := fmt.Sprintf("curgems.%d", nIndex)
	mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.CurGems[nIndex]}})
}

//保存主角的英雄ID
func (self *THeroMoudle) DB_SaveMainHeroID() {
	mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"curheros.0.id": self.CurHeros[0].ID}})
}

//保存称号信息
func (self *THeroMoudle) DB_SaveTitleInfo() {
	mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"titleid": self.TitleID}})
}

//保存称号信息
func (self *THeroMoudle) DB_SaveFashionInfo() {
	mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"fashionid": self.FashionID,
		"fashionlvl": self.FashionLvl}})
}

//保存上阵宠物信息
func (self *THeroMoudle) DB_SaveBattlePetAt(nIndex int) {
	FieldName := fmt.Sprintf("curpets.%d", nIndex)
	mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.CurPets[nIndex]}})
}

//保存额外属性信息
func (self *THeroMoudle) DB_SaveExtraProperty() {
	mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"extraprovalue": self.ExtraProValue,
		"extrapropercent": self.ExtraProPercent,
		"extracampkill":   self.ExtraCampKill,
		"extracampdef":    self.ExtraCampDef}})
}

//! 修改玩家公会技能等级
func (self *THeroMoudle) DB_SaveGuildSkillLevel() {
	mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"guildskilvl": self.GuildSkiLvl}})
}

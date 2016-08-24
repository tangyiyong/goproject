package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//保存上阵装备信息
func (self *THeroMoudle) DB_SaveBattleEquipAt(nIndex int) {
	FieldName := fmt.Sprintf("curequips.%d", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{FieldName: self.CurEquips[nIndex]}})
}

//保存上阵宝物信息
func (self *THeroMoudle) DB_SaveBattleGemAt(nIndex int) {
	FieldName := fmt.Sprintf("curgems.%d", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{FieldName: self.CurGems[nIndex]}})
}

//保存主角的英雄ID
func (self *THeroMoudle) DB_SaveMainHeroID() {
	FieldName := "curheros.0.id"
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{FieldName: self.CurHeros[0].ID}})
}

//保存公会技能等级信息
func (self *THeroMoudle) DB_SaveGuildSkill(nIndex int) {
	FieldName := fmt.Sprintf("guildskilvl.%d", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{FieldName: self.GuildSkiLvl[nIndex]}})
}

func (self *THeroMoudle) DB_SaveGuildSkillLst() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"guildskilvl": self.GuildSkiLvl}})
}

//保存称号信息
func (self *THeroMoudle) DB_SaveTitleInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"titleid": self.TitleID}})
}

//保存称号信息
func (self *THeroMoudle) DB_SaveFashionInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"fashionid": self.FashionID,
		"fashionlvl": self.FashionLvl}})
}

//保存上阵宠物信息
func (self *THeroMoudle) DB_SaveBattlePetAt(nIndex int) {
	FieldName := fmt.Sprintf("curpets.%d", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{FieldName: self.CurPets[nIndex]}})
}

//保存额外属性信息
func (self *THeroMoudle) DB_SaveExtraProperty() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"extraprovalue": self.ExtraProValue,
		"extrapropercent": self.ExtraProPercent,
		"extracampkill":   self.ExtraCampKill,
		"extracampdef":    self.ExtraCampDef}})
}

//保存额外属性信息
func (self *THeroMoudle) DB_SaveExtraPropertyAt(pid int, pvalue int, percent bool, camp int) {
	if camp > 0 {
		if pid == 6 {
			mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"extracampdef": self.ExtraCampDef}})
		} else if pid == 7 {
			mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"extracampkill": self.ExtraCampKill}})
		}
	} else {
		if percent == true {
			mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"extrapropercent": self.ExtraProPercent}})
		} else {
			mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHero", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"extraprovalue": self.ExtraProValue}})
		}
	}
}

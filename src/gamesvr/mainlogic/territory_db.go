package mainlogic

import (
	"fmt"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 更新领地信息
func (self *TTerritoryModule) DB_UpdateTerritory(index int, info *TTerritoryInfo) {
	filedName := fmt.Sprintf("territorylst.%d", index)
	mongodb.UpdateToDB("PlayerTerritory", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: *info}})
}

//! 更新领地技能
func (self *TTerritoryModule) DB_UpdateTerritorySkill(index int, level int) {
	filedName := fmt.Sprintf("territorylst.%d.skilllevel", index)
	mongodb.UpdateToDB("PlayerTerritory", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: level}})
}

//! 更新重置时间
func (self *TTerritoryModule) DB_UpdateResetTime() {
	mongodb.UpdateToDB("PlayerTerritory", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"resetday": self.ResetDay}})
}

//! 增加叛军信息
func (self *TTerritoryModule) DB_DB_AddTerritoryRiotInfo(index int, riot TTerritoryRiotData) {
	filedName := fmt.Sprintf("territorylst.%d.riotinfo", index)
	mongodb.UpdateToDB("PlayerTerritory", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{filedName: riot}})
}

//! 增加领地数量
func (self *TTerritoryModule) DB_AddTerritory(territory TTerritoryInfo) {
	mongodb.UpdateToDB("PlayerTerritory", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"territorylst": territory}})
}

//! 增加奖励
func (self *TTerritoryModule) DB_AddTerritoryAward(index int, award gamedata.ST_ItemData) {
	filedName := fmt.Sprintf("territorylst.%d.awarditem", index)
	mongodb.UpdateToDB("PlayerTerritory", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{filedName: award}})
}

//! 更新镇压暴动次数
func (self *TTerritoryModule) DB_UpdateRiotTimes() {
	mongodb.UpdateToDB("PlayerTerritory", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"suppressriottimes": self.SuppressRiotTimes}})
}

//! 更新镇压暴动信息
func (self *TTerritoryModule) DB_UpdateRiotInfo(index int, riotIndex int, info TTerritoryRiotData) {
	filedName := fmt.Sprintf("territorylst.%d.riotinfo.%d", index, riotIndex)
	mongodb.UpdateToDB("PlayerTerritory", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: info}})
}

package mainlogic

import (
	"appconfig"
	"fmt"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 更新领地信息
func (self *TTerritoryModule) UpdateTerritory(index int, info *TTerritoryInfo) {
	filedName := fmt.Sprintf("territorylst.%d", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTerritory", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: *info}})
}

//! 更新领地技能
func (self *TTerritoryModule) UpdateTerritorySkill(index int, level int) {
	filedName := fmt.Sprintf("territorylst.%d.skilllevel", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTerritory", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: level}})
}

//! 更新重置时间
func (self *TTerritoryModule) UpdateResetTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTerritory", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"resetday": self.ResetDay}})
}

//! 增加叛军信息
func (self *TTerritoryModule) AddTerritoryRiotInfo(index int, riot TTerritoryRiotData) {
	filedName := fmt.Sprintf("territorylst.%d.riotinfo", index)
	mongodb.AddToArray(appconfig.GameDbName, "PlayerTerritory", bson.M{"_id": self.PlayerID}, filedName, riot)
}

//! 增加领地数量
func (self *TTerritoryModule) AddTerritory(territory TTerritoryInfo) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerTerritory", bson.M{"_id": self.PlayerID}, "territorylst", territory)
}

//! 增加奖励
func (self *TTerritoryModule) AddTerritoryAward(index int, award gamedata.ST_ItemData) {
	filedName := fmt.Sprintf("territorylst.%d.awarditem", index)
	mongodb.AddToArray(appconfig.GameDbName, "PlayerTerritory", bson.M{"_id": self.PlayerID}, filedName, award)
}

//! 更新镇压暴动次数
func (self *TTerritoryModule) UpdateRiotTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTerritory", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"suppressriottimes": self.SuppressRiotTimes}})
}

//! 更新镇压暴动信息
func (self *TTerritoryModule) UpdateRiotInfo(index int, riotIndex int, info TTerritoryRiotData) {
	filedName := fmt.Sprintf("territorylst.%d.riotinfo.%d", index, riotIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTerritory", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: info}})
}

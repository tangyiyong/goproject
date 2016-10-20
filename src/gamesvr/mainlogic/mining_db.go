package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 存储玩家当前积分
func (self *TMiningModule) DB_SavePoint() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"point": self.Point}})
}

//! 存储玩家Buff信息
func (self *TMiningModule) DB_SaveBuff() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buff": self.Buff}})
}

//! 刷新翻牌奖励
func (self *TMiningModule) DB_SaveBossAward() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"bossaward": self.BossAward}})
}

//! 存储挂机信息
func (self *TMiningModule) DB_SaveGuajiInfo() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"guajitype": self.GuaJiType,
		"guajitime": self.GuajiTime}})
}

//! 挖掘矿洞
func (self *TMiningModule) DB_DigMining(index int32, value uint64) {
	filedName := fmt.Sprintf("mapdata.%d", index)
	self.MapCnt += 1
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: value,
		"mapcnt":  self.MapCnt,
		"lastpos": self.LastPos}})

}

//! 增加元素
func (self *TMiningModule) DB_AddElement(value int32) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"element": value}})
}

//! 删除元素
func (self *TMiningModule) DB_RemoveElement(value int32) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"element": value}})
}

//! 删除怪物
func (self *TMiningModule) DB_RemoveMonster(index TMiningMonster) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"monsterlst": index}})
}

//! 增加怪物
func (self *TMiningModule) DB_AddMonster(value TMiningMonster) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"monsterlst": value}})
}

//! 设置怪物
func (self *TMiningModule) DB_SetMonster(index int32, value TMiningMonster) {
	filedName := fmt.Sprintf("monsterlst.%d", index)
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: value}})
}

//! 设置怪物血量
func (self *TMiningModule) DB_SetMonsterLife(index int32, value int) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID, "monsterlst.index": index}, &bson.M{"$set": bson.M{
		"monsterlst.$.life": value}})
}

//! 重置地图上元素信息
func (self *TMiningModule) DB_ResetMapData() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"mapdata":    self.MapData,
		"element":    self.Element,
		"monsterlst": self.MonsterLst,
		"buff":       self.Buff,
		"point":      self.Point,
		"bossaward":  self.BossAward,
		"resettimes": self.ResetTimes,
		"statuscode": self.StatusCode}})
}

//! 改变矿洞状态码信息
func (self *TMiningModule) DB_UpdateMiningStatusCode() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"statuscode": self.StatusCode}})
}

//! 增加黑市商品
func (self *TMiningModule) DB_AddBuyRecord(itemid int32) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"buyrecord": itemid}})
}

//! 更改黑市购买标记
func (self *TMiningModule) DB_UpdateBuyRecord() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buyrecord": self.BuyRecord}})
}

//! 更改翻拍奖励标记
func (self *TMiningModule) DB_UpdateBossAwardMark(index int32) {
	filedName := fmt.Sprintf("bossaward.%d.status", index)
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: true}})
}

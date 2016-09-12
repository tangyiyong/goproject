package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 存储玩家当前积分
func (self *TMiningModule) DB_SavePlayerPoint() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"point": self.Point}})
}

//! 存储玩家Buff信息
func (self *TMiningModule) DB_SavePlayerBuff() {
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
		"guajitype":     self.GuaJiType,
		"guajicalctime": self.GuajiCalcTime}})
}

//! 存储行动力购买次数
func (self *TMiningModule) DB_SaveActionBuyTimes() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"actionbuytimes": self.ActionBuyTimes,
		"resetday":       self.ResetDay}})
}

//! 存储矿洞地图
func (self *TMiningModule) DB_SaveMiningMap() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"miningmap": self.MiningMap}})
}

//! 挖掘矿洞
func (self *TMiningModule) DB_DigMining(index int, value uint64) {
	filedName := fmt.Sprintf("miningmap.%d", index)
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: value,
		"lastpos": self.LastPos}})

}

//! 设置Buff次数
func (self *TMiningModule) DB_SubMiningBuffTimes(times int) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buff.times": times}})
}

//! 增加元素
func (self *TMiningModule) DB_AddElement(value int) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"element": value}})
}

//! 修改元素
func (self *TMiningModule) DB_SetElement(index int, value int) {
	filedName := fmt.Sprintf("element.%d", index)
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.Element}})
}

//! 删除元素
func (self *TMiningModule) DB_RemoveElement(index int) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"element": index}})
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
func (self *TMiningModule) DB_SetMonster(index int, value TMiningMonster) {
	filedName := fmt.Sprintf("monsterlst.%d", index)
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: value}})
}

//! 设置怪物血量
func (self *TMiningModule) DB_SetMonsterLife(index int, value int) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID, "monsterlst.index": index}, &bson.M{"$set": bson.M{
		"monsterlst.$.life": value}})
}

//! 重置地图上元素信息
func (self *TMiningModule) DB_ResetMapInfo() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"element":          self.Element,
		"monsterlst":       self.MonsterLst,
		"buff":             self.Buff,
		"point":            self.Point,
		"bossaward":        self.BossAward,
		"miningresettimes": self.MiningResetTimes,
		"statuscode":       self.StatusCode}})
}

//! 改变矿洞状态码信息
func (self *TMiningModule) DB_UpdateMiningStatusCode() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"statuscode": self.StatusCode}})
}

//! 增加黑市商品
func (self *TMiningModule) DB_AddBlackMarketMark(itemid int) {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"blackmarketbuymark": itemid}})
}

//! 更改黑市购买标记
func (self *TMiningModule) DB_UpdateBlackMarketMark() {
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"blackmarketbuymark": self.BlackMarketBuyMark}})
}

//! 更改翻拍奖励标记
func (self *TMiningModule) DB_UpdateBossAwardMark(index int) {
	filedName := fmt.Sprintf("bossaward.%d.status", index)
	mongodb.UpdateToDB("PlayerMining", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: true}})
}

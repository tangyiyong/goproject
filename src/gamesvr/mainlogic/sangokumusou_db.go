package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 更新是否结束标识
func (self *TSangokuMusouModule) UpdateIsEndMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"isend": self.IsEnd}})
}

//! 更新Buff
func (self *TSangokuMusouModule) UpdateAttr(index int, value int) {
	filedName := fmt.Sprintf("attrmarkuplst.%d.value", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: value}})
}

//! 添加Buff
func (self *TSangokuMusouModule) AddAttr(buff TSangokuMusouAttrData2) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, "attrmarkuplst", buff)
}

//! 更新通关记录
func (self *TSangokuMusouModule) UpdatePassCopyRecord() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"curstar":       self.CurStar,
		"canusestar":    self.CanUseStar,
		"passcopyid":    self.PassCopyID,
		"historystar":   self.HistoryStar,
		"historycopyid": self.HistoryCopyID}})
}

//! 增加通关信息
func (self *TSangokuMusouModule) AddPassCopyInfoLst(info TSangokuMusouCopyInfo) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, "copyinfolst", info)
}

//! 更新精英挑战次数
func (self *TSangokuMusouModule) UpdatePassEliteCopyRecord() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"passelitecopyid":  self.PassEliteCopyID,
		"elitebattletimes": self.EliteBattleTimes}})
}

//! 更新当前可使用星数
func (self *TSangokuMusouModule) UpdateCanUseStar() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"canusestar": self.CanUseStar}})
}

//! 更新当前无双秘藏以及购买状态
func (self *TSangokuMusouModule) UpdateTreasure() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"treasureid":    self.TreasureID,
		"isbuytreasure": self.IsBuyTreasure}})
}

//! 更新章节奖励领取标记
func (self *TSangokuMusouModule) UpdateChapterAwardMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"chapterawardmark": self.ChapterAwardMark}})
}

//! 更新章节Buff领取标记
func (self *TSangokuMusouModule) UpdateChapterBuffMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"chapterbuffmark": self.ChapterBuffMark}})
}

//! 重置普通挑战
func (self *TSangokuMusouModule) UpdateResetCopy() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"curstar":          self.CurStar,
		"canusestar":       self.CanUseStar,
		"passcopyid":       self.PassCopyID,
		"treasureid":       self.TreasureID,
		"isbuytreasure":    self.IsBuyTreasure,
		"battletimes":      self.BattleTimes,
		"attrmarkuplst":    self.AttrMarkupLst,
		"awardattrlst":     self.AwardAttrLst,
		"chapterawardmark": self.ChapterAwardMark,
		"chapterbuffmark":  self.ChapterBuffMark,
		"isend":            self.IsEnd,
		"copyinfolst":      self.CopyInfoLst}})
}

//! 更新重置时间戳
func (self *TSangokuMusouModule) UpdateResetTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"resetday":            self.ResetDay,
		"battletimes":         self.BattleTimes,
		"addelitebattletimes": self.AddEliteBattleTimes,
		"elitebattletimes":    self.EliteBattleTimes,
		"shoppinglst":         self.ShoppingLst}})
}

//! 更新精英挑战次数
func (self *TSangokuMusouModule) UpdateEliteBattleTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"addelitebattletimes": self.AddEliteBattleTimes,
		"elitebattletimes":    self.EliteBattleTimes}})
}

//! 更新购买次数
func (self *TSangokuMusouModule) UpdateStoreItemBuyTimes(id int, times int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID, "shoppinglst.id": id}, bson.M{"$set": bson.M{
		"shoppinglst.$.times": times}})
}

//! 增加购买信息
func (self *TSangokuMusouModule) AddStoreItemBuyInfo(info TStoreBuyData) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, "shoppinglst", info)
}

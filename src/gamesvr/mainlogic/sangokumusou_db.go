package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"
	"msg"

	"gopkg.in/mgo.v2/bson"
)

//! 更新是否结束标识
func (self *TSangokuMusouModule) DB_UpdateIsEndMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"isend": self.IsEnd}})
}

//! 更新Buff
func (self *TSangokuMusouModule) DB_UpdateAttr(index int, value int) {
	filedName := fmt.Sprintf("attrmarkuplst.%d.value", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: value}})
}

//! 添加Buff
func (self *TSangokuMusouModule) DB_AddAttr(buff TSangokuMusouAttrData2) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, "attrmarkuplst", buff)
}

//! 更新通关记录
func (self *TSangokuMusouModule) DB_UpdatePassCopyRecord() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"curstar":       self.CurStar,
		"canusestar":    self.CanUseStar,
		"passcopyid":    self.PassCopyID,
		"historystar":   self.HistoryStar,
		"historycopyid": self.HistoryCopyID}})
}

//! 增加通关信息
func (self *TSangokuMusouModule) DB_AddPassCopyInfoLst(info TSangokuMusouCopyInfo) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, "copyinfolst", info)
}

//! 更新精英挑战次数
func (self *TSangokuMusouModule) DB_UpdatePassEliteCopyRecord() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"passelitecopyid":  self.PassEliteCopyID,
		"elitebattletimes": self.EliteBattleTimes}})
}

//! 更新当前可使用星数
func (self *TSangokuMusouModule) DB_UpdateCanUseStar() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"canusestar": self.CanUseStar}})
}

//! 更新当前无双秘藏以及购买状态
func (self *TSangokuMusouModule) DB_UpdateTreasure() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"treasureid":    self.TreasureID,
		"isbuytreasure": self.IsBuyTreasure}})
}

//! 更新章节奖励领取标记
func (self *TSangokuMusouModule) DB_UpdateChapterAwardMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"chapterawardmark": self.ChapterAwardMark}})
}

//! 更新章节Buff领取标记
func (self *TSangokuMusouModule) DB_UpdateChapterBuffMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"chapterbuffmark": self.ChapterBuffMark}})
}

//! 重置普通挑战
func (self *TSangokuMusouModule) DB_UpdateResetCopy() {
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
func (self *TSangokuMusouModule) DB_UpdateResetTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"resetday":            self.ResetDay,
		"battletimes":         self.BattleTimes,
		"addelitebattletimes": self.AddEliteBattleTimes,
		"elitebattletimes":    self.EliteBattleTimes,
		"buyrecord":           self.BuyRecord}})
}

//! 更新精英挑战次数
func (self *TSangokuMusouModule) DB_UpdateEliteBattleTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"addelitebattletimes": self.AddEliteBattleTimes,
		"elitebattletimes":    self.EliteBattleTimes}})
}

//! 更新购买次数
func (self *TSangokuMusouModule) DB_UpdateStoreItemBuyTimes(index int, times int) {
	filedName := fmt.Sprintf("buyrecord.%d.times", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: times}})
}

//! 增加购买信息
func (self *TSangokuMusouModule) DB_AddStoreItemBuyInfo(info msg.MSG_BuyData) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerSangokuMusou", bson.M{"_id": self.PlayerID}, "buyrecord", info)
}

package mainlogic

import (
	"fmt"
	"mongodb"
	"msg"

	"gopkg.in/mgo.v2/bson"
)

func (self *THeroSoulsModule) DB_SaveHeroSoulsLst() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"herosoulslst": self.HeroSoulsLst,
		"targetindex":  self.TargetIndex}})
}

func (self *THeroSoulsModule) DB_SaveHeroSoulsStoreLst() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"herosoulsstorelst": self.HeroSoulsStoreLst}})
}

func (self *THeroSoulsModule) DB_SaveHeroSoulsRefreshMark() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"storemark": self.StoreMark}})
}

func (self *THeroSoulsModule) DB_UpdateHeroSoulsMark(index int) {
	filedName := fmt.Sprintf("herosoulslst.%d.isexist", index)
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: false}})
}

func (self *THeroSoulsModule) DB_AddHeroSoulsLink(link msg.MSG_HeroSoulsLink) {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"herosoulslink": link}})
}

func (self *THeroSoulsModule) DB_UnLockChapter() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"unlockchapter": self.UnLockChapter}})
}

func (self *THeroSoulsModule) DB_UpdateHeroSoulsLinkLevel(index int, level int) {
	filedName := fmt.Sprintf("herosoulslink.%d.level", index)
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: level}})
}

func (self *THeroSoulsModule) DB_UpdateSoulMapValue() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"soulmapvalue": self.SoulMapValue}})
}

func (self *THeroSoulsModule) DB_UpdateTargetIndex() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"targetindex": self.TargetIndex}})
}

func (self *THeroSoulsModule) DB_Reset() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buytimes":  self.BuyTimes,
		"storemark": self.StoreMark,
		"lefttimes": self.LeftTimes,
		"resetday":  self.ResetDay}})
}

func (self *THeroSoulsModule) DB_UpdateChallengeHeroSoulsTimes() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"lefttimes": self.LeftTimes}})
}

func (self *THeroSoulsModule) DB_BuyChallengeHeroSoulsTimes() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buytimes":  self.BuyTimes,
		"lefttimes": self.LeftTimes}})
}

func (self *THeroSoulsModule) DB_UpdateStoreGoodsStatus(index int, status bool) {
	filedName := fmt.Sprintf("herosoulsstorelst.%d.isbuy", index)
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: status}})
}

func (self *THeroSoulsModule) DB_UpdateSoulMapAchievement() {
	mongodb.UpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"achievement": self.Achievement}})
}

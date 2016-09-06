package mainlogic

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func (self *THeroSoulsModule) DB_SaveHeroSoulsLst() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"herosoulslst": self.HeroSoulsLst,
		"targetindex":  self.TargetIndex}})
}

func (self *THeroSoulsModule) DB_SaveHeroSoulsStoreLst() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"herosoulsstorelst": self.HeroSoulsStoreLst}})
}

func (self *THeroSoulsModule) DB_SaveHeroSoulsRefreshMark() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"refreshstoretimemark": self.RefreshStoreTimeMark}})
}

func (self *THeroSoulsModule) DB_UpdateHeroSoulsMark(index int) {
	filedName := fmt.Sprintf("herosoulslst.%d.isexist", index)
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: false}})
}

func (self *THeroSoulsModule) DB_AddHeroSoulsLink(link THeroSoulsLink) {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"herosoulslink": link}})
}

func (self *THeroSoulsModule) DB_UnLockChapter() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"unlockchapter": self.UnLockChapter}})
}

func (self *THeroSoulsModule) DB_UpdateHeroSoulsLinkLevel(index int, level int) {
	filedName := fmt.Sprintf("herosoulslink.%d.level", index)
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: level}})
}

func (self *THeroSoulsModule) DB_UpdateSoulMapValue() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"soulmapvalue": self.SoulMapValue}})
}

func (self *THeroSoulsModule) DB_UpdateTargetIndex() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"targetindex": self.TargetIndex}})
}

func (self *THeroSoulsModule) DB_Reset() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buychallengetimes":    self.BuyChallengeTimes,
		"refreshstoretimemark": self.RefreshStoreTimeMark,
		"challengetimes":       self.ChallengeTimes,
		"resetday":             self.ResetDay}})
}

func (self *THeroSoulsModule) DB_UpdateChallengeHeroSoulsTimes() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"challengetimes": self.ChallengeTimes}})
}

func (self *THeroSoulsModule) DB_BuyChallengeHeroSoulsTimes() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buychallengetimes": self.BuyChallengeTimes,
		"challengetimes":    self.ChallengeTimes}})
}

func (self *THeroSoulsModule) DB_UpdateStoreGoodsStatus(index int, status bool) {
	filedName := fmt.Sprintf("herosoulsstorelst.%d.isbuy", index)
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: status}})
}

func (self *THeroSoulsModule) DB_UpdateSoulMapAchievement() {
	GameSvrUpdateToDB("PlayerHeroSouls", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"achievement": self.Achievement}})
}

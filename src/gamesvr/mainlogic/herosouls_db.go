package mainlogic

import (
	"appconfig"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (self *THeroSoulsModule) DB_SaveHeroSoulsLst() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"herosoulslst": self.HeroSoulsLst,
		"targetindex":  self.TargetIndex}})
}

func (self *THeroSoulsModule) DB_SaveHeroSoulsStoreLst() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"herosoulsstorelst": self.HeroSoulsStoreLst}})
}

func (self *THeroSoulsModule) DB_SaveHeroSoulsRefreshMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"refreshstoretimemark": self.RefreshStoreTimeMark}})
}

func (self *THeroSoulsModule) DB_UpdateHeroSoulsMark(index int) {
	filedName := fmt.Sprintf("herosoulslst.%d.isexist", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: false}})
}

func (self *THeroSoulsModule) DB_AddHeroSoulsLink(link THeroSoulsLink) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, "herosoulslink", link)
}

func (self *THeroSoulsModule) DB_UnLockChapter() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"unlockchapter": self.UnLockChapter}})
}

func (self *THeroSoulsModule) DB_UpdateHeroSoulsLinkLevel(index int, level int) {
	filedName := fmt.Sprintf("herosoulslink.%d.level", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: level}})
}

func (self *THeroSoulsModule) DB_UpdateSoulMapValue() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"soulmapvalue": self.SoulMapValue}})
}

func (self *THeroSoulsModule) DB_UpdateTargetIndex() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"targetindex": self.TargetIndex}})
}

func (self *THeroSoulsModule) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"buychallengetimes":    self.BuyChallengeTimes,
		"refreshstoretimemark": self.RefreshStoreTimeMark,
		"challengetimes":       self.ChallengeTimes,
		"resetday":             self.ResetDay}})
}

func (self *THeroSoulsModule) DB_UpdateChallengeHeroSoulsTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"challengetimes": self.ChallengeTimes}})
}

func (self *THeroSoulsModule) DB_BuyChallengeHeroSoulsTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"buychallengetimes": self.BuyChallengeTimes,
		"challengetimes":    self.ChallengeTimes}})
}

func (self *THeroSoulsModule) DB_UpdateStoreGoodsStatus(index int, status bool) {
	filedName := fmt.Sprintf("herosoulsstorelst.%d.isbuy", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		filedName: status}})
}

func (self *THeroSoulsModule) DB_UpdateSoulMapAchievement() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerHeroSouls", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"achievement": self.Achievement}})
}

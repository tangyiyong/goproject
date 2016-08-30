package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (main_copy *TCopyMoudle) UpdateMainCopyAt(copyindex int) {
	fieldName := fmt.Sprintf("main.copyinfo.%d", copyindex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, bson.M{"$set": bson.M{
		"main.curcopyid":  main_copy.Main.CurCopyID,
		"main.curchapter": main_copy.Main.CurChapter,
		fieldName:         main_copy.Main.CopyInfo[copyindex]}})
}

func (main_copy *TCopyMoudle) UpdateEliteCopyAt(copyindex int) {
	fieldName := fmt.Sprintf("elite.copyinfo.%d", copyindex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, bson.M{"$set": bson.M{
		"elite.curcopyid":  main_copy.Elite.CurCopyID,
		"elite.curchapter": main_copy.Elite.CurChapter,
		fieldName:          main_copy.Elite.CopyInfo[copyindex]}})
}

func (daily_copy *TCopyMoudle) UpdateDailyCopyMask(index int, mask bool) {
	filedName := fmt.Sprintf("daily.copyinfo.%d.ischallenge", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": daily_copy.PlayerID}, bson.M{"$set": bson.M{
		filedName: mask}})
}

func (famous_copy *TCopyMoudle) UpdateFamousCopyBattleTimes(chapter int, copyIndex int, battleTimes int) {
	filedname := fmt.Sprintf("famous.chapter.%d.passedcopy.%d.battletimes", chapter, copyIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": famous_copy.PlayerID}, bson.M{"$set": bson.M{
		filedname: battleTimes}})
}

func (famous_copy *TCopyMoudle) UpdateFamousCopyTotalBattleTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": famous_copy.PlayerID}, bson.M{"$set": bson.M{
		"famous.battletimes": famous_copy.Famous.BattleTimes}})
}

func (famous_copy *TCopyMoudle) UpdateFamousCopyCurCopyID() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": famous_copy.PlayerID}, bson.M{"$set": bson.M{
		"famous.curcopyid": famous_copy.Famous.CurCopyID}})
}

func (famous_copy *TCopyMoudle) IncFamousCopy(chapter int, famousCopy TFamousCopy) {
	filedname := fmt.Sprintf("famous.chapter.%d.passedcopy", chapter)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": famous_copy.PlayerID}, bson.M{"$push": bson.M{
		filedname: famousCopy}})
}

func (self *TCopyMoudle) UpdateCopy() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"main":     self.Main,
		"elite":    self.Elite,
		"daily":    self.Daily,
		"famous":   self.Famous,
		"resetday": self.ResetDay}})
}

func (main_copy *TCopyMoudle) UpdateGuaJiAwardData(chapter int) {
	filedName := fmt.Sprintf("guaji.awarddata.%d", chapter)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, bson.M{"$set": bson.M{
		filedName: true}})
}

//! 存储主线关卡星级奖励
func (main_copy *TCopyMoudle) UpdateMainAward(chapter int) {
	filedName1 := fmt.Sprintf("main.chapter.%d.staraward", chapter)
	filedName2 := fmt.Sprintf("main.chapter.%d.sceneaward", chapter)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, bson.M{"$set": bson.M{
		filedName1: main_copy.Main.Chapter[chapter].StarAward,
		filedName2: main_copy.Main.Chapter[chapter].SceneAward}})
}

//! 存储精英关卡星级奖励
func (main_copy *TCopyMoudle) UpdateEliteAward(chapter int) {
	filedName1 := fmt.Sprintf("elite.chapter.%d.staraward", chapter)
	filedName2 := fmt.Sprintf("elite.chapter.%d.sceneaward", chapter)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, bson.M{"$set": bson.M{
		filedName1: main_copy.Elite.Chapter[chapter].StarAward,
		filedName2: main_copy.Elite.Chapter[chapter].SceneAward}})
}

//! 更新精英副本入侵时间信息
func (elite_copy *TCopyMoudle) UpdateEliteInvadeTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, bson.M{"$set": bson.M{
		"lastinvadetime": elite_copy.LastInvadeTime}})
}

//! 更新名将副本章节奖励
func (elite_copy *TCopyMoudle) UpdateFamousAward(chapter int) {
	filedName1 := fmt.Sprintf("famous.chapter.%d.chapteraward", chapter)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, bson.M{"$set": bson.M{
		filedName1: true}})
}

//! 增加章节信息
func (main_copy *TCopyMoudle) AddMainChapterInfo(chapter TMainChapter) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, "main.chapter", chapter)

	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, bson.M{"$set": bson.M{
		"main.curcopyid":  main_copy.Main.CurCopyID,
		"main.curchapter": main_copy.Main.CurChapter}})
}

func (elite_copy *TCopyMoudle) AddEliteChapterInfo(chapter TEliteChapter) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, "elite.chapter", chapter)

	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, bson.M{"$set": bson.M{
		"elite.curcopyid":  elite_copy.Elite.CurCopyID,
		"elite.curchapter": elite_copy.Elite.CurChapter}})
}

//! 增加关卡信息
func (main_copy *TCopyMoudle) AddMainCopyInfo(copyInfo TMainCopy) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, "main.copyinfo", copyInfo)

	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, bson.M{"$set": bson.M{
		"main.curcopyid":  main_copy.Main.CurCopyID,
		"main.curchapter": main_copy.Main.CurChapter}})
}

func (elite_copy *TCopyMoudle) AddEliteCopyInfo(copyInfo TEliteCopy) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, "elite.copyinfo", copyInfo)

	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, bson.M{"$set": bson.M{
		"elite.curcopyid":  elite_copy.Elite.CurCopyID,
		"elite.curchapter": elite_copy.Elite.CurChapter}})
}

//! 入侵增删
func (elite_copy *TCopyMoudle) AddEliteInvade(chapter int) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, "elite.invadechapter", chapter)
}

func (elite_copy *TCopyMoudle) RemoveEliteInvade(chapter int) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, "elite.invadechapter", chapter)
}

func (main_copy *TCopyMoudle) UpdateMainCopyInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": main_copy.PlayerID}, bson.M{"$set": bson.M{
		"main.curcopyid":  main_copy.Main.CurCopyID,
		"main.curchapter": main_copy.Main.CurChapter}})
}

func (elite_copy *TCopyMoudle) UpdateEliteCopyInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCopy", bson.M{"_id": elite_copy.PlayerID}, bson.M{"$set": bson.M{
		"elite.curcopyid":  elite_copy.Elite.CurCopyID,
		"elite.curchapter": elite_copy.Elite.CurChapter}})
}

package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (self *TCopyMoudle) DB_UpdateMainCopyAt(copyindex int) {
	fieldName := fmt.Sprintf("main.copylst.%d", copyindex)
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"main.curid":      self.Main.CurID,
		"main.curchapter": self.Main.CurChapter,
		fieldName:         self.Main.CopyLst[copyindex]}})
}

func (self *TCopyMoudle) DB_UpdateEliteCopyAt(copyindex int) {
	fieldName := fmt.Sprintf("elite.copylst.%d", copyindex)
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"elite.curid":      self.Elite.CurID,
		"elite.curchapter": self.Elite.CurChapter,
		fieldName:          self.Elite.CopyLst[copyindex]}})
}

func (self *TCopyMoudle) DB_UpdateDailyCopyMask(index int, mask bool) {
	filedName := fmt.Sprintf("daily.copyinfo.%d.ischallenge", index)
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: mask}})
}

func (self *TCopyMoudle) DB_AddFamousPassCopy(chapter int, cpoyid int) {
	filedname := fmt.Sprintf("famous.chapter.%d.passedcopy", chapter)
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{
		filedname: cpoyid}})
}

func (self *TCopyMoudle) DB_UpdateFamousExtra(chapter int) {
	filedname := fmt.Sprintf("famous.chapter.%d.extra", chapter)
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedname: true}})
}

func (self *TCopyMoudle) DB_UpdateFamousCopyData() {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"famous.curid":      self.Famous.CurID,
		"famous.curchapter": self.Famous.CurChapter,
		"famous.times":      self.Famous.Times}})
}

func (self *TCopyMoudle) DB_UpdateCopy() {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"main":     self.Main,
		"elite":    self.Elite,
		"daily":    self.Daily,
		"famous":   self.Famous,
		"resetday": self.ResetDay}})
}

//! 存储主线关卡星级奖励
func (self *TCopyMoudle) DB_UpdateMainAward(chapter int) {
	filedName1 := fmt.Sprintf("main.chapter.%d.staraward", chapter)
	filedName2 := fmt.Sprintf("main.chapter.%d.sceneaward", chapter)
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName1: self.Main.Chapter[chapter].StarAward,
		filedName2: self.Main.Chapter[chapter].SceneAward}})
}

//! 存储精英关卡星级奖励
func (self *TCopyMoudle) DB_UpdateEliteAward(chapter int) {
	filedName1 := fmt.Sprintf("elite.chapter.%d.staraward", chapter)
	filedName2 := fmt.Sprintf("elite.chapter.%d.sceneaward", chapter)
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName1: self.Elite.Chapter[chapter].StarAward,
		filedName2: self.Elite.Chapter[chapter].SceneAward}})
}

//! 更新精英副本入侵时间信息
func (self *TCopyMoudle) DB_UpdateEliteInvadeTime() {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"lastinvadetime": self.LastInvadeTime}})
}

//! 更新名将副本章节奖励
func (self *TCopyMoudle) DB_UpdateFamousAward(chapter int) {
	filedName1 := fmt.Sprintf("famous.chapter.%d.chapteraward", chapter)
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName1: true}})
}

//! 增加章节信息
func (self *TCopyMoudle) DB_AddMainChapterInfo(chapter TMainChapter) {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"main.chapter": chapter}})
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"main.curid":      self.Main.CurID,
		"main.curchapter": self.Main.CurChapter}})
}

func (self *TCopyMoudle) DB_AddEliteChapterInfo(chapter TEliteChapter) {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"elite.chapter": chapter}})
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"elite.curid":      self.Elite.CurID,
		"elite.curchapter": self.Elite.CurChapter}})
}

//! 增加关卡信息
func (self *TCopyMoudle) DB_AddMainCopyInfo(copyInfo TMainCopy) {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"main.copyinfo": copyInfo}})
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"main.curid":      self.Main.CurID,
		"main.curchapter": self.Main.CurChapter}})
}

func (self *TCopyMoudle) DB_AddEliteCopyInfo(copyInfo TEliteCopy) {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"elite.copyinfo": copyInfo}})
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"elite.curid":      self.Elite.CurID,
		"elite.curchapter": self.Elite.CurChapter}})
}

//! 入侵增删
func (self *TCopyMoudle) DB_AddEliteInvade(chapter int) {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"elite.invadechapter": chapter}})
}

func (self *TCopyMoudle) DB_RemoveEliteInvade(chapter int) {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"elite.invadechapter": chapter}})
}

func (self *TCopyMoudle) DB_UpdateMainCopyInfo() {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"main.curid":      self.Main.CurID,
		"main.curchapter": self.Main.CurChapter}})
}

func (self *TCopyMoudle) DB_UpdateEliteCopyInfo() {
	mongodb.UpdateToDB("PlayerCopy", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"elite.curid":      self.Elite.CurID,
		"elite.curchapter": self.Elite.CurChapter}})
}

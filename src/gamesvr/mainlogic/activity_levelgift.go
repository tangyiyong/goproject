package mainlogic

import (
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TLevelGiftInfo struct {
	GiftID   int   //! 礼包ID
	BuyTimes int   //! 当前可购买次数
	DeadLine int64 //! 过期时间
}

//! 等级礼包
type TActivityLevelGift struct {
	ActivityID int //! 活动ID

	GiftLst       []TLevelGiftInfo //! 等级礼包
	IsHaveNewItem bool             //! 红点显示规则

	VersionCode int32 //! 版本号
	ResetCode   int32 //! 迭代号

	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityLevelGift) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityLevelGift) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr

	self.GiftLst = []TLevelGiftInfo{}
	self.IsHaveNewItem = false

	self.activityModule.activityPtrs[activityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityLevelGift) Refresh(versionCode int32) {
	//! 检测物品过期
	self.CheckDeadLine()

	self.VersionCode = versionCode
	self.DB_Refresh()
}

func (self *TActivityLevelGift) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.GiftLst = []TLevelGiftInfo{}
	self.IsHaveNewItem = false
	self.DB_Reset()
}

func (self *TActivityLevelGift) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityLevelGift) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityLevelGift) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.IsHaveNewItem == true {
		return true
	}

	return false
}

//! 获取等级礼包信息
func (self *TActivityLevelGift) GetLevelGiftInfo(giftID int) *TLevelGiftInfo {
	length := len(self.GiftLst)
	for i := 0; i < length; i++ {
		if self.GiftLst[i].GiftID == giftID {
			return &self.GiftLst[i]
		}
	}

	return nil
}

//! 升级检测
func (self *TActivityLevelGift) CheckLevelUp(level int) {
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	giftLst := gamedata.GetLevelGiftLst(awardType)
	length := len(giftLst)

	now := time.Now().Unix()

	for i := 0; i < length; i++ {
		needLevel, _ := strconv.Atoi(giftLst[i].Level)
		if needLevel > level {
			continue
		}

		isExist := false
		for j := 0; j < len(self.GiftLst); j++ {
			if self.GiftLst[j].GiftID == giftLst[i].ID {
				isExist = true
				break
			}
		}

		if isExist == false {
			var gift TLevelGiftInfo
			gift.GiftID = giftLst[i].ID
			gift.BuyTimes = giftLst[i].BuyTimes
			gift.DeadLine = now + int64(giftLst[i].DeadLine)
			self.GiftLst = append(self.GiftLst, gift)
			self.IsHaveNewItem = true
			self.DB_AddGift(&gift)
		}
	}
}

//! 检测过期时间
func (self *TActivityLevelGift) CheckDeadLine() {
	now := time.Now().Unix()
	length := len(self.GiftLst)
	for i := 0; i < length; i++ {
		if self.GiftLst[i].DeadLine <= now || self.GiftLst[i].BuyTimes == 0 { //! 过期或者可购买次数为零
			self.DB_RemoveDeadGift(&self.GiftLst[i])

			if i == 0 {
				self.GiftLst = self.GiftLst[1:]
			} else if (i + 1) == len(self.GiftLst) {
				self.GiftLst = self.GiftLst[:i]
			} else {
				self.GiftLst = append(self.GiftLst[:i], self.GiftLst[i+1:]...)
			}

			length = len(self.GiftLst)
			i--
		}
	}
}

func (self *TActivityLevelGift) DB_Refresh() {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"levelgift.versioncode": self.VersionCode}})
}

func (self *TActivityLevelGift) DB_Reset() {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"levelgift.activityid":    self.ActivityID,
		"levelgift.resetcode":     self.ResetCode,
		"levelgift.ishavenewitem": self.IsHaveNewItem,
		"levelgift.giftlst":       self.GiftLst,
		"levelgift.versioncode":   self.VersionCode}})
}

func (self *TActivityLevelGift) DB_RemoveDeadGift(gift *TLevelGiftInfo) {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$pull": bson.M{"levelgift.giftlst": *gift}})
}

func (self *TActivityLevelGift) DB_AddGift(gift *TLevelGiftInfo) {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$push": bson.M{"levelgift.giftlst": *gift}})

	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"levelgift.ishavenewitem": self.IsHaveNewItem}})
}

func (self *TActivityLevelGift) DB_UpdateGiftLst() {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"levelgift.giftlst": self.GiftLst}})
}

func (self *TActivityLevelGift) DB_UpdateBuyTimes(id int, times int) {
	index := -1
	for i, v := range self.GiftLst {
		if v.GiftID == id {
			index = i
		}
	}

	if index < 0 {
		gamelog.Error("DB_UpdateBuyTimes Fail: Not find week gift id: %d", id)
		return
	}

	filedName := fmt.Sprintf("levelgift.giftlst.%d.buytimes", index)
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID},
		&bson.M{"$set": bson.M{filedName: times}})
}

func (self *TActivityLevelGift) DB_UpdateNewItemMark() {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID},
		&bson.M{"$set": bson.M{"levelgift.ishavenewitem": self.IsHaveNewItem}})
}

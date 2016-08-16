package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

type TRankGiftInfo struct {
	GiftID   int //! 礼包ID
	BuyTimes int //! 当前可购买次数
}

//! 等级礼包
type TActivityRankGift struct {
	ActivityID int //! 活动ID

	GiftLst       []TRankGiftInfo //! 等级礼包
	IsHaveNewItem bool            //! 红点显示规则

	VersionCode int //! 版本号
	ResetCode   int //! 迭代号

	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityRankGift) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityRankGift) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr

	self.GiftLst = []TRankGiftInfo{}
	self.IsHaveNewItem = false

	self.activityModule.activityPtrs[activityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityRankGift) Refresh(versionCode int) {
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

func (self *TActivityRankGift) End(versionCode int, resetCode int) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.GiftLst = []TRankGiftInfo{}
	self.IsHaveNewItem = false
	go self.DB_Reset()
}

func (self *TActivityRankGift) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityRankGift) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityRankGift) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.IsHaveNewItem == true {
		return true
	}

	return false
}

//! 获取排名礼包信息
func (self *TActivityRankGift) GetRankGiftInfo(giftID int) *TRankGiftInfo {
	length := len(self.GiftLst)
	for i := 0; i < length; i++ {
		if self.GiftLst[i].GiftID == giftID {
			return &self.GiftLst[i]
		}
	}

	return nil
}

//! 排名检测
func (self *TActivityRankGift) CheckRankUp(rank int) {
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	giftLst := gamedata.GetLevelGiftLst(awardType)
	length := len(giftLst)

	for i := 0; i < length; i++ {
		if giftLst[i].Level > rank {
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
			var gift TRankGiftInfo
			gift.GiftID = giftLst[i].ID
			gift.BuyTimes = giftLst[i].BuyTimes
			self.GiftLst = append(self.GiftLst, gift)
			self.IsHaveNewItem = true
			go self.DB_AddGift(&gift)
		}
	}
}

func (self *TActivityRankGift) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"rankgift.versioncode": self.VersionCode}})
	return true
}

func (self *TActivityRankGift) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"rankgift.activityid":    self.ActivityID,
		"rankgift.resetcode":     self.ResetCode,
		"rankgift.ishavenewitem": self.IsHaveNewItem,
		"rankgift.giftlst":       self.GiftLst,
		"rankgift.versioncode":   self.VersionCode}})
	return true
}

func (self *TActivityRankGift) DB_AddGift(gift *TRankGiftInfo) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, "rankgift.giftlst", *gift)

	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"rankgift.ishavenewitem": self.IsHaveNewItem}})
}

func (self *TActivityRankGift) DB_UpdateBuyTimes(id int, times int) {
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

	filedName := fmt.Sprintf("rankgift.giftlst.%d.buytimes", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID},
		bson.M{"$set": bson.M{filedName: times}})
}

func (self *TActivityRankGift) DB_UpdateNewItemMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"rankgift.ishavenewitem": self.IsHaveNewItem}})
}

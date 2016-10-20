package mainlogic

import (
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"strconv"
	"strings"
)

type TRankGiftInfo struct {
	GiftID   int32 //! 礼包ID
	BuyTimes int   //! 当前可购买次数
}

//! 等级礼包
type TActivityRankGift struct {
	ActivityID int32 //! 活动ID

	GiftLst       []TRankGiftInfo //! 等级礼包
	IsHaveNewItem bool            //! 红点显示规则

	VersionCode int32 //! 版本号
	ResetCode   int32 //! 迭代号

	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityRankGift) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityRankGift) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
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
func (self *TActivityRankGift) Refresh(versionCode int32) {
	self.VersionCode = versionCode
	self.DB_Refresh()
}

func (self *TActivityRankGift) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.GiftLst = []TRankGiftInfo{}
	self.IsHaveNewItem = false
	self.DB_Reset()
}

func (self *TActivityRankGift) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityRankGift) GetResetV() int32 {
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
func (self *TActivityRankGift) GetRankGiftInfo(giftID int32) *TRankGiftInfo {
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
		values := strings.Split(giftLst[i].Level, "|")
		rankValueLst := IntLst{}
		for _, n := range values {
			rank, _ := strconv.Atoi(n)
			rankValueLst = append(rankValueLst, rank)
		}

		if len(rankValueLst) != 2 { //! 当前名次没有奖励, 直接返回
			//gamelog.Error("GetLevelGiftLst Level Error: Can't split rank num %v", rankValueLst)
			return
		}

		if rankValueLst[1] >= rank && rankValueLst[0] <= rank {
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
			self.DB_AddGift(&gift)
		}
	}
}

func (self *TActivityRankGift) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"rankgift.versioncode": self.VersionCode}})
}

func (self *TActivityRankGift) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"rankgift.activityid":    self.ActivityID,
		"rankgift.resetcode":     self.ResetCode,
		"rankgift.ishavenewitem": self.IsHaveNewItem,
		"rankgift.giftlst":       self.GiftLst,
		"rankgift.versioncode":   self.VersionCode}})
}

func (self *TActivityRankGift) DB_AddGift(gift *TRankGiftInfo) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$push": bson.M{"rankgift.giftlst": *gift}})

	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"rankgift.ishavenewitem": self.IsHaveNewItem}})
}

func (self *TActivityRankGift) DB_UpdateBuyTimes(id int32, times int) {
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
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID},
		&bson.M{"$set": bson.M{filedName: times}})
}

func (self *TActivityRankGift) DB_UpdateNewItemMark() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"rankgift.ishavenewitem": self.IsHaveNewItem}})
}

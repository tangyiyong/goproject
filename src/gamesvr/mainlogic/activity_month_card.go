package mainlogic

import (
	"fmt"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 月卡活动
type TActivityMonthCard struct {
	ActivityID int //! 活动ID

	CardDays   []int  //! 月卡状态表
	CardStatus []bool //! 月卡领取状态

	VersionCode    int32            //! 版本号
	ResetCode      int32            //! 迭代号
	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityMonthCard) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityMonthCard) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr

	count := gamedata.GetMonthCardCount()
	self.CardDays = make([]int, count)
	self.CardStatus = make([]bool, count)
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityMonthCard) Refresh(versionCode int32) {
	for i, v := range self.CardStatus {
		if v == true {
			self.CardStatus[i] = false
		}

		//! 减去过去的天数
		self.CardDays[i] -= int(versionCode - self.VersionCode)
		if self.CardDays[i] < 0 {
			self.CardDays[i] = 0
		}
	}

	self.VersionCode = versionCode
	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityMonthCard) End(versionCode int32, resetCode int32) {
	self.ResetCode = resetCode
	self.VersionCode = versionCode
	self.DB_Reset()
}

func (self *TActivityMonthCard) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityMonthCard) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityMonthCard) RedTip() bool {
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	for _, v := range self.CardStatus {
		if v == false {
			return true
		}
	}

	return false
}

//! 重置
func (self *TActivityMonthCard) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"monthcard.activityid":  self.ActivityID,
		"monthcard.resetcode":   self.ResetCode,
		"monthcard.versioncode": self.VersionCode}})
}

//! 更新月卡状态与时间
func (self *TActivityMonthCard) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"monthcard.carddays":    self.CardDays,
		"monthcard.versioncode": self.VersionCode,
		"monthcard.cardstatus":  self.CardStatus}})
}

func (self *TActivityMonthCard) DB_UpdateCardStatus() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"monthcard.carddays":   self.CardDays,
		"monthcard.cardstatus": self.CardStatus}})
}

//! 更新月卡天数
func (self *TActivityMonthCard) DB_UpdateCardDays(index int, days int) {
	filedName := fmt.Sprintf("monthcard.carddays.%d", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		filedName: days}})
}

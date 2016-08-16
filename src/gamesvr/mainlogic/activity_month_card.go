package mainlogic

import (
	"appconfig"
	"fmt"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 月卡活动
type TActivityMonthCard struct {
	ActivityID int //! 活动ID

	CardDays   []int  //! 月卡状态表
	CardStatus []bool //! 月卡领取状态

	VersionCode    int              //! 版本号
	ResetCode      int              //! 迭代号
	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityMonthCard) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityMonthCard) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
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
func (self *TActivityMonthCard) Refresh(versionCode int) {
	for i, v := range self.CardStatus {
		if v == true {
			self.CardStatus[i] = false
		}

		//! 减去过去的天数
		self.CardDays[i] -= (versionCode - self.VersionCode)
		if self.CardDays[i] < 0 {
			self.CardDays[i] = 0
		}
	}

	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityMonthCard) End(versionCode int, resetCode int) {
	self.ResetCode = resetCode
	self.VersionCode = versionCode
	go self.DB_Reset()
}

func (self *TActivityMonthCard) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityMonthCard) GetResetV() int {
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
func (self *TActivityMonthCard) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"monthcard.activityid":  self.ActivityID,
		"monthcard.resetcode":   self.ResetCode,
		"monthcard.versioncode": self.VersionCode}})
	return true
}

//! 更新月卡状态与时间
func (self *TActivityMonthCard) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"monthcard.carddays":    self.CardDays,
		"monthcard.versioncode": self.VersionCode,
		"monthcard.cardstatus":  self.CardStatus}})
	return true
}

func (self *TActivityMonthCard) DB_UpdateCardStatus() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"monthcard.carddays":   self.CardDays,
		"monthcard.cardstatus": self.CardStatus}})
}

//! 更新月卡天数
func (self *TActivityMonthCard) DB_UpdateCardDays(index int, days int) {
	filedName := fmt.Sprintf("monthcard.carddays.%d", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: days}})
}

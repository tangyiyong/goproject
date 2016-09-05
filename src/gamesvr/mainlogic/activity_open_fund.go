package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var G_BuyFundNum int

func InitBuyOpenFundNum() bool {
	s := mongodb.GetDBSession()
	defer s.Close()
	count, err := s.DB(appconfig.GameDbName).C("PlayerActivity").Find(bson.M{"openfund.isbuyfund": true}).Count()
	if err != nil {
		if err != mgo.ErrNotFound {
			gamelog.Error("Init DB Error!!!")
			return false
		}
	}

	G_BuyFundNum = count
	return true
}

//! 开服基金活动
type TActivityOpenFund struct {
	ActivityID int //! 活动ID

	FundLevelMark Mark //! 等级奖励领取标记
	FundCountMark Mark //! 购买基金人数奖励领取
	IsBuyFund     bool //! 购买基金标记

	VersionCode    int32            //! 版本号
	ResetCode      int32            //! 迭代号
	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityOpenFund) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityOpenFund) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityOpenFund) Refresh(versionCode int32) {
	//! 此活动没有刷新
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityOpenFund) End(versionCode int32, resetCode int32) {
	self.FundLevelMark = 0
	self.FundCountMark = 0
	self.IsBuyFund = false
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	go self.DB_Reset()
}

func (self *TActivityOpenFund) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityOpenFund) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityOpenFund) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.IsBuyFund == false {
		return false
	}

	tempLst := IntLst{}
	for _, v := range gamedata.GT_OpenFundLst[0] {
		if self.activityModule.ownplayer.GetLevel() >= v.Count {
			tempLst.Add(v.ID)
		}
	}

	for _, v := range tempLst {
		if self.FundLevelMark.Get(uint32(v)) == false {
			return true
		}
	}

	tempLst = IntLst{}
	for _, v := range gamedata.GT_OpenFundLst[1] {
		if G_BuyFundNum >= v.Count {
			tempLst.Add(v.ID)
		}
	}

	for _, v := range tempLst {
		if self.FundCountMark.Get(uint32(v)) == false {
			return true
		}
	}

	return false
}

func (self *TActivityOpenFund) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"openfund.activityid":    self.ActivityID,
		"openfund.fundlevelmark": self.FundLevelMark,
		"openfund.fundcountmark": self.FundCountMark,
		"openfund.isbuyfund":     self.IsBuyFund,
		"openfund.versioncode":   self.VersionCode,
		"openfund.resetcode":     self.ResetCode}})
}

func (self *TActivityOpenFund) DB_Refresh() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"openfund.versioncode": self.VersionCode}})
}

//! 更改购买标记
func (self *TActivityOpenFund) UpdateBuyFundMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"openfund.isbuyfund": self.IsBuyFund}})
}

//! 更改基金全民奖励领取标记
func (self *TActivityOpenFund) UpdateFundCountMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"openfund.fundcountmark": self.FundCountMark}})
}

//! 更改基金等级奖励领取标记
func (self *TActivityOpenFund) UpdateFundLevelMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"openfund.fundlevelmark": self.FundLevelMark}})
}

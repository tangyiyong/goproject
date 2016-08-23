package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

const (
	SignPlus_Can_Not_Receive = 0
	SignPlus_Can_Receive     = 1
	SignPlus_Have_Received   = 2
)

//! 登录送礼活动
type TActivitySign struct {
	ActivityID int //! 活动ID

	SignDay        int                    //! 签到天数
	IsSign         bool                   //! 签到状态
	SignPlusAward  []gamedata.ST_ItemData //! 豪华签到奖励
	IsSignPlus     bool                   //! 豪华签到更新时间
	SignPlusStatus int                    //! 豪华签到状态
	VersionCode    int32                  //! 版本号
	ResetCode      int32                  //! 迭代号

	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivitySign) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivitySign) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.IsSign = false
	self.IsSignPlus = false
	self.activityModule.activityPtrs[activityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode

	self.SignDay = 0

	//! 刷新豪华签到信息
	self.RefreshSignPlusAward(false)
}

//! 刷新数据
func (self *TActivitySign) Refresh(versionCode int32) {
	//! 刷新签到标记
	self.IsSign = false
	self.IsSignPlus = false
	self.VersionCode = versionCode

	awardCount := gamedata.GetSignAwardCount()
	if self.SignDay > awardCount {
		self.SignDay = awardCount - 30
	}

	//! 获取奖励内容
	data := gamedata.GetSignPlusDataFromLevel(self.activityModule.ownplayer.GetLevel())
	if data == nil {
		gamelog.Error("RefreshSignPlusAward fail")
		return
	}

	if len(self.SignPlusAward) != 0 {
		self.SignPlusAward = []gamedata.ST_ItemData{}
	}

	awardLst := gamedata.GetItemsFromAwardIDEx(data.SignAward)
	self.SignPlusAward = append(self.SignPlusAward, awardLst...)

	//! 设置豪华签到状态
	self.SignPlusStatus = SignPlus_Can_Not_Receive

	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivitySign) End(versionCode int32, resetCode int32) {
	self.IsSign = false
	self.IsSignPlus = false
	self.SignDay = 0
	self.SignPlusAward = []gamedata.ST_ItemData{}
	self.SignPlusStatus = 0
	self.VersionCode = versionCode

	go self.DB_Reset()
}

func (self *TActivitySign) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivitySign) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivitySign) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.IsSign == false {
		return true
	}

	return true
}

//! 刷新豪华签到奖励
func (self *TActivitySign) RefreshSignPlusAward(updateDB bool) {

	//! 获取奖励内容
	data := gamedata.GetSignPlusDataFromLevel(self.activityModule.ownplayer.GetLevel())
	if data == nil {
		gamelog.Error("GetSignPlusDataFromLevel nil")
		return
	}

	if len(self.SignPlusAward) != 0 {
		self.SignPlusAward = []gamedata.ST_ItemData{}
	}

	awardLst := gamedata.GetItemsFromAwardIDEx(data.SignAward)
	self.SignPlusAward = append(self.SignPlusAward, awardLst...)

	//! 设置豪华签到状态
	self.SignPlusStatus = SignPlus_Can_Not_Receive

	//! 更新到数据库
	if updateDB == true {
		self.DB_UpdateSignPlusInfoToDatabase()
	}
}

//! 普通签到
func (self *TActivitySign) Sign() (bool, int, int) {

	//! 签到天数加一
	self.SignDay += 1
	self.IsSign = true

	//! 更新到数据库
	go self.DB_UpdateSignInfoToDatabase()

	//! 发放奖励
	awardData := gamedata.GetSignData(self.SignDay)
	if awardData == nil {
		gamelog.Error("Sign error: GetSignData return nil. SignDay: %d", self.SignDay)
		return false, 0, 0
	}

	//! 获取用户VIP等级
	playerVip := self.activityModule.ownplayer.GetVipLevel()

	//! 判断签到奖励是否有VIP加成
	//! 有VIP加成,判断是否满足加成条件
	if playerVip >= awardData.VipLevel {
		//! 满足,加倍领取
		if awardData.Multiple == 0 {
			awardData.Multiple = 1
			self.activityModule.ownplayer.BagMoudle.AddAwardItem(awardData.AwardItem, awardData.Count*awardData.Multiple)
		}

		return true, awardData.AwardItem, awardData.Count * awardData.Multiple
	}

	//! 不满足则普通领取
	self.activityModule.ownplayer.BagMoudle.AddAwardItem(awardData.AwardItem, awardData.Count)
	return true, awardData.AwardItem, awardData.Count
}

//! 豪华签到
func (self *TActivitySign) SignPlus() []gamedata.ST_ItemData {

	//! 设置领取状态
	self.SignPlusStatus = SignPlus_Have_Received
	self.IsSignPlus = true

	//! 更新奖励标记到数据库
	ret := self.DB_UpdateSignPlusInfoToDatabase()
	if ret == false {
		gamelog.Error("DB_UpdateSignPlusInfoToDatabase error.")
		return []gamedata.ST_ItemData{}
	}

	//! 发放奖励
	self.activityModule.ownplayer.BagMoudle.AddAwardItems(self.SignPlusAward)

	return self.SignPlusAward
}

//! 设置豪华签到可领取状态
func (self *TActivitySign) SetSignPlusStatus() {
	if self.SignPlusStatus == SignPlus_Have_Received {
		//! 已领取则不改变状态
		return
	}
	self.SignPlusStatus = SignPlus_Can_Receive
}

func (self *TActivitySign) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"sign.activityid":     self.ActivityID,
		"sign.issign":         self.IsSign,
		"sign.issignplus":     self.IsSignPlus,
		"sign.signday":        self.SignDay,
		"sign.signplusaward":  self.SignPlusAward,
		"sign.signplusstatus": self.SignPlusStatus,
		"sign.versioncode":    self.VersionCode,
		"sign.resetcode":      self.ResetCode}})
	return true
}

//! 更新豪华签到信息到数据库
func (self *TActivitySign) DB_UpdateSignPlusInfoToDatabase() bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"sign.signplusaward":  self.SignPlusAward,
		"sign.issignplus":     self.IsSignPlus,
		"sign.signplusstatus": self.SignPlusStatus}})
}

//! 更新普通签到信息到数据库
func (self *TActivitySign) DB_UpdateSignInfoToDatabase() bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"sign.signday":     self.SignDay,
		"sign.issign":      self.IsSign,
		"sign.versioncode": self.VersionCode}})
}

func (self *TActivitySign) DB_Refresh() bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"sign.signplusaward":  self.SignPlusAward,
		"sign.issignplus":     self.IsSignPlus,
		"sign.signday":        self.SignDay,
		"sign.issign":         self.IsSign,
		"sign.versioncode":    self.VersionCode,
		"sign.signplusstatus": self.SignPlusStatus}})
}

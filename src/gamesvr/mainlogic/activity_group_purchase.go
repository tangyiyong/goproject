package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 记录物品花费,活动结束补齐差价
type TActivityPurchaseCost struct {
	ItemID   int
	MoneyNum int
	Times    int //! 购买次数
}

//! 团购
type TActivityGroupPurchase struct {
	ActivityID          int                     //! 活动ID
	PurchaseCostLst     []TActivityPurchaseCost //! 个人花费信息
	Score               int                     //! 积分
	ScoreAwardMark      IntLst                  //! 积分奖励领取标记
	IsDifferenceReceive bool                    //! 差价领取标记
	VersionCode         int                     //! 更新号
	ResetCode           int                     //! 迭代号
	activityModule      *TActivityModule        //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityGroupPurchase) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityGroupPurchase) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.PurchaseCostLst = []TActivityPurchaseCost{}
	self.Score = 0
	self.ScoreAwardMark = IntLst{}
	self.IsDifferenceReceive = false

	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityGroupPurchase) Refresh(versionCode int) {
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityGroupPurchase) End(versionCode int, resetCode int) {
	self.PurchaseCostLst = []TActivityPurchaseCost{}
	self.Score = 0
	self.ScoreAwardMark = IntLst{}
	self.IsDifferenceReceive = false
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	//! 存储数据库
	go self.DB_Reset()
}

func (self *TActivityGroupPurchase) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityGroupPurchase) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityGroupPurchase) RedTip() bool {
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	isEnd, _ := G_GlobalVariables.IsActivityTime(self.ActivityID)
	if isEnd == false && self.IsDifferenceReceive == false && len(self.PurchaseCostLst) != 0 {
		return true //! 活动结束有返还货币的情况则返回红点
	}

	return false
}

func (self *TActivityGroupPurchase) GetGroupItemInfo(itemID int) (*TActivityPurchaseCost, int) {
	length := len(self.PurchaseCostLst)
	for i := 0; i < length; i++ {
		if self.PurchaseCostLst[i].ItemID == itemID {
			return &self.PurchaseCostLst[i], i
		}
	}

	var newRecord TActivityPurchaseCost
	newRecord.ItemID = itemID
	newRecord.Times = 0
	newRecord.MoneyNum = 0
	self.PurchaseCostLst = append(self.PurchaseCostLst, newRecord)
	go self.DB_AddNewPurchaseCostInfo(&newRecord)

	return &self.PurchaseCostLst[length], length
}

func (self *TActivityGroupPurchase) DB_AddNewPurchaseCostInfo(newRecord *TActivityPurchaseCost) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID},
		"grouppurchase.purchasecostlst", *newRecord)
}

func (self *TActivityGroupPurchase) DB_UpdatePurchaseCostInfo(index int) {
	filedName := fmt.Sprintf("grouppurchase.%d.purchasecostlst.times", index)
	filedName2 := fmt.Sprintf("grouppurchase.%d.purchasecostlst.moneynum", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.PurchaseCostLst[index].Times,
		filedName2: self.PurchaseCostLst[index].MoneyNum}})
}

func (self *TActivityGroupPurchase) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"grouppurchase.versioncode": self.VersionCode}})
	return true
}

//! 存储数据库
func (self *TActivityGroupPurchase) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"grouppurchase.activityid":          self.ActivityID,
		"grouppurchase.purchasecostlst":     self.PurchaseCostLst,
		"grouppurchase.score":               self.Score,
		"grouppurchase.scoreawardmark":      self.ScoreAwardMark,
		"grouppurchase.isdifferencereceive": self.IsDifferenceReceive,
		"grouppurchase.versioncode":         self.VersionCode,
		"grouppurchase.resetcode":           self.ResetCode}})
	return true
}

func (self *TActivityGroupPurchase) DB_SaveIdfferenceMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"grouppurchase.isdifferencereceive": self.IsDifferenceReceive}})
}

func (self *TActivityGroupPurchase) DB_SaveScore() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"grouppurchase.score": self.Score}})
}

func (self *TActivityGroupPurchase) DB_AddScoreAward(id int) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, "grouppurchase.scoreawardmark", id)
}
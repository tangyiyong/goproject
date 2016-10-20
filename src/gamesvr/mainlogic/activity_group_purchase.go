package mainlogic

import (
	"fmt"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 记录物品花费,活动结束补齐差价
type TActivityPurchaseCost struct {
	ItemID   int
	MoneyNum int
	Times    int //! 购买次数
}

//! 记录团购购买物品次数
type TActivityPurchaseOrder struct {
	ItemID int
	Times  int //! 购买次数
}

//! 团购
type TActivityGroupPurchase struct {
	ActivityID      int32                    //! 活动ID
	PurchaseCostLst []TActivityPurchaseCost  //! 个人花费信息
	ShoppingInfo    []TActivityPurchaseOrder //! 购物信息
	Score           int                      //! 积分
	ScoreAwardMark  IntLst                   //! 积分奖励领取标记
	VersionCode     int32                    //! 更新号
	ResetCode       int32                    //! 迭代号
	activityModule  *TActivityModule         //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityGroupPurchase) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityGroupPurchase) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.PurchaseCostLst = []TActivityPurchaseCost{}
	self.Score = 0
	self.ScoreAwardMark = IntLst{}

	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityGroupPurchase) Refresh(versionCode int32) {
	self.VersionCode = versionCode

	length := len(self.ShoppingInfo)
	for i := 0; i < length; i++ {
		self.ShoppingInfo[i].Times = 0
	}

	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityGroupPurchase) End(versionCode int32, resetCode int32) {
	self.ShoppingInfo = []TActivityPurchaseOrder{}
	self.PurchaseCostLst = []TActivityPurchaseCost{}
	self.Score = 0
	self.ScoreAwardMark = IntLst{}
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	//! 存储数据库
	self.DB_Reset()
}

func (self *TActivityGroupPurchase) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityGroupPurchase) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityGroupPurchase) RedTip() bool {
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	beginIndex, endIndex := gamedata.GetGroupPurchaseScoreAwardSection(awardType)
	for i := beginIndex; i < endIndex; i++ {
		scoreInfo := gamedata.GetGroupPurchaseScoreAward(i)

		if self.Score >= scoreInfo.NeedScore && self.ScoreAwardMark.IsExist(i) < 0 {
			return true //! 有未领的积分奖励
		}
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
	self.DB_AddNewPurchaseCostInfo(&newRecord)

	return &self.PurchaseCostLst[length], length
}

func (self *TActivityGroupPurchase) GetGroupItemShoppingInfo(itemID int) (*TActivityPurchaseOrder, int) {
	length := len(self.ShoppingInfo)
	for i := 0; i < length; i++ {
		if self.ShoppingInfo[i].ItemID == itemID {
			return &self.ShoppingInfo[i], i
		}
	}

	var newRecord TActivityPurchaseOrder
	newRecord.ItemID = itemID
	newRecord.Times = 0
	self.ShoppingInfo = append(self.ShoppingInfo, newRecord)
	self.DB_AddNewPurchaseOrderInfo(&newRecord)

	return &self.ShoppingInfo[length], length
}

func (self *TActivityGroupPurchase) DB_AddNewPurchaseCostInfo(newRecord *TActivityPurchaseCost) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID},
		&bson.M{"$push": bson.M{"grouppurchase.purchasecostlst": *newRecord}})
}

func (self *TActivityGroupPurchase) DB_AddNewPurchaseOrderInfo(newRecord *TActivityPurchaseOrder) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID},
		&bson.M{"$push": bson.M{"grouppurchase.shoppinginfo": *newRecord}})
}

func (self *TActivityGroupPurchase) DB_UpdatePurchaseCostInfo(index int) {
	filedName := fmt.Sprintf("grouppurchase.purchasecostlst.%d.times", index)
	filedName2 := fmt.Sprintf("grouppurchase.purchasecostlst.%d.moneynum", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		filedName:  self.PurchaseCostLst[index].Times,
		filedName2: self.PurchaseCostLst[index].MoneyNum}})
}

func (self *TActivityGroupPurchase) DB_UpdatePurchaseOrderInfo(index int) {
	filedName := fmt.Sprintf("grouppurchase.shoppinginfo.%d.times", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.ShoppingInfo[index].Times}})
}

func (self *TActivityGroupPurchase) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"grouppurchase.versioncode":  self.VersionCode,
		"grouppurchase.shoppinginfo": self.ShoppingInfo}})
}

//! 存储数据库
func (self *TActivityGroupPurchase) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"grouppurchase.activityid":      self.ActivityID,
		"grouppurchase.purchasecostlst": self.PurchaseCostLst,
		"grouppurchase.score":           self.Score,
		"grouppurchase.scoreawardmark":  self.ScoreAwardMark,
		"grouppurchase.versioncode":     self.VersionCode,
		"grouppurchase.shoppinginfo":    self.ShoppingInfo,
		"grouppurchase.resetcode":       self.ResetCode}})
}

func (self *TActivityGroupPurchase) DB_SaveScore() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"grouppurchase.score": self.Score}})
}

func (self *TActivityGroupPurchase) DB_AddScoreAward(id int) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$push": bson.M{"grouppurchase.scoreawardmark": id}})
}

func (self *TActivityGroupPurchase) DB_UpdateScoreAward() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{"grouppurchase.scoreawardmark": self.ScoreAwardMark}})
}

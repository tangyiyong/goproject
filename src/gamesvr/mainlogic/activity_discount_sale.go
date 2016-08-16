package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

type TDiscountSaleGoodsInfo struct {
	Index int //! 索引
	Times int //! 剩余购买次数
}

//! 折扣贩售
type TActivityDiscountSale struct {
	ActivityID     int                      //! 活动ID
	ShopLst        []TDiscountSaleGoodsInfo //! 商品列表
	VersionCode    int                      //! 版本号
	ResetCode      int                      //! 迭代号
	activityModule *TActivityModule         //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityDiscountSale) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityDiscountSale) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityDiscountSale) Refresh(versionCode int) {
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityDiscountSale) End(versionCode int, resetCode int) {
	self.ShopLst = []TDiscountSaleGoodsInfo{}
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	go self.DB_Reset()
}

func (self *TActivityDiscountSale) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityDiscountSale) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityDiscountSale) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	return false //! 折扣贩售没有红点
}

//! 添加一个商品数据
func (self *TActivityDiscountSale) AddItem(info TDiscountSaleGoodsInfo, index int) *TDiscountSaleGoodsInfo {
	self.ShopLst = append(self.ShopLst, info)
	go self.DB_AddShoppingInfo(index, &info)
	return &self.ShopLst[len(self.ShopLst)-1]
}

func (self *TActivityDiscountSale) DB_Refresh() bool {
	index := -1
	for i, v := range self.activityModule.DiscountSale {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("DiscountSale DB_Refresh fail. ActivityID: %d", self.ActivityID)
		return false
	}

	filedName := fmt.Sprintf("discount.%d.versioncode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: self.VersionCode}})
	return true
}

func (self *TActivityDiscountSale) DB_Reset() bool {
	index := -1
	for i, v := range self.activityModule.DiscountSale {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("DiscountSale DB_Reset fail. ActivityID: %d", self.ActivityID)
		return false
	}

	filedName := fmt.Sprintf("discount.%d.shoplst", index)
	filedName2 := fmt.Sprintf("discount.%d.versioncode", index)
	filedName3 := fmt.Sprintf("discount.%d.resetcode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.ShopLst,
		filedName2: self.VersionCode,
		filedName3: self.ResetCode}})
	return true
}

func (self *TActivityDiscountSale) DB_AddShoppingInfo(activityIndex int, info *TDiscountSaleGoodsInfo) {
	filedName := fmt.Sprintf("discountsale.%d.shoplst", activityIndex)
	mongodb.AddToArray(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, filedName, *info)
}

func (self *TActivityDiscountSale) DB_UpdateShoppingTimes(activityIndex int, index int, info *TDiscountSaleGoodsInfo) {
	filedName := fmt.Sprintf("discountsale.%d.shoplst.%d.times", activityIndex, index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: (*info).Times}})
}

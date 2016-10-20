package mainlogic

import (
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
type TActivityDiscount struct {
	ActivityID  int32                    //! 活动ID
	ShopLst     []TDiscountSaleGoodsInfo //! 商品列表
	VersionCode int32                    //! 版本号
	ResetCode   int32                    //! 迭代号
	modulePtr   *TActivityModule         //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityDiscount) SetModulePtr(mPtr *TActivityModule) {
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityDiscount) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityDiscount) Refresh(versionCode int32) {
	self.VersionCode = versionCode
	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityDiscount) End(versionCode int32, resetCode int32) {
	self.ShopLst = []TDiscountSaleGoodsInfo{}
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	self.DB_Reset()
}

func (self *TActivityDiscount) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityDiscount) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityDiscount) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	return false //! 折扣贩售没有红点
}

//! 添加一个商品数据
func (self *TActivityDiscount) AddItem(info TDiscountSaleGoodsInfo, index int) *TDiscountSaleGoodsInfo {
	self.ShopLst = append(self.ShopLst, info)
	self.DB_AddShoppingInfo(index, &info)
	return &self.ShopLst[len(self.ShopLst)-1]
}

func (self *TActivityDiscount) DB_Refresh() {
	index := -1
	for i, v := range self.modulePtr.DiscountSale {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("DiscountSale DB_Refresh fail. ActivityID: %d", self.ActivityID)
		return
	}

	filedName := fmt.Sprintf("discount.%d.versioncode", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.VersionCode}})
}

func (self *TActivityDiscount) DB_Reset() {
	index := -1
	for i, v := range self.modulePtr.DiscountSale {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("DiscountSale DB_Reset fail. ActivityID: %d", self.ActivityID)
		return
	}

	filedName := fmt.Sprintf("discount.%d.shoplst", index)
	filedName2 := fmt.Sprintf("discount.%d.versioncode", index)
	filedName3 := fmt.Sprintf("discount.%d.resetcode", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		filedName:  self.ShopLst,
		filedName2: self.VersionCode,
		filedName3: self.ResetCode}})
}

func (self *TActivityDiscount) DB_AddShoppingInfo(activityIndex int, info *TDiscountSaleGoodsInfo) {
	filedName := fmt.Sprintf("discountsale.%d.shoplst", activityIndex)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$push": bson.M{filedName: *info}})
}

func (self *TActivityDiscount) DB_UpdateShoppingTimes(activityIndex int, index int, info *TDiscountSaleGoodsInfo) {
	filedName := fmt.Sprintf("discountsale.%d.shoplst.%d.times", activityIndex, index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		filedName: (*info).Times}})
}

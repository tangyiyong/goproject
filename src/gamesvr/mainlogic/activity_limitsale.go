package mainlogic

import (
	"fmt"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"mongodb"
	"time"
	"utility"
)

type TLimitSaleInfo struct {
	ID     int  //! 商品ID
	Status bool //! 是否购买 true -> 已购买  false -> 未购买
}

//! 限时优惠活动
type TActivityLimitSale struct {
	ActivityID       int32            //! 活动ID
	Score            int              //! 当前积分
	ItemLst          []TLimitSaleInfo //! 当天优惠物品
	DiscountChargeID int              //! 优惠充值ID
	RefreshMark      bool             //! 刷新标记
	AwardMark        BitsType         //! 全民奖励领取标记
	WeekReset        uint32           //! 全民奖励刷新周
	VersionCode      int32            //! 版本号
	ResetCode        int32            //! 迭代号
	modulePtr        *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityLimitSale) SetModulePtr(mPtr *TActivityModule) {
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityLimitSale) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.modulePtr = mPtr

	self.Score = 0
	self.ItemLst = []TLimitSaleInfo{}
	self.AwardMark = 0
	self.WeekReset = utility.GetCurDay()
	self.DiscountChargeID = 0

	self.modulePtr.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode

	self.RefreshItem()
	self.RefreshMark = true
}

//! 刷新数据
func (self *TActivityLimitSale) Refresh(versionCode int32) {
	//! 刷新贩售物品
	self.RefreshItem()
	self.RefreshMark = true
	self.VersionCode = versionCode

	//! 如果积分满100分, 则清空
	if self.Score >= 100 {
		self.Score = 0
		self.DiscountChargeID = 0
	}

	if utility.IsSameWeek(self.WeekReset) == false {
		//! 刷新全民奖励
		self.AwardMark = 0
		self.WeekReset = utility.GetCurDay()
		G_GlobalVariables.LimitSaleNum = 0
		G_GlobalVariables.DB_UpdateLimitSaleNum()
	}

	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityLimitSale) End(versionCode int32, resetCode int32) {

	self.ResetCode = resetCode
	self.VersionCode = versionCode

	self.DB_Reset()
}

func (self *TActivityLimitSale) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityLimitSale) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityLimitSale) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.Score >= 100 {
		//! 积分满时显示红点
		self.RefreshMark = false
		self.DB_SaveRefreshMark()
		return true
	}

	if self.RefreshMark == true {
		self.RefreshMark = false
		self.DB_SaveRefreshMark()
		return true
	}

	return false
}

func (self *TActivityLimitSale) RefreshItem() {
	if len(self.ItemLst) != 0 {
		self.ItemLst = []TLimitSaleInfo{}
	}

	itemIDLst := gamedata.RandLimitSaleItem()
	for i := 0; i < len(itemIDLst); i++ {
		var item TLimitSaleInfo
		item.ID = itemIDLst[i]
		item.Status = false
		self.ItemLst = append(self.ItemLst, item)
	}

}

func (self *TActivityLimitSale) GetDiscountCharge() int {
	begin, end := gamedata.GetDiscountChargeIDSection()

	if self.DiscountChargeID < begin ||
		self.DiscountChargeID > end ||
		G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return 0
	}

	return self.DiscountChargeID
}

func (self *TActivityLimitSale) RandDiscountCharge() {
	if self.DiscountChargeID != 0 {
		return
	}

	beginid, endid := gamedata.GetDiscountChargeIDSection()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	self.DiscountChargeID = r.Intn(endid-beginid+1) + beginid
	self.DB_UpdateDiscountCharge()
}

func (self *TActivityLimitSale) DiscountChargeClear() {
	self.DiscountChargeID = 0
	self.Score = 0
	self.DB_UpdateDiscountCharge()
}

func (self *TActivityLimitSale) DB_UpdateDiscountCharge() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"limitsale.score":            self.Score,
		"limitsale.discountchargeid": self.DiscountChargeID}})
}

func (self *TActivityLimitSale) DB_UpdateScore() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"limitsale.score": self.Score}})
}

func (self *TActivityLimitSale) DB_UpdateStatus(index int) {
	filedName := fmt.Sprintf("limitsale.itemlst.%d.status", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.ItemLst[index].Status}})
}

func (self *TActivityLimitSale) DB_SaveRefreshMark() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"limitsale.refreshmark": self.RefreshMark}})
}

func (self *TActivityLimitSale) DB_UpdateAwardMark() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"limitsale.awardmark": self.AwardMark}})
}

func (self *TActivityLimitSale) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"limitsale.versioncode":      self.VersionCode,
		"limitsale.refreshmark":      self.RefreshMark,
		"limitsale.score":            self.Score,
		"limitsale.discountchargeid": self.DiscountChargeID,
		"limitsale.awardmark":        self.AwardMark,
		"limitsale.weekreset":        self.WeekReset,
		"limitsale.itemlst":          self.ItemLst}})
}

func (self *TActivityLimitSale) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"limitsale.versioncode":      self.VersionCode,
		"limitsale.resetcode":        self.ResetCode,
		"limitsale.refreshmark":      self.RefreshMark,
		"limitsale.score":            self.Score,
		"limitsale.discountchargeid": self.DiscountChargeID,
		"limitsale.awardmark":        self.AwardMark,
		"limitsale.weekreset":        self.WeekReset,
		"limitsale.itemlst":          self.ItemLst}})
}

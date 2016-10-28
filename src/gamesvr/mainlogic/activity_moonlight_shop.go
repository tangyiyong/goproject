package mainlogic

import (
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"utility"
)

const (
	MoonShop_Num = 6
)

// 月光集市
type TMoonShop struct {
	TMoonlightShopData
	ActivityID  int32            //! 活动ID
	VersionCode int32            //! 版本号
	ResetCode   int32            //! 迭代号
	modulePtr   *TActivityModule //! 指针
}
type TMoonlightShopData struct {
	ExchangeTimes []byte
	Goods         [MoonShop_Num]TMoonlightGoods
	RefreshTime   int32
	Score         int
	BuyTimes      int
	ScoreAward    BitsType64
}

type TMoonlightGoods struct {
	ID            int
	BuyTimes      byte
	Discount      byte // 百分比折扣
	DiscountTimes byte
}

//！ 活动框架代码
func (self *TMoonShop) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.modulePtr = mPtr
	self.VersionCode = vercode
	self.ResetCode = resetcode
	self.modulePtr.activityPtrs[self.ActivityID] = self

	self.Goods = [MoonShop_Num]TMoonlightGoods{}
	self.ExchangeTimes = make([]byte, len(gamedata.G_MoonShopExchg_List))
}
func (self *TMoonShop) SetModulePtr(mPtr *TActivityModule) {
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
}
func (self *TMoonShop) Refresh(versionCode int32) {
	self.VersionCode = versionCode

	self.ExchangeTimes = make([]byte, len(gamedata.G_MoonShopExchg_List))

	self.BuyTimes = 0

	self.DB_Refresh()
}
func (self *TMoonShop) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.ExchangeTimes = make([]byte, len(gamedata.G_MoonShopExchg_List))

	self.Goods = [MoonShop_Num]TMoonlightGoods{}
	self.RefreshTime = 0

	self.ScoreAward = 0

	self.DB_Reset()

}
func (self *TMoonShop) GetRefreshV() int32 {
	return self.VersionCode
}
func (self *TMoonShop) GetResetV() int32 {
	return self.ResetCode
}
func (self *TMoonShop) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.CanGetScoreAward() {
		return true
	}

	return false
}

//! 数据操作代码
func (self *TMoonlightShopData) GetShopDtad() *TMoonlightShopData {
	return self
}

//! DB相关
func (self *TMoonShop) DB_SaveExchangeTimes(nIndex int) {
	FieldName := fmt.Sprintf("moonshop.exchangetimes.%d", nIndex)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.ExchangeTimes[nIndex]}})
}
func (self *TMoonShop) DB_SaveAllExchangeTimes() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{"moonshop.exchangetimes": self.ExchangeTimes}})
}
func (self *TMoonShop) DB_SaveGoods(nIndex int) {
	FieldName := fmt.Sprintf("moonshop.goods.%d", nIndex)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.Goods[nIndex]}})
}
func (self *TMoonShop) DB_SaveAllGoods() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"moonshop.goods":       self.Goods,
		"moonshop.refreshtime": self.RefreshTime}})
}
func (self *TMoonShop) DB_Save_Score_Buytimes() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"moonshop.score":    self.Score,
		"moonshop.buytimes": self.BuyTimes}})
}
func (self *TMoonShop) DB_SaveScoreAwardFlag() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{"moonshop.scoreaward": self.ScoreAward}})
}
func (self *TMoonShop) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"moonshop.score":         self.Score,
		"moonshop.buytimes":      self.BuyTimes,
		"moonshop.exchangetimes": self.ExchangeTimes,
		"moonshop.activityid":    self.ActivityID,
		"moonshop.versioncode":   self.VersionCode,
		"moonshop.resetcode":     self.ResetCode}})
}
func (self *TMoonShop) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"moonshop.goods":         self.Goods,
		"moonshop.refreshtime":   self.RefreshTime,
		"moonshop.scoreaward":    self.ScoreAward,
		"moonshop.exchangetimes": self.ExchangeTimes,
		"moonshop.activityid":    self.ActivityID,
		"moonshop.versioncode":   self.VersionCode,
		"moonshop.resetcode":     self.ResetCode}})
}

//！ 逻辑代码
// 兑换月光币
func (self *TMoonShop) ExchangeMoonMoney(player *TPlayer, id int) bool {
	pExchangeInfo := gamedata.GetMoonShopExchgInfo(id)
	if pExchangeInfo == nil {
		return false
	}
	if self.ExchangeTimes[id] < pExchangeInfo.DailyTimes {
		if player.RoleMoudle.CostMoney(pExchangeInfo.CostMoneyID, pExchangeInfo.CostNum) {
			player.BagMoudle.AddNormalItem(pExchangeInfo.ItemID, pExchangeInfo.ItemNum)
			self.ExchangeTimes[id]++
			self.DB_SaveExchangeTimes(id)
			return true
		}
	}
	return false
}

//商品打折
func (self *TMoonShop) ReduceDiscount(player *TPlayer, goodsID int) (bool, byte) {
	for i := 0; i < MoonShop_Num; i++ {
		goods := &self.Goods[i]
		if goods.ID == goodsID {
			csv := gamedata.GetMoonGoodsInfo(goodsID)
			if goods.Discount > csv.MinDiscount && player.BagMoudle.RemoveNormalItem(gamedata.MoonShop_Money_ID, goods.getDiscountCost()) {
				goods.reduceDiscount(csv.MinDiscount)
				self.DB_SaveGoods(i)
				return true, goods.Discount
			}
			return false, goods.Discount
		}
	}
	return false, 100
}
func (self *TMoonlightGoods) getDiscountCost() int {
	if int(self.DiscountTimes) >= len(gamedata.MoonlightShop_Discount_Cost) {
		return 0
	}
	return int(gamedata.MoonlightShop_Discount_Cost[self.DiscountTimes])
}
func (self *TMoonlightGoods) reduceDiscount(min byte) {
	left := int(gamedata.MoonlightShop_Discount_OneTiems[0])
	right := int(gamedata.MoonlightShop_Discount_OneTiems[1])
	var reduce byte = byte(utility.RandBetween(left, right))
	self.Discount -= reduce
	if self.Discount < min {
		self.Discount = min
	}
	self.DiscountTimes++
}

//重新生成一批商品
func (self *TMoonShop) RefreshShop_Auto(player *TPlayer) bool {
	now := utility.GetCurTime()
	if now-self.RefreshTime >= int32(gamedata.MoonlightShop_Shop_Refresh_CD*60) {
		self.RefreshTime = now
		self.refreshShop()
		return true
	}
	return false
}
func (self *TMoonShop) RefreshShop_Buy(player *TPlayer) bool {
	cost := int(gamedata.MoonlightShop_Shop_Refresh_Cost)
	if player.BagMoudle.RemoveNormalItem(gamedata.MoonShop_Money_ID, cost) {
		self.Score += cost
		self.DB_Save_Score_Buytimes()
		self.refreshShop()
		return true
	}
	return false
}
func (self *TMoonShop) refreshShop() {
	IDList := gamedata.RandSelect_MoonlightGoods(self.ActivityID, MoonShop_Num)
	for i := 0; i < MoonShop_Num; i++ {
		goods := &self.Goods[i]
		goods.ID = IDList[i]
		goods.BuyTimes = 0
		goods.DiscountTimes = 0
		csv := gamedata.GetMoonGoodsInfo(goods.ID)
		goods.Discount = byte(utility.RandBetween(int(csv.MinDiscount), int(csv.MaxDiscount)))
	}
	self.DB_SaveAllGoods()
}

func (self *TMoonShop) BuyGoods(player *TPlayer, goodsID int) bool {
	if self.BuyTimes >= gamedata.MoonlightShop_BuyTimes_Max {
		return false
	}
	for i := 0; i < MoonShop_Num; i++ {
		goods := &self.Goods[i]
		if goods.ID == goodsID {
			csv := gamedata.GetMoonGoodsInfo(goodsID) // 经过上面的判断，此指针不会为nil
			if goods.BuyTimes >= csv.DailyTimes {
				return false
			}

			cost := csv.Price * int(goods.Discount) / 100
			if player.BagMoudle.RemoveNormalItem(gamedata.MoonShop_Money_ID, cost) {
				if player.BagMoudle.AddAwardItem(csv.ItemID, csv.ItemNum) {
					self.BuyTimes++
					goods.BuyTimes++
					self.Score += cost
					self.DB_Save_Score_Buytimes()
					self.DB_SaveGoods(i)
				}
			}
			gamelog.Error("moonshop BuyGoods Error: ItemID:%d ItemNum:%d cost:%d", csv.ItemID, csv.ItemNum, cost)
			return false
		}
	}
	return false
}
func (self *TMoonShop) GetScoreAward(player *TPlayer, awardID int) bool {
	pAwardInfo := gamedata.GetMoonShopAwardInfo(awardID)
	if pAwardInfo == nil || self.ScoreAward.Get(awardID) {
		return false
	}

	if self.Score >= pAwardInfo.NeedScore {
		if player.BagMoudle.AddAwardItem(pAwardInfo.ItemID, pAwardInfo.ItemNum) {
			self.ScoreAward.Set(awardID)
			self.DB_SaveScoreAwardFlag()
			return true
		}
	}
	return false
}
func (self *TMoonShop) CanGetScoreAward() bool {
	for i := 1; i < len(gamedata.G_MoonAward_List); i++ {
		csv := &gamedata.G_MoonAward_List[i]
		if !self.ScoreAward.Get(csv.ID) && self.Score >= csv.NeedScore {
			return true
		}
	}
	return false
}

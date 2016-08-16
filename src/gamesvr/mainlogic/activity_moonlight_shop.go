package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

const (
	MoonlightShop_Goods_Num = 6
)

// 月光集市
type TMoonlightShop struct {
	TMoonlightShopData

	ActivityID     int              //! 活动ID
	VersionCode    int              //! 版本号
	ResetCode      int              //! 迭代号
	activityModule *TActivityModule //! 指针
}
type TMoonlightShopData struct {
	ExchangeTimes   []byte
	Goods           [MoonlightShop_Goods_Num]TMoonlightGoods
	AutoRefreshTime int64
	Score           int
	BuyTimes        int
	ScoreAwardFlag  int64
}
type TMoonlightGoods struct {
	ID            int
	BuyTimes      byte
	Discount      byte // 百分比折扣
	DiscountTimes byte
}

//！ 活动框架代码
func (self *TMoonlightShop) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.VersionCode = vercode
	self.ResetCode = resetcode
	self.activityModule.activityPtrs[self.ActivityID] = self

	self.Goods = [MoonlightShop_Goods_Num]TMoonlightGoods{}
	self.ExchangeTimes = make([]byte, len(gamedata.G_MoonlightShopExchangeCsv))
}
func (self *TMoonlightShop) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}
func (self *TMoonlightShop) Refresh(versionCode int) {
	self.VersionCode = versionCode

	self.ExchangeTimes = make([]byte, len(gamedata.G_MoonlightShopExchangeCsv))

	self.BuyTimes = 0

	self.DB_Refresh()
}
func (self *TMoonlightShop) End(versionCode int, resetCode int) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.ExchangeTimes = make([]byte, len(gamedata.G_MoonlightShopExchangeCsv))

	self.Goods = [MoonlightShop_Goods_Num]TMoonlightGoods{}
	self.AutoRefreshTime = 0

	self.ScoreAwardFlag = 0

	self.DB_Reset()

}
func (self *TMoonlightShop) GetRefreshV() int {
	return self.VersionCode
}
func (self *TMoonlightShop) GetResetV() int {
	return self.ResetCode
}
func (self *TMoonlightShop) RedTip() bool {
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
func (self *TMoonlightShop) db_SaveExchangeTimes(nIndex int) {
	FieldName := fmt.Sprintf("moonlightshop.exchangetimes.%d", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{FieldName: self.ExchangeTimes[nIndex]}})
}
func (self *TMoonlightShop) db_SaveAllExchangeTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{"moonlightshop.exchangetimes": self.ExchangeTimes}})
}
func (self *TMoonlightShop) db_SaveGoods(nIndex int) {
	FieldName := fmt.Sprintf("moonlightshop.goods.%d", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{FieldName: self.Goods[nIndex]}})
}
func (self *TMoonlightShop) db_SaveAllGoods() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"moonlightshop.goods":           self.Goods,
		"moonlightshop.autorefreshtime": self.AutoRefreshTime}})
}
func (self *TMoonlightShop) db_Save_Score_Buytimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"moonlightshop.score":    self.Score,
		"moonlightshop.buytimes": self.BuyTimes}})
}
func (self *TMoonlightShop) db_SaveScoreAwardFlag() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{"moonlightshop.scoreawardflag": self.ScoreAwardFlag}})
}
func (self *TMoonlightShop) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"moonlightshop.score":         self.Score,
		"moonlightshop.buytimes":      self.BuyTimes,
		"moonlightshop.exchangetimes": self.ExchangeTimes,
		"moonlightshop.activityid":    self.ActivityID,
		"moonlightshop.versioncode":   self.VersionCode,
		"moonlightshop.resetcode":     self.ResetCode}})
	return true
}
func (self *TMoonlightShop) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"moonlightshop.goods":           self.Goods,
		"moonlightshop.autorefreshtime": self.AutoRefreshTime,
		"moonlightshop.scoreawardflag":  self.ScoreAwardFlag,
		"moonlightshop.exchangetimes":   self.ExchangeTimes,
		"moonlightshop.activityid":      self.ActivityID,
		"moonlightshop.versioncode":     self.VersionCode,
		"moonlightshop.resetcode":       self.ResetCode}})
	return true
}

//！ 逻辑代码
// 兑换月光币
func (self *TMoonlightShop) ExchangeToken(player *TPlayer, id int) bool {
	csv := gamedata.GetMoonlightShopExchangeCsv(id)
	if csv == nil {
		return false
	}
	if self.ExchangeTimes[id] < csv.DailyTimes {
		if player.RoleMoudle.CostMoney(csv.CostType, csv.CostNum) {
			player.BagMoudle.AddNormalItem(gamedata.MoonlightShop_Token_ItemID, csv.GetToken)
			self.ExchangeTimes[id]++
			self.db_SaveExchangeTimes(id)
			return true
		}
	}
	return false
}

//商品打折
func (self *TMoonlightShop) ReduceDiscount(player *TPlayer, goodsID int) (bool, byte) {
	for i := 0; i < MoonlightShop_Goods_Num; i++ {
		goods := &self.Goods[i]
		if goods.ID == goodsID {
			csv := gamedata.GetMoonlightGoodsCsv(goodsID) // 经过上面的判断，此指针不会为nil
			if goods.Discount > csv.MinDiscount && player.BagMoudle.RemoveNormalItem(gamedata.MoonlightShop_Token_ItemID, goods.getDiscountCost()) {
				goods.reduceDiscount(csv.MinDiscount)
				self.db_SaveGoods(i)
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
func (self *TMoonlightShop) RefreshShop_Auto(player *TPlayer) bool {
	now := time.Now().Unix()
	if now-self.AutoRefreshTime >= int64(gamedata.MoonlightShop_Shop_Refresh_CD*60) {
		self.AutoRefreshTime = now
		self.refreshShop()
		return true
	}
	return false
}
func (self *TMoonlightShop) RefreshShop_Buy(player *TPlayer) bool {
	cost := int(gamedata.MoonlightShop_Shop_Refresh_Cost)
	if player.BagMoudle.RemoveNormalItem(gamedata.MoonlightShop_Token_ItemID, cost) {
		self.Score += cost
		self.db_Save_Score_Buytimes()
		self.refreshShop()
		return true
	}
	return false
}
func (self *TMoonlightShop) refreshShop() {
	IDList := gamedata.RandSelect_MoonlightGoods(self.ActivityID, MoonlightShop_Goods_Num)
	for i := 0; i < MoonlightShop_Goods_Num; i++ {
		goods := &self.Goods[i]
		goods.ID = IDList[i]
		goods.BuyTimes = 0
		goods.DiscountTimes = 0
		csv := gamedata.GetMoonlightGoodsCsv(goods.ID)
		goods.Discount = byte(utility.RandBetween(int(csv.MinDiscount), int(csv.MaxDiscount)))
	}
	self.db_SaveAllGoods()
}

func (self *TMoonlightShop) BuyGoods(player *TPlayer, goodsID int) bool {
	if self.BuyTimes >= gamedata.MoonlightShop_BuyTimes_Max {
		return false
	}
	for i := 0; i < MoonlightShop_Goods_Num; i++ {
		goods := &self.Goods[i]
		if goods.ID == goodsID {
			csv := gamedata.GetMoonlightGoodsCsv(goodsID) // 经过上面的判断，此指针不会为nil
			if goods.BuyTimes >= csv.DailyTimes {
				return false
			}

			cost := csv.Price * int(goods.Discount) / 100
			if player.BagMoudle.RemoveNormalItem(gamedata.MoonlightShop_Token_ItemID, cost) {
				if player.BagMoudle.AddAwardItem(csv.ItemID, csv.ItemNum) {
					self.BuyTimes++
					goods.BuyTimes++
					self.Score += cost
					self.db_Save_Score_Buytimes()
					self.db_SaveGoods(i)
				}
			}
			gamelog.Error("MoonlightShop BuyGoods Error: ItemID:%d ItemNum:%d cost:%d", csv.ItemID, csv.ItemNum, cost)
			return false
		}
	}
	return false
}
func (self *TMoonlightShop) GetScoreAward(player *TPlayer, awardID int) bool {
	csv := gamedata.GetMoonlightShopAwardCsv(awardID)
	if csv == nil || self.scoreAwardFlag(awardID) {
		return false
	}
	if self.Score >= csv.NeedScore {
		if player.BagMoudle.AddAwardItem(csv.ItemID, csv.ItemNum) {
			self.setScoreAwardFlag(awardID, true)
			self.db_SaveScoreAwardFlag()
			return true
		}
	}
	return false
}
func (self *TMoonlightShop) CanGetScoreAward() bool {
	for i := 1; i < len(gamedata.G_MoonlightAwardCsv); i++ {
		csv := &gamedata.G_MoonlightAwardCsv[i]
		if !self.scoreAwardFlag(csv.ID) && self.Score >= csv.NeedScore {
			return true
		}
	}
	return false
}
func (self *TMoonlightShop) scoreAwardFlag(awardID int) bool {
	var num uint = uint(awardID)
	return self.ScoreAwardFlag&(1<<num) > 0
}
func (self *TMoonlightShop) setScoreAwardFlag(awardID int, flag bool) {
	var num uint = uint(awardID)
	if flag {
		self.ScoreAwardFlag |= (1 << num)
	} else {
		self.ScoreAwardFlag &= ^(1 << num)
	}
}

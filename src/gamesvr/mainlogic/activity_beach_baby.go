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
	BeachBaby_Goods_Num = 16
)

type TBeachBabyInfo struct {
	TBeachBabyGoodsData

	ActivityID     int              //! 活动ID
	VersionCode    int              //! 版本号
	ResetCode      int              //! 迭代号
	activityModule *TActivityModule //! 指针

	selectGoodsIDs []int
}

type TBeachBabyGoodsData struct {
	Goods               [BeachBaby_Goods_Num]TBeachBabyGoods
	AutoRefreshTime     int64
	Score               [2]int
	TotalScore          int
	IsGetTodayRankAward bool // 今日排行奖励
	IsGetTotalRankAward bool // 累计排行奖励
	FreeConchBit        int8 // 特定时间免费领取贝壳
}

type TBeachBabyGoods struct {
	ID      int
	IsOpen  bool
	IsValid bool
}

//！ 活动框架代码
func (self *TBeachBabyInfo) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.VersionCode = vercode
	self.ResetCode = resetcode
	mPtr.activityPtrs[self.ActivityID] = self
}
func (self *TBeachBabyInfo) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
}
func (self *TBeachBabyInfo) Refresh(versionCode int) {
	self.VersionCode = versionCode

	self.FreeConchBit = 0

	self.IsGetTodayRankAward = false
	self.SetTodayScore(0)

	self.DB_Refresh()
}
func (self *TBeachBabyInfo) End(versionCode int, resetCode int) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.FreeConchBit = 0

	self.TotalScore = 0
	self.SetYesterdayScore(0)
	self.SetTodayScore(0)

	self.IsGetTodayRankAward = false
	self.IsGetTotalRankAward = false

	self.DB_Refresh()
}
func (self *TBeachBabyInfo) GetRefreshV() int {
	return self.VersionCode
}
func (self *TBeachBabyInfo) GetResetV() int {
	return self.ResetCode
}
func (self *TBeachBabyInfo) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	//! 检查排行榜是否有名次
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(self.ActivityID)
	if isHandleTime {
		//! 检查昨日排行榜
		for _, v := range G_BeachBabyYesterdayRanker.List {
			if v.RankID == self.activityModule.PlayerID && self.IsGetTodayRankAward == false {
				return true
			}
		}
	} else {
		//! 检查总排行榜
		for _, v := range G_BeachBabyTotalRanker.List {
			if v.RankID == self.activityModule.PlayerID && self.IsGetTotalRankAward == false {
				return true
			}
		}
	}

	//! 领取免费贝壳
	if self.CanGetFreeConch() {
		return true
	}

	return false
}

//！数据操作代码
func (self *TBeachBabyGoodsData) GetBeachBabyDtad() *TBeachBabyGoodsData {
	return self
}

func (self *TBeachBabyInfo) AddScore(num int) {
	newScore := self.GetTodayScore() + num
	self.SetTodayScore(newScore)
	self.TotalScore += num
	self.db_SaveScore()

	G_BeachBabyTodayRanker.SetRankItem(self.activityModule.PlayerID, newScore)
	G_BeachBabyTotalRanker.SetRankItem(self.activityModule.PlayerID, self.TotalScore)
}
func (self *TBeachBabyInfo) GetTodayScore() int {
	return self.Score[utility.GetCurDayMod()]
}
func (self *TBeachBabyInfo) SetTodayScore(num int) {
	self.Score[utility.GetCurDayMod()] = num
}
func (self *TBeachBabyInfo) GetYesterdayScore() int {
	if utility.GetCurDayMod() == 1 {
		return self.Score[0]
	} else {
		return self.Score[1]
	}
}
func (self *TBeachBabyInfo) SetYesterdayScore(num int) {
	if utility.GetCurDayMod() == 1 {
		self.Score[0] = num
	} else {
		self.Score[1] = num
	}
}
func (self *TBeachBabyInfo) GetTotalScore() int {
	return self.TotalScore
}

//! DB相关
func (self *TBeachBabyInfo) db_SaveGoods(nIndex int) {
	FieldName := fmt.Sprintf("beachbaby.goods.%d", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{FieldName: self.Goods[nIndex]}})
}
func (self *TBeachBabyInfo) db_SaveAllGoods() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"beachbaby.goods":           self.Goods,
		"beachbaby.autorefreshtime": self.AutoRefreshTime}})
}
func (self *TBeachBabyInfo) db_SaveScore() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"beachbaby.score":      self.Score,
		"beachbaby.totalscore": self.TotalScore}})
}
func (self *TBeachBabyInfo) db_SaveRankAwardFlag() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"beachbaby.isgettodayrankaward": self.IsGetTodayRankAward,
		"beachbaby.isgettotalrankaward": self.IsGetTotalRankAward}})
}
func (self *TBeachBabyInfo) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"beachbaby.score":               self.Score,
		"beachbaby.totalscore":          self.TotalScore,
		"beachbaby.freeconchbit":        self.FreeConchBit,
		"beachbaby.isgettodayrankaward": self.IsGetTodayRankAward,
		"beachbaby.isgettotalrankaward": self.IsGetTotalRankAward,
		"beachbaby.activityid":          self.ActivityID,
		"beachbaby.versioncode":         self.VersionCode,
		"beachbaby.resetcode":           self.ResetCode}})
	return true
}

func (self *TBeachBabyInfo) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"beachbaby.score":               self.Score,
		"beachbaby.totalscore":          self.TotalScore,
		"beachbaby.freeconchbit":        self.FreeConchBit,
		"beachbaby.isgettodayrankaward": self.IsGetTodayRankAward,
		"beachbaby.isgettotalrankaward": self.IsGetTotalRankAward,
		"beachbaby.activityid":          self.ActivityID,
		"beachbaby.versioncode":         self.VersionCode,
		"beachbaby.resetcode":           self.ResetCode}})
	return true
}

//！ 逻辑代码
// 翻开贝壳
func (self *TBeachBabyInfo) OpenGoods(idx int) (ret gamedata.ST_ItemData, bGetItem bool) {
	if idx < 0 || idx >= len(self.Goods) {
		return ret, false
	}
	if self.Goods[idx].IsOpen {
		return ret, false
	}

	csv, openCnt := gamedata.GetBeachBabyGoodsCsv(self.Goods[idx].ID), self.getOpenGoodsCnt()
	if csv == nil {
		return ret, false
	}
	if openCnt >= BeachBaby_Goods_Num {
		gamelog.Warn("BeachBaby OpenGoods Error: goods all open!!!")
		return ret, false
	}

	price := int(gamedata.BeachBaby_OpenGoods_Cost[openCnt])

	if self.costOpenGoods(price) {
		self.Goods[idx].IsOpen = true

		ret.ItemID = csv.ItemID
		ret.ItemNum = csv.ItemNum

		otherIdx := self.getOtherSameOpenGoodsIdx(idx)
		if otherIdx >= 0 {
			self.Goods[idx].IsValid = false
			self.Goods[otherIdx].IsValid = false
			self.db_SaveGoods(otherIdx)
			self.activityModule.ownplayer.BagMoudle.AddAwardItem(csv.ItemID, csv.ItemNum)
			return ret, true
		}

		self.AddScore(price)
		self.db_SaveGoods(idx)

		self.checkRefreshGoods()

		return ret, false
	}

	return ret, false
}
func (self *TBeachBabyInfo) costOpenGoods(cost int) bool {
	if self.activityModule.ownplayer.BagMoudle.RemoveNormalItem(gamedata.BeachBaby_Token_ItemID, cost) {
		return true
	}
	if self.activityModule.ownplayer.RoleMoudle.CostMoney(gamedata.BeachBaby_CostMoneyType, cost) {
		return true
	}

	gamelog.Error("BeachBaby  costOpenGoods Error : Conch:%d, Diamond:%d, Cost:%d",
		self.activityModule.ownplayer.BagMoudle.GetNormalItemCount(gamedata.BeachBaby_Token_ItemID),
		self.activityModule.ownplayer.RoleMoudle.GetMoney(gamedata.BeachBaby_CostMoneyType),
		cost)
	return false
}
func (self *TBeachBabyInfo) getOpenGoodsCnt() int {
	openCnt := 0
	for i := 0; i < BeachBaby_Goods_Num; i++ {
		if self.Goods[i].IsOpen {
			openCnt++
		}
	}
	return openCnt
}
func (self *TBeachBabyInfo) getOtherSameOpenGoodsIdx(idx int) int {
	if !self.Goods[idx].IsOpen {
		return -1
	}
	goalID := self.Goods[idx].ID
	for i, v := range self.Goods {
		if i != idx && v.ID == goalID && v.IsOpen && v.IsValid {
			return i
		}
	}
	return -1
}

// 一键全开
func (self *TBeachBabyInfo) OpenAllGoods() (ret []gamedata.ST_ItemData, bSuccess bool) {
	openCnt, cost := self.getOpenGoodsCnt(), 0
	if openCnt >= BeachBaby_Goods_Num {
		gamelog.Warn("BeachBaby OpenAllGoods: goods all open!!!")
		return ret, false
	}
	for i := openCnt; i < BeachBaby_Goods_Num; i++ {
		cost += int(gamedata.BeachBaby_OpenGoods_Cost[i])
	}

	token := self.activityModule.ownplayer.BagMoudle.GetNormalItemCount(gamedata.BeachBaby_Token_ItemID)
	diamond := self.activityModule.ownplayer.RoleMoudle.GetMoney(gamedata.BeachBaby_CostMoneyType)
	if token+diamond < cost {
		gamelog.Warn("BeachBaby OpenAllGoods: money not enough!!! token:%d, diamond:%d, cost:%d", token, diamond, cost)
		return ret, false
	}

	// 扣贝壳/钱
	if token >= cost {
		bSuccess = self.activityModule.ownplayer.BagMoudle.RemoveNormalItem(gamedata.BeachBaby_Token_ItemID, cost)
	} else {
		bSuccess = self.activityModule.ownplayer.BagMoudle.RemoveNormalItem(gamedata.BeachBaby_Token_ItemID, token) &&
			self.activityModule.ownplayer.RoleMoudle.CostMoney(gamedata.BeachBaby_CostMoneyType, cost-token)
	}
	if !bSuccess {
		gamelog.Error("BeachBaby OpenAllGoods CostMoney Error!!! token:%d, diamond:%d, cost:%d", token, diamond, cost)
		return ret, false
	}

	self.AddScore(cost)

	// 给物品
	openGoods := make(map[int]int)
	for i := 0; i < BeachBaby_Goods_Num; i++ {
		goods := &self.Goods[i]
		if !goods.IsOpen {
			cnt, ok := openGoods[goods.ID]
			if ok {
				openGoods[goods.ID] = cnt + 1
			} else {
				openGoods[goods.ID] = 1
			}
		}
	}
	for goodsID, goodsCnt := range openGoods {
		cnt := (goodsCnt + 1) / 2
		csv := gamedata.GetBeachBabyGoodsCsv(goodsID)
		self.activityModule.ownplayer.BagMoudle.AddAwardItem(csv.ItemID, csv.ItemNum*cnt)
		ret = append(ret, gamedata.ST_ItemData{csv.ItemID, csv.ItemNum * cnt})
	}

	// 刷新
	self.refreshGoods()

	return ret, true
}

// 免费/钻石刷新
func (self *TBeachBabyInfo) Refresh_Auto() bool {
	now := time.Now().Unix()
	if now-self.AutoRefreshTime >= int64(gamedata.BeachBaby_Refresh_CD*60) {
		self.AutoRefreshTime = now
		self.refreshGoods()
		return true
	}
	return false
}
func (self *TBeachBabyInfo) Refresh_Buy() bool {
	cost := int(gamedata.BeachBaby_Refresh_Cost)
	if self.activityModule.ownplayer.RoleMoudle.CostMoney(gamedata.BeachBaby_CostMoneyType, cost) {
		self.AddScore(cost)
		self.refreshGoods()
		return true
	}
	return false
}
func (self *TBeachBabyInfo) refreshGoods() {
	IDList := self.getNewGoodsIDList()
	for i := 0; i < BeachBaby_Goods_Num; i++ {
		goods := &self.Goods[i]
		goods.ID = 0
		goods.IsOpen = true
		goods.IsValid = false
	}
	for i := 0; i < len(IDList); i++ {
		goods := &self.Goods[i]
		goods.ID = IDList[i]
		goods.IsOpen = false
		goods.IsValid = true
	}
	self.db_SaveAllGoods()
}
func (self *TBeachBabyInfo) checkRefreshGoods() {
	if self.getOpenGoodsCnt() >= BeachBaby_Goods_Num {
		self.refreshGoods()
	}
}

// 领取免费贝壳
func (self *TBeachBabyInfo) CanGetFreeConch() bool {
	hour := time.Now().Hour()
	for i, v := range gamedata.BeachBaby_GetFreeToken_Time {
		if v == byte(hour) && self.getFreeConchBit(i) == false {
			return true
		}
	}
	return false
}
func (self *TBeachBabyInfo) GetFreeConch() bool {
	if self.CanGetFreeConch() {
		itemID, itemCnt := gamedata.BeachBaby_Token_ItemID, int(gamedata.BeachBaby_Token_ItemID)
		return self.activityModule.ownplayer.BagMoudle.AddNormalItem(itemID, itemCnt) > 0
	}
	return false
}
func (self *TBeachBabyInfo) getFreeConchBit(idx int) bool {
	var num uint = uint(idx)
	return self.FreeConchBit&(1<<num) > 0
}
func (self *TBeachBabyInfo) setFreeConchBit(idx int, flag bool) {
	var num uint = uint(idx)
	if flag {
		self.FreeConchBit |= (1 << num)
	} else {
		self.FreeConchBit &= ^(1 << num)
	}
}

// 自选必被刷出的商品
func (self *TBeachBabyInfo) SelectGoodsID(ids []int) bool {
	if len(self.selectGoodsIDs)+len(ids) > len(gamedata.BeachBaby_SelectGoods_Cost) {
		gamelog.Warn("BeachBaby SelectGoodsID: select goods is up to limit! selectCnt:%d, limit:%d",
			len(self.selectGoodsIDs)+len(ids), len(gamedata.BeachBaby_SelectGoods_Cost))
		return false
	}
	for _, v := range self.selectGoodsIDs {
		for _, id := range ids {
			if v == id {
				gamelog.Warn("BeachBaby SelectGoodsID: ID:(%d) is selected already", id)
				return false
			}
		}
	}
	self.selectGoodsIDs = append(self.selectGoodsIDs, ids...)
	return true
}
func (self *TBeachBabyInfo) getNewGoodsIDList() []int {
	IDList := gamedata.RandSelect_BeachBabyGoods(self.ActivityID, BeachBaby_Goods_Num/2)

	if len(IDList) != BeachBaby_Goods_Num/2 {
		gamelog.Error("RandSelect_BeachBabyGoods len:%d", len(IDList))
	}

	cost := 0
	for i := 0; i < len(self.selectGoodsIDs); i++ {
		cost += int(gamedata.BeachBaby_SelectGoods_Cost[i])
	}
	if cost > 0 && self.activityModule.ownplayer.RoleMoudle.CostMoney(gamedata.BeachBaby_CostMoneyType, cost) {

		self.AddScore(cost)

		// 将自选商品替换进随机结果中
		for j, v := range self.selectGoodsIDs {
			isExist := false
			for i := 0; i < len(IDList); i++ {
				if IDList[i] == v {
					isExist = true
					break
				}
			}
			if isExist == false {
				IDList[j] = v
			}
		}
	}
	IDList = append(IDList, IDList...) // 产生双份
	utility.RandShuffle(IDList)
	return IDList
}

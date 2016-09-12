package mainlogic

import (
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"utility"
)

// 卡牌大师
type TCardMasterInfo struct {
	CardList            []uint16
	ExchangeTimes       []uint16
	Point               int
	FreeNormalDrawTimes byte
	JiFen               [2]int
	TotalJiFen          int
	IsGetTodayRankAward bool // 今日排行奖励
	IsGetTotalRankAward bool // 累计排行奖励

	ActivityID     int              //! 活动ID
	VersionCode    int32            //! 版本号
	ResetCode      int32            //! 迭代号
	activityModule *TActivityModule //! 指针
}

//！ 活动框架代码
func (self *TCardMasterInfo) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.VersionCode = vercode
	self.ResetCode = resetcode
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.CardList = make([]uint16, len(gamedata.G_CardCsv))
	self.ExchangeTimes = make([]uint16, len(gamedata.G_CMExchangeItemCsv))
}
func (self *TCardMasterInfo) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}
func (self *TCardMasterInfo) Refresh(versionCode int32) {
	self.VersionCode = versionCode

	self.FreeNormalDrawTimes = gamedata.CardMaster_FreeTimes

	self.ExchangeTimes = make([]uint16, len(gamedata.G_CMExchangeItemCsv))

	self.IsGetTodayRankAward = false
	self.SetTodayScore(0)

	self.DB_Refresh()
}
func (self *TCardMasterInfo) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.FreeNormalDrawTimes = gamedata.CardMaster_FreeTimes

	self.TotalJiFen = 0
	self.SetYesterdayScore(0)
	self.SetTodayScore(0)

	self.IsGetTodayRankAward = false
	self.IsGetTotalRankAward = false

	self.ExchangeTimes = make([]uint16, len(gamedata.G_CMExchangeItemCsv))

	self.DB_Refresh()
}

func (self *TCardMasterInfo) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TCardMasterInfo) GetResetV() int32 {
	return self.ResetCode
}

func (self *TCardMasterInfo) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	//! 检查排行榜是否有名次
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(self.ActivityID)
	if isHandleTime {
		//! 检查昨日排行榜
		for _, v := range G_CardMasterYesterdayRanker.List {
			if v.RankID == self.activityModule.PlayerID && self.IsGetTodayRankAward == false {
				return true
			}
		}
	} else {
		//! 检查总排行榜
		for _, v := range G_CardMasterTotalRanker.List {
			if v.RankID == self.activityModule.PlayerID && self.IsGetTotalRankAward == false {
				return true
			}
		}
	}
	return false
}

//！ 数据操作代码
func (self *TCardMasterInfo) AddCard(id int, count_ int) int {
	var count uint16 = uint16(count_)
	if count <= 0 || id <= 0 || id >= len(self.CardList) {
		gamelog.Error("AddItem Error : Invalid id :%d, count:%d", id, count)
		return 0
	}
	self.CardList[id] += count
	self.DB_SaveCardList()
	return int(self.CardList[id])
}
func (self *TCardMasterInfo) DelCard(id int, count_ int) bool {
	var count uint16 = uint16(count_)
	if count <= 0 || id <= 0 || id >= len(self.CardList) {
		gamelog.Error3("DelCard Error : Invalid id :%d, count:%d", id, count)
		return false
	}
	if self.CardList[id] < count {
		return false
	} else {
		self.CardList[id] -= count
		self.DB_SaveCardList()
		return true
	}
}
func (self *TCardMasterInfo) DelCards(cards []gamedata.ST_ItemData) bool {
	// 判断是否都能删（全都）
	MaxLen := len(self.CardList)
	for _, v := range cards {
		if v.ItemNum <= 0 || v.ItemID <= 0 || v.ItemID >= MaxLen {
			gamelog.Error("DelCards Error : Invalid id:%d, count:%d", v.ItemID, v.ItemNum)
			return false
		}
		if self.CardList[v.ItemID] < uint16(v.ItemNum) {
			gamelog.Error("DelCards Error : itemId:%d, haveCnt:%d, needCnt:%d", v.ItemID, self.CardList[v.ItemID], v.ItemNum)
			return false
		}
	}
	// 判断都通过，批量删除
	for _, v := range cards {
		self.CardList[v.ItemID] -= uint16(v.ItemNum)
	}

	self.DB_SaveCardList()
	return true
}
func (self *TCardMasterInfo) AddExchangeTimes(id int, count_ int) int {
	var count uint16 = uint16(count_)
	if count <= 0 || id <= 0 || id >= len(self.ExchangeTimes) {
		gamelog.Error("AddExchangeTimes Error : Invalid id :%d, count:%d", id, count)
		return 0
	}
	self.ExchangeTimes[id] += count
	self.DB_SaveExchangeTimes()
	return int(self.ExchangeTimes[id])
}
func (self *TCardMasterInfo) AddJiFen(num int) {
	newJiFen := self.GetTodayScore() + num
	self.SetTodayScore(newJiFen)
	self.TotalJiFen += num
	self.DB_SaveJiFen()

	G_CardMasterTodayRanker.SetRankItem(self.activityModule.PlayerID, newJiFen)
	G_CardMasterTotalRanker.SetRankItem(self.activityModule.PlayerID, self.TotalJiFen)

	//! 限时任务完成
	self.activityModule.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_CARD_MASTER_SCORE, num)
}
func (self *TCardMasterInfo) GetTodayScore() int {
	return self.JiFen[utility.GetCurDayMod()]
}
func (self *TCardMasterInfo) SetTodayScore(num int) {
	self.JiFen[utility.GetCurDayMod()] = num
}
func (self *TCardMasterInfo) GetYesterdayScore() int {
	if utility.GetCurDayMod() == 1 {
		return self.JiFen[0]
	} else {
		return self.JiFen[1]
	}
}
func (self *TCardMasterInfo) SetYesterdayScore(num int) {
	if utility.GetCurDayMod() == 1 {
		self.JiFen[0] = num
	} else {
		self.JiFen[1] = num
	}
}
func (self *TCardMasterInfo) GetTotalScore() int {
	return self.TotalJiFen
}

//! DB相关
func (self *TCardMasterInfo) DB_SaveCardList() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{"cardmaster.cardlist": self.CardList}})
}
func (self *TCardMasterInfo) DB_SavePoint() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{"cardmaster.point": self.Point}})
}
func (self *TCardMasterInfo) DB_SaveFreeTimes() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{"cardmaster.freenormaldrawtimes": self.FreeNormalDrawTimes}})
}
func (self *TCardMasterInfo) DB_SaveJiFen() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"cardmaster.jifen":      self.JiFen,
		"cardmaster.totaljifen": self.TotalJiFen}})
}
func (self *TCardMasterInfo) DB_SaveRankAwardFlag() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"cardmaster.isgettodayrankaward": self.IsGetTodayRankAward,
		"cardmaster.isgettotalrankaward": self.IsGetTotalRankAward}})
}
func (self *TCardMasterInfo) DB_SaveExchangeTimes() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{"cardmaster.exchangetimes": self.ExchangeTimes}})
}
func (self *TCardMasterInfo) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"cardmaster.exchangetimes":       self.ExchangeTimes,
		"cardmaster.freenormaldrawtimes": self.FreeNormalDrawTimes,
		"cardmaster.jifen":               self.JiFen,
		"cardmaster.totaljifen":          self.TotalJiFen,
		"cardmaster.isgettodayrankaward": self.IsGetTodayRankAward,
		"cardmaster.isgettotalrankaward": self.IsGetTotalRankAward,
		"cardmaster.activityid":          self.ActivityID,
		"cardmaster.versioncode":         self.VersionCode,
		"cardmaster.resetcode":           self.ResetCode}})
}

func (self *TCardMasterInfo) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"cardmaster.exchangetimes":       self.ExchangeTimes,
		"cardmaster.freenormaldrawtimes": self.FreeNormalDrawTimes,
		"cardmaster.jifen":               self.JiFen,
		"cardmaster.totaljifen":          self.TotalJiFen,
		"cardmaster.isgettodayrankaward": self.IsGetTodayRankAward,
		"cardmaster.isgettotalrankaward": self.IsGetTotalRankAward,
		"cardmaster.activityid":          self.ActivityID,
		"cardmaster.versioncode":         self.VersionCode,
		"cardmaster.resetcode":           self.ResetCode}})
}

//！ 逻辑代码
func (self *TCardMasterInfo) NormalDraw() []gamedata.ST_ItemData { // 普通抽卡
	if self.cost_NormalDraw() {
		self.AddJiFen(gamedata.CardMaster_NormalJiFen)
		List := gamedata.GetItemsFromAwardID(gamedata.CardMaster_NormalAwardID)
		for _, v := range List {
			self.AddCard(v.ItemID, v.ItemNum)
		}
		return List
	}
	return nil
}
func (self *TCardMasterInfo) cost_NormalDraw() bool {
	if self.FreeNormalDrawTimes > 0 {
		self.FreeNormalDrawTimes--
		self.DB_SaveFreeTimes()
		return true
	}
	if self.activityModule.ownplayer.BagMoudle.RemoveNormalItem(gamedata.CardMaster_RaffleTicket, 1) {
		return true
	}
	if self.activityModule.ownplayer.RoleMoudle.CostMoney(gamedata.CardMaster_CostType, gamedata.CardMaster_NormalCost) {
		return true
	}

	gamelog.Error("CardMaster  cost_NormalDraw Error : FreeNormalDrawTimes:%d, RaffleTicket:%d",
		self.FreeNormalDrawTimes,
		self.activityModule.ownplayer.BagMoudle.GetNormalItemCount(gamedata.CardMaster_RaffleTicket))
	return false
}

func (self *TCardMasterInfo) SpecialDraw(drawType byte) []gamedata.ST_ItemData {
	diamond, awardID, jifen := GetDrawCostData(drawType)
	if self.activityModule.ownplayer.RoleMoudle.CostMoney(gamedata.CardMaster_CostType, diamond) {
		self.AddJiFen(jifen)
		List := gamedata.GetItemsFromAwardID(awardID)
		for _, v := range List {
			self.AddCard(v.ItemID, v.ItemNum)
		}
		return List
	}
	return nil
}
func GetDrawCostData(drawType byte) (diamond int, awardID int, jifen int) {
	switch drawType {
	case 1: // 普通抽
		{
			diamond = gamedata.CardMaster_NormalCost
			awardID = gamedata.CardMaster_NormalAwardID
			jifen = gamedata.CardMaster_NormalJiFen
		}
	case 2: // 普通十连
		{
			diamond = gamedata.CardMaster_NormalCost_10
			awardID = gamedata.CardMaster_NormalAwardID_10
			jifen = gamedata.CardMaster_NormalJiFen * 10
		}
	case 3: // 高级
		{
			diamond = gamedata.CardMaster_SpecialCost
			awardID = gamedata.CardMaster_SpecialAwardID
			jifen = gamedata.CardMaster_SpecialJiFen
		}
	case 4: // 高级十连
		{
			diamond = gamedata.CardMaster_SpecialCost_10
			awardID = gamedata.CardMaster_SpecialAwardID_10
			jifen = gamedata.CardMaster_SpecialJiFen * 10
		}
	default:
		{
			diamond, awardID, jifen = 0, 0, 0
		}
	}
	return diamond, awardID, jifen
}
func (self *TCardMasterInfo) Card2Item(exchangeID int) bool {
	csv := gamedata.GetCMExchangeItemCsvInfo(exchangeID)
	if csv == nil {
		gamelog.Error("Card2Item ExchangeItemCsv nil(%d)", exchangeID)
		return false
	}
	if self.ExchangeTimes[exchangeID] >= csv.DailyTimes {
		gamelog.Error("CardMaster Card2Item Error: ExchangeTimes limit. cur:%d, limit:%d", self.ExchangeTimes[exchangeID], csv.DailyTimes)
		return false
	}
	if self.DelCards(csv.NeedCards) {
		self.AddExchangeTimes(exchangeID, 1)
		self.activityModule.ownplayer.BagMoudle.AddAwardItems(csv.Items)
		return true
	}
	return false
}
func (self *TCardMasterInfo) Card2Point(cards []gamedata.ST_ItemData) bool {
	var sumPoint int
	for _, v := range cards {
		csv := gamedata.GetCardCsvInfo(v.ItemID)
		if csv == nil || csv.PointSell <= 0 {
			return false
		}
		sumPoint += csv.PointSell * v.ItemNum
	}
	if self.DelCards(cards) {
		self.Point += sumPoint
		self.DB_SavePoint()
		return true
	}
	return false
}
func (self *TCardMasterInfo) Point2Card(cards []gamedata.ST_ItemData) bool {
	var sumPoint int
	for _, v := range cards {
		csv := gamedata.GetCardCsvInfo(v.ItemID)
		if csv == nil || csv.PointBuy <= 0 {
			return false
		}
		sumPoint += csv.PointBuy * v.ItemNum
	}
	if self.Point >= sumPoint {
		self.Point -= sumPoint
		self.DB_SavePoint()
		for _, v := range cards {
			self.AddCard(v.ItemID, v.ItemNum)
		}
		return true
	}
	return false
}

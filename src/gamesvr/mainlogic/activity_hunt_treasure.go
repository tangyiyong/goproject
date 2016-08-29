package mainlogic

import (
	"appconfig"
	"fmt"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"utility"
)

type THuntStoreItem struct {
	ID    int
	IsBuy bool //! 是否已经购买
}

//! 巡回活动
type TActivityHunt struct {
	ActivityID           int              //! 活动ID
	HuntAward            Mark             //! 巡回奖励
	Score                int              //! 积分
	TodayScore           [2]int           //! 按照奇偶交替表现分数
	HuntTurns            int              //! 巡回次数
	CurrentPos           int              //! 当前坐标
	StoreItemLst         []THuntStoreItem //! 商店物品
	IsHaveStore          bool             //! 是否有商店
	IsRecvTodayRankAward bool             //! 今日排行奖励领取标记
	IsRecvTotalRankAward bool             //! 总排行领取标记
	FreeTimes            int              //! 免费次数
	VersionCode          int32            //! 版本号
	ResetCode            int32            //! 迭代号
	activityModule       *TActivityModule //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityHunt) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityHunt) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self

	self.HuntAward = 0
	self.Score = 0
	self.TodayScore = [2]int{0, 0}
	self.FreeTimes = gamedata.HuntFreeTimes
	self.HuntTurns = 0
	self.CurrentPos = 0
	self.StoreItemLst = []THuntStoreItem{}
	self.IsHaveStore = false
	self.IsRecvTodayRankAward = false
	self.IsRecvTotalRankAward = false

	self.VersionCode = vercode
	self.ResetCode = resetcode

}

//! 刷新数据
func (self *TActivityHunt) Refresh(versionCode int32) {
	self.IsRecvTodayRankAward = false
	self.FreeTimes = gamedata.HuntFreeTimes
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityHunt) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	self.HuntAward = 0
	self.Score = 0
	self.TodayScore = [2]int{0, 0}
	self.HuntTurns = 0
	self.CurrentPos = 0
	self.FreeTimes = gamedata.HuntFreeTimes
	self.StoreItemLst = []THuntStoreItem{}
	self.IsRecvTodayRankAward = false
	self.IsRecvTotalRankAward = false
	self.IsHaveStore = false
	go self.DB_Reset()
}

func (self *TActivityHunt) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityHunt) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityHunt) GetTodayScore() int {
	return self.TodayScore[utility.GetCurDayMod()]
}
func (self *TActivityHunt) GetYesterdayScore() int {
	if utility.GetCurDayMod() == 1 {
		return self.TodayScore[0]
	} else {
		return self.TodayScore[1]
	}
}
func (self *TActivityHunt) GetTotalScore() int {
	return self.Score
}

func (self *TActivityHunt) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	//! 检查排行榜是否有名次
	isEnd, _ := G_GlobalVariables.IsActivityTime(self.ActivityID)
	if isEnd == true {
		if self.FreeTimes != 0 { //! 有免费次数则返回红点
			return true
		}
		//! 检查昨日排行榜
		rank := G_HuntTreasureYesterdayRanker.GetRankIndex(self.activityModule.PlayerID, self.GetYesterdayScore())
		if rank > 0 && rank <= 50 {
			return true
		}

	} else {
		//! 检查总排行榜
		totayRank := G_HuntTreasureTotalRanker.GetRankIndex(self.activityModule.PlayerID, self.Score)
		if totayRank > 0 && totayRank <= 50 {
			return true
		}
	}

	return false
}

func (self *TActivityHunt) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.activityid":           self.ActivityID,
		"hunttreasure.huntaward":            self.HuntAward,
		"hunttreasure.score":                self.Score,
		"hunttreasure.todayscore":           self.TodayScore,
		"hunttreasure.huntturns":            self.HuntTurns,
		"hunttreasure.currentpos":           self.CurrentPos,
		"hunttreasure.storeitemlst":         self.StoreItemLst,
		"hunttreasure.isrecvtodayrankaward": self.IsRecvTodayRankAward,
		"hunttreasure.isrecvtotalrankaward": self.IsRecvTotalRankAward,
		"hunttreasure.freetimes":            self.FreeTimes,
		"hunttreasure.versioncode":          self.VersionCode,
		"hunttreasure.ishavestore":          self.IsHaveStore,
		"hunttreasure.resetcode":            self.ResetCode}})
	return true
}

func (self *TActivityHunt) DB_UpdateHuntTodayRankAward() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.isrecvtodayrankaward": self.IsRecvTodayRankAward,
		"hunttreasure.freetimes":            self.FreeTimes,
		"hunttreasure.versioncode":          self.VersionCode}})
}

func (self *TActivityHunt) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.isrecvtodayrankaward": self.IsRecvTodayRankAward,
		"hunttreasure.freetimes":            self.FreeTimes,
		"hunttreasure.versioncode":          self.VersionCode}})
	return true
}

func (self *TActivityHunt) DB_UpdateHuntTotalRankAward() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.isrecvtotalrankaward": self.IsRecvTotalRankAward}})
}

func (self *TActivityHunt) DB_UpdateHuntStore() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.storeitemlst": self.StoreItemLst}})
}

func (self *TActivityHunt) DB_SaveHuntScore() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.todayscore": self.TodayScore,
		"hunttreasure.score":      self.Score}})
}

func (self *TActivityHunt) DB_SaveFreeTiems() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.freetimes": self.FreeTimes}})
}

func (self *TActivityHunt) DB_SaveHuntTurnsAwardMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.huntaward": self.HuntAward}})
}

func (self *TActivityHunt) DB_SaveHuntStatus() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.currentpos": self.CurrentPos,
		"hunttreasure.todayscore": self.TodayScore,
		"hunttreasure.score":      self.Score,
		"hunttreasure.huntturns":  self.HuntTurns}})
}

func (self *TActivityHunt) DB_ChangeHuntStoreItemMark(index int) {
	filedName := fmt.Sprintf("hunttreasure.storeitemlst.%d.isbuy", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: true}})
}

func (self *TActivityHunt) DB_SaveStoreMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"hunttreasure.ishavestore": self.IsHaveStore}})
}

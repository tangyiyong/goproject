package mainlogic

import (
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
	ActivityID   int32            //! 活动ID
	HuntAward    BitsType         //! 巡回奖励
	Score        int              //! 积分
	TodayScore   [2]int           //! 按照奇偶交替表现分数
	HuntTurns    int              //! 巡回次数
	CurrentPos   int              //! 当前坐标
	StoreItemLst []THuntStoreItem //! 商店物品
	IsHaveStore  bool             //! 是否有商店
	RankAward    [2]int8          //! 排行奖励领取标记 //0:表示今天，1:表示总榜
	FreeTimes    int              //! 免费次数
	VersionCode  int32            //! 版本号
	ResetCode    int32            //! 迭代号
	modulePtr    *TActivityModule //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityHunt) SetModulePtr(mPtr *TActivityModule) {
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityHunt) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self

	self.HuntAward = 0
	self.Score = 0
	self.TodayScore = [2]int{0, 0}
	self.FreeTimes = gamedata.HuntFreeTimes
	self.HuntTurns = 0
	self.CurrentPos = 0
	self.StoreItemLst = []THuntStoreItem{}
	self.IsHaveStore = false
	self.RankAward[0] = 0
	self.RankAward[1] = 0

	self.VersionCode = vercode
	self.ResetCode = resetcode

}

//! 刷新数据
func (self *TActivityHunt) Refresh(versionCode int32) {
	self.RankAward[0] = 0
	self.FreeTimes = gamedata.HuntFreeTimes
	self.VersionCode = versionCode
	self.DB_Refresh()
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
	self.RankAward[0] = 0
	self.RankAward[1] = 0
	self.IsHaveStore = false
	self.DB_Reset()
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
	isEnd := G_GlobalVariables.IsActivityTime(self.ActivityID)
	if isEnd == true {
		if self.FreeTimes != 0 { //! 有免费次数则返回红点
			return true
		}
		//! 检查昨日排行榜
		rank := G_HuntTreasureYesterdayRanker.GetRankIndex(self.modulePtr.PlayerID, self.GetYesterdayScore())
		if rank > 0 && rank <= 50 {
			return true
		}

	} else {
		//! 检查总排行榜
		totayRank := G_HuntTreasureTotalRanker.GetRankIndex(self.modulePtr.PlayerID, self.Score)
		if totayRank > 0 && totayRank <= 50 {
			return true
		}
	}

	return false
}

func (self *TActivityHunt) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.activityid":   self.ActivityID,
		"hunttreasure.huntaward":    self.HuntAward,
		"hunttreasure.score":        self.Score,
		"hunttreasure.todayscore":   self.TodayScore,
		"hunttreasure.huntturns":    self.HuntTurns,
		"hunttreasure.currentpos":   self.CurrentPos,
		"hunttreasure.storeitemlst": self.StoreItemLst,
		"hunttreasure.rankaward":    self.RankAward,
		"hunttreasure.freetimes":    self.FreeTimes,
		"hunttreasure.versioncode":  self.VersionCode,
		"hunttreasure.ishavestore":  self.IsHaveStore,
		"hunttreasure.resetcode":    self.ResetCode}})
}

func (self *TActivityHunt) DB_UpdateHuntRankAward() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.rankaward":   self.RankAward,
		"hunttreasure.freetimes":   self.FreeTimes,
		"hunttreasure.versioncode": self.VersionCode}})
}

func (self *TActivityHunt) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.rankaward":   self.RankAward,
		"hunttreasure.freetimes":   self.FreeTimes,
		"hunttreasure.versioncode": self.VersionCode}})
}

func (self *TActivityHunt) DB_UpdateHuntStore() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.storeitemlst": self.StoreItemLst}})
}

func (self *TActivityHunt) DB_SaveHuntScore() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.todayscore": self.TodayScore,
		"hunttreasure.score":      self.Score}})
}

func (self *TActivityHunt) DB_SaveFreeTiems() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.freetimes": self.FreeTimes}})
}

func (self *TActivityHunt) DB_SaveHuntTurnsAwardMark() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.huntaward": self.HuntAward}})
}

func (self *TActivityHunt) DB_SaveHuntStatus() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.currentpos": self.CurrentPos,
		"hunttreasure.todayscore": self.TodayScore,
		"hunttreasure.score":      self.Score,
		"hunttreasure.huntturns":  self.HuntTurns}})
}

func (self *TActivityHunt) DB_ChangeHuntStoreItemMark(index int) {
	filedName := fmt.Sprintf("hunttreasure.storeitemlst.%d.isbuy", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		filedName: true}})
}

func (self *TActivityHunt) DB_SaveStoreMark() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"hunttreasure.ishavestore": self.IsHaveStore}})
}

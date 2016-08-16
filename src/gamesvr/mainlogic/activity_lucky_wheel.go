package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"mongodb"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

//! 幸运轮盘
type TActivityWheel struct {
	ActivityID           int              //! 活动ID
	NormalAwardLst       []int            //! 奖品ID列表
	ExcitedAwardLst      []int            //! 豪华ID列表
	NormalFreeTimes      int              //! 普通轮盘免费次数
	ExcitedFreeTimes     int              //! 豪华轮盘免费次数
	TodayScore           [2]int           //! 奇偶分数交换使用
	TotalScore           int              //! 总分数
	IsRecvTodayRankAward bool             //! 今日排行奖励领取标记
	IsRecvTotalRankAward bool             //! 总排行领取标记
	VersionCode          int              //! 版本号
	ResetCode            int              //! 迭代号
	activityModule       *TActivityModule //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityWheel) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityWheel) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self

	activityInfo := gamedata.GetActivityInfo(activityID)

	self.NormalAwardLst = gamedata.GetLuckyWheelItemFromDay(activityInfo.AwardType, 1)
	self.ExcitedAwardLst = gamedata.GetLuckyWheelItemFromDay(activityInfo.AwardType, 2)

	self.NormalFreeTimes = gamedata.NormalWheelFreeTimes
	self.ExcitedFreeTimes = gamedata.ExcitedWheelFreeTimes
	self.IsRecvTodayRankAward = false
	self.IsRecvTotalRankAward = false
	self.TotalScore = 0
	self.TodayScore = [2]int{0, 0}
	self.VersionCode = vercode
	self.ResetCode = resetcode

}

//! 刷新数据
func (self *TActivityWheel) Refresh(versionCode int) {
	//! 数据变更
	self.IsRecvTodayRankAward = false
	self.NormalFreeTimes = gamedata.NormalWheelFreeTimes
	self.ExcitedFreeTimes = gamedata.ExcitedWheelFreeTimes
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityWheel) End(versionCode int, resetCode int) {
	self.NormalAwardLst = []int{}
	self.ExcitedAwardLst = []int{}
	self.NormalFreeTimes = 0
	self.ExcitedFreeTimes = 0
	self.IsRecvTodayRankAward = false
	self.IsRecvTotalRankAward = false
	self.TotalScore = 0
	self.TodayScore = [2]int{0, 0}
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	//! 奖金池清空
	G_GlobalVariables.DB_CleanMoneyPoor()

	go self.DB_Reset()
}

func (self *TActivityWheel) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityWheel) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityWheel) GetTodayScore() int {
	if utility.GetCurDayMod() == 1 {
		return self.TodayScore[1]
	} else {
		return self.TodayScore[0]
	}
}
func (self *TActivityWheel) GetYesterdayScore() int {
	if utility.GetCurDayMod() == 1 {
		return self.TodayScore[0]
	} else {
		return self.TodayScore[1]
	}
}
func (self *TActivityWheel) GetTotalScore() int {
	return self.TotalScore
}

func (self *TActivityWheel) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	//! 检查排行榜是否有名次
	isEnd, _ := G_GlobalVariables.IsActivityTime(self.ActivityID)
	if isEnd == true {
		if self.NormalFreeTimes != 0 || self.ExcitedFreeTimes != 0 {
			return true //! 拥有免费次数则返回红点
		}

		//! 检查昨日排行榜
		for _, v := range G_LuckyWheelYesterdayRanker.List {
			if v.RankID == self.activityModule.PlayerID && self.IsRecvTodayRankAward == false {
				return true
			}
		}
	} else {
		//! 检查总排行榜
		for _, v := range G_LuckyWheelTotalRanker.List {
			if v.RankID == self.activityModule.PlayerID && self.IsRecvTotalRankAward == false {
				return true
			}
		}
	}

	return false
}

//! 随机一个轮盘奖励
func (self *TActivityWheel) RandWheelAward(wheelType int) (itemID int, itemNum int, isSpecial int, index int) {

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	//! 随机奖励
	awardLst := []gamedata.ST_LuckyWheel{}
	totalWeight := 0

	moneyPercent := 0

	if wheelType == 1 {
		for _, v := range self.NormalAwardLst {
			award := gamedata.GetLuckyWheelItemFromID(v)
			if award == nil {
				gamelog.Error("GetLuckyWheelItemFromID Error: Invalid ID %d", v)
				return 0, 0, 0, 0
			}
			totalWeight += award.Weight
			awardLst = append(awardLst, *award)

			if award.IsSpecial == 1 {
				moneyPercent = award.ItemNum
			}
		}
	} else if wheelType == 2 {
		for _, v := range self.ExcitedAwardLst {
			award := gamedata.GetLuckyWheelItemFromID(v)
			if award == nil {
				gamelog.Error("GetLuckyWheelItemFromID Error: Invalid ID %d", v)
				return 0, 0, 0, 0
			}
			totalWeight += award.Weight
			awardLst = append(awardLst, *award)

			if award.IsSpecial == 1 {
				moneyPercent = award.ItemNum
			}
		}
	} else {
		return 0, 0, 0, 0
	}

	for {
		randomWeight := random.Intn(totalWeight)

		curWeight := 0
		for i, v := range awardLst {
			if curWeight <= randomWeight && randomWeight < curWeight+v.Weight {
				if v.IsSpecial == 1 {
					if wheelType == 1 && G_GlobalVariables.NormalMoneyPoor*moneyPercent/10000 < 100 {
						continue
					} else if wheelType == 2 && G_GlobalVariables.ExcitedMoneyPoor*moneyPercent/10000 < 100 {
						continue
					}
				}
				return v.ItemID, v.ItemNum, v.IsSpecial, i
			}

			curWeight += v.Weight
		}
	}

	return 0, 0, 0, 0
}

func (self *TActivityWheel) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"luckywheel.isrecvtodayrankaward": self.IsRecvTodayRankAward,
		"luckywheel.normalfreetimes":      self.NormalFreeTimes,
		"luckywheel.excitedfreetimes":     self.ExcitedFreeTimes,
		"luckywheel.versioncode":          self.VersionCode}})
	return true
}

func (self *TActivityWheel) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"luckywheel.activityid":           self.ActivityID,
		"luckywheel.normalawardlst":       self.NormalAwardLst,
		"luckywheel.excitedawardlst":      self.ExcitedAwardLst,
		"luckywheel.normalfreetimes":      self.NormalFreeTimes,
		"luckywheel.excitedfreetimes":     self.ExcitedFreeTimes,
		"luckywheel.versioncode":          self.VersionCode,
		"luckywheel.isrecvtotalrankaward": self.IsRecvTotalRankAward,
		"luckywheel.isrecvtodayrankaward": self.IsRecvTodayRankAward,
		"luckywheel.todayscore":           self.TodayScore,
		"luckywheel.totalscore":           self.TotalScore,
		"luckywheel.resetcode":            self.ResetCode}})
	return true
}

func (self *TActivityWheel) DB_SaveLuckyWheelScore() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"luckywheel.todayscore": self.TodayScore,
		"luckywheel.totalscore": self.TotalScore}})
}

func (self *TActivityWheel) DB_SaveLuckyWheelFreeTimes(wheelType int) {
	if wheelType == 1 {
		mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
			"luckywheel.normalfreetimes": self.NormalFreeTimes}})
	} else if wheelType == 2 {
		mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
			"luckywheel.excitedfreetimes": self.ExcitedFreeTimes}})
	}
}

func (self *TActivityWheel) DB_UpdateWheelTodayRankAward() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"luckywheel.isrecvtodayrankaward": self.IsRecvTodayRankAward}})
}

func (self *TActivityWheel) DB_UpdateWheelTotalRankAward() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"luckywheel.isrecvtotalrankaward": self.IsRecvTotalRankAward}})
}
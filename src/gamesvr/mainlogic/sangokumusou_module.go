package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"mongodb"
	"msg"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TSangokuMusouAttrData struct {
	ID       int
	CostStar int //! 消耗星数
	AttrID   int //! 属性ID
	Value    int //! 加成概率
}

type TSangokuMusouAttrData2 struct {
	AttrID int //! 属性ID
	Value  int //! 加成概率
}

type TSangokuMusouCopyInfo struct {
	CopyID  int //! 副本ID
	StarNum int //! 星数
}

type TSangokuMusouStarRank struct {
	PlayerID int64
	Star     int
}

type TSangokuMusouModule struct {
	PlayerID            int32                    `bson:"_id"`
	PassCopyID          int                      //! 当前通过的关卡ID
	CopyInfoLst         []TSangokuMusouCopyInfo  //! 通过关卡信息
	PassEliteCopyID     int                      //! 当前通过的精英关卡ID
	CurStar             int                      //! 当前星数
	HistoryStar         int                      //! 历史最高星数
	HistoryCopyID       int                      //! 历史通关最高关卡ID
	CanUseStar          int                      //! 当前可用星数
	BattleTimes         int                      //! 普通挑战次数
	EliteBattleTimes    int                      //! 精英挑战次数
	AddEliteBattleTimes int                      //! 精英副本增加次数
	ShoppingLst         []TStoreBuyData          //! 已购买物品列表
	AttrMarkupLst       []TSangokuMusouAttrData2 //! 属性加成列表
	AwardAttrLst        []TSangokuMusouAttrData  //! 章节属性奖励选择列表
	ChapterAwardMark    IntLst                   //! 章节奖励领取标记
	ChapterBuffMark     IntLst                   //! 章节奖励Buff领取标记
	TreasureID          int                      //! 无双迷藏
	IsBuyTreasure       bool                     //! 是否已购买无双秘藏
	IsEnd               bool                     //! 是否已经结束
	ResetDay            uint32                   //! 重置天数
	ownplayer           *TPlayer
}

func (self *TSangokuMusouModule) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

//! 玩家创建角色
func (self *TSangokuMusouModule) OnCreate(playerid int32) {
	//! 初始化信息
	self.IsEnd = false
	self.EliteBattleTimes = gamedata.SangokuMusouEliteFreeTimes
	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerSangokuMusou", self)
}

func (self *TSangokuMusouModule) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TSangokuMusouModule) OnNewDay(newday uint32) {
	self.ResetDay = newday

	self.BattleTimes = 0
	self.AddEliteBattleTimes = 0
	self.EliteBattleTimes = gamedata.SangokuMusouEliteFreeTimes

	//! 刷新物品购买次数
	for i, v := range self.ShoppingLst {
		item := gamedata.GetSangokumusouStoreInfo(v.ID)
		if item.ItemType != 4 && item.BuyTimes != 0 {
			self.ShoppingLst[i].Times = 0
		}
	}

	go self.UpdateResetTime()
}

//! 玩家销毁角色
func (self *TSangokuMusouModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (self *TSangokuMusouModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离线
func (self *TSangokuMusouModule) OnPlayerOffline(playerid int32) {

}

//! 预取玩家信息
func (self *TSangokuMusouModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerSangokuMusou").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerSangokuMusou Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 获取物品购买次数
func (self *TSangokuMusouModule) GetItemBuyTimes(id int) int {
	for _, v := range self.ShoppingLst {
		if v.ID == id {
			return v.Times
		}
	}

	return 0
}

//! 获取无双秘藏
func (self *TSangokuMusouModule) GetMusouTreasure() int {
	//! 判断当前玩家是否结束
	if self.TreasureID == 0 && self.IsBuyTreasure == false {
		self.TreasureID = gamedata.RandMusouTreasure(self.CurStar)

		go self.UpdateTreasure()
		return self.TreasureID
	}

	return self.TreasureID
}

func (self *TSangokuMusouModule) RedTip() bool {
	if self.BattleTimes == 0 {
		return true
	}

	if self.EliteBattleTimes == 0 {
		return true
	}

	for _, v := range gamedata.GT_SangokuMusou_Store {
		if v.ItemType == 4 && self.HistoryStar >= v.NeedStar && self.ownplayer.GetLevel() >= v.NeedLevel {
			isExist := false
			for _, n := range self.ShoppingLst {
				if n.ID == v.ID {
					isExist = true
					break
				}
			}
			if isExist == false {
				return true
			}
		}
	}

	return false
}

//! 通关三国无双普通关卡
func (self *TSangokuMusouModule) PassCopy(copyID int, starNum int, isVictory bool) (dropItem []msg.MSG_SangokuMusouDropItem) {
	//! 判断胜利失败
	if isVictory == false {
		self.IsEnd = true
	}

	if self.IsEnd == true {
		go self.UpdateIsEndMark()
		return
	}

	//! 增加当前星数
	self.CurStar += starNum
	self.CanUseStar += starNum

	//! 对比历史最高星数
	if self.HistoryStar < self.CurStar {
		self.HistoryStar = self.CurStar

		self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SGWS_STAR, self.HistoryStar)

		//! 对比排行数据,取决是否计入内存
		if G_SgwsStarRanker.SetRankItem(self.PlayerID, self.HistoryStar) > 0 {
			self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SGWS_RANK, self.HistoryStar)
		}
	}

	//! 记录通关关卡
	self.PassCopyID = copyID

	//! 对比历史通关最高关卡
	if self.HistoryCopyID < self.PassCopyID {
		self.HistoryCopyID = self.PassCopyID
	}

	//! 信息写入数据库
	go self.UpdatePassCopyRecord()

	var info TSangokuMusouCopyInfo
	info.CopyID = copyID
	info.StarNum = starNum

	self.CopyInfoLst = append(self.CopyInfoLst, info)

	//! 记录通关信息
	go self.AddPassCopyInfoLst(info)

	//! 获取副本信息
	copyData := gamedata.GetSangokuMusouChapterInfo(copyID)

	//! 获取奖励ID
	awardID := 0
	if starNum == 1 {
		awardID = copyData.Diffculty1
	} else if starNum == 2 {
		awardID = copyData.Diffculty2
	} else if starNum == 3 {
		awardID = copyData.Diffculty3
	}

	//! 掉落奖励
	dropAward := gamedata.GetItemsFromAwardID(awardID)
	award := []msg.MSG_SangokuMusouDropItem{}
	for _, v := range dropAward {
		item := msg.MSG_SangokuMusouDropItem{
			ItemID:   v.ItemID,
			ItemNum:  v.ItemNum,
			CritType: 0}

		award = append(award, item)
	}

	//! 计算奖励暴击
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < len(award); i++ {
		randValue := r.Intn(100)

		if randValue >= 0 && randValue < gamedata.LuckyCritPro {
			award[i].ItemNum = award[i].ItemNum * 1000 / gamedata.LuckyCritMultiple
			award[i].CritType = 3
		} else if randValue >= 0 && randValue < gamedata.BigCritPro {
			award[i].ItemNum = award[i].ItemNum * 1000 / gamedata.LuckyCritMultiple
			award[i].CritType = 2
		} else if randValue >= 0 && randValue < gamedata.CritPro {
			award[i].ItemNum = award[i].ItemNum * 1000 / gamedata.LuckyCritMultiple
			award[i].CritType = 1
		}
	}

	dropItem = append(dropItem, award...)

	//! 发放奖励
	for _, v := range dropItem {
		self.ownplayer.BagMoudle.AddAwardItem(v.ItemID, v.ItemNum)
	}

	return dropItem
}

//! 通关三国无双精英关卡
func (self *TSangokuMusouModule) PassEliteCopy(copyID int) bool {
	//! 挑战次数-1
	self.EliteBattleTimes -= 1

	//! 判断是否首胜
	isFirstVictory := false
	if copyID > self.PassEliteCopyID {
		isFirstVictory = true
		self.PassEliteCopyID = copyID
	}

	//! 获取精英挑战信息
	copyInfo := gamedata.GetSangokuMusouEliteCopyInfo(copyID)
	if copyInfo == nil {
		gamelog.Error("GetSGMSEliteCopyInfo fail. CopyID: %d", copyInfo.CopyID)
		return false
	}

	copyBase := gamedata.GetCopyBaseInfo(copyInfo.CopyID)
	if copyBase == nil {
		gamelog.Error("GetCopyBaseInfo fail. CopyID: %d", copyInfo.CopyID)
		return false
	}

	if isFirstVictory == true {
		firstAward := gamedata.GetItemsFromAwardID(copyBase.FirstAward)
		self.ownplayer.BagMoudle.AddAwardItems(firstAward)
	}

	//! 获取普通奖励
	normalAward := gamedata.GetItemsFromAwardID(copyBase.AwardID)
	self.ownplayer.BagMoudle.AddAwardItems(normalAward)

	go self.UpdatePassEliteCopyRecord()

	return isFirstVictory
}

//! 扫荡章节
func (self *TSangokuMusouModule) SweepChapter(chapter int) (dropItem [][]msg.MSG_SangokuMusouDropItem) {
	chapterInfo := gamedata.GetSGWSChapterCopyLst(chapter)

	self.PassCopyID = gamedata.GetSGWSChapterEndCopyID(chapter)

	//! 三星扫荡
	for i := 0; i < len(chapterInfo); i++ {
		self.CurStar += 3
		self.CanUseStar += 3

		var info TSangokuMusouCopyInfo
		info.CopyID = chapterInfo[i]
		info.StarNum = 3
		self.CopyInfoLst = append(self.CopyInfoLst, info)

		if self.CurStar > self.HistoryStar {
			self.HistoryStar = self.CurStar
		}

		//! 获取副本信息
		copyData := gamedata.GetSangokuMusouChapterInfo(info.CopyID)

		awardID := copyData.Diffculty3

		//! 掉落奖励
		dropAward := gamedata.GetItemsFromAwardID(awardID)

		award := []msg.MSG_SangokuMusouDropItem{}
		for _, v := range dropAward {
			item := msg.MSG_SangokuMusouDropItem{
				ItemID:   v.ItemID,
				ItemNum:  v.ItemNum,
				CritType: 0}

			award = append(award, item)
		}

		//! 计算暴击
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for i := 0; i < len(award); i++ {
			randValue := r.Intn(100)

			if randValue >= 0 && randValue < gamedata.LuckyCritPro {
				award[i].ItemNum = award[i].ItemNum * 1000 / gamedata.LuckyCritMultiple
				award[i].CritType = 3
			} else if randValue >= 0 && randValue < gamedata.BigCritPro {
				award[i].ItemNum = award[i].ItemNum * 1000 / gamedata.LuckyCritMultiple
				award[i].CritType = 2
			} else if randValue >= 0 && randValue < gamedata.CritPro {
				award[i].ItemNum = award[i].ItemNum * 1000 / gamedata.LuckyCritMultiple
				award[i].CritType = 1
			}
		}

		item := append([]msg.MSG_SangokuMusouDropItem{}, award...)

		dropItem = append(dropItem, item)

	}

	go self.UpdatePassCopyRecord()

	return dropItem
}

//! 重置普通关卡
func (self *TSangokuMusouModule) ResetCopy() {
	self.PassCopyID = 0
	self.CopyInfoLst = []TSangokuMusouCopyInfo{}
	self.TreasureID = 0
	self.IsBuyTreasure = false
	self.BattleTimes += 1
	self.AttrMarkupLst = []TSangokuMusouAttrData2{}
	self.AwardAttrLst = []TSangokuMusouAttrData{}
	self.CurStar = 0
	self.CanUseStar = 0
	self.ChapterAwardMark = IntLst{}
	self.ChapterBuffMark = IntLst{}
	self.IsEnd = false

	self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SGWS_RESET, 1)

	go self.UpdateResetCopy()
}

func (self *TSangokuMusouModule) IsShoppingInfoExist(id int) bool {
	for _, v := range self.ShoppingLst {
		if v.ID == id {
			return true
		}
	}
	return false
}

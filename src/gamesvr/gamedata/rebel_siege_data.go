package gamedata

import (
	"gamelog"
)

//! 根据怪物ID取得基础属性 + 等级加成属性 = 怪物最终属性
type ST_RebelSiege struct {
	CopyID        int //! 副本
	LifeValue     int //! 基础生命值
	Difficulty    int //! 难度 1->普通的 2->勇猛的 3->无双的
	RangeLevelMin int //! 玩家等级取值范围
	RangeLevelMax int //! 玩家等级取值范围
}

var GT_RebelSiegeLst []ST_RebelSiege

func InitRebelSiegeParser(total int) bool {
	GT_RebelSiegeLst = make([]ST_RebelSiege, total+1)
	return true
}

func ParseRebelSiegeRecord(rs *RecordSet) {
	copyID := CheckAtoi(rs.Values[0], 0)

	var rebel ST_RebelSiege
	rebel.CopyID = copyID
	rebel.LifeValue = rs.GetFieldInt("lifevalue")
	rebel.Difficulty = rs.GetFieldInt("difficulty")
	rebel.RangeLevelMin = rs.GetFieldInt("range_level_min")
	rebel.RangeLevelMax = rs.GetFieldInt("range_level_max")

	GT_RebelSiegeLst = append(GT_RebelSiegeLst, rebel)
}

//! 随机一个叛军
func RandRebel(playerLevel int) *ST_RebelSiege {
	//! 随机叛军品质
	proLst := []int{LowerRebelPro, MiddleRebelPro, SeniorRebelPro}
	randValue := r.Intn(100)
	curValue := 0
	quality := 0
	for i, v := range proLst {
		if randValue >= curValue && randValue < curValue+v {
			quality = i + 1
			break
		}
		curValue += v
	}

	for i, v := range GT_RebelSiegeLst {
		if playerLevel >= v.RangeLevelMin && playerLevel <= v.RangeLevelMax && v.Difficulty == quality {
			return &GT_RebelSiegeLst[i]
		}
	}
	return nil
}

//! 获取叛军信息
func GetRebelInfo(id int) *ST_RebelSiege {
	for i, v := range GT_RebelSiegeLst {
		if v.CopyID == id {
			return &GT_RebelSiegeLst[i]
		}
	}

	gamelog.Error("GetRebelInfo error: not found rebel")
	return nil
}

//! 功勋奖励
type ST_Exploit_Award struct {
	ID          int //! 唯一标识
	ItemID      int //! 物品ID
	ItemNum     int //! 物品数量
	NeedExploit int //! 需求功勋
	MinLevel    int //! 最小等级区间
	MaxLevel    int //! 最大等级区间
}

var GT_ExploitAwardLst []ST_Exploit_Award

func InitExploitAwardParser(total int) bool {
	GT_ExploitAwardLst = make([]ST_Exploit_Award, total+1)
	return true
}

func ParseExploitAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_ExploitAwardLst[id].ID = id
	GT_ExploitAwardLst[id].ItemID = rs.GetFieldInt("itemid")
	GT_ExploitAwardLst[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_ExploitAwardLst[id].NeedExploit = rs.GetFieldInt("needexploit")
	GT_ExploitAwardLst[id].MinLevel = rs.GetFieldInt("minlevel")
	GT_ExploitAwardLst[id].MaxLevel = rs.GetFieldInt("maxlevel")
}

//! 获取功绩奖励信息
func GetExploitAward(id int) *ST_Exploit_Award {
	if id >= len(GT_ExploitAwardLst) || id <= 0 {
		gamelog.Error("GetExploitAward Fail. Invalid id: %d", id)
		return nil
	}

	return &GT_ExploitAwardLst[id]
}

//! 根据等级返回奖励列表
func GetExploitAwardFromLevel(level int) (awardLst []*ST_Exploit_Award) {
	length := len(GT_ExploitAwardLst)
	for i := 0; i < length; i++ {
		if level >= GT_ExploitAwardLst[i].MinLevel && level <= GT_ExploitAwardLst[i].MaxLevel {
			awardLst = append(awardLst, &GT_ExploitAwardLst[i])
		}
	}

	return awardLst
}

//! 战功商店
type ST_Exploit_Store struct {
	ID           int //! 唯一标识
	ItemID       int //! 商品ID
	ItemNum      int //! 商品数量
	NeedMoneyID  int //! 需要金钱ID
	NeedMoneyNum int //! 需要金钱数量
	NeedItemID   int //! 需要物品ID
	NeedItemNum  int //! 需要物品数量
	NeedLevel    int //! 开放等级
}

var GT_ExploitStoreLst []ST_Exploit_Store

func InitExploitStoreParser(total int) bool {
	GT_ExploitStoreLst = make([]ST_Exploit_Store, total+1)
	return true
}

func ParseExploitStoreRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_ExploitStoreLst[id].ID = id
	GT_ExploitStoreLst[id].ItemID = rs.GetFieldInt("itemid")
	GT_ExploitStoreLst[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_ExploitStoreLst[id].NeedItemID = rs.GetFieldInt("needitemid")
	GT_ExploitStoreLst[id].NeedItemNum = rs.GetFieldInt("needitemnum")
	GT_ExploitStoreLst[id].NeedMoneyID = rs.GetFieldInt("needmoneyid")
	GT_ExploitStoreLst[id].NeedMoneyNum = rs.GetFieldInt("needmoneynum")
	GT_ExploitStoreLst[id].NeedLevel = rs.GetFieldInt("needlevel")
}

func GetExploitStoreItemIDLst(level int) (itemLst []int) {
	for _, v := range GT_ExploitStoreLst {
		if v.NeedLevel < level {
			itemLst = append(itemLst, v.ID)
		}
	}

	return itemLst
}

func GetExploitStoreItemInfo(id int) *ST_Exploit_Store {
	if id >= len(GT_ExploitStoreLst) || id <= 0 {
		gamelog.Error("GetExploitStoreItemInfo Fail. Invalid id: %d", id)
		return nil
	}

	return &GT_ExploitStoreLst[id]
}

//! 叛军围剿排行榜奖励
type ST_RebelRankAward struct {
	AwardID1    int //! 功勋排行榜奖励
	AwardID2    int //! 伤害排行榜奖励
	RankLowest  int //! 排名取值区间
	RankHighest int
}

var GT_RebelRankAwardLst []ST_RebelRankAward

func InitRebelRankAwardParser(total int) bool {
	//GT_RebelRankAwardLst = make([]ST_RebelRankAward, total+1)
	return true
}

func ParseRebelRankAwardRecord(rs *RecordSet) {
	var award ST_RebelRankAward
	award.AwardID1 = rs.GetFieldInt("awardid1")
	award.AwardID2 = rs.GetFieldInt("awardid2")
	award.RankLowest = rs.GetFieldInt("ranklowest")
	award.RankHighest = rs.GetFieldInt("rankhighest")
	GT_RebelRankAwardLst = append(GT_RebelRankAwardLst, award)
}

func GetRebelRankAward(rank int, rankType int) int {
	for i, v := range GT_RebelRankAwardLst {
		if rankType == 1 {
			if v.RankLowest <= rank && rank <= v.RankHighest {
				return GT_RebelRankAwardLst[i].AwardID1
			}
		} else if rankType == 2 {
			if v.RankLowest <= rank && rank <= v.RankHighest {
				return GT_RebelRankAwardLst[i].AwardID2
			}
		}
	}

	gamelog.Error("GetRebelRankAward Error：Invlid rank: %d  rantype: %d", rank, rankType)
	return 0
}

//! 叛军发现/击杀奖励
const (
	Find_Rebel = 1
	Kill_Rebel = 2
)

type ST_RebelActionAward struct {
	Action    int //! 类型: 1->发现叛军  2->击杀叛军
	Diffculty int //! 1->普通 2->勇猛 3->无双
	MoneyID   int //! 奖励货币
	MoneyNum  int //! 奖励货币数量
}

var GT_RebelActionAwardLst []ST_RebelActionAward

func InitRebelActionAwardParser(total int) bool {
	//GT_RebelActionAwardLst = make([]ST_RebelActionAward, total+1)
	return true
}

func ParseRebelActionAwardRecord(rs *RecordSet) {
	actionAward := ST_RebelActionAward{}
	action := CheckAtoi(rs.Values[0], 0)
	actionAward.Action = action
	actionAward.Diffculty = rs.GetFieldInt("diffculty")
	actionAward.MoneyID = rs.GetFieldInt("moneyid")
	actionAward.MoneyNum = rs.GetFieldInt("moneynum")
	GT_RebelActionAwardLst = append(GT_RebelActionAwardLst, actionAward)
}

func GetRebelActionAward(action int, diffculty int) *ST_ItemData {
	for _, v := range GT_RebelActionAwardLst {
		if v.Action == action && v.Diffculty == diffculty {
			return &ST_ItemData{v.MoneyID, v.MoneyNum}
		}
	}
	return nil
}

type ST_Rebel_Activity struct {
	Activity  int //! 活动类型: 1->消耗征讨令减半  2->功勋翻倍
	BeginTime int //! 活动时间
	EndTime   int //! 活动结束时间
}

var GT_RebelActivityLst []ST_Rebel_Activity

func InitRebelActivityParser(total int) bool {
	GT_RebelActivityLst = make([]ST_Rebel_Activity, total+1)
	return true
}

func ParseRebelActivityRecord(rs *RecordSet) {
	activity := CheckAtoi(rs.Values[0], 0)
	GT_RebelActivityLst[activity].Activity = activity
	GT_RebelActivityLst[activity].BeginTime = rs.GetFieldInt("begintime")
	GT_RebelActivityLst[activity].EndTime = rs.GetFieldInt("endtime")
}

func GetRebelOpenActivity(sec int) int {
	for _, v := range GT_RebelActivityLst {
		if sec >= v.BeginTime && sec < v.EndTime {
			return v.Activity
		}
	}

	return 0
}

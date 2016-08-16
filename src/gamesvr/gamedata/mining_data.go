package gamedata

import (
	"gamelog"
	"math/rand"
	"time"
)

//! 地图元素
const (
	MiningElement_Lower_Refining_Stone        = iota + 1 //! 低级精炼石
	MiningElement_Intermediate_Refining_Stone            //! 中级精炼石
	MiningElement_Advanced_Refining_Stone                //! 高级精炼石
	MiningElement_Ultimate_Refining_Stone                //! 极品精炼石
	MiningElement_Can_Not_Break_Obstacle                 //! 不可破坏障碍物
	MinintElement_Can_Break_Obstacle                     //! 可破坏障碍物
	MiningElement_Event                                  //! 事件

	//! 事
	MiningEvent_Action_Award   //! 行动力奖励
	MiningEvent_Black_Market   //! 黑市
	MiningEvent_Normal_Monster //! 普通怪
	MiningEvent_Elite_Monster  //! 精英怪
	MiningEvent_Boss           //! BOSS
	MiningEvent_Treasure       //! 宝箱
	MiningEvent_MagicBox       //! 魔盒
	MiningEvent_Scanning       //! 扫描
	MiningEvent_Question       //! 答题
	MiningEvent_Buff           //! Buff
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandNum(num int) int {
	return r.Intn(num)
}

//! 挖矿元素表
type ST_MiningElement struct {
	Element     int //! 元素
	Probability int //! 随机概率
	ItemID      int //! 物品ID
}

//! 精炼石随机表
type ST_MiningStoneNumberRandom struct {
	ItemID   int
	MinNum   int
	MaxNum   int
	MinLevel int
	MaxLevel int
}

//! 怪物随机表
type ST_MiningMonster struct {
	ID          int
	Event       int //! 事件ID
	CopyID      int //! 战斗副本ID
	MinLevel    int //! 等级区间取值
	MaxLevel    int
	MonsterLife int //! 怪物血量
}

var (
	GT_MiningStoneNumberRandomLst []ST_MiningStoneNumberRandom //! 矿洞精炼石随机表
	GT_MiningElementLst           []ST_MiningElement           //! 矿洞元素表
	GT_MiningEventLst             []ST_MiningEvent             //! 矿洞事件表
	GT_MiningEventBlackMarketLst  []ST_MiningEvent_BlackMarket //! 矿洞事件-黑市
	GT_MiningEventQuestionLst     []ST_MiningEvent_Question    //! 矿洞事件-答题
	GT_MiningEventBuffLst         []ST_MiningEvent_Buff        //! 矿洞事件-Buff
	GT_MiningEventTreasureLst     []ST_MiningEvent_Treasure    //! 矿洞事件-宝箱
	GT_MiningMonsterLst           []ST_MiningMonster           //! 矿洞事件-怪物

	GT_MiningAward []ST_MiningAward //! 矿洞结算奖励表
	GT_MiningGuaJi []ST_MiningGuaJi //! 矿洞挂机表
)

//! 矿洞精炼石随机表
func InitMiningStoneRandomParser(total int) bool {
	return true
}

func ParseMiningStoneRecord(rs *RecordSet) {
	var random ST_MiningStoneNumberRandom

	random.ItemID = rs.GetFieldInt("itemid")
	random.MinNum = rs.GetFieldInt("minnum")
	random.MaxNum = rs.GetFieldInt("maxnum")
	random.MaxLevel = rs.GetFieldInt("maxlevel")
	random.MinLevel = rs.GetFieldInt("minlevel")

	GT_MiningStoneNumberRandomLst = append(GT_MiningStoneNumberRandomLst, random)
}

//! 随机精炼石个数
func RandStoneNum(itemID int, level int) int {
	for _, v := range GT_MiningStoneNumberRandomLst {
		if level >= v.MinLevel && level <= v.MaxLevel {
			if v.ItemID == itemID && v.MaxNum != v.MinNum {
				randNum := GetRandNum(v.MaxNum - v.MinNum)
				return v.MinNum + randNum
			} else {
				return v.MaxNum
			}
		}

	}
	return 1
}

//! 矿洞元素表
func InitMiningElementParser(total int) bool {
	GT_MiningElementLst = make([]ST_MiningElement, total+1)
	return true
}

func ParseMiningElementRecord(rs *RecordSet) {
	element := CheckAtoi(rs.Values[0], 0)
	GT_MiningElementLst[element].Element = element
	GT_MiningElementLst[element].Probability = rs.GetFieldInt("probability")
	GT_MiningElementLst[element].ItemID = rs.GetFieldInt("itemid")
}

//! 获取矿洞元素信息
func GetMiningElementInfo(id int) (element *ST_MiningElement) {
	if id >= len(GT_MiningElementLst) || id <= 0 {
		gamelog.Error("GetMiningElementInfoLst fail. Invalid id: %d", id)
		return nil
	}

	element = &GT_MiningElementLst[id]
	return element
}

//! 随机一个元素
func RandMiningElement() (element int) {
	randValue := GetRandNum(1000)
	curPro := 0
	for _, v := range GT_MiningElementLst {
		if randValue >= curPro && randValue < curPro+v.Probability {
			return v.Element
		}

		curPro += v.Probability
	}
	return 0
}

//! 挖矿事件表
type ST_MiningEvent struct {
	Event       int //! 事件
	Probability int //! 随机概率
	Value1      int //! 配置值
	Value2      int //! 配置值  扫描 -> Value1 = 长 Value2 = 宽
}

func InitMiningEventParser(total int) bool {
	GT_MiningEventLst = make([]ST_MiningEvent, total+1)
	return true
}

func ParserMiningEventRecord(rs *RecordSet) {
	event := CheckAtoi(rs.Values[0], 0)
	GT_MiningEventLst[event].Event = event
	GT_MiningEventLst[event].Probability = rs.GetFieldInt("probability")
	GT_MiningEventLst[event].Value1 = rs.GetFieldInt("value1")
	GT_MiningEventLst[event].Value2 = rs.GetFieldInt("value2")
}

//! 获取矿洞事件信息
func GetMiningEventInfo(id int) (event *ST_MiningEvent) {

	id = id - MiningElement_Event
	if id > len(GT_MiningEventLst)-1 {
		gamelog.Error("GetMiningEventInfoLst fail. Invalid id: %d", id)
		return nil
	}
	event = &GT_MiningEventLst[id]
	return event
}

//! 随机一个事件
func RandMimingEvent() int {
	randValue := GetRandNum(1000)
	curPro := 0
	for _, v := range GT_MiningEventLst {
		if randValue >= curPro && randValue < curPro+v.Probability {
			return v.Event + MiningElement_Event
		}
		curPro += v.Probability
	}
	return 0
}

//! 挖矿事件-地下商店
type ST_MiningEvent_BlackMarket struct {
	ID        int
	ItemID    int //! 物品ID
	ItemNum   int //! 物品数量
	MoneyID   int //! 货币ID
	MoneyNum  int //! 价格
	NeedLevel int //! 需求等级
	Point     int //! 增加积分
	Discount  int //! 折扣
}

func InitMiningEventBlackMarketParser(total int) bool {
	GT_MiningEventBlackMarketLst = make([]ST_MiningEvent_BlackMarket, total+1)
	return true
}

func ParserMiningEventBlackMarketRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_MiningEventBlackMarketLst[id].ID = id
	GT_MiningEventBlackMarketLst[id].ItemID = rs.GetFieldInt("itemid")
	GT_MiningEventBlackMarketLst[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_MiningEventBlackMarketLst[id].MoneyID = rs.GetFieldInt("moneyid")
	GT_MiningEventBlackMarketLst[id].MoneyNum = rs.GetFieldInt("moneynum")
	GT_MiningEventBlackMarketLst[id].Point = rs.GetFieldInt("value")
	GT_MiningEventBlackMarketLst[id].NeedLevel = rs.GetFieldInt("needlevel")
	GT_MiningEventBlackMarketLst[id].Discount = rs.GetFieldInt("discount")
}

func RandBlackMarketGoosLst(num int, level int) []int {

	goodLst := []int{}
	for i := 0; i < num; i++ {
		randIndex := GetRandNum(len(GT_MiningEventBlackMarketLst))
		if GT_MiningEventBlackMarketLst[randIndex].NeedLevel < level {
			i -= 1
			continue
		}

		isExist := false
		for _, v := range goodLst {
			if GT_MiningEventBlackMarketLst[randIndex].ID == v {
				isExist = true
				break
			}
		}

		if isExist == true {
			i -= 1
			continue
		}

		goodLst = append(goodLst, GT_MiningEventBlackMarketLst[randIndex].ID)
	}
	return goodLst
}

//! 获取矿洞黑市信息
func GetMiningEventBlackMarketInfo(id int) (goods *ST_MiningEvent_BlackMarket) {
	if id > len(GT_MiningEventBlackMarketLst)-1 {
		gamelog.Error("GetMiningEventBlackMarketInfo fail. Invalid id: %d", id)
		return nil
	}
	goods = &GT_MiningEventBlackMarketLst[id]
	return goods
}

//! 挖矿事件-答题
type ST_MiningEvent_Question struct {
	QuestionID int       //! 问题ID
	Question   string    //! 问题内容
	Option     [4]string //! 问题选项
	Answer     int       //! 正确答案
}

func InitMiningEventQuestionParser(total int) bool {
	GT_MiningEventQuestionLst = make([]ST_MiningEvent_Question, total+1)
	return true
}

func ParserMiningEventQuestionRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_MiningEventQuestionLst[id].QuestionID = id
}

//! 获取矿洞行动力奖励信息
func GetMiningEventQuestionInfo(id int) (question *ST_MiningEvent_Question) {
	if id > len(GT_MiningEventQuestionLst)-1 {
		gamelog.Error("GetMiningEventQuestionInfo fail. Invalid id: %d", id)
		return nil
	}
	question = &GT_MiningEventQuestionLst[id]
	return question
}

//! 挖矿事件-Buff
type ST_MiningEvent_Buff struct {
	BuffType int //! 1->战力翻倍 2->资源翻倍 3->积分翻倍
	Value    int //! 倍数亦或其他
	Times    int //! 持续次数
}

func InitMiningEventBuffParser(total int) bool {
	return true
}

func ParserMiningEventBuffRecord(rs *RecordSet) {
	buffType := CheckAtoi(rs.Values[0], 0)

	var buff ST_MiningEvent_Buff
	buff.BuffType = buffType
	buff.Value = rs.GetFieldInt("value")
	buff.Times = rs.GetFieldInt("time")
	GT_MiningEventBuffLst = append(GT_MiningEventBuffLst, buff)
}

//! 获取矿洞Buff信息
func GetMiningEventBuffInfo(id int) (buff *ST_MiningEvent_Buff) {
	if id > len(GT_MiningEventBuffLst) {
		gamelog.Error("GetMiningEventBuffInfo fail. Invalid id: %d", id)
		return nil
	}
	buff = &GT_MiningEventBuffLst[id]
	return buff
}

//! 随机一个Buff
func RandMiningEventBuff() *ST_MiningEvent_Buff {
	return &GT_MiningEventBuffLst[GetRandNum(len(GT_MiningEventBuffLst))]
}

//! 挖矿事件-宝箱
type ST_MiningEvent_Treasure struct {
	Award       int //! 奖励ID
	Probability int //! 随机概率
}

func InitMiningEventTreasureParser(total int) bool {
	GT_MiningEventTreasureLst = make([]ST_MiningEvent_Treasure, total+1)
	return true
}

func ParserMiningEventTreasureRecord(rs *RecordSet) {
	award := CheckAtoi(rs.Values[0], 0)
	var awardTreasure ST_MiningEvent_Treasure
	awardTreasure.Award = award
	awardTreasure.Probability = rs.GetFieldInt("probability")

	GT_MiningEventTreasureLst = append(GT_MiningEventTreasureLst, awardTreasure)
}

//! 随机一个宝箱奖励
func RandMiningTreasure() int {
	randValue := GetRandNum(1000)
	curPro := 0
	for _, v := range GT_MiningEventTreasureLst {
		if randValue >= curPro && randValue < curPro+v.Probability {
			return v.Award
		}
		curPro += v.Probability
	}
	return 0
}

func GetMiningEventTreasureInfo(id int) (treasure *ST_MiningEvent_Treasure) {
	if id > len(GT_MiningEventTreasureLst)-1 {
		gamelog.Error("GetMiningEventTreasureInfo fail. Invalid id: %d", id)
		return nil
	}
	treasure = &GT_MiningEventTreasureLst[id]
	return treasure
}

func InitMiningEventMonsterPerser(total int) bool {
	GT_MiningMonsterLst = make([]ST_MiningMonster, total+1)
	return true
}

func ParseMiningEventMonsterRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_MiningMonsterLst[id].ID = id
	GT_MiningMonsterLst[id].Event = rs.GetFieldInt("event")
	GT_MiningMonsterLst[id].CopyID = rs.GetFieldInt("copyid")
	GT_MiningMonsterLst[id].MonsterLife = rs.GetFieldInt("life")
	GT_MiningMonsterLst[id].MaxLevel = rs.GetFieldInt("maxlevel")
	GT_MiningMonsterLst[id].MinLevel = rs.GetFieldInt("minlevel")
}

//! 随机矿洞怪物
func RandMiningMonster(event int, level int) int {
	for _, v := range GT_MiningMonsterLst {
		if v.Event == event {
			if level >= v.MinLevel && level <= v.MaxLevel {
				return v.ID
			}
		}
	}

	gamelog.Error("RandMiningMonster Error")
	return 0
}

func GetMiningMonsterLife(id int) int {
	if id > len(GT_MiningMonsterLst)-1 {
		gamelog.Error("GetMiningMonsterLife fail. Invalid id: %d", id)
		return 0
	}

	return GT_MiningMonsterLst[id].MonsterLife
}

func GetMonsterEventInfo(id int) *ST_MiningMonster {
	if id > len(GT_MiningMonsterLst)-1 {
		gamelog.Error("GetMiningMonsterLife fail. Invalid id: %d", id)
		return nil
	}
	return &GT_MiningMonsterLst[id]
}

//! 挖矿随机表
type ST_MiningRand struct {
	Value  int //! 数量
	Weight int //! 权重
}

var GT_MiningRand [2][]ST_MiningRand

func InitMiningRandParser(total int) bool {
	return true
}

func ParserMiningRandRecord(rs *RecordSet) {
	var award ST_MiningRand
	awardType := rs.GetFieldInt("type")
	award.Value = rs.GetFieldInt("value")
	award.Weight = rs.GetFieldInt("weight")
	GT_MiningRand[awardType-1] = append(GT_MiningRand[awardType-1], award)
}

//! 随机奖励 1->矿石数量 2->魔盒
func MiningRandAward(awardType int) int {
	totalWeight := 0
	for _, v := range GT_MiningRand[awardType-1] {
		totalWeight += v.Weight
	}

	//! 开始随机
	curWeight := 0
	randWeight := GetRandNum(totalWeight)
	for _, v := range GT_MiningRand[awardType-1] {
		if randWeight >= curWeight && randWeight < curWeight+v.Weight {
			return v.Value
		}

		curWeight += v.Weight
	}

	gamelog.Error("MiningRandAward error: Rand award fail type: %d", awardType)
	return 0
}

//! 挖矿奖励结算表
type ST_MiningAward struct {
	ID      int
	ItemID  int
	ItemNum int
}

func InitMiningAwardParser(total int) bool {
	// GT_MiningAward = make([]ST_MiningAward, total+1)
	return true
}

func ParserMiningAwardRecord(rs *RecordSet) {
	var award ST_MiningAward
	award.ID = rs.GetFieldInt("id")
	award.ItemID = rs.GetFieldInt("itemid")
	award.ItemNum = rs.GetFieldInt("itemnum")
	GT_MiningAward = append(GT_MiningAward, award)
}

func RandMiningAward(num int) []ST_MiningAward {
	if len(GT_MiningAward) <= 0 {
		return []ST_MiningAward{}
	}

	awardLst := []ST_MiningAward{}

	for i := 0; i < num; i++ {
		randIndex := GetRandNum(len(GT_MiningAward))
		awardLst = append(awardLst, GT_MiningAward[randIndex])
	}

	gamelog.Info("RandAward: %v", awardLst)
	return awardLst
}

//! 挖矿施工团队表
type ST_MiningGuaJi struct {
	ID            int
	Hour          int //! 挂机时间
	CostActionNum int //! 耗费行动力数量
	Award         int //! 奖励物品
}

func InitMiningGuaJiParser(total int) bool {
	GT_MiningGuaJi = make([]ST_MiningGuaJi, total+1)
	return true
}

func ParserMiningGuaJiRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_MiningGuaJi[id].ID = id
	GT_MiningGuaJi[id].Hour = rs.GetFieldInt("hour")
	GT_MiningGuaJi[id].CostActionNum = rs.GetFieldInt("costactionnum")
	GT_MiningGuaJi[id].Award = rs.GetFieldInt("award")
}

func GetMiningGuajiInfo(id int) *ST_MiningGuaJi {
	if id > len(GT_MiningGuaJi)-1 {
		gamelog.Error("GetMiningGuajiInfo fail. Invalid id: %d", id)
		return nil
	}

	return &GT_MiningGuaJi[id]
}

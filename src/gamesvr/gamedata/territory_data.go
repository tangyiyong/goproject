package gamedata

import (
	"gamelog"
)

type ST_TerritoryData struct {
	ID         int //! 领地ID
	CopyID     int //! 挑战副本ID
	FightValue int //! 推荐战力
}

type ST_TerritoryPatrolType struct {
	Type       int
	Time       int
	ActionType int
	ActionNum  int
}

type ST_TerritorySkillData struct {
	Level         int //! 技能等级
	SkillOpenTime int //! 需求时间
	CostMoneyID   int //! 花费货币ID
	CostMoneyNum  int //! 花费金额
	DoublePro     int //! 翻倍概率
}

type ST_TerritoryAwardData struct {
	ID      int
	ItemID  int
	ItemNum int
}

type ST_TerritoryAwardType struct {
	Type   int
	Time   int
	FuncID int
}

var (
	GT_TerritoryList          []ST_TerritoryData         //! 领地表
	GT_TerritorySkillList     [][5]ST_TerritorySkillData //! 领地技能表
	GT_TerritoryAwardList     [][]ST_TerritoryAwardData  //! 领地奖励表
	GT_TerritoryPatrolList    []ST_TerritoryPatrolType   //! 领地巡逻表
	GT_TerritoryAwardTypeList []ST_TerritoryAwardType    //! 领地奖励类型表
)

func InitTerritoryParser(total int) bool {
	GT_TerritoryList = make([]ST_TerritoryData, total+1)
	return true
}

func InitTerritoryAwardParser(total int) bool {
	if len(GT_TerritoryList) <= 0 {
		gamelog.Error("GT_TerritoryList is nil")
		return false
	}

	GT_TerritoryAwardList = make([][]ST_TerritoryAwardData, len(GT_TerritoryList))
	return true
}

func InitTerritorySkillParser(total int) bool {
	if len(GT_TerritoryList) <= 0 {
		gamelog.Error("GT_TerritoryList is nil")
		return false
	}

	GT_TerritorySkillList = make([][5]ST_TerritorySkillData, len(GT_TerritoryList))
	return true
}

func ParseTerritoryRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_TerritoryList[id].ID = id

	GT_TerritoryList[id].CopyID = rs.GetFieldInt("copyid")
}

func ParseTerritorySkillRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)

	var skillData ST_TerritorySkillData
	skillData.Level = rs.GetFieldInt("level")
	skillData.SkillOpenTime = rs.GetFieldInt("skillopentime")
	skillData.CostMoneyID = rs.GetFieldInt("costmoneyid")
	skillData.CostMoneyNum = rs.GetFieldInt("costmoneynum")
	skillData.DoublePro = rs.GetFieldInt("doublepro")

	GT_TerritorySkillList[id][skillData.Level-1] = skillData
}

func ParseTerritoryAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)

	var awardData ST_TerritoryAwardData
	awardData.ItemID = rs.GetFieldInt("itemid")
	awardData.ItemNum = rs.GetFieldInt("itemnum")

	GT_TerritoryAwardList[id] = append(GT_TerritoryAwardList[id], awardData)
}

//! 获取领地信息
func GetTerritoryData(id int) *ST_TerritoryData {
	if id > len(GT_TerritoryList)-1 {
		gamelog.Error("GetTerritoryData invlid id: %d", id)
		return nil
	}
	return &GT_TerritoryList[id]
}

//! 随机一个领地奖励
func RandTerritoryAward(id int) *ST_TerritoryAwardData {
	randIndex := r.Intn(len(GT_TerritoryAwardList[id]))
	return &GT_TerritoryAwardList[id][randIndex]
}

//! 获取领地技能信息
func GetTerritorySkillData(id int, level int) *ST_TerritorySkillData {
	if id >= len(GT_TerritoryList) || id <= 0 {
		gamelog.Error("GetTerritorySkillData invlid id: %d", id)
		return nil
	}
	return &GT_TerritorySkillList[id][level-1]
}

//! 初始化领地巡逻表
func InitTerritoryPatrolParser(total int) bool {
	GT_TerritoryPatrolList = make([]ST_TerritoryPatrolType, total+1)
	return true
}

func ParseTerritoryPatrolRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_TerritoryPatrolList[id].Type = id
	GT_TerritoryPatrolList[id].Time = rs.GetFieldInt("time")
	GT_TerritoryPatrolList[id].ActionType = rs.GetFieldInt("actiontype")
	GT_TerritoryPatrolList[id].ActionNum = rs.GetFieldInt("actionnum")
}

func GetPatrolTypeInfo(id int) *ST_TerritoryPatrolType {
	if id > len(GT_TerritoryPatrolList)-1 {
		gamelog.Error("GetPatrolTypeInfo fail: invalid id: %d", id)
		return nil
	}

	return &GT_TerritoryPatrolList[id]
}

//! 初始化领地奖励类型表
func InitTerritoryAwardTypeParser(total int) bool {
	GT_TerritoryAwardTypeList = make([]ST_TerritoryAwardType, total+1)
	return true
}

func ParseTerritoryAwardTypeRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_TerritoryAwardTypeList[id].Type = id
	GT_TerritoryAwardTypeList[id].Time = rs.GetFieldInt("time")
	GT_TerritoryAwardTypeList[id].FuncID = rs.GetFieldInt("func_id")
}

func GetTerritoryAwardType(id int) (int, int) {
	if id > len(GT_TerritoryPatrolList)-1 {
		gamelog.Error("GetTerritoryAwardType fail: invalid id: %d", id)
		return 0, 0
	}

	return GT_TerritoryAwardTypeList[id].Time, GT_TerritoryAwardTypeList[id].FuncID
}

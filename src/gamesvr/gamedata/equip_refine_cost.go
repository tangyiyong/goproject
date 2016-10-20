package gamedata

import (
	"gamelog"
)

type ST_EquipRefineCost struct {
	Level    int     //等级
	NeedExp  [10]int //依据品质
	TotalExp [10]int
}

type ST_ShenBinInfo struct {
	Pos        int
	Level      int //等级
	Propertyid int //属性ID
	PropertyV  int //属性值
	SkillID    int //神兵技能ID
}

const (
	SBST_DOUBLE_HURT   = 1 //双倍伤害
	SBST_RESTORE_HP    = 2 //恢复生命
	SBST_REBOUND_MAGIC = 3 //反弹法伤
	SBST_REBOUND_PHYIS = 4 //反弹物伤
	SBST_IGNORE_DEF    = 5 //无视防御
	SBST_SUCK_BLOOD    = 6 //吸血

)

type ST_ShenBinSkillInfo struct {
	ID    int //技能ID
	Type  int //技能类型
	Ratio int //概率
	Value int //数值
}

var (
	GT_EquipRefineCostList []ST_EquipRefineCost  = nil
	GT_ShenBinList         [7][51]ST_ShenBinInfo //最大精练等级装备50， 宝物20, 六个位置
	GT_ShenBinSkillList    []ST_ShenBinSkillInfo = nil
	GT_MinRefineLevel                            = 10000
)

func InitEquipRefineCostParser(total int) bool {
	GT_EquipRefineCostList = make([]ST_EquipRefineCost, total)

	return true
}

func ParseEquipRefineCostRecord(rs *RecordSet) {
	level := CheckAtoi(rs.Values[0], 0)
	GT_EquipRefineCostList[level].Level = level
	GT_EquipRefineCostList[level].NeedExp[0] = CheckAtoi(rs.Values[1], 1)
	GT_EquipRefineCostList[level].NeedExp[1] = CheckAtoi(rs.Values[2], 2)
	GT_EquipRefineCostList[level].NeedExp[2] = CheckAtoi(rs.Values[3], 3)
	GT_EquipRefineCostList[level].NeedExp[3] = CheckAtoi(rs.Values[4], 4)
	GT_EquipRefineCostList[level].NeedExp[4] = CheckAtoi(rs.Values[5], 5)
	GT_EquipRefineCostList[level].NeedExp[5] = CheckAtoi(rs.Values[6], 6)
	GT_EquipRefineCostList[level].NeedExp[6] = CheckAtoi(rs.Values[7], 7)
	GT_EquipRefineCostList[level].NeedExp[7] = CheckAtoi(rs.Values[8], 8)
	GT_EquipRefineCostList[level].NeedExp[8] = CheckAtoi(rs.Values[9], 9)
	GT_EquipRefineCostList[level].NeedExp[9] = CheckAtoi(rs.Values[10], 10)

	GT_EquipRefineCostList[level].TotalExp[0] = CheckAtoi(rs.Values[11], 11)
	GT_EquipRefineCostList[level].TotalExp[1] = CheckAtoi(rs.Values[12], 12)
	GT_EquipRefineCostList[level].TotalExp[2] = CheckAtoi(rs.Values[13], 13)
	GT_EquipRefineCostList[level].TotalExp[3] = CheckAtoi(rs.Values[14], 14)
	GT_EquipRefineCostList[level].TotalExp[4] = CheckAtoi(rs.Values[15], 15)
	GT_EquipRefineCostList[level].TotalExp[5] = CheckAtoi(rs.Values[16], 16)
	GT_EquipRefineCostList[level].TotalExp[6] = CheckAtoi(rs.Values[17], 17)
	GT_EquipRefineCostList[level].TotalExp[7] = CheckAtoi(rs.Values[18], 18)
	GT_EquipRefineCostList[level].TotalExp[8] = CheckAtoi(rs.Values[19], 19)
	GT_EquipRefineCostList[level].TotalExp[9] = CheckAtoi(rs.Values[20], 20)
	return
}

func GetEquipRefineCostInfo(level int) *ST_EquipRefineCost {
	if level >= len(GT_EquipRefineCostList) {
		gamelog.Error("GetEquipRefineCostInfo Error : Invalid level %d", level)
		return nil
	}

	return &GT_EquipRefineCostList[level]
}

func InitShenBinParser(total int) bool {
	return true
}

func ParseShenBinRecord(rs *RecordSet) {
	pos := rs.GetFieldInt("pos")
	level := rs.GetFieldInt("level")
	GT_ShenBinList[pos][level].Pos = pos
	GT_ShenBinList[pos][level].Level = level
	GT_ShenBinList[pos][level].Propertyid = rs.GetFieldInt("p_id")
	GT_ShenBinList[pos][level].PropertyV = rs.GetFieldInt("p_value")
	GT_ShenBinList[pos][level].SkillID = rs.GetFieldInt("skill_id")

	if GT_MinRefineLevel > level {
		GT_MinRefineLevel = level
	}
	return
}

func FinishShenBinParser() bool {
	for i := 1; i < 7; i++ {
		for j := 1; j < 51; j++ {
			if GT_ShenBinList[i][j].Level == 0 {
				GT_ShenBinList[i][j] = GT_ShenBinList[i][j-1]
			}
		}
	}

	return true
}

func GetShenBinInfo(pos int, level int) *ST_ShenBinInfo {
	if level < GT_MinRefineLevel {
		return nil
	}

	if pos <= 0 || pos > 6 || level > 50 {
		gamelog.Error("GetShenBinInfo Error : Invalid level %d and pos :%d", level, pos)
		return nil
	}

	return &GT_ShenBinList[pos][level]
}

func InitShenBinSkillParser(total int) bool {
	GT_ShenBinSkillList = make([]ST_ShenBinSkillInfo, total+1)
	return true
}

func ParseShenBinSkillRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_ShenBinSkillList[id].ID = id
	GT_ShenBinSkillList[id].Type = rs.GetFieldInt("type")
	GT_ShenBinSkillList[id].Ratio = rs.GetFieldInt("ratio")
	GT_ShenBinSkillList[id].Value = rs.GetFieldInt("value")
	return
}

package gamedata

import (
	"fmt"
	"gamelog"
	"strings"
)

type ST_EquipSuitBuff struct {
	PropertyID    int  //属性ID
	PropertyValue int  //属性值
	IsPercent     bool //是否百分比
}

type ST_EquipSuitInfo struct {
	SuitID int                 //套装ID
	Buffs  [4]ST_EquipSuitBuff //装备的Buff集

}

var (
	GT_EquipSuitList []ST_EquipSuitInfo = nil
)

func InitEquipSuitParser(total int) bool {

	GT_EquipSuitList = make([]ST_EquipSuitInfo, total+1)

	return true
}

func ParseEquipSuitRecord(rs *RecordSet) {
	SuitID := CheckAtoi(rs.Values[0], 0)

	GT_EquipSuitList[SuitID].SuitID = SuitID
	GT_EquipSuitList[SuitID].Buffs[0].PropertyID = rs.GetFieldInt("two_id")
	GT_EquipSuitList[SuitID].Buffs[0].PropertyValue = rs.GetFieldInt("two_value")
	GT_EquipSuitList[SuitID].Buffs[0].IsPercent = rs.GetFieldInt("two_percent") == 1

	GT_EquipSuitList[SuitID].Buffs[1].PropertyID = rs.GetFieldInt("three_id")
	GT_EquipSuitList[SuitID].Buffs[1].PropertyValue = rs.GetFieldInt("three_value")
	GT_EquipSuitList[SuitID].Buffs[1].IsPercent = rs.GetFieldInt("three_percent") == 1

	sv := strings.Split(rs.GetFieldString("four_id"), "|")
	if len(sv) <= 1 {
		panic(fmt.Sprintf("field: four_id Is wrong format!!!"))
	}

	GT_EquipSuitList[SuitID].Buffs[2].PropertyID = CheckAtoi(sv[0], 41)
	GT_EquipSuitList[SuitID].Buffs[3].PropertyID = CheckAtoi(sv[1], 42)

	sv = strings.Split(rs.GetFieldString("four_value"), "|")
	if len(sv) <= 1 {
		panic(fmt.Sprintf("field: four_value Is wrong format!!!"))
	}
	GT_EquipSuitList[SuitID].Buffs[2].PropertyValue = CheckAtoi(sv[0], 51)
	GT_EquipSuitList[SuitID].Buffs[3].PropertyValue = CheckAtoi(sv[1], 52)

	GT_EquipSuitList[SuitID].Buffs[2].IsPercent = (rs.GetFieldInt("four_percent") == 1)
	GT_EquipSuitList[SuitID].Buffs[3].IsPercent = GT_EquipSuitList[SuitID].Buffs[2].IsPercent

	return
}

func GetEquipSuitInfo(suitid int) *ST_EquipSuitInfo {
	if suitid >= len(GT_EquipSuitList) || suitid <= 0 {
		gamelog.Error("GetEquipSuitInfo Error : equipid suitid %d", suitid)
		return nil
	}

	return &GT_EquipSuitList[suitid]
}

func GetEquipSuitBuff(suitid int, equips int) []ST_EquipSuitBuff {
	pSuitInfo := GetEquipSuitInfo(suitid)
	if pSuitInfo == nil {
		gamelog.Error("GetEquipSuitBuff Error : Invalid suitid %d", suitid)
		return nil
	}

	var ret []ST_EquipSuitBuff

	if equips == 2 {
		return pSuitInfo.Buffs[:1]
	} else if equips == 3 {
		return pSuitInfo.Buffs[:2]
	} else if equips == 4 {
		return pSuitInfo.Buffs[:4]
	}

	return ret
}

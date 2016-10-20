package gamedata

import (
	"gamelog"
	"strings"
)

type ST_PropertyChange struct {
	PropertyID  int //属性的ID
	PropertyInc int //属性增加值
}

type ST_MasterItem struct {
	Level int //大师等级

	EquipStrengthLevel     int //装备强化等级
	EquipStrengthPropertys []ST_PropertyChange

	EquipRefineLevel     int //装备精炼等级
	EquipRefinePropertys []ST_PropertyChange

	GemStrengthLevel     int //宝物强化等级
	GemStrengthPropertys []ST_PropertyChange

	GemRefineLevel     int //宝物精炼等级
	GemRefinePropertys []ST_PropertyChange
}

const (
	MTYPE_Equip_Strength = 1 //装备强化
	MTYPE_Equip_Refine   = 2 //装备精炼
	MTYPE_Gem_Strength   = 3 //宝物强化
	MTYPE_Gem_Refine     = 4 //宝物精炼
)

var (
	GT_Master_List []ST_MasterItem = nil
)

func InitMasterParser(total int) bool {
	GT_Master_List = make([]ST_MasterItem, total+1)
	return true
}

func ParseMasterRecord(rs *RecordSet) {
	Level := CheckAtoi(rs.Values[0], 0)
	GT_Master_List[Level].Level = Level
	GT_Master_List[Level].EquipStrengthLevel = rs.GetFieldInt("equip_strength_level")
	GT_Master_List[Level].EquipRefineLevel = rs.GetFieldInt("equip_refine_level")
	GT_Master_List[Level].GemStrengthLevel = rs.GetFieldInt("gem_strength_level")
	GT_Master_List[Level].GemRefineLevel = rs.GetFieldInt("gem_refine_level")
	GT_Master_List[Level].EquipStrengthPropertys = ParseToPropertyChanges(rs.GetFieldString("equip_strength_propertys"))
	GT_Master_List[Level].EquipRefinePropertys = ParseToPropertyChanges(rs.GetFieldString("equip_refine_propertys"))
	GT_Master_List[Level].GemStrengthPropertys = ParseToPropertyChanges(rs.GetFieldString("gem_strength_propertys"))
	GT_Master_List[Level].GemRefinePropertys = ParseToPropertyChanges(rs.GetFieldString("gem_refine_propertys"))
}

func ParseToPropertyChanges(svalue string) (ret []ST_PropertyChange) {
	if svalue == "NULL" {
		return nil
	}

	var item ST_PropertyChange
	var sFix = strings.Trim(svalue, "()")
	slice := strings.Split(sFix, ")(")
	for i := 0; i < len(slice); i++ {
		pv := strings.Split(slice[i], "|")
		item.PropertyID = CheckAtoi(pv[0], 99)
		item.PropertyInc = CheckAtoi(pv[1], 99)
		ret = append(ret, item)
	}

	return
}

func GetMasterInfo(masterType int, level int) []ST_PropertyChange {
	var validLevel = 0
	for _, item := range GT_Master_List {
		if masterType == MTYPE_Equip_Strength {
			if level > item.EquipStrengthLevel {
				validLevel = item.Level
			} else {
				break
			}
		} else if masterType == MTYPE_Equip_Refine {
			if level > item.EquipRefineLevel {
				validLevel = item.Level
			} else {
				break
			}
		} else if masterType == MTYPE_Gem_Strength {
			if level > item.GemStrengthLevel {
				validLevel = item.Level
			} else {
				break
			}
		} else if masterType == MTYPE_Gem_Refine {
			if level > item.GemRefineLevel {
				validLevel = item.Level
			} else {
				break
			}
		} else {
			gamelog.Error("GetMaterInfo Error ")
		}
	}

	if masterType == MTYPE_Equip_Strength {
		return GT_Master_List[validLevel].EquipStrengthPropertys
	} else if masterType == MTYPE_Equip_Refine {
		return GT_Master_List[validLevel].EquipRefinePropertys
	} else if masterType == MTYPE_Gem_Strength {
		return GT_Master_List[validLevel].GemStrengthPropertys
	} else if masterType == MTYPE_Gem_Refine {
		return GT_Master_List[validLevel].GemRefinePropertys
	}

	return nil
}

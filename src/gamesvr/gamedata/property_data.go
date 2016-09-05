package gamedata

//属性配制表
//配制了所有属性

import (
	"gamelog"
)

type ST_PropertyInfo struct {
	ID          int   //属性ID
	FightFactor int32 //战力因子
}

var GT_Property_List []ST_PropertyInfo = nil

func InitPropertyParser(total int) bool {
	GT_Property_List = make([]ST_PropertyInfo, total+1)
	return true
}

func ParsePropertyRecord(rs *RecordSet) {
	propertyid := rs.GetFieldInt("id")
	if propertyid >= 20 {
		return
	}

	GT_Property_List[propertyid].ID = propertyid
	GT_Property_List[propertyid].FightFactor = int32(rs.GetFieldInt("fight_factor"))
}

func GetPropertyCount() int {
	return len(GT_Property_List) - 1
}

func GetPropertyInfo(propertyid int) *ST_PropertyInfo {
	if propertyid >= len(GT_Property_List) || propertyid == 0 {
		gamelog.Error("GetPropertyInfo Error: invalid propertyid :%d", propertyid)
		return nil
	}

	return &GT_Property_List[propertyid]
}

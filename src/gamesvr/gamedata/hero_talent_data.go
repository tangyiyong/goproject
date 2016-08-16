package gamedata

import (
	"gamelog"
)

const (
	TargetType_Self      = 1 //1:自己增加属性
	TargetType_Friend    = 2 //2:友军增加属性
	TargetType_Camp      = 3 //3:阵营增加属性
	TargetType_Camp_Kill = 4 //4:全体灭阵营
)

type ST_TalentItem struct {
	TalentID       int  //天赋ID
	TargetType     int  //1: 自己 2:友军  3: 阵营
	TargetCamp     int  //目标阵营
	PropertyID     int  //属性ID
	PropertyValue1 int  //普通值
	PropertyValue2 int  //神将天赋值
	IsPercent      bool //是否是百分比值
}

var (
	GT_Talent_List []ST_TalentItem = nil
)

func InitTalentParser(total int) bool {
	GT_Talent_List = make([]ST_TalentItem, total+1)
	return true
}

func ParseTalentRecord(rs *RecordSet) {
	TalentID := rs.GetFieldInt("id")
	GT_Talent_List[TalentID].TalentID = TalentID
	GT_Talent_List[TalentID].TargetType = rs.GetFieldInt("target")
	GT_Talent_List[TalentID].TargetCamp = rs.GetFieldInt("camp")
	GT_Talent_List[TalentID].PropertyID = rs.GetFieldInt("propertyid")
	GT_Talent_List[TalentID].PropertyValue1 = rs.GetFieldInt("propertyvalue1")
	GT_Talent_List[TalentID].PropertyValue2 = rs.GetFieldInt("propertyvalue2")
	GT_Talent_List[TalentID].IsPercent = (1 == rs.GetFieldInt("is_percent"))
}

func GetTalentInfo(talentID int) *ST_TalentItem {
	if talentID >= len(GT_Talent_List) || talentID <= 0 {
		gamelog.Error("GetTalentInfo Error: invalid talentID :%d", talentID)
		return nil
	}

	return &GT_Talent_List[talentID]
}

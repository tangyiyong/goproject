package gamedata

import (
	"gamelog"
)

//! 豪华签到奖励
type ST_SignPlusInfo struct {
	ID        int //! 唯一标识
	SignAward int //! 豪华签到奖励 取值Award表
	MinLevel  int //! 豪华签到等级范围
	MaxLevel  int //! 豪华签到等级范围
}

var GT_Sign_Plus_List []ST_SignPlusInfo = nil

//! 初始化签到分析器
func InitSignPlusParser(total int) bool {
	GT_Sign_Plus_List = make([]ST_SignPlusInfo, total+1)
	return true
}

//! 分析CSV
func ParseSignPlusRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_Sign_Plus_List[id].ID = id
	GT_Sign_Plus_List[id].SignAward = rs.GetFieldInt("signaward")
	GT_Sign_Plus_List[id].MinLevel = rs.GetFieldInt("minlevel")
	GT_Sign_Plus_List[id].MaxLevel = rs.GetFieldInt("maxlevel")
}

//! 提供接口
func GetSignPlusDataFromLevel(level int) *ST_SignPlusInfo {
	if level <= 0 {
		gamelog.Error("GetSignPlusDataFromLevel Error : Invalid level:%d", level)
		return nil
	}

	for i, v := range GT_Sign_Plus_List {
		if level >= v.MinLevel && level < v.MaxLevel {
			return &GT_Sign_Plus_List[i]
		}
	}
	return nil
}

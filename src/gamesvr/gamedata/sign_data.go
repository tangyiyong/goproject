package gamedata

import (
	"gamelog"
)

//! 签到奖励
type ST_SignInfo struct {
	DayID     int  //! 天数
	AwardItem int  //! 奖品
	Count     int  //! 奖励个数
	VipLevel  int8 //! 倍数奖励需求VIP等级
	Multiple  int  //! 领取奖励倍数
	Type      int  //! 种类
}

var GT_Sign_List []ST_SignInfo = nil

//! 初始化签到分析器
func InitSignParser(total int) bool {
	GT_Sign_List = make([]ST_SignInfo, total+1)
	return true
}

//! 分析CSV
func ParseSignRecord(rs *RecordSet) {
	dayID := rs.GetFieldInt("dayid")
	GT_Sign_List[dayID].DayID = dayID
	GT_Sign_List[dayID].AwardItem = rs.GetFieldInt("awarditem")
	GT_Sign_List[dayID].Count = rs.GetFieldInt("count")
	GT_Sign_List[dayID].VipLevel = int8(rs.GetFieldInt("viplevel"))
	GT_Sign_List[dayID].Multiple = rs.GetFieldInt("multiple")
	GT_Sign_List[dayID].Type = rs.GetFieldInt("type")
}

//! 提供接口
func GetSignData(day int) *ST_SignInfo {
	if day >= len(GT_Sign_List) || day <= 0 {
		gamelog.Error("GetSignData Error: invalid day %d", day)
		return nil
	}

	return &GT_Sign_List[day]
}

func GetSignAwardCount() int {
	return len(GT_Sign_List) - 1
}

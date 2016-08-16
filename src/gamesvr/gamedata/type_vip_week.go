package gamedata

import (
	"gamelog"
)

//! VIP每周礼包表
type ST_VipWeekGiftInfo struct {
	ID              int
	Award           int //! 礼包ID
	MoneyID         int //! 货币ID
	MoneyNum        int //! 数量
	BuyTimes        int //! 购买次数
	Range_level_min int //! 取值等级小
	Range_level_max int //! 取值等级大
}

var GT_VipWeekGiftLst []ST_VipWeekGiftInfo

func InitVipWeekParser(total int) bool {
	GT_VipWeekGiftLst = make([]ST_VipWeekGiftInfo, total+1)
	return true
}

func ParseVipWeekRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_VipWeekGiftLst[id].ID = id
	GT_VipWeekGiftLst[id].Award = rs.GetFieldInt("award")
	GT_VipWeekGiftLst[id].MoneyID = rs.GetFieldInt("moneyid")
	GT_VipWeekGiftLst[id].MoneyNum = rs.GetFieldInt("moneynum")
	GT_VipWeekGiftLst[id].BuyTimes = rs.GetFieldInt("buytimes")
	GT_VipWeekGiftLst[id].Range_level_min = rs.GetFieldInt("range_level_min")
	GT_VipWeekGiftLst[id].Range_level_max = rs.GetFieldInt("range_level_max")
}

func GetVipWeekItem(level int) (giftLst []ST_VipWeekGiftInfo) {
	for _, v := range GT_VipWeekGiftLst {
		if level >= v.Range_level_min && level <= v.Range_level_max {
			giftLst = append(giftLst, v)
		}
	}

	if len(giftLst) == 0 {
		gamelog.Error("GetVipWeekItem nil, level: %d", giftLst)
		return giftLst
	}

	return giftLst
}

func GetVipWeekItemFromID(id int) *ST_VipWeekGiftInfo {
	if id >= len(GT_VipWeekGiftLst) || id <= 0 {
		gamelog.Error("GetVipWeekItemFromID Error: invalid id %d", id)
		return nil
	}

	return &GT_VipWeekGiftLst[id]
}

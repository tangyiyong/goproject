package gamedata

import (
	"gamelog"
)

type ST_ChargeItem struct {
	ID           int //充值ID
	Type         int //类型 1->月卡 2->普通充值 3->优惠充值
	RenMinBi     int //充值的人民币
	Diamond      int //获得钻石数
	ExtraAward   int //额外奖励
	FirstAwardID int //首冲奖励
	DayNum       int //持续天数
}

var (
	GT_ChargeItemList []ST_ChargeItem //充值项列表
)

func InitChargeItemParser(total int) bool {
	GT_ChargeItemList = make([]ST_ChargeItem, total+1)
	return true
}
func ParseChargeItemRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	data := &GT_ChargeItemList[id]
	data.ID = id
	data.Type = rs.GetFieldInt("type")
	data.RenMinBi = rs.GetFieldInt("renminbi")
	data.Diamond = rs.GetFieldInt("diamond")
	data.ExtraAward = rs.GetFieldInt("extra_award")
	data.FirstAwardID = rs.GetFieldInt("first_award_id")
	data.DayNum = rs.GetFieldInt("day_num")
}

func GetChargeItem(id int) *ST_ChargeItem {
	if id >= len(GT_ChargeItemList) || id <= 0 {
		gamelog.Error("GetChargeItem Error : Invalid id  %d", id)
		return nil
	}

	return &GT_ChargeItemList[id]
}

func GetChargeItemCount() int {
	return len(GT_ChargeItemList)
}
func GetMonthCardCount() (ret int) {
	for _, v := range GT_ChargeItemList {
		if v.Type == 1 {
			ret = v.ID + 1
		}
	}
	gamelog.Info("GetMonthCardCount: %d", ret)
	return
}

//! 获取优惠充值取值区间
func GetDiscountChargeIDSection() (int, int) {
	minID, maxID := 0, 0
	for _, v := range GT_ChargeItemList {
		if v.Type == 3 {
			if minID == 0 {
				minID = v.ID
			}

			if v.ID > maxID {
				maxID = v.ID
			} else if v.ID < minID {
				minID = v.ID
			}

		}
	}
	return minID, maxID
}

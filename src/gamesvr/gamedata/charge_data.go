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
	AwardID      int //常规奖励ID
}

type ST_MonthCard struct {
	ID       int //月卡ID
	RenMinBi int //充值的人民币
	MoneyID  int //货币ID
	MoneyNum int //货币数
	DayNum   int //天数
}

var (
	GT_ChargeItemList []ST_ChargeItem //充值项列表
	GT_MonthCardList  []ST_MonthCard  //月卡列表
)

func InitChargeItemParser(total int) bool {
	GT_ChargeItemList = make([]ST_ChargeItem, total+1)
	return true
}

func InitMonthCardParser(total int) bool {
	GT_MonthCardList = make([]ST_MonthCard, total+1)
	return true
}

func ParseChargeItemRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	data := &GT_ChargeItemList[id]
	data.ID = id
	data.Type = rs.GetFieldInt("type")
	data.RenMinBi = rs.GetFieldInt("renminbi")
	data.Diamond = rs.GetFieldInt("diamond")
	data.FirstAwardID = rs.GetFieldInt("first_award_id")
	data.AwardID = rs.GetFieldInt("award_id")
	data.ExtraAward = rs.GetFieldInt("extra_award")
}

func ParseMonthCardRecord(rs *RecordSet) {
	cardid := rs.GetFieldInt("card_id")
	data := &GT_MonthCardList[cardid]
	data.ID = cardid
	data.RenMinBi = rs.GetFieldInt("renminbi")
	data.MoneyID = rs.GetFieldInt("money_id")
	data.MoneyNum = rs.GetFieldInt("money_num")
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

func GetMonthCardInfo(cardid int) *ST_MonthCard {
	if cardid >= len(GT_MonthCardList) || cardid <= 0 {
		gamelog.Error("GetMonthCardInfo Error : Invalid cardid  %d", cardid)
		return nil
	}

	return &GT_MonthCardList[cardid]
}

func GetMonthCardCount() int {
	return len(GT_MonthCardList)
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

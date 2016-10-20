package gamedata

import (
	"gamelog"
)

const (
	Summon_Normal = 1 //! 普通招贤
	Summon_Senior = 2 //! 高级招贤
)

type ST_SummonConifgInfo struct {
	SummonType   int //! 召唤类型
	CostItemID   int //! 花费道具类型
	CostItemNum  int //! 花费道具数量
	CostMoneyID  int //! 花费金钱类型
	CostMoneyNum int //! 花费金钱数量
	NeedPoint    int //! 到达必出武将所需积分
}

var GT_Summon_Config_List []ST_SummonConifgInfo

func InitSummonConfigParser(total int) bool {
	GT_Summon_Config_List = make([]ST_SummonConifgInfo, total+1)
	return true
}

func ParseSummonConfigRecord(rs *RecordSet) {
	summonType := CheckAtoi(rs.Values[0], 0)
	GT_Summon_Config_List[summonType].SummonType = summonType
	GT_Summon_Config_List[summonType].CostItemID = rs.GetFieldInt("costitemid")
	GT_Summon_Config_List[summonType].CostItemNum = rs.GetFieldInt("costitemnum")
	GT_Summon_Config_List[summonType].CostMoneyID = rs.GetFieldInt("costmoneyid")
	GT_Summon_Config_List[summonType].CostMoneyNum = rs.GetFieldInt("costmoneynum")
	GT_Summon_Config_List[summonType].NeedPoint = rs.GetFieldInt("needpoint")

}

func GetSummonConfig(summonType int) *ST_SummonConifgInfo {
	if summonType > len(GT_Summon_Config_List) || summonType <= 0 {
		gamelog.Error("GetSummonConfig Error: invalid summontype :%d", summonType)
		return nil
	}

	return &GT_Summon_Config_List[summonType]
}

func GetSummonInfoRandom(summonType int, number int) []ST_ItemData {
	if number <= 0 {
		gamelog.Error("GetSummonInfoRandom Error: Invalid Number:%d", number)
		return nil
	}

	if summonType == Summon_Normal {
		normalHeroLst := GetItemsAwardIDTimes(NormalSummonAwardID, number)
		if normalHeroLst == nil || len(normalHeroLst) <= 0 {
			gamelog.Error("GetSummonInfoRandom Error: invalid NormalSummonAwardID :%d", NormalSummonAwardID)
		}
		return normalHeroLst
	} else if summonType == Summon_Senior {
		seniorHeroLst := GetItemsAwardIDTimes(SeniorSummonAwardID, number)
		if seniorHeroLst == nil || len(seniorHeroLst) <= 0 {
			gamelog.Error("GetSummonInfoRandom Error: invalid SeniorSummonAwardID :%d", SeniorSummonAwardID)
		}
		return seniorHeroLst
	}

	return nil
}

//! 橙将随机
func GetSummonInfoOrangeRandom() int {
	awardLst := GetItemsFromAwardID(OrangeSummonAwardID)
	if awardLst == nil || len(awardLst) <= 0 {
		gamelog.Error("GetSummonInfoOrangeRandom Error: invalid OrangeSummonAwardID :%d", OrangeSummonAwardID)
		return 0
	}
	return awardLst[0].ItemID
}

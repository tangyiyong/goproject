package gamedata

//阵营战水晶配制表

import (
	"gamelog"
)

type ST_CrystalInfo struct {
	CrystalID    int    //水晶ID
	MoneyID      [2]int //收获货币ID
	MoneyNum     [2]int //收获货币Num
	CostMoneyID  int    //花费货币ID
	CostMoneyNum int    //花费货币数
}

var GT_Crystal_List []ST_CrystalInfo = nil

func InitCrystalParser(total int) bool {
	GT_Crystal_List = make([]ST_CrystalInfo, total+1)
	return true
}

func ParseCrystalRecord(rs *RecordSet) {
	crystalID := rs.GetFieldInt("id")
	GT_Crystal_List[crystalID].CrystalID = crystalID
	GT_Crystal_List[crystalID].MoneyID[0] = rs.GetFieldInt("money_id_1")
	GT_Crystal_List[crystalID].MoneyID[1] = rs.GetFieldInt("money_id_2")
	GT_Crystal_List[crystalID].MoneyNum[0] = rs.GetFieldInt("money_num_1")
	GT_Crystal_List[crystalID].MoneyNum[1] = rs.GetFieldInt("money_num_2")
	GT_Crystal_List[crystalID].CostMoneyID = rs.GetFieldInt("cost_money_id")
	GT_Crystal_List[crystalID].CostMoneyNum = rs.GetFieldInt("cost_money_num")
}

func GetCrystalInfo(crystalID int) *ST_CrystalInfo {
	if crystalID >= len(GT_Crystal_List) || crystalID == 0 {
		gamelog.Error("GetCrystalInfo Error: invalid crystalID :%d", crystalID)
		return nil
	}

	return &GT_Crystal_List[crystalID]
}

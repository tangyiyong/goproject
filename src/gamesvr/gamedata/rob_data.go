package gamedata

import (
	"fmt"
	"gamelog"
)

//! 夺宝配置表
type ST_RobConfig struct {
	CopyID int //! 副本基本信息

	//! 抢夺NPC
	Quality [4]int //! 品质
	RobPro  [4]int //! 抢夺概率

	//! 抢夺玩家
	PlayerGeneralPro       int //! 抢劫玩家普通概率
	GeneralLevelDifference int //! 普通概率等级差
	PlayerHighPro          int //! 抢劫玩家高等概率
	HighLevelDifference    int //! 高等概率等级差
	PlayerLowPro           int //! 抢劫玩家低等概率
	LowLevelDifference     int //! 低等概率等级差
}

var RobConfig ST_RobConfig

func InitRobParser(total int) bool {
	return true
}

func ParseRobRecord(rs *RecordSet) {
	RobConfig.CopyID = rs.GetFieldInt("copy_id")

	for i := 1; i <= 4; i++ {
		filedName := fmt.Sprintf("quality%d", i)
		RobConfig.Quality[i-1] = rs.GetFieldInt(filedName)

		filedName = fmt.Sprintf("robpro%d", i)
		RobConfig.RobPro[i-1] = rs.GetFieldInt(filedName)
	}

	RobConfig.PlayerGeneralPro = rs.GetFieldInt("playergeneralpro")
	RobConfig.GeneralLevelDifference = rs.GetFieldInt("generalleveldifference")
	RobConfig.PlayerHighPro = rs.GetFieldInt("playerhighpro")
	RobConfig.HighLevelDifference = rs.GetFieldInt("highleveldifference")
	RobConfig.PlayerLowPro = rs.GetFieldInt("playerlowpro")
	RobConfig.LowLevelDifference = rs.GetFieldInt("lowleveldifference")
}

func GetRobConfig() *ST_RobConfig {
	return &RobConfig
}

//! 夺宝熔炼表
type ST_TreasureMelting struct {
	GemID        int
	CostMoneyID  int
	CostMoneyNum int
}

var GT_TreasureMeltingLst []ST_TreasureMelting

func InitTreasureMeltingParser(total int) bool {
	//	GT_TreasureMeltingLst = make([]ST_TreasureMelting, total+1)
	return true
}

func ParseTreasureMeltingRecord(rs *RecordSet) {
	var info ST_TreasureMelting
	info.GemID = rs.GetFieldInt("gemid")
	info.CostMoneyID = rs.GetFieldInt("cost_money_id")
	info.CostMoneyNum = rs.GetFieldInt("cost_money_num")
	GT_TreasureMeltingLst = append(GT_TreasureMeltingLst, info)
}

func GetTreasureMeltingInfo(gemid int) (int, int) {
	for _, v := range GT_TreasureMeltingLst {
		if v.GemID == gemid {
			return v.CostMoneyID, v.CostMoneyNum
		}
	}

	gamelog.Error("GetTreasureMeltingInfo Error: Invalid gemid %d", gemid)
	return 0, 0
}

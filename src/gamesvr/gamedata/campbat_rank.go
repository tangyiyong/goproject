package gamedata

import (
	"gamelog"
)

//! 阵营战排行奖励
type ST_CampBatRank struct {
	ID             int //! 唯一标识
	CampKillID     int //! 阵营榜奖励
	CampDestroyID  int //阵营团灭
	TotalKillID    int //! 总表奖励
	TotalDestroyID int //! 总表奖励
	MinLevel       int //! 等级范围
	MaxLevel       int //! 等级范围
}

var GT_CampBatRank_List []ST_CampBatRank = nil

//! 阵营战排行奖励
func InitCampBatRankParser(total int) bool {
	GT_CampBatRank_List = make([]ST_CampBatRank, total+1)
	return true
}

//! 分析CSV
func ParseCampBatRankRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_CampBatRank_List[id].ID = id
	GT_CampBatRank_List[id].CampKillID = rs.GetFieldInt("camp_kill_rank")
	GT_CampBatRank_List[id].CampDestroyID = rs.GetFieldInt("camp_destroy_rank")
	GT_CampBatRank_List[id].TotalKillID = rs.GetFieldInt("total_kill_rank")
	GT_CampBatRank_List[id].TotalDestroyID = rs.GetFieldInt("total_destroy_rank")
	GT_CampBatRank_List[id].MinLevel = rs.GetFieldInt("rank_min")
	GT_CampBatRank_List[id].MaxLevel = rs.GetFieldInt("rank_max")
}

//! 提供接口
func GetCampBatRank(level int) *ST_CampBatRank {
	if level <= 0 {
		gamelog.Error("GetCampBatRank Error : Invalid level:%d", level)
		return nil
	}

	for i, v := range GT_CampBatRank_List {
		if level >= v.MinLevel && level < v.MaxLevel {
			return &GT_CampBatRank_List[i]
		}
	}
	return nil
}

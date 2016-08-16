package gamedata

import (
	"gamelog"
)

type ST_FoodWarRankAward struct {
	ID       int
	Rank_min int
	Rank_max int
	Award    int
}

var GT_FoodWarRankAward []ST_FoodWarRankAward

func InitFoodWarRankAwardParser(total int) bool {
	GT_FoodWarRankAward = make([]ST_FoodWarRankAward, total+1)
	return true
}

func ParseFoodWarRankAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_FoodWarRankAward[id].ID = id
	GT_FoodWarRankAward[id].Rank_min = rs.GetFieldInt("rank_min")
	GT_FoodWarRankAward[id].Rank_max = rs.GetFieldInt("rank_max")
	GT_FoodWarRankAward[id].Award = rs.GetFieldInt("award")
}

func GetFoodWarRankAward(rank int) int {
	for _, v := range GT_FoodWarRankAward {
		if rank >= v.Rank_min && rank <= v.Rank_max {
			return v.Award
		}
	}

	return 0
}

type ST_FoodWarAward struct {
	ID     int
	Target int
	Award  int
}

var GT_FoodWarAward []ST_FoodWarAward

func InitFoodWarAwardParser(total int) bool {
	GT_FoodWarAward = make([]ST_FoodWarAward, total+1)
	return true
}

func ParseFoodWarAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_FoodWarAward[id].ID = id
	GT_FoodWarAward[id].Target = rs.GetFieldInt("target")
	GT_FoodWarAward[id].Award = rs.GetFieldInt("award")
}

func GetFoodWarAward(id int) *ST_FoodWarAward {
	if id > len(GT_FoodWarAward)-1 {
		gamelog.Error("GetFoodWarRankAward Error: Invalid id %d", id)
	}

	return &GT_FoodWarAward[id]
}

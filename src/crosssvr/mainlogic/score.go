package mainlogic

//import (
//	"utility"
//)

var (
	G_ScoreRanker TRoleRanker //积分赛排行榜

)

func InitRankMgr() {
	//积分赛排行榜
	InitScoreRanker()
}

//等级排行榜
func InitScoreRanker() bool {
	G_ScoreRanker.InitRanker(20, 1000)

	return true
}

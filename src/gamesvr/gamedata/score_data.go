package gamedata

import (
	"gamelog"
)

//! 积分赛段位奖励
type ST_ScoreDw struct {
	DuanWei int //! 段位数
	Score   int //! 分数线
	AwardID int //! 奖励ID
}

type ST_ScoreTime struct {
	TimeID  int //! 时间ID
	Times   int //! 战斗位数
	AwardID int //! 奖励ID
}

type ST_ScoreRank struct {
	MinRank int //最小排名
	MaxRank int //最大排名
	AwardID int //! 奖励ID
}

var (
	GT_ScoreDw_List   []ST_ScoreDw   = nil
	GT_ScoreTimeAward []ST_ScoreTime = nil
	GT_ScoreRankAward []ST_ScoreRank = nil
)

//! 初始化积分赛分析器
func InitScoreDwParser(total int) bool {
	GT_ScoreDw_List = make([]ST_ScoreDw, total+1)
	return true
}

func ParseScoreDwRecord(rs *RecordSet) {
	duan := rs.GetFieldInt("duanwei")
	GT_ScoreDw_List[duan].DuanWei = duan
	GT_ScoreDw_List[duan].Score = rs.GetFieldInt("score")
	GT_ScoreDw_List[duan].AwardID = rs.GetFieldInt("award")
}

//! 初始化积分赛奖励分析器
func InitScoreAwardParser(total int) bool {
	GT_ScoreTimeAward = make([]ST_ScoreTime, 0)
	GT_ScoreRankAward = make([]ST_ScoreRank, 0)
	return true
}

func ParseScoreAwardRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	aType := rs.GetFieldInt("type")
	fightTime := rs.GetFieldInt("fight_time")
	minLevel := rs.GetFieldInt("min_level")
	maxLevel := rs.GetFieldInt("max_level")
	awardid := rs.GetFieldInt("award")
	if aType == 1 {
		GT_ScoreTimeAward = append(GT_ScoreTimeAward, ST_ScoreTime{id, fightTime, awardid})
	} else if aType == 2 {
		GT_ScoreRankAward = append(GT_ScoreRankAward, ST_ScoreRank{minLevel, maxLevel, awardid})
	}
}

//! 提供接口
func GetScoreDwAward(duan int) *ST_ScoreDw {
	if duan >= len(GT_ScoreDw_List) || duan <= 0 {
		gamelog.Error("GetScoreDwAward Error: invalid duan %d", duan)
		return nil
	}

	return &GT_ScoreDw_List[duan]
}

//获取当前的段位
func GetScoreDuanWei(score int) int {
	var duan int = 0
	for i := 0; i < len(GT_ScoreDw_List); i++ {
		if score >= GT_ScoreDw_List[i].Score {
			duan = GT_ScoreDw_List[i].DuanWei
		} else {
			return duan
		}
	}

	return duan
}

//获取积分赛参与次数奖励
func GetScoreTimeAward(timeid int) *ST_ScoreTime {
	for i := 0; i < len(GT_ScoreTimeAward); i++ {
		if timeid == GT_ScoreTimeAward[i].TimeID {
			return &GT_ScoreTimeAward[i]
		}
	}

	return nil
}

//获取排名奖励
func GetScoreRankAward(rank int) int {
	for i := 0; i < len(GT_ScoreRankAward); i++ {
		if GT_ScoreRankAward[i].MinRank <= rank && rank <= GT_ScoreRankAward[i].MaxRank {
			return GT_ScoreRankAward[i].AwardID
		}
	}

	return 0
}

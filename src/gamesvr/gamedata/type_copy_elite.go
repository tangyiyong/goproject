package gamedata

import (
	"fmt"
	"gamelog"
)

type ST_EliteAward struct {
	StarNum int //星数
	AwardID int //奖励ID
}

type ST_EliteSceneAward struct {
	Levels  int //! 关卡
	AwardID int //! 奖励ID
}

type ST_Elitechapter struct {
	Chapter     int                //章节
	StartID     int                //开始副本ID
	EndID       int                //结束副本ID
	StarAwards  [3]ST_EliteAward   //精英副本奖励
	SceneAwards ST_EliteSceneAward //场景奖励宝箱
	InvadeID    int                //入侵副本ID
}

var (
	GT_EliteChapterList []ST_Elitechapter //精英奖励
)

func InitEliteParser(total int) bool {
	GT_EliteChapterList = make([]ST_Elitechapter, total+1) //精英奖励
	return true
}

func ParseEliteRecord(rs *RecordSet) {

	chapter := CheckAtoi(rs.Values[0], 0)

	GT_EliteChapterList[chapter].Chapter = chapter
	GT_EliteChapterList[chapter].StartID = rs.GetFieldInt("start_copy_id")
	GT_EliteChapterList[chapter].EndID = rs.GetFieldInt("end_copy_id")
	GT_EliteChapterList[chapter].InvadeID = rs.GetFieldInt("invade_copy")
	for i := 1; i < 4; i++ {
		fieldName := fmt.Sprintf("star%d", i)

		var starAward ST_EliteAward
		starAward.StarNum = rs.GetFieldInt(fieldName)

		fieldName = fmt.Sprintf("award%d", i)
		starAward.AwardID = rs.GetFieldInt(fieldName)
		GT_EliteChapterList[chapter].StarAwards[i-1] = starAward
	}

	var sceneAward ST_EliteSceneAward
	fieldName := fmt.Sprintf("levels")
	sceneAward.Levels = rs.GetFieldInt(fieldName)

	fieldName = fmt.Sprintf("sceneaward")
	sceneAward.AwardID = rs.GetFieldInt(fieldName)
	GT_EliteChapterList[chapter].SceneAwards = sceneAward

}

func GetEliteChapterInfo(chapter int) *ST_Elitechapter {
	if chapter >= len(GT_EliteChapterList) || chapter <= 0 {
		gamelog.Error("GetElitechapterInfo Error: invalid chapter :%d", chapter)
		return nil
	}

	return &GT_EliteChapterList[chapter]
}

func GetEliteChapterCount() int {
	return len(GT_EliteChapterList) - 1
}

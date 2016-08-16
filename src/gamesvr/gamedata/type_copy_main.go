package gamedata

import (
	"fmt"
	"gamelog"
)

type ST_MainAward struct {
	StarNum int //星数
	AwardID int //奖励ID
}

type ST_MainSceneAward struct {
	Levels  int //! 关卡
	AwardID int //! 奖励ID
}

type ST_MainChapter struct {
	Chapter     int                  //章节
	StartID     int                  //开始副本ID
	EndID       int                  //结束副本ID
	StarAwards  [3]ST_MainAward      //主线副本奖励
	SceneAwards [3]ST_MainSceneAward //场景奖励宝箱
}

var (
	GT_MainChapterList []ST_MainChapter //主线奖励
)

func InitMainParser(total int) bool {
	GT_MainChapterList = make([]ST_MainChapter, total+1) //主线奖励
	return true
}

func ParseMainRecord(rs *RecordSet) {

	chapter := CheckAtoi(rs.Values[0], 0)

	GT_MainChapterList[chapter].Chapter = chapter
	GT_MainChapterList[chapter].StartID = rs.GetFieldInt("start_copy_id")
	GT_MainChapterList[chapter].EndID = rs.GetFieldInt("end_copy_id")

	for i := 1; i < 4; i++ {
		fieldName := fmt.Sprintf("star%d", i)

		var starAward ST_MainAward
		starAward.StarNum = rs.GetFieldInt(fieldName)

		fieldName = fmt.Sprintf("award%d", i)
		starAward.AwardID = rs.GetFieldInt(fieldName)
		GT_MainChapterList[chapter].StarAwards[i-1] = starAward

		var sceneAward ST_MainSceneAward
		fieldName = fmt.Sprintf("levels%d", i)
		sceneAward.Levels = rs.GetFieldInt(fieldName)

		fieldName = fmt.Sprintf("sceneaward%d", i)
		sceneAward.AwardID = rs.GetFieldInt(fieldName)
		GT_MainChapterList[chapter].SceneAwards[i-1] = sceneAward
	}

}

func GetMainChapterInfo(chapter int) *ST_MainChapter {
	if chapter >= len(GT_MainChapterList) || chapter <= 0 {
		gamelog.Error("GetMainchapterInfo Error: invalid chapter :%d", chapter)
		return nil
	}

	return &GT_MainChapterList[chapter]
}

func GetMainChapterCount() int {
	return len(GT_MainChapterList)
}

func GetChaperCopyStartID(chapter int, copyType int) int {
	if copyType == COPY_TYPE_Main {
		return GT_MainChapterList[chapter].StartID
	} else if copyType == COPY_TYPE_Elite {
		return GT_EliteChapterList[chapter].StartID
	} else {
		gamelog.Error("GetChaperCopyStartID Error: invalid copyType :%d", copyType)
	}

	return 0
}
func GetChaperCopyEndID(chapter int, copyType int) int {
	if copyType == COPY_TYPE_Main {
		return GT_MainChapterList[chapter].EndID
	} else if copyType == COPY_TYPE_Elite {
		return GT_EliteChapterList[chapter].EndID
	} else {
		gamelog.Error("GetChaperCopyEndID Error: invalid copyType :%d", copyType)
	}

	return 0
}

func IsChapterEnd(copyID int, chapter int, copyType int) bool {
	if chapter <= 0 || copyType <= 0 {
		gamelog.Error("IsChapterEnd Error : Invalid copyid:%d, chapter :%d, copytype:%d", copyID, chapter, copyType)
		return false
	}

	if copyID == GetChaperCopyEndID(chapter, copyType) {
		return true
	}

	return false
}

func GetNextCopy(copyID int, chapter int, copyType int) (int, int) {
	if chapter <= 0 || copyType <= 0 {
		gamelog.Error("GetNextCopy Error : Invalid copyid:%d, chapter :%d, copytype:%d", copyID, chapter, copyType)
		return 0, 0
	}

	var nextCopyID int = copyID
	var nextChapter int = chapter

	if copyID == 0 {
		nextCopyID = GetChaperCopyStartID(chapter, copyType)
		return nextCopyID, nextChapter
	}

	isEnd := IsChapterEnd(copyID, chapter, copyType)
	if isEnd == true {
		nextChapter = chapter + 1
		if copyType == COPY_TYPE_Main {
			chapterInfo := GetMainChapterInfo(nextChapter)
			nextCopyID = chapterInfo.StartID

		} else if copyType == COPY_TYPE_Famous {
			chapterInfo := GetFamousChapterInfo(nextChapter)
			nextCopyID = chapterInfo.StartID
		} else if copyType == COPY_TYPE_Elite {
			chapterInfo := GetEliteChapterInfo(nextChapter)
			nextCopyID = chapterInfo.StartID
		} else {
			gamelog.Error("GetNextCopy Error : copytype:%d", copyType)
		}
	} else {
		nextCopyID += 1
	}

	return nextCopyID, nextChapter
}

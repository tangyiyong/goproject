package gamedata

import (
	"gamelog"
)

type ST_FamousChapter struct {
	Chapter  int //! 章节
	StartID  int //! 开始副本ID
	EndID    int //! 结束副本ID
	SerialID int //! 连环计 为0则没有连环计
	PreCopy  int //! 前置关卡 通过该关卡后开启名将副本
	Award    int //! 章节通关宝箱
}

var GT_FamousChapterList []ST_FamousChapter //! 名将副本

func InitFamousParser(total int) bool {
	GT_FamousChapterList = make([]ST_FamousChapter, total+1)
	return true
}

func ParseFamousRecord(rs *RecordSet) {
	chapter := CheckAtoi(rs.Values[0], 0)

	GT_FamousChapterList[chapter].Chapter = chapter
	GT_FamousChapterList[chapter].StartID = rs.GetFieldInt("startid")
	GT_FamousChapterList[chapter].EndID = rs.GetFieldInt("endid")
	GT_FamousChapterList[chapter].SerialID = rs.GetFieldInt("serialid")
	GT_FamousChapterList[chapter].PreCopy = rs.GetFieldInt("precopy")
	GT_FamousChapterList[chapter].Award = rs.GetFieldInt("award")
}

func GetFamousChapterInfo(chapter int) *ST_FamousChapter {
	if chapter >= len(GT_FamousChapterList) {
		gamelog.Error("GetFamousChapterData Error: invalid chapter :%d", chapter)
		return nil
	}

	return &GT_FamousChapterList[chapter]
}

func IsSerialCopy(chapter int, copyID int) bool {
	if GT_FamousChapterList[chapter].SerialID == copyID {
		return true
	}
	return false
}

func GetFamousChapterCount() int {
	return len(GT_FamousChapterList)
}

package gamedata

import (
	"gamelog"
)

type ST_DailyCopy struct {
	ID         int //! 唯一标识
	ResType    int //! 资源ID
	CopyID     int //! 副本ID
	Difficulty int //! 难度
	Level      int //! 开启等级
	Type       int //! 副本类型 1->1,3,5  2->2,4,6
}

var (
	GT_DailyCopyList    []ST_DailyCopy
	GT_DailyResTypeList []int
)

func InitDailyParse(total int) bool {
	GT_DailyCopyList = make([]ST_DailyCopy, total+1)
	return true
}

func ParseDailyRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)

	GT_DailyCopyList[id].ID = id
	GT_DailyCopyList[id].ResType = rs.GetFieldInt("restype")
	GT_DailyCopyList[id].CopyID = rs.GetFieldInt("copyid")
	GT_DailyCopyList[id].Difficulty = rs.GetFieldInt("difficulty")
	GT_DailyCopyList[id].Level = rs.GetFieldInt("level")
	GT_DailyCopyList[id].Type = rs.GetFieldInt("type")

	var isExist bool
	for _, res := range GT_DailyResTypeList {
		if res == GT_DailyCopyList[id].ResType {
			isExist = true
			break
		}
	}

	if isExist == false {
		GT_DailyResTypeList = append(GT_DailyResTypeList, GT_DailyCopyList[id].ResType)
	}
}

func GetDailyCopyData(copyID int) *ST_DailyCopy {
	for i := 0; i < len(GT_DailyCopyList); i++ {
		if GT_DailyCopyList[i].CopyID == copyID {
			return &GT_DailyCopyList[i]
		}
	}

	gamelog.Error("GetDailyCopyData fail.")

	return nil
}

func GetDailyCopyDataFromLevel(level int) []ST_DailyCopy {
	var copyLst []ST_DailyCopy
	for _, v := range GT_DailyCopyList {
		if level >= v.Level {
			copyLst = append(copyLst, v)
		}
	}
	return copyLst
}

func GetDailyCopyDataFromType(dateType int) []ST_DailyCopy {
	var copyLst []ST_DailyCopy
	if dateType == 3 {
		for _, v := range GT_DailyCopyList {
			if v.ID == 0 {
				continue
			}
			copyLst = append(copyLst, v)
		}

		return copyLst
	}

	for _, v := range GT_DailyCopyList {
		if v.ID == 0 {
			continue
		}

		if dateType == v.Type {
			copyLst = append(copyLst, v)
		}
	}
	return copyLst
}

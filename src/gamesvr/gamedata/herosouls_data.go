package gamedata

import (
	"fmt"
	"gamelog"
)

//! 名将试炼
type ST_HeroSoulsTrials struct {
	ID       int //! 标识
	RandType int //! 随机库 1->低级库 2->普通库 3->高级库  4->更高级哭
	HeroID   int //! 将灵ID(道具)
	Weight   int //! 权重
}

var GT_HeroSoulsTrialsLst []ST_HeroSoulsTrials

func InitHeroSoulsTrialsParser(total int) bool {
	GT_HeroSoulsTrialsLst = make([]ST_HeroSoulsTrials, total+1)
	return true
}

func ParseHeroSoulsTrialRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_HeroSoulsTrialsLst[id].ID = id
	GT_HeroSoulsTrialsLst[id].HeroID = rs.GetFieldInt("item_id")
	GT_HeroSoulsTrialsLst[id].RandType = rs.GetFieldInt("rand_type")
	GT_HeroSoulsTrialsLst[id].Weight = rs.GetFieldInt("weight")
}

func RandHeroSouls() (randTrialsID []int, IDLst []int) {
	var randArray_Low []ST_HeroSoulsTrials
	var randArray_Normal []ST_HeroSoulsTrials
	var randArray_High []ST_HeroSoulsTrials
	var randArray_Max []ST_HeroSoulsTrials

	var totalWeight_Low, totalWeight_Normal, totalWeight_High, totalWeight_Max int

	//! 将随机库按品质区分
	for _, v := range GT_HeroSoulsTrialsLst {

		if v.RandType == 3 {
			randArray_High = append(randArray_High, v)
			totalWeight_High += v.Weight
		} else if v.RandType == 2 {
			randArray_Normal = append(randArray_Normal, v)
			totalWeight_Normal += v.Weight
		} else if v.RandType == 1 {
			randArray_Low = append(randArray_Low, v)
			totalWeight_Low += v.Weight
		} else if v.RandType == 4 {
			randArray_Max = append(randArray_Max, v)
			totalWeight_Max += v.Weight
		}
	}

	if len(randArray_High) <= 0 || len(randArray_Normal) <= 0 || len(randArray_Low) <= 0 || len(randArray_Max) <= 0 {
		gamelog.Error("RandHeroSouls Error: randArray is nil")
		return
	}

	//! 从品质中随机出对应的英雄将灵
	for i := 0; i < 2; i++ {
		randWeight := r.Intn(totalWeight_High)
		curWeight := 0
		index := 0
		for i, v := range randArray_High {
			if randWeight >= curWeight && randWeight < v.Weight+curWeight {
				randTrialsID = append(randTrialsID, v.HeroID)
				IDLst = append(IDLst, v.ID)
				totalWeight_High -= v.Weight
				index = i
				break
			}

			curWeight += v.Weight
		}

		if index == 0 {
			randArray_High = randArray_High[1:]
		} else if (index + 1) == len(randArray_High) {
			randArray_High = randArray_High[:index]
		} else {
			randArray_High = append(randArray_High[:index], randArray_High[index+1:]...)
		}
	}

	for i := 0; i < 3; i++ {
		randWeight := r.Intn(totalWeight_Normal)
		curWeight := 0
		index := 0
		for i, v := range randArray_Normal {
			if randWeight >= curWeight && randWeight < v.Weight+curWeight {
				randTrialsID = append(randTrialsID, v.HeroID)
				IDLst = append(IDLst, v.ID)
				totalWeight_Normal -= v.Weight
				index = i
				break
			}

			curWeight += v.Weight
		}

		if index == 0 {
			randArray_Normal = randArray_Normal[1:]
		} else if (index + 1) == len(randArray_Normal) {
			randArray_Normal = randArray_Normal[:index]
		} else {
			randArray_Normal = append(randArray_Normal[:index], randArray_Normal[index+1:]...)
		}
	}

	for i := 0; i < 1; i++ {
		randWeight := r.Intn(totalWeight_Low)
		curWeight := 0
		index := 0
		for i, v := range randArray_Low {
			if randWeight >= curWeight && randWeight < v.Weight+curWeight {
				randTrialsID = append(randTrialsID, v.HeroID)
				IDLst = append(IDLst, v.ID)
				totalWeight_Low -= v.Weight
				index = i
				break
			}

			curWeight += v.Weight
		}

		if index == 0 {
			randArray_Low = randArray_Low[1:]
		} else if (index + 1) == len(randArray_Low) {
			randArray_Low = randArray_Low[:index]
		} else {
			randArray_Low = append(randArray_Low[:index], randArray_Low[index+1:]...)
		}
	}

	for i := 0; i < 2; i++ {
		randWeight := r.Intn(totalWeight_Max)
		curWeight := 0
		index := 0
		for i, v := range randArray_Max {
			if randWeight >= curWeight && randWeight < v.Weight+curWeight {
				randTrialsID = append(randTrialsID, v.HeroID)
				IDLst = append(IDLst, v.ID)
				totalWeight_Max -= v.Weight
				index = i
				break
			}

			curWeight += v.Weight
		}

		if index == 0 {
			randArray_Max = randArray_Max[1:]
		} else if (index + 1) == len(randArray_Max) {
			randArray_Max = randArray_Max[:index]
		} else {
			randArray_Max = append(randArray_Max[:index], randArray_Max[index+1:]...)
		}
	}

	return randTrialsID, IDLst
}

func GetHeroSoulsTrialInfo(id int) *ST_HeroSoulsTrials {
	if id > len(GT_HeroSoulsTrialsLst)-1 {
		gamelog.Error("GetHeroSoulsTrialInfo Error: Invalid id %d", id)
		return nil
	}

	return &GT_HeroSoulsTrialsLst[id]
}

//! 将灵
type ST_HeroSoulsProperty struct {
	PropertyID    int
	PropertyValue int
	Camp          int
	Is_Percent    bool
	LevelUp       int
}

type ST_HeroSouls struct {
	ID       int
	HeroID   [3]int
	Property [2]ST_HeroSoulsProperty
}

var GT_HeroSoulsLst []ST_HeroSouls

func InitHeroSoulsParser(total int) bool {
	GT_HeroSoulsLst = make([]ST_HeroSouls, total+1)
	return true
}

func ParseHeroSoulsRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_HeroSoulsLst[id].ID = id

	for i := 0; i < 3; i++ {
		filedName := fmt.Sprintf("itemid%d", i+1)
		GT_HeroSoulsLst[id].HeroID[i] = rs.GetFieldInt(filedName)
	}

	for i := 0; i < 2; i++ {
		filedName := fmt.Sprintf("propertyid%d", i+1)
		GT_HeroSoulsLst[id].Property[i].PropertyID = rs.GetFieldInt(filedName)

		filedName = fmt.Sprintf("propertyvalue%d", i+1)
		GT_HeroSoulsLst[id].Property[i].PropertyValue = rs.GetFieldInt(filedName)

		filedName = fmt.Sprintf("camp%d", i+1)
		GT_HeroSoulsLst[id].Property[i].Camp = rs.GetFieldInt(filedName)

		filedName = fmt.Sprintf("is_percent%d", i+1)
		GT_HeroSoulsLst[id].Property[i].Is_Percent = (rs.GetFieldInt(filedName) == 1)

		filedName = fmt.Sprintf("level_up%d", i+1)
		GT_HeroSoulsLst[id].Property[i].LevelUp = rs.GetFieldInt(filedName)
	}

}

func GetHeroSoulsInfo(id int) *ST_HeroSouls {
	if id > len(GT_HeroSoulsLst)-1 {
		gamelog.Error("GetHeroSoulsInfo Error: Invalid id %d", id)
		return nil
	}

	return &GT_HeroSoulsLst[id]
}

//! 将灵章节
type ST_HeroSoulsChapter struct {
	Chapter       int //! 章节
	BeginSoulsID  int //! 起始将灵链接ID
	EndSoulsID    int //! 终止将灵链接ID
	UnLockChapter int //! 解锁需达成章节X
	UnlockCount   int //! 激活X条将灵链接
}

var GT_HeroSoulsChapter []ST_HeroSoulsChapter

func InitHeroSoulsChapterParser(total int) bool {
	GT_HeroSoulsChapter = make([]ST_HeroSoulsChapter, total+1)
	return true
}

func ParseHeroSoulsChapterRecord(rs *RecordSet) {
	chapter := CheckAtoi(rs.Values[0], 0)

	GT_HeroSoulsChapter[chapter].Chapter = chapter
	GT_HeroSoulsChapter[chapter].BeginSoulsID = rs.GetFieldInt("begin_id")
	GT_HeroSoulsChapter[chapter].EndSoulsID = rs.GetFieldInt("end_id")
	GT_HeroSoulsChapter[chapter].UnLockChapter = rs.GetFieldInt("unlock_chapter")
	GT_HeroSoulsChapter[chapter].UnlockCount = rs.GetFieldInt("unlock_count")
}

func GetHeroSoulsChapterCount() int {
	return len(GT_HeroSoulsChapter) - 1
}

func GetHeroSoulsChapterInfo(chapter int) *ST_HeroSoulsChapter {
	if chapter > GetHeroSoulsChapterCount() {
		gamelog.Error("GetHeroSoulsChapterInfo Error: Invalid chapter %d", chapter)
		return nil
	}

	return &GT_HeroSoulsChapter[chapter]
}

func GetHeroSoulsBelongChapter(soulsID int) int {
	for _, v := range GT_HeroSoulsChapter {
		if soulsID <= v.EndSoulsID && soulsID >= v.BeginSoulsID {
			return v.Chapter
		}
	}

	gamelog.Error("GetHeroSoulsBelongChapter Error: Invalid soulsID %d", soulsID)
	return 0
}

//! 将灵商店
type ST_HeroSoulsStore struct {
	HeroID   int
	MoneyID  int
	MoneyNum int
}

var GT_HeroSoulsStoreLst []ST_HeroSoulsStore

func InitHeroSoulsStoreParser(total int) bool {
	return true
}

func ParseHeroSoulsStoreRecrod(rs *RecordSet) {
	var goods ST_HeroSoulsStore
	goods.HeroID = rs.GetFieldInt("heroid")
	goods.MoneyID = rs.GetFieldInt("moneyid")
	goods.MoneyNum = rs.GetFieldInt("moneynum")

	GT_HeroSoulsStoreLst = append(GT_HeroSoulsStoreLst, goods)
}

func GetHeroSoulsStoreInfo(heroID int) *ST_HeroSoulsStore {
	for i, v := range GT_HeroSoulsStoreLst {
		if v.HeroID == heroID {
			return &GT_HeroSoulsStoreLst[i]
		}
	}

	gamelog.Error("GetHeroSoulsStoreInfo Error: Not find heroID: %d", heroID)
	return nil
}

func RandHeroSoulsStore(needNum int) (IDLst []int) {
	if needNum > len(GT_HeroSoulsStoreLst) {
		gamelog.Error("RandHeroSoulsStore Error: GT_HeroSoulsStoreLst length not enough %d", needNum)
		return IDLst
	}

	for i := 0; i < needNum; i++ {
		randIndex := r.Intn(len(GT_HeroSoulsStoreLst))

		isExist := false
		for _, v := range IDLst {
			if v == GT_HeroSoulsStoreLst[randIndex].HeroID {
				isExist = true
				break
			}
		}

		if isExist == false {
			IDLst = append(IDLst, GT_HeroSoulsStoreLst[randIndex].HeroID)
		} else {
			i -= 1
		}
	}

	return IDLst
}

//! 阵魂图加成
type ST_SoulMap struct {
	ID            int
	Souls         int  //! 需求阵魂值
	PropertyID    int  //! 加成属性ID
	PropertyValue int  //! 加成属性数量
	Is_percent    bool //! 是否为百分比加成
}

var GT_SoulMapLst []ST_SoulMap

func InitSoulMapParser(total int) bool {
	GT_SoulMapLst = make([]ST_SoulMap, total+1)
	return true
}

func ParseSoulMapRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_SoulMapLst[id].ID = id
	GT_SoulMapLst[id].Souls = rs.GetFieldInt("souls")
	GT_SoulMapLst[id].PropertyID = rs.GetFieldInt("propertyid")
	GT_SoulMapLst[id].PropertyValue = rs.GetFieldInt("propertyvalue")
	GT_SoulMapLst[id].Is_percent = (1 == rs.GetFieldInt("is_percent"))
}

func GetSoulMapInfo(id int) *ST_SoulMap {
	if id > len(GT_SoulMapLst)-1 {
		gamelog.Error("GetSoulMapInfo Error: Invalid id %d", id)
		return nil
	}

	return &GT_SoulMapLst[id]
}

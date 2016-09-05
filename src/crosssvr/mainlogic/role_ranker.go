package mainlogic

import (
	"sort"
	"sync"
)

type TRoleRankItem struct {
	RankID     int32  //角色ID
	RankValue  int    //排名的值
	Quality    int8   //品质
	HeroID     int    //英雄ID
	Level      int    //等级
	FightValue int32  //战力
	SvrID      int32  //服务器ID
	SvrName    string //服务器名
	SvrIp      string //服务器IP地址
	RoleName   string //角色名
}

type TRoleRankList []*TRoleRankItem

type TRoleRanker struct {
	sync.Mutex
	List    TRoleRankList
	ShowNum int //显示数量
	RankNum int //排行数量
}

func (ranker *TRoleRanker) InitRanker(show int, rank int) {
	ranker.ShowNum = show
	ranker.RankNum = rank
	ranker.List = make([]*TRoleRankItem, rank)
	for i := 0; i < len(ranker.List); i++ {
		ranker.List[i] = &TRoleRankItem{0, 0, 0, 0, 0, 0, 0, "", "", ""}
	}
}

func (ranker *TRoleRanker) SetRankItem(rankid int32, rankvalue int, level int,
	fightvalue int32, heroid int, svrid int32, quality int8, svrname string, rolename string) int {
	ranker.Lock()
	defer ranker.Unlock()
	nCount := len(ranker.List)
	MinValue := ranker.List[nCount-1].RankValue
	if rankvalue <= MinValue {
		return -1
	}

	targetIndex := sort.Search(nCount, func(i int) bool {
		if ranker.List[i].RankValue <= rankvalue {
			return true
		}
		return false
	})

	myIndex := nCount - 1
	for i := targetIndex; i < nCount; i++ {
		if ranker.List[i].RankID == rankid || ranker.List[i].RankID == 0 {
			ranker.List[i].RankID = rankid
			ranker.List[i].RankValue = rankvalue
			ranker.List[i].HeroID = heroid
			ranker.List[i].Level = level
			ranker.List[i].FightValue = fightvalue
			ranker.List[i].SvrID = svrid
			ranker.List[i].SvrName = svrname
			ranker.List[i].RoleName = rolename
			ranker.List[i].Quality = quality
			myIndex = i
			break
		}
	}

	if myIndex == targetIndex {
		return targetIndex
	}

	for i := myIndex; i > targetIndex; i-- {
		ranker.List[i].RankID = ranker.List[i-1].RankID
		ranker.List[i].RankValue = ranker.List[i-1].RankValue
		ranker.List[i].HeroID = ranker.List[i-1].HeroID
		ranker.List[i].Level = ranker.List[i-1].Level
		ranker.List[i].FightValue = ranker.List[i-1].FightValue
		ranker.List[i].SvrID = ranker.List[i-1].SvrID
		ranker.List[i].SvrName = ranker.List[i-1].SvrName
		ranker.List[i].RoleName = ranker.List[i-1].RoleName
	}

	ranker.List[targetIndex].RankID = rankid
	ranker.List[targetIndex].RankValue = rankvalue
	ranker.List[targetIndex].HeroID = heroid
	ranker.List[targetIndex].Level = level
	ranker.List[targetIndex].FightValue = fightvalue
	ranker.List[targetIndex].SvrID = svrid
	ranker.List[targetIndex].SvrName = svrname
	ranker.List[targetIndex].RoleName = rolename
	return targetIndex
}

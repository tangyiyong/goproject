package gamedata

import (
	"gamelog"
)

type ST_HeroFriendInfo struct {
	ID        int
	Level     int //等级
	Propertys [5][2]int
}

var (
	GT_HeroFriend_List []ST_HeroFriendInfo = nil
)

func InitHeroFriendParser(total int) bool {
	GT_HeroFriend_List = make([]ST_HeroFriendInfo, total+1)

	return true
}

func ParseHeroFriendRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_HeroFriend_List[id].Level = rs.GetFieldInt("level")
	GT_HeroFriend_List[id].Propertys[0] = ParseTo2IntSlice(rs.GetFieldString("property_id_1"))
	GT_HeroFriend_List[id].Propertys[1] = ParseTo2IntSlice(rs.GetFieldString("property_id_2"))
	GT_HeroFriend_List[id].Propertys[2] = ParseTo2IntSlice(rs.GetFieldString("property_id_3"))
	GT_HeroFriend_List[id].Propertys[3] = ParseTo2IntSlice(rs.GetFieldString("property_id_4"))
	GT_HeroFriend_List[id].Propertys[4] = ParseTo2IntSlice(rs.GetFieldString("property_id_5"))

	return
}

func GetHeroFriendInfo(level int) (pInfo *ST_HeroFriendInfo) {
	if level <= 0 {
		gamelog.Error("GetHeroFriendInfo Error: Invalid level:%d", level)
		return
	}
	pInfo = nil
	for i := 0; i < len(GT_HeroFriend_List); i++ {
		pInfo = &GT_HeroFriend_List[i]
		if level < GT_HeroFriend_List[i].Level {
			break
		}
	}

	if pInfo != nil {
		if level < pInfo.Level {
			pInfo = nil
			return
		}
	}

	return
}

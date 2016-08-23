package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TArenaRankInfo struct {
	PlayerID int32
	IsRobot  bool
}

type TArenaRankInfoLst [5000]TArenaRankInfo

var (
	G_Rank_List TArenaRankInfoLst
)

//! 初始化竞技场前五千名
func InitArenaMgr() bool {
	arenaInfoLst := []TArenaModule{}
	s := mongodb.GetDBSession()
	defer s.Close()
	err := s.DB(appconfig.GameDbName).C("PlayerArena").Find(bson.M{"currentrank": bson.M{"$lt": 5000}}).All(&arenaInfoLst)
	if err != nil {
		if err != mgo.ErrNotFound {
			gamelog.Error("Init DB Error!!!")
			return false
		}
	}

	for _, v := range arenaInfoLst {
		G_Rank_List[v.CurrentRank-1].IsRobot = false
		G_Rank_List[v.CurrentRank-1].PlayerID = v.PlayerID
	}

	for i := int32(0); i < 5000; i++ {
		if G_Rank_List[i].PlayerID == 0 {
			pRobot := gamedata.GetRobot(i%10 + 1)
			if pRobot == nil {
				gamelog.Error("GetRobot error: robotID: %d", i%10+1)
				G_Rank_List[i].PlayerID = 0
				G_Rank_List[i].IsRobot = true
				continue
			}
			G_Rank_List[i].PlayerID = pRobot.RobotID
			G_Rank_List[i].IsRobot = true
		}
	}

	return true
}

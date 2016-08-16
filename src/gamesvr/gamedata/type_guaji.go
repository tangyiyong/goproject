package gamedata

import (
	"gamelog"
)

type ST_HangUp struct {
	BossID     int //副本ID
	FightValue int //每秒获取的金币数
	Level      int //每秒获取的经验数
	ProduceID  int
	ProduceNum int
	CDTime     int //间隔时间
}

var (
	GT_HangUpList []ST_HangUp //主线奖励
)

func InitHangUpParser(total int) bool {
	GT_HangUpList = make([]ST_HangUp, total+1)

	return true
}

func ParseHangUpRecord(rs *RecordSet) {
	BossID := CheckAtoi(rs.Values[0], 0)
	GT_HangUpList[BossID].BossID = BossID
	GT_HangUpList[BossID].FightValue = rs.GetFieldInt("fightvalue")
	GT_HangUpList[BossID].Level = rs.GetFieldInt("level")
	GT_HangUpList[BossID].ProduceID = rs.GetFieldInt("produce_item_id")
	GT_HangUpList[BossID].ProduceNum = rs.GetFieldInt("produce_item_num")
	GT_HangUpList[BossID].CDTime = rs.GetFieldInt("cd_time")
}

func GetHangUpInfo(bossid int) *ST_HangUp {
	if bossid >= len(GT_HangUpList) || bossid == 0 {
		gamelog.Error("GetHangUpInfo Error: invalid bossid :%d", bossid)
		return nil
	}

	return &GT_HangUpList[bossid]
}

package gamedata

import (
	"gamelog"
	"strings"
)

type ST_RobotHero struct {
	HeroID    int     //英雄ID
	Level     int     //英雄等级
	Propertys [11]int //六个英雄属性
}

type ST_Robot struct {
	RobotID    int32           //机器人ID
	Level      int             //机器人等级
	Name       string          //名字
	FightValue int32           //战力
	Quality    int8            //品质
	Heros      [6]ST_RobotHero //六个英雄
}

var (
	GT_Robot_List []ST_Robot = nil
)

func InitRobotParser(total int) bool {
	GT_Robot_List = make([]ST_Robot, total+1)
	return true
}

func ParseRobotRecord(rs *RecordSet) {
	RobotID := int32(rs.GetFieldInt("id"))
	if (RobotID <= 0) || (RobotID >= int32(len(GT_Robot_List))) {
		gamelog.Error("ParseRobotRecord Error: Invalid RobotID :%d", RobotID)
		return
	}
	GT_Robot_List[RobotID].RobotID = RobotID

	GT_Robot_List[RobotID].Level = rs.GetFieldInt("level")
	GT_Robot_List[RobotID].FightValue = int32(rs.GetFieldInt("fightvalue"))
	GT_Robot_List[RobotID].Name = rs.GetFieldString("name")
	GT_Robot_List[RobotID].Quality = int8(rs.GetFieldInt("quality"))
	var heroindex int = 0
	for i := 5; i < 11; i++ {
		if rs.Values[i] == "NULL" {
			break
		}
		slice := strings.Split(rs.Values[i], "|")
		GT_Robot_List[RobotID].Heros[heroindex].HeroID = CheckAtoi(slice[0], i)
		GT_Robot_List[RobotID].Heros[heroindex].Level = CheckAtoi(slice[1], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[0] = CheckAtoi(slice[2], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[1] = CheckAtoi(slice[3], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[2] = CheckAtoi(slice[4], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[3] = CheckAtoi(slice[5], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[4] = CheckAtoi(slice[6], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[5] = CheckAtoi(slice[7], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[6] = CheckAtoi(slice[8], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[7] = CheckAtoi(slice[9], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[8] = CheckAtoi(slice[10], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[9] = CheckAtoi(slice[11], i)
		GT_Robot_List[RobotID].Heros[heroindex].Propertys[10] = CheckAtoi(slice[12], i)
		heroindex += 1
	}

	return
}

func GetRobot(robotid int32) *ST_Robot {
	if robotid >= int32(len(GT_Robot_List)) || robotid <= 0 {
		gamelog.Error("GetRobot Error: invalid robotid :%d", robotid)
		return nil
	}

	if GT_Robot_List[robotid].RobotID != robotid {
		gamelog.Error("GetRobot Error: invalid robotid2 :%d", robotid)
		return nil
	}

	return &GT_Robot_List[robotid]
}

//! 随机机器人
func RandRobot(minLevel int) *ST_Robot {

	robotLst := []ST_Robot{}
	length := len(GT_Robot_List)
	for i := 0; i < length; i++ {
		if GT_Robot_List[i].Level >= minLevel {
			robotLst = append(robotLst, GT_Robot_List[i])
		}
	}

	if len(robotLst) <= 0 {
		gamelog.Error("RandRobot error: list is nil minlevel: %d", minLevel)
		return nil
	}

	randIndex := r.Intn(len(robotLst) - 1)
	randIndex += 1

	return &robotLst[randIndex]
}

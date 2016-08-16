package gamedata

import (
	"gamelog"
)

type ST_HeroLevelInfo struct {
	Quality      int     //品质
	Level        int     //等级
	NeedExp      int     //升级所需要的经验
	MainNeedExp  int     //主角升级需要的经验
	TotalNeedExp int     //累积所需经验
	TotalMoney   int     //累积所需货币
	MoneyID      int     //需要的货币ID
	MoneyNum     int     //需要的货币系数*经验值
	Propertys    [11]int //属性值
}

var (
	GT_HeroLevel_List [11][101]ST_HeroLevelInfo
	G_HeroMaxLevel    = 0
)

func InitHeroLevelParser(total int) bool {
	//GT_HeroLevel_List = make([]ST_HeroLevelInfo, total+1)
	return true
}

func ParseHeroLevelRecord(rs *RecordSet) {
	Quality := rs.GetFieldInt("quality")
	Level := rs.GetFieldInt("level")

	GT_HeroLevel_List[Quality][Level].Quality = Quality
	GT_HeroLevel_List[Quality][Level].Level = Level
	GT_HeroLevel_List[Quality][Level].NeedExp = rs.GetFieldInt("needexp")
	GT_HeroLevel_List[Quality][Level].MainNeedExp = rs.GetFieldInt("main_needexp")
	GT_HeroLevel_List[Quality][Level].TotalNeedExp = rs.GetFieldInt("total_needexp")
	GT_HeroLevel_List[Quality][Level].TotalMoney = rs.GetFieldInt("total_money")
	GT_HeroLevel_List[Quality][Level].MoneyID = rs.GetFieldInt("money_id")
	GT_HeroLevel_List[Quality][Level].MoneyNum = rs.GetFieldInt("money_num")
	GT_HeroLevel_List[Quality][Level].Propertys[0] = rs.GetFieldInt("p1")
	GT_HeroLevel_List[Quality][Level].Propertys[1] = rs.GetFieldInt("p2")
	GT_HeroLevel_List[Quality][Level].Propertys[2] = rs.GetFieldInt("p3")
	GT_HeroLevel_List[Quality][Level].Propertys[3] = rs.GetFieldInt("p4")
	GT_HeroLevel_List[Quality][Level].Propertys[4] = rs.GetFieldInt("p5")

	if GT_HeroLevel_List[Quality][Level].NeedExp <= 0 {
		panic("field needexp is not a valid equip id!")
	}

	if GT_HeroLevel_List[Quality][Level].MoneyID <= 0 {
		panic("field money_id is not a valid money id!")
	}

	if GT_HeroLevel_List[Quality][Level].MoneyNum <= 0 {
		panic("field money_num  can't be zero!")
	}

	if Level > G_HeroMaxLevel {
		G_HeroMaxLevel = Level
	}

	return
}

func GetHeroLevelInfo(quality int, level int) *ST_HeroLevelInfo {
	if quality >= 8 || quality <= 0 {
		gamelog.Error("GetHeroLevelInfo Error : Invalid quality :%d", quality)
		return nil
	}

	if level > G_HeroMaxLevel || level <= 0 {
		gamelog.Error("GetHeroLevelInfo Error : Invalid level :%d", level)
		return nil
	}

	return &GT_HeroLevel_List[quality][level]
}

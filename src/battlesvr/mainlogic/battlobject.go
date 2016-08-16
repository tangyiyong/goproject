package mainlogic

import (
	"battlesvr/gamedata"
)

type TBuffItem struct {
	ID      int   //BuffID
	EndTime int64 //buff结束时间
}

type TSkillItem struct {
	ID   int   //技能ID
	Time int64 //技能施放时间
}

type THeroObj struct {
	ObjectID        int          //实例ID
	HeroID          int          //英雄ID
	SkiLvl          int          //技能等级
	Camp            int          //英雄的阵营
	PropertyValue   [11]int      //数值属性
	PropertyPercent [11]int      //百分比属性
	CampDef         [5]int       //抗阵营属性
	CampKill        [5]int       //灭阵营属性
	Position        [5]float32   //英雄的坐标(x,y,z,d,v, x,y,z 主向 速度)
	CurProperty     [11]int      //英雄当前的属性
	CurHp           int          //当前的生命值
	AttackPID       int          //攻击属性ID
	SkillState      TSkillItem   //技能施放状态
	BuffLst         [4]TBuffItem //英雄受的BUFF状态
}

type TBattleObj struct {
	PlayerID int //玩家ID
	Level    int //玩家等级
	BatCamp  int //战斗阵营
	HeroObj  [6]THeroObj

	//以下为功能属性
	MoveEndTime int64         //搬水晶结束时间
	SeriesKill  int           //连续杀人数
	SkillState  [4]TSkillItem //四个玩家可以施放的技能
}

func (self *THeroObj) CalcCurProperty(initcurhp bool) {
	if self.HeroID == 0 {
		return
	}

	var i int = 0
	for ; i < 7; i++ {
		self.CurProperty[i] = self.PropertyValue[i] + self.PropertyValue[i]*self.PropertyPercent[i]/1000
	}

	i = 7
	for ; i < 11; i++ {
		self.CurProperty[i] = self.PropertyValue[i] + self.PropertyPercent[i]
	}

	if initcurhp == true {
		self.CurHp = self.CurProperty[0]
	}

	return
}

func (self *TBattleObj) IsAllDie() bool {
	for i := 0; i < 6; i++ {
		if self.HeroObj[i].HeroID <= 0 {
			break
		}

		if self.HeroObj[i].CurHp > 0 {
			return false
		}
	}

	return true
}

func (self *TBattleObj) IsTeamIn(rc *gamedata.TRect) bool {
	for i := 0; i < 6; i++ {
		if self.HeroObj[i].HeroID <= 0 {
			break
		}

		if self.HeroObj[i].CurHp <= 0 {
			continue
		}

		if self.HeroObj[i].Position[0] > 0 {
			return true
		}
	}

	return false
}

func (self *TBattleObj) GetNewSkill() int {
	self.SkillState[3].ID = gamedata.RandSkill()
	return self.SkillState[3].ID
}

func (self *TBattleObj) InitSkillState() bool {
	for i := 0; i < 3; i++ {
		self.SkillState[i].ID = gamedata.GetSceneInfo().SkillFix[i]
	}

	self.SkillState[3].ID = gamedata.RandSkill()

	return true
}

//玩家重生
func (self *TBattleObj) Revive() bool {

	return true
}

func (self *THeroObj) AddBuff(id int) {

	return
}

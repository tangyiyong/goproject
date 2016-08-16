package gamedata

func LoadConfig() bool {
	//加载战场数据
	LoadScene()

	//加载技能配制
	LoadSkills()

	//加载BUFF配制
	LoadBuffs()

	return true
}

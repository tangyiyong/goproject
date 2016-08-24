package gamedata

type ST_BaGuaJingConfig struct {
	BaseMoneyID1  int //! 基础消耗货币1
	BaseMoneyNum1 int
	BaseMoneyID2  int //! 基础消耗货币2
	BaseMoneyNum2 int

	CrossFactionNeedLevel int //! 跨阵营转换等级限制

	LimitQuality                 int8 //! 品质限制
	QualityHeroExchangeNeedLevel int  //! 转换红色武将等级限制
}

var GT_BGJConfig ST_BaGuaJingConfig

func InitBaGuaJingParser(total int) bool {
	return true
}

func ParseBaGuaJingRecord(rs *RecordSet) {
	GT_BGJConfig.BaseMoneyID1 = rs.GetFieldInt("base_money_id1")
	GT_BGJConfig.BaseMoneyNum1 = rs.GetFieldInt("base_money_num1")
	GT_BGJConfig.BaseMoneyID2 = rs.GetFieldInt("base_money_id2")
	GT_BGJConfig.BaseMoneyNum2 = rs.GetFieldInt("base_money_num2")
	GT_BGJConfig.CrossFactionNeedLevel = rs.GetFieldInt("cross_faction_need_level")
	GT_BGJConfig.LimitQuality = int8(rs.GetFieldInt("limit_quality"))
	GT_BGJConfig.QualityHeroExchangeNeedLevel = rs.GetFieldInt("quality_hero_exchange_need_level")
}

func GetBaGuaJingConfigData() *ST_BaGuaJingConfig {
	return &GT_BGJConfig
}

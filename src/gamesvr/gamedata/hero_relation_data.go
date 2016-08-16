package gamedata

import (
	"gamelog"
)

const (
	RTYPE_HERO  = 1 //和英雄
	RTYPE_EQIUP = 2 //和装备
	RTYPE_GEM   = 3 //和宝物
)

type ST_RelationBuff struct {
	ID            int    //RelationBuffID
	PropertyID    [2]int //属性ID
	PropertyValue [2]int //普通值
	IsPercent     bool   //是否是百分比值
}

type ST_RelationItem struct {
	RelationType   int    //羁绊类型 1:和英雄发生的， 2: 和装备发生的; 3:和宝物发生的
	TargetIDs      []int  //如果Type:1 ， 则是英雄的ID集, 否则是装备的ID集或者是宝物ID集
	QualityLimit   [2]int //品质限制
	RelationBuffID int    //天赋ID
}

type ST_HeroRelationInfo struct {
	HeroID    int               //英雄ID
	Relations []ST_RelationItem //羁绊表
}

var (
	GT_HeroRelation_List     []ST_HeroRelationInfo = nil
	GT_HeroRelationBuff_List []ST_RelationBuff     = nil
)

func InitHeroRelationParser(total int) bool {
	GT_HeroRelation_List = make([]ST_HeroRelationInfo, 1500)
	return true
}

func InitHeroRelationBuffParser(total int) bool {
	GT_HeroRelationBuff_List = make([]ST_RelationBuff, total+1)
	return true
}

func ParseHeroRelationRecord(rs *RecordSet) {
	heroID := rs.GetFieldInt("hero_id")
	GT_HeroRelation_List[heroID].HeroID = heroID

	var item ST_RelationItem
	item.RelationBuffID = rs.GetFieldInt("relation_id")
	item.RelationType = rs.GetFieldInt("relation_type")
	item.QualityLimit = ParseTo2IntSlice(rs.GetFieldString("quality"))
	item.TargetIDs = ParseToIntSlice(rs.GetFieldString("target_id"))
	GT_HeroRelation_List[heroID].Relations = append(GT_HeroRelation_List[heroID].Relations, item)

	return
}

func ParseHeroRelationBuffRecord(rs *RecordSet) {
	BuffID := rs.GetFieldInt("id")
	GT_HeroRelationBuff_List[BuffID].ID = BuffID
	GT_HeroRelationBuff_List[BuffID].PropertyID[0] = rs.GetFieldInt("propertyid1")
	GT_HeroRelationBuff_List[BuffID].PropertyValue[0] = rs.GetFieldInt("propertyvalue1")
	GT_HeroRelationBuff_List[BuffID].PropertyID[1] = rs.GetFieldInt("propertyid2")
	GT_HeroRelationBuff_List[BuffID].PropertyValue[1] = rs.GetFieldInt("propertyvalue2")
	GT_HeroRelationBuff_List[BuffID].IsPercent = (1 == rs.GetFieldInt("is_percent"))

	return
}

func GetHeroRelationInfo(heroID int) *ST_HeroRelationInfo {
	if heroID >= len(GT_HeroRelation_List) || heroID <= 0 {
		gamelog.Error("GetHeroRelationInfo Error : Invalid heroID :%d", heroID)
		return nil
	}

	if heroID != GT_HeroRelation_List[heroID].HeroID {
		gamelog.Error("GetHeroRelationInfo Error : Invalid heroID2 :%d", heroID)
		return nil
	}

	return &GT_HeroRelation_List[heroID]
}

func GetHeroRelationItems(heroID int) []ST_RelationItem {
	if heroID >= len(GT_HeroRelation_List) || heroID <= 0 {
		gamelog.Error("GetHeroRelationItems Error : Invalid heroID :%d", heroID)
		return nil
	}

	return GT_HeroRelation_List[heroID].Relations
}

func GetHeroRelationBuff(relationbuff int) *ST_RelationBuff {
	if relationbuff >= len(GT_HeroRelationBuff_List) || relationbuff <= 0 {
		gamelog.Error("GetHeroRelationBuff Error : Invalid relationbuff :%d", relationbuff)
		return nil
	}

	return &GT_HeroRelationBuff_List[relationbuff]
}

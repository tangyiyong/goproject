package gamedata

import (
	"encoding/xml"
	"gamelog"
	"io/ioutil"
	"os"
	"sort"
	"utility"
)

type ST_BuffInfo struct {
	ID       int32      `xml:"id,attr"`       //技能ID
	CD       int32      `xml:"cd,attr"`       //技能CD
	Duration int32      `xml:"duration,attr"` //技能时长
	Radius   int32      `xml:"radius,attr"`   //半径
	BuffID   int32      `xml:"buffid,attr"`   //BuffID
	Hurts    []ST_Hurts `xml:"Hurt"`
}

type ST_BuffMgr struct {
	Buffs []ST_BuffInfo `xml:"Buff"`
}

var G_BuffMgr ST_BuffMgr

func LoadBuffs() bool {
	filepath := utility.GetCurrPath() + "battle/buff.xml"
	file, err := os.Open(filepath)
	if err != nil {
		gamelog.Error("LoadBuffs Error: %v", err)
		return false
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		gamelog.Error("LoadBuffs Error: %v", err)
		return false
	}

	err = xml.Unmarshal(data, &G_BuffMgr)
	if err != nil {
		gamelog.Error("LoadBuffs Error: %v", err)
		return false
	}

	return true
}

func (self *ST_BuffMgr) GetBuffInfo(id int32) *ST_BuffInfo {
	i := sort.Search(len(self.Buffs), func(i int) bool { return self.Buffs[i].ID >= id })
	if i < len(self.Buffs) && self.Buffs[i].ID == id {
		return &self.Buffs[i]
	}

	return nil
}

func RandSkill() int32 {
	var id int32 = 31000021
	rvalue := utility.Rand()
	nCount := len(G_SceneInfo.SkillRand)
	nIndex := rvalue % nCount
	if nCount > 0 {
		id = G_SceneInfo.SkillRand[nIndex]
	}
	return id
}

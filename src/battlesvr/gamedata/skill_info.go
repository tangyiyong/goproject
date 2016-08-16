package gamedata

import (
	"encoding/xml"
	"gamelog"
	"io/ioutil"
	"os"
	"sort"
	"utility"
)

type ST_Hurts struct {
	Level   int `xml:"level,attr"`
	Percent int `xml:"percent,attr"`
	Fixed   int `xml:"fixed,attr"`
}

type ST_SkillInfo struct {
	ID       int        `xml:"id,attr"`       //技能ID
	CD       int        `xml:"cd,attr"`       //技能CD
	Duration int        `xml:"duration,attr"` //技能时长
	Radius   int        `xml:"radius,attr"`   //半径
	BuffID   int        `xml:"buffid,attr"`   //BuffID
	Hurts    []ST_Hurts `xml:"Hurt"`
}

type ST_SkillMgr struct {
	Skills []ST_SkillInfo `xml:"Skill"`
}

var G_SkillMgr ST_SkillMgr

func LoadSkills() bool {
	filepath := utility.GetCurrPath() + "battle/skill.xml"
	file, err := os.Open(filepath)
	if err != nil {
		gamelog.Error("LoadSkills Error: %v", err)
		return false
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		gamelog.Error("LoadSkills Error: %v", err)
		return false
	}

	err = xml.Unmarshal(data, &G_SkillMgr)
	if err != nil {
		gamelog.Error("LoadSkills Error: %v", err)
		return false
	}

	return true
}

func GetSkillInfo(id int) *ST_SkillInfo {
	i := sort.Search(len(G_SkillMgr.Skills), func(i int) bool { return G_SkillMgr.Skills[i].ID >= id })
	if i < len(G_SkillMgr.Skills) && G_SkillMgr.Skills[i].ID == id {
		return &G_SkillMgr.Skills[i]
	}

	return nil
}

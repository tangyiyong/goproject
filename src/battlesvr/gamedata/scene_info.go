package gamedata

import (
	"encoding/xml"
	"gamelog"
	"io/ioutil"
	"os"
	"utility"
)

type TRect struct {
	Left   float32 `xml:"left,attr"`
	Right  float32 `xml:"right,attr"`
	Top    float32 `xml:"top,attr"`
	Bottom float32 `xml:"bottom,attr"`
}

type TPoint struct {
	X float32 `xml:"x,attr"`
	Y float32 `xml:"y,attr"`
	Z float32 `xml:"z,attr"`
}

type TCamp struct {
	BatCamp   int32  `xml:"batcamp,attr"` //阵营战阵营
	SafeRect  TRect  `xml:"Safe"`
	MoveBegin TRect  `xml:"MoveBegin"`
	MoveEnd   TRect  `xml:"MoveEnd"`
	BornPt    TPoint `xml:"BornPt"`
}

type TSceneInfo struct {
	SceneRect TRect   `xml:"Rect"`
	Camps     []TCamp `xml:"Camp"`
	SkillFix  []int32 `xml:"SkillFix"`
	SkillRand []int32 `xml:"SkillRand"`
	HeroSpeed float32 `xml:"HeroSpeed"`
}

var G_SceneInfo TSceneInfo

func GetCampHeroPos(batcamp int32) (pos [5]float32) {
	if batcamp <= 0 || batcamp > int32(len(G_SceneInfo.Camps)) {
		gamelog.Error("GetCampHeroPos Error : Invalid BatCamp:%d", batcamp)
		return
	}

	pos[0] = G_SceneInfo.Camps[batcamp-1].BornPt.X + 0.9*float32(utility.Rand()%3)
	pos[1] = G_SceneInfo.Camps[batcamp-1].BornPt.Y
	pos[2] = G_SceneInfo.Camps[batcamp-1].BornPt.Z + 0.9*float32(utility.Rand()%3)

	return
}

func LoadScene() bool {
	filepath := utility.GetCurrPath() + "battle/scene.xml"
	file, err := os.Open(filepath)
	if err != nil {
		gamelog.Error("LoadSecene Error: %v", err)
		return false
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		gamelog.Error("LoadSecene Error: %v", err)
		return false
	}

	err = xml.Unmarshal(data, &G_SceneInfo)
	if err != nil {
		gamelog.Error("LoadSecene Error: %v", err)
		return false
	}

	return true
}

func GetSceneInfo() *TSceneInfo {

	return &G_SceneInfo
}

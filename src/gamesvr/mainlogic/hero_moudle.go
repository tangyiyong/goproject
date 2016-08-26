package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

const (
	BATTLE_NUM = 6 //上阵英雄数
	BACK_NUM   = 6 //援军英雄数
	EQUIP_NUM  = BATTLE_NUM * 4
	GEM_NUM    = BATTLE_NUM * 2
)

type THeroMoudle struct {
	PlayerID    int32                 `bson:"_id"` //玩家ID
	CurHeros    [BATTLE_NUM]THeroData //上阵英雄
	BackHeros   [BATTLE_NUM]THeroData //援军英雄
	CurEquips   [EQUIP_NUM]TEquipData //上阵装备
	CurGems     [GEM_NUM]TGemData     //上阵宝物
	CurPets     [BATTLE_NUM]TPetData  //上阵宠物
	GuildSkiLvl [11]int8              //公会技能等级
	TitleID     int                   //称号ID
	FashionID   int                   //时装ID
	FashionLvl  int                   //时装等级

	//其它系统添加的固定增加属性
	//宠物图鉴,  时装图鉴， 将灵， 阵图
	ExtraProValue   [11]int32 //增加的数值属性
	ExtraProPercent [11]int32 //增加的百分比属性
	ExtraCampDef    [5]int32  //抗阵营属性  6:号属性
	ExtraCampKill   [5]int32  //灭阵营属性  7:号属性

	ownplayer *TPlayer //父player指针
}

func (self *THeroMoudle) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//OnCreate 响应角色创建
func (self *THeroMoudle) OnCreate(playerid int32) {
	self.PlayerID = playerid
	if self.CurHeros[0].ID <= 0 {
		gamelog.Error("Create Hero Moudle Failed, Hero is 0 !")
		return
	}

	self.CurHeros[0].Level = 1
	self.CurHeros[0].CurExp = 0
	self.CurHeros[0].Quality = gamedata.GetHeroQuality(self.CurHeros[0].ID)
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerHero", self)
}

//OnDestroy player销毁
func (self *THeroMoudle) OnDestroy(playerid int32) {

}

//OnPlayerOnline player进入游戏
func (self *THeroMoudle) OnPlayerOnline(playerid int32) {

}

//OnPlayerOffline player 离开游戏
func (self *THeroMoudle) OnPlayerOffline(playerid int32) {

}

//OnLoad 从数据库中加载
func (self *THeroMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) bool {
	s := mongodb.GetDBSession()
	defer s.Close()
	var bRet = true

	err := s.DB(appconfig.GameDbName).C("PlayerHero").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerHero Load Error :%s， PlayerID: %d", err.Error(), playerid)
		bRet = false
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid

	return bRet
}

type THeroResult struct {
	HeroID           int       //英雄ID
	Quality          int8      //品质
	Camp             int8      //阵营
	PropertyValues   [11]int32 //属性数值
	PropertyPercents [11]int32 //属性增加百分比
	CampDef          [5]int32  //抗阵营属性
	CampKill         [5]int32  //灭阵营属性
}

func (self *THeroMoudle) CalcGem(HeroResults []THeroResult, heroIndex int) bool {
	var begin = heroIndex * 2
	var end = heroIndex*2 + 2

	var minStengthLevel = 100
	var minRefineLevel = 100
	for i := begin; i < end; i++ {
		var pGemData = &self.CurGems[i]
		if pGemData.ID == 0 {
			continue
		}

		pGemInfo := gamedata.GetGemInfo(pGemData.ID)
		pStrengthenInfo := gamedata.GetStrengthInfo(pGemInfo.Quality)

		//计算宝物强化影响
		pid1 := pGemInfo.StrengthPropertys[0]
		pid2 := pGemInfo.StrengthPropertys[1]
		pid1StrengthInc := int32(pStrengthenInfo.PropertyInc[pGemInfo.Position-1][0]*pGemData.StrengLevel + pGemInfo.BasePropertys[0])
		pid2StrengthInc := int32(pStrengthenInfo.PropertyInc[pGemInfo.Position-1][1]*pGemData.StrengLevel + pGemInfo.BasePropertys[1])
		if pid1 == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += pid1StrengthInc
			HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += pid1StrengthInc
		} else if pid1 == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += pid1StrengthInc
			HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += pid1StrengthInc
		} else {
			HeroResults[heroIndex].PropertyValues[pid1-1] += pid1StrengthInc
		}

		if pid2 == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += pid2StrengthInc
			HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += pid2StrengthInc
		} else if pid2 == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += pid2StrengthInc
			HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += pid2StrengthInc
		} else {
			HeroResults[heroIndex].PropertyValues[pid2-1] += pid2StrengthInc
		}

		//计算宝物精炼影响
		pRefineInfo := gamedata.GetRefineInfo(pGemInfo.Quality)
		if pRefineInfo == nil {
			gamelog.Error("CalcGem Error Invalid pGemInfo.Quality :%d", pGemInfo.Quality)
			return false
		}
		pid1 = pGemInfo.RefinePropertys[0]
		pid2 = pGemInfo.RefinePropertys[1]
		pid1RefineInc := int32(pRefineInfo.PropertyInc[pGemInfo.Position-1][0] * pGemData.RefineLevel)
		pid2RefineInc := int32(pRefineInfo.PropertyInc[pGemInfo.Position-1][1] * pGemData.RefineLevel)
		if pid1 == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += pid1RefineInc
			HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += pid1RefineInc
		} else if pid1 == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += pid1RefineInc
			HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += pid1RefineInc
		} else {
			HeroResults[heroIndex].PropertyValues[pid1-1] += pid1RefineInc
		}

		if pid2 == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += pid2RefineInc
			HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += pid2RefineInc
		} else if pid2 == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += pid2RefineInc
			HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += pid2RefineInc
		} else {
			HeroResults[heroIndex].PropertyValues[pid2-1] += pid2RefineInc
		}

		if pGemData.RefineLevel < minRefineLevel {
			minRefineLevel = pGemData.RefineLevel
		}

		if pGemData.StrengLevel < minStengthLevel {
			minStengthLevel = pGemData.StrengLevel
		}
	}

	//计算装备大师影响
	//先计算宝物强化大师
	masterPropertys := gamedata.GetMasterInfo(gamedata.MTYPE_Gem_Strength, minStengthLevel)
	for _, m := range masterPropertys {
		HeroResults[heroIndex].PropertyValues[m.PropertyID-1] += int32(m.PropertyInc)
	}
	//再计算宝物精炼大师
	masterPropertys = gamedata.GetMasterInfo(gamedata.MTYPE_Gem_Refine, minRefineLevel)
	for _, m := range masterPropertys {
		HeroResults[heroIndex].PropertyValues[m.PropertyID-1] += int32(m.PropertyInc)
	}

	return true
}

func (self *THeroMoudle) CalcPet(HeroResults []THeroResult, heroIndex int) bool {
	pPetData := &self.CurPets[heroIndex]
	if pPetData.ID <= 0 {
		return false
	}

	pPetInfo := gamedata.GetPetInfo(pPetData.ID)
	if pPetInfo == nil {
		gamelog.Error("CalcPet Error pPetInfo == nil", pPetData.ID)
		return false
	}

	pPetLevelInfo := gamedata.GetPetLevelInfo(pPetData.ID, pPetData.Level)
	if pPetLevelInfo == nil {
		gamelog.Error("CalcPet Error Invalid pPetLevelInfo, petid:%d, level:%d", pPetData.ID, pPetData.Level)
		return false
	}

	pPetStarInfo := gamedata.GetPetStarInfo(pPetInfo.Quality, pPetData.Star)
	if pPetStarInfo == nil {
		gamelog.Error("CalcPet Error Invalid pPetStarInfo, Quality:%d, Star:%d", pPetInfo.Quality, pPetData.Star)
		return false
	}

	for i, v := range pPetLevelInfo.Propertys {
		HeroResults[heroIndex].PropertyValues[i] += int32(v * pPetStarInfo.PropertyTrans / 1000)
	}

	pPetGodInfo := gamedata.GetPetGodInfo(pPetData.ID, pPetData.God)
	if pPetGodInfo == nil {
		gamelog.Error("CalcPet Error Invalid pPetGodInfo, PetID:%d, God:%d", pPetInfo.PetID, pPetData.God)
		return false
	}

	for _, v := range pPetGodInfo.Propertys {
		if v.PropertyID > 0 {
			if v.IsPercent {
				if v.PropertyID == gamedata.AttackPropertyID {
					HeroResults[heroIndex].PropertyPercents[gamedata.AttackMagicID-1] += int32(v.Value)
					HeroResults[heroIndex].PropertyPercents[gamedata.AttackPhysicID-1] += int32(v.Value)
				} else if v.PropertyID == gamedata.DefencePropertyID {
					HeroResults[heroIndex].PropertyPercents[gamedata.DefenceMagicID-1] += int32(v.Value)
					HeroResults[heroIndex].PropertyPercents[gamedata.DefencePhysicID-1] += int32(v.Value)
				} else {
					HeroResults[heroIndex].PropertyPercents[v.PropertyID-1] += int32(v.Value)
				}
			} else {
				if v.PropertyID == gamedata.AttackPropertyID {
					HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += int32(v.Value)
					HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += int32(v.Value)
				} else if v.PropertyID == gamedata.DefencePropertyID {
					HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += int32(v.Value)
					HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += int32(v.Value)
				} else {
					HeroResults[heroIndex].PropertyValues[v.PropertyID-1] += int32(v.Value)
				}
			}
		}
	}

	return true
}

func (self *THeroMoudle) CalcEquip(HeroResults []THeroResult, heroIndex int) bool {
	var begin = heroIndex * 4
	var end = heroIndex*4 + 4

	var SuitMap map[int]int
	SuitMap = make(map[int]int, 1)

	var minStengthLevel = 100
	var minRefineLevel = 100
	for i := begin; i < end; i++ {
		var pEquipData = &self.CurEquips[i]
		if pEquipData.ID == 0 {
			continue
		}

		pEquipInfo := gamedata.GetEquipmentInfo(pEquipData.ID)
		pStrengthenInfo := gamedata.GetStrengthInfo(pEquipInfo.Quality)
		nStrengInc := int32(pStrengthenInfo.PropertyInc[pEquipInfo.Position-1][0]*pEquipData.StrengLevel + pEquipInfo.BaseProperty)
		//计算装备强化影响
		if pEquipInfo.StrengthProperty == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += nStrengInc
			HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += nStrengInc
		} else if pEquipInfo.StrengthProperty == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += nStrengInc
			HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += nStrengInc
		} else {
			HeroResults[heroIndex].PropertyValues[pEquipInfo.StrengthProperty-1] += nStrengInc
		}

		//计算装备精炼影响
		pRefineInfo := gamedata.GetRefineInfo(pEquipInfo.Quality)
		if pRefineInfo == nil {
			gamelog.Error("CalcEquip Error Invalid pEquipInfo.Quality :%d", pEquipInfo.Quality)
			return false
		}
		pidRefine1 := pEquipInfo.RefinePropertys[0]
		pidRefine2 := pEquipInfo.RefinePropertys[1]
		nRefineInc1 := int32(pRefineInfo.PropertyInc[pEquipInfo.Position-1][0] * pEquipData.RefineLevel)
		nRefineInc2 := int32(pRefineInfo.PropertyInc[pEquipInfo.Position-1][1] * pEquipData.RefineLevel)
		if pidRefine1 == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += nRefineInc1
			HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += nRefineInc1
		} else if pidRefine1 == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += nRefineInc1
			HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += nRefineInc1
		} else {
			HeroResults[heroIndex].PropertyValues[pidRefine1-1] += nRefineInc1
		}

		if pidRefine2 == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += nRefineInc2
			HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += nRefineInc2
		} else if pidRefine2 == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += nRefineInc2
			HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += nRefineInc2
		} else {
			HeroResults[heroIndex].PropertyValues[pidRefine2-1] += nRefineInc2
		}

		//计算装备升星影响
		pEquipStarInfo := gamedata.GetEquipStarInfo(pEquipInfo.Quality, pEquipInfo.Position, pEquipData.Star)
		if pEquipStarInfo != nil {
			var StarProInc = int32(pEquipStarInfo.NeedExp/pEquipStarInfo.AddExp*pEquipStarInfo.AddProperty + pEquipStarInfo.SumProperty)
			if pEquipStarInfo.PropertyID == gamedata.AttackPropertyID {
				HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += StarProInc
				HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += StarProInc
			} else if pEquipStarInfo.PropertyID == gamedata.DefencePropertyID {
				HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += StarProInc
				HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += StarProInc
			} else {
				HeroResults[heroIndex].PropertyValues[pEquipStarInfo.PropertyID-1] += StarProInc
			}

		}

		//统计套装数目
		Num, ok := SuitMap[pEquipInfo.SuitID]
		if ok {
			SuitMap[pEquipInfo.SuitID] = Num + 1
		} else {
			SuitMap[pEquipInfo.SuitID] = 1
		}

		if pEquipData.RefineLevel < minRefineLevel {
			minRefineLevel = pEquipData.RefineLevel
		}

		if pEquipData.StrengLevel < minStengthLevel {
			minStengthLevel = pEquipData.StrengLevel
		}
	}

	//计算套装影响
	for suitid, num := range SuitMap {
		if num <= 1 {
			continue
		}

		SuitBuffs := gamedata.GetEquipSuitBuff(suitid, num)
		for i := 0; i < len(SuitBuffs); i++ {
			self.CalcEquipSuitBuff(HeroResults, heroIndex, &SuitBuffs[i])
		}
	}
	//计算装备大师影响
	//先计算装备强化大师
	masterPropertys := gamedata.GetMasterInfo(gamedata.MTYPE_Equip_Strength, minStengthLevel)
	for _, m := range masterPropertys {
		HeroResults[heroIndex].PropertyValues[m.PropertyID-1] += int32(m.PropertyInc)
	}

	//再计算装备精炼大师
	masterPropertys = gamedata.GetMasterInfo(gamedata.MTYPE_Equip_Refine, minRefineLevel)
	for _, m := range masterPropertys {
		HeroResults[heroIndex].PropertyValues[m.PropertyID-1] += int32(m.PropertyInc)
	}
	return true
}

func (self *THeroMoudle) CalcEquipSuitBuff(HeroResults []THeroResult, heroIndex int, pSuitBuff *gamedata.ST_EquipSuitBuff) bool {
	if pSuitBuff == nil {
		gamelog.Error("CalcEquipSuitBuff Error pSuitBuff is Nil")
		return true
	}

	if pSuitBuff.IsPercent {
		if pSuitBuff.PropertyID == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyPercents[gamedata.AttackMagicID-1] += int32(pSuitBuff.PropertyValue)
			HeroResults[heroIndex].PropertyPercents[gamedata.AttackPhysicID-1] += int32(pSuitBuff.PropertyValue)
		} else if pSuitBuff.PropertyID == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyPercents[gamedata.DefenceMagicID-1] += int32(pSuitBuff.PropertyValue)
			HeroResults[heroIndex].PropertyPercents[gamedata.DefencePhysicID-1] += int32(pSuitBuff.PropertyValue)
		} else {
			HeroResults[heroIndex].PropertyPercents[pSuitBuff.PropertyID-1] += int32(pSuitBuff.PropertyValue)
		}
	} else {
		if pSuitBuff.PropertyID == gamedata.AttackPropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += int32(pSuitBuff.PropertyValue)
			HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += int32(pSuitBuff.PropertyValue)
		} else if pSuitBuff.PropertyID == gamedata.DefencePropertyID {
			HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += int32(pSuitBuff.PropertyValue)
			HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += int32(pSuitBuff.PropertyValue)
		} else {
			HeroResults[heroIndex].PropertyValues[pSuitBuff.PropertyID-1] += int32(pSuitBuff.PropertyValue)
		}
	}

	return true
}

func (self *THeroMoudle) CalcTalentItem(HeroResults []THeroResult, heroIndex int, talentid int) bool {
	var pItem = gamedata.GetTalentInfo(talentid)
	if pItem == nil {
		gamelog.Error("CalcTalentItem Error pItem is Nil")
		return true
	}

	var Value int32 = 0
	if HeroResults[heroIndex].Quality > 4 {
		Value = int32(pItem.PropertyValue2)
	} else {
		Value = int32(pItem.PropertyValue1)
	}

	if pItem.TargetType == gamedata.TargetType_Self {
		if pItem.IsPercent {
			if pItem.PropertyID == gamedata.AttackPropertyID {
				HeroResults[heroIndex].PropertyPercents[gamedata.AttackMagicID-1] += Value
				HeroResults[heroIndex].PropertyPercents[gamedata.AttackPhysicID-1] += Value
			} else if pItem.PropertyID == gamedata.DefencePropertyID {
				HeroResults[heroIndex].PropertyPercents[gamedata.DefenceMagicID-1] += Value
				HeroResults[heroIndex].PropertyPercents[gamedata.DefencePhysicID-1] += Value
			} else {
				HeroResults[heroIndex].PropertyPercents[pItem.PropertyID-1] += Value
			}
		} else {
			if pItem.PropertyID == gamedata.AttackPropertyID {
				HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += Value
				HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += Value
			} else if pItem.PropertyID == gamedata.DefencePropertyID {
				HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += Value
				HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += Value
			} else {
				HeroResults[heroIndex].PropertyValues[pItem.PropertyID-1] += Value
			}
		}
	} else if pItem.TargetType == gamedata.TargetType_Friend {
		for j := 0; j < BATTLE_NUM; j++ {
			if HeroResults[j].HeroID == 0 {
				continue
			}
			if pItem.IsPercent {
				if pItem.PropertyID == gamedata.AttackPropertyID {
					HeroResults[j].PropertyPercents[gamedata.AttackMagicID-1] += Value
					HeroResults[j].PropertyPercents[gamedata.AttackPhysicID-1] += Value
				} else if pItem.PropertyID == gamedata.DefencePropertyID {
					HeroResults[j].PropertyPercents[gamedata.DefenceMagicID-1] += Value
					HeroResults[j].PropertyPercents[gamedata.DefencePhysicID-1] += Value
				} else {
					HeroResults[j].PropertyPercents[pItem.PropertyID-1] += Value
				}
			} else {
				if pItem.PropertyID == gamedata.AttackPropertyID {
					HeroResults[j].PropertyValues[gamedata.AttackMagicID-1] += Value
					HeroResults[j].PropertyValues[gamedata.AttackPhysicID-1] += Value
				} else if pItem.PropertyID == gamedata.DefencePropertyID {
					HeroResults[j].PropertyValues[gamedata.DefenceMagicID-1] += Value
					HeroResults[j].PropertyValues[gamedata.DefencePhysicID-1] += Value
				} else {
					HeroResults[j].PropertyValues[pItem.PropertyID-1] += Value
				}
			}
		}
	} else if pItem.TargetType == gamedata.TargetType_Camp {
		for j := 0; j < BATTLE_NUM; j++ {
			if HeroResults[j].HeroID == 0 {
				continue
			}
			if HeroResults[j].Camp == pItem.TargetCamp {
				if pItem.IsPercent {
					if pItem.PropertyID == gamedata.AttackPropertyID {
						HeroResults[j].PropertyPercents[gamedata.AttackMagicID-1] += Value
						HeroResults[j].PropertyPercents[gamedata.AttackPhysicID-1] += Value
					} else if pItem.PropertyID == gamedata.DefencePropertyID {
						HeroResults[j].PropertyPercents[gamedata.DefenceMagicID-1] += Value
						HeroResults[j].PropertyPercents[gamedata.DefencePhysicID-1] += Value
					} else {
						HeroResults[j].PropertyPercents[pItem.PropertyID-1] += Value
					}
				} else {
					if pItem.PropertyID == gamedata.AttackPropertyID {
						HeroResults[j].PropertyValues[gamedata.AttackMagicID-1] += Value
						HeroResults[j].PropertyValues[gamedata.AttackPhysicID-1] += Value
					} else if pItem.PropertyID == gamedata.DefencePropertyID {
						HeroResults[j].PropertyValues[gamedata.DefenceMagicID-1] += Value
						HeroResults[j].PropertyValues[gamedata.DefencePhysicID-1] += Value
					} else {
						HeroResults[j].PropertyValues[pItem.PropertyID-1] += Value
					}
				}
			}
		}
	} else if pItem.TargetType == gamedata.TargetType_Camp_Kill {

	}

	return true
}

func (self *THeroMoudle) CalcRelationBuff(HeroResults []THeroResult, heroIndex int, buffid int) bool {
	var pItem = gamedata.GetHeroRelationBuff(buffid)
	if pItem == nil {
		gamelog.Error("CalcRelationBuff Error Invalid buffid :%d", buffid)
		return false
	}
	for i := 0; i < 2; i++ {
		if pItem.PropertyID[i] != 0 {
			if pItem.IsPercent {
				if pItem.PropertyID[i] == gamedata.AttackPropertyID {
					HeroResults[heroIndex].PropertyPercents[gamedata.AttackMagicID-1] += int32(pItem.PropertyValue[i])
					HeroResults[heroIndex].PropertyPercents[gamedata.AttackPhysicID-1] += int32(pItem.PropertyValue[i])
				} else if pItem.PropertyID[i] == gamedata.DefencePropertyID {
					HeroResults[heroIndex].PropertyPercents[gamedata.DefenceMagicID-1] += int32(pItem.PropertyValue[i])
					HeroResults[heroIndex].PropertyPercents[gamedata.DefencePhysicID-1] += int32(pItem.PropertyValue[i])
				} else {
					HeroResults[heroIndex].PropertyPercents[pItem.PropertyID[i]-1] += int32(pItem.PropertyValue[i])
				}
			} else {
				if pItem.PropertyID[i] == gamedata.AttackPropertyID {
					HeroResults[heroIndex].PropertyValues[gamedata.AttackMagicID-1] += int32(pItem.PropertyValue[i])
					HeroResults[heroIndex].PropertyValues[gamedata.AttackPhysicID-1] += int32(pItem.PropertyValue[i])
				} else if pItem.PropertyID[i] == gamedata.DefencePropertyID {
					HeroResults[heroIndex].PropertyValues[gamedata.DefenceMagicID-1] += int32(pItem.PropertyValue[i])
					HeroResults[heroIndex].PropertyValues[gamedata.DefencePhysicID-1] += int32(pItem.PropertyValue[i])
				} else {
					HeroResults[heroIndex].PropertyValues[pItem.PropertyID[i]-1] += int32(pItem.PropertyValue[i])
				}
			}
		}
	}

	return true
}

func (self *THeroMoudle) IsRelationMatch(pRelationItem *gamedata.ST_RelationItem, heroIndex int) bool {
	if (self.CurHeros[heroIndex].Quality < pRelationItem.QualityLimit[0]) ||
		(self.CurHeros[heroIndex].Quality > pRelationItem.QualityLimit[1]) {
		return false
	}

	if pRelationItem.RelationType == gamedata.RTYPE_HERO {
		for _, heroid := range pRelationItem.TargetIDs {
			for j := 0; j < BATTLE_NUM; j++ {
				if heroid == self.CurHeros[j].ID {
					break
				}
			}
			for j := 0; j < BACK_NUM; j++ {
				if heroid == self.BackHeros[j].ID {
					break
				}
			}
			return false
		}
	} else if pRelationItem.RelationType == gamedata.RTYPE_EQIUP {
		var begin = heroIndex * 4
		var end = heroIndex*4 + 4
		for i := begin; i < end; i++ {
			var pEquipData = &self.CurEquips[i]
			if pEquipData.ID == pRelationItem.TargetIDs[0] {
				return true
			}
			return false
		}
	} else if pRelationItem.RelationType == gamedata.RTYPE_GEM {
		var begin = heroIndex * 2
		var end = heroIndex*2 + 2
		for i := begin; i < end; i++ {
			var pGemData = &self.CurGems[i]
			if pGemData.ID == pRelationItem.TargetIDs[0] {
				return true
			}
			return false
		}
	}

	return true
}

//计算英灵的战力影响
func (self *THeroMoudle) CalcExtraProperty(HeroResults []THeroResult) bool {
	for k := 0; k < BATTLE_NUM; k++ {
		if HeroResults[k].HeroID == 0 {
			continue
		}

		for pid := 0; pid < 11; pid++ {
			HeroResults[k].PropertyPercents[pid] += self.ExtraProPercent[pid]
			HeroResults[k].PropertyValues[pid] += self.ExtraProValue[pid]
		}
	}
	return true
}

//计算时装的属性影响
func (self *THeroMoudle) CalcFashion(HeroResults []THeroResult) bool {
	if self.FashionID <= 0 {
		return false
	}

	pFashionLvlInfo := gamedata.GetFashionLevelInfo(self.FashionID, self.FashionLvl)
	if pFashionLvlInfo == nil {
		gamelog.Error("CalcFashion Error : Invalid ID:%d and lvl:%d", self.FashionID, self.FashionLvl)
		return false
	}

	for k := 0; k < BATTLE_NUM; k++ {
		if HeroResults[k].HeroID == 0 {
			continue
		}

		for pid := 0; pid < 5; pid++ {
			HeroResults[k].PropertyValues[pid] += int32(pFashionLvlInfo.PropertyValues[pid])
			HeroResults[k].PropertyPercents[pid] += int32(pFashionLvlInfo.PropertyPercents[pid])
		}
	}
	return true
}

//计算援军系统的战力影响
func (self *THeroMoudle) CalcHeroFriend(HeroResults []THeroResult) bool {
	minLevel := 200
	for i := 0; i < BACK_NUM; i++ {
		if self.BackHeros[i].ID == 0 {
			return false
		}

		if minLevel > self.BackHeros[i].Level {
			minLevel = self.BackHeros[i].Level
		}
	}

	pHeroFriendInfo := gamedata.GetHeroFriendInfo(minLevel)
	if pHeroFriendInfo == nil {
		return false
	}

	for k := 0; k < BATTLE_NUM; k++ {
		if HeroResults[k].HeroID == 0 {
			continue
		}

		for pid := 0; pid < 5; pid++ {
			HeroResults[k].PropertyPercents[pid] += int32(pHeroFriendInfo.Propertys[pid][1])
			HeroResults[k].PropertyValues[pid] += int32(pHeroFriendInfo.Propertys[pid][0])
		}
	}
	return true
}

//计算称号的战力影响
func (self *THeroMoudle) CalcTitle(HeroResults []THeroResult) bool {
	if self.TitleID <= 0 {
		return false
	}

	pTitleInfo := gamedata.GetTitleInfo(self.TitleID)
	if pTitleInfo == nil {
		gamelog.Error("CalcTitle Error Invalid TitleID :%d", self.TitleID)
		return false
	}

	if pTitleInfo.IsAll {
		//! 作用于全体
		for i := 0; i < 3; i++ {
			if pTitleInfo.Property[i].IsPercent == true {
				//! 若为百分比加成
				for k := 0; k < BATTLE_NUM; k++ {
					if HeroResults[k].HeroID != 0 {
						//! 判断是否为全属性加成
						if pTitleInfo.Property[i].PropertyID == 22 {
							for b := 0; b < 5; b++ {
								HeroResults[k].PropertyPercents[b] += int32(pTitleInfo.Property[i].Value)
							}
						} else {
							HeroResults[k].PropertyPercents[pTitleInfo.Property[i].PropertyID-1] += int32(pTitleInfo.Property[i].Value)
						}

					}
				}
			} else {
				//! 非百分比加成
				for k := 0; k < BATTLE_NUM; k++ {
					if HeroResults[k].HeroID != 0 {
						if pTitleInfo.Property[i].PropertyID == 22 {
							for b := 0; b < 5; b++ {
								HeroResults[k].PropertyValues[b] += int32(pTitleInfo.Property[i].Value)
							}
						} else {
							HeroResults[k].PropertyValues[pTitleInfo.Property[i].PropertyID-1] += int32(pTitleInfo.Property[i].Value)
						}

					}
				}
			}
		}

	} else {
		//! 作用于主英雄
		for i := 0; i < 3; i++ {
			if pTitleInfo.Property[i].IsPercent == true {
				//! 若为百分比加成
				if pTitleInfo.Property[i].PropertyID == 22 {
					for b := 0; b < 5; b++ {
						HeroResults[0].PropertyPercents[b] += int32(pTitleInfo.Property[i].Value)
					}
				} else {
					HeroResults[0].PropertyPercents[pTitleInfo.Property[i].PropertyID-1] += int32(pTitleInfo.Property[i].Value)
				}

			} else {
				//! 非百分比加成
				if pTitleInfo.Property[i].PropertyID == 22 {
					for b := 0; b < 5; b++ {
						HeroResults[0].PropertyValues[b] += int32(pTitleInfo.Property[i].Value)
					}
				} else {
					HeroResults[0].PropertyValues[pTitleInfo.Property[i].PropertyID-1] += int32(pTitleInfo.Property[i].Value)
				}
			}
		}
	}

	return true
}

//计算角色的战力
func (self *THeroMoudle) CalcFightValue(HeroResults []THeroResult) int {
	if HeroResults == nil {
		HeroResults = make([]THeroResult, BATTLE_NUM)
	}
	for heroIndex := 0; heroIndex < BATTLE_NUM; heroIndex++ {
		if self.CurHeros[heroIndex].ID == 0 {
			continue
		}

		pHeroInfo := gamedata.GetHeroInfo(self.CurHeros[heroIndex].ID)
		if pHeroInfo == nil {
			gamelog.Error("CalcFightValue Error Invalid HeroID :%d", self.CurHeros[heroIndex].ID)
			return 0
		}

		if self.CurHeros[heroIndex].Quality == 0 {
			self.CurHeros[heroIndex].Quality = pHeroInfo.Quality
		}

		HeroResults[heroIndex].Camp = pHeroInfo.Camp
		HeroResults[heroIndex].HeroID = pHeroInfo.HeroID
		HeroResults[heroIndex].Quality = self.CurHeros[heroIndex].Quality
	}

	for heroIndex := 0; heroIndex < BATTLE_NUM; heroIndex++ {
		var pCurHeroData = &self.CurHeros[heroIndex]
		if pCurHeroData.ID == 0 {
			continue
		}

		pHeroInfo := gamedata.GetHeroInfo(pCurHeroData.ID)
		//计算等级的影响****************************************
		pHeroLevelInfo := gamedata.GetHeroLevelInfo(pCurHeroData.Quality, pCurHeroData.Level)
		if pHeroLevelInfo == nil {
			gamelog.Error("CalcFightValue Error : Invalid Level :%d", pCurHeroData.Level)
			return 0
		}
		for pid := 0; pid < 5; pid++ {
			HeroResults[heroIndex].PropertyValues[pid] += int32(pHeroInfo.BasePropertys[pid])
			HeroResults[heroIndex].PropertyValues[pid] += int32(pHeroLevelInfo.Propertys[pid])
		}
		//计算等级的影响****************************************

		//计算突破的影响
		pBreakInfo := gamedata.GetHeroBreakInfo(pCurHeroData.BreakLevel)
		if pBreakInfo != nil {
			for pid := 0; pid < 5; pid++ {
				HeroResults[heroIndex].PropertyPercents[pid] += int32(pBreakInfo.IncPercent) //突破加的都是百分比
			}
		}

		//计算培养的数值影响
		for pid := 0; pid < 5; pid++ {
			HeroResults[heroIndex].PropertyValues[pid] += int32(pCurHeroData.Cultures[pid])
		}

		//计算雕文影响
		for pid := 0; pid < 30; pid++ {
			HeroResults[heroIndex].PropertyValues[pid%6] += int32(pCurHeroData.DiaoWenPtys[pid])
		}

		//计算天命影响 (天命都是影响百分比)
		DestinyLevel := pCurHeroData.DestinyState >> 24 & 0x000F
		DestinyIndex := pCurHeroData.DestinyState >> 16 & 0x000F
		var pLastInfo *gamedata.ST_DestinyItem = nil
		if DestinyLevel > 1 {
			pLastInfo = gamedata.GetHeroDestinyInfo(int(DestinyLevel - 1))
			if pLastInfo == nil {
				gamelog.Error("CalcFightValue Error : pLastInfo is nil , cur destiny level:%d", DestinyLevel, DestinyLevel)
			}
		}

		pDestinyInfo := gamedata.GetHeroDestinyInfo(int(DestinyLevel))
		if pDestinyInfo == nil {
			gamelog.Error("CalcFightValue Error : Invalid DestinyLevel :%d, State:%d", DestinyLevel, pCurHeroData.DestinyState)
			return 0
		}

		var lastProInc int32 = 0
		if pLastInfo != nil {
			lastProInc = int32(pLastInfo.PropertyInc)
		}

		var curProInc int32 = int32(pDestinyInfo.PropertyInc)

		if DestinyIndex >= 4 {
			HeroResults[heroIndex].PropertyPercents[0] += curProInc
			HeroResults[heroIndex].PropertyPercents[1] += curProInc
			HeroResults[heroIndex].PropertyPercents[2] += curProInc
			HeroResults[heroIndex].PropertyPercents[3] += curProInc
			HeroResults[heroIndex].PropertyPercents[4] += curProInc
		} else if DestinyIndex == 3 {
			HeroResults[heroIndex].PropertyPercents[0] += curProInc
			HeroResults[heroIndex].PropertyPercents[1] += curProInc
			HeroResults[heroIndex].PropertyPercents[2] += curProInc
			HeroResults[heroIndex].PropertyPercents[3] += curProInc
			HeroResults[heroIndex].PropertyPercents[4] += lastProInc
		} else if DestinyIndex == 2 {
			HeroResults[heroIndex].PropertyPercents[0] += curProInc
			HeroResults[heroIndex].PropertyPercents[1] += curProInc
			HeroResults[heroIndex].PropertyPercents[2] += lastProInc
			HeroResults[heroIndex].PropertyPercents[3] += curProInc
			HeroResults[heroIndex].PropertyPercents[4] += lastProInc
		} else if DestinyIndex == 1 {
			HeroResults[heroIndex].PropertyPercents[0] += curProInc
			HeroResults[heroIndex].PropertyPercents[1] += lastProInc
			HeroResults[heroIndex].PropertyPercents[2] += lastProInc
			HeroResults[heroIndex].PropertyPercents[3] += lastProInc
			HeroResults[heroIndex].PropertyPercents[4] += lastProInc
		}

		///计算觉醒影响  (培养加的都是值)
		for _, wakeitem := range pCurHeroData.WakeItem {
			if wakeitem != 0 {
				pWakeitemInfo := gamedata.GetItemInfo(wakeitem)
				if pWakeitemInfo.Type != gamedata.TYPE_WAKE {
					gamelog.Error("CalcFightValue Error : wakeitem type is not wakeitem :%d", wakeitem)
				} else if pWakeitemInfo != nil {
					HeroResults[heroIndex].PropertyValues[0] += int32(pWakeitemInfo.Propertys[0])
					HeroResults[heroIndex].PropertyValues[1] += int32(pWakeitemInfo.Propertys[1])
					HeroResults[heroIndex].PropertyValues[2] += int32(pWakeitemInfo.Propertys[2])
					HeroResults[heroIndex].PropertyValues[3] += int32(pWakeitemInfo.Propertys[1])
					HeroResults[heroIndex].PropertyValues[4] += int32(pWakeitemInfo.Propertys[2])
				}
			}
		}

		if pCurHeroData.WakeLevel > 1 {
			pWakeLevelInfo := gamedata.GetWakeLevelItem(pCurHeroData.WakeLevel - 1)
			for pid := 0; pid < 11; pid++ {
				HeroResults[heroIndex].PropertyValues[pid] += int32(pWakeLevelInfo.PropertyValues[pid])
				HeroResults[heroIndex].PropertyPercents[pid] += int32(pWakeLevelInfo.PropertyPercents[pid])
			}
		}

		//计算英雄的缘分
		pRelations := gamedata.GetHeroRelationItems(pCurHeroData.ID)
		for r := 0; r < len(pRelations); r++ {
			if self.IsRelationMatch(&pRelations[r], heroIndex) {
				self.CalcRelationBuff(HeroResults, heroIndex, pRelations[r].RelationBuffID)
			}
		}

		//计算突破天赋影响
		pBreakTalent := gamedata.GetHeroBreakTalentInfo(pCurHeroData.ID)
		if pBreakTalent != nil {
			for c := int8(0); c < 15; c++ {
				if pBreakTalent.Talents[c] <= 0 {
					break
				}

				if c >= pCurHeroData.BreakLevel {
					break
				}

				self.CalcTalentItem(HeroResults, heroIndex, pBreakTalent.Talents[c])
			}
		}

		//计算化神的影响
		if pCurHeroData.GodLevel > 0 {
			pHeroGodInfo := gamedata.GetHeroGodInfo(pCurHeroData.GodLevel)
			if pHeroGodInfo != nil {
				HeroResults[heroIndex].PropertyValues[0] += int32(pHeroGodInfo.Propertys[0])
				HeroResults[heroIndex].PropertyValues[1] += int32(pHeroGodInfo.Propertys[1])
				HeroResults[heroIndex].PropertyValues[2] += int32(pHeroGodInfo.Propertys[2])
				HeroResults[heroIndex].PropertyValues[3] += int32(pHeroGodInfo.Propertys[3])
				HeroResults[heroIndex].PropertyValues[4] += int32(pHeroGodInfo.Propertys[4])
			}
		}

		//计算装备影响
		self.CalcEquip(HeroResults, heroIndex)

		//计算宝物的数值影响
		self.CalcGem(HeroResults, heroIndex)

		//计算宠物数值影响
		self.CalcPet(HeroResults, heroIndex)
	}

	//计算称号对战力的影响
	self.CalcTitle(HeroResults)

	//计算其它系统的影响
	self.CalcExtraProperty(HeroResults)

	//计算英雄援军系统的影响
	self.CalcHeroFriend(HeroResults)

	//计算时装系统的影响
	self.CalcFashion(HeroResults)

	//计算公会技能的数值影响
	if self.ownplayer != nil && self.ownplayer.pSimpleInfo.GuildID != 0 {
		for k := 0; k < BATTLE_NUM; k++ {
			if HeroResults[k].HeroID != 0 {
				for pid := 0; pid < 11; pid++ {
					proValue := gamedata.GetGuildSkillValue(int(self.GuildSkiLvl[pid]), pid)
					HeroResults[k].PropertyValues[pid] += int32(proValue)
				}
			}
		}
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	//**
	//**               以下是计算最终战力
	//**
	//<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	for k := 0; k < BATTLE_NUM; k++ {
		if HeroResults[k].HeroID != 0 {
			for p := 0; p < 7; p++ {
				HeroResults[k].PropertyValues[p] += HeroResults[k].PropertyValues[p] * HeroResults[k].PropertyPercents[p] / 1000
			}
			for p := 7; p < 11; p++ {
				HeroResults[k].PropertyValues[p] += HeroResults[k].PropertyPercents[p]
			}
		}
	}

	var fightvalue int = 0
	for k := 0; k < BATTLE_NUM; k++ {
		if HeroResults[k].HeroID != 0 {
			pHeroInfo := gamedata.GetHeroInfo(HeroResults[k].HeroID)
			if (pHeroInfo.AttackType == 1) || (pHeroInfo.AttackType == 3) {
				HeroResults[k].PropertyValues[gamedata.AttackMagicID-1] = 0
			} else {
				HeroResults[k].PropertyValues[gamedata.AttackPhysicID-1] = 0
			}

			//生命
			fightvalue += int(float32(HeroResults[k].PropertyValues[0]) * gamedata.GetPropertyInfo(1).FightFactor)
			//物理攻击
			fightvalue += int(float32(HeroResults[k].PropertyValues[1]) * gamedata.GetPropertyInfo(2).FightFactor)
			//物理防御
			fightvalue += int(float32(HeroResults[k].PropertyValues[2]) * gamedata.GetPropertyInfo(3).FightFactor)
			//魔法攻击
			fightvalue += int(float32(HeroResults[k].PropertyValues[3]) * gamedata.GetPropertyInfo(4).FightFactor)
			//魔法防御
			fightvalue += int(float32(HeroResults[k].PropertyValues[4]) * gamedata.GetPropertyInfo(5).FightFactor)
			//伤害减免
			fightvalue += int(float32(HeroResults[k].PropertyValues[2]/2+HeroResults[k].PropertyValues[4]/2) / gamedata.GetPropertyInfo(6).FightFactor / 10 * float32(HeroResults[k].PropertyValues[5]))
			//伤害加成
			fightvalue += int(float32(HeroResults[k].PropertyValues[1]+HeroResults[k].PropertyValues[3]) / gamedata.GetPropertyInfo(7).FightFactor / 10 * float32(HeroResults[k].PropertyValues[6]))
			//闪避率
			fightvalue += int(float32(HeroResults[k].PropertyValues[4]) / gamedata.GetPropertyInfo(8).FightFactor / 10 * float32(HeroResults[k].PropertyValues[7]))
			//命中率
			fightvalue += int(float32(HeroResults[k].PropertyValues[0]) / gamedata.GetPropertyInfo(9).FightFactor / 10 * float32(HeroResults[k].PropertyValues[8]))
			//暴击率
			fightvalue += int(float32(HeroResults[k].PropertyValues[1]+HeroResults[k].PropertyValues[3]) / gamedata.GetPropertyInfo(10).FightFactor / 10 * float32(HeroResults[k].PropertyValues[9]))
			//抗暴率
			fightvalue += int(float32(HeroResults[k].PropertyValues[2]) / gamedata.GetPropertyInfo(11).FightFactor / 10 * float32(HeroResults[k].PropertyValues[10]))
		}
	}

	return fightvalue
}

func (self *THeroMoudle) UpdateHeroLevel(pHeroData *THeroData) bool {
	var bUpdate = false
	pHeroInfo := gamedata.GetHeroInfo(pHeroData.ID)
	if pHeroInfo == nil {
		gamelog.Error("UpdateHeroLevel Error : Invalid HeroID:%d", pHeroData.ID)
		return false
	}
	for {
		pStHeroLevelInfo := gamedata.GetHeroLevelInfo(pHeroInfo.Quality, pHeroData.Level)
		if pStHeroLevelInfo == nil {
			gamelog.Error("AddHeroExp Error: Invalid HeroID")
			break
		}

		if pHeroData.CurExp < pStHeroLevelInfo.NeedExp {
			break
		}

		if (pHeroData.Level + 1) > self.ownplayer.GetLevel() {
			pHeroData.CurExp = pStHeroLevelInfo.NeedExp
			bUpdate = true
			break
		}

		bUpdate = true
		pHeroData.CurExp -= pStHeroLevelInfo.NeedExp
		pHeroData.Level += 1
	}

	return bUpdate
}

//给主角英雄增加经验
func (self *THeroMoudle) AddMainHeroExp(exp int) int {
	self.CurHeros[0].CurExp += exp
	self.ownplayer.DB_SaveHeroLevelExp(POSTYPE_BATTLE, 0)
	return self.CurHeros[0].CurExp
}

//修改英雄的品质信息
func (self *THeroMoudle) ChangeMainQuality(value int8) bool {
	self.CurHeros[0].Quality = value
	self.ownplayer.DB_SaveHeroQuality(POSTYPE_BATTLE, 0)
	G_SimpleMgr.Set_HeroQuality(self.PlayerID, value)
	return true
}

//给其它的英雄增加经验
func (self *THeroMoudle) AddHeroExp(postype int, heroindex int, exp int) bool {
	if heroindex <= 0 {
		gamelog.Error("AddHeroExp Error: Invalid heroIndex :%d", heroindex)
		return false
	}

	var pHeroData *THeroData
	if postype == POSTYPE_BATTLE {
		pHeroData = &self.CurHeros[heroindex]
	} else if postype == POSTYPE_BAG {
		pHeroData = &self.ownplayer.BagMoudle.HeroBag.Heros[heroindex]
	} else if postype == POSTYPE_BACK {
		pHeroData = &self.BackHeros[heroindex]
	}

	if pHeroData == nil {
		gamelog.Error("AddHeroExp Error: Invalid pHeroData == nil")
		return false
	}

	pHeroData.CurExp += exp
	self.UpdateHeroLevel(pHeroData)
	self.ownplayer.DB_SaveHeroLevelExp(postype, heroindex)

	return true
}

//获取指定位置的上阵英雄
func (self *THeroMoudle) GetBattleHeroByPos(pos int) *THeroData {
	if pos < 0 || pos >= 8 {
		gamelog.Error("GetBattleHeroByPos Error : Invalid pos :%d", pos)
		return nil
	}

	return &self.CurHeros[pos]
}

//设置指定位置的上阵英雄
func (self *THeroMoudle) SetBattleHeroByPos(pos int, pHero *THeroData) bool {
	if pos < 0 || pos >= 8 {
		gamelog.Error("SetBattleHeroByPos Error : Invalid pos :%d", pos)
		return false
	}

	self.CurHeros[pos] = *pHero

	return true
}

//获取指定位置的援军英雄
func (self *THeroMoudle) GetBackHeroByPos(pos int) *THeroData {
	return &self.BackHeros[pos]
}

//设置指定位置的援军英雄
func (self *THeroMoudle) SetBackHeroByPos(pos int, pHero *THeroData) bool {
	if pos < 0 || pos >= len(self.BackHeros) {
		gamelog.Error("SetBackHeroByPos Error : Invalid pos :%d", pos)
		return false
	}

	self.BackHeros[pos] = *pHero

	return true
}

//增加公会技能等级
func (self *THeroMoudle) AddGuildSkillProLevel(pid int) bool {
	if pid == gamedata.AttackPropertyID {
		self.GuildSkiLvl[gamedata.AttackPhysicID-1] += 1
		self.DB_SaveGuildSkill(gamedata.AttackPhysicID - 1)
	} else {
		self.GuildSkiLvl[pid-1] += 1
		self.DB_SaveGuildSkill(pid - 1)
	}

	return true
}

func (self *THeroMoudle) ClearGuildSkillProLevel() {
	for i := 0; i < len(self.GuildSkiLvl); i++ {
		self.GuildSkiLvl[i] = 0
	}

	self.DB_SaveGuildSkillLst()
}

func (self *THeroMoudle) AddExtraProperty(pid int, pvalue int32, percent bool, camp int) {
	if pid <= 0 || pid > 11 {
		gamelog.Error("AddExtraProperty Error : Invalid Pid:%d", pid)
		return
	}

	if camp <= 0 {
		if percent == true {
			if pid == gamedata.AttackPropertyID {
				self.ExtraProPercent[gamedata.AttackPhysicID-1] += pvalue
				self.ExtraProPercent[gamedata.AttackMagicID-1] += pvalue
			} else if pid == gamedata.DefencePropertyID {
				self.ExtraProPercent[gamedata.DefencePhysicID-1] += pvalue
				self.ExtraProPercent[gamedata.DefenceMagicID-1] += pvalue
			} else {
				self.ExtraProPercent[pid-1] += pvalue
			}

		} else {
			if pid == gamedata.AttackPropertyID {
				self.ExtraProValue[gamedata.AttackPhysicID-1] += pvalue
				self.ExtraProValue[gamedata.AttackMagicID-1] += pvalue
			} else if pid == gamedata.DefencePropertyID {
				self.ExtraProValue[gamedata.DefencePhysicID-1] += pvalue
				self.ExtraProValue[gamedata.DefenceMagicID-1] += pvalue
			} else {
				self.ExtraProValue[pid-1] += pvalue
			}
		}
	} else {
		if pid == 6 {
			self.ExtraCampDef[pid-1] += pvalue
		} else if pid == 7 {
			self.ExtraCampKill[pid-1] += pvalue
		} else {
			gamelog.Error("AddExtraProperty Error : if camp:%d != 0, pid :%d should be 6 or 7!")
		}
	}

}

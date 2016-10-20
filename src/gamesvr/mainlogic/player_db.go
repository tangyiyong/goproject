package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (self *TPlayer) DB_SaveHeroAt(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curheros.%d", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.CurHeros[nIndex]}})
	} else if posType == POSTYPE_BACK {
		FieldName := fmt.Sprintf("backheros.%d", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.BackHeros[nIndex]}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("herobag.heros.%d", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.BagMoudle.HeroBag.Heros[nIndex]}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroLevelExp(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldExp := fmt.Sprintf("curheros.%d.curexp", nIndex)
		FieldLevel := fmt.Sprintf("curheros.%d.level", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid},
			&bson.M{"$set": bson.M{FieldExp: self.HeroMoudle.CurHeros[nIndex].CurExp,
				FieldLevel: self.HeroMoudle.CurHeros[nIndex].Level}})
	} else if posType == POSTYPE_BACK {
		FieldExp := fmt.Sprintf("backheros.%d.curexp", nIndex)
		FieldLevel := fmt.Sprintf("backheros.%d.level", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid},
			&bson.M{"$set": bson.M{FieldExp: self.HeroMoudle.BackHeros[nIndex].CurExp,
				FieldLevel: self.HeroMoudle.BackHeros[nIndex].Level}})
	} else if posType == POSTYPE_BAG {
		FieldExp := fmt.Sprintf("herobag.heros.%d.curexp", nIndex)
		FieldLevel := fmt.Sprintf("herobag.heros.%d.level", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid},
			&bson.M{"$set": bson.M{FieldExp: self.BagMoudle.HeroBag.Heros[nIndex].CurExp,
				FieldLevel: self.BagMoudle.HeroBag.Heros[nIndex].Level}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroBreakLevel(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curheros.%d.breaklevel", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.CurHeros[nIndex].BreakLevel}})
	} else if posType == POSTYPE_BACK {
		FieldName := fmt.Sprintf("backheros.%d.breaklevel", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.BackHeros[nIndex].BreakLevel}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("herobag.heros.%d.breaklevel", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.BagMoudle.HeroBag.Heros[nIndex].BreakLevel}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroGodLevel(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		godlvl := fmt.Sprintf("curheros.%d.godlevel", nIndex)
		quality := fmt.Sprintf("curheros.%d.quality", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{godlvl: self.HeroMoudle.CurHeros[nIndex].GodLevel,
			quality: self.HeroMoudle.CurHeros[nIndex].Quality}})
	} else if posType == POSTYPE_BACK {
		godlvl := fmt.Sprintf("backheros.%d.godlevel", nIndex)
		quality := fmt.Sprintf("backheros.%d.quality", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{godlvl: self.HeroMoudle.BackHeros[nIndex].GodLevel,
			quality: self.HeroMoudle.BackHeros[nIndex].Quality}})
	} else if posType == POSTYPE_BAG {
		godlvl := fmt.Sprintf("herobag.heros.%d.godlevel", nIndex)
		quality := fmt.Sprintf("herobag.heros.%d.quality", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{godlvl: self.BagMoudle.HeroBag.Heros[nIndex].GodLevel,
			quality: self.BagMoudle.HeroBag.Heros[nIndex].Quality}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroWakeLevel(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		WakelevelName := fmt.Sprintf("curheros.%d.wakelevel", nIndex)
		WakeItems := fmt.Sprintf("curheros.%d.wakeitem", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{WakelevelName: self.HeroMoudle.CurHeros[nIndex].WakeLevel,
			WakeItems: self.HeroMoudle.CurHeros[nIndex].WakeItem}})
	} else if posType == POSTYPE_BACK {
		WakelevelName := fmt.Sprintf("backheros.%d.wakelevel", nIndex)
		WakeItems := fmt.Sprintf("backheros.%d.wakeitem", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{WakelevelName: self.HeroMoudle.BackHeros[nIndex].WakeLevel,
			WakeItems: self.HeroMoudle.BackHeros[nIndex].WakeItem}})
	} else if posType == POSTYPE_BAG {
		WakelevelName := fmt.Sprintf("herobag.heros.%d.wakelevel", nIndex)
		WakeItems := fmt.Sprintf("herobag.heros.%d.wakeitem", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{WakelevelName: self.BagMoudle.HeroBag.Heros[nIndex].WakeLevel,
			WakeItems: self.BagMoudle.HeroBag.Heros[nIndex].WakeItem}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroWakeItem(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curheros.%d.wakeitem", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.CurHeros[nIndex].WakeItem}})
	} else if posType == POSTYPE_BACK {
		FieldName := fmt.Sprintf("backheros.%d.wakeitem", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.BackHeros[nIndex].WakeItem}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("herobag.heros.%d.wakeitem", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.BagMoudle.HeroBag.Heros[nIndex].WakeItem}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroCulture(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		cul := fmt.Sprintf("curheros.%d.cultures", nIndex)
		culcost := fmt.Sprintf("curheros.%d.culturescost", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{cul: self.HeroMoudle.CurHeros[nIndex].Cultures,
			culcost: self.HeroMoudle.CurHeros[nIndex].CulturesCost}})
	} else if posType == POSTYPE_BACK {
		cul := fmt.Sprintf("backheros.%d.cultures", nIndex)
		culcost := fmt.Sprintf("backheros.%d.culturescost", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{cul: self.HeroMoudle.BackHeros[nIndex].Cultures,
			culcost: self.HeroMoudle.BackHeros[nIndex].CulturesCost}})
	} else if posType == POSTYPE_BAG {
		cul := fmt.Sprintf("herobag.heros.%d.cultures", nIndex)
		culcost := fmt.Sprintf("herobag.heros.%d.culturescost", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{cul: self.BagMoudle.HeroBag.Heros[nIndex].Cultures,
			culcost: self.BagMoudle.HeroBag.Heros[nIndex].CulturesCost}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroDestiny(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldLevelName := fmt.Sprintf("curheros.%d.destinystate", nIndex)
		FieldExpName := fmt.Sprintf("curheros.%d.destinytime", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldLevelName: self.HeroMoudle.CurHeros[nIndex].DestinyState,
			FieldExpName: self.HeroMoudle.CurHeros[nIndex].DestinyTime}})
	} else if posType == POSTYPE_BACK {
		FieldLevelName := fmt.Sprintf("backheros.%d.destinystate", nIndex)
		FieldExpName := fmt.Sprintf("backheros.%d.destinytime", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldLevelName: self.HeroMoudle.BackHeros[nIndex].DestinyState,
			FieldExpName: self.HeroMoudle.CurHeros[nIndex].DestinyTime}})
	} else if posType == POSTYPE_BAG {
		FieldLevelName := fmt.Sprintf("herobag.heros.%d.destinystate", nIndex)
		FieldExpName := fmt.Sprintf("herobag.heros.%d.destinytime", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldLevelName: self.BagMoudle.HeroBag.Heros[nIndex].DestinyState,
			FieldExpName: self.HeroMoudle.CurHeros[nIndex].DestinyTime}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroXiLian(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		diaoquality := fmt.Sprintf("curheros.%d.diaowenquality", nIndex)
		ptys := fmt.Sprintf("curheros.%d.diaowenptys", nIndex)
		backs := fmt.Sprintf("curheros.%d.diaowenback", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{
			ptys:        self.HeroMoudle.CurHeros[nIndex].DiaoWenPtys,
			diaoquality: self.HeroMoudle.CurHeros[nIndex].DiaoWenQuality,
			backs:       self.HeroMoudle.CurHeros[nIndex].DiaoWenBack}})
	} else if posType == POSTYPE_BACK {
		diaoquality := fmt.Sprintf("backheros.%d.diaowenquality", nIndex)
		ptys := fmt.Sprintf("backheros.%d.diaowenptys", nIndex)
		backs := fmt.Sprintf("backheros.%d.diaowenback", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{
			ptys:        self.HeroMoudle.BackHeros[nIndex].DiaoWenPtys,
			diaoquality: self.HeroMoudle.BackHeros[nIndex].DiaoWenQuality,
			backs:       self.HeroMoudle.BackHeros[nIndex].DiaoWenBack}})
	} else if posType == POSTYPE_BAG {
		diaoquality := fmt.Sprintf("herobag.heros.%d.diaowenquality", nIndex)
		ptys := fmt.Sprintf("herobag.heros.%d.diaowenptys", nIndex)
		backs := fmt.Sprintf("herobag.heros.%d.diaowenback", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{
			ptys:        self.BagMoudle.HeroBag.Heros[nIndex].DiaoWenPtys,
			diaoquality: self.BagMoudle.HeroBag.Heros[nIndex].DiaoWenQuality,
			backs:       self.BagMoudle.HeroBag.Heros[nIndex].DiaoWenBack}})
	}
	return true
}

func (self *TPlayer) DB_SaveHeroQuality(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curheros.%d.quality", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.CurHeros[nIndex].Quality}})
	} else if posType == POSTYPE_BACK {
		FieldName := fmt.Sprintf("backheros.%d.quality", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.BackHeros[nIndex].Quality}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("herobag.heros.%d.quality", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.BagMoudle.HeroBag.Heros[nIndex].Quality}})
	}
	return true
}

//保存装备数据
func (self *TPlayer) DB_SaveEquipAt(posType int, nIndex int) {
	if posType == POSTYPE_BATTLE {
		self.HeroMoudle.DB_SaveBattleEquipAt(nIndex)
	} else if posType == POSTYPE_BAG {
		self.BagMoudle.DB_SaveBagEquipAt(nIndex)
	}
	return
}

//保存装备的强化等级
func (self *TPlayer) DB_SaveEquipStrength(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curequips.%d.strenglevel", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.HeroMoudle.CurEquips[nIndex].StrengLevel}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("equipbag.equips.%d.strenglevel", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{FieldName: self.BagMoudle.EquipBag.Equips[nIndex].StrengLevel}})
	}
	return true
}

func (self *TPlayer) DB_SaveEquipRefine(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		refinelvl := fmt.Sprintf("curequips.%d.refinelevel", nIndex)
		refineexp := fmt.Sprintf("curequips.%d.refineexp", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{refinelvl: self.HeroMoudle.CurEquips[nIndex].RefineLevel,
			refineexp: self.HeroMoudle.CurEquips[nIndex].RefineExp}})
	} else if posType == POSTYPE_BAG {
		refinelvl := fmt.Sprintf("equipbag.equips.%d.refinelevel", nIndex)
		refineexp := fmt.Sprintf("equipbag.equips.%d.refineexp", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{refinelvl: self.BagMoudle.EquipBag.Equips[nIndex].RefineLevel,
			refineexp: self.BagMoudle.EquipBag.Equips[nIndex].RefineExp}})
	}
	return true
}

func (self *TPlayer) DB_SaveEquipStar(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		star := fmt.Sprintf("curequips.%d.star", nIndex)
		starexp := fmt.Sprintf("curequips.%d.starexp", nIndex)
		starluck := fmt.Sprintf("curequips.%d.starluck", nIndex)
		starcost := fmt.Sprintf("curequips.%d.starcost", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{star: self.HeroMoudle.CurEquips[nIndex].Star,
			starexp:  self.HeroMoudle.CurEquips[nIndex].StarExp,
			starluck: self.HeroMoudle.CurEquips[nIndex].StarLuck,
			starcost: self.HeroMoudle.CurEquips[nIndex].StarCost}})
	} else if posType == POSTYPE_BAG {
		star := fmt.Sprintf("equipbag.equips.%d.star", nIndex)
		starexp := fmt.Sprintf("equipbag.equips.%d.starexp", nIndex)
		starluck := fmt.Sprintf("equipbag.equips.%d.starluck", nIndex)
		starcost := fmt.Sprintf("equipbag.equips.%d.starcost", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{star: self.BagMoudle.EquipBag.Equips[nIndex].Star,
			starexp:  self.BagMoudle.EquipBag.Equips[nIndex].StarExp,
			starluck: self.BagMoudle.EquipBag.Equips[nIndex].StarLuck,
			starcost: self.BagMoudle.EquipBag.Equips[nIndex].StarCost}})
	}
	return true
}

//保存宝物数据
func (self *TPlayer) DB_SaveGemAt(posType int, nIndex int) {
	if posType == POSTYPE_BATTLE {
		self.HeroMoudle.DB_SaveBattleGemAt(nIndex)
	} else if posType == POSTYPE_BAG {
		self.BagMoudle.DB_SaveBagGemAt(nIndex)
	}
	return
}

//保存宠物数据
func (self *TPlayer) DB_SavePetAt(posType int, nIndex int) {
	if posType == POSTYPE_BATTLE {
		self.HeroMoudle.DB_SaveBattlePetAt(nIndex)
	} else if posType == POSTYPE_BAG {
		self.BagMoudle.DB_SaveBagPetAt(nIndex)
	}
	return
}

func (self *TPlayer) DB_SavePetLevel(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		levelName := fmt.Sprintf("curpets.%d.level", nIndex)
		ExpName := fmt.Sprintf("curpets.%d.exp", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{levelName: self.HeroMoudle.CurPets[nIndex].Level,
			ExpName: self.HeroMoudle.CurPets[nIndex].Exp}})
	} else if posType == POSTYPE_BAG {
		levelName := fmt.Sprintf("petbag.pets.%d.level", nIndex)
		ExpName := fmt.Sprintf("petbag.pets.%d.exp", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{levelName: self.BagMoudle.PetBag.Pets[nIndex].Level,
			ExpName: self.BagMoudle.PetBag.Pets[nIndex].Exp}})
	}
	return true
}

func (self *TPlayer) DB_SavePetStar(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		starName := fmt.Sprintf("curpets.%d.star", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{starName: self.HeroMoudle.CurPets[nIndex].Star}})
	} else if posType == POSTYPE_BAG {
		starName := fmt.Sprintf("petbag.pets.%d.star", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{starName: self.BagMoudle.PetBag.Pets[nIndex].Star}})
	}
	return true
}

func (self *TPlayer) DB_SavePetGod(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		levelName := fmt.Sprintf("curpets.%d.god", nIndex)
		ExpName := fmt.Sprintf("curpets.%d.godexp", nIndex)
		mongodb.UpdateToDB("PlayerHero", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{levelName: self.HeroMoudle.CurPets[nIndex].God,
			ExpName: self.HeroMoudle.CurPets[nIndex].GodExp}})
	} else if posType == POSTYPE_BAG {
		levelName := fmt.Sprintf("petbag.pets.%d.god", nIndex)
		ExpName := fmt.Sprintf("petbag.pets.%d.godexp", nIndex)
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.playerid}, &bson.M{"$set": bson.M{levelName: self.BagMoudle.PetBag.Pets[nIndex].God,
			ExpName: self.BagMoudle.PetBag.Pets[nIndex].GodExp}})
	}
	return true
}

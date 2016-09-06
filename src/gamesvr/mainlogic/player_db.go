package mainlogic

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func (player *TPlayer) DB_SaveHeroAt(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curheros.%d", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.CurHeros[nIndex]}})
	} else if posType == POSTYPE_BACK {
		FieldName := fmt.Sprintf("backheros.%d", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.BackHeros[nIndex]}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("herobag.heros.%d", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.BagMoudle.HeroBag.Heros[nIndex]}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroLevelExp(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldExp := fmt.Sprintf("curheros.%d.curexp", nIndex)
		FieldLevel := fmt.Sprintf("curheros.%d.level", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid},
			&bson.M{"$set": bson.M{FieldExp: player.HeroMoudle.CurHeros[nIndex].CurExp,
				FieldLevel: player.HeroMoudle.CurHeros[nIndex].Level}})
	} else if posType == POSTYPE_BACK {
		FieldExp := fmt.Sprintf("backheros.%d.curexp", nIndex)
		FieldLevel := fmt.Sprintf("backheros.%d.level", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid},
			&bson.M{"$set": bson.M{FieldExp: player.HeroMoudle.BackHeros[nIndex].CurExp,
				FieldLevel: player.HeroMoudle.BackHeros[nIndex].Level}})
	} else if posType == POSTYPE_BAG {
		FieldExp := fmt.Sprintf("herobag.heros.%d.curexp", nIndex)
		FieldLevel := fmt.Sprintf("herobag.heros.%d.level", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid},
			&bson.M{"$set": bson.M{FieldExp: player.BagMoudle.HeroBag.Heros[nIndex].CurExp,
				FieldLevel: player.BagMoudle.HeroBag.Heros[nIndex].Level}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroBreakLevel(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curheros.%d.breaklevel", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.CurHeros[nIndex].BreakLevel}})
	} else if posType == POSTYPE_BACK {
		FieldName := fmt.Sprintf("backheros.%d.breaklevel", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.BackHeros[nIndex].BreakLevel}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("herobag.heros.%d.breaklevel", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.BagMoudle.HeroBag.Heros[nIndex].BreakLevel}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroGodLevel(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		godlvl := fmt.Sprintf("curheros.%d.godlevel", nIndex)
		quality := fmt.Sprintf("curheros.%d.quality", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{godlvl: player.HeroMoudle.CurHeros[nIndex].GodLevel,
			quality: player.HeroMoudle.CurHeros[nIndex].Quality}})
	} else if posType == POSTYPE_BACK {
		godlvl := fmt.Sprintf("backheros.%d.godlevel", nIndex)
		quality := fmt.Sprintf("backheros.%d.quality", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{godlvl: player.HeroMoudle.BackHeros[nIndex].GodLevel,
			quality: player.HeroMoudle.BackHeros[nIndex].Quality}})
	} else if posType == POSTYPE_BAG {
		godlvl := fmt.Sprintf("herobag.heros.%d.godlevel", nIndex)
		quality := fmt.Sprintf("herobag.heros.%d.quality", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{godlvl: player.BagMoudle.HeroBag.Heros[nIndex].GodLevel,
			quality: player.BagMoudle.HeroBag.Heros[nIndex].Quality}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroWakeLevel(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		WakelevelName := fmt.Sprintf("curheros.%d.wakelevel", nIndex)
		WakeItems := fmt.Sprintf("curheros.%d.wakeitem", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{WakelevelName: player.HeroMoudle.CurHeros[nIndex].WakeLevel,
			WakeItems: player.HeroMoudle.CurHeros[nIndex].WakeItem}})
	} else if posType == POSTYPE_BACK {
		WakelevelName := fmt.Sprintf("backheros.%d.wakelevel", nIndex)
		WakeItems := fmt.Sprintf("backheros.%d.wakeitem", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{WakelevelName: player.HeroMoudle.BackHeros[nIndex].WakeLevel,
			WakeItems: player.HeroMoudle.BackHeros[nIndex].WakeItem}})
	} else if posType == POSTYPE_BAG {
		WakelevelName := fmt.Sprintf("herobag.heros.%d.wakelevel", nIndex)
		WakeItems := fmt.Sprintf("herobag.heros.%d.wakeitem", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{WakelevelName: player.BagMoudle.HeroBag.Heros[nIndex].WakeLevel,
			WakeItems: player.BagMoudle.HeroBag.Heros[nIndex].WakeItem}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroWakeItem(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curheros.%d.wakeitem", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.CurHeros[nIndex].WakeItem}})
	} else if posType == POSTYPE_BACK {
		FieldName := fmt.Sprintf("backheros.%d.wakeitem", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.BackHeros[nIndex].WakeItem}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("herobag.heros.%d.wakeitem", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.BagMoudle.HeroBag.Heros[nIndex].WakeItem}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroCulture(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		cul := fmt.Sprintf("curheros.%d.cultures", nIndex)
		culcost := fmt.Sprintf("curheros.%d.culturescost", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{cul: player.HeroMoudle.CurHeros[nIndex].Cultures,
			culcost: player.HeroMoudle.CurHeros[nIndex].CulturesCost}})
	} else if posType == POSTYPE_BACK {
		cul := fmt.Sprintf("backheros.%d.cultures", nIndex)
		culcost := fmt.Sprintf("backheros.%d.culturescost", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{cul: player.HeroMoudle.BackHeros[nIndex].Cultures,
			culcost: player.HeroMoudle.BackHeros[nIndex].CulturesCost}})
	} else if posType == POSTYPE_BAG {
		cul := fmt.Sprintf("herobag.heros.%d.cultures", nIndex)
		culcost := fmt.Sprintf("herobag.heros.%d.culturescost", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{cul: player.BagMoudle.HeroBag.Heros[nIndex].Cultures,
			culcost: player.BagMoudle.HeroBag.Heros[nIndex].CulturesCost}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroDestiny(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldLevelName := fmt.Sprintf("curheros.%d.destinystate", nIndex)
		FieldExpName := fmt.Sprintf("curheros.%d.destinytime", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldLevelName: player.HeroMoudle.CurHeros[nIndex].DestinyState,
			FieldExpName: player.HeroMoudle.CurHeros[nIndex].DestinyTime}})
	} else if posType == POSTYPE_BACK {
		FieldLevelName := fmt.Sprintf("backheros.%d.destinystate", nIndex)
		FieldExpName := fmt.Sprintf("backheros.%d.destinytime", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldLevelName: player.HeroMoudle.BackHeros[nIndex].DestinyState,
			FieldExpName: player.HeroMoudle.CurHeros[nIndex].DestinyTime}})
	} else if posType == POSTYPE_BAG {
		FieldLevelName := fmt.Sprintf("herobag.heros.%d.destinystate", nIndex)
		FieldExpName := fmt.Sprintf("herobag.heros.%d.destinytime", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldLevelName: player.BagMoudle.HeroBag.Heros[nIndex].DestinyState,
			FieldExpName: player.HeroMoudle.CurHeros[nIndex].DestinyTime}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroXiLian(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		diaoquality := fmt.Sprintf("curheros.%d.diaowenquality", nIndex)
		ptys := fmt.Sprintf("curheros.%d.diaowenptys", nIndex)
		backs := fmt.Sprintf("curheros.%d.diaowenback", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{
			ptys:        player.HeroMoudle.CurHeros[nIndex].DiaoWenPtys,
			diaoquality: player.HeroMoudle.CurHeros[nIndex].DiaoWenQuality,
			backs:       player.HeroMoudle.CurHeros[nIndex].DiaoWenBack}})
	} else if posType == POSTYPE_BACK {
		diaoquality := fmt.Sprintf("backheros.%d.diaowenquality", nIndex)
		ptys := fmt.Sprintf("backheros.%d.diaowenptys", nIndex)
		backs := fmt.Sprintf("backheros.%d.diaowenback", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{
			ptys:        player.HeroMoudle.BackHeros[nIndex].DiaoWenPtys,
			diaoquality: player.HeroMoudle.BackHeros[nIndex].DiaoWenQuality,
			backs:       player.HeroMoudle.BackHeros[nIndex].DiaoWenBack}})
	} else if posType == POSTYPE_BAG {
		diaoquality := fmt.Sprintf("herobag.heros.%d.diaowenquality", nIndex)
		ptys := fmt.Sprintf("herobag.heros.%d.diaowenptys", nIndex)
		backs := fmt.Sprintf("herobag.heros.%d.diaowenback", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{
			ptys:        player.BagMoudle.HeroBag.Heros[nIndex].DiaoWenPtys,
			diaoquality: player.BagMoudle.HeroBag.Heros[nIndex].DiaoWenQuality,
			backs:       player.BagMoudle.HeroBag.Heros[nIndex].DiaoWenBack}})
	}
	return true
}

func (player *TPlayer) DB_SaveHeroQuality(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curheros.%d.quality", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.CurHeros[nIndex].Quality}})
	} else if posType == POSTYPE_BACK {
		FieldName := fmt.Sprintf("backheros.%d.quality", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.BackHeros[nIndex].Quality}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("herobag.heros.%d.quality", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.BagMoudle.HeroBag.Heros[nIndex].Quality}})
	}
	return true
}

//保存装备数据
func (player *TPlayer) DB_SaveEquipAt(posType int, nIndex int) {
	if posType == POSTYPE_BATTLE {
		player.HeroMoudle.DB_SaveBattleEquipAt(nIndex)
	} else if posType == POSTYPE_BAG {
		player.BagMoudle.DB_SaveBagEquipAt(nIndex)
	}
	return
}

//保存装备的强化等级
func (player *TPlayer) DB_SaveEquipStrength(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		FieldName := fmt.Sprintf("curequips.%d.strenglevel", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.HeroMoudle.CurEquips[nIndex].StrengLevel}})
	} else if posType == POSTYPE_BAG {
		FieldName := fmt.Sprintf("equipbag.equips.%d.strenglevel", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{FieldName: player.BagMoudle.EquipBag.Equips[nIndex].StrengLevel}})
	}
	return true
}

func (player *TPlayer) DB_SaveEquipRefine(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		refinelvl := fmt.Sprintf("curequips.%d.refinelevel", nIndex)
		refineexp := fmt.Sprintf("curequips.%d.refineexp", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{refinelvl: player.HeroMoudle.CurEquips[nIndex].RefineLevel,
			refineexp: player.HeroMoudle.CurEquips[nIndex].RefineExp}})
	} else if posType == POSTYPE_BAG {
		refinelvl := fmt.Sprintf("equipbag.equips.%d.refinelevel", nIndex)
		refineexp := fmt.Sprintf("equipbag.equips.%d.refineexp", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{refinelvl: player.BagMoudle.EquipBag.Equips[nIndex].RefineLevel,
			refineexp: player.BagMoudle.EquipBag.Equips[nIndex].RefineExp}})
	}
	return true
}

func (player *TPlayer) DB_SaveEquipStar(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		star := fmt.Sprintf("curequips.%d.star", nIndex)
		starexp := fmt.Sprintf("curequips.%d.starexp", nIndex)
		starluck := fmt.Sprintf("curequips.%d.starluck", nIndex)
		starcost := fmt.Sprintf("curequips.%d.starcost", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{star: player.HeroMoudle.CurEquips[nIndex].Star,
			starexp:  player.HeroMoudle.CurEquips[nIndex].StarExp,
			starluck: player.HeroMoudle.CurEquips[nIndex].StarLuck,
			starcost: player.HeroMoudle.CurEquips[nIndex].StarCost}})
	} else if posType == POSTYPE_BAG {
		star := fmt.Sprintf("equipbag.equips.%d.star", nIndex)
		starexp := fmt.Sprintf("equipbag.equips.%d.starexp", nIndex)
		starluck := fmt.Sprintf("equipbag.equips.%d.starluck", nIndex)
		starcost := fmt.Sprintf("equipbag.equips.%d.starcost", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{star: player.BagMoudle.EquipBag.Equips[nIndex].Star,
			starexp:  player.BagMoudle.EquipBag.Equips[nIndex].StarExp,
			starluck: player.BagMoudle.EquipBag.Equips[nIndex].StarLuck,
			starcost: player.BagMoudle.EquipBag.Equips[nIndex].StarCost}})
	}
	return true
}

//保存宝物数据
func (player *TPlayer) DB_SaveGemAt(posType int, nIndex int) {
	if posType == POSTYPE_BATTLE {
		player.HeroMoudle.DB_SaveBattleGemAt(nIndex)
	} else if posType == POSTYPE_BAG {
		player.BagMoudle.DB_SaveBagGemAt(nIndex)
	}
	return
}

//保存宠物数据
func (player *TPlayer) DB_SavePetAt(posType int, nIndex int) {
	if posType == POSTYPE_BATTLE {
		player.HeroMoudle.DB_SaveBattlePetAt(nIndex)
	} else if posType == POSTYPE_BAG {
		player.BagMoudle.DB_SaveBagPetAt(nIndex)
	}
	return
}

func (player *TPlayer) DB_SavePetLevel(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		levelName := fmt.Sprintf("curpets.%d.level", nIndex)
		ExpName := fmt.Sprintf("curpets.%d.exp", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{levelName: player.HeroMoudle.CurPets[nIndex].Level,
			ExpName: player.HeroMoudle.CurPets[nIndex].Exp}})
	} else if posType == POSTYPE_BAG {
		levelName := fmt.Sprintf("petbag.pets.%d.level", nIndex)
		ExpName := fmt.Sprintf("petbag.pets.%d.exp", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{levelName: player.BagMoudle.PetBag.Pets[nIndex].Level,
			ExpName: player.BagMoudle.PetBag.Pets[nIndex].Exp}})
	}
	return true
}

func (player *TPlayer) DB_SavePetStar(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		starName := fmt.Sprintf("curpets.%d.star", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{starName: player.HeroMoudle.CurPets[nIndex].Star}})
	} else if posType == POSTYPE_BAG {
		starName := fmt.Sprintf("petbag.pets.%d.star", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{starName: player.BagMoudle.PetBag.Pets[nIndex].Star}})
	}
	return true
}

func (player *TPlayer) DB_SavePetGod(posType int, nIndex int) bool {
	if posType == POSTYPE_BATTLE {
		levelName := fmt.Sprintf("curpets.%d.god", nIndex)
		ExpName := fmt.Sprintf("curpets.%d.godexp", nIndex)
		GameSvrUpdateToDB("PlayerHero", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{levelName: player.HeroMoudle.CurPets[nIndex].God,
			ExpName: player.HeroMoudle.CurPets[nIndex].GodExp}})
	} else if posType == POSTYPE_BAG {
		levelName := fmt.Sprintf("petbag.pets.%d.god", nIndex)
		ExpName := fmt.Sprintf("petbag.pets.%d.godexp", nIndex)
		GameSvrUpdateToDB("PlayerBag", &bson.M{"_id": player.playerid}, &bson.M{"$set": bson.M{levelName: player.BagMoudle.PetBag.Pets[nIndex].God,
			ExpName: player.BagMoudle.PetBag.Pets[nIndex].GodExp}})
	}
	return true
}

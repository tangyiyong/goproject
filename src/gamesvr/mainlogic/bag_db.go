package mainlogic

import (
	"fmt"
	"gamelog"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//修改一个英雄ID
func (self *TBagMoudle) DB_UpdateHeroID(pos int, heroID int) {
	filedName := fmt.Sprintf("herobag.heros.%d.heroid", pos)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{filedName: heroID}})
}

//添加一个英雄列表
func (self *TBagMoudle) DB_AddHeroList(heros []THeroData, bCol bool) {
	count := len(heros)
	if count <= 0 {
		gamelog.Error("DB_AddHeroList Error :Invalid count :%d", count)
		return
	}

	if bCol == false {
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"herobag.heros": heros}})
	} else {
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"herobag.heros": heros,
			"colheros": []int{heros[0].ID}}})
	}
}

func (self *TBagMoudle) DB_RemoveHeroAt(nIndex int) bool {
	FieldName := fmt.Sprintf("herobag.heros.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$unset": bson.M{FieldName: 1}})
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"herobag.heros": nil}})
	return true
}

func (self *TBagMoudle) DB_RemoveHeros(nIndex []int) bool {
	var heros bson.M = make(map[string]interface{}, 1)
	for _, v := range nIndex {
		FieldName := fmt.Sprintf("herobag.heros.%d", v)
		heros[FieldName] = 1
	}

	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$unset": heros})
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"herobag.heros": nil}})
	return true
}

//装备包
func (self *TBagMoudle) DB_SaveBagEquipAt(nIndex int) {
	FieldName := fmt.Sprintf("equipbag.equips.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.EquipBag.Equips[nIndex]}})
}

func (self *TBagMoudle) DB_RemoveEquipAt(nIndex int) bool {
	FieldName := fmt.Sprintf("equipbag.equips.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$unset": bson.M{FieldName: 1}})
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"equipbag.equips": nil}})
	return true
}

func (self *TBagMoudle) DB_RemoveEquips(nIndex []int) bool {
	var equips bson.M = make(map[string]interface{}, 1)
	for _, v := range nIndex {
		FieldName := fmt.Sprintf("equipbag.equips.%d", v)
		equips[FieldName] = 1
	}

	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$unset": equips})
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"equipbag.equips": nil}})
	return true
}

//添加一个装备列表
func (self *TBagMoudle) DB_AddEquipsList(equips []TEquipData) {
	count := len(equips)
	if count <= 0 {
		gamelog.Error("DB_AddEquipList Error :Invalid count :%d", count)
		return
	}
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"equipbag.equips": equips}})
}

//宝物背包
func (self *TBagMoudle) DB_SaveBagGemAt(nIndex int) {
	FieldName := fmt.Sprintf("gembag.gems.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.GemBag.Gems[nIndex]}})
}

//添加一个宝物列表
func (self *TBagMoudle) DB_AddGemList(gems []TGemData) {
	count := len(gems)
	if count <= 0 {
		gamelog.Error("DB_AddGemList Error :Invalid count :%d", count)
		return
	}
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"gembag.gems": gems}})
}

func (self *TBagMoudle) DB_RemoveGemAt(nIndex int) bool {
	FieldName := fmt.Sprintf("gembag.gems.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$unset": bson.M{FieldName: 1}})
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"gembag.gems": nil}})
	return true
}

func (self *TBagMoudle) DB_RemoveGems(nIndex []int) bool {
	var gems bson.M = make(map[string]interface{}, 1)
	for _, v := range nIndex {
		FieldName := fmt.Sprintf("gembag.gems.%d", v)
		gems[FieldName] = 1
	}

	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$unset": gems})
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"gembag.gems": nil}})
	return true
}

//英雄碎片包
func (self *TBagMoudle) DB_SaveHeroPieceBagAt(nIndex int) {
	FieldName := fmt.Sprintf("heropiecebag.items.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.HeroPieceBag.Items[nIndex]}})
}

//装备碎片包
func (self *TBagMoudle) DB_SaveEquipPieceBagAt(nIndex int) {
	FieldName := fmt.Sprintf("equippiecebag.items.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.EquipPieceBag.Items[nIndex]}})
}

//宝物碎片包
func (self *TBagMoudle) DB_SaveGemPieceBagAt(nIndex int) {
	FieldName := fmt.Sprintf("gempiecebag.items.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.GemPieceBag.Items[nIndex]}})
}

//英雄碎片包
func (self *TBagMoudle) DB_RemoveHeroPiece(itemid int) {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"heropiecebag.items": bson.M{"itemid": itemid}}})
}

//装备碎片包
func (self *TBagMoudle) DB_RemoveEquipPiece(itemid int) {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"equippiecebag.items": bson.M{"itemid": itemid}}})
}

//宝物碎片包
func (self *TBagMoudle) DB_RemoveGemPiece(itemid int) {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"gempiecebag.items": bson.M{"itemid": itemid}}})
}

//道具背包
func (self *TBagMoudle) DB_SaveNormalItemBagAt(nIndex int) {
	FieldName := fmt.Sprintf("normalitembag.items.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.NormalItemBag.Items[nIndex]}})
}

//删除道具背包
func (self *TBagMoudle) DB_RemoveNormalItem(itemid int) {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"normalitembag.items": bson.M{"itemid": itemid}}})
}

//道具背包
func (self *TBagMoudle) DB_SaveWakeItemBagAt(nIndex int) {
	FieldName := fmt.Sprintf("wakeitembag.items.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.WakeItemBag.Items[nIndex]}})
}

//删除道具背包
func (self *TBagMoudle) DB_RemoveWakeItem(itemid int) {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"wakeitembag.items": bson.M{"itemid": itemid}}})
}

//宠物包
func (self *TBagMoudle) DB_SaveBagPetAt(nIndex int) {
	FieldName := fmt.Sprintf("petbag.pets.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.PetBag.Pets[nIndex]}})
}

func (self *TBagMoudle) DB_RemovePetAt(nIndex int) bool {
	FieldName := fmt.Sprintf("petbag.pets.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$unset": bson.M{FieldName: 1}})
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"petbag.pets": nil}})
	return true
}

//添加一个宠物列表
func (self *TBagMoudle) DB_AddPetList(pets []TPetData, bCol bool) {
	count := len(pets)
	if count <= 0 {
		gamelog.Error("DB_AddPetList Error :Invalid count :%d", count)
		return
	}

	if bCol == false {
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"petbag.pets": pets}})
	} else {
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"petbag.pets": pets,
			"colpets": []int{pets[0].ID}}})
	}
}

//保存宠物碎片包
func (self *TBagMoudle) DB_SavePetPieceBagAt(nIndex int) {
	FieldName := fmt.Sprintf("petpiecebag.items.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.PetPieceBag.Items[nIndex]}})
}

//宠物碎片包
func (self *TBagMoudle) DB_RemovePetPiece(itemid int) {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"petpiecebag.items": bson.M{"itemid": itemid}}})
}

//宠物背包
func (self *TBagMoudle) DB_SavePetBag() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"petbag.pets": self.PetBag.Pets}})
}

//保存将灵片包
func (self *TBagMoudle) DB_SaveHeroSoulBagAt(nIndex int) {
	FieldName := fmt.Sprintf("herosoulbag.items.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.HeroSoulBag.Items[nIndex]}})
}

//将灵包
func (self *TBagMoudle) DB_RemoveHeroSoul(itemid int) {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"herosoulbag.items": bson.M{"itemid": itemid}}})
}

//数据库中在英雄背包中最末尾添加一个英雄
func (self *TBagMoudle) DB_AddFashionAtLast() {
	nIndex := len(self.FashionBag.Fashions) - 1
	if nIndex < 0 {
		gamelog.Error("DB_AddFashionAtLast Error :Invalid nIndex :%d", nIndex)
		return
	}

	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"fashionbag.fashions": self.FashionBag.Fashions[nIndex]}})
}

//添加一个英雄列表
func (self *TBagMoudle) DB_AddFashionList(fashions []TFashionData) {
	count := len(fashions)
	if count <= 0 {
		gamelog.Error("DB_AddFashionList Error :Invalid count :%d", count)
		return
	}
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"fashionbag.fashions": fashions}})
}

func (self *TBagMoudle) DB_SaveFashionAt(nIndex int) {
	FieldName := fmt.Sprintf("fashionbag.fashions.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.FashionBag.Fashions[nIndex]}})
}

//英雄碎片包
func (self *TBagMoudle) DB_SaveFashionPieceBagAt(nIndex int) {
	FieldName := fmt.Sprintf("fashionpiecebag.items.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.FashionPieceBag.Items[nIndex]}})
}

//英雄碎片包
func (self *TBagMoudle) DB_RemoveFashionPiece(itemid int) {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"fashionpiecebag.items": bson.M{"itemid": itemid}}})
}

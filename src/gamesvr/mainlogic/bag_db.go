package mainlogic

import (
	"fmt"
	"gamelog"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (self *TBagMoudle) DB_SaveHeroBag() bool {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"herobag.heros": self.HeroBag.Heros}})
	return true
}

//数据库中在英雄背包中最末尾添加一个英雄
func (self *TBagMoudle) DB_AddHeroAtLast(bCol bool) {
	nIndex := len(self.HeroBag.Heros) - 1
	if nIndex < 0 {
		gamelog.Error("DB_AddHeroAtLast Error :Invalid nIndex :%d", nIndex)
		return
	}

	if bCol == false {
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"herobag.heros": self.HeroBag.Heros[nIndex]}})
	} else {
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"herobag.heros": self.HeroBag.Heros[nIndex],
			"colheros": self.HeroBag.Heros[nIndex].ID}})
	}

}

//修改一个英雄ID
func (self *TBagMoudle) DB_UpdateHeroID(pos int, heroID int) {
	filedName := fmt.Sprintf("herobag.heros.%d.heroid", pos)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{filedName: heroID}})
}

//添加一个英雄列表
func (self *TBagMoudle) DB_AddHeroList(heros []THeroData) {
	count := len(heros)
	if count <= 0 {
		gamelog.Error("DB_AddHeroList Error :Invalid count :%d", count)
		return
	}
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"herobag.heros": heros}})
}

func (self *TBagMoudle) DB_RemoveHeroAt(nIndex int) bool {
	FieldName := fmt.Sprintf("herobag.heros.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$unset": bson.M{FieldName: 1}})
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"herobag.heros": nil}})
	return true
}

//装备包
func (self *TBagMoudle) DB_SaveBagEquips() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"equipbag.equips": self.EquipBag.Equips}})
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

//数据库中在装备背包中最末尾添加一个装备
func (self *TBagMoudle) DB_AddEquipAtLast() {
	nIndex := len(self.EquipBag.Equips) - 1
	if nIndex < 0 {
		gamelog.Error("DB_AddEquipoAtLast Error :Invalid nIndex :%d", nIndex)
		return
	}
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"equipbag.equips": self.EquipBag.Equips[nIndex]}})
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
func (self *TBagMoudle) DB_SaveGemBag() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"gembag.gems": self.GemBag.Gems}})
}

//宝物背包
func (self *TBagMoudle) DB_SaveBagGemAt(nIndex int) {
	FieldName := fmt.Sprintf("gembag.gems.%d", nIndex)
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{FieldName: self.GemBag.Gems[nIndex]}})
}

func (self *TBagMoudle) DB_AddGemAtLast() {
	nIndex := len(self.GemBag.Gems) - 1
	if nIndex < 0 {
		gamelog.Error("DB_AddGemAtLast Error :Invalid nIndex :%d", nIndex)
		return
	}
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"gembag.gems": self.GemBag.Gems[nIndex]}})
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

//宝物碎片包
func (self *TBagMoudle) DB_SaveGemPieceBag() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"gempiecebag.items": self.GemPieceBag.Items}})
}

//道具背包
func (self *TBagMoudle) DB_SaveNormalItemBag() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"normalitembag.items": self.NormalItemBag.Items}})
}

//觉醒道具背包
func (self *TBagMoudle) DB_SaveWakeItemBag() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"wakeitembag.items": self.WakeItemBag.Items}})
}

//英雄碎片包
func (self *TBagMoudle) DB_SaveHeroPieceBag() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"heropiecebag.items": self.HeroPieceBag.Items}})
}

//装备碎片包
func (self *TBagMoudle) DB_SaveEquipPieceBag() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"equipbag.equips": self.EquipPieceBag.Items}})
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

//在宠物背包中最末尾添加一个宠物
func (self *TBagMoudle) DB_AddPetAtLast(bCol bool) {
	nIndex := len(self.PetBag.Pets) - 1
	if nIndex < 0 {
		gamelog.Error("DB_AddPetAtLast Error :Invalid nIndex :%d", nIndex)
		return
	}

	if bCol == false {
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"petbag.pets": self.PetBag.Pets[nIndex]}})
	} else {
		mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"petbag.pets": self.PetBag.Pets[nIndex],
			"colpets": self.PetBag.Pets[nIndex].ID}})
	}

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
func (self *TBagMoudle) DB_AddPetList(pets []TPetData) {
	count := len(pets)
	if count <= 0 {
		gamelog.Error("DB_AddPetList Error :Invalid count :%d", count)
		return
	}
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{"petbag.pets": pets}})
}

//宠物碎片包
func (self *TBagMoudle) DB_SavePetPieceBag() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"petpiecebag.items": self.PetPieceBag.Items}})
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

//将灵背包
func (self *TBagMoudle) DB_SaveBagHeroSoul() {
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"herosoulbag.items": self.HeroSoulBag.Items}})
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
	mongodb.UpdateToDB("PlayerBag", &bson.M{"_id": self.PlayerID}, &bson.M{"$pushAll": bson.M{FieldName: self.FashionBag.Fashions[nIndex]}})
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

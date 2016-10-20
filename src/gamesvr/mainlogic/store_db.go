package mainlogic

import (
	"fmt"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 修改商品标记
func (storemodule *TStoreModule) DB_UpdateShopItemStatusToDatabase(index int, storeType int) {
	if storeType == gamedata.StoreType_Hero {
		filedName := fmt.Sprintf("heroshopitemlst.%d.status", index)
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			filedName: 1}})
	} else if storeType == gamedata.StoreType_Awake {
		filedName := fmt.Sprintf("awakeshopitemlst.%d.status", index)
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			filedName: 1}})
	} else if storeType == gamedata.StoreType_Pet {
		filedName := fmt.Sprintf("patshopitemlst.%d.status", index)
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			filedName: 1}})
	}

}

//! 存取刷新数据到数据库
func (storemodule *TStoreModule) DB_SaveRefreshInfoToDatabase(storeType int) {
	if storeType == gamedata.StoreType_Hero {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"herorefreshcount":     storemodule.HeroRefreshCount,
			"herofreerefreshcount": storemodule.HeroFreeRefreshCount,
			"herofreerefreshtime":  storemodule.HeroFreeRefreshTime,
			"heroshopitemlst":      storemodule.HeroShopItemLst}})
	} else if storeType == gamedata.StoreType_Awake {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"awakerefreshcount":     storemodule.AwakeRefreshCount,
			"awakefreerefreshcount": storemodule.AwakeFreeRefreshCount,
			"awakefreerefreshtime":  storemodule.AwakeFreeRefreshTime,
			"awakeshopitemlst":      storemodule.AwakeShopItemLst}})
	} else if storeType == gamedata.StoreType_Pet {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"petrefreshcount":     storemodule.PetRefreshCount,
			"petfreerefreshcount": storemodule.PetFreeRefreshCount,
			"petfreerefreshtime":  storemodule.PetFreeRefreshTime,
			"petshopitemlst":      storemodule.PetShopItemLst}})
	}

}

//! 更新商品列表
func (storemodule *TStoreModule) DB_UpdateShopItemToDatabase(storeType int) {
	if storeType == gamedata.StoreType_Hero {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"heroshopitemlst": storemodule.HeroShopItemLst}})
	} else if storeType == gamedata.StoreType_Awake {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"awakeshopitemlst": storemodule.AwakeShopItemLst}})
	} else if storeType == gamedata.StoreType_Pet {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"patshopitemlst": storemodule.PetShopItemLst}})
	}

}

//! 更新免费刷新时间
func (storemodule *TStoreModule) DB_UpdateFreeRefreshtime(storeType int) {
	if storeType == gamedata.StoreType_Hero {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"herofreerefreshtime": storemodule.HeroFreeRefreshTime}})
	} else if storeType == gamedata.StoreType_Awake {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"awakefreerefreshtime": storemodule.AwakeFreeRefreshTime}})
	} else if storeType == gamedata.StoreType_Pet {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"petfreerefreshtime": storemodule.PetFreeRefreshTime}})
	}

}

//! 更新免费次数
func (storemodule *TStoreModule) DB_UpdateRefreshFreeCount(storeType int) {
	if storeType == gamedata.StoreType_Hero {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"herofreerefreshtime":  storemodule.HeroFreeRefreshTime,
			"herofreerefreshcount": storemodule.HeroFreeRefreshCount}})
	} else if storeType == gamedata.StoreType_Awake {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"awakefreerefreshtime":  storemodule.AwakeFreeRefreshTime,
			"awakefreerefreshcount": storemodule.AwakeFreeRefreshCount}})
	} else if storeType == gamedata.StoreType_Pet {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"petfreerefreshtime":  storemodule.PetFreeRefreshTime,
			"petfreerefreshcount": storemodule.PetFreeRefreshCount}})
	}

}

//! 更新次数
func (storemodule *TStoreModule) DB_UpdateRefreshCount(storeType int) {
	if storeType == gamedata.StoreType_Hero {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"herorefreshcount": storemodule.HeroRefreshCount}})
	} else if storeType == gamedata.StoreType_Awake {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"awakerefreshcount": storemodule.AwakeRefreshCount}})
	} else if storeType == gamedata.StoreType_Pet {
		mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
			"petrefreshcount": storemodule.PetRefreshCount}})
	}
}

//! 更新重置时间
func (storemodule *TStoreModule) DB_UpdateResetTime() {
	mongodb.UpdateToDB("HeroStore", &bson.M{"_id": storemodule.PlayerID}, &bson.M{"$set": bson.M{
		"resetday": storemodule.ResetDay}})
}

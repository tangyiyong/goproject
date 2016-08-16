package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"utility"
)

//请求背包数据
func Hand_GetBagData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagData_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagHeros : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagData_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		comdata := utility.CompressData(b)
		//gamelog.Error("Hand_GetBagHeros : orginalLen:%d, compressLen:%d", len(b), len(comdata))
		w.Write(comdata)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.Heros = pPlayer.BagMoudle.HeroBag.Heros
	response.EquipPieces = pPlayer.BagMoudle.EquipPieceBag.Items
	response.Equips = pPlayer.BagMoudle.EquipBag.Equips
	response.GemPieces = pPlayer.BagMoudle.GemPieceBag.Items
	response.Gems = pPlayer.BagMoudle.GemBag.Gems
	response.HeroPieces = pPlayer.BagMoudle.HeroPieceBag.Items
	response.PetPieces = pPlayer.BagMoudle.PetPieceBag.Items
	response.Pets = pPlayer.BagMoudle.PetBag.Pets
	response.WakeItems = pPlayer.BagMoudle.WakeItemBag.Items
	response.Normals = pPlayer.BagMoudle.NormalItemBag.Items
	response.HeroSouls = pPlayer.BagMoudle.HeroSoulBag.Items

	response.RetCode = msg.RE_SUCCESS

	return
}

//请求背包中的所有英雄
func Hand_GetBagHeros(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagHeros_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagHeros : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagHeros_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.Heros = pPlayer.BagMoudle.HeroBag.Heros

	response.RetCode = msg.RE_SUCCESS

	return
}

//请求背包中的所有装备
func Hand_GetBagEquips(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagEquip_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagEquips : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagEquip_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.Equips = pPlayer.BagMoudle.EquipBag.Equips

	response.RetCode = msg.RE_SUCCESS
}

//请求背包中的所有英雄碎片
func Hand_GetBagHerosPiece(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagHerosPiece_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagHerosPiece : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagHerosPiece_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.HeroPieces = pPlayer.BagMoudle.HeroPieceBag.Items

	response.RetCode = msg.RE_SUCCESS
}

//请求背包中的所有宝物碎片
func Hand_GetBagGemPiece(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagGemPiece_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagEquipPiece : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagGemPiece_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.GemPieces = pPlayer.BagMoudle.GemPieceBag.Items

	response.RetCode = msg.RE_SUCCESS
}

//请求背包中的所有装备碎片
func Hand_GetBagEquipPiece(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagEquipPiece_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagEquipPiece : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagEquipPiece_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.EquipPieces = pPlayer.BagMoudle.EquipPieceBag.Items

	response.RetCode = msg.RE_SUCCESS
}

//请求背包中的所有宝物碎片
func Hand_GetBagGems(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagGems_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagGems : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagGems_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.Gems = pPlayer.BagMoudle.GemBag.Gems

	response.RetCode = msg.RE_SUCCESS
}

//请求背包里的道具
func Hand_GetBagItems(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagItems_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagItems : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagItems_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//是否是合法的请求
	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.Items = pPlayer.BagMoudle.NormalItemBag.Items

	response.RetCode = msg.RE_SUCCESS
}

//请求背包里的觉醒道具
func Hand_GetBagWakeItems(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagWakeItems_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagWakeItems : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagWakeItems_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.Items = pPlayer.BagMoudle.WakeItemBag.Items

	response.RetCode = msg.RE_SUCCESS
}

func Hand_SellItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SellItem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_SellItem : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_SellItem_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//是否是合法的请求
	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	var moneyid int = 0
	var moneynum int = 0

	var tempPos = 10000
	if req.ItemType == gamedata.TYPE_HERO {
		//进行参数检查
		for _, item := range req.Items {
			if pPlayer.BagMoudle.HeroBag.Heros[item.Pos].HeroID != item.ID {
				gamelog.Error("Hand_SellItem Error Invalid Pos:%d and id:%d", item.Pos, item.ID)
				response.RetCode = msg.RE_INVALID_PARAM
				return
			}

			if item.Pos > tempPos {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem error :  Wrong Squence: %d", item.Pos)
				return
			}

			tempPos = item.Pos

			pHeroInfo := gamedata.GetHeroInfo(item.ID)
			if pHeroInfo == nil {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid heroid :%d", item.ID)
				return
			}

			if !pHeroInfo.CanSell {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error item cannot sell")
				return
			}

			moneyid = pHeroInfo.SellID
			moneynum += pHeroInfo.SellPrice
		}
		for _, item := range req.Items {
			pPlayer.BagMoudle.RemoveHeroAt(item.Pos)
		}
		pPlayer.BagMoudle.DB_SaveHeroBag()
	} else if req.ItemType == gamedata.TYPE_EQUIPMENT {
		//进行参数检查
		for _, item := range req.Items {
			if pPlayer.BagMoudle.EquipBag.Equips[item.Pos].EquipID != item.ID {
				gamelog.Error("Hand_SellItem Error Invalid Pos:%d and id:%d", item.Pos, item.ID)
				response.RetCode = msg.RE_INVALID_PARAM
				return
			}

			if item.Pos > tempPos {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem error :  Wrong Squence: %d", item.Pos)
				return
			}

			tempPos = item.Pos

			pEquipInfo := gamedata.GetEquipmentInfo(item.ID)
			if pEquipInfo == nil {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid equip :%d", item.ID)
				return
			}

			if !pEquipInfo.CanSell {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error item cannot sell")
				return
			}

			moneyid = pEquipInfo.SellID[0]
			moneynum += pEquipInfo.SellPrice[0]
		}
		for _, item := range req.Items {
			pPlayer.BagMoudle.RemoveEquipAt(item.Pos)
		}
		pPlayer.BagMoudle.DB_SaveBagEquips()
	} else if req.ItemType == gamedata.TYPE_GEM {
		//进行参数检查
		for _, item := range req.Items {
			if pPlayer.BagMoudle.GemBag.Gems[item.Pos].GemID != item.ID {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid Pos:%d and id:%d", item.Pos, item.ID)
				return
			}

			if item.Pos > tempPos {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem error :  Wrong Squence: %d", item.Pos)
				return
			}

			tempPos = item.Pos

			pGemInfo := gamedata.GetGemInfo(item.ID)
			if pGemInfo == nil {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid gemid :%d", item.ID)
				return
			}

			if !pGemInfo.CanSell {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error item cannot sell")
				return
			}

			moneyid = pGemInfo.SellID
			moneynum += pGemInfo.SellPrice
		}
		for _, item := range req.Items {
			pPlayer.BagMoudle.RemoveGemAt(item.Pos)
		}
		pPlayer.BagMoudle.DB_SaveGemBag()
	} else if req.ItemType == gamedata.TYPE_PET {
		//进行参数检查
		for _, item := range req.Items {
			if pPlayer.BagMoudle.PetBag.Pets[item.Pos].PetID != item.ID {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid Pos:%d and id:%d", item.Pos, item.ID)
				return
			}

			if item.Pos > tempPos {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem error :  Wrong Squence: %d", item.Pos)
				return
			}

			tempPos = item.Pos

			pPetInfo := gamedata.GetPetInfo(item.ID)
			if pPetInfo == nil {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid gemid :%d", item.ID)
				return
			}

			if !pPetInfo.CanSell {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error item cannot sell")
				return
			}

			moneyid = pPetInfo.SellID
			moneynum += pPetInfo.SellPrice
		}
		for _, item := range req.Items {
			pPlayer.BagMoudle.RemovePetAt(item.Pos)
		}
		pPlayer.BagMoudle.DB_SavePetBag()
	} else if req.ItemType == gamedata.TYPE_HERO_PIECE {
		//进行参数检查
		for _, item := range req.Items {
			if pPlayer.BagMoudle.HeroPieceBag.Items[item.Pos].ItemID != item.ID {
				gamelog.Error("Hand_SellItem Error Invalid Pos:%d and id:%d", item.Pos, item.ID)
				response.RetCode = msg.RE_INVALID_PARAM
				return
			}

			pItemInfo := gamedata.GetItemInfo(item.ID)
			if pItemInfo == nil {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid gemid :%d", item.ID)
				return
			}

			moneyid = pItemInfo.SellID
			moneynum += pItemInfo.SellPrice * pPlayer.BagMoudle.HeroPieceBag.Items[item.Pos].ItemNum
		}
		for _, item := range req.Items {
			pPlayer.BagMoudle.RemoveHeroPiece(item.ID, pPlayer.BagMoudle.HeroPieceBag.Items[item.Pos].ItemNum)
		}
	} else if req.ItemType == gamedata.TYPE_EQUIP_PIECE {
		//进行参数检查
		for _, item := range req.Items {
			if pPlayer.BagMoudle.EquipPieceBag.Items[item.Pos].ItemID != item.ID {
				gamelog.Error("Hand_SellItem Error Invalid Pos:%d and id:%d", item.Pos, item.ID)
				response.RetCode = msg.RE_INVALID_PARAM
				return
			}

			pItemInfo := gamedata.GetItemInfo(item.ID)
			if pItemInfo == nil {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid gemid :%d", item.ID)
				return
			}

			moneyid = pItemInfo.SellID
			moneynum += pItemInfo.SellPrice * pPlayer.BagMoudle.EquipPieceBag.Items[item.Pos].ItemNum
		}

		for _, item := range req.Items {
			pPlayer.BagMoudle.RemoveEquipPiece(item.ID, pPlayer.BagMoudle.EquipPieceBag.Items[item.Pos].ItemNum)
		}
	} else if req.ItemType == gamedata.TYPE_PET_PIECE {
		//进行参数检查
		for _, item := range req.Items {
			if pPlayer.BagMoudle.PetPieceBag.Items[item.Pos].ItemID != item.ID {
				gamelog.Error("Hand_SellItem Error Invalid Pos:%d and id:%d", item.Pos, item.ID)
				response.RetCode = msg.RE_INVALID_PARAM
				return
			}

			pItemInfo := gamedata.GetItemInfo(item.ID)
			if pItemInfo == nil {
				response.RetCode = msg.RE_INVALID_PARAM
				gamelog.Error("Hand_SellItem Error Invalid gemid :%d", item.ID)
				return
			}

			moneyid = pItemInfo.SellID
			moneynum += pItemInfo.SellPrice * pPlayer.BagMoudle.PetPieceBag.Items[item.Pos].ItemNum
		}
		for _, item := range req.Items {
			pPlayer.BagMoudle.RemovePetPiece(item.ID, pPlayer.BagMoudle.PetPieceBag.Items[item.Pos].ItemNum)
		}
	}

	pPlayer.RoleMoudle.AddMoney(moneyid, moneynum)
	response.MoneyID = moneyid
	response.MoneyNum = moneynum
	response.RetCode = msg.RE_SUCCESS

	return
}

func Hand_UseItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_UseItem_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UseItem : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_UseItem_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	pItemInfo := gamedata.GetItemInfo(req.ItemID)
	if pItemInfo == nil {
		gamelog.Error("Hand_UseItem Error : Invalid ItemID :%d", req.ItemID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if !pPlayer.BagMoudle.IsItemEnough(req.ItemID, req.ItemNum) {
		gamelog.Error("Hand_UseItem Error : Not Enough Item!, id:%d, num:%d", req.ItemID, req.ItemNum)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	switch pItemInfo.SubType {
	case gamedata.SUB_TYPE_MONEY: //货币道具，使用后直接增加货币
		{
			pPlayer.RoleMoudle.AddMoney(pItemInfo.Data1, pItemInfo.Data2*req.ItemNum)
			pPlayer.BagMoudle.RemoveNormalItem(req.ItemID, req.ItemNum)
		}
	case gamedata.SUB_TYPE_ACTION: //行动力道具，使用后直接增加行动力
		{
			pPlayer.RoleMoudle.AddAction(pItemInfo.Data1, pItemInfo.Data2*req.ItemNum)
			pPlayer.BagMoudle.RemoveNormalItem(req.ItemID, req.ItemNum)
		}
	case gamedata.SUB_TYPE_FREE_WAR: //免战道具，使用后增加免战时间
		{
			pPlayer.RobModule.AddFreeWarTime(pItemInfo.Data1 * req.ItemNum)
			pPlayer.BagMoudle.RemoveNormalItem(req.ItemID, req.ItemNum)
		}
	case gamedata.SUB_TYPE_GIFT_BAG: //礼包道具, 使用后获得礼包里的道具
		{
			if pItemInfo.UseType == 17 {
				awardItem := gamedata.GetAwardItemByIndex(pItemInfo.Data1, req.Index)
				var item msg.MSG_ItemData
				item.ID = awardItem.ItemID
				item.Num = awardItem.ItemNum * req.ItemNum
				response.Items = append(response.Items, item)
				pPlayer.BagMoudle.AddAwardItem(item.ID, item.Num)
				pPlayer.BagMoudle.RemoveNormalItem(req.ItemID, req.ItemNum)
			} else {
				awardLst := gamedata.GetItemsFromAwardID(pItemInfo.Data1)
				for _, v := range awardLst {
					var item msg.MSG_ItemData
					item.ID = v.ItemID
					item.Num = v.ItemNum
					response.Items = append(response.Items, item)
				}
				pPlayer.BagMoudle.RemoveNormalItem(req.ItemID, req.ItemNum)
				pPlayer.BagMoudle.AddAwardItems(awardLst)
			}
		}
	case gamedata.SUB_TYPE_CHARGE: //礼包道具, 使用后获得礼包里的道具
		{

		}
	default:
		{
			gamelog.Error("Hand_UseItem Error : Item cannot be use!, itemid:%d, sumtype:%d", req.ItemID, pItemInfo.SubType)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}
	}

	response.RetCode = msg.RE_SUCCESS

	return
}

//请求背包中的所有宠物
func Hand_GetBagPets(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagPets_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagPets : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagPets_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.Pets = pPlayer.BagMoudle.PetBag.Pets
	response.RetCode = msg.RE_SUCCESS
	return
}

//请求背包中的所有宠物碎片
func Hand_GetBagPetsPiece(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req MSG_GetBagPetsPiece_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetBagPetsPiece : Unmarshal error!!!!")
		return
	}

	var response MSG_GetBagPetsPiece_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pPlayer *TPlayer = nil
	pPlayer, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if pPlayer == nil {
		return
	}

	response.PetPieces = pPlayer.BagMoudle.PetPieceBag.Items
	response.RetCode = msg.RE_SUCCESS
}

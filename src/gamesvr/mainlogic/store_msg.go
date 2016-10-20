package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"utility"
)

//! 获取所有商店信息
func Hand_GetAllStoreData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetAllStoreData_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetAllStoreInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_GetAllStoreData_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 结算免费次数
	player.StoreModule.CheckReset(utility.GetCurTime())

	//! 神将
	//! 获取当前神将商店信息
	response.RetCode = msg.RE_SUCCESS
	for _, v := range player.StoreModule.HeroShopItemLst {
		good := msg.MSG_StoreItem{}
		good.ID = v.ID
		good.Status = v.Status
		response.Hero.GoodsInfoLst = append(response.Hero.GoodsInfoLst, good)
	}

	response.Hero.FreeCount = player.StoreModule.HeroFreeRefreshCount

	refreshTime := player.StoreModule.HeroFreeRefreshTime - utility.GetCurTime()
	if refreshTime < 0 {
		refreshTime = 0
	}

	response.Hero.FreeRefeshTime = refreshTime
	response.Hero.FreeCountLimit = gamedata.StoreFreeRefreshTimes
	response.Hero.RefreshCount = player.StoreModule.HeroRefreshCount

	//! 获取当前觉醒商店信息
	for _, v := range player.StoreModule.AwakeShopItemLst {
		good := msg.MSG_StoreItem{}
		good.ID = v.ID
		good.Status = v.Status
		response.Awake.GoodsInfoLst = append(response.Awake.GoodsInfoLst, good)
	}

	//! 获取当前神将商店信息
	response.Awake.FreeCount = player.StoreModule.AwakeFreeRefreshCount

	refreshTime = player.StoreModule.AwakeFreeRefreshTime - utility.GetCurTime()
	if refreshTime < 0 {
		refreshTime = 0
	}

	response.Awake.FreeRefeshTime = refreshTime
	response.Awake.FreeCountLimit = gamedata.StoreFreeRefreshTimes
	response.Awake.RefreshCount = player.StoreModule.AwakeRefreshCount

	//! 战宠商店
	for _, v := range player.StoreModule.PetShopItemLst {
		good := msg.MSG_StoreItem{}
		good.ID = v.ID
		good.Status = v.Status
		response.Pet.GoodsInfoLst = append(response.Pet.GoodsInfoLst, good)
	}

	//! 获取当前神将商店信息
	response.Pet.FreeCount = player.StoreModule.PetFreeRefreshCount

	refreshTime = player.StoreModule.PetFreeRefreshTime - utility.GetCurTime()
	if refreshTime < 0 {
		refreshTime = 0
	}

	response.Pet.FreeRefeshTime = refreshTime
	response.Pet.FreeCountLimit = gamedata.StoreFreeRefreshTimes
	response.Pet.RefreshCount = player.StoreModule.PetRefreshCount

}

//! 获取神将商店信息
func Hand_GetStoreData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetStoreData_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetHeroStoreInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_GetStoreData_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 结算免费次数
	player.StoreModule.CheckReset(utility.GetCurTime())

	if req.StoreType == gamedata.StoreType_Hero {

		//! 获取当前神将商店信息
		response.RetCode = msg.RE_SUCCESS
		for _, v := range player.StoreModule.HeroShopItemLst {
			good := msg.MSG_StoreItem{}
			good.ID = v.ID
			good.Status = v.Status
			response.GoodsInfoLst = append(response.GoodsInfoLst, good)
		}

		//! 获取当前神将商店信息
		response.RetCode = msg.RE_SUCCESS
		response.FreeCount = player.StoreModule.HeroFreeRefreshCount

		refreshTime := player.StoreModule.HeroFreeRefreshTime - utility.GetCurTime()
		if refreshTime < 0 {
			refreshTime = 0
		}

		response.FreeRefeshTime = refreshTime
		response.FreeCountLimit = gamedata.StoreFreeRefreshTimes
		response.RefreshCount = player.StoreModule.HeroRefreshCount

	} else if req.StoreType == gamedata.StoreType_Awake {

		//! 获取当前觉醒商店信息
		response.RetCode = msg.RE_SUCCESS
		for _, v := range player.StoreModule.AwakeShopItemLst {
			good := msg.MSG_StoreItem{}
			good.ID = v.ID
			good.Status = v.Status
			response.GoodsInfoLst = append(response.GoodsInfoLst, good)
		}

		//! 获取当前神将商店信息
		response.RetCode = msg.RE_SUCCESS
		response.FreeCount = player.StoreModule.AwakeFreeRefreshCount

		refreshTime := player.StoreModule.AwakeFreeRefreshTime - utility.GetCurTime()
		if refreshTime < 0 {
			refreshTime = 0
		}

		response.FreeRefeshTime = refreshTime
		response.FreeCountLimit = gamedata.StoreFreeRefreshTimes
		response.RefreshCount = player.StoreModule.AwakeRefreshCount

	} else if req.StoreType == gamedata.StoreType_Pet {
		//! 获取当前战宠商店信息
		response.RetCode = msg.RE_SUCCESS
		for _, v := range player.StoreModule.PetShopItemLst {
			good := msg.MSG_StoreItem{}
			good.ID = v.ID
			good.Status = v.Status
			response.GoodsInfoLst = append(response.GoodsInfoLst, good)
		}

		//! 获取当前神将商店信息
		response.RetCode = msg.RE_SUCCESS
		response.FreeCount = player.StoreModule.PetFreeRefreshCount

		refreshTime := player.StoreModule.PetFreeRefreshTime - utility.GetCurTime()
		if refreshTime < 0 {
			refreshTime = 0
		}

		response.FreeRefeshTime = refreshTime
		response.FreeCountLimit = gamedata.StoreFreeRefreshTimes
		response.RefreshCount = player.StoreModule.PetRefreshCount
	}

}

//! 刷新神将商店
func Hand_RefreshStore(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RefreshStore_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RefreshHeroStore Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_RefreshStore_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_HERO_STORE, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 判断是否满足刷新条件
	ret, errcode := player.StoreModule.CheckRefreshDemand(req.StoreType)
	if ret == false {
		response.RetCode = errcode
		return
	}

	//! 扣除条件
	response.CostType, response.CostNum = player.StoreModule.PaymentTerms(req.StoreType)

	//! 开始刷新
	player.StoreModule.RefreshGoods(req.StoreType)

	//! 存储数据库
	player.StoreModule.DB_UpdateShopItemToDatabase(req.StoreType)

	//! 返回消息赋值
	if req.StoreType == gamedata.StoreType_Hero {
		response.RetCode = msg.RE_SUCCESS
		response.RefreshCount = player.StoreModule.HeroRefreshCount

		refreshTime := player.StoreModule.HeroFreeRefreshTime - utility.GetCurTime()
		if refreshTime < 0 {
			refreshTime = 0
		}
		response.FreeRefeshTime = refreshTime

		response.FreeCount = player.StoreModule.HeroFreeRefreshCount
		response.FreeCountLimit = gamedata.StoreFreeRefreshTimes
		for _, v := range player.StoreModule.HeroShopItemLst {
			good := msg.MSG_StoreItem{}
			good.ID = v.ID
			good.Status = v.Status
			response.GoodsInfoLst = append(response.GoodsInfoLst, good)
		}

		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_STORE_REFRESH, 1)

	} else if req.StoreType == gamedata.StoreType_Awake {
		response.RetCode = msg.RE_SUCCESS
		response.RefreshCount = player.StoreModule.AwakeRefreshCount

		refreshTime := player.StoreModule.AwakeFreeRefreshTime - utility.GetCurTime()
		if refreshTime < 0 {
			refreshTime = 0
		}
		response.FreeRefeshTime = refreshTime

		response.FreeCount = player.StoreModule.AwakeFreeRefreshCount
		response.FreeCountLimit = gamedata.StoreFreeRefreshTimes
		for _, v := range player.StoreModule.AwakeShopItemLst {
			good := msg.MSG_StoreItem{}
			good.ID = v.ID
			good.Status = v.Status
			response.GoodsInfoLst = append(response.GoodsInfoLst, good)
		}
	} else if req.StoreType == gamedata.StoreType_Pet {
		response.RetCode = msg.RE_SUCCESS
		response.RefreshCount = player.StoreModule.PetRefreshCount

		refreshTime := player.StoreModule.PetFreeRefreshTime - utility.GetCurTime()
		if refreshTime < 0 {
			refreshTime = 0
		}
		response.FreeRefeshTime = refreshTime

		response.FreeCount = player.StoreModule.PetFreeRefreshCount
		response.FreeCountLimit = gamedata.StoreFreeRefreshTimes
		for _, v := range player.StoreModule.PetShopItemLst {
			good := msg.MSG_StoreItem{}
			good.ID = v.ID
			good.Status = v.Status
			response.GoodsInfoLst = append(response.GoodsInfoLst, good)
		}
	}

}

//! 购买商品
func Hand_Store_Buy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_StoreBuyItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RefreshHeroStore Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_StoreBuyItem_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_HERO_STORE, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 检测商品状态
	ret, errcode := player.StoreModule.CheckGoodsStatus(req.Index, req.StoreType)
	if ret == false {
		response.RetCode = errcode
		return
	}

	//! 扣除金币 发放奖励
	player.StoreModule.PayGoods(req.Index, req.StoreType)

	response.RetCode = msg.RE_SUCCESS
}

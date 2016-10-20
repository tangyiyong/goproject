package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 查询折扣贩售活动信息
func Hand_QueryActivityDiscountSaleInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryActivity_DisountSale_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryActivityDiscountSaleInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryActivity_DisountSale_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.ActivityModule.CheckReset()

	var activity *TActivityDiscount
	for i, v := range player.ActivityModule.DiscountSale {
		if v.ActivityID == req.ActivityID {
			activity = &player.ActivityModule.DiscountSale[i]
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_QueryActivityDiscountSaleInfo Error: Activity not exist %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if G_GlobalVariables.IsActivityOpen(activity.ActivityID) == false {
		gamelog.Error("IsActivityOpen Error: Activity is not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	response.RetCode = msg.RE_SUCCESS
	response.ActivityID = req.ActivityID
	response.AwardType = G_GlobalVariables.GetActivityAwardType(req.ActivityID)
	//! 折扣贩售
	response.ShopLst = []msg.TActivityDiscount{}
	for _, v := range activity.ShopLst {
		var info msg.TActivityDiscount
		info.Index = v.Index
		info.BuyTimes = v.Times
		response.ShopLst = append(response.ShopLst, info)
	}
}

//! 玩家请求购买折扣商品
func Hand_BuyDiscountSaleItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyDiscountItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyDiscountSaleItem : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_BuyDiscountItem_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测当前是否有此活动
	var activity *TActivityDiscount
	var activityIndex int
	for i, v := range player.ActivityModule.DiscountSale {
		if v.ActivityID == req.ActivityID {
			activity = &player.ActivityModule.DiscountSale[i]
			activityIndex = i
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_BuyDiscountSaleItem Error: Activity not exist %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	player.ActivityModule.CheckReset()

	if G_GlobalVariables.IsActivityOpen(activity.ActivityID) == false {
		gamelog.Error("IsActivityOpen Error: Activity is not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if req.Index <= 0 {
		gamelog.Error("Hand_BuyDiscountSaleItem Error: Invalid Index %d", req.Index)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 获取商品信息
	activityInfo := gamedata.GetActivityInfo(req.ActivityID)
	goodsLst := gamedata.GetDiscountSaleInfo(activityInfo.AwardType)
	if req.Index > len(goodsLst) {
		gamelog.Error("Hand_BuyDiscountSaleItem Error: invalid index %d", req.Index)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	goodsInfo := goodsLst[req.Index-1]

	//! 检测购买次数是否足够
	shoppingInfo := player.ActivityModule.GetItemShoppingInfo(req.ActivityID, req.Index)
	if shoppingInfo == nil {
		//! 未曾购买,添加记录
		var info TDiscountSaleGoodsInfo
		info.Index = req.Index
		info.Times = 0

		for i, v := range player.ActivityModule.DiscountSale {
			if v.ActivityID == req.ActivityID {
				shoppingInfo = player.ActivityModule.DiscountSale[i].AddItem(info, activityIndex)
			}
		}
	}

	if shoppingInfo.Times+req.Count > goodsInfo.Times {
		gamelog.Error("Hand_BuyDiscountSaleItem Error: Times is use up  Now: %d  Limit: %d", shoppingInfo.Times, goodsInfo.Times)
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 判断货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(goodsInfo.MoneyID, goodsInfo.MoneyNum*req.Count) == false {
		gamelog.Error("Hand_BuyDiscountSaleItem Error: Money not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(goodsInfo.MoneyID, goodsInfo.MoneyNum*req.Count)
	response.MoneyID = goodsInfo.MoneyID
	response.MoneyNum = goodsInfo.MoneyNum * req.Count
	awardLst := gamedata.GetItemsFromAwardID(goodsInfo.Award)

	//! 判断是否为多选
	if goodsInfo.IsSelect == 1 {
		//! 判断用户选择是否正确
		if req.Choice > 4 || req.Choice < 1 {
			gamelog.Error("Hand_BuyDiscountSaleItem Error:  invalid choice %d", req.Choice)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}

		item := awardLst[req.Choice-1]
		if item.ItemID == 0 {
			gamelog.Error("Hand_BuyDiscountSaleItem Error:  invalid choice %d", req.Choice)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}

		//! 给予玩家奖励
		player.BagMoudle.AddAwardItem(item.ItemID, item.ItemNum*req.Count)
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{item.ItemID, item.ItemNum * req.Count})
	} else {

		for _, v := range awardLst {
			if v.ItemID != 0 && v.ItemNum != 0 {
				player.BagMoudle.AddAwardItem(v.ItemID, v.ItemNum*req.Count)
				response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum * req.Count})
			}

		}
	}

	//! 增加购买次数
	shoppingInfo.Times += req.Count
	index := 0
	for _, v := range player.ActivityModule.DiscountSale {
		if v.ActivityID == req.ActivityID {
			for j, n := range v.ShopLst {
				if n.Index == req.Index {
					index = j
					break
				}
			}
		}
	}

	activity.DB_UpdateShoppingTimes(activityIndex, index, shoppingInfo)

	response.ActivityID = req.ActivityID
	response.BuyNum = req.Count
	response.Index = req.Index
	response.RetCode = msg.RE_SUCCESS
}

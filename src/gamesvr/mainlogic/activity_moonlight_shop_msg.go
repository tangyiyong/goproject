package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

// 月光集市

type MSG_MoonlightShop_GetInfo_Ack struct {
	RetCode int
	Shop    TMoonlightShopData
}
type MSG_MoonlightShop_RefreshShop_Buy_Ack struct {
	RetCode int
	Goods   [MoonlightShop_Goods_Num]TMoonlightGoods
}
type MSG_MoonlightShop_RefreshShop_Auto_Ack struct {
	RetCode int
	Goods   [MoonlightShop_Goods_Num]TMoonlightGoods
}

func Hand_MoonlightShop_GetInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MoonlightShop_GetInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MoonlightShop_GetInfo unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_MoonlightShop_GetInfo_Ack
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

	info := &player.ActivityModule.MoonlightShop

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(info.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	info.RefreshShop_Auto(player)

	response.Shop = *info.GetShopDtad()
	response.RetCode = msg.RE_SUCCESS
}
func Hand_MoonlightShop_ExchangeToken(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MoonlightShop_ExchangeToken_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MoonlightShop_ExchangeToken unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MoonlightShop_ExchangeToken_Ack
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

	info := &player.ActivityModule.MoonlightShop

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(info.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if info.ExchangeToken(player, req.ExchangeID) {
		response.RetCode = msg.RE_SUCCESS
	}
}
func Hand_MoonlightShop_ReduceDiscount(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MoonlightShop_ReduceDiscount_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MoonlightShop_ReduceDiscount unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MoonlightShop_ReduceDiscount_Ack
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

	info := &player.ActivityModule.MoonlightShop

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(info.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	ok, newDiscount := info.ReduceDiscount(player, req.GoodsID)
	if ok {
		response.Discount = newDiscount
		response.RetCode = msg.RE_SUCCESS
	}
}
func Hand_MoonlightShop_RefreshShop_Buy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MoonlightShop_RefreshShop_Buy_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MoonlightShop_RefreshShop_Buy unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_MoonlightShop_RefreshShop_Buy_Ack
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

	shop := &player.ActivityModule.MoonlightShop

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(shop.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if shop.RefreshShop_Buy(player) {
		response.RetCode = msg.RE_SUCCESS
	} else {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
	}
	response.Goods = shop.Goods
}
func Hand_MoonlightShop_RefreshShop_Auto(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MoonlightShop_RefreshShop_Auto_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MoonlightShop_RefreshShop_Auto unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_MoonlightShop_RefreshShop_Auto_Ack
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

	shop := &player.ActivityModule.MoonlightShop

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(shop.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if shop.RefreshShop_Auto(player) {
		response.RetCode = msg.RE_SUCCESS
	}
	response.Goods = shop.Goods
}
func Hand_MoonlightShop_BuyGoods(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MoonlightShop_BuyGoods_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MoonlightShop_BuyGoods unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MoonlightShop_BuyGoods_Ack
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

	shop := &player.ActivityModule.MoonlightShop

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(shop.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if shop.BuyGoods(player, req.GoodsID) {
		response.RetCode = msg.RE_SUCCESS
	}
}
func Hand_MoonlightShop_GetScoreAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MoonlightShop_GetScoreAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MoonlightShop_GetScoreAward unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MoonlightShop_GetScoreAward_Ack
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

	shop := &player.ActivityModule.MoonlightShop

	//! 检测当前是否有此活动
	if G_GlobalVariables.IsActivityOpen(shop.ActivityID) == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if shop.GetScoreAward(player, req.AwardID) {
		response.RetCode = msg.RE_SUCCESS
	}
}

package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

type MSG_BeachBaby_Info_Ack struct {
	RetCode int
	Data    TBeachBabyGoodsData
}
type MSG_BeachBaby_OpenAllGoods_Ack struct {
	RetCode int
	Items   []msg.MSG_ItemData
	Goods   [BeachBaby_Goods_Num]TBeachBabyGoods
}
type MSG_BeachBaby_Refresh_Auto_Ack struct {
	RetCode int
	Goods   [BeachBaby_Goods_Num]TBeachBabyGoods
}
type MSG_BeachBaby_Refresh_Buy_Ack struct {
	RetCode int
	Goods   [BeachBaby_Goods_Num]TBeachBabyGoods
}

func Hand_BeachBaby_Info(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BeachBaby_Info_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BeachBaby_Info unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_BeachBaby_Info_Ack
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

	info := &player.ActivityModule.BeachBaby

	//! 检测当前是否有此活动
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	info.Refresh_Auto()

	response.Data = *info.GetBeachBabyDtad()
	response.RetCode = msg.RE_SUCCESS
}
func Hand_BeachBaby_OpenGoods(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BeachBaby_OpenGoods_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BeachBaby_OpenGoods unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BeachBaby_OpenGoods_Ack
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

	info := &player.ActivityModule.BeachBaby

	//! 检测当前是否有此活动
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	info.Refresh_Auto()

	item, bGetItem := info.OpenGoods(req.Index)
	response.IsGetItem = bGetItem
	response.Item.ID = item.ItemID
	response.Item.Num = item.ItemNum
	if item.ItemID > 0 {
		response.RetCode = msg.RE_SUCCESS
	}
}
func Hand_BeachBaby_OpenAllGoods(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BeachBaby_OpenAllGoods_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BeachBaby_OpenAllGoods unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_BeachBaby_OpenAllGoods_Ack
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

	info := &player.ActivityModule.BeachBaby

	//! 检测当前是否有此活动
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	info.Refresh_Auto()

	items, bSuccess := info.OpenAllGoods()
	if bSuccess {
		response.RetCode = msg.RE_SUCCESS
		for _, v := range items {
			response.Items = append(response.Items, msg.MSG_ItemData{v.ItemID, v.ItemNum})
		}
		response.Goods = info.Goods
	}
}
func Hand_BeachBaby_Refresh_Auto(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BeachBaby_Refresh_Auto_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BeachBaby_Refresh_Auto unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_BeachBaby_Refresh_Auto_Ack
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

	info := &player.ActivityModule.BeachBaby

	//! 检测当前是否有此活动
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if info.Refresh_Auto() {
		response.RetCode = msg.RE_SUCCESS
	}
	response.Goods = info.Goods
}
func Hand_BeachBaby_Refresh_Buy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BeachBaby_Refresh_Buy_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BeachBaby_Refresh_Buy unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_BeachBaby_Refresh_Buy_Ack
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

	info := &player.ActivityModule.BeachBaby

	//! 检测当前是否有此活动
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if info.Refresh_Buy() {
		response.RetCode = msg.RE_SUCCESS
	} else {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
	}
	response.Goods = info.Goods
}
func Hand_BeachBaby_GetFreeConch(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BeachBaby_GetFreeConch_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BeachBaby_GetFreeConch unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BeachBaby_GetFreeConch_Ack
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

	info := &player.ActivityModule.BeachBaby

	//! 检测当前是否有此活动
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	info.Refresh_Auto()

	if info.GetFreeConch() {
		response.RetCode = msg.RE_SUCCESS
	}
}
func Hand_BeachBaby_SelectGoodsID(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BeachBaby_SelectGoodsID_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BeachBaby_SelectGoodsID unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BeachBaby_SelectGoodsID_Ack
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

	info := &player.ActivityModule.BeachBaby

	//! 检测当前是否有此活动
	isHandleTime, _ := G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	info.Refresh_Auto()

	if info.SelectGoodsID(req.IDs) {
		response.RetCode = msg.RE_SUCCESS
	}
}

func BeachBabyRankAward(awardType int, player *TPlayer, indexToday int) (awardItem []gamedata.ST_ItemData, retCode int) {
	info := &player.ActivityModule.BeachBaby
	switch awardType {
	case 1: //昨日
		{
			if info.IsGetTodayRankAward {
				retCode = msg.RE_ALREADY_RECEIVED
				return awardItem, retCode
			}

			rankIdx := G_BeachBabyYesterdayRanker.GetRankIndex(player.playerid, info.GetYesterdayScore())
			if rankIdx <= 0 {
				retCode = msg.RE_NOT_ENOUGH_RANK
				return awardItem, retCode
			}

			awardType := G_GlobalVariables.GetActivityAwardType(info.ActivityID)
			rankAward := gamedata.GetOperationalRankAwardFromRank(gamedata.Activity_Beach_Baby, awardType, rankIdx)
			if rankAward == nil {
				retCode = msg.RE_UNKNOWN_ERR
				gamelog.Error("GetOperationalRankAwardFromRank get nil, rank: %d", rankIdx)
				return awardItem, retCode
			}

			//! 奖励物品
			awardItem := gamedata.GetItemsFromAwardID(rankAward.TodayNormalRankAward)
			if info.GetYesterdayScore() >= gamedata.BeachBaby_TodayRank_Limit {
				awardItem = append(awardItem, gamedata.GetItemsFromAwardID(rankAward.TodayEliteRankAward)...)
			}
			player.BagMoudle.AddAwardItems(awardItem)

			for z := 0; z < len(awardItem); z++ {
				awardItem = append(awardItem, awardItem[z])
			}

			//! 更改标记
			info.IsGetTodayRankAward = true
			info.DB_SaveRankAwardFlag()
			retCode = msg.RE_SUCCESS
		}
	case 2: // 累计
		{
			//! 检测当前是否有此活动
			isHandleTime, _ := G_GlobalVariables.IsActivityTime(info.ActivityID)
			if isHandleTime {
				retCode = msg.RE_ACTIVITY_NOT_OVER
				return awardItem, retCode
			}

			if info.IsGetTotalRankAward {
				retCode = msg.RE_ALREADY_RECEIVED
				return awardItem, retCode
			}

			rankIdx := G_BeachBabyTotalRanker.GetRankIndex(player.playerid, info.GetTotalScore())
			if rankIdx <= 0 {
				retCode = msg.RE_NOT_ENOUGH_RANK
				return awardItem, retCode
			}

			awardType := G_GlobalVariables.GetActivityAwardType(info.ActivityID)
			rankAward := gamedata.GetOperationalRankAwardFromRank(gamedata.Activity_Beach_Baby, awardType, rankIdx)
			if rankAward == nil {
				retCode = msg.RE_UNKNOWN_ERR
				gamelog.Error("GetOperationalRankAwardFromRank get nil, rank: %d", rankIdx)
				return awardItem, retCode
			}

			//! 奖励物品
			awardItem := gamedata.GetItemsFromAwardID(rankAward.TotalNormalRankAward)
			if info.GetTotalScore() >= gamedata.BeachBaby_TotalRank_Limit {
				awardItem = append(awardItem, gamedata.GetItemsFromAwardID(rankAward.TotalEliteRankAward)...)
			}
			player.BagMoudle.AddAwardItems(awardItem)

			for z := 0; z < len(awardItem); z++ {
				awardItem = append(awardItem, awardItem[z])
			}

			//! 更改标记
			info.IsGetTotalRankAward = true
			info.DB_SaveRankAwardFlag()
			retCode = msg.RE_SUCCESS
		}
	}
	return awardItem, retCode
}

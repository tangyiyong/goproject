package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

func Hand_CardMaster_CardList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_CardMaster_CardList_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_CardMaster_CardList unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_CardMaster_CardList_Ack
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

	//! 检测当前是否有此活动
	isHandleTime:= G_GlobalVariables.IsActivityTime(player.ActivityModule.CardMaster.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	for i, v := range player.ActivityModule.CardMaster.CardList {
		if v > 0 {
			response.Cards = append(response.Cards, msg.MSG_ItemData{i, int(v)})
		}
	}
	for i, v := range player.ActivityModule.CardMaster.ExchangeTimes {
		if v > 0 {
			response.ExchangeTimes = append(response.ExchangeTimes, msg.MSG_ItemData{i, int(v)})
		}
	}
	response.FreeTimes = player.ActivityModule.CardMaster.FreeDrawTimes
	response.Score = player.ActivityModule.CardMaster.GetTodayScore()
	response.Point = player.ActivityModule.CardMaster.Point
	response.RetCode = msg.RE_SUCCESS
}
func Hand_CardMaster_Draw(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_CardMaster_Draw_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_CardMaster_Draw unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_CardMaster_Draw_Ack
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

	info := &player.ActivityModule.CardMaster

	//! 检测当前是否有此活动
	isHandleTime:= G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	var List []gamedata.ST_ItemData = nil
	switch req.Type {
	case 1:
		{
			List = info.NormalDraw()
		}
	case 2, 3, 4:
		{
			List = info.SpecialDraw(req.Type)
		}
	}

	if List != nil {
		for _, v := range List {
			response.Cards = append(response.Cards, msg.MSG_ItemData{v.ItemID, v.ItemNum})
		}
		response.RetCode = msg.RE_SUCCESS
	}
}
func Hand_CardMaster_Card2Item(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_CardMaster_Card2Item_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_CardMaster_Card2Item unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_CardMaster_Card2Item_Ack
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

	info := &player.ActivityModule.CardMaster

	//! 检测当前是否有此活动
	isHandleTime:= G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if info.Card2Item(req.ExchangeID) {
		response.RetCode = msg.RE_SUCCESS
	}
}
func Hand_CardMaster_Card2Point(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_CardMaster_Card2Point_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_CardMaster_Card2Point unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_CardMaster_Card2Point_Ack
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

	info := &player.ActivityModule.CardMaster

	//! 检测当前是否有此活动
	isHandleTime:= G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	var CardList []gamedata.ST_ItemData
	for _, v := range req.Cards {
		CardList = append(CardList, gamedata.ST_ItemData{v.ID, v.Num})
	}
	if info.Card2Point(CardList) {
		response.RetCode = msg.RE_SUCCESS
		response.Point = info.Point
	}
}
func Hand_CardMaster_Point2Card(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_CardMaster_Point2Card_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_CardMaster_Point2Card unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_CardMaster_Point2Card_Ack
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

	info := &player.ActivityModule.CardMaster

	//! 检测当前是否有此活动
	isHandleTime:= G_GlobalVariables.IsActivityTime(info.ActivityID)
	if isHandleTime == false {
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	var CardList []gamedata.ST_ItemData
	for _, v := range req.Cards {
		CardList = append(CardList, gamedata.ST_ItemData{v.ID, v.Num})
	}
	if info.Point2Card(CardList) {
		response.RetCode = msg.RE_SUCCESS
		response.Point = info.Point
	}
}

func CardMasterRankAward(awardType int, player *TPlayer, indexToday int) (awardItem []gamedata.ST_ItemData, retCode int) {
	info := &player.ActivityModule.CardMaster
	switch awardType {
	case 1: //昨日
		{
			if info.RankAward[0] == 1 {
				retCode = msg.RE_ALREADY_RECEIVED
				return awardItem, retCode
			}

			rankIdx := G_CardMasterYesterdayRanker.GetRankIndex(player.playerid, info.GetYesterdayScore())
			if rankIdx <= 0 {
				retCode = msg.RE_NOT_ENOUGH_RANK
				return awardItem, retCode
			}

			awardType := G_GlobalVariables.GetActivityAwardType(info.ActivityID)
			rankAward := gamedata.GetOperationalRankAwardFromRank(gamedata.Activity_Card_Master, awardType, rankIdx)
			if rankAward == nil {
				retCode = msg.RE_UNKNOWN_ERR
				gamelog.Error("GetOperationalRankAwardFromRank get nil, rank: %d", rankIdx)
				return awardItem, retCode
			}

			//! 奖励物品
			awardItem := gamedata.GetItemsFromAwardID(rankAward.TodayNormalRankAward)
			if info.GetYesterdayScore() >= gamedata.CardMaster_TodayRank_Limit {
				awardItem = append(awardItem, gamedata.GetItemsFromAwardID(rankAward.TodayEliteRankAward)...)
			}
			player.BagMoudle.AddAwardItems(awardItem)

			for z := 0; z < len(awardItem); z++ {
				awardItem = append(awardItem, awardItem[z])
			}

			//! 更改标记
			info.RankAward[0] = 1
			info.DB_SaveRankAwardFlag()
			retCode = msg.RE_SUCCESS
		}
	case 2: // 累计
		{
			//! 检测当前是否有此活动
			isHandleTime:= G_GlobalVariables.IsActivityTime(info.ActivityID)
			if isHandleTime {
				retCode = msg.RE_ACTIVITY_NOT_OVER
				return awardItem, retCode
			}

			if info.RankAward[1] == 1 {
				retCode = msg.RE_ALREADY_RECEIVED
				return awardItem, retCode
			}

			rankIdx := G_CardMasterTotalRanker.GetRankIndex(player.playerid, info.GetTotalScore())
			if rankIdx <= 0 {
				retCode = msg.RE_NOT_ENOUGH_RANK
				return awardItem, retCode
			}

			awardType := G_GlobalVariables.GetActivityAwardType(info.ActivityID)
			rankAward := gamedata.GetOperationalRankAwardFromRank(gamedata.Activity_Card_Master, awardType, rankIdx)
			if rankAward == nil {
				retCode = msg.RE_UNKNOWN_ERR
				gamelog.Error("GetOperationalRankAwardFromRank get nil, rank: %d", rankIdx)
				return awardItem, retCode
			}

			//! 奖励物品
			awardItem := gamedata.GetItemsFromAwardID(rankAward.TotalNormalRankAward)
			if info.GetTotalScore() >= gamedata.CardMaster_TotalRank_Limit {
				awardItem = append(awardItem, gamedata.GetItemsFromAwardID(rankAward.TotalEliteRankAward)...)
			}
			player.BagMoudle.AddAwardItems(awardItem)

			for z := 0; z < len(awardItem); z++ {
				awardItem = append(awardItem, awardItem[z])
			}

			//! 更改标记
			info.RankAward[1] = 1
			info.DB_SaveRankAwardFlag()
			retCode = msg.RE_SUCCESS
		}
	}
	return awardItem, retCode
}

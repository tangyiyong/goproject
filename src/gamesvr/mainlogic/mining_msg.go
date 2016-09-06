package mainlogic

import (
	"encoding/json"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"math"
	"msg"
	"net/http"
	"time"
)

//! 玩家查询挖矿信息
func Hand_GetMiningInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetMiningInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMiningInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetMiningInfo_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 获取信息
	for i, v := range player.MiningModule.MiningMap {
		response.MapData.DigStatus[i] = fmt.Sprintf("%v", v)
	}

	response.MapData.MonsterInfo = []msg.MSG_MiningMonster{}
	for _, v := range player.MiningModule.MonsterLst {
		response.MapData.MonsterInfo = append(response.MapData.MonsterInfo, msg.MSG_MiningMonster{v.Index, v.ID, v.Life})
	}

	response.MapData.Element = []int{}
	response.MapData.Element = append(response.MapData.Element, player.MiningModule.Element...)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求矿洞状态码信息
func Hand_GetMiningStatusCode(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetMiningStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMiningStatusCode Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetMiningStatus_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if req.StatusCode == player.MiningModule.StatusCode {
		response.IsVerified = true
	} else {
		response.IsVerified = false
	}

	//! 获取信息
	response.StatusCode = player.MiningModule.StatusCode
	response.LastPos.X = player.MiningModule.LastPos.X
	response.LastPos.Y = player.MiningModule.LastPos.Y
	response.Point = player.MiningModule.Point

	//! 获取Buff列表
	response.Buff.BuffType = player.MiningModule.Buff.BuffType
	response.Buff.Times = player.MiningModule.Buff.Times

	if player.MiningModule.GuajiCalcTime == 0 {
		response.GuajiStatus = false
	} else {
		response.GuajiStatus = true
	}

	response.GuajiType = player.MiningModule.GuaJiType
	response.GuajiTime = player.MiningModule.GuajiCalcTime

	response.ResetTimes = player.MiningModule.MiningResetTimes

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求获取某个点的信息
func Hand_MiningDig(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningDig_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningDig Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningDig_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 获取上次操作点的元素
	posLst := []TMiningPos{}
	element := player.MiningModule.Element.Get(player.MiningModule.LastPos.Y*gamedata.MiningMapLength + player.MiningModule.LastPos.X)
	if element == gamedata.MiningEvent_Scanning {
		//! 若为扫描事件,可视区域为3*3
		for x := 1; x <= 3; x++ {
			for y := 1; y <= 3; y++ {
				posLst = append(posLst, TMiningPos{player.MiningModule.LastPos.X + x, player.MiningModule.LastPos.Y + y})
			}
		}

	} else {
		//! 非扫描事件
		//! 获取可视范围
		posLst = player.MiningModule.GetVisualPosArena(player.MiningModule.LastPos.X, player.MiningModule.LastPos.Y)

	}

	gamelog.Error("PosLst: %v    LastPos: %v", posLst, player.MiningModule.LastPos)

	for _, v := range req.Pos {
		isExist := false
		for _, n := range posLst {
			if v.X == n.X && v.Y == n.Y {
				isExist = true
				break
			}
		}

		if isExist == false {
			gamelog.Error("GetVisualPosArena Error: invalid pos: x: %d  y: %d", v.X, v.Y)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}
	}

	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	for _, v := range req.Pos {
		//! 判断坐标是否合法
		if v.X >= gamedata.MiningMapLength ||
			v.Y >= gamedata.MiningMapLength ||
			v.X < 0 ||
			v.Y < 0 {
			gamelog.Error("Invalid: x %d y %d", v.X, v.Y)
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}

		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.MapData = append(response.MapData, mapData)
	}

	response.StatusCode = player.MiningModule.StatusCode
	response.RetCode = msg.RE_SUCCESS
}

//! 挖矿事件-行动力奖励
func Hand_MiningEvent_ActionAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_ActionAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningEvent_AwardAction Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_ActionAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningCostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningEvent_Action_Award {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	//! 获取事件信息
	eventInfo := gamedata.GetMiningEventInfo(element)
	if eventInfo == nil {
		response.RetCode = msg.RE_UNKNOWN_ERR
		gamelog.Error("Get Mining event info fail. Event: %d", element)
		return
	}

	if isDig == true {
		//! 已经触发过该事件
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, gamedata.MiningCostActionNum)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X
	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	//! 删除事件
	player.MiningModule.DeleteElement(index)

	//! 随机一个行动力奖励
	actionAward := gamedata.MiningRandAward(1)

	//! 判断buff
	value := 1
	if player.MiningModule.Buff.BuffType == 2 {
		value = player.MiningModule.Buff.Value
		player.MiningModule.Buff.Times -= 1
		if player.MiningModule.Buff.Times == 0 {
			player.MiningModule.Buff = TMiningBuff{}
			player.MiningModule.DB_SavePlayerBuff()
		} else {
			player.MiningModule.DB_SubMiningBuffTimes(1)
		}
	}

	//! 奖励玩家积分
	player.MiningModule.Point += (eventInfo.Value1 * value)
	player.MiningModule.DB_SavePlayerPoint()

	response.Point = player.MiningModule.Point

	//! 奖励玩家行动值
	player.RoleMoudle.AddAction(gamedata.MiningCostActionID, actionAward)

	//! 增加版本号
	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	//! 生成可视区域
	posLst := player.MiningModule.GetNewVisualPosArena(req.PlayerPos.X, req.PlayerPos.Y)
	for _, v := range posLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	//! 返回参数
	response.StatusCode = player.MiningModule.StatusCode
	response.AddActionID = gamedata.MiningCostActionID
	response.AddActionNum = actionAward
	response.RetCode = msg.RE_SUCCESS
}

//! 挖矿事件-黑市
func Hand_MiningEvent_BalckMarket(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_GetBlackMarket_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningEvent_BalckMarket Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_GetBlackMarket_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Error("Return: %s", b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		gamelog.Error("Error pos: x %v  y: %v", req.PlayerPos.X, req.PlayerPos.Y)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningCostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningEvent_Black_Market {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	//! 获取事件信息
	eventInfo := gamedata.GetMiningEventInfo(element)
	if eventInfo == nil {
		response.RetCode = msg.RE_UNKNOWN_ERR
		gamelog.Error("Get Mining event info fail. Event: %d", element)
		return
	}

	if isDig == true {
		//! 已经触发过该元素
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, gamedata.MiningCostActionNum)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X
	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	//! 删除事件
	player.MiningModule.DeleteElement(index)

	player.MiningModule.BlackMarketBuyMark = IntLst{}
	player.MiningModule.DB_UpdateBlackMarketMark()

	//! 随机黑市商品
	goodsLst := gamedata.RandBlackMarketGoosLst(2, player.GetLevel())
	response.GoodsLst = append(response.GoodsLst, goodsLst...)
	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	//! 生成可视区域
	posLst := player.MiningModule.GetNewVisualPosArena(req.PlayerPos.X, req.PlayerPos.Y)
	for _, v := range posLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
	response.StatusCode = player.MiningModule.StatusCode
}

//! 挖矿事件-购买黑市商品
func Hand_MiningEvent_BuyBlackMarketItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_BuyBlackMarket_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningEvent_BuyBlackMarketItem Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_BuyBlackMarket_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if player.MiningModule.BlackMarketBuyMark.IsExist(req.ID) >= 0 {
		gamelog.Error("Hand_MiningEvent_BuyBlackMarketItem Error: Aleady <buy></buy>")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 获取商品信息
	goodsInfo := gamedata.GetMiningEventBlackMarketInfo(req.ID)
	if goodsInfo == nil {
		gamelog.Error("Hand_MiningEvent_BuyBlackMarketItem Error: GetMiningEventBlackMarketInfo nil")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测玩家金钱是否足够
	if player.RoleMoudle.CheckMoneyEnough(goodsInfo.MoneyID, goodsInfo.MoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 获取事件信息
	eventInfo := gamedata.GetMiningEventInfo(gamedata.MiningEvent_Black_Market)
	if eventInfo == nil {
		response.RetCode = msg.RE_UNKNOWN_ERR
		gamelog.Error("Get Mining event info fail. Event: %d", gamedata.MiningEvent_Black_Market)
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(goodsInfo.MoneyID, goodsInfo.MoneyNum)

	//! 奖励物品
	player.BagMoudle.AddAwardItem(goodsInfo.ItemID, goodsInfo.ItemNum)

	//! 玩家获取积分
	//! 判断buff
	value := 1
	if player.MiningModule.Buff.BuffType == 2 {
		value = player.MiningModule.Buff.Value
		player.MiningModule.Buff.Times -= 1
		if player.MiningModule.Buff.Times == 0 {
			player.MiningModule.Buff = TMiningBuff{}
			player.MiningModule.DB_SavePlayerBuff()
		} else {
			player.MiningModule.DB_SubMiningBuffTimes(1)
		}
	}

	//! 奖励玩家积分
	player.MiningModule.Point += (goodsInfo.Point * value)
	player.MiningModule.DB_SavePlayerPoint()

	response.Point = player.MiningModule.Point
	player.MiningModule.Point += eventInfo.Value1

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	//! 更改标记
	player.MiningModule.BlackMarketBuyMark = append(player.MiningModule.BlackMarketBuyMark, req.ID)
	player.MiningModule.DB_AddBlackMarketMark(req.ID)
}

//! 挖矿事件-怪物信息
func Hand_MiningEvent_MonsterInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_Monster_Info_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_miningEvent_Monster Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_Monster_Info_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.Pos.X >= gamedata.MiningMapLength ||
		req.Pos.Y >= gamedata.MiningMapLength ||
		req.Pos.X < 0 ||
		req.Pos.Y < 0 {
		gamelog.Error("Hand_MiningEvent_Monster Error: invalid pos x: %d  y: %d", req.Pos.X, req.Pos.Y)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 得到地图该点信息
	_, element, errcode := player.MiningModule.GetMapPosData(req.Pos.X, req.Pos.Y, true)
	if element != gamedata.MiningEvent_Boss &&
		element != gamedata.MiningEvent_Normal_Monster &&
		element != gamedata.MiningEvent_Elite_Monster {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	//! 获取事件信息
	eventInfo := gamedata.GetMiningEventInfo(element)
	if eventInfo == nil {
		response.RetCode = msg.RE_UNKNOWN_ERR
		gamelog.Error("Get Mining event info fail. Event: %d", element)
		return
	}

	index := req.Pos.Y*gamedata.MiningMapLength + req.Pos.X

	//! 获取怪物信息
	monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
	if monsterInfo == nil {
		gamelog.Error("Hand_MiningEvent_Monster Error: GetMonsterInfo error")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	monster := gamedata.GetMonsterEventInfo(monsterInfo.ID)

	response.Life = monsterInfo.Life
	response.Level = player.MiningModule.MiningResetTimes
	response.CopyID = monster.CopyID
	response.TotalLife = int(float64(monster.MonsterLife) * math.Pow(1.2, float64(player.MiningModule.MiningResetTimes)))
	response.MonsterType = eventInfo.Event
	response.RetCode = msg.RE_SUCCESS
}

//! 挖矿事件-怪物
func Hand_MiningEvent_Monster(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_Monster_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_miningEvent_Monster Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_Monster_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		gamelog.Error("Hand_MiningEvent_Monster Error: invalid pos x: %d  y: %d", req.PlayerPos.X, req.PlayerPos.Y)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningEvent_Boss &&
		element != gamedata.MiningEvent_Normal_Monster &&
		element != gamedata.MiningEvent_Elite_Monster {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	//! 检查行动力是否足够
	needAction := gamedata.MiningCostActionNum
	if element == gamedata.MiningEvent_Boss {
		needAction = gamedata.MiningAttackBossAction
	}
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, needAction) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 获取事件信息
	eventInfo := gamedata.GetMiningEventInfo(element)
	if eventInfo == nil {
		response.RetCode = msg.RE_UNKNOWN_ERR
		gamelog.Error("Get Mining event info fail. Event: %d", element)
		return
	}

	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X

	//! 获取怪物信息
	monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
	if monsterInfo == nil {
		gamelog.Error("Hand_MiningEvent_Monster Error: GetMonsterInfo error")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if isDig == true && monsterInfo.Life <= 0 {
		//! 已经触发过该元素
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	if element == gamedata.MiningEvent_Boss {
		// if player.MiningModule.GetSchedule() < 70 {
		// 	//! 完成度不足以挑战Boss
		// 	response.RetCode = msg.RE_MAP_NOT_COMPLETION
		// 	return
		// }
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, needAction)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	if player.MiningModule.Buff.BuffType == 3 {
		player.MiningModule.Buff.Times -= 1
		if player.MiningModule.Buff.Times == 0 {
			player.MiningModule.Buff = TMiningBuff{}
			player.MiningModule.DB_SavePlayerBuff()
		} else {
			player.MiningModule.DB_SubMiningBuffTimes(1)
		}
	}

	//! 给怪物造成伤害
	monsterInfo.Life -= req.Damage
	if monsterInfo.Life < 0 {
		monsterInfo.Life = 0
	}

	//! 存储怪物剩余血量
	player.MiningModule.DB_SetMonsterLife(index, monsterInfo.Life)

	if monsterInfo.Life <= 0 {
		//! 怪物死亡
		monsterInfo.Life = 0
		response.IsKill = true

		//! 删除事件
		player.MiningModule.DeleteElement(index)

		//! 删除怪物
		if element != gamedata.MiningEvent_Boss {
			player.MiningModule.DeleteMonster(index)
		}

		//! 击杀获取积分
		//! 判断buff
		value := 1
		if player.MiningModule.Buff.BuffType == 2 {
			value = player.MiningModule.Buff.Value
			player.MiningModule.Buff.Times -= 1
			if player.MiningModule.Buff.Times == 0 {
				player.MiningModule.Buff = TMiningBuff{}
				player.MiningModule.DB_SavePlayerBuff()
			} else {
				player.MiningModule.DB_SubMiningBuffTimes(1)
			}
		}

		//! 奖励玩家积分
		player.MiningModule.Point += (eventInfo.Value1 * value)

		if element == gamedata.MiningEvent_Boss {
			player.MiningModule.Point += gamedata.MiningBossValue
		}

		response.Point = player.MiningModule.Point
		player.MiningModule.DB_SavePlayerPoint()

	} else {
		response.IsKill = false
	}

	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	if element != gamedata.MiningEvent_Boss {
		monsterData := gamedata.GetMonsterEventInfo(monsterInfo.ID)
		copyBase := gamedata.GetCopyBaseInfo(monsterData.CopyID)
		awardLst := gamedata.GetItemsFromAwardID(copyBase.AwardID)
		player.BagMoudle.AddAwardItems(awardLst)

		for _, v := range awardLst {
			var item msg.MSG_ItemData
			item.ID = v.ItemID
			item.Num = v.ItemNum
			response.DropItem = append(response.DropItem, item)
		}
	}

	//! 生成可视区域
	posLst := player.MiningModule.GetNewVisualPosArena(req.PlayerPos.X, req.PlayerPos.Y)
	for _, v := range posLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			gamelog.Info("MonsterInfo: %v", monsterInfo)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	response.Point = player.MiningModule.Point

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 挖矿事件-宝箱
func Hand_MiningEvent_Treasure(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_Treasure_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningEvent_Treasure Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_Treasure_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningCostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningEvent_Treasure {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	if isDig == true {
		//! 已经触发过该元素
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	//! 获取事件信息
	eventInfo := gamedata.GetMiningEventInfo(element)
	if eventInfo == nil {
		response.RetCode = msg.RE_UNKNOWN_ERR
		gamelog.Error("Get Mining event info fail. Event: %d", element)
		return
	}

	//! 随机一个宝箱奖励
	treasureAward := gamedata.RandMiningTreasure()
	if treasureAward == 0 {
		//! 随机宝箱奖励失败
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, gamedata.MiningCostActionNum)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X
	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	//! 删除事件
	player.MiningModule.DeleteElement(index)

	//! 给予奖励
	awardItems := gamedata.GetItemsFromAwardID(treasureAward)
	player.BagMoudle.AddAwardItems(awardItems)

	for _, v := range awardItems {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.AwardItem = append(response.AwardItem, item)
	}

	//! 增加积分
	//! 判断buff
	value := 1
	if player.MiningModule.Buff.BuffType == 2 {
		value = player.MiningModule.Buff.Value
		player.MiningModule.Buff.Times -= 1
		if player.MiningModule.Buff.Times == 0 {
			player.MiningModule.Buff = TMiningBuff{}
			player.MiningModule.DB_SavePlayerBuff()
		} else {
			player.MiningModule.DB_SubMiningBuffTimes(1)
		}
	}

	//! 奖励玩家积分
	player.MiningModule.Point += (eventInfo.Value1 * value)
	player.MiningModule.DB_SavePlayerPoint()

	response.Point = player.MiningModule.Point

	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	//! 生成可视区域
	posLst := player.MiningModule.GetNewVisualPosArena(req.PlayerPos.X, req.PlayerPos.Y)
	for _, v := range posLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 挖矿事件-魔盒
func Hand_MiningEvent_MagicBox(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_Box_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningEvent_Treasure Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_Box_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}
	//! 检查行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningCostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningEvent_MagicBox {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	if isDig == true {
		//! 已经触发过该元素
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, gamedata.MiningCostActionNum)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X
	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	//! 删除事件
	player.MiningModule.DeleteElement(index)

	//! 随机一个魔盒奖励
	response.RandPoint = gamedata.MiningRandAward(2)

	//! 修改玩家积分
	//! 判断buff
	value := 1
	if player.MiningModule.Buff.BuffType == 2 {
		value = player.MiningModule.Buff.Value
		player.MiningModule.Buff.Times -= 1
		if player.MiningModule.Buff.Times == 0 {
			player.MiningModule.Buff = TMiningBuff{}
			player.MiningModule.DB_SavePlayerBuff()
		} else {
			player.MiningModule.DB_SubMiningBuffTimes(1)
		}
	}

	//! 奖励玩家积分
	player.MiningModule.Point += (response.RandPoint * value)
	response.RandPoint *= value
	player.MiningModule.DB_SavePlayerPoint()

	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	//! 生成可视区域
	posLst := player.MiningModule.GetNewVisualPosArena(req.PlayerPos.X, req.PlayerPos.Y)
	for _, v := range posLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 挖矿事件-扫描
func Hand_MiningEvent_Scan(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_Scan_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningEvent_Treasure Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_Scan_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}
	//! 检查行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningCostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningEvent_Scanning {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	eventInfo := gamedata.GetMiningEventInfo(gamedata.MiningEvent_Scanning)

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}
	if isDig == true {
		//! 已经触发过该元素
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, gamedata.MiningCostActionNum)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X
	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	//! 删除事件
	player.MiningModule.DeleteElement(index)

	//! 增加积分
	//! 判断buff
	value := 1
	if player.MiningModule.Buff.BuffType == 1 {
		value = player.MiningModule.Buff.Value
		player.MiningModule.Buff.Times -= 1
		if player.MiningModule.Buff.Times == 0 {
			player.MiningModule.Buff = TMiningBuff{}
			player.MiningModule.DB_SavePlayerBuff()
		} else {
			player.MiningModule.DB_SubMiningBuffTimes(1)
		}
	}

	//! 奖励玩家积分
	player.MiningModule.Point += (eventInfo.Value1 * value)
	player.MiningModule.DB_SavePlayerPoint()

	response.Point = player.MiningModule.Point

	posLst := []TMiningPos{}
	for x := -3; x <= 3; x++ {
		for y := -3; y <= 3; y++ {
			if req.PlayerPos.X+x >= 0 && req.PlayerPos.X+x < gamedata.MiningMapLength &&
				req.PlayerPos.Y+y >= 0 && req.PlayerPos.Y+y < gamedata.MiningMapLength {
				posLst = append(posLst, TMiningPos{req.PlayerPos.X + x, req.PlayerPos.Y + y})
			}

		}
	}

	//! 去除已经可以看见的区域
	visualLst := []TMiningPos{}
	for _, v := range posLst {
		index := v.Y*gamedata.MiningMapLength + v.X

		//! 去除已挖掘
		if player.MiningModule.MiningMap.Get(index) == true {
			continue
		}

		//! 去除已有事件
		if player.MiningModule.Element.Get(index) != 0 {
			continue
		}

		visualLst = append(visualLst, TMiningPos{v.X, v.Y})
	}

	//! 生成信息
	for _, v := range visualLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	response.RetCode = msg.RE_SUCCESS
	response.StatusCode = player.MiningModule.StatusCode
}

//! 挖矿事件-答题
func Hand_MiningEvent_Question(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_Question_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningEvent_Question Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_Question_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}
	//! 检查行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningCostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningEvent_Question {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	if isDig == true {
		//! 已经触发过该元素
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	//! 检测题目得分
	if req.AddPoint > 15 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, gamedata.MiningCostActionNum)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X
	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	//! 删除事件
	player.MiningModule.DeleteElement(index)

	//! 增加玩家积分
	//! 判断buff
	eventInfo := gamedata.GetMiningEventInfo(gamedata.MiningEvent_Question)
	value := 1
	if player.MiningModule.Buff.BuffType == 2 {
		value = player.MiningModule.Buff.Value
		player.MiningModule.Buff.Times -= 1
		if player.MiningModule.Buff.Times == 0 {
			player.MiningModule.Buff = TMiningBuff{}
			player.MiningModule.DB_SavePlayerBuff()
		} else {
			player.MiningModule.DB_SubMiningBuffTimes(1)
		}
	}

	//! 奖励玩家积分
	player.MiningModule.Point += (eventInfo.Value1 * value)
	player.MiningModule.DB_SavePlayerPoint()

	response.Point = player.MiningModule.Point

	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	//! 生成可视区域
	posLst := player.MiningModule.GetNewVisualPosArena(req.PlayerPos.X, req.PlayerPos.Y)
	for _, v := range posLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	response.Point = player.MiningModule.Point
	response.RetCode = msg.RE_SUCCESS
}

//! 挖矿事件-Buff
func Hand_MiningEvent_Buff(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningEvent_Buff_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningEvent_Buff Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningEvent_Buff_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}
	//! 检查行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningCostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningEvent_Buff {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	if isDig == true {
		//! 已经触发过该元素
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, gamedata.MiningCostActionNum)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X
	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	//! 删除事件
	player.MiningModule.DeleteElement(index)

	//! 随机一个Buff值
	randBuffInfo := gamedata.RandMiningEventBuff()

	//! 添加对应记录
	player.MiningModule.AddBuff(randBuffInfo.BuffType, randBuffInfo.Times, randBuffInfo.Value)

	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	//! 判断buff
	eventInfo := gamedata.GetMiningEventInfo(gamedata.MiningEvent_Buff)
	// value := 1
	// if player.MiningModule.Buff.BuffType == 2 {
	// 	value = player.MiningModule.Buff.Value
	// 	player.MiningModule.Buff.Times -= 1
	// 	if player.MiningModule.Buff.Times == 0 {
	// 		player.MiningModule.Buff = TMiningBuff{}
	// 		player.MiningModule.DB_SavePlayerBuff()
	// 	} else {
	// 		player.MiningModule.DB_SubMiningBuffTimes(1)
	// 	}
	// }

	//! 奖励玩家积分
	player.MiningModule.Point += eventInfo.Value1
	player.MiningModule.DB_SavePlayerPoint()

	response.Point = player.MiningModule.Point

	//! 生成可视区域
	posLst := player.MiningModule.GetNewVisualPosArena(req.PlayerPos.X, req.PlayerPos.Y)
	for _, v := range posLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)
			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	response.Buff.BuffType = randBuffInfo.BuffType
	response.Buff.Value = randBuffInfo.Value
	response.Buff.Times = randBuffInfo.Times

	response.RetCode = msg.RE_SUCCESS
}

//! 挖矿事件-精炼石
func Hand_MiningElement_RefiningStone(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningElement_GetStone_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningElement_RefiningStone Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningElement_GetStone_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 判断坐标是否合法
	if req.PlayerPos.X >= gamedata.MiningMapLength ||
		req.PlayerPos.Y >= gamedata.MiningMapLength ||
		req.PlayerPos.X < 0 ||
		req.PlayerPos.Y < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Invalid Mining Pos : X: %d  Y: %d", req.PlayerPos.X, req.PlayerPos.Y)
		return
	}
	//! 检查行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningCostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 得到地图该点信息
	isDig, element, errcode := player.MiningModule.GetMapPosData(req.PlayerPos.X, req.PlayerPos.Y, true)
	if element != gamedata.MiningElement_Lower_Refining_Stone &&
		element != gamedata.MiningElement_Intermediate_Refining_Stone &&
		element != gamedata.MiningElement_Advanced_Refining_Stone &&
		element != gamedata.MiningElement_Ultimate_Refining_Stone &&
		element != gamedata.MinintElement_Can_Break_Obstacle {
		//! 地图信息不匹配
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	if errcode != msg.RE_SUCCESS {
		response.RetCode = errcode
		return
	}

	if isDig == true {
		//! 已经触发过该元素
		response.RetCode = msg.RE_ALREADY_DIG
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, gamedata.MiningCostActionNum)

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

	//! 记录已挖掘状态
	player.MiningModule.LastPos.X = req.PlayerPos.X
	player.MiningModule.LastPos.Y = req.PlayerPos.Y
	index := req.PlayerPos.Y*gamedata.MiningMapLength + req.PlayerPos.X

	player.MiningModule.MiningMap.Set(index)
	player.MiningModule.DB_DigMining(req.PlayerPos.Y, player.MiningModule.MiningMap[req.PlayerPos.Y])

	//! 删除事件
	player.MiningModule.DeleteElement(index)

	itemValue := 1
	if player.MiningModule.Buff.BuffType == 1 {
		itemValue = player.MiningModule.Buff.Value
		player.MiningModule.Buff.Times -= 1
		if player.MiningModule.Buff.Times == 0 {
			player.MiningModule.Buff = TMiningBuff{}
			player.MiningModule.DB_SavePlayerBuff()
		} else {
			player.MiningModule.DB_SubMiningBuffTimes(1)
		}
	}

	if element != gamedata.MinintElement_Can_Break_Obstacle {
		//! 奖励玩家物品
		elementInfo := gamedata.GetMiningElementInfo(element)

		itemNum := gamedata.RandStoneNum(elementInfo.ItemID, player.GetLevel())
		player.BagMoudle.AddAwardItem(elementInfo.ItemID, itemNum)
		response.ItemID = elementInfo.ItemID
		response.ItemNum = itemNum * itemValue

		value := 1
		if player.MiningModule.Buff.BuffType == 2 {
			value = player.MiningModule.Buff.Value
			player.MiningModule.Buff.Times -= 1
			if player.MiningModule.Buff.Times == 0 {
				player.MiningModule.Buff = TMiningBuff{}
				player.MiningModule.DB_SavePlayerBuff()
			} else {
				player.MiningModule.DB_SubMiningBuffTimes(1)
			}
		}

		//! 奖励玩家积分
		player.MiningModule.Point += (1 * value)
		player.MiningModule.DB_SavePlayerPoint()

	}

	//! 版本号改变
	player.MiningModule.AddMiningStatusCode()
	response.StatusCode = player.MiningModule.StatusCode

	response.Point = player.MiningModule.Point

	//! 回复成功
	response.RetCode = msg.RE_SUCCESS

	//! 生成可视区域
	posLst := player.MiningModule.GetNewVisualPosArena(req.PlayerPos.X, req.PlayerPos.Y)
	for _, v := range posLst {
		//! 获取该点地图信息
		_, element, errcode := player.MiningModule.GetMapPosData(v.X, v.Y, true)
		if errcode != msg.RE_SUCCESS {
			response.RetCode = errcode
			return
		}

		var mapData msg.MSG_MiningDigData
		mapData.Element = (v.Y*gamedata.MiningMapLength+v.X)<<16 + element

		if element == gamedata.MiningEvent_Elite_Monster ||
			element == gamedata.MiningEvent_Normal_Monster ||
			element == gamedata.MiningEvent_Boss {
			index := v.Y*gamedata.MiningMapLength + v.X
			monsterInfo, _ := player.MiningModule.GetMonsterInfo(index)

			mapData.Monster = monsterInfo.ID
			mapData.MonsterLife = monsterInfo.Life
		}

		response.VisualPos = append(response.VisualPos, mapData)
	}

	response.StatusCode = player.MiningModule.StatusCode
}

//! 玩家请求随机九种打完Boss翻牌奖励
func Hand_GetRandBossAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningGetAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetRandBossAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningGetAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	needPoint := gamedata.MiningBossValue
	if player.MiningModule.Point < needPoint {
		response.RetCode = msg.RE_NOT_ENOUGH_POINT
		return
	}

	//! 检测Boss是不是已经死亡
	monster, _ := player.MiningModule.GetBossInfo()
	if monster.Life > 0 {
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	//! 随机翻牌奖励
	for _, v := range player.MiningModule.BossAward {
		var award msg.MSG_MiningAward
		award.ID = v.ID
		award.ItemID = v.ItemID
		award.ItemNum = v.ItemNum
		award.Status = v.Status
		response.AwardLst = append(response.AwardLst, award)
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求选择翻牌奖励
func Hand_SelectBossAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningSelectAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SelectBossAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningSelectAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测积分
	needPoint := gamedata.MiningBossValue
	if player.MiningModule.Point < needPoint {
		response.RetCode = msg.RE_NOT_ENOUGH_POINT
		return
	}

	awardIndex := -1
	for i, v := range player.MiningModule.BossAward {
		if v.ID == req.SelectID {
			awardIndex = i
		}
	}

	if awardIndex < 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测Boss是不是已经死亡
	monster, _ := player.MiningModule.GetBossInfo()
	if monster.Life > 0 {
		response.RetCode = msg.RE_INVALID_EVENT
		gamelog.Error("url: %s RE_INVALID_EVENT", r.URL.String())
		return
	}

	//! 获取奖励内容
	award := player.MiningModule.GetBossAward(req.SelectID)
	if award == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 奖励玩家
	player.BagMoudle.AddAwardItem(award.ItemID, award.ItemNum)

	//! 记录领取
	player.MiningModule.BossAward[req.SelectID].Status = true
	player.MiningModule.DB_UpdateMiningStatusCode()

	//! 版本号改变
	player.MiningModule.AddMiningStatusCode()
	response.IsEnd = false

	player.MiningModule.DB_SavePlayerPoint()
	player.MiningModule.DB_UpdateMiningStatusCode()

	response.Point = player.MiningModule.Point

	//! 计算剩余积分
	player.MiningModule.Point -= gamedata.MiningBossValue
	if player.MiningModule.Point < gamedata.MiningBossValue {
		response.IsEnd = true

		//! 通关矿区
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PASS_MINING, 1)

		//! 重置地图
		player.MiningModule.ResetMiningMap()
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求挂机
func Hand_MiningGuaji(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningGuaji_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_MiningGuaji Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningGuaji_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测当前是否正在挂机
	if player.MiningModule.GuajiCalcTime != 0 {
		response.RetCode = msg.RE_REPEATED_GUAJI
		return
	}

	//! 获取挂机信息
	guajiInfo := gamedata.GetMiningGuajiInfo(req.ID)
	if guajiInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测行动力是否足够
	if player.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, guajiInfo.CostActionNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_ACTION
		return
	}

	//! 扣除行动力
	player.RoleMoudle.CostAction(gamedata.MiningCostActionID, guajiInfo.CostActionNum)

	//! 设置挂机
	player.MiningModule.SetGuaji(req.ID, guajiInfo.Hour)

	response.RetCode = msg.RE_SUCCESS
	response.GuajiCalcTime = player.MiningModule.GuajiCalcTime

	//! 获取体力值与体力恢复时间
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.MiningCostActionID)

}

//! 玩家请求查询挂机倒计时
func Hand_GetMiningGuajiTime(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_MiningGuajiTime_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMiningGuajiTime Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_MiningGuajiTime_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if player.MiningModule.GuajiCalcTime == 0 {
		response.GuajiStatus = false
	} else {
		response.GuajiStatus = true
	}

	response.GuajiTime = player.MiningModule.GuajiCalcTime

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求领取挂机奖励
func Hand_GetMiningGuajiAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetMiningGuajiAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMiningGuajiAward Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetMiningGuajiAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR

	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)

	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if player.MiningModule.GuaJiType == 0 {
		//! 当前并未挂机
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if player.MiningModule.GuajiCalcTime > time.Now().Unix() {
		//! 挂机时间未到
		response.RetCode = msg.RE_NOT_ENOUGH_PATROL_TIME
		return
	}

	//! 获取挂机奖励
	itemLst := player.MiningModule.GetGuajiAward()

	for _, v := range itemLst {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}
	response.RetCode = msg.RE_SUCCESS
}

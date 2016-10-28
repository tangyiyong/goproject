package mainlogic

import (
	"appconfig"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"io/ioutil"
	"msg"
	"net/http"
	"os"
	"strconv"
	"strings"
	"utility"
)

func Hand_SendAwardToPlayer(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_Send_Award_Player_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SendAwardToPlayer unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_Send_Award_Player_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_SendAwardToPlayer Error Invalid Gm request!!!")
		return
	}

	var data TAwardData
	data.TextType = Text_GM_Mail
	data.Value = append(data.Value, req.Value)
	for _, v := range req.ItemLst {
		data.ItemLst = append(data.ItemLst, gamedata.ST_ItemData{v.ID, v.Num})
	}

	SendAwardToPlayer(req.TargetID, &data)
	response.RetCode = msg.RE_SUCCESS
}
func Hand_AddSvrAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SvrAward_Add_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_AddSvrAward unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SvrAward_Add_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_AddSvrAward Error Invalid Gm request!!!")
		return
	}

	var data TAwardData
	data.TextType = Text_GM_Mail
	data.Value = append(data.Value, req.Value)
	for _, v := range req.ItemLst {
		data.ItemLst = append(data.ItemLst, gamedata.ST_ItemData{v.ID, v.Num})
	}

	G_GlobalVariables.AddSvrAward(&data)
	response.RetCode = msg.RE_SUCCESS
}
func Hand_DelSvrAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SvrAward_Del_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_DelSvrAward unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_SvrAward_Del_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_DelSvrAward Error Invalid Gm request!!!")
		return
	}

	G_GlobalVariables.DelSvrAward(req.ID)
	response.RetCode = msg.RE_SUCCESS
}

func Hand_UpdateGameData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	msglen := binary.LittleEndian.Uint32(buffer[:4])
	var req msg.MSG_UpdateGameData_Req
	if json.Unmarshal(buffer[4:4+msglen], &req) != nil {
		gamelog.Error("Hand_UpdateGameData : Unmarshal error!!!!")
		return
	}

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_UpdateGameData Error Invalid Gm request!!!")
		return
	}

	b, _ := utility.UnCompressData(buffer[4+msglen:])
	fileN := utility.GetCurrCsvPath() + req.TbName + ".csv"
	ioutil.WriteFile(fileN, b, 777)
	gamedata.ReloadOneFile(fileN)
	OnConfigChange(req.TbName)
	var response msg.MSG_UpdateGameData_Ack
	response.RetCode = msg.RE_SUCCESS
	ret, _ := json.Marshal(&response)
	w.Write(ret)
	return

}

func Hand_GetServerInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	var response msg.MSG_GetServerInfo_Ack
	response.RetCode = msg.RE_SUCCESS
	response.SvrID = int32(appconfig.GameSvrID)
	response.SvrName = appconfig.GameSvrName
	response.OnlineCnt = G_OnlineCnt
	response.MaxOnlineCnt = G_MaxOnlineCnt
	response.RegisterCnt = G_RegisterCnt

	//	var ms runtime.MemStats
	//	runtime.ReadMemStats(&ms)
	//	response.MemAlloc = ms.HeapAlloc / 1024 / 1024
	//	response.MemInuse = ms.HeapSys / 1024 / 1024
	//	response.MenObjNum = ms.HeapObjects

	ret, _ := json.Marshal(&response)

	w.Write(ret)
	return

}

var clientlog *os.File = nil

func Hand_SaveClientInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var err error
	if clientlog == nil {
		clientlog, err = os.OpenFile(utility.GetCurrPath()+"log/client.log", os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			gamelog.Error("Hand_SaveClientInfo Error : %s", err.Error())
			return
		}
	}

	clientlog.Write(buffer)
	return
}

func Hand_QueryAccountID(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_QueryAccountID_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryAccountID unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_QueryAccountID_Ack
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	response.AccountID = G_SimpleMgr.GetPlayerIDByName(req.Name)

	return
}

func Hand_QueryPlayerInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_QueryPlayerInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryPlayerInfo unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_QueryPlayerInfo_Ack
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.PlayerID == 0 {
		response.PlayerID = G_SimpleMgr.GetPlayerIDByName(req.PlayerName)
		if response.PlayerID == 0 {
			gamelog.Error("Hand_QueryPlayerInfo Error: Player not exist id: %d  name: %s", req.PlayerID, req.PlayerName)
			response.RetCode = msg.RE_ACCOUNT_NOT_EXIST
			return
		}
	} else {
		response.PlayerID = req.PlayerID
	}

	simpleInfo := G_SimpleMgr.GetSimpleInfoByID(response.PlayerID)
	if simpleInfo == nil {
		gamelog.Error("Hand_QueryPlayerInfo Error: Player not exist id: %d  name: %s", req.PlayerID, req.PlayerName)
		response.RetCode = msg.RE_ACCOUNT_NOT_EXIST
		return
	}

	player := GetPlayerByID(response.PlayerID)
	if player == nil {
		player = LoadPlayerFromDB(response.PlayerID)
	}

	response.PlayerName = simpleInfo.Name
	response.Level = simpleInfo.Level
	response.VIPLevel = simpleInfo.VipLevel
	response.LastLogoffTime = simpleInfo.LogoffTime
	response.IsOnline = simpleInfo.isOnline
	response.FightValue = simpleInfo.FightValue
	response.Charge = player.RoleMoudle.TotalCharge

	for i := 0; i < 14; i++ {
		response.Money[i] = player.RoleMoudle.GetMoney(i + 1)
	}

	response.Strength = player.RoleMoudle.GetAction(1)
	response.Action = player.RoleMoudle.GetAction(2)
	response.AttackTimes = player.RoleMoudle.GetAction(3)
	response.LastLoginIP = simpleInfo.LoginIP
	response.RetCode = msg.RE_SUCCESS
}

//! 剔除作弊玩家
func Hand_KickArenaRanker(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_KickCheatRanker_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_KickArenaRanker Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_KickCheatRanker_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Hand_KickArenaRanker Error Invalid Gm request!!!")
		return
	}

	//! 获取排名玩家信息
	index := -1
	for i := 0; i < len(G_Rank_List); i++ {
		if G_Rank_List[i].PlayerID == req.PlayerID {
			index = i
			break
		}
	}

	if index <= 0 {
		gamelog.Error("Hand_KickArenaRanker Error: Player not exist.")
		response.RetCode = msg.RE_ACCOUNT_NOT_EXIST
		return
	}

	var tempLst3 []TArenaRankInfo
	tempLst := G_Rank_List[index+1:]
	tempLst2 := G_Rank_List[:index]
	tempLst3 = append(tempLst3, tempLst2...)
	tempLst3 = append(tempLst3, tempLst...)

	robot := gamedata.RandRobot(0)
	if robot == nil {
		gamelog.Error("Rand Robot Error: robot is nil")
		return
	}

	var rankerInfo TArenaRankInfo
	rankerInfo.IsRobot = true
	rankerInfo.PlayerID = robot.RobotID
	tempLst3 = append(tempLst3, rankerInfo)

	for i := 0; i < len(G_Rank_List); i++ {
		G_Rank_List[i].PlayerID = tempLst3[i].PlayerID
		G_Rank_List[i].IsRobot = tempLst3[i].IsRobot
	}

	player := GetPlayerByID(req.PlayerID)
	if player == nil {
		player = LoadPlayerFromDB(req.PlayerID)
	}

	player.ArenaModule.CurrentRank = 5001
	player.ArenaModule.HistoryRank = 5001
	player.ArenaModule.DB_UpdateRankToDatabase()

	response.RetCode = msg.RE_SUCCESS
}

//! 修改活动表数据
func Hand_UpdateActivityList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_UpdateActivityList_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_UpdateActivityList Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_UpdateActivityList_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	response.RetCode = msg.RE_SUCCESS

	var newActivity []string
	var updateActivity map[int32]string
	updateActivity = make(map[int32]string)

	for _, v := range req.ActivityLst {
		if v.Change == 2 { //! 新加活动
			//! 检测存在
			if gamedata.GetActivityInfo(v.ID) != nil {
				gamelog.Error("Hand_UpdateActivityList Error: Activity ID aleady exist.%d", v.ID)
				response.RetCode = msg.RE_INVALID_PARAM
				return
			}

			//! 修改内存
			data := new(gamedata.ST_ActivityInfo)
			data.ID = v.ID
			data.Name = v.Name
			data.TimeType = v.TimeType
			data.CycleType = v.CycleType
			data.BeginTime = v.BeginTime
			data.EndTime = v.EndTime
			data.AwardTime = v.AwardTime
			data.ActType = v.Type
			data.AwardType = v.AwardType
			data.Status = v.Status
			data.Icon = v.Icon
			data.Inside = v.Inside
			data.Days = v.Days
			gamedata.GT_ActivityLst[v.ID] = data

			str := fmt.Sprintf("%d,%s,%s,%s,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d",
				data.ID, data.Name, v.Desc, v.Ad, v.CycleType, v.TimeType, v.BeginTime, v.EndTime,
				v.AwardTime, v.Type, v.AwardType, v.Status, v.Icon, v.Inside, v.Days)
			newActivity = append(newActivity, str)
		} else if v.Change == 1 { //! 修改活动
			//! 检测存在
			data := gamedata.GetActivityInfo(v.ID)
			if data == nil {
				gamelog.Error("Hand_UpdateActivityList Error: Activity ID not exist.")
				response.RetCode = msg.RE_INVALID_PARAM
				return
			}

			//! 修改内存中静态表数据
			data.Name = v.Name
			data.TimeType = v.TimeType
			data.CycleType = v.CycleType
			data.BeginTime = v.BeginTime
			data.EndTime = v.EndTime
			data.AwardTime = v.AwardTime
			data.ActType = v.Type
			data.AwardType = v.AwardType
			data.Status = v.Status
			data.Icon = v.Icon
			data.Inside = v.Inside
			data.Days = v.Days

			//! 修改内存中动态数据
			for i := 0; i < len(G_GlobalVariables.ActivityLst); i++ {
				if G_GlobalVariables.ActivityLst[i].ActivityID == v.ID {
					G_GlobalVariables.ActivityLst[i].Status = v.Status
					G_GlobalVariables.ActivityLst[i].award = v.AwardType
					G_GlobalVariables.DB_UpdateActivityStatus(i)
					break
				}
			}

			str := fmt.Sprintf("%d,%s,%s,%s,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d",
				data.ID, data.Name, v.Desc, v.Ad, v.CycleType, v.TimeType, v.BeginTime, v.EndTime,
				v.AwardTime, v.Type, v.AwardType, v.Status, v.Icon, v.Inside, v.Days)

			updateActivity[v.ID] = str
		}
	}

	csv, err := ioutil.ReadFile("csv/type_activity.csv")
	if err != nil {
		gamelog.Error("ReadFile Fail: %s", err.Error())
		return
	}

	os.Remove("csv/type_activity.csv")

	strLst := strings.Split(string(csv), "\r\n")

	for i, v := range strLst {
		paramLst := strings.Split(string(v), ",")
		if len(paramLst) < 1 {
			continue
		}

		id, _ := strconv.Atoi(paramLst[0])

		str, isExist := updateActivity[int32(id)]
		if isExist == true {
			strLst[i] = str
		}
	}

	file, err := os.OpenFile("csv/type_activity.csv", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		gamelog.Error("OpenFile Error: %s", err.Error())
		return
	}

	defer file.Close()

	write := bufio.NewWriter(file)
	for _, v := range strLst {
		write.WriteString(v + "\r\n")
	}

	//! 增加新活动
	for _, v := range newActivity {
		write.WriteString(v + "\r\n")
	}

	write.Flush()

	G_GlobalVariables.CheckActivityNew()

	G_GlobalVariables.UpdateActivity()

	response.RetCode = msg.RE_SUCCESS

}

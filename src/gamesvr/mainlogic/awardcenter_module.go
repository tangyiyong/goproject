package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//! 奖励信息
type TAwardData struct {
	ID       int                    //! ID
	TextType int                    //!
	Value    []string               //! 参数
	ItemLst  []gamedata.ST_ItemData //! 奖励内容
	Time     int64                  //! 发放奖励时间戳
}

//! 领奖中心模块
type TAwardCenterModule struct {
	PlayerID   int          `bson:"_id"`
	AwardLst   []TAwardData //! 奖励列表
	SvrAwardID int          //! 已领取的全服奖励

	ownplayer *TPlayer
}

//! 设置指针
func (self *TAwardCenterModule) SetPlayerPtr(playerid int, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//! 玩家创建角色
func (self *TAwardCenterModule) OnCreate(playerID int) {
	//! 初始化信息
	self.AwardLst = make([]TAwardData, 0)
	self.SvrAwardID = G_GlobalVariables.SvrAwardIncID
	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerAwardCenter", self)
}

//! 玩家销毁角色
func (self *TAwardCenterModule) OnDestroy(playerID int) {

}

//! 玩家进入游戏
func (self *TAwardCenterModule) OnPlayerOnline(playerID int) {

}

//! 玩家离线
func (self *TAwardCenterModule) OnPlayerOffline(playerID int) {

}

//! 预取玩家信息
func (self *TAwardCenterModule) OnPlayerLoad(playerid int, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerAwardCenter").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerAwardCenter Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}

	self.PlayerID = playerid
}

func (self *TAwardCenterModule) RedTip() bool {
	if len(self.AwardLst) > 0 {
		return true
	}

	return false
}

func SendAwardMail(playerID int, textType int, awardLst []gamedata.ST_ItemData, value []string) {
	var awardData TAwardData
	awardData.TextType = textType
	awardData.ItemLst = awardLst
	awardData.Time = time.Now().Unix()
	awardData.Value = value
	SendAwardToPlayer(playerID, &awardData)
}

func SendAwardToPlayer(playerid int, pAwardData *TAwardData) {
	if playerid <= 0 {
		gamelog.Error("SendAwardToPlayer Error :Invalid PlayerID: %d", playerid)
		return
	}

	pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(playerid)
	if pSimpleInfo == nil {
		gamelog.Error("GetSimpleInfoByID Error :Invalid PlayerID: %d", playerid)
		return
	}

	//如果玩家不在线，并且己经离线超过7天时间，则不发邮件
	if pSimpleInfo.isOnline == false && (time.Now().Unix()-pSimpleInfo.LogoffTime) > 604800 && pSimpleInfo.LogoffTime != 0 {
		return
	}

	pAwardData.ID = pSimpleInfo.AwardCenterID
	pSimpleInfo.AwardCenterID += 1
	G_SimpleMgr.DB_SetAwardCenterID(playerid, pSimpleInfo.AwardCenterID)

	pPlayer := GetPlayerByID(playerid)
	if pPlayer != nil {
		pPlayer.AwardCenterModule.AwardLst = append(pPlayer.AwardCenterModule.AwardLst, *pAwardData)
	}

	go DB_SaveAwardToPlayer(playerid, *pAwardData)
}

//! 删除奖励
func (self *TAwardCenterModule) RemoveAward(id int) {
	pos := 0
	for i, v := range self.AwardLst {
		if v.ID == id {
			pos = i
			break
		}
	}

	if pos == 0 {
		self.AwardLst = self.AwardLst[1:]
	} else if (pos + 1) == len(self.AwardLst) {
		self.AwardLst = self.AwardLst[:pos]
	} else {
		self.AwardLst = append(self.AwardLst[:pos], self.AwardLst[pos+1:]...)
	}

	self.RemoveDatabaseLst(id)
}

//! 获取奖励内容
func (self *TAwardCenterModule) GetAwardData(id int) *TAwardData {
	for i, v := range self.AwardLst {
		if v.ID == id {
			return &self.AwardLst[i]
		}
	}

	gamelog.Error("AwardCenterModule GetAwardData fail. ID: %d  Lst: %v", id, self.AwardLst)
	return nil
}

////! DB相关
//! 增加奖励项到数据库
func (self *TAwardCenterModule) AddToDatabaseLst(award TAwardData) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerAwardCenter", bson.M{"_id": self.PlayerID}, "awardlst", award)
}

func DB_SaveAwardToPlayer(playerid int, award TAwardData) {
	if playerid <= 0 {
		gamelog.Error3("DB_SaveAwardToPlayer error. Invalid PlayerID:%d", playerid)
		return
	}
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerAwardCenter", bson.M{"_id": playerid}, bson.M{"$push": bson.M{"awardlst": award}})
}

//! 删除奖励项到数据库
func (self *TAwardCenterModule) RemoveDatabaseLst(id int) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "PlayerAwardCenter", bson.M{"_id": self.PlayerID}, "awardlst", bson.M{"id": id})
}

func SendSvrAwardToPlayer(playerid int) {
	if playerid <= 0 {
		gamelog.Error("SendSvrAwardToPlayer Error :Invalid PlayerID: %d", playerid)
		return
	}

	pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(playerid)
	if pSimpleInfo == nil {
		gamelog.Error("GetSimpleInfoByID Error :Invalid PlayerID: %d", playerid)
		return
	}

	pPlayer := GetPlayerByID(playerid)
	if pPlayer != nil {
		if pPlayer.AwardCenterModule.SvrAwardID < G_GlobalVariables.SvrAwardIncID {
			for _, v := range G_GlobalVariables.SvrAwardList {
				if v.ID > pPlayer.AwardCenterModule.SvrAwardID {
					pPlayer.AwardCenterModule.SvrAwardID = v.ID
					SendAwardToPlayer(playerid, &v)
				}
			}
		}

		pPlayer.AwardCenterModule.db_UpdateSvrAwardID()
	}
}

// 更新玩家已领取的全服奖励
func (self *TAwardCenterModule) db_UpdateSvrAwardID() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerAwardCenter", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"svrawardid": self.SvrAwardID}})
}

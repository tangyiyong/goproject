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

//! 称号信息
type TitleInfo struct {
	TitleID int   //! 拥有称号ID
	EndTime int64 //! 结束时间
	Status  int   //! 0->未激活 1->已激活 2->已佩戴
}

//! 称号模块
type TTitleModule struct {
	PlayerID int `bson:"_id"`

	TitleLst    []TitleInfo //! 拥有称号ID
	EquiTitleID int

	ownplayer *TPlayer
}

func (self *TTitleModule) SetPlayerPtr(playerid int, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

func (self *TTitleModule) OnCreate(playerID int) {

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerTitle", self)
}

func (self *TTitleModule) OnDestroy(playerID int) {

}

func (self *TTitleModule) OnPlayerOnline(playerID int) {

}

//! 玩家离开游戏
func (self *TTitleModule) OnPlayerOffline(playerID int) {

}

//! 读取玩家
func (self *TTitleModule) OnPlayerLoad(playerid int, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerTitle").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("Title Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 增加一个称号
func (self *TTitleModule) AddTitle(titleID int) {
	titleInfo := gamedata.GetTitleInfo(titleID)
	if titleInfo == nil {
		gamelog.Error("GetTitleInfo Error: invalid titleID %d", titleID)
		return
	}

	var title TitleInfo
	title.TitleID = titleID
	title.EndTime = time.Now().Unix() + int64(titleInfo.Time)

	if titleInfo.Time == -1 {
		title.EndTime = -1
	}

	self.TitleLst = append(self.TitleLst, title)
	go self.DB_AddTitleInfo(title)
}

func (self *TTitleModule) RemoveTitle(index int) {
	go self.DB_RemoveTitleInfo(&self.TitleLst[index])

	//! 若为佩戴状态则先替换下称号
	if self.ownplayer.HeroMoudle.TitleID == self.TitleLst[index].TitleID {
		self.ownplayer.HeroMoudle.TitleID = 0
		go self.ownplayer.HeroMoudle.DB_SaveTitleInfo()
	}

	if index == 0 {
		self.TitleLst = self.TitleLst[1:]
	} else if (index + 1) == len(self.TitleLst) {
		self.TitleLst = self.TitleLst[:index]
	} else {
		self.TitleLst = append(self.TitleLst[:index], self.TitleLst[index+1:]...)
	}

}

//! 检测称号到期情况
func (self *TTitleModule) CheckTitleDeadLine() {
	now := time.Now().Unix()

	for i := 0; i < len(self.TitleLst); i++ {
		if now >= self.TitleLst[i].EndTime && self.TitleLst[i].EndTime >= 0 {
			self.RemoveTitle(i)
			i -= 1
		}
	}
}

//! 佩戴称号
func (self *TTitleModule) EquiTitle(titleID int) bool {
	//! 取下原先称号
	for i, v := range self.TitleLst {
		if v.Status == 2 {
			v.Status = 1
			go self.DB_UpdateTitleStatus(i, 1)
		}
	}

	//! 换上新称号
	isExist := false
	for i, v := range self.TitleLst {
		if v.TitleID == titleID {
			isExist = true
			self.TitleLst[i].Status = 2
			go self.DB_UpdateTitleStatus(i, 2)
		}
	}

	if isExist == false {
		return false
	}

	self.ownplayer.HeroMoudle.TitleID = titleID
	go self.ownplayer.HeroMoudle.DB_SaveTitleInfo()

	return true
}

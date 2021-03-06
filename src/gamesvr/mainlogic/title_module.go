package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

//! 称号信息
type TitleInfo struct {
	TitleID int   //! 拥有称号ID
	EndTime int32 //! 结束时间
	Status  int   //! 0->未激活 1->已激活 2->已佩戴
}

//! 称号模块
type TTitleModule struct {
	PlayerID int32 `bson:"_id"`

	TitleLst    []TitleInfo //! 拥有称号ID
	EquiTitleID int

	ownplayer *TPlayer
}

func (self *TTitleModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TTitleModule) OnCreate(playerid int32) {

	//! 插入数据库
	mongodb.InsertToDB("PlayerTitle", self)
}

func (self *TTitleModule) OnDestroy(playerid int32) {

}

func (self *TTitleModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TTitleModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *TTitleModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerTitle").Find(&bson.M{"_id": playerid}).One(self)
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
	title.EndTime = utility.GetCurTime() + int32(titleInfo.Time)

	if titleInfo.Time == -1 {
		title.EndTime = 0xFFFFFFF
	}

	self.TitleLst = append(self.TitleLst, title)
	self.DB_AddTitleInfo(title)
}

func (self *TTitleModule) RemoveTitle(index int) {
	self.DB_RemoveTitleInfo(&self.TitleLst[index])

	//! 若为佩戴状态则先替换下称号
	if self.ownplayer.HeroMoudle.TitleID == self.TitleLst[index].TitleID {
		self.ownplayer.HeroMoudle.TitleID = 0
		self.ownplayer.HeroMoudle.DB_SaveTitleInfo()
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
	now := utility.GetCurTime()

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
			self.DB_UpdateTitleStatus(i, 1)
		}
	}

	//! 换上新称号
	isExist := false
	for i, v := range self.TitleLst {
		if v.TitleID == titleID {
			isExist = true
			self.TitleLst[i].Status = 2
			self.DB_UpdateTitleStatus(i, 2)
		}
	}

	if isExist == false {
		return false
	}

	self.ownplayer.HeroMoudle.TitleID = titleID
	self.ownplayer.HeroMoudle.DB_SaveTitleInfo()

	return true
}

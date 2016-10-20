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

type TBuyInfo struct {
	ID       int
	Type     int
	BuyTimes int
}

type TGuildSkill struct {
	SkillID int
	Level   int
}

//! 公会模块
type TGuildModule struct {
	PlayerID       int32      `bson:"_id"`
	JiTian         int8       //! 祭天状态
	JiTianAwardLst IntLst     //! 祭天奖励领取
	HisContribute  int        //! 历史贡献
	ApplyGuildList Int32Lst   //! 申请帮派列表
	BuyItems       []TBuyInfo //! 商店购买信息
	ActBuyTimes    int        //! 行动力购买次数
	ActTimes       int        //! 军团副本战斗次数
	ActRcrTime     int32      //! 行动力恢复
	CopyAwardMark  Int32Lst   //! 章节通关奖励
	QuitTime       int32      //! 退出公会时间
	ResetDay       uint32     //! 重置天数
	ownplayer      *TPlayer
}

func (self *TGuildModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TGuildModule) OnCreate(playerid int32) {
	//! 初始化各类参数
	self.ActBuyTimes = 0
	self.JiTian = 0
	self.ActTimes = gamedata.GuildBattleInitTime
	self.ActRcrTime = 0
	self.ResetDay = utility.GetCurDay()
	//! 插入数据库
	mongodb.InsertToDB("PlayerGuild", self)
}

func (self *TGuildModule) OnDestroy(playerid int32) {

}

func (self *TGuildModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TGuildModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *TGuildModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerGuild").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerGuild Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 检测重置
func (self *TGuildModule) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TGuildModule) OnNewDay(newday uint32) {
	//! 重置参数
	self.ResetDay = newday
	self.JiTian = 0
	self.JiTianAwardLst = IntLst{}
	self.ActBuyTimes = 0
	self.ActTimes = gamedata.GuildBattleInitTime

	for i := 0; i < len(self.BuyItems); i++ {
		if self.BuyItems[i].Type == 1 {
			self.BuyItems[i].BuyTimes = 0
		}
	}

	self.DB_Reset()
}

//! 检测会长是否弃坑
func (self *TGuildModule) CheckGuildLeader() {
	//! 获取会长信息
	guild := GetGuildByID(self.ownplayer.pSimpleInfo.GuildID)
	boss := guild.GetGuildLeader()

	bossInfo := G_SimpleMgr.GetSimpleInfoByID(boss.PlayerID)
	if bossInfo.LogoffTime == 0 {
		return
	}

	if utility.GetCurTime()-bossInfo.LogoffTime < 14*24*60*60 {
		return
	}

	//! 弃坑后职位禅让
	poseLst := []int{Pose_Deputy, Pose_Old, Pose_Elite, Pose_Member}

	role := 0
	for _, v := range poseLst {
		num := guild.GetRoleNum(v)
		if num != 0 {
			role = v
			break
		}
	}

	if role == 0 {
		//! 公会空无一人,删除公会
		RemoveGuild(self.ownplayer.pSimpleInfo.GuildID)
		self.ActRcrTime = 0
		self.QuitTime = utility.GetCurTime()
		self.DB_ExitGuild()
		return
	}

	//! 获取该职业所有成员
	memberLst := []TMember{}
	for _, v := range guild.MemberList {
		if v.Role == role {
			memberLst = append(memberLst, v)
		}
	}

	playerId := memberLst[0].PlayerID
	contribute := memberLst[0].Contribute
	for _, v := range memberLst {
		if v.Contribute > contribute {
			playerId = v.PlayerID
			contribute = v.Contribute
		}
	}

	member := guild.GetGuildMember(playerId)
	member.Role = Pose_Boss
	guild.UpdateGuildMemeber(playerId, Pose_Boss, member.Contribute)

	//! 解除现会长身份
	boss.Role = Pose_Member
	guild.UpdateGuildMemeber(boss.PlayerID, Pose_Member, boss.Contribute)

}

//! 增加贡献
func (self *TGuildModule) AddContribution(contribution int) {
	self.HisContribute += contribution
	self.DB_UpdateHisContribution()
}

//! 刷新限时抢购商品信息
func (self *TGuildModule) RefreshFalshSale() {
	for i, v := range self.BuyItems {
		if v.Type == 3 {
			self.BuyItems[i].BuyTimes = 0
		}
	}

	self.DB_ResetBuyLst()
}

//! 行动力恢复
func (self *TGuildModule) RecoverAction() {
	if self.ownplayer.pSimpleInfo.GuildID == 0 {
		//! 玩家不存在公会
		return
	}

	if self.ActRcrTime == 0 {
		return
	}

	if self.ActTimes > 8 {
		self.ActRcrTime = 0
		return
	}

	action := (utility.GetCurTime() - self.ActRcrTime) / int32(gamedata.GuildBattleRecoverTime)
	self.ActRcrTime = self.ActRcrTime + action*int32(gamedata.GuildBattleRecoverTime)
	self.ActTimes += int(action)

	if self.ActTimes > 8 {
		self.ActTimes = gamedata.GuildBattleInitTime
		return
	}

	self.DB_UpdateBattleTimes()
}

//! 红点提示
func (self *TGuildModule) RedTip() bool {
	if self.ownplayer.pSimpleInfo.GuildID == 0 {
		return true
	}

	guild := GetGuildByID(self.ownplayer.pSimpleInfo.GuildID)
	if self.JiTian == 0 {
		return true
	}

	//! 祭天奖励
	for i := 0; i < len(gamedata.GT_GuildSacrificeAwardLst); i++ {
		if guild.SacrificeSchedule >= gamedata.GT_GuildSacrificeAwardLst[i].NeedSchedule {
			if self.JiTianAwardLst.IsExist(gamedata.GT_GuildSacrificeAwardLst[i].ID) < 0 {
				return true
			}
		}
	}

	//! 公会副本奖励
	isExist := false
	for _, v := range guild.CopyTreasure {
		if v.PlayerID == self.PlayerID {
			isExist = true
			break
		}
	}

	if isExist == false {
		return true
	}

	return false
}

package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TGuildShopInfo struct {
	ID       int
	Type     int
	BuyTimes int
}

type TGuildShopInfoLst []TGuildShopInfo

func (self *TGuildShopInfoLst) Get(id int) *TGuildShopInfo {
	for i, v := range *self {
		if v.ID == id {
			return &(*self)[i]
		}
	}

	return nil
}

func (self *TGuildShopInfoLst) Add(id int, itemType int, times int) int {
	item := TGuildShopInfo{id, itemType, times}
	(*self) = append((*self), item)
	return len((*self)) - 1
}

type TGuildSkill struct {
	SkillID int
	Level   int
}

//! 公会模块
type TGuildModule struct {
	PlayerID int32 `bson:"_id"`

	SacrificeStatus   int    //! 祭天状态
	SacrificeAwardLst IntLst //! 祭天奖励领取

	HistoryContribution int //! 历史贡献
	TodayContribution   int //! 今日贡献

	ApplyGuildList IntLst //! 申请帮派列表

	ShoppingLst TGuildShopInfoLst //! 商店购买信息

	ActionTimes       int   //! 军团副本行动力
	ActionRecoverTime int64 //! 行动力恢复

	CopyAwardMark IntLst //! 章节通关奖励

	ExitGuildTime int64 //! 退出公会时间

	SkillLst []TGuildSkill //! 工会技能信息

	ResetDay uint32 //! 重置天数

	ownplayer *TPlayer
}

func (self *TGuildModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TGuildModule) OnCreate(playerid int32) {
	//! 初始化各类参数
	self.SacrificeStatus = 0
	self.ActionTimes = 10
	hour := gamedata.GuildCopyBattleTimeBegin / 3600
	min := (gamedata.GuildCopyBattleTimeBegin - hour*3600) / 60
	sec := gamedata.GuildCopyBattleTimeBegin - hour*3600 - min*60

	beginTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, min, sec, 0, time.Now().Location())

	self.ActionRecoverTime = beginTime.Unix() + int64(gamedata.GuildActionRecoverTime)

	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerGuild", self)
}

func (self *TGuildModule) OnDestroy(playerid int32) {

}

func (self *TGuildModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TGuildModule) OnPlayerOffline(playerid int32) {

}

//! 获取公会列表
func (self *TGuildModule) GetGuildLst(index int) (guildLst []TGuild) {

	//! 排序公会
	SortGuild()

	//! 获取公会列表
	for i := index; i < index+5; i++ {
		if i >= len(G_Guild_Key_List) {
			break
		}

		guildLst = append(guildLst, *G_Guild_List[G_Guild_Key_List[i]])
	}

	return guildLst
}

//! 读取玩家
func (self *TGuildModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerGuild").Find(bson.M{"_id": playerid}).One(self)
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
	self.SacrificeStatus = 0
	self.TodayContribution = 0
	self.SacrificeAwardLst = IntLst{}

	for i := 0; i < len(self.ShoppingLst); i++ {
		if self.ShoppingLst[i].Type == 1 {
			self.ShoppingLst[i].BuyTimes = 0
		}
	}

	go self.DB_Reset()
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

	if time.Now().Unix()-bossInfo.LogoffTime < 14*24*60*60 {
		return
	}

	//! 弃坑后职位禅让
	poseLst := []int{Pose_Deputy, Pose_Old, Pose_Elite, Pose_Member}

	role := 0
	for _, v := range poseLst {
		num := guild.GetPoseNumber(v)
		if num != 0 {
			role = v
			break
		}
	}

	if role == 0 {
		//! 公会空无一人,删除公会
		RemoveGuild(self.ownplayer.pSimpleInfo.GuildID)
		self.ActionRecoverTime = 0
		self.ExitGuildTime = time.Now().Unix()
		go self.DB_ExitGuild()
		return
	}

	//! 获取该职业所有成员
	memberLst := []TMember{}
	for _, v := range guild.MemberList {
		if v.Pose == role {
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
	member.Pose = Pose_Boss
	go guild.UpdateGuildMemeber(playerId, Pose_Boss, member.Contribute)

	//! 解除现会长身份
	boss.Pose = Pose_Member
	go guild.UpdateGuildMemeber(boss.PlayerID, Pose_Member, boss.Contribute)

}

//! 增加贡献
func (self *TGuildModule) AddContribution(contribution int) {
	self.HistoryContribution += contribution
	self.TodayContribution += contribution
	go self.DB_AddGuildContribution()
}

//! 刷新限时抢购商品信息
func (self *TGuildModule) RefreshFalshSale() {
	for i, v := range self.ShoppingLst {
		if v.Type == 3 {
			self.ShoppingLst[i].BuyTimes = 0
		}
	}

	go self.DB_ResetBuyLst()
}

//! 行动力恢复
func (self *TGuildModule) RecoverAction() {
	now := time.Now().Unix()
	if now < self.ActionRecoverTime {
		//! 未到恢复时间
		return
	}

	//! 副本关闭时间
	hour := gamedata.GuildCopyBattleTimeEnd / 60 * 60
	min := (gamedata.GuildCopyBattleTimeEnd - hour*3600) / 60
	sec := gamedata.GuildCopyBattleTimeEnd - hour*3600 - min*60
	endTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, min, sec, 0, time.Now().Location())

	if now >= endTime.Unix() {
		//! 副本开始时间
		hour = gamedata.GuildCopyBattleTimeBegin / 60 * 60
		min = (gamedata.GuildCopyBattleTimeBegin - hour*3600) / 60
		sec = gamedata.GuildCopyBattleTimeBegin - hour*3600 - min*60

		beginTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour+1, min, sec, 0, time.Now().Location())
		beginTime.AddDate(0, 0, 1)

		self.ActionTimes = 10
		self.ActionRecoverTime = beginTime.Unix() + int64(gamedata.GuildActionRecoverTime)
		go self.DB_UpdateCopyAction()
		return
	}

	if self.ownplayer.pSimpleInfo.GuildID == 0 {
		//! 玩家不存在公会
		return
	}

	if self.ActionRecoverTime == 0 {
		return
	}

	action := 1
	interval := now - self.ActionRecoverTime
	action += int(interval / int64(gamedata.GuildActionRecoverTime))

	self.ActionRecoverTime = self.ActionRecoverTime + int64(action*gamedata.GuildActionRecoverTime)
	self.ActionTimes += action
	go self.DB_UpdateCopyAction()
}

//! 获取玩家当前公会技能等级
func (self *TGuildModule) GetPlayerGuildSKillLevel(skillID int) int {
	for _, v := range self.SkillLst {
		if v.SkillID == skillID {
			return v.Level
		}
	}

	return 0
}

//! 增加玩家公会技能等级
func (self *TGuildModule) AddPlayerGuildSkillLevel(skillID int) {
	isExist := false
	level := 0
	for i, v := range self.SkillLst {
		if v.SkillID == skillID {
			self.SkillLst[i].Level += 1
			level = self.SkillLst[i].Level
			isExist = true
			break
		}
	}

	if isExist == false {
		var skill TGuildSkill
		skill.Level = 1
		skill.SkillID = skillID
		self.SkillLst = append(self.SkillLst, skill)

		go self.DB_AddGuildSkillInfo(skill)
	} else {
		go self.DB_UpdateGuildSkillLevel(skillID, level)
	}

	moneyID, moneyNum := gamedata.GetGuildSkillNeedMoney(level+1, skillID)
	self.ownplayer.RoleMoudle.CostMoney(moneyID, moneyNum)
}

//! 红点提示
func (self *TGuildModule) RedTip() bool {
	//! 加入公会
	if self.ownplayer.pSimpleInfo.GuildID == 0 {
		return true
	}

	guild := GetGuildByID(self.ownplayer.pSimpleInfo.GuildID)

	//! 祭天
	if self.SacrificeStatus == 0 {
		return true
	}

	//! 祭天奖励
	for i := 0; i < len(gamedata.GT_GuildSacrificeAwardLst); i++ {
		if guild.SacrificeSchedule >= gamedata.GT_GuildSacrificeAwardLst[i].NeedSchedule {
			if self.SacrificeAwardLst.IsExist(gamedata.GT_GuildSacrificeAwardLst[i].ID) < 0 {
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

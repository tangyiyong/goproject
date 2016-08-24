package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"strconv"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	Text_GM_Mail              = 1  //! GM邮件
	Text_Arean_Def_SUCCESS    = 2  //! 竞技场防守成功
	Text_Arean_Def_Fail_Drop  = 3  //! 竞技场防守失败
	Text_Arean_Def_Fail       = 4  //! 不掉排名的防守失败
	Text_Guild_Change_Name    = 5  //! 公会改名
	Text_Friend_Change_Name   = 6  //! 好友改名
	Text_Recharge             = 7  //! 充值成功
	Text_ScoreRace_Ret        = 8  //! 积分赛战报邮件
	Text_Arean_Win            = 9  //! 竞技场排名
	Text_Rebel_Find           = 10 //! 发现
	Text_Rebel_Killed         = 11 //! 击杀
	Text_Rebel_Exploit        = 12 //! 战功排名
	Text_Rebel_Damage         = 13 //! 伤害排名
	Text_CompetitionRankAward = 14 //! 开服竞赛排名
	Text_FoodWar_Rank         = 15 //! 粮草战排名
	Text_Recommand_Camp       = 16 //! 推荐阵营奖励邮件
	Text_MonthFund            = 17 //! 月基金奖励邮件
	TextCampBatTodayKill      = 18 //阵营战今日击杀排名奖励邮件
	TextCampBatTodayDestroy   = 19 //阵营战今日团灭排名奖励邮件
	TextCampBatCampKill       = 20 //阵营战阵营击杀排名奖励邮件
	TextCampBatCampDestroy    = 21 //阵营战阵营团灭排名奖励邮件

	TextCampBatHorseLamp = 101 //阵营战连杀跑马灯

)

type TMailInfo struct {
	TextType   int
	MailTime   int64
	MailParams []string
}

type TScoreReport struct {
	Name       string //名字
	HeroID     int    //英雄ID
	FightValue int    //战力
	Time       int64  //时间点
	Score      int    //积分
	Attack     bool   //是攻击还是防守
}

//角色邮件基本数据表结构
type TMailMoudle struct {
	PlayerID  int32          `bson:"_id"`
	MailList  []TMailInfo    //邮件列表
	Reports   []TScoreReport //积分战报
	ownplayer *TPlayer       //父player指针
}

func (playermail *TMailMoudle) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	playermail.PlayerID = playerid
	playermail.ownplayer = pPlayer
}

func (playermail *TMailMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	playermail.PlayerID = playerid
	//创建数据库记录
	playermail.MailList = make([]TMailInfo, 0)
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerMail", playermail)
}

//玩家对象销毁
func (playermail *TMailMoudle) OnDestroy(playerid int32) {
	playermail = nil
}

//玩家进入游戏
func (playermail *TMailMoudle) OnPlayerOnline(playerid int32) {
	//
}

//OnPlayerOffline 玩家离开游戏
func (playermail *TMailMoudle) OnPlayerOffline(playerid int32) {
	//
}

//玩家离开游戏
func (playermail *TMailMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerMail").Find(bson.M{"_id": playerid}).One(playermail)
	if err != nil {
		gamelog.Error("PlayerMail Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	playermail.PlayerID = playerid
}

func (self *TMailMoudle) RedTip() bool {
	//! 邮件列表
	if len(self.MailList) > 0 {
		return true
	}

	return false
}

//发邮件给角色
func SendMailToPlayer(playerid int32, pMailInfo *TMailInfo) {
	pPlayer := GetPlayerByID(playerid)
	if pPlayer != nil {
		pPlayer.MailMoudle.MailList = append(pPlayer.MailMoudle.MailList, *pMailInfo)
		DB_SaveMailToPlayer(playerid, pMailInfo)
		return
	}

	pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(playerid)
	//如果玩家不在线，并且己经离线超过7天时间，则不发邮件
	if pSimpleInfo.isOnline == false && (time.Now().Unix()-pSimpleInfo.LogoffTime) > 604800 {
		return
	}

	DB_SaveMailToPlayer(playerid, pMailInfo)
}

//! 保存邮件到数据库
func DB_SaveMailToPlayer(playerid int32, pMailInfo *TMailInfo) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerMail", bson.M{"_id": playerid}, bson.M{"$push": bson.M{"maillist": *pMailInfo}})
}

//! 保存战报到数据库
func DB_SaveScoreResultToPlayer(playerid int32, pResult *TScoreReport) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerMail", bson.M{"_id": playerid}, bson.M{"$push": bson.M{"reports": *pResult}})
}

//! 清空邮件到数据库
func (self *TMailMoudle) DB_ClearAllMails() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerMail", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"maillist": self.MailList}})
}

//! 清空邮件到数据库
func (self *TMailMoudle) DB_ClearAllReports() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerMail", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"reports": self.Reports}})
}

//以下为各功能的发邮件方法
////////////////////////////////////////////////////////////
//1. 竞技场邮件
func SendArenaMail(playerid int32, targetname string, rank int, win int, isChangeRank bool) {
	var mail TMailInfo
	if win == 0 {
		mail.TextType = Text_Arean_Def_SUCCESS
		mail.MailTime = time.Now().Unix()
		mail.MailParams = make([]string, 1)
		mail.MailParams[0] = targetname
		SendMailToPlayer(playerid, &mail)
	} else {
		if isChangeRank == true {
			mail.TextType = Text_Arean_Def_Fail_Drop
			mail.MailTime = time.Now().Unix()
			mail.MailParams = make([]string, 2)
			mail.MailParams[0] = targetname
			mail.MailParams[1] = strconv.Itoa(rank)
			SendMailToPlayer(playerid, &mail)
		} else {
			mail.TextType = Text_Arean_Def_Fail
			mail.MailTime = time.Now().Unix()
			mail.MailParams = make([]string, 1)
			mail.MailParams[0] = targetname
			SendMailToPlayer(playerid, &mail)
		}
	}
}

//2. 充值邮件
func SendRechargeMail(playerid int32, money int) {
	var mail TMailInfo
	mail.TextType = Text_Recharge
	mail.MailTime = time.Now().Unix()
	mail.MailParams = make([]string, 1)
	mail.MailParams[0] = strconv.Itoa(money)
	SendMailToPlayer(playerid, &mail)
}

//3. 积分赛战报邮件
func SendScoreResultMail(playerid int32, name string, fight int, heroid int, attack bool, score int) {
	var result TScoreReport
	result.Name = name
	result.FightValue = fight
	result.HeroID = heroid
	result.Attack = attack
	result.Score = score
	result.Time = time.Now().Unix()
	pPlayer := GetPlayerByID(playerid)
	if pPlayer != nil {
		pPlayer.MailMoudle.Reports = append(pPlayer.MailMoudle.Reports, result)
	}

	DB_SaveScoreResultToPlayer(playerid, &result)
}

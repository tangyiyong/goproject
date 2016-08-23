package mainlogic

import (
	"gamelog"
	"gamesvr/gamedata"

	"appconfig"
	"mongodb"
	"sort"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2"
)

const (
	Pose_Boss   = 1 //会长
	Pose_Deputy = 2 //副会长
	Pose_Old    = 3 //元老
	Pose_Elite  = 4 //精英
	Pose_Member = 5 //会员
)

type TMember struct {
	PlayerID     int32 //! ID
	Pose         int   //! 位置
	Contribute   int   //! 军团贡献
	EnterTime    int64 //! 加入时间
	BattleTimes  int   //! 攻打军团副本次数
	BattleDamage int64 //! 攻打军团副本最高伤害
}

const (
	GuildEvent_Sacrifice      = 1
	GuildEvent_Sacrifice_Crit = 2
	GuildEvent_AddMember      = 3
	GuildEvent_ChangePose     = 4
	GuildEvent_ExpelMember    = 5
	GuildEvent_LevelUp        = 6
)

type GuildEvent struct {
	PlayerID   int32
	PlayerName string
	Type       int //! 祭天->类型
	Value      int //! 祭天->经验  升级->等级  职位->新职位
	Action     int //!
	Time       int64
}

type GuildCopyTreasure struct {
	CopyID     int
	Index      int
	AwardID    int
	PlayerID   int32
	PlayerName string
}

type GuildCopyTreasureLst []GuildCopyTreasure

func (self *GuildCopyTreasureLst) GetNum(ID int) int {
	num := 0
	for _, v := range *self {
		if v.AwardID == ID {
			num++
		}
	}
	return num
}

//! 通关章节记录
type PassAwardChapter struct {
	PassChapter int
	CopyID      int
	PassTime    int64
	PlayerName  string
}

//! 公会留言板
type TGuildMsgBoard struct {
	PlayerID int32
	Message  string
	Time     int64
}

type TGuildCopyLifeInfo struct {
	CopyID int
	Life   int64
}

type MemberLst []TMember

func (self MemberLst) Len() int {
	return len(self)
}

func (self MemberLst) Less(i int, j int) bool {
	if (self)[i].BattleDamage < (self)[j].BattleDamage {
		return true
	} else if (self)[i].BattleDamage == (self)[j].BattleDamage {
		return (self)[i].BattleTimes < (self)[j].BattleTimes
	}
	return false
}

func (self MemberLst) Swap(i int, j int) {
	temp := (self)[i]
	(self)[i] = (self)[j]
	(self)[j] = temp
}

//公会表结构
type TGuild struct {
	GuildID     int       `bson:"_id"` //! 军团ID
	Name        string    //! 军团名字
	bossName    string    //! 军团长姓名
	Icon        int       //! 军团Icon
	Notice      string    //! 军团公告
	Declaration string    //! 军团宣言
	Level       int       //! 军团等级
	CurExp      int       //! 军团经验
	MemberList  MemberLst //! 军团成员列表
	fightVaule  int       //! 公会战力值

	ApplyList []int32 //! 申请列表

	EventLst []GuildEvent //! 军团动态

	Sacrifice         int //! 工会祭天人数
	SacrificeSchedule int //! 公会祭天进度

	SkillLst []TGuildSkill //! 工会技能信息

	HistoryPassChapter int                   //! 公会副本历史通关
	PassChapter        int                   //! 当天公会挑战章节
	CampLife           [4]TGuildCopyLifeInfo //! 副本四大势力对应血量
	IsBack             bool                  //! 是否回退
	CopyTreasure       GuildCopyTreasureLst  //! 副本奖励
	AwardChapterLst    []PassAwardChapter    //! 通关章节时间记录

	MsgBoard []TGuildMsgBoard //! 公会留言板

	ResetDay uint32 //! 重置天数
}

type GuildMap map[int]*TGuild
type GuildKeyLst []int //! 用于排序的Key值Slice

var (
	G_Guild_List     GuildMap
	G_Guild_Key_List GuildKeyLst
	Guild_Map_Mutex  sync.Mutex
	G_CurGuildID     int
)

func (self *TGuild) GetCopyLifeInfo(copyID int) int64 {
	for _, v := range self.CampLife {
		if v.CopyID == copyID {
			return v.Life
		}
	}

	gamelog.Error("GetCopyLifeInfo Error: invalid copyid: %v", copyID)
	return 0
}

func (self *TGuild) SetCopyLife(copyID int, life int64) {
	for i, v := range self.CampLife {
		if v.CopyID == copyID {
			self.CampLife[i].Life = life
			break
		}
	}
}

func (self GuildKeyLst) Len() int {
	return len(self)
}

func (self GuildKeyLst) Less(i int, j int) bool {
	if G_Guild_List[(self)[i]].Level < G_Guild_List[(self)[j]].Level {
		return true
	} else if G_Guild_List[(self)[i]].Level == G_Guild_List[(self)[j]].Level {
		return G_Guild_List[(self)[i]].fightVaule < G_Guild_List[(self)[j]].fightVaule
	}
	return false
}

func (self GuildKeyLst) Swap(i int, j int) {
	self[i], self[j] = self[j], self[i]
}

//初始化工会管理器
func InitGuildMgr() bool {

	s := mongodb.GetDBSession()
	defer s.Close()

	guildLst := []TGuild{}

	err := s.DB(appconfig.GameDbName).C("Guild").Find(nil).Sort("+_id").All(&guildLst)
	if err != nil {
		if err == mgo.ErrNotFound {
			G_CurGuildID = 1
		} else {
			gamelog.Error("Init GuildMgr Failed Error : %s!!", err.Error())
			return false
		}
	}

	if len(guildLst) <= 0 {
		G_CurGuildID = 1
	} else {
		G_CurGuildID = guildLst[len(guildLst)-1].GuildID + 1
	}

	if G_Guild_List == nil {
		G_Guild_List = make(GuildMap)
	}

	//! 初始化公会会长名
	if len(guildLst) != 0 {

		for i, n := range guildLst {
			//! 初始化排行榜与公会会长姓名
			bossInfo := n.GetGuildLeader()
			player := G_SimpleMgr.GetSimpleInfoByID(bossInfo.PlayerID)
			if player == nil {
				player := GetPlayerByID(bossInfo.PlayerID)
				guildLst[i].bossName = player.RoleMoudle.Name
			} else {
				guildLst[i].bossName = player.Name
			}

			G_Guild_List[n.GuildID] = &guildLst[i]
			G_Guild_Key_List = append(G_Guild_Key_List, n.GuildID)
		}
	}

	return true
}

//创建一个新的工会
func CreateNewGuild(playerid int32, name string, icon int) *TGuild {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	var newGuild TGuild
	newGuild.GuildID = G_CurGuildID
	newGuild.MemberList = make([]TMember, 1)
	newGuild.MemberList[0].PlayerID = playerid
	newGuild.MemberList[0].Pose = Pose_Boss
	newGuild.MemberList[0].Contribute = 0
	newGuild.MemberList[0].EnterTime = time.Now().Unix()
	newGuild.Level = 1
	newGuild.CurExp = 0
	newGuild.Notice = ""
	newGuild.Declaration = ""
	newGuild.Name = name
	newGuild.Icon = icon
	newGuild.HistoryPassChapter = 1
	newGuild.PassChapter = 1
	newGuild.ResetDay = utility.GetCurDay()

	guildCopy := gamedata.GetGuildChapterInfo(newGuild.PassChapter)
	if guildCopy != nil {
		for i := 0; i < 4; i++ {
			newGuild.CampLife[i].CopyID = guildCopy.CopyID[i]
			newGuild.CampLife[i].Life = guildCopy.Life
		}
	}

	newGuild.IsBack = false

	G_Guild_List[G_CurGuildID] = &newGuild
	G_Guild_Key_List = append(G_Guild_Key_List, newGuild.GuildID)
	G_CurGuildID += 1

	G_GuildCopyRanker.SetRankItem(int32(newGuild.GuildID), newGuild.HistoryPassChapter)
	G_GuildLevelRanker.SetRankItem(int32(newGuild.GuildID), newGuild.Level)

	//! 插入数据库
	DB_CreateGuild(&newGuild)
	return &newGuild
}

//! 解散公会
func RemoveGuild(guildID int) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()
	delete(G_Guild_List, guildID)

	//! 删除对应Key值
	removePos := -1
	for i, v := range G_Guild_Key_List {
		if v == guildID {
			removePos = i
		}
	}

	if removePos == 0 {
		G_Guild_Key_List = G_Guild_Key_List[1:]
	} else if (removePos + 1) == len(G_Guild_Key_List) {
		G_Guild_Key_List = G_Guild_Key_List[:removePos]
	} else {
		G_Guild_Key_List = append(G_Guild_Key_List[:removePos], G_Guild_Key_List[removePos+1:]...)
	}

	DB_RemoveGuild(guildID)
}

//! 排序公会列表
func SortGuild() {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()
	//! 计算战力
	for i, _ := range G_Guild_List {
		G_Guild_List[i].fightVaule = 0
		for _, v := range G_Guild_List[i].MemberList {
			G_Guild_List[i].fightVaule += G_SimpleMgr.Get_FightValue(v.PlayerID)
		}
	}

	sort.Sort(sort.Reverse(&G_Guild_Key_List))
}

//! 排序公会输出
func (self *TGuild) SortDamage() {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	sort.Sort(sort.Reverse(&self.MemberList))
}

//获取一个工会
func GetGuildByID(guildid int) *TGuild {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	pGuild, ok := G_Guild_List[guildid]
	if pGuild == nil || !ok {
		gamelog.Error("GetGuildByID Error: have not guild's id is :%d", guildid)
		return nil
	}

	return pGuild
}

//! 获取一个公会
func GetGuildByName(name string) *TGuild {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	var pGuild *TGuild
	for i, v := range G_Guild_List {
		if v.Name == name {
			pGuild = G_Guild_List[i]
		}
	}

	return pGuild
}

//获取公会名
func GetGuildName(guildid int) string {
	if guildid == 0 {
		return ""
	}
	pGuild := GetGuildByID(guildid)
	if pGuild == nil {
		return ""
	}

	return pGuild.Name
}

//获取工会成员信息
func (pGuild *TGuild) GetGuildMember(playerid int32) *TMember {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	for i := 0; i < len(pGuild.MemberList); i++ {
		if pGuild.MemberList[i].PlayerID == playerid {
			return &pGuild.MemberList[i]
		}
	}

	return nil
}

//! 检测重置
func (self *TGuild) CheckReset() {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.ResetDay = utility.GetCurDay()
	self.Sacrifice = 0
	self.SacrificeSchedule = 0

	if self.IsBack == true {
		//! 回退章节
		if self.PassChapter != 1 {
			self.PassChapter = self.PassChapter - 1
		}
	}

	guildCopy := gamedata.GetGuildChapterInfo(self.PassChapter)

	for i := 0; i < len(self.CampLife); i++ {
		self.CampLife[i].CopyID = guildCopy.CopyID[i]
		self.CampLife[i].Life = guildCopy.Life
	}

	//! 清空奖励章节记录
	self.AwardChapterLst = []PassAwardChapter{}
	self.CopyTreasure = []GuildCopyTreasure{}

	for i, _ := range self.MemberList {
		self.MemberList[i].BattleDamage = 0
		self.MemberList[i].BattleTimes = 0
	}

	self.DB_Reset()
}

//! 获取公会可领取奖励ID
func (self *TGuild) GetAleadyRecvAwardIDLst(chapter int, camp int) map[int]int {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	awardLst := make(map[int]int)
	for _, n := range self.CopyTreasure {
		if n.CopyID == camp {
			awardLst[n.AwardID] += 1
		}
	}

	tresureTypeLst := gamedata.GetGuildChapterCampAwardInfo(chapter, camp)
	for _, v := range tresureTypeLst {
		award := gamedata.GetGuildCampAwardInfo(v)
		awardLst[v] = award.Limit - awardLst[v]
	}

	return awardLst
}

//获取工会团长信息
func (pGuild *TGuild) GetGuildLeader() *TMember {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	for i := 0; i < len(pGuild.MemberList); i++ {
		if pGuild.MemberList[i].Pose == Pose_Boss {
			return &pGuild.MemberList[i]
		}
	}

	gamelog.Error("GetGuildLeader Error : Guild Has No Leader!!")
	return nil
}

//添加工会成员信息
func (pGuild *TGuild) AddGuildMember(playerid int32) bool {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	var newMember TMember
	newMember.PlayerID = playerid
	newMember.Contribute = 0
	newMember.Pose = Pose_Member
	newMember.EnterTime = time.Now().Unix()

	pGuild.MemberList = append(pGuild.MemberList, newMember)

	//! 插入数据库
	go DB_GuildAddMember(pGuild.GuildID, &newMember)

	return true
}

//! 增加军团经验
func (self *TGuild) AddExp(exp int) {
	Guild_Map_Mutex.Lock()

	self.CurExp += exp
	go self.DB_UpdateGuildLevel()

	Guild_Map_Mutex.Unlock()

	//! 检查公会升级
	self.LevelUp()
}

//! 公会升级
func (self *TGuild) LevelUp() {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	//! 获取下一级所需经验
	guildData := gamedata.GetGuildBaseInfo(self.Level + 1)
	if guildData == nil {
		return
	}

	if self.CurExp < guildData.NeedExp {
		return
	}

	self.CurExp -= guildData.NeedExp
	self.Level += 1
	go self.DB_UpdateGuildLevel()
	G_GuildLevelRanker.SetRankItem(int32(self.GuildID), self.Level)
}

//! 增加军团祭天进度
func (self *TGuild) AddSacrifice(schedule int) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	//! 祭天人数加一
	self.Sacrifice += 1

	//! 祭天进度增加
	self.SacrificeSchedule += schedule

	self.DB_UpdateGuildSacrifice()
}

//! 获取公会技能等级
func (self *TGuild) GetGuildSkillLevel(skillID int) int {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	for _, v := range self.SkillLst {
		if v.SkillID == skillID {
			return v.Level
		}
	}

	return 0
}

//! 升级公会技能等级
func (self *TGuild) AddGuildSkillLevel(skillID int, costExp int) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

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

		self.DB_AddGuildSkillLimit(skill)
	} else {
		go self.DB_UpdateGuildSkillLimit(skillID, level)
	}

	self.CurExp -= costExp
	go self.DB_UpdateGuildLevel()

}

//删除工会成员信息
func (pGuild *TGuild) RemoveGuildMember(playerid int32) bool {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	pos := 0
	var removeMember TMember
	for i, v := range pGuild.MemberList {
		if v.PlayerID == playerid {
			pos = i
			removeMember = v
			break
		}
	}

	if removeMember.PlayerID == 0 {
		gamelog.Error("RemoveGuildMember Error: invalid playerID %v", playerid)
		return false
	}

	if pos == 0 {
		pGuild.MemberList = pGuild.MemberList[1:]
	} else if (pos + 1) == len(pGuild.MemberList) {
		pGuild.MemberList = pGuild.MemberList[:pos]
	} else {
		pGuild.MemberList = append(pGuild.MemberList[:pos], pGuild.MemberList[pos+1:]...)
	}

	//! 修改数据库
	go DB_GuildRemoveMember(pGuild.GuildID, &removeMember)

	return true
}

//! 获取公会角色数量
func (self *TGuild) GetPoseNumber(role int) int {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	number := 0
	for _, v := range self.MemberList {
		if v.Pose == role {
			number += 1
		}
	}

	return number
}

//更新工会成员贡献信息
func (pGuild *TGuild) UpdateGuildMemeber(playerid int32, pose int, contribute int) bool {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	for i := 0; i < len(pGuild.MemberList); i++ {
		if pGuild.MemberList[i].PlayerID == playerid {
			pGuild.MemberList[i].Pose = pose
			pGuild.MemberList[i].Contribute = contribute
			go DB_GuildUpdateMember(pGuild.GuildID, &pGuild.MemberList[i], i)
		}
	}

	return true
}

//! 增加申请列表
func (self *TGuild) AddApplyList(playerid int32) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()
	self.ApplyList = append(self.ApplyList, playerid)
	go DB_AddApplyList(self.GuildID, playerid)
}

//! 删除申请列表
func (self *TGuild) RemoveApplyList(playerid int32) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()
	pos := 0
	isFind := false
	for i, v := range self.ApplyList {
		if v == playerid {
			pos = i
			isFind = true
			break
		}
	}

	if isFind == false {
		gamelog.Error("RemoveApplyList Error: invalid playerid: %v", playerid)
		return
	}

	if pos == 0 {
		self.ApplyList = self.ApplyList[1:]
	} else if (pos + 1) == len(self.ApplyList) {
		self.ApplyList = self.ApplyList[:pos]
	} else {
		self.ApplyList = append(self.ApplyList[:pos], self.ApplyList[pos+1:]...)
	}

	go DB_RemoveApplyList(self.GuildID, playerid)
}

//! 增加军团动态
func (self *TGuild) AddGuildEvent(playerid int32, action int, value int, value2 int) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	//! 超过20条则删除第一条
	if len(self.EventLst) > 20 {
		go self.DB_RemoveGuildEvent(self.EventLst[0])
		self.EventLst = self.EventLst[1:]
	}

	//! 构建事件动态
	event := GuildEvent{}
	event.PlayerID = playerid
	event.Action = action
	event.Type = value2
	event.Value = value
	event.Time = time.Now().Unix()

	self.EventLst = append(self.EventLst, event)
	go self.DB_AddGuildEvent(event)
}

//! 扣除军团副本阵营血量
func (self *TGuild) SubCampLife(copyID int, damage int64, playerName string) (bool, bool) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	life := self.GetCopyLifeInfo(copyID)
	life -= damage

	self.SetCopyLife(copyID, life)

	isKilled := false

	if life <= 0 {
		life = 0

		//! 记录今日通关
		var passChapter PassAwardChapter
		passChapter.PassChapter = self.PassChapter
		passChapter.CopyID = copyID
		passChapter.PassTime = time.Now().Unix()
		passChapter.PlayerName = playerName
		self.AwardChapterLst = append(self.AwardChapterLst, passChapter)
		go self.DB_AddPassChapter(passChapter)

		isKilled = true
	}

	self.DB_CostCampLife(copyID, life)

	isVictory := true
	for _, v := range self.CampLife {
		if v.Life != 0 {
			isVictory = false
		}
	}

	return isVictory, isKilled
}

//! 进入下一章副本
func (self *TGuild) NextChapter() {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	self.HistoryPassChapter = self.PassChapter
	self.PassChapter += 1

	if self.PassChapter > gamedata.GetGuildChapterCount() {
		self.PassChapter -= 1
	}

	guildCopy := gamedata.GetGuildChapterInfo(self.PassChapter)

	for i := 0; i < 4; i++ {
		self.CampLife[i].CopyID = guildCopy.CopyID[i]
		self.CampLife[i].Life = guildCopy.Life
	}

	G_GuildCopyRanker.SetRankItem(int32(self.GuildID), self.PassChapter)

	go self.DB_UpdateChapter()
}

//! 记录玩家领取副本奖励
func (self *TGuild) PlayerRecvAward(playerid int32, playerName string, copyID int, index int, awardID int) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()
	var treasure GuildCopyTreasure
	treasure.CopyID = copyID
	treasure.AwardID = awardID
	treasure.Index = index
	treasure.PlayerID = playerid
	treasure.PlayerName = playerName
	self.CopyTreasure = append(self.CopyTreasure, treasure)

	go self.DB_AddRecvRecord(treasure)
}

//! 判断玩家是否领取该奖励
func (self *TGuild) IsRecvCampAward(playerid int32, copyID int, chapter int) bool {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	for _, v := range self.CopyTreasure {
		if v.PlayerID == playerid {
			award := gamedata.GetGuildCampAwardInfo(v.AwardID)
			if award.CopyID == copyID && award.Chapter == chapter {
				return true
			}
		}
	}
	return false
}

//! 新加留言板留言
func (self *TGuild) AddMsgBoard(playerid int32, message string) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	//! 保持公会留言板存在三十条留言
	if len(self.MsgBoard) > 30 {
		self.DB_RemoveGuildMsgBoard(self.MsgBoard[0])
		self.MsgBoard = self.MsgBoard[1:]
	}

	var msg TGuildMsgBoard
	msg.PlayerID = playerid
	msg.Message = message
	msg.Time = time.Now().Unix()

	self.MsgBoard = append(self.MsgBoard, msg)

	go self.DB_AddGuildMsgBoard(msg)
}

//! 删除留言板留言
func (self *TGuild) RemoveMsgBoard(playerid int32, time int64) {
	Guild_Map_Mutex.Lock()
	defer Guild_Map_Mutex.Unlock()

	removePos := -1
	removeMsg := TGuildMsgBoard{}
	for i, v := range self.MsgBoard {
		if v.PlayerID == playerid && v.Time == time {
			removePos = i
			removeMsg = v
			break
		}
	}

	if removePos < 0 {
		return
	}

	if removePos == 0 {
		self.MsgBoard = self.MsgBoard[1:]
	} else if (removePos + 1) == len(self.MsgBoard) {
		self.MsgBoard = self.MsgBoard[:removePos]
	} else {
		self.MsgBoard = append(self.MsgBoard[:removePos], self.MsgBoard[removePos+1:]...)
	}

	go self.DB_RemoveGuildMsgBoard(removeMsg)
}

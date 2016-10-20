package gamedata

import (
	"fmt"
	"gamelog"
	"strconv"
)

//! 公会基础表
type ST_GuildBase struct {
	Level          int //! 公会等级
	SacrificeTimes int //! 祭天次数
	MemberLimit    int //! 人数上限
	NeedExp        int //! 需要经验
}

var GT_GuildBaseLst []ST_GuildBase

func InitGuildParser(total int) bool {
	GT_GuildBaseLst = make([]ST_GuildBase, total+1)
	return true
}

func ParseGuildBaseRecord(rs *RecordSet) {
	level := CheckAtoi(rs.Values[0], 0)
	GT_GuildBaseLst[level].Level = level
	GT_GuildBaseLst[level].SacrificeTimes = rs.GetFieldInt("sacrifice_times")
	GT_GuildBaseLst[level].MemberLimit = rs.GetFieldInt("member_limit")
	GT_GuildBaseLst[level].NeedExp = rs.GetFieldInt("needexp")
}

func GetGuildBaseInfo(level int) *ST_GuildBase {
	if level > len(GT_GuildBaseLst)-1 {
		gamelog.Error("GetGuildBaseInfo Error: invalid level: %d", level)
		return nil
	}

	return &GT_GuildBaseLst[level]
}

func GetGuildLevelFromExp(exp int, oldlevel int) int {

	pGuildInfo := GetGuildBaseInfo(oldlevel)
	if pGuildInfo == nil {
		gamelog.Error("GetGuildLevelFromExp fail. Vip info is nil")
		return oldlevel
	}

	if pGuildInfo.NeedExp > exp {
		return oldlevel
	}

	i := oldlevel
	for {
		pGuildInfo = GetGuildBaseInfo(i)
		if pGuildInfo == nil {
			return oldlevel
		}

		if pGuildInfo.NeedExp > exp {
			return i
		}

		i += 1
	}

	gamelog.Error("GetGuildLevelFromExp fail.return 0")
	return 0
}

//! 公会角色表
const (
	Permission_Income       = iota //! 收人权限
	Permission_UpdateNotice        //! 改公告
	Permission_Kick                //! 踢人权限
	Permission_Research            //! 研究公会技能
	Permission_UpdateGuild         //! 升级工会
	Permission_ResetCopy           //! 重置副本章节
	Permission_Change              //! 修改职位
	Permission_Dissolution         //! 解散公会
	Permission_End
)

type ST_GuildRole struct {
	Role       int //! 身份
	Number     int
	Permission [8]int
}

var GT_GuildRoleLst []ST_GuildRole

func InitGuildRoleParser(total int) bool {
	GT_GuildRoleLst = make([]ST_GuildRole, total+1)
	return true
}

func ParseGuildRoleRecord(rs *RecordSet) {
	role := CheckAtoi(rs.Values[0], 0)

	GT_GuildRoleLst[role].Number = CheckAtoi(rs.Values[2], 2)
	for i := 0; i < Permission_End; i++ {
		GT_GuildRoleLst[role].Permission[i] = CheckAtoi(rs.Values[i+3], i+3)
	}
}

func GetMaxRoleNum(role int) int {
	if role >= Permission_End {
		gamelog.Error("GetMaxRoleNum Error: Invalid role %d", role)
		return 0
	}

	return GT_GuildRoleLst[role].Number
}

func HasPermission(role int, permission int) bool {
	roleInfo := GT_GuildRoleLst[role]

	if roleInfo.Permission[permission] == 1 {
		return true
	}
	return false
}

//! 公会祭天类型表
type ST_GuildSacrifice struct {
	ID           int //! 唯一ID
	Schedule     int //! 增加祭天进度
	Exp          int //! 增加军团经验
	MoneyID      int //! 增加货币ID
	MoneyNum     int //! 增加货币数量
	CostMoneyID  int //! 消耗货币ID
	CostMoneyNum int //! 消耗货币数量
}

var GT_GuildSacrificeLst []ST_GuildSacrifice

func InitGuildSacrificeParser(total int) bool {
	GT_GuildSacrificeLst = make([]ST_GuildSacrifice, total+1)
	return true
}

func ParseGuildSacrificeRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_GuildSacrificeLst[id].ID = id
	GT_GuildSacrificeLst[id].Schedule = rs.GetFieldInt("schedule")
	GT_GuildSacrificeLst[id].Exp = rs.GetFieldInt("exp")
	GT_GuildSacrificeLst[id].MoneyID = rs.GetFieldInt("moneyid")
	GT_GuildSacrificeLst[id].MoneyNum = rs.GetFieldInt("moneynum")
	GT_GuildSacrificeLst[id].CostMoneyID = rs.GetFieldInt("costmoneyid")
	GT_GuildSacrificeLst[id].CostMoneyNum = rs.GetFieldInt("costmoneynum")
}

func GetGuildSacrificeInfo(id int) *ST_GuildSacrifice {
	if id > len(GT_GuildSacrificeLst)-1 {
		gamelog.Error("GetGuildSacrificeInfo Error: invalid id: %d", id)
		return nil
	}

	return &GT_GuildSacrificeLst[id]
}

//! 公会祭天奖励表
type ST_GuildSacrificeAward struct {
	ID           int //! 唯一ID
	NeedSchedule int //! 需求进度
	Award        int //! 奖励数量
	Level        int //! 公会等级
}

var GT_GuildSacrificeAwardLst []ST_GuildSacrificeAward

func InitGuildSacrificeAwardParser(total int) bool {
	GT_GuildSacrificeAwardLst = make([]ST_GuildSacrificeAward, total+1)
	return true
}

func ParseGuildSacrificeAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_GuildSacrificeAwardLst[id].ID = id
	GT_GuildSacrificeAwardLst[id].NeedSchedule = rs.GetFieldInt("schedule")
	GT_GuildSacrificeAwardLst[id].Award = rs.GetFieldInt("award")
	GT_GuildSacrificeAwardLst[id].Level = rs.GetFieldInt("level")
}

func GetGuildSacrificeAwardInfo(id int) *ST_GuildSacrificeAward {
	if id > len(GT_GuildSacrificeAwardLst)-1 {
		gamelog.Error("GetGuildSacrificeAwardInfo Error: invalid id: %d", id)
		return nil
	}

	return &GT_GuildSacrificeAwardLst[id]
}

func GetGuildSacrificeAwardFromLevel(level int) []int {
	awardLst := []int{}
	for _, v := range GT_GuildSacrificeAwardLst {
		if v.Level == level {
			awardLst = append(awardLst, v.ID)
		}
	}

	return awardLst
}

//! 公会商店商品表
type ST_GuildStore struct {
	ID            int
	Type          int //! 1->道具 2->时装 3->奖励
	NeedLevel     int //! 需求公会级别
	Limit         int //! 限购次数
	CostMoneyID1  int //! 货币
	CostMoneyNum1 int
	CostMoneyID2  int
	CostMoneyNum2 int
	ItemID        int
	ItemNum       int
	Discount      int //! 折扣
}

var GT_GuildStoreLst []ST_GuildStore

func InitGuildStoreParser(total int) bool {
	GT_GuildStoreLst = make([]ST_GuildStore, total+1)
	return true
}

func ParseGuildStoreRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_GuildStoreLst[id].ID = id
	GT_GuildStoreLst[id].Type = rs.GetFieldInt("type")
	GT_GuildStoreLst[id].NeedLevel = rs.GetFieldInt("needlevel")
	GT_GuildStoreLst[id].Limit = rs.GetFieldInt("limit")
	GT_GuildStoreLst[id].CostMoneyID1 = rs.GetFieldInt("costmoneyid1")
	GT_GuildStoreLst[id].CostMoneyNum1 = rs.GetFieldInt("costmoneynum1")
	GT_GuildStoreLst[id].CostMoneyID2 = rs.GetFieldInt("costmoneyid2")
	GT_GuildStoreLst[id].CostMoneyNum2 = rs.GetFieldInt("costmoneynum2")
	GT_GuildStoreLst[id].ItemID = rs.GetFieldInt("itemid")
	GT_GuildStoreLst[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_GuildStoreLst[id].Discount = rs.GetFieldInt("discount")
}

func GetGuildItemInfo(id int) *ST_GuildStore {
	if id > len(GT_GuildStoreLst)-1 {
		gamelog.Error("GetGuildItemInfo Error: invalid id %v", id)
		return nil
	}

	return &GT_GuildStoreLst[id]
}

//! 公会副本表
type ST_GuildCopy struct {
	Chapter          int
	CopyID           [4]int
	Life             int64
	MoneyID          int //! 军团贡献ID
	Contribution_max int
	Contribution_min int
	Exp              int //! 军团经验
	Award            int //! 章节通关奖励
}

var GT_GuildCopyLst []ST_GuildCopy

func InitGuildCopyParser(total int) bool {
	GT_GuildCopyLst = make([]ST_GuildCopy, total+1)
	return true
}

func ParseGuildCopyRecord(rs *RecordSet) {
	chapter := CheckAtoi(rs.Values[0], 0)
	GT_GuildCopyLst[chapter].Chapter = chapter

	for i := 1; i <= 4; i++ {
		filedName := fmt.Sprintf("copyid%d", i)
		GT_GuildCopyLst[chapter].CopyID[i-1] = rs.GetFieldInt(filedName)
	}

	GT_GuildCopyLst[chapter].Life, _ = strconv.ParseInt(rs.Values[6], 10, 64)
	GT_GuildCopyLst[chapter].MoneyID = rs.GetFieldInt("moneyid")
	GT_GuildCopyLst[chapter].Contribution_max = rs.GetFieldInt("contribution_max")
	GT_GuildCopyLst[chapter].Contribution_min = rs.GetFieldInt("contribution_min")
	GT_GuildCopyLst[chapter].Exp = rs.GetFieldInt("exp")
	GT_GuildCopyLst[chapter].Award = rs.GetFieldInt("award")
}

func GetGuildChapterInfo(chapter int32) *ST_GuildCopy {
	if int(chapter) > len(GT_GuildCopyLst)-1 {
		gamelog.Error("GetGuildChapterInfo Error: invalid chapter %v", chapter)
		return nil
	}

	return &GT_GuildCopyLst[chapter]
}

func GetGuildChapterCount() int32 {
	return int32(len(GT_GuildCopyLst) - 1)
}

type ST_GuildCopy_Award struct {
	ID      int
	Chapter int32
	CopyID  int
	ItemID  int
	ItemNum int
	Limit   int
}

var GT_GuildCopyAward [][]ST_GuildCopy_Award

func InitGuildCopyAwardParser(total int) bool {
	GT_GuildCopyAward = make([][]ST_GuildCopy_Award, total/4+1)
	return true
}

func ParseGuildCopyAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	chapter := CheckAtoi(rs.Values[1], 1)
	var award ST_GuildCopy_Award
	award.ID = id
	award.Chapter = int32(chapter)
	award.CopyID = rs.GetFieldInt("copyid")
	award.ItemID = rs.GetFieldInt("itemid")
	award.ItemNum = rs.GetFieldInt("itemnum")
	award.Limit = rs.GetFieldInt("limit")
	GT_GuildCopyAward[chapter] = append(GT_GuildCopyAward[chapter], award)
}

//! 获取章节总奖励
func GetGuildChapterCampAwardInfo(chapter int32, copyID int) []int {
	idLst := []int{}
	for _, v := range GT_GuildCopyAward[chapter] {
		if copyID == v.CopyID {
			idLst = append(idLst, v.ID)
		}
	}

	return idLst
}

//! 根据ID获取奖励信息
func GetGuildCampAwardInfo(id int) *ST_GuildCopy_Award {
	for i, v := range GT_GuildCopyAward {
		for j, n := range v {
			if n.ID == id {
				return &GT_GuildCopyAward[i][j]
			}
		}
	}
	return nil
}

//! 获取随机阵营奖励
func RandGuildCampAward(chapter int32, copyID int, recvLst map[int]int) *ST_GuildCopy_Award {

	awardLst := []int{}

	for i, v := range recvLst {

		if v != 0 {
			awardLst = append(awardLst, i)
		}
	}

	if len(awardLst) == 0 {
		return nil
	}

	award := awardLst[r.Intn(len(awardLst))]

	return GetGuildCampAwardInfo(award)
}

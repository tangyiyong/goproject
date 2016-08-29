package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 创建新公会
func DB_CreateGuild(info *TGuild) {
	mongodb.InsertToDB(appconfig.GameDbName, "Guild", info)
}

//! 解散公会
func DB_RemoveGuild(guildID int32) {
	mongodb.RemoveFromDB(appconfig.GameDbName, "Guild", bson.M{"_id": guildID})
}

//! 添加工会成员
func DB_GuildAddMember(guildID int32, member *TMember) {
	mongodb.AddToArray(appconfig.GameDbName, "Guild", bson.M{"_id": guildID}, "memberlist", *member)
}

//! 删除工会成员
func DB_GuildRemoveMember(guildID int32, member *TMember) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "Guild", bson.M{"_id": guildID}, "memberlist", *member)
}

//! 修改工会成员信息
func DB_GuildUpdateMember(guildID int32, member *TMember, index int) {
	filedName := fmt.Sprintf("memberlist.%d", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": guildID}, bson.M{"$set": bson.M{
		filedName: *member}})
}

//! 增加帮派事件
func (self *TGuild) DB_AddGuildEvent(event GuildEvent) {
	mongodb.AddToArray(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, "eventlst", event)
}

//! 删除帮派事件
func (self *TGuild) DB_RemoveGuildEvent(event GuildEvent) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, "eventlst", event)
}

//! 增加申请名单
func DB_AddApplyList(guildID int32, playerid int32) {
	mongodb.AddToArray(appconfig.GameDbName, "Guild", bson.M{"_id": guildID}, "applylist", playerid)
}

//! 删除申请名单
func DB_RemoveApplyList(guildID int32, playerid int32) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "Guild", bson.M{"_id": guildID}, "applylist", playerid)
}

//! 增加申请帮派名单
func (self *TGuildModule) DB_AddApplyGuildList(guildID int32) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, "applyguildlist", guildID)
}

//! 删除申请帮派名单
func (self *TGuildModule) DB_RemoveApplyGuildList(guildID int32) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, "applyguildlist", guildID)
}

//! 更改祭天标记
func (self *TGuildModule) DB_UpdateSacrifice() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"sacrificestatus": self.SacrificeStatus}})
}

//! 清空玩家申请列表
func (self *TGuildModule) DB_CleanApplyList() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"applyguildlist":    []int{},
		"actionrecovertime": self.ActionRecoverTime}})
}

func (self *TGuildModule) DB_ResetApplyList() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"applyguildlist": []int{}}})
}

//! 退出帮派
func (self *TGuildModule) DB_ExitGuild() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"actionrecovertime": self.ActionRecoverTime,
		"exitguildtime":     self.ExitGuildTime}})
}

//! 隔天刷新
func (self *TGuildModule) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"sacrificestatus":   self.SacrificeStatus,
		"resetday":          self.ResetDay,
		"todaycontribution": self.TodayContribution,
		"sacrificeawardlst": self.SacrificeAwardLst,
		"shoppinglst":       self.ShoppingLst}})
}

func (self *TGuildModule) DB_ResetBuyLst() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"shoppinglst": self.ShoppingLst}})
}

func (self *TGuild) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{
		"sacrifice":         self.Sacrifice,
		"sacrificeschedule": self.SacrificeSchedule,
		"resetday":          self.ResetDay,
		"camplife":          self.CampLife,
		"passchapter":       self.PassChapter,
		"awardchapterlst":   self.AwardChapterLst,
		"copytreasure":      self.CopyTreasure,
		"memberlist":        self.MemberList}})
}

//! 更新公会贡献
func (self *TGuildModule) DB_AddGuildContribution() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"historycontribution": self.HistoryContribution,
		"todaycontribution":   self.TodayContribution}})
}

//! 更新公会等级
func (self *TGuild) DB_UpdateGuildLevel() {
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{
		"level":  self.Level,
		"curexp": self.CurExp}})
}

//! 增加军团祭天信息
func (self *TGuild) DB_UpdateGuildSacrifice() {
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{
		"sacrifice":         self.Sacrifice,
		"sacrificeschedule": self.SacrificeSchedule}})
}

//! 更新祭天奖励领取标记
func (self *TGuildModule) DB_AddSacrificeMark(awardID int) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, "sacrificeawardlst", awardID)
}

//! 增加购买次数
func (self *TGuildModule) DB_AddShoppingTimes(id int, times int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID, "shoppinglst.id": id}, bson.M{"$set": bson.M{
		"shoppinglst.$.buytimes": times}})
}

//! 增加购买信息
func (self *TGuildModule) DB_AddShoppingInfo(info TGuildShopInfo) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, "shoppinglst", info)
}

//! 更新刷新标记
func (self *TGuild) DB_UpdateRefreshMark(index int) {
	filedName := fmt.Sprintf("isrefresh.%d", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{filedName: true}})
}

//! 减少购买次数
func (self *TGuildModule) DB_SubFlashSaleTimes(id int, times int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID, "shoppinglst.id": id}, bson.M{"$set": bson.M{
		"flashsalelst.$.buytimes": times}})
}

//! 更新行动力
func (self *TGuildModule) DB_UpdateCopyAction() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"actiontimes":       self.ActionTimes,
		"actionrecovertime": self.ActionRecoverTime}})
}

//! 更新最高伤害与攻击次数
func (self *TGuild) DB_UpdateDamageAndTimes(playerid int32, battleTimes int, battleDamage int64) {
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID, "memberlist.playerid": playerid}, bson.M{"$set": bson.M{
		"memberlist.$.battletimes":  battleTimes,
		"memberlist.$.battledamage": battleDamage}})
}

//! 扣除公会副本阵营血量
func (self *TGuild) DB_CostCampLife(copyID int, life int64) {
	filedName := fmt.Sprintf("camplife.$.life")
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID, "camplife.copyid": copyID}, bson.M{"$set": bson.M{
		filedName: life}})
}

//! 增加通关章节记录
func (self *TGuild) DB_AddPassChapter(chapter PassAwardChapter) {
	mongodb.AddToArray(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, "awardchapterlst", chapter)
}

//! 下一章节
func (self *TGuild) DB_UpdateChapter() {
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{
		"passchapter":        self.PassChapter,
		"historypasschapter": self.HistoryPassChapter,
		"camplife":           self.CampLife}})

}

//! 增加领取记录
func (self *TGuild) DB_AddRecvRecord(treasure GuildCopyTreasure) {
	mongodb.AddToArray(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, "copytreasure", treasure)
}

//! 增加章节奖励领取记录
func (self *TGuildModule) DB_AddChapterAwardRecord(chapter int) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": self.PlayerID}, "copyawardmark", chapter)
}

//! 修改工会基础信息
func (self *TGuild) DB_UpdateGuildInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{
		"notice":      self.Notice,
		"declaration": self.Declaration,
		"icon":        self.Icon}})
}

//! 修改公会名字
func (self *TGuild) DB_UpdateGuildName() {
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{
		"name": self.Name}})
}

//! 将玩家踢出公会
func (self *TGuildModule) DB_KickPlayer(playerid int32) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerGuild", bson.M{"_id": playerid}, bson.M{"$set": bson.M{
		"role":              0,
		"guildid":           0,
		"actionrecovertime": 0}})
}

//! 移除公会留言板信息
func (self *TGuild) DB_RemoveGuildMsgBoard(msg TGuildMsgBoard) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, "msgboard", msg)
}

//! 增加公会留言板信息
func (self *TGuild) DB_AddGuildMsgBoard(msg TGuildMsgBoard) {
	mongodb.AddToArray(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, "msgboard", msg)
}

//! 公会副本回退状态
func (self *TGuild) DB_UpdateGuildBackStatus() {
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{
		"isback": self.IsBack}})
}

//! 修改公会技能等级
func (self *TGuild) DB_UpdateGuildSkillLimit(index int) {
	filedName := fmt.Sprintf("skilllst.%d", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "Guild", bson.M{"_id": self.GuildID}, bson.M{"$set": bson.M{
		filedName: self.SkillLst[index]}})
}

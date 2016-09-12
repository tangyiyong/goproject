package mainlogic

import (
	"fmt"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

type TFestivalTask struct {
	ID       int //! 任务ID
	TaskType int //! 任务类型
	Count    int //! 当前次数
	Need     int //! 需要次数
	Status   int //! 状态: 0->未完成 1->已完成 2->已领取
	Award    int //! 奖励
}

type TFestivalExchangeRecord struct {
	ID    int //! 兑换ID
	Times int //! 兑换次数
}

type TFestivalSale struct {
	ID       int
	BuyTimes int
}

//! 节日欢庆
type TActivityFestival struct {
	ActivityID int //! 活动ID

	TaskLst     []TFestivalTask           //! 任务链
	ExchangeLst []TFestivalExchangeRecord //! 兑换记录
	BuyLst      []TFestivalSale           //! 购买记录

	VersionCode    int32            //! 更新号
	ResetCode      int32            //! 迭代号
	activityModule *TActivityModule //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityFestival) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityFestival) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self

	self.TaskLst = []TFestivalTask{}
	self.ExchangeLst = []TFestivalExchangeRecord{}

	self.RefreshTask(false)
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityFestival) Refresh(versionCode int32) {
	//! 刷新兑换次数
	length := len(self.ExchangeLst)
	for i := 0; i < length; i++ {
		self.ExchangeLst[i].Times = 0
	}

	self.VersionCode = versionCode
	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityFestival) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	self.ExchangeLst = []TFestivalExchangeRecord{}
	self.TaskLst = []TFestivalTask{}
	self.DB_Reset()
}

func (self *TActivityFestival) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityFestival) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityFestival) RedTip() bool {
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	for _, v := range self.TaskLst {
		if v.Status == 1 {
			return true
		}
	}

	return false
}

func (self *TActivityFestival) GetTaskInfo(taskID int) (*TFestivalTask, int) {
	length := len(self.TaskLst)
	for i := 0; i < length; i++ {
		if self.TaskLst[i].ID == taskID {
			return &self.TaskLst[i], i
		}
	}

	return nil, -1
}

func (self *TActivityFestival) GetExchangeInfo(id int) (*TFestivalExchangeRecord, int) {
	length := len(self.ExchangeLst)
	for i := 0; i < length; i++ {
		if self.ExchangeLst[i].ID == id {
			return &self.ExchangeLst[i], i
		}
	}

	//! 不存在则创建
	var record TFestivalExchangeRecord
	record.ID = id
	record.Times = 0
	self.ExchangeLst = append(self.ExchangeLst, record)
	self.DB_AddNewExchangeRecord(record)
	return &self.ExchangeLst[length], length
}

func (self *TActivityFestival) RefreshTask(isSaveDB bool) {
	if len(self.TaskLst) != 0 {
		self.TaskLst = []TFestivalTask{}
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	taskLst := gamedata.GetFestivalTaskFromType(awardType)

	length := len(taskLst)
	for i := 0; i < length; i++ {
		var task TFestivalTask
		task.ID = taskLst[i].ID
		task.Count = 0
		task.Need = taskLst[i].Need
		task.Status = 0
		task.Award = taskLst[i].Award
		task.TaskType = taskLst[i].TaskType
		self.TaskLst = append(self.TaskLst, task)
	}

	if isSaveDB {
		self.DB_RefreshTask()
	}
}

func (self *TActivityFestival) GetFestivalSaleInfo(id int) *TFestivalSale {
	length := len(self.BuyLst)
	for i := 0; i < length; i++ {
		if self.BuyLst[i].ID == id {
			return &self.BuyLst[i]
		}
	}

	//! 若不存在则创建信息
	var saleInfo TFestivalSale
	saleInfo.ID = id
	saleInfo.BuyTimes = 0
	self.BuyLst = append(self.BuyLst, saleInfo)
	self.DB_AddBuyLst(saleInfo)
	return &self.BuyLst[length]
}

func (self *TActivityFestival) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"festival.activityid":  self.ActivityID,
		"festival.tasklst":     self.TaskLst,
		"festival.exchangelst": self.ExchangeLst,
		"festival.versioncode": self.VersionCode,
		"festival.resetcode":   self.ResetCode}})
}

func (self *TActivityFestival) DB_RefreshTask() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"festival.tasklst": self.TaskLst}})
}

func (self *TActivityFestival) DB_RefreshExchangeReocrd() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"festival.exchangelst": self.ExchangeLst}})
}

func (self *TActivityFestival) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"festival.exchangelst": self.ExchangeLst,
		"festival.versioncode": self.VersionCode}})
}

func (self *TActivityFestival) DB_UpdateTaskStatus(index int) {

	filedName := fmt.Sprintf("festival.tasklst.%d.status", index)
	filedName2 := fmt.Sprintf("festival.tasklst.%d.count", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		filedName:  self.TaskLst[index].Status,
		filedName2: self.TaskLst[index].Count}})
}

func (self *TActivityFestival) DB_AddNewExchangeRecord(record TFestivalExchangeRecord) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$push": bson.M{"festival.exchangelst": record}})
}

func (self *TActivityFestival) DB_UpdateExchangeTimes(index int, times int) {
	filedName := fmt.Sprintf("festival.exchangelst.%d.times", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		filedName: times}})
}

func (self *TActivityFestival) DB_AddBuyLst(saleInfo TFestivalSale) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$push": bson.M{
		"festival.buylst": saleInfo}})
}

func (self *TActivityFestival) DB_UpdateBuyTimes(saleInfo *TFestivalSale) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID, "festival.buylst.id": saleInfo.ID}, &bson.M{"$set": bson.M{
		"festival.buylst.$.buytimes": saleInfo.BuyTimes}})
}

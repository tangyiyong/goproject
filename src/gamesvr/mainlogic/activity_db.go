package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (self *TActivityModule) DB_UpdateLimitDailySchedule(activityIndex int, taskIndex int) {
	filedName := fmt.Sprintf("limitdaily.%d.tasklst.%d.count", activityIndex, taskIndex)
	filedName2 := fmt.Sprintf("limitdaily.%d.tasklst.%d.status", activityIndex, taskIndex)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName:  self.LimitDaily[activityIndex].TaskLst[taskIndex].Count,
		filedName2: self.LimitDaily[activityIndex].TaskLst[taskIndex].Status}})
}

func (self *TActivityModule) DB_UpdateTotalRecharge(index int, value int) {
	filedName := fmt.Sprintf("recharge.%d.rechargevalue", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: value}})
}

func (self *TActivityModule) DB_UpdateRechargeRecord(activityIndex int, index int, times int) {
	filedName := fmt.Sprintf("singlerecharge.%d.rechargerecord.%d.status", activityIndex, index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: times}})
}

func (self *TActivityModule) DB_UpdateSingelAward(activityIndex int, index int, times int) {
	filedName := fmt.Sprintf("singlerecharge.%d.singleawardlst.%d.times", activityIndex, index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: times}})
}

func (self *TActivityModule) DB_UpdateSingleRecharge(index int, info TSingleRechargeRecord) {
	filedName := fmt.Sprintf("singlerecharge.%d.rechargerecord", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{filedName: info}})
}

func (self *TActivityModule) DB_AddNewLoginActivity(activity TActivityLogin) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"login": activity}})
}

func (self *TActivityModule) DB_AddNewRechargeActivity(activity TActivityRecharge) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"recharge": activity}})
}

func (self *TActivityModule) DB_AddNewDiscountActivity(activity TActivityDiscountSale) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"discountsale": activity}})
}

func (self *TActivityModule) DB_AddNewSingleRechargeActivity(activity TActivitySingleRecharge) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"singlerecharge": activity}})
}

func (self *TActivityModule) DB_AddNewLimitDailyActivity(activity TActivityLimitDaily) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"limitdaily": activity}})
}

func (self *TActivityModule) DB_AddNewSevenDay(activity TActivitySevenDay) {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"sevenday": activity}})
}

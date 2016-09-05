package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"msg"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

//角色基本数据表结构
type THangUpMoudle struct {
	PlayerID    int32          `bson:"_id"` //玩家ID
	CurBossID   int            //当前的BossID
	StartTime   int64          //挂机开始时间
	GridNum     int            //格子数
	ExpItems    []int          //经验丹
	QuickTime   int            //快速战斗次数
	AddGridTime int            //增加格子的次数
	History     []msg.THisHang //挂机历史数据
	ResetDay    uint32         //重置时间
	ownplayer   *TPlayer       //父player指针
}

func (hang *THangUpMoudle) SetPlayerPtr(playerid int32, player *TPlayer) {
	hang.PlayerID = playerid
	hang.ownplayer = player
}

func (hang *THangUpMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	hang.PlayerID = playerid
	hang.StartTime = 0
	hang.ExpItems = make([]int, 0)
	hang.QuickTime = 0
	hang.History = make([]msg.THisHang, 0)
	hang.AddGridTime = 0
	hang.GridNum = gamedata.HangUpInitGridNum
	hang.ResetDay = utility.GetCurDay()

	//创建数据库记录
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerHang", hang)
}

//玩家对象销毁
func (hang *THangUpMoudle) OnDestroy(playerid int32) {
	hang = nil
}

//玩家进入游戏
func (hang *THangUpMoudle) OnPlayerOnline(playerid int32) {
}

//OnPlayerOffline 玩家离开游戏
func (hang *THangUpMoudle) OnPlayerOffline(playerid int32) {

}

//玩家离开游戏
func (hang *THangUpMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) bool {
	s := mongodb.GetDBSession()
	defer s.Close()
	var bRet = true
	err := s.DB(appconfig.GameDbName).C("PlayerHang").Find(bson.M{"_id": playerid}).One(hang)
	if err != nil {
		gamelog.Error("PlayerHang Load Error :%s， PlayerID: %d", err.Error(), playerid)
		bRet = false
	}

	if wg != nil {
		wg.Done()
	}
	hang.PlayerID = playerid
	return bRet
}

func (hang *THangUpMoudle) CheckReset() {
	curDay := utility.GetCurDay()
	if curDay == hang.ResetDay {
		return
	}

	hang.OnNewDay(curDay)
}

func (hang *THangUpMoudle) OnNewDay(newday uint32) {
	hang.ResetDay = newday
	hang.QuickTime = 0
	hang.DB_SaveQuickFightTime()
}

func (hang *THangUpMoudle) CalcHangUpRatio(roleFight int32, bossFight int) int {
	var r = float64(roleFight)
	var b = float64(bossFight)
	ratio := r / b
	if ratio > 1 {
		ratio = 1
	} else if ratio < 0.4 {
		ratio = 0.4
	}
	return int(ratio * 10000)
}

//计算收获
func (hang *THangUpMoudle) ReceiveHangUpProduce() bool {
	if hang.CurBossID == 0 {
		gamelog.Error("UpdateHangUpState : Invalid BossID:%d", hang.CurBossID)
		return false
	}

	pHangUpInfo := gamedata.GetHangUpInfo(hang.CurBossID)
	if pHangUpInfo == nil {
		gamelog.Error("UpdateHangUpState : Invalid BossID2:%d", hang.CurBossID)
		return false
	}

	var produce bool = false
	for int(time.Now().Unix()-hang.StartTime) > pHangUpInfo.CDTime {
		hang.StartTime = hang.StartTime + int64(pHangUpInfo.CDTime)
		if utility.Rand() < hang.CalcHangUpRatio(hang.ownplayer.GetFightValue(), pHangUpInfo.FightValue) {
			for j := 0; j < pHangUpInfo.ProduceNum; j++ {
				if len(hang.ExpItems) < hang.GridNum {
					hang.ExpItems = append(hang.ExpItems, pHangUpInfo.ProduceID)
				}
			}
			hang.History = append(hang.History, msg.THisHang{hang.CurBossID, pHangUpInfo.ProduceID,
				pHangUpInfo.ProduceNum, hang.StartTime})
		} else {
			hang.History = append(hang.History, msg.THisHang{hang.CurBossID, pHangUpInfo.ProduceID,
				0, hang.StartTime})
		}

		produce = true
	}

	return produce
}

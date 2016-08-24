package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"msg"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

const (
	Task_Unfinished = 0 //! 未完成
	Task_Finished   = 1 //! 已完成
	Task_Received   = 2 //! 已领取
)

//! 任务分类
const (
	Task_Kind_Daily       = 1 //! 日常任务
	Task_Kind_Achievement = 2 //! 成就任务
	Task_Kind_SevenDay    = 3 //! 七天活动
	Task_Kind_Limit_Daily = 4 //! 限时日常
)

//角色任务表结构
type TTaskInfo struct {
	TaskID     int //! 任务ID
	TaskType   int //! 任务类型
	TaskStatus int //! 任务状态 0-> 未完成 1-> 已完成  2-> 已领取
	TaskCount  int //! 任务次数
}

//! 角色成就表结构
type TAchievementInfo struct {
	ID         int //! 成就ID
	Type       int //! 成就类型
	TaskStatus int //! 成就达成状态 0-> 未完成 1-> 已完成  2-> 已领取
	TaskCount  int //! 成就达成次数
}

//! 积分奖励表
type TTaskMoudle struct {
	PlayerID         int32       `bson:"_id"` //! 玩家ID
	TaskScore        int         //! 日常任务积分
	ScoreAwardStatus IntLst      //! 记录已积分宝箱ID
	ScoreAwardID     []int       //! 日常任务宝箱
	TaskList         []TTaskInfo //! 任务列表

	AchievementList []TAchievementInfo //! 成就列表
	AchievedList    []int              //! 已达成成就

	ResetDay  uint32   //! 更新时间戳
	ownplayer *TPlayer //父player指针
}

//! 刷新任务
func (taskmodule *TTaskMoudle) RefreshTask(update bool) {
	if len(taskmodule.TaskList) > 0 {
		taskmodule.TaskList = []TTaskInfo{}
	}

	if len(taskmodule.ScoreAwardID) > 0 {
		taskmodule.ScoreAwardID = []int{}
	}

	taskmodule.TaskScore = 0
	taskmodule.ScoreAwardStatus = IntLst{}

	//! 获取对应等级日常任务
	level := taskmodule.ownplayer.GetLevel()
	dailyTaskLst := gamedata.GetDailyTask(level)

	//! 添加日常任务
	for _, v := range dailyTaskLst {
		var task TTaskInfo
		task.TaskID = v.TaskID
		task.TaskType = v.Type
		task.TaskStatus = Task_Unfinished //! 状态初始未完成
		task.TaskCount = 0                //! 次数初始为零

		taskmodule.TaskList = append(taskmodule.TaskList, task)
	}

	//! 获取对应等级的日常任务积分奖励宝箱ID
	taskmodule.ScoreAwardID = gamedata.GetTaskScoreAwardID(level)

	if update == true {
		//! 更新到数据库
		go taskmodule.UpdateDailyTaskInfo()
	}
}

func (taskmoudle *TTaskMoudle) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	taskmoudle.PlayerID = playerid
	taskmoudle.ownplayer = pPlayer
}

func (taskmoudle *TTaskMoudle) OnCreate(playerid int32) {

	//创建数据库记录
	//初始化各个成员数值
	taskmoudle.PlayerID = playerid

	//! 刷新日常任务
	taskmoudle.RefreshTask(false)

	//! 创建成就任务
	achieveLst := gamedata.GetAchievementTask(taskmoudle.ownplayer.GetLevel())

	for _, v := range achieveLst {
		var info TAchievementInfo
		if v.TaskID == 0 {
			continue
		}
		info.ID = v.TaskID
		info.Type = v.Type
		info.TaskStatus = 0
		info.TaskCount = 0

		taskmoudle.AchievementList = append(taskmoudle.AchievementList, info)
	}

	taskmoudle.ResetDay = utility.GetCurDay()

	//创建数据库记录
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerTask", taskmoudle)
}

//! 检测重置时间
func (self *TTaskMoudle) CheckReset() {
	if utility.IsSameDay(self.ResetDay) {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TTaskMoudle) RedTip() bool {
	length := len(self.TaskList)

	for i := 0; i < length; i++ {
		if self.TaskList[i].TaskStatus == 1 {
			return true
		}
	}

	length = len(self.AchievementList)
	for i := 0; i < length; i++ {
		if self.AchievementList[i].TaskStatus == 1 {
			return true
		}
	}

	awardLst := gamedata.GetTaskScoreAwardID(self.ownplayer.GetLevel())
	length = len(awardLst)
	for i := 0; i < length; i++ {
		scoreAward := gamedata.GetTaskScoreAwardData(awardLst[i])
		if self.TaskScore >= scoreAward.NeedScore && self.ScoreAwardStatus.IsExist(scoreAward.TaskAwardID) == -1 {
			return true
		}
	}

	return false
}

func (self *TTaskMoudle) OnNewDay(newday uint32) {
	//! 刷新日常任务与重置时间
	self.RefreshTask(true)
	self.ResetDay = newday
	go self.UpdateResetTime()
}

//玩家对象销毁
func (taskmoudle *TTaskMoudle) OnDestroy(playerid int32) {

}

//OnPlayerOnline 玩家进入游戏
func (taskmoudle *TTaskMoudle) OnPlayerOnline(playerid int32) {
	//taskmoudle.LoadPlayer(playerid)
}

//OnPlayerOffline 玩家离开游戏
func (taskmoudle *TTaskMoudle) OnPlayerOffline(playerid int32) {

}

func (taskmoudle *TTaskMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerTask").Find(bson.M{"_id": playerid}).One(taskmoudle)
	if err != nil {
		gamelog.Error("PlayerTask Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	taskmoudle.PlayerID = playerid
}

//! 增加任务进度
func (taskmodule *TTaskMoudle) AddPlayerTaskSchedule(taskType int, count int) {

	//! 检测任务进度
	taskmodule.CheckReset()
	kind := gamedata.GetTaskSubType(taskType)

	for _, taskKind := range kind {
		if taskKind == Task_Kind_Achievement {
			//! 成就任务
			for i, v := range taskmodule.AchievementList {
				if v.Type == taskType && v.TaskStatus == Task_Unfinished {
					//! 判断次数是否上限
					info := gamedata.GetAchievementTaskInfo(taskmodule.AchievementList[i].ID)
					if info == nil {
						gamelog.Error("GetAchievementTask failed. taskID: %v", v.ID)
						break
					}

					//! 增加完成进度
					//! 特殊处理
					if taskType == gamedata.TASK_HERO_EQUI_STRENGTH ||
						taskType == gamedata.TASK_HERO_EQUI_QUALITY ||
						taskType == gamedata.TASK_ARENA_RANK ||
						taskType == gamedata.TASK_SGWS_RANK ||
						taskType == gamedata.TASK_HERO_EQUI_REFINED ||
						taskType == gamedata.TASK_HERO_EQUI_REFINED_MAX ||
						taskType == gamedata.TASK_HERO_DESTINY_LEVEL ||
						taskType == gamedata.TASK_HERO_DESTINY_LEVEL_MAX ||
						taskType == gamedata.TASK_ATTACK_REBEL_DAMAGE ||
						taskType == gamedata.TASK_SCORE_RANK ||
						taskType == gamedata.TASK_REBEL_EXPLOIT ||
						taskType == gamedata.TASK_PASS_EPIC_COPY ||
						taskType == gamedata.TASK_SGWS_STAR ||
						taskType == gamedata.TASK_HERO_GEM_REFINED ||
						taskType == gamedata.TASK_HERO_GEM_REFINED_MAX ||
						taskType == gamedata.TASK_FIGHT_VALUE ||
						taskType == gamedata.TASK_CUR_HERO_BREAK ||
						taskType == gamedata.TASK_DIAOWEN_QUALITY ||
						taskType == gamedata.TASK_HERO_WAKE ||
						taskType == gamedata.TASK_FASHION_COMPOSE ||
						taskType == gamedata.TASK_PASS_MAIN_COPY_CHAPTER ||
						taskType == gamedata.TASK_PASS_ELITE_COPY_CHAPTER ||
						taskType == gamedata.TASK_CAMP_BATTLE_KILL ||
						taskType == gamedata.TASK_GUILD_LEVEL ||
						taskType == gamedata.TASK_PET_QUALITY ||
						taskType == gamedata.TASK_PET_STAR ||
						taskType == gamedata.TASK_CAMP_BATTLE_GROUP_KILL ||
						taskType == gamedata.TASK_PET_LEVEL ||
						taskType == gamedata.TASK_HERO_QUALITY {

						if count > taskmodule.AchievementList[i].TaskCount {
							taskmodule.AchievementList[i].TaskCount = count
						}

					} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_4) && info.Count == count {
						taskmodule.AchievementList[i].TaskCount += 1
					} else {
						//! 增加完成进度
						taskmodule.AchievementList[i].TaskCount += count
					}

					//! 判断完成进度
					if taskmodule.AchievementList[i].TaskCount >= info.Count && taskmodule.AchievementList[i].TaskStatus != Task_Received {
						//! 任务完成
						taskmodule.AchievementList[i].TaskStatus = Task_Finished
					}

					//! 更新数据
					taskmodule.UpdatePlayerAchievement(taskmodule.AchievementList[i].ID,
						taskmodule.AchievementList[i].TaskCount,
						taskmodule.AchievementList[i].TaskStatus)
				}
			}

		} else if taskKind == Task_Kind_Daily {
			//! 日常任务
			for i, _ := range taskmodule.TaskList {
				if taskmodule.TaskList[i].TaskType == taskType && taskmodule.TaskList[i].TaskStatus == Task_Unfinished {
					//! 判断次数是否上限
					info := gamedata.GetTaskInfo(taskmodule.TaskList[i].TaskID)
					if info == nil {
						gamelog.Error("GetTaskInfo failed. taskID: %v", taskmodule.TaskList[i].TaskID)
						break
					}

					if taskmodule.TaskList[i].TaskCount >= info.Count {
						//! 已达次数上限
						continue
					}

					//! 增加完成进度
					//! 特殊处理
					if taskType == gamedata.TASK_HERO_EQUI_STRENGTH ||
						taskType == gamedata.TASK_HERO_EQUI_QUALITY ||
						taskType == gamedata.TASK_ARENA_RANK ||
						taskType == gamedata.TASK_SGWS_RANK ||
						taskType == gamedata.TASK_HERO_EQUI_REFINED ||
						taskType == gamedata.TASK_HERO_EQUI_REFINED_MAX ||
						taskType == gamedata.TASK_HERO_DESTINY_LEVEL ||
						taskType == gamedata.TASK_HERO_DESTINY_LEVEL_MAX ||
						taskType == gamedata.TASK_ATTACK_REBEL_DAMAGE ||
						taskType == gamedata.TASK_REBEL_EXPLOIT ||
						taskType == gamedata.TASK_PASS_EPIC_COPY ||
						taskType == gamedata.TASK_SGWS_STAR ||
						taskType == gamedata.TASK_HERO_GEM_REFINED ||
						taskType == gamedata.TASK_HERO_GEM_REFINED_MAX ||
						taskType == gamedata.TASK_FIGHT_VALUE ||
						taskType == gamedata.TASK_SCORE_RANK ||
						taskType == gamedata.TASK_CUR_HERO_BREAK ||
						taskType == gamedata.TASK_DIAOWEN_QUALITY ||
						taskType == gamedata.TASK_HERO_WAKE ||
						taskType == gamedata.TASK_PET_LEVEL ||
						taskType == gamedata.TASK_FASHION_COMPOSE ||
						taskType == gamedata.TASK_PASS_MAIN_COPY_CHAPTER ||
						taskType == gamedata.TASK_PASS_ELITE_COPY_CHAPTER ||
						taskType == gamedata.TASK_CAMP_BATTLE_KILL ||
						taskType == gamedata.TASK_PET_STAR ||
						taskType == gamedata.TASK_PET_QUALITY ||
						taskType == gamedata.TASK_GUILD_LEVEL ||
						taskType == gamedata.TASK_CAMP_BATTLE_GROUP_KILL ||
						taskType == gamedata.TASK_HERO_QUALITY {

						if count > taskmodule.TaskList[i].TaskCount {
							taskmodule.TaskList[i].TaskCount = count
						}

					} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_4) && info.Count == count {
						taskmodule.TaskList[i].TaskCount += 1
					} else {
						//! 增加完成进度
						taskmodule.TaskList[i].TaskCount += count
					}

					//! 判断完成进度
					if taskmodule.TaskList[i].TaskCount >= info.Count {
						//! 任务完成
						taskmodule.TaskList[i].TaskStatus = Task_Finished
					}

					//! 更新数据
					taskmodule.UpdatePlayerTask(taskmodule.TaskList[i].TaskID,
						taskmodule.TaskList[i].TaskCount,
						taskmodule.TaskList[i].TaskStatus)
				}
			}
		} else if taskKind == Task_Kind_SevenDay {

			for n, v := range taskmodule.ownplayer.ActivityModule.SevenDay {
				if G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
					//! 七日活动
					for i, _ := range taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList {
						if taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskType == taskType && taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskStatus == Task_Unfinished {
							//! 判断次数是否上限
							info := gamedata.GetSevenTaskInfo(taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskID)
							if info == nil {
								gamelog.Error("GetTaskInfo failed. taskID: %v", taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskID)
								break
							}

							if taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskCount >= info.Count {
								//! 已达次数上限
								continue
							}

							//! 特殊处理
							if taskType == gamedata.TASK_HERO_EQUI_STRENGTH ||
								taskType == gamedata.TASK_HERO_EQUI_QUALITY ||
								taskType == gamedata.TASK_ARENA_RANK ||
								taskType == gamedata.TASK_SGWS_RANK ||
								taskType == gamedata.TASK_HERO_EQUI_REFINED ||
								taskType == gamedata.TASK_HERO_EQUI_REFINED_MAX ||
								taskType == gamedata.TASK_HERO_DESTINY_LEVEL ||
								taskType == gamedata.TASK_HERO_DESTINY_LEVEL_MAX ||
								taskType == gamedata.TASK_ATTACK_REBEL_DAMAGE ||
								taskType == gamedata.TASK_PET_LEVEL ||
								taskType == gamedata.TASK_REBEL_EXPLOIT ||
								taskType == gamedata.TASK_PASS_EPIC_COPY ||
								taskType == gamedata.TASK_SGWS_STAR ||
								taskType == gamedata.TASK_HERO_GEM_REFINED ||
								taskType == gamedata.TASK_HERO_GEM_REFINED_MAX ||
								taskType == gamedata.TASK_FIGHT_VALUE ||
								taskType == gamedata.TASK_CUR_HERO_BREAK ||
								taskType == gamedata.TASK_DIAOWEN_QUALITY ||
								taskType == gamedata.TASK_SCORE_RANK ||
								taskType == gamedata.TASK_GUILD_LEVEL ||
								taskType == gamedata.TASK_HERO_WAKE ||
								taskType == gamedata.TASK_FASHION_COMPOSE ||
								taskType == gamedata.TASK_PASS_MAIN_COPY_CHAPTER ||
								taskType == gamedata.TASK_PASS_ELITE_COPY_CHAPTER ||
								taskType == gamedata.TASK_CAMP_BATTLE_KILL ||
								taskType == gamedata.TASK_PET_QUALITY ||
								taskType == gamedata.TASK_PET_STAR ||
								taskType == gamedata.TASK_CAMP_BATTLE_GROUP_KILL ||
								taskType == gamedata.TASK_HERO_QUALITY {

								if count > taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskCount {
									taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskCount = count
								}
							} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
								taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
								taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
								taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
								taskType == gamedata.TASK_CAMP_HERO_FULL_4) && info.Count == count {
								taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskCount += 1
							} else {
								//! 增加完成进度
								taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskCount += count
							}

							//! 判断完成进度
							if taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskCount >= info.Count {
								//! 任务完成
								taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskStatus = Task_Finished
							}

							//! 更新数据
							taskmodule.ownplayer.ActivityModule.SevenDay[n].DB_UpdatePlayerSevenTask(
								taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskID,
								taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskCount,
								taskmodule.ownplayer.ActivityModule.SevenDay[n].TaskList[i].TaskStatus)
						}
					}
				}

			}
		} else if taskKind == Task_Kind_Limit_Daily { //! 限时日常
			//! 遍历目前所有开启活动
			for i, v := range taskmodule.ownplayer.ActivityModule.LimitDaily {
				for j, m := range v.TaskLst {

					if m.TaskType == taskType {
						//! 特殊处理任务
						if taskType == gamedata.TASK_COMPLETE_ALL_TASK && m.Status == 0 {
							if taskmodule.ownplayer.ActivityModule.LimitDaily[i].IsAllComplete() {
								taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count = 1
								taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Status = 1
								go taskmodule.ownplayer.ActivityModule.DB_UpdateLimitDailySchedule(i, j)
							}
							continue
						}

						//! 遍历限时日常活动
						if taskType == gamedata.TASK_HERO_EQUI_STRENGTH ||
							taskType == gamedata.TASK_HERO_EQUI_QUALITY ||
							taskType == gamedata.TASK_ARENA_RANK ||
							taskType == gamedata.TASK_SGWS_RANK ||
							taskType == gamedata.TASK_HERO_EQUI_REFINED ||
							taskType == gamedata.TASK_HERO_EQUI_REFINED_MAX ||
							taskType == gamedata.TASK_HERO_DESTINY_LEVEL ||
							taskType == gamedata.TASK_HERO_DESTINY_LEVEL_MAX ||
							taskType == gamedata.TASK_ATTACK_REBEL_DAMAGE ||
							taskType == gamedata.TASK_REBEL_EXPLOIT ||
							taskType == gamedata.TASK_PASS_EPIC_COPY ||
							taskType == gamedata.TASK_SGWS_STAR ||
							taskType == gamedata.TASK_HERO_GEM_REFINED ||
							taskType == gamedata.TASK_HERO_GEM_REFINED_MAX ||
							taskType == gamedata.TASK_GUILD_LEVEL ||
							taskType == gamedata.TASK_PET_LEVEL ||
							taskType == gamedata.TASK_HERO_QUALITY ||
							taskType == gamedata.TASK_CUR_HERO_BREAK ||
							taskType == gamedata.TASK_DIAOWEN_QUALITY ||
							taskType == gamedata.TASK_HERO_WAKE ||
							taskType == gamedata.TASK_PASS_MAIN_COPY_CHAPTER ||
							taskType == gamedata.TASK_PASS_ELITE_COPY_CHAPTER ||
							taskType == gamedata.TASK_FASHION_COMPOSE ||
							taskType == gamedata.TASK_SCORE_RANK ||
							taskType == gamedata.TASK_PET_STAR ||
							taskType == gamedata.TASK_PET_QUALITY ||
							taskType == gamedata.TASK_CAMP_BATTLE_KILL ||
							taskType == gamedata.TASK_CAMP_BATTLE_GROUP_KILL ||
							taskType == gamedata.TASK_FIGHT_VALUE && m.Status == 0 {

							if count > taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count {
								taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count = count
							}

						} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
							taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
							taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
							taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
							taskType == gamedata.TASK_CAMP_HERO_FULL_4) && m.Need == count && m.Status == 0 {
							taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count += 1
						} else {
							taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count += count
						}

						if taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count >= m.Need {
							taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Status = Task_Finished
						}

						// gamelog.Error("LimitDailyTask taskType: %d  TaskCount: %d   Status: %d", taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].TaskType,
						// 	taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count,
						// 	taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Status)

						taskInfo := taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j]

						if taskInfo.Status == Task_Received {
							continue
						}

						go taskmodule.ownplayer.ActivityModule.DB_UpdateLimitDailySchedule(i, j)
					}

				}
			}

			for j, m := range taskmodule.ownplayer.ActivityModule.Festival.TaskLst {

				if m.TaskType == taskType && m.Status == 0 {
					taskmodule.ownplayer.ActivityModule.Festival.TaskLst[j].Count += count
				} else if taskType == gamedata.TASK_HERO_EQUI_STRENGTH ||
					taskType == gamedata.TASK_HERO_EQUI_QUALITY ||
					taskType == gamedata.TASK_ARENA_RANK ||
					taskType == gamedata.TASK_SGWS_RANK ||
					taskType == gamedata.TASK_HERO_EQUI_REFINED ||
					taskType == gamedata.TASK_HERO_EQUI_REFINED_MAX ||
					taskType == gamedata.TASK_HERO_DESTINY_LEVEL ||
					taskType == gamedata.TASK_HERO_DESTINY_LEVEL_MAX ||
					taskType == gamedata.TASK_ATTACK_REBEL_DAMAGE ||
					taskType == gamedata.TASK_REBEL_EXPLOIT ||
					taskType == gamedata.TASK_PASS_EPIC_COPY ||
					taskType == gamedata.TASK_SGWS_STAR ||
					taskType == gamedata.TASK_HERO_GEM_REFINED ||
					taskType == gamedata.TASK_HERO_GEM_REFINED_MAX ||
					taskType == gamedata.TASK_HERO_QUALITY ||
					taskType == gamedata.TASK_HERO_WAKE ||
					taskType == gamedata.TASK_FASHION_COMPOSE ||
					taskType == gamedata.TASK_PASS_MAIN_COPY_CHAPTER ||
					taskType == gamedata.TASK_PASS_ELITE_COPY_CHAPTER ||
					taskType == gamedata.TASK_GUILD_LEVEL ||
					taskType == gamedata.TASK_DIAOWEN_QUALITY ||
					taskType == gamedata.TASK_CAMP_BATTLE_KILL ||
					taskType == gamedata.TASK_PET_LEVEL ||
					taskType == gamedata.TASK_PET_STAR ||
					taskType == gamedata.TASK_PET_QUALITY ||
					taskType == gamedata.TASK_CAMP_BATTLE_GROUP_KILL ||
					taskType == gamedata.TASK_FIGHT_VALUE {

					if count > taskmodule.ownplayer.ActivityModule.Festival.TaskLst[j].Count {
						taskmodule.ownplayer.ActivityModule.Festival.TaskLst[j].Count = count
					}

				} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
					taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
					taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
					taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
					taskType == gamedata.TASK_CAMP_HERO_FULL_4) && m.Need == count {
					taskmodule.ownplayer.ActivityModule.Festival.TaskLst[j].Count += 1
				}

				if taskmodule.ownplayer.ActivityModule.Festival.TaskLst[j].Count >= m.Need {
					taskmodule.ownplayer.ActivityModule.Festival.TaskLst[j].Status = Task_Finished
				}

				taskInfo := taskmodule.ownplayer.ActivityModule.Festival.TaskLst[j]

				if taskInfo.Status == Task_Received {
					continue
				}

				go taskmodule.ownplayer.ActivityModule.Festival.DB_UpdateTaskStatus(j)
			}
		}
	}

}

//! 检查任务进度
func (taskmodule *TTaskMoudle) CheckPlayerTask(taskID int) (result bool, errcode int) {

	result = false
	errcode = msg.RE_TASK_NOT_COMPLETE

	//! 检测任务是否非法
	data := gamedata.GetTaskInfo(taskID)
	if data == nil {
		errcode = msg.RE_UNKNOWN_ERR
		return
	}

	//! 根据ID获取玩家进度
	for _, v := range taskmodule.TaskList {
		if v.TaskID == taskID {
			if v.TaskStatus == Task_Received {
				errcode = msg.RE_ALREADY_RECEIVED
				break
			}

			if v.TaskStatus == Task_Finished {
				result = true
				errcode = msg.RE_SUCCESS
				break
			}
		}
	}

	return result, errcode
}

//! 发放任务完成奖励
func (taskmodule *TTaskMoudle) ReceiveTaskAward(taskID int) (bool, []gamedata.ST_ItemData) {
	data := gamedata.GetTaskInfo(taskID)

	//! 获取任务积分
	taskmodule.TaskScore += data.Score
	taskmodule.UpdatePlayerTaskScore(taskmodule.TaskScore)

	//! 发放物品奖励
	awardInfo := gamedata.GetItemsFromAwardID(data.AwardItem)
	return taskmodule.ownplayer.BagMoudle.AddAwardItems(awardInfo), awardInfo
}

//! 发放成就完成奖励
func (taskmodule *TTaskMoudle) ReceiveAchievementAward(achievementID int) bool {
	//! 获取奖励信息
	data := gamedata.GetAchievementTaskInfo(achievementID)

	awardInfo := gamedata.GetItemsFromAwardID(data.AwardItem)
	return taskmodule.ownplayer.BagMoudle.AddAwardItems(awardInfo)
}

//! 发放任务积分奖励
func (taskmodule *TTaskMoudle) ReceiveTaskScoreAward(scoreAwardID int) (bool, []gamedata.ST_ItemData) {
	//! 获取奖励信息
	data := gamedata.GetTaskScoreAwardData(scoreAwardID)
	if data == nil {
		return false, []gamedata.ST_ItemData{}
	}

	//! 发放物品奖励
	awardInfo := gamedata.GetItemsFromAwardID(data.Award)
	return taskmodule.ownplayer.BagMoudle.AddAwardItems(awardInfo), awardInfo
}

//! 判断玩家积分
func (taskmodule *TTaskMoudle) CheckTaskScore(scoreAwardID int) (result bool, errcode int) {
	result = false
	errcode = msg.RE_UNKNOWN_ERR

	//! 判断奖励信息是否存在
	exist := false
	for _, v := range taskmodule.ScoreAwardID {
		if v == scoreAwardID {
			exist = true
			break
		}
	}

	if exist == false {
		errcode = msg.RE_INVALID_PARAM
		return result, errcode
	}

	//! 获取奖励信息
	data := gamedata.GetTaskScoreAwardData(scoreAwardID)
	if data == nil {
		return result, errcode
	}

	//! 判断奖励条件
	if taskmodule.TaskScore < data.NeedScore {
		errcode = msg.RE_TASK_SCORE_NOT_ENOUGH
		return result, errcode
	}

	//! 判断是否已领取
	if taskmodule.ScoreAwardStatus.IsExist(scoreAwardID) >= 0 {
		errcode = msg.RE_ALREADY_RECEIVED
		return result, errcode
	}

	result = true
	return result, errcode
}

//! 检查成就是否达成
func (taskmodule *TTaskMoudle) CheckAchievement(achievementID int) (result bool, errcode int) {
	//! 检查用户是否有条件达成该成就
	exist := false
	result = false
	var info *TAchievementInfo
	for i, v := range taskmodule.AchievementList {
		if v.ID == achievementID {
			info = &taskmodule.AchievementList[i]
			exist = true
		}
	}

	if exist == false {
		return result, msg.RE_INVALID_PARAM
	}

	//! 获取奖励内容
	awardData := gamedata.GetAchievementTaskInfo(achievementID)
	if awardData == nil {
		return result, msg.RE_INVALID_PARAM
	}

	//! 检查进度
	if info.TaskCount < awardData.Count {
		return result, msg.RE_TASK_NOT_COMPLETE
	}

	result = true
	return result, msg.RE_SUCCESS
}

//! 查询替换成就
func (taskmodule *TTaskMoudle) UpdateNextAchievement(achievementID int) *TAchievementInfo {
	data := gamedata.GetAchievementTaskFromFrontTask(achievementID)

	//! 创建新任务
	frontTaskID := 0
	var newTask TAchievementInfo
	for i, _ := range taskmodule.AchievementList {
		if taskmodule.AchievementList[i].ID == achievementID {
			if data == nil {
				//! 已经没有新的成就任务,返回完成
				return &taskmodule.AchievementList[i]
			}

			//! 赋值新的成就任务
			frontTaskID = taskmodule.AchievementList[i].ID
			taskmodule.AchievementList[i].ID = data.TaskID
			taskmodule.AchievementList[i].Type = data.Type

			//! 判断新成就达成
			if taskmodule.AchievementList[i].TaskCount >= data.Count {
				taskmodule.AchievementList[i].TaskStatus = Task_Finished
			} else {
				taskmodule.AchievementList[i].TaskStatus = Task_Unfinished
			}

			//! 替换数据库成就
			go taskmodule.UpdateAchievement(&taskmodule.AchievementList[i], frontTaskID)

			newTask.ID = taskmodule.AchievementList[i].ID
			newTask.Type = taskmodule.AchievementList[i].Type
			newTask.TaskCount = taskmodule.AchievementList[i].TaskCount
			newTask.TaskStatus = taskmodule.AchievementList[i].TaskStatus
		}
	}
	return &newTask
}

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
	ID     int //! 任务ID
	Type   int //! 任务类型
	Status int //! 任务状态 0-> 未完成 1-> 已完成  2-> 已领取
	Count  int //! 任务次数
}

//! 积分奖励表
type TTaskMoudle struct {
	PlayerID         int32       `bson:"_id"` //! 玩家ID
	TaskScore        int         //! 日常任务积分
	ScoreAwardStatus IntLst      //! 记录已积分宝箱ID
	ScoreAwardID     []int       //! 日常任务宝箱
	TaskList         []TTaskInfo //! 任务列表
	AchieveList      []TTaskInfo //! 成就列表
	AchieveIDs       []int       //! 已达成成就

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

	for i := 0; i < len(gamedata.GT_Task_List); i++ {
		if level >= gamedata.GT_Task_List[i].NeedMinLevel && level < gamedata.GT_Task_List[i].NeedMaxLevel {
			if gamedata.GT_Task_List[i].TaskID == 0 {
				continue
			}

			var task TTaskInfo
			task.ID = gamedata.GT_Task_List[i].TaskID
			task.Type = gamedata.GT_Task_List[i].Type
			task.Status = Task_Unfinished //! 状态初始未完成
			task.Count = 0                //! 次数初始为零

			taskmodule.TaskList = append(taskmodule.TaskList, task)
		}
	}

	//! 获取对应等级的日常任务积分奖励宝箱ID
	taskmodule.ScoreAwardID = gamedata.GetTaskScoreAwardID(level)

	if update == true {
		//! 更新到数据库
		taskmodule.DB_UpdateDailyTaskInfo()
	}
}

func (taskmoudle *TTaskMoudle) SetPlayerPtr(playerid int32, player *TPlayer) {
	taskmoudle.PlayerID = playerid
	taskmoudle.ownplayer = player
}

func (taskmoudle *TTaskMoudle) OnCreate(playerid int32) {

	//创建数据库记录
	//初始化各个成员数值
	taskmoudle.PlayerID = playerid

	//! 刷新日常任务
	taskmoudle.RefreshTask(false)

	//! 创建成就任务
	for i := 0; i < len(gamedata.GT_Achievement_Lst); i++ {
		if gamedata.GT_Achievement_Lst[i].TaskID == 0 {
			continue
		}

		if 1 >= gamedata.GT_Achievement_Lst[i].NeedLevel && gamedata.GT_Achievement_Lst[i].FrontID == 0 {
			var info TTaskInfo
			info.ID = gamedata.GT_Achievement_Lst[i].TaskID
			info.Type = gamedata.GT_Achievement_Lst[i].Type
			taskmoudle.AchieveList = append(taskmoudle.AchieveList, info)
		}
	}

	taskmoudle.ResetDay = utility.GetCurDay()

	//创建数据库记录
	mongodb.InsertToDB("PlayerTask", taskmoudle)
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
		if self.TaskList[i].Status == 1 {
			return true
		}
	}

	length = len(self.AchieveList)
	for i := 0; i < length; i++ {
		if self.AchieveList[i].Status == 1 {
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
	self.ResetDay = newday
	self.RefreshTask(true)
}

//玩家对象销毁
func (taskmoudle *TTaskMoudle) OnDestroy(playerid int32) {

}

//OnPlayerOnline 玩家进入游戏
func (taskmoudle *TTaskMoudle) OnPlayerOnline(playerid int32) {
}

//OnPlayerOffline 玩家离开游戏
func (taskmoudle *TTaskMoudle) OnPlayerOffline(playerid int32) {

}

func (taskmoudle *TTaskMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerTask").Find(&bson.M{"_id": playerid}).One(taskmoudle)
	if err != nil {
		gamelog.Error("PlayerTask Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	taskmoudle.PlayerID = playerid
}

//! 增加任务进度
func (self *TTaskMoudle) AddPlayerTaskSchedule(taskType int, count int) {
	self.CheckReset()
	KindLst := gamedata.GetTaskSubType(taskType)
	if len(KindLst) <= 0 {
		gamelog.Error("AddPlayerTaskSchedule failed. Invalid taskType: %d", taskType)
		return
	}

	for kind := 0; kind < len(KindLst); kind++ {
		if KindLst[kind] == Task_Kind_Achievement { //! 成就任务
			for i := 0; i < len(self.AchieveList); i++ {
				if self.AchieveList[i].Type == taskType && self.AchieveList[i].Status == Task_Unfinished {
					info := gamedata.GetAchievementInfo(self.AchieveList[i].ID)
					if info == nil {
						gamelog.Error("AddPlayerTaskSchedule failed. AchieID: %v", self.AchieveList[i].ID)
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
						if count > self.AchieveList[i].Count {
							self.AchieveList[i].Count = count
						}
					} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_4) && info.Count == count {
						self.AchieveList[i].Count += 1
					} else {
						//! 增加完成进度
						self.AchieveList[i].Count += count
					}

					//! 判断完成进度
					if self.AchieveList[i].Count >= info.Count {
						//! 任务完成
						self.AchieveList[i].Status = Task_Finished
					}

					//! 更新数据
					self.DB_UpdateAchieve(i)
				}
			}

		} else if KindLst[kind] == Task_Kind_Daily {
			//! 日常任务
			for i := 0; i < len(self.TaskList); i++ {
				if self.TaskList[i].Type == taskType && self.TaskList[i].Status == Task_Unfinished {
					//! 判断次数是否上限
					info := gamedata.GetTaskInfo(self.TaskList[i].ID)
					if info == nil {
						gamelog.Error("GetTaskInfo failed. taskID: %v", self.TaskList[i].ID)
						break
					}

					if self.TaskList[i].Count >= info.Count {
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
						if count > self.TaskList[i].Count {
							self.TaskList[i].Count = count
						}

					} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
						taskType == gamedata.TASK_CAMP_HERO_FULL_4) && info.Count == count {
						self.TaskList[i].Count += 1
					} else {
						//! 增加完成进度
						self.TaskList[i].Count += count
					}

					//! 判断完成进度
					if self.TaskList[i].Count >= info.Count {
						//! 任务完成
						self.TaskList[i].Status = Task_Finished
					}

					//! 更新数据
					self.DB_UpdateTask(i)
				}
			}
		} else if KindLst[kind] == Task_Kind_SevenDay {
			for n, v := range self.ownplayer.ActivityModule.SevenDay {
				if G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
					//! 七日活动
					for i, _ := range self.ownplayer.ActivityModule.SevenDay[n].TaskList {
						if self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Type == taskType && self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Status == Task_Unfinished {
							//! 判断次数是否上限
							info := gamedata.GetSevenTaskInfo(self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].ID)
							if info == nil {
								gamelog.Error("GetTaskInfo failed. taskID: %v", self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].ID)
								break
							}

							openDay := GetOpenServerDay()
							if info.OpenDay > openDay {
								continue //! 未开启活动不计算进度
							}

							if self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Count >= info.Count {
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

								if count > self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Count {
									self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Count = count
								}
							} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
								taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
								taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
								taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
								taskType == gamedata.TASK_CAMP_HERO_FULL_4) && info.Count == count {
								self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Count += 1
							} else {
								//! 增加完成进度
								self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Count += count
							}

							//! 判断完成进度
							if self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Count >= info.Count {
								//! 任务完成
								self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Status = Task_Finished
							}

							//! 更新数据
							self.ownplayer.ActivityModule.SevenDay[n].DB_UpdatePlayerSevenTask(
								self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].ID,
								self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Count,
								self.ownplayer.ActivityModule.SevenDay[n].TaskList[i].Status)
						}
					}
				}

			}
		} else if KindLst[kind] == Task_Kind_Limit_Daily { //! 限时日常
			//! 遍历目前所有开启活动
			for i, v := range self.ownplayer.ActivityModule.LimitDaily {
				for j, m := range v.TaskLst {

					if m.TaskType == taskType {
						//! 特殊处理任务
						if taskType == gamedata.TASK_COMPLETE_ALL_TASK && m.Status == 0 {
							if self.ownplayer.ActivityModule.LimitDaily[i].IsAllComplete() {
								self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count = 1
								self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Status = 1
								self.ownplayer.ActivityModule.DB_UpdateLimitDailySchedule(i, j)
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

							if count > self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count {
								self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count = count
							}

						} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
							taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
							taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
							taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
							taskType == gamedata.TASK_CAMP_HERO_FULL_4) && m.Need == count && m.Status == 0 {
							self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count += 1
						} else {
							self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count += count
						}

						if self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count >= m.Need {
							self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Status = Task_Finished
						}

						// gamelog.Error("LimitDailyTask taskType: %d  TaskCount: %d   Status: %d", taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].TaskType,
						// 	taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Count,
						// 	taskmodule.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j].Status)

						taskInfo := self.ownplayer.ActivityModule.LimitDaily[i].TaskLst[j]

						if taskInfo.Status == Task_Received {
							continue
						}

						self.ownplayer.ActivityModule.DB_UpdateLimitDailySchedule(i, j)
					}

				}
			}

			for j, m := range self.ownplayer.ActivityModule.Festival.TaskLst {

				if m.TaskType == taskType && m.Status == 0 {
					self.ownplayer.ActivityModule.Festival.TaskLst[j].Count += count
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

					if count > self.ownplayer.ActivityModule.Festival.TaskLst[j].Count {
						self.ownplayer.ActivityModule.Festival.TaskLst[j].Count = count
					}

				} else if (taskType == gamedata.TASK_SINGLE_RECHARGE ||
					taskType == gamedata.TASK_CAMP_HERO_FULL_1 ||
					taskType == gamedata.TASK_CAMP_HERO_FULL_2 ||
					taskType == gamedata.TASK_CAMP_HERO_FULL_3 ||
					taskType == gamedata.TASK_CAMP_HERO_FULL_4) && m.Need == count {
					self.ownplayer.ActivityModule.Festival.TaskLst[j].Count += 1
				}

				if self.ownplayer.ActivityModule.Festival.TaskLst[j].Count >= m.Need {
					self.ownplayer.ActivityModule.Festival.TaskLst[j].Status = Task_Finished
				}

				taskInfo := self.ownplayer.ActivityModule.Festival.TaskLst[j]

				if taskInfo.Status == Task_Received {
					continue
				}

				self.ownplayer.ActivityModule.Festival.DB_UpdateTaskStatus(j)
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
		if v.ID == taskID {
			if v.Status == Task_Received {
				errcode = msg.RE_ALREADY_RECEIVED
				break
			}

			if v.Status == Task_Finished {
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
	taskmodule.DB_UpdateTaskScore(taskmodule.TaskScore)

	//! 发放物品奖励
	awardInfo := gamedata.GetItemsFromAwardID(data.AwardItem)
	return taskmodule.ownplayer.BagMoudle.AddAwardItems(awardInfo), awardInfo
}

//! 发放成就完成奖励
func (taskmodule *TTaskMoudle) ReceiveAchievementAward(achievementID int) bool {
	//! 获取奖励信息
	data := gamedata.GetAchievementInfo(achievementID)

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
		errcode = msg.RE_SCORE_NOT_ENOUGH
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
	var info *TTaskInfo
	for i, v := range taskmodule.AchieveList {
		if v.ID == achievementID {
			info = &taskmodule.AchieveList[i]
			exist = true
		}
	}

	if exist == false {
		return result, msg.RE_INVALID_PARAM
	}

	//! 获取奖励内容
	awardData := gamedata.GetAchievementInfo(achievementID)
	if awardData == nil {
		return result, msg.RE_INVALID_PARAM
	}

	//! 检查进度
	if info.Count < awardData.Count {
		return result, msg.RE_TASK_NOT_COMPLETE
	}

	result = true
	return result, msg.RE_SUCCESS
}

//! 查询替换成就
func (taskmodule *TTaskMoudle) UpdateNextAchievement(achievementID int) *TTaskInfo {
	data := gamedata.GetNextAchievement(achievementID)

	//! 创建新任务
	for i := 0; i < len(taskmodule.AchieveList); i++ {
		if taskmodule.AchieveList[i].ID == achievementID {
			if data == nil {
				//! 已经没有新的成就任务,返回完成
				return &taskmodule.AchieveList[i]
			}
			//! 赋值新的成就任务
			taskmodule.AchieveList[i].ID = data.TaskID
			taskmodule.AchieveList[i].Type = data.Type

			//! 判断新成就达成
			if taskmodule.AchieveList[i].Count >= data.Count {
				taskmodule.AchieveList[i].Status = Task_Finished
			} else {
				taskmodule.AchieveList[i].Status = Task_Unfinished
			}

			//! 替换数据库成就
			taskmodule.DB_UpdateAchieve(i)
			return &taskmodule.AchieveList[i]
		}
	}
	return nil
}

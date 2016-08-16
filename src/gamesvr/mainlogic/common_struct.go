package mainlogic

import (
	"gamesvr/gamedata"
)

//装备/英雄/宝物的三种位置类型
const (
	POSTYPE_BATTLE = 1 //上阵
	POSTYPE_BACK   = 2 //援军中
	POSTYPE_BAG    = 3 //背包中
)

//背包里的物品数据
type TItemData struct {
	ItemID  int
	ItemNum int
}

//THeroInfo is a sigle hero info
//单个英雄数据结构
type THeroData struct {
	HeroID         int     //英雄ID
	Level          int     //等级
	CurExp         int     //当前的经验值
	Quality        int     //品质， 只对主角有作用
	BreakLevel     int     //突破
	Cultures       [5]int  //培养
	CulturesCost   int     //培养消耗道具数量
	DiaoWenQuality [6]int  //雕文品质
	DiaoWenPtys    [30]int //雕文属性
	DiaoWenBack    [30]int //雕文等待替换属性
	DestinyState   int     //天命状态
	DestinyTime    int     //天命时间
	WakeLevel      int     //当前的觉醒等级
	WakeItem       [4]int  //四个觉醒道具
	GodLevel       int     //化神等级
}

func (self *THeroData) Init(heroid int) {
	self.HeroID = heroid
	self.Level = 1
	self.BreakLevel = 0
	self.DestinyState = 0x01000000
	self.DestinyTime = 0
	self.CurExp = 0
	self.Quality = gamedata.GetHeroQuality(heroid)
	self.DiaoWenQuality = [6]int{2, 2, 3, 3, 4, 5}
	self.GodLevel = 0
}

func (self *THeroData) Clear() {
	self.HeroID = 0
	self.Level = 1
	self.BreakLevel = 0
	self.DestinyState = 0x01000000
	self.DestinyTime = 0
	self.CurExp = 0
	self.Quality = 0
	self.DiaoWenQuality = [6]int{2, 2, 3, 3, 4, 5}
	self.GodLevel = 0
}

type TEquipData struct {
	EquipID         int //装备ID
	StrengLevel     int //强化等级
	RefineLevel     int //精炼等级
	RefineExp       int //当前精炼经验
	Star            int //升星等级
	StarExp         int //星级经验
	StarLuck        int //升星幸运值
	StarMoneyCost   int //升星银币消耗
	StarYuanBaoCost int //升星元宝消耗
	StarPieceCost   int //升星装备碎片消耗
}

func (self *TEquipData) Init(equipid int) {
	self.EquipID = equipid
	self.StrengLevel = 1
	self.RefineLevel = 0
	self.RefineExp = 0
	self.Star = 0
	self.StarExp = 0
	self.StarLuck = 0
}

func (self *TEquipData) Clear() {
	self.EquipID = 0
	self.StrengLevel = 0
	self.RefineLevel = 0
	self.RefineExp = 0
	self.Star = 0
	self.StarExp = 0
	self.StarLuck = 0
}

type TGemData struct {
	GemID       int //宝物ID
	StrengLevel int //强化等级
	StrengExp   int //当前强化经验
	RefineLevel int //精炼等级
}

func (self *TGemData) Init(gemid int) {
	self.GemID = gemid
	self.StrengLevel = 1
	self.StrengExp = 0
	self.RefineLevel = 0
}

func (self *TGemData) Clear() {
	self.GemID = 0
	self.StrengLevel = 0
	self.StrengExp = 0
	self.RefineLevel = 0
}

type TPetData struct {
	PetID  int //宠物ID
	Exp    int //宠物当前经验
	Level  int //宠物等级
	Star   int //宠物星级
	God    int //神炼等级
	GodExp int //神炼经验
}

func (self *TPetData) Init(petid int) {
	self.PetID = petid
	self.Level = 1
	self.Exp = 0
	self.Star = 0
	self.God = 0
	self.GodExp = 0
}

func (self *TPetData) Clear() {
	self.PetID = 0
	self.Level = 1
	self.Exp = 0
	self.Star = 0
	self.God = 0
	self.GodExp = 0
}

//时装数据
type TFashionData struct {
	ID    int //时装ID
	Level int //时装等级
}

func (self *TFashionData) Init(id int) {
	self.ID = id
	self.Level = 0
}

func (self *TFashionData) Clear() {
	self.ID = 0
	self.Level = 0
}

type TStoreBuyData struct {
	ID    int //! 物品ID
	Times int //! 购买次数
}

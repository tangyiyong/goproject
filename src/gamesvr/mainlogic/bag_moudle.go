package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

type TItemBag struct {
	Items []TItemData
}

type THeroBag struct {
	Heros []THeroData
}

type TEquipBag struct {
	Equips []TEquipData
}

type TGemBag struct {
	Gems []TGemData
}

type TPetBag struct {
	Pets []TPetData
}

type TFashionBag struct {
	Fashions []TFashionData
}

type TBagMoudle struct {
	PlayerID        int32       `bson:"_id"` //玩家ID
	HeroBag         THeroBag    //英雄包
	HeroPieceBag    TItemBag    //英雄碎片包
	EquipBag        TEquipBag   //装备包
	EquipPieceBag   TItemBag    //装备碎片包
	GemBag          TGemBag     //宝物背包
	GemPieceBag     TItemBag    //宝物碎片包
	NormalItemBag   TItemBag    //道具背包
	WakeItemBag     TItemBag    //觉醒道具背包
	PetBag          TPetBag     //宠物背包
	PetPieceBag     TItemBag    //宠物碎片背包
	HeroSoulBag     TItemBag    //将灵背包
	FashionBag      TFashionBag //时装背包
	FashionPieceBag TItemBag    //时装碎片包
	ColHeros        []int16     //收集过英雄列表
	ColPets         []int16     //收集过的宠物

	//以下属性非数据库属性
	ownplayer *TPlayer //父player指针
}

func (self *TBagMoudle) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	if pPlayer == nil {
		gamelog.Error("TBagMoudle SetPlayerPtr pPlayer is nil")
		return
	}

	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

//响应玩家创建
func (self *TBagMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	self.PlayerID = playerid
	//创建数据库记录
	self.HeroBag.Heros = make([]THeroData, 0)
	self.HeroPieceBag.Items = make([]TItemData, 0)
	self.EquipBag.Equips = make([]TEquipData, 0)
	self.EquipPieceBag.Items = make([]TItemData, 0)
	self.GemBag.Gems = make([]TGemData, 0)
	self.GemPieceBag.Items = make([]TItemData, 0)
	self.NormalItemBag.Items = make([]TItemData, 0)
	self.WakeItemBag.Items = make([]TItemData, 0)
	self.PetBag.Pets = make([]TPetData, 0)
	self.PetPieceBag.Items = make([]TItemData, 0)
	self.HeroSoulBag.Items = make([]TItemData, 0)

	self.InitAddItem()

	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerBag", self)
}

//玩家新建角色时给玩家的道具
func (self *TBagMoudle) InitAddItem() {
	var hero THeroData
	hero.Init(421)
	self.HeroBag.Heros = append(self.HeroBag.Heros, hero)
}

//玩家对象销毁
func (self *TBagMoudle) OnDestroy(playerid int32) {
}

//玩家进入游戏
func (self *TBagMoudle) OnPlayerOnline(playerid int32) {
}

//玩家离开游戏
func (self *TBagMoudle) OnPlayerOffline(playerid int32) {
	return
}

//玩家数据从数据库加载
func (self *TBagMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()
	err := s.DB(appconfig.GameDbName).C("PlayerBag").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerBag Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
	return
}

//添加奖励物品列表(不用增加数据库操作)
func (self *TBagMoudle) AddAwardItems(awarditems []gamedata.ST_ItemData) bool {
	if awarditems == nil {
		gamelog.Error("AddAwardItems Error : awarditems is nil!")
		return false
	}

	for _, item := range awarditems {
		self.AddAwardItem(item.ItemID, item.ItemNum)
	}

	return true
}

//添加单个奖励物品(不用增加数据库操作)
func (self *TBagMoudle) AddAwardItem(itemid int, num int) bool {
	if self.ownplayer == nil {
		gamelog.Error("AddAwardItem Error : ownplayer is nil !")
		return false
	}

	pItemInfo := gamedata.GetItemInfo(itemid)
	if pItemInfo == nil {
		gamelog.Error("AddAwardItem Error: Invalid itemid :%d", itemid)
		return false
	}

	if num <= 0 {
		gamelog.Error("AddAwardItem Error: Invalid num :%d", num)
		return false
	}

	switch pItemInfo.Type {
	case gamedata.TYPE_MONEY:
		{
			self.ownplayer.RoleMoudle.AddMoney(pItemInfo.Data1, num)
		}
	case gamedata.TYPE_ACTION:
		{
			self.ownplayer.RoleMoudle.AddAction(pItemInfo.Data1, num)
		}
	case gamedata.TYPE_HERO:
		{
			if num == 1 {
				self.AddHeroByID(pItemInfo.Data1, pItemInfo.Data2)
			} else {
				self.AddHeros(pItemInfo.Data1, pItemInfo.Data2, num)
			}
		}
	case gamedata.TYPE_EQUIPMENT:
		{
			if num == 1 {
				self.AddEqiupByID(pItemInfo.Data1)
			} else {
				self.AddEqiups(pItemInfo.Data1, num)
			}
		}
	case gamedata.TYPE_GEM:
		{
			if num == 1 {
				self.AddGemByID(pItemInfo.Data1)
			} else {
				self.AddGems(pItemInfo.Data1, num)
			}
		}
	case gamedata.TYPE_HERO_PIECE:
		{
			self.AddHeroPiece(itemid, num)
		}
	case gamedata.TYPE_EQUIP_PIECE:
		{
			self.AddEqiupPiece(itemid, num)
		}
	case gamedata.TYPE_GEM_PIECE:
		{
			self.AddGemPiece(itemid, num)
		}
	case gamedata.TYPE_WAKE:
		{
			self.AddWakeItem(itemid, num)
		}
	case gamedata.TYPE_NORMAL:
		{
			self.AddNormalItem(itemid, num)
		}
	case gamedata.TYPE_PET:
		{
			if num == 1 {
				self.AddPetByID(pItemInfo.Data1)
			} else {
				self.AddPets(pItemInfo.Data1, num)
			}
		}
	case gamedata.TYPE_PET_PIECE:
		{
			self.AddPetPiece(itemid, num)
		}
	case gamedata.TYPE_HEROSOUL:
		{
			self.AddHeroSoul(itemid, num)
		}
	case gamedata.TYPE_FASHION:
		{
			if num == 1 {
				self.AddFashionByID(pItemInfo.Data1)
			} else {
				self.AddFashions(pItemInfo.Data1, num)
			}
		}
	case gamedata.TYPE_FASHION_PIECE:
		{
			self.AddFashionPiece(itemid, num)
		}
	default:
		{
			return false
		}
	}

	return true
}

//获取指定位置索取引的英雄
func (self *TBagMoudle) GetBagHeroByPos(pos int) *THeroData {
	if (pos < 0) || (pos >= len(self.HeroBag.Heros)) {
		gamelog.Error("GetHeroByPos Error Invalid Pos :%d", pos)
		return nil
	}

	return &self.HeroBag.Heros[pos]
}

func (self *TBagMoudle) SetBagHeroByPos(pos int, pHero *THeroData) bool {
	if (pos < 0) || (pos >= len(self.HeroBag.Heros)) {
		gamelog.Error("SetHeroByPos Error Invalid Pos :%d", pos)
		return false
	}

	self.HeroBag.Heros[pos] = *pHero
	self.ownplayer.DB_SaveHeroAt(POSTYPE_BAG, pos)
	return true
}

//获取背包英雄的总数
func (self *TBagMoudle) GetBagHeroCount() int {
	if self.HeroBag.Heros == nil {
		gamelog.Error("GetHeroCount nil pointer")
		return 0
	}

	return len(self.HeroBag.Heros)
}

//添加一个出生英雄到背包
func (self *TBagMoudle) AddHeroByID(heroid int, level int) bool {
	if heroid <= 0 {
		gamelog.Error("AddHeroByID Error : Invalid heroid :%d", heroid)
		return false
	}

	var hero THeroData
	hero.Init(heroid)
	if level > 1 {
		hero.Level = level
	}

	self.HeroBag.Heros = append(self.HeroBag.Heros, hero)
	pHeroInfo := gamedata.GetHeroInfo(heroid)
	if pHeroInfo != nil && pHeroInfo.Quality >= 2 {
		bCol := true
		for i := 0; i < len(self.ColHeros); i++ {
			if self.ColHeros[i] == int16(heroid) {
				bCol = false
				break
			}
		}

		if bCol == true {
			self.ColHeros = append(self.ColHeros, int16(heroid))
		}

		self.DB_AddHeroAtLast(bCol)
	} else {
		self.DB_AddHeroAtLast(false)
	}

	campHeroCountLst := IntLst{0, 0, 0, 0}
	for i := 0; i < len(self.ColHeros); i++ {
		pHeroInfo = gamedata.GetHeroInfo(int(self.ColHeros[i]))
		if pHeroInfo != nil {
			switch pHeroInfo.Camp {
			case 1:
				{
					campHeroCountLst[0]++
				}
			case 2:
				{
					campHeroCountLst[1]++
				}
			case 3:
				{
					campHeroCountLst[2]++
				}
			case 4:
				{
					campHeroCountLst[3]++
				}
			}
		}
	}

	fullCampHeroLst := gamedata.GetCampHeroCount()
	for i := 0; i < 4; i++ {
		if campHeroCountLst[i] == fullCampHeroLst[i] {
			self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_CAMP_HERO_FULL_1+i, 1)
		}
	}

	return true
}

func (self *TBagMoudle) IsHeroBagFull() bool {
	nCount := len(self.HeroBag.Heros)
	maxCount := gamedata.GetFuncVipValue(gamedata.FUNC_HERO_BAG_CAPACITY, self.ownplayer.GetVipLevel())
	if nCount >= maxCount {
		return true
	}

	return false
}

func (self *TBagMoudle) IsEquipBagFull() bool {
	nCount := len(self.EquipBag.Equips)
	maxCount := gamedata.GetFuncVipValue(gamedata.FUNC_EQUIP_BAG_CAPACITY, self.ownplayer.GetVipLevel())
	if nCount >= maxCount {
		return true
	}

	return false
}

func (self *TBagMoudle) IsGembagFull() bool {
	nCount := len(self.GemBag.Gems)
	maxCount := gamedata.GetFuncVipValue(gamedata.FUNC_GEM_BAG_CAPACITY, self.ownplayer.GetVipLevel())
	if nCount >= maxCount {
		return true
	}

	return false
}

//添加多个英雄到背包
func (self *TBagMoudle) AddHeros(heroid int, level int, num int) bool {
	if (num <= 0) || (heroid <= 0) {
		gamelog.Error("AddHeros Invalid heroid: %d num:%d", heroid, num)
		return false
	}

	var heros []THeroData = make([]THeroData, num)
	for i := 0; i < num; i++ {
		heros[i].Init(heroid)
		heros[i].Level = level
	}

	self.HeroBag.Heros = append(self.HeroBag.Heros, heros...)
	self.DB_AddHeroList(heros)

	return true
}

//删除指定位置的英雄
func (self *TBagMoudle) RemoveHeroAt(pos int) bool {
	if pos >= len(self.HeroBag.Heros) {
		gamelog.Error("RemoveHeroAt Error Pos :%d can greater than bagnum", pos)
		return false
	}

	if pos == 0 {
		self.HeroBag.Heros = self.HeroBag.Heros[1:]
	} else if (pos + 1) == len(self.HeroBag.Heros) {
		self.HeroBag.Heros = self.HeroBag.Heros[:pos]
	} else {
		self.HeroBag.Heros = append(self.HeroBag.Heros[:pos], self.HeroBag.Heros[pos+1:]...)
	}
	return true
}

//获取指定位置的装备
func (self *TBagMoudle) GetEqiupByPos(pos int) *TEquipData {
	if (pos < 0) || (pos >= len(self.EquipBag.Equips)) {
		gamelog.Error("GetEqiupByPos Error Invalid Pos :%d", pos)
		return nil
	}
	return &self.EquipBag.Equips[pos]
}

//获取指定位置宠物
func (self *TBagMoudle) GetPetByPos(pos int) *TPetData {
	if (pos < 0) || (pos >= len(self.PetBag.Pets)) {
		gamelog.Error("GetPetByPos Error Invalid Pos :%d", pos)
		return nil
	}
	return &self.PetBag.Pets[pos]
}

//获取背包装备的总数
func (self *TBagMoudle) GetEqiupCount() int {
	if self.EquipBag.Equips == nil {
		gamelog.Error("GetEqiupCount nil pointer")
		return 0
	}

	return len(self.EquipBag.Equips)
}

//添加一个装备到背包
func (self *TBagMoudle) AddEqiupByID(equipid int) bool {
	if equipid <= 0 {
		gamelog.Error("AddEqiupByID Error : Invalid equipid : %d", equipid)
		return false
	}

	var equip TEquipData
	equip.Init(equipid)
	self.EquipBag.Equips = append(self.EquipBag.Equips, equip)
	self.DB_AddEquipAtLast()
	return true
}

//添加一个装备到背包
func (self *TBagMoudle) AddEqiups(equipid int, num int) bool {
	if (num <= 0) || (equipid <= 0) {
		gamelog.Error("AddEqiup Error : Invalid equipid:%d, num:%d", equipid, num)
		return false
	}

	var equips []TEquipData = make([]TEquipData, num)
	for i := 0; i < num; i++ {
		equips[i].Init(equipid)
	}

	self.EquipBag.Equips = append(self.EquipBag.Equips, equips...)
	self.DB_AddEquipsList(equips)
	return true
}

//添加一个装备到背包
func (self *TBagMoudle) AddEqiupData(pEquipData *TEquipData) bool {
	if pEquipData.ID <= 0 {
		gamelog.Error("AddEqiupData Error : Invalid EquipID:%d", pEquipData.ID)
		return false
	}
	self.EquipBag.Equips = append(self.EquipBag.Equips, *pEquipData)
	self.DB_AddEquipAtLast()
	return true
}

//删除指定位置的装备
func (self *TBagMoudle) RemoveEquipAt(pos int) bool {
	if pos >= len(self.EquipBag.Equips) {
		gamelog.Error("RemoveEquipAt Error Pos :%d can greater than bagnum", pos)
		return false
	}

	if pos == 0 {
		self.EquipBag.Equips = self.EquipBag.Equips[1:]
	} else if (pos + 1) == len(self.EquipBag.Equips) {
		self.EquipBag.Equips = self.EquipBag.Equips[:pos]
	} else {
		self.EquipBag.Equips = append(self.EquipBag.Equips[:pos], self.EquipBag.Equips[pos+1:]...)
	}
	return true
}

//获取指定位置的宝物
func (self *TBagMoudle) GetGemByPos(pos int) *TGemData {
	if (pos < 0) || (pos >= len(self.GemBag.Gems)) {
		gamelog.Error("GetGemByPos Error Invalid Pos :%d", pos)
		return nil
	}

	return &self.GemBag.Gems[pos]
}

//获取宝物的总数
func (self *TBagMoudle) GetGemCount() int {
	if self.GemBag.Gems == nil {
		gamelog.Error("GetGemCount nil pointer")
		return 0
	}
	return len(self.GemBag.Gems)
}

//统计指定索引以外的，可用宝物的个数
func (self *TBagMoudle) GetGemCountExcept(pos int, gemid int) (count int) {
	count = 0
	for i := 0; i < len(self.GemBag.Gems); i++ {
		if self.GemBag.Gems[i].ID == gemid {
			if i == pos {
				continue
			} else {
				if self.GemBag.Gems[i].RefineLevel <= 0 {
					count += 1
				}
			}
		}
	}

	return count
}

//添加一个宝物到背包
func (self *TBagMoudle) AddGemByID(gemid int) bool {
	if gemid <= 0 {
		gamelog.Error("AddGemByID Invalid gemid :%d", gemid)
		return false
	}

	var gem TGemData
	gem.Init(gemid)

	self.GemBag.Gems = append(self.GemBag.Gems, gem)
	self.DB_AddGemAtLast()
	return true
}

//添加一个装备到背包
func (self *TBagMoudle) AddGemData(pGemData *TGemData) bool {
	if pGemData.ID <= 0 {
		gamelog.Error("AddGemData Error : Invalid GemID", pGemData.ID)
		return false
	}
	self.GemBag.Gems = append(self.GemBag.Gems, *pGemData)
	self.DB_AddGemAtLast()
	return true
}

//添加一个宝物到背包
func (self *TBagMoudle) AddGems(gemid int, num int) bool {
	if gemid <= 0 {
		gamelog.Error("AddGems Error: Invalid gemid:%d, num:%d", gemid, num)
		return false
	}

	var gems []TGemData = make([]TGemData, num)
	for i := 0; i < num; i++ {
		gems[i].Init(gemid)
	}

	self.GemBag.Gems = append(self.GemBag.Gems, gems...)
	self.DB_AddGemList(gems)

	return true
}

//删除指定的定物
func (self *TBagMoudle) RemoveGemAt(pos int) bool {
	if pos >= len(self.GemBag.Gems) {
		gamelog.Error("RemoveGemAt Error Pos :%d can greater than bagnum", pos)
		return false
	}

	if pos == 0 {
		self.GemBag.Gems = self.GemBag.Gems[1:]
	} else if (pos + 1) == len(self.GemBag.Gems) {
		self.GemBag.Gems = self.GemBag.Gems[:pos]
	} else {
		self.GemBag.Gems = append(self.GemBag.Gems[:pos], self.GemBag.Gems[pos+1:]...)
	}
	return true
}

//获取指定碎片的个数
func (self *TBagMoudle) GetEqiupPieceCount(itemid int) int {
	for _, t := range self.EquipPieceBag.Items {
		if t.ItemID == itemid {
			return t.ItemNum
		}
	}
	return 0
}

//添加装备碎片
func (self *TBagMoudle) AddEqiupPiece(itemid int, count int) int {
	if count <= 0 {
		gamelog.Error("AddEqiupPiece Error : Invalid count :%d", count)
		return 0
	}

	for i := 0; i < len(self.EquipPieceBag.Items); i++ {
		if self.EquipPieceBag.Items[i].ItemID == itemid {
			self.EquipPieceBag.Items[i].ItemNum += count
			self.DB_SaveEquipPieceBagAt(i)
			return self.EquipPieceBag.Items[i].ItemNum
		}
	}
	self.EquipPieceBag.Items = append(self.EquipPieceBag.Items, TItemData{ItemID: itemid, ItemNum: count})
	self.DB_SaveEquipPieceBagAt(len(self.EquipPieceBag.Items) - 1)
	return count
}

//删除装备碎片
func (self *TBagMoudle) RemoveEquipPiece(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("RemoveEqiupPiece Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	for i := 0; i < len(self.EquipPieceBag.Items); i++ {
		if self.EquipPieceBag.Items[i].ItemID == itemid {
			self.EquipPieceBag.Items[i].ItemNum -= count
			if self.EquipPieceBag.Items[i].ItemNum <= 0 {
				self.DB_RemoveEquipPiece(itemid)
			} else {
				self.DB_SaveEquipPieceBagAt(i)
			}
			return true
		}
	}
	return false
}

//获取指定碎片的个数
func (self *TBagMoudle) GetGemPieceCount(itemid int) int {
	for _, t := range self.GemPieceBag.Items {
		if t.ItemID == itemid {
			return t.ItemNum
		}
	}
	return 0
}

//添加装备碎片
func (self *TBagMoudle) AddGemPiece(itemid int, count int) int {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("AddGemPiece Error : Invalid itemid :%d, count:%d", itemid, count)
		return 0
	}

	for i := 0; i < len(self.GemPieceBag.Items); i++ {
		if self.GemPieceBag.Items[i].ItemID == itemid {
			self.GemPieceBag.Items[i].ItemNum += count
			self.DB_SaveGemPieceBagAt(i)
			return self.GemPieceBag.Items[i].ItemNum
		}
	}
	self.GemPieceBag.Items = append(self.GemPieceBag.Items, TItemData{ItemID: itemid, ItemNum: count})
	self.DB_SaveGemPieceBagAt(len(self.GemPieceBag.Items) - 1)
	return count
}

//删除装备碎片
func (self *TBagMoudle) RemoveGemPiece(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("RemoveGemPiece Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	for i := 0; i < len(self.GemPieceBag.Items); i++ {
		if self.GemPieceBag.Items[i].ItemID == itemid {
			if self.GemPieceBag.Items[i].ItemNum < count {
				return false
			}
			self.GemPieceBag.Items[i].ItemNum -= count
			if self.GemPieceBag.Items[i].ItemNum <= 0 {
				if i == 0 {
					self.GemPieceBag.Items = self.GemPieceBag.Items[1:]
				} else if (i + 1) == len(self.GemPieceBag.Items) {
					self.GemPieceBag.Items = self.GemPieceBag.Items[:i]
				} else {
					self.GemPieceBag.Items = append(self.GemPieceBag.Items[:i], self.GemPieceBag.Items[i+1:]...)
				}
				self.DB_RemoveGemPiece(itemid)
			} else {
				self.DB_SaveGemPieceBagAt(i)
			}
			return true
		}
	}
	return false
}

//获取英雄碎片数
func (self *TBagMoudle) GetHeroPieceCount(itemid int) int {
	for _, t := range self.HeroPieceBag.Items {
		if t.ItemID == itemid {
			return t.ItemNum
		}
	}

	return 0
}

//添加英雄碎片数
func (self *TBagMoudle) AddHeroPiece(itemid int, count int) int {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("AddHeroPiece Error :Invalid itemid :%d, count:%d", itemid, count)
		return 0
	}

	for i := 0; i < len(self.HeroPieceBag.Items); i++ {
		if self.HeroPieceBag.Items[i].ItemID == itemid {
			self.HeroPieceBag.Items[i].ItemNum += count
			self.DB_SaveHeroPieceBagAt(i)
			return self.HeroPieceBag.Items[i].ItemNum
		}
	}

	self.HeroPieceBag.Items = append(self.HeroPieceBag.Items, TItemData{ItemID: itemid, ItemNum: count})
	self.DB_SaveHeroPieceBagAt(len(self.HeroPieceBag.Items) - 1)
	return count
}

//删除英雄碎片数
func (self *TBagMoudle) RemoveHeroPiece(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("RemoveHeroPiece Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	for i := 0; i < len(self.HeroPieceBag.Items); i++ {
		if self.HeroPieceBag.Items[i].ItemID == itemid {
			if self.HeroPieceBag.Items[i].ItemNum < count {
				return false
			}
			self.HeroPieceBag.Items[i].ItemNum -= count
			if self.HeroPieceBag.Items[i].ItemNum <= 0 {
				if i == 0 {
					self.HeroPieceBag.Items = self.HeroPieceBag.Items[1:]
				} else if (i + 1) == len(self.HeroPieceBag.Items) {
					self.HeroPieceBag.Items = self.HeroPieceBag.Items[:i]
				} else {
					self.HeroPieceBag.Items = append(self.HeroPieceBag.Items[:i], self.HeroPieceBag.Items[i+1:]...)
				}
				self.DB_RemoveHeroPiece(itemid)
			} else {
				self.DB_SaveHeroPieceBagAt(i)
			}
			return true
		}
	}

	return false
}

func (self *TBagMoudle) GetNormalItemCount(itemid int) int {
	for _, t := range self.NormalItemBag.Items {
		if t.ItemID == itemid {
			return t.ItemNum
		}
	}

	return 0
}

func (self *TBagMoudle) AddNormalItem(itemid int, count int) int {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("AddItem Error : Invalid itemid :%d, count:%d", itemid, count)
		return 0
	}

	for i := 0; i < len(self.NormalItemBag.Items); i++ {
		if self.NormalItemBag.Items[i].ItemID == itemid {
			self.NormalItemBag.Items[i].ItemNum += count
			self.DB_SaveNormalItemBagAt(i)
			return self.NormalItemBag.Items[i].ItemNum
		}
	}

	self.NormalItemBag.Items = append(self.NormalItemBag.Items, TItemData{ItemID: itemid, ItemNum: count})
	self.DB_SaveNormalItemBagAt(len(self.NormalItemBag.Items) - 1)
	return count
}

func (self *TBagMoudle) RemoveNormalItem(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error3("RemoveItem Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	for i := 0; i < len(self.NormalItemBag.Items); i++ {
		if self.NormalItemBag.Items[i].ItemID == itemid {
			if self.NormalItemBag.Items[i].ItemNum < count {
				return false
			}
			self.NormalItemBag.Items[i].ItemNum -= count
			if self.NormalItemBag.Items[i].ItemNum == 0 {
				if i == 0 {
					self.NormalItemBag.Items = self.NormalItemBag.Items[1:]
				} else if (i + 1) == len(self.NormalItemBag.Items) {
					self.NormalItemBag.Items = self.NormalItemBag.Items[:i]
				} else {
					self.NormalItemBag.Items = append(self.NormalItemBag.Items[:i], self.NormalItemBag.Items[i+1:]...)
				}
				self.DB_RemoveNormalItem(itemid)
			} else {
				self.DB_SaveNormalItemBagAt(i)
			}
			return true
		}
	}

	return false
}

func (self *TBagMoudle) IsWakeItemEnough(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("IsWakeItemEnough Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	hascount := self.GetWakeItemCount(itemid)
	if hascount >= count {
		return true
	}

	return false
}

func (self *TBagMoudle) GetWakeItemCount(itemid int) int {
	for _, t := range self.WakeItemBag.Items {
		if t.ItemID == itemid {
			return t.ItemNum
		}
	}
	return 0
}

func (self *TBagMoudle) AddWakeItem(itemid int, count int) int {
	if count <= 0 || itemid <= 0 {
		gamelog.Error3("AddWakeItem Error : Invalid itemid :%d, count:%d", itemid, count)
		return 0
	}

	for i := 0; i < len(self.WakeItemBag.Items); i++ {
		if self.WakeItemBag.Items[i].ItemID == itemid {
			self.WakeItemBag.Items[i].ItemNum += count
			self.DB_SaveWakeItemBagAt(i)
			return self.WakeItemBag.Items[i].ItemNum
		}
	}

	self.WakeItemBag.Items = append(self.WakeItemBag.Items, TItemData{ItemID: itemid, ItemNum: count})
	self.DB_SaveWakeItemBagAt(len(self.WakeItemBag.Items) - 1)
	return count
}

func (self *TBagMoudle) RemoveWakeItem(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("RemoveItem Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	for i := 0; i < len(self.WakeItemBag.Items); i++ {
		if self.WakeItemBag.Items[i].ItemID == itemid {
			if self.WakeItemBag.Items[i].ItemNum < count {
				return false
			}
			self.WakeItemBag.Items[i].ItemNum -= count
			if self.WakeItemBag.Items[i].ItemNum <= 0 {
				if i == 0 {
					self.WakeItemBag.Items = self.WakeItemBag.Items[1:]
				} else if (i + 1) == len(self.WakeItemBag.Items) {
					self.WakeItemBag.Items = self.WakeItemBag.Items[:i]
				} else {
					self.WakeItemBag.Items = append(self.WakeItemBag.Items[:i], self.WakeItemBag.Items[i+1:]...)
				}
				self.DB_RemoveWakeItem(itemid)
			} else {
				self.DB_SaveWakeItemBagAt(i)
			}

			return true
		}
	}

	return false
}

func (self *TBagMoudle) IsItemEnough(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("IsItemEnough Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	hascount := self.GetNormalItemCount(itemid)
	if hascount >= count {
		return true
	}

	return false
}

//添加一个宠物到背包
func (self *TBagMoudle) AddPetByID(petid int) bool {
	if petid <= 0 {
		gamelog.Error("AddPetByID Invalid petid :%d", petid)
		return false
	}

	var pet TPetData
	pet.Init(petid)

	self.PetBag.Pets = append(self.PetBag.Pets, pet)

	bCol := true
	for i := 0; i < len(self.ColPets); i++ {
		if self.ColPets[i] == int16(petid) {
			bCol = false
			break
		}
	}

	if bCol == true {
		self.ColPets = append(self.ColPets, int16(petid))
	}

	self.DB_AddPetAtLast(bCol)

	pPetInfo := gamedata.GetPetInfo(petid)
	if pPetInfo != nil {
		self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PET_QUALITY, pPetInfo.Quality)
	}

	return true
}

//添加一个宠物到背包
func (self *TBagMoudle) AddPetData(pPetData *TPetData) bool {
	if pPetData.ID <= 0 {
		gamelog.Error("AddPetData Error : Invalid PetID", pPetData.ID)
		return false
	}
	self.PetBag.Pets = append(self.PetBag.Pets, *pPetData)
	self.DB_AddPetAtLast(false)
	return true
}

//添加宠物到背包
func (self *TBagMoudle) AddPets(petid int, num int) bool {
	if petid <= 0 {
		gamelog.Error("AddPets Error: Invalid petid:%d, num:%d", petid, num)
		return false
	}

	var pets []TPetData = make([]TPetData, num)
	for i := 0; i < num; i++ {
		pets[i].Init(petid)
	}

	self.PetBag.Pets = append(self.PetBag.Pets, pets...)
	self.DB_AddPetList(pets)

	return true
}

//删除指定的宠物
func (self *TBagMoudle) RemovePetAt(pos int) bool {
	if pos >= len(self.PetBag.Pets) {
		gamelog.Error("RemovePetAt Error Pos :%d can greater than bagnum", pos)
		return false
	}

	if pos == 0 {
		self.PetBag.Pets = self.PetBag.Pets[1:]
	} else if (pos + 1) == len(self.PetBag.Pets) {
		self.PetBag.Pets = self.PetBag.Pets[:pos]
	} else {
		self.PetBag.Pets = append(self.PetBag.Pets[:pos], self.PetBag.Pets[pos+1:]...)
	}
	return true
}

//获取宠物碎片数
func (self *TBagMoudle) GetPetPieceCount(itemid int) int {
	for _, t := range self.PetPieceBag.Items {
		if t.ItemID == itemid {
			return t.ItemNum
		}
	}

	return 0
}

//添加宠物碎片数
func (self *TBagMoudle) AddPetPiece(itemid int, count int) int {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("AddPetPiece Error :Invalid itemid :%d, count:%d", itemid, count)
		return 0
	}

	for i := 0; i < len(self.PetPieceBag.Items); i++ {
		if self.PetPieceBag.Items[i].ItemID == itemid {
			self.PetPieceBag.Items[i].ItemNum += count
			self.DB_SavePetPieceBagAt(i)
			return self.PetPieceBag.Items[i].ItemNum
		}
	}

	self.PetPieceBag.Items = append(self.PetPieceBag.Items, TItemData{ItemID: itemid, ItemNum: count})
	self.DB_SavePetPieceBagAt(len(self.PetPieceBag.Items) - 1)
	return count
}

//删除宠物碎片数
func (self *TBagMoudle) RemovePetPiece(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("RemovePetPiece Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	for i := 0; i < len(self.PetPieceBag.Items); i++ {
		if self.PetPieceBag.Items[i].ItemID == itemid {
			if self.PetPieceBag.Items[i].ItemNum < count {
				return false
			}
			self.PetPieceBag.Items[i].ItemNum -= count
			if self.PetPieceBag.Items[i].ItemNum <= 0 {
				if i == 0 {
					self.PetPieceBag.Items = self.PetPieceBag.Items[1:]
				} else if (i + 1) == len(self.PetPieceBag.Items) {
					self.PetPieceBag.Items = self.PetPieceBag.Items[:i]
				} else {
					self.PetPieceBag.Items = append(self.PetPieceBag.Items[:i], self.PetPieceBag.Items[i+1:]...)
				}
				self.DB_RemovePetPiece(itemid)
			} else {
				self.DB_SavePetPieceBagAt(i)
			}
			return true
		}
	}

	return false
}

//获取将灵数
func (self *TBagMoudle) GetHeroSoulCount(itemid int) int {
	for _, t := range self.HeroSoulBag.Items {
		if t.ItemID == itemid {
			return t.ItemNum
		}
	}

	return 0
}

//添加将灵
func (self *TBagMoudle) AddHeroSoul(itemid int, count int) int {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("AddHeroSoul Error :Invalid itemid :%d, count:%d", itemid, count)
		return 0
	}

	for i := 0; i < len(self.HeroSoulBag.Items); i++ {
		if self.HeroSoulBag.Items[i].ItemID == itemid {
			self.HeroSoulBag.Items[i].ItemNum += count
			self.DB_SaveHeroSoulBagAt(i)
			return self.HeroSoulBag.Items[i].ItemNum
		}
	}

	self.HeroSoulBag.Items = append(self.HeroSoulBag.Items, TItemData{ItemID: itemid, ItemNum: count})
	self.DB_SaveHeroSoulBagAt(len(self.HeroSoulBag.Items) - 1)
	return count
}

//删除将灵
func (self *TBagMoudle) RemoveHeroSoul(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("RemoveHeroSoul Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	for i := 0; i < len(self.HeroSoulBag.Items); i++ {
		if self.HeroSoulBag.Items[i].ItemID == itemid {
			if self.HeroSoulBag.Items[i].ItemNum < count {
				return false
			}
			self.HeroSoulBag.Items[i].ItemNum -= count
			if self.HeroSoulBag.Items[i].ItemNum <= 0 {
				if i == 0 {
					self.HeroSoulBag.Items = self.HeroSoulBag.Items[1:]
				} else if (i + 1) == len(self.HeroSoulBag.Items) {
					self.HeroSoulBag.Items = self.HeroSoulBag.Items[:i]
				} else {
					self.HeroSoulBag.Items = append(self.HeroSoulBag.Items[:i], self.HeroSoulBag.Items[i+1:]...)
				}
				self.DB_RemoveHeroSoul(itemid)
			} else {
				self.DB_SaveHeroSoulBagAt(i)
			}
			return true
		}
	}

	return false
}

//添加一个时装到背包
func (self *TBagMoudle) AddFashionByID(id int) bool {
	if id <= 0 {
		gamelog.Error("AddFashionByID Invalid id :%d", id)
		return false
	}

	var fashion TFashionData
	fashion.Init(id)

	self.FashionBag.Fashions = append(self.FashionBag.Fashions, fashion)

	self.DB_AddFashionAtLast()

	return true
}

//添加宠物到背包
func (self *TBagMoudle) AddFashions(id int, num int) bool {
	if id <= 0 {
		gamelog.Error("AddFashions Error: Invalid id:%d, num:%d", id, num)
		return false
	}

	var fashions []TFashionData = make([]TFashionData, num)
	for i := 0; i < num; i++ {
		fashions[i].Init(id)
	}

	self.FashionBag.Fashions = append(self.FashionBag.Fashions, fashions...)
	self.DB_AddFashionList(fashions)

	return true
}

//删除指定的宠物
func (self *TBagMoudle) RemoveFashionByID(id int) bool {
	if id <= 0 {
		gamelog.Error("RemoveFashionByID Error: Invalid id:%d", id)
		return false
	}

	return true
}

//获取宠物碎片数
func (self *TBagMoudle) GetFashionPieceCount(itemid int) int {
	for _, t := range self.FashionPieceBag.Items {
		if t.ItemID == itemid {
			return t.ItemNum
		}
	}

	return 0
}

//添加宠物碎片数
func (self *TBagMoudle) AddFashionPiece(itemid int, count int) int {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("AddFashionPiece Error :Invalid itemid :%d, count:%d", itemid, count)
		return 0
	}

	for i := 0; i < len(self.FashionPieceBag.Items); i++ {
		if self.FashionPieceBag.Items[i].ItemID == itemid {
			self.FashionPieceBag.Items[i].ItemNum += count
			self.DB_SaveFashionPieceBagAt(i)
			return self.FashionPieceBag.Items[i].ItemNum
		}
	}

	self.FashionPieceBag.Items = append(self.FashionPieceBag.Items, TItemData{ItemID: itemid, ItemNum: count})
	self.DB_SaveFashionPieceBagAt(len(self.PetPieceBag.Items) - 1)
	return count
}

//删除宠物碎片数
func (self *TBagMoudle) RemoveFashionPiece(itemid int, count int) bool {
	if count <= 0 || itemid <= 0 {
		gamelog.Error("RemoveFashionPiece Error : Invalid itemid :%d, count:%d", itemid, count)
		return false
	}

	for i := 0; i < len(self.FashionPieceBag.Items); i++ {
		if self.FashionPieceBag.Items[i].ItemID == itemid {
			if self.FashionPieceBag.Items[i].ItemNum < count {
				return false
			}
			self.FashionPieceBag.Items[i].ItemNum -= count
			if self.FashionPieceBag.Items[i].ItemNum <= 0 {
				if i == 0 {
					self.FashionPieceBag.Items = self.FashionPieceBag.Items[1:]
				} else if (i + 1) == len(self.FashionPieceBag.Items) {
					self.FashionPieceBag.Items = self.FashionPieceBag.Items[:i]
				} else {
					self.FashionPieceBag.Items = append(self.FashionPieceBag.Items[:i], self.FashionPieceBag.Items[i+1:]...)
				}
				self.DB_RemoveFashionPiece(itemid)
			} else {
				self.DB_SaveFashionPieceBagAt(i)
			}
			return true
		}
	}

	return false
}

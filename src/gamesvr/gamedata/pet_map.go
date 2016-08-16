package gamedata

//宠物图鉴配制表

type ST_PetMapInfo struct {
	ID     int //宠物图鉴ID
	PetIds []int
	Buffs  [3]ST_PropertyBuff //宠物的Buff集
}

var (
	GT_PetMap_List []ST_PetMapInfo = nil
)

func InitPetMapParser(total int) bool {
	GT_PetMap_List = make([]ST_PetMapInfo, total+1)
	return true
}

func ParsePetMapRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_PetMap_List[id].ID = id

	petid := rs.GetFieldInt("pet_id1")
	GT_PetMap_List[id].PetIds = append(GT_PetMap_List[id].PetIds, petid)

	petid = rs.GetFieldInt("pet_id2")
	if petid > 0 {
		GT_PetMap_List[id].PetIds = append(GT_PetMap_List[id].PetIds, petid)
	}

	petid = rs.GetFieldInt("pet_id3")
	if petid > 0 {
		GT_PetMap_List[id].PetIds = append(GT_PetMap_List[id].PetIds, petid)
	}

	GT_PetMap_List[id].Buffs[0].PropertyID = rs.GetFieldInt("property1")
	GT_PetMap_List[id].Buffs[0].Value = rs.GetFieldInt("value1")
	GT_PetMap_List[id].Buffs[0].IsPercent = rs.GetFieldInt("is_percent1") == 1
	GT_PetMap_List[id].Buffs[1].PropertyID = rs.GetFieldInt("property2")
	GT_PetMap_List[id].Buffs[1].Value = rs.GetFieldInt("value2")
	GT_PetMap_List[id].Buffs[1].IsPercent = rs.GetFieldInt("is_percent2") == 1
	GT_PetMap_List[id].Buffs[2].PropertyID = rs.GetFieldInt("property3")
	GT_PetMap_List[id].Buffs[2].Value = rs.GetFieldInt("value3")
	GT_PetMap_List[id].Buffs[2].IsPercent = rs.GetFieldInt("is_percent3") == 1
}

func (self *ST_PetMapInfo) IsMapOK(pets []int16) bool {
	for i := 0; i < 3; i++ {
		bFind := false
		for j := 0; j < len(pets); j++ {
			if self.PetIds[i] == int(pets[j]) || self.PetIds[i] == 0 {
				bFind = true
				break
			}
		}

		if bFind == false {
			return false
		}
	}

	return true
}

package gamedata

import (
	"utility"
)

var (
	GT_Xing_List []string = nil
	GT_Ming_List []string = nil
)

func InitNameParser(total int) bool {
	GT_Xing_List = make([]string, total+1)
	GT_Ming_List = make([]string, total+1)
	return true
}

func ParseNameRecord(rs *RecordSet) {

	return
}

func GetRandName() string {
	var nXing = utility.Rand() % len(GT_Xing_List)
	var nMing = utility.Rand() % len(GT_Ming_List)
	return GT_Xing_List[nXing] + GT_Ming_List[nMing]
}

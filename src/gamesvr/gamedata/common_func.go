package gamedata

import (
	"strings"
)

func ParseToIntSlice(svalue string) (ret []int) {
	sv := strings.Split(svalue, "|")
	ret = make([]int, len(sv))
	for i := 0; i < len(sv); i++ {
		ret[i] = CheckAtoi(sv[i], 91)
	}
	return
}

func ParseTo2IntSlice(svalue string) (ret [2]int) {
	sv := strings.Split(svalue, "|")

	if len(sv) <= 1 {
		panic("ParseTo2IntSlice Error Invalid Length")
		return
	}

	ret[0] = CheckAtoi(sv[0], 92)
	ret[1] = CheckAtoi(sv[1], 93)

	return
}

func ParseTo2Int8Slice(svalue string) (ret [2]int8) {
	sv := strings.Split(svalue, "|")

	if len(sv) <= 1 {
		panic("ParseTo2IntSlice Error Invalid Length")
		return
	}

	ret[0] = int8(CheckAtoi(sv[0], 92))
	ret[1] = int8(CheckAtoi(sv[1], 93))

	return
}

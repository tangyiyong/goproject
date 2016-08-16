package utility

import (
	"math/rand"
	"time"
)

const MaxRandNum = 10000

var (
	randValueList [MaxRandNum]int16
	nCurIndex     = 0
)

func disOrder() {
	var nIndex int
	var nTemp int16
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < MaxRandNum; i++ {
		nIndex = rand.Int() % (i + 1)
		if nIndex != i {
			nTemp = randValueList[i]
			randValueList[i] = randValueList[nIndex]
			randValueList[nIndex] = nTemp
		}
	}
}

func initRandom() {
	for i := int16(0); i < MaxRandNum; i++ {
		randValueList[i] = i
	}

	disOrder()
}

func GetRandValue16() int16 {

	nCurIndex = (nCurIndex + 1) % MaxRandNum

	return randValueList[nCurIndex]
}

func GetRandValueInt() int {

	nCurIndex = (nCurIndex + 1) % MaxRandNum

	return int(randValueList[nCurIndex])
}

func Rand() int {
	return GetRandValueInt()
}

func HitRandTest(value int) bool {
	if value > GetRandValueInt() {
		return true
	}

	return false
}

func RandBetween(left, right int) int { // [a, b]
	return rand.Intn(right+1-left) + left
}

func RandShuffle(slice []int) {
	length := len(slice)
	for i := 0; i < length; i++ {
		ri := rand.Intn(length-i) + i
		temp := slice[ri]
		slice[i] = temp
		slice[ri] = slice[i]
	}
}

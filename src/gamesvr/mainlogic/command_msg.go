package mainlogic

import (
	"fmt"
	"gamelog"
	"strconv"
	"time"
)

func HandCmd_AddHeros(args []string) bool {
	PlayerID, err := strconv.Atoi(args[1])
	if err != nil {
		gamelog.Error("HandCmd_AddHeros Error : Invalid PlayerID :%d", PlayerID)
		return true
	}

	var player *TPlayer = GetPlayerByID(int32(PlayerID))
	if player == nil {
		gamelog.Error("HandCmd_AddHeros error : Cannot get player by ID : %d", PlayerID)
		return true
	}

	t1 := time.Now().UnixNano()
	for i := 5; i < 148; i++ {
		player.BagMoudle.AddAwardItem(i, 1)
	}

	//player.BagMoudle.RemoveHero(3)
	fmt.Println(args[0], "Finished Time:", time.Now().UnixNano()-t1)
	return true
}

func HandCmd_SetLogLevel(args []string) bool {
	level, err := strconv.Atoi(args[1])
	if err != nil {
		gamelog.Error("HandCmd_SetLogLevel Error : Invalid param :%s", args[1])
		return true
	}
	gamelog.SetLevel(level)
	return true
}

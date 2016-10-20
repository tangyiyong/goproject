package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"utility"
)

var GAME_DATA_PATH = utility.GetCurrPath() + "game.dat"

type TLocalData struct {
	OpenSvrTime int64
	Param       int32
	Param2      bool
	Param3      string
}

var (
	G_LocalArchive TLocalData //服务器本地数据
)

func ReadGameData() bool {
	var f *os.File
	var err error
	_, err = os.Stat(GAME_DATA_PATH)
	if !os.IsNotExist(err) {
		f, err = os.Open(GAME_DATA_PATH)
		if err != nil {
			return false
		}

		b, _ := ioutil.ReadAll(f)
		json.Unmarshal(b, &G_LocalArchive)
	}
	return true
}

func SaveGameData() bool {
	b, _ := json.Marshal(&G_LocalArchive)
	ioutil.WriteFile(GAME_DATA_PATH, b, 777)
	return true

	bstr := string(b)
	newstring := strings.Replace(bstr, ",\"", ",\n\"", -1)
	newstring = strings.Replace(newstring, "{\"", "{\n\"", -1)
	newstring = strings.Replace(newstring, "\"}", "\"\n}", -1)
	ioutil.WriteFile(GAME_DATA_PATH, []byte(newstring), 777)

	return true
}

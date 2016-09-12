package appconfig

import (
	"bufio"
	"os"
	"strings"
	"utility"
)

var CONFIG_PATH = utility.GetCurrPath() + "config.ini"

func LoadConfig() bool {
	bRet := parseConfigFile(CONFIG_PATH)
	if bRet {
		bRet = initGlobalVar()
	}

	return bRet
}

func parseConfigFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		panic("parseConfigFile Failed Error :" + err.Error())
		return false
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var eqpos = 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) <= 2 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		eqpos = 0
		for eqpos = 0; eqpos < len(line); eqpos++ {
			if line[eqpos] == '=' {
				break
			}
		}

		if eqpos >= len(line)-1 {
			continue
		}

		ParseConfigValue(strings.TrimSpace(line[0:eqpos]), strings.TrimSpace(line[eqpos+1:]))
	}

	return true
}

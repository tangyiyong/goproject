package main

import (
	// "appconfig"
	"fmt"
	"gamelog"
	"gamesvr/tcpclient"
	//"io"
	"bufio"
	"io/ioutil"
	"msg"
	"os"
	"strconv"
	"strings"
)

var (
	accountip string = "http://192.168.0.222:8081/"
	gamesvrip string = "http://172.17.2.249:8082/"
)

//消息:/gm_add_giftaward
type MSG_AddGiftAward_Req struct {
	SessionID  string
	SessionKey string
	ItemID     []int //物品ID
	ItemNum    []int //物品数量
}

type MSG_AddGiftAward_Ack struct {
	RetCode int
}

//消息:/gm_make_giftcode
type MSG_MakeGiftCode_Req struct {
	SessionID   string //GM SessionID
	SessionKey  string //GM SessionKey
	Platform    int32  //平台ID
	SvrID       int32  //服务器ID
	EndTime     int32  //结束时间
	GiftAwardID int32  //奖励ID
	GiftCodeNum int    //激活码数量
	IsAll       bool   //是否为全服发放
}

type MSG_MakeGiftCode_Ack struct {
	RetCode   int
	GiftCodes []string //激活码
}

type t_struct struct {
	player string
	id     int
}

func test() (str string, err error) {
	return
}

func main() {

	str, err := test()
	fmt.Println(str, " 111 ", err)

	return
	sFix := strings.Trim("(2|130000&130000|0)(110|13&13|0)(116|10&10|0)(117|12&12|0)", "()")
	fmt.Println(sFix)
	return

	buffer, err := ioutil.ReadFile("../bin/csv/type_activity.csv")
	if err != nil {
		fmt.Println(err.Error())
	}

	strLst := strings.Split(string(buffer), "\r\n")

	for i, v := range strLst {
		paramLst := strings.Split(string(v), ",")
		if len(paramLst) < 1 {
			continue
		}

		id, _ := strconv.Atoi(paramLst[0])
		fmt.Printf("%d   %s\r\n", id, paramLst[0])
		if id == 89 {
			strLst[i] = "88,异域Shit商人,异域,快点来充钱,3,1,11,14,1,16,1,1,11,2,0"
		}
	}

	os.Remove("../bin/csv/type_activity.csv")
	f, err := os.OpenFile("../bin/csv/type_activity.csv", os.O_RDWR|os.O_CREATE, 0660)
	if err == nil {
		w := bufio.NewWriter(f)

		for _, v := range strLst {
			fmt.Printf("%s\r\n", v)
			w.WriteString(v + "\r\n")
		}

		w.Flush()

		f.Close()
	}

	// f, err := os.OpenFile("../bin/csv/type_activity.csv", os.O_RDWR, 0660)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// defer f.Close()

	// var offset int64
	// r := bufio.NewReader(f)
	// for {
	// 	line, _, err := r.ReadLine()
	// 	if err == io.EOF {
	// 		break
	// 	}

	// 	strLst := strings.Split(string(line), ",")
	// 	id, _ := strconv.Atoi(strLst[0])

	// 	if id == 89 {
	// 		//! 删去该行
	// 		fs, _ := f.Stat()
	// 		size := fs.Size()
	// 		begin := make([]byte, offset)
	// 		f.Seek(0, os.SEEK_SET)
	// 		f.Read(begin)

	// 		offset += int64(len(line) + 2)
	// 		end := make([]byte, size-offset)
	// 		f.Seek(offset, 0)
	// 		f.Read(end)

	// newFile, err := os.OpenFile("../bin/csv/type_activity.csv.temp", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	//			defer newFile.Close()
	// w := bufio.NewWriter(f)
	// w.Write(begin)
	// w.Write([]byte("88,异域商人,异域商人(新）,快点来充钱,3,1,11,14,1,16,1,1,11,2,0\r\n"))
	// w.Write(end)
	// w.Flush()

	// f.Seek(0, 0)

	// f.Seek(offset, 0)
	// w = bufio.NewWriter(f)
	// w.WriteString(string("89,1111,11111111111111111111111111111111111,2222,5,3,0,0,0,4,1,1,28,1,0"))
	// w.Flush()

	// 		break
	// 	}

	// 	offset += (int64(len(line) + 2))
	// }

	// f.Close()
	// newFile.Close()

	// os.Remove("../bin/csv/type_activity.csv")
	// os.Rename("../bin/csv/type_activity.csv.temp", "../bin/csv/type_activity.csv")

	// r := csv.NewReader(file)
	// csvStr, _ := r.ReadAll()
	// fmt.Println(csvStr)

	// RegTcpMsgHandler()

	// InitPlayerMgr()
	// for i := 2; i < 3; i++ {
	// 	CreatePlayer(i)
	// }

	// StartTest()
	// utility.StartConsoleWait()
}

func RegTcpMsgHandler() {
	tcpclient.HandleFunc(msg.MSG_DISCONNECT, Hand_DisConnect)
	tcpclient.HandleFunc(msg.MSG_CONNECT, Hand_Connect)
	tcpclient.HandleFunc(msg.MSG_ENTER_ROOM_ACK, Hand_EnterRoomAck)
	tcpclient.HandleFunc(msg.MSG_ENTER_ROOM_NTY, Hand_NoneFunction)
	tcpclient.HandleFunc(msg.MSG_MOVE_STATE, Hand_NoneFunction)

}

func Hand_NoneFunction(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
}

func Hand_Connect(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_Connect")

	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		gamelog.Info("Hand_Connect Error: pClient == nil")
		return
	}

	player := pClient.ExtraData.(*TPlayer)
	if player == nil {
		gamelog.Info("Hand_Connect Error: player == nil")
		return
	}

	if pClient.ConType == tcpclient.CON_TYPE_BATSVR {
		player.userEnterRoom()
	} else {
		player.userCheckIn()
	}

	return
}

func Hand_DisConnect(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_DisConnect")
	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		return
	}

	return
}

func Hand_EnterRoomAck(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_EnterRoomAck")
	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		return
	}

	player := pClient.ExtraData.(*TPlayer)
	if player == nil {
		gamelog.Info("Hand_EnterRoomAck Error: player == nil")
		return
	}

	var req msg.MSG_EnterRoom_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_EnterRoomAck : Message Reader Error!!!!")
		return
	}

	player.Heros = req.Heros

	player.IsEnter = true

	player.PackNo = req.BeginMsgNo

	return
}

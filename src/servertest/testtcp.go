// main
package main

import (
	"appconfig"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"msg"
	"net"
	"time"
	"utility"
)

func TestTcp() {
	//加载配制文件
	appconfig.LoadConfig()

	for i := 1; i < 10; i++ {
		time.Sleep(time.Millisecond)
		go testProcess(i)
	}

	utility.StartConsoleWait()
}

func testProcess(i int) {
	conn, err := net.Dial("tcp", appconfig.ChatSvrAddr)
	if err != nil {
		fmt.Println("connect  error playerid : ", i)
		return
	}

	if conn == nil {
		return
	}

	go ReadProcess(conn)

	for {
		time.Sleep(1 * time.Second)
		var req msg.MSG_CheckIn_Req
		req.PlayerID = i
		req.GuildID = rand.Int() % 100
		b, _ := json.Marshal(&req)
		WriteMsg(conn, msg.MSG_CHECK_IN_REQ, b)

		time.Sleep(2 * time.Second)
		var req2 msg.MSG_Chat_Req
		req2.TargetChannel = msg.MSG_CHANNEL_GUILD
		req2.TargetGuildID = rand.Int() % 100
		req2.MsgContent = "this is the message !!!!"
		b2, _ := json.Marshal(&req2)
		WriteMsg(conn, msg.MSG_CHATMSG_REQ, b2)
	}

}

func WriteMsg(conn net.Conn, msgID int16, msgdata []byte) bool {
	msgLen := len(msgdata)

	msgbuffer := make([]byte, 6+msgLen)

	binary.LittleEndian.PutUint32(msgbuffer, uint32(msgLen))

	binary.LittleEndian.PutUint16(msgbuffer[4:], uint16(msgID))

	copy(msgbuffer[6:], msgdata)

	conn.Write(msgbuffer)

	return true
}

func ReadProcess(conn net.Conn) error {
	defer func() {
		conn.Close()
		fmt.Println("TCPConn fianlly been closed:")
	}()

	var err error
	var msgHeader = make([]byte, 6)
	var msgID int16
	var msgLen int32
	for {
		_, err = io.ReadAtLeast(conn, msgHeader, 6)
		if err != nil {
			return err
		}

		// parse len
		msgLen = int32(binary.LittleEndian.Uint16(msgHeader[:4]))
		msgID = int16(binary.LittleEndian.Uint16(msgHeader[4:]))

		// data
		msgData := make([]byte, msgLen)
		_, err = io.ReadAtLeast(conn, msgData, int(msgLen))
		if err != nil {
			return err
		}

		fmt.Println("msgID:", msgID, "msgLen: ", msgLen)
		if msgID == msg.MSG_CHATMSG_NOTIFY {
			var notify msg.MSG_Chat_Msg_Notify
			json.Unmarshal(msgData, &notify)
			fmt.Println("msgContent:", notify.MsgContent)
		}
	}

	return nil
}

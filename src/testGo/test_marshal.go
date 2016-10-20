// main
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"gamelog"
	"io"
	"io/ioutil"
	"msg"
	"net/http"
	"os"
	"utility"
)

func PostServerReq(url string, body io.Reader) ([]byte, error) {
	//TestLock.Lock()
	//t1 := time.Now().UnixNano()
	resp, err := http.Post(url, "text/HTML", body)
	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	resp.Body.Close()
	//fmt.Println("t:", time.Now().UnixNano()-t1)
	//TestLock.Unlock()

	return buffer, err
}

func TestMarshal() {
	gamelog.InitLogger("test", 0)

	filepath := "C:/Users/DJ/Desktop/新建文本文档 (3).txt"
	file, err := os.Open(filepath)
	if err != nil {
		gamelog.Error("TestMarshal Error: %v", err)
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		gamelog.Error("TestMarshal Error: %v", err)
		return
	}

	fileData := utility.CompressData(data)

	reqUrl := "http://127.0.0.1:8082/update_gamedata"
	var req msg.MSG_UpdateGameData_Req
	req.SessionID = "1"
	req.SessionKey = "dafd"
	req.TbName = "type_action"
	msgData, _ := json.Marshal(req)
	msgdatalen := len(msgData)
	totaldata := make([]byte, len(fileData)+len(msgData)+4)
	binary.LittleEndian.PutUint32(totaldata, uint32(msgdatalen))
	copy(totaldata[4:], msgData)
	copy(totaldata[4+msgdatalen:], fileData)

	PostServerReq(reqUrl, bytes.NewReader(totaldata))
}

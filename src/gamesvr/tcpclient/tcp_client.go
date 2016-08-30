package tcpclient

import (
	"gamelog"
	"net"
	"time"
)

const (
	CON_TYPE_CHAT   = 1
	CON_TYPE_BATSVR = 2
)

type TCPClient struct {
	Addr            string
	PendingWriteNum int
	TcpConn         *TCPConn
	Reconnect       int //重连次数
	ConType         int //连接类型 1:聊天服连接，2: 战斗服连接
	SvrID           int //连接服务器ID
	ExtraData       interface{}
}

type MsgHanler func(pTcpConn *TCPConn, extra int16, pdata []byte)

var (
	HandlerMap map[int16]func(pTcpConn *TCPConn, extra int16, pdata []byte)
)

func Init() bool {
	HandlerMap = make(map[int16]func(pTcpConn *TCPConn, extra int16, pdata []byte))
	return true
}

func HandleFunc(msgid int16, mh MsgHanler) {
	if HandlerMap == nil {
		HandlerMap = make(map[int16]func(pTcpConn *TCPConn, extra int16, pdata []byte), 100)
	}

	HandlerMap[msgid] = mh

	return
}

func (client *TCPClient) ConnectToSvr(addr string, reconnect int) {
	client.Addr = addr
	client.PendingWriteNum = 32
	client.TcpConn = nil
	client.Reconnect = reconnect

	//会断线后自动重连
	go client.connectRoutine()
}

func (client *TCPClient) connectRoutine() {
	recontime := client.Reconnect
	for {
		if recontime > 0 {
			recontime = recontime - 1
		} else {
			break
		}

		if client.connect() {
			if client.TcpConn != nil {
				go client.TcpConn.WriteRoutine()
				msgDispatcher(client.TcpConn, 1, 0, nil)
				client.TcpConn.ReadRoutine()
				recontime = client.Reconnect
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (client *TCPClient) connect() bool {
	conn, err := net.Dial("tcp", client.Addr)
	if err != nil {
		gamelog.Error("connect to %s error :%s", client.Addr, err.Error())
		return false
	}

	if conn == nil {
		return false
	}

	if client.TcpConn != nil {
		client.TcpConn.ResetConn(conn)
	} else {
		if client.ConType <= 0 {
			gamelog.Error("connect error invalid contype : %d", client.ConType)
			return false
		}

		client.TcpConn = newTCPConn(conn, client.PendingWriteNum)
		client.TcpConn.Data = client

	}

	return true
}

func (client *TCPClient) Close() {
	client.TcpConn.Close()
	client.TcpConn = nil
}

func msgDispatcher(pTcpConn *TCPConn, msgid int16, extra int16, pdata []byte) {
	msghandler, ok := HandlerMap[msgid]
	if !ok {
		gamelog.Error("msgid: [%d] need a handler!!!", msgid)
		return
	}

	msghandler(pTcpConn, extra, pdata)
}

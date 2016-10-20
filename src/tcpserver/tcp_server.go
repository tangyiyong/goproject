package tcpserver

import (
	"gamelog"
	"net"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type ConnSet map[net.Conn]bool
type MsgHanler func(pTcpConn *TCPConn, Extra int16, pdata []byte)

var MsgDispatcher func(pTcpConn *TCPConn, MsgID int16, Extra int16, pdata []byte)

type TCPServer struct {
	Addr            string
	MaxConnNum      int
	PendingWriteNum int
	listener        net.Listener
	mutexConns      sync.Mutex
	connset         ConnSet
	wgLn            sync.WaitGroup
	wgConns         sync.WaitGroup
}

var HandlerMap map[int16]func(pTcpConn *TCPConn, extra int16, pdata []byte)
var G_Server *TCPServer

func ServerRun(addr string, maxconn int) {
	G_Server = new(TCPServer)
	G_Server.Addr = addr
	G_Server.MaxConnNum = maxconn
	G_Server.PendingWriteNum = 32
	if MsgDispatcher == nil {
		MsgDispatcher = DefalutDispatcher
	}

	G_Server.Init()
	G_Server.Run()
	G_Server.Close()
	return
}

func (server *TCPServer) Init() bool {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		gamelog.Error("TCPServer Init failed  error :%s", err.Error())
		return false
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 5000
		gamelog.Info("Invalid MaxConnNum, reset to %d", server.MaxConnNum)
	}
	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 32
		gamelog.Info("Invalid PendingWriteNum, reset to %d", server.PendingWriteNum)
	}

	server.listener = ln
	server.connset = make(ConnSet)

	return true
}

func (server *TCPServer) Run() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()
	var tempDelay time.Duration
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				gamelog.Error("accept error: %s retrying in %d", err.Error(), tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			gamelog.Error("accept error: %s", err.Error())
			return
		}
		tempDelay = 0
		connNum := len(server.connset)
		server.mutexConns.Lock()
		if connNum >= server.MaxConnNum {
			server.mutexConns.Unlock()
			conn.Close()
			gamelog.Error("too many connections")
			continue
		}

		server.connset[conn] = true
		server.mutexConns.Unlock()
		server.wgConns.Add(1)
		gamelog.Info("Connect From: %s,  ConnNum: %d", conn.RemoteAddr().String(), connNum+1)
		tcpConn := newTCPConn(conn, server.PendingWriteNum)
		conn.(*net.TCPConn).SetKeepAlive(true)
		conn.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 60)
		tcpConn.OnNetClose = func() {
			// 清理tcp_server相关数据
			server.mutexConns.Lock()
			delete(server.connset, tcpConn.conn)
			close(tcpConn.writeChan) //服务端没有重连，直接关闭chan
			gamelog.Info("Connect Endded:ConnID:%d, ConnNum is:%d", tcpConn.ConnID, len(server.connset))
			server.mutexConns.Unlock()
			server.wgConns.Done()
		}
		go tcpConn.ReadRoutine()
		go tcpConn.WriteRoutine()
	}
}

func (server *TCPServer) Close() {
	server.listener.Close()
	server.wgLn.Wait()

	server.mutexConns.Lock()
	for conn := range server.connset {
		conn.Close()
	}

	server.connset = nil
	server.mutexConns.Unlock()
	server.wgConns.Wait()
	gamelog.Info("server been closed!!")
}

func DefalutDispatcher(pTcpConn *TCPConn, msgid int16, extra int16, pdata []byte) {
	msghandler, ok := HandlerMap[msgid]
	if !ok {
		gamelog.Error("msgid : %d have not a msg handler!!", msgid)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				gamelog.Error("MsgID %d Error  %s", msgid, debug.Stack())
			}
		}
	}()
	msghandler(pTcpConn, extra, pdata)
}

func HandleFunc(msgid int16, mh MsgHanler) {
	if HandlerMap == nil {
		HandlerMap = make(map[int16]func(pTcpConn *TCPConn, extra int16, pdata []byte), 100)
	}

	HandlerMap[msgid] = mh

	return
}

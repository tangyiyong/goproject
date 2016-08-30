package tcpclient

import (
	"bufio"
	"encoding/binary"
	"gamelog"
	"io"
	"msg"
	"net"
)

type TCPConn struct {
	conn      net.Conn
	reader    *bufio.Reader //包装conn减少conn.Read的io次数
	writeChan chan []byte
	closeFlag bool
	Data      interface{}
}

func newTCPConn(conn net.Conn, pendingWriteNum int) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.ResetConn(conn)
	tcpConn.writeChan = make(chan []byte, pendingWriteNum)
	tcpConn.Data = nil
	return tcpConn
}

//closeFlag仅在ResetConn、Close中设置，其它地方只读
func (tcpConn *TCPConn) ResetConn(conn net.Conn) {
	tcpConn.conn = conn
	tcpConn.reader = bufio.NewReader(conn)
	tcpConn.closeFlag = false
}
func (tcpConn *TCPConn) Close() {
	if tcpConn.closeFlag {
		return
	}
	tcpConn.conn.Close()
	tcpConn.doWrite(nil) //触发writeRoutine结束
	tcpConn.closeFlag = true
}

func (tcpConn *TCPConn) doWrite(b []byte) {
	select {
	case tcpConn.writeChan <- b: //chan满后再写即阻塞，select进入default分支报错
	default:
		gamelog.Error("doWrite: channel full")
		tcpConn.conn.(*net.TCPConn).SetLinger(0)
		tcpConn.conn.Close()
		// close(tcpConn.writeChan) //重连后chan里的数据得保留
	}
}

// b must not be modified by other goroutines
func (tcpConn *TCPConn) write(b []byte) {
	if tcpConn.closeFlag || b == nil {
		return
	}

	tcpConn.doWrite(b)
}

func (tcpConn *TCPConn) WriteMsg(msgID int16, extra int16, msgdata []byte) bool {
	msgLen := len(msgdata)
	msgbuffer := make([]byte, 8+msgLen)
	binary.LittleEndian.PutUint32(msgbuffer, uint32(msgLen))
	binary.LittleEndian.PutUint16(msgbuffer[4:6], uint16(extra))
	binary.LittleEndian.PutUint16(msgbuffer[6:8], uint16(msgID))
	copy(msgbuffer[8:], msgdata)
	tcpConn.write(msgbuffer)
	return true
}

func (tcpConn *TCPConn) WriteMsgContinue(playerid int32, msgID int16, extra int16, msgdata []byte) bool {
	msgLen := len(msgdata)
	msgbuffer := make([]byte, 16+msgLen)
	binary.LittleEndian.PutUint32(msgbuffer, uint32(msgLen+10))
	binary.LittleEndian.PutUint16(msgbuffer[4:6], uint16(0))
	binary.LittleEndian.PutUint16(msgbuffer[6:8], uint16(msg.MSG_GAME_TO_CLIENT))
	binary.LittleEndian.PutUint32(msgbuffer[8:12], uint32(playerid))
	binary.LittleEndian.PutUint16(msgbuffer[12:14], uint16(extra))
	binary.LittleEndian.PutUint16(msgbuffer[14:16], uint16(msgID))
	copy(msgbuffer[12:], msgdata)
	tcpConn.write(msgbuffer)

	return true
}

func (tcpConn *TCPConn) WriteMsgData(msgdata []byte) bool {
	tcpConn.write(msgdata)
	return true
}

func (tcpConn *TCPConn) ReadProcess() error {
	var err error
	var msgHeader = make([]byte, 8)

	var msgLen int32
	var msgID int16
	var Extra int16

	//循环结束，会在ReadRoutine中紧接着关闭tcpConn
	for {
		if tcpConn.closeFlag {
			break
		}

		_, err = io.ReadAtLeast(tcpConn.reader, msgHeader, 8)
		if err != nil {
			gamelog.Error("ReadAtLeast error: %s", err.Error())
			return err
		}

		// parse len
		msgLen = int32(binary.LittleEndian.Uint16(msgHeader[:4]))
		if msgLen <= 0 || msgLen > 10240 {
			gamelog.Error("ReadProcess error: Invalid msgLen :%d", msgLen)
			break
		}

		Extra = int16(binary.LittleEndian.Uint16(msgHeader[4:6]))
		msgID = int16(binary.LittleEndian.Uint16(msgHeader[6:8]))
		if msgID <= msg.MSG_BEGIN || msgID >= msg.MSG_END {
			gamelog.Error("ReadProcess error: Invalid msgID :%d", msgID)
			break
		}

		// data
		msgData := make([]byte, msgLen)
		_, err = io.ReadAtLeast(tcpConn.reader, msgData, int(msgLen))
		if err != nil {
			gamelog.Error("ReadAtLeast error: %s", err.Error())
			return err
		}

		msgDispatcher(tcpConn, msgID, Extra, msgData)
	}

	return nil
}

//连接的写协程
func (tcpConn *TCPConn) WriteRoutine() {
	for b := range tcpConn.writeChan {
		if b == nil {
			break
		}
		_, err := tcpConn.conn.Write(b)
		if err != nil {
			gamelog.Error("WriteRoutine error: %s", err.Error())
			break
		}
	}
	tcpConn.Close()
}

//连接的读协程
func (tcpConn *TCPConn) ReadRoutine() {
	tcpConn.ReadProcess()
	tcpConn.Close()

	msgDispatcher(tcpConn, msg.MSG_DISCONNECT, 0, nil)
}

//连接的读协程
func (tcpConn *TCPConn) IsConnected() bool {
	return tcpConn.closeFlag
}

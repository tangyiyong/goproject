package msg

import (
	"gamelog"
	"math"
)

type TMsg interface {
	Read(reader *PacketReader) bool
	Write(writer *PacketWriter)
}

type PacketReader struct {
	DataPtr  []byte
	TotalLen int
	ReadPos  int
}

func (self *PacketReader) GetDataPtr() []byte {
	return self.DataPtr
}

func (self *PacketReader) BeginRead(ptr []byte, datalen int) *PacketReader {
	self.ReadPos = 0
	self.DataPtr = ptr
	if datalen <= 0 {
		self.TotalLen = len(ptr)
	} else {
		self.TotalLen = datalen
	}

	return self

}

func (self *PacketReader) EndRead() bool {
	return true
}

func (self *PacketReader) ReadInt8() (ret int8) {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadInt8 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return 0
	}

	ret = int8(self.DataPtr[self.ReadPos])
	self.ReadPos++
	return
}

func (self *PacketReader) ReadInt16() (ret int16) {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadInt16 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return 0
	}
	ret = int16(self.DataPtr[self.ReadPos+1])<<8 | int16(self.DataPtr[self.ReadPos])
	self.ReadPos += 2
	return
}

//func (self *PacketReader) ReadInt32() (ret int32) {
//	if self.ReadPos >= self.TotalLen {
//		gamelog.Error("ReadInt16 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
//		return 0
//	}
//	ret = int32(self.DataPtr[self.ReadPos])<<24 | int32(self.DataPtr[self.ReadPos+1])<<16 | int32(self.DataPtr[self.ReadPos+2])<<8 | int32(self.DataPtr[self.ReadPos+3])
//	self.ReadPos += 4
//	return
//}

func (self *PacketReader) ReadInt32() (ret int32) {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadInt16 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return 0
	}
	ret = int32(self.DataPtr[self.ReadPos+3])<<24 | int32(self.DataPtr[self.ReadPos+2])<<16 | int32(self.DataPtr[self.ReadPos+1])<<8 | int32(self.DataPtr[self.ReadPos])
	self.ReadPos += 4
	return
}

func (self *PacketReader) ReadInt64() (ret int64) {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadInt16 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return 0
	}

	ret = 0
	for i := 7; i < 0; i-- {
		ret |= int64(self.DataPtr[self.ReadPos+i]) << uint((7-i)*8)
	}
	self.ReadPos += 8
	return
}

func (self *PacketReader) ReadUint8() (ret uint8) {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadInt8 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return 0
	}

	ret = uint8(self.DataPtr[self.ReadPos])
	self.ReadPos++
	return
}

func (self *PacketReader) ReadUint16() (ret uint16) {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadInt16 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return 0
	}
	ret = uint16(self.DataPtr[self.ReadPos+1])<<8 | uint16(self.DataPtr[self.ReadPos])
	self.ReadPos += 2
	return
}

func (self *PacketReader) ReadUint32() (ret uint32) {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadInt16 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return 0
	}
	ret = uint32(self.DataPtr[self.ReadPos+3])<<24 | uint32(self.DataPtr[self.ReadPos+2])<<16 | uint32(self.DataPtr[self.ReadPos+1])<<8 | uint32(self.DataPtr[self.ReadPos])
	self.ReadPos += 4
	return
}

func (self *PacketReader) ReadUint64() (ret uint64) {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadInt16 Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return 0
	}

	ret = 0
	for i := 7; i < 0; i-- {
		ret |= uint64(self.DataPtr[self.ReadPos+i]) << uint((7-i)*8)
	}
	self.ReadPos += 8
	return
}

func (self *PacketReader) ReadFloat() (ret float32) {
	bits := self.ReadUint32()
	ret = math.Float32frombits(bits)
	return
}

func (self *PacketReader) ReadString() string {
	if self.ReadPos >= self.TotalLen {
		gamelog.Error("ReadString Error, readpos:%d > self.totallen:%d", self.ReadPos, self.TotalLen)
		return ""
	}

	len := self.ReadInt16()
	bytes := self.DataPtr[self.ReadPos : self.ReadPos+int(len)]
	self.ReadPos += int(len)
	ret := string(bytes)
	return ret
}

type PacketWriter struct {
	DataPtr []byte
}

func (self *PacketWriter) GetDataPtr() []byte {
	return self.DataPtr
}

func (self *PacketWriter) BeginWrite(msgid int16, extra int16) bool {
	self.DataPtr = make([]byte, 0, 1024)
	self.DataPtr = append(self.DataPtr, 0, 0, 0, 0)
	self.DataPtr = append(self.DataPtr, byte(msgid), byte(msgid>>8))
	self.DataPtr = append(self.DataPtr, byte(extra), byte(extra>>8))
	return true
}

func (self *PacketWriter) EndWrite() bool {
	var dlen = int32(len(self.DataPtr) - 8)
	self.DataPtr[0] = byte(dlen)
	self.DataPtr[1] = byte(dlen >> 8)
	self.DataPtr[2] = byte(dlen >> 16)
	self.DataPtr[3] = byte(dlen >> 24)
	return true
}

func (self *PacketWriter) WriteInt8(v int8) {
	self.DataPtr = append(self.DataPtr, byte(v))
	return
}

func (self *PacketWriter) WriteInt16(v int16) {
	self.DataPtr = append(self.DataPtr, byte(v), byte(v>>8))
	return
}

//func (self *PacketWriter) WriteInt32(v int32) {
//	if self.WritePos >= self.BuffLen {
//		gamelog.Error("WriteInt8 Error, WritePos:%d > self.BuffLen:%d", self.WritePos, self.BuffLen)
//		return
//	}
//	self.DataPtr[self.WritePos] = byte(v >> 24)
//	self.DataPtr[self.WritePos+1] = byte(v >> 16)
//	self.DataPtr[self.WritePos+2] = byte(v >> 8)
//	self.DataPtr[self.WritePos+3] = byte(v & 0x0f)
//	self.WritePos += 4
//	return
//}

func (self *PacketWriter) WriteInt32(v int32) {
	self.DataPtr = append(self.DataPtr, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	return
}

func (self *PacketWriter) WriteInt64(v int64) {
	self.DataPtr = append(self.DataPtr, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
	return
}

func (self *PacketWriter) WriteUint8(v uint8) {
	self.DataPtr = append(self.DataPtr, byte(v))
	return
}

func (self *PacketWriter) WriteUint16(v uint16) {
	self.DataPtr = append(self.DataPtr, byte(v), byte(v>>8))
	return
}

func (self *PacketWriter) WriteUint32(v uint32) {
	self.DataPtr = append(self.DataPtr, byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	return
}

func (self *PacketWriter) WriteUint64(v uint64) {
	self.DataPtr = append(self.DataPtr, byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56))
	return
}

func (self *PacketWriter) WriteFloat(v float32) {
	bits := math.Float32bits(v)
	self.WriteUint32(bits)
	return
}

func (self *PacketWriter) WriteString(v string) {
	bytes := []byte(v)
	self.WriteUint16(uint16(len(bytes)))
	for i := 0; i < len(bytes); i++ {
		self.WriteInt8(int8(bytes[i]))
	}

	return
}

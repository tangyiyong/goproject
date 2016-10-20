package utility

type BitMap struct {
	data    []byte
	bitsize int
	maxpos  int
	count   int //有效位的数量
}

// SetBit 将 offset 位置的 bit 置为 value (0/1)
func (this *BitMap) SetBit(offset int) bool {
	index, pos := offset/8, uint8(offset%8)
	if this.bitsize < offset {
		return false
	}

	this.data[index] |= 0x01 << (7 - pos)
	if this.maxpos < offset {
		this.maxpos = offset
	}

	return true
}

func (this *BitMap) ClrBit(offset int) bool {
	index, pos := offset/8, uint8(offset%8)

	if this.bitsize < offset {
		return false
	}

	this.data[index] &^= 0x01 << (7 - pos)

	return true
}

func (this *BitMap) GetBit(offset int) bool {
	index, pos := offset/8, uint8(offset%8)
	if this.bitsize < offset {
		return false
	}

	nRet := this.data[index] & (0x01 << (7 - pos))
	return nRet != 0
}

func (this *BitMap) Init(maxbit int) bool {
	this.data = make([]byte, maxbit/8+1)
	this.bitsize = maxbit
	this.maxpos = 0
	return true
}

func (this *BitMap) GetPositiveCnt() int {
	return this.count
}

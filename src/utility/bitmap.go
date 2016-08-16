package utility

import (
	"fmt"
)

type Bitmap struct {
	data    []byte
	bitsize int
	maxpos  int
}

// SetBit 将 offset 位置的 bit 置为 value (0/1)
func (this *Bitmap) SetBit(offset int, value int8) bool {
	index, pos := offset/8, uint8(offset%8)

	if this.bitsize < offset {
		return false
	}

	if value == 0 {
		this.data[index] &^= 0x01 << (7 - pos)
	} else {
		this.data[index] |= 0x01 << (7 - pos)
		if this.maxpos < offset {
			this.maxpos = offset
		}
	}
	return true
}

func (this *Bitmap) GetBit(offset int) bool {
	index, pos := offset/8, uint8(offset%8)
	if this.bitsize < offset {
		return false
	}

	nRet := this.data[index] & (0x01 << (7 - pos))
	return nRet != 0
}

func (this *Bitmap) Init(maxbit int) bool {
	this.data = make([]byte, maxbit/8+1)
	this.bitsize = maxbit
	this.maxpos = 0

	return true
}

func (this *Bitmap) Print() {

	for i := 0; i < len(this.data); i++ {
		fmt.Printf("%08b-", this.data[i])
	}

	return
}

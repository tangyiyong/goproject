package mainlogic

//import (
//	"gamelog"
//	"sync"
//)

type ST_Rect struct {
	left, right, top, bottom float32
}

func (self *ST_Rect) Contained(x, z float32) bool {
	if x < self.left || x > self.right || z < self.top || z > self.bottom {
		return false
	}
	return true
}

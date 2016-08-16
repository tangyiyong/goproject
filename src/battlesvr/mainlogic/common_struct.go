package mainlogic

//import (
//	"gamelog"
//	"sync"
//)

type ST_Rect struct {
	left, right, top, bottom float32
}

func (rect *ST_Rect) Contained(x, z float32) bool {
	if x < rect.left || x > rect.right || z < rect.top || z > rect.bottom {
		return false
	}
	return true
}

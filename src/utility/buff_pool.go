package utility

var G_BuffPool chan []byte //消息队列

func InitBuffPool(poolsize int) {
	G_BuffPool = make(chan []byte, poolsize)
}

func AllocBuff() (buff []byte) {
	select {
	case buff = <-G_BuffPool:
	default:
		buff = make([]byte, 0, 2048)
	}

	return
}

func ReleaseBuff(buff []byte) {
	select {
	case G_BuffPool <- buff:
	default:
		buff = nil
	}

	return
}

//type DB_Func func()

//var G_DB_FuncList chan DB_Func //函数队列

//func AddDBFunc(f DB_Func) {
//	select {
//	case G_DB_FuncList = <-G_BuffPool:
//	default:

//	}

//	return
//}

//go func(){
//	for dbfunc := range G_DB_FuncList {
//		dbfunc()
//	}
//}()

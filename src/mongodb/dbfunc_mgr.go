package mongodb

type DB_Func func()

var G_DB_FuncList chan DB_Func //函数队列

func AddDBFunc(f DB_Func) bool {
	G_DB_FuncList <- f
	return true
}

func InitDbFuncMgr() {
	G_DB_FuncList = make(chan DB_Func, 1024)
	go func() {
		for f := range G_DB_FuncList {
			f()
		}
	}()
}

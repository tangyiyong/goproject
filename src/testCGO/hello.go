package main

/*
#include <stdlib.h>
#include "hello.h"
struct TestStruct 
{
    int32 a;
    int32 b;
};
*/
import "C"

import (
    "fmt"
)

func main() {
	//C.sayHi()
	var text string = "testest"
	C.hello(C.CString(text))
	//C.printf(C.CString("%s"), C.CString("your house!!!"))
	i := C.rand()
    
    fmt.Println(i)
}

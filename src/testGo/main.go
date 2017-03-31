// main
package main

import (
	"fmt"
	"utility"
)

func test(buffer []byte, xorcode [4]byte) (bool, []byte) {
	for i := 0; i < 10; i++ {
		buffer[i] = buffer[i] + 10
	}

	return true, buffer[:5]
}

func main() {

	//TestMongoDB()
	//TestMarshal()
	//TestMutiArray()
	//TestOther2()

	if utility.RunOnlyOne() {
		fmt.Println("true-true-true-true-true-true-true-true")
	} else {
		fmt.Println("false-false-false-false-false-false-false-false")
	}

	var buf []byte = make([]byte, 10)
	var code [4]byte

	for i := 0; i < 10; i++ {
		buf[i] = byte(i)
	}

	_, pbuff := test(buf, code)

	fmt.Println(pbuff)

	utility.StartConsoleWait()
}

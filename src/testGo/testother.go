// main
package main

import (
	"fmt"
	"gamelog"
)

func TestOther() {
	gamelog.InitLogger("test", 0)
}

type speaker interface {
	say(v int)
}

type KKK struct {
	a int
}

func (self *KKK) say(v int) {
	fmt.Println("kkk say")
	self.a = v
}

type AAA struct {
	b int
}

func (self *AAA) say(v int) {
	fmt.Println("AAA say")
	self.b = v
}

var testmap map[int]speaker

func addinterface(id int, ptr speaker) {
	testmap[id] = ptr
}

func interfacesay() {
	for _, v := range testmap {
		v.say(2)
	}
}

func TestOther2() {
	gamelog.InitLogger("test", 0)

	testmap = make(map[int]speaker)

	var k KKK
	var a AAA

	addinterface(1, &k)
	addinterface(2, &a)

	interfacesay()

	fmt.Println(k)
	fmt.Println(a)
}

package main

import "fmt"

type heater interface {
	heat()
}

type cooler interface {
	cool()
}

type someThing struct {
	name string
}

func (s someThing) heat() {
	fmt.Println("I am heating with name = " + s.name)
}

func (s someThing) cool() {
	fmt.Println("I am cooling with name = " + s.name)
}

func main() {
	fmt.Printf("reveal interface assignment and assertion ......\n\n")

	some := someThing{name: "stone"}

	var h heater = some
	h.heat()

	var c cooler = some // NOTE: var c cooler = h, compile failed
	c.cool()

	var any interface{} = some // NOTE: var any interface{} = h, can be compiled and has the same result

	v0, ok0 := any.(someThing)
	if ok0 {
		fmt.Printf("any, empty interface, assert struct someThing, type = %T, val = %v\n", v0, v0)
	}

	v1, ok1 := any.(cooler)
	if ok1 {
		fmt.Printf("any, empty interface, assert interface coolor, type = %T, val = %v\n", v1, v1)
	}

	v2, ok2 := any.(heater)
	if ok2 {
		fmt.Printf("any, empty interface, assert interface heator, type = %T, val = %v\n", v2, v2)
	}

	// any = c, compile OK
	// h = c, compile fail
	v3, ok3 := h.(cooler)
	if ok3 {
		fmt.Printf("h, interface heator, assert interface coolor, type = %T, val = %v\n", v3, v3)
	}

	tryDuckerWithPointer()

	tryAssignInterface()
}

package main

import "fmt"

type heator interface {
	heat()
}

type coolor interface {
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
	some := someThing{"stone"}

	var h heator = some
	h.heat()

	var c coolor = some // NOTE: var c cooler = h, can not be compiled because interface can not be receiver
	c.cool()

	var any interface{} = h // NOTE: var any interface{} = some, can be compiled and has the same result

	v0, ok0 := any.(someThing)
	if ok0 {
		fmt.Printf("any has concrete someThing, type = %T, val = %v\n", v0, v0)
	}

	v1, ok1 := any.(coolor)
	if ok1 {
		fmt.Printf("any has interface coolor, type = %T, val = %v\n", v1, v1)
	}

	v2, ok2 := any.(heator)
	if ok2 {
		fmt.Printf("any has interface heator, type = %T, val = %v\n", v2, v2)
	}

	v3, ok3 := h.(coolor)
	if ok3 {
		fmt.Printf("heator has interface coolor, type = %T, val = %v\n", v3, v3)
	}
}

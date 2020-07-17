package main

import "fmt"

type duckor interface {
	quark()
}

type duckStuff struct {
	duckName string
}

func (d *duckStuff) quark() {
	fmt.Printf("pointer->quark(), %v\n", d.duckName)
}

func tryDuckorWithPointer() {
	fmt.Printf("\ntyrDuckorWithPointer ....\n\n")

	d := duckStuff{duckName: "Stone"}

	fmt.Println("assign duckStuff to any, which is empty interface")
	var any interface{} = d

	v0, ok0 := any.(duckStuff)
	if ok0 {
		fmt.Printf("any, empty interfae, assert struct, type = %T, val = %v\n", v0, v0)
	}

	v1, ok1 := any.(duckor)
	if ok1 {
		fmt.Printf("any, empty interface, assert interface duckor, type = %T, val = %v\n", v1, v1)
	} else {
		fmt.Printf("any, empty interface, fail to assert interface duckor\n")
	}

	fmt.Println("\nassign *duckStuff to any, which is empty interface")
	any = &d // NOTE: any the second reference is the pointer of struct
	v2, ok2 := any.(duckor)
	if ok2 {
		fmt.Printf("any, empty interface, assert interface duckor, type = %T, val = %v\n", v2, v2)

		// var i duckor = any   // illegal in compile
		var i duckor = v2
		fmt.Printf("assign v2 to duckor interface, type = %T\n", i)
	}

	var i duckor = &d
	p1 := &i
	var any1 interface{} = p1
	p2 := i
	var any2 interface{} = p2
	fmt.Printf("p1 type = %T, p2 type = %T, any1 = %T, any2 = %T\n", p1, p2, any1, any2)
}

type itfer interface {
	Dummy()
}

type foo struct{}

func (f foo) Dummy() {}

type bar struct{}

func (b *bar) Dummy() {}

func tryAssignInterface() {
	fmt.Printf("\ntryAssignInterface .......\n\n")

	var f1 foo
	var f2 *foo = &foo{}

	var i1 itfer = f1
	fmt.Printf("i1 from f1, type = %T\n", i1)
	var i2 itfer = f2
	fmt.Printf("i2 from f2, type = %T\n", i2)

	var b1 bar
	var b2 *bar = &b1

	// var i3 itfer = b1	// NOTE: compile fail
	var i4 itfer = b2
	fmt.Printf("i4 from b2, type = %T\n", i4)
}

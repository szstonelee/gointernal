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
	fmt.Println("\ntyrDuckorWithPointer ....")

	d := duckStuff{duckName: "Stone"}

	fmt.Println("\nassign duckStuff to any, which is empty interface")
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
}

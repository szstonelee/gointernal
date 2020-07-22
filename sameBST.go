package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Tree is binary sorted tree with root not be nil
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

// New returns a new, random binary tree holding the values k, 2k, ..., 10k.
func newTree(k int, n int) *Tree {
	if k <= 0 || n <= 0 {
		return nil
	}

	var t *Tree
	for _, v := range rand.Perm(n) {
		t = insert(t, (1+v)*k)
	}
	return t
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{nil, v, nil}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func walk(t *Tree, ch chan int) {
	if t.Left != nil {
		walk(t.Left, ch)
	}
	ch <- t.Value
	if t.Right != nil {
		walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func same(t1, t2 *Tree) bool {
	c1 := make(chan int, 1<<10)
	c2 := make(chan int, 1<<10)

	go func() {
		walk(t1, c1)
		close(c1)
	}()
	go func() {
		walk(t2, c2)
		close(c2)
	}()

	for {
		select {
		case v1, ok1 := <-c1:
			v2, ok2 := <-c2
			if !ok1 {
				return !ok2
			}
			if !ok2 {
				return false
			}

			if v1 != v2 {
				return false
			}

		case v2, ok2 := <-c2:
			v1, ok1 := <-c1
			if !ok2 {
				return !ok1
			}
			if !ok1 {
				return false
			}

			if v1 != v2 {
				return false
			}
		}
	}

}

func main() {
	start1 := time.Now()
	t1 := newTree(1, 1<<23)
	t2 := newTree(1, 1<<23)
	fmt.Println("Tree construct time = ", time.Since(start1))

	start2 := time.Now()
	same := same(t1, t2)
	fmt.Println("Same compute time = ", time.Since(start2))
	fmt.Println("same = ", same)
}

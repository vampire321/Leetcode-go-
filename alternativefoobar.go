package main

import (
	"fmt"
	"sync"
)

type FooBar struct {
	n       int
	fooTurn chan struct{} // foo starts with 1 permit
	barTurn chan struct{} // bar starts with 0 permits
}

func NewFooBar(n int) *FooBar {
	fooCh := make(chan struct{}, 1)
	fooCh <- struct{}{} // give initial turn to foo

	return &FooBar{
		n:       n,
		fooTurn: fooCh,
		barTurn: make(chan struct{}), // unbuffered → bar waits
	}
}

func (fb *FooBar) Foo(printFoo func()) {
	for i := 0; i < fb.n; i++ {
		<-fb.fooTurn          // acquire turn (blocks if not foo's turn)
		printFoo()            // Do not change this line
		fb.barTurn <- struct{}{} // give turn to bar
	}
}

func (fb *FooBar) Bar(printBar func()) {
	for i := 0; i < fb.n; i++ {
		<-fb.barTurn          // acquire turn (blocks until foo finishes)
		printBar()            // Do not change this line
		fb.fooTurn <- struct{}{} // give turn back to foo
	}
}

// ======================
// Simple test harness (for your local practice only)
// ======================
func main() {
	n := 5 // change this to test different values
	fb := NewFooBar(n)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fb.Foo(func() { fmt.Print("foo") })
	}()

	go func() {
		defer wg.Done()
		fb.Bar(func() { fmt.Print("bar") })
	}()

	wg.Wait()
	fmt.Println() // newline
}
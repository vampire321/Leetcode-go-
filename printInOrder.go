package main

import (
	"fmt"
	"sync"
)

type Foo struct {
	done1 chan struct{}
	done2 chan struct{}
}

func NewFoo() *Foo {
	return &Foo{
		done1: make(chan struct{}),
		done2: make(chan struct{}),
	}
}

func (f *Foo) First(printFirst func()) {
	printFirst()
	close(f.done1) // signal that first is done
}

func (f *Foo) Second(printSecond func()) {
	<-f.done1      // wait for first to complete
	printSecond()
	close(f.done2) // signal that second is done
}

func (f *Foo) Third(printThird func()) {
	<-f.done2      // wait for second to complete
	printThird()
}

// ======================
// Simple test harness (for your local practice only)
// ======================
func main() {
	foo := NewFoo()
	var wg sync.WaitGroup
	wg.Add(3)

	// Launch three goroutines in random order
	go func() {
		defer wg.Done()
		foo.Third(func() { fmt.Print("third") })
	}()

	go func() {
		defer wg.Done()
		foo.First(func() { fmt.Print("first") })
	}()

	go func() {
		defer wg.Done()
		foo.Second(func() { fmt.Print("second") })
	}()

	wg.Wait()
	fmt.Println() // newline at the end
}
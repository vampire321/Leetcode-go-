package main

import (
	"fmt"
	"sync"
)

type ZeroEvenOdd struct {
	n      int
	zeroCh chan struct{} // zero has initial permit
	evenCh chan struct{}
	oddCh  chan struct{}
}

func NewZeroEvenOdd(n int) *ZeroEvenOdd {
	z := &ZeroEvenOdd{
		n:      n,
		zeroCh: make(chan struct{}, 1),
		evenCh: make(chan struct{}),
		oddCh:  make(chan struct{}),
	}
	z.zeroCh <- struct{}{} // Give initial turn to zero
	return z
}

func (zeo *ZeroEvenOdd) Zero(printZero func(int)) {
	for i := 0; i < zeo.n; i++ {
		<-zeo.zeroCh                // Wait for my turn
		printZero(0)                // Do not change this line

		// Signal next: odd after even-indexed zero, even after odd-indexed zero
		if i%2 == 0 {
			zeo.oddCh <- struct{}{} // Next is odd (1,3,5...)
		} else {
			zeo.evenCh <- struct{}{} // Next is even (2,4,6...)
		}
	}
}

func (zeo *ZeroEvenOdd) Even(printEven func(int)) {
	for i := 2; i <= zeo.n; i += 2 {
		<-zeo.evenCh          // Wait for my turn
		printEven(i)          // Do not change this line
		zeo.zeroCh <- struct{}{} // Give turn back to zero
	}
}

func (zeo *ZeroEvenOdd) Odd(printOdd func(int)) {
	for i := 1; i <= zeo.n; i += 2 {
		<-zeo.oddCh           // Wait for my turn
		printOdd(i)           // Do not change this line
		zeo.zeroCh <- struct{}{} // Give turn back to zero
	}
}

// ======================
// Simple test harness (for your local practice only)
// ======================
func oevenOdd() {
	n := 5
	zeo := NewZeroEvenOdd(n)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		zeo.Zero(func(x int) { fmt.Printf("%d ", x) })
	}()

	go func() {
		defer wg.Done()
		zeo.Odd(func(x int) { fmt.Printf("%d ", x) })
	}()

	go func() {
		defer wg.Done()
		zeo.Even(func(x int) { fmt.Printf("%d ", x) })
	}()

	wg.Wait()
	fmt.Println()
}
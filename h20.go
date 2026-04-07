package main

/*import (
	"fmt"
	"sync"
)
*/
type H2O struct {
	hSem chan struct{}
	oSem chan struct{}
}

func NewH2O() *H2O {
	h := &H2O{
		hSem: make(chan struct{}, 2),
		oSem: make(chan struct{}),
	}

	// allow 2 hydrogens initially
	h.hSem <- struct{}{}
	h.hSem <- struct{}{}

	return h
}

func (h2o *H2O) Hydrogen(releaseHydrogen func()) {
	<-h2o.hSem
	releaseHydrogen()

	h2o.oSem <- struct{}{}
}

func (h2o *H2O) Oxygen(releaseOxygen func()) {
	<-h2o.oSem
	<-h2o.oSem

	releaseOxygen()

	// reset for next molecule
	h2o.hSem <- struct{}{}
	h2o.hSem <- struct{}{}
}
/* dataflow: Experimenting with data flow programming.
Author: James Fairbanks
Date: 2013-12-23
License : BSD
*/
package main

import (
	"fmt"
	. "github.com/jpfairbanks/streams/dataflow/datachannel"
	"math/rand"
	"os"
	"time"
)

var seed int64 = time.Now().UnixNano()
var delay time.Duration = 10 << 22

var ifPositive, ifNegative, unitFilter Filter

func init() {
	fmt.Printf("randseed: %d\n", seed)
	//ifPositive: Return true if x is positive
	ifPositive = MkStep(0)
	//ifNegative: Return true if x is negative
	ifNegative = func(x Data) Data {
		return 1 - ifPositive(x)
	}
	unitFilter = MkSquareFilter(0, 1)
}

func main() {
	var exitStatus int
	var inch, posch, negch Dchan
	inch = make(Dchan)
	posch = make(Dchan)
	negch = make(Dchan)
	//go spew_ints(2<<4, inch)
	go randwalk_gen(delay, 0.0, 1.0, inch)
	pospipe := Filterpipe{ifPositive, posch}
	negpipe := Filterpipe{ifNegative, negch}
	go Split(inch, pospipe, negpipe)
	go Print("pos: %f\n", os.Stdout, posch)
	go Print("neg:%f\n", os.Stdout, negch)
	//go Drain(negch)
	for {
		time.Sleep(10000)
	}
	os.Exit(exitStatus)
}

//randwalk_gen: Produces Gaussian a random walk with a fixed mean and variance on ch
func randwalk_gen(delay time.Duration, mean float64, vari float64, ch Dchan) {
	x := 0.0
	r := rand.New(rand.NewSource(seed))
	fmt.Printf("delay: %v\n", delay)
	rwalk := func() Data {
		time.Sleep(delay)
		x += (vari * r.NormFloat64()) + mean
		return Data(x)
	}
	Spew(rwalk, ch)
}

//spew_ints: Make ints uniformly in [0,max).
func spew_ints(max int, output Dchan) {
	r := rand.New(rand.NewSource(seed))
	f := func() Data {
		x := r.Intn(max)
		return Data(x)
	}
	Spew(f, output)
}

func collect(sigch Dchan) {
	var accum Data
	accum = 0.0
	count := 0
	for x := range sigch {
		count += 1
		accum += x
		fmt.Printf("passed: %v\n", count)
		fmt.Printf("mean: %v\n", accum/Data(count))
	}
}

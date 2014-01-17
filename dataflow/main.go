/* dataflow: Experimenting with data flow programming.
Author: James Fairbanks
Date: 2013-12-23
License : BSD
*/
package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var seed int64 = time.Now().UnixNano()
var delay time.Duration = 10 << 22

type data float64
type dchan chan data
type filter func(data) data

var ifPositive, ifNegative, unitFilter filter

func init() {
	fmt.Printf("randseed: %d", seed)
	//ifPositive: Return true if x is positive
	ifPositive = mk_step(0)
	//ifNegative: Return true if x is negative
	ifNegative = func(x data) data {
		return 1 - ifPositive(x)
	}
	unitFilter = mk_square_filter(0, 1)
}

func main() {
	var exit_status int
	fmt.Print("Hello World\n")
	var inch, posch, negch dchan
	inch = make(dchan)
	posch = make(dchan)
	negch = make(dchan)
	//go spew_ints(2<<4, inch)
	go randwalk_gen(0.0, 1.0, inch)
	pospipe := filterpipe{ifPositive, posch}
	negpipe := filterpipe{ifNegative, negch}
	go split(inch, pospipe, negpipe)
	go Printer("pos: %f\n", os.Stdout, posch)
	go Printer("neg: %f\n", os.Stdout, negch)
	//go Drain(negch)
	for {
		time.Sleep(10000)
	}
	os.Exit(exit_status)
}

//identity: The identity filter leaves the data unchanged acts as a noop.
func identity(x data) data {
	return x
}

//apply: binds a filter to and input channel and an output channel.
func apply(f filter, inchan dchan, outchan dchan) {
	for x := range inchan {
		outchan <- f(x)
	}
}

//tee: copies an input channel to a variable number of outputs.
func tee(inchan dchan, outchan ...dchan) {
	for x := range inchan {
		L := len(outchan)
		for i := 0; i < L; i++ {
			outchan[i] <- x
		}
	}
}

//filterpipe: a struct matching a filter to a channel for the output.
type filterpipe struct {
	filt filter
	pipe dchan
}

//split: Act like tee but each filterpipe transmits only the input for which
//filt returns 1. No constraint on the pipes partitioning the data.
//If all filters always return 1, then this behaves like tee.
//If there are two filters f(x) and 1-f(x), and f(x) is in {0,1}, then
//split will partition the data into two streams.
func split(inchan dchan, fpipes ...filterpipe) {
	L := len(fpipes)
	for x := range inchan {
		for i := 0; i < L; i++ {
			fp := fpipes[i]
			if fp.filt(x) == 1 {
				fp.pipe <- x
			}
		}
	}
}

//Drain: Just drain a channel by throwing away the data that comes through.
func Drain(ch dchan) {
	for _ = range ch {
		continue
	}
}

//mk_step: Returns a function that acts as a threshold filter, returning true if x > t.
func mk_step(t data) filter {
	return func(x data) data {
		if x > t {
			return 1
		} else {
			return 0
		}
	}
}

//mk_shift: Returns a function that shifts each data point by a constant.
func mk_shift(k data) filter {
	return func(x data) data {
		return x + k
	}
}

//mk_scale: Returns a function that scales each data point by a constant.
func mk_scale(k data) filter {
	return func(x data) data {
		return x * k
	}
}

//mk_square_filter: Returns a filter returning 1 if x in [lowerbound, upperbound] else zero
func mk_square_filter(lowerbound data, upperbound data) filter {
	return func(x data) data {
		var bigenough, smallenough data
		bigenough, smallenough = 0, 0
		if x >= lowerbound {
			bigenough = 1
		}
		if x <= upperbound {
			smallenough = 1
		}
		return bigenough * smallenough
	}
}

//spew: take a function with no arguments and spew the output to a channel
//good for generating random numbers or using a closure.
func spew(f func() data, ch dchan) {
	for {
		ch <- f()
	}
}

//spew_ints: Make ints uniformly in [0,max).
func spew_ints(max int, output dchan) {
	r := rand.New(rand.NewSource(seed))
	f := func() data {
		x := r.Intn(max)
		return data(x)
	}
	spew(f, output)
}

////randwalk_gen_naive: Produces Gaussian a random walk with a fixed mean and variance on ch
//func randwalk_gen(mean float64, vari float64, ch chan data) {
//x := 0.0
//r := rand.New(rand.NewSource(seed))
//fmt.Printf("delay: %v\n", delay)
//for {
//time.Sleep(delay)
//x += (vari * r.NormFloat64()) + mean
//ch <- data(x)
//}
//}

//randwalk_gen: Produces Gaussian a random walk with a fixed mean and variance on ch
func randwalk_gen(mean float64, vari float64, ch chan data) {
	x := 0.0
	r := rand.New(rand.NewSource(seed))
	fmt.Printf("delay: %v\n", delay)
	rwalk := func() data {
		time.Sleep(delay)
		x += (vari * r.NormFloat64()) + mean
		return data(x)
	}
	spew(rwalk, ch)
}

//Printer: Use a fmtstring to print out the data that comes through the channel
//and write the resulting string to the io.Writer.
func Printer(fmtstring string, fp io.Writer, input dchan) {
	for x := range input {
		s := fmt.Sprintf(fmtstring, x)
		fp.Write([]byte(s))
	}
}

func apply_test(f filter, test filter, inch dchan, outch dchan, sigch dchan) {
	var x data
	for {
		x = <-inch
		//fmt.Printf("app: %v\n", x)
		outch <- f(x)
		if test(x) == 1 {
			sigch <- x
		}
	}
}

func collect(sigch dchan) {
	var accum data
	accum = 0.0
	count := 0
	for x := range sigch {
		count += 1
		accum += x
		fmt.Printf("passed: %v\n", count)
		fmt.Printf("mean: %v\n", accum/data(count))
	}
}

/* datachannel.go: Functions for manipulating filters and streams.
Author: James Fairbanks
Date: 2014-01-24
License : BSD
*/
package datachannel

import (
	"fmt"
	"io"
)

/*
 * Type Defs
 */
type Data float64
type Dchan chan Data
type Filter func(Data) Data

//filterpipe: a struct matching a filter to a channel for the output.
type Filterpipe struct {
	Filt Filter
	Pipe Dchan
}
type FilterpipeI interface {
	Read() Data
	Write(Data)
	Filt() Filter
}

/*
 * General Operators
 */

//Identity: The identity filter leaves the Data unchanged acts as a noop.
func Identity(x Data) Data {
	return x
}

//Apply: binds a Filter to and input channel and an output channel.
func Apply(f Filter, inchan Dchan, outchan Dchan) {
	for x := range inchan {
		outchan <- f(x)
	}
}

//Spew: take a function with no arguments and spew the output to a channel
//good for generating random numbers or using a closure.
func Spew(f func() Data, ch Dchan) {
	for {
		ch <- f()
	}
}

//Tee: copies an input channel to a variable number of outputs.
func Tee(inchan Dchan, outchan ...Dchan) {
	for x := range inchan {
		L := len(outchan)
		for i := 0; i < L; i++ {
			outchan[i] <- x
		}
	}
}

//Split: Act like tee but each filterpipe transmits only the input for which
//filt returns 1. No constraint on the pipes partitioning the Data.
//If all filters always return 1, then this behaves like tee.
//If there are two filters f(x) and 1-f(x), and f(x) is in {0,1}, then
//split will partition the Data into two streams.
func Split(inchan Dchan, fpipes ...Filterpipe) {
	L := len(fpipes)
	for x := range inchan {
		for i := 0; i < L; i++ {
			fp := fpipes[i]
			if fp.Filt(x) == 1 {
				fp.Pipe <- x
			}
		}
	}
}

//Drain: Just drain a channel by throwing away the Data that comes through.
// similar to  writing to /dev/null
func Drain(ch Dchan) {
	for _ = range ch {
		continue
	}
}

//Print: Use a fmtstring to print out the Data that comes through the channel
//and write the resulting string to the io.Writer.
func Print(fmtstring string, fp io.Writer, input Dchan) {
	for x := range input {
		s := fmt.Sprintf(fmtstring, x)
		fp.Write([]byte(s))
	}
}

package utils

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"time"
)

//CatFloat: Scan an IO stream for floats and put them into a channel
//log the errors to stderr and close the channel when the reader is done.
//You can check for the closed channel to know when to stop receiving.
func CatFloat(fp *bufio.Scanner, ch chan float64) error {
	var line string
	var number float64
	var casterr error
	for fp.Scan() {
		line = fp.Text()
		number, casterr = strconv.ParseFloat(line, 64)
		if casterr != nil {
			log.Println(casterr)
			log.Println(line)
		} else {
			ch <- number
		}
	}
	//we only need to check for a read error when we stop
	//because errors cause stoppage
	err := fp.Err()
	if err != nil {
		log.Print(err)
	}
	close(ch)
	return err
}

func TeeFloat(inchan chan float64, outchan1 chan float64, outchan2 chan float64) {
	for x := range inchan {
		outchan1 <- x
		outchan2 <- x
	}
}

func Zip(inchan []chan float64, outchan chan []float64) {
	closed := false
	numchans := len(inchan)
	buff := make([]float64, numchans)
	for !closed {
		for i := 0; i < numchans; i++ {
			buff[i] = <-inchan[i]
		}
		outchan <- buff
	}
}

func ReduceFloat(inchan chan []float64, f func([]float64) float64, outchan chan float64) {
	for buff := range inchan {
		outchan <- f(buff)
	}
}

/*
type FloatStream struct{
	ch chan float64
	fp io.Writer
}
*/

//Print: Start a print server that receives on ch and writes to Stdout
//Does not stop. You must call this in a go routine.
func Print(fp io.Writer, ch chan string) {
	var s string
	for {
		select {
		case s = <-ch:
			fp.Write([]byte(s))
			//outfile.Write(s)
		}
	}
}

//SkipPrint: just like Print but takes an interval prints every ith element.
//Anyone who can write to the channel can advance the counter.
func SkipPrint(fp io.Writer, interval int64, ch chan string) {
	var s string
	var state int64
	for {
		select {
		case s = <-ch:
			if state%interval == 0 {
				fp.Write([]byte(s))
			}
			state += 1
		}
	}
}

//PeriodicQuery: Runs a query function at periodic intervals and reports the result on a channel
//Will run forever. Combine with Print to have Duration based Print Server.
func PeriodicQuery(ts time.Duration, query func() string, ansch chan string) {
	tickch := time.Tick(ts)
	for _ = range tickch {
		answer := query()
		ansch <- answer
	}
}

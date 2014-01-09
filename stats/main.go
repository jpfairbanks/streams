/* main.go: streamstats is a package for streaming stats on the command line
Author: James Fairbanks
Date: 2013-08-30
Liscence: BSD
*/

package main

import "os"
import "log"
import "bufio"
import "github.com/jpfairbanks/estimation"
import "github.com/jpfairbanks/stream/utils"

func main() {
	var est estimation.Estimator
	est = new(estimation.TwomomentsChan)
	var stdin *bufio.Scanner
	stdin = bufio.NewScanner(os.Stdin)
	log.Print("Starting to read\n")
	//setup file reader
	inchan := make(chan float64)
	go utils.CatFloat(stdin, inchan)
	//output print server
	outchan := make(chan string)
	go utils.Print(os.Stdout, outchan)
	//main loop
	var number float64
	var ok bool
	for {
		select {
		case number, ok = <-inchan:
			if ok {
				est.Push(number)
			}
		}

		if !ok {
			//the channel has closed no new data is coming
			break
		}
		outchan <- est.String()
	}
}

//Spawn: Start a loop where we read off the in channel update the state and
//output the estimates on the outchan.
func Spawn(est estimation.Estimator,
	inchan chan float64, outchan chan string) {
	var number float64
	var ok bool
	for {
		select {
		case number, ok = <-inchan:
			if ok {
				est.Push(number)
			}
		}
		if !ok {
			//the channel has closed no new data is coming
			break
		}
		outchan <- est.String()
	}
}

type dtype float64

//Worker: An interface for a worker. We do something when we get a datum,
//and we do something before we write.
type Worker interface {
	PostRead(x dtype)
	PreWrite() dtype
}

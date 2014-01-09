package main

import (
	//	"fmt"
	"bufio"
	"github.com/jpfairbanks/estimation"
	"github.com/jpfairbanks/stream/utils"
	"log"
	"os"
	"testing"
	"time"
)

var datafile string = "normal_data.txt"

//getInput: Take a filename and return a channel of the floats contained in the file.
func getInput(filename string) (chan float64, error) {
	var fp *bufio.Scanner
	file, err := os.Open(filename)
	if err != nil {
		log.Println(file)
	}
	fp = bufio.NewScanner(file)
	log.Printf("Making input channel from file\n")
	inchan := make(chan float64)
	go utils.CatFloat(fp, inchan)
	return inchan, err
}

//TestGetInput: Make sure that we can take a file and turn it into a channel.
func TestGetInput(t *testing.T) {
	inchan, err := getInput(datafile)
	if err != nil {
		t.Error(err.Error())
	}
	log.Println("Made channel")
	for i := range inchan {
		t.Log(i)
	}
	log.Println("Drained channel")
}

//TestPrintMain: Run a test showing consume a stream of numbers
//Estimating the distribution and and printing out the estimate at regular intervals of time
func TestPrintMain(t *testing.T) {
	var est estimation.Estimator
	est = new(estimation.Twomoments)

	inchan, err := getInput(datafile)
	if err != nil {
		t.Error(err.Error())
	}
	//output print server
	outchan := make(chan string)
	go utils.Print(os.Stdout, outchan)
	//main loop
	Spawn(est, inchan, outchan)
}

//TestRange: test the estimator.Range class
func TestRange(t *testing.T) {
	var est estimation.Estimator
	est = estimation.NewRange()

	inchan, err := getInput(datafile)
	if err != nil {
		t.Error(err.Error())
	}
	//output print server
	outchan := make(chan string)
	go utils.Print(os.Stdout, outchan)
	//main loop
	Spawn(est, inchan, outchan)
}

//TestPeriodicQueryMain: Run a test showing consume a stream of numbers,
//estimating the distribution and and printing out the estimate at regular intervals of time.
func TestPeriodicQueryMain(t *testing.T) {
	var est estimation.Estimator
	est = new(estimation.Twomoments)
	inchan, err := getInput(datafile)
	if err != nil {
		t.Error(err.Error())
	}
	timedch := make(chan string)
	go utils.PeriodicQuery(time.Nanosecond,
		func() string { return est.String() },
		timedch)
	outchan := make(chan string)
	go utils.Print(os.Stdout, outchan)
	var number float64
	var ok bool
	var report string
	for {
		select {
		case number, ok = <-inchan:
			if ok {
				est.Push(number)
			}
		case report, ok = <-timedch:
			if ok {
				outchan <- report
			}
			if !ok {
				t.Logf("not ok in report\n")
			}
		}

		if !ok {
			break
		}
	}
}

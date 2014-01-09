/* gendata.go: take args and print out data on stdout
Author: James Fairbanks
Date: 2013-08-30
Liscence: BSD
*/
package main

import "fmt"
import "log"
import "math/rand"
import "flag"
import "time"
//import "bufio"

// logFlag : print a Flag pointer to stderr
func logFlag(fl *flag.Flag) {
	log.Println(fl)
}

var countp = flag.Int("n",10, "number of elements to produce")
var seedp  = flag.Int64("seed", 0, "A seed for the random number generator for deterministic testing")
var usetime  = flag.Bool("time", false, "Use the time as a seed, dominates --seed")
var dist = flag.String("distribution", "uniform", "type of variable to generate: uniform, normal, exponential")

func main() {
	var seed int64
	var rsource *rand.Rand
	var randfunc func() float64

	//handle command line flags
	flag.Parse()
	//flag.VisitAll(logFlag)

	//set up the Random number generator
	if *usetime {
		seed = time.Now().Unix()
	} else {
		seed = *seedp
	}
	rsource = rand.New(rand.NewSource(seed))

	// pick a distribution to use
	switch {
		case *dist == "exponential": randfunc = rsource.ExpFloat64
		case *dist == "normal": randfunc = rsource.NormFloat64
		default: randfunc = rsource.Float64
	}

	log.Printf("Starting to generate data\n")
	//shove data out the door
	for i:=0; i < *countp; i++{
		fmt.Println(randfunc())
	}
}

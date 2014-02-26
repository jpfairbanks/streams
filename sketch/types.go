package sketch

import (
	"fmt"

//"math/rand"
)

//Number: used for the values that go into arrays.
type Number int

//Position: used to index arrays and hash functions.
type Position int

type Vector []Number

//add: Elementwise addition, panics if dimensions do not match.
func (v Vector) add(w Vector) Vector {
	sum := make(Vector, len(v))
	if len(v) != len(w) {
		panic(fmt.Sprintf("Dimensions do not match %d != %d", len(v), len(w)))
	}
	for i := 0; i < len(v); i++ {
		sum[i] = v[i] + w[i]
	}
	return sum
}

type Hash interface {
	Apply(Position) Position
}

type NHash interface {
	Apply(Position) Number
}

type Datum struct {
	index Position
	c     Number
}

type Query struct {
	index  Position
	result Number
}

//Sketch: Base interface for a sketch data structure.
//An important thing to note is the constraints on the sketches that are necessary
//for the Combine method to work.
//In most cases they need to have the same random seeds or hash functions,
//and the same parameters for dimension and error rates.
type Sketch interface {
	Insert(Datum) error
	Query(Query) error
	Combine(Sketch) error
}

//CombinationError: Useless interface type
type CombinationError string

func (e CombinationError) Error() string {
	return string(e)
}

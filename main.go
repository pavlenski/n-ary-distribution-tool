package main

import (
	"fmt"
	"strings"
)

var arity = 2

func main() {
	input := "desi brate? kako si danas ti meni? nadam se da si dobro. inace, voleo bih da se druzimo nekad. nedostajes mi!"
	a := strings.FieldsFunc(input, split)
	for i := 0; i < len(a)-arity; i++ {
		fmt.Println(a[i : i+arity])
	}

	// the key of the map for the bag of words could be a string
	// the checking of its existence can either be done with contains (checking if the words of the key are the same)
	// or the words being

	// log2arity - sorting the words to check the key (using O(log2n) sort)
	// arity^2 - contains
}

func split(r rune) bool {
	return r == '.' || r == '!' || r == ',' || r == '?'
}

package main

import (
	"fmt"
	"strings"
)

func main() {
	input := "desi brate? kako si danas ti meni? nadam se da si dobro. inace, voleo bih da se druzimo nekad. nedostajes mi!"
	a := strings.FieldsFunc(input, split)
	for i := 0; i < len(a); i++ {
		fmt.Println(a[i])
	}
}

func split(r rune) bool {
	return r == '.' || r == '!' || r == ',' || r == '?'
}

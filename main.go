package main

import (
	"fmt"
)

func main() {
	scanner.src = []byte("for (;;) {}")
	tokens := Tokenize()
	fmt.Println(tokens)
	// fmt.Println(Compile())
}

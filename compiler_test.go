package main

import (
	"testing"
)

func TestCompile(t *testing.T) {
	src := `func func1() {}
        // Comment 1
        func func2() {
            // Comment 2
        }
        func func3() {
            for(int i = 0; i < 10; i++) {
                // Comment 3
            }
        }
        func func4(int a, double b) int {
            return a
        }
    `
	Compile([]byte(src))
}

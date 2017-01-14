package main

import (
	"io/ioutil"
	"testing"
)

func TestCompile(t *testing.T) {
	fi := "testdata/input.txt"
	bs, err := ioutil.ReadFile(fi)
	if err != nil {
		t.Fatalf("opening %q: %v", fi, err)
	}

	// TODO compare with output.txt
	c := Compiler{}
	// c.Init(input, output)
	c.Compile(bs)
}

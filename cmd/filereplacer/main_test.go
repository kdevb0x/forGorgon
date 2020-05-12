package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	bak1, err := ioutil.ReadFile("testdata/testroot/rootsub1/blah.txt")
	if err != nil {
		panic(err)
	}
	bak2, err := ioutil.ReadFile("testdata/testroot/bleep.txt")
	if err != nil {
		panic(err)
	}

	var targetPath1 = "testdata/testroot/rootsub1/blah.txt"
	var targetPath2 = "testdata/testroot/bleep.txt"
	var replacePath1 = "testdata/replaceRoot/blah.txt"
	var replacePath2 = "testdata/replaceRoot/another/bleep.txt"

	size1 := len(bak)

	size2 := len(bak2)

	cleanup := func() {
		_ := os.Remove(targetPath1)
		_ := os.Remove(targetPath2)
		_, _ := ioutil.WriteFile(targetPath1, bak1, os.ModePerm)
		_, _ := ioutil.WriteFile(targetPath2, bak2, os.ModePerm)
	}
	targetRoot = "testdata/testroot"
	replacementRoot = "testdata/replaceRoot"

	t1, err := os.Stat(targetPath1)
	if err != nil {
		t.Fatalf("%w", err)
	}
	t2, err := os.Stat(replacePath1)
	if err != nil {
		t.Fatalf("%w", err)
	}
	if t1.Size() != t2.Size() {
		t.Fail()
	}
	t.Cleanup(cleanup)
}

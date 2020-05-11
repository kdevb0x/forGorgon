package main

import (
	"fmt"
	"os"
	"testing"
)

var f = newFileIdx(10)

func TestSet(t *testing.T) {

}

func TestSearchFiles(t *testing.T) {
	var orig = []string{"testdata/testroot/rootsub1/blah.txt", "testdata/testroot/rootsub3/boo.txt", "testdata/testroot/bleep.txt"}
	cwd, err := os.Getwd()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	var testroot = cwd + string(os.PathSeparator) + "testdata"

	newidx, err := searchFiles(testroot, f)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	for k, v := range newidx.f {
		fmt.Printf("filename: %s, record: %#v\n", k, v)
		if !contains(orig, k) {
			t.Fail()
		}
	}
}

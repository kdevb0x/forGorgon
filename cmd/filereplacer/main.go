package main

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/spf13/pflag"
)

var (
	srcroot  string
	destroot string
)

func flags() {
	srcroot = *pflag.StringP("root", "r", "/", "root dir containing files to change")
	destroot = *pflag.StringP("dest", "d", "$PWD", "destination to root created files")
}

type result struct {
	// hit on result
	found bool

	// abspath of the hit
	abspath string

	origSize int64

	// the target that matched this result
	matched string
}

func (r result) replace(with string, newname string) error {
	if r.found {
		return errors.New("no target to replace; result() can only be called after r.found == true")
	}
	var b bytes.Buffer
	f, err := ioutil.ReadFile(r.abspath)
}

type bakWorker struct {
	root string
	// filenames to compare
	targets []string
	resultQ chan result
}

func worker(root string, targets []string, resultChan chan result) {

}
func main() {
	flags()
}

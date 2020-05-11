package main

import (
	"errors"
	"io/ioutil"
	"os"

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

	newSize int64

	// the target that matched this result
	matched string

	// expired will be true after r.replace is called and returned nil error
	expired bool
}

func (r result) replace(with string, newname string) error {
	if !r.found {
		return errors.New("no target to replace; result() can only be called after r.found == true")
	}
	f, err := ioutil.ReadFile(with)
	if err != nil {
		return err
	}

	if r.origSize > 0 { // why this check? (I forget)
		r.newSize = int64(len(f))

	}
	err = os.Rename(r.abspath, r.abspath+".bak")
	if err != nil {
		return err
	}

	/* n, err := os.Create(r.abspath)
	if err != nil {
		return err
	}
	*/

	if err := ioutil.WriteFile(r.abspath+"/"+newname, f, os.ModePerm); err != nil {
		return err
	}
	r.expired = true
	return nil

}

type bakWorker struct {
	root string
	// filenames to compare
	targets []string
	resultQ chan result
}

func spawnWorker(root string, targets []string, resultChan chan result) {

}

type idx struct {
	// filename
	k string

	// orig path
	o string

	// new path

	n string
}

type fileidx struct {
	f map[string][]string

	pool []bakWorker

	setq chan idx
}

func newFileIdx(size int) *fileidx {
	var i = &fileidx{f: make(map[string][]string), pool: make([]bakWorker, size), setq: make(chan idx, size)}
	go func() {
		for r := range i.setq {
			i.f[r.k] = []string{r.o, r.n}
		}
	}()
	return i
}

// Set concurrently sets a key in f, overwriting if it already exists.
// Set returns a channel that the value to set should be sent down.
func (f fileidx) Set() chan idx {
	var c = make(chan idx, 1)
	go func() {
		for range c {
			f.setq <- <-c
		}
	}()
	return c
}

func searchFiles(rootdir string, targets fileidx) (*fileidx, error) {
	subs, err := ioutil.ReadDir(rootdir)
	if err != nil {
		return nil, err
	}
	for _, s := range subs {
		if s.IsDir() {
			// go spawnWorker(s, target)
		}
	}

}

func main() {
	flags()
}

package main

import "context"

type idx struct {
	// filename
	filename string

	// orig path
	origPath string

	// new path

	newPath string
}

type fileidx struct {
	f map[string]idx

	pool []bakWorker

	// the workqueue for concurrent map writes
	setq chan idx
}

func newFileIdx(size int) *fileidx {
	var i = &fileidx{f: make(map[string]idx), pool: make([]bakWorker, size), setq: make(chan idx, size)}
	go func() {
		// the set loop
		go func() {
			for j := range i.setq {
				i.f[j.filename] = j
			}
		}()

		// the getloop
		go func() {

		}()
	}()
	return i
}

// Set concurrently sets a key in f, overwriting if it already exists.
// Set returns a channel that the value to set should be sent down.
func (f *fileidx) Set() chan idx {
	var c = make(chan idx, 1)
	go func() {
		for range c {
			f.setq <- <-c
		}
	}()
	return c
}

func (f fileidx) Get(key string) (chan idx, context.Context) {

}

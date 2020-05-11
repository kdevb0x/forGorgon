package main

// idx is a record of a file.
type idx struct {
	// filename
	filename string

	// orig path
	origPath string

	// new path

	newPath string
}

// fileidx is a map of filenames to idx
type fileidx struct {
	f map[string]idx

	pool []bakWorker

	// the workqueue for concurrent map writes
	// only need a set queue, reads can be done without locking.
	setq chan idx
}

func newFileIdx(queueSize int) *fileidx {
	var i = &fileidx{f: make(map[string]idx), pool: make([]bakWorker, 5, 10), setq: make(chan idx, queueSize)}
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

// shouldnt need this, since the changes are done sequentially.
/*
func (f fileidx) Get(key string) (chan result, context.Context) {
	if i, found := f.f[key]; found {

	}

}
*/

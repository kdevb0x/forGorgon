package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/pflag"
)

var (
	overwrite bool
	srcroot   string
	destroot  string
	args      []string
)

func flags() {
	srcroot = *pflag.StringP("root", "r", "/", "root dir containing files to change")
	destroot = *pflag.StringP("dest", "d", "$PWD", "destination to root created files")
	overwrite = *pflag.BoolP("overwrite", "f", false, "force overwriting old file without backing up (saves backup as origname.bak by default)")
	pflag.Parse()
	args = pflag.Args()
}

type result struct {
	// hit on result
	found bool

	idx *idx
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
	// save a copy of the original
	err = os.Rename(r.abspath, r.abspath+".bak")
	if err != nil {
		return err
	}

	// write the new file in the old ones place
	if err := ioutil.WriteFile(r.abspath+"/"+newname, f, os.ModePerm); err != nil {
		return err
	}

	r.expired = true

	if overwrite {
		// remove the backup
		return os.Remove(r.abspath + ".bak")
	}
	return nil

}

type bakWorker struct {
	root string
	// filenames to compare
	targets []string
	resultQ chan result
}

func spawnWorker(root string, targets *fileidx, resultChan chan result) {

}

func searchFiles(rootdir string, targets *fileidx) (*fileidx, error) {
	subs, err := ioutil.ReadDir(rootdir)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	var errchan = make(chan error)

	for _, s := range subs {
		if s.IsDir() {
			// go spawnWorker(s, target)
			wg.Add(1)
			go func() {
				_, err := searchFiles(filepath.Join(rootdir, s.Name()), targets)
				if err != nil {
					log.Println(err)
					errchan <- err
				}
				wg.Done()

			}()
		}
		c := targets.Set()
		c <- newIdx(s.Name(), filepath.Join(rootdir, s.Name()))
		close(c)
	}
	wg.Wait()
	return targets, nil
}

func main() {
	flags()
}

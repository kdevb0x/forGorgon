package main

import (
	"errors"
	"fmt"
	"io"
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

	idx idx
	// abspath of the hit
	abspath string

	origSize int64

	newSize int64

	// the target that matched this result
	matched string

	// expired will be true after r.replace is called and returned nil error
	expired bool
}

func newResult(idx idx) *result {
	return &result{idx: idx}
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

func matchTargets(targets fileidx) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if res, found := targets.Get(info.Name()); found {
				// err := res.replace(res.matched, info.Name())
				err := simpleReplace(filepath.Join(path, info.Name()), res.abspath, !overwrite)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
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
				if err := os.Chdir(filepath.Join(rootdir, s.Name())); err != nil {
					errchan <- err
					return
				}
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

func contains(s []string, target string) bool {
	for _, m := range s {
		if m == target {
			return true
		}
	}
	return false
}

// simpleReplace replaces the file 'target', with the file 'with'.
// Both target and with should be absolute paths.
func simpleReplace(target string, with string, saveBackup bool) error {
	if saveBackup {
		err := os.Rename(target, target+".bak")
		if err != nil {
			return err
		}
	}
	c, err := os.Open(with)
	if err != nil {
		return err
	}
	defer c.Close()
	cinfo, err := c.Stat()
	if err != nil {
		return err
	}
	csize := cinfo.Size()

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Truncate(csize)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, c)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flags()
	ix, err := searchFiles(srcroot, newFileIdx(10))
	if err != nil {
		panic(err)
	}
	for k, v := range ix.f {
		fmt.Println(k)
		fmt.Println(v)
	}
}

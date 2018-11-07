package glc

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"path/filepath"
	"bufio"
	"compress/gzip"
)

// GLC define the glog cleaner options:
//
//     path     - Log files will be clean to this directory
//     prefix   - Log files name prefix
//     interval - Log files clean scanning interval
//     reserve  - Log files reserve time
//
type GLC struct {
	path     string
	prefix   string
	interval time.Duration
	reserve  time.Duration
}

// InitOption define the glog cleaner init options for glc:
//
//     Path     - Log files will be clean to this directory
//     Prefix   - Log files name prefix
//     Interval - Log files clean scanning interval
//     Reserve  - Log files reserve time
//
type InitOption struct {
	Path     string
	Prefix   string
	Interval time.Duration
	Reserve  time.Duration
}

// NewGLC create a cleaner in a goroutine and do instantiation GLC by given
// init options.
func NewGLC(option InitOption) *GLC {
	c := new(GLC)
	c.path = option.Path
	c.interval = option.Interval
	c.prefix = option.Prefix
	c.reserve = option.Reserve

	go c.cleaner()
	return c
}

// Create gz files
func (c *GLC) createGZ(src, dest string) error {
	sf, err := os.Open(filepath.Join(c.path, src))
	if err != nil {
		return err
	}
	defer sf.Close()

	sfb := bufio.NewReader(sf)

	df, err := os.OpenFile(filepath.Join(c.path, dest), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil{
		return err
	}
	defer df.Close()

	gw := gzip.NewWriter(df)
	defer gw.Close()

	gwb := bufio.NewWriter(gw)
	defer gwb.Flush()

	if _, err = gwb.ReadFrom(sfb); err != nil {
		return err
	}

	return nil
}

// clean provides function to check path exists by given log files path.
func (c *GLC) clean() {
	exists, err := c.exists(c.path)
	if err != nil {
		glog.Errorln(err)
		return
	}
	if !exists {
		return
	}

	files, err := ioutil.ReadDir(c.path)
	if err != nil {
		glog.Errorln(err)
		return
	}
	c.check(files)
}

// exists returns whether the given file or directory exists or not
func (c *GLC) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Get read symlink files
func (c *GLC) getRealFile(files []os.FileInfo) (str string) {
	for _, f := range files {
		// Skip not symlink file
		if f.Mode()&os.ModeSymlink == 0 {
			continue
		}

		// Get real file name
		if s, e := filepath.EvalSymlinks(filepath.Join(c.path, f.Name())); e != nil {
			continue
		} else {
			str += filepath.Base(s)
		}

	}

	return
}

// check provides function to check log files name whether the deletion
// condition is satisfied.
func (c *GLC) check(files []os.FileInfo) {
	rf := c.getRealFile(files)

	for _, f := range files {
		// Skip directory
		if f.IsDir() {
			continue
		}

		// Skip not has prefix string
		if ! strings.HasPrefix(f.Name(), c.prefix) {
			continue
		}

		// Skip symlink files
		if f.Mode()&os.ModeSymlink != 0 {
			continue
		}

		// Skip symlink real files
		if strings.Contains(rf, f.Name()) {
			continue
		}

		// Delete old files
		if time.Since(f.ModTime()) > c.reserve {
			if err := os.Remove(filepath.Join(c.path, f.Name())); err != nil {
				glog.Error(err)
			}
		} else if time.Since(f.ModTime()) > time.Duration(time.Minute*3) && (! strings.HasSuffix(f.Name(), ".gz")) {
			if err := c.createGZ(f.Name(), f.Name()+".gz"); err != nil { // Compress log files
				glog.Error(err)
			} else {
				if err := os.Remove(filepath.Join(c.path, f.Name())); err != nil {
					glog.Error(err)
				}
			}
		}
	}
}

// drop check the log file creation time and delete the file if the conditions
// are met.
func (c *GLC) drop(f os.FileInfo) {
	if time.Since(f.ModTime()) > c.reserve {
		err := os.Remove(c.path + f.Name())
		if err != nil {
			glog.Errorln(err)
		}
	}
}

// cleaner provides regular cleaning function by given log files clean
// scanning interval.
func (c *GLC) cleaner() {
	for {
		c.clean()
		time.Sleep(c.interval)
	}
}

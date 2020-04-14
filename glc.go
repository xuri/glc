package glc

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
)

type severity int32 // sync/atomic int32

// identical to glog log severities
const (
	infoLog severity = iota
	warningLog
	errorLog
	fatalLog
	numSeverity = 4
)

var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

var empty = struct{}{}

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
	symlinks map[string]struct{}
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
	c.symlinks = make(map[string]struct{}, numSeverity)

	for i := 0; i < numSeverity; i++ {
		symlink := c.prefix + "." + severityName[i]
		c.symlinks[symlink] = empty
	}

	go c.cleaner()
	return c
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

// check provides function to check log files name whether the deletion
// condition is satisfied.
func (c *GLC) check(files []os.FileInfo) {
	excludes := make(map[string]struct{}, numSeverity)
	for _, f := range files {
		// skip the destination of symlink files
		if _, ok := c.symlinks[f.Name()]; ok && f.Mode()&os.ModeSymlink != 0 {
			if dst, err := os.Readlink(c.path + f.Name()); err == nil {
				excludes[dst] = empty
			}
		}
	}

	for _, f := range files {
		if _, ok := excludes[f.Name()]; ok {
			continue
		}

		prefix := strings.HasPrefix(f.Name(), c.prefix)
		str := strings.Split(f.Name(), `.`)
		if prefix && len(str) == 7 && str[3] == `log` {
			c.drop(f)
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

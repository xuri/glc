package glc

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGLC(t *testing.T) {
	files := []string{"glc.localhost.xuri.log.WARNING.20180312-144710.3877", "glc.localhost.xuri.log.WARNING.20180312-144710"}
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Error(err)
		return
	}
	path += `/`

	go func() {
		glc := NewGLC(InitOption{
			Path:     path,
			Prefix:   `glc`,
			Interval: time.Duration(time.Second),
			Reserve:  time.Duration(time.Second * 3),
		})
		glc.exists(files[0])
	}()

	for _, file := range files {
		fp, err := os.OpenFile(path+file, os.O_CREATE|os.O_RDWR, 0700)
		if err != nil {
			fp.Close()
			t.Error(err)
			continue
		}
		_, err = fp.WriteAt([]byte{0}, 10)
		if err != nil {
			fp.Close()
			t.Error(err)
			continue
		}
		fp.Close()
		time.Sleep(time.Second * 5)
	}
}

func TestBadPath(t *testing.T) {
	path := []string{"", "/usr/bin/nohup"}
	for _, p := range path {
		go func(p string) {
			NewGLC(InitOption{
				Path:     p,
				Prefix:   `glc`,
				Interval: time.Duration(time.Second),
				Reserve:  time.Duration(time.Second * 3),
			})
		}(p)
	}
	time.Sleep(time.Second * 5)
}

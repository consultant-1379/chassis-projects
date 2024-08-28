package utils

import (
	"errors"
	"fmt"
	//"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"path/filepath"
)

var (
	fileRepo = make(map[string]FileInfo)
)

type FsNotify struct {
	watcher  *fsnotify.Watcher
	fileRepo map[string]FileInfo
	done     chan bool
	sync.Mutex
}

type FileInfo interface {
	GetFileName() string
	Handler(name, op string)
}

func InitFsWatcher() *FsNotify {
	fs := &FsNotify{
		fileRepo: make(map[string]FileInfo),
		done:     make(chan bool),
	}
	var err error
	fs.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case event := <-fs.watcher.Events:
				fmt.Println("fsnotify event ", event)
				//dir, _ := path.Split(event.Name)
				dir, _ := filepath.Abs(filepath.Dir(event.Name))
				if f, ok := fs.fileRepo[dir]; ok {
					f.Handler(event.Name, event.Op.String())
				} else {
					fmt.Println("Can not find related filerepo ", dir)
				}

			case err := <-fs.watcher.Errors:
				fmt.Println("fsnotify error ", err)
			case <-fs.done:
				fmt.Println("fsnotify is exiting")
				return
			}
		}
	}()

	return fs
}

func (fs *FsNotify) AddFileToFsWatcher(f FileInfo) error {
	name := f.GetFileName()

	//if strings.HasSuffix(name, "/") == false {
	//	name = name + "/"
	//}

	if _, ok := fs.fileRepo[name]; ok {
		return errors.New("Has already existed")
	}

	if err := fs.watcher.Add(name); err != nil {
		return err
	}

	fs.fileRepo[name] = f

	return nil
}

func (fs *FsNotify) StopFsWatcher() {
	fs.done <- true
}

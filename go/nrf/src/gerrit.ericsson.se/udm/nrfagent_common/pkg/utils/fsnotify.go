package utils

import (
	"errors"
	"sync"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/fsnotify/fsnotify"
)

var (
	fileRepo = make(map[string]FileInfo)
)

type FsNofity struct {
	watcher  *fsnotify.Watcher
	fileRepo map[string]FileInfo
	done     chan bool
	sync.Mutex
}

type FileInfo interface {
	GetFileName() string
	Handler(name, op string)
}

func InitFsWatcher() *FsNofity {
	fs := &FsNofity{
		fileRepo: make(map[string]FileInfo),
		done:     make(chan bool),
	}
	var err error
	fs.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	//defer fs.watcher.Close()
	go func() {
		for {
			select {
			case event := <-fs.watcher.Events:
				//dir, _ := path.Split(event.Name)
				//dir := filepath.Dir(event.Name)
				//fmt.Println("InitFsWatcher: fileRepo is ", fs.fileRepo)
				log.Debugf("InitFsWatcher: fileRepo is %+v", fs.fileRepo)
				if f, ok := fs.fileRepo[event.Name]; ok {
					f.Handler(event.Name, event.Op.String())
					if event.Op == fsnotify.Remove {
						_ = fs.watcher.Remove(event.Name)
						_ = fs.watcher.Add(event.Name)
					}
				} else {
					//fmt.Println("InitFsWatcher: Can not find related filerepo ",  event.Name)
					log.Debugf("InitFsWatcher: Can not find related filerepo %s", event.Name)
				}

			case err := <-fs.watcher.Errors:
				//fmt.Println("InitFsWatcher: fsnotify error is ", err)
				log.Errorf("InitFsWatcher: fsnotify error %s", err)
			case <-fs.done:
				//fmt.Println("InitFsWatcher: fsnotify is exiting")
				log.Debugf("InitFsWatcher: fsnotify is exiting")
				return
			}
		}
	}()

	return fs
}

func (fs *FsNofity) AddFileToFsWatcher(f FileInfo) error {
	name := f.GetFileName()

	//	if strings.HasSuffix(name, "/") == false {
	//		name = name + "/"
	//	}

	if _, ok := fs.fileRepo[name]; ok {
		return errors.New("Has already existed")
	}

	if err := fs.watcher.Add(name); err != nil {
		return err
	}

	fs.fileRepo[name] = f

	return nil
}

func (fs *FsNofity) StopFsWatcher() {
	fs.done <- true
}

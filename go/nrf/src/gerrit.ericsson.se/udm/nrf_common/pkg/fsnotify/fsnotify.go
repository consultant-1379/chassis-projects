package fsnotify

import (
	"sync"

	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"

	"gerrit.ericsson.se/udm/common/pkg/log"
	fsnotify3pp "github.com/fsnotify/fsnotify"
)

var (
	fsnotifyIns *FsNofity
)

type FsNofity struct {
	watcher  *fsnotify3pp.Watcher
	fileRepo map[string]fileHandler
	done     chan bool
	sync.Mutex
}

func Init() error {
	fsnotifyIns = &FsNofity{
		fileRepo: make(map[string]fileHandler),
		done:     make(chan bool),
	}
	var err error
	fsnotifyIns.watcher, err = fsnotify3pp.NewWatcher()
	if err != nil {
		return err
	}

	err = AddFileToFsWatcher()
	if err != nil {
		return err
	}

	return nil
}

func AddFileToFsWatcher() error {
	for k := range configmap.ConfigMapMap {
		if err := fsnotifyIns.watcher.Add(k); err != nil {
			log.Errorf("add file %s to fsnotify error, %s", k, err.Error())
			return err
		}

		fsnotifyIns.fileRepo[k] = &configmapHandler{fileName: k}
	}

	return nil
}

func Run() {
	go func() {
		for {
			select {
			case event := <-fsnotifyIns.watcher.Events:
				log.Infof("fsnotify comes, fileName: %s, op: %s", event.Name, event.Op.String())
				if f, ok := fsnotifyIns.fileRepo[event.Name]; ok {
					f.Handler(event.Op.String())
					if event.Op == fsnotify3pp.Remove {
						err := fsnotifyIns.watcher.Add(event.Name)
						if err != nil {
							log.Warnf("readd %s to watcher error, %v", event.Name, err)
						}
					}
				} else {
					log.Infof("%s is not in filerepo", event.Name)
				}

			case err := <-fsnotifyIns.watcher.Errors:
				log.Debugf("fsnotify error, %v", err)
			case <-fsnotifyIns.done:
				log.Debug("fsnotify is exiting")
				return
			}
		}
	}()
}

func Stop() {
	fsnotifyIns.done <- true
}

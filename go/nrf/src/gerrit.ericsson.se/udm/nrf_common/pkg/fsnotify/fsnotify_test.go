/*
  Note: If the package name changed, please change the "path" in test function TestFsnotify as well
*/

package fsnotify

import (
	"os"
	"strings"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

func init() {
	log.SetLevel(log.FatalLevel)
}

type testHandler struct {
	handled bool
}

func (t *testHandler) Handler(op string) {
	t.handled = true
}

func TestFsnotify(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	if !strings.HasSuffix(gopath, "/") {
		gopath = gopath + "/"
	}

	path := gopath + "src/gerrit.ericsson.se/udm/nrf_common/pkg/fsnotify"
	fileName := path + "/test.data"

	defer func() { _ = os.Remove(fileName) }()
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	err = Init()
	if err != nil {
		t.Fatalf("Initialize fsnotify error")
	}

	if err = fsnotifyIns.watcher.Add(path); err != nil {
		t.Fatalf("add file to fsnotify error, %s", err.Error())
	}

	fsnotifyHandler := &testHandler{handled: false}

	fsnotifyIns.fileRepo[fileName] = fsnotifyHandler
	Run()
	time.Sleep(time.Second * time.Duration(1))

	_ = os.Remove(fileName)

	time.Sleep(time.Second * time.Duration(2))

	if !fsnotifyHandler.handled {
		t.Fatalf("fs notify does not work")
	}

	Stop()
}

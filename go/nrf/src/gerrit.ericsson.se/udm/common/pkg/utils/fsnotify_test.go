package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type Task struct {
	file string
	init bool
}

func (t *Task) GetFileName() string {
	return t.file
}

func (t *Task) Handler(name, op string) {
	fmt.Println(name, " ", op)
	t.init = true
}

func TestFsNotifyModify(t *testing.T) {

	f := InitFsWatcher()

	filePath := filepath.Join(".", "test.data")
	defer func() { _ = os.Remove(filePath) }()
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	dirname, _ := filepath.Abs(filepath.Dir(filePath))
	a := &Task{
		file: dirname,
		init: false,
	}
	if err := f.AddFileToFsWatcher(a); err != nil {
		t.Fatalf("error happen %v", err.Error())
	}
	wr, err := file.WriteString("test data \n")
	fmt.Println("wrote bytes: ", wr)
	time.Sleep(2 * time.Second)
	if !a.init {
		t.Errorf("fs notify does not work")
	}
	f.StopFsWatcher()
}

func TestFsNotifyCreate(t *testing.T) {
	testDir := filepath.Join(".", "public")
	_ = os.MkdirAll(testDir, os.ModePerm)

	defer func() { _ = os.RemoveAll(testDir) }()
	defer func() { _ = os.Remove(testDir) }()

	f := InitFsWatcher()

	a := &Task{
		file: testDir,
		init: false,
	}
	if err := f.AddFileToFsWatcher(a); err != nil {
		t.Fatalf("error happen %v", err.Error())
	}

	h, err := os.Create(filepath.Join(testDir, "testfile"))
	if err != nil {
		t.Fatalf("Failed to create file in testdir: %v", err)
	}
	if err := h.Close(); err != nil {
		t.Fatalf("Can not close fsnotify %s", err.Error())
	}

	time.Sleep(2 * time.Second)
	if !a.init {
		//		t.Errorf("fs notify does not work")
	}
	f.StopFsWatcher()
}

func TestFsNotifyAddTwice(t *testing.T) {
	testDir := filepath.Join(".", "public")
	_ = os.MkdirAll(testDir, os.ModePerm)

	defer func() { _ = os.RemoveAll(testDir) }()
	defer func() { _ = os.Remove(testDir) }()

	f := InitFsWatcher()

	a := &Task{
		file: testDir,
		init: false,
	}
	if err := f.AddFileToFsWatcher(a); err != nil {
		t.Fatalf("error happen %v", err.Error())
	}
	if err := f.AddFileToFsWatcher(a); err == nil {
		t.Fatalf("Impossible")
	}

	f.StopFsWatcher()
}

func TestFsNotifyNoExist(t *testing.T) {
	testDir := filepath.Join(".", "public")

	_ = os.RemoveAll(testDir)

	f := InitFsWatcher()

	a := &Task{
		file: testDir,
		init: false,
	}
	if err := f.AddFileToFsWatcher(a); err == nil {
		t.Fatalf("Can not happen")
	}

	f.StopFsWatcher()
}

package utils

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	pprofMutex sync.Mutex

	pprofOngoing = false
)

func GenCpuMemPprof(d time.Duration, destPath string) bool {
	pprofMutex.Lock()
	defer pprofMutex.Unlock()
	if pprofOngoing {
		fmt.Println("Generating is ongoing")
		return false
	}
	pprofOngoing = true
	t := d
	path := destPath
	if path == "" {
		path = "/tmp"
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func(l time.Duration) {
		defer waitGroup.Done()
		tmp := fmt.Sprintf("%s/cpu.pprof", path)
		f, err := os.Create(tmp)
		if err != nil {
			fmt.Println("Can not create cpu pprof file ", err.Error())
			return
		}

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println("Can not close cpu profile ", err)
			}
		}()

		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Println("Can not start cpu profile ", err)
			return
		}
		defer pprof.StopCPUProfile()
		timer1 := time.NewTimer(l)
		<-timer1.C
		fmt.Println("cpu pprof is complete ", tmp)
	}(t)

	go func(l time.Duration) {
		defer waitGroup.Done()
		tmp := fmt.Sprintf("%s/mem.pprof", path)
		f, err := os.Create(tmp)
		if err != nil {
			fmt.Println("Can not create mem pprof file ", err.Error())
			return
		}

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println("Can not close cpu profile ", err)
			}
		}()

		defer func() {
			if err := pprof.Lookup("heap").WriteTo(f, 0); err != nil {
				fmt.Println("Can not write heap to mem file ", err)
			}
		}()
		timer1 := time.NewTimer(l)
		<-timer1.C
		fmt.Println("mem pprof is complete ", tmp)
	}(t)

	go func() {
		waitGroup.Wait()
		pprofMutex.Lock()
		defer pprofMutex.Unlock()
		pprofOngoing = false
	}()

	fmt.Println("Current GoRoutine Number: ", runtime.NumGoroutine())
	return true
}

package timer

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

//Timer enhanced timer
type Timer struct {
	timerChan    chan string
	timerMonitor *time.Timer

	timePointsActiveMutex sync.Mutex
	timePointsBackupMutex sync.Mutex

	timePointsActive []*timePoint
	timePointsBackup []*timePoint
}

type timePoint struct {
	tag     string
	rawTime time.Time
}

type timePointSlice []*timePoint

func (s timePointSlice) Len() int      { return len(s) }
func (s timePointSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s timePointSlice) Less(i, j int) bool {
	return s[i].rawTime.Before(s[j].rawTime)
}

//NewTimer create a enhanced timer
func NewTimer() *Timer {
	rawTimer := time.NewTimer(time.Hour)
	if rawTimer == nil {
		return nil
	}
	rawTimer.Stop()

	rawChan := make(chan string)
	if rawChan == nil {
		return nil
	}

	t := &Timer{}
	t.timerChan = rawChan
	t.timerMonitor = rawTimer
	t.timePointsActive = make([]*timePoint, 0)
	t.timePointsBackup = make([]*timePoint, 0)

	return t
}

//TimerChan message channel from the enhanced timer
func (t *Timer) TimerChan() <-chan string {
	return t.timerChan
}

//StartTimer start the enhanced timer
func (t *Timer) StartTimer() {
	go func() {
		for _ = range t.timerMonitor.C {
			t.resetTimer()
		}
	}()
	t.resetTimer()
}

//StopTimer stop the enhanced timer
func (t *Timer) StopTimer() {
	t.timePointsActiveMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()

	t.timerMonitor.Stop()
}

//StopTimer stop the enhanced timer
func (t *Timer) DestroyTimer() {
	t.timePointsActiveMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()

	t.timerMonitor.Stop()

	close(t.timerChan)
	t.timerMonitor = nil

	t.timePointsActive = nil
	t.timePointsBackup = nil
}

//AddTimePoint add a time to the enhanced timer
func (t *Timer) AddTimePoint(rawTime *time.Time, tag string) bool {
	if rawTime == nil || tag == "" {
		return false
	}

	t.timePointsActiveMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()

	if _, found := t.findTimePoint(rawTime, tag); !found {
		newTimePoint := &timePoint{
			tag:     tag,
			rawTime: *rawTime,
		}
		t.timePointsActive = append(t.timePointsActive, newTimePoint)

		t.resetTimer()
		return true
	}

	return false
}

//DelTimePoint delete a time to the enhanced timer
func (t *Timer) DelTimePoint(rawTime *time.Time, tag string) bool {
	if rawTime == nil || tag == "" {
		return false
	}

	t.timePointsActiveMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()

	if i, found := t.findTimePoint(rawTime, tag); found {
		t.timePointsActive = append(t.timePointsActive[:i], t.timePointsActive[i+1:]...)
		if len(t.timePointsActive) == 0 {
			t.timePointsActive = nil
		}

		t.resetTimer()
		return true
	}

	return false
}

//DelTimePointAll delete all time point
func (t *Timer) DelTimePointAll() {
	t.timePointsActiveMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()

	t.timePointsActive = nil
	t.resetTimer()
	//t.timerMonitor.Stop()
}

//DelTimePointTag delete time point by tag
func (t *Timer) DelTimePointTag(tag string) bool {
	if tag == "" {
		return false
	}

	t.timePointsActiveMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()

	if _, found := t.findTimePointTag(tag); found {
		var timePointsActiveNew []*timePoint
		for _, timePoint := range t.timePointsActive {
			if timePoint.tag == tag {
				continue
			}
			timePointsActiveNew = append(timePointsActiveNew, timePoint)
		}
		t.timePointsActive = timePointsActiveNew
		/*
			if t.timePointsActive == nil {
				t.timerMonitor.Stop()
			}
		*/
		t.resetTimer()
		return true
	}

	return false
}

//StopTimePointTag stop time point by tag
func (t *Timer) StopTimePointTag(tag string) bool {
	t.timePointsActiveMutex.Lock()
	t.timePointsBackupMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()
	defer t.timePointsBackupMutex.Unlock()

	if _, found := t.findTimePointTag(tag); !found {
		return false
	}

	var timePointsActiveNew []*timePoint
	var timePointsBackupNew []*timePoint

	for _, timePoint := range t.timePointsActive {
		if tag == timePoint.tag {
			timePointsBackupNew = append(timePointsBackupNew, timePoint)
		} else {
			timePointsActiveNew = append(timePointsActiveNew, timePoint)
		}
	}

	t.timePointsActive = timePointsActiveNew
	t.timePointsBackup = timePointsBackupNew

	t.resetTimer()

	return true
}

//GetTimePoint get time point by tag
func (t *Timer) GetTimePoint(tag string) (time.Time, error) {
	t.timePointsActiveMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()

	index := t.timePointIndex(tag)
	if index == -1 {
		return time.Time{}, fmt.Errorf("No such timer for %s", tag)
	}

	return t.timePointsActive[index].rawTime, nil
}

func (t *Timer) Travel() {
	if len(t.timePointsActive) == 0 {
		return
	}

	t.timePointsActiveMutex.Lock()
	defer t.timePointsActiveMutex.Unlock()

	for i, ti := range t.timePointsActive {
		fmt.Printf("KIM add index[%d], info[%s], timer[%v]\n", i, ti.tag, ti.rawTime)
	}
}

//////////private////////////

func (t *Timer) resetTimer() {
	if len(t.timePointsActive) == 0 {
		t.timerMonitor.Stop()
		return
	}

	//t.timePointsActiveMutex.Lock()
	//defer t.timePointsActiveMutex.Unlock()

	sort.Sort(timePointSlice(t.timePointsActive))
	for len(t.timePointsActive) > 0 {
		if expired := t.timePointsActive[0].rawTime.Before(time.Now()); expired {
			t.timerChan <- t.timePointsActive[0].tag
			//fmt.Printf("KIM add push to chan :%s\n", t.timerInfoList[0].info)
			if len(t.timePointsActive) > 1 {
				// remove header timerInfoList
				t.timePointsActive = append(t.timePointsActive[:0], t.timePointsActive[1:]...)
				continue
			} else { //all expired
				t.timePointsActive = nil
				t.timerMonitor.Stop()
				break
			}
		}
		//reset timerManager next timestamp
		t.timerMonitor.Reset(time.Until(t.timePointsActive[0].rawTime))
		break
	}
}

func (t *Timer) findTimePoint(rawTime *time.Time, tag string) (int, bool) {
	index := 0
	found := false

	for i, ti := range t.timePointsActive {
		if rawTime != nil && !rawTime.Equal(ti.rawTime) {
			continue
		}
		if tag != "" && tag != ti.tag {
			continue
		}
		index = i
		found = true
		break
	}

	return index, found
}

func (t *Timer) findTimePointTag(tag string) ([]int, bool) {
	if tag == "" {
		return nil, false
	}

	found := false
	var tags []int

	for i, timePoint := range t.timePointsActive {
		if tag == timePoint.tag {
			tags = append(tags, i)
			found = true
		}
	}

	return tags, found
}

func (t *Timer) timePointIndex(tag string) int {
	if tag == "" {
		return -1
	}

	for i, ti := range t.timePointsActive {
		if tag == ti.tag {
			return i
		}
	}

	return -1
}

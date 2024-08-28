package dockerstat

import (
	"io/ioutil"
	"strconv"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

// IsOverloadWithMemory is to check if memory usage of container exceeds certain percent of limit memory
func IsOverloadWithMemory(thresholdPercent float64) bool {

	var currentPercent float64

	bytes, err := ioutil.ReadFile("/sys/fs/cgroup/memory/memory.usage_in_bytes")
	if err != nil {
		log.Warnf("Failed to read config file:memory.usage_in_bytes %s", err.Error())
	}
	memoryUsage, err := strconv.ParseInt(strings.Replace(string(bytes), "\n", "", -1), 10, 64)
	if err != nil {
		log.Warnf(" string convert into int error: %s", err.Error())
	}

	bytes, err = ioutil.ReadFile("/sys/fs/cgroup/memory/memory.limit_in_bytes")
	if err != nil {
		log.Warnf("Failed to read config file:memory.limit_in_bytes %s", err.Error())
	}
	memoryLimit, err := strconv.ParseInt(strings.Replace(string(bytes), "\n", "", -1), 10, 64)
	if err != nil {
		log.Warnf(" string convert into int error: %s", err.Error())
	}

	currentPercent = float64(memoryUsage) / float64(memoryLimit)
	log.Debugf("container memory usage:%v  limit memory:%v curPercent:%v thresholdPercent:%v", memoryUsage, memoryLimit, currentPercent, thresholdPercent)
	if currentPercent > thresholdPercent {
		return true
	}
	return false
}

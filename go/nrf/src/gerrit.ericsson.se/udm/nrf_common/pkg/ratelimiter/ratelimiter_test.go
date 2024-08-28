package ratelimiter

import (
	"testing"
	"time"
)

func TestStartMonitorExpiredLimiter(t *testing.T) {
	source := "abcdefd"
	createRateLimiter(source)
	go startMonitorExpiredLimiter(1, 2)
	time.Sleep(1 * time.Second)

	_, exist := rateLimiterMap[source]

	if !exist {
		t.Fatal("Should exist, but Not!")
	}

	// Update time by getRateLimiter
	getRateLimiter(source)

	time.Sleep(time.Duration(1100) * time.Millisecond)
	_, exist = rateLimiterMap[source]

	if !exist {
		t.Fatal("Should exist, but Not!")
	}

	// Sleep
	time.Sleep(time.Duration(2100) * time.Millisecond)
	_, exist = rateLimiterMap[source]
	if exist {
		t.Fatal("Should Not exist, but Exist!")
	}

}

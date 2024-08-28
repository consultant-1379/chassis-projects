package utils

import (
	"regexp"
)

func IsDigit(s string) bool {
	matched, _ := regexp.MatchString("^\\d+$", s)
	return matched
}

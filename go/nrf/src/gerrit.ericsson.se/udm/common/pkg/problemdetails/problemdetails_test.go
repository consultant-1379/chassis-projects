package problemdetails

import (
	"fmt"
	"testing"
)

func TestToString(t *testing.T) {
	problemDetailsIns := New()
	fmt.Println(problemDetailsIns.ToString())
}

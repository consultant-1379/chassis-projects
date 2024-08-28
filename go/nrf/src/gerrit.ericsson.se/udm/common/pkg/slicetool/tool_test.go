package slicetool

import (
	"testing"
)

func TestUniqueInt(t *testing.T) {
	sliceInt := []int{1, 2, 2, 3, 3, 4, 5, 4}
	sliceInt2 := UniqueInt(sliceInt)

	if len(sliceInt2) != 5 {
		t.Fatalf("UniqueInt return wrong slice")
	}
	mapItemExist := map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
	}

	for _, item := range sliceInt2 {
		mapItemExist[item] = true
	}

	ok := true
	for _, value := range mapItemExist {
		if !value {
			ok = false
			break
		}
	}

	if !ok {
		t.Fatalf("UniqueInt return wrong slice")
	}
}

func TestUniqueString(t *testing.T) {
	sliceString := []string{"a", "b", "a", "b", "c", "c", "c", "d", "e", "d"}
	sliceString2 := UniqueString(sliceString)

	if len(sliceString2) != 5 {
		t.Fatalf("UniqueString return wrong slice")
	}
	mapItemExist := map[string]bool{
		"a": false,
		"b": false,
		"c": false,
		"d": false,
		"e": false,
	}

	for _, item := range sliceString2 {
		mapItemExist[item] = true
	}

	ok := true
	for _, value := range mapItemExist {
		if !value {
			ok = false
			break
		}
	}

	if !ok {
		t.Fatalf("UniqueString return wrong slice")
	}
}

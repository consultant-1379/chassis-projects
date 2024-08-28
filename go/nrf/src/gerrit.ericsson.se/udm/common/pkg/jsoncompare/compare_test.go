package jsoncompare

import (
	"testing"
)

func TestEqual(t *testing.T) {
	s1 := `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	s2 := `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	if !Equal([]byte(s1), []byte(s2)) {
		t.Fatalf("s1 is equal to s2, but return false!")
	}

	s1 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	s2 = `{
		"key2": "value2",
		"key1": "value1",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	if !Equal([]byte(s1), []byte(s2)) {
		t.Fatalf("s1 is equal to s2, but return false!")
	}

	s1 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	s2 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key6": {
			"key61": "value61",
			"key62": 62
		},
		"key5": [1, 2, 3]
	}`

	if !Equal([]byte(s1), []byte(s2)) {
		t.Fatalf("s1 is equal to s2, but return false!")
	}

	s1 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	s2 = `{
		"key1":  "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61":  "value61",
			"key62": 62
		}
	}`

	if !Equal([]byte(s1), []byte(s2)) {
		t.Fatalf("s1 is equal to s2, but return false!")
	}

	s1 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	s2 = `{
		"key1": "value2",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61":  "value61",
			"key62": 62
		}
	}`

	if Equal([]byte(s1), []byte(s2)) {
		t.Fatalf("s1 is  not equal to s2, but return true!")
	}

	s1 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	s2 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [2, 1, 3],
		"key6": {
			"key61":  "value61",
			"key62": 62
		}
	}`

	if Equal([]byte(s1), []byte(s2)) {
		t.Fatalf("s1 is  not equal to s2, but return true!")
	}

	s1 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61": "value61",
			"key62": 62
		}
	}`

	s2 = `{
		"key1": "value1",
		"key2": "value2",
		"key3": 3,
		"key4": 4,
		"key5": [1, 2, 3],
		"key6": {
			"key61":  "value62",
			"key62": 62
		}
	}`

	if Equal([]byte(s1), []byte(s2)) {
		t.Fatalf("s1 is  not equal to s2, but return true!")
	}

}

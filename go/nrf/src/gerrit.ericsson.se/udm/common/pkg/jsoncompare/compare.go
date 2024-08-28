package jsoncompare

import (
	"encoding/json"
	"reflect"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

// Equal compares whether s1 is equal to s2
func Equal(s1, s2 []byte) bool {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal(s1, &o1)
	if err != nil {
		log.Errorf("Error mashalling string 1 :: %s", err.Error())
		return false
	}
	err = json.Unmarshal(s2, &o2)
	if err != nil {
		log.Errorf("Error mashalling string 2 :: %s", err.Error())
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

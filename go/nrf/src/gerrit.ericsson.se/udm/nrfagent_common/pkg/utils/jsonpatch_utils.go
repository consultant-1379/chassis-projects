package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type JsonPatch []PatchOP
type PatchOP struct {
	Op    string      `json:"op,omitempty"`
	Path  string      `json:"path,omitempty"`
	Value interface{} `json:"value,omitempty"`
	From  string      `json:"from,omitempty"`
}

const (
	// op type
	PATCH_OP_AD = "add"
	PATCH_OP_RM = "remove"
	PATCH_OP_RP = "replace"
	PATCH_OP_MV = "move"
	PATCH_OP_CP = "copy"
	PATCH_OP_TS = "test"
)

func isIgnoreKey(key string, ignoreMap map[string]bool) bool {
	if ignoreMap == nil {
		return false
	}
	return ignoreMap[key]
}

// a: original map
// b: current map
func GetJSONPatch(base string, a, b map[string]interface{}, ignoreMap map[string]bool) JsonPatch {
	var patch JsonPatch
	for key, bv := range b {
		av, ok := a[key]
		if isIgnoreKey(key, ignoreMap) {
			continue
		}
		// value was added
		if !ok {
			patch = append(patch, PatchOP{Op: PATCH_OP_AD, Path: base + key, Value: bv})
			continue
		}
		// If types have changed, replace completely
		if reflect.TypeOf(av) != reflect.TypeOf(bv) {
			patch = append(patch, PatchOP{Op: PATCH_OP_RP, Path: base + key, Value: bv})
			continue
		}
		// Types are the same, compare values
		switch at := av.(type) {
		case map[string]interface{}:
			bt := bv.(map[string]interface{})
			tmpPatch := GetJSONPatch(base+key+"/", at, bt, ignoreMap)
			for index := 0; index < len(tmpPatch); index++ {
				patch = append(patch, tmpPatch[index])
			}
		case string, float64, bool:
			if !matchesValue(av, bv) {
				patch = append(patch, PatchOP{Op: PATCH_OP_RP, Path: base + key, Value: bv})
			}
		case []interface{}:
			bt := bv.([]interface{})
			tmpPatch := getPatchFromArray(base+key+"/", at, bt, ignoreMap)
			for index := 0; index < len(tmpPatch); index++ {
				patch = append(patch, tmpPatch[index])
			}
		case nil:
			switch bv.(type) {
			case nil:
				// Both nil, fine.
			default:
				patch = append(patch, PatchOP{Op: PATCH_OP_RP, Path: base + key, Value: bv})
			}
		default:
			panic(fmt.Sprintf("Unknown type:%T in key %s", av, base+key))
		}
	}
	// Now add all deleted values as nil
	for key := range a {
		_, found := b[key]
		if !found && !isIgnoreKey(key, ignoreMap) {
			patch = append(patch, PatchOP{Op: PATCH_OP_RM, Path: base + key})
		}
	}
	return patch
}

//To quickly implement, just take assumption as all addition append to the tail,
//for deletion, all items located at the right of the deletion position will be
//taken as been changed, as the result, the last one will be taken as been deleted
//a: original array
//b: current array
func getPatchFromArray(base string, a, b []interface{}, ignoreMap map[string]bool) JsonPatch {
	var patch JsonPatch
	if a == nil && b == nil {
		return patch
	} else if a == nil && b != nil || a != nil && b == nil {
		patch = append(patch, PatchOP{Op: PATCH_OP_RP, Path: strings.TrimSuffix(base, "/"), Value: b})
	}
	minLen := len(b)
	if minLen > len(a) {
		minLen = len(a)
	}
	for index := 0; index < minLen; index++ {
		//handle the changed items
		av, bv, key := a[index], b[index], strconv.Itoa(index)
		if isIgnoreKey(key, ignoreMap) {
			continue
		}
		switch at := av.(type) {
		case map[string]interface{}:
			bt := bv.(map[string]interface{})
			tmpPatch := GetJSONPatch(base+key+"/", at, bt, ignoreMap)
			for i := 0; i < len(tmpPatch); i++ {
				patch = append(patch, tmpPatch[i])
			}
		case string, float64, bool:
			if !matchesValue(av, bv) {
				patch = append(patch, PatchOP{Op: PATCH_OP_RP, Path: base + key, Value: bv})
			}
		case []interface{}:
			bt := bv.([]interface{})
			tmpPatch := getPatchFromArray(base+key+"/", at, bt, ignoreMap)
			for i := 0; i < len(tmpPatch); i++ {
				patch = append(patch, tmpPatch[i])
			}
		case nil:
			switch bv.(type) {
			case nil:
				// Both nil, fine.
			default:
				patch = append(patch, PatchOP{Op: PATCH_OP_RP, Path: base + key, Value: bv})
			}
		default:
			panic(fmt.Sprintf("Unknown type:%T in key %s", av, base+key))
		}
	}
	for index := minLen; index < len(a); index++ {
		//handle the removed items in case of a longer than b
		patch = append(patch, PatchOP{Op: PATCH_OP_RM, Path: base + strconv.Itoa(index)})
	}
	for index := minLen; index < len(b); index++ {
		//handle the added items in case of b longer than a
		patch = append(patch, PatchOP{Op: PATCH_OP_AD, Path: base + strconv.Itoa(index), Value: b[index]})
	}
	return patch
}

func matchesValue(av, bv interface{}) bool {
	if reflect.TypeOf(av) != reflect.TypeOf(bv) {
		return false
	}
	switch at := av.(type) {
	case string:
		bt := bv.(string)
		if bt == at {
			return true
		}
	case float64:
		bt := bv.(float64)
		if bt == at {
			return true
		}
	case bool:
		bt := bv.(bool)
		if bt == at {
			return true
		}
	default:
		return false
	}
	return false
}

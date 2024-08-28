package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
)

const (
	// OpAdd is 'add' operation of nfprofile patch
	OpAdd = "ADD"
	// OpReplace is 'replace' operation of nfprofile patch
	OpReplace = "REPLACE"
	// OpRemove is 'remove' operation of nfprofile patch
	OpRemove = "REMOVE"
)

var (
	// OpMap is operation list of nfprofile's patch
	OpMap = map[string]string{OpAdd: "add", OpReplace: "replace", OpRemove: "remove"}
)

func isIgnoreKey(key string, ignoreMap map[string]bool) bool {

	if ignoreMap == nil {
		return false
	}

	return ignoreMap[key]
}

// GetJSONPatch is to get nrf profile patch
func GetJSONPatch(org, cur string, ignoreMap map[string]bool) (bool, []*nrfschema.NfProfilePatchData) {
	log.Debugf("GetJSONPatch: orginal: %v", org)
	log.Debugf("GetJSONPatch: current: %v", cur)
	var obj interface{}
	//get the registered profile with map[string]interface{} type
	err := json.Unmarshal([]byte(org), &obj)
	if err != nil {
		log.Errorf("GetNrfProfilePatch: Failed to Unmarshal registered NfProfile, error: %v.", err.Error())
		return false, nil
	}
	orgRaw, ok := obj.(map[string]interface{})
	if !ok {
		log.Errorf("GetNrfProfilePatch: Failed to format type of registered NfProfile, the type is : %v", reflect.TypeOf(obj))
		return false, nil
	}
	//get the current profile with map[string]interface{} type
	err = json.Unmarshal([]byte(cur), &obj)
	if err != nil {
		log.Errorf("GetNrfProfilePatch: Failed to Unmarshal current NfProfile, error: %v.", err.Error())
		return false, nil
	}
	curRaw, ok := obj.(map[string]interface{})
	if !ok {
		log.Errorf("GetNrfProfilePatch: Failed to format type of current NfProfile, the type is : %v", reflect.TypeOf(obj))
		return false, nil
	}

	jsonPatch := getJSONPatch("/", orgRaw, curRaw, ignoreMap)
	if len(jsonPatch) == 0 {
		//no need to do update
		log.Warnf("doPatch: no change for nfProfile")
		return false, nil
	}
	log.Debugf("GetNrfProfilePatch: json patch: %+v", jsonPatch)

	return true, jsonPatch
}

// getJSONPatch is to get json path; para: a means original map,   b means current map
func getJSONPatch(base string, a, b map[string]interface{}, ignoreMap map[string]bool) []*nrfschema.NfProfilePatchData {
	var patch []*nrfschema.NfProfilePatchData
	for key, bv := range b {
		av, ok := a[key]
		if isIgnoreKey(key, ignoreMap) {
			continue
		}
		// value was added
		if !ok {
			patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpAdd, Path: base + key, NewValue: bv})
			continue
		}
		// If types have changed, replace completely
		if reflect.TypeOf(av) != reflect.TypeOf(bv) {
			patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpReplace, Path: base + key, OrigValue: av, NewValue: bv})
			continue
		}
		// Types are the same, compare values
		switch at := av.(type) {
		case map[string]interface{}:
			bt := bv.(map[string]interface{})
			tmpPatch := getJSONPatch(base+key+"/", at, bt, ignoreMap)
			for index := 0; index < len(tmpPatch); index++ {
				patch = append(patch, tmpPatch[index])
			}
		case string, float64, bool:
			if !matchesValue(av, bv) {
				patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpReplace, Path: base + key, OrigValue: av, NewValue: bv})
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
				patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpReplace, Path: base + key, OrigValue: av, NewValue: bv})
			}
		default:
			panic(fmt.Sprintf("Unknown type:%T in key %s", av, base+key))
		}
	}
	// Now add all deleted values as nil
	for key, av := range a {
		_, found := b[key]
		if !found && !isIgnoreKey(key, ignoreMap) {
			patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpRemove, Path: base + key, OrigValue: av})
		}
	}
	return patch
}

//To quickly implement, just take assumption as all addition append to the tail,
//for deletion, all items located at the right of the deletion position will be
//taken as been changed, as the result, the last one will be taken as been deleted
//a: original array
//b: current array
func getPatchFromArray(base string, a, b []interface{}, ignoreMap map[string]bool) []*nrfschema.NfProfilePatchData {

	var patch []*nrfschema.NfProfilePatchData
	if a == nil && b == nil {
		return patch
	} else if a == nil && b != nil || a != nil && b == nil {
		patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpReplace, Path: strings.TrimSuffix(base, "/"), OrigValue: a, NewValue: b})
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
			tmpPatch := getJSONPatch(base+key+"/", at, bt, ignoreMap)
			for i := 0; i < len(tmpPatch); i++ {
				patch = append(patch, tmpPatch[i])
			}
		case string, float64, bool:
			if !matchesValue(av, bv) {
				patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpReplace, Path: base + key, OrigValue: av, NewValue: bv})
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
				patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpReplace, Path: base + key, OrigValue: av, NewValue: bv})
			}
		default:
			panic(fmt.Sprintf("Unknown type:%T in key %s", av, base+key))
		}
	}
	for index := minLen; index < len(a); index++ {
		//handle the removed items in case of a longer than b
		patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpRemove, Path: base + strconv.Itoa(index), OrigValue: a[index]})
	}
	for index := minLen; index < len(b); index++ {
		//handle the added items in case of b longer than a
		patch = append(patch, &nrfschema.NfProfilePatchData{Op: OpAdd, Path: base + strconv.Itoa(index), NewValue: b[index]})
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

// ConvertPatchFormat is to convert profileChanges of NotificationData into nfprofile patch
func ConvertPatchFormat(nfProfilePatchData []*nrfschema.NfProfilePatchData) []*nrfschema.NrfInfoPatchData {
	var patch []*nrfschema.NrfInfoPatchData
	var patchItem *nrfschema.NrfInfoPatchData

	for _, nfProfile := range nfProfilePatchData {
		if nfProfile.Op != OpRemove {
			patchItem = &nrfschema.NrfInfoPatchData{Op: OpMap[nfProfile.Op], Path: nfProfile.Path, Value: nfProfile.NewValue}
		} else {
			patchItem = &nrfschema.NrfInfoPatchData{Op: OpMap[nfProfile.Op], Path: nfProfile.Path}
		}
		patch = append(patch, patchItem)
	}
	return patch
}

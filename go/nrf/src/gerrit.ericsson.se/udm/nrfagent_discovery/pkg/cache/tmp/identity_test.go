package cache

import (
	"encoding/json"
	"testing"

	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
)

var (
	validSupi          = `{"Start": "000010", "End": "000100", "Pattern": "^imsi-[0-5]{5,15}$|^nai-,+$|^suci-[2-5]{5,15}$"}`
	invalidSupi        = `{"Start": "000010", "End": "", "Pattern": "^imsi-[0-5]{5,15}$|^nai-,+$|^suci-[2-5]{5,15}$"}`
	invalidSupiPattern = `{"Start": "000010", "End": "", "Pattern": ""}`
)

func TestCheck(t *testing.T) {
	var supi identity
	err := json.Unmarshal([]byte(validSupi), &supi)
	if err != nil {
		t.Errorf("TestCheck: Unmarshal fail, err: %s", err)
	}
	ok := supi.check()
	if !ok {
		t.Errorf("TestCheck: fail to check the supi validity, supi value: %+v", supi)
	}
	err = json.Unmarshal([]byte(invalidSupi), &supi)
	if err != nil {
		t.Errorf("TestCheck: Unmarshal fail, err: %s", err)
	}
	ok = supi.check()
	if ok {
		t.Errorf("TestCheck: fail to check the supi validity, supi value: %+v", supi)
	}
}

func TestCover(t *testing.T) {
	util.PreComplieRegexp()
	validString := "imsi-000050"
	invalidString := "imsi-006000"
	var supi identity
	err := json.Unmarshal([]byte(validSupi), &supi)
	if err != nil {
		t.Errorf("TestCover: Unmarshal fail, err: %s", err)
	}
	ok := supi.cover(validString)
	if !ok {
		t.Errorf("TestCover: supi Cover validString fail")
	}
	ok = supi.cover(invalidString)
	if ok {
		t.Errorf("TestCover: supi Cover invalidString fail")
	}
}

func TestValidCheck(t *testing.T) {
	validString := "imsi-000010"
	var supi identity
	err := json.Unmarshal([]byte(validSupi), &supi)
	if err != nil {
		t.Errorf("TestvalidCheck: Unmarshal fail, err: %s", err)
	}
	ok := supi.validCheck(validString)
	if !ok {
		t.Errorf("TestvalidCheck: supi validCheck validString fail")
	}
}

func TestMeetRangeCheck(t *testing.T) {
	var supi identity
	err := json.Unmarshal([]byte(validSupi), &supi)
	if err != nil {
		t.Errorf("TestMeetRangeCheck: Unmarshal fail, err: %s", err)
	}
	ok := supi.meetRangeCheck()
	if !ok {
		t.Errorf("TestMeetRangeCheck: supi meetRangeCheck validSupi fail")
	}
	err = json.Unmarshal([]byte(invalidSupi), &supi)
	if err != nil {
		t.Errorf("TestMeetRangeCheck: Unmarshal fail, err: %s", err)
	}
	ok = supi.meetRangeCheck()
	if ok {
		t.Errorf("TestMeetRangeCheck: supi meetRangeCheck invalidSupi fail")
	}
}

func TestMeetPatternCheck(t *testing.T) {
	var supi identity
	err := json.Unmarshal([]byte(validSupi), &supi)
	if err != nil {
		t.Errorf("TesPatternCheck: Unmarshal fail, err: %s", err)
	}
	ok := supi.meetPatternCheck()
	if !ok {
		t.Errorf("TesPatternCheck: supi meetPatternCheck validSupi fail")
	}
	err = json.Unmarshal([]byte(invalidSupiPattern), &supi)
	if err != nil {
		t.Errorf("TesPatternCheck: Unmarshal fail, err: %s", err)
	}
	ok = supi.meetPatternCheck()
	if ok {
		t.Errorf("TesPatternCheck: supi meetPatternCheck invalidSupiPattern fail")
	}
}

func TestRangeCheck(t *testing.T) {
	var validNumber = "imsi-000050"
	var invalidNumber = "imsi-0020"
	var supi identity
	err := json.Unmarshal([]byte(validSupi), &supi)
	if err != nil {
		t.Errorf("TestRangeCheck: Unmarshal fail, err: %s", err)
	}
	ok := supi.rangeCheck(validNumber)
	if !ok {
		t.Errorf("TestRangeCheck: supi rangeCheck validNumber fail")
	}
	ok = supi.rangeCheck(invalidNumber)
	if ok {
		t.Errorf("TestRangeCheck: supi rangeCheck invalidNumber fail")
	}
}

func TestPatternCheck(t *testing.T) {
	var validNumber = "imsi-000050"
	var invalidNumber = "imsi-007000"
	var supi identity
	err := json.Unmarshal([]byte(validSupi), &supi)
	if err != nil {
		t.Errorf("TestPatternCheck: Unmarshal fail, err: %s", err)
	}
	ok := supi.patternCheck(validNumber)
	if !ok {
		t.Errorf("TestPatternCheck: supi patternCheck validNumber fail")
	}
	ok = supi.patternCheck(invalidNumber)
	if ok {
		t.Errorf("TestPatternCheck: supi patternCheck invalidNumber fail")
	}
}

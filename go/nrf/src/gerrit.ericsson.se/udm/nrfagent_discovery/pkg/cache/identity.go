package cache

import (
	"regexp"
	"strconv"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/utils"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/util"
)

//Identity is for both supi and gpsi
type identity struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

//SupiRange is for SupiRanges
type supiRange struct {
	SupiRanges []identity `json:"supiRanges,omitempty"`
}

//NfInfo is for NF info
type nfInfo struct {
	GroupID    string     `json:"groupId,omitempty"`
	SupiRanges []identity `json:"supiRanges,omitempty"`
	GpsiRanges []identity `json:"GpsiRanges,omitempty"`
}

//UdrInfo is for udrInfo
type udrInfo struct {
	UdrInfo *nfInfo `json:"udrInfo"`
}

//UdmInfo is for udmInfo
type udmInfo struct {
	UdmInfo *nfInfo `json:"udmInfo"`
}

//AusfInfo is for ausfInfo
type ausfInfo struct {
	AusfInfo *nfInfo `json:"ausfInfo"`
}

//PcfInfo is for pcfInfo
type pcfInfo struct {
	PcfInfo *supiRange `json:"pcfInfo"`
}

//check is for check identity format
func (identity *identity) check() bool {
	result := identity.Start != "" && identity.End != ""
	if result {
		return true
	}
	result = identity.Start == "" && identity.End == ""
	if result {
		return true
	}

	return false
}

//cover is to check if identity cover number.
func (identity *identity) cover(number string) bool {
	var result bool = false

	if identity.meetPatternCheck() {
		result = identity.patternCheck(number)
		if result {
			return true
		}
	}

	if identity.meetRangeCheck() {
		result = identity.rangeCheck(number)
	}

	return result
}

func (identity *identity) validCheck(number string) bool {
	//supi format '^(imsi-[0-9]{5,15}|nai-.+|.+)$', gpsi format '^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+|.+)$'
	if !utils.IsDigit(identity.Start) || !utils.IsDigit(identity.End) {
		log.Errorf("The start or end of Identity range is not digit")
		return false
	}

	return true
}

func (identity *identity) meetRangeCheck() bool {
	return identity.Start != "" && identity.End != ""
}

func (identity *identity) meetPatternCheck() bool {
	return identity.Pattern != ""
}

func (identity *identity) rangeCheck(number string) bool {
	if !identity.validCheck(number) {
		log.Warnf("Identity search number format is invalid")
		return false
	}

	//Both Supi and Gpsi Ranges format is "[0-9]{5,15}"
	re := util.Compile[consts.SupiRanges]
	supiInt64, _ := strconv.ParseInt(re.FindString(number), 10, 64)

	result := false
	start, _ := strconv.ParseInt(identity.Start, 10, 64)
	end, _ := strconv.ParseInt(identity.End, 10, 64)
	if supiInt64 >= start && supiInt64 <= end {
		result = true
	}
	log.Debugf("The identity is %v, digital[%d], range  is %v-%v, and the matched result is %v\n", number, supiInt64, identity.Start, identity.End, result)

	return result
}

func (identity *identity) patternCheck(number string) bool {
	reg := regexp.MustCompile(identity.Pattern)
	return reg.Match([]byte(number))
}

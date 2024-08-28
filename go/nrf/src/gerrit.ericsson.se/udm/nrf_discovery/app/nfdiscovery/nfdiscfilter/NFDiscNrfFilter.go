package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"github.com/buger/jsonparser"
	"strings"
)

//NFNRFInfoFilter to process nrfprofile filter
type NFNRFInfoFilter struct {

}

//filter is to match nrfprofile
func (a *NFNRFInfoFilter) filter(nfprofile []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	if !filterInfo.KVDBSearch && queryForm.GetExistFlag(constvalue.SearchDataSnssais) && len(queryForm.GetValue()[constvalue.SearchDataSnssais]) > 0 {
		if !a.isMatchedSnssais(queryForm.GetNRFDiscListSnssais(constvalue.SearchDataSnssais), nfprofile) {
			log.Debugf("No Matched nrfProfile with Snssais: %s", queryForm.GetNRFDiscListSnssais(constvalue.SearchDataSnssais))
			return false
		}
	}
	return true
}

//isMatchedSnssais is to match snssais in nrfprofile
func (a *NFNRFInfoFilter) isMatchedSnssais(snssais string, nfProfile []byte) bool {
	matched := false
	if len(snssais) > 0 {
		_, err := jsonparser.ArrayEach([]byte(snssais), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
			if matched {
				return
			}
			sstID, parseErr := jsonparser.GetInt(value, constvalue.SearchDataSnssaiSst)
			sdID, parseErr1 := jsonparser.GetString(value, constvalue.SearchDataSnssaiSd)
			if parseErr != nil || parseErr1 != nil {
				log.Errorf("sst or sd parse error, err=%v, err1=%v", parseErr, parseErr1)
			}
			_, err3 := jsonparser.ArrayEach(nfProfile, func(value2 []byte, dataType jsonparser.ValueType, offset int, err2 error) {
				if matched {
					return
				}
				sstID2, parseErr := jsonparser.GetInt(value2, constvalue.Sst)
				sdID2, parseErr1 := jsonparser.GetString(value2, constvalue.Sd)
				if parseErr != nil || parseErr1 != nil {
					log.Errorf("sst or sd parse error, err=%v, err1=%v", parseErr, parseErr1)
				}
				sdID = strings.ToLower(sdID)
				sdID2 = strings.ToLower(sdID2)
				if sstID == sstID2 && sdID == sdID2 {
					matched = true
					return
				}
			}, constvalue.Snssais)
			if err3 != nil {
				matched = false
				return
			}
		})

		if err != nil {
			return false
		}

	}

	return matched
}
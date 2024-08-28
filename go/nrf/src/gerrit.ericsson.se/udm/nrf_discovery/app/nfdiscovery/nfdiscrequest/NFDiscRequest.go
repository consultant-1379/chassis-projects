package nfdiscrequest

import (
	"net/url"
	"github.com/buger/jsonparser"
	"fmt"
	"net/http"
	"strings"
	"net"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"strconv"
	"sort"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"encoding/json"
	"gerrit.ericsson.se/udm/common/pkg/httpclient"
)


//NFParaMap is all parameters disc supported
var NFParaMap = map[string]bool{
	constvalue.SearchDataTargetNfType:        true,
	constvalue.SearchDataRequesterNfType:     true,
	constvalue.SearchDataServiceName:         true,
	constvalue.SearchDataRequesterNFInstFQDN: true,
	constvalue.SearchDataTargetPlmnList:      true,
	constvalue.SearchDataRequesterPlmnList:   true,
	constvalue.SearchDataTargetInstID:        true,
	constvalue.SearchDataTargetNFFQDN:        true,
	constvalue.SearchDataHnrfURI:             true,
	constvalue.SearchDataSnssais:             true,
	constvalue.SearchDataNsiList:             true,
	constvalue.SearchDataDnn:                 true,
	constvalue.SearchDataSmfServingArea:      true,
	constvalue.SearchDataTai:                 true,
	constvalue.SearchDataAmfRegionID:         true,
	constvalue.SearchDataAmfSetID:            true,
	constvalue.SearchDataGuami:               true,
	constvalue.SearchDataSupi:                true,
	constvalue.SearchDataUEIPv4Addr:          true,
	constvalue.SearchDataIPDoamin:            true,
	constvalue.SearchDataUEIPv6Prefix:        true,
	constvalue.SearchDataPGWInd:              true,
	constvalue.SearchDataPGW:                 true,
	constvalue.SearchDataGpsi:                true,
	constvalue.SearchDataExterGroupID:        true,
	constvalue.SearchDataDataSet:             true,
	constvalue.SearchDataRoutingIndic:        true,
	constvalue.SearchDataGroupIDList:         true,
	constvalue.SearchDataDnaiList:            true,
	constvalue.SearchDataUpfIwkEpsInd:        true,
	constvalue.SearchDataChfSupportedPlmn:    true,
	constvalue.SearchDataPreferredLocality:   true,
	constvalue.SearchDataAccessType:          true,
	constvalue.SearchDataSupportedFeatures:   true,
	constvalue.SearchDatacomplexQuery:        true,
}

//DiscGetPara NRF Discovery Request Parameters
type DiscGetPara struct {
	value url.Values
	flag  map[string]bool
	localCacheKey string
}

func (f *DiscGetPara) setExistFlag(SearchData string, IsExist bool) {
	f.flag[SearchData] = IsExist
}
//GetValue to get all parameters
func (f *DiscGetPara) GetValue() url.Values{
	return f.value
}
//SetFlag to set parameters flag wheter it is exist
func (f *DiscGetPara) SetFlag(key string, flag bool){
	f.flag[key] = flag
}
//SetValue to set parameter value
func (f *DiscGetPara) SetValue(key string, value []string) {
	f.value[key] = value
}
//InitMember to initial DiscGetPara
func (f *DiscGetPara) InitMember(value url.Values) {
	f.value = value
	f.flag = make(map[string]bool)
}
//GetExistFlag to get parameter flag wheter exist
func (f *DiscGetPara) GetExistFlag(SearchData string) bool {
	return f.flag[SearchData]
}
func (f *DiscGetPara) validateStringTypeForMan(SearchData string) error {
	if len(f.value[SearchData]) <= 0 {
		return fmt.Errorf(constvalue.MadatoryFieldNotExistFormat, SearchData, "SearchRequest")
	}

	if len(f.value[SearchData]) > 1 {
		return fmt.Errorf(constvalue.FieldMultipleValue, SearchData)
	}

	if f.value[SearchData][0] == "" {
		return fmt.Errorf(constvalue.FieldEmptyValue, SearchData)
	}
	f.setExistFlag(SearchData, true)
	return nil
}

func (f *DiscGetPara) validateStringTypeForOpt(SearchData string, Parttern bool) error {
	if len(f.value[SearchData]) <= 0 {
		return nil
	}

	if len(f.value[SearchData]) > 1 {
		return fmt.Errorf(constvalue.FieldMultipleValue, SearchData)
	}

	if f.value[SearchData][0] == "" {
		return fmt.Errorf(constvalue.FieldEmptyValue, SearchData)
	}

	if Parttern {
		//matched, _ := regexp.MatchString(Parttern, f.value[SearchData][0])
		matched := nfdiscutil.Compile[SearchData].MatchString(f.value[SearchData][0])
		if !matched {
			return fmt.Errorf("Parameter %s value can't match parttern", SearchData)
		}
	}

	f.setExistFlag(SearchData, true)
	return nil
}

//validateAccessType is to validate access-type
func (f *DiscGetPara) validateAccessType(SearchData string) error {
	if len(f.value[SearchData]) <= 0 {
		return nil
	}

	if len(f.value[SearchData]) > 1 {
		return fmt.Errorf(constvalue.FieldMultipleValue, SearchData)
	}

	if f.value[SearchData][0] == "" {
		return fmt.Errorf(constvalue.FieldEmptyValue, SearchData)
	}

	if f.value[SearchData][0] != constvalue.Access3GPP && f.value[SearchData][0] != constvalue.NonAccess3GPP {
		return fmt.Errorf("invalid access-type %s, should be 3GPP_ACCESS or NON_3GPP_ACCESS", f.value[SearchData][0])
	}
	f.setExistFlag(SearchData, true)
	return nil
}

func (f *DiscGetPara) validateSupportedFeature() error {
	return f.validateStringTypeForOpt(constvalue.SearchDataSupportedFeatures, true)
}

func (f *DiscGetPara) validateHnrfURI() error {
	//TODO URI regex match
	return f.validateStringTypeForOpt(constvalue.SearchDataHnrfURI, true)
}

func (f *DiscGetPara) validateGroupID() error {
	return f.validateStringTypeForOpt(constvalue.SearchDataExterGroupID, true)
}

func (f *DiscGetPara) validateRoutingIndicator() error {
	return f.validateStringTypeForOpt(constvalue.SearchDataRoutingIndic, true)
}

func (f *DiscGetPara) validateAMFRegionID() error {
	return f.validateStringTypeForOpt(constvalue.SearchDataAmfRegionID, false)
}

func (f *DiscGetPara) validateAMFSetID() error {
	return f.validateStringTypeForOpt(constvalue.SearchDataAmfSetID, false)
}

func (f *DiscGetPara) validateDataSet() error {
	return f.validateStringTypeForOpt(constvalue.SearchDataDataSet, true)
}

func (f *DiscGetPara) validatePlmnType(SearchData string) error {
	if len(f.value[SearchData]) <= 0 {
		return nil
	}

	if len(f.value[SearchData]) > 1 {
		return fmt.Errorf(constvalue.FieldMultipleValue, SearchData)
	}

	plmn := f.value[SearchData][0]

	if plmn == "" {
		return fmt.Errorf(constvalue.FieldEmptyValue, SearchData)
	}

	var mcc string
	var err error
	if SearchData == constvalue.SearchDataGuami || SearchData == constvalue.SearchDataTai {
		mcc, err = jsonparser.GetString([]byte(plmn), "plmnId", constvalue.SearchDataMcc)
	} else {
		mcc, err = jsonparser.GetString([]byte(plmn), constvalue.SearchDataMcc)
	}
	if err != nil {
		return fmt.Errorf(constvalue.MadatoryFieldNotExistFormat, "mcc", SearchData)
	}

	matched := nfdiscutil.Compile[constvalue.SearchDataMcc].MatchString(mcc)
	if !matched {
		return fmt.Errorf("invalid format for mcc in %s", SearchData)
	}
	var mnc string
	if SearchData == constvalue.SearchDataGuami || SearchData == constvalue.SearchDataTai {
		mnc, err = jsonparser.GetString([]byte(plmn), "plmnId", constvalue.SearchDataMnc)
	} else {
		mnc, err = jsonparser.GetString([]byte(plmn), constvalue.SearchDataMnc)
	}
	if err != nil {
		return fmt.Errorf(constvalue.MadatoryFieldNotExistFormat, "mnc", SearchData)
	}

	matched = nfdiscutil.Compile[constvalue.SearchDataMnc].MatchString(mnc)
	if !matched {
		return fmt.Errorf("invalid format for mnc in %s", SearchData)
	}
	f.setExistFlag(SearchData, true)
	return nil
}

//isJSON is to check if the string is a valid json
func isJSON(s string) bool {
	var js interface{}

	return json.Unmarshal([]byte(s), &js) == nil
}

func (f *DiscGetPara) validatePlmnListType(SearchData string) error {
	if _, ok := f.value[SearchData]; ok {
		if len(f.value[SearchData]) <= 0 {
			return fmt.Errorf(constvalue.FieldEmptyValue, SearchData)
		}
		for _, plmn := range (f.value[SearchData]) {
			if plmn != "" {
				ok := true
				errorInfo := ""
				if !isJSON(plmn) {
					return fmt.Errorf("parse %v json wrong", SearchData)
				}
				if !strings.Contains(plmn, "[") && !strings.Contains(plmn, "]") {
					plmn = "[" + plmn + "]"
				}
				_, err := jsonparser.ArrayEach([]byte(plmn), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
					if !ok {
						return
					}
					mcc, err := jsonparser.GetString(value, constvalue.SearchDataMcc)
					if err != nil {
						ok = false
						errorInfo = fmt.Sprintf(constvalue.MadatoryFieldNotExistFormat, "mcc", SearchData)
					}
					matched := nfdiscutil.Compile[constvalue.SearchDataMcc].MatchString(mcc)
					if !matched {
						ok = false
						errorInfo = fmt.Sprintf("invalid format for mcc in %s", SearchData)
					}
					mnc, err2 := jsonparser.GetString(value, constvalue.SearchDataMnc)
					if err2 != nil {
						ok = false
						errorInfo = fmt.Sprintf(constvalue.MadatoryFieldNotExistFormat, "mnc", SearchData)
					}
					matched = nfdiscutil.Compile[constvalue.SearchDataMnc].MatchString(mnc)
					if !matched {
						ok = false
						errorInfo = fmt.Sprintf("invalid format for mnc in %s", SearchData)
					}
				})
				if err != nil {
					return fmt.Errorf("parse %v array wrong: %v", SearchData, err)
				}

				if !ok {
					return fmt.Errorf("%s", errorInfo)
				}
			} else {
				return fmt.Errorf(constvalue.FieldEmptyValue, SearchData)
			}
		}
	} else {
		return nil
	}

	f.setExistFlag(SearchData, true)
	return nil
}

func (f *DiscGetPara) validateGuamiType(SearchData string) error {
	err := f.validatePlmnType(SearchData)
	if err == nil && f.GetExistFlag(SearchData) {
		amfID, err2 := jsonparser.GetString([]byte(f.value[SearchData][0]), constvalue.SearchDataAmfID)
		if err2 != nil {
			f.setExistFlag(SearchData, false)
			return fmt.Errorf(constvalue.MadatoryFieldNotExistFormat, "amfid", "Guami")
		}
		matched := nfdiscutil.Compile[constvalue.SearchDataAmfID].MatchString(amfID)
		if !matched {
			f.setExistFlag(SearchData, false)
			return fmt.Errorf("Invalid format for amfid in Guami")
		}

		return nil
	}

	return err
}

func (f *DiscGetPara) validateTaiType(SearchData string) error {
	err := f.validatePlmnType(SearchData)
	if err == nil && f.GetExistFlag(SearchData) {
		tac, err2 := jsonparser.GetString([]byte(f.value[SearchData][0]), constvalue.SearchDataTac)
		if err2 != nil {
			f.setExistFlag(SearchData, false)
			return fmt.Errorf(constvalue.MadatoryFieldNotExistFormat, "Tac", "Tai")
		}

		matched := nfdiscutil.Compile[constvalue.SearchDataTac].MatchString(tac)
		if !matched {
			f.setExistFlag(SearchData, false)
			return fmt.Errorf("invalid format for amfid in Gumai")
		}

		return nil
	}

	return err
}

func (f *DiscGetPara) validateIPV4Addr(SearchData string) error {
	//TODO ipv4 addr
	return f.validateStringTypeForOpt(SearchData, true)
}

func (f *DiscGetPara) validateIPV6Prefix(SearchData string) error {
	err := f.validateStringTypeForOpt(SearchData, false)
	if err != nil || f.GetExistFlag(SearchData) == false {
		return err
	}
	_, _, err = net.ParseCIDR(f.value[SearchData][0])
	if err != nil {
		f.setExistFlag(SearchData, false)
		return err
	}

	return nil
}

func (f *DiscGetPara) validateNFType(SearchData string) error {
	err := f.validateStringTypeForMan(SearchData)
	if err != nil {
		return err
	}

	/*if !checkNfTypes(f.value[SearchData][0]) {
		f.setExistFlag(SearchData, false)
		return fmt.Errorf("%s %s is invaild.", SearchData, f.value[SearchData][0])
	}*/
	f.setExistFlag(SearchData, true)
	return nil

}

func (f *DiscGetPara) validateGpsiType(SearchData string) error {
	return f.validateStringTypeForOpt(SearchData, true)
}

func (f *DiscGetPara) validateListSnssais(SearchData string) error {
	if len(f.value[SearchData]) > 0 {
		for _, item := range f.value[SearchData] {
			if item != "" {
				ok := true
				errorInfo := ""
				if !isJSON(item) {
					return fmt.Errorf("parse %v json wrong", SearchData)
				}
				if !strings.Contains(item, "[") && !strings.Contains(item, "]") {
					item = "[" + item + "]"
				}
				_, err := jsonparser.ArrayEach([]byte(item), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
					if !ok {
						return
					}
					sstID, err := jsonparser.GetInt(value, constvalue.SearchDataSnssaiSst)
					if err != nil {
						ok = false
						errorInfo = fmt.Sprintf("Get sst from snssais wrong: %v", err)
						return
					}
					if sstID < 0 || sstID > 255 {
						ok = false
						errorInfo = fmt.Sprintf("sst out of range")
						return
					}
					//sd is option, but if exist, need meet the pattern
					sd, err := jsonparser.GetString(value, constvalue.SearchDataSnssaiSd)
					if err == nil {
						matched := nfdiscutil.Compile[constvalue.SearchDataSnssaiSd].MatchString(sd)
						if !matched {
							ok = false
							errorInfo = fmt.Sprintf("sd exist, but format not right")
							return
						}
					}
				})
				if err != nil {
					return fmt.Errorf("parse snssais array wrong: %v", err)
				}

				if !ok {
					return fmt.Errorf("%s", errorInfo)
				}
			} else {
				return fmt.Errorf(constvalue.FieldEmptyValue, SearchData)
			}
		}
		f.setExistFlag(SearchData, true)
	}

	return nil
}

func (f *DiscGetPara) validateDnnType(SearchData string) error {
	return f.validateStringTypeForOpt(SearchData, false)
}

func (f *DiscGetPara) validateNFInstIDType(SearchData string) error {
	return f.validateStringTypeForOpt(SearchData, false)
}

func (f *DiscGetPara) validateSmfServingArea(SearchData string) error {
	return f.validateStringTypeForOpt(SearchData, false)
}

func (f *DiscGetPara)validateIPDoamin(SearchData string) error {
	return f.validateStringTypeForOpt(SearchData, false)
}

func (f *DiscGetPara) validatePreferredLocality(SearchData string) error {
	return f.validateStringTypeForOpt(SearchData, false)
}

func (f *DiscGetPara) validateListString(SearchData string) error {
	if _, exist := f.value[SearchData]; exist {
		if len(f.value[SearchData]) <= 0 {
			return fmt.Errorf(constvalue.ArrayFileldExistEmptyValue, SearchData)
		}

		for _, item := range f.value[SearchData] {
			if item == "" {
				return fmt.Errorf(constvalue.FieldEmptyValue, SearchData)
			}
		}

		f.setExistFlag(SearchData, true)
	}
	return nil
}

func (f *DiscGetPara) validateListStringType(SearchData string) error {
	switch SearchData {
	case constvalue.SearchDataServiceName:
		return f.validateListString(constvalue.SearchDataServiceName)
	case constvalue.SearchDataNsiList:
		return f.validateListString(constvalue.SearchDataNsiList)
	case constvalue.SearchDataGroupIDList:
		return f.validateListString(constvalue.SearchDataGroupIDList)
	case constvalue.SearchDataDnaiList:
		return f.validateListString(constvalue.SearchDataDnaiList)
	default:
		return fmt.Errorf("Invalid ListString Parameter %s", SearchData)
	}
}

func (f *DiscGetPara) validateSupiType(SearchData string) error {
	if len(f.value[SearchData]) <= 0 {
		return nil
	}

	if len(f.value[SearchData]) > 1 {
		return fmt.Errorf(constvalue.FieldMultipleValue, SearchData)
	}

	//supiRegex := "^imsi-[0-9]{5,15}$|nai-.+$"
	supi := f.value[SearchData][0]

	matched := nfdiscutil.Compile[constvalue.SearchDataSupi].MatchString(supi)
	if !matched {
		return fmt.Errorf("The %s is %s, doesn't match the regex", SearchData, supi)
	}
	f.setExistFlag(SearchData, true)
	return nil
}

func (f *DiscGetPara) validateFQDNType(SearchData string) error {
	switch SearchData {
	case constvalue.SearchDataRequesterNFInstFQDN:
		//after modify need modify FT
		return f.validateStringTypeForOpt(SearchData, false)
	case constvalue.SearchDataPGW:
		return f.validateStringTypeForOpt(SearchData, false)
	case constvalue.SearchDataTargetNFFQDN:
		return f.validateStringTypeForOpt(SearchData, false)
	default:
		return fmt.Errorf("Invalid FQDN type Parameter %s ", SearchData)
	}
}

func (f *DiscGetPara) validateIfNoneMatch() error {
	return f.validateStringTypeForOpt(constvalue.SearchDataIfNoneMatch, false)
}

func (f *DiscGetPara) validateBoolType(SearchData string) error {
	if len(f.value[SearchData]) <= 0 {
		return nil
	}

	if len(f.value[SearchData]) > 1 {
		return fmt.Errorf(constvalue.FieldMultipleValue, SearchData)
	}

	if f.value[SearchData][0] != constvalue.BoolTrueString && f.value[SearchData][0] != constvalue.BoolFalseString {
		return fmt.Errorf("%s value only true or false", SearchData)
	}

	f.setExistFlag(SearchData, true)
	return nil
}

//this func should be invoked at last on func ValidateNRFDiscovery
func (f *DiscGetPara) validateException() error {
	//if pgw-ind value is false and pgw exist, should return 400 Bad Request
	if f.GetExistFlag(constvalue.SearchDataPGW) {
		value, err := f.GetNRFDiscBoolValue(constvalue.SearchDataPGWInd)
		if err == nil && false == value {
			return fmt.Errorf("pgw-ind value is false and pgw exist, should return 400 Bad Request")
		}
	}

	return nil
}
//GetNRFDiscBoolValue to get bool type parameter value
func (f *DiscGetPara) GetNRFDiscBoolValue(SearchData string) (bool, error) {
	if !f.GetExistFlag(SearchData) {
		return false, fmt.Errorf("Field not exist")
	}

	if constvalue.BoolTrueString == f.value[SearchData][0] {
		return true, nil
	}

	return false, nil
}

func (f *DiscGetPara) validateComplexQuery() error{
	if _, exist := f.value[constvalue.SearchDatacomplexQuery]; exist{
		return  fmt.Errorf("Not Support Parameter complexQuery")
	}
	return nil
}

//validateIsSupportedParam is to refuse the abnormal params
func (f *DiscGetPara) validateIsSupportedParam() (string, error) {
	isSupported := true
	var invalidParams string
	for key := range f.value {
		if NFParaMap[key] != true {
			if invalidParams == "" {
				invalidParams = key
			} else {
				invalidParams += "," + key
			}
			isSupported = false
		}
	}
	if !isSupported {
		return invalidParams, fmt.Errorf("Not Support Parameter")
	}
	return invalidParams, nil
}

func (f *DiscGetPara) setProblemDetails(SearchData string, err error) *problemdetails.ProblemDetails {
	var problemDetail *problemdetails.ProblemDetails
	invalidpara := &problemdetails.InvalidParam{}
	problemDetail = &problemdetails.ProblemDetails{}
	problemDetail.InvalidParams = append(problemDetail.InvalidParams, invalidpara)

	problemDetail.Title = err.Error()
	problemDetail.InvalidParams[0].Param = SearchData
	problemDetail.InvalidParams[0].Reason = err.Error()

	if SearchData == constvalue.SearchDatacomplexQuery {
		problemDetail.Cause = constvalue.UnSupportedQueryParameter
	}

	return problemDetail
}

//ValidateNRFDiscovery to validate request's parameter
func (f *DiscGetPara) ValidateNRFDiscovery() (*problemdetails.ProblemDetails) {
	invalidParam, err := f.validateIsSupportedParam()
	if err != nil {
		return f.setProblemDetails(invalidParam, err)
	}

	err = f.validateListStringType(constvalue.SearchDataServiceName)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataServiceName, err)
	}

	err = f.validateBoolType(constvalue.SearchDataPGWInd)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataPGWInd, err)
	}

	err = f.validateBoolType(constvalue.SearchDataUpfIwkEpsInd)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataUpfIwkEpsInd, err)
	}

	err = f.validateListStringType(constvalue.SearchDataGroupIDList)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataGroupIDList, err)
	}

	err = f.validateNFType(constvalue.SearchDataTargetNfType)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataTargetNfType, err)
	}

	err = f.validateNFType(constvalue.SearchDataRequesterNfType)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataRequesterNfType, err)
	}

	err = f.validateFQDNType(constvalue.SearchDataRequesterNFInstFQDN)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataRequesterNFInstFQDN, err)
	}

	err = f.validateFQDNType(constvalue.SearchDataTargetNFFQDN)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataTargetNFFQDN, err)
	}

	err = f.validateNFInstIDType(constvalue.SearchDataTargetInstID)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataTargetInstID, err)
	}

	err = f.validatePlmnListType(constvalue.SearchDataRequesterPlmnList)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataRequesterPlmnList, err)
	}

	err = f.validatePlmnListType(constvalue.SearchDataTargetPlmnList)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataTargetPlmnList, err)
	}

	err = f.validatePlmnType(constvalue.SearchDataChfSupportedPlmn)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataChfSupportedPlmn, err)
	}

	err = f.validateListSnssais(constvalue.SearchDataSnssais)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataSnssais, err)
	}

	err = f.validateDnnType(constvalue.SearchDataDnn)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataDnn, err)
	}

	err = f.validateSmfServingArea(constvalue.SearchDataSmfServingArea)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataSmfServingArea, err)
	}

	err = f.validateSupiType(constvalue.SearchDataSupi)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataSupi, err)
	}

	err = f.validateAMFRegionID()
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataAmfRegionID, err)
	}

	err = f.validateAMFSetID()
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataAmfSetID, err)
	}

	err = f.validateGuamiType(constvalue.SearchDataGuami)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataGuami, err)
	}

	err = f.validateTaiType(constvalue.SearchDataTai)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataTai, err)
	}

	err = f.validateIPV4Addr(constvalue.SearchDataUEIPv4Addr)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataUEIPv4Addr, err)
	}

	err = f.validateIPDoamin(constvalue.SearchDataIPDoamin)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataIPDoamin, err)
	}

	err = f.validateIPV6Prefix(constvalue.SearchDataUEIPv6Prefix)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataUEIPv6Prefix, err)
	}

	err = f.validateFQDNType(constvalue.SearchDataPGW)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataPGW, err)
	}

	err = f.validateGpsiType(constvalue.SearchDataGpsi)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataGpsi, err)
	}

	err = f.validateGroupID()
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataExterGroupID, err)
	}

	err = f.validateDataSet()
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataDataSet, err)
	}

	err = f.validateListStringType(constvalue.SearchDataDnaiList)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataDnaiList, err)
	}

	err = f.validateRoutingIndicator()
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataRoutingIndic, err)
	}

	err = f.validateHnrfURI()
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataHnrfURI, err)
	}

	err = f.validateListStringType(constvalue.SearchDataNsiList)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataNsiList, err)
	}

	err = f.validateAccessType(constvalue.SearchDataAccessType)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataAccessType, err)
	}

	err = f.validatePreferredLocality(constvalue.SearchDataPreferredLocality)
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataPreferredLocality, err)
	}

	err = f.validateSupportedFeature()
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataSupportedFeatures, err)
	}

	err = f.validateComplexQuery()
	if err != nil {
		return  f.setProblemDetails(constvalue.SearchDatacomplexQuery, err)
	}
	//after validate pgw and pgw-ind, then invoke this func
	err = f.validateException()
	if err != nil {
		return f.setProblemDetails(constvalue.SearchDataPGWInd, err)
	}

	if cm.DiscLocalCacheEnable {
		log.Debugf("LocalCache Enable, generator Cache Key")
		f.generatorCacheKey()
	}
	return nil
}

func (f *DiscGetPara) getNRFDiscStringValue(SearchData string) string {
	if !f.GetExistFlag(SearchData) {
		return ""
	}
	return f.value[SearchData][0]
}

//GetNRFDiscHnrfURI to get parameter hnrf-uri value
func (f *DiscGetPara) GetNRFDiscHnrfURI() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataHnrfURI)
}

//GetNRFDiscSupportedFeatures to get parameter supportfeatures value
func (f *DiscGetPara) GetNRFDiscSupportedFeatures() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataSupportedFeatures)
}

//GetNRFDiscRoutingIndicator to get parameter routingindicator value
func (f *DiscGetPara) GetNRFDiscRoutingIndicator() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataRoutingIndic)
}

//GetNRFDiscExterGroupID to get externalgroupidentity value
func (f *DiscGetPara) GetNRFDiscExterGroupID() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataExterGroupID)
}

//GetNRFDiscDataSet to get parameter dataset value
func (f *DiscGetPara) GetNRFDiscDataSet() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataDataSet)
}

//GetNRFDiscGspi to get parameter gspi value
func (f *DiscGetPara) GetNRFDiscGspi() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataGpsi)
}

//GetNRFDiscIPV4Addr to get ipv4 type's parameter value
func (f *DiscGetPara) GetNRFDiscIPV4Addr(SearchData string) string {
	return f.getNRFDiscStringValue(SearchData)
}

//GetNRFDiscIPDomain to get ip-domain value
func (f *DiscGetPara) GetNRFDiscIPDomain(SearchData string) string {
	return f.getNRFDiscStringValue(SearchData)
}

//GetNRFDiscIPV6Prefix to get ipv6prefix type's parameter value
func (f *DiscGetPara) GetNRFDiscIPV6Prefix(SearchData string) string {
	return f.getNRFDiscStringValue(SearchData)
}

func (f *DiscGetPara) getNRFDiscListString(SearchData string) []string {
	if !f.GetExistFlag(SearchData) {
		return nil
	}

	var retlist []string
	for _, item := range f.value[SearchData] {
		str := strings.Replace(item, " ", "", -1)
		array := strings.Split(str, ",")
		for _, item2 := range array {
			retlist = append(retlist, item2)
		}
	}

	return retlist
}

//GetNRFDiscServiceName to get parameter servicename value
func (f *DiscGetPara) GetNRFDiscServiceName() []string {
	return f.getNRFDiscListString(constvalue.SearchDataServiceName)
}

//GetNRFDiscNsiList to get parameter nsilist value
func (f *DiscGetPara) GetNRFDiscNsiList() []string {
	return f.getNRFDiscListString(constvalue.SearchDataNsiList)
}
//GetNRFDiscGroupIDList to get parmeter grouid list
func (f *DiscGetPara) GetNRFDiscGroupIDList() []string {
	return f.getNRFDiscListString(constvalue.SearchDataGroupIDList)
}

//GetNRFDiscDnaiList to get parmeter grouid list
func (f *DiscGetPara) GetNRFDiscDnaiList() []string {
	return f.getNRFDiscListString(constvalue.SearchDataDnaiList)
}

//GetNRFDiscAccessType to get parameter access-type value
func (f *DiscGetPara) GetNRFDiscAccessType() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataAccessType)
}

//GetNRFDiscPreferredLocality is to get parameter preferred-locality value
func (f *DiscGetPara) GetNRFDiscPreferredLocality() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataPreferredLocality)
}

//GetNRFDiscAMFRegionID to get parameter amfregionid value
func (f *DiscGetPara) GetNRFDiscAMFRegionID() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataAmfRegionID)
}

//GetNRFDiscAMFSetID to get parameter amfsetid value
func (f *DiscGetPara) GetNRFDiscAMFSetID() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataAmfSetID)
}

//GetNRFDiscSupiValue to get parameter supi value
func (f *DiscGetPara) GetNRFDiscSupiValue() string {
	if !f.GetExistFlag(constvalue.SearchDataSupi) {
		return ""
	}
	return f.value[constvalue.SearchDataSupi][0]
}

//GetNRFDiscNFInstIDValue to get parameter instanceid value
func (f *DiscGetPara) GetNRFDiscNFInstIDValue() string {

	return f.getNRFDiscStringValue(constvalue.SearchDataTargetInstID)
}

//GetNRFDiscNFTypeValue to get nftype value
func (f *DiscGetPara) GetNRFDiscNFTypeValue(SearchData string) string {
	return f.getNRFDiscStringValue(SearchData)
}

//GetNRFDiscPlmnValue to get plmn value
func (f *DiscGetPara) GetNRFDiscPlmnValue(SearchData string) (string, string) {
	if !f.GetExistFlag(SearchData) {
		return "", ""
	}
	plmn := f.value[SearchData][0]
	var mcc string
	var mnc string
	var err, err1 error
	if SearchData == constvalue.SearchDataGuami || SearchData == constvalue.SearchDataTai {
		mcc, err = jsonparser.GetString([]byte(plmn), "plmnId", constvalue.SearchDataMcc)
		mnc, err1 = jsonparser.GetString([]byte(plmn), "plmnId", constvalue.SearchDataMnc)
		log.Debugf("PLMN: %s %s %s %s %s", plmn, SearchData, constvalue.SearchDataTai, mcc, mnc)
	} else {
		mcc, err = jsonparser.GetString([]byte(plmn), constvalue.SearchDataMcc)
		mnc, err1 = jsonparser.GetString([]byte(plmn), constvalue.SearchDataMnc)
	}
	if err != nil || err1 != nil {
		log.Debugf("parse plmn fail, err=%v, err1=%v", err, err1)
	}

	//if len(mnc) == 2 {
	//	mnc = "0" + mnc
	//}
	switch SearchData {
	//case constvalue.SearchDataRequesterPlmn:
	//	return mcc, mnc
	//case constvalue.SearchDataTargetPlmn:
	//case constvalue.SearchDataTai:
	case constvalue.SearchDataTai, constvalue.SearchDataGuami, constvalue.SearchDataChfSupportedPlmn:
		return mcc + mnc, ""
	default:
		return "", ""
	}

}

//GetNRFDiscPlmnListValue to get plmnList value
func (f *DiscGetPara) GetNRFDiscPlmnListValue(SearchData string) ([]string) {
	var plmnList []string
	if !f.GetExistFlag(SearchData) {
		return plmnList
	}
	for _, plmn := range (f.value[SearchData]) {
		plmnIds, valueType, _, err := jsonparser.Get([]byte(plmn))
		if err == nil {
			if valueType == jsonparser.Array {
				_, parseErr := jsonparser.ArrayEach(plmnIds, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					mcc, err1 := jsonparser.GetString(value, constvalue.SearchDataMcc)
					mnc, err2 := jsonparser.GetString(value, constvalue.SearchDataMnc)
					if err1 != nil || err2 != nil {
						log.Debugf("parse mnc, mcc error, err1=%v,err2=%v", err1, err2)
					}
					plmnList = append(plmnList, mcc + mnc)
				})
				if parseErr != nil {
					log.Debugf("arrayEach plmnIds error, err=%v", parseErr)
				}
			} else if valueType == jsonparser.Object {
				mcc, err1 := jsonparser.GetString(plmnIds, constvalue.SearchDataMcc)
				mnc, err2 := jsonparser.GetString(plmnIds, constvalue.SearchDataMnc)
				plmnList = append(plmnList, mcc + mnc)
				if err1 != nil || err2 != nil {
					log.Debugf("parse mnc, mcc error, err1=%v,err2=%v", err1, err2)
				}
			}
		}
	}
	return plmnList
}

//GetNRFDiscGuamiType to get guami value
func (f *DiscGetPara) GetNRFDiscGuamiType() (string, string) {
	if !f.GetExistFlag(constvalue.SearchDataGuami) {
		return "", ""
	}
	first, _ := f.GetNRFDiscPlmnValue(constvalue.SearchDataGuami)
	amfID, err := jsonparser.GetString([]byte(f.value[constvalue.SearchDataGuami][0]), constvalue.SearchDataAmfID)
	if err != nil {
		log.Debugf("parse guami error, err=%v", err)
	}
	return first, amfID
}

//GetNRFDiscTaiType to get tai value
func (f *DiscGetPara) GetNRFDiscTaiType() (string, string) {
	if !f.GetExistFlag(constvalue.SearchDataTai) {
		return "", ""
	}

	first, _ := f.GetNRFDiscPlmnValue(constvalue.SearchDataTai)
	log.Debugf("first: %s", first)
	tac, err := jsonparser.GetString([]byte(f.value[constvalue.SearchDataTai][0]), constvalue.SearchDataTac)
	if err != nil {
		log.Debugf("parse tai error, err=%v", err)
	}
	return first, tac
}

//GetNRFDiscDnnValue to get dnn value
func (f *DiscGetPara) GetNRFDiscDnnValue() string {

	return f.getNRFDiscStringValue(constvalue.SearchDataDnn)
}

//GetNRFDiscSmfServingArea to get smfservingarea value
func (f *DiscGetPara) GetNRFDiscSmfServingArea() string {

	return f.getNRFDiscStringValue(constvalue.SearchDataSmfServingArea)
}

//GetNRFDiscRequesterNFInstFQDN to get requesternfinstancefqdn value
func (f *DiscGetPara) GetNRFDiscRequesterNFInstFQDN() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataRequesterNFInstFQDN)
}

//GetNRFDiscPGW to get pgw value
func (f *DiscGetPara) GetNRFDiscPGW() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataPGW)
}

//GetNRFDiscTargetNFFQDN to get target nf fqdn value
func (f *DiscGetPara) GetNRFDiscTargetNFFQDN() string {
	return f.getNRFDiscStringValue(constvalue.SearchDataTargetNFFQDN)
}

//GetNRFDiscListSnssais is to get snssias from searchdata
func (f *DiscGetPara) GetNRFDiscListSnssais(SearchData string) string {
	var snssais string
	if f.GetExistFlag(SearchData) && len(f.value[SearchData]) > 0 {
		for _, item := range f.value[SearchData] {
			if !strings.Contains(item, "[") && !strings.Contains(item, "]") {
				item = "[" + item + "]"
			}
			_, err := jsonparser.ArrayEach([]byte(item), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
				sstID, parseErr1 := jsonparser.GetInt(value, constvalue.SearchDataSnssaiSst)
				sdID, parseErr2 := jsonparser.GetString(value, constvalue.SearchDataSnssaiSd)
				if parseErr1 != nil || parseErr2 != nil {
					log.Debugf("parse sst or st error, err1=%v,err2=%v", parseErr1, parseErr2)
				}
				if snssais != "" {
					snssais = fmt.Sprintf("%s,{\"sst\": %d,\"sd\": \"%s\"}", snssais, sstID, sdID)
				} else {
					snssais = fmt.Sprintf("{\"sst\": %d,\"sd\": \"%s\"}", sstID, sdID)
				}

			})

			if err != nil {
				return ""
			}
		}
		snssais = "[" + snssais + "]"
	}

	return snssais
}

//GetNRFDiscListSnssaisForCache to get snssais value for cache
func (f *DiscGetPara) GetNRFDiscListSnssaisForCache(SearchData string) []string {
        var snssais []string
	if f.GetExistFlag(SearchData) && len(f.value[SearchData]) > 0 {
		for _, item := range f.value[SearchData] {
			if !strings.Contains(item, "[") && !strings.Contains(item, "]") {
				item = "[" + item + "]"
			}
			_, err := jsonparser.ArrayEach([]byte(item), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
				sstID, parseErr1 := jsonparser.GetInt(value, constvalue.SearchDataSnssaiSst)
				sdID, parseErr2:= jsonparser.GetString(value, constvalue.SearchDataSnssaiSd)
				if parseErr1 != nil || parseErr2 != nil {
					log.Debugf("parse sst or st error, err1=%v,err2=%v", parseErr1, parseErr2)
				}
				sstIDStr := strconv.Itoa(int(sstID))
				snssais = append(snssais, (sstIDStr+strings.ToLower(sdID)))
			})

			if err != nil {
				return nil
			}
		}
	}

	return snssais
}

func (f *DiscGetPara) generatorCacheKey() {
	var parameter []string

	for k, v := range f.flag {
		log.Debugf("parameter key: %v value: %v", k, v)
		if v {
			parameter = append(parameter, k)
		}
	}

	sort.Strings(parameter)
	for _, para := range parameter{
		f.localCacheKey = f.localCacheKey + para
		switch para {
		case constvalue.SearchDataServiceName:
			name := f.GetNRFDiscServiceName()
			if name != nil {
				sort.Strings(name)
				for _, v := range name {
					f.localCacheKey = f.localCacheKey + v
				}
			}
		case constvalue.SearchDataPGWInd:
			pgwind, err := f.GetNRFDiscBoolValue(constvalue.SearchDataPGWInd)
			if err == nil {
				if pgwind {
					f.localCacheKey = f.localCacheKey + "true"
				} else {
					f.localCacheKey = f.localCacheKey + "false"
				}
			}
		case constvalue.SearchDataGroupIDList:
			grouidList := f.GetNRFDiscGroupIDList()
			if grouidList != nil {
				sort.Strings(grouidList)
				for _, groupid := range grouidList {
					f.localCacheKey = f.localCacheKey + groupid
				}
			}
		case constvalue.SearchDataTargetNfType:
			f.localCacheKey = f.localCacheKey + f.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)
		case constvalue.SearchDataRequesterNfType:
			f.localCacheKey = f.localCacheKey + f.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType)
		case constvalue.SearchDataRequesterNFInstFQDN:
			if "" != f.GetNRFDiscRequesterNFInstFQDN(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscRequesterNFInstFQDN()
			}
		case constvalue.SearchDataTargetNFFQDN:
			if "" != f.GetNRFDiscTargetNFFQDN() {
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscTargetNFFQDN()
			}
		case constvalue.SearchDataTargetInstID:
			if "" != f.GetNRFDiscNFInstIDValue() {
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscNFInstIDValue()
			}
		case constvalue.SearchDataRequesterPlmnList:
			plmnList := f.GetNRFDiscPlmnListValue(constvalue.SearchDataRequesterPlmnList)
			if len(plmnList) > 0 {
				sort.Strings(plmnList)
				for _, plmn := range (plmnList) {
					f.localCacheKey = f.localCacheKey + plmn
				}
			}
		case constvalue.SearchDataTargetPlmnList:
			plmnList := f.GetNRFDiscPlmnListValue(constvalue.SearchDataTargetPlmnList)
			if len(plmnList) > 0 {
				sort.Strings(plmnList)
				for _, plmn := range (plmnList) {
					f.localCacheKey = f.localCacheKey + plmn
				}
			}
		case constvalue.SearchDataSnssais:
			snssaisList := f.GetNRFDiscListSnssaisForCache(constvalue.SearchDataSnssais)
			if snssaisList != nil {
				sort.Strings(snssaisList)
				for _, snssais := range snssaisList {
					f.localCacheKey = f.localCacheKey + snssais
				}
			}
		case constvalue.SearchDataDnn:
			if "" != f.GetNRFDiscDnnValue(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscDnnValue()
			}
		case constvalue.SearchDataSmfServingArea:
			if "" != f.GetNRFDiscSmfServingArea() {
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscSmfServingArea()
			}
		case constvalue.SearchDataSupi:
			if "" != f.GetNRFDiscSupiValue(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscSupiValue()
			}
		case constvalue.SearchDataAmfRegionID:
			if "" != f.GetNRFDiscAMFRegionID(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscAMFRegionID()
			}
		case constvalue.SearchDataAmfSetID:
			if "" != f.GetNRFDiscAMFSetID() {
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscAMFSetID()
			}
		case constvalue.SearchDataGuami:
			plmn, amfid := f.GetNRFDiscGuamiType()
			if plmn != "" && amfid != ""{
				f.localCacheKey = f.localCacheKey + plmn
				f.localCacheKey = f.localCacheKey + strings.ToLower(amfid)
			}
		case constvalue.SearchDataTai:
			plmn, tac := f.GetNRFDiscTaiType()
			if plmn != "" && tac != ""{
				f.localCacheKey = f.localCacheKey + plmn
				f.localCacheKey = f.localCacheKey + strings.ToLower(tac)
			}
		case constvalue.SearchDataUEIPv4Addr:
			if "" != f.GetNRFDiscIPV4Addr(constvalue.SearchDataUEIPv4Addr){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscIPV4Addr(constvalue.SearchDataUEIPv4Addr)
			}
		case constvalue.SearchDataIPDoamin:
			if "" != f.GetNRFDiscIPDomain(constvalue.SearchDataIPDoamin){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscIPDomain(constvalue.SearchDataIPDoamin)
			}
		case constvalue.SearchDataUEIPv6Prefix:
			if "" != f.GetNRFDiscIPV6Prefix(constvalue.SearchDataUEIPv6Prefix){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscIPV6Prefix(constvalue.SearchDataUEIPv6Prefix)
			}
		case constvalue.SearchDataPGW:
			if "" != f.GetNRFDiscPGW(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscPGW()
			}
		case constvalue.SearchDataGpsi:
			if "" != f.GetNRFDiscGspi(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscGspi()
			}
		case constvalue.SearchDataExterGroupID:
			if "" != f.GetNRFDiscExterGroupID(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscExterGroupID()
			}
		case constvalue.SearchDataDataSet:
			if "" != f.GetNRFDiscDataSet(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscDataSet()
			}
		case constvalue.SearchDataRoutingIndic:
			if "" != f.GetNRFDiscRoutingIndicator(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscRoutingIndicator()
			}
		case constvalue.SearchDataDnaiList:
			dnaiList := f.GetNRFDiscDnaiList()
			if dnaiList != nil {
				sort.Strings(dnaiList)
				for _, dnai := range(dnaiList) {
					f.localCacheKey = f.localCacheKey + dnai
				}
			}
		case constvalue.SearchDataUpfIwkEpsInd:
			upfIwkEpsInd, err := f.GetNRFDiscBoolValue(constvalue.SearchDataUpfIwkEpsInd)
			if err == nil {
				if upfIwkEpsInd {
					f.localCacheKey = f.localCacheKey + "true"
				} else {
					f.localCacheKey = f.localCacheKey + "false"
				}
			}
		case constvalue.SearchDataChfSupportedPlmn:
			plmn, _ := f.GetNRFDiscPlmnValue(constvalue.SearchDataChfSupportedPlmn)
			if plmn != "" {
				f.localCacheKey = f.localCacheKey + plmn
			}
		case constvalue.SearchDataPreferredLocality:
			if "" != f.GetNRFDiscPreferredLocality(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscPreferredLocality()
			}
		case constvalue.SearchDataAccessType:
			if "" != f.GetNRFDiscAccessType(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscAccessType()
			}
		case constvalue.SearchDataHnrfURI:
			if "" != f.GetNRFDiscHnrfURI(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscHnrfURI()
			}
		case constvalue.SearchDataNsiList:
			nsiList := f.GetNRFDiscNsiList()
			if nsiList != nil {
				sort.Strings(nsiList)
				for _, nsi := range nsiList{
					f.localCacheKey = f.localCacheKey + nsi
				}
			}
		case constvalue.SearchDataSupportedFeatures:
			if "" != f.GetNRFDiscSupportedFeatures(){
				f.localCacheKey = f.localCacheKey + f.GetNRFDiscSupportedFeatures()
			}
		}

	}
}

//GetLocalCacheKey to get cache key
func (f *DiscGetPara) GetLocalCacheKey() string {
	return f.localCacheKey
}

//GetNRFDiscIfNoneMatch to get header ifnonematch value
func GetNRFDiscIfNoneMatch(req *http.Request) []string {
	value := req.Header.Get(constvalue.SearchDataIfNoneMatch)
	if value == "" {
		return nil
	}

	value = strings.Replace(value, "W/", "", -1)
	valueList := strings.Split(value, ",")
	return valueList

}

//GetNRFDiscForward to get header forward value
func GetNRFDiscForward(req *http.Request) string {
	forwarded := req.Header.Get(constvalue.SearchDataForward)
	return strings.Replace(forwarded, " ", "", -1)
}

//GetNRFDiscReqCacheControl to get  request header cache-control value
func  GetNRFDiscReqCacheControl(req *http.Request) []string {

	value := req.Header.Get(constvalue.SearchDataCacheControl)
	log.Debugf("Request Cache-Control Header Field-value: %s", value)
	valList := strings.Split(value, ",")
	var retList []string
	for _, v := range valList {
		vv := strings.Replace(v, " ", "", -1)
		retList = append(retList, vv)
	}
	return retList
}

//GetNRFDiscRespCacheControl to get response header cache-control value
func GetNRFDiscRespCacheControl(resp *httpclient.HttpRespData) []string {
	value := resp.Header.Get(constvalue.SearchDataCacheControl)
	log.Debugf("Response Cache-Control Header Field-value: %s", value)
	valList := strings.Split(value, ",")
	var retList []string
	for _, v := range valList {
		vv := strings.Replace(v, " ", "", -1)
		retList = append(retList, vv)
	}
	return retList
}
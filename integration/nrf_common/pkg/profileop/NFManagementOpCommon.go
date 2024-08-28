package profileop

import (
	"com/dbproxy"
	"com/dbproxy/nfmessage/common"
	"com/dbproxy/nfmessage/nfprofile"
	"com/dbproxy/nfmessage/nrfaddress"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/jsoncompare"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/common/pkg/slicetool"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"

	"github.com/buger/jsonparser"
)

const (
	// ProfileNf indicates geting nf profile successfully
	ProfileNf = 1
	// ProfileNrf indicates geting nrf profile successfully
	ProfileNrf = 2
	// ProfileNotFound indicates not geting nrf profile
	ProfileNotFound = 3
	// ProfileError indicates geting profile failure
	ProfileError = 4
)

const (
	//SupiRangeDelete supiRange deleted
	SupiRangeDelete = 0
	//SupiRangeAdd supiRange add
	SupiRangeAdd = 1
	//SupiRangeChange supiRange changed
	SupiRangeChange = 2
	//SupiRangeSame supiRange same
	SupiRangeSame = 3
)

const (
	//GpsiRangeDelete gpsiRange deleted
	GpsiRangeDelete = 0
	//GpsiRangeAdd gpsiRange add
	GpsiRangeAdd = 1
	//GpsiRangeChange gpsiRange changed
	GpsiRangeChange = 2
	//GpsiRangeSame gpsiRange same
	GpsiRangeSame = 3
)

// ValidateOtherRules is to validate NF profile
func ValidateOtherRules(nfProfile []byte) *problemdetails.ProblemDetails {
	profile := &nrfschema.TNFProfile{}
	err := json.Unmarshal(nfProfile, profile)
	if err != nil {
		log.Warnf("Unmarshal nf profile error, %v", err)
		return &problemdetails.ProblemDetails{
			Title: "not a valid nf profile",
		}
	}

	problemDetails := profile.Validate()
	if problemDetails != nil {
		return problemDetails
	}

	return nil
}

// GetNFInstanceID is to get instance ID from body
func GetNFInstanceID(body []byte) (string, *problemdetails.ProblemDetails) {
	nfInstanceId, err := jsonparser.GetString(body, constvalue.NfInstanceId)
	if err != nil {
		errorInfo := fmt.Sprintf(constvalue.MadatoryFieldNotExistFormat, constvalue.NfInstanceId, constvalue.NfProfile)
		return "", &problemdetails.ProblemDetails{
			Title: "not a valid nf profile",
			InvalidParams: []*problemdetails.InvalidParam{
				&problemdetails.InvalidParam{
					Param:  constvalue.NfInstanceId,
					Reason: errorInfo,
				},
			},
		}
	}
	return nfInstanceId, nil
}

// GetNFType is to get NF type from body
func GetNFType(body []byte) (string, *problemdetails.ProblemDetails) {
	nfType, err := jsonparser.GetString(body, constvalue.NfType)
	if err != nil {
		errorInfo := fmt.Sprintf(constvalue.MadatoryFieldNotExistFormat, constvalue.NfType, constvalue.NfProfile)
		return "", &problemdetails.ProblemDetails{
			Title: "not a valid nf profile",
			InvalidParams: []*problemdetails.InvalidParam{
				&problemdetails.InvalidParam{
					Param:  constvalue.NfType,
					Reason: errorInfo,
				},
			}}
	}
	return nfType, nil
}

// GetNFStatus is to get NF status from body
func GetNFStatus(body []byte) (string, *problemdetails.ProblemDetails) {
	nfStatus, err := jsonparser.GetString(body, constvalue.NfStatus)
	if err != nil {
		errorInfo := "Can not find nfStatus in nfProfile"
		return "", &problemdetails.ProblemDetails{
			Title: "not a valid nf profile",
			InvalidParams: []*problemdetails.InvalidParam{
				&problemdetails.InvalidParam{
					Param:  constvalue.NfStatus,
					Reason: errorInfo,
				},
			},
		}
	}
	return nfStatus, nil
}

//GetServiceName is for get serviceName from nfprofile
func GetServiceName(nfProfile []byte) ([]string, *problemdetails.ProblemDetails) {

	nfServices, dataType, _, err := jsonparser.Get(nfProfile, constvalue.NfServices)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		errorInfo := fmt.Sprintf("parsing failed for %s", constvalue.NfServices)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	var serviceNames []string
	errorInfo := ""

	_, err = jsonparser.ArrayEach(nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if errorInfo != "" {
			return
		}
		serviceName, err := jsonparser.GetString(value, constvalue.NFServiceName)
		if err != nil {
			errorInfo = fmt.Sprintf(constvalue.MadatoryFieldNotExistFormat, constvalue.NFServiceName, constvalue.NfServices)
			return
		}
		serviceNames = append(serviceNames, serviceName)
	})

	if err != nil {
		errorInfo = fmt.Sprintf("parsing array fail for %s", constvalue.NfServices)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	if errorInfo != "" {
		return nil, &problemdetails.ProblemDetails{
			Title: "not a valid nf profile",
			InvalidParams: []*problemdetails.InvalidParam{
				&problemdetails.InvalidParam{
					Param:  constvalue.NFServiceName,
					Reason: errorInfo,
				},
			},
		}
	}

	return serviceNames, nil
}

// GenerateExpiredTime is to generate expired time
func GenerateExpiredTime(validityPeriodInSecond int) (expiredTimeInMilsecond uint64) {
	//	The value of uint64 expired_time is described as this : the difference, measured in milliseconds, between the current time and midnight, January 1, 1970 UTC.
	return (uint64(time.Now().Unix()) + uint64(validityPeriodInSecond)) * 1000
}

// GenerateLastUpdateTime is to generate last time
func GenerateLastUpdateTime() uint64 {
	return uint64(time.Now().Unix()) * 1000
}

// GenerateProfileUpdateTime is to generate nfprofile update time
func GenerateProfileUpdateTime() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

// GetNFInstanceProfileInfo is to get NF profile by instance ID
func GetNFInstanceProfileInfo(instanceID string) (string, error) {
	nfProfileKey := &nfprofile.NFProfileGetRequest_TargetNfInstanceId{
		TargetNfInstanceId: instanceID,
	}
	getNFProfileReq := &nfprofile.NFProfileGetRequest{
		Data: nfProfileKey,
	}

	getNFProfileRsp, err := dbmgmt.GetNFProfile(getNFProfileReq)
	if err != nil {
		return "", err
	}
	if getNFProfileRsp.Code != dbmgmt.DbGetSuccess {
		return "", err
	}
	return getNFProfileRsp.GetNfProfile()[0], nil
}

// ContructNfInstanceURI is to construct nf instance uri
func ContructNfInstanceURI(instanceId string) string {
	return fmt.Sprintf("%s%s/%s", cm.GetMgmtIngressAddress(), constvalue.NfInstancesResouceURL, instanceId)
}

// ContructNfInstancesURI is to construct NF uri
func ContructNfInstancesURI() string {
	return fmt.Sprintf("%s%s", cm.GetMgmtIngressAddress(), constvalue.NfInstancesResouceURL)
}

// CurrentTimeInMilsecond is get current time
func CurrentTimeInMilsecond() uint64 {
	return uint64(time.Now().Unix()) * 1000
}

//not used function
func checkNfStatus(nfStatus string) bool {
	_, ok := constvalue.NFStatusMap[nfStatus]
	return ok
}

func validateNfService(nfServices []byte, dataType jsonparser.ValueType) (string, bool) {
	ok := true
	errorInfo := ""
	switch dataType {
	case jsonparser.Array:
		_, err := jsonparser.ArrayEach(nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			fields := [4]string{constvalue.NFServiceInstanceId, constvalue.NFServiceName, constvalue.NFServiceVersions, constvalue.NFServiceScheme}
			for _, field := range fields {
				_, err := jsonparser.GetString(value, field)
				if err != nil {
					errorInfo = fmt.Sprintf("Lack of Mandatory field %s for %s", field, constvalue.NfServices)
					ok = false
					return
				}
			}
		})

		if err != nil {
			ok = false
			errorInfo = fmt.Sprintf("parsing arry fail for %s", constvalue.NfServices)
		}
	case jsonparser.NotExist:
		errorInfo = fmt.Sprintf("%s is not existed", constvalue.NfServices)
		ok = false

	default:
		errorInfo = fmt.Sprintf("The type of %s is not Array", constvalue.NfServices)
		ok = false
	}
	return errorInfo, ok
}

//GetAddrFromURL is to get address from url, such as "http://192.168.1.1:80/test" return "192.168.1.1:80"
func GetAddrFromURL(urlAddr string) string {
	urlStruct, err := url.Parse(urlAddr)
	if err != nil {
		return ""
	}
	return urlStruct.Host
}

// LastSubString is to get last sub string
func LastSubString(s string, sep string) string {
	index := strings.LastIndex(s, sep)
	if index == -1 {
		return ""
	}

	lengthOfSubstr := len([]rune(sep))
	start := index + lengthOfSubstr

	rs := []rune(s)
	end := len(rs)
	substr := string(rs[start:end])

	if substr == "" {
		return ""
	}
	return substr
}

//GetChangedServiceName is for get changed serviceName for nfprofile
func GetChangedServiceName(oldNfProfile []byte, newNfProfile []byte) ([]string, error) {
	oldServiceNames, err := GetServiceName(oldNfProfile)
	if err != nil {
		return nil, fmt.Errorf(err.Title)
	}

	newServiceNames, err := GetServiceName(newNfProfile)
	if err != nil {
		return nil, fmt.Errorf(err.Title)
	}

	if oldServiceNames == nil && newServiceNames == nil {
		return nil, nil
	}

	if oldServiceNames == nil && newServiceNames != nil {
		return slicetool.UniqueString(newServiceNames), nil
	}

	if oldServiceNames != nil && newServiceNames == nil {
		return slicetool.UniqueString(oldServiceNames), nil
	}

	var intersectionServiceNames, unionServiceNames []string

	for _, oldServiceName := range oldServiceNames {
		found := false
		for _, newServiceName := range newServiceNames {
			if newServiceName == oldServiceName {
				found = true
				break
			}
		}
		if found {
			intersectionServiceNames = append(intersectionServiceNames, oldServiceName)
		} else {
			unionServiceNames = append(unionServiceNames, oldServiceName)
		}
	}

	for _, newServiceName := range newServiceNames {
		found := false
		for _, oldServiceName := range oldServiceNames {
			if oldServiceName == newServiceName {
				found = true
				break
			}
		}
		if !found {
			unionServiceNames = append(unionServiceNames, newServiceName)
		}
	}

	for _, item := range slicetool.UniqueString(intersectionServiceNames) {
		if isServiceChanged(item, oldNfProfile, newNfProfile) {
			unionServiceNames = append(unionServiceNames, item)
		}
	}

	return slicetool.UniqueString(unionServiceNames), nil

}

//GetUnionServiceName is for get union serviceName
func GetUnionServiceName(oldNfProfile []byte, newNfProfile []byte) ([]string, error) {
	oldServiceNames, err := GetServiceName(oldNfProfile)
	if err != nil {
		return nil, fmt.Errorf(err.Title)
	}

	newServiceNames, err := GetServiceName(newNfProfile)
	if err != nil {
		return nil, fmt.Errorf(err.Title)
	}

	if oldServiceNames == nil && newServiceNames == nil {
		return nil, nil
	}

	if oldServiceNames == nil && newServiceNames != nil {
		return slicetool.UniqueString(newServiceNames), nil
	}

	if oldServiceNames != nil && newServiceNames == nil {
		return slicetool.UniqueString(oldServiceNames), nil
	}

	for _, oldServiceName := range oldServiceNames {
		found := false
		for _, newServiceName := range newServiceNames {
			if newServiceName == oldServiceName {
				found = true
				break
			}
		}
		if !found {
			newServiceNames = append(newServiceNames, oldServiceName)
		}
	}

	return slicetool.UniqueString(newServiceNames), nil

}

func isServiceChanged(serviceName string, oldNfProfile []byte, newNfProfile []byte) bool {
	var oldServices, newServices [][]byte

	_, err := jsonparser.ArrayEach(oldNfProfile, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceNameInProfile, err := jsonparser.GetString(value, constvalue.NFServiceName)
		if err != nil {
			return
		}
		if serviceNameInProfile == serviceName {
			oldServices = append(oldServices, value)
		}
	}, constvalue.NfServices)

	if err != nil {
		oldServices = nil
	}

	_, err = jsonparser.ArrayEach(newNfProfile, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceNameInProfile, err := jsonparser.GetString(value, constvalue.NFServiceName)
		if err != nil {
			return
		}
		if serviceNameInProfile == serviceName {
			newServices = append(newServices, value)
		}
	}, constvalue.NfServices)

	if err != nil {
		newServices = nil
	}

	if oldServices == nil && newServices == nil {
		return false
	}

	if len(oldServices) != len(newServices) {
		return true
	}

	for _, oldItem := range oldServices {
		found := false
		for _, newItem := range newServices {
			if jsoncompare.Equal(oldItem, newItem) {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}

	for _, newItem := range newServices {
		found := false
		for _, oldItem := range oldServices {
			if jsoncompare.Equal(oldItem, newItem) {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}

	return false
}

//IsProfileChanged is for check if nfprofile changed
func IsProfileChanged(oldNfProfile []byte, newNfProfile []byte) bool {
	return !jsoncompare.Equal(oldNfProfile, newNfProfile)
}

// IsProfileCommonPartChanged is for check whether value except for nfServices is changed
func IsProfileCommonPartChanged(oldNfProfile []byte, newNfProfile []byte) bool {
	// jsonparser.Delete will damage source structure, so use a copy as input
	oldNfProfileStr, newNfProfileStr := string(oldNfProfile), string(newNfProfile)
	oldCommonPart := jsonparser.Delete([]byte(oldNfProfileStr), constvalue.NfServices)
	newCommonPart := jsonparser.Delete([]byte(newNfProfileStr), constvalue.NfServices)
	return !jsoncompare.Equal(oldCommonPart, newCommonPart)
}

// CheckSupportedFields is to check supported filed
func CheckSupportedFields(queryForm url.Values, fields ...string) []string {
	var invalidParameters []string
	for key := range queryForm {
		isFound := false
		for _, field := range fields {
			if key == field {
				isFound = true
				break
			}
		}

		if !isFound {
			invalidParameters = append(invalidParameters, key)
		}
	}
	return invalidParameters
}

// ConstructNRFAddressWithPlmnID is to construnct NRF address
func ConstructNRFAddressWithPlmnID(plmnId string) []dbmgmt.NRFAddress {

	mcc := plmnId[:3]
	mnc := plmnId[3:]
	if len(mnc) == 2 {
		mnc = "0" + mnc
	}

	fqdn := fmt.Sprintf(constvalue.NRFFqdnFormat, mnc, mcc)
	nrfAddrs := []dbmgmt.NRFAddress{}
	nrfAddr := dbmgmt.NRFAddress{
		Scheme: constvalue.RemoteDefaultScheme,
		Fqdn:   fqdn,
		Port:   constvalue.RemoteDefaultPort,
	}
	nrfAddrs = append(nrfAddrs, nrfAddr)
	return nrfAddrs
}

// GetNRFAddressFromDB is to get NRF address from db
func GetNRFAddressFromDB(plmnID string) []dbmgmt.NRFAddress {
	plmnList := []string{}
	plmnList = append(plmnList, plmnID)
	index := &nrfaddress.NRFAddressGetIndex{
		NrfAddressKey1: plmnList,
	}

	filter := &nrfaddress.NRFAddressFilter{
		Index: index,
	}

	nrfAddressFilter := &nrfaddress.NRFAddressGetRequest_Filter{
		Filter: filter,
	}

	nrfAddressGetRequest := &nrfaddress.NRFAddressGetRequest{
		Data: nrfAddressFilter,
	}

	nrfAddressResponse, err := dbmgmt.GetNRFAddress(nrfAddressGetRequest)
	if err != nil {
		log.Warningf("Failed to get NRF address with plmnId %s, error is %v", plmnID, err)
		return []dbmgmt.NRFAddress{}
	}

	if nrfAddressResponse.Code != dbmgmt.DbGetSuccess {
		log.Warningf("Failed to get NRF address with plmnId %s, the db ret code is %d", plmnID, nrfAddressResponse.Code)
		return []dbmgmt.NRFAddress{}

	}

	plmnNrfAddressProv := &dbmgmt.PlmnNRFAddressProv{}
	nrfAddrsProv := []dbmgmt.NRFAddressProv{}

	for _, nrfAddressData := range nrfAddressResponse.NrfAddressData {
		if err := json.Unmarshal(nrfAddressData, plmnNrfAddressProv); err != nil {
			log.Warnf("Unmarshal NrfAddressData error. %v", err)
			continue
		}

		nrfAddrsProv = append(nrfAddrsProv, plmnNrfAddressProv.NrfAddresses...)
	}

	//mapping provision NRFAddressProv to NRFAddress
	var nrfAddr dbmgmt.NRFAddress
	nrfAddrs := []dbmgmt.NRFAddress{}
	for _, nrfAddressProv := range nrfAddrsProv {
		nrfAddr.Scheme = nrfAddressProv.Scheme
		nrfAddr.Port = nrfAddressProv.Port
		if nrfAddressProv.Address.Fqdn != "" {
			nrfAddr.Fqdn = nrfAddressProv.Address.Fqdn
		} else if nrfAddressProv.Address.Ipv4Address != "" {
			nrfAddr.Fqdn = nrfAddressProv.Address.Ipv4Address
		} else if nrfAddressProv.Address.Ipv6Address != "" {
			nrfAddr.Fqdn = nrfAddressProv.Address.Ipv6Address
		} else {
			continue
		}
		nrfAddrs = append(nrfAddrs, nrfAddr)
	}
	return nrfAddrs
}

// GetNfInfoChangedCode is to get NFInfo changed code
func GetNfInfoChangedCode(oldNfProfile []byte, newNfProfile []byte) int {
	// Return code:
	// 0 : not change,
	// 1 : add,
	// 2 : change
	// 3 : delete
	oldNfInfo := GetNfInfo(oldNfProfile)
	newNfInfo := GetNfInfo(newNfProfile)

	if oldNfInfo == nil && newNfInfo == nil {
		return 0
	}

	if oldNfInfo == nil && newNfInfo != nil {
		return 1
	}

	if oldNfInfo != nil && newNfInfo == nil {
		return 3
	}

	if oldNfInfo != nil && newNfInfo != nil {
		if !jsoncompare.Equal(oldNfInfo, newNfInfo) {
			return 2
		}
		return 0
	}

	return 0
}

// GetNfInfo Only return nfInfo which should be summarize in nrfInfo
func GetNfInfo(nfProfile []byte) []byte {
	nfType, _ := GetNFType(nfProfile)
	if _, ok := constvalue.NFInfoMap[nfType]; !ok {
		return nil
	}
	nfInfo, dataType, _, err := jsonparser.Get(nfProfile, constvalue.NFInfoMap[nfType])
	if dataType == jsonparser.NotExist || err != nil {
		return nil
	}

	return nfInfo
}

//GetNfProfilesFromDbResponse is to extract nfprofiles from the db response
func GetNfProfilesFromDbResponse(nfProfileResponse *nfprofile.NFProfileGetResponse) ([]string, error) {

	if nfProfileResponse == nil {
		return nil, fmt.Errorf("nil Pointer")
	}

	if nfProfileResponse.GetCode() == dbmgmt.DbDataNotExist {
		return nil, nil
	}

	if nfProfileResponse.GetCode() != dbmgmt.DbGetSuccess {
		return nil, fmt.Errorf("DB error")
	}
	nfProfilesInfo := nfProfileResponse.GetNfProfile()
	fragmentInfo := nfProfileResponse.GetFragmentInfo()
	if fragmentInfo != nil {
		for fragmentInfo.TransmittedNumber < fragmentInfo.TotalNumber {
			nfProfileFragmentInfo := &nfprofile.NFProfileGetRequest_FragmentSessionId{
				FragmentSessionId: fragmentInfo.FragmentSessionId,
			}

			nfProfileGetRequest := &nfprofile.NFProfileGetRequest{
				Data: nfProfileFragmentInfo,
			}

			getResp, err := dbmgmt.GetNFProfile(nfProfileGetRequest)

			if err != nil {
				return nil, err
			}

			if getResp.GetCode() == dbmgmt.DbDataNotExist {
				break
			}

			if getResp.GetCode() != dbmgmt.DbGetSuccess {
				return nil, fmt.Errorf("DB error")
			}

			fragmentInfo = getResp.GetFragmentInfo()
			nfProfilesInfo = append(nfProfilesInfo, getResp.GetNfProfile()...)

			if fragmentInfo == nil {
				break
			}
		}
	}
	return nfProfilesInfo, nil
}

// NFResponseHander is a common response handler for NF operations, such as NFUpdate, NFDeregister, NFUnsubscribe
func NFResponseHander(rw http.ResponseWriter, req *http.Request, statuscode int, contentType string, body string) {
	if contentType != "" {
		rw.Header().Set("Content-Type", contentType)
	}

	if internalconf.HTTPWithXVersion {
		rw.Header().Set("X-Version", cm.ServiceVersion)
	}

	if statuscode == http.StatusTooManyRequests {
		retryAfter := fmt.Sprintf("%d", utils.RandomInt(internalconf.RetryAfterRangeStart, internalconf.RetryAfterRangeEnd))
		rw.Header().Set("Retry-After", retryAfter)
	}

	rw.WriteHeader(statuscode)
	if body != "" {
		_, err := rw.Write([]byte(body))
		if err != nil {
			log.Warnf("%v", err)
		}
	}
}

//ConstructPlmnID returns the mcc+mnc
func ConstructPlmnID(mcc, mnc string) string {
	return mcc + mnc
}

// RequestFromNRFProv is to get sub string
func RequestFromNRFProv(host string, innerPort int) bool {
	_, port, err := net.SplitHostPort(host)
	if err != nil {
		log.Warnf("RequestFromNRFProv split HostPort:%s failure", host)
		return false
	}
	if port != strconv.Itoa(innerPort) {
		return false
	}
	return true
}

// RequestFromNRFMgmt is to get sub string
func RequestFromNRFMgmt(host string, innerPort int) bool {
	_, port, err := net.SplitHostPort(host)

	if err != nil {
		log.Warnf("RequestFromNRFMgmt split HostPort:%s failure", host)
		return false
	}
	if port != strconv.Itoa(innerPort) {
		return false
	}
	return true
}

//ConstructNFprofileForGRPC construct a json including NF profile which is to be put into DB
func ConstructNFprofileForGRPC(expiredTime uint64, lastUpdateTime uint64, profileUpdateTime uint64, provisionFlag int32, body string, nfProfile *nrfschema.TNFProfile, overrideInfo string, supiVersion int64, gpsiVersion int64) string {
	if nfProfile == nil {
		log.Warnf("nfProfile is a nil pointer.")
		return ""
	}

	nfProfileMd5 := nfProfile.GenerateMd5()
	if nfProfileMd5 == "" {
		log.Warnf("Generate md5 for NF profile of nfInstance %s failed.", nfProfile.NfInstanceId)
		return ""
	}

	var nfServiceMd5 map[string]string
	if nfProfile.NfServices != nil {
		nfServiceMd5 = make(map[string]string)
		for _, item := range nfProfile.NfServices {
			if item.ServiceInstanceId != "" {
				serviceMd5 := item.GenerateMd5()
				if serviceMd5 == "" {
					log.Warnf("Generate md5 for NF service of serviceInstanceId %s failed.", item.ServiceInstanceId)
					return ""
				}
				nfServiceMd5[item.ServiceInstanceId] = serviceMd5
			}
		}
	}

	md5sum := fmt.Sprintf(`"nfProfile": "%s"`, nfProfileMd5)

	if nfServiceMd5 != nil {
		for serviceInstanceId, serviceMd5 := range nfServiceMd5 {
			md5sum = fmt.Sprintf(`%s, "%s": "%s"`, md5sum, serviceInstanceId, serviceMd5)
		}
	}

	//helperInfo := nfProfile.CreateHelperInfo(0)

	var targetNFProfile string
	if len(overrideInfo) == 0 {
		targetNFProfile = fmt.Sprintf(constvalue.TargetNFProfile, expiredTime, lastUpdateTime, profileUpdateTime, provisionFlag, md5sum, body, supiVersion, gpsiVersion)
	} else {
		targetNFProfile = fmt.Sprintf(constvalue.TargetNFProfileWithOverride, expiredTime, lastUpdateTime, profileUpdateTime, provisionFlag, md5sum, body, overrideInfo, supiVersion, gpsiVersion)
	}
	//fmt.Sprintf(`{"helper":{%s}}`, helperInfo)
	return targetNFProfile
}

//BuildStringSearchParameter build expression for grpc get request
func BuildStringSearchParameter(path, value string, operation int32) *common.MetaExpression {

	searchAttribute := &common.SearchAttribute{
		Name:      path,
		Operation: operation,
	}

	str := &common.StringValue{
		Value: value,
	}

	str_value := &common.SearchValue_Str{
		Str: str,
	}

	searchValue := &common.SearchValue{
		Data: str_value,
	}

	searchParameter := &common.SearchParameter{
		Attribute: searchAttribute,
		Value:     searchValue,
	}

	return BuildSearchParameter(searchParameter)

}

//BuildSearchParameter build expression for grpc get request
func BuildSearchParameter(searchParameter *common.SearchParameter) *common.MetaExpression {

	metaExpressionSearchParameter := &common.MetaExpression_SearchParameter{
		SearchParameter: searchParameter,
	}

	metaExpression := &common.MetaExpression{
		Data: metaExpressionSearchParameter,
	}

	return metaExpression
}

//ConstructGRPCGetRequestFilter construct a grpc get request filter
func ConstructGRPCGetRequestFilter(expiredTimeStart, expiredTimeEnd, lastUpdateTimeStart, lastUpdateTimeEnd uint64, provisioned int32, provVersion *nfprofile.ProvVersion, searchExpression *common.SearchExpression) *nfprofile.NFProfileGetRequest_Filter {
	nfProfileFilter := &nfprofile.NFProfileFilter{
		Provisioned: provisioned,
	}

	if expiredTimeStart < expiredTimeEnd {
		expiredTimeRange := &nfprofile.Range{
			Start: expiredTimeStart,
			End:   expiredTimeEnd,
		}
		nfProfileFilter.ExpiredTimeRange = expiredTimeRange
	}

	if lastUpdateTimeStart < lastUpdateTimeEnd {
		lastUpdateTimeRange := &nfprofile.Range{
			Start: lastUpdateTimeStart,
			End:   lastUpdateTimeEnd,
		}
		nfProfileFilter.LastUpdateTimeRange = lastUpdateTimeRange
	}
	supiVersion := provVersion.GetSupiVersion()
	gpsiVersion := provVersion.GetGpsiVersion()
	if supiVersion > 0 || gpsiVersion > 0 {
		provVersion := &nfprofile.ProvVersion{
			SupiVersion: supiVersion,
			GpsiVersion: gpsiVersion,
		}
		nfProfileFilter.ProvVersion = provVersion
	}

	nfProfileFilter.SearchExpression = searchExpression

	nfProfileFilterData := &nfprofile.NFProfileGetRequest_Filter{
		Filter: nfProfileFilter,
	}

	return nfProfileFilterData
}

//ConstructGRPCGetRequestFilter construct a grpc get request filter
func ConstructNFProfileCountGetRequest(expiredTimeStart, expiredTimeEnd, lastUpdateTimeStart, lastUpdateTimeEnd uint64, provisioned int32, searchExpression *common.SearchExpression) *nfprofile.NFProfileCountGetRequest {
	nfProfileFilter := &nfprofile.NFProfileFilter{
		Provisioned: provisioned,
	}

	if expiredTimeStart < expiredTimeEnd {
		expiredTimeRange := &nfprofile.Range{
			Start: expiredTimeStart,
			End:   expiredTimeEnd,
		}
		nfProfileFilter.ExpiredTimeRange = expiredTimeRange
	}

	if lastUpdateTimeStart < lastUpdateTimeEnd {
		lastUpdateTimeRange := &nfprofile.Range{
			Start: lastUpdateTimeStart,
			End:   lastUpdateTimeEnd,
		}
		nfProfileFilter.LastUpdateTimeRange = lastUpdateTimeRange
	}

	nfProfileFilter.SearchExpression = searchExpression

	nfProfileCountGetRequest := &nfprofile.NFProfileCountGetRequest{
		Filter: nfProfileFilter,
	}

	return nfProfileCountGetRequest
}

// GetNFInfoByFilter retrieves nfInfo(e.g. ausfInfo, amfInfo, ...) from DB
func GetNFInfoByFilter(filters map[string]*NFProfileQueryFilter, paths map[string]string) ([]string, error) {
	if len(paths) < 1 {
		return nil, fmt.Errorf("parameter path is necessary")
	}

	queryReq := &dbproxy.QueryRequest{
		RegionName: configmap.DBNfprofileRegionName,
	}
	for key, value := range filters {
		path, ok := paths[key]
		if !ok || len(path) < 1 {
			continue
		}
		oql := constructNFProfileOQL(value, path, false)
		queryReq.Query = append(queryReq.Query, oql)
	}

	if len(queryReq.Query) < 1 {
		return nil, fmt.Errorf("no oql applied")
	}

	nfProfileResponse, err := dbmgmt.QueryWithFilter(queryReq)
	if err != nil {
		return nil, fmt.Errorf("DB error, %v", err)
	}

	if nfProfileResponse.Code != dbmgmt.DbDataNotExist && nfProfileResponse.Code != dbmgmt.DbGetSuccess {
		return nil, fmt.Errorf("DB error, DB code is %d", nfProfileResponse.Code)
	}

	if nfProfileResponse.Code == dbmgmt.DbDataNotExist {
		fmt.Println("DbDataNotExist")
		return nil, nil
	}

	return nfProfileResponse.Value, nil
}

// GetNFProfileByFilter retrieves NF Profiles from DB
func GetNFProfileByFilter(filter *NFProfileQueryFilter) ([]string, error) {
	oql := constructNFProfileOQL(filter, "value", false)

	queryReq := &dbproxy.QueryRequest{
		RegionName: configmap.DBNfprofileRegionName,
		Query:      []string{oql},
	}

	nfProfileResponse, err := dbmgmt.QueryWithFilter(queryReq)
	if err != nil {
		return nil, fmt.Errorf("DB error, %v", err)
	}

	if nfProfileResponse.Code != dbmgmt.DbDataNotExist && nfProfileResponse.Code != dbmgmt.DbGetSuccess {
		return nil, fmt.Errorf("DB error, DB code is %d", nfProfileResponse.Code)
	}

	if nfProfileResponse.Code == dbmgmt.DbDataNotExist {
		return nil, nil
	}

	return nfProfileResponse.Value, nil
}

// GetNFProfileCountByFilter retrieves count of NF profiles from DB
func GetNFProfileCountByFilter(filter *NFProfileQueryFilter) (int, error) {
	oql := constructNFProfileOQL(filter, "", true)

	queryReq := &dbproxy.QueryRequest{
		RegionName: configmap.DBNfprofileRegionName,
		Query:      []string{oql},
	}

	nfProfileResponse, err := dbmgmt.QueryWithFilter(queryReq)
	if err != nil {
		return 0, fmt.Errorf("DB error, %v", err)
	}

	if nfProfileResponse.Code != dbmgmt.DbDataNotExist && nfProfileResponse.Code != dbmgmt.DbGetSuccess {
		return 0, fmt.Errorf("DB error, DB code is %d", nfProfileResponse.Code)
	}

	if nfProfileResponse.Code == dbmgmt.DbDataNotExist {
		return 0, nil
	}

	count := 0
	if len(nfProfileResponse.Value) > 0 {
		var err error
		count, err = strconv.Atoi(nfProfileResponse.Value[0])
		if err != nil {
			return 0, fmt.Errorf("strconv.Atoi error, %v", err)
		}
	}

	return count, nil
}

// constructNFProfileOQL constructs OQL for DB query
// filter: query criteria
// queryContent: what to retrieve
// count: a flag indicating whether retrieve the count or not, when set to true, queryContent is ignored
func constructNFProfileOQL(filter *NFProfileQueryFilter, queryContent string, count bool) string {
	var oql string

	if !count {
		oql = "SELECT " + queryContent + " FROM /ericsson-nrf-nfprofiles.entrySet"
	} else {
		oql = "SELECT COUNT(*) FROM /ericsson-nrf-nfprofiles.entrySet"
	}

	andFlag := false

	if filter != nil {
		if filter.ExpiredTimeStart < filter.ExpiredTimeEnd {
			if !andFlag {
				oql += " WHERE value.expiredTime >= " + fmt.Sprintf("%dL", filter.ExpiredTimeStart) +
					" AND value.expiredTime < " + fmt.Sprintf("%dL", filter.ExpiredTimeEnd)
				andFlag = true
			} else {
				oql += " AND value.expiredTime >= " + fmt.Sprintf("%dL", filter.ExpiredTimeStart) +
					" AND value.expiredTime < " + fmt.Sprintf("%dL", filter.ExpiredTimeEnd)
			}

		}

		if filter.LastUpdateTimeStart < filter.LastUpdateTimeEnd {
			if !andFlag {
				oql += " WHERE value.lastUpdateTime >= " + fmt.Sprintf("%dL", filter.LastUpdateTimeStart) +
					" AND value.lastUpdateTime < " + fmt.Sprintf("%dL", filter.LastUpdateTimeEnd)
				andFlag = true
			} else {
				oql += " AND value.lastUpdateTime >= " + fmt.Sprintf("%dL", filter.LastUpdateTimeStart) +
					" AND value.lastUpdateTime < " + fmt.Sprintf("%dL", filter.LastUpdateTimeEnd)
			}
		}

		if filter.Provisioned > 0 {
			if !andFlag {
				oql += " WHERE value.provisioned = " + fmt.Sprintf("%d", filter.Provisioned)
				andFlag = true
			} else {
				oql += " AND value.provisioned = " + fmt.Sprintf("%d", filter.Provisioned)
			}
		}

		if filter.ProvVersion != nil {
			if !andFlag {
				oql += " WHERE (value.provSupiVersion >= " + fmt.Sprintf("%dL", filter.ProvVersion.SupiVersion) +
					" OR value.provGpsiVersion >= " + fmt.Sprintf("%dL", filter.ProvVersion.GpsiVersion) + ")"
				andFlag = true
			} else {
				oql += " AND (value.provSupiVersion >= " + fmt.Sprintf("%dL", filter.ProvVersion.SupiVersion) +
					" OR value.provGpsiVersion >= " + fmt.Sprintf("%dL", filter.ProvVersion.GpsiVersion) + ")"
			}
		}

		for _, item := range filter.QueryList {
			if !andFlag {
				oql += " WHERE " + item.Key + " = " + item.Value
				andFlag = true
			} else {
				oql += " AND " + item.Key + " = " + item.Value
			}
		}
	}

	return oql
}

//RebuildNfServiceOverrideInfo is for rebuild nfService overrideInfo
func RebuildNfServiceOverrideInfo(oldMapper map[string]int, newMapper map[string]int, overrideInfo []nrfschema.OverrideInfo) []nrfschema.OverrideInfo {
	newOverrideInfo := make([]nrfschema.OverrideInfo, 0)

	oldMapperReverse := make(map[int]string, 0)
	newMapperReverse := make(map[int]string, 0)

	for k, v := range oldMapper {
		oldMapperReverse[v] = k
	}

	for k, v := range newMapper {
		newMapperReverse[v] = k
	}

	for _, overrideItem := range overrideInfo {
		overrideItemNew := nrfschema.OverrideInfo{}

		path := overrideItem.Path
		match, _ := regexp.MatchString("^/nfServices/\\d+/", path)
		if match {
			paths := strings.Split(path, "/")
			log.Debugf("items index: %v\n", paths[2])
			index, err := strconv.Atoi(paths[2])
			if err != nil {
				log.Errorf("Convert %s to integer failure, will remove the overrideAttribute, err:%s", paths[2], err.Error())
				continue
			}
			serviceID, ok := oldMapperReverse[index]
			if !ok {
				log.Warnf("Origin nfPrifile overrideInfo nfService array index not correct, will remove the overrideAttribute")
				continue
			}

			indexNew, ok := newMapper[serviceID]
			if !ok {
				log.Warnf("Update nfProfile delete nfService %s, will remove the related overrideAttributes", serviceID)
				continue
			}

			if index != indexNew {
				log.Warnf("Update nfProfile nfService index is changed, will update the related overrideInfo index")
				paths[2] = fmt.Sprintf("%d", indexNew)
				path = strings.Join(paths, "/")
			}

			overrideItemNew.Path = path
			overrideItemNew.Action = overrideItem.Action
			overrideItemNew.Value = overrideItem.Value

			newOverrideInfo = append(newOverrideInfo, overrideItemNew)
		} else {
			overrideItemNew.Path = path
			overrideItemNew.Action = overrideItem.Action
			overrideItemNew.Value = overrideItem.Value

			newOverrideInfo = append(newOverrideInfo, overrideItemNew)
		}
	}
	return newOverrideInfo
}

//SupiRangeInjection is for get SupiRange from nfProfile
func SupiRangeInjection(nfProfile *nrfschema.TNFProfile) []nrfschema.TSupiRange {
	switch nfProfile.NfType {
	case constvalue.NfTypeUDM, constvalue.NfTypeNRFUDM:
		if nfProfile.UdmInfo != nil {
			return nfProfile.UdmInfo.SupiRanges
		}
	case constvalue.NfTypeUDR, constvalue.NfTypeNRFUDR:
		if nfProfile.UdrInfo != nil {
			return nfProfile.UdrInfo.SupiRanges
		}
	case constvalue.NfTypeAUSF, constvalue.NfTypeNRFAUSF:
		if nfProfile.AusfInfo != nil {
			return nfProfile.AusfInfo.SupiRanges
		}
	case constvalue.NfTypePCF, constvalue.NfTypeNRFPCF:
		if nfProfile.PcfInfo != nil {
			return nfProfile.PcfInfo.SupiRanges
		}
	case constvalue.NfTypeCHF, constvalue.NfTypeNRFCHF:
		if nfProfile.ChfInfo != nil {
			return nfProfile.ChfInfo.SupiRangeList
		}
	default:
		return nil
	}
	return nil
}

//GpsiRangeInjection is for get GpsiRange from nfProfile
func GpsiRangeInjection(nfProfile *nrfschema.TNFProfile) []nrfschema.TIdentityRange {
	switch nfProfile.NfType {
	case constvalue.NfTypeUDM, constvalue.NfTypeNRFUDM:
		if nfProfile.UdmInfo != nil {
			return nfProfile.UdmInfo.GpsiRanges
		}
	case constvalue.NfTypeUDR, constvalue.NfTypeNRFUDR:
		if nfProfile.UdrInfo != nil {
			return nfProfile.UdrInfo.GpsiRanges
		}
	case constvalue.NfTypeCHF, constvalue.NfTypeNRFCHF:
		if nfProfile.ChfInfo != nil {
			return nfProfile.ChfInfo.GpsiRangeList
		}
	case constvalue.NfTypePCF, constvalue.NfTypeNRFPCF:
		if nfProfile.PcfInfo != nil {
			return nfProfile.PcfInfo.GpsiRanges
		}
	default:
		return nil
	}
	return nil
}

// CaculateHeartBeatTimer is to Caculate HeartBeat timer
func CaculateHeartBeatTimer(heartBeatTimerInProfile, defaultHeartBeatTimer int) int {
	minRange := int(math.Max(float64(defaultHeartBeatTimer-internalconf.HeartBeatTimerOffset), float64(internalconf.HeartBeatTimerMin)))
	maxRange := defaultHeartBeatTimer + internalconf.HeartBeatTimerOffset

	if minRange <= heartBeatTimerInProfile && heartBeatTimerInProfile <= maxRange {
		return heartBeatTimerInProfile
	}
	return defaultHeartBeatTimer
}

//SplitNrfInfo is to split nrfInfo in nrfprofile
func SplitNrfInfo(nrfInstanceId string, nrfInfo *nrfschema.TNrfInfo) []string {
	var nfProfiles []string
	if nrfInfo == nil {
		return nfProfiles
	}
	if nrfInfo.ServedAmfInfo != nil && len(nrfInfo.ServedAmfInfo) > 0 {
		amfProfiles := constructNFProfileForAmf(nrfInstanceId, nrfInfo.ServedAmfInfo)
		nfProfiles = append(nfProfiles, amfProfiles...)
	}
	if nrfInfo.ServedAusfInfo != nil && len(nrfInfo.ServedAusfInfo) > 0 {
		ausfProfiles := constructNFProfileForAusf(nrfInstanceId, nrfInfo.ServedAusfInfo)
		nfProfiles = append(nfProfiles, ausfProfiles...)
	}
	if nrfInfo.ServedBsfInfo != nil && len(nrfInfo.ServedBsfInfo) > 0 {
		bsfProfiles := constructNFProfileForBsf(nrfInstanceId, nrfInfo.ServedBsfInfo)
		nfProfiles = append(nfProfiles, bsfProfiles...)
	}
	if nrfInfo.ServedChfInfo != nil && len(nrfInfo.ServedChfInfo) > 0 {
		chfProfiles := constructNFProfileForChf(nrfInstanceId, nrfInfo.ServedChfInfo)
		nfProfiles = append(nfProfiles, chfProfiles...)
	}
	if nrfInfo.ServedPcfInfo != nil && len(nrfInfo.ServedPcfInfo) > 0 {
		pcfProfiles := constructNFProfileForPcf(nrfInstanceId, nrfInfo.ServedPcfInfo)
		nfProfiles = append(nfProfiles, pcfProfiles...)
	}
	if nrfInfo.ServedSmfInfo != nil && len(nrfInfo.ServedSmfInfo) > 0 {
		smfProfiles := constructNFProfileForSmf(nrfInstanceId, nrfInfo.ServedSmfInfo)
		nfProfiles = append(nfProfiles, smfProfiles...)
	}
	if nrfInfo.ServedUdmInfo != nil && len(nrfInfo.ServedUdmInfo) > 0 {
		udmProfiles := constructNFProfileForUdm(nrfInstanceId, nrfInfo.ServedUdmInfo)
		nfProfiles = append(nfProfiles, udmProfiles...)
	}
	if nrfInfo.ServedUdrInfo != nil && len(nrfInfo.ServedUdrInfo) > 0 {
		udrProfiles := constructNFProfileForUdr(nrfInstanceId, nrfInfo.ServedUdrInfo)
		nfProfiles = append(nfProfiles, udrProfiles...)
	}
	if nrfInfo.ServedUpfInfo != nil && len(nrfInfo.ServedUpfInfo) > 0 {
		upfProfiles := constructNFProfileForUpf(nrfInstanceId, nrfInfo.ServedUpfInfo)
		nfProfiles = append(nfProfiles, upfProfiles...)
	}
	return nfProfiles
}

//ConstructNrfInfoFromArray is to construct nrfinfo with nfprofile array
func ConstructNrfInfoFromArray(nfProfileArray []string) string {
	var nfProfileStr string
	for _, nfprofile := range nfProfileArray {
		if nfProfileStr == "" {
			nfProfileStr = nfprofile
		} else {
			nfProfileStr += "," + nfprofile
		}
	}
	return "[" + nfProfileStr + "]"
}

//constructNFProfileForAmf is to construct amf profile from nrfInfo servedAmfInfo
func constructNFProfileForAmf(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TAmfInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		amfInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal amfInfo error, err=%v", err)
			continue
		}
		nfProfile.AmfInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypeAMF
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypeAMF, constvalue.AmfInfo, amfInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//constructNFProfileForAusf is to construct ausf profile from nrfInfo servedAusfInfo
func constructNFProfileForAusf(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TAusfInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		ausfInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal ausfInfo error, err=%v", err)
			continue
		}
		nfProfile.AusfInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypeAUSF
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypeAUSF, constvalue.AusfInfo, ausfInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//constructNFProfileForBsf is to construct bsf profile from nrfInfo servedBsfInfo
func constructNFProfileForBsf(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TBsfInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		bsfInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal bsfInfo error, err=%v", err)
			continue
		}
		nfProfile.BsfInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypeBSF
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypeBSF, constvalue.BsfInfo, bsfInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//constructNFProfileForChf is to construct chf profile from nrfInfo servedChfInfo
func constructNFProfileForChf(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TChfInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		chfInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal chfInfo error, err=%v", err)
			continue
		}
		nfProfile.ChfInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypeCHF
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypeCHF, constvalue.ChfInfo, chfInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//constructNFProfileForPcf is to construct pcf profile from nrfInfo servedPcfInfo
func constructNFProfileForPcf(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TPcfInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		pcfInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal pcfInfo error, err=%v", err)
			continue
		}
		nfProfile.PcfInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypePCF
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypePCF, constvalue.PcfInfo, pcfInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//constructNFProfileForSmf is to construct smf profile from nrfInfo servedSmfInfo
func constructNFProfileForSmf(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TSmfInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		smfInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal smfInfo error, err=%v", err)
		}
		nfProfile.SmfInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypeSMF
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypeSMF, constvalue.SmfInfo, smfInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//constructNFProfileForUdm is to construct udm profile from nrfInfo servedUdmInfo
func constructNFProfileForUdm(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TUdmInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		udmInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal udmInfo error, err=%v", err)
		}
		nfProfile.UdmInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypeUDM
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypeUDM, constvalue.UdmInfo, udmInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//constructNFProfileForUdr is to construct udr profile from nrfInfo servedUdrInfo
func constructNFProfileForUdr(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TUdrInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		udrInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal udrInfo error, err=%v", err)
		}
		nfProfile.UdrInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypeUDR
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypeUDR, constvalue.UdrInfo, udrInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//constructNFProfileForUpf is to construct upf profile from nrfInfo servedUpfInfo
func constructNFProfileForUpf(nrfInstanceId string, servedNfInfo map[string]*nrfschema.TUpfInfo) []string {
	var nfProfiles []string
	for nfInstanceId, nfInfo := range servedNfInfo {
		nfProfile := &nrfschema.TNFProfile{}
		upfInfoStr, err := json.Marshal(nfInfo)
		if err != nil {
			log.Errorf("marshal upfInfo error, err=%v", err)
		}
		nfProfile.UpfInfo = nfInfo
		nfProfile.NfInstanceId = nfInstanceId
		nfProfile.NfType = constvalue.NfTypeUPF
		profileUpdateTime := GenerateProfileUpdateTime()
		nfHelper := nfProfile.CreateHelperInfo(profileUpdateTime)
		supiVersion, gpsiVersion := getSupiGpsiVersion(nfProfile, profileUpdateTime)
		nfProfileStr := fmt.Sprintf(constvalue.SplitNFProfile, nfInstanceId, nrfInstanceId, constvalue.NfTypeUPF, constvalue.UpfInfo, upfInfoStr, nfHelper, supiVersion, gpsiVersion)
		nfProfiles = append(nfProfiles, nfProfileStr)
	}
	return nfProfiles
}

//ConstructNrfInfoFromRegion is to collect nfInfo from db region ericsson-nrf-regionnfinfo by nrfInstanceId and construct nrfInfo
func ConstructNrfInfoFromRegion(nrfInstanceId string) (*nrfschema.TNrfInfo, error) {
	nrfInfo := &nrfschema.TNrfInfo{}
	searchOql := fmt.Sprintf("select value.body from /ericsson-nrf-regionnfinfo.entrySet where value.body.nrfInstanceId='%s'", nrfInstanceId)
	queryRequest := &dbproxy.QueryRequest{
		RegionName: configmap.DBRegionNfInfoRegionName,
		Query:      []string{searchOql},
	}
	nfInfoResp, err := dbmgmt.QueryWithFilter(queryRequest)
	if err != nil {
		log.Errorf("fail to get region nfInfo by nrfInstanceId=%s, error=%v", nrfInstanceId, err)
		return nil, err
	}
	if nfInfoResp.GetCode() == dbmgmt.DbDataNotExist {
		log.Debugf("not found any nfProfile in nrfInfo by nrfInstanceId=%s", nrfInstanceId)
		return nil, nil
	}
	if nfInfoResp.GetCode() != dbmgmt.DbGetSuccess {
		log.Errorf("fail to get nfInfo, return code=%v", nfInfoResp.GetCode())
		return nil, fmt.Errorf("fail to get nfInfo in nrfInfo by nrfInstanceId=%s, return code=%d", nrfInstanceId, nfInfoResp.GetCode())
	}
	nfInfos := nfInfoResp.GetValue()
	for _, item := range nfInfos {
		insertNfInfoToNrfInfo(nrfInfo, item)
	}
	return nrfInfo, nil
}

//insertNfInfoToNrfInfo is to insert a nfInfo into nrfInfo
func insertNfInfoToNrfInfo(nrfInfo *nrfschema.TNrfInfo, nfProfile string) {
	nfProfileStruct := &nrfschema.TNFProfile{}
	err := json.Unmarshal([]byte(nfProfile), nfProfileStruct)
	if err != nil {
		log.Debugf("fail to unmarshal nfProfile, error=%v", err)
		return
	}
	switch nfProfileStruct.NfType {
	case constvalue.NfTypeUDR:
		if nrfInfo.ServedUdrInfo == nil {
			nrfInfo.ServedUdrInfo = make(map[string]*nrfschema.TUdrInfo)
		}
		nrfInfo.ServedUdrInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.UdrInfo
	case constvalue.NfTypeUDM:
		if nrfInfo.ServedUdmInfo == nil {
			nrfInfo.ServedUdmInfo = make(map[string]*nrfschema.TUdmInfo)
		}
		nrfInfo.ServedUdmInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.UdmInfo
	case constvalue.NfTypeAUSF:
		if nrfInfo.ServedAusfInfo == nil {
			nrfInfo.ServedAusfInfo = make(map[string]*nrfschema.TAusfInfo)
		}
		nrfInfo.ServedAusfInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.AusfInfo
	case constvalue.NfTypeAMF:
		if nrfInfo.ServedAmfInfo == nil {
			nrfInfo.ServedAmfInfo = make(map[string]*nrfschema.TAmfInfo)
		}
		nrfInfo.ServedAmfInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.AmfInfo
	case constvalue.NfTypeSMF:
		if nrfInfo.ServedSmfInfo == nil {
			nrfInfo.ServedSmfInfo = make(map[string]*nrfschema.TSmfInfo)
		}
		nrfInfo.ServedSmfInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.SmfInfo
	case constvalue.NfTypeUPF:
		if nrfInfo.ServedUpfInfo == nil {
			nrfInfo.ServedUpfInfo = make(map[string]*nrfschema.TUpfInfo)
		}
		nrfInfo.ServedUpfInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.UpfInfo
	case constvalue.NfTypePCF:
		if nrfInfo.ServedPcfInfo == nil {
			nrfInfo.ServedPcfInfo = make(map[string]*nrfschema.TPcfInfo)
		}
		nrfInfo.ServedPcfInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.PcfInfo
	case constvalue.NfTypeBSF:
		if nrfInfo.ServedBsfInfo == nil {
			nrfInfo.ServedBsfInfo = make(map[string]*nrfschema.TBsfInfo)
		}
		nrfInfo.ServedBsfInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.BsfInfo
	case constvalue.NfTypeCHF:
		if nrfInfo.ServedChfInfo == nil {
			nrfInfo.ServedChfInfo = make(map[string]*nrfschema.TChfInfo)
		}
		nrfInfo.ServedChfInfo[nfProfileStruct.NfInstanceId] = nfProfileStruct.ChfInfo
	default:
		log.Warningf("unknown nftype=%s", nfProfileStruct.NfType)
	}
}

//getSupiGpsiVersion is to get supi,gpsi version from nfprofile
func getSupiGpsiVersion(nfProfile *nrfschema.TNFProfile, profileUpdateTime uint64) (int64, int64) {
	var supiVersion int64
	var gpsiVersion int64
	supiExist := SupiRangeInfoExist(nfProfile)
	gpsiExist := GpsiRangeInfoExist(nfProfile)
	if supiExist {
		supiVersion = int64(profileUpdateTime)
	} else {
		supiVersion = 0
	}
	if gpsiExist {
		gpsiVersion = int64(profileUpdateTime)
	} else {
		gpsiVersion = 0
	}
	return supiVersion, gpsiVersion
}

//SupiRangeInfoExist check whether supiInfo exist
func SupiRangeInfoExist(nfProfile *nrfschema.TNFProfile) bool {
	switch nfProfile.NfType {
	case constvalue.NfTypeUDM:
		if nfProfile.UdmInfo != nil {
			return len(nfProfile.UdmInfo.SupiRanges) != 0
		}
	case constvalue.NfTypeUDR:
		if nfProfile.UdrInfo != nil {
			return len(nfProfile.UdrInfo.SupiRanges) != 0
		}
	case constvalue.NfTypeAUSF:
		if nfProfile.AusfInfo != nil {
			return len(nfProfile.AusfInfo.SupiRanges) != 0
		}
	case constvalue.NfTypePCF:
		if nfProfile.PcfInfo != nil {
			return len(nfProfile.PcfInfo.SupiRanges) != 0
		}
	case constvalue.NfTypeCHF:
		if nfProfile.ChfInfo != nil {
			return len(nfProfile.ChfInfo.SupiRangeList) != 0
		}
	}

	return false
}

//GpsiRangeInfoExist check whether gpsiInfo exist
func GpsiRangeInfoExist(nfProfile *nrfschema.TNFProfile) bool {
	switch nfProfile.NfType {
	case constvalue.NfTypeUDM:
		if nfProfile.UdmInfo != nil {
			return len(nfProfile.UdmInfo.GpsiRanges) != 0
		}
	case constvalue.NfTypeUDR:
		if nfProfile.UdrInfo != nil {
			return len(nfProfile.UdrInfo.GpsiRanges) != 0
		}
	case constvalue.NfTypeCHF:
		if nfProfile.ChfInfo != nil {
			return len(nfProfile.ChfInfo.GpsiRangeList) != 0
		}
	case constvalue.NfTypePCF:
		if nfProfile.PcfInfo != nil {
			return len(nfProfile.PcfInfo.GpsiRanges) != 0
		}
	}

	return false
}

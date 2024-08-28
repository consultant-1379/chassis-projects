package profileop

import (
	"com/dbproxy"
	"com/dbproxy/nfmessage/common"
	"com/dbproxy/nfmessage/nfprofile"
	"com/dbproxy/nfmessage/nrfaddress"
	"com/dbproxy/nfmessage/nrfprofile"
	"encoding/json"
	"fmt"
	"math"
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
		Scheme: cm.NrfCommon.RemoteDefaultSetting.Scheme,
		Fqdn:   fqdn,
		Port:   cm.NrfCommon.RemoteDefaultSetting.Port,
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

// ConstructNRFProfileIndex is to construct NRF profile index
func ConstructNRFProfileIndex(nrfProfile []byte) (*nrfprofile.NRFProfileIndex, *problemdetails.ProblemDetails) {
	profileIndex := &nrfprofile.NRFProfileIndex{}

	//Construct the common key for NRFProfile Index
	constructNrfCommonIndex(profileIndex)

	//Construct the specific key for nfProfile nrfInfo Index
	nrfInfo, problemDetails := getNrfInfo(nrfProfile)
	if problemDetails != nil {
		return nil, problemDetails
	}

	problemDetails = constructNrfAmfIndex(nrfInfo, profileIndex)
	if problemDetails != nil {
		log.Warnf(problemDetails.Title)
	}
	problemDetails = constructNrfAusfIndex(nrfInfo, profileIndex)
	if problemDetails != nil {
		log.Warnf(problemDetails.Title)
	}
	problemDetails = constructNrfPcfIndex(nrfInfo, profileIndex)
	if problemDetails != nil {
		log.Warnf(problemDetails.Title)
	}
	problemDetails = constructNrfSmfIndex(nrfInfo, profileIndex)
	if problemDetails != nil {
		log.Warnf(problemDetails.Title)
	}
	problemDetails = constructNrfUdmIndex(nrfInfo, profileIndex)
	if problemDetails != nil {
		log.Warnf(problemDetails.Title)
	}

	return profileIndex, nil
}

func constructNrfCommonIndex(profileIndex *nrfprofile.NRFProfileIndex) *problemdetails.ProblemDetails {
	profileIndex.Key1 = math.MaxInt64 - 1

	return nil
}

//nrf profileIndex
func constructNrfAmfIndex(nrfInfo []byte, profileIndex *nrfprofile.NRFProfileIndex) *problemdetails.ProblemDetails {
	guamiList, taiList, regionIDList, setIDList, problemDetails := getNrfInfoAmfProperties(nrfInfo)
	if problemDetails != nil {
		return problemDetails
	}

	if guamiList != nil {
		profileIndex.AmfKey1 = guamiList
	}

	if taiList != nil {
		profileIndex.AmfKey2 = taiList
	}

	if regionIDList != nil {
		profileIndex.AmfKey3 = regionIDList
	}

	if setIDList != nil {
		profileIndex.AmfKey4 = setIDList
	}

	return nil
}

func constructNrfSmfIndex(nrfInfo []byte, profileIndex *nrfprofile.NRFProfileIndex) *problemdetails.ProblemDetails {
	dnnList, pgwFqdnList, taiList, problemDetails := getNrfInfoSmfProperties(nrfInfo)
	if problemDetails != nil {
		return problemDetails
	}

	if dnnList != nil {
		profileIndex.SmfKey1 = dnnList
	}

	if pgwFqdnList != nil {
		profileIndex.SmfKey2 = pgwFqdnList
	}

	if taiList != nil {
		profileIndex.SmfKey3 = taiList
	}

	return nil
}

func constructNrfUdmIndex(nrfInfo []byte, profileIndex *nrfprofile.NRFProfileIndex) *problemdetails.ProblemDetails {
	groupID, routingIndicator, problemDetails := getNrfInfoUdmProperties(nrfInfo)
	if problemDetails != nil {
		return problemDetails
	}

	if groupID != nil {
		profileIndex.UdmKey1 = groupID
	}

	if routingIndicator != nil {
		profileIndex.UdmKey2 = routingIndicator
	}

	return nil
}

func constructNrfAusfIndex(nrfInfo []byte, profileIndex *nrfprofile.NRFProfileIndex) *problemdetails.ProblemDetails {
	groupIDList, routingIndicatorList, problemDetails := getNrfInfoAusfProperties(nrfInfo)
	if problemDetails != nil {
		return problemDetails
	}

	if groupIDList != nil {
		profileIndex.AusfKey1 = groupIDList
	}

	if routingIndicatorList != nil {
		profileIndex.AusfKey2 = routingIndicatorList
	}

	return nil
}

func constructNrfPcfIndex(nrfInfo []byte, profileIndex *nrfprofile.NRFProfileIndex) *problemdetails.ProblemDetails {
	dnnList, groupIDList, problemDetails := getNrfInfoPcfProperties(nrfInfo)
	if problemDetails != nil {
		return problemDetails
	}

	if dnnList != nil {
		profileIndex.PcfKey1 = dnnList
	}

	if groupIDList != nil {
		profileIndex.PcfKey2 = groupIDList
	}

	return nil
}

func getNrfInfo(nrfProfile []byte) ([]byte, *problemdetails.ProblemDetails) {
	nrfInfo, dataType, _, err := jsonparser.Get(nrfProfile, constvalue.NrfInfo)
	if dataType == jsonparser.NotExist {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("NrfProfile less %s configuration", constvalue.NrfInfo),
		}
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s", constvalue.NrfInfo),
		}
	}

	return nrfInfo, nil
}

func getPcfInfoSum(nrfInfo []byte) ([]byte, *problemdetails.ProblemDetails) {
	pcfInfoSum, dataType, _, err := jsonparser.Get(nrfInfo, constvalue.PcfInfoSum)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing nrfInfo failed for %s", constvalue.PcfInfoSum),
		}
	}

	return pcfInfoSum, nil
}

func getAusfInfoSum(nrfInfo []byte) ([]byte, *problemdetails.ProblemDetails) {
	ausfInfoSum, dataType, _, err := jsonparser.Get(nrfInfo, constvalue.AusfInfoSum)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing nrfInfo failed for %s", constvalue.AusfInfoSum),
		}
	}

	return ausfInfoSum, nil
}

func getUdmInfoSum(nrfInfo []byte) ([]byte, *problemdetails.ProblemDetails) {
	udmInfoSum, dataType, _, err := jsonparser.Get(nrfInfo, constvalue.UdmInfoSum)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing nrfInfo failed for %s", constvalue.UdmInfoSum),
		}
	}

	return udmInfoSum, nil
}

func getSmfInfoSum(nrfInfo []byte) ([]byte, *problemdetails.ProblemDetails) {
	smfInfoSum, dataType, _, err := jsonparser.Get(nrfInfo, constvalue.SmfInfoSum)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing nrfInfo failed for %s", constvalue.SmfInfoSum),
		}
	}

	return smfInfoSum, nil
}

func getAmfInfoSum(nrfInfo []byte) ([]byte, *problemdetails.ProblemDetails) {
	amfInfoSum, dataType, _, err := jsonparser.Get(nrfInfo, constvalue.AmfInfoSum)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing nrfInfo failed for %s", constvalue.AmfInfoSum),
		}
	}

	return amfInfoSum, nil
}

//groupId, routingIndicator
func getAusfProperties(nfProfile []byte) (string, string, *problemdetails.ProblemDetails) {
	ausfInfo, dataType, _, err := jsonparser.Get(nfProfile, constvalue.AusfInfo)
	if dataType == jsonparser.NotExist {
		return "", "", nil
	}

	if err != nil {
		return "", "", &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s", constvalue.AusfInfo),
		}
	}

	groupId, problemDetails := getAusfGroupId(ausfInfo)
	if problemDetails != nil {
		return "", "", problemDetails
	}

	routingIndicator, problemDetails := getAusfRoutingIndicator(ausfInfo)
	if problemDetails != nil {
		return "", "", problemDetails
	}

	return groupId, routingIndicator, nil
}

//groupIDList, routingIndicatorList
func getNrfInfoAusfProperties(nrfInfo []byte) ([]*nrfprofile.NRFKeyStruct, []*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	ausfInfoSum, problemDetails := getAusfInfoSum(nrfInfo)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	groupIDList, problemDetails := getAusfInfoSumGroupIDList(ausfInfoSum)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	routingIndicatorList, problemDetails := getAusfInfoSumRoutingIndicatorList(ausfInfoSum)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	return groupIDList, routingIndicatorList, nil
}

//groupIDList, routingIndicatorList
func getNrfInfoUdmProperties(nrfInfo []byte) ([]*nrfprofile.NRFKeyStruct, []*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	udmInfoSum, problemDetails := getUdmInfoSum(nrfInfo)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	groupIDList, problemDetails := getUdmInfoSumGroupIDList(udmInfoSum)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	routingIndicatorList, problemDetails := getUdmInfoSumRoutingIndicatorList(udmInfoSum)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	return groupIDList, routingIndicatorList, nil
}

//dnnList, groupIDList
func getNrfInfoPcfProperties(nrfInfo []byte) ([]*nrfprofile.NRFKeyStruct, []*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	pcfInfoSum, problemDetails := getPcfInfoSum(nrfInfo)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	dnnList, problemDetails := getPcfInfoSumDnnList(pcfInfoSum)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	groupIDList, problemDetails := getPcfInfoSumGroupIDList(pcfInfoSum)
	if problemDetails != nil {
		return nil, nil, problemDetails
	}

	return dnnList, groupIDList, nil
}

//dnnList, pgwFqdnList, taiList
func getNrfInfoSmfProperties(nrfInfo []byte) ([]*nrfprofile.NRFKeyStruct, []*nrfprofile.NRFKeyStruct, []*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	smfInfoSum, problemDetails := getSmfInfoSum(nrfInfo)
	if problemDetails != nil {
		return nil, nil, nil, problemDetails
	}

	dnnList, problemDetails := getSmfInfoSumDnnList(smfInfoSum)
	if problemDetails != nil {
		return nil, nil, nil, problemDetails
	}

	pgwFqdnList, problemDetails := getSmfInfoPgwFqdnList(smfInfoSum)
	if problemDetails != nil {
		return nil, nil, nil, problemDetails
	}

	taiList, problemDetails := getSmfInfoTaiList(smfInfoSum)
	if problemDetails != nil {
		return nil, nil, nil, problemDetails
	}

	return dnnList, pgwFqdnList, taiList, nil
}

//guamiList, taiList, regionIDList, setIDList
func getNrfInfoAmfProperties(nrfInfo []byte) ([]*nrfprofile.NRFKeyStruct, []*nrfprofile.NRFKeyStruct, []*nrfprofile.NRFKeyStruct, []*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	amfInfoSum, problemDetails := getAmfInfoSum(nrfInfo)
	if problemDetails != nil {
		return nil, nil, nil, nil, problemDetails
	}

	guamiList, problemDetails := getAmfInfoGuamiList(amfInfoSum)
	if problemDetails != nil {
		return nil, nil, nil, nil, problemDetails
	}

	taiList, problemDetails := getAmfInfoTaiList(amfInfoSum)
	if problemDetails != nil {
		return nil, nil, nil, nil, problemDetails
	}

	regionIDList, problemDetails := getAmfInfoAmfRegionIDList(amfInfoSum)
	if problemDetails != nil {
		return nil, nil, nil, nil, problemDetails
	}

	setIDList, problemDetails := getAmfInfoAmfSetIDList(amfInfoSum)
	if problemDetails != nil {
		return nil, nil, nil, nil, problemDetails
	}

	return guamiList, taiList, regionIDList, setIDList, nil
}

//groupId, routingIndicator
func getUdmProperties(nfProfile []byte) (string, string, *problemdetails.ProblemDetails) {
	udmInfo, dataType, _, err := jsonparser.Get(nfProfile, constvalue.UdmInfo)
	if dataType == jsonparser.NotExist {
		return "", "", nil
	}

	if err != nil {
		return "", "", &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s", constvalue.UdmInfo),
		}
	}

	groupId, problemDetails := getUdmGroupId(udmInfo)
	if problemDetails != nil {
		return "", "", problemDetails
	}

	routingIndicator, problemDetails := getUdmRoutingIndicator(udmInfo)
	if problemDetails != nil {
		return "", "", problemDetails
	}

	return groupId, routingIndicator, nil
}

//groupId
func getUdrProperties(nfProfile []byte) (string, *problemdetails.ProblemDetails) {
	udrInfo, dataType, _, err := jsonparser.Get(nfProfile, constvalue.UdrInfo)
	if dataType == jsonparser.NotExist {
		return "", nil
	}

	if err != nil {
		return "", &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s", constvalue.UdrInfo),
		}
	}

	groupId, problemDetails := getUdrGroupId(udrInfo)
	if problemDetails != nil {
		return "", problemDetails
	}

	return groupId, nil
}

func getAmfSetId(amfInfo []byte) (string, *problemdetails.ProblemDetails) {
	amfSetId, err := jsonparser.GetString(amfInfo, constvalue.AmfSetID)
	if err != nil {
		errorInfo := fmt.Sprintf("Can not find %s in %s", constvalue.AmfSetID, constvalue.AmfInfo)
		return "", &problemdetails.ProblemDetails{
			Title: "not a valid nf profile",
			InvalidParams: []*problemdetails.InvalidParam{
				&problemdetails.InvalidParam{
					Param:  constvalue.AmfInfo,
					Reason: errorInfo,
				},
			},
		}
	}
	return amfSetId, nil
}

func getAmfRegionId(amfInfo []byte) (string, *problemdetails.ProblemDetails) {
	amfRegionId, err := jsonparser.GetString(amfInfo, constvalue.AmfRegionID)
	if err != nil {
		errorInfo := fmt.Sprintf("Can not find %s in %s", constvalue.AmfRegionID, constvalue.AmfInfo)
		return "", &problemdetails.ProblemDetails{
			Title: "not a valid nf profile",
			InvalidParams: []*problemdetails.InvalidParam{
				&problemdetails.InvalidParam{
					Param:  constvalue.AmfInfo,
					Reason: errorInfo,
				},
			},
		}
	}
	return amfRegionId, nil
}

func getAusfGroupId(ausfInfo []byte) (string, *problemdetails.ProblemDetails) {
	groupId, err := jsonparser.GetString(ausfInfo, constvalue.GroupID)
	if err != nil {
		return "", nil
	}
	return groupId, nil
}

func getAusfRoutingIndicator(ausfInfo []byte) (string, *problemdetails.ProblemDetails) {
	routingIndicator, err := jsonparser.GetString(ausfInfo, constvalue.RoutingIndicator)
	if err != nil {
		return "", nil
	}
	return routingIndicator, nil
}

func getPcfInfoSumDnnList(pcfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	dnnList, dataType, _, err := jsonparser.Get(pcfInfoSum, constvalue.DnnList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.PcfInfoSum, constvalue.DnnList),
		}
	}

	var dnnArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(dnnList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var dnn nrfprofile.NRFKeyStruct
		dnn.SubKey1 = string(value)
		dnnArray = append(dnnArray, &dnn)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.PcfInfoSum, constvalue.DnnList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return dnnArray, nil
}

func getPcfInfoSumGroupIDList(pcfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	groupIDList, dataType, _, err := jsonparser.Get(pcfInfoSum, constvalue.GroupIDList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.PcfInfoSum, constvalue.GroupIDList),
		}
	}

	var groupIDArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(groupIDList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var groupID nrfprofile.NRFKeyStruct
		groupID.SubKey1 = string(value)
		groupIDArray = append(groupIDArray, &groupID)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.PcfInfoSum, constvalue.GroupIDList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return groupIDArray, nil
}

func getAusfInfoSumGroupIDList(ausfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	groupIDList, dataType, _, err := jsonparser.Get(ausfInfoSum, constvalue.GroupIDList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.AusfInfoSum, constvalue.GroupIDList),
		}
	}

	var groupIDArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(groupIDList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var groupID nrfprofile.NRFKeyStruct
		groupID.SubKey1 = string(value)
		groupIDArray = append(groupIDArray, &groupID)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.AusfInfoSum, constvalue.GroupIDList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return groupIDArray, nil
}

func getAusfInfoSumRoutingIndicatorList(ausfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	routingIndicatorList, dataType, _, err := jsonparser.Get(ausfInfoSum, constvalue.RoutingIndicatorList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.AusfInfoSum, constvalue.RoutingIndicatorList),
		}
	}

	var routingIndicatorArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(routingIndicatorList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var routingIndicator nrfprofile.NRFKeyStruct
		routingIndicator.SubKey1 = string(value)
		routingIndicatorArray = append(routingIndicatorArray, &routingIndicator)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.AusfInfoSum, constvalue.RoutingIndicatorList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return routingIndicatorArray, nil
}

func getUdmInfoSumGroupIDList(udmInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	groupIDList, dataType, _, err := jsonparser.Get(udmInfoSum, constvalue.GroupIDList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.UdmInfoSum, constvalue.GroupIDList),
		}
	}

	var groupIDArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(groupIDList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var groupID nrfprofile.NRFKeyStruct
		groupID.SubKey1 = string(value)
		groupIDArray = append(groupIDArray, &groupID)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.UdmInfoSum, constvalue.GroupIDList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return groupIDArray, nil
}

func getUdmInfoSumRoutingIndicatorList(udmInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	routingIndicatorList, dataType, _, err := jsonparser.Get(udmInfoSum, constvalue.RoutingIndicatorList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.UdmInfoSum, constvalue.RoutingIndicatorList),
		}
	}

	var routingIndicatorArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(routingIndicatorList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var routingIndicator nrfprofile.NRFKeyStruct
		routingIndicator.SubKey1 = string(value)
		routingIndicatorArray = append(routingIndicatorArray, &routingIndicator)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.UdmInfoSum, constvalue.RoutingIndicatorList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return routingIndicatorArray, nil
}

func getSmfInfoSumDnnList(smfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	dnnList, dataType, _, err := jsonparser.Get(smfInfoSum, constvalue.DnnList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.SmfInfoSum, constvalue.DnnList),
		}
	}

	var dnnArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(dnnList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var dnn nrfprofile.NRFKeyStruct
		dnn.SubKey1 = string(value)
		dnnArray = append(dnnArray, &dnn)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.SmfInfoSum, constvalue.DnnList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return dnnArray, nil
}

func getSmfInfoPgwFqdnList(smfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	pgwFqdnList, dataType, _, err := jsonparser.Get(smfInfoSum, constvalue.PgwFqdnList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.SmfInfoSum, constvalue.PgwFqdnList),
		}
	}

	var pgwFqdnArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(pgwFqdnList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var pgwFqdn nrfprofile.NRFKeyStruct
		pgwFqdn.SubKey1 = string(value)
		pgwFqdnArray = append(pgwFqdnArray, &pgwFqdn)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.SmfInfoSum, constvalue.PgwFqdnList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return pgwFqdnArray, nil
}

func getSmfInfoTaiList(smfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	taiList, dataType, _, err := jsonparser.Get(smfInfoSum, constvalue.TaiList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.SmfInfoSum, constvalue.TaiList),
		}
	}

	var taiArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(taiList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var tai nrfprofile.NRFKeyStruct
		mcc, _ := jsonparser.GetString(value, "plmnId", "mcc")
		mnc, _ := jsonparser.GetString(value, "plmnId", "mnc")
		if len(mnc) == 2 {
			mnc = "0" + mnc
		}
		amfID, _ := jsonparser.GetString(value, "amfId")

		tai.SubKey1 = mcc + mnc
		tai.SubKey2 = amfID

		taiArray = append(taiArray, &tai)
	})
	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.SmfInfoSum, constvalue.TaiList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return taiArray, nil
}

func getAmfInfoGuamiList(amfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	guamiList, dataType, _, err := jsonparser.Get(amfInfoSum, constvalue.GuamiList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.AmfInfoSum, constvalue.GuamiList),
		}
	}

	var guamiArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(guamiList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var guami nrfprofile.NRFKeyStruct
		mcc, _ := jsonparser.GetString(value, "plmnId", "mcc")
		mnc, _ := jsonparser.GetString(value, "plmnId", "mnc")
		if len(mnc) == 2 {
			mnc = "0" + mnc
		}
		amfID, _ := jsonparser.GetString(value, "amfId")

		guami.SubKey1 = mcc + mnc
		guami.SubKey2 = amfID

		guamiArray = append(guamiArray, &guami)
	})
	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.AmfInfoSum, constvalue.GuamiList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return guamiArray, nil
}

func getAmfInfoTaiList(amfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	taiList, dataType, _, err := jsonparser.Get(amfInfoSum, constvalue.TaiList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.AmfInfoSum, constvalue.TaiList),
		}
	}

	var taiArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(taiList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var tai nrfprofile.NRFKeyStruct
		mcc, _ := jsonparser.GetString(value, "plmnId", "mcc")
		mnc, _ := jsonparser.GetString(value, "plmnId", "mnc")
		if len(mnc) == 2 {
			mnc = "0" + mnc
		}
		tac, _ := jsonparser.GetString(value, "tac")

		tai.SubKey1 = mcc + mnc
		tai.SubKey2 = tac

		taiArray = append(taiArray, &tai)
	})
	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.AmfInfoSum, constvalue.TaiList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return taiArray, nil
}

func getAmfInfoAmfRegionIDList(amfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	amfRegionIDList, dataType, _, err := jsonparser.Get(amfInfoSum, constvalue.AmfRegionIDList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.AmfInfoSum, constvalue.AmfRegionIDList),
		}
	}

	var amfRegionIDArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(amfRegionIDList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var amfRegionID nrfprofile.NRFKeyStruct
		amfRegionID.SubKey1 = string(value)
		amfRegionIDArray = append(amfRegionIDArray, &amfRegionID)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.AmfInfoSum, constvalue.AmfRegionIDList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return amfRegionIDArray, nil
}

func getAmfInfoAmfSetIDList(amfInfoSum []byte) ([]*nrfprofile.NRFKeyStruct, *problemdetails.ProblemDetails) {
	amfSetIDList, dataType, _, err := jsonparser.Get(amfInfoSum, constvalue.AmfSetIDList)
	if dataType == jsonparser.NotExist {
		return nil, nil
	}

	if err != nil {
		return nil, &problemdetails.ProblemDetails{
			Title: fmt.Sprintf("parsing failed for %s.%s", constvalue.AmfInfoSum, constvalue.AmfSetIDList),
		}
	}

	var amfSetIDArray []*nrfprofile.NRFKeyStruct
	_, err = jsonparser.ArrayEach(amfSetIDList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var amfSetID nrfprofile.NRFKeyStruct
		amfSetID.SubKey1 = string(value)
		amfSetIDArray = append(amfSetIDArray, &amfSetID)
	})

	if err != nil {
		errorInfo := fmt.Sprintf("parsing array fail for %s.%s", constvalue.AmfInfoSum, constvalue.AmfSetIDList)
		return nil, &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	return amfSetIDArray, nil
}

func getSmfPgwFqdn(smfInfo []byte) (string, *problemdetails.ProblemDetails) {
	pgwFqdn, err := jsonparser.GetString(smfInfo, constvalue.PgwFqdn)
	if err != nil {
		return "", nil
	}
	return pgwFqdn, nil
}

func getUdmGroupId(udmInfo []byte) (string, *problemdetails.ProblemDetails) {
	groupId, err := jsonparser.GetString(udmInfo, constvalue.GroupID)
	if err != nil {
		return "", nil
	}
	return groupId, nil
}

func getUdmRoutingIndicator(udmInfo []byte) (string, *problemdetails.ProblemDetails) {
	routingIndicator, err := jsonparser.GetString(udmInfo, constvalue.RoutingIndicator)
	if err != nil {
		return "", nil
	}
	return routingIndicator, nil
}

func getUdrGroupId(udrInfo []byte) (string, *problemdetails.ProblemDetails) {
	groupId, err := jsonparser.GetString(udrInfo, constvalue.GroupID)
	if err != nil {
		return "", nil
	}
	return groupId, nil
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
	return strings.Contains(host, strconv.Itoa(innerPort))
}

// RequestFromNRFMgmt is to get sub string
func RequestFromNRFMgmt(host string, innerPort int) bool {
	return strings.Contains(host, strconv.Itoa(innerPort))
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

	helperInfo := nfProfile.CreateHelperInfo(profileUpdateTime)

	var targetNFProfile string
	if len(overrideInfo) == 0 {
		targetNFProfile = fmt.Sprintf(constvalue.TargetNFProfile, expiredTime, lastUpdateTime, profileUpdateTime, provisionFlag, md5sum, helperInfo, body, supiVersion, gpsiVersion)
	} else {
		targetNFProfile = fmt.Sprintf(constvalue.TargetNFProfileWithOverride, expiredTime, lastUpdateTime, profileUpdateTime, provisionFlag, md5sum, helperInfo, body, overrideInfo, supiVersion, gpsiVersion)
	}

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
func GetNFInfoByFilter(filter *NFProfileQueryFilter, path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("parameter path is necessary")
	}

	oql := constructNFProfileOQL(filter, path, false)

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

	count, err := strconv.Atoi(nfProfileResponse.Value[0])
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi error, %v", err)
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
					" AND value.expiredTime <= " + fmt.Sprintf("%dL", filter.ExpiredTimeEnd)
				andFlag = true
			} else {
				oql += " AND value.expiredTime >= " + fmt.Sprintf("%dL", filter.ExpiredTimeStart) +
					" AND value.expiredTime <= " + fmt.Sprintf("%dL", filter.ExpiredTimeEnd)
			}

		}

		if filter.LastUpdateTimeStart < filter.LastUpdateTimeEnd {
			if !andFlag {
				oql += " WHERE value.lastUpdateTime >= " + fmt.Sprintf("%dL", filter.LastUpdateTimeStart) +
					" AND value.lastUpdateTime <= " + fmt.Sprintf("%dL", filter.LastUpdateTimeEnd)
				andFlag = true
			} else {
				oql += " AND value.lastUpdateTime >= " + fmt.Sprintf("%dL", filter.LastUpdateTimeStart) +
					" AND value.lastUpdateTime <= " + fmt.Sprintf("%dL", filter.LastUpdateTimeEnd)
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
	case constvalue.NfTypeUDM:
		if nfProfile.UdmInfo != nil {
			return nfProfile.UdmInfo.SupiRanges
		}
	case constvalue.NfTypeUDR:
		if nfProfile.UdrInfo != nil {
			return nfProfile.UdrInfo.SupiRanges
		}
	case constvalue.NfTypeAUSF:
		if nfProfile.AusfInfo != nil {
			return nfProfile.AusfInfo.SupiRanges
		}
	case constvalue.NfTypePCF:
		if nfProfile.PcfInfo != nil {
			return nfProfile.PcfInfo.SupiRanges
		}
	default:
		return nil
	}
	return nil
}

//GpsiRangeInjection is for get GpsiRange from nfProfile
func GpsiRangeInjection(nfProfile *nrfschema.TNFProfile) []nrfschema.TIdentityRange {
	switch nfProfile.NfType {
	case constvalue.NfTypeUDM:
		if nfProfile.UdmInfo != nil {
			return nfProfile.UdmInfo.GpsiRanges
		}
	case constvalue.NfTypeUDR:
		if nfProfile.UdrInfo != nil {
			return nfProfile.UdrInfo.GpsiRanges
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

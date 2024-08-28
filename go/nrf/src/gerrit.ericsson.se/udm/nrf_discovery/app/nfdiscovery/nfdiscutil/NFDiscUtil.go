package nfdiscutil

import (
	"regexp"
	"sort"
	"strings"

	"bytes"
	"com/dbproxy"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/client"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/provprofile"
	"github.com/buger/jsonparser"
)

var (
	//Compile to compile partern into memory
	Compile map[string]*regexp.Regexp
)

//MatchResult is for match function result value
type MatchResult int32

const (
	//ResultError is for MatchResult error
	ResultError MatchResult = 0
	//ResultFoundMatch is for MatchResult found and match
	ResultFoundMatch MatchResult = 1
	//ResultFoundNotMatch is for MatchResult found and not match
	ResultFoundNotMatch MatchResult = 2
)

//PreComplieRegexp to compile pattern into memory
func PreComplieRegexp() {
	Compile = make(map[string]*regexp.Regexp)
	re, err := regexp.Compile("^[A-Fa-f0-9]*$")
	Compile[constvalue.SearchDataSupportedFeatures] = re

	re1, err1 := regexp.Compile("^((http|https)://).*$")
	Compile[constvalue.SearchDataHnrfURI] = re1

	re2, err2 := regexp.Compile("^[A-Fa-f0-9]{8}-[0-9]{3}-[0-9]{2,3}-([A-Fa-f0-9][A-Fa-f0-9]){1,10}$")
	Compile[constvalue.SearchDataExterGroupID] = re2

	re3, err3 := regexp.Compile("^(SUBSCRIPTION|POLICY|EXPOSURE|APPLICATION)$")
	Compile[constvalue.SearchDataDataSet] = re3

	re4, err4 := regexp.Compile("^[0-9]{3}$")
	Compile[constvalue.SearchDataMcc] = re4

	re5, err5 := regexp.Compile("^[0-9]{2,3}$")
	Compile[constvalue.SearchDataMnc] = re5

	re6, err6 := regexp.Compile("^[A-Fa-f0-9]{6}$")
	Compile[constvalue.SearchDataAmfID] = re6

	re7, err7 := regexp.Compile("(^[A-Fa-f0-9]{4}$)|(^[A-Fa-f0-9]{6}$)")
	Compile[constvalue.SearchDataTac] = re7

	re8, err8 := regexp.Compile("^((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$")
	Compile[constvalue.SearchDataUEIPv4Addr] = re8

	re9, err9 := regexp.Compile("^(msisdn-[0-9]{5,15}|extid-[^@]+@[^@]+|.+)$")
	Compile[constvalue.SearchDataGpsi] = re9

	re10, err10 := regexp.Compile("^(imsi-[0-9]{5,15}|nai-.+|.+)$")
	Compile[constvalue.SearchDataSupi] = re10

	re11, err11 := regexp.Compile("[0-9]{5,15}")
	Compile[constvalue.GpsiRanges] = re11

	re12, err12 := regexp.Compile("[0-9]{5,15}")
	Compile[constvalue.SupiRanges] = re12

	re13, err13 := regexp.Compile("imsi-[0-9]{5,15}|suci-[0-9]{5,15}")
	Compile[constvalue.SupiFormat] = re13

	re14, err14 := regexp.Compile("^[A-Fa-f0-9]{6}$")
	Compile[constvalue.SearchDataSnssaiSd] = re14

	re15, err15 := regexp.Compile("^[0-9]{1,4}$")
	Compile[constvalue.SearchDataRoutingIndic] = re15
	if err != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil || err8 != nil || err9 != nil || err10 != nil || err11 != nil || err12 != nil || err13 != nil || err14 != nil || err15 != nil {
		log.Debugf("err=%v, err1=%v, err2=%v, err3=%v, err4=%v, err5=%v, err6=%v, err7=%v, err8=%v, err9=%v, err10=%v, err11=%v, err12=%v, err13=%v, err14=%v, err15=%v", err, err1, err2, err3, err4, err5, err6, err7, err8, err9, err10, err11, err12, err13, err14, err15)
	}
}

//GetNFProfileMD5Sum to get md5sum from nfprofile
func GetNFProfileMD5Sum(md5sum, nfServices []byte) string {
	nfprofileSum, err := jsonparser.GetString(md5sum, "nfProfile")
	if err != nil {
		log.Errorf("No nfProfile field in md5sum")
		return ""
	}

	serviceMd5SumMap := make(map[string]string, 1)
	var keys []string
	exist := true
	_, err1 := jsonparser.ArrayEach(nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		nfServiceInstanceID, err1 := jsonparser.GetString(value, constvalue.NFServiceInstanceId)
		if err1 != nil {
			exist = false
			return
		}
		serviceMd5Sum, err := jsonparser.GetString(md5sum, nfServiceInstanceID)
		if err != nil {
			exist = false
			return
		}
		keys = append(keys, nfServiceInstanceID)
		serviceMd5SumMap[nfServiceInstanceID] = serviceMd5Sum
	})

	if exist == false || err1 != nil {
		log.Errorf("NFService MD5Sum is missed")
		return ""
	}

	sort.Strings(keys)
	for _, k := range keys {
		nfprofileSum = nfprofileSum + serviceMd5SumMap[k]
	}

	return nfprofileSum
}

func getGpsiProfile(gpsiStr string, gpsiSearchResultList *[]provprofile.GpsiSearchResult) (uint32, error) {
	request := &dbproxy.QueryRequest{}

	request.RegionName = "ericsson-nrf-gpsiprefixprofiles"
	request.Query = append(request.Query, gpsiStr)

	response, err := dbmgmt.QueryWithKey(request)
	if err != nil {
		return dbmgmt.DbInvalidData, fmt.Errorf("Get GpsiprefixProfile DB error: %v", err.Error())
	}

	if response.Code != dbmgmt.DbGetSuccess && response.Code != dbmgmt.DbDataNotExist {
		return response.Code, fmt.Errorf("Fail to get GpsiprefixProfiles, error code %d", response.Code)
	}

	if response.Code == dbmgmt.DbDataNotExist {
		return response.Code, fmt.Errorf("GpsiprefixProfile Not Found by gpsi %s", gpsiStr)
	}
	gpsiLenStr := strconv.Itoa(len(gpsiStr))

	for _, item := range response.Value {
		item = strings.Replace(item, " ", "", -1)
		value := strings.Split(item, "_")
		if 4 == len(value) && (gpsiLenStr == value[0] || "0" == value[0]) {
			var gpsiSearchResult = provprofile.GpsiSearchResult{}
			gpsiSearchResult.NfType = strings.Split(value[3], "+")
			if len(gpsiSearchResult.NfType) == 0 {
				continue
			}
			gpsiSearchResult.ValueType = value[1]
			gpsiSearchResult.ValueID = value[2]
			*gpsiSearchResultList = append(*gpsiSearchResultList, gpsiSearchResult)
		}
	}

	return dbmgmt.DbGetSuccess, nil
}

//GetGpsiGroupIDfromDB to get groupid by gpsi
func GetGpsiGroupIDfromDB(targetNfType string, gpsi string) ([]string, []string) {
	var groupID []string
	var instanceID []string
	gpsiSearchResultList := []provprofile.GpsiSearchResult{}
	if false == strings.Contains(gpsi, "msisdn-") {
		log.Debugf("getGpsiGroupIDfromDB return emptry if not contain msisdn-")
		return groupID, instanceID
	}

	//get msisdn value from ^msisdn-[0-9]{5,15}
	msisdn := string([]byte(gpsi)[7:])
	//_, err := provprofile.GetGpsiProfile(msisdn, &gpsiSearchResultList)
	_, err := getGpsiProfile(msisdn, &gpsiSearchResultList)
	if err != nil {
		log.Warnf("getGpsiGroupIDfromDB GetGpsiProfile err: %s", err.Error())
	}
	for _, item := range gpsiSearchResultList {
		nfTypeMatch := false
		for _, nfType := range item.NfType {
			if nfType == targetNfType {
				nfTypeMatch = true
				break
			}
		}
		if !nfTypeMatch {
			continue
		}
		if item.ValueType == provprofile.PrefixTypeGroupID {
			groupID = append(groupID, item.ValueID)
		} else if item.ValueType == provprofile.PrefixTypeNFInstanceID {
			instanceID = append(instanceID, item.ValueID)
		}
	}
	log.Debugf("GroupID List: %v, instanceID List: %v", groupID, instanceID)
	return groupID, instanceID
}

func getImsiProfile(imsiStr string, imsiSearchResultList *[]provprofile.ImsiSearchResult) (uint32, error) {
	request := &dbproxy.QueryRequest{}
	request.RegionName = "ericsson-nrf-imsiprefixprofiles"
	request.Query = append(request.Query, imsiStr)
	response, err := dbmgmt.QueryWithKey(request)
	if err != nil {
		return dbmgmt.DbInvalidData, fmt.Errorf("Get imsiprefixProfile DB error: %v", err.Error())
	}

	if response.Code != dbmgmt.DbGetSuccess && response.Code != dbmgmt.DbDataNotExist {
		return response.Code, fmt.Errorf("Fail to get imsiprefixProfiles, error code %d", response.Code)
	}

	if response.Code == dbmgmt.DbDataNotExist {
		return response.Code, fmt.Errorf("imsiprefixProfile Not Found by imsi %s", imsiStr)
	}
	imsiLenStr := strconv.Itoa(len(imsiStr))
	for _, item := range response.Value {
		item = strings.Replace(item, " ", "", -1)
		value := strings.Split(item, "_")
		if 4 == len(value) && (imsiLenStr == value[0] || "0" == value[0]) {
			var imsiSearchResult = provprofile.ImsiSearchResult{}
			imsiSearchResult.NfType = strings.Split(value[3], "+")
			if len(imsiSearchResult.NfType) == 0 {
				continue
			}
			imsiSearchResult.ValueType = value[1]
			imsiSearchResult.ValueID = value[2]
			*imsiSearchResultList = append(*imsiSearchResultList, imsiSearchResult)
		}
	}
	return dbmgmt.DbGetSuccess, nil
}

//GetGroupIDfromDB to get groupid by supi
func GetGroupIDfromDB(targetNfType string, supi string) ([]string, []string) {
	var groupID []string
	var instanceID []string
	imsiSearchResultList := []provprofile.ImsiSearchResult{}
	if false == strings.Contains(supi, "imsi-") {
		log.Debugf("getGroupIDfromDB return emptry if not contain imsi-")
		return groupID, instanceID
	}

	//get imsi value from ^imsi-[0-9]{5,15}
	imsi := string([]byte(supi)[5:])
	//_, err := provprofile.GetImsiProfile(imsi, &imsiSearchResultList)
	_, err := getImsiProfile(imsi, &imsiSearchResultList)
	if err != nil {
		log.Debugf("getGroupIDfromDB GetImsiProfile err: %s", err.Error())
	}
	for _, item := range imsiSearchResultList {
		nfTypeMatch := false
		for _, nfType := range item.NfType {
			if nfType == targetNfType {
				nfTypeMatch = true
				break
			}
		}
		if !nfTypeMatch {
			continue
		}
		if item.ValueType == provprofile.PrefixTypeGroupID {
			groupID = append(groupID, item.ValueID)
		} else if item.ValueType == provprofile.PrefixTypeNFInstanceID {
			instanceID = append(instanceID, item.ValueID)
		}
	}
	log.Debugf("GroupID List: %v, instanceID List: %v", groupID, instanceID)
	return groupID, instanceID
}

//IsAllowedNfType is to match requesterNfType in nfProfile/nfService field allowedNfTypes
func IsAllowedNfType(serviceOrProfile []byte, requesterNfType string, searchData string) bool {
	ok := false
	_, err := jsonparser.ArrayEach(serviceOrProfile, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if ok {
			return
		}
		nfType := string(value[:])
		if nfType == requesterNfType {
			ok = true
			return
		}

	}, searchData)
	if err != nil {
		ok = false
	}
	return ok
}

//IsAllowedPLMN is to match plmn in nfProfile/nfService field allowedPlmns
func IsAllowedPLMN(serviceOrProfile []byte, requesterPlmnList []string, searchData string) bool {
	ok := false
	_, err := jsonparser.ArrayEach(serviceOrProfile, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if ok {
			return
		}
		mcc, err2 := jsonparser.GetString(value, constvalue.Mcc)
		if err2 != nil {
			return
		}
		mnc, err2 := jsonparser.GetString(value, constvalue.Mnc)
		if err2 != nil {
			return
		}
		for _, requesterPlmn := range requesterPlmnList {
			if requesterPlmn == mcc+mnc {
				ok = true
				return
			}
		}
	}, searchData)
	if err != nil {
		ok = false
	}
	return ok
}

//IsAllowedNfFQDN is to match fqdn in nfProfile/nfService field allowedNfDomains
func IsAllowedNfFQDN(serviceOrProfile []byte, fqdn string, searchData string) bool {
	ok := false
	_, err := jsonparser.ArrayEach(serviceOrProfile, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if ok {
			return
		}
		domainPattern := string(value[:])
		domainPattern = strings.Replace(domainPattern, `\\`, `\`, -1)
		matched, err1 := regexp.MatchString(domainPattern, fqdn)
		if err1 != nil {
			log.Debugf("fqdn regex match error, err=%v", err1)
		}
		if matched {
			ok = true
			return
		}

	}, searchData)
	if err != nil {
		ok = false
	}
	log.Debugf("regexpresult = %t", ok)
	return ok

}

//IsPlmnMatchHomePlmn is to match plmnList with home plmn list
func IsPlmnMatchHomePlmn(plmnList []string) bool {
	var homePlmnList []string
	if len(plmnList) > 0 {
		for _, plmn := range cm.NfProfile.PlmnID {
			if len(plmn.Mnc) == 2 {
				homePlmnList = append(homePlmnList, plmn.Mcc+"0"+plmn.Mnc)
			} else {
				homePlmnList = append(homePlmnList, plmn.Mcc+plmn.Mnc)
			}
		}
		log.Debugf("NRF Home Plmn List: %s, target Plmn List: %s", homePlmnList, plmnList)
		for _, homePlmn := range homePlmnList {
			for _, targetPlmn := range plmnList {
				if len(targetPlmn) == 5 {
					plmnArray := []rune(targetPlmn)
					targetPlmn = string(plmnArray[0:3]) + "0" + string(plmnArray[3:])
				}
				if homePlmn == targetPlmn {
					return true
				}
			}
		}
		return false
	}

	return true
}

//DiscHTTPDo when discovery as proxy, send request(https/http) to NRF, support self redirect ,not use golang redirect
func DiscHTTPDo(method, url string, header httpclient.NHeader, body io.Reader, selfURL []string) (resp *httpclient.HttpRespData, err error) {
	var res *httpclient.HttpRespData
	var e error
	log.Debugf("Forward HTTP Request Header : %v", header)
	cm.Mutex.RLock()
	if strings.HasPrefix(url, "https") {
		res, e = client.NoRedirect_https.HttpDoProcRedirect("GET", url, header, bytes.NewBufferString(""), cm.DiscNRFSelfAPIURI)
	} else {
		res, e = client.NoRedirect_h2c.HttpDoProcRedirect("GET", url, header, bytes.NewBufferString(""), cm.DiscNRFSelfAPIURI)
	}
	cm.Mutex.RUnlock()
	return res, e
}

//GetRequestParam is to get params from request url, Example: input=/nnrf-disc/v1/nf-instances?service-names=namf-auth, output=/nf-instances?service-names=namf-auth
func GetRequestParam(request string) string {
	index := strings.Index(request, "/nf-instances")
	return request[index:]
}

//GetRequestURIRoot is to get uri root path from request url, Example: input=http://10.111.137.76:3000/nnrf-disc/v1/nf-instances?service-names=namf-auth, output=http://10.111.137.76:3000/nnrf-disc/v1/
func GetRequestURIRoot(url string) string {
	index := strings.Index(url, "nf-instances")
	return url[0:index]
}

//GetRequestURIVersion is to get version from request url, Example: input=http://10.111.137.76:3000/nnrf-disc/v1/nf-instances?service-names=namf-auth, output=v1
func GetRequestURIVersion(url string) string {
	index1 := strings.Index(url, "nnrf-disc/")
	index2 := strings.Index(url, "/nf-instances")
	if index1 == -1 {
		return ""
	}
	if index2 == -1 {
		return url[index1+10:]
	}
	return url[index1+10 : index2]
}

//FilterAddrWithVersion is to filter addr version same with request url version
func FilterAddrWithVersion(addrs []string, url string) []string {
	var supportedAddrs []string
	for _, value := range addrs {
		if GetRequestURIVersion(value) != "" && GetRequestURIVersion(url) != "" && GetRequestURIVersion(value) == GetRequestURIVersion(url) {
			supportedAddrs = append(supportedAddrs, value)
		}
	}
	return supportedAddrs
}

//StatusCodeDirectReturn is some response code direct return
func StatusCodeDirectReturn(code int) bool {
	directReturnCode := map[int]bool{
		http.StatusOK:                    true,
		http.StatusTemporaryRedirect:     true,
		http.StatusBadRequest:            true,
		http.StatusForbidden:             true,
		http.StatusNotFound:              true,
		http.StatusLengthRequired:        true,
		http.StatusRequestEntityTooLarge: true,
		http.StatusUnsupportedMediaType:  true,
	}
	return directReturnCode[code]
}

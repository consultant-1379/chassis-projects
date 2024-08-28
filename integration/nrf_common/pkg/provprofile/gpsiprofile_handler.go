package provprofile

import (
	"com/dbproxy/nfmessage/gpsiprefixprofile"
	"com/dbproxy/nfmessage/gpsiprofile"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
	"github.com/deckarep/golang-set"
	"github.com/gorilla/mux"
)

var (
	gpsiNftypesMap = map[string]bool{
		constvalue.NfTypeUDM:    true,
		constvalue.NfTypeUDR:    true,
		constvalue.NfTypeCHF:    true,
		constvalue.NfTypePCF:    true,
		constvalue.NfTypeNRFUDM: true,
		constvalue.NfTypeNRFUDR: true,
		constvalue.NfTypeNRFCHF: true,
		constvalue.NfTypeNRFPCF: true,
	}
)

//GpsiProfileHandler to process gpsiprofile request
type GpsiProfileHandler struct {
	context *ProfileContext

	originProfile []byte
	gpsiProfile   *nrfschema.GpsiProfile

	profileType   uint32
	prefixTypeStr string
	gpsiVersion   uint64
}

//Init to process GpsiProfileHandler initial
func (p *GpsiProfileHandler) Init(rw http.ResponseWriter, req *http.Request, sequenceID, profileID string, profile []byte, version uint64) {
	p.context = &ProfileContext{}
	p.context.Init(rw, req, sequenceID, profileID)

	p.gpsiProfile = &nrfschema.GpsiProfile{}
	p.originProfile = profile
	p.gpsiVersion = version
}

//SetIsRegister to flag management or provistion invoke
func (p *GpsiProfileHandler) SetIsRegister(isRegister bool) {
	p.context.IsRegister = isRegister
}

//GetProfileID get gpsiprofile ID
func (p *GpsiProfileHandler) GetProfileID() string {
	return p.gpsiProfile.GpsiProfileID
}

//GetContext to get ProfileContext
func (p *GpsiProfileHandler) GetContext() *ProfileContext {
	return p.context
}

//PostHandler to process gpsiprofile POST request
func (p *GpsiProfileHandler) PostHandler() {
	var err error
	p.context.logcontent.RequestDescription = fmt.Sprintf(`{"GpsiProfile":%s}`, string(p.originProfile))
	if err := json.Unmarshal(p.originProfile, p.gpsiProfile); err != nil {
		errorInfo := fmt.Sprintf("Umarshal GpsiProfile error. %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("Unmarshal GpsiProfile error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}
	if validate, errInfo := p.validateGpsiProfile(p.gpsiProfile); !validate {
		errorInfo := fmt.Sprintf("GpsiProfile is not validate.%s", errInfo)
		p.context.problemDetails.Title = fmt.Sprintf("GpsiProfile is not validate.%s", errInfo)
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return

	}

	//profileType, prefixTypeStr, gpsiVersion :=
	p.getProfileInfo()
	if p.profileType == profileTypeInstanceID {
		if p.gpsiVersion <= 0 {
			errorInfo := fmt.Sprintf("POST GpsiProfile from mgmt must have positive Gpsi-Version.")
			p.context.problemDetails.Title = fmt.Sprintf("POST GpsiProfile from mgmt must have positive Gpsi-Version")
			p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
			p.context.statusCode = http.StatusBadRequest
			return
		}
		p.gpsiProfile.GpsiProfileID = p.gpsiProfile.GroupID
	} else {
		p.gpsiProfile.GpsiProfileID = GenerateID(p.originProfile)
	}

	p.context.body, err = json.Marshal(p.gpsiProfile)
	if err != nil {
		errorInfo := fmt.Sprintf("Marshal GpsiProfile error. %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("Marshal GpsiProfile error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}

	gpsiprefixList := []*gpsiprefixprofile.GpsiprefixProfile{}
	err = p.getGpsiprefixProfileList(p.gpsiProfile, &gpsiprefixList, p.prefixTypeStr)
	if err != nil {
		errorInfo := fmt.Sprintf("check gpsirange error: %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("check gpsirange format error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}

	groupIDList := []string{}
	groupIDList = append(groupIDList, p.gpsiProfile.GroupID)
	index := &gpsiprofile.GpsiProfileIndex{
		NfType:      p.gpsiProfile.NfType,
		GroupIndex:  groupIDList,
		ProfileType: p.profileType,
	}
	putReq := &gpsiprofile.GpsiProfilePutRequest{
		GpsiProfileId:   p.gpsiProfile.GpsiProfileID,
		Index:           index,
		GpsiVersion:     p.gpsiVersion,
		GpsiProfileData: p.context.body,
		GpsiPrefixPut:   gpsiprefixList,
	}

	putResp, err := dbmgmt.PutGpsiProfile(putReq)
	if err != nil {
		errorInfo := fmt.Sprintf("Put GpsiProfile into DB error: %v", err)
		p.context.problemDetails.Title = "DB error"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}

	if putResp.GetCode() != dbmgmt.DbPutSuccess {
		errorInfo := fmt.Sprintf("Put GpsiProfile into DB error, error code %d", putResp.GetCode())
		p.context.problemDetails.Title = "DB error"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}

	p.context.logcontent.ResponseDescription = fmt.Sprintf(`{"GpsiProfileId":"%s"}`, p.gpsiProfile.GpsiProfileID)
	p.context.statusCode = http.StatusCreated
}

//PutHandler to process gpsiprofile PUT request
func (p *GpsiProfileHandler) PutHandler() {
	p.context.logcontent.RequestDescription = fmt.Sprintf(`{"GpsiProfile":%s}`, string(p.originProfile))
	var err error
	if err = json.Unmarshal(p.originProfile, p.gpsiProfile); err != nil {
		errorInfo := fmt.Sprintf("Unmarshal GpsiProfile error. %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("Unmarshal GpsiProfile error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}

	if validate, errInfo := p.validateGpsiProfile(p.gpsiProfile); !validate {
		errorInfo := fmt.Sprintf("GpsiProfile is not validate.%v", errInfo)
		p.context.problemDetails.Title = fmt.Sprintf("GpsiProfile is not validate.%s", errInfo)
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return

	}
	//profileType, prefixTypeStr, gpsiVersion :=
	p.getProfileInfo()
	if p.profileType == profileTypeInstanceID && p.gpsiVersion <= 0 {
		errorInfo := fmt.Sprintf("PUT GpsiProfile from mgmt must have positive Gpsi-Version.")
		p.context.problemDetails.Title = fmt.Sprintf("PUT GpsiProfile from mgmt must have not positive Gpsi-Version")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}
	var gpsiProfileID string
	if p.context.IsRegister {
		gpsiProfileID = p.context.profileID
	} else {
		gpsiProfileID = mux.Vars(p.context.req)[constvalue.GroupProfileIDName]
	}
	ok, code, errorInfo, detailInfo, gpsiProfileInfo := p.nrfProvGpsiProfileProber(gpsiProfileID)
	if !ok && !(code == http.StatusNotFound && p.profileType == profileTypeInstanceID) {
		p.context.problemDetails.Title = detailInfo
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = code
		return
	}

	origGpsiProfile := &nrfschema.GpsiProfile{}
	if len(gpsiProfileInfo) != 0 {
		//If the message is already handled by timeoutSyncGpsiProfile, need not handler it again.
		if p.gpsiVersion != 0 && p.gpsiVersion == gpsiProfileInfo[0].GetGpsiVersion() {
			p.context.statusCode = http.StatusAlreadyReported
			return
		}
		//Provision interface not allow to change gpsi profile from nfProfile
		if p.gpsiVersion == 0 && gpsiProfileInfo[0].GetGpsiVersion() != 0 {
			errorInfo := fmt.Sprintf("Gpsi profile %s is not allow to change via provision", gpsiProfileID)
			p.context.problemDetails.Title = errorInfo
			p.context.logcontent.ResponseDescription = errorInfo
			p.context.statusCode = http.StatusBadRequest
			return
		}

		resBody := gpsiProfileInfo[0].GetGpsiProfileData()
		if err = json.Unmarshal(resBody, origGpsiProfile); err != nil {
			log.Warningf("NrfProvGpsiProfilePutHandler: unmarshal orig body failure.")
		}
	}
	p.gpsiProfile.GpsiProfileID = gpsiProfileID
	p.context.body, err = json.Marshal(p.gpsiProfile)
	if err != nil {
		errorInfo := fmt.Sprintf("Marshal GpsiProfile error. %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("Marshal GpsiProfile error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}
	retCode, err := p.updateGpsiProfileDB(origGpsiProfile, p.gpsiProfile, p.profileType, p.prefixTypeStr, p.gpsiVersion, p.context.body)
	if err != nil {
		p.context.problemDetails.Title = err.Error()
		p.context.logcontent.ResponseDescription = err.Error()
		p.context.statusCode = retCode
		return
	}

	p.context.logcontent.ResponseDescription = fmt.Sprintf(`{"GpsiProfileId":"%s"}`, p.gpsiProfile.GpsiProfileID)
	p.context.statusCode = http.StatusOK
}

//DeleteHandler to process gpsiprofile DELETE request
func (p *GpsiProfileHandler) DeleteHandler() {
	var gpsiProfileID string
	if p.context.IsRegister {
		gpsiProfileID = p.context.profileID
	} else {
		gpsiProfileID = mux.Vars(p.context.req)[constvalue.GroupProfileIDName]
	}
	ok, code, errorInfo, detailInfo, gpsiProfileInfo := p.nrfProvGpsiProfileProber(gpsiProfileID)
	if !ok {
		p.context.problemDetails.Title = detailInfo
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = code
		return
	}

	//profileType, prefixTypeStr, _ :=
	p.getProfileInfo()
	resBody := gpsiProfileInfo[0].GetGpsiProfileData()
	if err := json.Unmarshal(resBody, p.gpsiProfile); err != nil {
		log.Warningf("NrfProvGpsiProfilePutHandler: unmarshal orig body failure.")
	}
	//Provision interface not allow to delete gpsi profile from nfProfile
	if p.profileType == profileTypeGroupID && gpsiProfileInfo[0].GetGpsiVersion() > 0 {
		errorInfo := fmt.Sprintf("Gpsi profile %s is not allow to delete via provision", gpsiProfileID)
		p.context.problemDetails.Title = errorInfo
		p.context.logcontent.ResponseDescription = errorInfo
		p.context.statusCode = http.StatusBadRequest
		return
	}

	gpsiprefixList := []*gpsiprefixprofile.GpsiprefixProfile{}
	err := p.getGpsiprefixProfileList(p.gpsiProfile, &gpsiprefixList, p.prefixTypeStr)
	if err != nil {
		log.Warningf("NrfProvGpsiProfilePutHandler: GetGpsiprefixProfileList error: %s", err.Error())
	}

	p.context.logcontent.RequestDescription = fmt.Sprintf(`{"gpsiProfileId":"%s"}`, gpsiProfileID)
	p.context.logcontent.ResponseDescription = ""

	gpsiProfileDelRequest := &gpsiprofile.GpsiProfileDelRequest{
		GpsiProfileId:    gpsiProfileID,
		GpsiPrefixDelete: gpsiprefixList,
	}

	deleteResp, err := dbmgmt.DeleteGpsiProfile(gpsiProfileDelRequest)
	if err != nil {

		errorInfo := fmt.Sprintf("Delete GpsiProfile DB error: %v", err)
		p.context.problemDetails.Title = "DB error"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}

	if dbmgmt.DbDeleteSuccess != deleteResp.GetCode() && dbmgmt.DbDataNotExist != deleteResp.GetCode() {

		errorInfo := fmt.Sprintf("Fail to delete GpsiProfiles, error code %d", deleteResp.Code)
		p.context.problemDetails.Title = "Fail to delete GpsiProfiles"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return

	}

	if dbmgmt.DbDataNotExist == deleteResp.GetCode() {
		errorInfo := fmt.Sprintf("GpsiProfileId %s doesn't exist.", gpsiProfileID)
		p.context.problemDetails.Title = errorInfo
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, p.context.problemDetails.Title)
		p.context.statusCode = http.StatusNotFound
		return
	}

	p.context.logcontent.ResponseDescription = fmt.Sprintf(`"successful"`)
	p.context.statusCode = http.StatusNoContent
}

//GetHandler to process gpsiprofile GET request
func (p *GpsiProfileHandler) GetHandler() {

}

func (p *GpsiProfileHandler) validateGpsiProfile(gpsiprofile *nrfschema.GpsiProfile) (bool, string) {
	if gpsiprofile.GpsiProfileID != "" {
		log.Errorf("validateGpsiProfile : receive gpsi profile can not have gpsiprofileId: %s", gpsiprofile.GroupID)
		return false, fmt.Sprint("GpsiProfileID should be null, but not")
	}

	if gpsiprofile.GroupID == "" {
		return false, fmt.Sprint("GroupID is null")
	}
	if len(gpsiprofile.NfType) == 0 {
		return false, fmt.Sprint("NfType is null")
	}

	for _, v := range gpsiprofile.NfType {
		_, ok := gpsiNftypesMap[v]
		if !ok {
			return false, fmt.Sprint("Gpsi NfType is not validate")
		}
	}
	if len(gpsiprofile.GpsiRanges) == 0 {
		return false, fmt.Sprint("GpsiRanges is null")
	}
	regStart := Compile[constvalue.Start]
	regEnd := Compile[constvalue.End]
	for _, r := range gpsiprofile.GpsiRanges {
		if r.Start != "" {
			matched := regStart.MatchString(r.Start)
			if !matched {
				return false, fmt.Sprint("gpsiRanges is not validate")
			}
		}

		if r.End != "" {
			matched := regEnd.MatchString(r.End)
			if !matched {
				return false, fmt.Sprint("gpsiRanges is not validate")
			}
		}

		if (r.Start != "" && r.End != "" && r.Pattern == "") || (r.Start == "" && r.End == "" && r.Pattern != "") {
			continue
		} else {
			return false, fmt.Sprint("gpsiRanges is not validate")
		}
	}

	return true, ""
}

func (p *GpsiProfileHandler) getProfileInfo() {
	p.profileType = profileTypeGroupID
	p.prefixTypeStr = PrefixTypeGroupID
	if p.context.IsRegister {
		p.profileType = profileTypeInstanceID
		p.prefixTypeStr = PrefixTypeNFInstanceID
	} else {
		if p.context.req.Header.Get("Gpsi-Version") == "" {
			p.gpsiVersion = 0
			return
		}

		version, err := strconv.ParseUint(p.context.req.Header.Get("Gpsi-Version"), 10, 32)
		if err != nil {
			log.Errorf("getVersionTag versionStr(%s) ParseUint error(%s)", p.context.req.Header.Get("Gpsi-Version"), err.Error())
			p.gpsiVersion = 0
			return
		}
		p.gpsiVersion = version
	}
	return
}

func (p *GpsiProfileHandler) getGpsiprefixProfileList(origProfile *nrfschema.GpsiProfile, profileList *[]*gpsiprefixprofile.GpsiprefixProfile, idType string) error {
	if nil == profileList || nil == origProfile {
		err := fmt.Errorf("Generate gpsiprefix failure")
		return err
	}
	gpsiPrefixBodySet, err := p.extractGpsiPrefixSlice(origProfile.GpsiRanges)
	if nil != err {
		return err
	}

	var nfTypeStr = ""
	for _, nfType := range origProfile.NfType {
		if nfTypeStr == "" {
			nfTypeStr = nfType
		} else {
			nfTypeStr = nfTypeStr + NfTypeSeparator + nfType
		}
	}

	gpsiPrefixBodyList := gpsiPrefixBodySet.ToSlice()
	for _, v := range gpsiPrefixBodyList {
		gpsiPrefixBody, ok := v.(PrefixBody)
		if !ok {
			err := fmt.Errorf("gpsiPrefixBodyList vaule is invalid")
			return err
		}
		length := strconv.Itoa(gpsiPrefixBody.Lenth)
		valueInfo := length + ValueInfoSeparator + idType + ValueInfoSeparator + origProfile.GroupID + ValueInfoSeparator + nfTypeStr
		gpsiprefixProfilePtr := new(gpsiprefixprofile.GpsiprefixProfile)
		var err1 error
		gpsiprefixProfilePtr.GpsiPrefix, err1 = strconv.ParseUint(gpsiPrefixBody.Prefix, 10, 64)
		if err1 != nil {
			log.Warnf("ParseUint fail: %v", err1)
		}
		gpsiprefixProfilePtr.ValueInfo = valueInfo
		*profileList = append(*profileList, gpsiprefixProfilePtr)
	}
	log.Debugf("GetGpsiprefixProfileList profileList: %v", profileList)
	return nil
}

func (p *GpsiProfileHandler) extractGpsiPrefixSlice(gpsi []nrfschema.GpsiRange) (mapset.Set, error) {
	gpsiPrefixSet := mapset.NewSet()

	for _, v := range gpsi {
		if v.Start != "" && v.End != "" && v.Pattern != "" {
			err := fmt.Errorf("GpsiRange Pattern can not coexist with Start and End,invalid Start is %s, End is %s, Pattern is %s", v.Start, v.End, v.Pattern)
			return nil, err
		}

		if v.Start != "" && v.End != "" {
			//gpsiStartEndPattern: `^\d{` + gpsiMinLenStr + `,15}$`
			re := Compile[GpsiStartEnd]
			start := re.FindString(v.Start)
			end := re.FindString(v.End)
			if start == "" || end == "" {
				err := fmt.Errorf("GpsiRange Start and End format is invalid, invalid Start is %s, End is %s. Start and End only contain Number and the lenth range from %d-15", v.Start, v.End, gpsiMinLen)
				return nil, err
			}
			startEndgpsiPrefixSet, err := processStartEnd(v.Start, v.End, gpsiMinLen)
			if err != nil {
				return nil, err
			}
			gpsiPrefixSet = gpsiPrefixSet.Union(startEndgpsiPrefixSet)
		} else if v.Pattern != "" {
			patternGpsiPrefixSet, err := processPattern(v.Pattern, Msisdn, gpsiMinLen)
			if err != nil {
				return nil, err
			}
			gpsiPrefixSet = gpsiPrefixSet.Union(patternGpsiPrefixSet)
		}
	}
	return gpsiPrefixSet, nil
}

func (p *GpsiProfileHandler) nrfProvGpsiProfileProber(gpsiProfileID string) (bool, int, string, string, []*gpsiprofile.GpsiProfileInfo) {
	var emptyInfo = []*gpsiprofile.GpsiProfileInfo{}
	id := &gpsiprofile.GpsiProfileGetRequest_GpsiProfileId{
		GpsiProfileId: gpsiProfileID,
	}
	gpsiProfileGetRequest := &gpsiprofile.GpsiProfileGetRequest{
		Data: id,
	}
	gpsiProfileResponse, err := dbmgmt.GetGpsiProfile(gpsiProfileGetRequest)
	if err != nil {
		errorInfo := fmt.Sprintf("Get GpsiProfile DB error: %v", err)
		detailInfo := "DB error"
		return false, http.StatusInternalServerError, errorInfo, detailInfo, emptyInfo
	}

	if gpsiProfileResponse.Code != dbmgmt.DbGetSuccess && gpsiProfileResponse.Code != dbmgmt.DbDataNotExist {
		errorInfo := fmt.Sprintf("Fail to get GpsiProfiles, error code %d", gpsiProfileResponse.Code)
		detailInfo := "Fail to get GpsiProfiles"
		return false, http.StatusInternalServerError, errorInfo, detailInfo, emptyInfo
	}

	if gpsiProfileResponse.Code == dbmgmt.DbDataNotExist {
		errorInfo := fmt.Sprintf("GpsiProfile Not Found by GpsiProfileId %s", gpsiProfileID)
		detailInfo := "GpsiProfile Not Found"
		return false, http.StatusNotFound, errorInfo, detailInfo, emptyInfo
	}

	if len(gpsiProfileResponse.GpsiProfileInfo) == 0 {
		errorInfo := fmt.Sprintf("requested Gpsi Profile not found")
		detailInfo := "requested Gpsi Profile not found"
		return false, http.StatusNotFound, errorInfo, detailInfo, emptyInfo
	}

	return true, http.StatusOK, "", "", gpsiProfileResponse.GpsiProfileInfo
}

func (p *GpsiProfileHandler) updateGpsiProfileDB(origGpsiProfile *nrfschema.GpsiProfile, newGpsiProfile *nrfschema.GpsiProfile, profileType uint32, prefixType string, gpsiVersion uint64, newBody []byte) (int, error) {
	gpsiprefixPutList := []*gpsiprefixprofile.GpsiprefixProfile{}
	gpsiprefixDelList := []*gpsiprefixprofile.GpsiprefixProfile{}
	err := p.getChangeGpsiprefixProfileList(origGpsiProfile, newGpsiProfile, &gpsiprefixPutList, &gpsiprefixDelList, prefixType)
	if err != nil {
		retErr := fmt.Errorf("check gpsirange error: %s", err.Error())
		return http.StatusBadRequest, retErr
	}

	groupIDList := []string{}
	groupIDList = append(groupIDList, newGpsiProfile.GroupID)
	index := &gpsiprofile.GpsiProfileIndex{
		NfType:      newGpsiProfile.NfType,
		GroupIndex:  groupIDList,
		ProfileType: profileType,
	}

	subPutReq := &gpsiprofile.GpsiProfilePutRequest{
		GpsiProfileId:    newGpsiProfile.GpsiProfileID,
		GpsiProfileData:  newBody,
		Index:            index,
		GpsiVersion:      gpsiVersion,
		GpsiPrefixDelete: gpsiprefixDelList,
		GpsiPrefixPut:    gpsiprefixPutList,
	}

	subPutResp, err := dbmgmt.PutGpsiProfile(subPutReq)
	if err != nil {
		retErr := fmt.Errorf("Replace GpsiProfile into DB error: %v", err)
		return http.StatusInternalServerError, retErr
	}

	if subPutResp.GetCode() != dbmgmt.DbPutSuccess {
		retErr := fmt.Errorf("Replace GpsiProfile into DB error, error code %d", subPutResp.GetCode())
		return http.StatusInternalServerError, retErr
	}
	return http.StatusOK, nil
}

func (p *GpsiProfileHandler) getChangeGpsiprefixProfileList(origProfile *nrfschema.GpsiProfile, newProfile *nrfschema.GpsiProfile, putList *[]*gpsiprefixprofile.GpsiprefixProfile, delList *[]*gpsiprefixprofile.GpsiprefixProfile, idType string) error {
	if nil == putList || nil == delList || nil == origProfile || nil == newProfile {
		err := fmt.Errorf("Generate gpsiprefix failure for change gpsiprofile")
		return err
	}
	if p.gpsiRangesSliceEqual(origProfile.GpsiRanges, newProfile.GpsiRanges) && origProfile.GroupID == newProfile.GroupID {
		log.Debugf("GetChangeGpsiprefixProfileList gpsi slice is equal, Orig:%v, New: %v", origProfile.GpsiRanges, newProfile.GpsiRanges)
		return nil
	}
	err := p.getGpsiprefixProfileList(origProfile, delList, idType)
	if nil != err {
		return err
	}
	err = p.getGpsiprefixProfileList(newProfile, putList, idType)
	if nil != err {
		return err
	}
	return p.removeEqualGpsiprefixProfile(putList, delList)
}

func (p *GpsiProfileHandler) gpsiRangesSliceEqual(origGpsi []nrfschema.GpsiRange, newGpsi []nrfschema.GpsiRange) bool {
	lenOrig := len(origGpsi)
	lenNew := len(newGpsi)
	if 0 == lenOrig && 0 == lenNew {
		return true
	}
	if lenOrig != lenNew {
		return false
	}
	for i, gpsi := range origGpsi {
		if gpsi.End != newGpsi[i].End || gpsi.Start != newGpsi[i].Start || gpsi.Pattern != newGpsi[i].Pattern {
			return false
		}
	}
	return true
}

func (p *GpsiProfileHandler) removeEqualGpsiprefixProfile(putList *[]*gpsiprefixprofile.GpsiprefixProfile, delList *[]*gpsiprefixprofile.GpsiprefixProfile) error {
	if nil == putList || nil == delList {
		err := fmt.Errorf("Remove gpsiprefix failure for change gpsiprofile")
		return err
	}

	var rmPutIndexList, rmDelIndexList = []int{}, []int{}
	gpsiprefixMap := make(map[uint64]PrefixMapValue)
	for i, delProfile := range *delList {
		var mapValue = PrefixMapValue{}
		mapValue.Index = i
		mapValue.ValueInfo = delProfile.ValueInfo
		gpsiprefixMap[delProfile.GpsiPrefix] = mapValue
	}

	for idx, putProfile := range *putList {
		mapVal, exist := gpsiprefixMap[putProfile.GpsiPrefix]

		if exist && mapVal.ValueInfo == putProfile.ValueInfo {
			var removeDelIndex = []int{mapVal.Index}
			rmDelIndexList = append(removeDelIndex, rmDelIndexList...)
			var removePutIndex = []int{idx}
			rmPutIndexList = append(removePutIndex, rmPutIndexList...)
		}
	}

	//remove putProfile from the back to front.
	for _, rmIndex := range rmPutIndexList {
		(*putList) = append((*putList)[:rmIndex], (*putList)[rmIndex+1:]...)
	}

	//remove delProfile from the  back to front.
	sort.Ints(rmDelIndexList[:])
	for idx := len(rmDelIndexList) - 1; idx >= 0; idx-- {
		rmIdx := rmDelIndexList[idx]
		(*delList) = append((*delList)[:rmIdx], (*delList)[rmIdx+1:]...)
	}
	return nil
}

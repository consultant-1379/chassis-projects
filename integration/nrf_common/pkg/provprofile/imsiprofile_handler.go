package provprofile

import (
	"fmt"
	"encoding/json"
	"net/http"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"strconv"
	"com/dbproxy/nfmessage/imsiprefixprofile"
	"com/dbproxy/nfmessage/groupprofile"
	"github.com/gorilla/mux"
	"github.com/deckarep/golang-set"
	"sort"
)

//ImsiProfilHandler to process imsiprofile
type ImsiProfilHandler struct {
	context *ProfileContext

	originProfile []byte
	groupProfile *nrfschema.GroupProfile

	profileType uint32
	prefixTypeStr string
	supiVersion   uint64
}

//Init to process ImsiProfileHandler initial
func (p *ImsiProfilHandler) Init(rw http.ResponseWriter, req *http.Request, sequenceID, profileID string, profile []byte, version uint64){
	p.context = &ProfileContext{}
	p.context.Init(rw, req, sequenceID, profileID)

	p.groupProfile = &nrfschema.GroupProfile{}
	p.originProfile = profile
	p.supiVersion = version
}

//SetIsRegister to flag management or provistion invoke
func (p *ImsiProfilHandler) SetIsRegister( isRegister bool) {
	p.context.IsRegister = isRegister
}

//GetProfileID get imsiprofile ID
func (p *ImsiProfilHandler) GetProfileID() string {
	return p.groupProfile.GroupProfileID
}

//GetContext to get profilecontext
func (p *ImsiProfilHandler) GetContext() *ProfileContext{
	return p.context
}
//PostHandler to process imsiprofile POST request
func (p *ImsiProfilHandler) PostHandler(){
	p.context.logcontent.RequestDescription = fmt.Sprintf(`{"GroupProfile": %s}`, string(p.originProfile))
        var err error
	if err = json.Unmarshal(p.originProfile, p.groupProfile); err != nil {
		errorInfo := fmt.Sprintf("Unmarshal GroupProfile error. %v", err)
		p.context.problemDetails.Title = "Unmarshal GroupProfile error"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}

	if validate, errInfo := p.validateGroupProfile(p.groupProfile); !validate {
		errorInfo := fmt.Sprintf("GroupProfile is not validate.%s", errInfo)
		p.context.problemDetails.Title = fmt.Sprintf("GroupProfile is not validate.%s", errInfo)
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}


	p.getProfileInfo()
	if p.profileType == profileTypeInstanceID {
		if p.supiVersion <= 0 {
			errorInfo := fmt.Sprintf("POST GroupProfile from mgmt must have positive Supi-Version.")
			p.context.problemDetails.Title = fmt.Sprintf("POST GroupProfile from mgmt must have positive Supi-Version")
			p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
			p.context.statusCode = http.StatusBadRequest
			return
		}
		p.groupProfile.GroupProfileID = p.groupProfile.GroupID
	} else {
		p.groupProfile.GroupProfileID = GenerateID(p.originProfile)
	}

	p.context.body, err = json.Marshal(p.groupProfile)
	if err != nil {
		errorInfo := fmt.Sprintf("Marshal GroupProfile error. %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("Marshal GroupProfile error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}

	imsiprefixList := []*imsiprefixprofile.ImsiprefixProfile{}
	err = p.getImsiprefixProfileList(p.groupProfile, &imsiprefixList, p.prefixTypeStr)
	if err != nil {
		errorInfo := fmt.Sprintf("check supirange error: %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("check supirange format error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}

	groupIDList := []string{}
	groupIDList = append(groupIDList, p.groupProfile.GroupID)
	index := &groupprofile.GroupProfileIndex{
		NfType:      p.groupProfile.NfType,
		GroupIndex:  groupIDList,
		ProfileType: p.profileType,
	}
	putReq := &groupprofile.GroupProfilePutRequest{
		GroupProfileId:   p.groupProfile.GroupProfileID,
		Index:            index,
		SupiVersion:      p.supiVersion,
		GroupProfileData: p.context.body,
		ImsiPrefixPut:    imsiprefixList,
	}

	putResp, err := dbmgmt.PutGroupProfile(putReq)
	if err != nil {
		errorInfo := fmt.Sprintf("Put GroupProfile into DB error: %v", err)
		p.context.problemDetails.Title = "DB error"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}

	if putResp.GetCode() != dbmgmt.DbPutSuccess {
		errorInfo := fmt.Sprintf("Put GroupProfile into DB error, error code %d", putResp.GetCode())
		p.context.problemDetails.Title = "DB error"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}

	p.context.logcontent.ResponseDescription = fmt.Sprintf(`{"GroupProfileId":"%s"}`, p.groupProfile.GroupProfileID)
	p.context.statusCode = http.StatusCreated
	return

}

//PutHandler to process imsiprofile PUT request
func (p *ImsiProfilHandler) PutHandler(){
	p.context.logcontent.RequestDescription = fmt.Sprintf(`{"GroupProfile":%s}`, string(p.originProfile))
        var err error
	if err = json.Unmarshal(p.originProfile, p.groupProfile); err != nil {
		errorInfo := fmt.Sprintf("Unmarshal GroupProfile error. %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("Unmarshal GroupProfile error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}

	if validate, errInfo := p.validateGroupProfile(p.groupProfile); !validate {
		errorInfo := fmt.Sprintf("GroupProfile is not validate.%s", errInfo)
		p.context.problemDetails.Title = fmt.Sprintf("GroupProfile is not validate.%s", errInfo)
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode =http.StatusBadRequest
		return
	}
	p.getProfileInfo()
	if p.profileType == profileTypeInstanceID && p.supiVersion <= 0 {
		errorInfo := fmt.Sprintf("PUT GroupProfile from mgmt must have positive Supi-Version.")
		p.context.problemDetails.Title = fmt.Sprintf("PUT GroupProfile from mgmt must have not positive Supi-Version")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusBadRequest
		return
	}
	var groupProfileID string
	if p.context.IsRegister {
		groupProfileID = p.context.profileID
	} else {
		groupProfileID = mux.Vars(p.context.req)[constvalue.GroupProfileIDName]
        }
	ok, code, errorInfo, detailInfo, groupProfileInfo := p.nrfProvGroupProfileProber(groupProfileID)
	if !ok && !(code == http.StatusNotFound && p.profileType == profileTypeInstanceID) {
		p.context.problemDetails.Title = detailInfo
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = code
		return
	}
	origGroupProfile := &nrfschema.GroupProfile{}
	if len(groupProfileInfo) != 0 {
		//If the message is already handled by timeoutSyncGroupProfile, need not handler it again.
		if p.supiVersion != 0 && p.supiVersion == groupProfileInfo[0].GetSupiVersion() {
			p.context.statusCode = http.StatusAlreadyReported
			return
		}
		//Provision interface not allow to change group profile from nfProfile
		if p.supiVersion == 0 && groupProfileInfo[0].GetSupiVersion() != 0 {
			errorInfo := fmt.Sprintf("Group profile %s is not allow to change via provision", groupProfileID)
			p.context.problemDetails.Title = errorInfo
			p.context.logcontent.ResponseDescription = errorInfo
			p.context.statusCode = http.StatusBadRequest
			return
		}
		resBody := groupProfileInfo[0].GetGroupProfileData()
		if err = json.Unmarshal(resBody, origGroupProfile); err != nil {
			log.Warningf("NrfProvGroupProfilePutHandler: unmarshal orig body failure.")
		}
	}
	p.groupProfile.GroupProfileID = groupProfileID
	p.context.body, err = json.Marshal(p.groupProfile)
	if err != nil {
		errorInfo := fmt.Sprintf("Marshal GroupProfile error. %v", err)
		p.context.problemDetails.Title = fmt.Sprintf("Marshal GroupProfile error")
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}
	retCode, err := p.updateGroupProfileDB(origGroupProfile, p.groupProfile, p.profileType, p.prefixTypeStr, p.supiVersion, p.context.body)
	if err != nil {
		p.context.problemDetails.Title = err.Error()
		p.context.logcontent.ResponseDescription = err.Error()
		p.context.statusCode = retCode
		return
	}

	p.context.logcontent.ResponseDescription = fmt.Sprintf(`{"GroupProfileId":"%s"}`, p.groupProfile.GroupProfileID)
	p.context.statusCode = http.StatusOK
}

//DeleteHandler to process imsiprofile DELETE request
func (p *ImsiProfilHandler) DeleteHandler(){
	var groupProfileID string
	if p.context.IsRegister{
		groupProfileID = p.context.profileID
	} else {
		groupProfileID = mux.Vars(p.context.req)[constvalue.GroupProfileIDName]
        }
	ok, code, errorInfo, detailInfo, groupProfileInfo := p.nrfProvGroupProfileProber(groupProfileID)
	if !ok {
		p.context.problemDetails.Title = detailInfo
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = code
		return
	}

	p.getProfileInfo()
	resBody := groupProfileInfo[0].GetGroupProfileData()
	if err := json.Unmarshal(resBody, p.groupProfile); err != nil {
		log.Warningf("NrfProvGroupProfilePutHandler: unmarshal orig body failure.")
	}

	//Provision interface not allow to delete group profile from nfProfile
	if p.profileType == profileTypeGroupID && groupProfileInfo[0].GetSupiVersion() > 0 {
		errorInfo := fmt.Sprintf("Group profile %s is not allow to delete via provision", groupProfileID)
		p.context.problemDetails.Title = errorInfo
		p.context.logcontent.ResponseDescription = errorInfo
		p.context.statusCode = http.StatusBadRequest
		return
	}
	imsiprefixList := []*imsiprefixprofile.ImsiprefixProfile{}
	err := p.getImsiprefixProfileList(p.groupProfile, &imsiprefixList, p.prefixTypeStr)
	if err != nil {
		log.Warningf("NrfProvGroupProfilePutHandler: GetImsiprefixProfileList error: %s", err.Error())
	}

	p.context.logcontent.RequestDescription = fmt.Sprintf(`{"groupProfileId":"%s"}`, groupProfileID)
	p.context.logcontent.ResponseDescription = ""

	groupProfileDelRequest := &groupprofile.GroupProfileDelRequest{
		GroupProfileId:   groupProfileID,
		ImsiPrefixDelete: imsiprefixList,
	}

	deleteResp, err := dbmgmt.DeleteGroupProfile(groupProfileDelRequest)
	if err != nil {

		errorInfo := fmt.Sprintf("Delete GroupProfile DB error: %v", err)
		p.context.problemDetails.Title = "DB error"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return
	}

	if dbmgmt.DbDeleteSuccess != deleteResp.GetCode() && dbmgmt.DbDataNotExist != deleteResp.GetCode() {

		errorInfo := fmt.Sprintf("Fail to delete GroupProfiles, error code %d", deleteResp.Code)
		p.context.problemDetails.Title = "Fail to delete GroupProfiles"
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, errorInfo)
		p.context.statusCode = http.StatusInternalServerError
		return

	}

	if dbmgmt.DbDataNotExist == deleteResp.GetCode() {
		errorInfo := fmt.Sprintf("GroupProfileId %s doesn't exist.", groupProfileID)
		p.context.problemDetails.Title = errorInfo
		p.context.logcontent.ResponseDescription = fmt.Sprintf(`"%s"`, p.context.problemDetails.Title)
		p.context.statusCode = http.StatusNotFound
		return
	}

	p.context.logcontent.ResponseDescription = fmt.Sprintf(`"successful"`)
	p.context.statusCode = http.StatusNoContent
}

//GetHandler to process imsiprofile GET request
func (p *ImsiProfilHandler) GetHandler(){

}

func (p *ImsiProfilHandler)getImsiprefixProfileList(origProfile *nrfschema.GroupProfile, profileList *[]*imsiprefixprofile.ImsiprefixProfile, idType string) error {
	if nil == profileList || nil == origProfile {
		err := fmt.Errorf("Generate imsiprefix failure")
		return err
	}
	imsiPrefixBodySet, err := p.extractImsiPrefixSlice(origProfile.SupiRanges)
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

	imsiPrefixBodyList := imsiPrefixBodySet.ToSlice()
	for _, v := range imsiPrefixBodyList {
		imsiPrefixBody, ok := v.(PrefixBody)
		if !ok {
			err := fmt.Errorf("imsiPrefixBodyList vaule is invalid")
			return err
		}
		length := strconv.Itoa(imsiPrefixBody.Lenth)
		valueInfo := length + ValueInfoSeparator + idType + ValueInfoSeparator + origProfile.GroupID + ValueInfoSeparator + nfTypeStr
		imsiprefixProfilePtr := new(imsiprefixprofile.ImsiprefixProfile)
		var err1 error
		imsiprefixProfilePtr.ImsiPrefix, err1 = strconv.ParseUint(imsiPrefixBody.Prefix, 10, 64)
		if err1 != nil {
			log.Warnf("ParseUint Fail: %v", err1)
		}
		imsiprefixProfilePtr.ValueInfo = valueInfo
		*profileList = append(*profileList, imsiprefixProfilePtr)
	}
	log.Debugf("GetImsiprefixProfileList profileList: %v", profileList)
	return nil
}
func (p *ImsiProfilHandler) getProfileInfo(){
        p.profileType = profileTypeGroupID
	p.prefixTypeStr = PrefixTypeGroupID
	if p.context.IsRegister {
		p.profileType = profileTypeInstanceID
		p.prefixTypeStr = PrefixTypeNFInstanceID
	} else {
		if p.context.req.Header.Get("Supi-Version") == "" {
			p.supiVersion = 0
			return
		}

		version, err := strconv.ParseUint(p.context.req.Header.Get("Supi-Version"), 10, 32)
		if err != nil {
			log.Errorf("getVersionTag versionStr(%s) ParseUint error(%s)", p.context.req.Header.Get("Supi-Version"), err.Error())
			p.supiVersion = 0
			return
		}
		p.supiVersion = version
	}
	return
}

func (p *ImsiProfilHandler)validateGroupProfile(groupprofile *nrfschema.GroupProfile) (bool, string) {
	if groupprofile.GroupProfileID != "" {
		log.Errorf("validateGroupProfile : receive group profile can not have groupprofileId: %s", groupprofile.GroupID)
		return false, fmt.Sprint("GroupProfileID should be null, but not")
	}

	if groupprofile.GroupID == "" {
		return false, fmt.Sprint("GroupID is null")
	}
	if len(groupprofile.NfType) == 0 {
		return false, fmt.Sprint("NfType is null")
	}
	supiNftypesMap := map[string]bool{
		"UDM":     true,
		"UDR":     true,
		"PCF":     true,
		"AUSF":    true,
		"CHF":     true,
		"NRFUDM":  true,
		"NRFUDR":  true,
		"NRFPCF":  true,
		"NRFAUSF": true,
		"NRFCHF":  true,
	}

	for _, v := range groupprofile.NfType {
		_, ok := supiNftypesMap[v]
		if !ok {
			return false, fmt.Sprint("Supi NfType is not validate")
		}
	}
	if len(groupprofile.SupiRanges) == 0 {
		return false, fmt.Sprint("SupiRanges is null")
	}
	regStart := Compile[constvalue.Start]
	regEnd := Compile[constvalue.End]
	for _, r := range groupprofile.SupiRanges {
		if r.Start != "" {
			matched := regStart.MatchString(r.Start)
			if !matched {
				return false, fmt.Sprint("supiRanges is not validate")
			}
		}

		if r.End != "" {
			matched := regEnd.MatchString(r.End)
			if !matched {
				return false, fmt.Sprint("supiRanges is not validate")
			}
		}

		if (r.Start != "" && r.End != "" && r.Pattern == "") || (r.Start == "" && r.End == "" && r.Pattern != "") {
			continue
		} else {
			return false, fmt.Sprint("supiRanges is not validate")
		}
	}

	return true, ""
}


func (p *ImsiProfilHandler)extractImsiPrefixSlice(supi []nrfschema.SupiRange) (mapset.Set, error) {
	imsiPrefixSet := mapset.NewSet()

	for _, v := range supi {
		if v.Start != "" && v.End != "" && v.Pattern != "" {
			err := fmt.Errorf("SupiRange Pattern can not coexist with Start and End,invalid Start is %s, End is %s, Pattern is %s", v.Start, v.End, v.Pattern)
			return nil, err
		}

		if v.Start != "" && v.End != "" {
			//imsiStartEndPattern: `^\d{` + imsiMinLenStr + `,15}$`
			re := Compile[ImsiStartEnd]
			start := re.FindString(v.Start)
			end := re.FindString(v.End)
			if start == "" || end == "" {
				err := fmt.Errorf("SupiRange Start and End format is invalid, invalid Start is %s, End is %s. Start and End only contain Number and the lenth range from %d-15", v.Start, v.End, imsiMinLen)
				return nil, err
			}
			startEndimsiPrefixSet, err := processStartEnd(v.Start, v.End, imsiMinLen)
			if err != nil {
				return nil, err
			}
			imsiPrefixSet = imsiPrefixSet.Union(startEndimsiPrefixSet)
		} else if v.Pattern != "" {
			patternImsiPrefixSet, err := processPattern(v.Pattern, Imsi, imsiMinLen)
			if err != nil {
				return nil, err
			}
			imsiPrefixSet = imsiPrefixSet.Union(patternImsiPrefixSet)
		}
	}
	return imsiPrefixSet, nil
}



func (p *ImsiProfilHandler)nrfProvGroupProfileProber(groupProfileID string) (bool, int, string, string, []*groupprofile.GroupProfileInfo) {
	var emptyInfo = []*groupprofile.GroupProfileInfo{}
	id := &groupprofile.GroupProfileGetRequest_GroupProfileId{
		GroupProfileId: groupProfileID,
	}
	groupProfileGetRequest := &groupprofile.GroupProfileGetRequest{
		Data: id,
	}
	groupProfileResponse, err := dbmgmt.GetGroupProfile(groupProfileGetRequest)
	if err != nil {
		errorInfo := fmt.Sprintf("Get GroupProfile DB error: %v", err)
		detailInfo := "DB error"
		return false, http.StatusInternalServerError, errorInfo, detailInfo, emptyInfo
	}

	if groupProfileResponse.Code != dbmgmt.DbGetSuccess && groupProfileResponse.Code != dbmgmt.DbDataNotExist {
		errorInfo := fmt.Sprintf("Fail to get GroupProfiles, error code %d", groupProfileResponse.Code)
		detailInfo := "Fail to get GroupProfiles"
		return false, http.StatusInternalServerError, errorInfo, detailInfo, emptyInfo
	}

	if groupProfileResponse.Code == dbmgmt.DbDataNotExist {
		errorInfo := fmt.Sprintf("GroupProfile Not Found by GroupProfileId %s", groupProfileID)
		detailInfo := "GroupProfile Not Found"
		return false, http.StatusNotFound, errorInfo, detailInfo, emptyInfo
	}

	if len(groupProfileResponse.GroupProfileInfo) == 0 {
		errorInfo := fmt.Sprintf("requested Group Profile not found")
		detailInfo := "requested Group Profile not found"
		return false, http.StatusNotFound, errorInfo, detailInfo, emptyInfo
	}

	return true, http.StatusOK, "", "", groupProfileResponse.GroupProfileInfo
}


func (p *ImsiProfilHandler)updateGroupProfileDB(origGroupProfile *nrfschema.GroupProfile, newGroupProfile *nrfschema.GroupProfile, profileType uint32, prefixType string, supiVersion uint64, newBody []byte) (int, error) {
	imsiprefixPutList := []*imsiprefixprofile.ImsiprefixProfile{}
	imsiprefixDelList := []*imsiprefixprofile.ImsiprefixProfile{}
	err := p.GetChangeImsiprefixProfileList(origGroupProfile, newGroupProfile, &imsiprefixPutList, &imsiprefixDelList, prefixType)
	if err != nil {
		retErr := fmt.Errorf("check supirange error: %s", err.Error())
		return http.StatusBadRequest, retErr
	}

	groupIDList := []string{}
	groupIDList = append(groupIDList, newGroupProfile.GroupID)
	index := &groupprofile.GroupProfileIndex{
		NfType:      newGroupProfile.NfType,
		GroupIndex:  groupIDList,
		ProfileType: profileType,
	}

	subPutReq := &groupprofile.GroupProfilePutRequest{
		GroupProfileId:   newGroupProfile.GroupProfileID,
		GroupProfileData: newBody,
		Index:            index,
		SupiVersion:      supiVersion,
		ImsiPrefixDelete: imsiprefixDelList,
		ImsiPrefixPut:    imsiprefixPutList,
	}

	subPutResp, err := dbmgmt.PutGroupProfile(subPutReq)
	if err != nil {
		retErr := fmt.Errorf("Replace GroupProfile into DB error: %v", err)
		return http.StatusInternalServerError, retErr
	}

	if subPutResp.GetCode() != dbmgmt.DbPutSuccess {
		retErr := fmt.Errorf("Replace GroupProfile into DB error, error code %d", subPutResp.GetCode())
		return http.StatusInternalServerError, retErr
	}
	return http.StatusOK, nil
}


//GetChangeImsiprefixProfileList get imsi prefix change
func (p *ImsiProfilHandler)GetChangeImsiprefixProfileList(origProfile *nrfschema.GroupProfile, newProfile *nrfschema.GroupProfile, putList *[]*imsiprefixprofile.ImsiprefixProfile, delList *[]*imsiprefixprofile.ImsiprefixProfile, idType string) error {
	if nil == putList || nil == delList || nil == origProfile || nil == newProfile {
		err := fmt.Errorf("Generate imsiprefix failure for change groupprofile")
		return err
	}
	if p.supiRangesSliceEqual(origProfile.SupiRanges, newProfile.SupiRanges) && origProfile.GroupID == newProfile.GroupID {
		log.Debugf("GetChangeImsiprefixProfileList supi slice is equal, Orig:%v, New: %v", origProfile.SupiRanges, newProfile.SupiRanges)
		return nil
	}
	err := p.getImsiprefixProfileList(origProfile, delList, idType)
	if nil != err {
		return err
	}
	err = p.getImsiprefixProfileList(newProfile, putList, idType)
	if nil != err {
		return err
	}
	return p.removeEqualImsiprefixProfile(putList, delList)
}

func (p *ImsiProfilHandler)removeEqualImsiprefixProfile(putList *[]*imsiprefixprofile.ImsiprefixProfile, delList *[]*imsiprefixprofile.ImsiprefixProfile) error {
	if nil == putList || nil == delList {
		err := fmt.Errorf("Remove imsiprefix failure for change groupprofile")
		return err
	}

	var rmPutIndexList = []int{}
	var rmDelIndexList = []int{}
	imsiprefixMap := make(map[uint64]PrefixMapValue)
	for i, delProfile := range *delList {
		var mapValue = PrefixMapValue{}
		mapValue.Index = i
		mapValue.ValueInfo = delProfile.ValueInfo
		imsiprefixMap[delProfile.ImsiPrefix] = mapValue
	}

	for idx, putProfile := range *putList {
		mapVal, exist := imsiprefixMap[putProfile.ImsiPrefix]

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
func (p *ImsiProfilHandler)supiRangesSliceEqual(origSupi []nrfschema.SupiRange, newSupi []nrfschema.SupiRange) bool {
	lenOrig := len(origSupi)
	lenNew := len(newSupi)
	if 0 == lenOrig && 0 == lenNew {
		return true
	}
	if lenOrig != lenNew {
		return false
	}
	for i, supi := range origSupi {
		if supi.End != newSupi[i].End || supi.Start != newSupi[i].Start || supi.Pattern != newSupi[i].Pattern {
			return false
		}
	}
	return true
}



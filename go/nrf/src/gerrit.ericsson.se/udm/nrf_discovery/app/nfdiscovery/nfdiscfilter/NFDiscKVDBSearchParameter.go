package nfdiscfilter

import (
	"strings"
	"github.com/buger/jsonparser"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"com/dbproxy/nfmessage/nfprofile"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

func buildStringSearchParameter(mapping configmap.SearchMapping, value string, operation int32) *MetaExpression {
	searchExpression := &SearchExpression{
		from: mapping.From,
		where: mapping.Where,
		value: value,
		operation: operation,
		valueType: constvalue.ValueString,
	}
	metaExpression := &MetaExpression{
		commonExpression:searchExpression,
		expressionType: constvalue.TypeSearchExpression,
	}
	return metaExpression
}

func buildIntegerSearchParameter(mapping configmap.SearchMapping, value uint64, operation int32) *MetaExpression {
	searchExpression := &SearchExpression{
		from: mapping.From,
		where: mapping.Where,
		value: value,
		operation: operation,
		valueType: constvalue.ValueNum,
	}
	metaExpression := &MetaExpression{
		commonExpression:searchExpression,
		expressionType: constvalue.TypeSearchExpression,
	}
	return metaExpression
}

func buildAndExpression(metaExpressionList []*MetaExpression) *MetaExpression {
	if len(metaExpressionList) == 0 {
		return nil
	}
	andExpression := &AndExpression{
		andMetaExpression: metaExpressionList,
	}
	metaAndExpression := &MetaExpression{
		commonExpression: andExpression,
		expressionType: constvalue.TypeAndExpression,
	}
	return metaAndExpression
}
func buildORExpression(metaExpressionList []*MetaExpression) *MetaExpression {
	if len(metaExpressionList) == 0 {
		return nil
	}
	orExpression := &ORExpression{
		orMetaExpression: metaExpressionList,
	}
	metaORExpression := &MetaExpression{
		commonExpression: orExpression,
		expressionType: constvalue.TypeORExpression,
	}
	return metaORExpression
}

func createTaiExpressionForAbsence(nfType, parameter string) *MetaExpression {
	tacMapping := getParamSearchPath(nfType, parameter + constvalue.PathList + constvalue.PathTac)
	tacExpression := buildStringSearchParameter(tacMapping, constvalue.EmptyTac, constvalue.EQ)

	patternMapping := getParamSearchPath(nfType, parameter + constvalue.PathRangeList + constvalue.PathTac + constvalue.PathPattern)
	tacRangeExpression := buildStringSearchParameter(patternMapping, constvalue.EmptyTacRangePattern, constvalue.EQ)

	var taiExpressionList []*MetaExpression
	taiExpressionList = append(taiExpressionList, tacExpression, tacRangeExpression)
	return buildAndExpression(taiExpressionList)
}

func createTaiExpression(nfType string, parameter string, plmnid string, tac string) *MetaExpression {

	var metaExpressionListTaiList []*MetaExpression
	var metaExpressionListTai []*MetaExpression
	var metaExpressionListTaiRange []*MetaExpression

	taiListParameter := parameter + constvalue.PathList
	taiPlmnExpression := createPlmnFilter(nfType, taiListParameter, plmnid)
	metaExpressionListTaiList = append(metaExpressionListTaiList, taiPlmnExpression)

	tacMapping := getParamSearchPath(nfType, parameter + constvalue.PathList + constvalue.PathTac)
	tacExpression := buildStringSearchParameter(tacMapping, tac, constvalue.EQ)
	metaExpressionListTaiList = append(metaExpressionListTaiList, tacExpression)
	metaExpressionListTai = append(metaExpressionListTai, buildAndExpression(metaExpressionListTaiList))

	taiRangeListParameter := parameter + constvalue.PathRangeList
	taiPlmnExpression = createPlmnFilter(nfType, taiRangeListParameter, plmnid)
	metaExpressionListTaiRange = append(metaExpressionListTaiRange, taiPlmnExpression)

	tacParameter := parameter + constvalue.PathRangeList + constvalue.PathTac
	tacRangeExpression := createRangeExpression(nfType, tacParameter, tac, tac)
	metaExpressionListTaiRange = append(metaExpressionListTaiRange, tacRangeExpression)

	metaExpressionListTai = append(metaExpressionListTai, buildAndExpression(metaExpressionListTaiRange))

	return buildORExpression(metaExpressionListTai)
}

func createSupiExpressionForAbsence(nfType, parameter string) *MetaExpression {
	patternMapping := getParamSearchPath(nfType, parameter+constvalue.PathAbsencePattern)
	return buildStringSearchParameter(patternMapping, constvalue.MatchAll, constvalue.EQ)
}

func  createGpsiExpressionForAbsence(nfType, parameter string) *MetaExpression {
	mapping := getParamSearchPath(nfType, parameter + constvalue.PathAbsencePattern)
	return buildStringSearchParameter(mapping, constvalue.MatchAll, constvalue.EQ)
}

func createSupiGroupIDExpressionForAbsence(nfType, parameter string) *MetaExpression{
	var metaExpressionList []*MetaExpression
	groupIDMapping := getParamSearchPath(nfType, constvalue.SearchDataGroupIDList)
	if groupIDMapping.Parameter != "" {
		metaExpressionList = append(metaExpressionList, buildStringSearchParameter(groupIDMapping, constvalue.EmptyGroupID, constvalue.EQ))
	}

	patternMapping := getParamSearchPath(nfType, parameter+constvalue.PathAbsencePattern)
	metaExpressionList = append(metaExpressionList, buildStringSearchParameter(patternMapping, constvalue.MatchAll, constvalue.EQ))

	return buildAndExpression(metaExpressionList)
}

func createGpsiGroupIDExpressionForAbsence(nfType , parameter string) *MetaExpression{
	var metaExpressionList []*MetaExpression
	groupIDMapping := getParamSearchPath(nfType, constvalue.SearchDataGroupIDList)
	if groupIDMapping.Parameter != "" {
		metaExpressionList = append(metaExpressionList, buildStringSearchParameter(groupIDMapping, constvalue.EmptyGroupID, constvalue.EQ))
	}


	mapping := getParamSearchPath(nfType, parameter+constvalue.PathAbsencePattern)
	metaExpressionList = append(metaExpressionList, buildStringSearchParameter(mapping, constvalue.MatchAll, constvalue.EQ))

	return buildAndExpression(metaExpressionList)
}

func createPlmnExpressionForAbsence(nfType , parameter string) *MetaExpression{
	mapping := getParamSearchPath(nfType, parameter + constvalue.PathPattern)
	return buildStringSearchParameter(mapping, constvalue.EmptyPlmnRangePattern, constvalue.EQ)
}

func createRangeAbsenceExpression(nfType, parameter string) *MetaExpression{
	mapping := getParamSearchPath(nfType, parameter + constvalue.PathPattern)
	return buildStringSearchParameter(mapping, constvalue.EmptyExternalIDPattern, constvalue.EQ)
}

func createGroupIDExpression(nfType, parameter string, groupID []string ) *MetaExpression {
        var metaExpressionList []*MetaExpression

	groupIDMapping := getParamSearchPath(nfType, parameter)
	if groupIDMapping.Parameter != "" {
		for _, v := range groupID {
			groudIDExpression := buildStringSearchParameter(groupIDMapping, v, constvalue.EQ)
			metaExpressionList = append(metaExpressionList, groudIDExpression)
		}
	}

	return buildORExpression(metaExpressionList)
}

func createInstanceExpression(instances []string) *MetaExpression {
	mapping := getParamSearchPath(constvalue.Common, constvalue.SearchDataTargetInstID)
	var metaExpressionList []*MetaExpression
	for _, v := range instances {
		instExpression := buildStringSearchParameter(mapping, v, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, instExpression)
	}

	return buildORExpression(metaExpressionList)
}

func createGroupIDInstanceIDExpression(value, nfType, filterType string) *MetaExpression {
	var groupID []string
	var instanceID []string
	if filterType == constvalue.SearchDataSupi {
		groupID, instanceID = nfdiscutil.GetGroupIDfromDB(nfType, value)
	} else if filterType == constvalue.SearchDataGpsi {
		groupID, instanceID = nfdiscutil.GetGpsiGroupIDfromDB(nfType, value)
	} else {
		return nil
	}

	var groupIDExpressionList []*MetaExpression
	if filterType == constvalue.SearchDataSupi {
		if nfType != constvalue.NfTypeCHF {
			groupIDExpressionList = append(groupIDExpressionList, createSupiGroupIDExpressionForAbsence(nfType, filterType))
		} else {
			groupIDExpressionList = append(groupIDExpressionList, createSupiExpressionForAbsence(nfType, filterType))
		}
	} else {
		if nfType != constvalue.NfTypeCHF {
			groupIDExpressionList = append(groupIDExpressionList, createGpsiGroupIDExpressionForAbsence(nfType, filterType))
		} else {
			groupIDExpressionList = append(groupIDExpressionList, createGpsiExpressionForAbsence(nfType, filterType))
		}
	}

	if len(groupID) > 0 && len(instanceID) > 0 {
		groupIDExpressionList = append(groupIDExpressionList, createGroupIDExpression(nfType, constvalue.SearchDataGroupIDList, groupID))
		groupIDExpressionList = append(groupIDExpressionList, createInstanceExpression(instanceID))
		return buildORExpression(groupIDExpressionList)
	} else if len(groupID) > 0 {
		groupIDExpressionList = append(groupIDExpressionList, createGroupIDExpression(nfType, constvalue.SearchDataGroupIDList, groupID))
		return buildORExpression(groupIDExpressionList)
	} else if len(instanceID) > 0 {
		groupIDExpressionList = append(groupIDExpressionList, createInstanceExpression(instanceID))
		return buildORExpression(groupIDExpressionList)
	}

	return nil
}

func createSupiExpression(nfType, parameter, supi string) *MetaExpression {
	imsi := string([]byte(supi)[5:])
	return createRangeExpression(nfType, parameter, imsi, supi)
}


func createGpsiExpression(nfType, parameter, gpsi string) *MetaExpression {

	msisdn := string([]byte(gpsi)[7:])
	return createRangeExpression(nfType, parameter, msisdn, gpsi)
}


func createRangeExpression(nfType, parameter, startendValue string, patternValue string) *MetaExpression {

	startendExpression := createStartEndExpression(nfType, parameter, startendValue)
	patternMapping := getParamSearchPath(nfType, parameter+constvalue.PathPattern)
	patternExpression := buildStringSearchParameter(patternMapping, patternValue, constvalue.REGEX)

	var metaExpressionList []*MetaExpression
	metaExpressionList = append(metaExpressionList, startendExpression)
	metaExpressionList = append(metaExpressionList, patternExpression)

	return buildORExpression(metaExpressionList)
}

func createStartEndExpression(nfType, parameter, value string) *MetaExpression {

	var metaStartExpressionList []*MetaExpression
	startExpression := buildStringSearchParameter(getParamSearchPath(nfType, parameter+constvalue.PathStart), value, constvalue.LE)
	startLenExpression := buildIntegerSearchParameter(getParamSearchPath(nfType, parameter+constvalue.PathStartLength), uint64(len(value)), constvalue.LE)
	metaStartExpressionList = append(metaStartExpressionList, startExpression)
	metaStartExpressionList = append(metaStartExpressionList, startLenExpression)
	startExpressionTotal := buildAndExpression(metaStartExpressionList)

	var metaEndExpressionList []*MetaExpression
	endExpression := buildStringSearchParameter(getParamSearchPath(nfType, parameter+constvalue.PathEnd), value, constvalue.GE)
	endLenExpression := buildIntegerSearchParameter(getParamSearchPath(nfType, parameter+constvalue.PathEndLength), uint64(len(value)), constvalue.GE)
	metaEndExpressionList = append(metaEndExpressionList, endExpression)
	metaEndExpressionList = append(metaEndExpressionList, endLenExpression)
	endExpressionTotal := buildAndExpression(metaEndExpressionList)

	var metaExpressionList []*MetaExpression
	metaExpressionList = append(metaExpressionList, startExpressionTotal)
	metaExpressionList = append(metaExpressionList, endExpressionTotal)

	return buildAndExpression(metaExpressionList)
}

func createPlmnFilter(nfType, parameter, plmn string) *MetaExpression {

	plmnArray := []rune(plmn)
	mcc := string(plmnArray[0:3])
	mnc := string(plmnArray[3:])

	mccMapping := getParamSearchPath(nfType, parameter + constvalue.PathPlmnMcc)
	mccExpression := buildStringSearchParameter(mccMapping, mcc, constvalue.EQ)
	mncMapping := getParamSearchPath(nfType, parameter + constvalue.PathPlmnMnc)
	mncExpression := buildStringSearchParameter(mncMapping, mnc, constvalue.EQ)

	var metaExpressionList []*MetaExpression
	metaExpressionList = append(metaExpressionList, mccExpression, mncExpression)

	return buildAndExpression(metaExpressionList)
}

func createSnssaiFilter(queryForm *nfdiscrequest.DiscGetPara, nfType, parameter string) *MetaExpression {
	sdMapping := getParamSearchPath(nfType, parameter + constvalue.PathSd)
	sstMapping := getParamSearchPath(nfType, parameter + constvalue.PathSst)
	var metaExpressionList []*MetaExpression
	if queryForm.GetExistFlag(constvalue.SearchDataSnssais) && len(queryForm.GetValue()[constvalue.SearchDataSnssais]) > 0 {
		for _, item := range queryForm.GetValue()[constvalue.SearchDataSnssais] {
			if !strings.Contains(item, "[") && !strings.Contains(item, "]") {
				item = "[" + item + "]"
			}
			_, err := jsonparser.ArrayEach([]byte(item), func(value []byte, dataType jsonparser.ValueType, offset int, err2 error) {
				sstID, parseErr := jsonparser.GetInt(value, constvalue.SearchDataSnssaiSst)
				sstExpression := buildIntegerSearchParameter(sstMapping, uint64(sstID), constvalue.EQ)
				var subMetaExpressionList []*MetaExpression
				subMetaExpressionList = append(subMetaExpressionList, sstExpression)
				sdID, parseErr1 := jsonparser.GetString(value, constvalue.SearchDataSnssaiSd)
				if parseErr != nil || parseErr1 != nil {
					log.Debugf("sst or sd parse error, err=%v, err1=%v", parseErr, parseErr1)
				}
				if sdID != "" {
					var sdExpressionList []*MetaExpression
					sdID = strings.ToLower(sdID)
					sdExpression := buildStringSearchParameter(sdMapping, sdID, constvalue.EQ)
					sdExpressionAbsence := buildStringSearchParameter(sdMapping, constvalue.EmptySd, constvalue.EQ)
					sdExpressionList = append(sdExpressionList, sdExpression, sdExpressionAbsence)
					sdExpressionTotal := buildORExpression(sdExpressionList)
					subMetaExpressionList = append(subMetaExpressionList, sdExpressionTotal)
				}
				andExpression := buildAndExpression(subMetaExpressionList)

				metaExpressionList = append(metaExpressionList, andExpression)

			})

			if err != nil {
				return nil
			}
		}
	}

	if len(metaExpressionList) == 0 {
		return nil
	}
	return buildORExpression(metaExpressionList)
}

func setCustomInfoFilter(queryForm *nfdiscrequest.DiscGetPara, nfProfileFilter *nfprofile.NFProfileFilter) {
	/*nfProfileFilter.ExpiredTimeRange = &nfprofile.Range{
		Start: uint64(time.Now().Unix()) * 1000,
		End:   math.MaxInt64,
	}*/
	nfProfileFilter.Provisioned = 0
}

func createGuamiFilter(nfType, parameter, plmnid, amfid string ) *MetaExpression {
	var metaExpressionList []*MetaExpression
	var metaGuamiList []*MetaExpression
	var metaBackupFailureList []*MetaExpression
	var metaBackupRemovalList []*MetaExpression
	//to filter amfInfo.guamiList
	guamiPlmnExpression := createPlmnFilter(nfType, parameter, plmnid)
	metaGuamiList = append(metaGuamiList, guamiPlmnExpression)

	amfIDMapping := getParamSearchPath(nfType, parameter + constvalue.PathAmfID)
	amfIDExpression := buildStringSearchParameter(amfIDMapping, amfid, constvalue.EQ)
	metaGuamiList = append(metaGuamiList, amfIDExpression)

	metaExpressionList = append(metaExpressionList, buildAndExpression(metaGuamiList))

	//to filter backupInfoAmfFailure
	backupFailureParam := parameter + constvalue.PathBackfailure
	guamiPlmnExpression = createPlmnFilter(nfType, backupFailureParam, plmnid)
	metaBackupFailureList = append(metaBackupFailureList, guamiPlmnExpression)

	backupFailureAmfIDMapping := getParamSearchPath(nfType, parameter + constvalue.PathBackfailure + constvalue.PathAmfID)
	amfIDExpression = buildStringSearchParameter(backupFailureAmfIDMapping, amfid, constvalue.EQ)
	metaBackupFailureList = append(metaBackupFailureList, amfIDExpression)

	metaExpressionList = append(metaExpressionList, buildAndExpression(metaBackupFailureList))

	//to filter backpInfoAmfRemoval
	backupRemovalParam := parameter + constvalue.PathBackremoval
	guamiPlmnExpression = createPlmnFilter(nfType, backupRemovalParam, plmnid)
	metaBackupRemovalList = append(metaBackupRemovalList, guamiPlmnExpression)

	backupRemovalAmfIDMapping := getParamSearchPath(nfType, parameter + constvalue.PathBackremoval + constvalue.PathAmfID)
	amfIDExpression = buildStringSearchParameter(backupRemovalAmfIDMapping, amfid, constvalue.EQ)
	metaBackupRemovalList = append(metaBackupRemovalList, amfIDExpression)

	metaExpressionList = append(metaExpressionList, buildAndExpression(metaBackupRemovalList))

	return buildORExpression(metaExpressionList)
}

func createAbsenceInfoSnssaiExpression(nfType, parameter string) *MetaExpression {
	//var metaExpressionList []*MetaExpression

	sdMapping := getParamSearchPath(nfType, parameter + constvalue.PathSd)
	sstMapping := getParamSearchPath(nfType, parameter + constvalue.PathSst)
	var subMetaEmptyExpressionList []*MetaExpression
	sstExpression := buildIntegerSearchParameter(sstMapping, constvalue.EmptySst, constvalue.EQ)
	sdExpression := buildStringSearchParameter(sdMapping, constvalue.EmptySd, constvalue.EQ)
	subMetaEmptyExpressionList = append(subMetaEmptyExpressionList, sstExpression, sdExpression)
	emptyAndExpression := buildAndExpression(subMetaEmptyExpressionList)

	//andExpression := createSnssaiFilter(queryForm, "body.sNssais")
	//metaExpressionList = append(metaExpressionList, emptyAndExpression, andExpression)

	return emptyAndExpression
}

//createDnaiListFilter is to generate dnai expression filter
func createDnaiListFilter(dnaiList []string, mapping configmap.SearchMapping) *MetaExpression {
	var metaExpressionList []*MetaExpression
	for _, dnai := range (dnaiList) {
		dnaiExpression := buildStringSearchParameter(mapping, dnai, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, dnaiExpression)
	}

	return buildORExpression(metaExpressionList)
}

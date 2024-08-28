package provprofile

import (
	"com/dbproxy/nfmessage/gpsiprefixprofile"
	"fmt"
	"strconv"
	"strings"

	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
)

//GpsiSearchResult gpsi search result
type GpsiSearchResult struct {
	ValueType string
	ValueID   string
	NfType    []string
}

// GetGpsiProfile is for get groupid or nf-instance ID via gpsi
func GetGpsiProfile(gpsiStr string, gpsiSearchResultList *[]GpsiSearchResult) (uint32, error) {
	if nil == gpsiSearchResultList {
		err := fmt.Errorf("getGpsiProfile gpsiprefixValue ptr is nil")
		return dbmgmt.DbInvalidData, err
	}
	gpsi, err := strconv.ParseUint(gpsiStr, 10, 64)
	if err != nil {
		errParse := fmt.Errorf("Gpsi string parse to long error: %v", err)
		return dbmgmt.DbInvalidData, errParse
	}

	gpsiprefixProfileGetRequest := &gpsiprefixprofile.GpsiprefixProfileGetRequest{
		SearchGpsi: gpsi,
	}

	gpsiprefixProfileResponse, err := dbmgmt.GetGpsiprefixProfile(gpsiprefixProfileGetRequest)
	if err != nil {

		errDB := fmt.Errorf("Get GpsiprefixProfile DB error: %v", err)
		return dbmgmt.DbInvalidData, errDB
	}

	if gpsiprefixProfileResponse.Code != dbmgmt.DbGetSuccess && gpsiprefixProfileResponse.Code != dbmgmt.DbDataNotExist {

		err = fmt.Errorf("Fail to get GpsiprefixProfiles, error code %d", gpsiprefixProfileResponse.Code)
		return gpsiprefixProfileResponse.Code, err

	}

	if gpsiprefixProfileResponse.Code == dbmgmt.DbDataNotExist {

		err = fmt.Errorf("GpsiprefixProfile Not Found by gpsi %s", gpsiStr)
		return gpsiprefixProfileResponse.Code, err
	}

	gpsiLenStr := strconv.Itoa(len(gpsiStr))

	for _, item := range gpsiprefixProfileResponse.ValueInfo {
		value := strings.Split(item, ValueInfoSeparator)
		if 4 == len(value) && (gpsiLenStr == value[0] || "0" == value[0]) {
			var gpsiSearchResult = GpsiSearchResult{}
			gpsiSearchResult.NfType = strings.Split(value[3], NfTypeSeparator)
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

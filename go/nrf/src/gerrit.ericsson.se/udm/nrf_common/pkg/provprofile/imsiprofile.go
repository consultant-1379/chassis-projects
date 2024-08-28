package provprofile

import (
	"com/dbproxy/nfmessage/imsiprefixprofile"
	"fmt"
	"strconv"
	"strings"

	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
)

//ImsiSearchResult imsi search result
type ImsiSearchResult struct {
	ValueType string
	ValueID   string
	NfType    []string
}

const (
	//PrefixTypeGroupID is for type groupID
	PrefixTypeGroupID = "gid"
	//PrefixTypeNFInstanceID is for type NFInstanceID
	PrefixTypeNFInstanceID = "nf"
	//ValueInfoSeparator is for ValueInfo separate symbol
	ValueInfoSeparator = "_"
	//NfTypeSeparator is for []nfType separate symbol
	NfTypeSeparator = "+"
)

// GetImsiProfile is for get groupid or nf-instance ID via imsi
func GetImsiProfile(imsiStr string, imsiSearchResultList *[]ImsiSearchResult) (uint32, error) {
	if nil == imsiSearchResultList {
		err := fmt.Errorf("getImsiProfile imsiprefixValue ptr is nil")
		return dbmgmt.DbInvalidData, err
	}
	imsi, err := strconv.ParseUint(imsiStr, 10, 64)
	if err != nil {
		errParse := fmt.Errorf("Imsi string parse to long error: %v", err)
		return dbmgmt.DbInvalidData, errParse
	}

	imsiprefixProfileGetRequest := &imsiprefixprofile.ImsiprefixProfileGetRequest{
		SearchImsi: imsi,
	}

	imsiprefixProfileResponse, err := dbmgmt.GetImsiprefixProfile(imsiprefixProfileGetRequest)
	if err != nil {

		errDB := fmt.Errorf("Get ImsiprefixProfile DB error: %v", err)
		return dbmgmt.DbInvalidData, errDB
	}

	if imsiprefixProfileResponse.Code != dbmgmt.DbGetSuccess && imsiprefixProfileResponse.Code != dbmgmt.DbDataNotExist {
		err = fmt.Errorf("Fail to get ImsiprefixProfiles, error code %d", imsiprefixProfileResponse.Code)
		return imsiprefixProfileResponse.Code, err
	}

	if imsiprefixProfileResponse.Code == dbmgmt.DbDataNotExist {
		err = fmt.Errorf("ImsiprefixProfile Not Found by imsi %s", imsiStr)
		return imsiprefixProfileResponse.Code, err
	}

	imsiLenStr := strconv.Itoa(len(imsiStr))

	for _, item := range imsiprefixProfileResponse.ValueInfo {
		value := strings.Split(item, ValueInfoSeparator)
		if 4 == len(value) && (imsiLenStr == value[0] || "0" == value[0]) {
			var imsiSearchResult = ImsiSearchResult{}
			imsiSearchResult.NfType = strings.Split(value[3], NfTypeSeparator)
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

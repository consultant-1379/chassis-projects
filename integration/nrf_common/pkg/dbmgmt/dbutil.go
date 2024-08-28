package dbmgmt

import (
	"fmt"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"github.com/buger/jsonparser"
)

// GetExpiredTime is to get expire time of data
func GetExpiredTime(data []byte) (uint64, error) {
	expiredTime, err := jsonparser.GetInt(data, constvalue.ExpiredTime)
	if err != nil {
		return 0, fmt.Errorf("Fail to get expiredTime from nfprofile")
	}
	return uint64(expiredTime), nil
}

// GetLastUpdateTime is to get last update time of data
func GetLastUpdateTime(data []byte) (int64, error) {
	lastUpdateTime, err := jsonparser.GetInt(data, constvalue.LastUpdateTime)
	if err != nil {
		return 0, fmt.Errorf("Fail to get lastUpdateTime from nfprofile")
	}
	return lastUpdateTime, nil
}

// GetProvisionedFlag is to get provisionedFlag of data
func GetProvisionedFlag(data []byte) (int64, error) {
	provisionedFlag, err := jsonparser.GetInt(data, constvalue.ProvisionedFlag)
	if err != nil {
		return 0, fmt.Errorf("Fail to get provisionedFlag from nfprofile")
	}
	return provisionedFlag, nil
}

// GetMd5Sum is to get Md5Sum of data
func GetMd5Sum(data []byte) ([]byte, error) {
	md5sum, valueType, _, err := jsonparser.Get(data, constvalue.MD5SUM)
	if valueType == jsonparser.NotExist || err != nil {
		return []byte(""), fmt.Errorf("Fail to get md5sum from nfprofile")
	}
	return md5sum, nil
}

// GetBody is to get body of data
func GetBody(data []byte) ([]byte, error) {
	body, valueType, _, err := jsonparser.Get(data, constvalue.BODY)
	if valueType == jsonparser.NotExist || err != nil {
		return []byte(""), fmt.Errorf("Fail to get body from nfprofile")
	}
	return body, nil
}

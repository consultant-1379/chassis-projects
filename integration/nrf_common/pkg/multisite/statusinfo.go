package multisite

import (
	"encoding/json"
)

const (
        region string = "ericsson-nrf-multisiteinfo"
)

// StatusInfo struct for data in DB
type StatusInfo struct {
	InstanceID string `json:"id"`
	LastUpdateTime string `json:"lastUpdateTime"`
	Fqdn string `json:"fqdn"`
	Weight float32 `json:"weight"`
}

func encodeMultiSiteInfo(data StatusInfo) (string, error) {
        buf, err := json.Marshal(data)
        if err != nil {
                return "", err
        }
	return string(buf), nil
}

func decodeMultiSiteInfo(value string) (StatusInfo, error) {
	var data StatusInfo
	err := json.Unmarshal([]byte(value), &data)
	return data, err
}

package witness

import (
	"encoding/json"
)

const (
        region string = "ericsson-nrf-witnessinfo"
)

// WitnessInfo struct for witness data in DB
type WitnessInfo struct {
	InstanceID string `json:"id,omitempty"`
	Fqdn string `json:"fqdn,omitempty"`
	LastUpdateTime string `json:"lastUpdateTime"`
}

func encodeWitnessInfo(data WitnessInfo) (string, error) {
        buf, err := json.Marshal(data)
        if err != nil {
                return "", err
        }
	return string(buf), nil
}

func decodeWitnessInfo(value string) (WitnessInfo, error) {
	var data WitnessInfo
	err := json.Unmarshal([]byte(value), &data)
	return data, err
}

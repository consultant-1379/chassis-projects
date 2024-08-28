package nrfschema

import (
	"com/dbproxy/nfmessage/subscription"
	"fmt"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GenerateGrpcPutKey generate a SubKeyStruct used in subscriptionPutIndex
func (s *TSnssai) GenerateGrpcPutKey() *subscription.SubKeyStruct {
	subKey := &subscription.SubKeyStruct{
		SubKey1: fmt.Sprintf("%d", s.Sst),
	}

	if s.Sd != "" {
		subKey.SubKey2 = s.Sd
	} else {
		subKey.SubKey2 = constvalue.Wildcard
	}

	return subKey
}

// GenerateGrpcGetKey generate a SubKeyStructs in subscriptionGetIndex
func (s *TSnssai) GenerateGrpcGetKey() []*subscription.SubKeyStruct {
	var subKeys []*subscription.SubKeyStruct

	if s.Sd != "" {
		subKeys = append(subKeys, &subscription.SubKeyStruct{
			SubKey1: fmt.Sprintf("%d", s.Sst),
			SubKey2: s.Sd,
		})

		subKeys = append(subKeys, &subscription.SubKeyStruct{
			SubKey1: fmt.Sprintf("%d", s.Sst),
			SubKey2: constvalue.Wildcard,
		})
	} else {
		subKeys = append(subKeys, &subscription.SubKeyStruct{
			SubKey1: fmt.Sprintf("%d", s.Sst),
			SubKey2: constvalue.Wildcard,
		})
	}

	return subKeys
}

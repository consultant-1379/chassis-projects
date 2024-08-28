package nrfschema

import (
	"fmt"

	"crypto/md5"
	"encoding/json"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidIPEndPointIndexs return invalid ipEndPoints index
func (n *TNFService) GetInvalidIPEndPointIndexs() []string {
	var invalidIPEndPointIndexs []string

	if n.IpEndPoints != nil {
		index := 0
		for _, ipEndPoint := range n.IpEndPoints {
			if !ipEndPoint.IsValid() {
				invalidIPEndPointIndex := fmt.Sprintf("%s[%d]", constvalue.NFServiceIPEndPoints, index)
				invalidIPEndPointIndexs = append(invalidIPEndPointIndexs, invalidIPEndPointIndex)

			}
			index++
		}
	}

	return invalidIPEndPointIndexs
}

// GetInvalidChfServiceInfoIndex return invalid ChfServiceInfo index
func (n *TNFService) GetInvalidChfServiceInfoIndex() string {
	if n.ChfServiceInfo != nil {
		if !n.ChfServiceInfo.IsValid() {
			return constvalue.NFServiceChfServiceInfo
		}
	}

	return ""
}

// GetInvalidServiceNameIndex returns invalid serviceName index
func (n *TNFService) GetInvalidServiceNameIndex(nfType string) string {
	if _, ok := constvalue.ServiceNameNFTypeMap[n.ServiceName]; !ok {
		return ""
	}

	if nfType == constvalue.ServiceNameNFTypeMap[n.ServiceName] {
		return ""
	}

	return constvalue.NFServiceName
}

// GenerateMd5 generate md5 for NF service
func (n TNFService) GenerateMd5() string {
	n.InterPlmnFqdn = ""
	n.AllowedPlmns = nil
	n.AllowedNfTypes = nil
	n.AllowedNfDomains = nil
	n.AllowedNssais = nil

	body, err := json.Marshal(n)
	if err != nil {
		log.Warnf("Marshal NF service of serviceInstanceId %s failed.", n.ServiceInstanceId)
		return ""
	}

	eTag := md5.Sum(body)
	etagStr := fmt.Sprintf("%x", eTag)
	return etagStr
}

package nrfschema

import (
	"fmt"
	"regexp"
	"strings"

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

// IsAllowedNfType is to check whether nfType is allowed
func (n *TNFService) IsAllowedNfType(nfType string) bool {

	if len(n.AllowedNfTypes) == 0 {
		return true
	}

	if nfType == "" {
		return false
	}

	for _, allowedNfType := range n.AllowedNfTypes {
		if allowedNfType.(string) == nfType {
			return true
		}
	}
	return false
}

// IsAllowedPlmn is to check whether plmnId is allowed
func (n *TNFService) IsAllowedPlmn(plmnID *TPlmnId) bool {

	if len(n.AllowedPlmns) == 0 {
		return true
	}

	if plmnID == nil {
		return false
	}

	for _, allowedPlmnID := range n.AllowedPlmns {
		if allowedPlmnID.Mcc == plmnID.Mcc && allowedPlmnID.Mnc == plmnID.Mnc {
			return true
		}
	}
	return false
}

// IsAllowedNfDomain is to check whether domain is allowed
func (n *TNFService) IsAllowedNfDomain(nfDomain string) bool {
	if len(n.AllowedNfDomains) == 0 {
		return true
	}

	if nfDomain == "" {
		return false
	}

	for _, allowedNfDomainPattern := range n.AllowedNfDomains {

		allowedNfDomainPattern = strings.Replace(allowedNfDomainPattern, `\\`, `\`, -1)
		matched, err := regexp.MatchString(allowedNfDomainPattern, nfDomain)
		if err != nil {
			continue
		}

		if matched {
			return true
		}
	}
	return false
}

// IsAllowedNssai is to check whether domain is allowed
func (n *TNFService) IsAllowedNssai(nssai *TSnssai) bool {

	if len(n.AllowedNssais) == 0 {
		return true
	}

	if nssai == nil {
		return false
	}

	for _, allowedNssai := range n.AllowedNssais {
		if allowedNssai.Sd == "" {
			if allowedNssai.Sst == nssai.Sst {
				return true
			}
		} else {
			if allowedNssai.Sst == nssai.Sst && allowedNssai.Sd == nssai.Sd {
				return true
			}
		}
	}
	return false
}

package nrfschema

import (
	"com/dbproxy/nfmessage/subscription"
	"fmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func (b *TBsfInfo) createNfInfo() string {
	var dnnList, ipDomainList string
	if b.DnnList != nil && len(b.DnnList) > 0 {
		for _, v := range b.DnnList {
			if dnnList == "" {
				dnnList = fmt.Sprintf(`"%v"`, v)
			} else {
				dnnList = fmt.Sprintf(`%s,"%v"`, dnnList, v)
			}
		}
		dnnList = fmt.Sprintf(`"dnnList":[%s]`, dnnList)
	} else {
		dnnList = fmt.Sprintf(`"dnnList":["%s"]`, constvalue.EmptyDnn)
	}

	if b.IpDomainList != nil && len(b.IpDomainList) > 0 {
		for _, v := range b.IpDomainList {
			if ipDomainList == "" {
				ipDomainList = fmt.Sprintf(`"%v"`, v)
			} else {
				ipDomainList = fmt.Sprintf(`%s,"%v"`, ipDomainList, v)
			}
		}
		ipDomainList = fmt.Sprintf(`"ipDomainList":[%s]`, ipDomainList)
	} else {
		ipDomainList = fmt.Sprintf(`"ipDomainList":["%s"]`, constvalue.EmptyIPDomain)
	}

	return fmt.Sprintf(`"bsfInfo":{%s,%s}`, dnnList, ipDomainList)

}

// GenerateNfGroupCond generate NfGroupCond
func (b *TBsfInfo) GenerateNfGroupCond() *subscription.SubKeyStruct {
	return nil
}

// IsEqual is to check if NFInfo is equal
func (b *TBsfInfo) IsEqual(c TNfInfo) bool {

	a := c.(*TBsfInfo)

	if len(a.DnnList) != len(b.DnnList) {
		return false
	}

	if len(a.IpDomainList) != len(b.IpDomainList) {
		return false
	}

	if len(a.Ipv4AddressRanges) != len(b.Ipv4AddressRanges) {
		return false
	}

	if len(a.Ipv6PrefixRanges) != len(b.Ipv6PrefixRanges) {
		return false
	}

	for k, item := range a.DnnList {
		if item != b.DnnList[k] {
			return false
		}
	}

	for k, item := range a.IpDomainList {
		if item != b.IpDomainList[k] {
			return false
		}
	}

	for k, item := range a.Ipv4AddressRanges {
		bb := b.Ipv4AddressRanges[k]
		if item.Start != bb.Start || item.End != bb.End {
			return false
		}
	}

	for k, item := range a.Ipv6PrefixRanges {
		bb := b.Ipv6PrefixRanges[k]
		if item.Start != bb.Start || item.End != bb.End {
			return false
		}
	}

	return true
}

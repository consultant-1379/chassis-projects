package nrfschema

import (
	"fmt"
)

func (b *TBsfInfo) createNfInfo() string {
	var bsfInfo, dnnList, ipDomainList string
	if b.DnnList != nil && len(b.DnnList) > 0 {
		for _, v := range b.DnnList {
			if dnnList == "" {
				dnnList = fmt.Sprintf(`"%v"`, v)
			} else {
				dnnList = fmt.Sprintf(`%s,"%v"`, dnnList, v)
			}
		}
		dnnList = fmt.Sprintf(`"dnnList":[%s]`, dnnList)
	}
	if dnnList != "" {
		bsfInfo = dnnList
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
	}
	if ipDomainList != "" {
		bsfInfo = fmt.Sprintf(`%s,%s`, bsfInfo, ipDomainList)
	}

	return fmt.Sprintf(`"bsfInfo":{%s}`, bsfInfo)

}

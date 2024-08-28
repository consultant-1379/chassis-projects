package nrfschema

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

// Validate validates rules defined in <<3GPP TS 29.510>>
func (n *TNFProfile) Validate() *problemdetails.ProblemDetails {
	problemDetails := n.ValidateCommon()
	if problemDetails != nil {
		return problemDetails
	}

	problemDetails = n.ValidateService()
	if problemDetails != nil {
		return problemDetails
	}

	switch n.NfType {
	case constvalue.NfTypeAMF:
		problemDetails = n.ValidateAmf()
		if problemDetails != nil {
			return problemDetails
		}
	case constvalue.NfTypeAUSF:
		problemDetails = n.ValidateAusf()
		if problemDetails != nil {
			return problemDetails
		}
	case constvalue.NfTypeCHF:
		problemDetails = n.ValidateChf()
		if problemDetails != nil {
			return problemDetails
		}
	case constvalue.NfTypePCF:
		problemDetails = n.ValidatePcf()
		if problemDetails != nil {
			return problemDetails
		}
	case constvalue.NfTypeSMF:
		problemDetails = n.ValidateSmf()
		if problemDetails != nil {
			return problemDetails
		}
	case constvalue.NfTypeUDM:
		problemDetails = n.ValidateUdm()
		if problemDetails != nil {
			return problemDetails
		}
	case constvalue.NfTypeUDR:
		problemDetails = n.ValidateUdr()
		if problemDetails != nil {
			return problemDetails
		}
	case constvalue.NfTypeUPF:
		problemDetails = n.ValidateUpf()
		if problemDetails != nil {
			return problemDetails
		}
	}

	return nil
}

// ValidateCommon validates common part of NF profile
func (n *TNFProfile) ValidateCommon() *problemdetails.ProblemDetails {
	ok := false
	if n.Fqdn != "" || n.Ipv4Addresses != nil || n.Ipv6Addresses != nil {
		ok = true
	}

	if !ok {
		return &problemdetails.ProblemDetails{
			Title: "not a valid nf profile",
			InvalidParams: []*problemdetails.InvalidParam{
				&problemdetails.InvalidParam{
					Param:  constvalue.Fqdn,
					Reason: constvalue.NfProfileRule1,
				},
				&problemdetails.InvalidParam{
					Param:  constvalue.Ipv4Addresses,
					Reason: constvalue.NfProfileRule1,
				},
				&problemdetails.InvalidParam{
					Param:  constvalue.Ipv6Addresses,
					Reason: constvalue.NfProfileRule1,
				},
			},
		}
	}

	problemdetails := n.ValidatePlmnLit()

	if problemdetails != nil {
		return problemdetails
	}

	return nil
}

func (n *TNFProfile) ValidatePlmnLit() *problemdetails.ProblemDetails {
	if n.PlmnList == nil {
		return nil
	}

	isLocalPlmn := false
	for _, localPlmn := range cm.NfProfile.PlmnID {
		localPlmnID := localPlmn.GetPlmnID()
		for _, requestPlmn := range n.PlmnList {
			if requestPlmn.GetPlmnID() == localPlmnID {
				isLocalPlmn = true
				break
			}
		}

		if isLocalPlmn {
			break
		}
	}

	if !isLocalPlmn {
		localPlmnByte, _ := json.Marshal(cm.NfProfile.PlmnID)
		localPlmnStr := strings.Replace(string(localPlmnByte), `"`, "", -1)

		return &problemdetails.ProblemDetails{
			Title:  constvalue.ForbiddenUnlocalTitle,
			Detail: fmt.Sprintf("The local plmn list is %s", string(localPlmnStr)),
		}
	}

	return nil
}

// ValidateService validates ipEndPoints
func (n *TNFProfile) ValidateService() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam
	if n.NfServices != nil {
		index := 0
		for _, nfService := range n.NfServices {
			invalidIPEndPointIndexs := nfService.GetInvalidIPEndPointIndexs()
			if invalidIPEndPointIndexs != nil {
				for _, invalidIPEndPointIndex := range invalidIPEndPointIndexs {
					invalidParam := &problemdetails.InvalidParam{
						Param:  fmt.Sprintf("%s[%d].%s", constvalue.NfServices, index, invalidIPEndPointIndex),
						Reason: constvalue.NfProfileRule2,
					}
					invalidParams = append(invalidParams, invalidParam)
				}
			}

			InvalidChfServiceInfoIndex := nfService.GetInvalidChfServiceInfoIndex()
			if InvalidChfServiceInfoIndex != "" {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s[%d].%s", constvalue.NfServices, index, InvalidChfServiceInfoIndex),
					Reason: constvalue.NfProfileRule8,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

			InvalidServiceNameIndex := nfService.GetInvalidServiceNameIndex(n.NfType)
			if InvalidServiceNameIndex != "" {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s[%d].%s", constvalue.NfServices, index, InvalidServiceNameIndex),
					Reason: fmt.Sprintf(constvalue.NfProfileRule10, nfService.ServiceName, n.NfType),
				}
				invalidParams = append(invalidParams, invalidParam)
			}

			index++
		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// ValidateAmf validates rule for AMF
func (n *TNFProfile) ValidateAmf() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam

	if n.AmfInfo != nil {
		invalidTaiRangeIndexs := n.AmfInfo.GetInvalidTaiRangeIndexs()
		if invalidTaiRangeIndexs != nil {
			for _, invalidTaiRangeIndex := range invalidTaiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.AmfInfo, invalidTaiRangeIndex),
					Reason: constvalue.NfProfileRule6,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}

		invalidN2InterfaceAmfInfoIndex := n.AmfInfo.GetInvalidN2InterfaceAmfInfoIndex()
		if invalidN2InterfaceAmfInfoIndex != "" {
			invalidParam := &problemdetails.InvalidParam{
				Param:  fmt.Sprintf("%s.%s", constvalue.AmfInfo, invalidN2InterfaceAmfInfoIndex),
				Reason: constvalue.NfProfileRule7,
			}
			invalidParams = append(invalidParams, invalidParam)
		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// ValidateAusf validates rule for AUSF
func (n *TNFProfile) ValidateAusf() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam

	if n.AusfInfo != nil {
		invalidSupiRangeIndexs := n.AusfInfo.GetInvalidSupiRangeIndexs()
		if invalidSupiRangeIndexs != nil {
			for _, invalidSupiRangeIndex := range invalidSupiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.AusfInfo, invalidSupiRangeIndex),
					Reason: constvalue.NfProfileRule3,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// ValidateChf validates rule for CHF
func (n *TNFProfile) ValidateChf() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam

	if n.ChfInfo != nil {
		invalidSupiRangeIndexs := n.ChfInfo.GetInvalidSupiRangeIndexs()
		if invalidSupiRangeIndexs != nil {
			for _, invalidSupiRangeIndex := range invalidSupiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.ChfInfo, invalidSupiRangeIndex),
					Reason: constvalue.NfProfileRule3,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}

		invalidGpsiRangeIndexs := n.ChfInfo.GetInvalidGpsiRangeIndexs()
		if invalidGpsiRangeIndexs != nil {
			for _, invalidGpsiRangeIndex := range invalidGpsiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.ChfInfo, invalidGpsiRangeIndex),
					Reason: constvalue.NfProfileRule4,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}

		invalidPlmnRangeIndexs := n.ChfInfo.GetInvalidPlmnRangeIndexs()
		if invalidPlmnRangeIndexs != nil {
			for _, invalidPlmnRangeIndex := range invalidPlmnRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.ChfInfo, invalidPlmnRangeIndex),
					Reason: constvalue.NfProfileRule9,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// ValidatePcf validates rule for PCF
func (n *TNFProfile) ValidatePcf() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam

	if n.PcfInfo != nil {
		invalidSupiRangeIndexs := n.PcfInfo.GetInvalidSupiRangeIndexs()
		if invalidSupiRangeIndexs != nil {
			for _, invalidSupiRangeIndex := range invalidSupiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.PcfInfo, invalidSupiRangeIndex),
					Reason: constvalue.NfProfileRule3,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// ValidateSmf validates rule for SMF
func (n *TNFProfile) ValidateSmf() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam

	if n.SmfInfo != nil {
		invalidTaiRangeIndexs := n.SmfInfo.GetInvalidTaiRangeIndexs()
		if invalidTaiRangeIndexs != nil {
			for _, invalidTaiRangeIndex := range invalidTaiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.SmfInfo, invalidTaiRangeIndex),
					Reason: constvalue.NfProfileRule6,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// ValidateUdm validates rule for UDM
func (n *TNFProfile) ValidateUdm() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam

	if n.UdmInfo != nil {
		invalidSupiRangeIndexs := n.UdmInfo.GetInvalidSupiRangeIndexs()
		if invalidSupiRangeIndexs != nil {
			for _, invalidSupiRangeIndex := range invalidSupiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.UdmInfo, invalidSupiRangeIndex),
					Reason: constvalue.NfProfileRule3,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}

		invalidGpsiRangeIndexs := n.UdmInfo.GetInvalidGpsiRangeIndexs()
		if invalidGpsiRangeIndexs != nil {
			for _, invalidGpsiRangeIndex := range invalidGpsiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.UdmInfo, invalidGpsiRangeIndex),
					Reason: constvalue.NfProfileRule4,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}

		invalidEGIRangeIndexs := n.UdmInfo.GetInvalidEGIRangeIndexs()
		if invalidEGIRangeIndexs != nil {
			for _, invalidEGIRangeIndex := range invalidEGIRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.UdmInfo, invalidEGIRangeIndex),
					Reason: constvalue.NfProfileRule4,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// ValidateUdr validates rule for UDR
func (n *TNFProfile) ValidateUdr() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam

	if n.UdrInfo != nil {
		invalidSupiRangeIndexs := n.UdrInfo.GetInvalidSupiRangeIndexs()
		if invalidSupiRangeIndexs != nil {
			for _, invalidSupiRangeIndex := range invalidSupiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.UdrInfo, invalidSupiRangeIndex),
					Reason: constvalue.NfProfileRule3,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}

		invalidGpsiRangeIndexs := n.UdrInfo.GetInvalidGpsiRangeIndexs()
		if invalidGpsiRangeIndexs != nil {
			for _, invalidGpsiRangeIndex := range invalidGpsiRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.UdrInfo, invalidGpsiRangeIndex),
					Reason: constvalue.NfProfileRule4,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}

		invalidEGIRangeIndexs := n.UdrInfo.GetInvalidEGIRangeIndexs()
		if invalidEGIRangeIndexs != nil {
			for _, invalidEGIRangeIndex := range invalidEGIRangeIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.UdrInfo, invalidEGIRangeIndex),
					Reason: constvalue.NfProfileRule4,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// ValidateUpf validates rule for UPF
func (n *TNFProfile) ValidateUpf() *problemdetails.ProblemDetails {
	var invalidParams []*problemdetails.InvalidParam

	if n.UpfInfo != nil {
		invalidInterfaceUpfInfoIndexs := n.UpfInfo.GetInvalidInterfaceUpfInfoIndexs()
		if invalidInterfaceUpfInfoIndexs != nil {
			for _, invalidInterfaceUpfInfoIndex := range invalidInterfaceUpfInfoIndexs {
				invalidParam := &problemdetails.InvalidParam{
					Param:  fmt.Sprintf("%s.%s", constvalue.UpfInfo, invalidInterfaceUpfInfoIndex),
					Reason: constvalue.NfProfileRule5,
				}
				invalidParams = append(invalidParams, invalidParam)
			}

		}
	}

	if invalidParams != nil {
		return &problemdetails.ProblemDetails{
			Title:         "not a valid nf profile",
			InvalidParams: invalidParams,
		}
	}

	return nil
}

// GetServiceNames returns a list of serviceName from nfServices
func (n *TNFProfile) GetServiceNames() []string {
	var serviceNames []string

	if n.NfServices != nil && len(n.NfServices) > 0 {
		for _, item := range n.NfServices {
			if item.ServiceName != "" {
				serviceNames = append(serviceNames, item.ServiceName)
			}
		}
	}
	return serviceNames
}

// GenerateMd5 generates md5 for NF profile
func (n TNFProfile) GenerateMd5() string {
	n.NfServices = nil
	body, err := json.Marshal(n)
	if err != nil {
		log.Warnf("Marshal NF profile of nfInstance %s failed.", n.NfInstanceId)
		return ""
	}

	eTag := md5.Sum(body)
	etagStr := fmt.Sprintf("%x", eTag)
	return etagStr
}

// CreateHelperInfo creates help information for NF profile
func (n *TNFProfile) CreateHelperInfo(profileUpdateTime uint64) string {
	bodyCommon := n.createProfileCommonInfo(profileUpdateTime)
	sNssais := n.createSnssaisHelperInfo()
	nfService := n.createNfServiceHelperInfo()
	specificTypeInfo := ""

	if n.NfType == constvalue.NfTypeUDM {
		specificTypeInfo = n.createUdmInfo()

	} else if n.NfType == constvalue.NfTypeAMF {
		specificTypeInfo = n.createAmfInfo()

	} else if n.NfType == constvalue.NfTypeSMF {
		specificTypeInfo = n.createSmfInfo()

	} else if n.NfType == constvalue.NfTypeAUSF {
		specificTypeInfo = n.createAusfInfo()

	} else if n.NfType == constvalue.NfTypePCF {
		specificTypeInfo = n.createPcfInfo()

	} else if n.NfType == constvalue.NfTypeUDR {

		specificTypeInfo = n.createUdrInfo()
	} else if n.NfType == constvalue.NfTypeUPF {
		specificTypeInfo = n.createUpfInfo()
	} else if n.NfType == constvalue.NfTypeCHF {
		specificTypeInfo = n.createChfInfo()
	} else if n.NfType == constvalue.NfTypeBSF {
		specificTypeInfo = n.createBsfInfo()
	}

	helper := ""
	if nfService != "" {
		helper = fmt.Sprintf("%s", nfService)
	}
	if specificTypeInfo != "" {
		if helper == "" {
			helper = fmt.Sprintf("%s", specificTypeInfo)
		}
		helper = fmt.Sprintf("%s,%s", helper, specificTypeInfo)
	}
	if sNssais != "" {
		if helper == "" {
			helper = fmt.Sprintf("%s", sNssais)
		}
		helper = fmt.Sprintf("%s,%s", helper, sNssais)
	}
	if bodyCommon != "" {
		if helper == "" {
			helper = fmt.Sprintf("%s", bodyCommon)
		}
		helper = fmt.Sprintf("%s,%s", helper, bodyCommon)
	}

	return helper

}

func (n *TNFProfile) createNfServiceHelperInfo() string {

	allowedDomains := make(map[string]bool)
	allowedTypes := make(map[string]bool)
	allowedPlmns := make(map[string]bool)

	if n.NfServices != nil && len(n.NfServices) > 0 {

		for _, item := range n.NfServices {

			if item.AllowedNfDomains != nil && len(item.AllowedNfDomains) > 0 {
				for _, domain := range item.AllowedNfDomains {
					allowedDomains[domain] = true
				}
			} else {
				allowedDomains[constvalue.EmptyAllowedDomain] = true
			}

			if item.AllowedNfTypes != nil && len(item.AllowedNfTypes) > 0 {
				for _, nftype := range item.AllowedNfTypes {
					allowedTypes[nftype.(string)] = true
				}
			} else {
				allowedTypes[constvalue.EmptyAllowedNfType] = true
			}

			if item.AllowedPlmns != nil && len(item.AllowedPlmns) > 0 {
				for _, plmn := range item.AllowedPlmns {
					allowedPlmns[plmn.Mcc+plmn.Mnc] = true
				}
			} else {
				allowedPlmns[constvalue.EmptyAllowedPlmnMcc+constvalue.EmptyAllowedPlmnMnc] = true
			}

		}

	} else {
		allowedDomains[constvalue.EmptyAllowedDomain] = true
		allowedTypes[constvalue.EmptyAllowedNfType] = true
		allowedPlmns[constvalue.EmptyAllowedPlmnMcc+constvalue.EmptyAllowedPlmnMnc] = true
	}

	allowedDomainStr := ""
	for k := range allowedDomains {
		if allowedDomainStr == "" {
			allowedDomainStr = fmt.Sprintf(`"%s"`, k)
		} else {
			allowedDomainStr = allowedDomainStr + "," + fmt.Sprintf(`"%s"`, k)
		}
	}

	allowedNfTypeStr := ""
	for k := range allowedTypes {
		if allowedNfTypeStr == "" {
			allowedNfTypeStr = fmt.Sprintf(`"%s"`, k)
		} else {
			allowedNfTypeStr = allowedNfTypeStr + "," + fmt.Sprintf(`"%s"`, k)
		}
	}

	allowedPlmnStr := ""
	for k := range allowedPlmns {
		plmn := []rune(k)
		mcc := string(plmn[0:3])
		mnc := string(plmn[3:])
		if allowedPlmnStr == "" {
			allowedPlmnStr = fmt.Sprintf(`{"mcc":"%s", "mnc":"%s"}`, mcc, mnc)
		} else {
			allowedPlmnStr = allowedPlmnStr + "," + fmt.Sprintf(`{"mcc":"%s", "mnc":"%s"}`, mcc, mnc)
		}
	}

	allowedNfDomain := fmt.Sprintf(`"allowedNfDomains": [%s]`, allowedDomainStr)
	allowedNfType := fmt.Sprintf(`"allowedNfTypes": [%s]`, allowedNfTypeStr)
	allowedPlmn := fmt.Sprintf(`"allowedPlmns": [%s]`, allowedPlmnStr)
	nfService := fmt.Sprintf(`"nfServices": {%s,%s,%s}`, allowedNfDomain, allowedNfType, allowedPlmn)

	return nfService
}

func (n *TNFProfile) createSnssaisHelperInfo() string {
	if n.SNssais != nil && len(n.SNssais) > 0 {
		sNssaisList := ""
		for _, item := range n.SNssais {
			sst := item.Sst
			sd := strings.ToLower(item.Sd)
			if sd == "" {
				sd = constvalue.EmptySd
			}
			if sNssaisList != "" {
				sNssaisList += ","
			}
			sNssaisList += fmt.Sprintf(`{"sst":%d,"sd":"%s"}`, sst, sd)
		}
		return fmt.Sprintf(`"sNssais":[%s]`, sNssaisList)
	}

	return fmt.Sprintf(`"sNssais":[{"sst":%d}]`, constvalue.EmptySst)
}

func (n *TNFProfile) createProfileCommonInfo(profileUpdateTime uint64) string {
	var bodyCommon string
	if n.NfType != "" {
		if bodyCommon == "" {
			bodyCommon = fmt.Sprintf(`"nfType":"%s"`, n.NfType)
		} else {
			bodyCommon += fmt.Sprintf(`,"nfType":"%s"`, n.NfType)
		}
	}
	if n.NfStatus != "" {
		if bodyCommon == "" {
			bodyCommon = fmt.Sprintf(`"nfStatus":"%s"`, n.NfStatus)
		} else {
			bodyCommon += fmt.Sprintf(`,"nfStatus":"%s"`, n.NfStatus)
		}
	}
	if n.NfInstanceId != "" {
		if bodyCommon == "" {
			bodyCommon = fmt.Sprintf(`"nfInstanceId":"%s"`, n.NfInstanceId)
		} else {
			bodyCommon += fmt.Sprintf(`,"nfInstanceId":"%s"`, n.NfInstanceId)
		}
	}
	if n.Fqdn != "" {
		if bodyCommon == "" {
			bodyCommon = fmt.Sprintf(`"fqdn":"%s"`, n.Fqdn)
		} else {
			bodyCommon += fmt.Sprintf(`,"fqdn":"%s"`, n.Fqdn)
		}
	}
	if n.NsiList != nil && len(n.NsiList) > 0 {
		if bodyCommon == "" {
			bodyCommon = fmt.Sprintf(`"nsiList":"%v"`, n.NsiList)
		} else {
			bodyCommon += fmt.Sprintf(`,"nsiList":"%v"`, n.NsiList)
		}
	}
	if bodyCommon == "" {
		bodyCommon = fmt.Sprintf(`"profileUpdateTime":"%v"`, profileUpdateTime)
	} else {
		bodyCommon += fmt.Sprintf(`,"profileUpdateTime":"%v"`, profileUpdateTime)
	}
	return bodyCommon
}

func (n *TNFProfile) createUdmInfo() string {
	if n.UdmInfo == nil {
		supiRanges := fmt.Sprintf(`"supiRanges":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
		gpsiRanges := fmt.Sprintf(`"gpsiRanges":[{"pattern":"%s"}]`, constvalue.EmptyGpsiRangePattern)
		externalID := fmt.Sprintf(`"externalGroupIdentifiersRanges":[{"pattern":"%s"}]`, constvalue.EmptyExternalIDPattern)
		groupID := fmt.Sprintf(`"groupId":"%s"`, constvalue.EmptyGroupID)
		supiMatchAll := fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
		gpsiMatchAll := fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.MatchAll)
		return fmt.Sprintf(`"udmInfo":{%s,%s, %s, %s,%s,%s}`, supiRanges, gpsiRanges, externalID, groupID, supiMatchAll, gpsiMatchAll)
	}

	return n.UdmInfo.createNfInfo()
}

func (n *TNFProfile) createAmfInfo() string {
	if n.AmfInfo == nil {
		plmnid := fmt.Sprintf(`{"mcc":"%s","mnc":"%s"}`, constvalue.EmptyMcc, constvalue.EmptyMnc)
		taiList := fmt.Sprintf(`"taiList":[{"plmnId":%s, "tac": "%s"}]`, plmnid, constvalue.EmptyTac)

		tacRange := fmt.Sprintf(`{"pattern":"%s"}`, constvalue.EmptyTacRangePattern)
		taiRangeList := fmt.Sprintf(`"taiRangeList":[{"plmnId":%s, "tacRangeList": [%s]}]`, plmnid, tacRange)

		return fmt.Sprintf(`"amfInfo":{%s,%s}`, taiList, taiRangeList)
	}

	return n.AmfInfo.createNfInfo()
}

func (n *TNFProfile) createSmfInfo() string {
	if n.SmfInfo == nil {
		plmnid := fmt.Sprintf(`{"mcc":"%s","mnc":"%s"}`, constvalue.EmptyMcc, constvalue.EmptyMnc)
		taiList := fmt.Sprintf(`"taiList":[{"plmnId":%s, "tac": "%s"}]`, plmnid, constvalue.EmptyTac)
		tacRange := fmt.Sprintf(`{"pattern":"%s"}`, constvalue.EmptyTacRangePattern)
		taiRangeList := fmt.Sprintf(`"taiRangeList":[{"plmnId":%s, "tacRangeList": [%s]}]`, plmnid, tacRange)
		sNssaiSmffInfoList := fmt.Sprintf(`"sNssaiSmfInfoList":[{"sNssai":{"sst":%d,"sd":"%s"}}]`, constvalue.EmptySst, constvalue.EmptySd)
		return fmt.Sprintf(`"smfInfo":{%s,%s, %s}`, taiList, taiRangeList, sNssaiSmffInfoList)
	}

	return n.SmfInfo.createNfInfo()
}

func (n *TNFProfile) createAusfInfo() string {
	if n.AusfInfo == nil {
		supiRanges := fmt.Sprintf(`"supiRanges":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
		groupID := fmt.Sprintf(`"groupId":"%s"`, constvalue.EmptyGroupID)
		supiMatchAll := fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
		return fmt.Sprintf(`"ausfInfo":{%s, %s,%s}`, supiRanges, groupID, supiMatchAll)
	}

	return n.AusfInfo.createNfInfo()
}

func (n *TNFProfile) createPcfInfo() string {
	if n.PcfInfo == nil {

		supiRanges := fmt.Sprintf(`"supiRanges":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
		groupID := fmt.Sprintf(`"groupId":"%s"`, constvalue.EmptyGroupID)
		supiMatchAll := fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
		return fmt.Sprintf(`"pcfInfo":{%s, %s,%s}`, supiRanges, groupID, supiMatchAll)
	}

	return n.PcfInfo.createNfInfo()
}

func (n *TNFProfile) createUdrInfo() string {
	if n.UdrInfo == nil {
		supiRanges := fmt.Sprintf(`"supiRanges":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
		gpsiRanges := fmt.Sprintf(`"gpsiRanges":[{"pattern":"%s"}]`, constvalue.EmptyGpsiRangePattern)
		externalID := fmt.Sprintf(`"externalGroupIdentifiersRanges":[{"pattern":"%s"}]`, constvalue.EmptyExternalIDPattern)
		groupID := fmt.Sprintf(`"groupId":"%s"`, constvalue.EmptyGroupID)
		supiMatchAll := fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
		gpsiMatchAll := fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.MatchAll)
		return fmt.Sprintf(`"udrInfo":{%s,%s, %s,%s,%s,%s}`, supiRanges, gpsiRanges, externalID, groupID, supiMatchAll, gpsiMatchAll)
	}

	return n.UdrInfo.createNfInfo()
}

func (n *TNFProfile) createUpfInfo() string {
	if n.UpfInfo == nil {
		sNssaiUpfInfoList := fmt.Sprintf(`"sNssaiUpfInfoList":[{"sNssai":{"sst":%d,"sd":"%s"},"dnnUpfInfoList":[{"dnaiList":["%s"]}]}]`, constvalue.EmptySst, constvalue.EmptySd, constvalue.EmptyDnai)
		return fmt.Sprintf(`"upfInfo":{%s}`, sNssaiUpfInfoList)
	}
	return n.UpfInfo.createNfInfo()
}

func (n *TNFProfile) createChfInfo() string {
	if n.ChfInfo == nil {
		supiRangeList := fmt.Sprintf(`"supiRangeList":[{"pattern":"%s"}]`, constvalue.EmptySupiRangePattern)
		gpsiRangeList := fmt.Sprintf(`"gpsiRangeList":[{"pattern":"%s"}]`, constvalue.EmptyGpsiRangePattern)
		plmnRangeList := fmt.Sprintf(`"plmnRangeList":[{"pattern":"%s"}]`, constvalue.EmptyPlmnRangePattern)
		supiMatchAll := fmt.Sprintf(`"supiMatchAll":"%s"`, constvalue.MatchAll)
		gpsiMatchAll := fmt.Sprintf(`"gpsiMatchAll":"%s"`, constvalue.MatchAll)
		return fmt.Sprintf(`"chfInfo":{%s,%s,%s,%s,%s}`, supiRangeList, gpsiRangeList, plmnRangeList, supiMatchAll, gpsiMatchAll)
	}
	return n.ChfInfo.createNfInfo()
}

func (n *TNFProfile) createBsfInfo() string {
	if n.BsfInfo == nil {
		return ""
	}
	return n.BsfInfo.createNfInfo()
}

// Equal is used to check whether two nfprofiles are the same
func (n *TNFProfile) Equal(nfProfile *TNFProfile) bool {
	return reflect.DeepEqual(n, nfProfile)
}

// IsNfInfoExist check whether nfInfo should be summarized in nrfInfo
func (n *TNFProfile) IsNfInfoExist() bool {
	if _, ok := constvalue.NFInfoMap[n.NfType]; !ok {
		return false
	}

	ok := false

	switch n.NfType {
	case constvalue.NfTypeAMF:
		if n.AmfInfo != nil {
			ok = true
		}
	case constvalue.NfTypeAUSF:
		if n.AusfInfo != nil {
			ok = true
		}
	case constvalue.NfTypePCF:
		if n.PcfInfo != nil {
			ok = true
		}
	case constvalue.NfTypeSMF:
		if n.SmfInfo != nil {
			ok = true
		}
	case constvalue.NfTypeUDM:
		if n.UdmInfo != nil {
			ok = true
		}
	case constvalue.NfTypeNRF:
		if n.NrfInfo != nil {
			ok = true
		}
	}

	return ok
}

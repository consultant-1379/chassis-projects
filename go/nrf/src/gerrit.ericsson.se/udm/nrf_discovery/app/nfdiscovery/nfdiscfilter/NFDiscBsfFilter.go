package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"github.com/buger/jsonparser"
	"strings"
	"strconv"
	"net"
	"encoding/hex"
)

//NFBSFInfoFilter to process bsfinfo filter in nfprofile
type NFBSFInfoFilter struct {

}

func (a *NFBSFInfoFilter) filter(nfInfo []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	if queryForm.GetNRFDiscIPV4Addr(constvalue.SearchDataUEIPv4Addr) != "" {
		log.Debugf("Search nfProfile with ue-ipv4-addr : %s", queryForm.GetNRFDiscIPV4Addr(constvalue.SearchDataUEIPv4Addr))
		if len(nfInfo) > 0 {
			if !a.isMatchedIPv4Addr(queryForm.GetNRFDiscIPV4Addr(constvalue.SearchDataUEIPv4Addr), nfInfo) {
				log.Debugf("No Matched nfProfile with ue-ipv4-addr : %s", queryForm.GetNRFDiscIPV4Addr(constvalue.SearchDataUEIPv4Addr))
				return false
			}
			log.Debugf("Matched nfProfile is Found with ue-ipv4-addr :%s", queryForm.GetNRFDiscIPV4Addr(constvalue.SearchDataUEIPv4Addr))
		}
	}

	if queryForm.GetNRFDiscIPV6Prefix(constvalue.SearchDataUEIPv6Prefix) != "" {
		log.Debugf("Search nfProfile with ue-ipv6-prefix : %s", queryForm.GetNRFDiscIPV6Prefix(constvalue.SearchDataUEIPv6Prefix))
		if len(nfInfo) > 0 {
			if !a.isMatchedIPv6AddrPrefix(queryForm.GetNRFDiscIPV6Prefix(constvalue.SearchDataUEIPv6Prefix), nfInfo) {
				log.Debugf("No Matched nfProfile with ue-ipv6-prefix : %s", queryForm.GetNRFDiscIPV6Prefix(constvalue.SearchDataUEIPv6Prefix))
				return false
			}
			log.Debugf("Matched nfProfile is Found with ue-ipv6-prefix :%s", queryForm.GetNRFDiscIPV6Prefix(constvalue.SearchDataUEIPv6Prefix))
		}

	}
	return true
}

func (a *NFBSFInfoFilter) fullIPv6Addres(ip net.IP, prefix int) string {

	dst := make([]byte, hex.EncodedLen(len(ip)))
	_ = hex.Encode(dst, ip)
	if 0 == prefix%4 {
		return (string(dst[0:(prefix / 4)]))
	}
	return (string(dst[0:(prefix/4 + 1)]))
}


func (a *NFBSFInfoFilter) isMatchedIPv6PrefixRange(prefix, start, end string) bool {

	_, addrPrefix, err1 := net.ParseCIDR(prefix)
	if err1 != nil {
		log.Debugf("ipv6 prefix parse error, err=%v", err1)
		return false
	}
	_, sPrefix, err2 := net.ParseCIDR(start)
	if err2 != nil {
		log.Debugf("ipv6 start parse error, err=%v", err2)
		return false
	}
	_, ePrefix, err3 := net.ParseCIDR(end)
	if err3 != nil {
		log.Debugf("ipv6 end parse error, err=%v", err3)
		return false
	}
	addrLen, _ := addrPrefix.Mask.Size()
	sLen, _ := sPrefix.Mask.Size()
	eLen, _ := ePrefix.Mask.Size()

	//minLen := a.minPrefixLen(addrLen, sLen, eLen)
	addr := a.fullIPv6Addres(net.ParseIP(addrPrefix.IP.String()), addrLen)
	s := a.fullIPv6Addres(net.ParseIP(sPrefix.IP.String()), sLen)
	e := a.fullIPv6Addres(net.ParseIP(ePrefix.IP.String()), eLen)
	addrS := addr
	addrE := addr

	if addrLen < sLen {
		num := (len(s) - len(addr))
		for i := 0; i < num; i++ {
			addrS += "0"
		}
	} else if addrLen > sLen {
		num := len(addr) - len(s)
		for i := 0; i < num; i++ {
			s += "0"
		}
	}

	if addrLen < eLen {
		num := len(e) - len(addr)
		for i := 0; i < num; i++ {
			addrE += "f"
		}
	} else if (addrLen > eLen) {
		num := len(addr) - len(e)
		for i := 0; i < num ; i++{
			e += "f"
		}
	}

        log.Debugf("IPv6Prefix start: %s, start: %s, IPv6Prefix end: %s, end: %s", addrS, s, addrE, e)
	if  (1 != strings.Compare(s, addrS)) && (1 != strings.Compare(addrE, e)) {
		return true
	}

	return false
}

func (a *NFBSFInfoFilter) isMatchedIPv6AddrPrefix(ipv6AddrPrefix string, nfProfile []byte) bool {
	ret := false


	_, err := jsonparser.ArrayEach(nfProfile, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if ret {
			return
		}

		s, err := jsonparser.GetString(value, "start")
		if err != nil {
			return
		}

		e, err := jsonparser.GetString(value, "end")
		if err != nil {
			return
		}

                if a.isMatchedIPv6PrefixRange(ipv6AddrPrefix, s, e) {
			ret = true
		}

	}, constvalue.IPv6PrefixRanges)

	if err != nil {
		ret = false
		log.Debugf("Parsering array fail for %s, error: %v", constvalue.IPv4AddressRanges, err)
	}

	return ret
}

func (a *NFBSFInfoFilter)isIPv4AddrInRange(start []string, end []string, ipv4 []string) bool {
	for i := 0; i < 4; i++ {
		s, err1 := strconv.ParseInt(start[i], 10, 32)
		if err1 != nil {
			log.Debugf("ipv4 start parseint error, err=%v", err1)
			return false
		}
		e, err2 := strconv.ParseInt(end[i], 10, 32)
		if err2 != nil {
			log.Debugf("ipv4 end parseint error, err=%v", err2)
			return false
		}
		i, err3 := strconv.ParseInt(ipv4[i], 10, 32)
		if err3 != nil {
			log.Debugf("ipv4 parseint error, err=%v", err3)
			return false
		}

		if i >= s && i <= e {
			continue
		} else {
			return false
		}
	}

	return true
}

func (a *NFBSFInfoFilter)  isMatchedIPv4Addr(ipv4Addr string, nfInfo []byte) bool {
	ret := false
	log.Debugf("ipv4addr: %s", string(nfInfo))
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if ret {
			return
		}

		s, err := jsonparser.GetString(value, "start")
		if err != nil {
			return
		}

		e, err := jsonparser.GetString(value, "end")
		if err != nil {
			return
		}

		start := strings.Split(s, ".")
		end := strings.Split(e, ".")
		ipv4 := strings.Split(ipv4Addr, ".")
		if len(start) != 4 || len(end) != 4 || len(ipv4) != 4 {
			return
		}

		if a.isIPv4AddrInRange(start, end, ipv4) {
			ret = true
			return
		}

	}, constvalue.IPv4AddressRanges)

	if err != nil {
		ret = false
		log.Debugf("Parsering array fail for %s, error: %v", constvalue.IPv4AddressRanges, err)
	}

	return ret
}

//isMatchedIPDomain is to match ip-domain in bsfInfo
func (a *NFBSFInfoFilter) isMatchedIPDomain(ipDomain string, nfInfo []byte) bool {
	matched := false
	log.Debugf("nfInfo: %s", string(nfInfo))
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if matched {
			return
		}
		log.Debugf("ipdomain In Profile value: %s", string(value[:]))
		ipDomainInProfile := string(value[:])
		if ipDomain == ipDomainInProfile {
			matched = true
			return
		}
	}, constvalue.IPDomainList)

	if err != nil {
		matched = false
	}

	return matched
}

func (a *NFBSFInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFBSFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression
	dnn := queryForm.GetNRFDiscDnnValue()
	if "" != dnn {
		dnnPath := getParamSearchPath(constvalue.NfTypeBSF, constvalue.SearchDataDnn)
		dnnExpression := buildStringSearchParameter(dnnPath, dnn, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, dnnExpression)
	}

	ipDomain := queryForm.GetNRFDiscIPDomain(constvalue.SearchDataIPDoamin)
	if "" != ipDomain {
		ipDomainPath := getParamSearchPath(constvalue.NfTypeBSF, constvalue.SearchDataIPDoamin)
		ipDomainExpression := buildStringSearchParameter(ipDomainPath, ipDomain, constvalue.EQ)
		metaExpressionList = append(metaExpressionList, ipDomainExpression)

	}
	return metaExpressionList
}

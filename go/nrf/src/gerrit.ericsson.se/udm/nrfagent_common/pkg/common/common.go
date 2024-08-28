package common

import (
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/buger/jsonparser"

	"gerrit.ericsson.se/udm/common/pkg/cmproxy"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/msgbus"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

var (
	selfUUID string
	//ntfMsgbus *msgbus.MessageBus
	discMsgbus *msgbus.MessageBus
)

const (
	defaultTimeDeltaForSlave = 2 * time.Second
)

func GetSelfUUID() string {
	return selfUUID
}

func SetSelfUUID(id string) {
	selfUUID = id
}

func GetDiscMsgbus() *msgbus.MessageBus {
	return discMsgbus
}

func SetDiscMsgbus(msb *msgbus.MessageBus) {
	discMsgbus = msb
}

//CmNrfAgentLogHandler cm nrfAgent log handler
func CmNrfAgentLogHandler(Event, ConfigurationName, format string, RawData []byte) {
	log.Infof("cmNrfAgentLogHandler: %s, %s, %s", Event, format, string(RawData))
	if format != cmproxy.NtfFormatFull {
		log.Warnf("notification format:%s is not recommended", format)
		return
	}

	var rawData []byte
	var err error
	rawData, _, _, err = jsonparser.Get(RawData, "nrfagent-service-log")
	if err != nil {
		log.Errorf("Failed to run jsonparser.Get() nrfagent-service-log, %s", err.Error())
		return
	}

	switch {
	case cm.Opts.WorkMode == consts.AppWorkmodeREG:
		structs.UpdateNfServiceLog(rawData, "nrf_reg_agent")
	case cm.Opts.WorkMode == consts.AppWorkmodeNTF:
		structs.UpdateNfServiceLog(rawData, "nrf_ntf_agent")
	case cm.Opts.WorkMode == consts.AppWorkmodeDISC:
		structs.UpdateNfServiceLog(rawData, "nrf_disc_agent")
	}

	return
}

//CmGetTargetNfProfile is for get TargetNfProfile
func CmGetTargetNfProfile() ([]structs.TargetNf, bool) {
	var targetNfSet []structs.TargetNf
	targetNfProfiles := structs.GetTargetNfProfiles()
	if len(targetNfProfiles) == 0 {
		log.Errorf("CmGetTargetNfProfile: targetProfile is NULL.")
		return nil, false
	}

	for _, targetNfProfile := range targetNfProfiles {
		if targetNfProfile.RequesterNfType == "" ||
			targetNfProfile.TargetNfType == "" ||
			targetNfProfile.TargetServiceNames == nil ||
			len(targetNfProfile.TargetServiceNames) == 0 {
			log.Warnf("notification cmGetTargetNfProfile load failure RequesterNfType(%s), TargNfType(%s) ServiceNames(%+v)",
				targetNfProfile.RequesterNfType, targetNfProfile.TargetNfType, targetNfProfile.TargetServiceNames)
			continue
		}
		//		if len(targetNfProfile.TargetServiceNames) > 1 {
		//			log.Errorf("getTargetNF: Do not support multi Services, TargetNfType=%s", targetNfProfile.TargetNfType)
		//			continue
		//		}
		var targetNfInfo structs.TargetNf
		targetNfInfo.RequesterNfType = targetNfProfile.RequesterNfType
		targetNfInfo.TargetNfType = targetNfProfile.TargetNfType
		targetNfInfo.TargetServiceNames = targetNfProfile.TargetServiceNames
		targetNfSet = append(targetNfSet, targetNfInfo)
	}
	if len(targetNfSet) == 0 {
		return nil, false
	}

	return targetNfSet, true
}

//ConvertIpv6ToIpv4
func ConvertIpv6ToIpv4(v6Address string) string {
	if v6Addr := net.ParseIP(v6Address); v6Addr != nil {
		return v6Addr[12:16].String()
	}
	return ""
}

//ConvertIpv6ToIpv4InSearchResult convert Ipv6Address to Ipv4Address in SearchResult
func ConvertIpv6ToIpv4InSearchResult(content []byte, ipv6Supported bool) ([]byte, error) {
	if !ipv6Supported {
		return content, nil
	}

	var searchResult structs.SearchResult
	if err := json.Unmarshal(content, &searchResult); err != nil {
		log.Errorf("ConvertIpv6ToIpv4InSearchResult: invalid search result %+v", string(content))
		return nil, err
	}

	log.Infof("ConvertIpv6ToIpv4InSearchResult: convert Ipv6Address to Ipv4Address")
	for nfProfileIdx := range searchResult.NfInstances {
		searchResultNfProfile := &searchResult.NfInstances[nfProfileIdx]

		//LSV14 using it
		//		for len(searchResultNfProfile.Ipv6Address) > 0 {
		//			if v4Address := ConvertIpv6ToIpv4(searchResultNfProfile.Ipv6Address[0]); v4Address != "" {
		//				// convert ipv6address to ipv4address in NfProfile
		//				searchResultNfProfile.Ipv4Address = append(searchResultNfProfile.Ipv4Address, v4Address)
		//			} else {
		//				log.Errorf("ConvertIpv6ToIpv4InSearchResult: invalid Ipv6Address %s", searchResultNfProfile.Ipv6Address[0])
		//			}
		//			if len(searchResultNfProfile.Ipv6Address) > 1 {
		//				searchResultNfProfile.Ipv6Address = append(searchResultNfProfile.Ipv6Address[:0], searchResultNfProfile.Ipv6Address[1:]...)
		//			} else {
		//				searchResultNfProfile.Ipv6Address = nil
		//			}
		//		}
		for nfServiceIdx := range searchResultNfProfile.NfSrvList {
			nfService := &searchResultNfProfile.NfSrvList[nfServiceIdx]
			for ipEndPointIdx := range nfService.IPEndPoints {
				ipEndPoint := &nfService.IPEndPoints[ipEndPointIdx]
				if ipEndPoint.Ipv6Address != "" {
					if v4Address := ConvertIpv6ToIpv4(ipEndPoint.Ipv6Address); v4Address != "" {
						// convert ipv6address to ipv4address in NfProfile, and remove ipv6address field
						ipEndPoint.Ipv4Address = v4Address
						ipEndPoint.Ipv6Address = ""
					} else {
						log.Errorf("ConvertIpv6ToIpv4InSearchResult: invalid Ipv6Address %s", ipEndPoint.Ipv6Address)
					}
				}
			}
		}
		//Additional Code introduced by the mistake in NTF MSB schema
		for len(searchResultNfProfile.Ipv6Addresses) > 0 {
			if v4Address := ConvertIpv6ToIpv4(searchResultNfProfile.Ipv6Addresses[0]); v4Address != "" {
				// convert ipv6address to ipv4address in NfProfile
				searchResultNfProfile.Ipv4Addresses = append(searchResultNfProfile.Ipv4Addresses, v4Address)
			} else {
				log.Errorf("ConvertIpv6ToIpv4InSearchResult: invalid Ipv6Address %s", searchResultNfProfile.Ipv6Addresses[0])
			}
			if len(searchResultNfProfile.Ipv6Addresses) > 1 {
				searchResultNfProfile.Ipv6Addresses = append(searchResultNfProfile.Ipv6Addresses[:0], searchResultNfProfile.Ipv6Addresses[1:]...)
			} else {
				searchResultNfProfile.Ipv6Addresses = nil
			}
		}
	}
	return json.Marshal(searchResult)
}

func DispatchSubscrInfoToMessageBus(subscriptionInfo structs.SubscriptionInfo) bool {
	validityTime := subscriptionInfo.ValidityTime.Add(-defaultTimeDeltaForSlave)
	subscriptionInfo.ValidityTime = validityTime
	syncSubscrInfoMsg := structs.DiscDiscInnerMsg{
		EventType:       consts.EventTypeSyncSubscrInfo,
		AgentProducerID: GetSelfUUID(),
		SubscrInfo:      subscriptionInfo,
	}
	jsonBuf, err := json.Marshal(syncSubscrInfoMsg)
	if err != nil {
		log.Errorf("Failed to Marshal Disc message, Error: %s", err.Error())
		return false
	}

	innerTopicName := consts.MsgbusTopicNamePrefix + consts.DiscDiscInner
	return sendToMessageBus(innerTopicName, string(jsonBuf))
}

func IsRoamNotifcation(reqNfType string) bool {
	ok := strings.HasSuffix(reqNfType, consts.RoamSuffix)
	if !ok {
		log.Debugf("Regexp match roam for nfType:%s fail", reqNfType)
	}

	return ok
}

func GetReqNfTypeForRoam(reqNfTypeWithRoam string) string {
	return strings.Trim(reqNfTypeWithRoam, consts.RoamSuffix)
}

////////////private/////////////

func sendToMessageBus(topicName string, messageData string) bool {
	log.Debugf("Send message to message bus :Topic: %s message: %s", topicName, messageData)
	if discMsgbus := GetDiscMsgbus(); discMsgbus != nil {
		err := discMsgbus.SendMessage(topicName, messageData)
		if err != nil {
			log.Errorf("%s:Failed to send message to message bus, %s", topicName, err.Error())
			return false
		}
	} else {
		log.Errorf("message bus is not initialized")
		return false
	}
	return true
}

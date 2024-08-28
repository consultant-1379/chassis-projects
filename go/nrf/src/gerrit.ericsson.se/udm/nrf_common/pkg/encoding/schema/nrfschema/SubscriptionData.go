package nrfschema

import (
	"com/dbproxy/nfmessage/subscription"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

//ValidateNotificationURI : nfStatusNotificationUri must be valid
func (s *TSubscriptionData) ValidateNotificationURI() *problemdetails.ProblemDetails {
	notificationURI, err := url.ParseRequestURI(s.NfStatusNotificationUri)
	if err != nil || (notificationURI.Scheme != "http" && notificationURI.Scheme != "https") {
		return &problemdetails.ProblemDetails{
			Title: "not a valid subscriptionData",
			InvalidParams: []*problemdetails.InvalidParam{
				&problemdetails.InvalidParam{
					Param:  constvalue.SubDataNotificationUri,
					Reason: fmt.Sprintf("invalid nfStatusNotificationUri %s", s.NfStatusNotificationUri),
				},
			},
		}
	}

	return nil
}

/*ValidateSubscrCond : SubscrCond shall be oneOf the following objects:
  - NfInstanceIdCond
  - NfTypeCond
  - ServiceNameCond
  - AmfCond
  - GuamiListCond
  - NetworkSliceCond
  - NfGroupCond
*/
func (s *TSubscriptionData) ValidateSubscrCond() *problemdetails.ProblemDetails {

	if s.SubscrCond != nil {
		invalidReason := s.SubscrCond.Validate()
		if invalidReason != "" {
			return &problemdetails.ProblemDetails{
				Title: "not a valid subscriptionData",
				InvalidParams: []*problemdetails.InvalidParam{
					&problemdetails.InvalidParam{
						Param:  constvalue.SubDataSubscrCond,
						Reason: invalidReason,
					},
				},
			}
		}
	} else {
		if s.ReqNfType == "" {
			return &problemdetails.ProblemDetails{
				Title: "not a valid subscriptionData",
				InvalidParams: []*problemdetails.InvalidParam{
					&problemdetails.InvalidParam{
						Param:  constvalue.SubDataSubscrCond,
						Reason: constvalue.StatusSubscribeRule7,
					},
				},
			}
		}

		match := false
		for _, allowedSubscriptionAllNF := range cm.NrfPolicy.ManagementService.Subscription.AllowedSubscriptionAllNFs {
			//allowedSubscriptionAllNF.AllowedNfType is mandatory
			if allowedSubscriptionAllNF.AllowedNfType != "" && allowedSubscriptionAllNF.AllowedNfType == s.ReqNfType {
				match = true
				if allowedSubscriptionAllNF.AllowedNfDomains != "" {
					domainPattern := strings.Replace(allowedSubscriptionAllNF.AllowedNfDomains, `\\`, `\`, -1)
					match, _ = regexp.MatchString(domainPattern, s.ReqNfFqdn)
				}
			}
			if match {
				break
			}
		}
		if !match {
			invalidReason := fmt.Sprintf(constvalue.StatusSubscribeRule8, s.ReqNfType)
			return &problemdetails.ProblemDetails{
				Title: "not a valid subscriptionData",
				InvalidParams: []*problemdetails.InvalidParam{
					&problemdetails.InvalidParam{
						Param:  constvalue.SubDataSubscrCond,
						Reason: invalidReason,
					},
				},
			}
		}

	}

	return nil
}

// ValidateValidityTime : validity time cannot be before Now
func (s *TSubscriptionData) ValidateValidityTime() *problemdetails.ProblemDetails {
	if s.ValidityTime != "" {
		t, err := time.Parse(time.RFC3339, s.ValidityTime)
		if err != nil || t.Unix() <= time.Now().Unix() {
			return &problemdetails.ProblemDetails{
				Title: "not a valid subscriptionData",
				InvalidParams: []*problemdetails.InvalidParam{
					&problemdetails.InvalidParam{
						Param:  constvalue.SubDataValidityTime,
						Reason: fmt.Sprintf(constvalue.StatusSubscribeRule5, time.Now().Format(time.RFC3339)),
					},
				},
			}
		}
	}

	return nil
}

//ValidateNotifCondition : attributes "monitoredAttributes" and "unmonitoredAttributes" shall not be included simultaneously
func (s *TSubscriptionData) ValidateNotifCondition() *problemdetails.ProblemDetails {
	if s.NotifCondition != nil {
		if !s.NotifCondition.IsValid() {
			return &problemdetails.ProblemDetails{
				Title: "not a valid subscriptionData",
				InvalidParams: []*problemdetails.InvalidParam{
					&problemdetails.InvalidParam{
						Param:  constvalue.SubDataNotifCondition,
						Reason: constvalue.StatusSubscribeRule6,
					},
				},
			}

		}
	}

	return nil
}

// Validate validate subscriptionData
func (s *TSubscriptionData) Validate() *problemdetails.ProblemDetails {
	problemDetails := s.ValidateNotificationURI()
	if problemDetails != nil {
		return problemDetails
	}

	problemDetails = s.ValidateSubscrCond()
	if problemDetails != nil {
		return problemDetails
	}

	problemDetails = s.ValidateValidityTime()
	if problemDetails != nil {
		return problemDetails
	}

	problemDetails = s.ValidateNotifCondition()
	if problemDetails != nil {
		return problemDetails
	}

	return nil
}

// IsLocalPlmn judge whether the subscribe request is local
func (s *TSubscriptionData) IsLocalPlmn() bool {
	if s.PlmnId == nil {
		return true
	}

	isLocalPlmn := false
	for _, plmn := range cm.NfProfile.PlmnID {
		if s.PlmnId.GetPlmnID() == plmn.GetPlmnID() {
			isLocalPlmn = true
			break
		}
	}

	return isLocalPlmn
}

// ConstructSubscriptionIndex construct subscriptionIndex
func (s *TSubscriptionData) ConstructSubscriptionIndex() *subscription.SubscriptionPutIndex {
	subscriptionIndex := &subscription.SubscriptionPutIndex{
		NfStatusNotificationUri: s.NfStatusNotificationUri,
	}

	if s.SubscrCond == nil {
		subscriptionIndex.NoCond = constvalue.NoSubscrCond
		subscriptionIndex.NfInstanceId = constvalue.Wildcard
		subscriptionIndex.NfType = constvalue.Wildcard
		subscriptionIndex.ServiceName = constvalue.Wildcard
		subscriptionIndex.AmfCond = &subscription.SubKeyStruct{
			SubKey1: constvalue.Wildcard,
			SubKey2: constvalue.Wildcard,
		}
		subscriptionIndex.GuamiList = []*subscription.SubKeyStruct{
			&subscription.SubKeyStruct{
				SubKey1: constvalue.Wildcard,
				SubKey2: constvalue.Wildcard,
			},
		}
		subscriptionIndex.SnssaiList = []*subscription.SubKeyStruct{
			&subscription.SubKeyStruct{
				SubKey1: constvalue.Wildcard,
				SubKey2: constvalue.Wildcard,
			},
		}
		subscriptionIndex.NsiList = []string{constvalue.Wildcard}
		subscriptionIndex.NfGroupCond = &subscription.SubKeyStruct{
			SubKey1: constvalue.Wildcard,
			SubKey2: constvalue.Wildcard,
		}
	} else {
		subscriptionIndex.NoCond = constvalue.Wildcard

		if s.SubscrCond.NfInstanceID != "" {
			subscriptionIndex.NfInstanceId = s.SubscrCond.NfInstanceID
		} else {
			subscriptionIndex.NfInstanceId = constvalue.Wildcard
		}

		if s.SubscrCond.ServiceName != "" {
			subscriptionIndex.ServiceName = s.SubscrCond.ServiceName
		} else {
			subscriptionIndex.ServiceName = constvalue.Wildcard
		}

		if s.SubscrCond.NfGroupID != "" {
			subscriptionIndex.NfType = constvalue.Wildcard
			subscriptionIndex.NfGroupCond = &subscription.SubKeyStruct{
				SubKey1: s.SubscrCond.NfGroupID,
				SubKey2: s.SubscrCond.NfType,
			}
		} else {
			subscriptionIndex.NfGroupCond = &subscription.SubKeyStruct{
				SubKey1: constvalue.Wildcard,
				SubKey2: constvalue.Wildcard,
			}
			if s.SubscrCond.NfType != "" {
				subscriptionIndex.NfType = s.SubscrCond.NfType
			} else {
				subscriptionIndex.NfType = constvalue.Wildcard
			}
		}

		if s.SubscrCond.AmfSetID != "" && s.SubscrCond.AmfRegionID != "" {
			subscriptionIndex.AmfCond = &subscription.SubKeyStruct{
				SubKey1: s.SubscrCond.AmfSetID,
				SubKey2: s.SubscrCond.AmfRegionID,
			}
		} else if s.SubscrCond.AmfSetID != "" && s.SubscrCond.AmfRegionID == "" {
			subscriptionIndex.AmfCond = &subscription.SubKeyStruct{
				SubKey1: s.SubscrCond.AmfSetID,
				SubKey2: constvalue.Wildcard,
			}
		} else if s.SubscrCond.AmfSetID == "" && s.SubscrCond.AmfRegionID != "" {
			subscriptionIndex.AmfCond = &subscription.SubKeyStruct{
				SubKey1: constvalue.Wildcard,
				SubKey2: s.SubscrCond.AmfRegionID,
			}
		} else if s.SubscrCond.AmfSetID == "" && s.SubscrCond.AmfRegionID == "" {
			subscriptionIndex.AmfCond = &subscription.SubKeyStruct{
				SubKey1: constvalue.Wildcard,
				SubKey2: constvalue.Wildcard,
			}
		}

		var guamiArray []*subscription.SubKeyStruct
		if s.SubscrCond.GuamiList == nil || len(s.SubscrCond.GuamiList) == 0 {
			guami := &subscription.SubKeyStruct{
				SubKey1: constvalue.Wildcard,
				SubKey2: constvalue.Wildcard,
			}
			guamiArray = append(guamiArray, guami)
		} else {
			for _, item := range s.SubscrCond.GuamiList {
				guami := item.GenerateGrpcKey()
				guamiArray = append(guamiArray, guami)
			}
		}
		subscriptionIndex.GuamiList = guamiArray

		var snssaiArray []*subscription.SubKeyStruct
		if s.SubscrCond.SnssaiList == nil || len(s.SubscrCond.SnssaiList) == 0 {
			snssai := &subscription.SubKeyStruct{
				SubKey1: constvalue.Wildcard,
				SubKey2: constvalue.Wildcard,
			}
			snssaiArray = append(snssaiArray, snssai)
		} else {
			for _, item := range s.SubscrCond.SnssaiList {
				snssai := item.GenerateGrpcPutKey()
				snssaiArray = append(snssaiArray, snssai)
			}
		}
		subscriptionIndex.SnssaiList = snssaiArray

		var nsiList []string
		if s.SubscrCond.NsiList != nil {
			for _, item := range s.SubscrCond.NsiList {
				if item != "" {
					nsiList = append(nsiList, item)
				}
			}
		}

		if nsiList == nil {
			nsiList = append(nsiList, constvalue.Wildcard)
		}

		subscriptionIndex.NsiList = nsiList
	}

	subscriptionIndex.ValidityTime = uint64(s.GenerateExpiredTimeInMilSec())

	return subscriptionIndex
}

// GenerateValidatyDateTime returns the ValidatyDateTime of subscriptionData
func (s *TSubscriptionData) GenerateValidatyDateTime() string {

	cmValidateTimeInSeconds := int64(cm.ValidityPeriodOfSubscription) + time.Now().Unix()
	if s.ValidityTime == "" {
		return time.Unix(cmValidateTimeInSeconds, 0).Format(time.RFC3339)
	}

	t, err := time.Parse(time.RFC3339, s.ValidityTime)
	if err != nil {
		return time.Unix(cmValidateTimeInSeconds, 0).Format(time.RFC3339)
	}
	if t.Unix() <= cmValidateTimeInSeconds {
		return s.ValidityTime
	}

	return time.Unix(cmValidateTimeInSeconds, 0).Format(time.RFC3339)
}

// GenerateExpiredTimeInMilSec returns the expiredTime of subscriptionData in millisecond
func (s *TSubscriptionData) GenerateExpiredTimeInMilSec() int64 {
	validityDateTime := s.GenerateValidatyDateTime()
	s.ValidityTime = validityDateTime
	timeInSecond, _ := time.Parse(time.RFC3339, validityDateTime)
	return timeInSecond.Unix() * 1000
}

package subscribe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"github.com/buger/jsonparser"
)

const (
	httpContentTypeJSON        = "application/json"
	httpHeaderJSONPatchJSON    = "application/json-patch+json"
	httpContentTypeProblemJSON = "application/problem+json"
	httpResponseFormat         = `{"title": "%s"}`
)

const (
	defaultValidityTime      = 876576 * time.Hour
	defaultTimeDelta         = 5 * time.Second
	defaultTimeDeltaForSlave = 2 * time.Second
)

/*
var (
	cacheManager *cache.CacheManager
)

func init() {
	cacheManager = cache.Instance()
}
*/

func UnsubscribeByNfType(targetNfs []structs.TargetNf) {
	/*
		if worker.IsKeepCacheMode() {
			log.Info("keep cache mode, not send message to NRF.")
			return
		}
	*/
	if len(targetNfs) == 0 {
		log.Warn("No targetNfs for nf, please check configmap status")
		return
	}

	/*
		targetNfs, ok := cacheManager.GetTargetNfs(requesterNfType)
		if !ok {
			log.Errorf("Failed to get targetNfProfiles for nfType[%s], please check configmap status", requesterNfType)
			return
		}
	*/

	//for _, targetNf := range targetNfs {
	//	unsubscribeSubscription(targetNf.RequesterNfType, targetNf.TargetNfType)
	//}
}

/*
func UnsubscribeByNfInstanceID(requesterNfType, targetNfType, nfInstanceID string) {
	subscriptionID, ok := cacheManager.GetNfProfileSubscriptionID(requesterNfType, targetNfType, nfInstanceID)
	if !ok {
		log.Warnf("No such subscription for nfProfile:%s, will skip do unsubscribe by nfInstanceID", nfInstanceID)
		return
	}
	log.Infof("nfProfile:%s subscriptionID:%s", nfInstanceID, subscriptionID)

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON

	if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
		resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", subscriptionID, hdr, bytes.NewBuffer([]byte("")))
		if err != nil {
			log.Errorf("failed to send DELETE subscription request to NRF, %s", err.Error())
		} else {
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				log.Errorf("Failed to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
			} else {
				log.Infof("Success to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
			}
		}
	} else {
		log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
	}

	cacheManager.DelSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
	cacheManager.DelSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
}
*/
////////////////private//////////////

/*
func unsubscribeSubscription(requesterNfType, targetNfType string) {
	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON

	ok, subscriptionIDURLs := cache.Instance().GetSubscriptionIDs(requesterNfType, targetNfType)
	for ok && len(subscriptionIDURLs) > 0 {
		subscriptionID := subscriptionIDURLs[0]
		log.Infof("subscriptionID of %s: %+v", requesterNfType, subscriptionID)

		if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
			resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", subscriptionID, hdr, bytes.NewBuffer([]byte("")))
			if err != nil {
				log.Errorf("failed to send DELETE subscription request to NRF, %s", err.Error())
			} else {
				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
					log.Errorf("Failed to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
				} else {
					log.Infof("Success to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
				}
			}
		} else {
			log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
		}

		cacheManager.DelSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
		cacheManager.DelSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
		cacheManager.UpdateSubscriptionStorage()

		ok, subscriptionIDURLs = cacheManager.GetSubscriptionIDs(requesterNfType, targetNfType)
	}

	ok, roamsubscriptionIDURLs := cacheManager.GetRoamingSubscriptionIDs(requesterNfType, targetNfType)
	for ok && len(roamsubscriptionIDURLs) > 0 {
		roamsubscriptionID := roamsubscriptionIDURLs[0]
		log.Infof("roaming subscriptionID of %s: %+v", requesterNfType, roamsubscriptionID)

		if election.IsActiveLeader(strconv.Itoa(cm.Opts.PortHTTPWithoutTLS), consts.DiscoveryAgentReadinessProbe) {
			resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", roamsubscriptionID, hdr, bytes.NewBuffer([]byte("")))
			if err != nil {
				log.Errorf("failed to send DELETE roaming subscription request to NRF, %s", err.Error())
			} else {
				if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
					log.Errorf("Failed to DELETE roaming subscription(%s), StatusCode(%d)", roamsubscriptionID, resp.StatusCode)
				} else {
					log.Infof("Success to DELETE roaming subscription(%s), StatusCode(%d)", roamsubscriptionID, resp.StatusCode)
				}
			}
		} else {
			log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
		}

		cacheManager.DelRoamingSubscriptionInfo(requesterNfType, targetNfType, roamsubscriptionID)
		cacheManager.DelRoamingSubscriptionMonitor(requesterNfType, targetNfType, roamsubscriptionID)

		ok, roamsubscriptionIDURLs = cacheManager.GetRoamingSubscriptionIDs(requesterNfType, targetNfType)
	}
}
*/

//func SubscribeExecutor(oneSubsData *structs.OneSubscriptionData) (string, time.Time, error) {
func SubscribePostExecutor(subscriptionData *structs.SubscriptionData) (*httpclient.HttpRespData, error) {
	/*
		if IsKeepCacheMode() {
			log.Info("NRF Disc Agent is Master, Woke Mode is KeepCache Mode, So no need to do subscription")
			return "", time.Time{}, nil
		}

		fqdn, exists := cache.Instance().GetRequesterFqdn(oneSubsData.RequesterNfType)
		if !exists {
			log.Warnf("%s nf instance was deregistered", oneSubsData.RequesterNfType)
			return "", time.Time{}, nil
		}
	*/
	//var subscribeData []byte
	//subscribeData = util.BuildSubscriptionPostData(oneSubsData, fqdn)

	if subscriptionData == nil {
		errInfo := fmt.Sprintf("%s", "POST subscriptionData is nil")
		return nil, fmt.Errorf("%s", errInfo)
	}

	subscriptionDataRaw, err := json.Marshal(*subscriptionData)
	if err != nil {
		log.Errorf("Marshal subscriptionData fail, err:%s", err.Error())
		return nil, err
	}

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON
	resp, err := client.HTTPDoToNrfMgmt("h2", "POST", "subscriptions", hdr, bytes.NewBuffer(subscriptionDataRaw))
	if err != nil {
		log.Errorf("Failed to send subscription request to NRF, %s", err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		log.Errorf("Subscribe from NRF by POST method fail, statusCode:%d, body:%s", resp.StatusCode, string(resp.Body))
		return nil, fmt.Errorf("Subscribe from NRF for %s fail", subscriptionData.ReqNfType)
	}

	return resp, nil
	/*
		subscriptionID, timeStamp, err := subscribeResponseParser(resp)
		if err != nil {
			return "", time.Time{}, err
		}

		return subscriptionID, timeStamp, nil
	*/
}

func SubscribePatchExecutor(subscriptionID string, patchItems []structs.PatchItem) (*httpclient.HttpRespData, error) {
	if len(patchItems) == 0 {
		errInfo := fmt.Sprintf("Subscribe PATCH data is empty for subscriptionID:%s", subscriptionID)
		return nil, fmt.Errorf("%s", errInfo)
	}

	patchData, err := json.Marshal(patchItems)
	if err != nil {
		log.Errorf("Marshal patch items failure, err:%s", err.Error())
		return nil, err
	}

	hdr := make(map[string]string)
	hdr["Content-Type"] = httpHeaderJSONPatchJSON
	var resp *httpclient.HttpRespData

	resp, err = client.HTTPDoToNrfMgmt("h2", "PATCH", subscriptionID, hdr, bytes.NewBuffer(patchData))
	if err != nil {
		log.Errorf("Send subscribe request to NRF fail,err:%s", err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		errInfo := fmt.Sprintf("Prolong subscription[%s] by PATCH fail, code[%d], body[%s]", subscriptionID, resp.StatusCode, string(resp.Body))
		return nil, fmt.Errorf("%s", errInfo)
	}
	log.Infof("Subscribe PATCH response body from NRF:[%s]", string(resp.Body))

	return resp, nil
	/*
		var validityTime time.Time
		if resp.StatusCode == http.StatusNoContent {
			now := time.Now()
			prolongValue := time.Duration(targetNf.SubscriptionValidityTime) * time.Second
			nextTimeStamp := now.Add(prolongValue)
			validityTime = nextTimeStamp
		} else if resp.StatusCode == http.StatusOK {
			timestamp, err := jsonparser.GetString(resp.Body, "validityTime")
			if err != nil {
				log.Warnf("Get validityTime in subscription boby fail, err=%s", err.Error())
				return nil, nil
			}
			validityTime, err = time.Parse(time.RFC3339, timestamp)
			if err != nil {
				log.Warnf("Parse timestamp[%s] by RFC3339 fail, err=%s", timestamp, err.Error())
				return nil, nil
			}
		}

		validityTime = validityTime.Add(-defaultTimeDelta)
		return nil, &validityTime
	*/
}

func UnSubscribeExecutor(subscriptionID string) error {
	hdr := make(map[string]string)
	hdr["Content-Type"] = httpContentTypeJSON

	log.Infof("Do unsubscription by subscriptionID:%s", subscriptionID)

	resp, err := client.HTTPDoToNrfMgmt("h2", "DELETE", subscriptionID, hdr, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Errorf("failed to send DELETE subscription request to NRF, %s", err.Error())
		return err
	} else {
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			log.Errorf("Failed to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
			return fmt.Errorf("Failed to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
		} else {
			log.Infof("Success to DELETE subscription(%s), StatusCode(%d)", subscriptionID, resp.StatusCode)
			return nil
		}
	}
	//} else {
	//	log.Infof("Slaver discovery agent has no need to send unsubscription request to NRF")
	//}

	//cacheManager.DelSubscriptionInfo(requesterNfType, targetNfType, subscriptionID)
	//cacheManager.DelSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
	//cacheManager.UpdateSubscriptionStorage()

	//ok, subscriptionIDURLs = cacheManager.GetSubscriptionIDs(requesterNfType, targetNfType)
}

func SubscribePostRespParser(resp *httpclient.HttpRespData) (string, time.Time, error) {
	var errInfo string
	if resp == nil {
		errInfo = "subscribe POST response is nil"
		return "", time.Time{}, fmt.Errorf("%s", errInfo)
	}

	location := resp.Location
	log.Debugf("subscribe response location:%s", location)

	subURL := strings.Split(location, "//")
	if len(subURL) < 2 {
		errInfo = "subscriptionID in location is error"
		return "", time.Time{}, fmt.Errorf("%s", errInfo)
	}
	subURLSuffix := strings.Split(subURL[1], "/")
	if len(subURLSuffix) < 5 {
		errInfo = "subscriptionID in location is error"
		return "", time.Time{}, fmt.Errorf("%s", errInfo)
	}
	subscriptionID := subURLSuffix[3] + "/" + subURLSuffix[4]

	defaultTime := time.Now().Add(defaultValidityTime)
	vt, err := jsonparser.GetString(resp.Body, "validityTime")
	if err != nil {
		log.Warnf("No validityTime in subscription boby")
		return subscriptionID, defaultTime, nil
	}
	validityTime, err := time.Parse(time.RFC3339, vt)
	if err != nil {
		log.Warnf("validityTime %s is invalid when parse by %s", vt, time.RFC3339)
		return subscriptionID, defaultTime, nil
	}

	return subscriptionID, validityTime.Add(-defaultTimeDelta), nil
}

func SubscribePatchRespParser(resp *httpclient.HttpRespData, expect time.Time) (time.Time, error) {
	var errInfo string
	if resp == nil {
		errInfo = "subscribe PATCH response is nil"
		return time.Time{}, fmt.Errorf("%s", errInfo)
	}

	var validityTime time.Time
	if resp.StatusCode == http.StatusNoContent {
		validityTime = expect
	} else if resp.StatusCode == http.StatusOK {
		timestamp, err := jsonparser.GetString(resp.Body, "validityTime")
		if err != nil {
			return time.Time{}, err
		}
		validityTime, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			log.Warnf("Parse timestamp[%s] by RFC3339 fail, err=%s", timestamp, err.Error())
			return time.Time{}, err
		}
	}

	validityTime = validityTime.Add(-defaultTimeDelta)

	return validityTime, nil
}

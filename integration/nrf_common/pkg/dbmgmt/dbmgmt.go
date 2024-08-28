package dbmgmt

import (
	"com/dbproxy"
	"com/dbproxy/nfmessage/cachenfprofile"
	"com/dbproxy/nfmessage/gpsiprefixprofile"
	"com/dbproxy/nfmessage/gpsiprofile"
	"com/dbproxy/nfmessage/groupprofile"
	"com/dbproxy/nfmessage/imsiprefixprofile"
	"com/dbproxy/nfmessage/nfprofile"
	"com/dbproxy/nfmessage/nrfaddress"
	"com/dbproxy/nfmessage/nrfprofile"
	"com/dbproxy/nfmessage/subscription"
	"fmt"
	"math/rand"
	"time"

	"io"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//PlmnNRFAddress plmnNRF address info
type PlmnNRFAddress struct {
	NrfAddressID string       `json:"nrfAddressId"`
	ServingPLMN  PlmnID       `json:"servingPlmn"`
	NrfAddresses []NRFAddress `json:"nrfAddresses"`
}

//PlmnID plmn info
type PlmnID struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"mnc"`
}

//NRFAddress nrfaddress info
type NRFAddress struct {
	Scheme string `json:"scheme"`
	Fqdn   string `json:"fqdn"`
	Port   int    `json:"port"`
}

//PlmnNRFAddressProv plmnNRF address provision info
type PlmnNRFAddressProv struct {
	NrfAddressID string           `json:"nrfAddressId"`
	ServingPLMN  PlmnID           `json:"servingPlmn"`
	NrfAddresses []NRFAddressProv `json:"nrfAddresses"`
}

//NRFAddressProv nrfaddress provision info
type NRFAddressProv struct {
	Scheme  string  `json:"scheme"`
	Address NRFAddr `json:"address"`
	Port    int     `json:"port"`
}

//NRFAddr NRFAddr info
type NRFAddr struct {
	Fqdn        string `json:"fqdn,omitempty"`
	Ipv4Address string `json:"ipv4Address,omitempty"`
	Ipv6Address string `json:"ipv6Address,omitempty"`
}

var (
	dbConn *grpc.ClientConn
	dbFlag bool
	//dbclient           dbproxy.NFDataManagementServiceClient
	dbclientPool       []dbproxy.NFDataManagementServiceClient
	poolLength         int
	ctxTimeoutInSecond time.Duration
)

// InitDB is to init db proxy
func InitDB(hostColonPort string) {
	// Set up  connections to the server.
	var err error
	dbclientPool = make([]dbproxy.NFDataManagementServiceClient, internalconf.DbproxyConnectionNum)
	connectionNumIndex := 0
	for i := 0; i < 10000; i++ {
		dbConn, err = grpc.Dial(hostColonPort, grpc.WithReadBufferSize(131072), grpc.WithWriteBufferSize(1048576), grpc.WithInsecure())
		if err != nil {
			log.Errorf("did not connect: %v", err)
		} else {
			dbclientPool[connectionNumIndex] = dbproxy.NewNFDataManagementServiceClient(dbConn)
			connectionNumIndex++
			if connectionNumIndex == internalconf.DbproxyConnectionNum {
				dbFlag = true
				break

			}
		}
	}
	poolLength = len(dbclientPool)
	ctxTimeoutInSecond = time.Duration(internalconf.DbproxyGrpcCtxTimeout)
}

// Close the connection to the gRPC server.
func Close() {
	if err := dbConn.Close(); err != nil {
		log.Errorf("close connect error: %v", err)
	}
}

// Put is for nfprofile,subscription,nfscription registering
func Put(PutReqData *dbproxy.PutRequest) (*dbproxy.PutResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeoutInSecond)
	defer cancel()

	putRespone := &dbproxy.PutResponse{}

	putRequest := &dbproxy.NFRequest_PutRequest{
		PutRequest: PutReqData,
	}

	reqData := &dbproxy.NFRequest{
		Data: putRequest,
	}

	request := &dbproxy.NFMessage_Request{
		Request: reqData,
	}

	nfMsg := &dbproxy.NFMessage{
		Data: request,
	}

	dbclient := dbclientPool[rand.Intn(poolLength)]
	ret, err := dbclient.Execute(ctx, nfMsg)
	if err != nil {
		log.Errorf("calling put fail : %v", err)
		return putRespone, err
	}

	respone := &dbproxy.NFResponse{}
	switch ret.GetData().(type) {
	case *dbproxy.NFMessage_Response:
		respone = ret.GetResponse()

	case *dbproxy.NFMessage_ProtocolError:
		retCode := ret.GetProtocolError().GetCode()
		return putRespone, fmt.Errorf("protocol error,retcode=%d", retCode)
	default:
		log.Errorf("db-proxy return invalid messages type, expect response type for put")
		return putRespone, fmt.Errorf("db-proxy return invalid messages type")

	}

	switch respone.GetData().(type) {
	case *dbproxy.NFResponse_PutResponse:
		putRespone = respone.GetPutResponse()
	default:
		log.Errorf("db-proxy return invalid respone type,expect putResponse type")
		return putRespone, fmt.Errorf("db-proxy return invalid respone type")

	}

	return putRespone, nil
}

// Get is for nfprofile,subscription,nfscription getting
func Get(getReqData *dbproxy.GetRequest) (*dbproxy.GetResponse, error) {
	return GetWithTimer(getReqData, ctxTimeoutInSecond)
}

// GetWithTimer is for nfprofile,subscription,nfscription getting with special timer
func GetWithTimer(getReqData *dbproxy.GetRequest, timer time.Duration) (*dbproxy.GetResponse, error) {

	getRespone := &dbproxy.GetResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timer)
	defer cancel()

	getRequest := &dbproxy.NFRequest_GetRequest{
		GetRequest: getReqData,
	}

	reqData := &dbproxy.NFRequest{
		Data: getRequest,
	}

	request := &dbproxy.NFMessage_Request{
		Request: reqData,
	}
	nfMsg := &dbproxy.NFMessage{
		Data: request,
	}

	dbclient := dbclientPool[rand.Intn(poolLength)]
	ret, err := dbclient.Execute(ctx, nfMsg)
	if err != nil {
		log.Errorf("calling get fail : %v", err)
		return getRespone, err
	}

	respone := &dbproxy.NFResponse{}
	switch ret.GetData().(type) {
	case *dbproxy.NFMessage_Response:
		respone = ret.GetResponse()

	case *dbproxy.NFMessage_ProtocolError:
		retCode := ret.GetProtocolError().GetCode()
		return getRespone, fmt.Errorf("protocol error,retcode=%d", retCode)
	default:
		log.Errorf("db-proxy return invalid messages type, expect response type for get")
		return getRespone, fmt.Errorf("db-proxy return invalid messages type")
	}

	switch respone.GetData().(type) {
	case *dbproxy.NFResponse_GetResponse:
		getRespone = respone.GetGetResponse()
	default:
		log.Errorf("db-proxy return invalid respone type, expect getRespone type")
		return getRespone, fmt.Errorf("db-proxy return invalid respone type")
	}

	return getRespone, nil
}

//TransferParameter is transfer cm or other configuration to dbproxy
func TransferParameter(paraReq *dbproxy.ParaRequest) (*dbproxy.ParaResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeoutInSecond)
	defer cancel()
	dbclient := dbclientPool[rand.Intn(poolLength)]
	paraResp, err := dbclient.TransferParameter(ctx, paraReq)
	if err != nil {
		log.Errorf("calling get fail : %v", err)
		return paraResp, err
	}

	return paraResp, nil
}

//QueryWithKey is to use key to search in profile
func QueryWithKey(queryReq *dbproxy.QueryRequest) (*dbproxy.QueryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeoutInSecond)
	defer cancel()
	dbclient := dbclientPool[rand.Intn(poolLength)]
	var startTime int64
	if internalconf.EnableTimeStatistics {
		startTime = time.Now().UnixNano() / 1000000
		queryReq.TraceEnabled = true
	}

	stream, err := dbclient.QueryByKey(ctx, queryReq)
	if err != nil {
		log.Errorf("calling get fail : %v", err)
		return nil, err
	}

	var responseValue []string
	var code uint32
	var traceInfo *dbproxy.TraceInfo
	for {
		c, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Debug("queryWithKey receive streaming message end")
				break
			}
			log.Errorf("queryWithKey fail to receive streaming message, err=%v", err)
			break
		}
		responseValue = append(responseValue, c.GetValue()...)
		code = c.GetCode()
		traceInfo = c.GetTraceInfo()
	}
	queryResponse := &dbproxy.QueryResponse{Value: responseValue, Code: code}
	if internalconf.EnableTimeStatistics {
		if traceInfo != nil {
			enterGrpcTime := traceInfo.ArrivalTime
			leaveGrpcTime := traceInfo.DepartureTime
			endTime := time.Now().UnixNano() / 1000000
			if queryReq.RegionName == configmap.DBGpsiprefixProfileRegionName || queryReq.RegionName == configmap.DBImsiprefixProfileRegionName {
				DBLatency.groupIDRequestChannel <- Latency{startGrpcTime: startTime, endGrpcTime: endTime, enterDBProxyTime: enterGrpcTime, leaveDBProxyTime: leaveGrpcTime}
			} else if queryReq.RegionName == configmap.DBNfprofileRegionName {
				DBLatency.instIDRequestChannel <- Latency{startGrpcTime: startTime, endGrpcTime: endTime, enterDBProxyTime: enterGrpcTime, leaveDBProxyTime: leaveGrpcTime}
			}
		}
	}

	return queryResponse, nil
}

//QueryWithFilter is to use filter OQL to get profile from db
func QueryWithFilter(queryReq *dbproxy.QueryRequest) (*dbproxy.QueryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeoutInSecond)
	defer cancel()
	dbclient := dbclientPool[rand.Intn(poolLength)]
	var startTime int64
	if internalconf.EnableTimeStatistics {
		startTime = time.Now().UnixNano() / 1000000
		queryReq.TraceEnabled = true
	}
	stream, err := dbclient.QueryByFilter(ctx, queryReq)
	if err != nil {
		log.Errorf("calling get fail : %v", err)
		return nil, err
	}

	var responseValue []string
	var code uint32
	var traceInfo *dbproxy.TraceInfo
	for {
		c, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("queryWithFilter fail to receive streaming message, err=%v", err)
			break
		}
		responseValue = append(responseValue, c.GetValue()...)
		code = c.GetCode()
		traceInfo = c.GetTraceInfo()
	}
	queryResponse := &dbproxy.QueryResponse{Value: responseValue, Code: code}
	if internalconf.EnableTimeStatistics {
		if traceInfo != nil {
			enterGrpcTime := traceInfo.ArrivalTime
			leaveGrpcTime := traceInfo.DepartureTime
			endTime := time.Now().UnixNano() / 1000000
			DBLatency.nfProfileFilterChannel <- Latency{startGrpcTime: startTime, endGrpcTime: endTime, enterDBProxyTime: enterGrpcTime, leaveDBProxyTime: leaveGrpcTime}
		}
	}

	return queryResponse, nil
}

// Delete is for nfprofile,subscription,nfscription deregistering
func Delete(delReqData *dbproxy.DelRequest) (*dbproxy.DelResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeoutInSecond)
	defer cancel()

	delRespone := &dbproxy.DelResponse{}

	delRequest := &dbproxy.NFRequest_DelRequest{
		DelRequest: delReqData,
	}

	reqData := &dbproxy.NFRequest{
		Data: delRequest,
	}

	request := &dbproxy.NFMessage_Request{
		Request: reqData,
	}

	nfMsg := &dbproxy.NFMessage{
		Data: request,
	}

	dbclient := dbclientPool[rand.Intn(poolLength)]
	ret, err := dbclient.Execute(ctx, nfMsg)
	if err != nil {
		log.Errorf("calling delete fail : %v", err)
		return delRespone, err
	}

	respone := &dbproxy.NFResponse{}
	switch ret.GetData().(type) {
	case *dbproxy.NFMessage_Response:
		respone = ret.GetResponse()

	case *dbproxy.NFMessage_ProtocolError:
		retCode := ret.GetProtocolError().GetCode()
		return delRespone, fmt.Errorf("protocol error,retcode=%d", retCode)
	default:
		log.Errorf("db-proxy return invalid messages type,expect respone type for delete")
		return delRespone, fmt.Errorf("db-proxy return invalid messages type")
	}

	switch respone.GetData().(type) {
	case *dbproxy.NFResponse_DelResponse:
		delRespone = respone.GetDelResponse()
	default:
		log.Errorf("db-proxy return invalid respone type, expect delResponse type")
		return delRespone, fmt.Errorf("db-proxy return invalid respone type")
	}

	return delRespone, nil
}

// PutNFProfile is for NF profile register
func PutNFProfile(putRequest *nfprofile.NFProfilePutRequest) (*nfprofile.NFProfilePutResponse, error) {

	nfProfilePutRequest := &dbproxy.PutRequest_NfProfilePutRequest{
		NfProfilePutRequest: putRequest,
	}

	putReqData := &dbproxy.PutRequest{
		Data: nfProfilePutRequest,
	}

	ret, err := Put(putReqData)
	if err == nil {
		return ret.GetNfProfilePutResponse(), nil
	}

	return nil, err
}

// PutNFProfileWithID is for NF profile register
func PutNFProfileWithID(nfInstanceID, nfProfile string, nfHelperInfo string) (*nfprofile.NFProfilePutResponse, error) {
	putReq := &nfprofile.NFProfilePutRequest{NfInstanceId: nfInstanceID, NfProfile: nfProfile, NfHelperInfo: nfHelperInfo}
	return PutNFProfile(putReq)
}

// GetNFProfile is for NF profile get
func GetNFProfile(getRequest *nfprofile.NFProfileGetRequest) (*nfprofile.NFProfileGetResponse, error) {

	nfProfileGetRequest := &dbproxy.GetRequest_NfProfileGetRequest{
		NfProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: nfProfileGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetNfProfileGetResponse(), nil
	}

	return nil, err
}

// GetNFProfileByID is for NF profile get
func GetNFProfileByID(nfInstanceID string) (*nfprofile.NFProfileGetResponse, error) {
	nfProfileKey := &nfprofile.NFProfileGetRequest_TargetNfInstanceId{
		TargetNfInstanceId: nfInstanceID,
	}
	getNFProfileReq := &nfprofile.NFProfileGetRequest{
		Data: nfProfileKey,
	}

	return GetNFProfile(getNFProfileReq)
}

// GetNFProfileCount is for NF profile count get
func GetNFProfileCount(getRequest *nfprofile.NFProfileCountGetRequest) (*nfprofile.NFProfileCountGetResponse, error) {

	nfProfileGetRequest := &dbproxy.GetRequest_NfProfileCountGetRequest{
		NfProfileCountGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: nfProfileGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetNfProfileCountGetResponse(), nil
	}

	return nil, err
}

// DeleteNFProfile is for NF profile delete
func DeleteNFProfile(delRequest *nfprofile.NFProfileDelRequest) (*nfprofile.NFProfileDelResponse, error) {

	nfProfileDelRequest := &dbproxy.DelRequest_NfProfileDelRequest{
		NfProfileDelRequest: delRequest,
	}

	delReqData := &dbproxy.DelRequest{
		Data: nfProfileDelRequest,
	}

	ret, err := Delete(delReqData)
	if err == nil {
		return ret.GetNfProfileDelResponse(), nil
	}

	return nil, err
}

// PutSubscription is for subscription register
func PutSubscription(subscriptionReq *subscription.SubscriptionPutRequest) (*subscription.SubscriptionPutResponse, error) {

	subscriptionPutRequest := &dbproxy.PutRequest_SubscriptionPutRequest{
		SubscriptionPutRequest: subscriptionReq,
	}

	putReqData := &dbproxy.PutRequest{
		Data: subscriptionPutRequest,
	}

	ret, err := Put(putReqData)
	if err == nil {
		return ret.GetSubscriptionPutResponse(), nil
	}

	return nil, err
}

// GetSubscription is for subscription get
func GetSubscription(subscriptionReq *subscription.SubscriptionGetRequest) (*subscription.SubscriptionGetResponse, error) {

	subscriptionGetRequest := &dbproxy.GetRequest_SubscriptionGetRequest{
		SubscriptionGetRequest: subscriptionReq,
	}

	getReqData := &dbproxy.GetRequest{
		Data: subscriptionGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetSubscriptionGetResponse(), nil
	}

	return nil, err
}

// DeleteSubscription is for subscription delete
func DeleteSubscription(subscriptionReq *subscription.SubscriptionDelRequest) (*subscription.SubscriptionDelResponse, error) {

	subscriptionDelRequest := &dbproxy.DelRequest_SubscriptionDelRequest{
		SubscriptionDelRequest: subscriptionReq,
	}

	delReqData := &dbproxy.DelRequest{
		Data: subscriptionDelRequest,
	}

	ret, err := Delete(delReqData)
	if err == nil {
		return ret.GetSubscriptionDelResponse(), nil
	}

	return nil, err
}

// PutNRFAddress is for puting NRF address
func PutNRFAddress(putRequest *nrfaddress.NRFAddressPutRequest) (*nrfaddress.NRFAddressPutResponse, error) {

	nrfAddressPutRequest := &dbproxy.PutRequest_NrfAddressPutRequest{
		NrfAddressPutRequest: putRequest,
	}

	putReqData := &dbproxy.PutRequest{
		Data: nrfAddressPutRequest,
	}

	ret, err := Put(putReqData)
	if err == nil {
		return ret.GetNrfAddressPutResponse(), nil
	}

	return nil, err
}

// GetNRFAddress is for geting NRF address
func GetNRFAddress(getRequest *nrfaddress.NRFAddressGetRequest) (*nrfaddress.NRFAddressGetResponse, error) {

	nrfAddressGetRequest := &dbproxy.GetRequest_NrfAddressGetRequest{
		NrfAddressGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: nrfAddressGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetNrfAddressGetResponse(), nil
	}

	return nil, err
}

// DeleteNRFAddress is for deleting NRF address
func DeleteNRFAddress(delRequest *nrfaddress.NRFAddressDelRequest) (*nrfaddress.NRFAddressDelResponse, error) {

	nrfAddressDelRequest := &dbproxy.DelRequest_NrfAddressDelRequest{
		NrfAddressDelRequest: delRequest,
	}

	delReqData := &dbproxy.DelRequest{
		Data: nrfAddressDelRequest,
	}

	ret, err := Delete(delReqData)
	if err == nil {
		return ret.GetNrfAddressDelResponse(), nil
	}

	return nil, err
}

// PutGroupProfile is for puting group profile
func PutGroupProfile(putRequest *groupprofile.GroupProfilePutRequest) (*groupprofile.GroupProfilePutResponse, error) {

	groupProfilePutRequest := &dbproxy.PutRequest_GroupProfilePutRequest{
		GroupProfilePutRequest: putRequest,
	}

	putReqData := &dbproxy.PutRequest{
		Data: groupProfilePutRequest,
	}

	ret, err := Put(putReqData)
	if err == nil {
		return ret.GetGroupProfilePutResponse(), nil
	}

	return nil, err
}

// GetGroupProfile is for geting group profile
func GetGroupProfile(getRequest *groupprofile.GroupProfileGetRequest) (*groupprofile.GroupProfileGetResponse, error) {

	groupProfileGetRequest := &dbproxy.GetRequest_GroupProfileGetRequest{
		GroupProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: groupProfileGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetGroupProfileGetResponse(), nil
	}
	return nil, err
}

//GetGroupProfileWithIndex is to get group profile with Index
func GetGroupProfileWithIndex(getRequest *groupprofile.GroupProfileGetRequest) (*groupprofile.GroupProfileGetResponse, error) {

	groupProfileGetRequest := &dbproxy.GetRequest_GroupProfileGetRequest{
		GroupProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: groupProfileGetRequest,
	}

	ctxLongTimeoutInSecond := time.Duration(constvalue.DbproxyGrpcCtxLongTimeout)

	ret, err := GetWithTimer(getReqData, ctxLongTimeoutInSecond)
	if err == nil {
		return ret.GetGroupProfileGetResponse(), nil
	}

	return nil, err
}

//DeleteGroupProfile is to delete group profile
func DeleteGroupProfile(delRequest *groupprofile.GroupProfileDelRequest) (*groupprofile.GroupProfileDelResponse, error) {

	groupProfileDelRequest := &dbproxy.DelRequest_GroupProfileDelRequest{
		GroupProfileDelRequest: delRequest,
	}

	delReqData := &dbproxy.DelRequest{
		Data: groupProfileDelRequest,
	}

	ret, err := Delete(delReqData)
	if err == nil {
		return ret.GetGroupProfileDelResponse(), nil
	}

	return nil, err
}

//GetImsiprefixProfile is for get imsiprefix profile
func GetImsiprefixProfile(getRequest *imsiprefixprofile.ImsiprefixProfileGetRequest) (*imsiprefixprofile.ImsiprefixProfileGetResponse, error) {

	imsiprefixProfileGetRequest := &dbproxy.GetRequest_ImsiprefixProfileGetRequest{
		ImsiprefixProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: imsiprefixProfileGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetImsiprefixProfileGetResponse(), nil
	}
	return nil, err
}

//PutGpsiProfile is for create or modify group profile.
func PutGpsiProfile(putRequest *gpsiprofile.GpsiProfilePutRequest) (*gpsiprofile.GpsiProfilePutResponse, error) {

	gpsiProfilePutRequest := &dbproxy.PutRequest_GpsiProfilePutRequest{
		GpsiProfilePutRequest: putRequest,
	}

	putReqData := &dbproxy.PutRequest{
		Data: gpsiProfilePutRequest,
	}

	ret, err := Put(putReqData)
	if err == nil {
		return ret.GetGpsiProfilePutResponse(), nil
	}
	return nil, err

}

//GetGpsiProfile is for get GpsiProfile
func GetGpsiProfile(getRequest *gpsiprofile.GpsiProfileGetRequest) (*gpsiprofile.GpsiProfileGetResponse, error) {

	gpsiProfileGetRequest := &dbproxy.GetRequest_GpsiProfileGetRequest{
		GpsiProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: gpsiProfileGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetGpsiProfileGetResponse(), nil
	}
	return nil, err
}

//GetGpsiProfileWithIndex is to get gpsi profile with Index
func GetGpsiProfileWithIndex(getRequest *gpsiprofile.GpsiProfileGetRequest) (*gpsiprofile.GpsiProfileGetResponse, error) {

	gpsiProfileGetRequest := &dbproxy.GetRequest_GpsiProfileGetRequest{
		GpsiProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: gpsiProfileGetRequest,
	}

	ctxLongTimeoutInSecond := time.Duration(constvalue.DbproxyGrpcCtxLongTimeout)

	ret, err := GetWithTimer(getReqData, ctxLongTimeoutInSecond)
	if err == nil {
		return ret.GetGpsiProfileGetResponse(), nil
	}
	return nil, err
}

//DeleteGpsiProfile is for delete GpsiProfile
func DeleteGpsiProfile(delRequest *gpsiprofile.GpsiProfileDelRequest) (*gpsiprofile.GpsiProfileDelResponse, error) {

	gpsiProfileDelRequest := &dbproxy.DelRequest_GpsiProfileDelRequest{
		GpsiProfileDelRequest: delRequest,
	}

	delReqData := &dbproxy.DelRequest{
		Data: gpsiProfileDelRequest,
	}

	ret, err := Delete(delReqData)
	if err == nil {
		return ret.GetGpsiProfileDelResponse(), nil
	}
	return nil, err
}

//GetGpsiprefixProfile is for get Gpsiprefix profile
func GetGpsiprefixProfile(getRequest *gpsiprefixprofile.GpsiprefixProfileGetRequest) (*gpsiprefixprofile.GpsiprefixProfileGetResponse, error) {

	GpsiprefixProfileGetRequest := &dbproxy.GetRequest_GpsiprefixProfileGetRequest{
		GpsiprefixProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: GpsiprefixProfileGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetGpsiprefixProfileGetResponse(), nil
	}
	return nil, err
}

// PutNRFProfile is for NRF profile register
func PutNRFProfile(putRequest *nrfprofile.NRFProfilePutRequest) (*nrfprofile.NRFProfilePutResponse, error) {

	nrfProfilePutRequest := &dbproxy.PutRequest_NrfProfilePutRequest{
		NrfProfilePutRequest: putRequest,
	}

	putReqData := &dbproxy.PutRequest{
		Data: nrfProfilePutRequest,
	}

	ret, err := Put(putReqData)
	if err == nil {
		return ret.GetNrfProfilePutResponse(), nil
	}
	return nil, err

}

// GetNRFProfile is for NRF profile get
func GetNRFProfile(getRequest *nrfprofile.NRFProfileGetRequest) (*nrfprofile.NRFProfileGetResponse, error) {

	nrfProfileGetRequest := &dbproxy.GetRequest_NrfProfileGetRequest{
		NrfProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: nrfProfileGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetNrfProfileGetResponse(), nil
	}
	return nil, err

}

// GetNRFProfileByID is for NRF profile get
func GetNRFProfileByID(nfInstanceID string) (*nrfprofile.NRFProfileGetResponse, error) {

	nrfProfileKey := &nrfprofile.NRFProfileGetRequest_NrfInstanceId{
		NrfInstanceId: nfInstanceID,
	}
	getNRFProfileReq := &nrfprofile.NRFProfileGetRequest{
		Data: nrfProfileKey,
	}

	return GetNRFProfile(getNRFProfileReq)
}

// DeleteNRFProfile is for NRF profile delete
func DeleteNRFProfile(delRequest *nrfprofile.NRFProfileDelRequest) (*nrfprofile.NRFProfileDelResponse, error) {

	nrfProfileDelRequest := &dbproxy.DelRequest_NrfProfileDelRequest{
		NrfProfileDelRequest: delRequest,
	}

	delReqData := &dbproxy.DelRequest{
		Data: nrfProfileDelRequest,
	}

	ret, err := Delete(delReqData)
	if err == nil {
		return ret.GetNrfProfileDelResponse(), nil
	}
	return nil, err

}

//SetDbclientPool is for set dbclientPool for UT
func SetDbclientPool(newDbclientPool []dbproxy.NFDataManagementServiceClient) {
	dbclientPool = newDbclientPool
}

//GetCacheNFProfile is for discovery to get local cache
func GetCacheNFProfile(getRequest *cachenfprofile.CacheNFProfileGetRequest) (*cachenfprofile.CacheNFProfileGetResponse, error) {
	cacheNFProfileGetRequest := &dbproxy.GetRequest_CacheNfProfileGetRequest{
		CacheNfProfileGetRequest: getRequest,
	}

	getReqData := &dbproxy.GetRequest{
		Data: cacheNFProfileGetRequest,
	}

	ret, err := Get(getReqData)
	if err == nil {
		return ret.GetCacheNfProfileGetResponse(), nil
	}

	return nil, err
}

//PutCacheNFProfile is for discovery to put local cache
func PutCacheNFProfile(putRequest *cachenfprofile.CacheNFProfilePutRequest) (*cachenfprofile.CacheNFProfilePutResponse, error) {
	cacheProfilePutRequest := &dbproxy.PutRequest_CacheNfProfilePutRequest{
		CacheNfProfilePutRequest: putRequest,
	}

	putReqData := &dbproxy.PutRequest{
		Data: cacheProfilePutRequest,
	}

	ret, err := Put(putReqData)
	if err == nil {
		return ret.GetCacheNfProfilePutResponse(), nil
	}
	return nil, err
}

// Insert insert data to region with key, value
func Insert(region, key, value string) (*dbproxy.InsertResponse, error) {
	insertRequest := &dbproxy.InsertRequest{
		RegionName: region,
	}

	item := &dbproxy.KVItem{
		Key:   key,
		Value: value,
	}
	insertRequest.Item = append(insertRequest.Item, item)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeoutInSecond)
	defer cancel()
	dbclient := dbclientPool[rand.Intn(poolLength)]

	insertRespone, err := dbclient.Insert(ctx, insertRequest)
	if err == nil {
		return insertRespone, err
	}

	return nil, err
}

// Remove remove data in region by key
func Remove(region string, key []string) (uint32, error) {
	removeRequest := &dbproxy.RemoveRequest{
		RegionName: region,
	}
	removeRequest.Key = append(removeRequest.Key, key...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeoutInSecond)
	defer cancel()
	dbclient := dbclientPool[rand.Intn(poolLength)]

	removeResponse, err := dbclient.Remove(ctx, removeRequest)
	if err != nil {
		log.Errorf("calling get fail : %v", err)
		return DbRPCError, err
	}

	if removeResponse.Code != DbDeleteSuccess {
		return removeResponse.Code, fmt.Errorf("Fail to remove data from DB, error code %d", removeResponse.Code)
	}

	return DbDeleteSuccess, nil
}

// GetByKey get value from region by key
func GetByKey(region, key string) (string, error) {
	queryRequest := &dbproxy.QueryRequest{
		RegionName: region,
	}
	queryRequest.Query = append(queryRequest.Query, key)

	queryRespone, err := QueryWithKey(queryRequest)
	if err != nil {
		log.Errorf("query by key failed : %v", err)
		return "", err
	}

	if queryRespone.Code != DbGetSuccess {
		return "", fmt.Errorf("Fail to get data from DB, error code %d", queryRespone.Code)
	}

	values := queryRespone.GetValue()
	if len(values) == 0 {
		return "", fmt.Errorf("Get empty data from DB")
	}

	return values[0], nil
}

// GetByOQL get value from region by OQL
func GetByOQL(region, oql string) ([]string, error) {
	var queryArray []string
	queryArray = append(queryArray, oql)
	queryReq := &dbproxy.QueryRequest{
		RegionName: region,
		Query:      queryArray,
	}

	queryResponse, err := QueryWithFilter(queryReq)
	if err != nil {
		return nil, err
	}
	if queryResponse.Code != DbGetSuccess {
		return make([]string, 0), nil
	}

	return queryResponse.Value, nil
}

//PatchNrfProfile is used to do patch nrfprofile by grpc
func PatchNrfProfile(patchReq *dbproxy.PatchRequest) (*dbproxy.PatchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*ctxTimeoutInSecond)
	defer cancel()
	dbclient := dbclientPool[rand.Intn(poolLength)]
	patchResp, err := dbclient.PatchNrfProfile(ctx, patchReq)
	if err != nil {
		log.Errorf("calling patch fail : %v", err)
		return nil, err
	}
	return patchResp, nil
}

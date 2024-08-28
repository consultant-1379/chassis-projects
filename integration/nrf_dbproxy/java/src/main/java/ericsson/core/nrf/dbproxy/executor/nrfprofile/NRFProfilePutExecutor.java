package ericsson.core.nrf.dbproxy.executor.nrfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.KeyAggregation;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileIndexProto.NRFKeyStruct;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfilePutRequestProto.NRFProfilePutRequest;
import ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfilePutHelper;
import java.util.HashMap;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFProfilePutExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(NRFProfilePutExecutor.class);

  private static NRFProfilePutExecutor instance;

  static {
    instance = null;
  }

  private NRFProfilePutExecutor() {
    super(NRFProfilePutHelper.getInstance());
  }

  public static synchronized NRFProfilePutExecutor getInstance() {
    if (null == instance) {
      instance = new NRFProfilePutExecutor();
    }

    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    NRFProfilePutRequest putRequest = request.getRequest().getPutRequest()
        .getNrfProfilePutRequest();
    String nrfInstanceId = putRequest.getNrfInstanceId();
    NRFProfile nrfProfile = createNRFProfile(putRequest);
    int code = ClientCacheService.getInstance()
        .put(Code.NRFPROFILE_INDICE, nrfInstanceId, nrfProfile);
    return new ExecutionResult(code);
  }

  private NRFProfile createNRFProfile(NRFProfilePutRequest request) {
    NRFProfile nrfProfile = new NRFProfile();
    nrfProfile.setRaw_data(request.getRawNrfProfile());

    nrfProfile.setAMFKey1(processNRFKeyStructList(request.getIndex().getAmfKey1List()));
    nrfProfile.setAMFKey2(processNRFKeyStructList(request.getIndex().getAmfKey2List()));
    nrfProfile.setAMFKey3(processNRFKeyStructList(request.getIndex().getAmfKey3List()));
    nrfProfile.setAMFKey4(processNRFKeyStructList(request.getIndex().getAmfKey4List()));

    nrfProfile.setSMFKey1(processNRFKeyStructList(request.getIndex().getSmfKey1List()));
    nrfProfile.setSMFKey2(processNRFKeyStructList(request.getIndex().getSmfKey2List()));
    nrfProfile.setSMFKey3(processNRFKeyStructList(request.getIndex().getSmfKey3List()));

    nrfProfile.setUDMKey1(processNRFKeyStructList(request.getIndex().getUdmKey1List()));
    nrfProfile.setUDMKey2(processNRFKeyStructList(request.getIndex().getUdmKey2List()));

    nrfProfile.setAUSFKey1(processNRFKeyStructList(request.getIndex().getAusfKey1List()));
    nrfProfile.setAUSFKey2(processNRFKeyStructList(request.getIndex().getAusfKey2List()));

    nrfProfile.setPCFKey1(processNRFKeyStructList(request.getIndex().getPcfKey1List()));
    nrfProfile.setPCFKey2(processNRFKeyStructList(request.getIndex().getPcfKey2List()));

    nrfProfile.setKey1(request.getIndex().getKey1());
    nrfProfile.setKey3(request.getIndex().getKey3());

    LOGGER.debug("NRF Profile : {}", nrfProfile.toString());

    return nrfProfile;
  }

  private HashMap processNRFKeyStructList(List<NRFKeyStruct> nrfKeyStructs) {
    HashMap result = new HashMap();
    int id = 0;
    for (NRFKeyStruct nks : nrfKeyStructs) {

      KeyAggregation ka = new KeyAggregation();
      ka.setSubKey1(nks.getSubKey1());
      ka.setSubKey2(nks.getSubKey2());
      ka.setSubKey3(nks.getSubKey3());
      ka.setSubKey4(nks.getSubKey4());
      ka.setSubKey5(nks.getSubKey5());

      result.put(id, ka);
      id++;
    }
    return result;
  }
}

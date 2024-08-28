package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileDelRequestProto.GpsiProfileDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileDelResponseProto.GpsiProfileDelResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GpsiProfileDelHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(GpsiProfileDelHelper.class);

  private static GpsiProfileDelHelper instance;

  private GpsiProfileDelHelper() {
  }

  public static synchronized GpsiProfileDelHelper getInstance() {
    if (null == instance) {
      instance = new GpsiProfileDelHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    GpsiProfileDelRequest request = message.getRequest().getDelRequest().getGpsiProfileDelRequest();
    String gpsiProfileId = request.getGpsiProfileId();
    if (gpsiProfileId.isEmpty()) {
      LOGGER.error("gpsiProfileId field is empty in GpsiProfileDelRequest");
      return Code.EMPTY_GPSI_PROFILE_ID;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    GpsiProfileDelResponse gpsiProfileDelResponse = GpsiProfileDelResponse.newBuilder()
        .setCode(code).build();
    DelResponse delResponse = DelResponse.newBuilder()
        .setGpsiProfileDelResponse(gpsiProfileDelResponse).build();
    return createNFMessage(delResponse);
  }
}

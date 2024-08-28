package ericsson.core.nrf.dbproxy.helper.nfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.FragmentSessionManagement;
import ericsson.core.nrf.dbproxy.common.FragmentUtil;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.config.AttributeConfig;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.AndExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.MetaExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.ORExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.SearchExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchAttribute;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchParameter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchValue;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.NFProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.ProvVersion;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.Range;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileGetRequestProto.NFProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileGetResponseProto.NFProfileGetResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import java.util.ArrayList;
import java.util.List;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NFProfileGetHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NFProfileGetHelper.class);

  private static NFProfileGetHelper instance;

  private NFProfileGetHelper() {
  }

  public static synchronized NFProfileGetHelper getInstance() {
    if (null == instance) {
      instance = new NFProfileGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {
    NFProfileGetRequest request = message.getRequest().getGetRequest().getNfProfileGetRequest();
    if (!request.getTargetNfInstanceId().isEmpty()) {
      return validateTargetNfInstanceId(request.getTargetNfInstanceId());
    } else if (request.hasFilter()) {
      return validateFilter(request.getFilter());
    } else if (!request.getFragmentSessionId().isEmpty()) {
      return Code.VALID;
    } else {
      LOGGER.error("Empty NFProfileGetRequest is received");
      return Code.NFMESSAGE_PROTOCOL_ERROR;
    }
  }

  private int validateTargetNfInstanceId(String targetNfInstanceId) {
    int code = Code.VALID;
    if (targetNfInstanceId.length() > Code.KEY_MAX_LENGTH) {
      code = Code.NF_INSTANCE_ID_LENGTH_EXCEED_MAX;
    }
    return code;
  }

  private int validateFilter(NFProfileFilter filter) {
    boolean emptyFilter = true;

    if (filter.hasExpiredTimeRange()) {
      emptyFilter = false;
      int code = validateRange(filter.getExpiredTimeRange());
      if (code != Code.VALID) {
        return code;
      }
    }

    if (filter.hasLastUpdateTimeRange()) {
      emptyFilter = false;
      int code = validateRange(filter.getLastUpdateTimeRange());
      if (code != Code.VALID) {
        return code;
      }
    }


    int code = validateProvisioned(filter.getProvisioned());
    if (code != Code.VALID) {
      return code;
    }

    if (filter.hasProvVersion()) {
      emptyFilter = false;
      code = validateProvVersion(filter.getProvVersion());
      if (code != Code.VALID) {
        return code;
      }
    }

    if (filter.hasSearchExpression()) {
      emptyFilter = false;
      code = validateSearchExpression(filter.getSearchExpression());

      if (code != Code.VALID) {
        return code;
      }
    }

    return emptyFilter ? Code.EMPTY_NFPROFILE_FILTER : Code.VALID;
  }

  private int validateRange(Range range) {
    int code = Code.VALID;
    if (range.getStart() > range.getEnd()) {
      code = Code.INVALID_RANGE;
    }
    return code;
  }

  private int validateProvisioned(int provisioned) {
    int code = Code.VALID;
    if (provisioned < Code.REGISTERED_PROVISIONED || provisioned > Code.PROVISIONED_ONLY) {
      code = Code.INVALID_PROVISIONED;
    }
    return code;
  }

  private int validateProvVersion(ProvVersion provVersion) {
    int code = Code.VALID;
    if (provVersion.getSupiVersion() < 0 || provVersion.getGpsiVersion() < 0) {
      code = Code.INVALID_PROV_VERSION;
    }
    return code;
  }

  private int validateSearchExpression(SearchExpression searchExpression) {
    if (searchExpression.hasAndExpression()) {
      return validateAndExpression(searchExpression.getAndExpression());
    }

    if (searchExpression.hasOrExpression()) {
      return validateORExpression(searchExpression.getOrExpression());
    }

    return Code.EMPTY_SEARCH_EXPRESSION;
  }

  private int validateAndExpression(AndExpression andExpression) {
    if (andExpression.getMetaExpressionCount() == 0) {
      return Code.EMPTY_AND_EXPRESSION;
    }

    int code = Code.VALID;
    for (MetaExpression metaExpression : andExpression.getMetaExpressionList()) {
      if (metaExpression.hasSearchParameter()) {
        code = validateSearchParameter(metaExpression.getSearchParameter());
      } else if (metaExpression.hasAndExpression()) {
        code = validateAndExpression(metaExpression.getAndExpression());
      } else if (metaExpression.hasOrExpression()) {
        code = validateORExpression(metaExpression.getOrExpression());
      } else {
        code = Code.EMPTY_META_EXPRESSION;
      }

      if (code != Code.VALID) {
        break;
      }
    }

    return code;
  }

  private int validateORExpression(ORExpression orExpression) {
    if (orExpression.getMetaExpressionCount() == 0) {
      return Code.EMPTY_OR_EXPRESSION;
    }

    int code = Code.VALID;
    for (MetaExpression metaExpression : orExpression.getMetaExpressionList()) {
      if (metaExpression.hasSearchParameter()) {
        code = validateSearchParameter(metaExpression.getSearchParameter());
      } else if (metaExpression.hasAndExpression()) {
        code = validateAndExpression(metaExpression.getAndExpression());
      } else if (metaExpression.hasOrExpression()) {
        code = validateORExpression(metaExpression.getOrExpression());
      } else {
        code = Code.EMPTY_META_EXPRESSION;
      }

      if (code != Code.VALID) {
        break;
      }
    }

    return code;
  }

  private int validateSearchParameter(SearchParameter searchParameter) {
    if (!searchParameter.hasAttribute()) {
      return Code.SEARCH_ATTRIBUTE_MISSED;
    }

    SearchAttribute attribute = searchParameter.getAttribute();
    if (attribute.getName().isEmpty()) {
      return Code.EMPTY_ATTRIBUTE_NAME;
    }
    if (attribute.getOperation() < Code.OPERATOR_LT
        || attribute.getOperation() > Code.OPERATOR_REGEX) {
      return Code.INVALID_ATTRIBUTE_OPERATOR;
    }
    if (AttributeConfig.getInstance().get(attribute.getName()) == null) {
      return Code.ATTRIBUTE_NOT_KNOWN;
    }

    if (!searchParameter.hasValue()) {
      return Code.SEARCH_VALUE_MISSED;
    }
    SearchValue value = searchParameter.getValue();
    if (!value.hasNum() && !value.hasStr()) {
      return Code.EMPTY_SEARCH_VALUE;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    NFProfileGetResponse nfProfileGetResponse = NFProfileGetResponse.newBuilder().setCode(code)
        .build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setNfProfileGetResponse(nfProfileGetResponse).build();
    return createNFMessage(getResponse);
  }

  @Override
  public NFMessage createResponse(ExecutionResult executionResult) {
    if (executionResult.getCode() != Code.SUCCESS) {
      return createResponse(executionResult.getCode());
    } else {
      SearchResult searchResult = (SearchResult) executionResult;
      if (searchResult.isFragmented()) {
        FragmentResult fragmentResult = (FragmentResult) searchResult;
        if (fragmentResult.getFragmentSessionID().isEmpty()) {
          int firstTransmitNum = FragmentUtil
              .transmitNumPerTime(fragmentResult, Code.NFPROFILE_INDICE);
          if (FragmentSessionManagement.getInstance().put(fragmentResult, firstTransmitNum)) {
            FragmentResult item = new FragmentResult();
            item.addAll(fragmentResult.getItems().subList(0, firstTransmitNum));
            item.setFragmentSessionID(fragmentResult.getFragmentSessionID());
            item.setTotalNumber(fragmentResult.getTotalNumber());
            item.setTransmittedNumber(fragmentResult.getTransmittedNumber());
            return createResponse(item);
          } else {
            return createResponse(Code.INTERNAL_ERROR);
          }
        } else {
          String fragmentSessionId = fragmentResult.getFragmentSessionID();
          int totalNumber = fragmentResult.getTotalNumber();
          int transmittedNumber = fragmentResult.getTransmittedNumber();
          FragmentInfo fragmentInfo = FragmentInfo.newBuilder()
              .setFragmentSessionId(fragmentSessionId).setTotalNumber(totalNumber)
              .setTransmittedNumber(transmittedNumber).build();

          List<String> nfProfileList = getNFProfile(fragmentResult);

          NFProfileGetResponse nfProfileGetResponse = NFProfileGetResponse.newBuilder()
              .setCode(fragmentResult.getCode()).addAllNfProfile(nfProfileList)
              .setFragmentInfo(fragmentInfo).build();
          GetResponse getResponse = GetResponse.newBuilder()
              .setNfProfileGetResponse(nfProfileGetResponse).build();
          return createNFMessage(getResponse);
        }
      } else {
        List<String> nfProfileList = getNFProfile(searchResult);
        NFProfileGetResponse nfProfileGetResponse = NFProfileGetResponse.newBuilder()
            .setCode(searchResult.getCode()).addAllNfProfile(nfProfileList).build();
        GetResponse getResponse = GetResponse.newBuilder()
            .setNfProfileGetResponse(nfProfileGetResponse).build();
        return createNFMessage(getResponse);
      }
    }
  }

  private List<String> getNFProfile(SearchResult searchResult) {
    List<String> nfProfileList = new ArrayList<>();
    for (Object obj : searchResult.getItems()) {
      try {
        nfProfileList.add(JSONFormatter.toJSON((PdxInstance) obj));
      } catch (Exception e) {
        LOGGER.error("Fail to format to JSON, " + e.toString());
      }
    }
    return nfProfileList;
  }
}

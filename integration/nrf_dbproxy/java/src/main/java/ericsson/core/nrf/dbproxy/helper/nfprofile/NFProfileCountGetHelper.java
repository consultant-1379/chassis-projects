package ericsson.core.nrf.dbproxy.helper.nfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.config.AttributeConfig;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.AndExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.MetaExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.ORExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.SearchExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchAttribute;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchParameter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchValue;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileCountGetRequestProto.NFProfileCountGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileCountGetResponseProto.NFProfileCountGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.NFProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.Range;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NFProfileCountGetHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NFProfileCountGetHelper.class);

  private static NFProfileCountGetHelper instance;

  private NFProfileCountGetHelper() {
  }

  public static synchronized NFProfileCountGetHelper getInstance() {
    if (null == instance) {
      instance = new NFProfileCountGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {
    NFProfileCountGetRequest request = message.getRequest().getGetRequest()
        .getNfProfileCountGetRequest();
    if (request.hasFilter()) {
      return validateFilter(request.getFilter());
    } else {
      LOGGER.error("Empty NFProfileCountGetRequest is received");
      return Code.NFMESSAGE_PROTOCOL_ERROR;
    }
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

    NFProfileCountGetResponse nfProfileCountGetResponse = NFProfileCountGetResponse.newBuilder()
        .setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setNfProfileCountGetResponse(nfProfileCountGetResponse).build();
    return createNFMessage(getResponse);
  }

  @Override
  public NFMessage createResponse(ExecutionResult executionResult) {
    if (executionResult.getCode() != Code.SUCCESS) {
      return createResponse(executionResult.getCode());
    }

    SearchResult searchResult = (SearchResult) executionResult;

    int count = getNFProfileCount(searchResult);
    if (count == -1) {
      return createResponse(Code.INTERNAL_ERROR);
    }

    NFProfileCountGetResponse nfProfileCountGetResponse = NFProfileCountGetResponse.newBuilder()
        .setCode(searchResult.getCode()).setCount(count).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setNfProfileCountGetResponse(nfProfileCountGetResponse).build();
    return createNFMessage(getResponse);
  }

  private int getNFProfileCount(SearchResult searchResult) {
    int count = -1;

    if (searchResult.getItems().size() < 1) {
      return count;
    }

    try {
      count = Integer.parseInt(String.valueOf(searchResult.getItems().get(0)));
    } catch (Exception e) {
      LOGGER.error("Fail to format to int, " + e.toString());
    }

    return count;
  }
}

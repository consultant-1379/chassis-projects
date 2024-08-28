package ericsson.core.nrf.dbproxy.executor.nfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.config.Attribute;
import ericsson.core.nrf.dbproxy.config.AttributeConfig;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.AndExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.MetaExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.ORExpression;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchAttribute;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchParameter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.SearchValue;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.NFProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileGetRequestProto.NFProfileGetRequest;
import ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileGetHelper;
import java.util.ArrayList;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NFProfileGetExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(NFProfileGetExecutor.class);

  private static final String OQL_1 = " AND ";
  private static final String OQL_2 = ".compareTo('";
  private static NFProfileGetExecutor instance;

  static {
    instance = null;
  }

  private String fromPrefix;
  private String select;
  private String where;

  private NFProfileGetExecutor() {
    super(NFProfileGetHelper.getInstance());

    fromPrefix = "/" + Code.NFPROFILE_INDICE + ".entrySet";
    select = "select DISTINCT value FROM ";
    where = " where ";

  }

  public static synchronized NFProfileGetExecutor getInstance() {
    if (null == instance) {
      instance = new NFProfileGetExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    NFProfileGetRequest getRequest = request.getRequest().getGetRequest().getNfProfileGetRequest();
    if (!getRequest.getTargetNfInstanceId().isEmpty()) {
      return ClientCacheService.getInstance()
          .getByID(Code.NFPROFILE_INDICE, getRequest.getTargetNfInstanceId());
    } else if (getRequest.hasFilter()) {
      String queryString = getQueryString(getRequest.getFilter());
      return ClientCacheService.getInstance().getByFilter(Code.NFPROFILE_INDICE, queryString);
    } else {
      return ClientCacheService.getInstance()
          .getByFragSessionId(Code.NFPROFILE_INDICE, getRequest.getFragmentSessionId());
    }
  }

  private String getQueryString(NFProfileFilter filter) {
    String[] expressions = {fromPrefix, ""};
    if (filter.hasSearchExpression() && filter.getSearchExpression().hasAndExpression()) {
      if (filter.getSearchExpression().hasAndExpression()) {
        expressions = constructAndExpression(filter.getSearchExpression().getAndExpression(), true,
            false, true);
      } else if (filter.getSearchExpression().hasOrExpression()) {
        expressions = constructOrExpression(filter.getSearchExpression().getOrExpression());
      }
    }

    String query = buildQuery(expressions);

    return addCustomInfo(filter, query);
  }

  private String buildQuery(String[] expressions) {
    String query = "";
    if (!expressions[1].isEmpty()) {
      if (expressions[0].isEmpty()) {
        query = select + fromPrefix + where + expressions[1];
      } else if (expressions[0].indexOf("select") == -1) {
        query = select + fromPrefix + ", " + expressions[0] + where + expressions[1];
      } else {
        query = select + expressions[0] + where + expressions[1];
      }
    }

    return query;
  }

  private String addCustomInfo(NFProfileFilter filter, String query) {
    boolean exist = false;
    StringBuilder subWhere = new StringBuilder(this.where);
    if (filter.hasExpiredTimeRange()) {
      long start = filter.getExpiredTimeRange().getStart();
      long end = filter.getExpiredTimeRange().getEnd();
      subWhere.append(
          "value.expiredTime >= " + Long.toString(start) + "L AND value.expiredTime <= " + Long
              .toString(end) + "L");
      exist = true;
    }

    if (filter.hasLastUpdateTimeRange()) {
      long start = filter.getLastUpdateTimeRange().getStart();
      long end = filter.getLastUpdateTimeRange().getEnd();
      if (exist) {
        subWhere.append(OQL_1);
      }
      subWhere.append(
          "value.lastUpdateTime >= " + Long.toString(start) + "L AND value.lastUpdateTime <= "
              + Long.toString(end) + "L");
      exist = true;
    }

    if (filter.getProvisioned() == Code.REGISTERED_ONLY
        || filter.getProvisioned() == Code.PROVISIONED_ONLY) {
      if (exist) {
        subWhere.append(OQL_1);
      }
      subWhere.append("value.provisioned = " + Integer.toString(filter.getProvisioned()));
    }

    if (filter.hasProvVersion()) {
      long supiVersion = filter.getProvVersion().getSupiVersion();
      long gpsiVersion = filter.getProvVersion().getGpsiVersion();
      if (exist) {
        subWhere.append(OQL_1);
      }
      subWhere.append("(value.provSupiVersion >= " + Long.toString(supiVersion)
          + "L OR value.provGpsiVersion >= " + Long.toString(gpsiVersion) + "L)");
      exist = true;
    }

    if (!exist) {
      return query;
    }

    if (query.isEmpty()) {
      return select + fromPrefix + subWhere.toString();
    } else {
      return select + "(" + query + ") value" + subWhere.toString();
    }
  }

  public String[] constructAndExpression(AndExpression andExpression, boolean inAndExpression,
      boolean innerExpressionExist, boolean isFirstAnd) {
    String[] expressions = {"", ""};
    List<String> fromExpressionList = new ArrayList<String>();
    for (MetaExpression metaExpression : andExpression.getMetaExpressionList()) {
      String[] subExpressions = {"", ""};
      if (metaExpression.hasSearchParameter()) {
        subExpressions = constructSearchParameterExpression(metaExpression.getSearchParameter());
      } else if (metaExpression.hasAndExpression()) {
        if (expressions[1].isEmpty()) {
          subExpressions = constructAndExpression(metaExpression.getAndExpression(),
              inAndExpression, true, false);
        } else {
          subExpressions = constructAndExpression(metaExpression.getAndExpression(),
              inAndExpression, false, false);
        }
      } else if (metaExpression.hasOrExpression()) {
        subExpressions = constructOrExpression(metaExpression.getOrExpression());
      } else {
        LOGGER.error("Empty MetaExpression in the AndExpression = " + andExpression.toString());
      }

      if (subExpressions[1].isEmpty()) {
        continue;
      }

      if (inAndExpression && !innerExpressionExist && isFirstAnd) {
        String innerQuery = buildQuery(expressions);
        if (innerQuery.isEmpty()) {
          expressions = subExpressions;
        } else {
          if (subExpressions[0].isEmpty()) {
            expressions[0] = "(" + innerQuery + ") value";
          } else {
            expressions[0] = "(" + innerQuery + ") value, " + subExpressions[0];
          }
          expressions[1] = subExpressions[1];
        }
      } else {
        if (!subExpressions[0].isEmpty()) {
          String[] fromExpressions = subExpressions[0].split(",");
          for (String from : fromExpressions) {
            if (from.isEmpty()) {
              continue;
            }
            if (fromExpressionList.contains(from)) {
              continue;
            }
            fromExpressionList.add(from);
            if (expressions[0].isEmpty()) {
              expressions[0] = from;
            } else {
              expressions[0] += "," + from;
            }
          }
        }

        if (!expressions[1].isEmpty()) {
          expressions[1] += OQL_1;
        }
        expressions[1] += subExpressions[1];
      }
    }

    if (!inAndExpression && !expressions[1].isEmpty()) {
      expressions[1] = "(" + expressions[1] + ")";
    }

    return expressions;
  }

  public String[] constructOrExpression(ORExpression orExpression) {
    String[] expressions = {"", ""};
    List<String> fromExpressionList = new ArrayList<String>();
    for (MetaExpression metaExpression : orExpression.getMetaExpressionList()) {
      String[] subExpressions = {"", ""};
      if (metaExpression.hasSearchParameter()) {
        subExpressions = constructSearchParameterExpression(metaExpression.getSearchParameter());
      } else if (metaExpression.hasAndExpression()) {
        subExpressions = constructAndExpression(metaExpression.getAndExpression(), false, true,
            false);
      } else if (metaExpression.hasOrExpression()) {
        subExpressions = constructOrExpression(metaExpression.getOrExpression());
      } else {
        LOGGER.error("Empty MetaExpression in the ORExpression = " + orExpression.toString());
      }

      if (subExpressions[1].isEmpty()) {
        continue;
      }

      if (!subExpressions[0].isEmpty()) {
        String[] fromExpressions = subExpressions[0].split(",");
        for (String from : fromExpressions) {
          if (from.isEmpty()) {
            continue;
          }
          if (fromExpressionList.contains(from)) {
            continue;
          }
          fromExpressionList.add(from);
          if (expressions[0].isEmpty()) {
            expressions[0] = from;
          } else {
            expressions[0] += "," + from;
          }
        }
      }

      if (!expressions[1].isEmpty()) {
        expressions[1] += " OR ";
      }
      expressions[1] += subExpressions[1];
    }

    if (!expressions[1].isEmpty()) {
      expressions[1] = "(" + expressions[1] + ")";
    }

    return expressions;
  }


  private String[] constructSearchParameterExpression(SearchParameter searchParameter) {
    SearchAttribute searchAttribute = searchParameter.getAttribute();
    Attribute attribute = AttributeConfig.getInstance().get(searchAttribute.getName());
    String name = attribute.getWhere();
    int operation = searchAttribute.getOperation();

    String[] expressions = {"", ""};
    SearchValue searchValue = searchParameter.getValue();
    if (searchValue.hasNum()) {
      String op = "";
      switch (operation) {
        case Code.OPERATOR_LT:
          op = " < ";
          break;
        case Code.OPERATOR_LE:
          op = " <= ";
          break;
        case Code.OPERATOR_EQ:
          op = " = ";
          break;
        case Code.OPERATOR_GE:
          op = " >= ";
          break;
        case Code.OPERATOR_GT:
          op = " > ";
          break;
        default:
          LOGGER.warn("Invalid operation = " + Long.toString(operation) + ", ignore this attribute "
              + searchAttribute.getName());
          return expressions;
      }
      expressions[0] = attribute.getFrom();
      expressions[1] = name + op + Long.toString(searchValue.getNum().getValue()) + "L";
    } else if (searchValue.hasStr()) {
      String str = searchValue.getStr().getValue();
      switch (operation) {
        case Code.OPERATOR_LT:
          expressions[1] = name + OQL_2 + str + "') < 0";
          break;
        case Code.OPERATOR_LE:
          expressions[1] = name + OQL_2 + str + "') <= 0";
          break;
        case Code.OPERATOR_EQ:
          expressions[1] = name + OQL_2 + str + "') = 0";
          break;
        case Code.OPERATOR_GE:
          expressions[1] = name + OQL_2 + str + "') >= 0";
          break;
        case Code.OPERATOR_GT:
          expressions[1] = name + OQL_2 + str + "') > 0";
          break;
        case Code.OPERATOR_REGEX:
          expressions[1] = "'" + str + "'.matches(" + name + ".toString()) = true";
          break;
        default:
          LOGGER
              .warn("Invalid operation = " + Long.toString(operation) + ", ignore this attribute");
          return expressions;
      }
      expressions[0] = attribute.getFrom();
    } else {
      LOGGER.debug("Empty search value for the attribute " + searchAttribute.getName());
    }

    return expressions;
  }
}

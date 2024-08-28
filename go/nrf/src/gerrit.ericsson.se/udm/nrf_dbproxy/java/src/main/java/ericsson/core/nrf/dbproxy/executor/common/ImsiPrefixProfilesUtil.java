package ericsson.core.nrf.dbproxy.executor.common;

import java.util.List;
import java.util.ArrayList;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.clientcache.schema.ImsiprefixProfiles;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileProto.ImsiprefixProfile;

public class ImsiPrefixProfilesUtil
{
    private static final Logger logger = LogManager.getLogger(ImsiPrefixProfilesUtil.class);

    private ImsiPrefixProfilesUtil(){}
    public static int AddImsiprefixProfile(ImsiprefixProfile imsiprefixProfile)
    {
        ExecutionResult getResult = ClientCacheService.getInstance().getByID(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix());
        if ( getResult.getCode() == Code.DATA_NOT_EXIST) {
            ImsiprefixProfiles imsiprefix_profiles = createImsiprefixProfiles(imsiprefixProfile.getImsiPrefix(), imsiprefixProfile.getValueInfo());
            int codePut = ClientCacheService.getInstance().put(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix(), imsiprefix_profiles);
            if (codePut != Code.CREATED) {
                logger.error("AddImsiprefixProfile: imsiprefix {} created failure, code:{}", imsiprefixProfile.getImsiPrefix(), codePut);
                return codePut;
            }
        } else if (getResult.getCode() == Code.SUCCESS) {
            SearchResult search_result = (SearchResult)getResult;
            List<ImsiprefixProfiles> imsiProfilesList = getImsiProfilesList(search_result.getItems());
            ImsiprefixProfiles imsiProfiles = imsiProfilesList.get(0);
            imsiProfiles.addValueInfo(imsiprefixProfile.getValueInfo());
            logger.debug("Begin add imsiProfiles: {}",imsiProfiles.toString());
            int codePut = ClientCacheService.getInstance().put(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix(), imsiProfiles);
            if (codePut != Code.CREATED) {
                logger.error("AddImsiprefixProfile: expend imsiprefix {} failure, code:{}", imsiprefixProfile.getImsiPrefix(), codePut);
                return codePut;
            }

        } else {
            logger.error("AddImsiprefixProfile: imsi prefix get failure.");
            return getResult.getCode();
        }
        return Code.CREATED;
    }

    public static int DelImsiprefixProfile(ImsiprefixProfile imsiprefixProfile)
    {
        ExecutionResult getResult = ClientCacheService.getInstance().getByID(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix());
        if ( getResult.getCode() == Code.DATA_NOT_EXIST) {
            logger.error("DelImsiprefixProfile imsiprefix not exist: {} ", imsiprefixProfile.getImsiPrefix());
        } else if (getResult.getCode() == Code.SUCCESS) {
            SearchResult search_result = (SearchResult)getResult;
            List<ImsiprefixProfiles> imsiProfilesList = getImsiProfilesList(search_result.getItems());
            ImsiprefixProfiles imsiProfiles = imsiProfilesList.get(0);
            imsiProfiles.rmValueInfo(imsiprefixProfile.getValueInfo());
            logger.debug("Begin change imsiProfiles: {}",imsiProfiles.toString());
            if (imsiProfiles.getValueInfo().size() != 0 ) {
                int codePut = ClientCacheService.getInstance().put(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix(), imsiProfiles);
                if (codePut != Code.CREATED) {
                    logger.error("DelImsiprefixProfile: reduce imsiprefix {} failure, code:{}", imsiprefixProfile.getImsiPrefix(), codePut);
                    return codePut;
                }
            } else {
                int codeDel = ClientCacheService.getInstance().delete(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix(), false);
                if (codeDel != Code.SUCCESS) {
                    logger.error("DelImsiprefixProfile: delete imsiprefix {} failure, code:{}", imsiprefixProfile.getImsiPrefix(), codeDel);
                    return codeDel;
                }
            }

        } else {
            logger.error("DelImsiprefixProfile: imsi prefix get failure.");
            return getResult.getCode();
        }
        return Code.SUCCESS;
    }

    private static ImsiprefixProfiles createImsiprefixProfiles(Long imsi_prefix, String value_info)
    {
        ImsiprefixProfiles imsiprefix_profiles = new ImsiprefixProfiles();

        imsiprefix_profiles.setImsiprefix(imsi_prefix);
        imsiprefix_profiles.addValueInfo(value_info);

        logger.debug("Imsiprefix Profile : {}", imsiprefix_profiles.toString());

        return imsiprefix_profiles;
    }

    private static List<ImsiprefixProfiles> getImsiProfilesList(List<Object> items)
    {
        List<ImsiprefixProfiles> imsiprofiles_list = new ArrayList<>();
        for(Object obj : items) {
            ImsiprefixProfiles imsiProfiles = (ImsiprefixProfiles)obj;
            imsiprofiles_list.add(imsiProfiles);
            logger.debug("getValueInfoList: value_info {}", imsiProfiles.toString());
        }
        return imsiprofiles_list;
    }
}

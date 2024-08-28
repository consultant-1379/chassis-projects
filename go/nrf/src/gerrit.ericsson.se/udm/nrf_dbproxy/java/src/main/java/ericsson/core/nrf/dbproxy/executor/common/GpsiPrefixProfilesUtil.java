package ericsson.core.nrf.dbproxy.executor.common;

import java.util.List;
import java.util.ArrayList;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiprefixProfiles;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileProto.GpsiprefixProfile;

public class GpsiPrefixProfilesUtil
{
    private static final Logger logger = LogManager.getLogger(GpsiPrefixProfilesUtil.class);

    private GpsiPrefixProfilesUtil() {

    }
    public static int AddGpsiprefixProfile(GpsiprefixProfile gpsiprefixProfile)
    {
        ExecutionResult getResult = ClientCacheService.getInstance().getByID(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix());
        if ( getResult.getCode() == Code.DATA_NOT_EXIST) {
            GpsiprefixProfiles gpsiprefix_profiles = createGpsiprefixProfiles(gpsiprefixProfile.getGpsiPrefix(), gpsiprefixProfile.getValueInfo());
            int codePut = ClientCacheService.getInstance().put(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix(), gpsiprefix_profiles);
            if (codePut != Code.CREATED) {
                logger.error("AddGpsiprefixProfile: gpsiprefix {} created failure, code:{}", gpsiprefixProfile.getGpsiPrefix(), codePut);
                return codePut;
            }
        } else if (getResult.getCode() == Code.SUCCESS) {
            SearchResult search_result = (SearchResult)getResult;
            List<GpsiprefixProfiles> gpsiProfilesList = getGpsiProfilesList(search_result.getItems());
            GpsiprefixProfiles gpsiProfiles = gpsiProfilesList.get(0);
            gpsiProfiles.addValueInfo(gpsiprefixProfile.getValueInfo());
            logger.debug("Begin add gpsiProfiles: {}",gpsiProfiles.toString());
            int codePut = ClientCacheService.getInstance().put(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix(), gpsiProfiles);
            if (codePut != Code.CREATED) {
                logger.error("AddGpsiprefixProfile: expend gpsiprefix {} failure, code:{}", gpsiprefixProfile.getGpsiPrefix(), codePut);
                return codePut;
            }

        } else {
            logger.error("AddGpsiprefixProfile: gpsi prefix get failure.");
            return getResult.getCode();
        }
        return Code.CREATED;
    }

    public static int DelGpsiprefixProfile(GpsiprefixProfile gpsiprefixProfile)
    {
        ExecutionResult getResult = ClientCacheService.getInstance().getByID(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix());
        if ( getResult.getCode() == Code.DATA_NOT_EXIST) {
            logger.error("DelGpsiprefixProfile gpsiprefix not exist: {} ", gpsiprefixProfile.getGpsiPrefix());
        } else if (getResult.getCode() == Code.SUCCESS) {
            SearchResult search_result = (SearchResult)getResult;
            List<GpsiprefixProfiles> gpsiProfilesList = getGpsiProfilesList(search_result.getItems());
            GpsiprefixProfiles gpsiProfiles = gpsiProfilesList.get(0);
            gpsiProfiles.rmValueInfo(gpsiprefixProfile.getValueInfo());
            logger.debug("Begin change gpsiProfiles: {}",gpsiProfiles.toString());
            if (gpsiProfiles.getValueInfo().size() != 0 ) {
                int codePut = ClientCacheService.getInstance().put(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix(), gpsiProfiles);
                if (codePut != Code.CREATED) {
                    logger.error("DelGpsiprefixProfile: reduce gpsiprefix {} failure, code:{}", gpsiprefixProfile.getGpsiPrefix(), codePut);
                    return codePut;
                }
            } else {
                int codeDel = ClientCacheService.getInstance().delete(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix(), false);
                if (codeDel != Code.SUCCESS) {
                    logger.error("DelGpsiprefixProfile: delete gpsiprefix {} failure, code:{}", gpsiprefixProfile.getGpsiPrefix(), codeDel);
                    return codeDel;
                }
            }

        } else {
            logger.error("DelGpsiprefixProfile: gpsi prefix get failure.");
            return getResult.getCode();
        }
        return Code.SUCCESS;
    }

    private static GpsiprefixProfiles createGpsiprefixProfiles(Long gpsi_prefix, String value_info)
    {
        GpsiprefixProfiles gpsiprefix_profiles = new GpsiprefixProfiles();

        gpsiprefix_profiles.setGpsiprefix(gpsi_prefix);
        gpsiprefix_profiles.addValueInfo(value_info);

        logger.debug("Gpsiprefix Profile : {}", gpsiprefix_profiles.toString());

        return gpsiprefix_profiles;
    }

    private static List<GpsiprefixProfiles> getGpsiProfilesList(List<Object> items)
    {
        List<GpsiprefixProfiles> gpsiprofiles_list = new ArrayList<>();
        for(Object obj : items) {
            GpsiprefixProfiles gpsiProfiles = (GpsiprefixProfiles)obj;
            gpsiprofiles_list.add(gpsiProfiles);
            logger.debug("getValueInfoList: value_info {}", gpsiProfiles.toString());
        }
        return gpsiprofiles_list;
    }
}

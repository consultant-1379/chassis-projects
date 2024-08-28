package com.ericsson.nrf;

import com.ericsson.adp.kvdbag.adminmgrapi.model.AppJar;
import com.ericsson.adp.kvdbag.adminmgrapi.model.AppJarCategory;
import com.ericsson.adp.kvdbag.adminmgrapi.model.GfshCommand;
import com.ericsson.nrf.common.Constants;
import com.ericsson.nrf.config.ConfigLoader;
import com.ericsson.nrf.config.IndexBean;
import com.ericsson.nrf.config.RegionBean;
import com.ericsson.nrf.handler.ActionsHandler;
import com.ericsson.nrf.handler.AdminMgrClient;
import com.ericsson.nrf.handler.AppJarsHandler;
import com.ericsson.nrf.handler.RegionHandler;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.List;
import java.util.concurrent.CountDownLatch;

public class ConfigureGeode {
    private static final Logger LOGGER = LogManager.getLogger(ConfigureGeode.class);

    private static final String JAR_VERSION = "1.0.0";
    private static AppJarsHandler appJarsHandler;
    private static ActionsHandler actionsHandler;
    private static RegionHandler regionHandler;
    private static final int THREAD_NUM = 2;
    private static final CountDownLatch countDownLatch = new CountDownLatch(THREAD_NUM);

    private static AdminMgrClient adminMgrClient = new AdminMgrClient();

    public static void main(String[] args) {
        appJarsHandler = new AppJarsHandler(adminMgrClient);
        actionsHandler = new ActionsHandler(adminMgrClient);
        regionHandler = new RegionHandler(adminMgrClient);

        // Load region file and index file
        ConfigLoader.getInstance().parseRegionConfigFile();
        ConfigLoader.getInstance().parseIndexConfigFile();

        LOGGER.info("Check if Admin Manager pod is up and running");
        while(!adminMgrClient.waitUntilAdminMgrIsUpAndRunning()) {
            LOGGER.info("Admin Manager pod is not running");
        }

        // In this example, assumption is that application jar and/or multisite features are enabled in values.yaml,
        // so installation of kvdb with command "helm install kvdb_chart ..."
        // will automaticaly start only locator and admin-mgr pods
//        LOGGER.info("Execute start locators command which will wait until all locators are started");
//        LOGGER.info("START locators action returned: " + actionsHandler.executeStartActionOnAllLocators());

        // CONFIGURING MULTI SITE
//        LOGGER.info("Configuring multisite (use local site only)");
//        sitesHandler.configureGeoRedWithOneSite();


        // CONFIGURE PDX
        LOGGER.info("Configure PDX (read serialized and use persistence)");
        LOGGER.info("Configure PDX command output: "
                + adminMgrClient.executeGfshCommand("configure pdx --disk-store=DEFAULT --read-serialized=true --auto-serializable-classes=ericsson.core.nrf.dbproxy.clientcache.schema.*"));

        LOGGER.info("---------------------Start to deploy custom jar into kvdb---------------------------");
        int retry = 0;
        while (!deployJars()) {
            retry++;
            if (retry > 3) {
                LOGGER.error("deploy jar failed, progress will exit");
                System.exit(1);
            }
        }

        // START SERVERS
        LOGGER.info("Now when application jars are deployed and pdx is configured start servers as well");
        LOGGER.info("START all members returned: " + actionsHandler.executeStartActionOnAllMembers());
        // Create regions from configuration file region.json
        if (!createRegions()) {
            System.exit(1);
        }

        // Check if all the region created are exist in DB
        if (!checkRegionsExist()) {
            System.exit(1);
        }

        // Create indexes from configuration file index.json
        if (!createIndexes()) {
            System.exit(1);
        }
        try {
            countDownLatch.await();
            LOGGER.debug("---------------------------------successfully create indexes from configuration file index.json--------------------------------------");
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        // Check if all the index created are exist in DB
        if (!checkIndexesExist()) {
            System.exit(1);
        }
    }

    /**
     * upload and deploy the app jars
     *
     * @return execute success of failure
     */
    private static boolean deployJars() {
        // UPLOAD AND DEPLOY APP JARS
        LOGGER.info("Upload and deploy all application jars");

        LOGGER.info("Check that partitioner named " + Constants.PARTITIONER_JAR_FILE + " if is already deployed");
        if (appJarsHandler.checkAppJarExists(Constants.PARTITIONER_JAR_FILE)) {
            LOGGER.info(Constants.PARTITIONER_JAR_FILE + " should not have been deployed");
            return true;
        }

        LOGGER.info("Upload the " + Constants.PARTITIONER_JAR_FILE);
        String partitionerName = Constants.JARS_LOCATION + Constants.PARTITIONER_JAR_FILE;
        AppJar partitioner = appJarsHandler.uploadAppJar(partitionerName, JAR_VERSION,
                AppJarCategory.PARTITION_RESOLVER, "First partition resolver");

        LOGGER.info("Deploy partition resolver jar");
        GfshCommand deployPartitionerCommand = appJarsHandler.deployAppJar(partitioner);
        LOGGER.info("Deploy partitioner command output: " + deployPartitionerCommand.getOutput());
        if (deployPartitionerCommand.getStatusCode() != Constants.ADM_MGR_OK_EXEC_STATUS_CODE) {
            LOGGER.error("The execution of " + deployPartitionerCommand.getCommand() + " failed. Error: "
                    + deployPartitionerCommand.getOutput());
            return false;
        }
        return true;
    }

    /**
     * create regions via admin mgr service
     *
     * @return create regions successful or fail
     */
    private static boolean createRegions() {
        LOGGER.debug("----------------------start to create regions from configuration file region.json-----------------------------");
        // CREATE REGIONS
        List<RegionBean> regionList = ConfigLoader.getInstance().getRegionList();
        RegionBean bean;
        for (int i = 0; i < regionList.size(); i++) {
            bean = regionList.get(i);
            int retry = 0;
            // create region and retry 3 times at most if failure
            while (!regionHandler.createRegion(bean.getName(), bean.getType(), bean.getAdditional())) {
                retry++;
                if (retry > 3) {
                    return false;
                }
            }
        }
        LOGGER.debug("------------------------successfully create regions from configuration file region.json---------------------------");
        return true;
    }

    /**
     * Check if all the regions in configuration file are created successfully
     *
     * @return all the regions exist or not
     */
    private static boolean checkRegionsExist() {
        LOGGER.debug("start to check if all regions from configuration file region.json are exist in DB");
        List<RegionBean> indexList = ConfigLoader.getInstance().getRegionList();
        RegionBean bean;
        String regionsInDB = adminMgrClient.executeGfshCommand("list regions");
        for (int i = 0; i < indexList.size(); i++) {
            bean = indexList.get(i);
            int retryCheck = 0;
            while (!regionsInDB.contains(bean.getName())) {
                retryCheck++;
                if (retryCheck > 3) {
                    return false;
                }
                int retry = 0;
                // create region and retry 3 times at most if failure
                while (!regionHandler.createRegion(bean.getName(), bean.getType(), bean.getAdditional())) {
                    retry++;
                    if (retry > 3) {
                        return false;
                    }
                }
                regionsInDB = adminMgrClient.executeGfshCommand("list regions");
            }

        }
        LOGGER.debug("successfully check all regions from configuration file region.json are exist in DB");
        return true;
    }

    /**
     * create indexes via admin mgr service, start two threads to create indexes to decrease the time cost
     *
     * @return create indexes successful or fail
     */
    private static boolean createIndexes() {
        LOGGER.debug("-----------------------------------start to create indexes from configuration file index.json-----------------------------------------");
        List<IndexBean> indexList = ConfigLoader.getInstance().getIndexList();
        int midIndex=indexList.size()/2;
        new Thread(() -> {
            IndexBean bean;
            for (int i = 0; i < midIndex; i++) {
                bean = indexList.get(i);
                int retry = 0;
                // create index and retry 3 times at most if failure
                while (!regionHandler.createIndex(bean.getName(), bean.getExpression(), bean.getRegion(), bean.getType())) {
                    retry++;
                    if (retry > 3) {
                        break;
                    }
                }
            }
            countDownLatch.countDown();
        }).start();
        new Thread(() -> {
            IndexBean bean;
            for (int i = midIndex; i < indexList.size(); i++) {
                bean = indexList.get(i);
                int retry = 0;
                // create index and retry 3 times at most if failure
                while (!regionHandler.createIndex(bean.getName(), bean.getExpression(), bean.getRegion(), bean.getType())) {
                    retry++;
                    if (retry > 3) {
                        break;
                    }
                }
            }
            countDownLatch.countDown();
        }).start();

        return true;
    }

    /**
     * Check if all the indexes in the configuration file are created successfully
     *
     * @return if indexes exist or not
     */
    private static boolean checkIndexesExist() {
        LOGGER.debug("start to check if all the indexes in configuration file index.json are existed in DB");
        List<IndexBean> indexList = ConfigLoader.getInstance().getIndexList();
        IndexBean bean;
        String indexInDB = adminMgrClient.executeGfshCommand("list indexes");
        for (int i = 0; i < indexList.size(); i++) {
            bean = indexList.get(i);
            int retryCheck = 0;
            while (!indexInDB.contains(bean.getName())) {
                retryCheck++;
                if (retryCheck > 3) {
                    return false;
                }
                int retry = 0;
                // create index and retry 3 times at most if failure
                while (!regionHandler.createIndex(bean.getName(), bean.getExpression(), bean.getRegion(), bean.getType())) {
                    retry++;
                    if (retry > 3) {
                        return false;
                    }
                }
                indexInDB = adminMgrClient.executeGfshCommand("list indexes");
            }

        }
        LOGGER.debug("successfully check all the indexes in configuration file index.json are existed in DB");
        return true;
    }
}

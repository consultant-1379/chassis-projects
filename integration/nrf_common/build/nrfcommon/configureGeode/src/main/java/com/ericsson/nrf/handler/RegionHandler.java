package com.ericsson.nrf.handler;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class RegionHandler {

    /**
     * Logger.
     */
    private static final Logger LOGGER = LogManager.getLogger(RegionHandler.class);

    /**
     * Admin Mgr Client.
     */
    private AdminMgrClient adminMgrClient;

    /**
     * Constructor.
     *
     * @param adminMgrClient REST client for interacting with admin-mgr process.
     */
    public RegionHandler(final AdminMgrClient adminMgrClient) {
        this.adminMgrClient = adminMgrClient;
    }

    /**
     * Execute create region gfsh command and check output od list region command contains created region.
     *
     * @param regionName name of region to create
     * @param regionType type of region to create
     * @return true if region is created
     * false otherwise
     */
    public final boolean createRegion(final String regionName, final String regionType) {
        return this.createRegion(regionName, regionType, "");
    }

    /**
     * Execute create region gfsh command and check output of list region command contains created region.
     *
     * @param regionName       name of region to create
     * @param regionType       type of region to create
     * @param additionalParams string with all additional parameters and its values
     * @return true if region is created
     * false otherwise
     */
    public final boolean createRegion(final String regionName,
                                      final String regionType,
                                      final String additionalParams) {
//        if (this.checkIfRegionExists(regionName)) {
//            LOGGER.debug("region=" + regionName + " already exist");
//            return true;
//        }
        String createRegionCommandString = "create region --name=" + regionName
                + " --type=" + regionType
                + " " + additionalParams;
        LOGGER.info("Executing gfsh command: " + createRegionCommandString);
        return this.adminMgrClient.isGfshCommandExecuted(createRegionCommandString);

//        try {
//            Thread.sleep(Constants.MILLISECONDS_AFTER_CREATE_REGION);
//        } catch (InterruptedException exception) {
//            LOGGER.error("Unexpected exception occurred.", exception);
//        }
//
//        return this.checkIfRegionExists(regionName);
    }

    /**
     * Check does output of list region command contains a region with the given name.
     *
     * @param regionName name of region to check
     * @return true if region found
     * false otherwise
     */
//    public final boolean checkIfRegionExists(final String regionName) {
//        return this.adminMgrClient.executeGfshCommand("list regions").contains(regionName);
//    }

    /**
     * Execute create index gfsh command and check output of list indexes command contains created index.
     *
     * @param regionName      name of index to create
     * @param indexExpression expression of index to create
     * @param regionName      region name
     * @param indexType       type of index to create
     * @return true if index is created
     * false otherwise
     */
    public final boolean createIndex(final String indexName,
                                     final String indexExpression,
                                     final String regionName,
                                     final String indexType) {
//        if (this.checkIfIndexExists(indexName)) {
//            LOGGER.debug("index=" + indexName + " already exist");
//            return true;
//        }
        String createIndexCommandString = "create index"
                + " --name=" + indexName
                + " --expression=" + indexExpression
                + " --region=\"" + regionName
                + "\" --type=" + indexType;
        LOGGER.info("Executing gfsh command: " + createIndexCommandString);
        return this.adminMgrClient.isGfshCommandExecuted(createIndexCommandString);

//        try {
//            Thread.sleep(Constants.MILLISECONDS_AFTER_CREATE_REGION);
//        } catch (InterruptedException exception) {
//            LOGGER.error("Unexpected exception occurred while creating index.", exception);
//        }
//
//        return this.checkIfIndexExists(indexName);
    }

    /**
     * Check does output of list indexes command contains a index with the given name.
     *
     * @param indexName name of index to check
     * @return true if index found
     * false otherwise
     */
//    public final boolean checkIfIndexExists(final String indexName) {
//        return this.adminMgrClient.executeGfshCommand("list indexes").contains(indexName);
//    }

}
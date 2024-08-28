package com.ericsson.nrf.handler;

import com.ericsson.adp.kvdbag.adminmgrapi.ApiException;
import com.ericsson.adp.kvdbag.adminmgrapi.ApiResponse;
import com.ericsson.adp.kvdbag.adminmgrapi.model.AppJar;
import com.ericsson.adp.kvdbag.adminmgrapi.model.AppJarCategory;
import com.ericsson.adp.kvdbag.adminmgrapi.model.GfshCommand;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import javax.ws.rs.core.Response;
import java.io.File;

/**
 * Helper for operations with Application Jars.
 */
public class AppJarsHandler {

    /**
     * Deploy JAR command.
     */
    static final String DEPLOY_JAR_COMMAND = "deploy --jar=";
    /**
     * Logger.
     */
    private static final Logger LOGGER = LogManager.getLogger(AppJarsHandler.class);
    /**
     * Admin Mgr Client.
     */
    private AdminMgrClient adminMgrClient;

    /**
     * Constructor.
     *
     * @param adminMgrClient REST client for interacting with admin-mgr process.
     */
    public AppJarsHandler(final AdminMgrClient adminMgrClient) {
        this.adminMgrClient = adminMgrClient;
    }

    /**
     * Upload app jar with the given name.
     *
     * @param appJarName name with full path of the app jar to be uploaded
     * @param version    of app jar to deploy
     * @param category   of app jar to deploy
     * @param metadata   optional metadata of app jar
     * @return uploaded jar
     */
    public final AppJar uploadAppJar(final String appJarName, final String version,
                                     final AppJarCategory category, final String metadata) {
        AppJar appJar = null;
        try {
            File appJarFile = new File(appJarName);
            ApiResponse<AppJar> appJarApiResponse =
                    this.adminMgrClient.getDefaultApi().uploadAppJarWithHttpInfo(appJarFile.getName(), appJarFile,
                            version, category, metadata);

            appJar = appJarApiResponse.getData();
            LOGGER.info("Uploaded app jar named: {}", appJarName);
        } catch (ApiException exception) {
            LOGGER.error("Error occurred while trying to upload app jar named "
                    + appJarName, exception);
        }
        return appJar;
    }

    /**
     * Deploys the provided app jar.
     *
     * @param appJar app jar to be deployed
     * @return gfshCommand result of deploying jar command
     */
    public final GfshCommand deployAppJar(final AppJar appJar) {
        GfshCommand deployJarCommand = new GfshCommand();
        String deployJarCommandString = DEPLOY_JAR_COMMAND + appJar.getUploadDirectory()
                + File.separator + appJar.getName();

        deployJarCommand.setCommand(deployJarCommandString);
        return this.adminMgrClient.scheduleAndWait(deployJarCommand);
    }

    /**
     * Check if application jar is deployed successfully.
     *
     * @param appJarName name of the app jar
     * @return true if app jar is deployed successfully
     * false otherwise
     */
    public final boolean verifyAppJarDeployedSuccessfully(final String appJarName) {
        boolean deployedSuccessfully = true;
        try {
            ApiResponse<AppJar> apiResponse =
                    this.adminMgrClient.getDefaultApi().getAppJarInfoWithHttpInfo(appJarName);

            AppJar appJar = apiResponse.getData();
            LOGGER.info("Got the following application jar info: {}", appJar);
            if (!appJar.getDeployed()) {
                LOGGER.info("App jar deployment status in not Done!");
                deployedSuccessfully = false;
            }
        } catch (ApiException exception) {
            LOGGER.info("Exception occurred while trying to get deployment status", exception);
            deployedSuccessfully = false;
        }
        return deployedSuccessfully;
    }

    /**
     * Get app jar info, if ApiException is caught with Code 404 (Not Found) and summary describing AppJar is not found
     * then AppJar doesn't exist.
     *
     * @param appJarToCheck name of App Jar to check
     * @return true if app jar exists
     * false otherwise
     */
    public final boolean checkAppJarExists(final String appJarToCheck) {
        boolean appJarExists = true;
        try {
            this.adminMgrClient.getDefaultApi().getAppJarInfoWithHttpInfo(appJarToCheck);
        } catch (ApiException exception) {
            LOGGER.info("Exception occurred while trying to get app jar info");
            if (exception.getCode() == Response.Status.NOT_FOUND.getStatusCode()) {
                appJarExists = false;
            }
        }
        return appJarExists;
    }


}

package com.ericsson.nrf.handler;

import com.ericsson.adp.kvdbag.adminmgrapi.ApiClient;
import com.ericsson.adp.kvdbag.adminmgrapi.ApiException;
import com.ericsson.adp.kvdbag.adminmgrapi.ApiResponse;
import com.ericsson.adp.kvdbag.adminmgrapi.client.DefaultApi;
import com.ericsson.adp.kvdbag.adminmgrapi.model.*;
import com.ericsson.nrf.common.Constants;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.time.OffsetDateTime;

public class AdminMgrClient {

    /**
     * Timeout for waiting until admin mgr is up, in seconds.
     */
    public static final Integer ADMIN_MGR_STARTUP_TIMEOUT = 100;
    /**
     * AdminMgr API base endpoint path.
     */
    public static final String REST_API_PREFIX = "/kvdb-ag/management/v1";
    /**
     * Port for AdminMgr pod.
     */
    public static final int ADMIN_MGR_PORT = 8080;
    /**
     * Milliseconds between GFSH command execution status checks.
     */
    protected static final int MILLISECONDS_BETWEEN_CHECKS = 200;
    /**
     * AdminMgr service name.
     */
    private static final String ADMIN_MGR_SVC = System.getenv("ADMIN_MGR_SVC");
    /**
     * Logger.
     */
    private static final Logger LOGGER = LogManager.getLogger(AdminMgrClient.class);
    /**
     * Timeout (in milliseconds) to be used when establishing a TCP connection.
     */
    private static final int CONNECTION_TIMEOUT = 5000;
    /**
     * Timeout (in milliseconds) to be used when reading server data.
     */
    private static final int READ_TIMEOUT = 10000;
    /**
     * If the execution of deploy command last more than this value, it is considered timeout-ed.
     * The value is in seconds.
     */
    private static final long GFSH_COMMAND_TIMEOUT = 40;
    /**
     * Default KVDB admin-mgr API.
     */
    private DefaultApi defaultApi;

    /**
     * Constructor for AdminMgrClient.
     * Creating defaultApi with AdminMgr path.
     */
    public AdminMgrClient() {
        ApiClient apiClient = new ApiClient();
        String adminMgrBasePath = "http://" + ADMIN_MGR_SVC + ":"
                + ADMIN_MGR_PORT + REST_API_PREFIX;
        LOGGER.info("adminMgrBasePath: " + adminMgrBasePath);
        apiClient.setConnectTimeout(CONNECTION_TIMEOUT);
        apiClient.setReadTimeout(READ_TIMEOUT);
        apiClient.setBasePath(adminMgrBasePath);
        this.defaultApi = new DefaultApi(apiClient);
    }

    /**
     * Get KVDB admin-mgr API.
     *
     * @return default KVDB admin-mgr API
     */
    protected final DefaultApi getDefaultApi() {
        return defaultApi;
    }

    /**
     * Try to send get request to Admin Mgr api and wait until it respond
     *
     * @return true if admin mgr pod is up and running, false otherwise
     */
    public boolean waitUntilAdminMgrIsUpAndRunning() {
        long startTime = OffsetDateTime.now().toEpochSecond();

        while (OffsetDateTime.now().toEpochSecond() - startTime < ADMIN_MGR_STARTUP_TIMEOUT) {
            LOGGER.info("Waiting for admin mgr pod to be up and running...");
            try {
                // trying any get request to check if we will get response
                ApiResponse<GfshCommands> listGfshCommandsApiResponse =
                        this.defaultApi.getCommandListWithHttpInfo();

                LOGGER.info("list gfsh commands response=" + listGfshCommandsApiResponse.getStatusCode() + "------" + listGfshCommandsApiResponse.getData().toString());
                // admin-mgr is up and running if exception was not thrown
                return true;

            } catch (Exception e) {
                LOGGER.info("Exception occurred while executing get request to check admin-mgr api state", e.getMessage());
            }
        }

        return false;
    }

    /**
     * Schedules the given GFSH command for execution and waits for it to complete execution.
     *
     * @param gfshCommand command to be executed
     * @return executed GfshCommand
     */
    public final GfshCommand scheduleAndWait(final GfshCommand gfshCommand) {
        LOGGER.info("Executing GfshCommand...");
        GfshCommand updatedGfshCommand = null;
        try {
            ApiResponse<GfshCommandId> gfshCommandIdApiResponse = this.defaultApi.addCommandWithHttpInfo(gfshCommand);

            GfshCommandId gfshCommandId = gfshCommandIdApiResponse.getData();

            boolean executed = false;

            do {
                LOGGER.info("Waiting for " + gfshCommand.getCommand() + " to be executed");
                Thread.sleep(MILLISECONDS_BETWEEN_CHECKS);

                ApiResponse<GfshCommand> gfshCommandApiResponse =
                        this.defaultApi.getCommandInfoWithHttpInfo(gfshCommandId.getCommandId());
                updatedGfshCommand = gfshCommandApiResponse.getData();

                long secondsSinceScheduling =
                        OffsetDateTime.now().toEpochSecond()
                                - updatedGfshCommand.getReceivedTimestamp().toEpochSecond();
                LOGGER.info("Seconds since scheduling command: " + secondsSinceScheduling);

                if (updatedGfshCommand.getExecutionStatus() == ExecutionStatus.EXECUTED) {
                    executed = true;
                }

                if (secondsSinceScheduling > GFSH_COMMAND_TIMEOUT && !executed) {
                    LOGGER.error("Timeout occurred while executing " + gfshCommand.getCommand());
                    break;
                }
            } while (!executed);

        } catch (InterruptedException exception) {
            LOGGER.error("Execution interrupted!", exception);
        } catch (ApiException exception) {
            LOGGER.error("Error occurred while executing " + gfshCommand.getCommand()
                    + " exception: " + exception, exception);
        }
        return updatedGfshCommand;
    }


    /**
     * Execute gfsh command and check if the command is executed successfully.
     *
     * @param commandToExecute command to execute
     * @return string containing output of the command
     */
    public final String executeGfshCommand(final String commandToExecute) {
        GfshCommand gfshCommand = new GfshCommand();

        gfshCommand.setCommand(commandToExecute);

        gfshCommand = this.scheduleAndWait(gfshCommand);
        if (gfshCommand.getStatusCode() != Constants.ADM_MGR_OK_EXEC_STATUS_CODE) {
            LOGGER.error("The execution of " + gfshCommand.getCommand() + "failed. Error: "
                    + gfshCommand.getOutput());
        }
        return gfshCommand.getOutput();
    }

    /**
     * Execute gfsh command and check if the command is executed successfully
     *
     * @param commandToExecute command to execute
     * @return Whether the command execute successfully
     */
    public final boolean isGfshCommandExecuted(final String commandToExecute) {
        GfshCommand gfshCommand = new GfshCommand();

        gfshCommand.setCommand(commandToExecute);

        gfshCommand = this.scheduleAndWait(gfshCommand);
        if (gfshCommand.getStatusCode() != Constants.ADM_MGR_OK_EXEC_STATUS_CODE) {
            LOGGER.error("The execution of " + gfshCommand.getCommand() + "failed. Error: "
                    + gfshCommand.getOutput());
            return false;
        }
        return true;
    }

    /**
     * @return Information about configured PVCs (data and/or queue) for server pods.
     */
    public final PvcsInfo getPvcInfo() {
        try {
            return this.defaultApi.getPvcsInfoWithHttpInfo().getData();
        } catch (ApiException exception) {
            LOGGER.error("Error occurred while trying to get PVC information", exception);
        }
        return new PvcsInfo();
    }

}
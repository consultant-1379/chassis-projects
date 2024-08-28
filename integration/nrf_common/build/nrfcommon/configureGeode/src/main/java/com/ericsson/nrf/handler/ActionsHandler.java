package com.ericsson.nrf.handler;

import com.ericsson.adp.kvdbag.adminmgrapi.ApiException;
import com.ericsson.adp.kvdbag.adminmgrapi.ApiResponse;
import com.ericsson.adp.kvdbag.adminmgrapi.model.Action;
import com.ericsson.adp.kvdbag.adminmgrapi.model.ActionCategory;
import com.ericsson.adp.kvdbag.adminmgrapi.model.ActionId;
import com.ericsson.adp.kvdbag.adminmgrapi.model.ExecutionStatus;
import com.ericsson.nrf.common.Constants;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.time.OffsetDateTime;


/**
 * Helper for operations with Actions.
 */
public class ActionsHandler {

    /**
     * Timeout for actions seconds.
     */
    public static final Integer ACTION_TIMEOUT = 600;
    /**
     * Value allowed for parameter targetMembers.
     */
    public static final String ALL_TARGET_MEMBERS = "ALL";
    /**
     * Using this constant as targetMembers of an action will target all locators.
     */
    public static final String LOCATORS_TARGET_MEMBERS = "LOCATORS";
    /**
     * Logger.
     */
    private static final Logger LOGGER = LogManager.getLogger(ActionsHandler.class);
    /**
     * Admin Mgr Client.
     */
    private AdminMgrClient adminMgrClient;

    /**
     * Constructor.
     *
     * @param adminMgrClient REST client for interacting with admin-mgr process.
     */
    public ActionsHandler(final AdminMgrClient adminMgrClient) {
        this.adminMgrClient = adminMgrClient;
    }

    /**
     * Execute START action on all locators.
     *
     * @return action output
     */
    public final String executeStartActionOnAllLocators() {
        LOGGER.info("Executing START action on all locators");
        Action action = new Action();

        action.setCategory(ActionCategory.START);
        action.setTargetMembers(LOCATORS_TARGET_MEMBERS);

        return this.executeAction(action);
    }

    /**
     * Execute START action on all members.
     *
     * @return action output
     */
    public final String executeStartActionOnAllMembers() {
        LOGGER.info("Executing START action on all members");
        Action action = new Action();

        action.setCategory(ActionCategory.START);
        action.setTargetMembers(ALL_TARGET_MEMBERS);

        return this.executeAction(action);
    }

    /**
     * Schedules the given action for execution and waits for it to complete execution.
     *
     * @param action action to be executed
     * @return executed action
     */
    private Action scheduleAndWait(final Action action) {
        LOGGER.info("Executing Action...");
        Action updatedAction = null;
        try {
            ApiResponse<ActionId> actionIdApiResponse = this.adminMgrClient.getDefaultApi()
                    .executeActionWithHttpInfo(action);
            ActionId actionId = actionIdApiResponse.getData();
            boolean executed = false;

            do {
                LOGGER.info("Waiting for " + action.getCategory() + " action to be executed");
                Thread.sleep(AdminMgrClient.MILLISECONDS_BETWEEN_CHECKS);

                ApiResponse<Action> actionApiResponse =
                        this.adminMgrClient.getDefaultApi().getActionInfoWithHttpInfo(actionId.getActionId());
                updatedAction = actionApiResponse.getData();

                long secondsSinceScheduling =
                        OffsetDateTime.now().toEpochSecond()
                                - updatedAction.getReceivedTimestamp().toEpochSecond();
                LOGGER.info("Seconds since scheduling command: " + secondsSinceScheduling);

                if (updatedAction.getExecutionStatus() == ExecutionStatus.EXECUTED) {
                    executed = true;
                }

                if (secondsSinceScheduling > ACTION_TIMEOUT && !executed) {
                    LOGGER.error("Timeout occurred while executing action " + action.getCategory()
                            + " on members " + action.getTargetMembers().toString());
                    break;
                }
            } while (!executed);

        } catch (InterruptedException exception) {
            Thread.currentThread().interrupt();
            LOGGER.error("Execution interrupted!", exception);
        } catch (ApiException exception) {
            LOGGER.error("Error occurred while executing " + action.getCategory() + " action."
                    + " Exception: " + exception, exception);
        }
        return updatedAction;
    }

    /**
     * Schedule and wait for action result. Check action result.
     *
     * @param action action to execute
     * @return action output
     */
    private String executeAction(final Action action) {
        Action resultAction = this.scheduleAndWait(action);
        if (resultAction.getStatusCode() != Constants.ADM_MGR_OK_EXEC_STATUS_CODE) {
            LOGGER.error("The execution of action " + resultAction.getCategory() + " on members "
                    + resultAction.getTargetMembers().toString() + " failed. Error: "
                    + resultAction.getOutput());
        }

        return resultAction.getOutput();
    }

}

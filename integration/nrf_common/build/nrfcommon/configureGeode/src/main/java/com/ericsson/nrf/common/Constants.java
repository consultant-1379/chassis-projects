package com.ericsson.nrf.common;

public class Constants {
    /**
     * Milliseconds before create region check.
     */
    public static final int MILLISECONDS_AFTER_CREATE_REGION = 2000;

    /**
     * Full path of app jars directory.
     */
    public static final String JARS_LOCATION = "/bin/";

    /**
     * Name of Partitioner Resolver jar file.
     */
    public static final String PARTITIONER_JAR_FILE = "kvdb-1.0.jar";

    /**
     * Status code of successfully executed Admin Manager API command.
     */
    public static final int ADM_MGR_OK_EXEC_STATUS_CODE = 0;

    /**
     * Default hidden constructor.
     */
    private Constants() {
    }
}

package ericsson.core.nrf.dbproxy.config;

import java.util.List;
import java.util.ArrayList;
import java.util.Arrays;
import java.io.File;
import java.io.IOException;
import java.net.InetAddress;
import org.json.JSONObject;
import org.json.JSONException;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.commons.io.FileUtils;

import ericsson.core.nrf.dbproxy.common.Code;

public class GeodeConfig
{

    private static final Logger logger = LogManager.getLogger(GeodeConfig.class);

    private static final String LOG_1 = ", new name = ";
    private String kvdb_locator_name;
    private List<String> kvdb_locator_ip_list;
    private int kvdb_locator_port;
    private String[] region_list;

    private String pool_name;
    private int free_connection_timeout;
    private int idle_timeout;
    private int load_conditioning_interval;
    private int max_connections;
    private int min_connections;
    private int ping_interval;
    private boolean pr_single_hop_enabled;
    private int read_timeout;
    private int retry_attempts;
    private int socket_buffer_size;
    private int socket_connect_timeout;
    private boolean subscription_enabled;
    private boolean thread_local_connections;
    private static int max_transmit_fragment_size;
    private static int remote_cache_put_count;
    private static int remote_cache_clear_interval;
    public GeodeConfig()
    {
        kvdb_locator_ip_list = new ArrayList<>();
        kvdb_locator_port = 0;
        kvdb_locator_name = "";
        resetInternalValues();
    }

    private void resetInternalValues()
    {
        pool_name = "connection-pool";
        free_connection_timeout = 3000;
        idle_timeout = 5000;
        load_conditioning_interval = 120000;
        max_connections = 5000;
        min_connections = 10;
        ping_interval = 10000;
        pr_single_hop_enabled = false;
        read_timeout = 3000;
        retry_attempts = -1;
        socket_buffer_size = 32768;
        socket_connect_timeout = 3000;
        subscription_enabled = true;
        thread_local_connections = false;
        max_transmit_fragment_size = 3145728;
        remote_cache_clear_interval = 1000;
        remote_cache_put_count = 100;
    }

    public boolean initialize()
    {
        try {
            String db_info_conf = EnvironmentConfig.getInstance().getDBInfoConf();
            String db_content= FileUtils.readFileToString(new File(db_info_conf),"UTF-8");
            JSONObject db_jsonObject = new JSONObject(db_content);

            kvdb_locator_name = db_jsonObject.getString("locator-server-name");
            String locator_port = db_jsonObject.getString("locator-server-port");

            String region_names = db_jsonObject.getString("region-names");
            region_list = region_names.split(";");

            InetAddress[] address = InetAddress.getAllByName(kvdb_locator_name);
            for(int i = 0; i < address.length; i++) {
                kvdb_locator_ip_list.add(address[i].getHostAddress());
            }

            kvdb_locator_port = Integer.parseInt(locator_port);

            String internal_conf = EnvironmentConfig.getInstance().getInternalConf();
            String content= FileUtils.readFileToString(new File(internal_conf),"UTF-8");
            JSONObject jsonObject = new JSONObject(content);

            pool_name = jsonObject.getString("name");
            free_connection_timeout = Integer.parseInt(jsonObject.getString("free-connection-timeout"));
            idle_timeout = Integer.parseInt(jsonObject.getString("idle-timeout"));
            load_conditioning_interval = Integer.parseInt(jsonObject.getString("load-conditioning-interval"));
            max_connections = Integer.parseInt(jsonObject.getString("max-connections"));
            min_connections = Integer.parseInt(jsonObject.getString("min-connections"));
            ping_interval = Integer.parseInt(jsonObject.getString("ping-interval"));
            pr_single_hop_enabled = ((jsonObject.getString("pr-single-hop-enabled")).equals("true"));
            read_timeout = Integer.parseInt(jsonObject.getString("read-timeout"));
            retry_attempts = Integer.parseInt(jsonObject.getString("retry-attempts"));
            socket_buffer_size = Integer.parseInt(jsonObject.getString("socket-buffer-size"));
            socket_connect_timeout = Integer.parseInt(jsonObject.getString("socket-connect-timeout"));
            subscription_enabled = ((jsonObject.getString("subscription-enabled")).equals("true"));
            thread_local_connections = ((jsonObject.getString("thread-local-connections")).equals("true"));
            max_transmit_fragment_size = Integer.parseInt(jsonObject.getString("max-transmit-fragment-size"));
            remote_cache_put_count = Integer.parseInt(jsonObject.getString("remote-cache-put-threshold"));
            remote_cache_clear_interval = Integer.parseInt(jsonObject.getString("remote-cache-clear-interval"));

        } catch (IOException | JSONException | NullPointerException | SecurityException | NumberFormatException e) {
            logger.error(e.toString());
            logger.warn("Fail to parse " + EnvironmentConfig.getInstance().getInternalConf());
            resetInternalValues();
        } catch (Exception e) {
            logger.error(e.toString());
            logger.warn("Fail to parse " + EnvironmentConfig.getInstance().getInternalConf());
            resetInternalValues();
        }

        logger.debug(toString());

	return validate();
    }

    private boolean validate() {
	
	if(kvdb_locator_name.isEmpty())
	{
	    logger.error("kvdb_locator_name is empty");
	    return false;
	}

	if(kvdb_locator_ip_list.size() == 0)
	{
	    logger.error("No available kvdb locator ip");
	    return false;
	}

	if(kvdb_locator_port <= 0)
	{
	    logger.error("kvdb_locator_port " + Integer.toString(kvdb_locator_port) + " is invalid");
	    return false;
	}

        if(region_list.length == 0)
        {
            logger.error("region_list is empty");
            return false;
        }


	logger.trace("Geode Configuration Validation is Successful");
	return true;
    }

    public String getLocatorName()
    {
        return kvdb_locator_name;
    }

    public void addLocatorIP(String ip)
    {
        kvdb_locator_ip_list.add(ip);
    }

    public List<String> getLocatorIPList()
    {
        return kvdb_locator_ip_list;
    }

    public void setLocatorPort(int port)
    {
        kvdb_locator_port = port;
    }

    public int getLocatorPort()
    {
        return kvdb_locator_port;
    }

    public String[] getRegionList()
    {
        return region_list;
    }

    public String getPoolName()
    {
        return pool_name;
    }

    public int getFreeConnectionTimeout()
    {
        return free_connection_timeout;
    }

    public int getIdleTimeout()
    {
        return idle_timeout;
    }

    public int getLoadConditioningInterval()
    {
        return load_conditioning_interval;
    }

    public int getMaxConnections()
    {
        return max_connections;
    }

    public static int getRemoteCachePutCount(){
        return remote_cache_put_count;
    }

    public static int getRemoteCacheClearInterval(){
        return remote_cache_clear_interval;
    }

    public int getMinConnections()
    {
        return min_connections;
    }

    public int getPingInterval()
    {
        return ping_interval;
    }

    public boolean isPrSingleHopEnabled()
    {
        return pr_single_hop_enabled;
    }

    public int getReadTimeout()
    {
        return read_timeout;
    }

    public int getRetryAttempts()
    {
        return retry_attempts;
    }

    public int getSocketBufferSize()
    {
        return socket_buffer_size;
    }

    public int getSocketConnectTimeout()
    {
        return socket_connect_timeout;
    }

    public boolean isSubscriptionEnabled()
    {
        return subscription_enabled;
    }

    public boolean isThreadLocalConnections()
    {
        return thread_local_connections;
    }

    public static int getMaxTransmitFragmentSize() {
        return max_transmit_fragment_size;
    }

    public String toString()
    {
        StringBuilder sb = new StringBuilder("");
        sb.append("kvdb_locator_name = " + kvdb_locator_name + kvdb_locator_ip_list.toString() + ",");
        sb.append("kvdb_locator_port = " + Integer.toString(kvdb_locator_port) + ",");
        sb.append("region_name = " + Arrays.toString(region_list) + ",");
        sb.append("pool_name = " + pool_name + ",");
        sb.append("free_connection_timeout = " + Integer.toString(free_connection_timeout) + ",");
        sb.append("idle_timeout = " +  Integer.toString(idle_timeout) + ",");
        sb.append("load_conditioning_interval = " + Integer.toString(load_conditioning_interval) + ",");
        sb.append("max_connections = " + Integer.toString(max_connections) + ",");
        sb.append("min_connections = " + Integer.toString(min_connections) + ",");
        sb.append("ping_interval = " + Integer.toString(ping_interval) + ",");
        sb.append("pr_single_hop_enabled = " + pr_single_hop_enabled + ",");
        sb.append("read_timeout = " + Integer.toString(read_timeout) + ",");
        sb.append("retry_attempts = " + Integer.toString(retry_attempts) + ",");
        sb.append("socket_buffer_size = " + Integer.toString(socket_buffer_size) + ",");
        sb.append("socket_connect_timeout = " + Integer.toString(socket_connect_timeout) + ",");
        sb.append("subscription_enabled = " + subscription_enabled + ",");
        sb.append("thread_local_connections = " + thread_local_connections + ",");
        sb.append("max_transmit_fragment_size = " + max_transmit_fragment_size + ",");
        sb.append("remote_cache_put_threshold = " + remote_cache_put_count + ",");
        sb.append("remote_cache_clear_interval = " + remote_cache_clear_interval);

        return sb.toString();
    }

    public int compare(GeodeConfig geode_config)
    {

        if(kvdb_locator_name.compareTo(geode_config.getLocatorName()) != 0) {
            logger.warn("KVDB locator name is changed, old name = " + kvdb_locator_name + LOG_1 + geode_config.getLocatorName());
            return Code.KVDB_LOCATOR_NAME_CHANGED;
        }

        if(kvdb_locator_port != geode_config.getLocatorPort()) {
            logger.warn("KVDB locator port is changed, old port = " + Integer.toString(kvdb_locator_port) + ", new port = " + Integer.toString(geode_config.getLocatorPort()));
            return Code.KVDB_LOCATOR_PORT_CHANGED;
        }

	if (!Arrays.equals(region_list, geode_config.getRegionList())) {
            logger.warn("region list is changed, old = " + region_list + LOG_1 + geode_config.getRegionList());
            return Code.KVDB_REGION_NAME_CHANGED;
        }


        boolean found = false;
        for(String locator_ip : kvdb_locator_ip_list) {
            if(geode_config.getLocatorIPList().contains(locator_ip)) {
                found = true;
                break;
            }
        }

        if(found == false) {
            logger.warn("KVDB locator ip is changed, old ip list = " + kvdb_locator_ip_list.toString() + ", new ip list = " + geode_config.getLocatorIPList().toString());
            return Code.KVDB_LOCATOR_IP_CHANGED;
        }

        logger.debug("Geode configuration is same as before");
        return Code.NOT_CHANGED;
    }
}

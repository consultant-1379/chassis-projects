package ericsson.core.nrf.dbproxy.common;

import java.util.UUID;
import java.util.Random;
import java.util.concurrent.ConcurrentHashMap;
import java.util.Map.Entry;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.DBProxyServer;

public class FragmentSessionManagement
{
    private static final Logger logger = LogManager.getLogger(FragmentSessionManagement.class);

    private static FragmentSessionManagement instance;
    private static int MAX_RETRY_COUNT = 10;

    private ConcurrentHashMap<String, FragmentResult> fragment_result_map;
    private Random random;

    private FragmentSessionManagement()
    {
        fragment_result_map = new ConcurrentHashMap<String, FragmentResult>();
        random = new Random(UUID.randomUUID().hashCode());
    }

    public static synchronized FragmentSessionManagement getInstance()
    {
        if(null == instance) {
            instance = new FragmentSessionManagement();
        }
        return instance;
    }

    public void initialize()
    {
        Thread monitor = new Thread(() -> {

            while(DBProxyServer.getInstance().isRunning()) {

                FragmentSessionManagement.getInstance().validate();

                try {
                    Thread.sleep(1000);
                } catch (Exception e) {
                    logger.error(e.toString());
                }
            }

            logger.trace("Fragment Session Management Thread exits now");
        });

        monitor.start();
    }

    public void validate()
    {

        for(Entry<String, FragmentResult> entry : fragment_result_map.entrySet()) {

            FragmentResult fragment_result = entry.getValue();

            if(fragment_result.getExpiredTime() <= System.currentTimeMillis()) {
                String fragment_session_id = entry.getKey();
                fragment_result_map.remove(fragment_session_id);
                logger.warn("Fragment Session expires with fragment_session_id = {}", fragment_session_id);
            }
        }
    }

    public boolean put(FragmentResult fragment_result, int firstTransmitNum) {
        if (null == fragment_result) {
            return false;
        }

        if (fragment_result.getItems().size() < firstTransmitNum) {
            return false;
        }

        String fragment_session_id = "" + random.nextLong();

        int try_count = 0;
        while (try_count < MAX_RETRY_COUNT && fragment_result_map.containsKey(fragment_session_id)) {
            logger.warn("duplicate fragment_session_id = {}", fragment_session_id);
            fragment_session_id = "" + random.nextLong();
            try_count++;
        }

        if (try_count == MAX_RETRY_COUNT) {
            logger.error("Fail to create the correct fragment session id");
            return false;
        }

        fragment_result.setFragmentSessionID(fragment_session_id);
        fragment_result.setTotalNumber(fragment_result.getItems().size());
        fragment_result.setTransmittedNumber(firstTransmitNum);
        fragment_result.setExpiredTime(System.currentTimeMillis() + Code.FRAGMENT_SESSION_ACTIVE_TIME);
        fragment_result_map.put(fragment_session_id, fragment_result);

        return true;
    }

    public ExecutionResult get(String region_name, String fragment_session_id)
    {
        FragmentResult item = fragment_result_map.get(fragment_session_id);

        if(null == item) {
            logger.debug("No data found by fragment_session_id = {}", fragment_session_id);
            return new ExecutionResult(Code.DATA_NOT_EXIST);
        }


        int total_number = item.getTotalNumber();
        int transmitted_number = item.getTransmittedNumber();
        int block_number = FragmentUtil.transmitNumPerTime(item, region_name);
        FragmentResult fragment_result = new FragmentResult();

        if(transmitted_number + block_number <= total_number) {

            item.setTransmittedNumber(transmitted_number + block_number);
            item.setExpiredTime(System.currentTimeMillis() + Code.FRAGMENT_SESSION_ACTIVE_TIME);
            fragment_result.addAll(item.getItems().subList(transmitted_number, transmitted_number + block_number));

        } else {

            item.setTransmittedNumber(total_number);
            fragment_result.addAll(item.getItems().subList(transmitted_number, total_number));

            fragment_result_map.remove(fragment_session_id);
        }

        fragment_result.setFragmentSessionID(fragment_session_id);
        fragment_result.setTotalNumber(item.getTotalNumber());
        fragment_result.setTransmittedNumber(item.getTransmittedNumber());

        return fragment_result;
    }
}

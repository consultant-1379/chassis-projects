package ericsson.core.nrf.dbproxy.common;

import ericsson.core.nrf.dbproxy.DBProxyServer;
import java.util.Map.Entry;
import java.util.Random;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class FragmentSessionManagement {

  private static final Logger LOGGER = LogManager.getLogger(FragmentSessionManagement.class);

  private static FragmentSessionManagement instance;
  private static int maxRetryCount = 10;

  private ConcurrentHashMap<String, FragmentResult> fragment_result_map;
  private Random random;

  private FragmentSessionManagement() {
    fragment_result_map = new ConcurrentHashMap<String, FragmentResult>();
    random = new Random(UUID.randomUUID().hashCode());
  }

  public static synchronized FragmentSessionManagement getInstance() {
    if (null == instance) {
      instance = new FragmentSessionManagement();
    }
    return instance;
  }

  public void initialize() {
    Thread monitor = new Thread(() -> {

      while (DBProxyServer.getInstance().isRunning()) {

        FragmentSessionManagement.getInstance().validate();

        try {
          Thread.sleep(1000);
        } catch (Exception e) {
          LOGGER.error(e.toString());
        }
      }

      LOGGER.debug("Fragment Session Management Thread exits now");
    });

    monitor.start();
  }

  public void validate() {

    for (Entry<String, FragmentResult> entry : fragment_result_map.entrySet()) {

      FragmentResult fragmentResult = entry.getValue();

      if (fragmentResult.getExpiredTime() <= System.currentTimeMillis()) {
        String fragmentSessionId = entry.getKey();
        fragment_result_map.remove(fragmentSessionId);
        LOGGER.warn("Fragment Session expires with fragmentSessionId = {}", fragmentSessionId);
      }
    }
  }

  public boolean put(FragmentResult fragmentResult, int firstTransmitNum) {
    if (null == fragmentResult) {
      return false;
    }

    if (fragmentResult.getItems().size() < firstTransmitNum) {
      return false;
    }

    String fragmentSessionId = "" + random.nextLong();

    int tryCount = 0;
    while (tryCount < maxRetryCount && fragment_result_map.containsKey(fragmentSessionId)) {
      LOGGER.warn("duplicate fragment_session_id = {}", fragmentSessionId);
      fragmentSessionId = "" + random.nextLong();
      tryCount++;
    }

    if (tryCount == maxRetryCount) {
      LOGGER.error("Fail to create the correct fragment session id");
      return false;
    }

    fragmentResult.setFragmentSessionID(fragmentSessionId);
    fragmentResult.setTotalNumber(fragmentResult.getItems().size());
    fragmentResult.setTransmittedNumber(firstTransmitNum);
    fragmentResult.setExpiredTime(System.currentTimeMillis() + Code.FRAGMENT_SESSION_ACTIVE_TIME);
    fragment_result_map.put(fragmentSessionId, fragmentResult);

    return true;
  }

  public ExecutionResult get(String regionName, String fragmentSessionId) {
    FragmentResult item = fragment_result_map.get(fragmentSessionId);

    if (null == item) {
      LOGGER.debug("No data found by fragmentSessionId = {}", fragmentSessionId);
      return new ExecutionResult(Code.DATA_NOT_EXIST);
    }

    int totalNumber = item.getTotalNumber();
    int transmittedNumber = item.getTransmittedNumber();
    int blockNumber = FragmentUtil.transmitNumPerTime(item, regionName);
    FragmentResult fragmentResult = new FragmentResult();

    if (transmittedNumber + blockNumber <= totalNumber) {

      item.setTransmittedNumber(transmittedNumber + blockNumber);
      item.setExpiredTime(System.currentTimeMillis() + Code.FRAGMENT_SESSION_ACTIVE_TIME);
      fragmentResult
          .addAll(item.getItems().subList(transmittedNumber, transmittedNumber + blockNumber));

    } else {

      item.setTransmittedNumber(totalNumber);
      fragmentResult.addAll(item.getItems().subList(transmittedNumber, totalNumber));

      fragment_result_map.remove(fragmentSessionId);
    }

    fragmentResult.setFragmentSessionID(fragmentSessionId);
    fragmentResult.setTotalNumber(item.getTotalNumber());
    fragmentResult.setTransmittedNumber(item.getTransmittedNumber());

    return fragmentResult;
  }
}

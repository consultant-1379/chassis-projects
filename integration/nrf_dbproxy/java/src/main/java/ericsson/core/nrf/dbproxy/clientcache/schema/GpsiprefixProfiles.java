package ericsson.core.nrf.dbproxy.clientcache.schema;

import java.util.HashMap;

public class GpsiprefixProfiles {

  private Long gpsi_prefix;
  private HashMap value_info;

  public GpsiprefixProfiles() {
    this.gpsi_prefix = 0L;
    this.value_info = new HashMap();
  }

  public Long getGpsiprefix() {
    return gpsi_prefix;
  }

  public void setGpsiprefix(Long gpsiPrefix) {
    this.gpsi_prefix = gpsiPrefix;
  }

  public void addValueInfo(String valueInfo) {
    this.value_info.put(valueInfo, 1);
  }

  public void rmValueInfo(String valueInfo) {
    this.value_info.remove(valueInfo);
  }

  public HashMap getValueInfo() {
    return this.value_info;
  }

  public String toString() {
    return "GpsiprefixProfiles{" +
        "gpsi_prefix=" + gpsi_prefix.toString() +
        ", value_info=" + value_info.toString() +
        '}';
  }
}

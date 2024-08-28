package ericsson.core.nrf.dbproxy.clientcache.schema;

import java.util.HashMap;

public class ImsiprefixProfiles {

  private Long imsi_prefix;
  private HashMap value_info;

  public ImsiprefixProfiles() {
    this.imsi_prefix = 0L;
    this.value_info = new HashMap();
  }

  public Long getImsiprefix() {
    return imsi_prefix;
  }

  public void setImsiprefix(Long imsiPrefix) {
    this.imsi_prefix = imsiPrefix;
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
    return "ImsiprefixProfiles{" +
        "imsi_prefix=" + imsi_prefix.toString() +
        ", value_info=" + value_info.toString() +
        '}';
  }
}

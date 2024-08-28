package com.ericsson.geode.functionservice;

public class CacheSort {
    private String key;
    private long put_time;
    private long expiry_time;

    public CacheSort(long put_time, long expiry_time, String key){
      this.put_time = put_time;
      this.expiry_time = expiry_time;
      this.key = key;
    }

  public long getPut_time() {
    return put_time;
  }

  public void setPut_time(long put_time) {
    this.put_time = put_time;
  }

  public long getExpiry_time() {
    return expiry_time;
  }

  public void setExpiry_time(long expiry_time) {
    this.expiry_time = expiry_time;
  }

  public String getKey() {
    return key;
  }

  public void setKey(String key) {
    this.key = key;
  }
}

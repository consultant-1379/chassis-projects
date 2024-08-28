package ericsson.core.nrf.dbproxy.clientcache.schema;

import com.google.protobuf.ByteString;

public class CacheNFProfile {

  private ByteString raw_data;
  private long expiry_time;
  private long put_time;

  public CacheNFProfile() {
    this.expiry_time = 0;
    this.put_time = 0;
  }

  public ByteString getRaw_data() {
    return raw_data;
  }

  public void setRaw_data(ByteString rawData) {
    this.raw_data = rawData;
  }

  public long getExpiry_time() {
    return expiry_time;
  }

  public void setExpiry_time(long expiryTime) {
    this.expiry_time = expiryTime;
  }

  public long getPut_time() {
    return put_time;
  }

  public void setPut_time(long putTime) {
    this.put_time = putTime;
  }

  @Override
  public String toString() {
    return "CacheNFProfile{" +
        "raw_data=" + raw_data.toString() +
        ", expiry_time=" + expiry_time +
        ", put_time=" + put_time +
        '}';
  }
}

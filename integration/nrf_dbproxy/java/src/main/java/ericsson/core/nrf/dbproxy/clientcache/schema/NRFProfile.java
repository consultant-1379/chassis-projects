package ericsson.core.nrf.dbproxy.clientcache.schema;

import com.google.protobuf.ByteString;
import java.util.HashMap;

public class NRFProfile {

  private ByteString raw_data;

  private long key1;
  private long key3;

  private HashMap amf_key1;
  private HashMap amf_key2;
  private HashMap amf_key3;
  private HashMap amf_key4;
  private HashMap smf_key1;
  private HashMap smf_key2;
  private HashMap smf_key3;
  private HashMap udm_key1;
  private HashMap udm_key2;
  private HashMap ausf_key1;
  private HashMap ausf_key2;
  private HashMap pcf_key1;
  private HashMap pcf_key2;


  public NRFProfile() {
    key1 = 0;
    key3 = 0;
    amf_key1 = new HashMap();
    amf_key2 = new HashMap();
    amf_key3 = new HashMap();
    amf_key4 = new HashMap();
    smf_key1 = new HashMap();
    smf_key2 = new HashMap();
    smf_key3 = new HashMap();
    udm_key1 = new HashMap();
    udm_key2 = new HashMap();
    ausf_key1 = new HashMap();
    ausf_key2 = new HashMap();
    pcf_key1 = new HashMap();
    pcf_key2 = new HashMap();
  }

  public ByteString getRaw_data() {
    return raw_data;
  }

  public void setRaw_data(ByteString rawData) {
    this.raw_data = rawData;
  }

  public synchronized void setAMFKey1(HashMap map) {
    amf_key1.putAll(map);
  }

  public synchronized void setAMFKey2(HashMap map) {
    amf_key2.putAll(map);
  }

  public synchronized void setAMFKey3(HashMap map) {
    amf_key3.putAll(map);
  }

  public synchronized void setAMFKey4(HashMap map) {
    amf_key4.putAll(map);
  }

  public synchronized void setSMFKey1(HashMap map) {
    smf_key1.putAll(map);
  }

  public synchronized void setSMFKey2(HashMap map) {
    smf_key2.putAll(map);
  }

  public synchronized void setSMFKey3(HashMap map) {
    smf_key3.putAll(map);
  }

  public synchronized void setUDMKey1(HashMap map) {
    udm_key1.putAll(map);
  }

  public synchronized void setUDMKey2(HashMap map) {
    udm_key2.putAll(map);
  }

  public synchronized void setAUSFKey1(HashMap map) {
    ausf_key1.putAll(map);
  }

  public synchronized void setAUSFKey2(HashMap map) {
    ausf_key2.putAll(map);
  }

  public synchronized void setPCFKey1(HashMap map) {
    pcf_key1.putAll(map);
  }

  public synchronized void setPCFKey2(HashMap map) {
    pcf_key2.putAll(map);
  }

  public long getExpireTime() {
    return key1;
  }

  public void setKey1(long key1) {
    this.key1 = key1;
  }

  public void setKey3(long key3) {
    this.key3 = key3;
  }

  public long getProvFlag() {
    return key3;
  }
}

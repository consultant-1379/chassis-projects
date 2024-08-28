package ericsson.core.nrf.dbproxy.clientcache.schema;

import java.util.HashMap;
import java.util.Iterator;

public class KeyAggregationCollection {

  private String sub_key1;
  private String sub_key2;
  private HashMap sub_key3;

  public KeyAggregationCollection() {

    sub_key1 = sub_key2 = "";
    sub_key3 = new HashMap();
  }

  public void setSubKey1(String value) {
    sub_key1 = value;
  }

  public void setSubKey2(String value) {
    sub_key2 = value;
  }

  public synchronized void addIntoSubKey3List(int id, KeyAggregation ka) {
    sub_key3.put(id, ka);
  }


  public String toString() {

    StringBuilder sb = new StringBuilder("");

    sb.append("sub_key1 = {" + sub_key1 + "}, ");
    sb.append("sub_key2 = {" + sub_key2 + "}, ");

    Iterator iter = sub_key3.entrySet().iterator();
    while (iter.hasNext()) {
      HashMap.Entry entry = (HashMap.Entry) iter.next();
      KeyAggregation ka = (KeyAggregation) entry.getValue();
      sb.append("sub_key3 = {" + ka.toString() + "}, ");
    }

    return sb.toString();
  }

}

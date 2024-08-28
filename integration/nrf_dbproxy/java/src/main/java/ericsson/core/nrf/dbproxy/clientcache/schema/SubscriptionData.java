package ericsson.core.nrf.dbproxy.clientcache.schema;

import com.google.protobuf.ByteString;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;

public class SubscriptionData {

  private String subscription_id;
  private ByteString data;
  private String noCond;
  private String nfStatusNotificationUri;
  private String nfInstanceId;
  private String nfType;
  private String serviceName;

  private HashMap amfCond;
  private HashMap guamiList;
  private HashMap snssaiList;

  private List<String> nsiList;

  private HashMap nfGroupCond;

  private long validityTime;

  public SubscriptionData() {
    this.subscription_id = "";
    this.data = ByteString.EMPTY;
    this.noCond = "";
    this.nfStatusNotificationUri = "";
    this.nfInstanceId = "";
    this.nfType = "";
    this.serviceName = "";

    this.amfCond = new HashMap();
    this.guamiList = new HashMap();
    this.snssaiList = new HashMap();

    this.nsiList = new ArrayList<String>();
    this.nfGroupCond = new HashMap();

    this.validityTime = 0;
  }

  public String getSubscriptionID() {

    return this.subscription_id;
  }

  public void setSubscriptionID(String subscriptionId) {

    this.subscription_id = subscriptionId;
  }

  public ByteString getData() {

    return this.data;
  }

  public void setData(ByteString data) {

    this.data = data;
  }

  public String getNoCond() {

    return this.noCond;
  }

  public void setNoCond(String noCond) {

    this.noCond = noCond;
  }

  public String getNfStatusNotificationUri() {

    return this.nfStatusNotificationUri;
  }

  public void setNfStatusNotificationUri(String nfStatusNotificationUri) {

    this.nfStatusNotificationUri = nfStatusNotificationUri;
  }

  public String getNfInstanceId() {

    return this.nfInstanceId;
  }

  public void setNfInstanceId(String nfInstanceId) {

    this.nfInstanceId = nfInstanceId;
  }

  public String getNfType() {

    return this.nfType;
  }

  public void setNfType(String nfType) {

    this.nfType = nfType;
  }

  public void setServiceName(String serviceName) {

    this.serviceName = serviceName;
  }

  public synchronized void addAmfCond(int id, KeyAggregation ka) {
    amfCond.put(id, ka);
  }

  public synchronized void addGuamiList(int id, KeyAggregation ka) {
    guamiList.put(id, ka);
  }

  public synchronized void addSnssaiList(int id, KeyAggregation ka) {
    snssaiList.put(id, ka);
  }

  public synchronized void addNsiList(String nsi) {
    nsiList.add(nsi);
  }

  public synchronized void addNfGroupCond(int id, KeyAggregation ka) {
    nfGroupCond.put(id, ka);
  }

  public void setValidityTime(long value) {
    validityTime = value;
  }

  public String toString() {

    StringBuilder sb = new StringBuilder("");
    sb.append("subscription_id = {" + this.subscription_id + "}, ");
    sb.append("data = {" + this.data.toString() + "}, ");
    sb.append("noCond = {" + this.noCond + "}, ");
    sb.append("nfStatusNotificationUri = {" + this.nfStatusNotificationUri + "}, ");
    sb.append("nfInstanceId = {" + this.nfInstanceId + "}, ");
    sb.append("nfType = {" + this.nfType + "}, ");
    sb.append("serviceName = {" + this.serviceName + "}, ");

    Iterator iter = amfCond.entrySet().iterator();
    while (iter.hasNext()) {
      HashMap.Entry entry = (HashMap.Entry) iter.next();
      KeyAggregation ka = (KeyAggregation) entry.getValue();
      sb.append("amfCond = {" + ka.toString() + "}, ");
    }

    iter = guamiList.entrySet().iterator();
    while (iter.hasNext()) {
      HashMap.Entry entry = (HashMap.Entry) iter.next();
      KeyAggregation ka = (KeyAggregation) entry.getValue();
      sb.append("guamiList = {" + ka.toString() + "}, ");
    }

    iter = snssaiList.entrySet().iterator();
    while (iter.hasNext()) {
      HashMap.Entry entry = (HashMap.Entry) iter.next();
      KeyAggregation ka = (KeyAggregation) entry.getValue();
      sb.append("snssaiList = {" + ka.toString() + "}, ");
    }

    Iterator<String> iter2 = nsiList.iterator();
    while (iter2.hasNext()) {
      String nsi = iter2.next();
      sb.append("nsiList = {" + nsi + "}, ");
    }

    iter = nfGroupCond.entrySet().iterator();
    while (iter.hasNext()) {
      HashMap.Entry entry = (HashMap.Entry) iter.next();
      KeyAggregation ka = (KeyAggregation) entry.getValue();
      sb.append("nfGroupCond = {" + ka.toString() + "}, ");
    }

    sb.append("validityTime = {" + Long.toString(validityTime) + "}, ");

    return sb.toString();
  }
}

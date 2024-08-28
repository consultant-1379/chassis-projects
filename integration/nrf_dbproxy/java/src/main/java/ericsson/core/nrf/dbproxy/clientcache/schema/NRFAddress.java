package ericsson.core.nrf.dbproxy.clientcache.schema;

import com.google.protobuf.ByteString;

public class NRFAddress {

  private String nrf_address_id;
  private ByteString data;
  private String key1;
  private String key2;
  private String key3;
  private String key4;
  private String key5;

  public NRFAddress() {

    this.nrf_address_id = "";
    this.data = ByteString.EMPTY;
    this.key1 = this.key2 = this.key3 = this.key4 = this.key5 = "";
  }

  public String getNRFAddressID() {

    return this.nrf_address_id;
  }

  public void setNRFAddressID(String nrfAddressId) {

    this.nrf_address_id = nrfAddressId;
  }

  public ByteString getData() {

    return this.data;
  }

  public void setData(ByteString data) {

    this.data = data;
  }

  public String getKey1() {

    return this.key1;
  }

  public void setKey1(String key1) {

    this.key1 = key1;
  }

  public String getKey2() {

    return this.key2;
  }

  public void setKey2(String key2) {

    this.key2 = key2;
  }

  public String getKey3() {

    return this.key3;
  }

  public void setKey3(String key3) {

    this.key3 = key3;
  }

  public String getKey4() {

    return this.key4;
  }

  public void setKey4(String key4) {

    this.key4 = key4;
  }

  public String getKey5() {

    return this.key5;
  }

  public void setKey5(String key5) {

    this.key5 = key5;
  }

  public String toString() {

    String str = "nrf_address_id = {" + this.nrf_address_id + "}, ";
    str += "data = {" + this.data.toString() + "}, ";
    str += "key1 = {" + this.key1 + "}, ";
    str += "key2 = {" + this.key2 + "}, ";
    str += "key3 = {" + this.key3 + "}, ";
    str += "key4 = {" + this.key4 + "}, ";
    str += "key5 = {" + this.key5 + "}";

    return str;
  }
}

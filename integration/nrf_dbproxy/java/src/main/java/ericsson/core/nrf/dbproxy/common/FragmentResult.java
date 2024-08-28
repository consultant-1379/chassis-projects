package ericsson.core.nrf.dbproxy.common;

public class FragmentResult extends SearchResult {

  private String fragment_session_id;
  private int total_number;
  private int transmitted_number;
  private long result_expired_time;

  public FragmentResult() {
    super(true);
    fragment_session_id = "";
    total_number = 0;
    transmitted_number = 0;
    result_expired_time = 0;
  }

  public String getFragmentSessionID() {
    return fragment_session_id;
  }

  public void setFragmentSessionID(String value) {
    fragment_session_id = value;
  }

  public int getTotalNumber() {
    return total_number;
  }

  public void setTotalNumber(int value) {
    total_number = value;
  }

  public int getTransmittedNumber() {
    return transmitted_number;
  }

  public void setTransmittedNumber(int value) {
    transmitted_number = value;
  }

  public long getExpiredTime() {
    return result_expired_time;
  }

  public void setExpiredTime(long value) {
    result_expired_time = value;
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder("");
    sb.append("code = " + Integer.toString(code) + ",");
    sb.append("isFragmented = " + isFragmented + ",");
    sb.append("result size = " + Integer.toString(items.size()));
    sb.append("fragment_session_id = " + fragment_session_id + ",");
    sb.append("total_number = " + total_number + ",");
    sb.append("transmitted_number = " + Integer.toString(transmitted_number));
    return sb.toString();
  }
}

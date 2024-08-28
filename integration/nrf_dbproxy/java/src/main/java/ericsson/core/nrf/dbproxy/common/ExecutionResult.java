package ericsson.core.nrf.dbproxy.common;

public class ExecutionResult {

  protected int code;

  public ExecutionResult(int value) {
    code = value;
  }

  public int getCode() {
    return code;
  }

  public void setCode(int value) {
    code = value;
  }

  public String toString() {
    StringBuilder sb = new StringBuilder("");
    sb.append("code = " + Integer.toString(code));
    return sb.toString();
  }

}

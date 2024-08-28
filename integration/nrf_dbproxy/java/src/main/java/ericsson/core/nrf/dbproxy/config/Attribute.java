package ericsson.core.nrf.dbproxy.config;

public class Attribute {

  String parameter;
  String path;
  String from;
  String where;
  boolean exist_check;

  public Attribute(String parameter, String path, String from, String where, boolean existCheck) {
    this.parameter = parameter;
    this.path = path;
    this.from = from;
    this.where = where;
    this.exist_check = existCheck;
  }

  public String getPath() {
    return path;
  }

  public String getFrom() {
    return from;
  }

  public String getWhere() {
    return where;
  }

  public boolean isExistCheck() {
    return exist_check;
  }

  public String toString() {
    return "parameter = [" + parameter + "], path = [" + path + "], from = [" + from
        + "], where = [" + where + "], exist_check = [" + exist_check + "]";
  }
}

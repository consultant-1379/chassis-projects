package ericsson.core.nrf.dbproxy.common;

import java.util.ArrayList;
import java.util.List;

public class SearchResult extends ExecutionResult {

  protected List<Object> items;
  protected boolean isFragmented;

  public SearchResult(boolean isFragmented) {
    super(Code.SUCCESS);
    items = new ArrayList<>();
    this.isFragmented = isFragmented;
  }

  public void add(Object object) {
    items.add(object);
  }

  public void addAll(List<Object> objs) {
    items.addAll(objs);
  }

  public List<Object> getItems() {
    return items;
  }

  public boolean isFragmented() {
    return isFragmented;
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder("");
    sb.append("code = " + Integer.toString(code) + ",");
    sb.append("isFragmented = " + isFragmented + ",");
    sb.append("result size = " + Integer.toString(items.size()));
    return sb.toString();
  }
}

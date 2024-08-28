package ericsson.core.nrf.dbproxy.fsnotify;

public interface EventHandler {

  public static final int FILE_CREATE = 1;
  public static final int FILE_CHANGE = 2;
  public static final int FILE_REMOVE = 3;
  public static final int DIRECTORY_CREATE = 4;
  public static final int DIRECTORY_CHANGE = 5;
  public static final int DIRECTORY_REMOVE = 6;

  public void handle(String eventName, int op);
}

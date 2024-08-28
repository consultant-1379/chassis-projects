package ericsson.core.nrf.dbproxy.fsnotify;

import java.util.concurrent.TimeUnit;
import org.apache.commons.io.monitor.FileAlterationMonitor;
import org.apache.commons.io.monitor.FileAlterationObserver;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class FileMonitor {

  private static final Logger LOGGER = LogManager.getLogger(FileMonitor.class);
  private String filePath;
  private EventHandler eventHandler;
  private int interval;

  public FileMonitor(String filePath, EventHandler eventHandler, int interval) {
    this.filePath = filePath;
    this.eventHandler = eventHandler;
    this.interval = interval;
  }

  public void start() throws Exception {
    LOGGER.info("File Monitor starts, to monitor directory " + this.filePath);
    long checkInterval = TimeUnit.MILLISECONDS.toMillis(this.interval);
    FileAlterationObserver observer = new FileAlterationObserver(this.filePath);
    observer.addListener(new FileListener(this.eventHandler));
    FileAlterationMonitor monitor = new FileAlterationMonitor(checkInterval, observer);
    monitor.start();
  }

}

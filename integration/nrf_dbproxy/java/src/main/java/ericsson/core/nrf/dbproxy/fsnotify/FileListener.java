package ericsson.core.nrf.dbproxy.fsnotify;

import java.io.File;
import org.apache.commons.io.monitor.FileAlterationListenerAdaptor;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

class FileListener extends FileAlterationListenerAdaptor {
  private static final Logger LOGGER = LogManager.getLogger(FileListener.class);
  private EventHandler eventHandler;

  public FileListener(EventHandler eventHandler) {
    super();
    this.eventHandler = eventHandler;
  }

  @Override
  public void onFileCreate(File file) {
    LOGGER.debug("file " + file.getName() + " created");
    eventHandler.handle(file.getName(), EventHandler.FILE_CREATE);
  }

  @Override
  public void onFileChange(File file) {
    LOGGER.debug("file " + file.getName() + " changed");
    eventHandler.handle(file.getName(), EventHandler.FILE_CHANGE);
  }

  @Override
  public void onFileDelete(File file) {
    LOGGER.debug("file " + file.getName() + " removed");
    eventHandler.handle(file.getName(), EventHandler.FILE_REMOVE);
  }

  @Override
  public void onDirectoryCreate(File directory) {
    LOGGER.debug("directory " + directory.getAbsolutePath() + " created");
    eventHandler.handle(directory.getAbsolutePath(), EventHandler.DIRECTORY_CREATE);
  }

  @Override
  public void onDirectoryChange(File directory) {
    LOGGER.debug("directory " + directory.getAbsolutePath() + " changed");
    eventHandler.handle(directory.getAbsolutePath(), EventHandler.DIRECTORY_CHANGE);
  }

  @Override
  public void onDirectoryDelete(File directory) {
    LOGGER.debug("directory " + directory.getAbsolutePath() + " removed");
    eventHandler.handle(directory.getAbsolutePath(), EventHandler.DIRECTORY_REMOVE);
  }

}

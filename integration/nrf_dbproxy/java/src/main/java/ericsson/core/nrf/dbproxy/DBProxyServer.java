package ericsson.core.nrf.dbproxy;

import static java.util.concurrent.ForkJoinPool.defaultForkJoinWorkerThreadFactory;

import com.google.common.util.concurrent.UncaughtExceptionHandlers;
import ericsson.core.nrf.dbproxy.common.FragmentSessionManagement;
import ericsson.core.nrf.dbproxy.config.ConfigLoader;
import ericsson.core.nrf.dbproxy.config.EnvironmentConfig;
import ericsson.core.nrf.dbproxy.functionservice.RemoteCacheClearThread;
import ericsson.core.nrf.dbproxy.functionservice.RemoteCacheMonitorThread;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceImpl;
import ericsson.core.nrf.dbproxy.log.LogLevelControl;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import java.io.File;
import java.io.IOException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.ForkJoinPool;
import java.util.concurrent.ForkJoinWorkerThread;
import java.util.concurrent.atomic.AtomicInteger;
import org.apache.commons.io.FileUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.json.JSONObject;
import sun.misc.Signal;
import sun.misc.SignalHandler;

public class DBProxyServer {

  private static final Logger LOGGER = LogManager.getLogger(DBProxyServer.class);
  private static DBProxyServer instance;
  private Server server;
  private boolean running;

  private DBProxyServer() {
    this.running = true;
  }

  public static synchronized DBProxyServer getInstance() {
    if (null == instance) {
      instance = new DBProxyServer();
    }
    return instance;
  }

  public static void main(String[] args) {
    final DBProxyServer dbProxyServer = DBProxyServer.getInstance();

    dbProxyServer.start();

    dbProxyServer.blockUntilShutdown();

    if (dbProxyServer.isRunning()) {
      dbProxyServer.shutdown();
    }
  }

  public boolean isRunning() {
    return this.running;
  }

  ExecutorService getExecutor(int asyncThreads) {
    return new ForkJoinPool(asyncThreads,
        new ForkJoinPool.ForkJoinWorkerThreadFactory() {
          final AtomicInteger num = new AtomicInteger();

          @Override
          public ForkJoinWorkerThread newThread(ForkJoinPool pool) {
            ForkJoinWorkerThread thread = defaultForkJoinWorkerThreadFactory.newThread(pool);
            thread.setDaemon(true);
            thread.setName("server-worker-" + "-" + num.getAndIncrement());
            return thread;
          }
        }, UncaughtExceptionHandlers.systemExit(), true /* async */);
  }

  private Boolean isAsyncExecutorEnabled() {
    boolean isAsyncExecutorEnabled = true;
    try {
      String internalConf = EnvironmentConfig.getInstance().getInternalConf();
      String content = null;
      content = FileUtils.readFileToString(new File(internalConf), "UTF-8");
      JSONObject jsonObject = new JSONObject(content);
      if (jsonObject.has("async-executor-enabled")) {
        isAsyncExecutorEnabled = jsonObject.getString("async-executor-enabled").equals("true");
      }
    } catch (IOException e) {
      LOGGER.error(e.toString());
    }
    return isAsyncExecutorEnabled;
  }

  private int getAsyncThreadNum() {
    int asyncThreads = Runtime.getRuntime().availableProcessors();
    try {
      String internalConf = EnvironmentConfig.getInstance().getInternalConf();
      String content = null;
      content = FileUtils.readFileToString(new File(internalConf), "UTF-8");
      JSONObject jsonObject = new JSONObject(content);
      if (jsonObject.has("async-threadnum-factor")) {
        String factorstr = jsonObject.getString("async-threadnum-factor");
        double factor = Double.parseDouble(factorstr);
        asyncThreads = (int) (asyncThreads * factor);
      }
    } catch (IOException e) {
      LOGGER.error(e.toString());
    }
    return asyncThreads;
  }

  private void start() {
    LOGGER.debug("DBProxy Server is starting");

    DBProxySignalHandler signalHandler = new DBProxySignalHandler();
    Signal.handle(new Signal("TERM"), signalHandler);
    Signal.handle(new Signal("INT"), signalHandler);

    if (!LogLevelControl.getInstance().initialize()) {
      LOGGER.warn("Fail to initialize log level controller");
    }
    FragmentSessionManagement.getInstance().initialize();
    ConfigLoader.getInstance().start();

    RemoteCacheMonitorThread.getInstance().start();
    RemoteCacheClearThread.getInstance().start();

    boolean listenOk = false;
    int asyncThreadnum = getAsyncThreadNum();
    if (asyncThreadnum <= 0) {
      asyncThreadnum = 1;
    }
    LOGGER.debug("async thread num = " + asyncThreadnum);
    while (!listenOk) {
      try {
        if (isAsyncExecutorEnabled()) {
          LOGGER.debug("enabled async executor");
          server = ServerBuilder.forPort(EnvironmentConfig.getInstance().getGRPCPort())
              .addService(new NFDataManagementServiceImpl())
              .executor(getExecutor(asyncThreadnum))
              .build().start();
        } else {
          LOGGER.debug("disabled async executor");
          server = ServerBuilder.forPort(EnvironmentConfig.getInstance().getGRPCPort())
              .addService(new NFDataManagementServiceImpl())
              .build().start();
        }
        listenOk = true;
      } catch (IOException | IllegalStateException e) {
        LOGGER.error(e.toString());
        LOGGER.error("Fail to listen on port = {}", EnvironmentConfig.getInstance().getGRPCPort());
      }

      if (!listenOk) {
        LOGGER.error("Sleep 5 seconds and try to listen again");
        try {
          Thread.sleep(5000);
        } catch (Exception e) {
          LOGGER.error(e.toString());
        }
      }
    }

    LOGGER.debug(
        "DBProxy Server has been started successfully, listening on port " + EnvironmentConfig
            .getInstance().getGRPCPort());

    Runtime.getRuntime().addShutdownHook(new Thread(DBProxyServer.this::shutdown));

  }

  private void shutdown() {
    this.running = false;
    if (server != null) {
      LOGGER.warn("Shutdown GRPC server");
      server.shutdown();
      server = null;
    }

  }

  private void blockUntilShutdown() {
    if (server != null) {
      try {
        server.awaitTermination();
      } catch (InterruptedException e) {
        LOGGER.error(e.toString());
        LOGGER.error("Exception occurs while waiting for DBProxy Server Termination");
        Thread.currentThread().interrupt();
      }
    }
  }

  class DBProxySignalHandler implements SignalHandler {

    public void handle(Signal signal) {
      LOGGER.warn("Catch signal = " + signal.getName());
      LOGGER.warn("Sleep 10 seconds before shutdown DBProxy Server");
      try {
        Thread.sleep(10000);
      } catch (Exception e) {
        LOGGER.error(e.toString());
      }

      DBProxyServer.getInstance().shutdown();
    }
  }
}

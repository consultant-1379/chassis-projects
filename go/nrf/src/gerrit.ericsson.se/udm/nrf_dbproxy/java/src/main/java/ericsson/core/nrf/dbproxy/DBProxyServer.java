package ericsson.core.nrf.dbproxy;

import com.google.common.util.concurrent.UncaughtExceptionHandlers;
import ericsson.core.nrf.dbproxy.config.GeodeConfig;
import ericsson.core.nrf.dbproxy.functionservice.RemoteCacheClearThread;
import ericsson.core.nrf.dbproxy.functionservice.RemoteCacheMonitorThread;
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

import ericsson.core.nrf.dbproxy.common.FragmentSessionManagement;
import ericsson.core.nrf.dbproxy.config.ConfigLoader;
import ericsson.core.nrf.dbproxy.config.EnvironmentConfig;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceImpl;

import static java.util.concurrent.ForkJoinPool.defaultForkJoinWorkerThreadFactory;

public class DBProxyServer
{
    private static final Logger logger = LogManager.getLogger(DBProxyServer.class);

    private Server server;
    private static DBProxyServer instance;
    private boolean running;

    private DBProxyServer()
    {
        this.running = true;
    }

    public static synchronized DBProxyServer getInstance()
    {
        if(null == instance) {
            instance = new DBProxyServer();
        }
        return instance;
    }

    public boolean isRunning()
    {
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
            String internal_conf = EnvironmentConfig.getInstance().getInternalConf();
            String content= null;
            content = FileUtils.readFileToString(new File(internal_conf),"UTF-8");
            JSONObject jsonObject = new JSONObject(content);
            if (jsonObject.has("async-executor-enabled")) {
                isAsyncExecutorEnabled = jsonObject.getString("async-executor-enabled").equals("true");
            }
        } catch (IOException e) {
            logger.error(e.toString());
        }
        return isAsyncExecutorEnabled;
    }

    private int getAsyncThreadNum() {
        int asyncThreads = Runtime.getRuntime().availableProcessors();
        try {
            String internal_conf = EnvironmentConfig.getInstance().getInternalConf();
            String content= null;
            content = FileUtils.readFileToString(new File(internal_conf),"UTF-8");
            JSONObject jsonObject = new JSONObject(content);
            if (jsonObject.has("async-threadnum-factor")) {
                String factorstr = jsonObject.getString("async-threadnum-factor");
                double factor = Double.parseDouble(factorstr);
                asyncThreads = (int) (asyncThreads * factor);
            }
        } catch (IOException e) {
            logger.error(e.toString());
        }
        return asyncThreads;
    }

    private void start()
    {
        logger.debug("DBProxy Server is starting");

        DBProxySignalHandler signal_handler = new DBProxySignalHandler();
        Signal.handle(new Signal("TERM"), signal_handler);
        Signal.handle(new Signal("INT"), signal_handler);

        FragmentSessionManagement.getInstance().initialize();
        ConfigLoader.getInstance().start();

        RemoteCacheMonitorThread.getInstance().start();
        RemoteCacheClearThread.getInstance().start();

        boolean listen_ok = false;
        int asyncThreadnum = getAsyncThreadNum();
        if (asyncThreadnum <= 0) {
            asyncThreadnum = 1;
        }
        logger.debug("async thread num = " + asyncThreadnum);
        while(listen_ok == false) {
            try {
                if (isAsyncExecutorEnabled()) {
                    logger.debug("enabled async executor");
                    server = ServerBuilder.forPort(EnvironmentConfig.getInstance().getGRPCPort())
                            .addService(new NFDataManagementServiceImpl())
                            .executor(getExecutor(asyncThreadnum))
                            .build().start();
                } else {
                    logger.debug("disabled async executor");
                    server = ServerBuilder.forPort(EnvironmentConfig.getInstance().getGRPCPort())
                            .addService(new NFDataManagementServiceImpl())
                            .build().start();
                }
                listen_ok = true;
            } catch (IOException | IllegalStateException e) {
                logger.error(e.toString());
                logger.error("Fail to listen on port = {}", EnvironmentConfig.getInstance().getGRPCPort());
            }

            if(listen_ok == false) {
                logger.error("Sleep 5 seconds and try to listen again");
                try {
                    Thread.sleep(5000);
                } catch (Exception e) {
                    logger.error(e.toString());
                }
            }
        }

        logger.debug("DBProxy Server has been started successfully, listening on port " + EnvironmentConfig.getInstance().getGRPCPort());

        Runtime.getRuntime().addShutdownHook(new Thread(DBProxyServer.this::shutdown));

    }

    private void shutdown()
    {
        this.running = false;
        if (server != null) {
            logger.warn("Shutdown GRPC server");
            server.shutdown();
	    server = null;
        }

    }

    private void blockUntilShutdown()
    {
        if (server != null) {
            try {
                server.awaitTermination();
            } catch (InterruptedException e) {
                logger.error(e.toString());
                logger.error("Exception occurs while waiting for DBProxy Server Termination");
                Thread.currentThread().interrupt();
            }
        }
    }

    public static void main(String[] args)
    {
        final DBProxyServer db_proxy_server = DBProxyServer.getInstance();

        db_proxy_server.start();

        db_proxy_server.blockUntilShutdown();

        if(db_proxy_server.isRunning()) {
            db_proxy_server.shutdown();
        }
    }

    class DBProxySignalHandler implements SignalHandler
    {
        public void handle(Signal signal)
        {
            logger.warn("Catch signal = " + signal.getName());
            logger.warn("Sleep 10 seconds before shutdown DBProxy Server");
            try {
                Thread.sleep(10000);
            } catch (Exception e) {
                logger.error(e.toString());
            }

            DBProxyServer.getInstance().shutdown();
        }
    }
}

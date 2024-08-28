package com.ericsson.geode.functionservice;

import java.util.Map;
import java.util.Map.Entry;
import java.util.Set;
import java.util.concurrent.locks.Lock;
import org.apache.geode.cache.Region;
import org.apache.geode.cache.execute.Function;
import org.apache.geode.cache.execute.FunctionContext;
import org.apache.geode.cache.execute.RegionFunctionContext;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;

public class RemoteCacheDistributeLock implements Function {
  public static final String ID = RemoteCacheDistributeLock.class.getSimpleName();
  private static final String distributedLockKey = "DistributedLock";

  @Override
  public String getId() {
    return ID;
  }

  @Override
  public void execute(FunctionContext context) {
    RegionFunctionContext regionContext = (RegionFunctionContext) context;
    Region<String, PdxInstance> region = regionContext.getDataSet();
    Map<String, String> args = (Map<String, String>)regionContext.getArguments();
    String oper = "";
    String host = "";
    Set<Entry<String, String>> entryS = args.entrySet();
    for(Map.Entry<String, String> entry : entryS) {
        oper = entry.getKey();
        host = entry.getValue();
        break;
    }
     
    try {
      Lock distributedLock = region.getRegionDistributedLock();
      distributedLock.lock();
      Object obj = region.get(distributedLockKey);
      if (oper.equals(ClearCode.UnLock)) {
        if (null != obj) {
          PdxInstance value = (PdxInstance)obj;
          if(host.equals(value.getField(ClearCode.HostName))){
            region.put(distributedLockKey,
                JSONFormatter.fromJSON("{\"hostid\":\"" + ClearCode.ClearFinish + "\", \"expiry_time\":" + ClearCode.OccupyDistributedLockTime + "}"));
          } else if (!ClearCode.ClearFinish.equals(value.getField(ClearCode.HostName))){
             distributedLock.unlock();
             context.getResultSender().lastResult(ClearCode.GetLockFail);
             return;
          }
        }
        distributedLock.unlock();
        context.getResultSender().lastResult(ClearCode.GetLockSucc);
        return;
      } else {
        if (null == obj) {
          region.put(distributedLockKey,
              JSONFormatter.fromJSON("{\"hostid\":\"" + host + "\", \"expiry_time\":" + ClearCode.OccupyDistributedLockTime + "}"));
          distributedLock.unlock();
          context.getResultSender().lastResult(ClearCode.GetLockSucc);
          return;
        } else {
          PdxInstance value = (PdxInstance) obj;
          if (host.equals(value.getField(ClearCode.HostName))) {
            region.put(distributedLockKey,
                JSONFormatter.fromJSON("{\"hostid\":\"" + host + "\", \"expiry_time\":" + ClearCode.OccupyDistributedLockTime + "}"));
            distributedLock.unlock();
            context.getResultSender().lastResult(ClearCode.GetLockSucc);
            return;
          } else if (value.getField(ClearCode.HostName).equals(ClearCode.ClearFinish)) {
            region.put(distributedLockKey,
                JSONFormatter.fromJSON("{\"hostid\":\"" + host + "\", \"expiry_time\":" + ClearCode.OccupyDistributedLockTime + "}"));
            distributedLock.unlock();
            context.getResultSender().lastResult(ClearCode.GetLockSucc);
            return;
          } else {
            distributedLock.unlock();
            context.getResultSender().lastResult(ClearCode.GetLockFail);
            return;
          }
        }
      }
    } catch (Exception e){
      context.getResultSender().lastResult(ClearCode.GetLockFail);
      return ;
    }
  }

}

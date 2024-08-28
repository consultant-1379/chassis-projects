package com.ericsson.geode.functionservice;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.apache.geode.cache.Cache;
import org.apache.geode.cache.Region;
import org.apache.geode.cache.execute.Function;
import org.apache.geode.cache.execute.FunctionContext;
import org.apache.geode.cache.execute.RegionFunctionContext;
import org.apache.geode.pdx.PdxInstance;

public class RemoteCacheClearService implements Function {

    public static final String ID = RemoteCacheClearService.class.getSimpleName();

    @Override
    public String getId() {
        return ID;
    }

    @Override
    public void execute(FunctionContext context) {
        RegionFunctionContext regionContext = (RegionFunctionContext) context;
        Cache cache = regionContext.getCache();
        Region<String, PdxInstance> region = regionContext.getDataSet();
        int capacity = (int)regionContext.getArguments();
        Map<String, List<CacheSort>> cacheMap = new HashMap<String, List<CacheSort>>();
        Set<String> keySet= region.keySet();

        if (keySet.size() <= capacity || keySet.size() == 0) {
            context.getResultSender().lastResult(ClearCode.ClearSuccess);
            return ;
        }
        for (String k: keySet){

            PdxInstance value = region.get(k);
            if (null != value){
                long put_time = ((Number) value.getField(ClearCode.PutTime)).longValue();
                long expiry_time = ((Number) value.getField(ClearCode.ExpiryTime)).longValue();
                String from = (String)value.getField(ClearCode.From);

                if (cacheMap.containsKey(from)) {
                    cacheMap.get(from).add(new CacheSort(put_time, (put_time + expiry_time), k));
                } else {
                    List<CacheSort> mapValue = new ArrayList<CacheSort>();
                    mapValue.add(new CacheSort(put_time, (put_time + expiry_time), k));
                    cacheMap.put(from, mapValue);
                }
            }
        }

        for (Map.Entry<String, List<CacheSort>> entry : cacheMap.entrySet()){
            entry.getValue().sort((o1, o2)-> o1.getExpiry_time() > o2.getExpiry_time()
                ? 1 : o1.getExpiry_time() < o2.getExpiry_time()
                ? -1 :o1.getPut_time() > o2.getPut_time()
                ? -1 : o1.getPut_time() < o2.getPut_time()
                ? 1: 0);
        }
        int peerNRFNum = cacheMap.size();
        int peerNRFCapacity = capacity/peerNRFNum;
        if (peerNRFCapacity == 0) {
            peerNRFCapacity = 1;
        }

        try {
            Collection keys=new ArrayList();
            for (Map.Entry<String, List<CacheSort>> entry : cacheMap.entrySet()){
                int len = entry.getValue().size();
                if (len  > peerNRFCapacity) {
                    for (int i = 0; i < (len - peerNRFCapacity); i++) {
                         if (((Number)region.get(entry.getValue().get(i).getKey()).getField(ClearCode.PutTime)).longValue() == entry.getValue().get(i).getPut_time()) {

                             keys.add(entry.getValue().get(i).getKey());
                         }
                     }
                 }
            }
            cache.getCacheTransactionManager().begin();
            region.removeAll(keys);
            cache.getCacheTransactionManager().commit();

            for (Map.Entry<String, List<CacheSort>> entry : cacheMap.entrySet()){
                List<CacheSort> value = entry.getValue();
                for (CacheSort v : value){
                    v =  null;
                }
                value = null;
            }
            cacheMap.clear();
            cacheMap = null;
            keys = null;
            keySet = null;

            context.getResultSender().lastResult(ClearCode.ClearSuccess);
            return;
       } catch (Exception e){
            for (Map.Entry<String, List<CacheSort>> entry : cacheMap.entrySet()){
                List<CacheSort> value = entry.getValue();
                for (CacheSort v : value){
                    v =  null;
                }
                value = null;
            }
            cacheMap.clear();
            cacheMap = null;
            keySet = null;

            context.getResultSender().lastResult(ClearCode.ClearFail);
            return ;
       }
    }
}

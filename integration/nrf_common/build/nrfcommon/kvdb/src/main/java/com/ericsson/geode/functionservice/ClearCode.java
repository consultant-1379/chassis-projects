package com.ericsson.geode.functionservice;

public final class ClearCode
{
    private ClearCode() {
    }
    
    public static final int ClearSuccess = 0;
    public static final int ClearFail    = 1;
    public static final String PutTime   = "put_time";
    public static final String ExpiryTime= "expiry_time";
    public static final String From      = "from";

    public static final int GetLockSucc  = 10;
    public static final int GetLockFail  = 11;
    public static final String ClearFinish = "CLEARFINISH";
    public static final String HostName  = "hostid";
    public static final int OccupyDistributedLockTime = 10;
    public static final int UpdataDistributeInterval = 3;
    public static final String Lock      = "lock";
    public static final String UnLock    = "unlock";
}

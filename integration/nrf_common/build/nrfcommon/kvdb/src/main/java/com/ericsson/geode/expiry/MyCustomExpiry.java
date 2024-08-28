package com.ericsson.geode.expiry;

import org.apache.geode.cache.*;
import org.apache.geode.pdx.PdxInstance;

public class MyCustomExpiry implements CustomExpiry, Declarable {

    public ExpirationAttributes getExpiry(Region.Entry entry) {
        PdxInstance pdxInstance = (PdxInstance) entry.getValue();
        long expireTime = ((Number) pdxInstance.getField("expiry_time")).longValue();
        if (expireTime < 5) {
            expireTime = 5;
        }
        ExpirationAttributes CUSTOM_EXPIRY = new ExpirationAttributes((int) expireTime, ExpirationAction.DESTROY);
        return CUSTOM_EXPIRY;
    }
}

package ericsson.core.nrf.dbproxy.clientcache.schema;

import java.util.HashMap;

public class ImsiprefixProfiles
{
    private Long imsi_prefix;
    private HashMap  value_info;
    public ImsiprefixProfiles()
    {
        this.imsi_prefix = 0L;
        this.value_info = new HashMap();
    }

    public Long getImsiprefix()
    {
        return imsi_prefix;
    }

    public void setImsiprefix(Long imsi_prefix)
    {
        this.imsi_prefix = imsi_prefix;
    }

    public void addValueInfo(String value_info)
    {
        this.value_info.put(value_info, 1);
    }

    public void rmValueInfo(String value_info)
    {
        this.value_info.remove(value_info);
    }

    public HashMap getValueInfo() {
        return this.value_info;
    }
    public String toString()
    {
        return "ImsiprefixProfiles{" +
               "imsi_prefix=" + imsi_prefix.toString() +
               ", value_info=" + value_info.toString() +
               '}';
    }
}

package ericsson.core.nrf.dbproxy.clientcache.schema;

public class KeyAggregation
{

    private String sub_key1;
    private String sub_key2;
    private String sub_key3;
    private String sub_key4;
    private String sub_key5;

    public KeyAggregation()
    {

        sub_key1 = sub_key2 = sub_key3 = sub_key4 = sub_key5 = "";
    }

    public void setSubKey1(String value)
    {
        sub_key1 = value;
    }

    public void setSubKey2(String value)
    {
        sub_key2 = value;
    }

    public void setSubKey3(String value)
    {
        sub_key3 = value;
    }

    public void setSubKey4(String value)
    {
        sub_key4 = value;
    }

    public void setSubKey5(String value)
    {
        sub_key5 = value;
    }

    public String toString()
    {

        StringBuilder sb = new StringBuilder("");

        sb.append("sub_key1 = {" + sub_key1 + "}, ");
        sb.append("sub_key2 = {" + sub_key2 + "}, ");
        sb.append("sub_key3 = {" + sub_key3 + "}, ");
        sb.append("sub_key4 = {" + sub_key4 + "}, ");
        sb.append("sub_key5 = {" + sub_key5 + "}");

        return sb.toString();
    }

}

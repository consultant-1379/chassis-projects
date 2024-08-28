package ericsson.core.nrf.dbproxy.dbproxy_client;

import java.io.File;
import org.apache.commons.io.FileUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.json.JSONObject;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;

public class JSONUtil
{
	private static final Logger logger = LogManager.getLogger(JSONUtil.class);
	private JSONUtil(){}
    public static String[] readFile(String path)
    {
	String[] json = {"", ""};
	try
	{
	    String content = FileUtils.readFileToString(new File(path), "UTF-8");
	    JSONObject jsonObject = new JSONObject(content);
	    
            json[0] = jsonObject.getString("nfInstanceId");
	    json[1] = content;

    	    String jsonCustomer = "{Age: 10}";
	    JSONFormatter.fromJSON(jsonCustomer);
	}
	catch(Exception e)
	{
	    logger.error(e.toString());
	}
	return json;
    }
}

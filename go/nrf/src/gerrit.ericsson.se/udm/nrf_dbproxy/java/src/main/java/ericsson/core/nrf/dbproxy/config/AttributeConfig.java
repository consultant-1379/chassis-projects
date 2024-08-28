package ericsson.core.nrf.dbproxy.config;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;
import java.util.Iterator;

import org.json.JSONObject;
import org.json.JSONArray;
import org.json.JSONException;
import org.apache.commons.io.FileUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class AttributeConfig
{
    private static final Logger logger = LogManager.getLogger(AttributeConfig.class);

    private static AttributeConfig instance = null;

    HashMap attributes;
    private AttributeConfig()
    {
	attributes = new HashMap();
    }

    public static synchronized AttributeConfig getInstance() {
        if (null == instance) {
            instance = new AttributeConfig();
        }
        return instance;
    }

    public boolean load() {
        try {
            String attributes_conf = EnvironmentConfig.getInstance().getAttributesConf();
            String content = FileUtils.readFileToString(new File(attributes_conf), "UTF-8");
            JSONObject jsonObject = new JSONObject(content);
            Iterator iterator = jsonObject.keys();
            while (iterator.hasNext()) {
                JSONArray jsonArray = jsonObject.getJSONArray((String) iterator.next());
                for (int i = 0; i < jsonArray.length(); i++) {
                    JSONObject obj = (JSONObject) jsonArray.get(i);

                    String parameter = obj.getString("parameter");
                    String attribute_path = obj.getString("path");
                    String from = obj.getString("from");
                    String where = obj.getString("where");
                    boolean exist_check = obj.getBoolean("exist_check");

                    if (attributes.containsKey(attribute_path)) {
                        logger.error("Duplicate attribute path = " + attribute_path);
                        return false;
                    } else
                        attributes.put(attribute_path, new Attribute(parameter, attribute_path, from, where, exist_check));
                }
            }

        } catch (IOException | JSONException | NullPointerException | SecurityException | NumberFormatException e) {
            logger.error(e.toString());
            logger.warn("Fail to parse " + EnvironmentConfig.getInstance().getAttributesConf());
            return false;
        } catch (Exception e) {
            logger.error(e.toString());
            logger.warn("Fail to parse " + EnvironmentConfig.getInstance().getAttributesConf());
            return false;
        }

        for (Object obj : attributes.values()) {
            logger.trace(((Attribute) obj).toString());
        }

        return true;
    }

    public Attribute get(String attribute_path)
    {
	Object obj = attributes.get(attribute_path);
	if( null == obj) return null;
	return (Attribute)obj;
    }
}

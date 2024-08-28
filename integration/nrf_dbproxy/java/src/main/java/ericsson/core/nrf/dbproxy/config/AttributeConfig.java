package ericsson.core.nrf.dbproxy.config;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;
import java.util.Iterator;
import org.apache.commons.io.FileUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class AttributeConfig {

  private static final Logger LOGGER = LogManager.getLogger(AttributeConfig.class);

  private static AttributeConfig instance;

  static {
    instance = null;
  }

  HashMap attributes;

  private AttributeConfig() {
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
      String attributesConf = EnvironmentConfig.getInstance().getAttributesConf();
      String content = FileUtils.readFileToString(new File(attributesConf), "UTF-8");
      JSONObject jsonObject = new JSONObject(content);
      Iterator iterator = jsonObject.keys();
      while (iterator.hasNext()) {
        JSONArray jsonArray = jsonObject.getJSONArray((String) iterator.next());
        for (int i = 0; i < jsonArray.length(); i++) {
          JSONObject obj = (JSONObject) jsonArray.get(i);

          String parameter = obj.getString("parameter");
          String attributePath = obj.getString("path");
          String from = obj.getString("from");
          String where = obj.getString("where");
          boolean existCheck = obj.getBoolean("exist_check");

          if (attributes.containsKey(attributePath)) {
            LOGGER.error("Duplicate attribute path = " + attributePath);
            return false;
          } else {
            attributes.put(attributePath,
                new Attribute(parameter, attributePath, from, where, existCheck));
          }
        }
      }

    } catch (IOException | JSONException | NullPointerException | SecurityException | NumberFormatException e) {
      LOGGER.error(e.toString());
      LOGGER.warn("Fail to parse " + EnvironmentConfig.getInstance().getAttributesConf());
      return false;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      LOGGER.warn("Fail to parse " + EnvironmentConfig.getInstance().getAttributesConf());
      return false;
    }

    for (Object obj : attributes.values()) {
      LOGGER.debug(((Attribute) obj).toString());
    }

    return true;
  }

  public Attribute get(String attributePath) {
    Object obj = attributes.get(attributePath);
    if (null == obj) {
      return null;
    }
    return (Attribute) obj;
  }
}

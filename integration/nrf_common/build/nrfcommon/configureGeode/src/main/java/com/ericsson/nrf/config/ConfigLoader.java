package com.ericsson.nrf.config;

import com.ericsson.nrf.common.ConfigConst;
import org.apache.commons.io.IOUtils;
import org.json.JSONArray;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;

public class ConfigLoader {
    private static ConfigLoader instance;

    private List<RegionBean> regionList;
    private List<IndexBean> indexList;

    public ConfigLoader() {
        regionList = new ArrayList<>();
        indexList = new ArrayList<>();
    }

    public static synchronized ConfigLoader getInstance() {
        if (instance == null) {
            instance = new ConfigLoader();
        }
        return instance;
    }

    public List<RegionBean> getRegionList() {
        return regionList;
    }

    public List<IndexBean> getIndexList() {
        return indexList;
    }

    public void parseRegionConfigFile() {
        InputStream inputStream = getClass().getResourceAsStream(ConfigConst.REGION_FILE);
        try {
            String regions = IOUtils.toString(inputStream);
            JSONArray jsonArray = new JSONArray(regions);
            for (int i = 0; i < jsonArray.length(); i++) {
                JSONObject jsonObject = jsonArray.getJSONObject(i);
                RegionBean bean = new RegionBean();
                bean.setName(jsonObject.getString(ConfigConst.REGION_NAME));
                bean.setType(jsonObject.getString(ConfigConst.REGION_TYPE));
                bean.setAdditional(jsonObject.getString(ConfigConst.REGION_ADDITIONAL));
                regionList.add(bean);
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public void parseIndexConfigFile() {
        InputStream inputStream = getClass().getResourceAsStream(ConfigConst.INDEX_FILE);
        try {
            String indexes = IOUtils.toString(inputStream);
            JSONArray jsonArray = new JSONArray(indexes);
            for (int i = 0; i < jsonArray.length(); i++) {
                JSONObject jsonObject = jsonArray.getJSONObject(i);
                IndexBean bean = new IndexBean();
                bean.setName(jsonObject.getString(ConfigConst.INDEX_NAME));
                bean.setType(jsonObject.getString(ConfigConst.INDEX_TYPE));
                bean.setExpression(jsonObject.getString(ConfigConst.INDEX_EXPRESSION));
                bean.setRegion(jsonObject.getString(ConfigConst.INDEX_REGION));
                indexList.add(bean);
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}

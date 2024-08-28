package cmproxy

import (
	"encoding/json"
	"errors"
	"net/http"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

type cmConfiguration struct {
	Name  string          `json:"name"`
	Title string          `json:"title"`
	Data  json.RawMessage `json:"data"`
}

func doConfigurations(method, url string) (*httpclient.HttpRespData, error) {
	return cmHTTPClient.HttpDoJsonBody(method, url, nil)
}

func getConfiguration(configName string) (*httpclient.HttpRespData, error) {
	resp, err := doConfigurations("GET", getCmmURL("configurations", configName))
	if err != nil {
		log.Errorf("failed to GET configurations %s, %s", configName, err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		e := errors.New(string(resp.Body))
		log.Errorf("failed to GET configurations, %s", e.Error())
		return nil, e
	}
	return resp, nil
}

//GetConfiguration get configurations form CM Mediator immediately
func GetConfiguration(configName string) ([]byte, error) {
	resp, err := doConfigurations("GET", getCmmURL("configurations", configName))
	if err != nil {
		log.Errorf("failed to GET configurations %s, %s", configName, err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		e := errors.New(string(resp.Body))
		log.Errorf("failed to GET configurations, %s", e.Error())
		return nil, e
	}

	var configuration cmConfiguration
	err = json.Unmarshal(resp.Body, &configuration)
	if err != nil {
		log.Errorf("failed to Unmarshal configuration %s, %s, %s.", configName, string(resp.Body), err.Error())
		return nil, err
	}
	return configuration.Data, nil
}

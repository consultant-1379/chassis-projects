package kvdbclient

import (
	"bytes"
	"errors"
	"fmt"
	"time"
	"gerrit.ericsson.se/udm/common/pkg/log"
	typeKVDB "gerrit.ericsson.se/udm/nrf_common/pkg/kvdbclient/encoding"
	"net/http"
	"strings"
)

// KVDBClient struct
type KVDBClient struct {
	httpClient http.Client
	gfshURL    string
}

var instance *KVDBClient

// GetInstance for KVDBClient
func GetInstance() *KVDBClient {
	if instance == nil {
		instance = newKVDBClient()
	}
	return instance
}

func newKVDBClient() *KVDBClient {
	c := http.Client{}
	gfshURL := "http://eric-nrf-kvdb-ag-admin-mgr:8080/kvdb-ag/management/v1/gfsh-commands"
	return &KVDBClient{c, gfshURL}
}

// SendGFSHCommand to execute gfsh command
func (c *KVDBClient) SendGFSHCommand(command string) (string, error) {
	cmd := typeKVDB.GFSHCommand{
		Command: command,
	}

	buf, err := typeKVDB.EncodeGFSHCommand(cmd)
	if err != nil {
		log.Error(err)
		return "", err
	}

	reqURL := c.gfshURL
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(buf))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := c.httpClient
	response, err := client.Do(req)
        if err != nil {
                return "", err
        }

	if response.StatusCode != 202 {
		log.Warningf("HTTP error response %v", response.StatusCode)
		return "", errors.New("HTTP error response " + fmt.Sprintf("%v", response.StatusCode))
	}

	result, err := typeKVDB.DecodeGFSHCommandID(response.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return result.CommandId, nil
}

// GetGFSHCommandResult to get gfsh execute result
func (c *KVDBClient) GetGFSHCommandResult(id string) (*typeKVDB.GFSHCommand, error) {
	// wait 500ms for excuting command
	time.Sleep(500 * time.Millisecond)

	reqURL := c.gfshURL + "/" + id
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	client := c.httpClient
	response, err := client.Do(req)
        if err != nil {
                return nil, err
        }

	if response.StatusCode != 200 {
		log.Warningf("HTTP error response %v", response.StatusCode)
		return nil, errors.New("HTTP error response " + fmt.Sprintf("%v", response.StatusCode))
	}

	return typeKVDB.DecodeGFSHCommand(response.Body)
}

// GetQueryResult to get result for query
func (c *KVDBClient) GetQueryResult(id string) ([]string, error) {
	var values []string
	retry := 3
	for retry > 0 {
		result, err := c.GetGFSHCommandResult(id)
		if err != nil || result.ExecutionStatus != "EXECUTED" {
			retry--
			continue
		}

		tmp := strings.Split(result.Output, "-\n")
		if len(tmp) != 2 {
			return values, errors.New("unexpected result")
		}
		result.Output = strings.TrimSpace(tmp[1])
		values = strings.Split(result.Output, "\n")
		return values, nil
	}
	return values, nil
}

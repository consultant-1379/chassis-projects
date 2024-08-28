package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NrfAgentMsg struct {
	NfCMProfileName   string `json:"nfCMProfileName,omitempty"`
	NfProfilePath     string `json:"nfProfilePath,omitempty"`
	HeartbeatInterval uint64 `json:"heartbeatInterval,omitempty"`
}

var httpClient = &http.Client{
	Transport: &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (conn net.Conn, e error) {
			return net.Dial(network, addr)
		},
	},
	Timeout: time.Second * 2}

func addSuffix(uri string) string {
	if !strings.HasSuffix(uri, "/") {
		uri += "/"
	}
	return uri
}

func sendHeartbeat(serviceName string) error {
	config.nrfAgentUri = addSuffix(config.nrfAgentUri)
	target := config.nrfAgentUri + "nrf-register-agent/v1/nf-status/NEF/" + serviceName
	message := new(NrfAgentMsg)
	message.HeartbeatInterval = config.heartbeatInterval
	message.NfProfilePath = config.nfProfilePath
	message.NfCMProfileName = config.nfCMProfileName
	bytesData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	logger.Debugf("sending to %s, body: %s", target, string(bytesData))
	body := bytes.NewBuffer(bytesData)
	req, err := http.NewRequest(http.MethodPut, target, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if resp, err := httpClient.Do(req); err != nil {
		return err
	} else {
		defer resp.Body.Close()
		switch resp.StatusCode {
		case 200:
			logger.Debugf("successfully sent HB of the service( %s )", serviceName)
			return nil
		default:
			return errors.New("NRFAgent responded with StatusCode = " + strconv.Itoa(resp.StatusCode) + " for HB of the service( " + serviceName + " )")
		}
	}
}

package adpcm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"gerrit.ericsson.se/HSS/5G/cmproxy/statistics"
	"golang.org/x/net/http2"
)

const (
	separator  = "/"
	root       = "/"
	empty      = ""
	noSuchData = "No such data"
)

type result struct {
	Data interface{} `json:"data"`
}

type adpcm struct {
	cmUri         string
	data          result
	notifEndpoint string
}

// GetValue : get cfg value from key
// The key should be started with '/'
func (cm *adpcm) GetValue(key string) (string, error) {

	err := errors.New(noSuchData)
	element := cm.data.Data
	keyin := strings.TrimLeft(key, root)
	if keyin != empty {
		keys := strings.Split(keyin, separator)

		for _, k := range keys {
			switch t := element.(type) {
			case map[string]interface{}:
				if data, ok := t[k]; ok {
					element = data
				} else {
					return empty, err
				}
			default:
				return empty, err
			}
		}
	}

	if b, e := json.Marshal(element); e == nil {
		return string(b), nil
	}

	return empty, err
}

func NewAdpCm(cmUri, notifEndpoint string) *adpcm {

	cm := &adpcm{cmUri: cmUri, notifEndpoint: notifEndpoint}

	if cm.reloadData() {
		return cm
	} else {
		return nil
	}
}

func (cm *adpcm) reloadData() bool {

	rsp, err := http.Get(cm.cmUri)
	statistics.Statistics.NumberOfAdpReads = statistics.Statistics.NumberOfAdpReads + 1
	statistics.Statistics.LastAdpUpdate = time.Now().String()

	if err != nil {
		log.Printf("Get %s failed. Error is %s\n", cm.cmUri, err.Error())
		return false
	}

	if rsp.StatusCode != http.StatusOK {
		log.Printf("Get %s failed. Status is %s\n", cm.cmUri, rsp.Status)
		return false
	}

	d, errRead := ioutil.ReadAll(rsp.Body)
	if errRead != nil {
		log.Printf("Reading data of %s failed. Error is %s\n", cm.cmUri, errRead.Error())
		return false
	}

	var res result
	if err := json.Unmarshal(d, &res); err == nil {
		cm.data = res
		return true
	} else {
		log.Printf("Unmarshal data %s Failed, Error is %s\n", string(d), err.Error())
		return false
	}
}

func (cm *adpcm) MonitorToReLoad(msg []byte) {

	var res result

	if err := json.Unmarshal(msg, &res); err != nil {
		log.Println(err)
		return
	}

	//Reload data
	if cm.reloadData() {
		tr := &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			AllowHTTP:       true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}
		c := http.Client{Transport: tr}

		req, err := http.NewRequest(http.MethodPost, cm.notifEndpoint, bytes.NewBuffer([]byte{}))

		if err != nil {
			log.Println(err)
		}

		if rsp, err := c.Do(req); err != nil {
			log.Println(err)
		} else {
			if rsp.StatusCode != http.StatusOK {
				log.Println("status is " + rsp.Status)
			}
		}
	} else {
		log.Println("reload data failed")
	}
}

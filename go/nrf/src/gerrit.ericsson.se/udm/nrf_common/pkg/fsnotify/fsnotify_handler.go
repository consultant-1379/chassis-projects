package fsnotify

import (
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
)

type fileHandler interface {
	Handler(op string)
}

type configmapHandler struct {
	fileName string
}

func (c *configmapHandler) Handler(op string) {
	if op != "REMOVE" {
		return
	}

	if configuration, ok := configmap.ConfigMapMap[c.fileName]; ok {
		configuration.ReloadConf()
	}
}

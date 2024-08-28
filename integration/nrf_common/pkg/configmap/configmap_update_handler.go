package configmap

type configmapUpdateHandler struct {
}

func (c *configmapUpdateHandler) Handler(fileName string, op string) {
	if op != "REMOVE" {
		return
	}

	if configuration, ok := ConfigMapMap[fileName]; ok {
		configuration.ReloadConf()
	}
}

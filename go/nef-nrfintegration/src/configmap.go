package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type Config struct {
	namespace         string
	prefix            string
	nrfAgentUri       string
	nfCMProfileName   string
	nfProfilePath     string
	serviceNameMap    map[string]string
	heartbeatInterval uint64
	retryInterval     uint64
	retryTimes        uint64
}

func requireNotEmpty(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic(name + " is empty in ConfigMap")
	}
	logger.Info(name, " = ", v)
	return v
}

var (
	ErrSmallHeartbeatInterval = errors.New("retryInterval * retryTimes >= heartbeatInterval is not allowed")
)

func (config *Config) LoadConfigMap() error {
	logger.Info("Loading ConfigMap")
	config.prefix = requireNotEmpty("GLOBAL_PREFIX")
	config.namespace = requireNotEmpty("NAMESPACE")
	config.nrfAgentUri = requireNotEmpty("NRF_AGENT_URI")
	config.nfCMProfileName = requireNotEmpty("NF_CM_PROFILE_NAME")
	config.nfProfilePath = requireNotEmpty("NF_PROFILE_PATH")

	var err error
	if err = json.Unmarshal([]byte(requireNotEmpty("SERVICE_NAME_MAP")), &config.serviceNameMap); err != nil {
		return err
	}

	if config.heartbeatInterval, err = strconv.ParseUint(requireNotEmpty("HEARTBEAT_INTERVAL"), 10, 32); err != nil {
		return err
	}

	if config.retryInterval, err = strconv.ParseUint(requireNotEmpty("RETRY_INTERVAL"), 10, 32); err != nil {
		return err
	}

	if config.retryTimes, err = strconv.ParseUint(requireNotEmpty("RETRY_TIMES"), 10, 32); err != nil {
		return err
	}

	if config.retryTimes*config.retryInterval >= config.heartbeatInterval {
		return ErrSmallHeartbeatInterval
	}
	return nil
}

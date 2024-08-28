package siteleader

import (
	"gerrit.ericsson.se/udm/common/pkg/clusterleader"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/multisite"
)

// IsLeader is used to judge whether caller is a site leader
func IsLeader() bool {
	if clusterleader.IsLeader() {
		if configmap.MultisiteEnabled {
			return multisite.GetMonitor().IsLeaderSite()
		}
		return true
	}
	return false
}

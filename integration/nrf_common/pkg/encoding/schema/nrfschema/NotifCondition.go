package nrfschema

// IsValid check whether tacRange is valid
func (n *TNotifCondition) IsValid() bool {
	if n.MonitoredAttributes == nil || n.UnmonitoredAttributes == nil {
		return true
	}

	return false
}

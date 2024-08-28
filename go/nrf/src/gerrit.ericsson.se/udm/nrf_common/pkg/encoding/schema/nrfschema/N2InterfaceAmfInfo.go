package nrfschema

// IsValid check whether N2InterfaceAmfInfo is valid
func (n *TN2InterfaceAmfInfo) IsValid() bool {
	if n.Ipv4EndpointAddress != nil || n.Ipv6EndpointAddress != nil {
		return true
	}

	return false
}

package nrfschema

// IsValid check whether InterfaceUpfInfoItem is valid
func (i *TInterfaceUpfInfoItem) IsValid() bool {
	if i.EndpointFqdn != "" || i.Ipv4EndpointAddresses != nil || i.Ipv6EndpointAddresses != nil {
		return true
	}

	return false
}

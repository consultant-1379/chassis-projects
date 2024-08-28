package nrfschema

// IsValid check whether IpEndPoint is valid
func (i *TIpEndPoint) IsValid() bool {
	if i.Ipv4Address == "" || i.Ipv6Address == "" {
		return true
	}

	return false
}

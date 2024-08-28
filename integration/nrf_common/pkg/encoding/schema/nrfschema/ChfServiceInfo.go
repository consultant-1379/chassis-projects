package nrfschema

// IsValid check whether ChfServiceInfo is valid
func (c *TChfServiceInfo) IsValid() bool {
	if c.PrimaryChfServiceInstance == "" || c.SecondaryChfServiceInstance == "" {
		return true
	}

	return false
}

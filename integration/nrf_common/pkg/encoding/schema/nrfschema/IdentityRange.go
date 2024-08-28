package nrfschema

// IsValid check whether IdentityRange is valid
func (i *TIdentityRange) IsValid() bool {
	if i.Start != "" && i.End != "" && i.Pattern == "" {
		return true
	}

	if i.Start == "" && i.End == "" && i.Pattern != "" {
		return true
	}

	return false
}

// IsEqual check whether IdentityRange is equal
func (i *TIdentityRange) IsEqual(c *TIdentityRange) bool {
	if i.Pattern != c.Pattern || i.Start != c.Start || i.End != c.End {
		return false
	}

	return true
}

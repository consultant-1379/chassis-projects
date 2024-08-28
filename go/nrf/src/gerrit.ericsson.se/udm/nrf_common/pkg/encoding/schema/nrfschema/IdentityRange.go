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

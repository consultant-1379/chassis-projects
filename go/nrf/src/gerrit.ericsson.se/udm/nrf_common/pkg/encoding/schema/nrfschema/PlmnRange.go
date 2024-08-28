package nrfschema

// IsValid check whether PlmnRange is valid
func (t *TPlmnRange) IsValid() bool {
	if t.Start != "" && t.End != "" && t.Pattern == "" {
		return true
	}

	if t.Start == "" && t.End == "" && t.Pattern != "" {
		return true
	}

	return false
}

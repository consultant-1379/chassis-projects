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

// IsEqual check whether PlmnRange is equal
func (t *TPlmnRange) IsEqual(c *TPlmnRange) bool {
	if t.Pattern != c.Pattern || t.Start != c.Start || t.End != c.End {
		return false
	}

	return true
}

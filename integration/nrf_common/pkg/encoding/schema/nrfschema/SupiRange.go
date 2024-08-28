package nrfschema

// IsValid check whether SupiRange is valid
func (t *TSupiRange) IsValid() bool {
	if t.Start != "" && t.End != "" && t.Pattern == "" {
		return true
	}

	if t.Start == "" && t.End == "" && t.Pattern != "" {
		return true
	}

	return false
}

// IsEqual check whether SupiRange is equal
func (t *TSupiRange) IsEqual(c *TSupiRange) bool {
	if t.Pattern != c.Pattern || t.Start != c.Start || t.End != c.End {
		return false
	}

	return true
}

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

package nrfschema

// IsValid check whether tacRange is valid
func (t *TTacRange) IsValid() bool {
	if t.Start != "" && t.End != "" && t.Pattern == "" {
		return true
	}

	if t.Start == "" && t.End == "" && t.Pattern != "" {
		return true
	}

	return false
}

// IsEqual check whether tacRange is equal
func (t *TTacRange) IsEqual(c *TTacRange) bool {
	if t.Pattern != c.Pattern || t.Start != c.Start || t.End != c.End {
		return false
	}

	return true
}

package nrfschema

import (
	"fmt"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidTacRangeIndexs return invalid tacRangeList index
func (t *TTaiRange) GetInvalidTacRangeIndexs() []string {
	var invalidTacRangeIndexs []string

	if t.TacRangeList != nil {
		index := 0
		for _, tacRange := range t.TacRangeList {
			if !tacRange.IsValid() {
				invalidTacRangeIndex := fmt.Sprintf("%s[%d]", constvalue.TacRangeList, index)
				invalidTacRangeIndexs = append(invalidTacRangeIndexs, invalidTacRangeIndex)
			}
			index++
		}
	}

	return invalidTacRangeIndexs
}

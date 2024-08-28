package nrfschema

import (
	"fmt"
)

// GetPlmnID returns the mcc+mnc
func (p *TPlmnId) GetPlmnID() string {
	return p.Mcc + p.Mnc
}

//ToString returns the information of PlmnID
func (p *TPlmnId) ToString() string {
	return fmt.Sprintf("{mcc: %s, mnc: %s}", p.Mcc, p.Mnc)
}

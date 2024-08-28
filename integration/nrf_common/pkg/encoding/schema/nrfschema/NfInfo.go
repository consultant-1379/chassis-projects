package nrfschema

import (
	"com/dbproxy/nfmessage/subscription"
)

// TNfInfo is a interface which converges a set of functions implemented by some struct,
// such as TAusfInfo, TUdmInfo or TUdrInfo
type TNfInfo interface {
	GenerateNfGroupCond() *subscription.SubKeyStruct
	IsEqual(TNfInfo) bool
}

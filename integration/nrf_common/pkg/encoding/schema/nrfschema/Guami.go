package nrfschema

import (
	"com/dbproxy/nfmessage/subscription"
)

//GenerateGrpcKey generate a SubKeyStruct used in subscriptionIndex
func (g *TGuami) GenerateGrpcKey() *subscription.SubKeyStruct {
	return &subscription.SubKeyStruct{
		SubKey1: g.PlmnId.GetPlmnID(),
		SubKey2: g.AmfId,
	}
}

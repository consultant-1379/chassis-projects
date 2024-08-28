package cm

import (
	"testing"
)

func TestTNrfPolicyToUpper(t *testing.T) {
	policy := TNrfPolicy{
		ManagementService: &TNrfManagementServicePolicy{
			Subscription: &TSubscriptionPolicy{
				AllowedSubscriptionAllNFs: []TAllowedSubscriptionAllNFs{
					TAllowedSubscriptionAllNFs{
						AllowedNfType:    "amf",
						AllowedNfDomains: "xxxx1",
					},
					TAllowedSubscriptionAllNFs{
						AllowedNfType:    "ausf",
						AllowedNfDomains: "xxxx2",
					},
					TAllowedSubscriptionAllNFs{
						AllowedNfType:    "PCF",
						AllowedNfDomains: "xxxx3",
					},
				},
			},
		},
	}

	policy.toUpper()

	nfTypeMap := make(map[string]bool)

	for _, item := range policy.ManagementService.Subscription.AllowedSubscriptionAllNFs {
		nfTypeMap[item.AllowedNfType] = true
	}

	if len(nfTypeMap) != 3 {
		t.Fatalf("TNrfPolicy.toUpper didn't return right value !")
	}

	if !nfTypeMap["AMF"] || !nfTypeMap["AUSF"] || !nfTypeMap["PCF"] {
		t.Fatalf("TNrfPolicy.toUpper didn't return right value !")
	}

}

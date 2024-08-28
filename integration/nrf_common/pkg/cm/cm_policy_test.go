package cm

import (
	"testing"
)

func TestParseConfForPolicy(t *testing.T) {
	var nrfPolicy TNrfPolicy
	nrfPolicy.ParseConf()
	if GetNRFPolicy().ManagementService == nil {
		t.Fatalf("NrfPolicy.ManagementService don't set default !")
	}

	if GetNRFPolicy().ManagementService.Subscription == nil {
		t.Fatalf("NrfPolicy.ManagementService.Subscription don't set default !")
	}
}

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

func TestGetPriorityGroup(t *testing.T) {
	// case1: if sbi-message-priority-policy is not configured, return default
	nrfPolicy := TNrfPolicy{}
	nrfPolicy.atomicSetNRFPolicy()
	if GetPriorityGroup() == nil {
		t.Fatal("sbi-message-priority-policy is not configured, GetPriorityGroup() should return default, but not !")
	}

	// case2: if both priority-group-medium-start and priority-group-low-start are not configured, use default
	nrfPolicy = TNrfPolicy{
		MessagePriorityPolicy: &TMessagePriorityPolicy{},
	}
	nrfPolicy.atomicSetNRFPolicy()
	priorityGroup := GetPriorityGroup()
	if priorityGroup == nil {
		t.Fatal("sbi-message-priority-policy is configured, GetPriorityGroup() should not return nil, but did !")
	}

	if len(priorityGroup) != 3 {
		t.Fatal("GetPriorityGroup() should return 3 items, but not !")
	}

	priorityMap := make(map[int]bool, 0)
	for _, item := range priorityGroup {
		priorityMap[item.Level] = true
	}

	if len(priorityMap) != 3 {
		t.Fatal("length of priorityMap should be 3")
	}

	for k, _ := range priorityMap {
		if k != 1 && k != 2 && k != 3 {
			t.Fatal("the priority group level is wrong")
		}
	}

	for _, item := range priorityGroup {
		if item.Level == 1 {
			if item.Start != 0 || item.End != 7 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}

		if item.Level == 2 {
			if item.Start != 8 || item.End != 15 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}

		if item.Level == 3 {
			if item.Start != 16 || item.End != 31 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}
	}

	// case3: if priority-group-medium-start is not configured, use default
	lowStart := 15
	nrfPolicy = TNrfPolicy{
		MessagePriorityPolicy: &TMessagePriorityPolicy{
			LowStart: &lowStart,
		},
	}
	nrfPolicy.atomicSetNRFPolicy()
	priorityGroup = GetPriorityGroup()
	if priorityGroup == nil {
		t.Fatal("sbi-message-priority-policy is configured, GetPriorityGroup() should not return nil, but did !")
	}

	if len(priorityGroup) != 3 {
		t.Fatal("GetPriorityGroup() should return 3 items, but not !")
	}

	priorityMap = make(map[int]bool, 0)
	for _, item := range priorityGroup {
		priorityMap[item.Level] = true
	}

	if len(priorityMap) != 3 {
		t.Fatal("length of priorityMap should be 3")
	}

	for k, _ := range priorityMap {
		if k != 1 && k != 2 && k != 3 {
			t.Fatal("the priority group level is wrong")
		}
	}

	for _, item := range priorityGroup {
		if item.Level == 1 {
			if item.Start != 0 || item.End != 7 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}

		if item.Level == 2 {
			if item.Start != 8 || item.End != 14 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}

		if item.Level == 3 {
			if item.Start != 15 || item.End != 31 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}
	}

	// case4: if priority-group-low-start is not configured, use default
	mediumStart := 13
	nrfPolicy = TNrfPolicy{
		MessagePriorityPolicy: &TMessagePriorityPolicy{
			MediumStart: &mediumStart,
		},
	}
	nrfPolicy.atomicSetNRFPolicy()
	priorityGroup = GetPriorityGroup()
	if priorityGroup == nil {
		t.Fatal("sbi-message-priority-policy is configured, GetPriorityGroup() should not return nil, but did !")
	}

	if len(priorityGroup) != 3 {
		t.Fatal("GetPriorityGroup() should return 3 items, but not !")
	}

	priorityMap = make(map[int]bool, 0)
	for _, item := range priorityGroup {
		priorityMap[item.Level] = true
	}

	if len(priorityMap) != 3 {
		t.Fatal("length of priorityMap should be 3")
	}

	for k, _ := range priorityMap {
		if k != 1 && k != 2 && k != 3 {
			t.Fatal("the priority group level is wrong")
		}
	}

	for _, item := range priorityGroup {
		if item.Level == 1 {
			if item.Start != 0 || item.End != 12 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}

		if item.Level == 2 {
			if item.Start != 13 || item.End != 15 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}

		if item.Level == 3 {
			if item.Start != 16 || item.End != 31 {
				t.Fatal("GetPriorityGroup didn't return right priority group")
			}
		}
	}

	// case5: if priority-group-low-start < priority-group-medium-start, return default
	lowStart = 10
	mediumStart = 13
	nrfPolicy = TNrfPolicy{
		MessagePriorityPolicy: &TMessagePriorityPolicy{
			LowStart:    &lowStart,
			MediumStart: &mediumStart,
		},
	}
	nrfPolicy.atomicSetNRFPolicy()
	priorityGroup = GetPriorityGroup()
	if priorityGroup == nil {
		t.Fatal("sbi-message-priority-policy is not configured rightly, GetPriorityGroup() should return default, but not !")
	}

}

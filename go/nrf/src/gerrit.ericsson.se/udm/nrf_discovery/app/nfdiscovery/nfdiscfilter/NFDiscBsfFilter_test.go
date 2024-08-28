package nfdiscfilter

import (
	"testing"
)


func TestIsMatchedIPv4Addr(t *testing.T) {
	bsfFilter := &NFBSFInfoFilter{}
	bspInfo := []byte(`{
    "ipv4AddressRanges": [
      {
        "start": "172.16.4.1",
        "end": "172.16.4.254"
      }
    ],
    "ipv6PrefixRanges": [
      {
        "start": "1030:C9B4:FF1b::/96",
        "end": "1030:C9B4:FF1b::/64"
      }
    ]
  }`)
	ok := bsfFilter.isMatchedIPv4Addr("172.16.4.2", bspInfo)
	if !ok {
		t.Fatal("IPv4Addr should be matched, but not !")
	}
	ok1 := bsfFilter.isMatchedIPv4Addr("172.16.4.255", bspInfo)
	if ok1 {
		t.Fatal("IPv4Addr should NOT be matched, but matched !")
	}
}


func TestIsMatchedIPv6AddrPrefix(t *testing.T) {
	bsfFilter := &NFBSFInfoFilter{}
	bspInfo := []byte(`{
    "ipv4AddressRanges": [
      {
        "start": "172.16.4.1",
        "end": "172.16.4.254"
      }
    ],
    "ipv6PrefixRanges": [
      {
        "start": "1030:C9B4::/32",
        "end": "1030:C9B4:FF1b::/32"
      }
    ]
  }`)
	ok := bsfFilter.isMatchedIPv6AddrPrefix("1030:C9B4:FF11::/32", bspInfo)
	if !ok {
		t.Fatal("IPv6Addr should be matched, but not !")
	}
	ok1 := bsfFilter.isMatchedIPv6AddrPrefix("1030:C9B4:FF1c::/32", bspInfo)
	if !ok1 {
		t.Fatal("IPv6Addr should be matched, but not !")
	}

	ok2 := bsfFilter.isMatchedIPv6AddrPrefix("1030:C9B5:FF1c::/32", bspInfo)
	if ok2 {
		t.Fatal("IPv6Addr shoudl not matched, but matched")
	}
}

func TestIsMatchedIPDomain(t *testing.T)  {
	bsfFilter := &NFBSFInfoFilter{}
	bsfInfo := []byte(`{
    	"ipDomainList": [
      		"www.example.com"
    	]
    	}`)
	bsfInfo2 := []byte(`{

    	}`)
	ok := bsfFilter.isMatchedIPDomain("www.example.com", bsfInfo)
	if !ok {
		t.Fatal("ip-domain should be matched, but not !")
	}
	ok1 := bsfFilter.isMatchedIPDomain("www.test.com", bsfInfo)
	if ok1 {
		t.Fatal("ip-domain should NOT be matched, but matched !")
	}
	ok2 := bsfFilter.isMatchedIPDomain("www.test.com", bsfInfo2)
	if ok2 {
		t.Fatal("ip-domain should not be matched, but natched !")
	}
}
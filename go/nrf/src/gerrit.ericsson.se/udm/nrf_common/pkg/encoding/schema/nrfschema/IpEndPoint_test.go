package nrfschema

import (
	"encoding/json"
	"testing"
)

func TestIsIpEndPointValid(t *testing.T) {
	//right IpEndPoint
	body := []byte(`{
		"port": 80
	}`)

	ipEndPoint := &TIpEndPoint{}
	err := json.Unmarshal(body, ipEndPoint)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !ipEndPoint.IsValid() {
		t.Fatalf("TIpEndPoint.IsValid didn't return right value!")
	}

	//right IpEndPoint
	body = []byte(`{
		"ipv4Address": "10.10.10.10",
		"port": 80
	}`)

	ipEndPoint = &TIpEndPoint{}
	err = json.Unmarshal(body, ipEndPoint)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !ipEndPoint.IsValid() {
		t.Fatalf("TIpEndPoint.IsValid didn't return right value!")
	}

	//right IpEndPoint
	body = []byte(`{
		"ipv6Address": "1030::C9B4:FF12:48AA:1A2B",
		"port": 80
	}`)

	ipEndPoint = &TIpEndPoint{}
	err = json.Unmarshal(body, ipEndPoint)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if !ipEndPoint.IsValid() {
		t.Fatalf("TIpEndPoint.IsValid didn't return right value!")
	}

	//wrong IpEndPoint
	body = []byte(`{
		"ipv4Address": "10.10.10.10",
		"ipv6Address": "1030::C9B4:FF12:48AA:1A2B",
		"port": 80
	}`)

	ipEndPoint = &TIpEndPoint{}
	err = json.Unmarshal(body, ipEndPoint)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	if ipEndPoint.IsValid() {
		t.Fatalf("TIpEndPoint.IsValid didn't return right value!")
	}
}

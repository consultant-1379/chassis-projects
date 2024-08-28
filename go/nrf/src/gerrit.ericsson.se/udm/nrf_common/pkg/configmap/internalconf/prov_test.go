package internalconf

import (
	"encoding/json"
	"testing"
)

func TestParseProvInternalConf(t *testing.T) {
	jsonContent := []byte(`{
    "httpServer": {
        "httpWithXVersion": true
    },
    "provision": {
        "syncNFProfileTimer": 1000
    }
}`)

	internalConf := &InternalProvConf{}

	err := json.Unmarshal(jsonContent, internalConf)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	internalConf.ParseConf()

	if !HTTPWithXVersion {
		t.Fatalf("InternalProvConf.ParseConf parse error")
	}

	if SyncNFProfileTimer != 1000 {
		t.Fatalf("InternalProvConf.ParseConf parse error")
	}
}

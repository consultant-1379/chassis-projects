package utils

import (
	"testing"
)

func TestJsonFormatter(t *testing.T) {
	format1 := []byte(`{
  "validityPeriod" : 1152921504606846981 ,
  "nfInstances" : [{"fqdn":"seliius03696.seli.gic.ericsson.se","nfInstanceId":"udm-5g-01"}]}
`)
	format2 := []byte(`{
  "validityPeriod" : 1152921504606846981,
  "nfInstances" : [
   {
      "fqdn":"seliius03696.seli.gic.ericsson.se",
	  "nfInstanceId":"udm-5g-01"
	}
   ]
  }
`)
	format3 := []byte(`{
  "validityPeriod" : 1152921504606846981}`)

	f1, _ := JsonFormatter(format1)
	f2, _ := JsonFormatter(format2)
	f3, _ := JsonFormatter(format3)

	if string(f1[:]) != string(f2[:]) {
		t.Fatalf("Should be the same, but NOT")
	}

	if string(f1[:]) == string(f3[:]) {
		t.Fatalf("Should NOT be the same, but the same")
	}

}

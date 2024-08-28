package swagger

import (
	"fmt"
	"testing"
)

var specJSON = `{
	 "paths" : {
		"/ue-authentications" : {
			"post" : {
				"summary" : "Authentication Initiation Request",
				"operationId" : "AI",
				"requestBody" : {
				  "content" : {
					"application/json" : {
					  "schema" : {
						"$ref" : "#/components/schemas/AuthenticationInfo"
					  }
					}
				  }
				}
			}
		},
		"/ue-authentications/{authCtxId}/5g-aka-confirmation" : {
		}
	},
	"components" : {
		"schemas" : {
		  "XXX" : {
			"type" : "object",
			"properties" : {
			  "authType" : {
				"type" : "string",
				"enum" : [ "5G-AKA", "EAP-AKA-PRIME" ]
			  },
			  "eapPayload" : {
			    "$ref" : "#/components/schemas/EapPayload"
			  }
			},
			"required" : [ "authType" ]
		  },

			"EapPayload" : {
					"type" : "string",
					"format" : "byte"
			}
		}
	}

}`

func TestSpecInString(t *testing.T) {
	a, err := DecodeSpec(specJSON)
	if err != nil {
		t.Error(err.Error())
	}
	for _, v := range a {
		fmt.Println(v.String())
	}
}

func TestDecodeSpecFileInYaml(t *testing.T) {
	DecodeSpecFile("test.yaml",
		"tmp", "test", "v1")
}

func TestDecodeSpecFileInJson(t *testing.T) {
	DecodeSpecFile("test.json",
		"tmp", "test", "v2")
}

func TestDecodeSpecFileInYaml2(t *testing.T) {
	DecodeSpecFile("nrf.yaml",
		"tmp", "test", "v1")
}

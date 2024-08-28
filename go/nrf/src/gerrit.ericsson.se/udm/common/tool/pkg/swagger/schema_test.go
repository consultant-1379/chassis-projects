package swagger

import (
	"fmt"
	"testing"
)

func TestDecodeSchema1(t *testing.T) {
	p := `{
		"type" : "object",
		"properties" : {
			"authType" : {
			"type" : "string",
			"enum" : [ "5G-AKA", "EAP-AKA-PRIME" ]
			}
		},
		"required" : [ "authType" ]
	}`

	r := DecodeSchema("AuthType", p)

	fmt.Println(r)
}

func TestDecodeSchema2(t *testing.T) {
	p := `{
			"type" : "string",
			"format" : "byte"
		  }`

	r := DecodeSchema("EapPayload", p)

	fmt.Println(r.String())

	//want := "\ntype EapPayload string"
	//
	//if r != want {
	//	t.Errorf("want %v, get %v", want, r)
	//}
}

func TestDecodeSchema3(t *testing.T) {
	p := `{
		"type" : "object",
		"properties" : {
			"servingNetworkName" : { 
				"type" : "string"
			},
			"accessType" : {
			  "$ref" : "#/components/schemas/AccessType"
			}
		},
		"required" : [ "accessType" ]
	}`

	r := DecodeSchema("AuthenticationInfoRequest", p)

	fmt.Println(r.String())
}

func TestDecodeSchema4(t *testing.T) {
	p := `{
        "required": [
          "type"
        ],
        "properties": {
          "type": {
            "type": "string"
          },
          "title": {
            "type": "string"
          },
          "status": {
            "type": "integer"
          },
          "detail": {
            "type": "string"
          },
          "instance": {
            "type": "string"
          }
        }
      }`

	r := DecodeSchema("ProblemDetails", p)

	fmt.Println(r.String())
}

func TestDecodeSchema5(t *testing.T) {
	p := `{
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "rand": {
              "$ref": "#/components/schemas/Rand"
            }
		}}}`

	r := DecodeSchema("Av5GAka", p)

	fmt.Println(r.String())
}

func TestDecodeSchema6(t *testing.T) {
	p := `{
			"required": [
				"nfInstanceID",
				"nfType"
			],
			"properties": {
				"nfInstanceID": {
					"type": "string"
				},
				"nfType": {
					"$ref": "#/components/schemas/NFType"
				},
				"plmn": {
					"type": "string"
				},
				"sNssai": {
					"$ref": "#/components/schemas/SingleNssai"
				},
				"fqdn": {
					"type": "string"
				},
				"ipAddress": {
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"capacity": {
					"type": "integer"
				},
				"nfServiceList": {
					"type": "array",
					"items": {
						"$ref": "#/components/schemas/NFService"
					}
				}
			}}`

	r := DecodeSchema("NFProfile", p)

	fmt.Println(r.String())
}

func TestDecodeSchema7(t *testing.T) {
	p := `{
			"properties": {

				"nfServiceList": {
					"type": "array",
					"items": {
						"$ref": "#/components/schemas/NFService"
					}
				}
			}}`

	r := DecodeSchema("NFProfile", p)

	fmt.Println(r.String())
}

func TestDecodeSchema8(t *testing.T) {
	p := `{
			"required": [
				"nfInstanceID",
				"nfType"
			],
			"properties": {
				"fqdn": {
					"type": "string"
				},
				"ipAddress": {
					"type": "array",
					"items": {
						"type": "string"
					}
				},
				"capacity": {
					"type": "integer"
				}
			}}`

	r := DecodeSchema("NFProfile", p)

	fmt.Println(r.String())
}

func TestDecodeSchema9(t *testing.T) {
	p := `{
		"properties": {
			"notificationBody": {
				"oneOf": [
					{
						"$ref": "#/components/schemas/Registration"
					},
					{
						"$ref": "#/components/schemas/ProfileChange"
					},
					{
						"$ref": "#/components/schemas/Deregistration"
					}
				]
			}
		}
	}`

	r := DecodeSchema("NotificationData", p)

	fmt.Println(r.String())
}

func TestDecodeSchema10(t *testing.T) {
	p := `{
		"properties": {
			"value": {
				"type": "object"
			}
		}
	}`

	r := DecodeSchema("PatchItem", p)

	fmt.Println(r.String())
}

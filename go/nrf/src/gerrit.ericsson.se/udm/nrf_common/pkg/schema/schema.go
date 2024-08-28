package schema

import (
	"fmt"
	"os"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"github.com/xeipuuv/gojsonschema"
)

var (
	schemaNfProfile         *gojsonschema.Schema
	schemaPatchDocument     *gojsonschema.Schema
	schemaSubscriptionData  *gojsonschema.Schema
	schemaSubscriptionPatch *gojsonschema.Schema
)

// LoadManagementSchema loads schema file of NF Profile and NfSubscriptionData and ...
func LoadManagementSchema() error {
	prefix := "file://"
	schemaDir := os.Getenv("SCHEMA_DIR")
	schemaNfProfileFileName := prefix + schemaDir + "/" + os.Getenv("SCHEMA_NF_PROFILE")
	schemaPatchDocumentFileName := prefix + schemaDir + "/" + os.Getenv("SCHEMA_PATCH_DOCUMENT")
	schemaSubscriptionDataFileName := prefix + schemaDir + "/" + os.Getenv("SCHEMA_SUBSCRIPTIONDATA")
	schemaSubscriptionPatchFileName := prefix + schemaDir + "/" + os.Getenv("SCHEMA_SUBSCRIPTIONPATCH")

	schemaNfProfileLoader := gojsonschema.NewReferenceLoader(schemaNfProfileFileName)
	schemaPatchDocumentLoader := gojsonschema.NewReferenceLoader(schemaPatchDocumentFileName)
	schemaSubscriptionDataLoader := gojsonschema.NewReferenceLoader(schemaSubscriptionDataFileName)
	schemaSubscriptionPatchLoader := gojsonschema.NewReferenceLoader(schemaSubscriptionPatchFileName)

	var err error
	schemaNfProfile, err = gojsonschema.NewSchema(schemaNfProfileLoader)
	if err != nil {
		log.Errorf("create schema for NFProfile failed. error info: %v", err)
		return fmt.Errorf("create schema for NFProfile failed. error info: %v", err)
	}

	schemaPatchDocument, err = gojsonschema.NewSchema(schemaPatchDocumentLoader)
	if err != nil {
		log.Errorf("create schema for PatchDocument failed. error info: %v", err)
		return fmt.Errorf("create schema for PatchDocument failed. error info: %v", err)
	}

	schemaSubscriptionData, err = gojsonschema.NewSchema(schemaSubscriptionDataLoader)
	if err != nil {
		log.Errorf("create schema for SubscriptionData failed. error info: %v", err)
		return fmt.Errorf("create schema for SubscriptionData failed. error info: %v", err)
	}

	schemaSubscriptionPatch, err = gojsonschema.NewSchema(schemaSubscriptionPatchLoader)
	if err != nil {
		log.Errorf("create schema for SubscriptionPatch failed. error info: %v", err)
		return fmt.Errorf("create schema for SubscriptionPatch failed. error info: %v", err)
	}

	return nil

}

func validateJSONContent(schema *gojsonschema.Schema, jsonLoader gojsonschema.JSONLoader) *problemdetails.ProblemDetails {
	result, err := schema.Validate(jsonLoader)
	if err != nil {
		errorInfo := fmt.Sprintf("%v", err)
		log.Errorf("schema.Validate failed. error info: %s", errorInfo)
		return &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
	}

	if !result.Valid() {
		var invalidParams []*problemdetails.InvalidParam
		mapInvalidParam := make(map[string]string)
		errorInfo := ""
		for _, item := range result.Errors() {
			desc := fmt.Sprintf("%s", item)
			errorInfo = errorInfo + fmt.Sprintf("- %s\n", desc)
			var key, value string
			index0 := strings.Index(desc, ":")
			if -1 != index0 {
				key = desc[0:strings.Index(desc, ":")]
				value = strings.TrimLeft(desc[index0+1:len(desc)], " ")
			} else {
				key = desc
				value = ""
			}
			mapInvalidParam[key] = value
		}

		log.Errorf("not a valid json. see errors: %s", errorInfo)

		for k, v := range mapInvalidParam {
			invalidParam := &problemdetails.InvalidParam{
				Param:  k,
				Reason: v,
			}
			invalidParams = append(invalidParams, invalidParam)
		}

		return &problemdetails.ProblemDetails{
			InvalidParams: invalidParams,
		}
	}
	return nil

}

// ValidateNfProfile validate the NF Profile
func ValidateNfProfile(nfProfile string) *problemdetails.ProblemDetails {
	jsonLoader := gojsonschema.NewStringLoader(nfProfile)
	return validateJSONContent(schemaNfProfile, jsonLoader)
}

// ValidatePatchDocument validate the Patch data
func ValidatePatchDocument(patchDocument string) *problemdetails.ProblemDetails {
	jsonLoader := gojsonschema.NewStringLoader(patchDocument)
	return validateJSONContent(schemaPatchDocument, jsonLoader)
}

// ValidateSubscriptionData validate the Subscription data
func ValidateSubscriptionData(subscriptionData string) *problemdetails.ProblemDetails {
	jsonLoader := gojsonschema.NewStringLoader(subscriptionData)
	return validateJSONContent(schemaSubscriptionData, jsonLoader)
}

// ValidateSubscriptionPatch validate the Subscription patch data
func ValidateSubscriptionPatch(subscriptionPatchData string) *problemdetails.ProblemDetails {
	jsonLoader := gojsonschema.NewStringLoader(subscriptionPatchData)
	return validateJSONContent(schemaSubscriptionPatch, jsonLoader)
}

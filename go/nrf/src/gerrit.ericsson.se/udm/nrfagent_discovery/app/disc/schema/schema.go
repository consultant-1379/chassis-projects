package schema

import (
	"fmt"
	"os"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/xeipuuv/gojsonschema"
)

var (
	schemaNfProfile     *gojsonschema.Schema
	schemaPatchDocument *gojsonschema.Schema
)

// LoadDiscoverSchema loads schema file of NF Profile
func LoadDiscoverSchema() error {
	prefix := "file://"
	schemaDir := os.Getenv("SCHEMA_DIR")
	schemaNfProfileFileName := prefix + schemaDir + "/" + os.Getenv("SCHEMA_NF_PROFILE")
	schemaPatchDocumentFileName := prefix + schemaDir + "/" + os.Getenv("SCHEMA_PATCH_DOCUMENT")

	schemaNfProfileLoader := gojsonschema.NewReferenceLoader(schemaNfProfileFileName)
	schemaPatchDocumentLoader := gojsonschema.NewReferenceLoader(schemaPatchDocumentFileName)

	var err error
	schemaNfProfile, err = gojsonschema.NewSchema(schemaNfProfileLoader)
	if err != nil {
		log.Errorf("create schema for nfProfileInSearchResult failed. error info: %v", err)
		return fmt.Errorf("create schema for nfProfileInSearchResult failed. error info: %v", err)
	}

	schemaPatchDocument, err = gojsonschema.NewSchema(schemaPatchDocumentLoader)
	if err != nil {
		log.Errorf("create schema for PatchDocument failed. error info: %v", err)
		return fmt.Errorf("create schema for PatchDocument failed. error info: %v", err)
	}

	log.Debugf("LoadDiscoverSchema: load schema successfully")
	return nil

}

func validateJSONContent(schema *gojsonschema.Schema, jsonLoader gojsonschema.JSONLoader) error {
	result, err := schema.Validate(jsonLoader)
	if err != nil {
		errorInfo := fmt.Sprintf("%v", err)
		log.Errorf("validateJSONContent:schema validate failed. error info: %s", errorInfo)
		return err
	}

	if !result.Valid() {
		errorInfo := ""
		for _, item := range result.Errors() {
			desc := fmt.Sprintf("%s", item)
			errorInfo = errorInfo + fmt.Sprintf("- %s,", desc)
		}
		log.Errorf("validateJSONContent: not a valid json. see errors: %s", errorInfo)
		return fmt.Errorf("validateJSONContent: not a valid json. see errors: %s", errorInfo)
	}
	return nil

}

// ValidateNfProfile validate the NF Profile
func ValidateNfProfile(nfProfile string) error {
	jsonLoader := gojsonschema.NewStringLoader(nfProfile)
	return validateJSONContent(schemaNfProfile, jsonLoader)
}

// ValidatePatchDocument validate the Patch data
func ValidatePatchDocument(patchDocument string) error {
	jsonLoader := gojsonschema.NewStringLoader(patchDocument)
	return validateJSONContent(schemaPatchDocument, jsonLoader)
}

// SetSchemaNfProfile for setting schemaNfProfile
func SetSchemaNfProfile(schema *gojsonschema.Schema) bool {
	if schema == nil {
		return false
	}
	schemaNfProfile = schema
	return true
}

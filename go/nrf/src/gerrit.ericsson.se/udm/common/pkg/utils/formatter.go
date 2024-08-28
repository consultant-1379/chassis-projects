package utils

import (
	"bytes"
	"encoding/json"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

func ToPrettyJSON(src []byte) string {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, src, "", "\t")
	if error != nil {
		log.Error("JSON parse error: ", error)
		return ""
	}

	return string(prettyJSON.Bytes())
}

func JsonFormatter(src []byte) ([]byte, error) {

	var data map[string]interface{}
	br := bytes.NewReader([]byte(src))
	jd := json.NewDecoder(br)
	jd.UseNumber()
	if err := jd.Decode(&data); err != nil {
		log.Error("JSON parse error: ", err)
		return []byte{}, err
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error("JSON MarshalIndent error: ", err)
		return []byte{}, err
	}
	b2 := append(b, '\n')
	return b2, nil
}

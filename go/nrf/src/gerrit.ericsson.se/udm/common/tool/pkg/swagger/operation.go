package swagger

import (
	"encoding/json"
	"fmt"
)

func decodeOperation(op, s string) error {

	//fmt.Printf("op(%v), s(%v)", op, s)

	var data map[string]*json.RawMessage
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		fmt.Errorf("error: %v", err.Error())
		return err
	}

	//fmt.Printf("operationId (%v)", string(*data["operationId"]))

	//decode requestBody
	//rb := data["requestBody"]
	//
	//type Fields struct {
	//	Content struct {
	//		Format struct {
	//			Schema string `json:"schema"`
	//		} `json:"application/json"`
	//	}  `json:"content"`
	//}

	//decode responses

	return nil
}

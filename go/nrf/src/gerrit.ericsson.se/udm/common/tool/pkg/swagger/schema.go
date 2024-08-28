package swagger

import (
	"encoding/json"
	"fmt"
	"strings"
)

var JsonType2GoType = map[string]string{
	"string":      "string",
	"integer":     "int",
	"boolean":     "bool",
	"array":       "array",
	"object":      "struct",
	"number":      "float", //todo, not support number yet
	"interface{}": "interface{}",
}

var GoType2JsonType = map[string]string{
	"string":      "string",
	"int":         "integer",
	"bool":        "boolean",
	"struct":      "object",
	"interface{}": "interface{}",
}

// GetGoType return the Golang Type of the input Json/openapi Type
func GetGoType(jsonType string) string {

	if t, ok := JsonType2GoType[jsonType]; ok {
		return t
	} else {
		return fmt.Sprintf("unknown json type %s", t)
	}

	return ""
}

// GetJsonType return the Json/openapi type of the input Golang type
func GetJsonType(goType string) string {

	if t, ok := GoType2JsonType[goType]; ok {
		return t
	} else {
		return fmt.Sprintf("unknown json type %s", t)
	}

	return ""
}

// OpenApiField stands for an attribute/field in openapi struct
//   <Name> <Type> <Tag>
type OpenApiField struct {
	Name         string // name of the filed
	Type         string //
	Tag          string
	IsStructType bool // true if the type refers to a struct
	IsMust       bool // true if Name is in the required list
	IsArray      bool // true if the filed has array decorator:
	//			case 1) the field is array
	//			case 2) the Type is array  --> set it through scanning the type definition
	IsMultiType bool // true if field is oneof or anyof
}

// OpenApiStruct stands for an openapi type definition
// 		type <Name> struct  { <fields> }
//      type <Name> <Type>
type OpenApiStruct struct {
	Name string // name of the type definition

	Type string // built in type or "struct"

	Fields []OpenApiField // Fields are valid only if the type is a struct

	IsArrayType bool // decorator: the type is an array instead of an simple type

}

// IsJsonStructBlob return true is the input json string refer to an openapi struct definition
func IsJsonStuctBlob(b string) bool {

	if strings.ContainsAny(b, "properties") {
		return true
	}

	return false
}

// DecodeSchema() decode the input schema of json format, then return an OpenApiStruct object.
func DecodeSchema(name string, p string) OpenApiStruct {

	d, err := DecodeJsonMap([]byte(p))
	if err != nil {
		fmt.Println(err.Error())
	}

	t := ""
	if v, ok := d["type"]; ok {
		t = strings.Trim(string(*v), " \"")
	} else {
		t = "object"
	}

	if _, ok := d["oneOf"]; ok {
		t = "interface{}"
	}
	if _, ok := d["anyOf"]; ok {
		t = "interface{}"
	}

	switch t {
	case "object":
		return decodeObject(name, p)
	case "array":
		if IsJsonStuctBlob(p) {
			return decodeArray(name, p)
		} else {
			fmt.Println("not an valid struct string")
		}

	default:
	}

	// It is a trivial struct
	r := OpenApiStruct{
		Name: addTypeDecorator(name),
		Type: GetGoType(t),
	}

	return r
}

// NewOpenApiStruct() return a OpenApiStruct object
func NewOpenApiStruct(name string, fields map[string]OpenApiField) OpenApiStruct {
	var t OpenApiStruct

	t.Name = addTypeDecorator(name)
	t.Type = "struct"

	for _, v := range fields {
		t.Fields = append(t.Fields, v)
	}

	return t
}

// String() OpenApiStruct to string
func (t *OpenApiStruct) String() string {
	if t.Type != "struct" {
		str := fmt.Sprintf("\n// %v type\n", t.Name)
		str += fmt.Sprintf("type %v %v \n", t.Name, t.Type)
		return str
	}

	var str string
	if t.Fields != nil {
		str = fmt.Sprintf("// %v struct\n", t.Name)
		str += fmt.Sprintf("type %v struct { \n", t.Name)

		for _, f := range t.Fields {
			theType := f.Type

			if f.IsStructType && !f.IsMust && !f.IsArray {
				theType = "*" + theType
			}

			if f.IsArray {
				theType = "[]" + theType
			}

			str += fmt.Sprintf(
				"\t%-20v %-20v %v\n", f.Name, theType, f.Tag)

		}

		str += "\n}"
	}

	return str
}

// decodeArray return an OpenApiStruct object for openapi array schema
func decodeArray(name, p string) OpenApiStruct {
	var data struct {
		Type  string          `json:"type"`
		Items json.RawMessage `json:"items"`
	}

	err := json.Unmarshal([]byte(p), &data)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return OpenApiStruct{}
	}

	r := DecodeSchema(name, string(data.Items))

	r.IsArrayType = true

	return r
}

// decodeObject return an OpenApiStruct object for non-array openapi schema
func decodeObject(name string, p string) OpenApiStruct {
	var data struct {
		Type       string                      `json:"type"`
		Properties map[string]*json.RawMessage `json:"properties"`
		Required   []string                    `json:"required"`
	}
	err := json.Unmarshal([]byte(p), &data)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return OpenApiStruct{}
	}

	var fields = make(map[string]OpenApiField)
	for k, v := range data.Properties {
		//t,isArray := decodeFieldType(string(*v))
		f := decodeFieldType(k, string(*v))

		//fmt.Println(f)
		if _, ok := GoType2JsonType[f.Type]; !ok {
			f.IsStructType = true
		}

		theJsonObject := k
		if has(theJsonObject, data.Required) {
			f.IsMust = true
			f.Tag = fmt.Sprintf("`json:\"%v\"`", theJsonObject)
		} else {
			f.IsMust = false
			f.Tag = fmt.Sprintf("`json:\"%v,omitempty\"`", theJsonObject)
		}

		fields[k] = f
	}

	return NewOpenApiStruct(name, fields)
}

// decodeFieldType() return the field's type
func decodeFieldType(name, s string) OpenApiField {
	var data map[string]*json.RawMessage
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return OpenApiField{}
	}

	theName := addFieldNameDecorator(name)
	var f = OpenApiField{
		Name:         theName,
		Type:         "",
		Tag:          "",
		IsStructType: false,
		IsMust:       false,
		IsArray:      false,
		IsMultiType:  false,
	}

	if t, ok := data["type"]; ok {
		tmp := strings.Trim(string(*t), "\"")
		if t1 := GetGoType(tmp); t1 != "array" {
			if tmp == "object" {
				fmt.Println("field with object type (generic type) --> go type: interface{}")
				f.Type = addTypeDecorator("interface{}")
				return f
			} else {
				f.Type = addTypeDecorator(t1)
				return f
			}
		} else {
			f.IsArray = true
		}
	} else {
		f.Type = addTypeDecorator("interface{}")
	}

	for k, v := range data {
		tmp := string(*v)
		if k == "items" {
			if f.IsArray {
				tmpF := decodeFieldType(name, tmp)
				tmpF.IsArray = true
				return tmpF
			}
		}

		if k == "$ref" {
			tmp = strings.Trim(tmp, "\"")
			s := strings.Split(tmp, "/")
			f.Type = addTypeDecorator(s[len(s)-1])
			return f
		}

		if k == "oneOf" || k == "anyOf" {
			f.Type = addTypeDecorator("interface{}")
			return f
		}
	}

	return f
}

func decodeOneofType(k, s string) OpenApiField {
	var data []*json.RawMessage
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		fmt.Printf("error: %v\n", err.Error())
		return OpenApiField{}
	}

	theName := addFieldNameDecorator(k)
	var f = OpenApiField{
		Name:         theName,
		Type:         "",
		Tag:          "",
		IsStructType: false,
		IsMust:       false,
		IsArray:      false,
		IsMultiType:  true,
	}

	for _, subType := range data {
		tmpF := decodeFieldType(k, string(*subType))
		f.Type += tmpF.Type
		f.Type += " "
	}

	return f
}

// addFieldNameDecorator() change the first c to uppercase, add prefix "A" if the first c is a digit
func addFieldNameDecorator(s string) string {

	s = strings.Trim(s, " ")
	s = strings.Trim(s, "_")

	r := strings.ToUpper(s[0:1]) + s[1:]

	d := r[0]
	if d >= '0' && d <= '9' {
		r = "A" + r
	}

	return r
}

// addTypeDecorator() add "T" for stuct
func addTypeDecorator(s string) string {

	s = strings.Trim(s, " ")
	if _, ok := GoType2JsonType[s]; ok {
		return s
	}

	if _, ok := JsonType2GoType[s]; ok {
		return s
	}

	return "T" + s
}

// has() return true if s is in the given list; otherwise return false
func has(s string, list []string) bool {
	s = strings.Trim(s, " ")
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func DecodeJsonMap(v []byte) (map[string]*json.RawMessage, error) {
	var data map[string]*json.RawMessage
	err := json.Unmarshal(v, &data)
	if err != nil {
		fmt.Errorf("error: %v", err.Error())
		return nil, err
	} else {
		//for k, v := range data {
		//	fmt.Printf("key[%v] value[%v]\n", k, string(*v))
		//}
	}

	return data, nil
}

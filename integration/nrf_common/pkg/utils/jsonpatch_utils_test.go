package utils

import (
	"encoding/json"
	//"fmt"
	"reflect"
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/encoding/schema/nrfschema"
)

func TestGetJSONPatch(t *testing.T) {
	type argsT struct {
		base string
		a    map[string]interface{}
		b    map[string]interface{}
	}

	ignoreMap := map[string]bool{
		"load": true,
	}
	var jsonPatch []*nrfschema.NrfInfoPatchData
	tests := []struct {
		name string
		args argsT
		want []*nrfschema.NrfInfoPatchData
	}{
		// TODO: Add test cases.
		{"addOneItemAtLevelOne",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test"}, map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2"}}},
			[]*nrfschema.NrfInfoPatchData{{"add", "/allowedPlmns", "", []interface{}{"test1", "test2"}}}},
		{"addOneItemAtLevelTwo",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2"}}, map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2", "test3"}}},
			[]*nrfschema.NrfInfoPatchData{{"add", "/allowedPlmns/2", "", "test3"}}},
		{"RemoveOneItemAtLevelOne",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test"}, nil},
			[]*nrfschema.NrfInfoPatchData{{"remove", "/ApiPrefix", "", nil}}},
		{"RemoveOneItemAtLevelOne2",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test"}, map[string]interface{}{}},
			[]*nrfschema.NrfInfoPatchData{{"remove", "/ApiPrefix", "", nil}}},
		{"ReplaceOneItemAtLevelTwo",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2"}}, map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test2", "test3"}}},
			[]*nrfschema.NrfInfoPatchData{{"replace", "/allowedPlmns/0", "", "test2"}, {"replace", "/allowedPlmns/1", "", "test3"}}},
		{"ChangeOrderAtLevelTwo",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2"}}, map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test2", "test3", "test1"}}},
			[]*nrfschema.NrfInfoPatchData{{"replace", "/allowedPlmns/0", "", "test2"}, {"replace", "/allowedPlmns/1", "", "test3"}, {"add", "/allowedPlmns/2", "", "test1"}}},
		{"UpdateOneItemUnderLevelTwo",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test"}, map[string]interface{}{"ServiceList": map[string]interface{}{"ServciceID": "/auth01"}}},
			[]*nrfschema.NrfInfoPatchData{{"add", "/ServiceList", "", map[string]interface{}{"ServciceID": "/auth01"}}, {"remove", "/ApiPrefix", "", nil}}},
		{"UpdateOneItemUnderLevelTwo2",
			argsT{"/", map[string]interface{}{"ServiceList": []interface{}{"auth01", "auth02"}}, map[string]interface{}{"ServiceList": []interface{}{"auth01", map[string]interface{}{"ServciceID": "/auth01"}}}},
			[]*nrfschema.NrfInfoPatchData{{"replace", "/ServiceList/1", "", map[string]interface{}{"ServciceID": "/auth01"}}}},
		{"AddLoadUnderLevel2",
			argsT{"/", map[string]interface{}{"ServiceList": map[string]interface{}{"serviceName": "auth01"}}, map[string]interface{}{"ServiceList": map[string]interface{}{"serviceName": "auth01", "load": 3}}},
			jsonPatch},
		{"UpdateLoadUnderLevel2",
			argsT{"/", map[string]interface{}{"ServiceList": map[string]interface{}{"load": 2}},
				map[string]interface{}{"ServiceList": map[string]interface{}{"load": 3}}},
			jsonPatch},
		{"RemoveLoadUnderLevel2",
			argsT{"/", map[string]interface{}{"ServiceList": map[string]interface{}{"serviceName": "auth01", "load": 4}}, map[string]interface{}{"ServiceList": map[string]interface{}{"serviceName": "auth01"}}},
			jsonPatch},
		{"AddLoadUnderLevel1",
			argsT{"/", map[string]interface{}{"serviceName": "auth01"},
				map[string]interface{}{"serviceName": "auth01", "load": 5}},
			jsonPatch},
		{"RemoveLoadUnderLevel1",
			argsT{"/", map[string]interface{}{"serviceName": "auth01", "load": 0}, map[string]interface{}{"serviceName": "auth01"}},
			jsonPatch},
		{"UpdateLoadUnderLevel1",
			argsT{"/", map[string]interface{}{"load": 2},
				map[string]interface{}{"load": 3}},
			jsonPatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getJSONPatch(tt.args.base, tt.args.a, tt.args.b, ignoreMap)
			newGot := ConvertPatchFormat(got)
			if !reflect.DeepEqual(newGot, tt.want) {
				t.Errorf("JsonPatch = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetJSONPatch1(t *testing.T) {
	OK, patchData := GetJSONPatch(`{"id": 1,"names": ["n3"] }`, `{"id": 1, "names": ["n1", "n2"]}`, nil)
	if !OK {
		t.Fatalf("GetJSONPatch error")
	}

	jsonStr, err := json.Marshal(patchData)
	if err != nil {
		t.Fatalf("Marshal error, %v", err)
	}
	//fmt.Printf("patch=%s\n", string(jsonStr))
	if string(jsonStr) != `[{"op":"REPLACE","path":"/names/0","origValue":"n3","newValue":"n1"},{"op":"ADD","path":"/names/1","newValue":"n2"}]` {
		t.Fatalf("get JSON patch fail")
	}

	newPatchData := ConvertPatchFormat(patchData)
	newJsonStr, err := json.Marshal(newPatchData)
	if err != nil {
		t.Fatalf("Marshal error, %v", err)
	}
	if string(newJsonStr) != `[{"op":"replace","path":"/names/0","value":"n1"},{"op":"add","path":"/names/1","value":"n2"}]` {
		t.Fatalf("get JSON patch fail")
	}
}

func TestConvertPatchFormat(t *testing.T) {
	OK, patchData := GetJSONPatch(`{"id": 1,"names": ["n1", "n2"] }`, `{"id": 1, "names": ["n3"]}`, nil)
	if !OK {
		t.Fatalf("GetJSONPatch error")
	}

	jsonStr, err := json.Marshal(patchData)
	if err != nil {
		t.Fatalf("Marshal error, %v", err)
	}
	if string(jsonStr) != `[{"op":"REPLACE","path":"/names/0","origValue":"n1","newValue":"n3"},{"op":"REMOVE","path":"/names/1","origValue":"n2"}]` {
		t.Fatalf("get JSON patch fail")
	}

	newPatchData := ConvertPatchFormat(patchData)
	newJsonStr, err := json.Marshal(newPatchData)
	if err != nil {
		t.Fatalf("Marshal error, %v", err)
	}
	if string(newJsonStr) != `[{"op":"replace","path":"/names/0","value":"n3"},{"op":"remove","path":"/names/1"}]` {
		t.Fatalf("get JSON patch fail")
	}
}

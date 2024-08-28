package utils

import (
	"reflect"
	"testing"
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
	var jsonPatch JsonPatch
	tests := []struct {
		name string
		args argsT
		want JsonPatch
	}{
		// TODO: Add test cases.
		{"addOneItemAtLevelOne",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test"}, map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2"}}},
			JsonPatch{{"add", "/allowedPlmns", []interface{}{"test1", "test2"}, ""}}},
		{"addOneItemAtLevelTwo",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2"}}, map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2", "test3"}}},
			JsonPatch{{"add", "/allowedPlmns/2", "test3", ""}}},
		{"RemoveOneItemAtLevelOne",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test"}, nil},
			JsonPatch{{"remove", "/ApiPrefix", nil, ""}}},
		{"RemoveOneItemAtLevelOne2",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test"}, map[string]interface{}{}},
			JsonPatch{{"remove", "/ApiPrefix", nil, ""}}},
		{"ReplaceOneItemAtLevelTwo",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2"}}, map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test2", "test3"}}},
			JsonPatch{{"replace", "/allowedPlmns/0", "test2", ""}, {"replace", "/allowedPlmns/1", "test3", ""}}},
		{"ChangeOrderAtLevelTwo",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test1", "test2"}}, map[string]interface{}{"ApiPrefix": "/test", "allowedPlmns": []interface{}{"test2", "test3", "test1"}}},
			JsonPatch{{"replace", "/allowedPlmns/0", "test2", ""}, {"replace", "/allowedPlmns/1", "test3", ""}, {"add", "/allowedPlmns/2", "test1", ""}}},
		{"UpdateOneItemUnderLevelTwo",
			argsT{"/", map[string]interface{}{"ApiPrefix": "/test"}, map[string]interface{}{"ServiceList": map[string]interface{}{"ServciceID": "/auth01"}}},
			JsonPatch{{"add", "/ServiceList", map[string]interface{}{"ServciceID": "/auth01"}, ""}, {"remove", "/ApiPrefix", nil, ""}}},
		{"UpdateOneItemUnderLevelTwo2",
			argsT{"/", map[string]interface{}{"ServiceList": []interface{}{"auth01", "auth02"}}, map[string]interface{}{"ServiceList": []interface{}{"auth01", map[string]interface{}{"ServciceID": "/auth01"}}}},
			JsonPatch{{"replace", "/ServiceList/1", map[string]interface{}{"ServciceID": "/auth01"}, ""}}},
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
			if got := GetJSONPatch(tt.args.base, tt.args.a, tt.args.b, ignoreMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonPatch = %v, want %v", got, tt.want)
			}
		})
	}
}

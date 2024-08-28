package adpcm

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_adpcm_GetValue(t *testing.T) {
	jsonstr := `
	{"name":"smfreg","title":"smfreg config","data":{
		  "service":{"port":"9001"},
		  "common":{
			"cmport":"9080",
			"pmport":"9100",
			"logport":"9090",
			"dbproxyEndpoint":"eric-udm-dbproxy:9001",
			"healthEndpoint":"localhost:9040",
			"healthPeriod":5}}}`

	var res result
	json.Unmarshal([]byte(jsonstr), &res)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Write([]byte(jsonstr))
	}))

	cm := NewAdpCm(ts.URL, "http://localhost:9002")

	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"cmport", args{"/common/cmport"}, `"9080"`},
		{"common", args{"/common"}, `{"cmport":"9080","dbproxyEndpoint":"eric-udm-dbproxy:9001","healthEndpoint":"localhost:9040","healthPeriod":5,"logport":"9090","pmport":"9100"}`},
		{"root", args{"/"}, `{"common":{"cmport":"9080","dbproxyEndpoint":"eric-udm-dbproxy:9001","healthEndpoint":"localhost:9040","healthPeriod":5,"logport":"9090","pmport":"9100"},"service":{"port":"9001"}}`},
		{"healthPeriod", args{"/common/healthPeriod"}, "5"},
		{"notexisted", args{"/notexisted"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := cm.GetValue(tt.args.key); got != tt.want {
				t.Errorf("adpcm.GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
	ts.Close()
}

func TestNewAdpCm(t *testing.T) {
	jsonstr := `
	{"name":"smfreg","title":"smfreg config","data":{
		  "service":{"port":"9001"},
		  "common":{
			"cmport":"9010",
			"pmport":"9020",
			"logport":"9030",
			"dbproxyEndpoint":"eric-udm-dbproxy:9001",
			"healthEndpoint":"localhost:9040",
			"healthPeriod":5}}}`

	var res result
	json.Unmarshal([]byte(jsonstr), &res)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Write([]byte(jsonstr))
	}))

	type args struct {
		cmUri         string
		notifEndpoint string
	}
	tests := []struct {
		name string
		args args
		want *adpcm
	}{
		{"success", args{cmUri: ts.URL, notifEndpoint: "http://localhost:9000"}, &adpcm{cmUri: ts.URL, data: res, notifEndpoint: "http://localhost:9000"}},
		{"fail", args{cmUri: "", notifEndpoint: "http://localhost:9003"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAdpCm(tt.args.cmUri, tt.args.notifEndpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAdpCm() = %v, want %v", got, tt.want)
			}
		})
	}
	ts.Close()
}

func Test_adpcm_MonitorToReLoad(t *testing.T) {
	msg1 := `{"baseETag": "3725128f5bc7aec71d69e3e7e846abc4", "configETag": "3e594abacc716650cb9fbfba0019038d", "data": {"ericsson-udm:udm-function": {"hssIwk": {"enabled": true, "hssUri": "http://testinjector:9501"}, "roaming": {"homePlmn": [{"plmnId": "111222"}, {"plmnId": "33344"}], "roamingFunctionEnabled": false, "visitedPlmn": {"visitedPlmn": [{"plmnId": "555666"}, {"plmnId": "77788"}], "allowed": true}}, "arpf": {"a4key": "0123456789ABCDEF0123456789ABCDEF", "fSet": [{"r4": 64, "r5": 96, "r1": 64, "r2": 0, "r3": 32, "id": 0, "Op": "CDC202D5123E20F62B6D676AC72CB318"}, {"r4": 64, "r5": 96, "r1": 64, "r2": 0, "r3": 32, "id": 1, "Op": "00112233445566778899AABBCCDDEEFF"}, {"r4": 64, "r5": 96, "r1": 64, "r2": 0, "r3": 32, "id": 2, "Op": "AE54B06C38DD2A4B9947A33FEE1008BC"}, {"r4": 64, "r5": 96, "r1": 64, "r2": 0, "r3": 32, "id": 3, "Op": "AE54B06C38DD2A4B9947A33FEE1008BC"}]}, "udmUri": "1234", "ProxyInterworking": {"AlertCallbackUri": "hss.ericsson.se", "enabled": false}}}, "configName": "ericsson-udm", "event": "configUpdated"}`
	jsonstr := `
	{"name":"smfreg","title":"smfreg config","data":
	{"ericsson-udm:udm-function":
	{"hssIwk": {"enabled": true, "hssUri": "http://testinjector:9501"},
	 "roaming": {"homePlmn": [{"plmnId": "111222"}, {"plmnId": "33344"}],
				 "roamingFunctionEnabled": false,
				 "visitedPlmn": {"visitedPlmn": [{"plmnId": "555666"}, {"plmnId": "77788"}],
				                 "allowed": true}},
	 "arpf": {"a4key": "0123456789ABCDEF0123456789ABCDEF",
			  "fSet": [{"r4": 64, "r5": 96, "r1": 64, "r2": 0, "r3": 32, "id": 0, "Op": "CDC202D5123E20F62B6D676AC72CB318"},
					   {"r4": 64, "r5": 96, "r1": 64, "r2": 0, "r3": 32, "id": 1, "Op": "00112233445566778899AABBCCDDEEFF"},
					   {"r4": 64, "r5": 96, "r1": 64, "r2": 0, "r3": 32, "id": 2, "Op": "AE54B06C38DD2A4B9947A33FEE1008BC"},
					   {"r4": 64, "r5": 96, "r1": 64, "r2": 0, "r3": 32, "id": 3, "Op": "AE54B06C38DD2A4B9947A33FEE1008BC"}]},
	 "udmUri": "1234", "ProxyInterworking": {"AlertCallbackUri": "hss.ericsson.se", "enabled": false}}}}`

	type args struct {
		msg []byte
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Write([]byte(jsonstr))
	}))

	notif := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Write([]byte(jsonstr))
	}))
	tests := []struct {
		name          string
		url           string
		notifEndpoint string
		want          string
		args          args
	}{
		{"success", ts.URL, notif.URL, `{"enabled":true,"hssUri":"http://testinjector:9501"}`, args{msg: []byte(msg1)}},
		{"fail1", ts.URL, "", `{"enabled":true,"hssUri":"http://testinjector:9501"}`, args{msg: []byte(msg1)}},
		{"fail2", "", "", ``, args{msg: []byte(msg1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &adpcm{cmUri: tt.url}
			cm.MonitorToReLoad(tt.args.msg)
			v, _ := cm.GetValue("/ericsson-udm:udm-function/hssIwk")
			if v != tt.want {
				t.Errorf("%v, want %v", v, tt.want)
			}
			log.Println(v)
		})
	}
}

func Test_adpcm_reloadData(t *testing.T) {
	type fields struct {
		cmUri         string
		data          result
		notifEndpoint string
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
	}))
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"fail", fields{cmUri: ts.URL, notifEndpoint: ""}, false},
		{"fail", fields{cmUri: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &adpcm{
				cmUri:         tt.fields.cmUri,
				data:          tt.fields.data,
				notifEndpoint: tt.fields.notifEndpoint,
			}
			if got := cm.reloadData(); got != tt.want {
				t.Errorf("adpcm.reloadData() = %v, want %v", got, tt.want)
			}
		})
	}
}

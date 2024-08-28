package service

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"gerrit.ericsson.se/HSS/5G/cmproxy/service/cmadapter"
	cmgrpc "gerrit.ericsson.se/HSS/5G/protocols/sidecar/cm"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var value string

type mockAdapter struct{}

func (cm *mockAdapter) GetValue(key string) (string, error) {
	return value, nil
}

func (cm *mockAdapter) MonitorToReLoad(msg []byte) {

}

func TestServer_Read(t *testing.T) {
	type fields struct {
		cmadapter  cmadapter.Adapter
		grpcServer *grpc.Server
	}
	type args struct {
		ctx context.Context
		in  *cmgrpc.CmRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *cmgrpc.CmResponse
		wantErr bool
	}{
		{"success", fields{new(mockAdapter), nil}, args{nil, &cmgrpc.CmRequest{"/cmport"}}, &cmgrpc.CmResponse{value}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				cmadapter:  tt.fields.cmadapter,
				grpcServer: tt.fields.grpcServer,
			}
			got, err := s.Read(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Server.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {

	s := &Server{cmadapter: new(mockAdapter), grpcServer: grpc.NewServer()}

	a := strconv.Itoa(9080)
	go s.Start(a)

	timeout, _ := time.ParseDuration("5s")

	cgrpc, _ := grpc.Dial("localhost:9080", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeout))
	defer cgrpc.Close()
}

func TestNewServer(t *testing.T) {
	jsonstr := `
	{"name":"smfreg","title":"smfreg config","data":{
		"smfreg":{
		  "service":{"port":"9001"},
		  "common":{
			"cmport":"9080",
			"pmport":"9100",
			"logport":"9090",
			"dbproxyEndpoint":"eric-udm-dbproxy:9001",
			"healthEndpoint":"localhost:9040",
			"healthPeriod":5}}}}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		w.Write([]byte(jsonstr))
	}))

	type args struct {
		cmUri         string
		notifEndpoint string
	}
	tests := []struct {
		name  string
		args  args
		isNil bool
	}{
		{"success", args{ts.URL, "http://localhost:9081"}, false},
		{"fail", args{"", "http://localhost:9081"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServer(tt.args.cmUri, tt.args.notifEndpoint)

			if tt.isNil {
				if got != nil {
					t.Errorf("NewServer() failed")
				}
			} else {
				if got == nil {
					t.Errorf("NewServer() failed")
				}
			}
		})
	}
	ts.Close()
}

func TestServer_Stop(t *testing.T) {

	tests := []struct {
		name string
	}{
		{"success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{cmadapter: nil, grpcServer: grpc.NewServer()}
			s.Stop()
		})
	}
}

func TestServer_GetCB(t *testing.T) {
	mock := &mockAdapter{}
	type fields struct {
		cmadapter  cmadapter.Adapter
		grpcServer *grpc.Server
	}
	tests := []struct {
		name   string
		fields fields
		want   func([]byte)
	}{
		{"success", fields{cmadapter: mock}, mock.MonitorToReLoad},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				cmadapter:  tt.fields.cmadapter,
				grpcServer: tt.fields.grpcServer,
			}
			if got := s.GetCB(); got == nil {
				t.Errorf("Server.GetCB() = nil")
			}
		})
	}
}

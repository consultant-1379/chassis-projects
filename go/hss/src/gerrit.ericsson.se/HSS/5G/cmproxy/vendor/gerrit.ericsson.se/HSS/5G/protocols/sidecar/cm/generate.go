//go:generate protoc --proto_path=. --plugin=$GOPATH/bin/protoc-gen-go --go_out=plugins=grpc,:. cm.proto
package ericsson_udm_sidecar_cm_v1

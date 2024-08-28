// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/gpsiprefixprofile/GpsiprefixProfileGetResponse.proto

package gpsiprefixprofile // import "com/dbproxy/nfmessage/gpsiprefixprofile"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GpsiprefixProfileGetResponse struct {
	Code                 uint32   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	ValueInfo            []string `protobuf:"bytes,2,rep,name=value_info,json=valueInfo,proto3" json:"value_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GpsiprefixProfileGetResponse) Reset()         { *m = GpsiprefixProfileGetResponse{} }
func (m *GpsiprefixProfileGetResponse) String() string { return proto.CompactTextString(m) }
func (*GpsiprefixProfileGetResponse) ProtoMessage()    {}
func (*GpsiprefixProfileGetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_GpsiprefixProfileGetResponse_6192c32fddb684ab, []int{0}
}
func (m *GpsiprefixProfileGetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GpsiprefixProfileGetResponse.Unmarshal(m, b)
}
func (m *GpsiprefixProfileGetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GpsiprefixProfileGetResponse.Marshal(b, m, deterministic)
}
func (dst *GpsiprefixProfileGetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GpsiprefixProfileGetResponse.Merge(dst, src)
}
func (m *GpsiprefixProfileGetResponse) XXX_Size() int {
	return xxx_messageInfo_GpsiprefixProfileGetResponse.Size(m)
}
func (m *GpsiprefixProfileGetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GpsiprefixProfileGetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GpsiprefixProfileGetResponse proto.InternalMessageInfo

func (m *GpsiprefixProfileGetResponse) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *GpsiprefixProfileGetResponse) GetValueInfo() []string {
	if m != nil {
		return m.ValueInfo
	}
	return nil
}

func init() {
	proto.RegisterType((*GpsiprefixProfileGetResponse)(nil), "grpc.GpsiprefixProfileGetResponse")
}

func init() {
	proto.RegisterFile("nfmessage/gpsiprefixprofile/GpsiprefixProfileGetResponse.proto", fileDescriptor_GpsiprefixProfileGetResponse_6192c32fddb684ab)
}

var fileDescriptor_GpsiprefixProfileGetResponse_6192c32fddb684ab = []byte{
	// 198 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xb2, 0xcb, 0x4b, 0xcb, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x2f, 0x28, 0xce, 0x2c, 0x28, 0x4a, 0x4d, 0xcb, 0xac,
	0x28, 0x28, 0xca, 0x4f, 0xcb, 0xcc, 0x49, 0xd5, 0x77, 0x87, 0x8b, 0x04, 0x40, 0x44, 0xdc, 0x53,
	0x4b, 0x82, 0x52, 0x8b, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0xf5, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85,
	0x58, 0xd2, 0x8b, 0x0a, 0x92, 0x95, 0x02, 0xb9, 0x64, 0xf0, 0xa9, 0x15, 0x12, 0xe2, 0x62, 0x49,
	0xce, 0x4f, 0x49, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0d, 0x02, 0xb3, 0x85, 0x64, 0xb9, 0xb8,
	0xca, 0x12, 0x73, 0x4a, 0x53, 0xe3, 0x33, 0xf3, 0xd2, 0xf2, 0x25, 0x98, 0x14, 0x98, 0x35, 0x38,
	0x83, 0x38, 0xc1, 0x22, 0x9e, 0x79, 0x69, 0xf9, 0x4e, 0x1d, 0x8c, 0x5c, 0x56, 0xa9, 0x45, 0x99,
	0xc9, 0xc5, 0xc5, 0xf9, 0x79, 0x7a, 0xc9, 0xf9, 0x45, 0xa9, 0x7a, 0x79, 0x45, 0x69, 0x7a, 0x29,
	0x49, 0x05, 0x45, 0xf9, 0x15, 0x95, 0x7a, 0x20, 0x5b, 0xf5, 0xe0, 0x4e, 0xd7, 0xc3, 0x70, 0xba,
	0x93, 0x22, 0x3e, 0xf7, 0x04, 0x80, 0x9c, 0x1e, 0xa5, 0x9e, 0x9c, 0x9f, 0xab, 0x0f, 0x35, 0x50,
	0x1f, 0x4f, 0x30, 0x24, 0xb1, 0x81, 0xbd, 0x6a, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x93, 0xb0,
	0x68, 0x4a, 0x2c, 0x01, 0x00, 0x00,
}

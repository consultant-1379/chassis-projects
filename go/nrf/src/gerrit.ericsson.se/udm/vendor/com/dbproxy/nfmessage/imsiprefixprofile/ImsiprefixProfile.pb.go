// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/imsiprefixprofile/ImsiprefixProfile.proto

package imsiprefixprofile // import "com/dbproxy/nfmessage/imsiprefixprofile"

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

type ImsiprefixProfile struct {
	ImsiPrefix           uint64   `protobuf:"varint,1,opt,name=imsi_prefix,json=imsiPrefix,proto3" json:"imsi_prefix,omitempty"`
	ValueInfo            string   `protobuf:"bytes,2,opt,name=value_info,json=valueInfo,proto3" json:"value_info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImsiprefixProfile) Reset()         { *m = ImsiprefixProfile{} }
func (m *ImsiprefixProfile) String() string { return proto.CompactTextString(m) }
func (*ImsiprefixProfile) ProtoMessage()    {}
func (*ImsiprefixProfile) Descriptor() ([]byte, []int) {
	return fileDescriptor_ImsiprefixProfile_cee329b50c2b6d7a, []int{0}
}
func (m *ImsiprefixProfile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImsiprefixProfile.Unmarshal(m, b)
}
func (m *ImsiprefixProfile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImsiprefixProfile.Marshal(b, m, deterministic)
}
func (dst *ImsiprefixProfile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImsiprefixProfile.Merge(dst, src)
}
func (m *ImsiprefixProfile) XXX_Size() int {
	return xxx_messageInfo_ImsiprefixProfile.Size(m)
}
func (m *ImsiprefixProfile) XXX_DiscardUnknown() {
	xxx_messageInfo_ImsiprefixProfile.DiscardUnknown(m)
}

var xxx_messageInfo_ImsiprefixProfile proto.InternalMessageInfo

func (m *ImsiprefixProfile) GetImsiPrefix() uint64 {
	if m != nil {
		return m.ImsiPrefix
	}
	return 0
}

func (m *ImsiprefixProfile) GetValueInfo() string {
	if m != nil {
		return m.ValueInfo
	}
	return ""
}

func init() {
	proto.RegisterType((*ImsiprefixProfile)(nil), "grpc.ImsiprefixProfile")
}

func init() {
	proto.RegisterFile("nfmessage/imsiprefixprofile/ImsiprefixProfile.proto", fileDescriptor_ImsiprefixProfile_cee329b50c2b6d7a)
}

var fileDescriptor_ImsiprefixProfile_cee329b50c2b6d7a = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0xce, 0x4b, 0xcb, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0xcf, 0xcc, 0x2d, 0xce, 0x2c, 0x28, 0x4a, 0x4d, 0xcb, 0xac,
	0x28, 0x28, 0xca, 0x4f, 0xcb, 0xcc, 0x49, 0xd5, 0xf7, 0x84, 0x8b, 0x04, 0x40, 0x44, 0xf4, 0x0a,
	0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x58, 0xd2, 0x8b, 0x0a, 0x92, 0x95, 0x82, 0xb9, 0x04, 0x31, 0x14,
	0x08, 0xc9, 0x73, 0x71, 0x83, 0xcc, 0x89, 0x87, 0x88, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x04,
	0x71, 0x81, 0x84, 0x02, 0xc0, 0x22, 0x42, 0xb2, 0x5c, 0x5c, 0x65, 0x89, 0x39, 0xa5, 0xa9, 0xf1,
	0x99, 0x79, 0x69, 0xf9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x9c, 0x60, 0x11, 0xcf, 0xbc,
	0xb4, 0x7c, 0xa7, 0x5a, 0x2e, 0xab, 0xd4, 0xa2, 0xcc, 0xe4, 0xe2, 0xe2, 0xfc, 0x3c, 0xbd, 0xe4,
	0xfc, 0xa2, 0x54, 0xbd, 0xbc, 0xa2, 0x34, 0xbd, 0x94, 0xa4, 0x82, 0xa2, 0xfc, 0x8a, 0x4a, 0x3d,
	0x90, 0xb5, 0x7a, 0x70, 0x07, 0xeb, 0x61, 0x38, 0xd8, 0x49, 0x0c, 0xc3, 0x41, 0x01, 0x20, 0x07,
	0x47, 0xa9, 0x27, 0xe7, 0xe7, 0xea, 0x43, 0x4d, 0xd1, 0xc7, 0xe3, 0xe3, 0x24, 0x36, 0xb0, 0x07,
	0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x24, 0xb5, 0x93, 0xf3, 0x17, 0x01, 0x00, 0x00,
}

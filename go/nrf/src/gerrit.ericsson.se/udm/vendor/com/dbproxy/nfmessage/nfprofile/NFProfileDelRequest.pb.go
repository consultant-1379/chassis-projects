// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/nfprofile/NFProfileDelRequest.proto

package nfprofile // import "com/dbproxy/nfmessage/nfprofile"

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

type NFProfileDelRequest struct {
	NfInstanceId         string   `protobuf:"bytes,1,opt,name=nf_instance_id,json=nfInstanceId,proto3" json:"nf_instance_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NFProfileDelRequest) Reset()         { *m = NFProfileDelRequest{} }
func (m *NFProfileDelRequest) String() string { return proto.CompactTextString(m) }
func (*NFProfileDelRequest) ProtoMessage()    {}
func (*NFProfileDelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_NFProfileDelRequest_15c3cb549118086a, []int{0}
}
func (m *NFProfileDelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NFProfileDelRequest.Unmarshal(m, b)
}
func (m *NFProfileDelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NFProfileDelRequest.Marshal(b, m, deterministic)
}
func (dst *NFProfileDelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NFProfileDelRequest.Merge(dst, src)
}
func (m *NFProfileDelRequest) XXX_Size() int {
	return xxx_messageInfo_NFProfileDelRequest.Size(m)
}
func (m *NFProfileDelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NFProfileDelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NFProfileDelRequest proto.InternalMessageInfo

func (m *NFProfileDelRequest) GetNfInstanceId() string {
	if m != nil {
		return m.NfInstanceId
	}
	return ""
}

func init() {
	proto.RegisterType((*NFProfileDelRequest)(nil), "grpc.NFProfileDelRequest")
}

func init() {
	proto.RegisterFile("nfmessage/nfprofile/NFProfileDelRequest.proto", fileDescriptor_NFProfileDelRequest_15c3cb549118086a)
}

var fileDescriptor_NFProfileDelRequest_15c3cb549118086a = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xcd, 0x4b, 0xcb, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0xcf, 0x4b, 0x2b, 0x28, 0xca, 0x4f, 0xcb, 0xcc, 0x49, 0xd5,
	0xf7, 0x73, 0x0b, 0x80, 0xb0, 0x5c, 0x52, 0x73, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0xf4,
	0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x58, 0xd2, 0x8b, 0x0a, 0x92, 0x95, 0xac, 0xb9, 0x84, 0xb1,
	0x28, 0x11, 0x52, 0xe1, 0xe2, 0xcb, 0x4b, 0x8b, 0xcf, 0xcc, 0x2b, 0x2e, 0x49, 0xcc, 0x4b, 0x4e,
	0x8d, 0xcf, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0xe2, 0xc9, 0x4b, 0xf3, 0x84, 0x0a,
	0x7a, 0xa6, 0x38, 0xe5, 0x73, 0x19, 0xa5, 0x16, 0x65, 0x26, 0x17, 0x17, 0xe7, 0xe7, 0xe9, 0x25,
	0xe7, 0x17, 0xa5, 0xea, 0xe5, 0x15, 0xa5, 0xe9, 0xa5, 0x24, 0x15, 0x14, 0xe5, 0x57, 0x54, 0xea,
	0x81, 0x8c, 0xd7, 0x83, 0x3b, 0x49, 0x0f, 0xee, 0x24, 0x27, 0x09, 0x2c, 0x16, 0x06, 0x80, 0x9c,
	0x14, 0x25, 0x9f, 0x9c, 0x9f, 0xab, 0x0f, 0xd5, 0xaf, 0x8f, 0xc5, 0x37, 0x49, 0x6c, 0x60, 0xa7,
	0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb4, 0xcd, 0xdb, 0x26, 0xeb, 0x00, 0x00, 0x00,
}

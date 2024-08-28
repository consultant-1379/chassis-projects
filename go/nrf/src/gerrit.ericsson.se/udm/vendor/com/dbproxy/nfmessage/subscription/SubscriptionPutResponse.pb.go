// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/subscription/SubscriptionPutResponse.proto

package subscription // import "com/dbproxy/nfmessage/subscription"

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

type SubscriptionPutResponse struct {
	Code                 uint32   `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscriptionPutResponse) Reset()         { *m = SubscriptionPutResponse{} }
func (m *SubscriptionPutResponse) String() string { return proto.CompactTextString(m) }
func (*SubscriptionPutResponse) ProtoMessage()    {}
func (*SubscriptionPutResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_SubscriptionPutResponse_d2f48a316bb2df27, []int{0}
}
func (m *SubscriptionPutResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscriptionPutResponse.Unmarshal(m, b)
}
func (m *SubscriptionPutResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscriptionPutResponse.Marshal(b, m, deterministic)
}
func (dst *SubscriptionPutResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscriptionPutResponse.Merge(dst, src)
}
func (m *SubscriptionPutResponse) XXX_Size() int {
	return xxx_messageInfo_SubscriptionPutResponse.Size(m)
}
func (m *SubscriptionPutResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscriptionPutResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SubscriptionPutResponse proto.InternalMessageInfo

func (m *SubscriptionPutResponse) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func init() {
	proto.RegisterType((*SubscriptionPutResponse)(nil), "grpc.SubscriptionPutResponse")
}

func init() {
	proto.RegisterFile("nfmessage/subscription/SubscriptionPutResponse.proto", fileDescriptor_SubscriptionPutResponse_d2f48a316bb2df27)
}

var fileDescriptor_SubscriptionPutResponse_d2f48a316bb2df27 = []byte{
	// 161 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0xc9, 0x4b, 0xcb, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x2e, 0x4d, 0x2a, 0x4e, 0x2e, 0xca, 0x2c, 0x28, 0xc9,
	0xcc, 0xcf, 0xd3, 0x0f, 0x46, 0xe2, 0x04, 0x94, 0x96, 0x04, 0xa5, 0x16, 0x17, 0xe4, 0xe7, 0x15,
	0xa7, 0xea, 0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b, 0xb1, 0xa4, 0x17, 0x15, 0x24, 0x2b, 0xe9, 0x72,
	0x89, 0xe3, 0x50, 0x26, 0x24, 0xc4, 0xc5, 0x92, 0x9c, 0x9f, 0x92, 0x2a, 0xc1, 0xa8, 0xc0, 0xa8,
	0xc1, 0x1b, 0x04, 0x66, 0x3b, 0x55, 0x72, 0x99, 0xa6, 0x16, 0x65, 0x26, 0x17, 0x17, 0xe7, 0xe7,
	0xe9, 0x25, 0xe7, 0x17, 0xa5, 0xea, 0xe5, 0x15, 0xa5, 0xe9, 0xa5, 0x24, 0x15, 0x14, 0xe5, 0x57,
	0x54, 0xea, 0x81, 0x0c, 0xd4, 0x83, 0xbb, 0x45, 0x0f, 0xd9, 0x2d, 0x4e, 0x32, 0x38, 0x6c, 0x09,
	0x00, 0xb9, 0x25, 0x4a, 0x29, 0x39, 0x3f, 0x57, 0x1f, 0x6a, 0x8c, 0x3e, 0x76, 0xdf, 0x24, 0xb1,
	0x81, 0x9d, 0x6d, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xf0, 0x3f, 0x92, 0xc5, 0xee, 0x00, 0x00,
	0x00,
}

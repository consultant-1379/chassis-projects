// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/nrfaddress/NRFAddressPutRequest.proto

package nrfaddress // import "com/dbproxy/nfmessage/nrfaddress"

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

type NRFAddressPutRequest struct {
	NrfAddressId         string           `protobuf:"bytes,1,opt,name=nrf_address_id,json=nrfAddressId,proto3" json:"nrf_address_id,omitempty"`
	Index                *NRFAddressIndex `protobuf:"bytes,2,opt,name=index,proto3" json:"index,omitempty"`
	NrfAddressData       []byte           `protobuf:"bytes,3,opt,name=nrf_address_data,json=nrfAddressData,proto3" json:"nrf_address_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *NRFAddressPutRequest) Reset()         { *m = NRFAddressPutRequest{} }
func (m *NRFAddressPutRequest) String() string { return proto.CompactTextString(m) }
func (*NRFAddressPutRequest) ProtoMessage()    {}
func (*NRFAddressPutRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_NRFAddressPutRequest_8bf2ae12123e86f8, []int{0}
}
func (m *NRFAddressPutRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NRFAddressPutRequest.Unmarshal(m, b)
}
func (m *NRFAddressPutRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NRFAddressPutRequest.Marshal(b, m, deterministic)
}
func (dst *NRFAddressPutRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NRFAddressPutRequest.Merge(dst, src)
}
func (m *NRFAddressPutRequest) XXX_Size() int {
	return xxx_messageInfo_NRFAddressPutRequest.Size(m)
}
func (m *NRFAddressPutRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NRFAddressPutRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NRFAddressPutRequest proto.InternalMessageInfo

func (m *NRFAddressPutRequest) GetNrfAddressId() string {
	if m != nil {
		return m.NrfAddressId
	}
	return ""
}

func (m *NRFAddressPutRequest) GetIndex() *NRFAddressIndex {
	if m != nil {
		return m.Index
	}
	return nil
}

func (m *NRFAddressPutRequest) GetNrfAddressData() []byte {
	if m != nil {
		return m.NrfAddressData
	}
	return nil
}

func init() {
	proto.RegisterType((*NRFAddressPutRequest)(nil), "grpc.NRFAddressPutRequest")
}

func init() {
	proto.RegisterFile("nfmessage/nrfaddress/NRFAddressPutRequest.proto", fileDescriptor_NRFAddressPutRequest_8bf2ae12123e86f8)
}

var fileDescriptor_NRFAddressPutRequest_8bf2ae12123e86f8 = []byte{
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xcf, 0x4b, 0xcb, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0xcf, 0x2b, 0x4a, 0x4b, 0x4c, 0x49, 0x29, 0x4a, 0x2d, 0x2e,
	0xd6, 0xf7, 0x0b, 0x72, 0x73, 0x84, 0x30, 0x03, 0x4a, 0x4b, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b,
	0x4b, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x58, 0xd2, 0x8b, 0x0a, 0x92, 0xa5, 0xb4, 0x08,
	0x68, 0xf3, 0xcc, 0x4b, 0x49, 0xad, 0x80, 0xe8, 0x50, 0x9a, 0xcc, 0xc8, 0x25, 0x82, 0xcd, 0x40,
	0x21, 0x15, 0x2e, 0xbe, 0xbc, 0xa2, 0xb4, 0x78, 0xa8, 0xee, 0xf8, 0xcc, 0x14, 0x09, 0x46, 0x05,
	0x46, 0x0d, 0xce, 0x20, 0x9e, 0xbc, 0xa2, 0x34, 0x98, 0x39, 0x29, 0x42, 0xda, 0x5c, 0xac, 0x99,
	0x20, 0xd3, 0x24, 0x98, 0x14, 0x18, 0x35, 0xb8, 0x8d, 0x44, 0xf5, 0x40, 0x0e, 0xd0, 0x43, 0xb3,
	0x2a, 0x08, 0xa2, 0x46, 0x48, 0x83, 0x4b, 0x00, 0xd9, 0xc8, 0x94, 0xc4, 0x92, 0x44, 0x09, 0x66,
	0x05, 0x46, 0x0d, 0x9e, 0x20, 0x3e, 0x84, 0xa1, 0x2e, 0x89, 0x25, 0x89, 0x4e, 0x45, 0x5c, 0xc6,
	0xa9, 0x45, 0x99, 0xc9, 0xc5, 0xc5, 0xf9, 0x79, 0x7a, 0xc9, 0xf9, 0x45, 0xa9, 0x7a, 0x79, 0x45,
	0x69, 0x7a, 0x29, 0x49, 0x05, 0x45, 0xf9, 0x15, 0x95, 0x10, 0x2b, 0xe0, 0x5e, 0xd4, 0x43, 0x78,
	0xd1, 0x49, 0x12, 0x9b, 0x4f, 0x02, 0x40, 0xfe, 0x8c, 0x52, 0x48, 0xce, 0xcf, 0xd5, 0x87, 0x9a,
	0x80, 0x35, 0x58, 0x93, 0xd8, 0xc0, 0x01, 0x62, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xbc, 0xdd,
	0x61, 0xb4, 0x75, 0x01, 0x00, 0x00,
}

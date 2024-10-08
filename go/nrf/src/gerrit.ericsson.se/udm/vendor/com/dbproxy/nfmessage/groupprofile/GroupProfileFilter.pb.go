// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/groupprofile/GroupProfileFilter.proto

package groupprofile // import "com/dbproxy/nfmessage/groupprofile"

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

type GroupProfileFilter struct {
	AndOperation         bool               `protobuf:"varint,1,opt,name=and_operation,json=andOperation,proto3" json:"and_operation,omitempty"`
	Index                *GroupProfileIndex `protobuf:"bytes,2,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GroupProfileFilter) Reset()         { *m = GroupProfileFilter{} }
func (m *GroupProfileFilter) String() string { return proto.CompactTextString(m) }
func (*GroupProfileFilter) ProtoMessage()    {}
func (*GroupProfileFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_GroupProfileFilter_7ae8ae734e1363b3, []int{0}
}
func (m *GroupProfileFilter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GroupProfileFilter.Unmarshal(m, b)
}
func (m *GroupProfileFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GroupProfileFilter.Marshal(b, m, deterministic)
}
func (dst *GroupProfileFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupProfileFilter.Merge(dst, src)
}
func (m *GroupProfileFilter) XXX_Size() int {
	return xxx_messageInfo_GroupProfileFilter.Size(m)
}
func (m *GroupProfileFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupProfileFilter.DiscardUnknown(m)
}

var xxx_messageInfo_GroupProfileFilter proto.InternalMessageInfo

func (m *GroupProfileFilter) GetAndOperation() bool {
	if m != nil {
		return m.AndOperation
	}
	return false
}

func (m *GroupProfileFilter) GetIndex() *GroupProfileIndex {
	if m != nil {
		return m.Index
	}
	return nil
}

func init() {
	proto.RegisterType((*GroupProfileFilter)(nil), "grpc.GroupProfileFilter")
}

func init() {
	proto.RegisterFile("nfmessage/groupprofile/GroupProfileFilter.proto", fileDescriptor_GroupProfileFilter_7ae8ae734e1363b3)
}

var fileDescriptor_GroupProfileFilter_7ae8ae734e1363b3 = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xcf, 0x4b, 0xcb, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x2f, 0xca, 0x2f, 0x2d, 0x28, 0x28, 0xca, 0x4f, 0xcb,
	0xcc, 0x49, 0xd5, 0x77, 0x07, 0x71, 0x02, 0x20, 0x1c, 0xb7, 0xcc, 0x9c, 0x92, 0xd4, 0x22, 0xbd,
	0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x96, 0xf4, 0xa2, 0x82, 0x64, 0x29, 0x3d, 0x22, 0xb4, 0x79,
	0xe6, 0xa5, 0xa4, 0x56, 0x40, 0x74, 0x29, 0x65, 0x70, 0x09, 0x61, 0x9a, 0x28, 0xa4, 0xcc, 0xc5,
	0x9b, 0x98, 0x97, 0x12, 0x9f, 0x5f, 0x90, 0x5a, 0x94, 0x58, 0x92, 0x99, 0x9f, 0x27, 0xc1, 0xa8,
	0xc0, 0xa8, 0xc1, 0x11, 0xc4, 0x93, 0x98, 0x97, 0xe2, 0x0f, 0x13, 0x13, 0xd2, 0xe5, 0x62, 0xcd,
	0x04, 0x99, 0x24, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x6d, 0x24, 0xae, 0x07, 0x72, 0x80, 0x1e, 0x86,
	0x45, 0x41, 0x10, 0x55, 0x4e, 0x25, 0x5c, 0xa6, 0xa9, 0x45, 0x99, 0xc9, 0xc5, 0xc5, 0xf9, 0x79,
	0x7a, 0xc9, 0xf9, 0x45, 0xa9, 0x7a, 0x79, 0x45, 0x69, 0x7a, 0x29, 0x49, 0x05, 0x45, 0xf9, 0x15,
	0x95, 0x10, 0xad, 0x70, 0xa7, 0xeb, 0x21, 0x3b, 0xdd, 0x49, 0x1c, 0xd3, 0x81, 0x01, 0x20, 0xb7,
	0x47, 0x29, 0x25, 0xe7, 0xe7, 0xea, 0x43, 0x4d, 0xc0, 0x11, 0x5c, 0x49, 0x6c, 0x60, 0x6f, 0x1a,
	0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x82, 0xbb, 0x87, 0x41, 0x4f, 0x01, 0x00, 0x00,
}

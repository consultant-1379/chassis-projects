// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/cachenfprofile/CacheNFProfileGetRequest.proto

package cachenfprofile // import "com/dbproxy/nfmessage/cachenfprofile"

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

type CacheNFProfileGetRequest struct {
	CacheNfInstanceId    string   `protobuf:"bytes,1,opt,name=cache_nf_instance_id,json=cacheNfInstanceId,proto3" json:"cache_nf_instance_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CacheNFProfileGetRequest) Reset()         { *m = CacheNFProfileGetRequest{} }
func (m *CacheNFProfileGetRequest) String() string { return proto.CompactTextString(m) }
func (*CacheNFProfileGetRequest) ProtoMessage()    {}
func (*CacheNFProfileGetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_CacheNFProfileGetRequest_1ff675ad572fef55, []int{0}
}
func (m *CacheNFProfileGetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CacheNFProfileGetRequest.Unmarshal(m, b)
}
func (m *CacheNFProfileGetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CacheNFProfileGetRequest.Marshal(b, m, deterministic)
}
func (dst *CacheNFProfileGetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CacheNFProfileGetRequest.Merge(dst, src)
}
func (m *CacheNFProfileGetRequest) XXX_Size() int {
	return xxx_messageInfo_CacheNFProfileGetRequest.Size(m)
}
func (m *CacheNFProfileGetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CacheNFProfileGetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CacheNFProfileGetRequest proto.InternalMessageInfo

func (m *CacheNFProfileGetRequest) GetCacheNfInstanceId() string {
	if m != nil {
		return m.CacheNfInstanceId
	}
	return ""
}

func init() {
	proto.RegisterType((*CacheNFProfileGetRequest)(nil), "grpc.CacheNFProfileGetRequest")
}

func init() {
	proto.RegisterFile("nfmessage/cachenfprofile/CacheNFProfileGetRequest.proto", fileDescriptor_CacheNFProfileGetRequest_1ff675ad572fef55)
}

var fileDescriptor_CacheNFProfileGetRequest_1ff675ad572fef55 = []byte{
	// 188 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0xcf, 0x4b, 0xcb, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x4e, 0x4c, 0xce, 0x48, 0xcd, 0x4b, 0x2b, 0x28, 0xca,
	0x4f, 0xcb, 0xcc, 0x49, 0xd5, 0x77, 0x06, 0x71, 0xfd, 0xdc, 0x02, 0x20, 0x5c, 0xf7, 0xd4, 0x92,
	0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x96, 0xf4,
	0xa2, 0x82, 0x64, 0x25, 0x6f, 0x2e, 0x09, 0x5c, 0xea, 0x84, 0xf4, 0xb9, 0x44, 0xc0, 0x46, 0xc6,
	0xe7, 0xa5, 0xc5, 0x67, 0xe6, 0x15, 0x97, 0x24, 0xe6, 0x25, 0xa7, 0xc6, 0x67, 0xa6, 0x48, 0x30,
	0x2a, 0x30, 0x6a, 0x70, 0x06, 0x09, 0x82, 0xe5, 0xfc, 0xd2, 0x3c, 0xa1, 0x32, 0x9e, 0x29, 0x4e,
	0x75, 0x5c, 0xe6, 0xa9, 0x45, 0x99, 0xc9, 0xc5, 0xc5, 0xf9, 0x79, 0x7a, 0xc9, 0xf9, 0x45, 0xa9,
	0x7a, 0x79, 0x45, 0x69, 0x7a, 0x29, 0x49, 0x05, 0x45, 0xf9, 0x15, 0x95, 0x7a, 0x20, 0xeb, 0xf4,
	0xe0, 0x8e, 0xd5, 0x43, 0x75, 0xac, 0x93, 0x2c, 0x2e, 0x57, 0x04, 0x80, 0x1c, 0x1b, 0xa5, 0x92,
	0x9c, 0x9f, 0xab, 0x0f, 0x35, 0x49, 0x1f, 0x97, 0x8f, 0x93, 0xd8, 0xc0, 0x3e, 0x33, 0x06, 0x04,
	0x00, 0x00, 0xff, 0xff, 0xa2, 0x3b, 0x8e, 0x5f, 0x14, 0x01, 0x00, 0x00,
}

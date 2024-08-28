// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/nrfprofile/NRFProfileIndex.proto

package nrfprofile // import "com/dbproxy/nfmessage/nrfprofile"

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

type NRFProfileIndex struct {
	// for exprire time(Put) / start expire time(Get)
	Key1 uint64 `protobuf:"varint,1,opt,name=key1,proto3" json:"key1,omitempty"`
	//                      / end expire time(Get)
	Key2 uint64 `protobuf:"varint,2,opt,name=key2,proto3" json:"key2,omitempty"`
	// for register type(Put) / for register type 1:register 2:provision(Get)
	Key3 uint64 `protobuf:"varint,3,opt,name=key3,proto3" json:"key3,omitempty"`
	// for amfInfoSum
	AmfKey1 []*NRFKeyStruct `protobuf:"bytes,4,rep,name=amf_key1,json=amfKey1,proto3" json:"amf_key1,omitempty"`
	AmfKey2 []*NRFKeyStruct `protobuf:"bytes,5,rep,name=amf_key2,json=amfKey2,proto3" json:"amf_key2,omitempty"`
	AmfKey3 []*NRFKeyStruct `protobuf:"bytes,6,rep,name=amf_key3,json=amfKey3,proto3" json:"amf_key3,omitempty"`
	AmfKey4 []*NRFKeyStruct `protobuf:"bytes,7,rep,name=amf_key4,json=amfKey4,proto3" json:"amf_key4,omitempty"`
	// for smfInfoSum
	SmfKey1 []*NRFKeyStruct `protobuf:"bytes,8,rep,name=smf_key1,json=smfKey1,proto3" json:"smf_key1,omitempty"`
	SmfKey2 []*NRFKeyStruct `protobuf:"bytes,9,rep,name=smf_key2,json=smfKey2,proto3" json:"smf_key2,omitempty"`
	SmfKey3 []*NRFKeyStruct `protobuf:"bytes,10,rep,name=smf_key3,json=smfKey3,proto3" json:"smf_key3,omitempty"`
	// for udmInfoSum
	UdmKey1 []*NRFKeyStruct `protobuf:"bytes,11,rep,name=udm_key1,json=udmKey1,proto3" json:"udm_key1,omitempty"`
	UdmKey2 []*NRFKeyStruct `protobuf:"bytes,12,rep,name=udm_key2,json=udmKey2,proto3" json:"udm_key2,omitempty"`
	// for ausfInfoSum
	AusfKey1 []*NRFKeyStruct `protobuf:"bytes,13,rep,name=ausf_key1,json=ausfKey1,proto3" json:"ausf_key1,omitempty"`
	AusfKey2 []*NRFKeyStruct `protobuf:"bytes,14,rep,name=ausf_key2,json=ausfKey2,proto3" json:"ausf_key2,omitempty"`
	// for pcfInfoSum
	PcfKey1              []*NRFKeyStruct `protobuf:"bytes,15,rep,name=pcf_key1,json=pcfKey1,proto3" json:"pcf_key1,omitempty"`
	PcfKey2              []*NRFKeyStruct `protobuf:"bytes,16,rep,name=pcf_key2,json=pcfKey2,proto3" json:"pcf_key2,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *NRFProfileIndex) Reset()         { *m = NRFProfileIndex{} }
func (m *NRFProfileIndex) String() string { return proto.CompactTextString(m) }
func (*NRFProfileIndex) ProtoMessage()    {}
func (*NRFProfileIndex) Descriptor() ([]byte, []int) {
	return fileDescriptor_NRFProfileIndex_35d735e5a0675006, []int{0}
}
func (m *NRFProfileIndex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NRFProfileIndex.Unmarshal(m, b)
}
func (m *NRFProfileIndex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NRFProfileIndex.Marshal(b, m, deterministic)
}
func (dst *NRFProfileIndex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NRFProfileIndex.Merge(dst, src)
}
func (m *NRFProfileIndex) XXX_Size() int {
	return xxx_messageInfo_NRFProfileIndex.Size(m)
}
func (m *NRFProfileIndex) XXX_DiscardUnknown() {
	xxx_messageInfo_NRFProfileIndex.DiscardUnknown(m)
}

var xxx_messageInfo_NRFProfileIndex proto.InternalMessageInfo

func (m *NRFProfileIndex) GetKey1() uint64 {
	if m != nil {
		return m.Key1
	}
	return 0
}

func (m *NRFProfileIndex) GetKey2() uint64 {
	if m != nil {
		return m.Key2
	}
	return 0
}

func (m *NRFProfileIndex) GetKey3() uint64 {
	if m != nil {
		return m.Key3
	}
	return 0
}

func (m *NRFProfileIndex) GetAmfKey1() []*NRFKeyStruct {
	if m != nil {
		return m.AmfKey1
	}
	return nil
}

func (m *NRFProfileIndex) GetAmfKey2() []*NRFKeyStruct {
	if m != nil {
		return m.AmfKey2
	}
	return nil
}

func (m *NRFProfileIndex) GetAmfKey3() []*NRFKeyStruct {
	if m != nil {
		return m.AmfKey3
	}
	return nil
}

func (m *NRFProfileIndex) GetAmfKey4() []*NRFKeyStruct {
	if m != nil {
		return m.AmfKey4
	}
	return nil
}

func (m *NRFProfileIndex) GetSmfKey1() []*NRFKeyStruct {
	if m != nil {
		return m.SmfKey1
	}
	return nil
}

func (m *NRFProfileIndex) GetSmfKey2() []*NRFKeyStruct {
	if m != nil {
		return m.SmfKey2
	}
	return nil
}

func (m *NRFProfileIndex) GetSmfKey3() []*NRFKeyStruct {
	if m != nil {
		return m.SmfKey3
	}
	return nil
}

func (m *NRFProfileIndex) GetUdmKey1() []*NRFKeyStruct {
	if m != nil {
		return m.UdmKey1
	}
	return nil
}

func (m *NRFProfileIndex) GetUdmKey2() []*NRFKeyStruct {
	if m != nil {
		return m.UdmKey2
	}
	return nil
}

func (m *NRFProfileIndex) GetAusfKey1() []*NRFKeyStruct {
	if m != nil {
		return m.AusfKey1
	}
	return nil
}

func (m *NRFProfileIndex) GetAusfKey2() []*NRFKeyStruct {
	if m != nil {
		return m.AusfKey2
	}
	return nil
}

func (m *NRFProfileIndex) GetPcfKey1() []*NRFKeyStruct {
	if m != nil {
		return m.PcfKey1
	}
	return nil
}

func (m *NRFProfileIndex) GetPcfKey2() []*NRFKeyStruct {
	if m != nil {
		return m.PcfKey2
	}
	return nil
}

type NRFKeyStruct struct {
	SubKey1              string   `protobuf:"bytes,1,opt,name=sub_key1,json=subKey1,proto3" json:"sub_key1,omitempty"`
	SubKey2              string   `protobuf:"bytes,2,opt,name=sub_key2,json=subKey2,proto3" json:"sub_key2,omitempty"`
	SubKey3              string   `protobuf:"bytes,3,opt,name=sub_key3,json=subKey3,proto3" json:"sub_key3,omitempty"`
	SubKey4              string   `protobuf:"bytes,4,opt,name=sub_key4,json=subKey4,proto3" json:"sub_key4,omitempty"`
	SubKey5              string   `protobuf:"bytes,5,opt,name=sub_key5,json=subKey5,proto3" json:"sub_key5,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NRFKeyStruct) Reset()         { *m = NRFKeyStruct{} }
func (m *NRFKeyStruct) String() string { return proto.CompactTextString(m) }
func (*NRFKeyStruct) ProtoMessage()    {}
func (*NRFKeyStruct) Descriptor() ([]byte, []int) {
	return fileDescriptor_NRFProfileIndex_35d735e5a0675006, []int{1}
}
func (m *NRFKeyStruct) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NRFKeyStruct.Unmarshal(m, b)
}
func (m *NRFKeyStruct) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NRFKeyStruct.Marshal(b, m, deterministic)
}
func (dst *NRFKeyStruct) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NRFKeyStruct.Merge(dst, src)
}
func (m *NRFKeyStruct) XXX_Size() int {
	return xxx_messageInfo_NRFKeyStruct.Size(m)
}
func (m *NRFKeyStruct) XXX_DiscardUnknown() {
	xxx_messageInfo_NRFKeyStruct.DiscardUnknown(m)
}

var xxx_messageInfo_NRFKeyStruct proto.InternalMessageInfo

func (m *NRFKeyStruct) GetSubKey1() string {
	if m != nil {
		return m.SubKey1
	}
	return ""
}

func (m *NRFKeyStruct) GetSubKey2() string {
	if m != nil {
		return m.SubKey2
	}
	return ""
}

func (m *NRFKeyStruct) GetSubKey3() string {
	if m != nil {
		return m.SubKey3
	}
	return ""
}

func (m *NRFKeyStruct) GetSubKey4() string {
	if m != nil {
		return m.SubKey4
	}
	return ""
}

func (m *NRFKeyStruct) GetSubKey5() string {
	if m != nil {
		return m.SubKey5
	}
	return ""
}

func init() {
	proto.RegisterType((*NRFProfileIndex)(nil), "grpc.NRFProfileIndex")
	proto.RegisterType((*NRFKeyStruct)(nil), "grpc.NRFKeyStruct")
}

func init() {
	proto.RegisterFile("nfmessage/nrfprofile/NRFProfileIndex.proto", fileDescriptor_NRFProfileIndex_35d735e5a0675006)
}

var fileDescriptor_NRFProfileIndex_35d735e5a0675006 = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xc1, 0x8a, 0xea, 0x30,
	0x14, 0x86, 0xf1, 0x5a, 0xb5, 0x46, 0xef, 0xf5, 0x52, 0x66, 0x91, 0xd9, 0x89, 0x2b, 0x19, 0x98,
	0x96, 0x49, 0xeb, 0x0b, 0xb8, 0x10, 0x06, 0x41, 0xa4, 0xb3, 0x9b, 0x8d, 0xb4, 0x69, 0x2a, 0x32,
	0x93, 0xa6, 0x24, 0x16, 0xf4, 0x41, 0xe6, 0x4d, 0xe6, 0x01, 0x87, 0xa6, 0x31, 0x0d, 0x32, 0x25,
	0xbb, 0x72, 0xfe, 0xef, 0x9c, 0xfc, 0xa7, 0x9c, 0x1f, 0x3c, 0x15, 0x39, 0x25, 0x42, 0x24, 0x47,
	0x12, 0x14, 0x3c, 0x2f, 0x39, 0xcb, 0x4f, 0x9f, 0x24, 0xd8, 0xc5, 0x9b, 0x7d, 0xf3, 0xf9, 0x5a,
	0x64, 0xe4, 0xe2, 0x97, 0x9c, 0x9d, 0x99, 0xe7, 0x1c, 0x79, 0x89, 0x17, 0xdf, 0x03, 0x30, 0xbb,
	0xd3, 0x3d, 0x0f, 0x38, 0x1f, 0xe4, 0xfa, 0x02, 0x7b, 0xf3, 0xde, 0xd2, 0x89, 0xe5, 0xb7, 0xaa,
	0x21, 0xf8, 0x47, 0xd7, 0x90, 0xaa, 0x85, 0xb0, 0xaf, 0x6b, 0xa1, 0xf7, 0x0c, 0xdc, 0x84, 0xe6,
	0x07, 0xd9, 0xef, 0xcc, 0xfb, 0xcb, 0x09, 0xf2, 0xfc, 0xfa, 0x21, 0x7f, 0x17, 0x6f, 0xb6, 0xe4,
	0xfa, 0x76, 0xe6, 0x15, 0x3e, 0xc7, 0xa3, 0x84, 0xe6, 0xdb, 0x7a, 0x6c, 0x8b, 0x23, 0x38, 0xb0,
	0xe1, 0xc8, 0xc0, 0x43, 0x38, 0xb4, 0xe1, 0xa6, 0x99, 0x08, 0x8e, 0x6c, 0x78, 0x54, 0xe3, 0xe2,
	0xe6, 0xdd, 0xed, 0xc6, 0x45, 0xeb, 0x5d, 0xdc, 0xbc, 0x8f, 0x6d, 0x38, 0x32, 0xf0, 0x10, 0x02,
	0x1b, 0x2e, 0xbd, 0x57, 0x19, 0x6d, 0xcc, 0x4c, 0xba, 0xf1, 0x2a, 0xa3, 0x37, 0x33, 0x0a, 0x47,
	0x70, 0x6a, 0xc3, 0x91, 0x17, 0x80, 0x71, 0x52, 0x09, 0xb5, 0xeb, 0xdf, 0x4e, 0xde, 0xad, 0x21,
	0x39, 0xdf, 0x68, 0x40, 0xf0, 0x9f, 0xb5, 0x41, 0xae, 0x5b, 0x62, 0xf5, 0xc0, 0xac, 0xdb, 0x50,
	0x89, 0xf5, 0xcf, 0x54, 0x38, 0x82, 0xff, 0x6d, 0x38, 0x5a, 0x7c, 0xf5, 0xc0, 0xd4, 0x54, 0xbc,
	0x47, 0xe0, 0x8a, 0x2a, 0x3d, 0xe8, 0xbb, 0x1d, 0xc7, 0x23, 0x51, 0xa5, 0x72, 0x74, 0x2b, 0x35,
	0xe7, 0xab, 0x25, 0x64, 0x48, 0xcd, 0x15, 0x6b, 0x29, 0x34, 0xa4, 0x08, 0x3a, 0xa6, 0x14, 0x19,
	0xd2, 0x0a, 0x0e, 0x4c, 0x69, 0xb5, 0xa6, 0x20, 0x24, 0xfc, 0x84, 0x85, 0x60, 0x85, 0x8f, 0x19,
	0x27, 0x7e, 0xc1, 0x73, 0x3f, 0x4b, 0x4b, 0xce, 0x2e, 0xd7, 0x66, 0x1f, 0x9d, 0x50, 0xbf, 0x4d,
	0xe8, 0xfa, 0xe1, 0x2e, 0x82, 0xfb, 0x3a, 0xa1, 0xef, 0x73, 0xcc, 0x68, 0xa0, 0x9a, 0x83, 0xdf,
	0x92, 0x9d, 0x0e, 0x65, 0x94, 0xc3, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xdd, 0x5d, 0x1d, 0x94,
	0xf8, 0x03, 0x00, 0x00,
}

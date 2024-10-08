// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/gpsiprofile/GpsiProfileGetRequest.proto

package gpsiprofile // import "com/dbproxy/nfmessage/gpsiprofile"

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

type GpsiProfileGetRequest struct {
	// Types that are valid to be assigned to Data:
	//	*GpsiProfileGetRequest_GpsiProfileId
	//	*GpsiProfileGetRequest_Filter
	//	*GpsiProfileGetRequest_FragmentSessionId
	Data                 isGpsiProfileGetRequest_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *GpsiProfileGetRequest) Reset()         { *m = GpsiProfileGetRequest{} }
func (m *GpsiProfileGetRequest) String() string { return proto.CompactTextString(m) }
func (*GpsiProfileGetRequest) ProtoMessage()    {}
func (*GpsiProfileGetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_GpsiProfileGetRequest_862252e8687393ca, []int{0}
}
func (m *GpsiProfileGetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GpsiProfileGetRequest.Unmarshal(m, b)
}
func (m *GpsiProfileGetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GpsiProfileGetRequest.Marshal(b, m, deterministic)
}
func (dst *GpsiProfileGetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GpsiProfileGetRequest.Merge(dst, src)
}
func (m *GpsiProfileGetRequest) XXX_Size() int {
	return xxx_messageInfo_GpsiProfileGetRequest.Size(m)
}
func (m *GpsiProfileGetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GpsiProfileGetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GpsiProfileGetRequest proto.InternalMessageInfo

type isGpsiProfileGetRequest_Data interface {
	isGpsiProfileGetRequest_Data()
}

type GpsiProfileGetRequest_GpsiProfileId struct {
	GpsiProfileId string `protobuf:"bytes,1,opt,name=gpsi_profile_id,json=gpsiProfileId,proto3,oneof"`
}

type GpsiProfileGetRequest_Filter struct {
	Filter *GpsiProfileFilter `protobuf:"bytes,2,opt,name=filter,proto3,oneof"`
}

type GpsiProfileGetRequest_FragmentSessionId struct {
	FragmentSessionId string `protobuf:"bytes,3,opt,name=fragment_session_id,json=fragmentSessionId,proto3,oneof"`
}

func (*GpsiProfileGetRequest_GpsiProfileId) isGpsiProfileGetRequest_Data() {}

func (*GpsiProfileGetRequest_Filter) isGpsiProfileGetRequest_Data() {}

func (*GpsiProfileGetRequest_FragmentSessionId) isGpsiProfileGetRequest_Data() {}

func (m *GpsiProfileGetRequest) GetData() isGpsiProfileGetRequest_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *GpsiProfileGetRequest) GetGpsiProfileId() string {
	if x, ok := m.GetData().(*GpsiProfileGetRequest_GpsiProfileId); ok {
		return x.GpsiProfileId
	}
	return ""
}

func (m *GpsiProfileGetRequest) GetFilter() *GpsiProfileFilter {
	if x, ok := m.GetData().(*GpsiProfileGetRequest_Filter); ok {
		return x.Filter
	}
	return nil
}

func (m *GpsiProfileGetRequest) GetFragmentSessionId() string {
	if x, ok := m.GetData().(*GpsiProfileGetRequest_FragmentSessionId); ok {
		return x.FragmentSessionId
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*GpsiProfileGetRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _GpsiProfileGetRequest_OneofMarshaler, _GpsiProfileGetRequest_OneofUnmarshaler, _GpsiProfileGetRequest_OneofSizer, []interface{}{
		(*GpsiProfileGetRequest_GpsiProfileId)(nil),
		(*GpsiProfileGetRequest_Filter)(nil),
		(*GpsiProfileGetRequest_FragmentSessionId)(nil),
	}
}

func _GpsiProfileGetRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*GpsiProfileGetRequest)
	// data
	switch x := m.Data.(type) {
	case *GpsiProfileGetRequest_GpsiProfileId:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.GpsiProfileId)
	case *GpsiProfileGetRequest_Filter:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Filter); err != nil {
			return err
		}
	case *GpsiProfileGetRequest_FragmentSessionId:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.FragmentSessionId)
	case nil:
	default:
		return fmt.Errorf("GpsiProfileGetRequest.Data has unexpected type %T", x)
	}
	return nil
}

func _GpsiProfileGetRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*GpsiProfileGetRequest)
	switch tag {
	case 1: // data.gpsi_profile_id
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Data = &GpsiProfileGetRequest_GpsiProfileId{x}
		return true, err
	case 2: // data.filter
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GpsiProfileFilter)
		err := b.DecodeMessage(msg)
		m.Data = &GpsiProfileGetRequest_Filter{msg}
		return true, err
	case 3: // data.fragment_session_id
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Data = &GpsiProfileGetRequest_FragmentSessionId{x}
		return true, err
	default:
		return false, nil
	}
}

func _GpsiProfileGetRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*GpsiProfileGetRequest)
	// data
	switch x := m.Data.(type) {
	case *GpsiProfileGetRequest_GpsiProfileId:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.GpsiProfileId)))
		n += len(x.GpsiProfileId)
	case *GpsiProfileGetRequest_Filter:
		s := proto.Size(x.Filter)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GpsiProfileGetRequest_FragmentSessionId:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.FragmentSessionId)))
		n += len(x.FragmentSessionId)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*GpsiProfileGetRequest)(nil), "grpc.GpsiProfileGetRequest")
}

func init() {
	proto.RegisterFile("nfmessage/gpsiprofile/GpsiProfileGetRequest.proto", fileDescriptor_GpsiProfileGetRequest_862252e8687393ca)
}

var fileDescriptor_GpsiProfileGetRequest_862252e8687393ca = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x4f, 0x4b, 0xc3, 0x30,
	0x18, 0xc6, 0x57, 0x1d, 0x05, 0x23, 0x22, 0x56, 0xc4, 0xb1, 0xd3, 0xf4, 0xd4, 0x8b, 0xa9, 0x53,
	0x3f, 0x41, 0x0f, 0x6e, 0xbb, 0x8d, 0x7a, 0xf3, 0x52, 0xba, 0xe6, 0x4d, 0x08, 0xac, 0x79, 0xe3,
	0xfb, 0x66, 0xa0, 0x5f, 0xca, 0xcf, 0x28, 0x5d, 0xeb, 0x1f, 0xb0, 0xb0, 0x5b, 0xe0, 0xc9, 0xf3,
	0x7b, 0x92, 0x9f, 0x98, 0x3b, 0xdd, 0x00, 0x73, 0x65, 0x20, 0x33, 0x9e, 0xad, 0x27, 0xd4, 0x76,
	0x0b, 0xd9, 0xc2, 0xb3, 0x5d, 0x77, 0xe7, 0x05, 0x84, 0x02, 0xde, 0x76, 0xc0, 0x41, 0x7a, 0xc2,
	0x80, 0xc9, 0xd8, 0x90, 0xaf, 0xa7, 0x77, 0x07, 0x8b, 0xcf, 0x76, 0x1b, 0x80, 0xba, 0xd2, 0xed,
	0x67, 0x24, 0xae, 0x06, 0xa1, 0x49, 0x2a, 0xce, 0x5b, 0x40, 0xd9, 0x13, 0x4a, 0xab, 0x26, 0xd1,
	0x2c, 0x4a, 0x4f, 0x96, 0xa3, 0xe2, 0xcc, 0xfc, 0x36, 0x56, 0x2a, 0x99, 0x8b, 0x58, 0xef, 0x99,
	0x93, 0xa3, 0x59, 0x94, 0x9e, 0x3e, 0x5c, 0xcb, 0xf6, 0x25, 0xf2, 0xdf, 0xe4, 0x72, 0x54, 0xf4,
	0x17, 0x93, 0x7b, 0x71, 0xa9, 0xa9, 0x32, 0x0d, 0xb8, 0x50, 0x32, 0x30, 0x5b, 0x74, 0xed, 0xc0,
	0x71, 0x3f, 0x70, 0xf1, 0x1d, 0xbe, 0x74, 0xd9, 0x4a, 0xe5, 0xb1, 0x18, 0xab, 0x2a, 0x54, 0xf9,
	0x4e, 0x3c, 0x01, 0xd9, 0x9a, 0x19, 0x9d, 0xac, 0x91, 0x40, 0x3a, 0xd2, 0x52, 0x6d, 0x3c, 0xe1,
	0xfb, 0x47, 0xb7, 0xfb, 0x23, 0x40, 0xfe, 0x11, 0x90, 0x4f, 0x07, 0x7f, 0xb9, 0x6e, 0x25, 0xbc,
	0xde, 0xd4, 0xd8, 0x64, 0x3d, 0x23, 0x1b, 0xf4, 0xb7, 0x89, 0xf7, 0xba, 0x1e, 0xbf, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x0c, 0x09, 0xc6, 0x8b, 0x98, 0x01, 0x00, 0x00,
}

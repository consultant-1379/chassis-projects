// Code generated by protoc-gen-go. DO NOT EDIT.
// source: nfmessage/nrfaddress/NRFAddressGetRequest.proto

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

type NRFAddressGetRequest struct {
	// Types that are valid to be assigned to Data:
	//	*NRFAddressGetRequest_NrfAddressId
	//	*NRFAddressGetRequest_Filter
	Data                 isNRFAddressGetRequest_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *NRFAddressGetRequest) Reset()         { *m = NRFAddressGetRequest{} }
func (m *NRFAddressGetRequest) String() string { return proto.CompactTextString(m) }
func (*NRFAddressGetRequest) ProtoMessage()    {}
func (*NRFAddressGetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_NRFAddressGetRequest_b15d08c4fcd6d1c6, []int{0}
}
func (m *NRFAddressGetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NRFAddressGetRequest.Unmarshal(m, b)
}
func (m *NRFAddressGetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NRFAddressGetRequest.Marshal(b, m, deterministic)
}
func (dst *NRFAddressGetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NRFAddressGetRequest.Merge(dst, src)
}
func (m *NRFAddressGetRequest) XXX_Size() int {
	return xxx_messageInfo_NRFAddressGetRequest.Size(m)
}
func (m *NRFAddressGetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NRFAddressGetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NRFAddressGetRequest proto.InternalMessageInfo

type isNRFAddressGetRequest_Data interface {
	isNRFAddressGetRequest_Data()
}

type NRFAddressGetRequest_NrfAddressId struct {
	NrfAddressId string `protobuf:"bytes,1,opt,name=nrf_address_id,json=nrfAddressId,proto3,oneof"`
}

type NRFAddressGetRequest_Filter struct {
	Filter *NRFAddressFilter `protobuf:"bytes,2,opt,name=filter,proto3,oneof"`
}

func (*NRFAddressGetRequest_NrfAddressId) isNRFAddressGetRequest_Data() {}

func (*NRFAddressGetRequest_Filter) isNRFAddressGetRequest_Data() {}

func (m *NRFAddressGetRequest) GetData() isNRFAddressGetRequest_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *NRFAddressGetRequest) GetNrfAddressId() string {
	if x, ok := m.GetData().(*NRFAddressGetRequest_NrfAddressId); ok {
		return x.NrfAddressId
	}
	return ""
}

func (m *NRFAddressGetRequest) GetFilter() *NRFAddressFilter {
	if x, ok := m.GetData().(*NRFAddressGetRequest_Filter); ok {
		return x.Filter
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*NRFAddressGetRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _NRFAddressGetRequest_OneofMarshaler, _NRFAddressGetRequest_OneofUnmarshaler, _NRFAddressGetRequest_OneofSizer, []interface{}{
		(*NRFAddressGetRequest_NrfAddressId)(nil),
		(*NRFAddressGetRequest_Filter)(nil),
	}
}

func _NRFAddressGetRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*NRFAddressGetRequest)
	// data
	switch x := m.Data.(type) {
	case *NRFAddressGetRequest_NrfAddressId:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.NrfAddressId)
	case *NRFAddressGetRequest_Filter:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Filter); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("NRFAddressGetRequest.Data has unexpected type %T", x)
	}
	return nil
}

func _NRFAddressGetRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*NRFAddressGetRequest)
	switch tag {
	case 1: // data.nrf_address_id
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Data = &NRFAddressGetRequest_NrfAddressId{x}
		return true, err
	case 2: // data.filter
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(NRFAddressFilter)
		err := b.DecodeMessage(msg)
		m.Data = &NRFAddressGetRequest_Filter{msg}
		return true, err
	default:
		return false, nil
	}
}

func _NRFAddressGetRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*NRFAddressGetRequest)
	// data
	switch x := m.Data.(type) {
	case *NRFAddressGetRequest_NrfAddressId:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.NrfAddressId)))
		n += len(x.NrfAddressId)
	case *NRFAddressGetRequest_Filter:
		s := proto.Size(x.Filter)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*NRFAddressGetRequest)(nil), "grpc.NRFAddressGetRequest")
}

func init() {
	proto.RegisterFile("nfmessage/nrfaddress/NRFAddressGetRequest.proto", fileDescriptor_NRFAddressGetRequest_b15d08c4fcd6d1c6)
}

var fileDescriptor_NRFAddressGetRequest_b15d08c4fcd6d1c6 = []byte{
	// 218 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xcf, 0x4b, 0xcb, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0xcf, 0x2b, 0x4a, 0x4b, 0x4c, 0x49, 0x29, 0x4a, 0x2d, 0x2e,
	0xd6, 0xf7, 0x0b, 0x72, 0x73, 0x84, 0x30, 0xdd, 0x53, 0x4b, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b,
	0x4b, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x58, 0xd2, 0x8b, 0x0a, 0x92, 0xa5, 0xb4, 0x09,
	0x68, 0x73, 0xcb, 0xcc, 0x29, 0x49, 0x2d, 0x82, 0x68, 0x51, 0xaa, 0xe0, 0x12, 0xc1, 0x66, 0xa0,
	0x90, 0x1a, 0x17, 0x5f, 0x5e, 0x51, 0x5a, 0x3c, 0x54, 0x77, 0x7c, 0x66, 0x8a, 0x04, 0xa3, 0x02,
	0xa3, 0x06, 0xa7, 0x07, 0x43, 0x10, 0x4f, 0x5e, 0x51, 0x1a, 0x54, 0xbd, 0x67, 0x8a, 0x90, 0x01,
	0x17, 0x5b, 0x1a, 0xd8, 0x3c, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x6e, 0x23, 0x31, 0x3d, 0x90, 0x1b,
	0xf4, 0xd0, 0x6d, 0xf3, 0x60, 0x08, 0x82, 0xaa, 0x73, 0x62, 0xe3, 0x62, 0x49, 0x49, 0x2c, 0x49,
	0x74, 0x2a, 0xe2, 0x32, 0x4e, 0x2d, 0xca, 0x4c, 0x2e, 0x2e, 0xce, 0xcf, 0xd3, 0x4b, 0xce, 0x2f,
	0x4a, 0xd5, 0xcb, 0x2b, 0x4a, 0xd3, 0x4b, 0x49, 0x2a, 0x28, 0xca, 0xaf, 0xa8, 0x84, 0x18, 0x02,
	0xf7, 0x87, 0x1e, 0xc2, 0x1f, 0x4e, 0x92, 0xd8, 0x9c, 0x1b, 0x00, 0xf2, 0x4b, 0x94, 0x42, 0x72,
	0x7e, 0xae, 0x3e, 0xd4, 0x04, 0xac, 0x61, 0x97, 0xc4, 0x06, 0xf6, 0xb4, 0x31, 0x20, 0x00, 0x00,
	0xff, 0xff, 0x6a, 0x18, 0xc2, 0x1d, 0x5a, 0x01, 0x00, 0x00,
}

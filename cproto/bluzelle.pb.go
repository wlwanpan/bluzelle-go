// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bluzelle.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type BznMsgType int32

const (
	BznMsgType_BZN_MSG_UNDEFINED BznMsgType = 0
	BznMsgType_BZN_MSG_PBFT      BznMsgType = 1
)

var BznMsgType_name = map[int32]string{
	0: "BZN_MSG_UNDEFINED",
	1: "BZN_MSG_PBFT",
}

var BznMsgType_value = map[string]int32{
	"BZN_MSG_UNDEFINED": 0,
	"BZN_MSG_PBFT":      1,
}

func (x BznMsgType) String() string {
	return proto.EnumName(BznMsgType_name, int32(x))
}

func (BznMsgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_37e4e006cca7cf40, []int{0}
}

type BznMsg struct {
	// Types that are valid to be assigned to Msg:
	//	*BznMsg_Db
	//	*BznMsg_Json
	//	*BznMsg_AuditMessage
	//	*BznMsg_Pbft
	Msg                  isBznMsg_Msg `protobuf_oneof:"msg"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *BznMsg) Reset()         { *m = BznMsg{} }
func (m *BznMsg) String() string { return proto.CompactTextString(m) }
func (*BznMsg) ProtoMessage()    {}
func (*BznMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_37e4e006cca7cf40, []int{0}
}

func (m *BznMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BznMsg.Unmarshal(m, b)
}
func (m *BznMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BznMsg.Marshal(b, m, deterministic)
}
func (m *BznMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BznMsg.Merge(m, src)
}
func (m *BznMsg) XXX_Size() int {
	return xxx_messageInfo_BznMsg.Size(m)
}
func (m *BznMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_BznMsg.DiscardUnknown(m)
}

var xxx_messageInfo_BznMsg proto.InternalMessageInfo

type isBznMsg_Msg interface {
	isBznMsg_Msg()
}

type BznMsg_Db struct {
	Db *DatabaseMsg `protobuf:"bytes,10,opt,name=db,proto3,oneof"`
}

type BznMsg_Json struct {
	Json string `protobuf:"bytes,11,opt,name=json,proto3,oneof"`
}

type BznMsg_AuditMessage struct {
	AuditMessage *AuditMessage `protobuf:"bytes,12,opt,name=audit_message,json=auditMessage,proto3,oneof"`
}

type BznMsg_Pbft struct {
	Pbft *PbftMsg `protobuf:"bytes,13,opt,name=pbft,proto3,oneof"`
}

func (*BznMsg_Db) isBznMsg_Msg() {}

func (*BznMsg_Json) isBznMsg_Msg() {}

func (*BznMsg_AuditMessage) isBznMsg_Msg() {}

func (*BznMsg_Pbft) isBznMsg_Msg() {}

func (m *BznMsg) GetMsg() isBznMsg_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *BznMsg) GetDb() *DatabaseMsg {
	if x, ok := m.GetMsg().(*BznMsg_Db); ok {
		return x.Db
	}
	return nil
}

func (m *BznMsg) GetJson() string {
	if x, ok := m.GetMsg().(*BznMsg_Json); ok {
		return x.Json
	}
	return ""
}

func (m *BznMsg) GetAuditMessage() *AuditMessage {
	if x, ok := m.GetMsg().(*BznMsg_AuditMessage); ok {
		return x.AuditMessage
	}
	return nil
}

func (m *BznMsg) GetPbft() *PbftMsg {
	if x, ok := m.GetMsg().(*BznMsg_Pbft); ok {
		return x.Pbft
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*BznMsg) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _BznMsg_OneofMarshaler, _BznMsg_OneofUnmarshaler, _BznMsg_OneofSizer, []interface{}{
		(*BznMsg_Db)(nil),
		(*BznMsg_Json)(nil),
		(*BznMsg_AuditMessage)(nil),
		(*BznMsg_Pbft)(nil),
	}
}

func _BznMsg_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*BznMsg)
	// msg
	switch x := m.Msg.(type) {
	case *BznMsg_Db:
		b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Db); err != nil {
			return err
		}
	case *BznMsg_Json:
		b.EncodeVarint(11<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Json)
	case *BznMsg_AuditMessage:
		b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.AuditMessage); err != nil {
			return err
		}
	case *BznMsg_Pbft:
		b.EncodeVarint(13<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Pbft); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("BznMsg.Msg has unexpected type %T", x)
	}
	return nil
}

func _BznMsg_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*BznMsg)
	switch tag {
	case 10: // msg.db
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DatabaseMsg)
		err := b.DecodeMessage(msg)
		m.Msg = &BznMsg_Db{msg}
		return true, err
	case 11: // msg.json
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Msg = &BznMsg_Json{x}
		return true, err
	case 12: // msg.audit_message
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(AuditMessage)
		err := b.DecodeMessage(msg)
		m.Msg = &BznMsg_AuditMessage{msg}
		return true, err
	case 13: // msg.pbft
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PbftMsg)
		err := b.DecodeMessage(msg)
		m.Msg = &BznMsg_Pbft{msg}
		return true, err
	default:
		return false, nil
	}
}

func _BznMsg_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*BznMsg)
	// msg
	switch x := m.Msg.(type) {
	case *BznMsg_Db:
		s := proto.Size(x.Db)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *BznMsg_Json:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Json)))
		n += len(x.Json)
	case *BznMsg_AuditMessage:
		s := proto.Size(x.AuditMessage)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *BznMsg_Pbft:
		s := proto.Size(x.Pbft)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type WrappedBznMsg struct {
	// This is stored as a serialized string because we need to sign it, and serialization is not guarenteed to be deterministic
	Payload              []byte     `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Type                 BznMsgType `protobuf:"varint,2,opt,name=type,proto3,enum=pb.BznMsgType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *WrappedBznMsg) Reset()         { *m = WrappedBznMsg{} }
func (m *WrappedBznMsg) String() string { return proto.CompactTextString(m) }
func (*WrappedBznMsg) ProtoMessage()    {}
func (*WrappedBznMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_37e4e006cca7cf40, []int{1}
}

func (m *WrappedBznMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WrappedBznMsg.Unmarshal(m, b)
}
func (m *WrappedBznMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WrappedBznMsg.Marshal(b, m, deterministic)
}
func (m *WrappedBznMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WrappedBznMsg.Merge(m, src)
}
func (m *WrappedBznMsg) XXX_Size() int {
	return xxx_messageInfo_WrappedBznMsg.Size(m)
}
func (m *WrappedBznMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_WrappedBznMsg.DiscardUnknown(m)
}

var xxx_messageInfo_WrappedBznMsg proto.InternalMessageInfo

func (m *WrappedBznMsg) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *WrappedBznMsg) GetType() BznMsgType {
	if m != nil {
		return m.Type
	}
	return BznMsgType_BZN_MSG_UNDEFINED
}

func init() {
	proto.RegisterEnum("pb.BznMsgType", BznMsgType_name, BznMsgType_value)
	proto.RegisterType((*BznMsg)(nil), "pb.bzn_msg")
	proto.RegisterType((*WrappedBznMsg)(nil), "pb.wrapped_bzn_msg")
}

func init() { proto.RegisterFile("bluzelle.proto", fileDescriptor_37e4e006cca7cf40) }

var fileDescriptor_37e4e006cca7cf40 = []byte{
	// 282 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0xb1, 0x4e, 0xf3, 0x30,
	0x14, 0x85, 0xe3, 0xfc, 0xfd, 0xa9, 0xb8, 0x75, 0x43, 0x6a, 0x81, 0x64, 0x75, 0x8a, 0x22, 0x86,
	0x88, 0x21, 0x03, 0x0c, 0x30, 0x47, 0x6d, 0x09, 0x43, 0x23, 0x08, 0xb0, 0xb0, 0x58, 0xb6, 0x6c,
	0x22, 0x50, 0xd2, 0x58, 0x75, 0x2a, 0xd4, 0x3e, 0x11, 0x8f, 0x89, 0xec, 0x26, 0x12, 0x6c, 0x3e,
	0xdf, 0xd1, 0x39, 0xf7, 0xfa, 0x42, 0x20, 0xea, 0xdd, 0x41, 0xd5, 0xb5, 0x4a, 0xf5, 0xb6, 0xed,
	0x5a, 0xe2, 0x6b, 0x31, 0x0f, 0x24, 0xef, 0xb8, 0xe0, 0xa6, 0x67, 0xf3, 0x09, 0xdf, 0xc9, 0x8f,
	0xae, 0x17, 0xa0, 0xc5, 0x7b, 0xff, 0x8e, 0xbf, 0x11, 0x8c, 0xc5, 0x61, 0xc3, 0x1a, 0x53, 0x91,
	0x18, 0x7c, 0x29, 0x28, 0x44, 0x28, 0x99, 0x5c, 0x87, 0xa9, 0x16, 0xe9, 0x50, 0x62, 0xdd, 0xdc,
	0x2b, 0x7d, 0x29, 0xc8, 0x39, 0x8c, 0x3e, 0x4d, 0xbb, 0xa1, 0x93, 0x08, 0x25, 0xa7, 0xb9, 0x57,
	0x3a, 0x45, 0xee, 0x60, 0xea, 0x06, 0xb0, 0x46, 0x19, 0xc3, 0x2b, 0x45, 0xb1, 0x2b, 0x99, 0xd9,
	0x92, 0x3f, 0x46, 0xee, 0x95, 0xd8, 0x81, 0xf5, 0x51, 0x93, 0x18, 0x46, 0x76, 0x1b, 0x3a, 0x75,
	0x01, 0x6c, 0x03, 0x56, 0xf7, 0x13, 0x9d, 0x97, 0xfd, 0x87, 0x7f, 0x8d, 0xa9, 0xe2, 0x27, 0x38,
	0xfb, 0xda, 0x72, 0xad, 0x95, 0x64, 0xc3, 0xc6, 0x14, 0xc6, 0x9a, 0xef, 0xeb, 0x96, 0x4b, 0x8a,
	0x22, 0x94, 0xe0, 0x72, 0x90, 0xe4, 0x12, 0x46, 0xdd, 0x5e, 0x2b, 0xea, 0x47, 0x28, 0x09, 0x8e,
	0xbf, 0xe9, 0x43, 0xcc, 0xf2, 0xd2, 0xb9, 0x57, 0xb7, 0x80, 0x7f, 0x53, 0x72, 0x01, 0xb3, 0xec,
	0xad, 0x60, 0xeb, 0xe7, 0x7b, 0xf6, 0x5a, 0x2c, 0x96, 0xab, 0x87, 0x62, 0xb9, 0x08, 0x3d, 0x12,
	0x02, 0x1e, 0xf0, 0x63, 0xb6, 0x7a, 0x09, 0x91, 0x38, 0x71, 0xd7, 0xbb, 0xf9, 0x09, 0x00, 0x00,
	0xff, 0xff, 0xd3, 0x76, 0xc2, 0x18, 0x7c, 0x01, 0x00, 0x00,
}
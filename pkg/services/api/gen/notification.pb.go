// Code generated by protoc-gen-go. DO NOT EDIT.
// source: notification.proto

package gen

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Notification struct {
	Summary              string               `protobuf:"bytes,1,opt,name=Summary,proto3" json:"Summary,omitempty"`
	Date                 *timestamp.Timestamp `protobuf:"bytes,2,opt,name=Date,proto3" json:"Date,omitempty"`
	User                 string               `protobuf:"bytes,3,opt,name=User,proto3" json:"User,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Notification) Reset()         { *m = Notification{} }
func (m *Notification) String() string { return proto.CompactTextString(m) }
func (*Notification) ProtoMessage()    {}
func (*Notification) Descriptor() ([]byte, []int) {
	return fileDescriptor_736a457d4a5efa07, []int{0}
}

func (m *Notification) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Notification.Unmarshal(m, b)
}
func (m *Notification) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Notification.Marshal(b, m, deterministic)
}
func (m *Notification) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Notification.Merge(m, src)
}
func (m *Notification) XXX_Size() int {
	return xxx_messageInfo_Notification.Size(m)
}
func (m *Notification) XXX_DiscardUnknown() {
	xxx_messageInfo_Notification.DiscardUnknown(m)
}

var xxx_messageInfo_Notification proto.InternalMessageInfo

func (m *Notification) GetSummary() string {
	if m != nil {
		return m.Summary
	}
	return ""
}

func (m *Notification) GetDate() *timestamp.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

func (m *Notification) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func init() {
	proto.RegisterType((*Notification)(nil), "proto.Notification")
}

func init() { proto.RegisterFile("notification.proto", fileDescriptor_736a457d4a5efa07) }

var fileDescriptor_736a457d4a5efa07 = []byte{
	// 149 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xca, 0xcb, 0x2f, 0xc9,
	0x4c, 0xcb, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x05, 0x53, 0x52, 0xf2, 0xe9, 0xf9, 0xf9, 0xe9, 0x39, 0xa9, 0xfa, 0x60, 0x5e, 0x52, 0x69, 0x9a,
	0x7e, 0x49, 0x66, 0x6e, 0x6a, 0x71, 0x49, 0x62, 0x6e, 0x01, 0x44, 0x9d, 0x52, 0x0e, 0x17, 0x8f,
	0x1f, 0x92, 0x6e, 0x21, 0x09, 0x2e, 0xf6, 0xe0, 0xd2, 0xdc, 0xdc, 0xc4, 0xa2, 0x4a, 0x09, 0x46,
	0x05, 0x46, 0x0d, 0xce, 0x20, 0x18, 0x57, 0x48, 0x8f, 0x8b, 0xc5, 0x25, 0xb1, 0x24, 0x55, 0x82,
	0x49, 0x81, 0x51, 0x83, 0xdb, 0x48, 0x4a, 0x0f, 0x62, 0xb2, 0x1e, 0xcc, 0x64, 0xbd, 0x10, 0x98,
	0xc9, 0x41, 0x60, 0x75, 0x42, 0x42, 0x5c, 0x2c, 0xa1, 0xc5, 0xa9, 0x45, 0x12, 0xcc, 0x60, 0x63,
	0xc0, 0xec, 0x24, 0x36, 0xb0, 0x6a, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x37, 0xac, 0x4e,
	0x75, 0xb2, 0x00, 0x00, 0x00,
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/gocongress.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	proto/gocongress.proto
	proto/user.proto

It has these top-level messages:
	Upstream
	ItemRange
	SendLawRequest
	SendLawReply
	SendVoteRequest
	SendVoteReply
	StatusRequest
	StatusReply
	UpdateRequest
	UpdateReply
	User
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Upstream struct {
	Name            string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Description     string `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	BodyName        string `protobuf:"bytes,3,opt,name=body_name,json=bodyName" json:"body_name,omitempty"`
	BodyDescription string `protobuf:"bytes,4,opt,name=body_description,json=bodyDescription" json:"body_description,omitempty"`
	BodyUrl         string `protobuf:"bytes,5,opt,name=body_url,json=bodyUrl" json:"body_url,omitempty"`
	BodyEmail       string `protobuf:"bytes,6,opt,name=body_email,json=bodyEmail" json:"body_email,omitempty"`
}

func (m *Upstream) Reset()                    { *m = Upstream{} }
func (m *Upstream) String() string            { return proto1.CompactTextString(m) }
func (*Upstream) ProtoMessage()               {}
func (*Upstream) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Upstream) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Upstream) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Upstream) GetBodyName() string {
	if m != nil {
		return m.BodyName
	}
	return ""
}

func (m *Upstream) GetBodyDescription() string {
	if m != nil {
		return m.BodyDescription
	}
	return ""
}

func (m *Upstream) GetBodyUrl() string {
	if m != nil {
		return m.BodyUrl
	}
	return ""
}

func (m *Upstream) GetBodyEmail() string {
	if m != nil {
		return m.BodyEmail
	}
	return ""
}

type ItemRange struct {
	Ident     string `protobuf:"bytes,1,opt,name=ident" json:"ident,omitempty"`
	Username  string `protobuf:"bytes,2,opt,name=username" json:"username,omitempty"`
	StartDate string `protobuf:"bytes,3,opt,name=start_date,json=startDate" json:"start_date,omitempty"`
	EndDate   string `protobuf:"bytes,4,opt,name=end_date,json=endDate" json:"end_date,omitempty"`
	Limit     int32  `protobuf:"varint,5,opt,name=limit" json:"limit,omitempty"`
}

func (m *ItemRange) Reset()                    { *m = ItemRange{} }
func (m *ItemRange) String() string            { return proto1.CompactTextString(m) }
func (*ItemRange) ProtoMessage()               {}
func (*ItemRange) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ItemRange) GetIdent() string {
	if m != nil {
		return m.Ident
	}
	return ""
}

func (m *ItemRange) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *ItemRange) GetStartDate() string {
	if m != nil {
		return m.StartDate
	}
	return ""
}

func (m *ItemRange) GetEndDate() string {
	if m != nil {
		return m.EndDate
	}
	return ""
}

func (m *ItemRange) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type SendLawRequest struct {
	ItemRange *ItemRange `protobuf:"bytes,1,opt,name=item_range,json=itemRange" json:"item_range,omitempty"`
}

func (m *SendLawRequest) Reset()                    { *m = SendLawRequest{} }
func (m *SendLawRequest) String() string            { return proto1.CompactTextString(m) }
func (*SendLawRequest) ProtoMessage()               {}
func (*SendLawRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SendLawRequest) GetItemRange() *ItemRange {
	if m != nil {
		return m.ItemRange
	}
	return nil
}

type SendLawReply struct {
	NewItems []string `protobuf:"bytes,1,rep,name=new_items,json=newItems" json:"new_items,omitempty"`
	Updated  []string `protobuf:"bytes,2,rep,name=updated" json:"updated,omitempty"`
	Err      string   `protobuf:"bytes,3,opt,name=err" json:"err,omitempty"`
}

func (m *SendLawReply) Reset()                    { *m = SendLawReply{} }
func (m *SendLawReply) String() string            { return proto1.CompactTextString(m) }
func (*SendLawReply) ProtoMessage()               {}
func (*SendLawReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SendLawReply) GetNewItems() []string {
	if m != nil {
		return m.NewItems
	}
	return nil
}

func (m *SendLawReply) GetUpdated() []string {
	if m != nil {
		return m.Updated
	}
	return nil
}

func (m *SendLawReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type SendVoteRequest struct {
	ItemRange *ItemRange `protobuf:"bytes,1,opt,name=item_range,json=itemRange" json:"item_range,omitempty"`
}

func (m *SendVoteRequest) Reset()                    { *m = SendVoteRequest{} }
func (m *SendVoteRequest) String() string            { return proto1.CompactTextString(m) }
func (*SendVoteRequest) ProtoMessage()               {}
func (*SendVoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SendVoteRequest) GetItemRange() *ItemRange {
	if m != nil {
		return m.ItemRange
	}
	return nil
}

type SendVoteReply struct {
	NewItems []string `protobuf:"bytes,1,rep,name=new_items,json=newItems" json:"new_items,omitempty"`
	Updated  []string `protobuf:"bytes,2,rep,name=updated" json:"updated,omitempty"`
	Err      string   `protobuf:"bytes,3,opt,name=err" json:"err,omitempty"`
}

func (m *SendVoteReply) Reset()                    { *m = SendVoteReply{} }
func (m *SendVoteReply) String() string            { return proto1.CompactTextString(m) }
func (*SendVoteReply) ProtoMessage()               {}
func (*SendVoteReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SendVoteReply) GetNewItems() []string {
	if m != nil {
		return m.NewItems
	}
	return nil
}

func (m *SendVoteReply) GetUpdated() []string {
	if m != nil {
		return m.Updated
	}
	return nil
}

func (m *SendVoteReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type StatusRequest struct {
	ItemRange *ItemRange `protobuf:"bytes,1,opt,name=item_range,json=itemRange" json:"item_range,omitempty"`
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto1.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *StatusRequest) GetItemRange() *ItemRange {
	if m != nil {
		return m.ItemRange
	}
	return nil
}

type StatusReply struct {
	NewItems []string `protobuf:"bytes,1,rep,name=new_items,json=newItems" json:"new_items,omitempty"`
	Existing []string `protobuf:"bytes,2,rep,name=existing" json:"existing,omitempty"`
	Err      string   `protobuf:"bytes,3,opt,name=err" json:"err,omitempty"`
}

func (m *StatusReply) Reset()                    { *m = StatusReply{} }
func (m *StatusReply) String() string            { return proto1.CompactTextString(m) }
func (*StatusReply) ProtoMessage()               {}
func (*StatusReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *StatusReply) GetNewItems() []string {
	if m != nil {
		return m.NewItems
	}
	return nil
}

func (m *StatusReply) GetExisting() []string {
	if m != nil {
		return m.Existing
	}
	return nil
}

func (m *StatusReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type UpdateRequest struct {
	Upstream *Upstream `protobuf:"bytes,1,opt,name=upstream" json:"upstream,omitempty"`
	Err      string    `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *UpdateRequest) Reset()                    { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string            { return proto1.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()               {}
func (*UpdateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *UpdateRequest) GetUpstream() *Upstream {
	if m != nil {
		return m.Upstream
	}
	return nil
}

func (m *UpdateRequest) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type UpdateReply struct {
	Err string `protobuf:"bytes,1,opt,name=err" json:"err,omitempty"`
}

func (m *UpdateReply) Reset()                    { *m = UpdateReply{} }
func (m *UpdateReply) String() string            { return proto1.CompactTextString(m) }
func (*UpdateReply) ProtoMessage()               {}
func (*UpdateReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *UpdateReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto1.RegisterType((*Upstream)(nil), "proto.Upstream")
	proto1.RegisterType((*ItemRange)(nil), "proto.ItemRange")
	proto1.RegisterType((*SendLawRequest)(nil), "proto.SendLawRequest")
	proto1.RegisterType((*SendLawReply)(nil), "proto.SendLawReply")
	proto1.RegisterType((*SendVoteRequest)(nil), "proto.SendVoteRequest")
	proto1.RegisterType((*SendVoteReply)(nil), "proto.SendVoteReply")
	proto1.RegisterType((*StatusRequest)(nil), "proto.StatusRequest")
	proto1.RegisterType((*StatusReply)(nil), "proto.StatusReply")
	proto1.RegisterType((*UpdateRequest)(nil), "proto.UpdateRequest")
	proto1.RegisterType((*UpdateReply)(nil), "proto.UpdateReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GoCongressService service

type GoCongressServiceClient interface {
	SendVote(ctx context.Context, in *SendVoteRequest, opts ...grpc.CallOption) (*SendVoteReply, error)
	SendLaw(ctx context.Context, in *SendLawRequest, opts ...grpc.CallOption) (*SendLawReply, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error)
}

type goCongressServiceClient struct {
	cc *grpc.ClientConn
}

func NewGoCongressServiceClient(cc *grpc.ClientConn) GoCongressServiceClient {
	return &goCongressServiceClient{cc}
}

func (c *goCongressServiceClient) SendVote(ctx context.Context, in *SendVoteRequest, opts ...grpc.CallOption) (*SendVoteReply, error) {
	out := new(SendVoteReply)
	err := grpc.Invoke(ctx, "/proto.GoCongressService/SendVote", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goCongressServiceClient) SendLaw(ctx context.Context, in *SendLawRequest, opts ...grpc.CallOption) (*SendLawReply, error) {
	out := new(SendLawReply)
	err := grpc.Invoke(ctx, "/proto.GoCongressService/SendLaw", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goCongressServiceClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error) {
	out := new(StatusReply)
	err := grpc.Invoke(ctx, "/proto.GoCongressService/Status", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GoCongressService service

type GoCongressServiceServer interface {
	SendVote(context.Context, *SendVoteRequest) (*SendVoteReply, error)
	SendLaw(context.Context, *SendLawRequest) (*SendLawReply, error)
	Status(context.Context, *StatusRequest) (*StatusReply, error)
}

func RegisterGoCongressServiceServer(s *grpc.Server, srv GoCongressServiceServer) {
	s.RegisterService(&_GoCongressService_serviceDesc, srv)
}

func _GoCongressService_SendVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoCongressServiceServer).SendVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GoCongressService/SendVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoCongressServiceServer).SendVote(ctx, req.(*SendVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoCongressService_SendLaw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendLawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoCongressServiceServer).SendLaw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GoCongressService/SendLaw",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoCongressServiceServer).SendLaw(ctx, req.(*SendLawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoCongressService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoCongressServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GoCongressService/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoCongressServiceServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GoCongressService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GoCongressService",
	HandlerType: (*GoCongressServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendVote",
			Handler:    _GoCongressService_SendVote_Handler,
		},
		{
			MethodName: "SendLaw",
			Handler:    _GoCongressService_SendLaw_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _GoCongressService_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/gocongress.proto",
}

func init() { proto1.RegisterFile("proto/gocongress.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 500 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x5d, 0x8b, 0xd3, 0x4c,
	0x14, 0xc7, 0x9f, 0xb4, 0x9b, 0x36, 0x39, 0x7d, 0x6a, 0xeb, 0x58, 0x97, 0x18, 0x11, 0x4b, 0xae,
	0x14, 0x61, 0x17, 0x56, 0x41, 0xf0, 0xca, 0x97, 0x15, 0x11, 0x64, 0x2f, 0x52, 0xaa, 0xde, 0x85,
	0x6c, 0x73, 0x28, 0x03, 0xc9, 0x24, 0xce, 0x4c, 0xac, 0xfd, 0x10, 0x7e, 0x25, 0x2f, 0xfc, 0x64,
	0x32, 0x27, 0x93, 0x34, 0x54, 0x41, 0x50, 0xaf, 0x32, 0xe7, 0x7f, 0xce, 0xf9, 0xe7, 0x77, 0xe6,
	0x05, 0x4e, 0x2b, 0x59, 0xea, 0xf2, 0x7c, 0x5b, 0x6e, 0x4a, 0xb1, 0x95, 0xa8, 0xd4, 0x19, 0x09,
	0xcc, 0xa5, 0x4f, 0xf4, 0xdd, 0x01, 0x6f, 0x5d, 0x29, 0x2d, 0x31, 0x2d, 0x18, 0x83, 0x13, 0x91,
	0x16, 0x18, 0x38, 0x4b, 0xe7, 0x81, 0x1f, 0xd3, 0x9a, 0x2d, 0x61, 0x92, 0xa1, 0xda, 0x48, 0x5e,
	0x69, 0x5e, 0x8a, 0x60, 0x40, 0xa9, 0xbe, 0xc4, 0xee, 0x82, 0x7f, 0x5d, 0x66, 0xfb, 0x84, 0x5a,
	0x87, 0x94, 0xf7, 0x8c, 0x70, 0x65, 0xda, 0x1f, 0xc2, 0x9c, 0x92, 0x7d, 0x8f, 0x13, 0xaa, 0x99,
	0x19, 0xfd, 0xb2, 0xe7, 0x73, 0x07, 0xa8, 0x2d, 0xa9, 0x65, 0x1e, 0xb8, 0x54, 0x32, 0x36, 0xf1,
	0x5a, 0xe6, 0xec, 0x1e, 0x00, 0xa5, 0xb0, 0x48, 0x79, 0x1e, 0x8c, 0x28, 0x49, 0x3f, 0x7d, 0x6d,
	0x84, 0xe8, 0xab, 0x03, 0xfe, 0x5b, 0x8d, 0x45, 0x9c, 0x8a, 0x2d, 0xb2, 0x05, 0xb8, 0x3c, 0x43,
	0xa1, 0xed, 0x18, 0x4d, 0xc0, 0x42, 0xf0, 0x6a, 0x85, 0x92, 0x20, 0x9b, 0x21, 0xba, 0xd8, 0xd8,
	0x2b, 0x9d, 0x4a, 0x9d, 0x64, 0xa9, 0x6e, 0x47, 0xf0, 0x49, 0xb9, 0x4c, 0x35, 0x1a, 0x30, 0x14,
	0x59, 0x93, 0x6c, 0xd8, 0xc7, 0x28, 0x32, 0x4a, 0x2d, 0xc0, 0xcd, 0x79, 0xc1, 0x35, 0x01, 0xbb,
	0x71, 0x13, 0x44, 0x2f, 0xe0, 0xc6, 0x0a, 0x45, 0xf6, 0x2e, 0xdd, 0xc5, 0xf8, 0xa9, 0x46, 0xa5,
	0xd9, 0x39, 0x00, 0xd7, 0x58, 0x24, 0xd2, 0x10, 0x12, 0xd8, 0xe4, 0x62, 0xde, 0x9c, 0xc4, 0x59,
	0x47, 0x1e, 0xfb, 0xbc, 0x5d, 0x46, 0x1f, 0xe0, 0xff, 0xce, 0xa2, 0xca, 0xf7, 0x66, 0x93, 0x05,
	0xee, 0x12, 0x53, 0xa0, 0x02, 0x67, 0x39, 0x34, 0xfc, 0x02, 0x77, 0xa6, 0x57, 0xb1, 0x00, 0xc6,
	0x75, 0x65, 0xf0, 0xb2, 0x60, 0x40, 0xa9, 0x36, 0x64, 0x73, 0x18, 0xa2, 0x94, 0x76, 0x24, 0xb3,
	0x8c, 0x5e, 0xc2, 0xcc, 0x18, 0xbf, 0x2f, 0x35, 0xfe, 0x31, 0xdc, 0x47, 0x98, 0x1e, 0x3c, 0xfe,
	0x29, 0xdd, 0x73, 0x98, 0xae, 0x74, 0xaa, 0x6b, 0xf5, 0x17, 0x6c, 0x93, 0xd6, 0xe1, 0xb7, 0x64,
	0x21, 0x78, 0xf8, 0x85, 0x2b, 0xcd, 0xc5, 0xd6, 0xa2, 0x75, 0xf1, 0x2f, 0xd8, 0xae, 0x60, 0xba,
	0x26, 0xf0, 0x96, 0xed, 0x11, 0x78, 0xb5, 0x7d, 0x3a, 0x96, 0x6c, 0x66, 0xc9, 0xda, 0x17, 0x15,
	0x77, 0x05, 0xad, 0xdf, 0xe0, 0xe0, 0x77, 0x1f, 0x26, 0xad, 0x9f, 0x21, 0xb5, 0x05, 0x4e, 0x57,
	0x70, 0xf1, 0xcd, 0x81, 0x9b, 0x6f, 0xca, 0x57, 0xf6, 0xdd, 0xae, 0x50, 0x7e, 0xe6, 0x1b, 0x64,
	0xcf, 0xc0, 0x6b, 0x37, 0x9f, 0x9d, 0xda, 0xff, 0x1d, 0x9d, 0x68, 0xb8, 0xf8, 0x49, 0xaf, 0xf2,
	0x7d, 0xf4, 0x1f, 0x7b, 0x0a, 0x63, 0x7b, 0xab, 0xd8, 0xed, 0x5e, 0xc9, 0xe1, 0xa2, 0x86, 0xb7,
	0x8e, 0xe5, 0xa6, 0xf1, 0x09, 0x8c, 0x9a, 0x5d, 0x65, 0x9d, 0x75, 0xff, 0x98, 0x42, 0x76, 0xa4,
	0x52, 0xd7, 0xf5, 0x88, 0xc4, 0xc7, 0x3f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x67, 0x3b, 0xf7, 0x94,
	0x84, 0x04, 0x00, 0x00,
}

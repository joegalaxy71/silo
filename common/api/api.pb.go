// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Void struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Void) Reset()         { *m = Void{} }
func (m *Void) String() string { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()    {}
func (*Void) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *Void) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Void.Unmarshal(m, b)
}
func (m *Void) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Void.Marshal(b, m, deterministic)
}
func (m *Void) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Void.Merge(m, src)
}
func (m *Void) XXX_Size() int {
	return xxx_messageInfo_Void.Size(m)
}
func (m *Void) XXX_DiscardUnknown() {
	xxx_messageInfo_Void.DiscardUnknown(m)
}

var xxx_messageInfo_Void proto.InternalMessageInfo

type Gps struct {
	Gp                   []*Gp    `protobuf:"bytes,1,rep,name=gp,proto3" json:"gp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Gps) Reset()         { *m = Gps{} }
func (m *Gps) String() string { return proto.CompactTextString(m) }
func (*Gps) ProtoMessage()    {}
func (*Gps) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *Gps) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Gps.Unmarshal(m, b)
}
func (m *Gps) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Gps.Marshal(b, m, deterministic)
}
func (m *Gps) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Gps.Merge(m, src)
}
func (m *Gps) XXX_Size() int {
	return xxx_messageInfo_Gps.Size(m)
}
func (m *Gps) XXX_DiscardUnknown() {
	xxx_messageInfo_Gps.DiscardUnknown(m)
}

var xxx_messageInfo_Gps proto.InternalMessageInfo

func (m *Gps) GetGp() []*Gp {
	if m != nil {
		return m.Gp
	}
	return nil
}

type Gp struct {
	GpId                 string   `protobuf:"bytes,1,opt,name=gpId,proto3" json:"gpId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Gp) Reset()         { *m = Gp{} }
func (m *Gp) String() string { return proto.CompactTextString(m) }
func (*Gp) ProtoMessage()    {}
func (*Gp) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *Gp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Gp.Unmarshal(m, b)
}
func (m *Gp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Gp.Marshal(b, m, deterministic)
}
func (m *Gp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Gp.Merge(m, src)
}
func (m *Gp) XXX_Size() int {
	return xxx_messageInfo_Gp.Size(m)
}
func (m *Gp) XXX_DiscardUnknown() {
	xxx_messageInfo_Gp.DiscardUnknown(m)
}

var xxx_messageInfo_Gp proto.InternalMessageInfo

func (m *Gp) GetGpId() string {
	if m != nil {
		return m.GpId
	}
	return ""
}

type Job struct {
	IdJob                string   `protobuf:"bytes,1,opt,name=IdJob,proto3" json:"IdJob,omitempty"`
	Type                 string   `protobuf:"bytes,2,opt,name=Type,proto3" json:"Type,omitempty"`
	Content              string   `protobuf:"bytes,3,opt,name=Content,proto3" json:"Content,omitempty"`
	GpId                 string   `protobuf:"bytes,4,opt,name=GpId,proto3" json:"GpId,omitempty"`
	Active               string   `protobuf:"bytes,5,opt,name=Active,proto3" json:"Active,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}
func (*Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Job.Unmarshal(m, b)
}
func (m *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Job.Marshal(b, m, deterministic)
}
func (m *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Job.Merge(m, src)
}
func (m *Job) XXX_Size() int {
	return xxx_messageInfo_Job.Size(m)
}
func (m *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Job.DiscardUnknown(m)
}

var xxx_messageInfo_Job proto.InternalMessageInfo

func (m *Job) GetIdJob() string {
	if m != nil {
		return m.IdJob
	}
	return ""
}

func (m *Job) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Job) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Job) GetGpId() string {
	if m != nil {
		return m.GpId
	}
	return ""
}

func (m *Job) GetActive() string {
	if m != nil {
		return m.Active
	}
	return ""
}

type Jobs struct {
	Job                  []*Job   `protobuf:"bytes,1,rep,name=job,proto3" json:"job,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Jobs) Reset()         { *m = Jobs{} }
func (m *Jobs) String() string { return proto.CompactTextString(m) }
func (*Jobs) ProtoMessage()    {}
func (*Jobs) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *Jobs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Jobs.Unmarshal(m, b)
}
func (m *Jobs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Jobs.Marshal(b, m, deterministic)
}
func (m *Jobs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Jobs.Merge(m, src)
}
func (m *Jobs) XXX_Size() int {
	return xxx_messageInfo_Jobs.Size(m)
}
func (m *Jobs) XXX_DiscardUnknown() {
	xxx_messageInfo_Jobs.DiscardUnknown(m)
}

var xxx_messageInfo_Jobs proto.InternalMessageInfo

func (m *Jobs) GetJob() []*Job {
	if m != nil {
		return m.Job
	}
	return nil
}

type Results struct {
	Result               []*Result `protobuf:"bytes,1,rep,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Results) Reset()         { *m = Results{} }
func (m *Results) String() string { return proto.CompactTextString(m) }
func (*Results) ProtoMessage()    {}
func (*Results) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *Results) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Results.Unmarshal(m, b)
}
func (m *Results) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Results.Marshal(b, m, deterministic)
}
func (m *Results) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Results.Merge(m, src)
}
func (m *Results) XXX_Size() int {
	return xxx_messageInfo_Results.Size(m)
}
func (m *Results) XXX_DiscardUnknown() {
	xxx_messageInfo_Results.DiscardUnknown(m)
}

var xxx_messageInfo_Results proto.InternalMessageInfo

func (m *Results) GetResult() []*Result {
	if m != nil {
		return m.Result
	}
	return nil
}

type Entry struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}

func (m *Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entry.Unmarshal(m, b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
}
func (m *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(m, src)
}
func (m *Entry) XXX_Size() int {
	return xxx_messageInfo_Entry.Size(m)
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Entry) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Result struct {
	IdJob                string   `protobuf:"bytes,1,opt,name=IdJob,proto3" json:"IdJob,omitempty"`
	Type                 string   `protobuf:"bytes,2,opt,name=Type,proto3" json:"Type,omitempty"`
	Entries              []*Entry `protobuf:"bytes,3,rep,name=Entries,proto3" json:"Entries,omitempty"`
	GpId                 string   `protobuf:"bytes,4,opt,name=GpId,proto3" json:"GpId,omitempty"`
	Elapsed              int64    `protobuf:"varint,5,opt,name=Elapsed,proto3" json:"Elapsed,omitempty"`
	GpName               string   `protobuf:"bytes,6,opt,name=GpName,proto3" json:"GpName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{7}
}

func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetIdJob() string {
	if m != nil {
		return m.IdJob
	}
	return ""
}

func (m *Result) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Result) GetEntries() []*Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *Result) GetGpId() string {
	if m != nil {
		return m.GpId
	}
	return ""
}

func (m *Result) GetElapsed() int64 {
	if m != nil {
		return m.Elapsed
	}
	return 0
}

func (m *Result) GetGpName() string {
	if m != nil {
		return m.GpName
	}
	return ""
}

type Status struct {
	Code                 bool     `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{8}
}

func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetCode() bool {
	if m != nil {
		return m.Code
	}
	return false
}

func (m *Status) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*Void)(nil), "api.Void")
	proto.RegisterType((*Gps)(nil), "api.Gps")
	proto.RegisterType((*Gp)(nil), "api.Gp")
	proto.RegisterType((*Job)(nil), "api.Job")
	proto.RegisterType((*Jobs)(nil), "api.Jobs")
	proto.RegisterType((*Results)(nil), "api.Results")
	proto.RegisterType((*Entry)(nil), "api.Entry")
	proto.RegisterType((*Result)(nil), "api.Result")
	proto.RegisterType((*Status)(nil), "api.Status")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 401 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xc1, 0x8a, 0xdb, 0x30,
	0x10, 0x8d, 0x2d, 0x47, 0x8e, 0x67, 0xdb, 0x8b, 0x28, 0xad, 0xc8, 0xa1, 0x04, 0x75, 0xa1, 0x4b,
	0x0f, 0x81, 0x6e, 0xa1, 0xf7, 0xa5, 0x2c, 0x26, 0x3e, 0x94, 0xa2, 0x94, 0xdc, 0xed, 0x58, 0x18,
	0x17, 0x27, 0x12, 0x96, 0x1c, 0xc8, 0xcf, 0xf4, 0xd6, 0xff, 0x2c, 0x1a, 0x59, 0x90, 0x42, 0x2f,
	0x7b, 0x9b, 0xa7, 0x79, 0x7a, 0x7a, 0xf3, 0x34, 0x50, 0xd4, 0xa6, 0xdf, 0x9a, 0x51, 0x3b, 0xcd,
	0x48, 0x6d, 0x7a, 0x41, 0x21, 0x3b, 0xe8, 0xbe, 0x15, 0xef, 0x81, 0x94, 0xc6, 0xb2, 0x77, 0x90,
	0x76, 0x86, 0x27, 0x1b, 0xf2, 0x70, 0xf7, 0x98, 0x6f, 0x3d, 0xb7, 0x34, 0x32, 0xed, 0x8c, 0xe0,
	0x90, 0x96, 0x86, 0x31, 0xc8, 0x3a, 0xb3, 0x6b, 0x79, 0xb2, 0x49, 0x1e, 0x0a, 0x89, 0xb5, 0x98,
	0x80, 0x54, 0xba, 0x61, 0x6f, 0x60, 0xb9, 0x6b, 0x2b, 0xdd, 0xcc, 0xbd, 0x00, 0xfc, 0x85, 0x9f,
	0x57, 0xa3, 0x78, 0x1a, 0x2e, 0xf8, 0x9a, 0x71, 0xc8, 0xbf, 0xe9, 0xb3, 0x53, 0x67, 0xc7, 0x09,
	0x1e, 0x47, 0xe8, 0xd9, 0xa5, 0x97, 0xcf, 0x02, 0xdb, 0xd7, 0xec, 0x2d, 0xd0, 0xa7, 0xa3, 0xeb,
	0x2f, 0x8a, 0x2f, 0xf1, 0x74, 0x46, 0x42, 0x40, 0x56, 0xe9, 0xc6, 0xb2, 0x35, 0x90, 0x5f, 0xf8,
	0xaa, 0xb7, 0xbc, 0x42, 0xcb, 0x95, 0x6e, 0xa4, 0x3f, 0x14, 0x5b, 0xc8, 0xa5, 0xb2, 0xd3, 0xe0,
	0x2c, 0xfb, 0x00, 0x74, 0xc4, 0x72, 0x66, 0xde, 0x21, 0x33, 0x74, 0xe5, 0xdc, 0x12, 0x9f, 0x61,
	0xf9, 0x7c, 0x76, 0xe3, 0xd5, 0x1b, 0xf9, 0x5e, 0x9f, 0x54, 0x9c, 0xd3, 0xd7, 0x7e, 0xc0, 0x43,
	0x3d, 0x4c, 0x71, 0x96, 0x00, 0xc4, 0xef, 0x04, 0x68, 0x50, 0x79, 0x41, 0x02, 0xf7, 0x90, 0xfb,
	0x77, 0x7a, 0x65, 0x39, 0x41, 0x37, 0x80, 0x6e, 0xf0, 0x6d, 0x19, 0x5b, 0xff, 0x4d, 0x83, 0x43,
	0xfe, 0x3c, 0xd4, 0xc6, 0xaa, 0x16, 0xe3, 0x20, 0x32, 0x42, 0x9f, 0x53, 0x69, 0xd0, 0x34, 0x0d,
	0x39, 0x05, 0x24, 0xbe, 0x02, 0xdd, 0xbb, 0xda, 0x4d, 0xa8, 0x77, 0xd4, 0x6d, 0x18, 0x6a, 0x25,
	0xb1, 0xf6, 0x7a, 0x27, 0x65, 0x6d, 0xdd, 0x45, 0x83, 0x11, 0x3e, 0xfe, 0x49, 0x80, 0xee, 0xd5,
	0x78, 0x51, 0x23, 0xdb, 0x40, 0x5e, 0x2a, 0x87, 0x69, 0xaf, 0xe6, 0x9d, 0xb0, 0xeb, 0x22, 0x46,
	0x6d, 0xc5, 0x82, 0xdd, 0x03, 0x94, 0xca, 0x3d, 0x0d, 0x03, 0x92, 0x42, 0xcb, 0xaf, 0xd5, 0xbf,
	0xac, 0x8f, 0x50, 0xfc, 0x98, 0xdc, 0x9c, 0xd6, 0xed, 0x07, 0xac, 0x03, 0x08, 0x3e, 0xc5, 0x82,
	0x7d, 0x82, 0xd7, 0x41, 0x2e, 0xfe, 0xde, 0x8d, 0xe2, 0xab, 0x9b, 0x7b, 0x56, 0x2c, 0x1a, 0x8a,
	0xcb, 0xfc, 0xe5, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x46, 0x2b, 0x44, 0xd9, 0x02, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServerClient interface {
	GetJobs(ctx context.Context, in *Gps, opts ...grpc.CallOption) (*Jobs, error)
	GetAllJobs(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Jobs, error)
	PutResult(ctx context.Context, in *Result, opts ...grpc.CallOption) (*Status, error)
	GetAllResults(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Results, error)
}

type serverClient struct {
	cc *grpc.ClientConn
}

func NewServerClient(cc *grpc.ClientConn) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) GetJobs(ctx context.Context, in *Gps, opts ...grpc.CallOption) (*Jobs, error) {
	out := new(Jobs)
	err := c.cc.Invoke(ctx, "/api.Server/GetJobs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetAllJobs(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Jobs, error) {
	out := new(Jobs)
	err := c.cc.Invoke(ctx, "/api.Server/GetAllJobs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) PutResult(ctx context.Context, in *Result, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/api.Server/PutResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetAllResults(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Results, error) {
	out := new(Results)
	err := c.cc.Invoke(ctx, "/api.Server/GetAllResults", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServer is the server API for Server service.
type ServerServer interface {
	GetJobs(context.Context, *Gps) (*Jobs, error)
	GetAllJobs(context.Context, *Void) (*Jobs, error)
	PutResult(context.Context, *Result) (*Status, error)
	GetAllResults(context.Context, *Void) (*Results, error)
}

// UnimplementedServerServer can be embedded to have forward compatible implementations.
type UnimplementedServerServer struct {
}

func (*UnimplementedServerServer) GetJobs(ctx context.Context, req *Gps) (*Jobs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJobs not implemented")
}
func (*UnimplementedServerServer) GetAllJobs(ctx context.Context, req *Void) (*Jobs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllJobs not implemented")
}
func (*UnimplementedServerServer) PutResult(ctx context.Context, req *Result) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutResult not implemented")
}
func (*UnimplementedServerServer) GetAllResults(ctx context.Context, req *Void) (*Results, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllResults not implemented")
}

func RegisterServerServer(s *grpc.Server, srv ServerServer) {
	s.RegisterService(&_Server_serviceDesc, srv)
}

func _Server_GetJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Gps)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Server/GetJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetJobs(ctx, req.(*Gps))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetAllJobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetAllJobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Server/GetAllJobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetAllJobs(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_PutResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Result)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).PutResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Server/PutResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).PutResult(ctx, req.(*Result))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetAllResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetAllResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Server/GetAllResults",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetAllResults(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

var _Server_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Server",
	HandlerType: (*ServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetJobs",
			Handler:    _Server_GetJobs_Handler,
		},
		{
			MethodName: "GetAllJobs",
			Handler:    _Server_GetAllJobs_Handler,
		},
		{
			MethodName: "PutResult",
			Handler:    _Server_PutResult_Handler,
		},
		{
			MethodName: "GetAllResults",
			Handler:    _Server_GetAllResults_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: daemon.proto

package daemon

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

type CreateGameRequest struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ScenarioNo           uint32   `protobuf:"varint,3,opt,name=scenarioNo,proto3" json:"scenarioNo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateGameRequest) Reset()         { *m = CreateGameRequest{} }
func (m *CreateGameRequest) String() string { return proto.CompactTextString(m) }
func (*CreateGameRequest) ProtoMessage()    {}
func (*CreateGameRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{0}
}

func (m *CreateGameRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateGameRequest.Unmarshal(m, b)
}
func (m *CreateGameRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateGameRequest.Marshal(b, m, deterministic)
}
func (m *CreateGameRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateGameRequest.Merge(m, src)
}
func (m *CreateGameRequest) XXX_Size() int {
	return xxx_messageInfo_CreateGameRequest.Size(m)
}
func (m *CreateGameRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateGameRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateGameRequest proto.InternalMessageInfo

func (m *CreateGameRequest) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *CreateGameRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateGameRequest) GetScenarioNo() uint32 {
	if m != nil {
		return m.ScenarioNo
	}
	return 0
}

type CreateGameResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateGameResponse) Reset()         { *m = CreateGameResponse{} }
func (m *CreateGameResponse) String() string { return proto.CompactTextString(m) }
func (*CreateGameResponse) ProtoMessage()    {}
func (*CreateGameResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{1}
}

func (m *CreateGameResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateGameResponse.Unmarshal(m, b)
}
func (m *CreateGameResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateGameResponse.Marshal(b, m, deterministic)
}
func (m *CreateGameResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateGameResponse.Merge(m, src)
}
func (m *CreateGameResponse) XXX_Size() int {
	return xxx_messageInfo_CreateGameResponse.Size(m)
}
func (m *CreateGameResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateGameResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateGameResponse proto.InternalMessageInfo

func (m *CreateGameResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type StopGameRequest struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopGameRequest) Reset()         { *m = StopGameRequest{} }
func (m *StopGameRequest) String() string { return proto.CompactTextString(m) }
func (*StopGameRequest) ProtoMessage()    {}
func (*StopGameRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{2}
}

func (m *StopGameRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopGameRequest.Unmarshal(m, b)
}
func (m *StopGameRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopGameRequest.Marshal(b, m, deterministic)
}
func (m *StopGameRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopGameRequest.Merge(m, src)
}
func (m *StopGameRequest) XXX_Size() int {
	return xxx_messageInfo_StopGameRequest.Size(m)
}
func (m *StopGameRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StopGameRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StopGameRequest proto.InternalMessageInfo

func (m *StopGameRequest) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

type StopGameResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopGameResponse) Reset()         { *m = StopGameResponse{} }
func (m *StopGameResponse) String() string { return proto.CompactTextString(m) }
func (*StopGameResponse) ProtoMessage()    {}
func (*StopGameResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{3}
}

func (m *StopGameResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopGameResponse.Unmarshal(m, b)
}
func (m *StopGameResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopGameResponse.Marshal(b, m, deterministic)
}
func (m *StopGameResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopGameResponse.Merge(m, src)
}
func (m *StopGameResponse) XXX_Size() int {
	return xxx_messageInfo_StopGameResponse.Size(m)
}
func (m *StopGameResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StopGameResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StopGameResponse proto.InternalMessageInfo

func (m *StopGameResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type ListGamesResponse struct {
	Games                []*CreateGameRequest `protobuf:"bytes,1,rep,name=games,proto3" json:"games,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ListGamesResponse) Reset()         { *m = ListGamesResponse{} }
func (m *ListGamesResponse) String() string { return proto.CompactTextString(m) }
func (*ListGamesResponse) ProtoMessage()    {}
func (*ListGamesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{4}
}

func (m *ListGamesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListGamesResponse.Unmarshal(m, b)
}
func (m *ListGamesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListGamesResponse.Marshal(b, m, deterministic)
}
func (m *ListGamesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListGamesResponse.Merge(m, src)
}
func (m *ListGamesResponse) XXX_Size() int {
	return xxx_messageInfo_ListGamesResponse.Size(m)
}
func (m *ListGamesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListGamesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListGamesResponse proto.InternalMessageInfo

func (m *ListGamesResponse) GetGames() []*CreateGameRequest {
	if m != nil {
		return m.Games
	}
	return nil
}

type EmptyRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyRequest) Reset()         { *m = EmptyRequest{} }
func (m *EmptyRequest) String() string { return proto.CompactTextString(m) }
func (*EmptyRequest) ProtoMessage()    {}
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{5}
}

func (m *EmptyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyRequest.Unmarshal(m, b)
}
func (m *EmptyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyRequest.Marshal(b, m, deterministic)
}
func (m *EmptyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyRequest.Merge(m, src)
}
func (m *EmptyRequest) XXX_Size() int {
	return xxx_messageInfo_EmptyRequest.Size(m)
}
func (m *EmptyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyRequest proto.InternalMessageInfo

type ListScenariosResponse struct {
	Scenarios            []*ListScenariosResponse_Scenario `protobuf:"bytes,1,rep,name=scenarios,proto3" json:"scenarios,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                          `json:"-"`
	XXX_unrecognized     []byte                            `json:"-"`
	XXX_sizecache        int32                             `json:"-"`
}

func (m *ListScenariosResponse) Reset()         { *m = ListScenariosResponse{} }
func (m *ListScenariosResponse) String() string { return proto.CompactTextString(m) }
func (*ListScenariosResponse) ProtoMessage()    {}
func (*ListScenariosResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{6}
}

func (m *ListScenariosResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListScenariosResponse.Unmarshal(m, b)
}
func (m *ListScenariosResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListScenariosResponse.Marshal(b, m, deterministic)
}
func (m *ListScenariosResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListScenariosResponse.Merge(m, src)
}
func (m *ListScenariosResponse) XXX_Size() int {
	return xxx_messageInfo_ListScenariosResponse.Size(m)
}
func (m *ListScenariosResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListScenariosResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListScenariosResponse proto.InternalMessageInfo

func (m *ListScenariosResponse) GetScenarios() []*ListScenariosResponse_Scenario {
	if m != nil {
		return m.Scenarios
	}
	return nil
}

type ListScenariosResponse_Network struct {
	Challenges           []string `protobuf:"bytes,1,rep,name=challenges,proto3" json:"challenges,omitempty"`
	Vlan                 string   `protobuf:"bytes,2,opt,name=vlan,proto3" json:"vlan,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListScenariosResponse_Network) Reset()         { *m = ListScenariosResponse_Network{} }
func (m *ListScenariosResponse_Network) String() string { return proto.CompactTextString(m) }
func (*ListScenariosResponse_Network) ProtoMessage()    {}
func (*ListScenariosResponse_Network) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{6, 0}
}

func (m *ListScenariosResponse_Network) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListScenariosResponse_Network.Unmarshal(m, b)
}
func (m *ListScenariosResponse_Network) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListScenariosResponse_Network.Marshal(b, m, deterministic)
}
func (m *ListScenariosResponse_Network) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListScenariosResponse_Network.Merge(m, src)
}
func (m *ListScenariosResponse_Network) XXX_Size() int {
	return xxx_messageInfo_ListScenariosResponse_Network.Size(m)
}
func (m *ListScenariosResponse_Network) XXX_DiscardUnknown() {
	xxx_messageInfo_ListScenariosResponse_Network.DiscardUnknown(m)
}

var xxx_messageInfo_ListScenariosResponse_Network proto.InternalMessageInfo

func (m *ListScenariosResponse_Network) GetChallenges() []string {
	if m != nil {
		return m.Challenges
	}
	return nil
}

func (m *ListScenariosResponse_Network) GetVlan() string {
	if m != nil {
		return m.Vlan
	}
	return ""
}

type ListScenariosResponse_Scenario struct {
	Networks             []*ListScenariosResponse_Network `protobuf:"bytes,1,rep,name=networks,proto3" json:"networks,omitempty"`
	Duration             uint32                           `protobuf:"varint,2,opt,name=duration,proto3" json:"duration,omitempty"`
	Difficulty           string                           `protobuf:"bytes,3,opt,name=difficulty,proto3" json:"difficulty,omitempty"`
	Story                string                           `protobuf:"bytes,4,opt,name=story,proto3" json:"story,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *ListScenariosResponse_Scenario) Reset()         { *m = ListScenariosResponse_Scenario{} }
func (m *ListScenariosResponse_Scenario) String() string { return proto.CompactTextString(m) }
func (*ListScenariosResponse_Scenario) ProtoMessage()    {}
func (*ListScenariosResponse_Scenario) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{6, 1}
}

func (m *ListScenariosResponse_Scenario) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListScenariosResponse_Scenario.Unmarshal(m, b)
}
func (m *ListScenariosResponse_Scenario) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListScenariosResponse_Scenario.Marshal(b, m, deterministic)
}
func (m *ListScenariosResponse_Scenario) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListScenariosResponse_Scenario.Merge(m, src)
}
func (m *ListScenariosResponse_Scenario) XXX_Size() int {
	return xxx_messageInfo_ListScenariosResponse_Scenario.Size(m)
}
func (m *ListScenariosResponse_Scenario) XXX_DiscardUnknown() {
	xxx_messageInfo_ListScenariosResponse_Scenario.DiscardUnknown(m)
}

var xxx_messageInfo_ListScenariosResponse_Scenario proto.InternalMessageInfo

func (m *ListScenariosResponse_Scenario) GetNetworks() []*ListScenariosResponse_Network {
	if m != nil {
		return m.Networks
	}
	return nil
}

func (m *ListScenariosResponse_Scenario) GetDuration() uint32 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *ListScenariosResponse_Scenario) GetDifficulty() string {
	if m != nil {
		return m.Difficulty
	}
	return ""
}

func (m *ListScenariosResponse_Scenario) GetStory() string {
	if m != nil {
		return m.Story
	}
	return ""
}

func init() {
	proto.RegisterType((*CreateGameRequest)(nil), "CreateGameRequest")
	proto.RegisterType((*CreateGameResponse)(nil), "CreateGameResponse")
	proto.RegisterType((*StopGameRequest)(nil), "StopGameRequest")
	proto.RegisterType((*StopGameResponse)(nil), "StopGameResponse")
	proto.RegisterType((*ListGamesResponse)(nil), "ListGamesResponse")
	proto.RegisterType((*EmptyRequest)(nil), "EmptyRequest")
	proto.RegisterType((*ListScenariosResponse)(nil), "ListScenariosResponse")
	proto.RegisterType((*ListScenariosResponse_Network)(nil), "ListScenariosResponse.Network")
	proto.RegisterType((*ListScenariosResponse_Scenario)(nil), "ListScenariosResponse.Scenario")
}

func init() { proto.RegisterFile("daemon.proto", fileDescriptor_3ec90cbc4aa12fc6) }

var fileDescriptor_3ec90cbc4aa12fc6 = []byte{
	// 414 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xc1, 0x8e, 0xd3, 0x30,
	0x10, 0x4d, 0xb6, 0xed, 0x6e, 0x32, 0x6c, 0x20, 0x1d, 0x58, 0x14, 0xe5, 0xb0, 0x54, 0xe6, 0x92,
	0x03, 0xb2, 0x60, 0x39, 0x80, 0x90, 0x7a, 0x02, 0x84, 0x84, 0xd0, 0x1e, 0xdc, 0x0f, 0x40, 0xa6,
	0x75, 0x43, 0x44, 0x12, 0x87, 0xd8, 0x05, 0xf5, 0x67, 0xf8, 0x0c, 0xbe, 0x87, 0x4f, 0x41, 0x4e,
	0xe2, 0x24, 0xb4, 0x05, 0x6e, 0x9e, 0x97, 0xf1, 0x9b, 0xf1, 0x7b, 0x2f, 0x70, 0xb9, 0xe1, 0xa2,
	0x90, 0x25, 0xad, 0x6a, 0xa9, 0x25, 0xf9, 0x08, 0xf3, 0xd7, 0xb5, 0xe0, 0x5a, 0xbc, 0xe3, 0x85,
	0x60, 0xe2, 0xeb, 0x4e, 0x28, 0x8d, 0x21, 0x4c, 0x34, 0x4f, 0x23, 0x77, 0xe1, 0x26, 0x3e, 0x33,
	0x47, 0x44, 0x98, 0x96, 0xbc, 0x10, 0xd1, 0x59, 0x03, 0x35, 0x67, 0xbc, 0x06, 0x50, 0x6b, 0x51,
	0xf2, 0x3a, 0x93, 0xb7, 0x32, 0x9a, 0x2c, 0xdc, 0x24, 0x60, 0x23, 0xe4, 0xfd, 0xd4, 0x9b, 0x85,
	0x57, 0x84, 0x02, 0x8e, 0x07, 0xa8, 0x4a, 0x96, 0x4a, 0x60, 0x04, 0x17, 0x85, 0x50, 0x8a, 0xa7,
	0xa2, 0x9b, 0x62, 0x4b, 0xf2, 0x18, 0xee, 0xad, 0xb4, 0xac, 0xfe, 0xb9, 0x0e, 0x79, 0x02, 0xe1,
	0xd0, 0xf4, 0x5f, 0xca, 0x25, 0xcc, 0x3f, 0x64, 0x4a, 0x9b, 0x6e, 0xd5, 0xb7, 0x27, 0x30, 0x4b,
	0x0d, 0x10, 0xb9, 0x8b, 0x49, 0x72, 0xe7, 0x06, 0xe9, 0x91, 0x0c, 0xac, 0x6d, 0x20, 0x77, 0xe1,
	0xf2, 0x6d, 0x51, 0xe9, 0x7d, 0x07, 0x93, 0x9f, 0x67, 0x70, 0x65, 0xf8, 0x56, 0xdd, 0x53, 0x07,
	0xce, 0x25, 0xf8, 0xf6, 0xfd, 0x96, 0xf7, 0x11, 0x3d, 0xd9, 0x4a, 0x2d, 0xc2, 0x86, 0x1b, 0xf1,
	0x12, 0x2e, 0x6e, 0x85, 0xfe, 0x2e, 0xeb, 0x2f, 0x46, 0xdb, 0xf5, 0x67, 0x9e, 0xe7, 0xa2, 0x4c,
	0xbb, 0x15, 0x7d, 0x36, 0x42, 0x8c, 0x1f, 0xdf, 0x72, 0x5e, 0x5a, 0x3f, 0xcc, 0x39, 0xfe, 0xe1,
	0x82, 0x67, 0x69, 0xf1, 0x15, 0x78, 0x65, 0xcb, 0x65, 0x37, 0xb9, 0xfe, 0xcb, 0x26, 0xdd, 0x48,
	0xd6, 0xf7, 0x63, 0x0c, 0xde, 0x66, 0x57, 0x73, 0x9d, 0xc9, 0x76, 0x40, 0xc0, 0xfa, 0xda, 0x2c,
	0xb6, 0xc9, 0xb6, 0xdb, 0x6c, 0xbd, 0xcb, 0xf5, 0xbe, 0x31, 0xdd, 0x67, 0x23, 0x04, 0x1f, 0xc0,
	0x4c, 0x69, 0x59, 0xef, 0xa3, 0x69, 0xf3, 0xa9, 0x2d, 0xda, 0x28, 0xdc, 0xfc, 0x72, 0xe1, 0xfc,
	0x4d, 0x13, 0x3e, 0x7c, 0x01, 0x30, 0xe8, 0x8d, 0x27, 0xc4, 0x8f, 0xef, 0xd3, 0xe3, 0xd8, 0x10,
	0x07, 0x9f, 0x81, 0x67, 0x9d, 0xc7, 0x90, 0x1e, 0x24, 0x25, 0x9e, 0xd3, 0xc3, 0x58, 0x10, 0x07,
	0x9f, 0x82, 0xdf, 0xdb, 0x8f, 0x01, 0x1d, 0x7b, 0x19, 0x23, 0x3d, 0x4a, 0x06, 0x71, 0xf0, 0x25,
	0x04, 0x7f, 0x68, 0x75, 0x78, 0xeb, 0xe1, 0x69, 0x29, 0x89, 0xf3, 0xe9, 0xbc, 0xf9, 0xab, 0x9e,
	0xff, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x6d, 0x1d, 0xed, 0xf4, 0x65, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DaemonClient is the client API for Daemon service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DaemonClient interface {
	CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameResponse, error)
	StopGame(ctx context.Context, in *StopGameRequest, opts ...grpc.CallOption) (*StopGameResponse, error)
	ListGames(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*ListGamesResponse, error)
	ListScenarios(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*ListScenariosResponse, error)
}

type daemonClient struct {
	cc *grpc.ClientConn
}

func NewDaemonClient(cc *grpc.ClientConn) DaemonClient {
	return &daemonClient{cc}
}

func (c *daemonClient) CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameResponse, error) {
	out := new(CreateGameResponse)
	err := c.cc.Invoke(ctx, "/Daemon/CreateGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) StopGame(ctx context.Context, in *StopGameRequest, opts ...grpc.CallOption) (*StopGameResponse, error) {
	out := new(StopGameResponse)
	err := c.cc.Invoke(ctx, "/Daemon/StopGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) ListGames(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*ListGamesResponse, error) {
	out := new(ListGamesResponse)
	err := c.cc.Invoke(ctx, "/Daemon/ListGames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) ListScenarios(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*ListScenariosResponse, error) {
	out := new(ListScenariosResponse)
	err := c.cc.Invoke(ctx, "/Daemon/ListScenarios", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaemonServer is the server API for Daemon service.
type DaemonServer interface {
	CreateGame(context.Context, *CreateGameRequest) (*CreateGameResponse, error)
	StopGame(context.Context, *StopGameRequest) (*StopGameResponse, error)
	ListGames(context.Context, *EmptyRequest) (*ListGamesResponse, error)
	ListScenarios(context.Context, *EmptyRequest) (*ListScenariosResponse, error)
}

// UnimplementedDaemonServer can be embedded to have forward compatible implementations.
type UnimplementedDaemonServer struct {
}

func (*UnimplementedDaemonServer) CreateGame(ctx context.Context, req *CreateGameRequest) (*CreateGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGame not implemented")
}
func (*UnimplementedDaemonServer) StopGame(ctx context.Context, req *StopGameRequest) (*StopGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopGame not implemented")
}
func (*UnimplementedDaemonServer) ListGames(ctx context.Context, req *EmptyRequest) (*ListGamesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGames not implemented")
}
func (*UnimplementedDaemonServer) ListScenarios(ctx context.Context, req *EmptyRequest) (*ListScenariosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListScenarios not implemented")
}

func RegisterDaemonServer(s *grpc.Server, srv DaemonServer) {
	s.RegisterService(&_Daemon_serviceDesc, srv)
}

func _Daemon_CreateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).CreateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/CreateGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).CreateGame(ctx, req.(*CreateGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_StopGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).StopGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/StopGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).StopGame(ctx, req.(*StopGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_ListGames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).ListGames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/ListGames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).ListGames(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_ListScenarios_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).ListScenarios(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/ListScenarios",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).ListScenarios(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Daemon_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Daemon",
	HandlerType: (*DaemonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGame",
			Handler:    _Daemon_CreateGame_Handler,
		},
		{
			MethodName: "StopGame",
			Handler:    _Daemon_StopGame_Handler,
		},
		{
			MethodName: "ListGames",
			Handler:    _Daemon_ListGames_Handler,
		},
		{
			MethodName: "ListScenarios",
			Handler:    _Daemon_ListScenarios_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "daemon.proto",
}

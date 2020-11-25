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

type Network struct {
	Challenges           []string `protobuf:"bytes,1,rep,name=challenges,proto3" json:"challenges,omitempty"`
	Vlan                 string   `protobuf:"bytes,2,opt,name=vlan,proto3" json:"vlan,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Network) Reset()         { *m = Network{} }
func (m *Network) String() string { return proto.CompactTextString(m) }
func (*Network) ProtoMessage()    {}
func (*Network) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{6}
}

func (m *Network) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Network.Unmarshal(m, b)
}
func (m *Network) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Network.Marshal(b, m, deterministic)
}
func (m *Network) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Network.Merge(m, src)
}
func (m *Network) XXX_Size() int {
	return xxx_messageInfo_Network.Size(m)
}
func (m *Network) XXX_DiscardUnknown() {
	xxx_messageInfo_Network.DiscardUnknown(m)
}

var xxx_messageInfo_Network proto.InternalMessageInfo

func (m *Network) GetChallenges() []string {
	if m != nil {
		return m.Challenges
	}
	return nil
}

func (m *Network) GetVlan() string {
	if m != nil {
		return m.Vlan
	}
	return ""
}

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
	return fileDescriptor_3ec90cbc4aa12fc6, []int{7}
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

type ListScenariosResponse_Scenario struct {
	Networks             []*Network `protobuf:"bytes,1,rep,name=networks,proto3" json:"networks,omitempty"`
	Duration             uint32     `protobuf:"varint,2,opt,name=duration,proto3" json:"duration,omitempty"`
	Difficulty           string     `protobuf:"bytes,3,opt,name=difficulty,proto3" json:"difficulty,omitempty"`
	Story                string     `protobuf:"bytes,4,opt,name=story,proto3" json:"story,omitempty"`
	Id                   uint32     `protobuf:"varint,5,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ListScenariosResponse_Scenario) Reset()         { *m = ListScenariosResponse_Scenario{} }
func (m *ListScenariosResponse_Scenario) String() string { return proto.CompactTextString(m) }
func (*ListScenariosResponse_Scenario) ProtoMessage()    {}
func (*ListScenariosResponse_Scenario) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{7, 0}
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

func (m *ListScenariosResponse_Scenario) GetNetworks() []*Network {
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

func (m *ListScenariosResponse_Scenario) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type ListScenarioChallengesReq struct {
	ScenarioId           uint32   `protobuf:"varint,1,opt,name=scenario_id,json=scenarioId,proto3" json:"scenario_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListScenarioChallengesReq) Reset()         { *m = ListScenarioChallengesReq{} }
func (m *ListScenarioChallengesReq) String() string { return proto.CompactTextString(m) }
func (*ListScenarioChallengesReq) ProtoMessage()    {}
func (*ListScenarioChallengesReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{8}
}

func (m *ListScenarioChallengesReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListScenarioChallengesReq.Unmarshal(m, b)
}
func (m *ListScenarioChallengesReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListScenarioChallengesReq.Marshal(b, m, deterministic)
}
func (m *ListScenarioChallengesReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListScenarioChallengesReq.Merge(m, src)
}
func (m *ListScenarioChallengesReq) XXX_Size() int {
	return xxx_messageInfo_ListScenarioChallengesReq.Size(m)
}
func (m *ListScenarioChallengesReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ListScenarioChallengesReq.DiscardUnknown(m)
}

var xxx_messageInfo_ListScenarioChallengesReq proto.InternalMessageInfo

func (m *ListScenarioChallengesReq) GetScenarioId() uint32 {
	if m != nil {
		return m.ScenarioId
	}
	return 0
}

type ListScenarioChallengesResp struct {
	Chals                []*Network `protobuf:"bytes,1,rep,name=chals,proto3" json:"chals,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ListScenarioChallengesResp) Reset()         { *m = ListScenarioChallengesResp{} }
func (m *ListScenarioChallengesResp) String() string { return proto.CompactTextString(m) }
func (*ListScenarioChallengesResp) ProtoMessage()    {}
func (*ListScenarioChallengesResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ec90cbc4aa12fc6, []int{9}
}

func (m *ListScenarioChallengesResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListScenarioChallengesResp.Unmarshal(m, b)
}
func (m *ListScenarioChallengesResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListScenarioChallengesResp.Marshal(b, m, deterministic)
}
func (m *ListScenarioChallengesResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListScenarioChallengesResp.Merge(m, src)
}
func (m *ListScenarioChallengesResp) XXX_Size() int {
	return xxx_messageInfo_ListScenarioChallengesResp.Size(m)
}
func (m *ListScenarioChallengesResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ListScenarioChallengesResp.DiscardUnknown(m)
}

var xxx_messageInfo_ListScenarioChallengesResp proto.InternalMessageInfo

func (m *ListScenarioChallengesResp) GetChals() []*Network {
	if m != nil {
		return m.Chals
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateGameRequest)(nil), "CreateGameRequest")
	proto.RegisterType((*CreateGameResponse)(nil), "CreateGameResponse")
	proto.RegisterType((*StopGameRequest)(nil), "StopGameRequest")
	proto.RegisterType((*StopGameResponse)(nil), "StopGameResponse")
	proto.RegisterType((*ListGamesResponse)(nil), "ListGamesResponse")
	proto.RegisterType((*EmptyRequest)(nil), "EmptyRequest")
	proto.RegisterType((*Network)(nil), "Network")
	proto.RegisterType((*ListScenariosResponse)(nil), "ListScenariosResponse")
	proto.RegisterType((*ListScenariosResponse_Scenario)(nil), "ListScenariosResponse.Scenario")
	proto.RegisterType((*ListScenarioChallengesReq)(nil), "ListScenarioChallengesReq")
	proto.RegisterType((*ListScenarioChallengesResp)(nil), "ListScenarioChallengesResp")
}

func init() { proto.RegisterFile("daemon.proto", fileDescriptor_3ec90cbc4aa12fc6) }

var fileDescriptor_3ec90cbc4aa12fc6 = []byte{
	// 483 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x8e, 0x93, 0x38, 0xb5, 0xa7, 0x4d, 0x71, 0x06, 0x82, 0x8c, 0x91, 0xda, 0x68, 0xe1, 0x90,
	0x03, 0x5a, 0x41, 0x39, 0xc0, 0xa1, 0x39, 0x15, 0x84, 0x88, 0x50, 0x0f, 0xdb, 0x07, 0xa8, 0x96,
	0x78, 0x1b, 0x2c, 0x62, 0xaf, 0xeb, 0xdd, 0x80, 0xf2, 0x26, 0x3c, 0x04, 0x2f, 0xc7, 0x1b, 0xa0,
	0xf5, 0x3f, 0xf9, 0xa1, 0xb7, 0x99, 0xcf, 0x33, 0xb3, 0xdf, 0x7c, 0xf3, 0xc9, 0x70, 0x12, 0x72,
	0x11, 0xcb, 0x84, 0xa6, 0x99, 0xd4, 0x92, 0xdc, 0xc2, 0xe8, 0x2a, 0x13, 0x5c, 0x8b, 0x4f, 0x3c,
	0x16, 0x4c, 0xdc, 0xaf, 0x85, 0xd2, 0xe8, 0x41, 0x4f, 0xf3, 0xa5, 0x6f, 0x4d, 0xac, 0xa9, 0xcb,
	0x4c, 0x88, 0x08, 0xfd, 0x84, 0xc7, 0xc2, 0xef, 0xe6, 0x50, 0x1e, 0xe3, 0x19, 0x80, 0x5a, 0x88,
	0x84, 0x67, 0x91, 0xbc, 0x96, 0x7e, 0x6f, 0x62, 0x4d, 0x87, 0xac, 0x85, 0xcc, 0xfb, 0x8e, 0xed,
	0x8d, 0x09, 0x05, 0x6c, 0x3f, 0xa0, 0x52, 0x99, 0x28, 0x81, 0x3e, 0x1c, 0xc5, 0x42, 0x29, 0xbe,
	0x14, 0xe5, 0x2b, 0x55, 0x4a, 0x5e, 0xc0, 0xa3, 0x1b, 0x2d, 0xd3, 0xff, 0xd2, 0x21, 0xaf, 0xc0,
	0x6b, 0x8a, 0x1e, 0x1c, 0x39, 0x83, 0xd1, 0x97, 0x48, 0x69, 0x53, 0xad, 0xea, 0xf2, 0x29, 0xd8,
	0x4b, 0x03, 0xf8, 0xd6, 0xa4, 0x37, 0x3d, 0xbe, 0x40, 0xba, 0x23, 0x03, 0x2b, 0x0a, 0xc8, 0x29,
	0x9c, 0x7c, 0x8c, 0x53, 0xbd, 0x29, 0x61, 0x32, 0x83, 0xa3, 0x6b, 0xa1, 0x7f, 0xca, 0xec, 0xbb,
	0x91, 0x60, 0xf1, 0x8d, 0xaf, 0x56, 0x22, 0x59, 0x96, 0x93, 0x5c, 0xd6, 0x42, 0x8c, 0x6c, 0x3f,
	0x56, 0x3c, 0xa9, 0x64, 0x33, 0x31, 0xf9, 0x63, 0xc1, 0xd8, 0xd0, 0xb9, 0x29, 0x95, 0x6a, 0x28,
	0xcd, 0xc0, 0xad, 0xe4, 0xab, 0x68, 0x9d, 0xd3, 0xbd, 0xa5, 0xb4, 0x42, 0x58, 0xd3, 0x11, 0xfc,
	0xb2, 0xc0, 0xa9, 0x70, 0x7c, 0x09, 0x4e, 0x52, 0x90, 0xac, 0x46, 0x39, 0xb4, 0x64, 0xcd, 0xea,
	0x2f, 0x18, 0x80, 0x13, 0xae, 0x33, 0xae, 0x23, 0x59, 0x70, 0x1c, 0xb2, 0x3a, 0x37, 0xbb, 0x85,
	0xd1, 0xdd, 0x5d, 0xb4, 0x58, 0xaf, 0xf4, 0x26, 0x3f, 0xaf, 0xcb, 0x5a, 0x08, 0x3e, 0x01, 0x5b,
	0x69, 0x99, 0x6d, 0xfc, 0x7e, 0xfe, 0xa9, 0x48, 0xf0, 0x14, 0xba, 0x51, 0xe8, 0xdb, 0xf9, 0xac,
	0x6e, 0x14, 0xce, 0xfb, 0xce, 0xc0, 0x1b, 0x93, 0x4b, 0x78, 0xd6, 0xde, 0xe3, 0xaa, 0x56, 0x88,
	0x89, 0x7b, 0x3c, 0x87, 0xe3, 0x6a, 0x89, 0xdb, 0x28, 0xcc, 0x8f, 0xd7, 0x32, 0xd2, 0xe7, 0x90,
	0x5c, 0x42, 0x70, 0xa8, 0x5b, 0xa5, 0x78, 0x06, 0xb6, 0x51, 0x7c, 0x77, 0xcd, 0x02, 0xbe, 0xf8,
	0xdd, 0x85, 0xc1, 0x87, 0xdc, 0xf2, 0xf8, 0x0e, 0xa0, 0xb9, 0x32, 0xee, 0x39, 0x79, 0xf0, 0x98,
	0xee, 0x9a, 0x95, 0x74, 0xf0, 0x0d, 0x38, 0x95, 0xdf, 0xd0, 0xa3, 0x5b, 0xfe, 0x0c, 0x46, 0x74,
	0xdb, 0x8c, 0xa4, 0x83, 0xaf, 0xc1, 0xad, 0x4d, 0x87, 0x43, 0xda, 0x76, 0x50, 0x80, 0x74, 0xc7,
	0x8f, 0xa4, 0x83, 0xef, 0x61, 0xf8, 0xcf, 0xb1, 0xb7, 0xbb, 0x9e, 0xee, 0xf7, 0x02, 0xe9, 0xe0,
	0xbc, 0xe9, 0x34, 0xe2, 0x28, 0x0c, 0xe8, 0x41, 0xb9, 0x83, 0xe7, 0xf4, 0xb0, 0x98, 0xa4, 0xf3,
	0x75, 0x90, 0xff, 0x17, 0xde, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xa0, 0x9b, 0xd4, 0x6b, 0x27,
	0x04, 0x00, 0x00,
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
	ListScenChals(ctx context.Context, in *ListScenarioChallengesReq, opts ...grpc.CallOption) (*ListScenarioChallengesResp, error)
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

func (c *daemonClient) ListScenChals(ctx context.Context, in *ListScenarioChallengesReq, opts ...grpc.CallOption) (*ListScenarioChallengesResp, error) {
	out := new(ListScenarioChallengesResp)
	err := c.cc.Invoke(ctx, "/Daemon/ListScenChals", in, out, opts...)
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
	ListScenChals(context.Context, *ListScenarioChallengesReq) (*ListScenarioChallengesResp, error)
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
func (*UnimplementedDaemonServer) ListScenChals(ctx context.Context, req *ListScenarioChallengesReq) (*ListScenarioChallengesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListScenChals not implemented")
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

func _Daemon_ListScenChals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListScenarioChallengesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).ListScenChals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Daemon/ListScenChals",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).ListScenChals(ctx, req.(*ListScenarioChallengesReq))
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
		{
			MethodName: "ListScenChals",
			Handler:    _Daemon_ListScenChals_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "daemon.proto",
}

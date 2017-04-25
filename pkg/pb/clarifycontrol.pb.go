// Code generated by protoc-gen-go.
// source: clarifycontrol.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	clarifycontrol.proto

It has these top-level messages:
	Node
	Service
	Job
	ServiceLocationReply
	NodeStatusReply
	DrainRequest
	DrainReply
	StopReply
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type ClarifyStatus int32

const (
	ClarifyStatus_UNKNOWN            ClarifyStatus = 0
	ClarifyStatus_JOB_STOPPED        ClarifyStatus = 1
	ClarifyStatus_NODE_DRAINED       ClarifyStatus = 2
	ClarifyStatus_NODE_ALLOC_MIXED   ClarifyStatus = 3
	ClarifyStatus_NODE_ALLOC_STARTED ClarifyStatus = 4
	ClarifyStatus_NODE_UNALLOCATED   ClarifyStatus = 5
)

var ClarifyStatus_name = map[int32]string{
	0: "UNKNOWN",
	1: "JOB_STOPPED",
	2: "NODE_DRAINED",
	3: "NODE_ALLOC_MIXED",
	4: "NODE_ALLOC_STARTED",
	5: "NODE_UNALLOCATED",
}
var ClarifyStatus_value = map[string]int32{
	"UNKNOWN":            0,
	"JOB_STOPPED":        1,
	"NODE_DRAINED":       2,
	"NODE_ALLOC_MIXED":   3,
	"NODE_ALLOC_STARTED": 4,
	"NODE_UNALLOCATED":   5,
}

func (x ClarifyStatus) String() string {
	return proto.EnumName(ClarifyStatus_name, int32(x))
}
func (ClarifyStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Node struct {
	Hostname string `protobuf:"bytes,1,opt,name=hostname" json:"hostname,omitempty"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Node) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

type Service struct {
	ServiceName string `protobuf:"bytes,1,opt,name=serviceName" json:"serviceName,omitempty"`
	ServiceTag  string `protobuf:"bytes,2,opt,name=serviceTag" json:"serviceTag,omitempty"`
	Healthy     bool   `protobuf:"varint,3,opt,name=healthy" json:"healthy,omitempty"`
}

func (m *Service) Reset()                    { *m = Service{} }
func (m *Service) String() string            { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()               {}
func (*Service) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Service) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *Service) GetServiceTag() string {
	if m != nil {
		return m.ServiceTag
	}
	return ""
}

func (m *Service) GetHealthy() bool {
	if m != nil {
		return m.Healthy
	}
	return false
}

type Job struct {
	JobName string `protobuf:"bytes,1,opt,name=jobName" json:"jobName,omitempty"`
}

func (m *Job) Reset()                    { *m = Job{} }
func (m *Job) String() string            { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()               {}
func (*Job) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Job) GetJobName() string {
	if m != nil {
		return m.JobName
	}
	return ""
}

type ServiceLocationReply struct {
	Locations []*ServiceLocationReply_ServiceLocation `protobuf:"bytes,1,rep,name=locations" json:"locations,omitempty"`
}

func (m *ServiceLocationReply) Reset()                    { *m = ServiceLocationReply{} }
func (m *ServiceLocationReply) String() string            { return proto.CompactTextString(m) }
func (*ServiceLocationReply) ProtoMessage()               {}
func (*ServiceLocationReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ServiceLocationReply) GetLocations() []*ServiceLocationReply_ServiceLocation {
	if m != nil {
		return m.Locations
	}
	return nil
}

type ServiceLocationReply_ServiceLocation struct {
	ServiceName string `protobuf:"bytes,1,opt,name=serviceName" json:"serviceName,omitempty"`
	ServiceHost string `protobuf:"bytes,2,opt,name=serviceHost" json:"serviceHost,omitempty"`
	ServicePort int64  `protobuf:"varint,3,opt,name=servicePort" json:"servicePort,omitempty"`
}

func (m *ServiceLocationReply_ServiceLocation) Reset()         { *m = ServiceLocationReply_ServiceLocation{} }
func (m *ServiceLocationReply_ServiceLocation) String() string { return proto.CompactTextString(m) }
func (*ServiceLocationReply_ServiceLocation) ProtoMessage()    {}
func (*ServiceLocationReply_ServiceLocation) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{3, 0}
}

func (m *ServiceLocationReply_ServiceLocation) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *ServiceLocationReply_ServiceLocation) GetServiceHost() string {
	if m != nil {
		return m.ServiceHost
	}
	return ""
}

func (m *ServiceLocationReply_ServiceLocation) GetServicePort() int64 {
	if m != nil {
		return m.ServicePort
	}
	return 0
}

type NodeStatusReply struct {
	Details []*NodeStatusReply_NodeDetails `protobuf:"bytes,1,rep,name=details" json:"details,omitempty"`
	Error   string                         `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
}

func (m *NodeStatusReply) Reset()                    { *m = NodeStatusReply{} }
func (m *NodeStatusReply) String() string            { return proto.CompactTextString(m) }
func (*NodeStatusReply) ProtoMessage()               {}
func (*NodeStatusReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *NodeStatusReply) GetDetails() []*NodeStatusReply_NodeDetails {
	if m != nil {
		return m.Details
	}
	return nil
}

func (m *NodeStatusReply) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type NodeStatusReply_NodeDetails struct {
	Node              *Node         `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"`
	Status            ClarifyStatus `protobuf:"varint,2,opt,name=status,enum=pb.ClarifyStatus" json:"status,omitempty"`
	CoordinatorLeader bool          `protobuf:"varint,3,opt,name=coordinatorLeader" json:"coordinatorLeader,omitempty"`
	ReceiverLeader    bool          `protobuf:"varint,4,opt,name=receiverLeader" json:"receiverLeader,omitempty"`
}

func (m *NodeStatusReply_NodeDetails) Reset()                    { *m = NodeStatusReply_NodeDetails{} }
func (m *NodeStatusReply_NodeDetails) String() string            { return proto.CompactTextString(m) }
func (*NodeStatusReply_NodeDetails) ProtoMessage()               {}
func (*NodeStatusReply_NodeDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 0} }

func (m *NodeStatusReply_NodeDetails) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *NodeStatusReply_NodeDetails) GetStatus() ClarifyStatus {
	if m != nil {
		return m.Status
	}
	return ClarifyStatus_UNKNOWN
}

func (m *NodeStatusReply_NodeDetails) GetCoordinatorLeader() bool {
	if m != nil {
		return m.CoordinatorLeader
	}
	return false
}

func (m *NodeStatusReply_NodeDetails) GetReceiverLeader() bool {
	if m != nil {
		return m.ReceiverLeader
	}
	return false
}

type DrainRequest struct {
	Node    *Node `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"`
	Enabled bool  `protobuf:"varint,2,opt,name=enabled" json:"enabled,omitempty"`
}

func (m *DrainRequest) Reset()                    { *m = DrainRequest{} }
func (m *DrainRequest) String() string            { return proto.CompactTextString(m) }
func (*DrainRequest) ProtoMessage()               {}
func (*DrainRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *DrainRequest) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *DrainRequest) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

type DrainReply struct {
	Node   *Node         `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"`
	Status ClarifyStatus `protobuf:"varint,2,opt,name=status,enum=pb.ClarifyStatus" json:"status,omitempty"`
}

func (m *DrainReply) Reset()                    { *m = DrainReply{} }
func (m *DrainReply) String() string            { return proto.CompactTextString(m) }
func (*DrainReply) ProtoMessage()               {}
func (*DrainReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *DrainReply) GetNode() *Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *DrainReply) GetStatus() ClarifyStatus {
	if m != nil {
		return m.Status
	}
	return ClarifyStatus_UNKNOWN
}

type StopReply struct {
	Job    *Job          `protobuf:"bytes,1,opt,name=job" json:"job,omitempty"`
	Status ClarifyStatus `protobuf:"varint,2,opt,name=status,enum=pb.ClarifyStatus" json:"status,omitempty"`
}

func (m *StopReply) Reset()                    { *m = StopReply{} }
func (m *StopReply) String() string            { return proto.CompactTextString(m) }
func (*StopReply) ProtoMessage()               {}
func (*StopReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *StopReply) GetJob() *Job {
	if m != nil {
		return m.Job
	}
	return nil
}

func (m *StopReply) GetStatus() ClarifyStatus {
	if m != nil {
		return m.Status
	}
	return ClarifyStatus_UNKNOWN
}

func init() {
	proto.RegisterType((*Node)(nil), "pb.Node")
	proto.RegisterType((*Service)(nil), "pb.Service")
	proto.RegisterType((*Job)(nil), "pb.Job")
	proto.RegisterType((*ServiceLocationReply)(nil), "pb.ServiceLocationReply")
	proto.RegisterType((*ServiceLocationReply_ServiceLocation)(nil), "pb.ServiceLocationReply.ServiceLocation")
	proto.RegisterType((*NodeStatusReply)(nil), "pb.NodeStatusReply")
	proto.RegisterType((*NodeStatusReply_NodeDetails)(nil), "pb.NodeStatusReply.NodeDetails")
	proto.RegisterType((*DrainRequest)(nil), "pb.DrainRequest")
	proto.RegisterType((*DrainReply)(nil), "pb.DrainReply")
	proto.RegisterType((*StopReply)(nil), "pb.StopReply")
	proto.RegisterEnum("pb.ClarifyStatus", ClarifyStatus_name, ClarifyStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ClarifyControl service

type ClarifyControlClient interface {
	NodeStatus(ctx context.Context, in *Job, opts ...grpc.CallOption) (*NodeStatusReply, error)
	ServiceLocation(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceLocationReply, error)
	Drain(ctx context.Context, in *DrainRequest, opts ...grpc.CallOption) (*DrainReply, error)
	Stop(ctx context.Context, in *Job, opts ...grpc.CallOption) (*StopReply, error)
	Leader(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceLocationReply, error)
}

type clarifyControlClient struct {
	cc *grpc.ClientConn
}

func NewClarifyControlClient(cc *grpc.ClientConn) ClarifyControlClient {
	return &clarifyControlClient{cc}
}

func (c *clarifyControlClient) NodeStatus(ctx context.Context, in *Job, opts ...grpc.CallOption) (*NodeStatusReply, error) {
	out := new(NodeStatusReply)
	err := grpc.Invoke(ctx, "/pb.ClarifyControl/NodeStatus", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clarifyControlClient) ServiceLocation(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceLocationReply, error) {
	out := new(ServiceLocationReply)
	err := grpc.Invoke(ctx, "/pb.ClarifyControl/ServiceLocation", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clarifyControlClient) Drain(ctx context.Context, in *DrainRequest, opts ...grpc.CallOption) (*DrainReply, error) {
	out := new(DrainReply)
	err := grpc.Invoke(ctx, "/pb.ClarifyControl/Drain", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clarifyControlClient) Stop(ctx context.Context, in *Job, opts ...grpc.CallOption) (*StopReply, error) {
	out := new(StopReply)
	err := grpc.Invoke(ctx, "/pb.ClarifyControl/Stop", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clarifyControlClient) Leader(ctx context.Context, in *Service, opts ...grpc.CallOption) (*ServiceLocationReply, error) {
	out := new(ServiceLocationReply)
	err := grpc.Invoke(ctx, "/pb.ClarifyControl/Leader", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ClarifyControl service

type ClarifyControlServer interface {
	NodeStatus(context.Context, *Job) (*NodeStatusReply, error)
	ServiceLocation(context.Context, *Service) (*ServiceLocationReply, error)
	Drain(context.Context, *DrainRequest) (*DrainReply, error)
	Stop(context.Context, *Job) (*StopReply, error)
	Leader(context.Context, *Service) (*ServiceLocationReply, error)
}

func RegisterClarifyControlServer(s *grpc.Server, srv ClarifyControlServer) {
	s.RegisterService(&_ClarifyControl_serviceDesc, srv)
}

func _ClarifyControl_NodeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClarifyControlServer).NodeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ClarifyControl/NodeStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClarifyControlServer).NodeStatus(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClarifyControl_ServiceLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Service)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClarifyControlServer).ServiceLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ClarifyControl/ServiceLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClarifyControlServer).ServiceLocation(ctx, req.(*Service))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClarifyControl_Drain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DrainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClarifyControlServer).Drain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ClarifyControl/Drain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClarifyControlServer).Drain(ctx, req.(*DrainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClarifyControl_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Job)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClarifyControlServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ClarifyControl/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClarifyControlServer).Stop(ctx, req.(*Job))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClarifyControl_Leader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Service)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClarifyControlServer).Leader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ClarifyControl/Leader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClarifyControlServer).Leader(ctx, req.(*Service))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClarifyControl_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.ClarifyControl",
	HandlerType: (*ClarifyControlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NodeStatus",
			Handler:    _ClarifyControl_NodeStatus_Handler,
		},
		{
			MethodName: "ServiceLocation",
			Handler:    _ClarifyControl_ServiceLocation_Handler,
		},
		{
			MethodName: "Drain",
			Handler:    _ClarifyControl_Drain_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _ClarifyControl_Stop_Handler,
		},
		{
			MethodName: "Leader",
			Handler:    _ClarifyControl_Leader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "clarifycontrol.proto",
}

func init() { proto.RegisterFile("clarifycontrol.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 617 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xe3, 0x38, 0xad, 0xd3, 0x71, 0x9b, 0xba, 0xdb, 0x08, 0x19, 0x83, 0x68, 0xb4, 0x07,
	0x14, 0xa0, 0xb2, 0x44, 0xe0, 0x02, 0xb7, 0xb4, 0x4e, 0x45, 0x4b, 0x70, 0x82, 0x93, 0x0a, 0x6e,
	0x95, 0xff, 0x2c, 0xd4, 0x95, 0xeb, 0x0d, 0xeb, 0x6d, 0x51, 0xcf, 0x3c, 0x01, 0xef, 0xc1, 0xfb,
	0xf0, 0x32, 0x1c, 0x90, 0xd7, 0xeb, 0xc6, 0x4d, 0x8a, 0x28, 0x12, 0xb7, 0xcc, 0x37, 0x3f, 0x4f,
	0xe6, 0xdb, 0xd9, 0x1d, 0x68, 0x87, 0x89, 0xcf, 0xe2, 0x4f, 0x57, 0x21, 0x4d, 0x39, 0xa3, 0x89,
	0x3d, 0x63, 0x94, 0x53, 0x54, 0x9f, 0x05, 0x18, 0x43, 0xc3, 0xa5, 0x11, 0x41, 0x16, 0x34, 0x4f,
	0x69, 0xc6, 0x53, 0xff, 0x9c, 0x98, 0x4a, 0x47, 0xe9, 0xae, 0x79, 0xd7, 0x31, 0x26, 0xa0, 0x4d,
	0x08, 0xbb, 0x8c, 0x43, 0x82, 0x3a, 0xa0, 0x67, 0xc5, 0x4f, 0x77, 0x4e, 0x56, 0x25, 0xf4, 0x08,
	0x40, 0x86, 0x53, 0xff, 0xb3, 0x59, 0x17, 0x40, 0x45, 0x41, 0x26, 0x68, 0xa7, 0xc4, 0x4f, 0xf8,
	0xe9, 0x95, 0xa9, 0x76, 0x94, 0x6e, 0xd3, 0x2b, 0x43, 0xbc, 0x03, 0xea, 0x11, 0x0d, 0x72, 0xe0,
	0x8c, 0x06, 0x95, 0xf2, 0x65, 0x88, 0x7f, 0x2a, 0xd0, 0x96, 0x8d, 0x0c, 0x69, 0xe8, 0xf3, 0x98,
	0xa6, 0x1e, 0x99, 0x25, 0x57, 0xe8, 0x00, 0xd6, 0x12, 0x29, 0x64, 0xa6, 0xd2, 0x51, 0xbb, 0x7a,
	0xaf, 0x6b, 0xcf, 0x02, 0xfb, 0x36, 0x78, 0x49, 0x9c, 0x7f, 0x6a, 0x7d, 0x85, 0xcd, 0x85, 0xec,
	0x1d, 0x0c, 0xcf, 0x89, 0x37, 0x34, 0xe3, 0xd2, 0x71, 0x55, 0xaa, 0x10, 0x63, 0xca, 0xb8, 0xb0,
	0xad, 0x7a, 0x55, 0x09, 0x7f, 0xaf, 0xc3, 0x66, 0x3e, 0x86, 0x09, 0xf7, 0xf9, 0x45, 0x56, 0x98,
	0x7a, 0x05, 0x5a, 0x44, 0xb8, 0x1f, 0x27, 0xa5, 0xa5, 0x9d, 0xdc, 0xd2, 0x02, 0x25, 0x62, 0xa7,
	0xc0, 0xbc, 0x92, 0x47, 0x6d, 0x58, 0x21, 0x8c, 0x51, 0x26, 0x9b, 0x29, 0x02, 0xeb, 0x87, 0x02,
	0x7a, 0x05, 0x47, 0x0f, 0xa1, 0x91, 0xd2, 0xa8, 0xf0, 0xa4, 0xf7, 0x9a, 0x65, 0x75, 0x4f, 0xa8,
	0xe8, 0x09, 0xac, 0x66, 0xe2, 0x7f, 0x44, 0x91, 0x56, 0x6f, 0x2b, 0xcf, 0xef, 0x17, 0xd7, 0x48,
	0x36, 0x20, 0x01, 0xb4, 0x0b, 0x5b, 0x21, 0xa5, 0x2c, 0x8a, 0x53, 0x9f, 0x53, 0x36, 0x24, 0x7e,
	0x44, 0x98, 0x1c, 0xee, 0x72, 0x02, 0x3d, 0x86, 0x16, 0x23, 0x21, 0x89, 0x2f, 0x49, 0x89, 0x36,
	0x04, 0xba, 0xa0, 0xe2, 0x03, 0x58, 0x77, 0x98, 0x1f, 0xa7, 0x1e, 0xf9, 0x72, 0x41, 0x32, 0xfe,
	0x97, 0x76, 0x4d, 0xd0, 0x48, 0xea, 0x07, 0x09, 0x89, 0x44, 0xbf, 0x4d, 0xaf, 0x0c, 0xf1, 0x31,
	0x80, 0xac, 0x93, 0x9f, 0xea, 0xff, 0x32, 0x8d, 0xdf, 0xc3, 0xda, 0x84, 0xd3, 0x59, 0x51, 0xf5,
	0x3e, 0xa8, 0x67, 0x34, 0x90, 0x45, 0xb5, 0xfc, 0xa3, 0x23, 0x1a, 0x78, 0xb9, 0xf6, 0x0f, 0x25,
	0x9f, 0x7e, 0x53, 0x60, 0xe3, 0x46, 0x06, 0xe9, 0xa0, 0x1d, 0xbb, 0x6f, 0xdd, 0xd1, 0x07, 0xd7,
	0xa8, 0xa1, 0x4d, 0xd0, 0x8f, 0x46, 0x7b, 0x27, 0x93, 0xe9, 0x68, 0x3c, 0x1e, 0x38, 0x86, 0x82,
	0x0c, 0x58, 0x77, 0x47, 0xce, 0xe0, 0xc4, 0xf1, 0xfa, 0x87, 0xee, 0xc0, 0x31, 0xea, 0xa8, 0x0d,
	0x86, 0x50, 0xfa, 0xc3, 0xe1, 0x68, 0xff, 0xe4, 0xdd, 0xe1, 0xc7, 0x81, 0x63, 0xa8, 0xe8, 0x1e,
	0xa0, 0x8a, 0x3a, 0x99, 0xf6, 0xbd, 0xe9, 0xc0, 0x31, 0x1a, 0xd7, 0xf4, 0xb1, 0x2b, 0x32, 0xfd,
	0x5c, 0x5d, 0xe9, 0xfd, 0x52, 0xa0, 0x25, 0xbb, 0xd8, 0x2f, 0xd6, 0x05, 0xda, 0x05, 0x98, 0xdf,
	0x3b, 0x54, 0xfa, 0xb3, 0xb6, 0x6f, 0xb9, 0x90, 0xb8, 0x86, 0x5e, 0x2f, 0xbf, 0x22, 0xbd, 0xf2,
	0x1a, 0x2d, 0xf3, 0x4f, 0x4f, 0x13, 0xd7, 0xd0, 0x33, 0x58, 0x11, 0xc3, 0x42, 0x46, 0x0e, 0x55,
	0xe7, 0x6f, 0xb5, 0x2a, 0x4a, 0x01, 0x77, 0xa0, 0x91, 0x8f, 0x60, 0xde, 0xd0, 0x86, 0xa8, 0x5c,
	0x4e, 0x05, 0xd7, 0xd0, 0x73, 0x58, 0x95, 0xb7, 0xee, 0xae, 0x1d, 0xec, 0xbd, 0x84, 0x07, 0x21,
	0x3d, 0xb7, 0xc3, 0x84, 0x50, 0x5b, 0x6e, 0x4d, 0xfb, 0x7a, 0x6d, 0x06, 0x7b, 0xdb, 0x37, 0x8f,
	0x66, 0x9c, 0x2f, 0xd2, 0xb1, 0x12, 0xac, 0x8a, 0x8d, 0xfa, 0xe2, 0x77, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xbc, 0x4a, 0xd5, 0x44, 0x69, 0x05, 0x00, 0x00,
}

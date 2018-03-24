// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bid.proto

/*
Package sonm is a generated protocol buffer package.

It is generated from these files:
	bid.proto
	bigint.proto
	capabilities.proto
	container.proto
	deal.proto
	dwh.proto
	hub.proto
	insonmnia.proto
	locator.proto
	marketplace.proto
	miner.proto
	nat.proto
	net.proto
	node.proto
	rendezvous.proto
	timestamp.proto
	volume.proto

It has these top-level messages:
	Geo
	Resources
	Slot
	Order
	BigInt
	Capabilities
	CPUDevice
	RAMDevice
	GPUDevice
	NetworkSpec
	Container
	Deal
	OrdersListRequest
	OrdersListReply
	DealsListRequest
	DealsListReply
	DealChangeRequestsReply
	DealChangeRequest
	ListReply
	HubStartTaskRequest
	HubJoinNetworkRequest
	HubStartTaskReply
	HubStatusReply
	DealRequest
	ApproveDealRequest
	GetDevicePropertiesReply
	SetDevicePropertiesRequest
	SlotsReply
	GetAllSlotsReply
	AddSlotRequest
	RemoveSlotRequest
	GetRegisteredWorkersReply
	TaskListReply
	CPUDeviceInfo
	GPUDeviceInfo
	DevicesReply
	InsertSlotRequest
	PullTaskRequest
	DealInfoReply
	Empty
	ID
	TaskID
	PingReply
	CPUUsage
	MemoryUsage
	NetworkUsage
	ResourceUsage
	InfoReply
	TaskStatusReply
	AvailableResources
	StatusMapReply
	ContainerRestartPolicy
	TaskLogsRequest
	TaskLogsChunk
	DiscoverHubRequest
	TaskResourceRequirements
	Chunk
	Progress
	AnnounceRequest
	ResolveRequest
	ResolveReply
	GetOrdersRequest
	GetOrdersReply
	GetProcessingReply
	TouchOrdersRequest
	MinerHandshakeRequest
	MinerHandshakeReply
	MinerStartRequest
	MinerStartReply
	TaskInfo
	Endpoints
	MinerStatusMapRequest
	SaveRequest
	Addr
	SocketAddr
	JoinNetworkRequest
	TaskListRequest
	DealListRequest
	DealListReply
	DealStatusReply
	ConnectRequest
	PublishRequest
	RendezvousReply
	RendezvousState
	RendezvousMeeting
	ResolveMetaReply
	Timestamp
	Volume
*/
package sonm

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

type OrderType int32

const (
	OrderType_ANY OrderType = 0
	OrderType_BID OrderType = 1
	OrderType_ASK OrderType = 2
)

var OrderType_name = map[int32]string{
	0: "ANY",
	1: "BID",
	2: "ASK",
}
var OrderType_value = map[string]int32{
	"ANY": 0,
	"BID": 1,
	"ASK": 2,
}

func (x OrderType) String() string {
	return proto.EnumName(OrderType_name, int32(x))
}
func (OrderType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Geo represent GeoIP results for node
type Geo struct {
	Country string  `protobuf:"bytes,1,opt,name=country" json:"country,omitempty"`
	City    string  `protobuf:"bytes,2,opt,name=city" json:"city,omitempty"`
	Lat     float32 `protobuf:"fixed32,3,opt,name=lat" json:"lat,omitempty"`
	Lon     float32 `protobuf:"fixed32,4,opt,name=lon" json:"lon,omitempty"`
}

func (m *Geo) Reset()                    { *m = Geo{} }
func (m *Geo) String() string            { return proto.CompactTextString(m) }
func (*Geo) ProtoMessage()               {}
func (*Geo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Geo) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *Geo) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *Geo) GetLat() float32 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Geo) GetLon() float32 {
	if m != nil {
		return m.Lon
	}
	return 0
}

type Resources struct {
	// CPU core count
	CpuCores uint64 `protobuf:"varint,1,opt,name=cpuCores" json:"cpuCores,omitempty"`
	// RAM, in bytes
	RamBytes uint64 `protobuf:"varint,2,opt,name=ramBytes" json:"ramBytes,omitempty"`
	// GPU devices count
	GpuCount GPUCount `protobuf:"varint,3,opt,name=gpuCount,enum=sonm.GPUCount" json:"gpuCount,omitempty"`
	// todo: discuss
	// storage volume, in Megabytes
	Storage uint64 `protobuf:"varint,4,opt,name=storage" json:"storage,omitempty"`
	// Inbound network traffic (the higher value), in bytes
	NetTrafficIn uint64 `protobuf:"varint,5,opt,name=netTrafficIn" json:"netTrafficIn,omitempty"`
	// Outbound network traffic (the higher value), in bytes
	NetTrafficOut uint64 `protobuf:"varint,6,opt,name=netTrafficOut" json:"netTrafficOut,omitempty"`
	// Allowed network connections
	NetworkType NetworkType `protobuf:"varint,7,opt,name=networkType,enum=sonm.NetworkType" json:"networkType,omitempty"`
	// Other properties/benchmarks. The higher means better.
	Properties map[string]float64 `protobuf:"bytes,8,rep,name=properties" json:"properties,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"fixed64,2,opt,name=value"`
}

func (m *Resources) Reset()                    { *m = Resources{} }
func (m *Resources) String() string            { return proto.CompactTextString(m) }
func (*Resources) ProtoMessage()               {}
func (*Resources) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Resources) GetCpuCores() uint64 {
	if m != nil {
		return m.CpuCores
	}
	return 0
}

func (m *Resources) GetRamBytes() uint64 {
	if m != nil {
		return m.RamBytes
	}
	return 0
}

func (m *Resources) GetGpuCount() GPUCount {
	if m != nil {
		return m.GpuCount
	}
	return GPUCount_NO_GPU
}

func (m *Resources) GetStorage() uint64 {
	if m != nil {
		return m.Storage
	}
	return 0
}

func (m *Resources) GetNetTrafficIn() uint64 {
	if m != nil {
		return m.NetTrafficIn
	}
	return 0
}

func (m *Resources) GetNetTrafficOut() uint64 {
	if m != nil {
		return m.NetTrafficOut
	}
	return 0
}

func (m *Resources) GetNetworkType() NetworkType {
	if m != nil {
		return m.NetworkType
	}
	return NetworkType_NO_NETWORK
}

func (m *Resources) GetProperties() map[string]float64 {
	if m != nil {
		return m.Properties
	}
	return nil
}

type Slot struct {
	// Buyer’s rating. Got from Buyer’s profile for BID orders rating_supplier.
	BuyerRating int64 `protobuf:"varint,1,opt,name=buyerRating" json:"buyerRating,omitempty"`
	// Supplier’s rating. Got from Supplier’s profile for ASK orders.
	SupplierRating int64 `protobuf:"varint,2,opt,name=supplierRating" json:"supplierRating,omitempty"`
	// Geo represent Worker's position
	Geo *Geo `protobuf:"bytes,3,opt,name=geo" json:"geo,omitempty"`
	// Hardware resources requirements
	Resources *Resources `protobuf:"bytes,4,opt,name=resources" json:"resources,omitempty"`
	// Duration is resource rent duration in seconds
	Duration   uint64   `protobuf:"varint,5,opt,name=duration" json:"duration,omitempty"`
	Benchmarks []uint64 `protobuf:"varint,6,rep,packed,name=benchmarks" json:"benchmarks,omitempty"`
}

func (m *Slot) Reset()                    { *m = Slot{} }
func (m *Slot) String() string            { return proto.CompactTextString(m) }
func (*Slot) ProtoMessage()               {}
func (*Slot) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Slot) GetBuyerRating() int64 {
	if m != nil {
		return m.BuyerRating
	}
	return 0
}

func (m *Slot) GetSupplierRating() int64 {
	if m != nil {
		return m.SupplierRating
	}
	return 0
}

func (m *Slot) GetGeo() *Geo {
	if m != nil {
		return m.Geo
	}
	return nil
}

func (m *Slot) GetResources() *Resources {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *Slot) GetDuration() uint64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *Slot) GetBenchmarks() []uint64 {
	if m != nil {
		return m.Benchmarks
	}
	return nil
}

type Order struct {
	// Order ID, UUIDv4
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Buyer's EtherumID
	BuyerID string `protobuf:"bytes,2,opt,name=buyerID" json:"buyerID,omitempty"`
	// Supplier's is EtherumID
	SupplierID string `protobuf:"bytes,3,opt,name=supplierID" json:"supplierID,omitempty"`
	// Order type (Bid or Ask)
	OrderType OrderType `protobuf:"varint,5,opt,name=orderType,enum=sonm.OrderType" json:"orderType,omitempty"`
	// Slot describe resource requiements
	Slot *Slot `protobuf:"bytes,6,opt,name=slot" json:"slot,omitempty"`
	// PricePerSecond specifies order price for ordered resources per second.
	PricePerSecond *BigInt `protobuf:"bytes,7,opt,name=pricePerSecond" json:"pricePerSecond,omitempty"`
}

func (m *Order) Reset()                    { *m = Order{} }
func (m *Order) String() string            { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()               {}
func (*Order) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Order) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Order) GetBuyerID() string {
	if m != nil {
		return m.BuyerID
	}
	return ""
}

func (m *Order) GetSupplierID() string {
	if m != nil {
		return m.SupplierID
	}
	return ""
}

func (m *Order) GetOrderType() OrderType {
	if m != nil {
		return m.OrderType
	}
	return OrderType_ANY
}

func (m *Order) GetSlot() *Slot {
	if m != nil {
		return m.Slot
	}
	return nil
}

func (m *Order) GetPricePerSecond() *BigInt {
	if m != nil {
		return m.PricePerSecond
	}
	return nil
}

func init() {
	proto.RegisterType((*Geo)(nil), "sonm.Geo")
	proto.RegisterType((*Resources)(nil), "sonm.Resources")
	proto.RegisterType((*Slot)(nil), "sonm.Slot")
	proto.RegisterType((*Order)(nil), "sonm.Order")
	proto.RegisterEnum("sonm.OrderType", OrderType_name, OrderType_value)
}

func init() { proto.RegisterFile("bid.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 557 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x93, 0xdd, 0x8a, 0xdb, 0x3c,
	0x10, 0x86, 0x3f, 0xff, 0xe4, 0xc7, 0xe3, 0x7c, 0x49, 0x2a, 0x7a, 0x60, 0x52, 0xd8, 0x86, 0x50,
	0x96, 0xb0, 0xd0, 0x1c, 0x78, 0x7b, 0x50, 0x0a, 0xa5, 0x34, 0x4d, 0x59, 0x42, 0x61, 0x13, 0x94,
	0x2d, 0xa5, 0x87, 0x8e, 0xad, 0x75, 0x45, 0x12, 0xc9, 0xc8, 0x72, 0x4b, 0xee, 0xa0, 0x57, 0x58,
	0x7a, 0x39, 0x45, 0xe3, 0xd8, 0xf1, 0xe6, 0x4c, 0xf3, 0xcc, 0x2b, 0x6b, 0x66, 0xde, 0x31, 0x78,
	0x5b, 0x9e, 0xcc, 0x32, 0x25, 0xb5, 0x24, 0x6e, 0x2e, 0xc5, 0x61, 0xd4, 0xdb, 0xf2, 0x94, 0x0b,
	0x5d, 0xb2, 0xd1, 0x80, 0x0b, 0x43, 0x05, 0x8f, 0x4a, 0x30, 0xf9, 0x06, 0xce, 0x1d, 0x93, 0x24,
	0x80, 0x4e, 0x2c, 0x0b, 0xa1, 0xd5, 0x31, 0xb0, 0xc6, 0xd6, 0xd4, 0xa3, 0x55, 0x48, 0x08, 0xb8,
	0x31, 0xd7, 0xc7, 0xc0, 0x46, 0x8c, 0x67, 0x32, 0x04, 0x67, 0x1f, 0xe9, 0xc0, 0x19, 0x5b, 0x53,
	0x9b, 0x9a, 0x23, 0x12, 0x29, 0x02, 0xf7, 0x44, 0xa4, 0x98, 0xfc, 0x76, 0xc0, 0xa3, 0x2c, 0x97,
	0x85, 0x8a, 0x59, 0x4e, 0x46, 0xd0, 0x8d, 0xb3, 0xe2, 0x93, 0x54, 0x2c, 0xc7, 0x07, 0x5c, 0x5a,
	0xc7, 0x26, 0xa7, 0xa2, 0xc3, 0xfc, 0xa8, 0x59, 0x8e, 0xaf, 0xb8, 0xb4, 0x8e, 0xc9, 0x0d, 0x74,
	0x53, 0xa3, 0x2b, 0x44, 0xf9, 0x5c, 0x3f, 0xec, 0xcf, 0x4c, 0x03, 0xb3, 0xbb, 0xf5, 0x57, 0xa4,
	0xb4, 0xce, 0x9b, 0x1e, 0x72, 0x2d, 0x55, 0x94, 0x32, 0xac, 0xc3, 0xa5, 0x55, 0x48, 0x26, 0xd0,
	0x13, 0x4c, 0x3f, 0xa8, 0xe8, 0xf1, 0x91, 0xc7, 0x4b, 0x11, 0xb4, 0x30, 0xfd, 0x84, 0x91, 0x57,
	0xf0, 0xff, 0x39, 0x5e, 0x15, 0x3a, 0x68, 0xa3, 0xe8, 0x29, 0x24, 0xb7, 0xe0, 0x0b, 0xa6, 0x7f,
	0x49, 0xb5, 0x7b, 0x38, 0x66, 0x2c, 0xe8, 0x60, 0x49, 0xcf, 0xca, 0x92, 0xee, 0xcf, 0x09, 0xda,
	0x54, 0x91, 0x0f, 0x00, 0x99, 0x92, 0x19, 0x53, 0x9a, 0xb3, 0x3c, 0xe8, 0x8e, 0x9d, 0xa9, 0x1f,
	0xbe, 0x2c, 0xef, 0xd4, 0x13, 0x9a, 0xad, 0x6b, 0xc5, 0x67, 0x33, 0x77, 0xda, 0xb8, 0x32, 0x7a,
	0x0f, 0x83, 0x8b, 0xb4, 0x19, 0xf8, 0x8e, 0x55, 0x66, 0x99, 0x23, 0x79, 0x0e, 0xad, 0x9f, 0xd1,
	0xbe, 0x60, 0x38, 0x43, 0x8b, 0x96, 0xc1, 0x3b, 0xfb, 0xad, 0x35, 0xf9, 0x6b, 0x81, 0xbb, 0xd9,
	0x4b, 0x4d, 0xc6, 0xe0, 0x6f, 0x8b, 0x23, 0x53, 0x34, 0xd2, 0x5c, 0xa4, 0x78, 0xd9, 0xa1, 0x4d,
	0x44, 0xae, 0xa1, 0x9f, 0x17, 0x59, 0xb6, 0xe7, 0xb5, 0xc8, 0x46, 0xd1, 0x05, 0x25, 0x2f, 0xc0,
	0x49, 0x99, 0x44, 0x4b, 0xfc, 0xd0, 0x3b, 0x59, 0xc2, 0x24, 0x35, 0x94, 0xbc, 0x06, 0x4f, 0x55,
	0x7d, 0xa1, 0x15, 0x7e, 0x38, 0xb8, 0x68, 0x97, 0x9e, 0x15, 0xc6, 0xff, 0xa4, 0x50, 0x91, 0xe6,
	0xb2, 0x72, 0xa6, 0x8e, 0xc9, 0x15, 0xc0, 0x96, 0x89, 0xf8, 0xc7, 0x21, 0x52, 0xbb, 0x3c, 0x68,
	0x8f, 0x9d, 0xa9, 0x4b, 0x1b, 0x64, 0xf2, 0xc7, 0x82, 0xd6, 0x4a, 0x25, 0x4c, 0x91, 0x3e, 0xd8,
	0x3c, 0x39, 0xcd, 0xc3, 0xe6, 0x89, 0xd9, 0x06, 0x6c, 0x6c, 0xb9, 0x38, 0xad, 0x6e, 0x15, 0x9a,
	0x6f, 0x56, 0xdd, 0x2c, 0x17, 0xd8, 0x82, 0x47, 0x1b, 0xc4, 0x94, 0x2f, 0xcd, 0x27, 0xd1, 0xe1,
	0x16, 0x3a, 0x7c, 0x2a, 0x7f, 0x55, 0x61, 0x7a, 0x56, 0x90, 0x2b, 0x70, 0xf3, 0xbd, 0x2c, 0xf7,
	0xc5, 0x0f, 0xa1, 0x54, 0x9a, 0x71, 0x53, 0xe4, 0xe4, 0x0d, 0xf4, 0x33, 0xc5, 0x63, 0xb6, 0x66,
	0x6a, 0xc3, 0x62, 0x29, 0x12, 0xdc, 0x1a, 0x3f, 0xec, 0x95, 0xca, 0x39, 0x4f, 0x97, 0x42, 0xd3,
	0x0b, 0xcd, 0xcd, 0x35, 0x78, 0xf5, 0x6b, 0xa4, 0x03, 0xce, 0xc7, 0xfb, 0xef, 0xc3, 0xff, 0xcc,
	0x61, 0xbe, 0x5c, 0x0c, 0x2d, 0x24, 0x9b, 0x2f, 0x43, 0x7b, 0xdb, 0xc6, 0xdf, 0xf8, 0xf6, 0x5f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x9f, 0xda, 0xb0, 0xe0, 0xf8, 0x03, 0x00, 0x00,
}

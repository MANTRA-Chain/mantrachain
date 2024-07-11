// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mantrachain/rewards/v1beta1/snapshot.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Snapshot
type Snapshot struct {
	Id            uint64                                      `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	PairId        uint64                                      `protobuf:"varint,2,opt,name=pairId,proto3" json:"pairId,omitempty"`
	Pools         []*SnapshotPool                             `protobuf:"bytes,3,rep,name=pools,proto3" json:"pools,omitempty"`
	PoolIdToIdx   map[uint64]uint64                           `protobuf:"bytes,4,rep,name=poolIdToIdx,proto3" json:"poolIdToIdx,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Distributed   bool                                        `protobuf:"varint,5,opt,name=distributed,proto3" json:"distributed,omitempty"`
	DistributedAt *time.Time                                  `protobuf:"bytes,6,opt,name=distributedAt,proto3,stdtime" json:"distributedAt,omitempty"`
	Remaining     github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,7,rep,name=remaining,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"remaining"`
}

func (m *Snapshot) Reset()         { *m = Snapshot{} }
func (m *Snapshot) String() string { return proto.CompactTextString(m) }
func (*Snapshot) ProtoMessage()    {}
func (*Snapshot) Descriptor() ([]byte, []int) {
	return fileDescriptor_90c05cf638b78cad, []int{0}
}
func (m *Snapshot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Snapshot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Snapshot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Snapshot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Snapshot.Merge(m, src)
}
func (m *Snapshot) XXX_Size() int {
	return m.Size()
}
func (m *Snapshot) XXX_DiscardUnknown() {
	xxx_messageInfo_Snapshot.DiscardUnknown(m)
}

var xxx_messageInfo_Snapshot proto.InternalMessageInfo

func (m *Snapshot) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Snapshot) GetPairId() uint64 {
	if m != nil {
		return m.PairId
	}
	return 0
}

func (m *Snapshot) GetPools() []*SnapshotPool {
	if m != nil {
		return m.Pools
	}
	return nil
}

func (m *Snapshot) GetPoolIdToIdx() map[uint64]uint64 {
	if m != nil {
		return m.PoolIdToIdx
	}
	return nil
}

func (m *Snapshot) GetDistributed() bool {
	if m != nil {
		return m.Distributed
	}
	return false
}

func (m *Snapshot) GetDistributedAt() *time.Time {
	if m != nil {
		return m.DistributedAt
	}
	return nil
}

func (m *Snapshot) GetRemaining() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.Remaining
	}
	return nil
}

// SnapshotCount
type SnapshotCount struct {
	PairId uint64 `protobuf:"varint,1,opt,name=pairId,proto3" json:"pairId,omitempty"`
	Count  uint64 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (m *SnapshotCount) Reset()         { *m = SnapshotCount{} }
func (m *SnapshotCount) String() string { return proto.CompactTextString(m) }
func (*SnapshotCount) ProtoMessage()    {}
func (*SnapshotCount) Descriptor() ([]byte, []int) {
	return fileDescriptor_90c05cf638b78cad, []int{1}
}
func (m *SnapshotCount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SnapshotCount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SnapshotCount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SnapshotCount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotCount.Merge(m, src)
}
func (m *SnapshotCount) XXX_Size() int {
	return m.Size()
}
func (m *SnapshotCount) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotCount.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotCount proto.InternalMessageInfo

func (m *SnapshotCount) GetPairId() uint64 {
	if m != nil {
		return m.PairId
	}
	return 0
}

func (m *SnapshotCount) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

// SnapshotStartId
type SnapshotStartId struct {
	PairId     uint64 `protobuf:"varint,1,opt,name=pairId,proto3" json:"pairId,omitempty"`
	SnapshotId uint64 `protobuf:"varint,2,opt,name=snapshotId,proto3" json:"snapshotId,omitempty"`
}

func (m *SnapshotStartId) Reset()         { *m = SnapshotStartId{} }
func (m *SnapshotStartId) String() string { return proto.CompactTextString(m) }
func (*SnapshotStartId) ProtoMessage()    {}
func (*SnapshotStartId) Descriptor() ([]byte, []int) {
	return fileDescriptor_90c05cf638b78cad, []int{2}
}
func (m *SnapshotStartId) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SnapshotStartId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SnapshotStartId.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SnapshotStartId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotStartId.Merge(m, src)
}
func (m *SnapshotStartId) XXX_Size() int {
	return m.Size()
}
func (m *SnapshotStartId) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotStartId.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotStartId proto.InternalMessageInfo

func (m *SnapshotStartId) GetPairId() uint64 {
	if m != nil {
		return m.PairId
	}
	return 0
}

func (m *SnapshotStartId) GetSnapshotId() uint64 {
	if m != nil {
		return m.SnapshotId
	}
	return 0
}

func init() {
	proto.RegisterType((*Snapshot)(nil), "mantrachain.rewards.v1beta1.Snapshot")
	proto.RegisterMapType((map[uint64]uint64)(nil), "mantrachain.rewards.v1beta1.Snapshot.PoolIdToIdxEntry")
	proto.RegisterType((*SnapshotCount)(nil), "mantrachain.rewards.v1beta1.SnapshotCount")
	proto.RegisterType((*SnapshotStartId)(nil), "mantrachain.rewards.v1beta1.SnapshotStartId")
}

func init() {
	proto.RegisterFile("mantrachain/rewards/v1beta1/snapshot.proto", fileDescriptor_90c05cf638b78cad)
}

var fileDescriptor_90c05cf638b78cad = []byte{
	// 510 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xcd, 0x26, 0x4d, 0x08, 0x1b, 0x15, 0x2a, 0x2b, 0x42, 0x56, 0x40, 0x4e, 0x94, 0x53, 0x00,
	0x75, 0xad, 0xb6, 0x02, 0x21, 0x24, 0x40, 0x49, 0x01, 0x29, 0x48, 0x40, 0xe5, 0xe6, 0x80, 0xb8,
	0xad, 0xbd, 0x8b, 0xb3, 0x6a, 0xbc, 0x63, 0x79, 0x37, 0xa5, 0xf9, 0x8b, 0x7e, 0x07, 0x5f, 0x92,
	0x63, 0x8f, 0x9c, 0x5a, 0x94, 0x5c, 0xf9, 0x08, 0x64, 0x7b, 0x9d, 0xba, 0x48, 0x54, 0x9c, 0x32,
	0x33, 0x7a, 0xef, 0x65, 0xe7, 0x3d, 0x0f, 0x7e, 0x12, 0x51, 0xa9, 0x13, 0x1a, 0x4c, 0xa9, 0x90,
	0x6e, 0xc2, 0xbf, 0xd3, 0x84, 0x29, 0xf7, 0x74, 0xcf, 0xe7, 0x9a, 0xee, 0xb9, 0x4a, 0xd2, 0x58,
	0x4d, 0x41, 0x93, 0x38, 0x01, 0x0d, 0xd6, 0xc3, 0x12, 0x96, 0x18, 0x2c, 0x31, 0xd8, 0x4e, 0x3b,
	0x84, 0x10, 0x32, 0x9c, 0x9b, 0x56, 0x39, 0xa5, 0x33, 0xb8, 0x4d, 0x3e, 0x80, 0x28, 0x02, 0x69,
	0x90, 0xdd, 0x10, 0x20, 0x9c, 0x71, 0x37, 0xeb, 0xfc, 0xf9, 0x37, 0x57, 0x8b, 0x88, 0x2b, 0x4d,
	0xa3, 0xd8, 0x00, 0x9c, 0x00, 0x54, 0x04, 0xca, 0xf5, 0xa9, 0xe2, 0x25, 0x09, 0x61, 0x04, 0xfa,
	0xbf, 0x6b, 0xb8, 0x79, 0x6c, 0x1e, 0x6c, 0xdd, 0xc3, 0x55, 0xc1, 0x6c, 0xd4, 0x43, 0x83, 0x2d,
	0xaf, 0x2a, 0x98, 0xf5, 0x00, 0x37, 0x62, 0x2a, 0x92, 0x31, 0xb3, 0xab, 0xd9, 0xcc, 0x74, 0xd6,
	0x1b, 0x5c, 0x8f, 0x01, 0x66, 0xca, 0xae, 0xf5, 0x6a, 0x83, 0xd6, 0xfe, 0x63, 0x72, 0xcb, 0x8a,
	0xa4, 0x50, 0x3f, 0x02, 0x98, 0x79, 0x39, 0xcf, 0xfa, 0x82, 0x5b, 0x69, 0x31, 0x66, 0x13, 0x18,
	0xb3, 0x33, 0x7b, 0x2b, 0x93, 0x79, 0xfe, 0x5f, 0x32, 0xe4, 0xe8, 0x9a, 0xf8, 0x4e, 0xea, 0x64,
	0xe1, 0x95, 0xa5, 0xac, 0x1e, 0x6e, 0x31, 0xa1, 0x74, 0x22, 0xfc, 0xb9, 0xe6, 0xcc, 0xae, 0xf7,
	0xd0, 0xa0, 0xe9, 0x95, 0x47, 0xd6, 0x07, 0xbc, 0x5d, 0x6a, 0x87, 0xda, 0x6e, 0xf4, 0xd0, 0xa0,
	0xb5, 0xdf, 0x21, 0xb9, 0x95, 0xa4, 0xb0, 0x92, 0x4c, 0x0a, 0x2b, 0x47, 0xcd, 0xe5, 0x65, 0x17,
	0x9d, 0x5f, 0x75, 0x91, 0x77, 0x93, 0x6a, 0x01, 0xbe, 0x9b, 0xf0, 0x88, 0x0a, 0x29, 0x64, 0x68,
	0xdf, 0xc9, 0xb6, 0x78, 0x44, 0x72, 0xc7, 0x49, 0xea, 0xf8, 0xe6, 0xf5, 0x6f, 0x79, 0x70, 0x08,
	0x42, 0x8e, 0x0e, 0x96, 0x97, 0xdd, 0xca, 0x8f, 0xab, 0xee, 0xd3, 0x50, 0xe8, 0xe9, 0xdc, 0x27,
	0x01, 0x44, 0xae, 0x49, 0x28, 0xff, 0xd9, 0x55, 0xec, 0xc4, 0xd5, 0x8b, 0x98, 0xab, 0x82, 0xa3,
	0xbc, 0xeb, 0xff, 0xe8, 0xbc, 0xc6, 0x3b, 0x7f, 0xef, 0x6f, 0xed, 0xe0, 0xda, 0x09, 0x5f, 0x98,
	0xd8, 0xd2, 0xd2, 0x6a, 0xe3, 0xfa, 0x29, 0x9d, 0xcd, 0xb9, 0x89, 0x2d, 0x6f, 0x5e, 0x56, 0x5f,
	0xa0, 0xfe, 0x2b, 0xbc, 0x5d, 0x18, 0x79, 0x08, 0x73, 0xa9, 0x4b, 0x11, 0xa3, 0x1b, 0x11, 0xb7,
	0x71, 0x3d, 0x48, 0x01, 0x85, 0x44, 0xd6, 0xf4, 0xc7, 0xf8, 0x7e, 0x41, 0x3f, 0xd6, 0x34, 0xd1,
	0x63, 0xf6, 0x4f, 0x01, 0x07, 0xe3, 0xe2, 0x10, 0x36, 0xdf, 0x4f, 0x69, 0x32, 0xfa, 0xbc, 0x5c,
	0x39, 0xe8, 0x62, 0xe5, 0xa0, 0x5f, 0x2b, 0x07, 0x9d, 0xaf, 0x9d, 0xca, 0xc5, 0xda, 0xa9, 0xfc,
	0x5c, 0x3b, 0x95, 0xaf, 0xcf, 0x4a, 0xde, 0x7c, 0x1c, 0x7e, 0x9a, 0x78, 0xc3, 0xdd, 0xf7, 0x42,
	0x52, 0x19, 0x70, 0xb7, 0x7c, 0x17, 0x67, 0x9b, 0xcb, 0xc8, 0xec, 0xf2, 0x1b, 0x59, 0x70, 0x07,
	0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x51, 0x2e, 0xf7, 0x6b, 0x9c, 0x03, 0x00, 0x00,
}

func (m *Snapshot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Snapshot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Snapshot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Remaining) > 0 {
		for iNdEx := len(m.Remaining) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Remaining[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSnapshot(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.DistributedAt != nil {
		n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(*m.DistributedAt, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.DistributedAt):])
		if err1 != nil {
			return 0, err1
		}
		i -= n1
		i = encodeVarintSnapshot(dAtA, i, uint64(n1))
		i--
		dAtA[i] = 0x32
	}
	if m.Distributed {
		i--
		if m.Distributed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if len(m.PoolIdToIdx) > 0 {
		for k := range m.PoolIdToIdx {
			v := m.PoolIdToIdx[k]
			baseI := i
			i = encodeVarintSnapshot(dAtA, i, uint64(v))
			i--
			dAtA[i] = 0x10
			i = encodeVarintSnapshot(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintSnapshot(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Pools) > 0 {
		for iNdEx := len(m.Pools) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pools[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSnapshot(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.PairId != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.PairId))
		i--
		dAtA[i] = 0x10
	}
	if m.Id != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SnapshotCount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SnapshotCount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SnapshotCount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Count != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.Count))
		i--
		dAtA[i] = 0x10
	}
	if m.PairId != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.PairId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SnapshotStartId) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SnapshotStartId) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SnapshotStartId) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SnapshotId != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.SnapshotId))
		i--
		dAtA[i] = 0x10
	}
	if m.PairId != 0 {
		i = encodeVarintSnapshot(dAtA, i, uint64(m.PairId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintSnapshot(dAtA []byte, offset int, v uint64) int {
	offset -= sovSnapshot(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Snapshot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovSnapshot(uint64(m.Id))
	}
	if m.PairId != 0 {
		n += 1 + sovSnapshot(uint64(m.PairId))
	}
	if len(m.Pools) > 0 {
		for _, e := range m.Pools {
			l = e.Size()
			n += 1 + l + sovSnapshot(uint64(l))
		}
	}
	if len(m.PoolIdToIdx) > 0 {
		for k, v := range m.PoolIdToIdx {
			_ = k
			_ = v
			mapEntrySize := 1 + sovSnapshot(uint64(k)) + 1 + sovSnapshot(uint64(v))
			n += mapEntrySize + 1 + sovSnapshot(uint64(mapEntrySize))
		}
	}
	if m.Distributed {
		n += 2
	}
	if m.DistributedAt != nil {
		l = github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.DistributedAt)
		n += 1 + l + sovSnapshot(uint64(l))
	}
	if len(m.Remaining) > 0 {
		for _, e := range m.Remaining {
			l = e.Size()
			n += 1 + l + sovSnapshot(uint64(l))
		}
	}
	return n
}

func (m *SnapshotCount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PairId != 0 {
		n += 1 + sovSnapshot(uint64(m.PairId))
	}
	if m.Count != 0 {
		n += 1 + sovSnapshot(uint64(m.Count))
	}
	return n
}

func (m *SnapshotStartId) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PairId != 0 {
		n += 1 + sovSnapshot(uint64(m.PairId))
	}
	if m.SnapshotId != 0 {
		n += 1 + sovSnapshot(uint64(m.SnapshotId))
	}
	return n
}

func sovSnapshot(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSnapshot(x uint64) (n int) {
	return sovSnapshot(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Snapshot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSnapshot
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Snapshot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Snapshot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PairId", wireType)
			}
			m.PairId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PairId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pools", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSnapshot
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSnapshot
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pools = append(m.Pools, &SnapshotPool{})
			if err := m.Pools[len(m.Pools)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolIdToIdx", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSnapshot
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSnapshot
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PoolIdToIdx == nil {
				m.PoolIdToIdx = make(map[uint64]uint64)
			}
			var mapkey uint64
			var mapvalue uint64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSnapshot
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSnapshot
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSnapshot
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipSnapshot(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthSnapshot
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.PoolIdToIdx[mapkey] = mapvalue
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Distributed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Distributed = bool(v != 0)
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DistributedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSnapshot
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSnapshot
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DistributedAt == nil {
				m.DistributedAt = new(time.Time)
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(m.DistributedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Remaining", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSnapshot
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSnapshot
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Remaining = append(m.Remaining, types.DecCoin{})
			if err := m.Remaining[len(m.Remaining)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSnapshot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSnapshot
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SnapshotCount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSnapshot
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SnapshotCount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SnapshotCount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PairId", wireType)
			}
			m.PairId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PairId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Count |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSnapshot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSnapshot
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SnapshotStartId) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSnapshot
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SnapshotStartId: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SnapshotStartId: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PairId", wireType)
			}
			m.PairId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PairId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SnapshotId", wireType)
			}
			m.SnapshotId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SnapshotId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSnapshot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSnapshot
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSnapshot(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSnapshot
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSnapshot
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSnapshot
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSnapshot
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSnapshot
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSnapshot        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSnapshot          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSnapshot = fmt.Errorf("proto: unexpected end of group")
)

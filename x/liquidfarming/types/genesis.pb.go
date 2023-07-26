// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mantrachain/liquidfarming/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
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

// GenesisState defines the liquidfarming module's genesis state.
type GenesisState struct {
	Params                     Params                       `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	LastRewardsAuctionIdRecord []LastRewardsAuctionIdRecord `protobuf:"bytes,2,rep,name=last_rewards_auction_id_record,json=lastRewardsAuctionIdRecord,proto3" json:"last_rewards_auction_id_record"`
	LiquidFarms                []LiquidFarm                 `protobuf:"bytes,3,rep,name=liquid_farms,json=liquidFarms,proto3" json:"liquid_farms"`
	RewardsAuctions            []RewardsAuction             `protobuf:"bytes,4,rep,name=rewards_auctions,json=rewardsAuctions,proto3" json:"rewards_auctions"`
	Bids                       []Bid                        `protobuf:"bytes,5,rep,name=bids,proto3" json:"bids"`
	WinningBidRecords          []WinningBidRecord           `protobuf:"bytes,6,rep,name=winning_bid_records,json=winningBidRecords,proto3" json:"winning_bid_records"`
	LastRewardsAuctionEndTime  *time.Time                   `protobuf:"bytes,7,opt,name=last_rewards_auction_end_time,json=lastRewardsAuctionEndTime,proto3,stdtime" json:"last_rewards_auction_end_time,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a8882fe876b0475, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

type LastRewardsAuctionIdRecord struct {
	PoolId    uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	AuctionId uint64 `protobuf:"varint,2,opt,name=auction_id,json=auctionId,proto3" json:"auction_id,omitempty"`
}

func (m *LastRewardsAuctionIdRecord) Reset()         { *m = LastRewardsAuctionIdRecord{} }
func (m *LastRewardsAuctionIdRecord) String() string { return proto.CompactTextString(m) }
func (*LastRewardsAuctionIdRecord) ProtoMessage()    {}
func (*LastRewardsAuctionIdRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a8882fe876b0475, []int{1}
}
func (m *LastRewardsAuctionIdRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LastRewardsAuctionIdRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LastRewardsAuctionIdRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LastRewardsAuctionIdRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LastRewardsAuctionIdRecord.Merge(m, src)
}
func (m *LastRewardsAuctionIdRecord) XXX_Size() int {
	return m.Size()
}
func (m *LastRewardsAuctionIdRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_LastRewardsAuctionIdRecord.DiscardUnknown(m)
}

var xxx_messageInfo_LastRewardsAuctionIdRecord proto.InternalMessageInfo

// WinningBidRecord defines a custom winning bid record that is required to be recorded
// in genesis state.
type WinningBidRecord struct {
	AuctionId  uint64 `protobuf:"varint,1,opt,name=auction_id,json=auctionId,proto3" json:"auction_id,omitempty"`
	WinningBid Bid    `protobuf:"bytes,2,opt,name=winning_bid,json=winningBid,proto3" json:"winning_bid"`
}

func (m *WinningBidRecord) Reset()         { *m = WinningBidRecord{} }
func (m *WinningBidRecord) String() string { return proto.CompactTextString(m) }
func (*WinningBidRecord) ProtoMessage()    {}
func (*WinningBidRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a8882fe876b0475, []int{2}
}
func (m *WinningBidRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WinningBidRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WinningBidRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WinningBidRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WinningBidRecord.Merge(m, src)
}
func (m *WinningBidRecord) XXX_Size() int {
	return m.Size()
}
func (m *WinningBidRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_WinningBidRecord.DiscardUnknown(m)
}

var xxx_messageInfo_WinningBidRecord proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GenesisState)(nil), "mantrachain.liquidfarming.v1beta1.GenesisState")
	proto.RegisterType((*LastRewardsAuctionIdRecord)(nil), "mantrachain.liquidfarming.v1beta1.LastRewardsAuctionIdRecord")
	proto.RegisterType((*WinningBidRecord)(nil), "mantrachain.liquidfarming.v1beta1.WinningBidRecord")
}

func init() {
	proto.RegisterFile("mantrachain/liquidfarming/v1beta1/genesis.proto", fileDescriptor_8a8882fe876b0475)
}

var fileDescriptor_8a8882fe876b0475 = []byte{
	// 545 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x93, 0xad, 0xeb, 0xc0, 0x9d, 0xc4, 0x30, 0x48, 0x84, 0x4a, 0x4b, 0xc7, 0x0e, 0x68,
	0x1c, 0xe6, 0xa8, 0x9b, 0xb8, 0x80, 0x90, 0x68, 0x25, 0x36, 0x4d, 0x62, 0x08, 0x85, 0x01, 0x12,
	0x07, 0x22, 0xa7, 0xf6, 0x32, 0x4b, 0x89, 0x5d, 0x6c, 0x97, 0xc2, 0x27, 0x00, 0x71, 0xda, 0x47,
	0xd8, 0xc7, 0xe0, 0x23, 0xec, 0xb8, 0x23, 0x27, 0x40, 0xed, 0x85, 0x8f, 0x81, 0xe2, 0xa4, 0x6b,
	0x13, 0x98, 0x9a, 0xdd, 0xe2, 0x97, 0xf7, 0xff, 0x3d, 0xbf, 0xf7, 0xfe, 0x32, 0xf0, 0x12, 0xcc,
	0xb5, 0xc4, 0xbd, 0x63, 0xcc, 0xb8, 0x17, 0xb3, 0x0f, 0x03, 0x46, 0x8e, 0xb0, 0x4c, 0x18, 0x8f,
	0xbc, 0x8f, 0xed, 0x90, 0x6a, 0xdc, 0xf6, 0x22, 0xca, 0xa9, 0x62, 0x0a, 0xf5, 0xa5, 0xd0, 0x02,
	0xde, 0x9b, 0x11, 0xa0, 0x82, 0x00, 0xe5, 0x82, 0xe6, 0xed, 0x48, 0x44, 0xc2, 0x64, 0x7b, 0xe9,
	0x57, 0x26, 0x6c, 0xb6, 0x22, 0x21, 0xa2, 0x98, 0x7a, 0xe6, 0x14, 0x0e, 0x8e, 0x3c, 0xcd, 0x12,
	0xaa, 0x34, 0x4e, 0xfa, 0x79, 0xc2, 0xc3, 0xf9, 0x57, 0x29, 0xd6, 0xcb, 0x64, 0x68, 0xbe, 0xac,
	0x8f, 0x25, 0x4e, 0xf2, 0x06, 0x36, 0xbe, 0x2f, 0x81, 0x95, 0xbd, 0xac, 0xa5, 0x57, 0x1a, 0x6b,
	0x0a, 0xf7, 0x40, 0x3d, 0x4b, 0x70, 0xec, 0x75, 0x7b, 0xb3, 0xb1, 0xfd, 0x00, 0xcd, 0x6d, 0x11,
	0xbd, 0x34, 0x82, 0x6e, 0xed, 0xec, 0x67, 0xcb, 0xf2, 0x73, 0x39, 0xfc, 0x62, 0x03, 0x37, 0xc6,
	0x4a, 0x07, 0x92, 0x0e, 0xb1, 0x24, 0x2a, 0xc0, 0x83, 0x9e, 0x66, 0x82, 0x07, 0x8c, 0x04, 0x92,
	0xf6, 0x84, 0x24, 0xce, 0xc2, 0xfa, 0xe2, 0x66, 0x63, 0xfb, 0x49, 0x85, 0x0a, 0xcf, 0xb1, 0xd2,
	0x7e, 0xc6, 0xe9, 0x64, 0x98, 0x7d, 0xe2, 0x1b, 0x48, 0x5e, 0xb5, 0x19, 0x5f, 0x9a, 0x01, 0xdf,
	0x80, 0x95, 0x8c, 0x1a, 0xa4, 0x58, 0xe5, 0x2c, 0x9a, 0xb2, 0x5b, 0x55, 0xca, 0x9a, 0xe8, 0x2e,
	0x96, 0x49, 0x5e, 0xa6, 0x11, 0x5f, 0x44, 0x14, 0x0c, 0xc1, 0x6a, 0xa9, 0x37, 0xe5, 0xd4, 0x0c,
	0xbb, 0x5d, 0x81, 0x5d, 0xbc, 0x6c, 0xce, 0xbf, 0x21, 0x0b, 0x51, 0x05, 0x9f, 0x82, 0x5a, 0xc8,
	0x88, 0x72, 0x96, 0x0c, 0xf7, 0x7e, 0x05, 0x6e, 0x97, 0x4d, 0x66, 0x62, 0x94, 0x90, 0x81, 0x5b,
	0x43, 0xc6, 0x39, 0xe3, 0x51, 0x10, 0x5e, 0x8c, 0x5e, 0x39, 0x75, 0x03, 0xdc, 0xa9, 0x00, 0x7c,
	0x9b, 0xa9, 0xbb, 0xac, 0x38, 0xf1, 0x9b, 0xc3, 0x52, 0x3c, 0x1d, 0xc8, 0xda, 0x7f, 0x37, 0x4e,
	0x39, 0x09, 0x52, 0x7f, 0x3b, 0xcb, 0xc6, 0x52, 0x4d, 0x94, 0x99, 0x1f, 0x4d, 0xcc, 0x8f, 0x0e,
	0x27, 0xe6, 0xef, 0xd6, 0x4e, 0x7e, 0xb5, 0x6c, 0xff, 0xee, 0xbf, 0xdb, 0x7c, 0xc6, 0x49, 0x9a,
	0xf5, 0xe8, 0xda, 0xd7, 0xd3, 0x96, 0xf5, 0xe7, 0xb4, 0x65, 0x6d, 0xbc, 0x07, 0xcd, 0xcb, 0x6d,
	0x01, 0xef, 0x80, 0xe5, 0xbe, 0x10, 0x71, 0xc0, 0x88, 0x31, 0x72, 0xcd, 0xaf, 0xa7, 0xc7, 0x7d,
	0x02, 0xd7, 0x00, 0x98, 0x3a, 0xd1, 0x59, 0x30, 0xff, 0xae, 0xe3, 0x89, 0x7a, 0x86, 0xff, 0xcd,
	0x06, 0xab, 0xe5, 0xde, 0x4b, 0x6a, 0xbb, 0xa4, 0x86, 0x07, 0xa0, 0x31, 0x33, 0x6c, 0x43, 0xbf,
	0xea, 0xd6, 0xc0, 0x74, 0xae, 0xd3, 0xcb, 0x74, 0x5f, 0x9f, 0x8d, 0x5c, 0xfb, 0x7c, 0xe4, 0xda,
	0xbf, 0x47, 0xae, 0x7d, 0x32, 0x76, 0xad, 0xf3, 0xb1, 0x6b, 0xfd, 0x18, 0xbb, 0xd6, 0xbb, 0xc7,
	0x11, 0xd3, 0xc7, 0x83, 0x10, 0xf5, 0x44, 0xe2, 0x1d, 0x74, 0x5e, 0x1c, 0xfa, 0x9d, 0xad, 0x5d,
	0xc6, 0x31, 0xef, 0xd1, 0xc2, 0x6b, 0xf6, 0xa9, 0xf4, 0x1a, 0xe8, 0xcf, 0x7d, 0xaa, 0xc2, 0xba,
	0x59, 0xc1, 0xce, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcb, 0x7b, 0x60, 0xc4, 0xf9, 0x04, 0x00,
	0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastRewardsAuctionEndTime != nil {
		n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(*m.LastRewardsAuctionEndTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.LastRewardsAuctionEndTime):])
		if err1 != nil {
			return 0, err1
		}
		i -= n1
		i = encodeVarintGenesis(dAtA, i, uint64(n1))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.WinningBidRecords) > 0 {
		for iNdEx := len(m.WinningBidRecords) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.WinningBidRecords[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Bids) > 0 {
		for iNdEx := len(m.Bids) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Bids[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.RewardsAuctions) > 0 {
		for iNdEx := len(m.RewardsAuctions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RewardsAuctions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.LiquidFarms) > 0 {
		for iNdEx := len(m.LiquidFarms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LiquidFarms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.LastRewardsAuctionIdRecord) > 0 {
		for iNdEx := len(m.LastRewardsAuctionIdRecord) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LastRewardsAuctionIdRecord[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *LastRewardsAuctionIdRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LastRewardsAuctionIdRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LastRewardsAuctionIdRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.AuctionId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.AuctionId))
		i--
		dAtA[i] = 0x10
	}
	if m.PoolId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *WinningBidRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WinningBidRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WinningBidRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.WinningBid.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.AuctionId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.AuctionId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.LastRewardsAuctionIdRecord) > 0 {
		for _, e := range m.LastRewardsAuctionIdRecord {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.LiquidFarms) > 0 {
		for _, e := range m.LiquidFarms {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.RewardsAuctions) > 0 {
		for _, e := range m.RewardsAuctions {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Bids) > 0 {
		for _, e := range m.Bids {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.WinningBidRecords) > 0 {
		for _, e := range m.WinningBidRecords {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.LastRewardsAuctionEndTime != nil {
		l = github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.LastRewardsAuctionEndTime)
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func (m *LastRewardsAuctionIdRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovGenesis(uint64(m.PoolId))
	}
	if m.AuctionId != 0 {
		n += 1 + sovGenesis(uint64(m.AuctionId))
	}
	return n
}

func (m *WinningBidRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AuctionId != 0 {
		n += 1 + sovGenesis(uint64(m.AuctionId))
	}
	l = m.WinningBid.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastRewardsAuctionIdRecord", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LastRewardsAuctionIdRecord = append(m.LastRewardsAuctionIdRecord, LastRewardsAuctionIdRecord{})
			if err := m.LastRewardsAuctionIdRecord[len(m.LastRewardsAuctionIdRecord)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidFarms", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LiquidFarms = append(m.LiquidFarms, LiquidFarm{})
			if err := m.LiquidFarms[len(m.LiquidFarms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardsAuctions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardsAuctions = append(m.RewardsAuctions, RewardsAuction{})
			if err := m.RewardsAuctions[len(m.RewardsAuctions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bids", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bids = append(m.Bids, Bid{})
			if err := m.Bids[len(m.Bids)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WinningBidRecords", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.WinningBidRecords = append(m.WinningBidRecords, WinningBidRecord{})
			if err := m.WinningBidRecords[len(m.WinningBidRecords)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastRewardsAuctionEndTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.LastRewardsAuctionEndTime == nil {
				m.LastRewardsAuctionEndTime = new(time.Time)
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(m.LastRewardsAuctionEndTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *LastRewardsAuctionIdRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: LastRewardsAuctionIdRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LastRewardsAuctionIdRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionId", wireType)
			}
			m.AuctionId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuctionId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *WinningBidRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: WinningBidRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WinningBidRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionId", wireType)
			}
			m.AuctionId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuctionId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WinningBid", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.WinningBid.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
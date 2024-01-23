// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mantrachain/rewards/v1beta1/common.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

type ProviderPair struct {
	PairId                uint64                                      `protobuf:"varint,1,opt,name=pairId,proto3" json:"pairId,omitempty"`
	LastClaimedSnapshotId uint64                                      `protobuf:"varint,2,opt,name=lastClaimedSnapshotId,proto3" json:"lastClaimedSnapshotId,omitempty"`
	OwedRewards           github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,3,rep,name=owedRewards,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"owedRewards" yaml:"owed_rewards"`
	LastDepositTime       *time.Time                                  `protobuf:"bytes,4,opt,name=lastDepositTime,proto3,stdtime" json:"lastDepositTime,omitempty" yaml:"last_deposit_time"`
}

func (m *ProviderPair) Reset()         { *m = ProviderPair{} }
func (m *ProviderPair) String() string { return proto.CompactTextString(m) }
func (*ProviderPair) ProtoMessage()    {}
func (*ProviderPair) Descriptor() ([]byte, []int) {
	return fileDescriptor_274771e038f60916, []int{0}
}
func (m *ProviderPair) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProviderPair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProviderPair.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProviderPair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProviderPair.Merge(m, src)
}
func (m *ProviderPair) XXX_Size() int {
	return m.Size()
}
func (m *ProviderPair) XXX_DiscardUnknown() {
	xxx_messageInfo_ProviderPair.DiscardUnknown(m)
}

var xxx_messageInfo_ProviderPair proto.InternalMessageInfo

type SnapshotPool struct {
	PoolId          uint64                                      `protobuf:"varint,1,opt,name=poolId,proto3" json:"poolId,omitempty"`
	RewardsPerToken github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,2,rep,name=rewardsPerToken,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"rewardsPerToken" yaml:"rewards_per_token"`
}

func (m *SnapshotPool) Reset()         { *m = SnapshotPool{} }
func (m *SnapshotPool) String() string { return proto.CompactTextString(m) }
func (*SnapshotPool) ProtoMessage()    {}
func (*SnapshotPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_274771e038f60916, []int{1}
}
func (m *SnapshotPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SnapshotPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SnapshotPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SnapshotPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotPool.Merge(m, src)
}
func (m *SnapshotPool) XXX_Size() int {
	return m.Size()
}
func (m *SnapshotPool) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotPool.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotPool proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ProviderPair)(nil), "mantrachain.rewards.v1beta1.ProviderPair")
	proto.RegisterType((*SnapshotPool)(nil), "mantrachain.rewards.v1beta1.SnapshotPool")
}

func init() {
	proto.RegisterFile("mantrachain/rewards/v1beta1/common.proto", fileDescriptor_274771e038f60916)
}

var fileDescriptor_274771e038f60916 = []byte{
	// 486 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x4f, 0x6b, 0xd4, 0x40,
	0x18, 0xc6, 0x33, 0xed, 0xd2, 0x43, 0xb6, 0x50, 0x88, 0x7f, 0x88, 0xab, 0x24, 0x4b, 0xf0, 0xb0,
	0x20, 0xcd, 0xd0, 0xaa, 0x97, 0xde, 0xba, 0x2d, 0x42, 0x05, 0x75, 0x49, 0xf7, 0xe4, 0x25, 0x4c,
	0x92, 0x31, 0x3b, 0x36, 0x93, 0x37, 0xcc, 0x4c, 0x5b, 0xfb, 0x0d, 0xc4, 0x53, 0x6f, 0x5e, 0x3d,
	0x8a, 0x1f, 0x44, 0xf6, 0xd8, 0xa3, 0x20, 0x6c, 0x65, 0xf7, 0x1b, 0xf4, 0x13, 0x48, 0x32, 0x13,
	0x09, 0xc5, 0x83, 0x9e, 0x92, 0x97, 0x79, 0xde, 0xe7, 0x7d, 0xde, 0xdf, 0x30, 0xf6, 0x88, 0x93,
	0x52, 0x09, 0x92, 0xce, 0x08, 0x2b, 0xb1, 0xa0, 0xe7, 0x44, 0x64, 0x12, 0x9f, 0xed, 0x24, 0x54,
	0x91, 0x1d, 0x9c, 0x02, 0xe7, 0x50, 0x86, 0x95, 0x00, 0x05, 0xce, 0xc3, 0x8e, 0x32, 0x34, 0xca,
	0xd0, 0x28, 0x07, 0x77, 0x73, 0xc8, 0xa1, 0xd1, 0xe1, 0xfa, 0x4f, 0xb7, 0x0c, 0xfc, 0x1c, 0x20,
	0x2f, 0x28, 0x6e, 0xaa, 0xe4, 0xf4, 0x1d, 0x56, 0x8c, 0x53, 0xa9, 0x08, 0xaf, 0x8c, 0xe0, 0x41,
	0x0a, 0x92, 0x83, 0x8c, 0x75, 0xa7, 0x2e, 0xcc, 0x91, 0xa7, 0x2b, 0x9c, 0x10, 0x49, 0x3b, 0x81,
	0x98, 0x89, 0x13, 0xfc, 0x5c, 0xb3, 0x37, 0x27, 0x02, 0xce, 0x58, 0x46, 0xc5, 0x84, 0x30, 0xe1,
	0xdc, 0xb7, 0x37, 0x2a, 0xc2, 0xc4, 0x51, 0xe6, 0xa2, 0x21, 0x1a, 0xf5, 0x22, 0x53, 0x39, 0xcf,
	0xec, 0x7b, 0x05, 0x91, 0xea, 0xa0, 0x20, 0x8c, 0xd3, 0xec, 0xb8, 0x24, 0x95, 0x9c, 0x81, 0x3a,
	0xca, 0xdc, 0xb5, 0x46, 0xf6, 0xf7, 0x43, 0xe7, 0x13, 0xb2, 0xfb, 0x70, 0x4e, 0xb3, 0x48, 0x2f,
	0xea, 0xae, 0x0f, 0xd7, 0x47, 0xfd, 0xdd, 0x47, 0xa1, 0xc9, 0x58, 0xa7, 0x6a, 0x97, 0x0f, 0x0f,
	0x69, 0x7a, 0x00, 0xac, 0x1c, 0xbf, 0x9c, 0x2f, 0x7c, 0xeb, 0x66, 0xe1, 0xdf, 0xb9, 0x20, 0xbc,
	0xd8, 0x0b, 0xea, 0xf6, 0xd8, 0x80, 0x0a, 0xbe, 0x5d, 0xfb, 0x4f, 0x72, 0xa6, 0x66, 0xa7, 0x49,
	0x98, 0x02, 0x37, 0xab, 0x9a, 0xcf, 0xb6, 0xcc, 0x4e, 0xb0, 0xba, 0xa8, 0xa8, 0x6c, 0xad, 0x64,
	0xd4, 0x1d, 0xee, 0xbc, 0xb7, 0xb7, 0xea, 0x94, 0x87, 0xb4, 0x02, 0xc9, 0xd4, 0x94, 0x71, 0xea,
	0xf6, 0x86, 0x68, 0xd4, 0xdf, 0x1d, 0x84, 0x9a, 0x70, 0xd8, 0x12, 0x0e, 0xa7, 0x2d, 0xe1, 0xf1,
	0xe3, 0xf9, 0xc2, 0x47, 0x37, 0x0b, 0xdf, 0xd5, 0x69, 0x6a, 0x83, 0x38, 0xd3, 0x0e, 0x71, 0x7d,
	0x0f, 0xc1, 0xe5, 0xb5, 0x8f, 0xa2, 0xdb, 0xc6, 0x7b, 0xbd, 0x8f, 0x5f, 0x7c, 0x2b, 0xf8, 0x8e,
	0xec, 0xcd, 0x96, 0xc6, 0x04, 0xa0, 0x68, 0xe8, 0x02, 0x14, 0x1d, 0xba, 0x4d, 0xe5, 0x7c, 0x46,
	0xf6, 0x96, 0xd9, 0x71, 0x42, 0xc5, 0x14, 0x4e, 0x68, 0xe9, 0xae, 0xfd, 0x03, 0xab, 0x37, 0x86,
	0x95, 0x49, 0x67, 0x2c, 0xe2, 0x8a, 0x8a, 0x58, 0xd5, 0x26, 0xff, 0x0d, 0xec, 0x76, 0x0a, 0xbd,
	0xc8, 0xf8, 0xf8, 0xeb, 0xd2, 0x43, 0xf3, 0xa5, 0x87, 0xae, 0x96, 0x1e, 0xfa, 0xb5, 0xf4, 0xd0,
	0xe5, 0xca, 0xb3, 0xae, 0x56, 0x9e, 0xf5, 0x63, 0xe5, 0x59, 0x6f, 0x9f, 0x77, 0x26, 0xbc, 0xda,
	0x7f, 0x3d, 0x8d, 0xf6, 0xb7, 0x5f, 0xb0, 0x92, 0x94, 0x29, 0xc5, 0xdd, 0x77, 0xf1, 0xe1, 0xcf,
	0xcb, 0x68, 0x86, 0x26, 0x1b, 0x0d, 0xee, 0xa7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xd5,
	0xb4, 0x8f, 0x3d, 0x03, 0x00, 0x00,
}

func (this *ProviderPair) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProviderPair)
	if !ok {
		that2, ok := that.(ProviderPair)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.PairId != that1.PairId {
		return false
	}
	if this.LastClaimedSnapshotId != that1.LastClaimedSnapshotId {
		return false
	}
	if len(this.OwedRewards) != len(that1.OwedRewards) {
		return false
	}
	for i := range this.OwedRewards {
		if !this.OwedRewards[i].Equal(&that1.OwedRewards[i]) {
			return false
		}
	}
	if that1.LastDepositTime == nil {
		if this.LastDepositTime != nil {
			return false
		}
	} else if !this.LastDepositTime.Equal(*that1.LastDepositTime) {
		return false
	}
	return true
}
func (this *SnapshotPool) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SnapshotPool)
	if !ok {
		that2, ok := that.(SnapshotPool)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.PoolId != that1.PoolId {
		return false
	}
	if len(this.RewardsPerToken) != len(that1.RewardsPerToken) {
		return false
	}
	for i := range this.RewardsPerToken {
		if !this.RewardsPerToken[i].Equal(&that1.RewardsPerToken[i]) {
			return false
		}
	}
	return true
}
func (m *ProviderPair) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProviderPair) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProviderPair) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastDepositTime != nil {
		n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(*m.LastDepositTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.LastDepositTime):])
		if err1 != nil {
			return 0, err1
		}
		i -= n1
		i = encodeVarintCommon(dAtA, i, uint64(n1))
		i--
		dAtA[i] = 0x22
	}
	if len(m.OwedRewards) > 0 {
		for iNdEx := len(m.OwedRewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.OwedRewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCommon(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.LastClaimedSnapshotId != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.LastClaimedSnapshotId))
		i--
		dAtA[i] = 0x10
	}
	if m.PairId != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.PairId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SnapshotPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SnapshotPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SnapshotPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RewardsPerToken) > 0 {
		for iNdEx := len(m.RewardsPerToken) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RewardsPerToken[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCommon(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.PoolId != 0 {
		i = encodeVarintCommon(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommon(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommon(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProviderPair) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PairId != 0 {
		n += 1 + sovCommon(uint64(m.PairId))
	}
	if m.LastClaimedSnapshotId != 0 {
		n += 1 + sovCommon(uint64(m.LastClaimedSnapshotId))
	}
	if len(m.OwedRewards) > 0 {
		for _, e := range m.OwedRewards {
			l = e.Size()
			n += 1 + l + sovCommon(uint64(l))
		}
	}
	if m.LastDepositTime != nil {
		l = github_com_cosmos_gogoproto_types.SizeOfStdTime(*m.LastDepositTime)
		n += 1 + l + sovCommon(uint64(l))
	}
	return n
}

func (m *SnapshotPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovCommon(uint64(m.PoolId))
	}
	if len(m.RewardsPerToken) > 0 {
		for _, e := range m.RewardsPerToken {
			l = e.Size()
			n += 1 + l + sovCommon(uint64(l))
		}
	}
	return n
}

func sovCommon(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommon(x uint64) (n int) {
	return sovCommon(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProviderPair) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
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
			return fmt.Errorf("proto: ProviderPair: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProviderPair: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PairId", wireType)
			}
			m.PairId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return fmt.Errorf("proto: wrong wireType = %d for field LastClaimedSnapshotId", wireType)
			}
			m.LastClaimedSnapshotId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastClaimedSnapshotId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OwedRewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OwedRewards = append(m.OwedRewards, types.DecCoin{})
			if err := m.OwedRewards[len(m.OwedRewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastDepositTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.LastDepositTime == nil {
				m.LastDepositTime = new(time.Time)
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(m.LastDepositTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
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
func (m *SnapshotPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
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
			return fmt.Errorf("proto: SnapshotPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SnapshotPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardsPerToken", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardsPerToken = append(m.RewardsPerToken, types.DecCoin{})
			if err := m.RewardsPerToken[len(m.RewardsPerToken)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
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
func skipCommon(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommon
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
					return 0, ErrIntOverflowCommon
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
					return 0, ErrIntOverflowCommon
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
				return 0, ErrInvalidLengthCommon
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommon
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommon
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommon        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommon          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommon = fmt.Errorf("proto: unexpected end of group")
)

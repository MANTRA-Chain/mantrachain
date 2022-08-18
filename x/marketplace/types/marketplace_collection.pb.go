// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: marketplace/v1/marketplace_collection.proto

package types

import (
	bytes "bytes"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MarketplaceCollection struct {
	Index                                   []byte                                        `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty" yaml:"index"`
	MarketplaceIndex                        []byte                                        `protobuf:"bytes,2,opt,name=marketplace_index,json=marketplaceIndex,proto3" json:"marketplace_index,omitempty" yaml:"marketplace_index"`
	CollectionCreator                       string                                        `protobuf:"bytes,3,opt,name=collection_creator,json=collectionCreator,proto3" json:"collection_creator,omitempty" yaml:"collection_creator"`
	CollectionId                            string                                        `protobuf:"bytes,4,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty" yaml:"collection_id"`
	InitiallyNftCollectionOwnerNftsForSale  bool                                          `protobuf:"varint,5,opt,name=initially_nft_collection_owner_nfts_for_sale,json=initiallyNftCollectionOwnerNftsForSale,proto3" json:"initially_nft_collection_owner_nfts_for_sale,omitempty" yaml:"initially_nft_collection_owner_nfts_for_sale"`
	InitiallyNftCollectionOwnerNftsMinPrice *github_com_cosmos_cosmos_sdk_types.Coin      `protobuf:"bytes,6,opt,name=initially_nft_collection_owner_nfts_min_price,json=initiallyNftCollectionOwnerNftsMinPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin" json:"initially_nft_collection_owner_nfts_min_price,omitempty" yaml:"initially_nft_collection_owner_nfts_min_price"`
	NftsEarningsOnSale                      []*MarketplaceEarning                         `protobuf:"bytes,7,rep,name=nfts_earnings_on_sale,json=nftsEarningsOnSale,proto3" json:"nfts_earnings_on_sale,omitempty" yaml:"nfts_earnings_on_sale"`
	NftsEarningsOnYieldReward               []*MarketplaceEarning                         `protobuf:"bytes,8,rep,name=nfts_earnings_on_yield_reward,json=nftsEarningsOnYieldReward,proto3" json:"nfts_earnings_on_yield_reward,omitempty" yaml:"nfts_earnings_on_yield_reward"`
	InitiallyNftsVaultLockPercentage        *github_com_cosmos_cosmos_sdk_types.Int       `protobuf:"bytes,9,opt,name=initially_nfts_vault_lock_percentage,json=initiallyNftsVaultLockPercentage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"initially_nfts_vault_lock_percentage,omitempty" yaml:"initially_nfts_vault_lock_percentage"`
	Creator                                 github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,10,opt,name=creator,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"creator,omitempty" yaml:"creator"`
}

func (m *MarketplaceCollection) Reset()         { *m = MarketplaceCollection{} }
func (m *MarketplaceCollection) String() string { return proto.CompactTextString(m) }
func (*MarketplaceCollection) ProtoMessage()    {}
func (*MarketplaceCollection) Descriptor() ([]byte, []int) {
	return fileDescriptor_d566027cad7fc202, []int{0}
}
func (m *MarketplaceCollection) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketplaceCollection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketplaceCollection.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketplaceCollection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketplaceCollection.Merge(m, src)
}
func (m *MarketplaceCollection) XXX_Size() int {
	return m.Size()
}
func (m *MarketplaceCollection) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketplaceCollection.DiscardUnknown(m)
}

var xxx_messageInfo_MarketplaceCollection proto.InternalMessageInfo

func (m *MarketplaceCollection) GetIndex() []byte {
	if m != nil {
		return m.Index
	}
	return nil
}

func (m *MarketplaceCollection) GetMarketplaceIndex() []byte {
	if m != nil {
		return m.MarketplaceIndex
	}
	return nil
}

func (m *MarketplaceCollection) GetCollectionCreator() string {
	if m != nil {
		return m.CollectionCreator
	}
	return ""
}

func (m *MarketplaceCollection) GetCollectionId() string {
	if m != nil {
		return m.CollectionId
	}
	return ""
}

func (m *MarketplaceCollection) GetInitiallyNftCollectionOwnerNftsForSale() bool {
	if m != nil {
		return m.InitiallyNftCollectionOwnerNftsForSale
	}
	return false
}

func (m *MarketplaceCollection) GetNftsEarningsOnSale() []*MarketplaceEarning {
	if m != nil {
		return m.NftsEarningsOnSale
	}
	return nil
}

func (m *MarketplaceCollection) GetNftsEarningsOnYieldReward() []*MarketplaceEarning {
	if m != nil {
		return m.NftsEarningsOnYieldReward
	}
	return nil
}

func (m *MarketplaceCollection) GetCreator() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Creator
	}
	return nil
}

func init() {
	proto.RegisterType((*MarketplaceCollection)(nil), "LimeChain.mantrachain.marketplace.v1.MarketplaceCollection")
}

func init() {
	proto.RegisterFile("marketplace/v1/marketplace_collection.proto", fileDescriptor_d566027cad7fc202)
}

var fileDescriptor_d566027cad7fc202 = []byte{
	// 679 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xcb, 0x6a, 0xdb, 0x4c,
	0x14, 0xce, 0xfc, 0xf9, 0x73, 0x53, 0xd2, 0x92, 0x88, 0x04, 0x94, 0x34, 0x91, 0x84, 0x08, 0x89,
	0x21, 0x8d, 0x44, 0x2e, 0xd0, 0x52, 0xe8, 0x22, 0x36, 0x2d, 0x18, 0x72, 0x43, 0x85, 0x42, 0x0b,
	0x41, 0x8c, 0x47, 0x63, 0x67, 0xb0, 0x34, 0x63, 0x66, 0x26, 0x17, 0x3f, 0x44, 0xa1, 0x9b, 0xbe,
	0x41, 0x17, 0xdd, 0xf4, 0x3d, 0xb2, 0xcc, 0xb2, 0x64, 0x21, 0x8a, 0xfd, 0x06, 0x5a, 0x76, 0x55,
	0x74, 0xb1, 0xad, 0xd4, 0x09, 0x35, 0x74, 0xa5, 0xd1, 0x99, 0xef, 0x7c, 0xe7, 0x9b, 0xf9, 0xce,
	0x19, 0x65, 0x2b, 0x84, 0xbc, 0x89, 0x65, 0x2b, 0x80, 0x08, 0x3b, 0x97, 0x3b, 0x4e, 0xe1, 0xd7,
	0x43, 0x2c, 0x08, 0x30, 0x92, 0x84, 0x51, 0xbb, 0xc5, 0x99, 0x64, 0xea, 0xfa, 0x21, 0x09, 0x71,
	0xe5, 0x1c, 0x12, 0x6a, 0x87, 0x90, 0x4a, 0x0e, 0x51, 0xbe, 0xee, 0xe7, 0xd8, 0x97, 0x3b, 0x2b,
	0x3a, 0x62, 0x22, 0x64, 0xc2, 0xa9, 0x41, 0x91, 0x50, 0xd6, 0xb0, 0x84, 0x3b, 0x0e, 0x62, 0x24,
	0x67, 0x59, 0x59, 0x6c, 0xb0, 0x06, 0x4b, 0x97, 0x4e, 0xb2, 0xca, 0xa3, 0xcf, 0xfe, 0x10, 0x82,
	0x58, 0x18, 0xf6, 0x0a, 0x5b, 0x77, 0x33, 0xca, 0xd2, 0xd1, 0x60, 0xbf, 0xd2, 0x17, 0xa6, 0x6e,
	0x28, 0x13, 0x84, 0xfa, 0xf8, 0x5a, 0x03, 0x26, 0x28, 0xcd, 0x95, 0xe7, 0xe3, 0xc8, 0x98, 0x6b,
	0xc3, 0x30, 0x78, 0x65, 0xa5, 0x61, 0xcb, 0xcd, 0xb6, 0xd5, 0xaa, 0xb2, 0x50, 0x3c, 0x5a, 0x96,
	0xf3, 0x5f, 0x9a, 0xb3, 0x1a, 0x47, 0x86, 0x96, 0xe5, 0x0c, 0x41, 0x2c, 0x77, 0xbe, 0x10, 0xab,
	0xa6, 0x54, 0x87, 0x8a, 0x3a, 0xb8, 0x19, 0x0f, 0x71, 0x0c, 0x25, 0xe3, 0xda, 0xb8, 0x09, 0x4a,
	0x33, 0xe5, 0xb5, 0x38, 0x32, 0x96, 0x33, 0xae, 0x61, 0x8c, 0xe5, 0x2e, 0x0c, 0x82, 0x95, 0x2c,
	0xa6, 0xbe, 0x56, 0x9e, 0x14, 0x90, 0xc4, 0xd7, 0xfe, 0x4f, 0x89, 0xb4, 0x38, 0x32, 0x16, 0x87,
	0x88, 0x88, 0x6f, 0xb9, 0x73, 0x83, 0xff, 0xaa, 0xaf, 0x7e, 0x01, 0xca, 0x73, 0x42, 0x89, 0x24,
	0x30, 0x08, 0xda, 0x1e, 0xad, 0xcb, 0x82, 0x6b, 0x1e, 0xbb, 0xa2, 0x98, 0x27, 0x41, 0xe1, 0xd5,
	0x19, 0xf7, 0x04, 0x0c, 0xb0, 0x36, 0x61, 0x82, 0xd2, 0x74, 0xf9, 0x45, 0x1c, 0x19, 0x7b, 0xbd,
	0x7b, 0x1a, 0x3d, 0xdb, 0x72, 0x37, 0xfa, 0xf0, 0xe3, 0xba, 0x1c, 0xf8, 0x70, 0x92, 0x60, 0x8f,
	0xeb, 0x52, 0xbc, 0x65, 0xfc, 0x1d, 0x0c, 0xb0, 0xda, 0x01, 0xca, 0xf6, 0x28, 0xcc, 0x21, 0xa1,
	0x5e, 0x8b, 0x13, 0x84, 0xb5, 0x49, 0x13, 0x94, 0x66, 0x77, 0x97, 0xed, 0xac, 0x7b, 0xec, 0xa4,
	0x7b, 0xec, 0xbc, 0x7b, 0xec, 0x0a, 0x23, 0xb4, 0xdc, 0xb8, 0x89, 0x0c, 0x70, 0x17, 0x19, 0x9b,
	0x0d, 0x22, 0xcf, 0x2f, 0x6a, 0x36, 0x62, 0xa1, 0x93, 0xb7, 0x5a, 0xf6, 0xd9, 0x16, 0x7e, 0xd3,
	0x91, 0xed, 0x16, 0x16, 0x69, 0x42, 0x1c, 0x19, 0xfb, 0xa3, 0x1f, 0xb1, 0x2f, 0xc4, 0x72, 0x37,
	0xff, 0x72, 0xc6, 0x23, 0x42, 0x4f, 0x13, 0xa4, 0xfa, 0x09, 0x28, 0x4b, 0x69, 0x32, 0x86, 0x9c,
	0x12, 0xda, 0x10, 0x1e, 0xa3, 0xd9, 0x2d, 0x4f, 0x99, 0xe3, 0xa5, 0xd9, 0xdd, 0x97, 0xf6, 0x28,
	0x03, 0x63, 0x17, 0x3a, 0xfb, 0x4d, 0x46, 0x54, 0x36, 0xe3, 0xc8, 0x58, 0xcd, 0xc4, 0x3f, 0x58,
	0xc0, 0x72, 0xd5, 0x24, 0x9e, 0xc3, 0xc5, 0x09, 0x4d, 0x2f, 0xfd, 0x2b, 0x50, 0xd6, 0x86, 0xe0,
	0x6d, 0x82, 0x03, 0xdf, 0xe3, 0xf8, 0x0a, 0x72, 0x5f, 0x9b, 0xfe, 0x47, 0x5d, 0xa5, 0x38, 0x32,
	0xd6, 0x1f, 0xd1, 0x55, 0x2c, 0x64, 0xb9, 0xcb, 0xf7, 0xf5, 0x7d, 0x48, 0x36, 0xdd, 0x74, 0x4f,
	0xfd, 0x0e, 0x94, 0xf5, 0x7b, 0x96, 0x08, 0xef, 0x12, 0x5e, 0x04, 0xd2, 0x0b, 0x18, 0x6a, 0x7a,
	0x2d, 0xcc, 0x11, 0xa6, 0x12, 0x36, 0xb0, 0x36, 0x93, 0x8e, 0xc2, 0x59, 0xee, 0xfb, 0xc6, 0x08,
	0xbe, 0x57, 0xa9, 0x8c, 0x23, 0x63, 0xeb, 0x01, 0xdb, 0x1f, 0xa9, 0x61, 0xb9, 0x66, 0xd1, 0x6d,
	0xf1, 0x3e, 0x01, 0x1d, 0x32, 0xd4, 0x3c, 0xed, 0x43, 0xd4, 0x33, 0x65, 0xaa, 0x37, 0xe5, 0x4a,
	0xfa, 0x62, 0x54, 0xe2, 0xc8, 0x78, 0x9a, 0x0f, 0x67, 0x3e, 0xda, 0xbf, 0x22, 0x63, 0x7b, 0x04,
	0x7d, 0x07, 0x08, 0x1d, 0xf8, 0x3e, 0xc7, 0x42, 0xb8, 0x3d, 0xce, 0xb2, 0xfb, 0xad, 0xa3, 0x83,
	0x9b, 0x8e, 0x0e, 0x6e, 0x3b, 0x3a, 0xf8, 0xd9, 0xd1, 0xc1, 0xe7, 0xae, 0x3e, 0x76, 0xdb, 0xd5,
	0xc7, 0x7e, 0x74, 0xf5, 0xb1, 0x8f, 0xfb, 0x05, 0xd6, 0xbe, 0x6b, 0x4e, 0xc1, 0x35, 0xe7, 0xba,
	0xf8, 0x68, 0x67, 0x75, 0x6a, 0x93, 0xe9, 0xbb, 0xb9, 0xf7, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x67,
	0x99, 0xc6, 0x9b, 0xdf, 0x05, 0x00, 0x00,
}

func (this *MarketplaceCollection) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MarketplaceCollection)
	if !ok {
		that2, ok := that.(MarketplaceCollection)
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
	if !bytes.Equal(this.Index, that1.Index) {
		return false
	}
	if !bytes.Equal(this.MarketplaceIndex, that1.MarketplaceIndex) {
		return false
	}
	if this.CollectionCreator != that1.CollectionCreator {
		return false
	}
	if this.CollectionId != that1.CollectionId {
		return false
	}
	if this.InitiallyNftCollectionOwnerNftsForSale != that1.InitiallyNftCollectionOwnerNftsForSale {
		return false
	}
	if that1.InitiallyNftCollectionOwnerNftsMinPrice == nil {
		if this.InitiallyNftCollectionOwnerNftsMinPrice != nil {
			return false
		}
	} else if !this.InitiallyNftCollectionOwnerNftsMinPrice.Equal(*that1.InitiallyNftCollectionOwnerNftsMinPrice) {
		return false
	}
	if len(this.NftsEarningsOnSale) != len(that1.NftsEarningsOnSale) {
		return false
	}
	for i := range this.NftsEarningsOnSale {
		if !this.NftsEarningsOnSale[i].Equal(that1.NftsEarningsOnSale[i]) {
			return false
		}
	}
	if len(this.NftsEarningsOnYieldReward) != len(that1.NftsEarningsOnYieldReward) {
		return false
	}
	for i := range this.NftsEarningsOnYieldReward {
		if !this.NftsEarningsOnYieldReward[i].Equal(that1.NftsEarningsOnYieldReward[i]) {
			return false
		}
	}
	if that1.InitiallyNftsVaultLockPercentage == nil {
		if this.InitiallyNftsVaultLockPercentage != nil {
			return false
		}
	} else if !this.InitiallyNftsVaultLockPercentage.Equal(*that1.InitiallyNftsVaultLockPercentage) {
		return false
	}
	if !bytes.Equal(this.Creator, that1.Creator) {
		return false
	}
	return true
}
func (m *MarketplaceCollection) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketplaceCollection) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarketplaceCollection) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintMarketplaceCollection(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x52
	}
	if m.InitiallyNftsVaultLockPercentage != nil {
		{
			size := m.InitiallyNftsVaultLockPercentage.Size()
			i -= size
			if _, err := m.InitiallyNftsVaultLockPercentage.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintMarketplaceCollection(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
	}
	if len(m.NftsEarningsOnYieldReward) > 0 {
		for iNdEx := len(m.NftsEarningsOnYieldReward) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.NftsEarningsOnYieldReward[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketplaceCollection(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.NftsEarningsOnSale) > 0 {
		for iNdEx := len(m.NftsEarningsOnSale) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.NftsEarningsOnSale[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMarketplaceCollection(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.InitiallyNftCollectionOwnerNftsMinPrice != nil {
		{
			size := m.InitiallyNftCollectionOwnerNftsMinPrice.Size()
			i -= size
			if _, err := m.InitiallyNftCollectionOwnerNftsMinPrice.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintMarketplaceCollection(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.InitiallyNftCollectionOwnerNftsForSale {
		i--
		if m.InitiallyNftCollectionOwnerNftsForSale {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if len(m.CollectionId) > 0 {
		i -= len(m.CollectionId)
		copy(dAtA[i:], m.CollectionId)
		i = encodeVarintMarketplaceCollection(dAtA, i, uint64(len(m.CollectionId)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.CollectionCreator) > 0 {
		i -= len(m.CollectionCreator)
		copy(dAtA[i:], m.CollectionCreator)
		i = encodeVarintMarketplaceCollection(dAtA, i, uint64(len(m.CollectionCreator)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.MarketplaceIndex) > 0 {
		i -= len(m.MarketplaceIndex)
		copy(dAtA[i:], m.MarketplaceIndex)
		i = encodeVarintMarketplaceCollection(dAtA, i, uint64(len(m.MarketplaceIndex)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintMarketplaceCollection(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMarketplaceCollection(dAtA []byte, offset int, v uint64) int {
	offset -= sovMarketplaceCollection(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MarketplaceCollection) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovMarketplaceCollection(uint64(l))
	}
	l = len(m.MarketplaceIndex)
	if l > 0 {
		n += 1 + l + sovMarketplaceCollection(uint64(l))
	}
	l = len(m.CollectionCreator)
	if l > 0 {
		n += 1 + l + sovMarketplaceCollection(uint64(l))
	}
	l = len(m.CollectionId)
	if l > 0 {
		n += 1 + l + sovMarketplaceCollection(uint64(l))
	}
	if m.InitiallyNftCollectionOwnerNftsForSale {
		n += 2
	}
	if m.InitiallyNftCollectionOwnerNftsMinPrice != nil {
		l = m.InitiallyNftCollectionOwnerNftsMinPrice.Size()
		n += 1 + l + sovMarketplaceCollection(uint64(l))
	}
	if len(m.NftsEarningsOnSale) > 0 {
		for _, e := range m.NftsEarningsOnSale {
			l = e.Size()
			n += 1 + l + sovMarketplaceCollection(uint64(l))
		}
	}
	if len(m.NftsEarningsOnYieldReward) > 0 {
		for _, e := range m.NftsEarningsOnYieldReward {
			l = e.Size()
			n += 1 + l + sovMarketplaceCollection(uint64(l))
		}
	}
	if m.InitiallyNftsVaultLockPercentage != nil {
		l = m.InitiallyNftsVaultLockPercentage.Size()
		n += 1 + l + sovMarketplaceCollection(uint64(l))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovMarketplaceCollection(uint64(l))
	}
	return n
}

func sovMarketplaceCollection(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMarketplaceCollection(x uint64) (n int) {
	return sovMarketplaceCollection(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MarketplaceCollection) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMarketplaceCollection
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
			return fmt.Errorf("proto: MarketplaceCollection: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketplaceCollection: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = append(m.Index[:0], dAtA[iNdEx:postIndex]...)
			if m.Index == nil {
				m.Index = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MarketplaceIndex", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MarketplaceIndex = append(m.MarketplaceIndex[:0], dAtA[iNdEx:postIndex]...)
			if m.MarketplaceIndex == nil {
				m.MarketplaceIndex = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollectionCreator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CollectionCreator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CollectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitiallyNftCollectionOwnerNftsForSale", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
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
			m.InitiallyNftCollectionOwnerNftsForSale = bool(v != 0)
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitiallyNftCollectionOwnerNftsMinPrice", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
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
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.InitiallyNftCollectionOwnerNftsMinPrice == nil {
				m.InitiallyNftCollectionOwnerNftsMinPrice = &github_com_cosmos_cosmos_sdk_types.Coin{}
			}
			if err := m.InitiallyNftCollectionOwnerNftsMinPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NftsEarningsOnSale", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
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
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NftsEarningsOnSale = append(m.NftsEarningsOnSale, &MarketplaceEarning{})
			if err := m.NftsEarningsOnSale[len(m.NftsEarningsOnSale)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NftsEarningsOnYieldReward", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
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
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NftsEarningsOnYieldReward = append(m.NftsEarningsOnYieldReward, &MarketplaceEarning{})
			if err := m.NftsEarningsOnYieldReward[len(m.NftsEarningsOnYieldReward)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitiallyNftsVaultLockPercentage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Int
			m.InitiallyNftsVaultLockPercentage = &v
			if err := m.InitiallyNftsVaultLockPercentage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMarketplaceCollection
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMarketplaceCollection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = append(m.Creator[:0], dAtA[iNdEx:postIndex]...)
			if m.Creator == nil {
				m.Creator = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMarketplaceCollection(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMarketplaceCollection
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
func skipMarketplaceCollection(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMarketplaceCollection
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
					return 0, ErrIntOverflowMarketplaceCollection
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
					return 0, ErrIntOverflowMarketplaceCollection
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
				return 0, ErrInvalidLengthMarketplaceCollection
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMarketplaceCollection
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMarketplaceCollection
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMarketplaceCollection        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMarketplaceCollection          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMarketplaceCollection = fmt.Errorf("proto: unexpected end of group")
)

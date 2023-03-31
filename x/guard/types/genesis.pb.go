// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mantrachain/guard/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// GenesisState defines the guard module's genesis state.
type GenesisState struct {
	Params                 Params                `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	AccountPrivilegesList  []*AccountPrivileges  `protobuf:"bytes,2,rep,name=account_privileges_list,json=accountPrivilegesList,proto3" json:"account_privileges_list,omitempty"`
	GuardTransferCoins     []byte                `protobuf:"bytes,3,opt,name=guard_transfer_coins,json=guardTransferCoins,proto3" json:"guard_transfer_coins,omitempty"`
	RequiredPrivilegesList []*RequiredPrivileges `protobuf:"bytes,4,rep,name=required_privileges_list,json=requiredPrivilegesList,proto3" json:"required_privileges_list,omitempty"`
	// this line is used by starport scaffolding # genesis/proto/state
	LockedList []*Locked `protobuf:"bytes,5,rep,name=lockedList,proto3" json:"lockedList,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ea5a8529cb82408, []int{0}
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

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetAccountPrivilegesList() []*AccountPrivileges {
	if m != nil {
		return m.AccountPrivilegesList
	}
	return nil
}

func (m *GenesisState) GetGuardTransferCoins() []byte {
	if m != nil {
		return m.GuardTransferCoins
	}
	return nil
}

func (m *GenesisState) GetRequiredPrivilegesList() []*RequiredPrivileges {
	if m != nil {
		return m.RequiredPrivilegesList
	}
	return nil
}

func (m *GenesisState) GetLockedList() []*Locked {
	if m != nil {
		return m.LockedList
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "mantrachain.guard.v1.GenesisState")
}

func init() {
	proto.RegisterFile("mantrachain/guard/v1/genesis.proto", fileDescriptor_0ea5a8529cb82408)
}

var fileDescriptor_0ea5a8529cb82408 = []byte{
	// 353 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x41, 0x4f, 0xf2, 0x30,
	0x18, 0xc7, 0x37, 0xe0, 0xe5, 0x50, 0x38, 0x2d, 0xbc, 0xba, 0x10, 0x33, 0x91, 0x8b, 0xbb, 0xb8,
	0x09, 0xdc, 0x8c, 0x17, 0x30, 0xd1, 0x0b, 0x18, 0x32, 0x39, 0x79, 0x59, 0x4a, 0xa9, 0xa5, 0x91,
	0xb5, 0xb3, 0xed, 0x88, 0x7e, 0x0b, 0x3f, 0x16, 0x47, 0x8e, 0x9e, 0x8c, 0x81, 0x0f, 0xa2, 0xa1,
	0x5b, 0xcc, 0x22, 0xd3, 0x5b, 0xfb, 0x3c, 0xbf, 0xff, 0xd3, 0x5f, 0xfa, 0x80, 0x76, 0x04, 0x99,
	0x12, 0x10, 0xcd, 0x21, 0x65, 0x3e, 0x49, 0xa0, 0x98, 0xf9, 0xcb, 0x8e, 0x4f, 0x30, 0xc3, 0x92,
	0x4a, 0x2f, 0x16, 0x5c, 0x71, 0xab, 0x91, 0x63, 0x3c, 0xcd, 0x78, 0xcb, 0x4e, 0xb3, 0x41, 0x38,
	0xe1, 0x1a, 0xf0, 0x77, 0xa7, 0x94, 0x6d, 0x9e, 0x14, 0xce, 0x43, 0x3c, 0x8a, 0x38, 0xfb, 0x13,
	0x89, 0xa1, 0x80, 0x51, 0xf6, 0x62, 0xfb, 0xb3, 0x04, 0xea, 0x37, 0xa9, 0xc3, 0x9d, 0x82, 0x0a,
	0x5b, 0x17, 0xa0, 0x9a, 0x02, 0xb6, 0xd9, 0x32, 0xdd, 0x5a, 0xf7, 0xc8, 0x2b, 0x72, 0xf2, 0xc6,
	0x9a, 0x19, 0x54, 0x56, 0xef, 0xc7, 0x46, 0x90, 0x25, 0xac, 0x10, 0x1c, 0x42, 0x84, 0x78, 0xc2,
	0x54, 0x18, 0x0b, 0xba, 0xa4, 0x0b, 0x4c, 0xb0, 0x0c, 0x17, 0x54, 0x2a, 0xbb, 0xd4, 0x2a, 0xbb,
	0xb5, 0xee, 0x69, 0xf1, 0xb0, 0x7e, 0x1a, 0x1a, 0x7f, 0x67, 0x82, 0xff, 0xf0, 0x67, 0x69, 0x48,
	0xa5, 0xb2, 0xce, 0x41, 0x43, 0x87, 0x42, 0x25, 0x20, 0x93, 0x0f, 0x58, 0x84, 0x88, 0x53, 0x26,
	0xed, 0x72, 0xcb, 0x74, 0xeb, 0x81, 0xa5, 0x7b, 0x93, 0xac, 0x75, 0xb5, 0xeb, 0x58, 0x53, 0x60,
	0x0b, 0xfc, 0x94, 0x50, 0x81, 0x67, 0x7b, 0x4e, 0x15, 0xed, 0xe4, 0x16, 0x3b, 0x05, 0x59, 0x2a,
	0x27, 0x75, 0x20, 0xf6, 0x6a, 0xda, 0xea, 0x12, 0x80, 0x05, 0x47, 0x8f, 0x78, 0xb6, 0xbb, 0xd9,
	0xff, 0xf4, 0xd4, 0x5f, 0xbe, 0x6d, 0xa8, 0xb9, 0x20, 0xc7, 0x0f, 0x46, 0xab, 0x8d, 0x63, 0xae,
	0x37, 0x8e, 0xf9, 0xb1, 0x71, 0xcc, 0xd7, 0xad, 0x63, 0xac, 0xb7, 0x8e, 0xf1, 0xb6, 0x75, 0x8c,
	0xfb, 0x1e, 0xa1, 0x6a, 0x9e, 0x4c, 0x3d, 0xc4, 0x23, 0x7f, 0xd4, 0xbf, 0x9d, 0x04, 0xfd, 0xb3,
	0x6b, 0xca, 0x20, 0x43, 0xd8, 0xcf, 0x2f, 0xf6, 0x39, 0x5b, 0xad, 0x7a, 0x89, 0xb1, 0x9c, 0x56,
	0xf5, 0x5e, 0x7b, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xde, 0x08, 0x9c, 0x01, 0x6f, 0x02, 0x00,
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
	if len(m.LockedList) > 0 {
		for iNdEx := len(m.LockedList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LockedList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.RequiredPrivilegesList) > 0 {
		for iNdEx := len(m.RequiredPrivilegesList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RequiredPrivilegesList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.GuardTransferCoins) > 0 {
		i -= len(m.GuardTransferCoins)
		copy(dAtA[i:], m.GuardTransferCoins)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.GuardTransferCoins)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.AccountPrivilegesList) > 0 {
		for iNdEx := len(m.AccountPrivilegesList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AccountPrivilegesList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.AccountPrivilegesList) > 0 {
		for _, e := range m.AccountPrivilegesList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = len(m.GuardTransferCoins)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.RequiredPrivilegesList) > 0 {
		for _, e := range m.RequiredPrivilegesList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.LockedList) > 0 {
		for _, e := range m.LockedList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
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
				return fmt.Errorf("proto: wrong wireType = %d for field AccountPrivilegesList", wireType)
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
			m.AccountPrivilegesList = append(m.AccountPrivilegesList, &AccountPrivileges{})
			if err := m.AccountPrivilegesList[len(m.AccountPrivilegesList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GuardTransferCoins", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GuardTransferCoins = append(m.GuardTransferCoins[:0], dAtA[iNdEx:postIndex]...)
			if m.GuardTransferCoins == nil {
				m.GuardTransferCoins = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequiredPrivilegesList", wireType)
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
			m.RequiredPrivilegesList = append(m.RequiredPrivilegesList, &RequiredPrivileges{})
			if err := m.RequiredPrivilegesList[len(m.RequiredPrivilegesList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockedList", wireType)
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
			m.LockedList = append(m.LockedList, &Locked{})
			if err := m.LockedList[len(m.LockedList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
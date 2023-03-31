// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mantrachain/guard/v1/params.proto

package types

import (
	fmt "fmt"
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

// Params defines the parameters for the module.
type Params struct {
	AdminAccount                            string `protobuf:"bytes,1,opt,name=admin_account,json=adminAccount,proto3" json:"admin_account,omitempty"`
	AccountPrivilegesTokenCollectionCreator string `protobuf:"bytes,2,opt,name=account_privileges_token_collection_creator,json=accountPrivilegesTokenCollectionCreator,proto3" json:"account_privileges_token_collection_creator,omitempty"`
	AccountPrivilegesTokenCollectionId      string `protobuf:"bytes,3,opt,name=account_privileges_token_collection_id,json=accountPrivilegesTokenCollectionId,proto3" json:"account_privileges_token_collection_id,omitempty"`
	DefaultPrivileges                       []byte `protobuf:"bytes,4,opt,name=default_privileges,json=defaultPrivileges,proto3" json:"default_privileges,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7f2d974ff897741, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetAdminAccount() string {
	if m != nil {
		return m.AdminAccount
	}
	return ""
}

func (m *Params) GetAccountPrivilegesTokenCollectionCreator() string {
	if m != nil {
		return m.AccountPrivilegesTokenCollectionCreator
	}
	return ""
}

func (m *Params) GetAccountPrivilegesTokenCollectionId() string {
	if m != nil {
		return m.AccountPrivilegesTokenCollectionId
	}
	return ""
}

func (m *Params) GetDefaultPrivileges() []byte {
	if m != nil {
		return m.DefaultPrivileges
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "mantrachain.guard.v1.Params")
}

func init() { proto.RegisterFile("mantrachain/guard/v1/params.proto", fileDescriptor_e7f2d974ff897741) }

var fileDescriptor_e7f2d974ff897741 = []byte{
	// 280 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xc1, 0x4a, 0xc3, 0x30,
	0x18, 0xc7, 0x9b, 0x29, 0x03, 0xc3, 0x3c, 0x18, 0x3c, 0xf4, 0x14, 0xe6, 0x04, 0x1d, 0xc8, 0x5a,
	0xc6, 0x9e, 0xa0, 0x0e, 0x04, 0x0f, 0x93, 0x51, 0x76, 0x12, 0xa1, 0x7c, 0x4b, 0x63, 0x17, 0x6c,
	0x93, 0x92, 0xa6, 0x45, 0x5f, 0x42, 0x7c, 0x2c, 0x8f, 0x3b, 0x7a, 0x94, 0xf6, 0x45, 0xc4, 0xac,
	0xcc, 0xde, 0xdc, 0xf5, 0xfb, 0x7e, 0xbf, 0xdf, 0xe1, 0x8f, 0x2f, 0x32, 0x90, 0x46, 0x03, 0xdb,
	0x80, 0x90, 0x7e, 0x52, 0x82, 0x8e, 0xfd, 0x6a, 0xea, 0xe7, 0xa0, 0x21, 0x2b, 0xbc, 0x5c, 0x2b,
	0xa3, 0xc8, 0x79, 0x07, 0xf1, 0x2c, 0xe2, 0x55, 0xd3, 0xd1, 0x7b, 0x0f, 0xf7, 0x97, 0x16, 0x23,
	0x97, 0xf8, 0x14, 0xe2, 0x4c, 0xc8, 0x08, 0x18, 0x53, 0xa5, 0x34, 0x2e, 0x1a, 0xa2, 0xf1, 0x49,
	0x38, 0xb0, 0xc7, 0x60, 0x77, 0x23, 0x4f, 0xf8, 0xa6, 0x7d, 0x47, 0xb9, 0x16, 0x95, 0x48, 0x79,
	0xc2, 0x8b, 0xc8, 0xa8, 0x17, 0x2e, 0x23, 0xa6, 0xd2, 0x94, 0x33, 0x23, 0x94, 0x8c, 0x98, 0xe6,
	0x60, 0x94, 0x76, 0x7b, 0x36, 0x71, 0xdd, 0x2a, 0xcb, 0xbd, 0xb1, 0xfa, 0x15, 0xe6, 0x7b, 0x7e,
	0xbe, 0xc3, 0x49, 0x88, 0xaf, 0x0e, 0xa9, 0x8b, 0xd8, 0x3d, 0xb2, 0xe1, 0xd1, 0x7f, 0xe1, 0xfb,
	0x98, 0x4c, 0x30, 0x89, 0xf9, 0x33, 0x94, 0x69, 0xb7, 0xe9, 0x1e, 0x0f, 0xd1, 0x78, 0x10, 0x9e,
	0xb5, 0x9f, 0x3f, 0xff, 0x76, 0xf1, 0x59, 0x53, 0xb4, 0xad, 0x29, 0xfa, 0xae, 0x29, 0xfa, 0x68,
	0xa8, 0xb3, 0x6d, 0xa8, 0xf3, 0xd5, 0x50, 0xe7, 0x71, 0x96, 0x08, 0xb3, 0x29, 0xd7, 0x1e, 0x53,
	0x99, 0xbf, 0x08, 0x1e, 0x56, 0x61, 0x30, 0xb9, 0x13, 0x12, 0x24, 0xe3, 0x7e, 0x77, 0xfd, 0xd7,
	0x76, 0x7f, 0xf3, 0x96, 0xf3, 0x62, 0xdd, 0xb7, 0xe3, 0xcf, 0x7e, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x4b, 0xb6, 0x77, 0x74, 0xa1, 0x01, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DefaultPrivileges) > 0 {
		i -= len(m.DefaultPrivileges)
		copy(dAtA[i:], m.DefaultPrivileges)
		i = encodeVarintParams(dAtA, i, uint64(len(m.DefaultPrivileges)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.AccountPrivilegesTokenCollectionId) > 0 {
		i -= len(m.AccountPrivilegesTokenCollectionId)
		copy(dAtA[i:], m.AccountPrivilegesTokenCollectionId)
		i = encodeVarintParams(dAtA, i, uint64(len(m.AccountPrivilegesTokenCollectionId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.AccountPrivilegesTokenCollectionCreator) > 0 {
		i -= len(m.AccountPrivilegesTokenCollectionCreator)
		copy(dAtA[i:], m.AccountPrivilegesTokenCollectionCreator)
		i = encodeVarintParams(dAtA, i, uint64(len(m.AccountPrivilegesTokenCollectionCreator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.AdminAccount) > 0 {
		i -= len(m.AdminAccount)
		copy(dAtA[i:], m.AdminAccount)
		i = encodeVarintParams(dAtA, i, uint64(len(m.AdminAccount)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AdminAccount)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.AccountPrivilegesTokenCollectionCreator)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.AccountPrivilegesTokenCollectionId)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.DefaultPrivileges)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AdminAccount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AdminAccount = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountPrivilegesTokenCollectionCreator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccountPrivilegesTokenCollectionCreator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountPrivilegesTokenCollectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccountPrivilegesTokenCollectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultPrivileges", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DefaultPrivileges = append(m.DefaultPrivileges[:0], dAtA[iNdEx:postIndex]...)
			if m.DefaultPrivileges == nil {
				m.DefaultPrivileges = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
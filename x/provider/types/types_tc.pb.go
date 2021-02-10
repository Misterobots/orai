// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: x/provider/types/types_tc.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
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

type TestCase struct {
	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Contract string `protobuf:"bytes,2,opt,name=contract,proto3" json:"contract,omitempty"`
	// Owner is the address who is allowed to make further changes to the data source.
	Owner       github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,3,opt,name=owner,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"owner,omitempty"`
	Description string                                        `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Fees        github_com_cosmos_cosmos_sdk_types.Coins      `protobuf:"bytes,5,rep,name=fees,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"fees" json:"transaction_fee"`
}

func (m *TestCase) Reset()         { *m = TestCase{} }
func (m *TestCase) String() string { return proto.CompactTextString(m) }
func (*TestCase) ProtoMessage()    {}
func (*TestCase) Descriptor() ([]byte, []int) {
	return fileDescriptor_d4399028e7eeafc9, []int{0}
}
func (m *TestCase) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TestCase) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TestCase.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TestCase) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestCase.Merge(m, src)
}
func (m *TestCase) XXX_Size() int {
	return m.Size()
}
func (m *TestCase) XXX_DiscardUnknown() {
	xxx_messageInfo_TestCase.DiscardUnknown(m)
}

var xxx_messageInfo_TestCase proto.InternalMessageInfo

func (m *TestCase) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TestCase) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *TestCase) GetOwner() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *TestCase) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *TestCase) GetFees() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Fees
	}
	return nil
}

func init() {
	proto.RegisterType((*TestCase)(nil), "oraichain.orai.provider.TestCase")
}

func init() { proto.RegisterFile("x/provider/types/types_tc.proto", fileDescriptor_d4399028e7eeafc9) }

var fileDescriptor_d4399028e7eeafc9 = []byte{
	// 350 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x31, 0x4f, 0xc2, 0x40,
	0x14, 0xc7, 0x5b, 0x28, 0x06, 0x8b, 0x53, 0x63, 0xb4, 0x32, 0xb4, 0x0d, 0x53, 0x63, 0xc2, 0x5d,
	0xd0, 0x8d, 0x0d, 0x30, 0x31, 0x71, 0x24, 0x4e, 0x2e, 0xe4, 0x7a, 0x3d, 0xe0, 0x34, 0xdc, 0x6b,
	0xee, 0x4e, 0x94, 0x6f, 0xe1, 0xe8, 0xe8, 0xe4, 0xe0, 0x27, 0x61, 0x64, 0x74, 0x42, 0x03, 0xdf,
	0xc0, 0xd1, 0xc9, 0xb4, 0x87, 0x84, 0x38, 0xb9, 0xb4, 0xbf, 0xfb, 0xe7, 0xde, 0xff, 0xbd, 0xff,
	0x3d, 0x37, 0x7c, 0xc4, 0x99, 0x84, 0x29, 0x4f, 0x99, 0xc4, 0x7a, 0x96, 0x31, 0x65, 0xbe, 0x03,
	0x4d, 0x51, 0x26, 0x41, 0x83, 0x77, 0x0c, 0x92, 0x70, 0x3a, 0x26, 0x5c, 0xa0, 0x9c, 0xd0, 0xef,
	0xed, 0x7a, 0x40, 0x41, 0x4d, 0x40, 0xe1, 0x84, 0x28, 0x86, 0xa7, 0xad, 0x84, 0x69, 0xd2, 0xc2,
	0x14, 0xb8, 0x30, 0x85, 0xf5, 0xc3, 0x11, 0x8c, 0xa0, 0x40, 0x9c, 0x93, 0x51, 0x1b, 0xaf, 0x25,
	0xb7, 0x7a, 0xcd, 0x94, 0xee, 0x11, 0xc5, 0x3c, 0xcf, 0x75, 0x04, 0x99, 0x30, 0xdf, 0x8e, 0xec,
	0x78, 0xbf, 0x5f, 0xb0, 0x57, 0x77, 0xab, 0x14, 0x84, 0x96, 0x84, 0x6a, 0xbf, 0x54, 0xe8, 0xdb,
	0xb3, 0x77, 0xe9, 0x56, 0xe0, 0x41, 0x30, 0xe9, 0x97, 0x23, 0x3b, 0x3e, 0xe8, 0xb6, 0xbe, 0x97,
	0x61, 0x73, 0xc4, 0xf5, 0xf8, 0x3e, 0x41, 0x14, 0x26, 0x78, 0x33, 0x90, 0xf9, 0x35, 0x55, 0x7a,
	0x67, 0xb2, 0xa0, 0x0e, 0xa5, 0x9d, 0x34, 0x95, 0x4c, 0xa9, 0xbe, 0xa9, 0xf7, 0x22, 0xb7, 0x96,
	0x32, 0x45, 0x25, 0xcf, 0x34, 0x07, 0xe1, 0x3b, 0x45, 0x9f, 0x5d, 0xc9, 0x9b, 0xb9, 0xce, 0x90,
	0x31, 0xe5, 0x57, 0xa2, 0x72, 0x5c, 0x3b, 0x3b, 0x41, 0xc6, 0x14, 0xe5, 0x61, 0xd1, 0x26, 0x2c,
	0xea, 0x01, 0x17, 0xdd, 0xab, 0xf9, 0x32, 0xb4, 0xbe, 0x96, 0xe1, 0xd1, 0xad, 0x02, 0xd1, 0x6e,
	0x68, 0x49, 0x84, 0x22, 0x34, 0xf7, 0x18, 0x0c, 0x19, 0x6b, 0xbc, 0x7d, 0x84, 0xf1, 0x3f, 0x46,
	0xcc, 0xad, 0x54, 0xbf, 0x68, 0xd9, 0x76, 0x9e, 0x5f, 0x42, 0xbb, 0x7b, 0x31, 0x5f, 0x05, 0xf6,
	0x62, 0x15, 0xd8, 0x9f, 0xab, 0xc0, 0x7e, 0x5a, 0x07, 0xd6, 0x62, 0x1d, 0x58, 0xef, 0xeb, 0xc0,
	0xba, 0x39, 0xdd, 0xf1, 0xdb, 0x2e, 0xa7, 0x20, 0xfc, 0x77, 0x99, 0xc9, 0x5e, 0xf1, 0xea, 0xe7,
	0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x88, 0x94, 0x8f, 0xe7, 0x01, 0x00, 0x00,
}

func (m *TestCase) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TestCase) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TestCase) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Fees) > 0 {
		for iNdEx := len(m.Fees) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Fees[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypesTc(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintTypesTc(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTypesTc(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintTypesTc(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintTypesTc(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypesTc(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypesTc(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TestCase) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovTypesTc(uint64(l))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovTypesTc(uint64(l))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTypesTc(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovTypesTc(uint64(l))
	}
	if len(m.Fees) > 0 {
		for _, e := range m.Fees {
			l = e.Size()
			n += 1 + l + sovTypesTc(uint64(l))
		}
	}
	return n
}

func sovTypesTc(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypesTc(x uint64) (n int) {
	return sovTypesTc(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TestCase) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypesTc
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
			return fmt.Errorf("proto: TestCase: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TestCase: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypesTc
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
				return ErrInvalidLengthTypesTc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypesTc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypesTc
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
				return ErrInvalidLengthTypesTc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypesTc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypesTc
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
				return ErrInvalidLengthTypesTc
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypesTc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = append(m.Owner[:0], dAtA[iNdEx:postIndex]...)
			if m.Owner == nil {
				m.Owner = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypesTc
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
				return ErrInvalidLengthTypesTc
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypesTc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Fees", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypesTc
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
				return ErrInvalidLengthTypesTc
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypesTc
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Fees = append(m.Fees, types.Coin{})
			if err := m.Fees[len(m.Fees)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypesTc(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypesTc
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypesTc
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
func skipTypesTc(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypesTc
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
					return 0, ErrIntOverflowTypesTc
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
					return 0, ErrIntOverflowTypesTc
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
				return 0, ErrInvalidLengthTypesTc
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypesTc
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypesTc
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypesTc        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypesTc          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypesTc = fmt.Errorf("proto: unexpected end of group")
)

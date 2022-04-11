// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/cosmos/stake/staking.proto

package types

import (
	bytes "bytes"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/anypb"
	_ "google.golang.org/protobuf/types/known/durationpb"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type IBCParams struct {
	// unbonding_time is the time duration of unbonding.
	UnbondingTime time.Duration `protobuf:"bytes,1,opt,name=unbonding_time,json=unbondingTime,proto3,stdduration" json:"unbonding_time" yaml:"unbonding_time"`
	// max_validators is the maximum number of validators.
	MaxValidators uint32 `protobuf:"varint,2,opt,name=max_validators,json=maxValidators,proto3" json:"max_validators,omitempty" yaml:"max_validators"`
	// max_entries is the max entries for either unbonding delegation or redelegation (per pair/trio).
	MaxEntries uint32 `protobuf:"varint,3,opt,name=max_entries,json=maxEntries,proto3" json:"max_entries,omitempty" yaml:"max_entries"`
	// historical_entries is the number of historical entries to persist.
	HistoricalEntries uint32 `protobuf:"varint,4,opt,name=historical_entries,json=historicalEntries,proto3" json:"historical_entries,omitempty" yaml:"historical_entries"`
	// bond_denom defines the bondable coin denomination.
	BondDenom            string   `protobuf:"bytes,5,opt,name=bond_denom,json=bondDenom,proto3" json:"bond_denom,omitempty" yaml:"bond_denom"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IBCParams) Reset()      { *m = IBCParams{} }
func (*IBCParams) ProtoMessage() {}
func (*IBCParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_c9791118f0634123, []int{0}
}
func (m *IBCParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IBCParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IBCParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IBCParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IBCParams.Merge(m, src)
}
func (m *IBCParams) XXX_Size() int {
	return m.Size()
}
func (m *IBCParams) XXX_DiscardUnknown() {
	xxx_messageInfo_IBCParams.DiscardUnknown(m)
}

var xxx_messageInfo_IBCParams proto.InternalMessageInfo

func (m *IBCParams) GetUnbondingTime() time.Duration {
	if m != nil {
		return m.UnbondingTime
	}
	return 0
}

func (m *IBCParams) GetMaxValidators() uint32 {
	if m != nil {
		return m.MaxValidators
	}
	return 0
}

func (m *IBCParams) GetMaxEntries() uint32 {
	if m != nil {
		return m.MaxEntries
	}
	return 0
}

func (m *IBCParams) GetHistoricalEntries() uint32 {
	if m != nil {
		return m.HistoricalEntries
	}
	return 0
}

func (m *IBCParams) GetBondDenom() string {
	if m != nil {
		return m.BondDenom
	}
	return ""
}

func init() {
	proto.RegisterType((*IBCParams)(nil), "cosmos.staking.v1beta1.IBCParams")
}

func init() { proto.RegisterFile("proto/cosmos/stake/staking.proto", fileDescriptor_c9791118f0634123) }

var fileDescriptor_c9791118f0634123 = []byte{
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x4f, 0x8f, 0x93, 0x40,
	0x18, 0xc6, 0x3b, 0xeb, 0x9f, 0xd8, 0xd9, 0x74, 0x93, 0x25, 0xee, 0x06, 0x36, 0x91, 0x41, 0x4e,
	0x3d, 0x28, 0x64, 0xd5, 0xc4, 0x64, 0x0f, 0xc6, 0xe0, 0x7a, 0x30, 0xf1, 0x60, 0x88, 0xf1, 0xe0,
	0xa5, 0x19, 0x60, 0x64, 0x27, 0x65, 0x66, 0x1a, 0x66, 0x68, 0xe8, 0x37, 0xe9, 0xb1, 0x1f, 0xa7,
	0x47, 0x3f, 0x01, 0x9a, 0x7a, 0xf1, 0xcc, 0xcd, 0x9b, 0x81, 0x81, 0x92, 0x76, 0x2f, 0x30, 0xcf,
	0xfb, 0x7b, 0x9e, 0x87, 0x90, 0x79, 0xa1, 0xb3, 0xc8, 0x85, 0x12, 0x7e, 0x2c, 0x24, 0x13, 0xd2,
	0x97, 0x0a, 0xcf, 0x49, 0xfb, 0xa4, 0x3c, 0xf5, 0x5a, 0x64, 0x5c, 0x6a, 0xe6, 0xf5, 0xd3, 0xe5,
	0x75, 0x44, 0x14, 0xbe, 0xbe, 0x7a, 0x9a, 0x8a, 0x54, 0xe8, 0x74, 0x73, 0xd2, 0xee, 0x2b, 0x2b,
	0x15, 0x22, 0xcd, 0x88, 0xdf, 0xaa, 0xa8, 0xf8, 0xe1, 0x63, 0xbe, 0xea, 0x90, 0x7d, 0x8c, 0x92,
	0x22, 0xc7, 0x8a, 0x0a, 0xde, 0x71, 0x74, 0xcc, 0x15, 0x65, 0x44, 0x2a, 0xcc, 0x16, 0xda, 0xe0,
	0xfe, 0x3b, 0x81, 0xe3, 0x4f, 0xc1, 0x87, 0x2f, 0x38, 0xc7, 0x4c, 0x1a, 0x31, 0x3c, 0x2b, 0x78,
	0x24, 0x78, 0x42, 0x79, 0x3a, 0x6b, 0xac, 0x26, 0x70, 0xc0, 0xf4, 0xf4, 0x95, 0xe5, 0xe9, 0x1e,
	0xaf, 0xef, 0xf1, 0x6e, 0xbb, 0xef, 0x04, 0xcf, 0xb7, 0x15, 0x1a, 0xd5, 0x15, 0xba, 0x58, 0x61,
	0x96, 0xdd, 0xb8, 0x87, 0x71, 0x77, 0xfd, 0x0b, 0x81, 0x70, 0xb2, 0x1f, 0x7e, 0xa5, 0x8c, 0x18,
	0xef, 0xe1, 0x19, 0xc3, 0xe5, 0x6c, 0x89, 0x33, 0x9a, 0x60, 0x25, 0x72, 0x69, 0x9e, 0x38, 0x60,
	0x3a, 0x09, 0xac, 0xa1, 0xe5, 0x90, 0xbb, 0xe1, 0x84, 0xe1, 0xf2, 0xdb, 0x5e, 0x1b, 0x6f, 0xe1,
	0x69, 0xe3, 0x20, 0x5c, 0xe5, 0x94, 0x48, 0xf3, 0x41, 0x1b, 0xbf, 0xac, 0x2b, 0x64, 0x0c, 0xf1,
	0x0e, 0xba, 0x21, 0x64, 0xb8, 0xfc, 0xa8, 0x85, 0xf1, 0x19, 0x1a, 0x77, 0x54, 0x2a, 0x91, 0xd3,
	0x18, 0x67, 0xfb, 0xfc, 0xc3, 0x36, 0xff, 0xac, 0xae, 0x90, 0xa5, 0xf3, 0xf7, 0x3d, 0x6e, 0x78,
	0x3e, 0x0c, 0xfb, 0xb6, 0x37, 0x10, 0x36, 0xff, 0x35, 0x4b, 0x08, 0x17, 0xcc, 0x7c, 0xe4, 0x80,
	0xe9, 0x38, 0xb8, 0xa8, 0x2b, 0x74, 0xae, 0x5b, 0x06, 0xe6, 0x86, 0xe3, 0x46, 0xdc, 0x36, 0xe7,
	0x9b, 0x27, 0xeb, 0x0d, 0x1a, 0xfd, 0xdd, 0x20, 0x10, 0xbc, 0xdb, 0xee, 0x6c, 0xf0, 0x73, 0x67,
	0x83, 0xdf, 0x3b, 0x1b, 0xac, 0xff, 0xd8, 0xa3, 0xef, 0x2f, 0x52, 0xaa, 0xee, 0x8a, 0xc8, 0x8b,
	0x05, 0xeb, 0xd7, 0x47, 0xbf, 0x5e, 0xca, 0x64, 0xee, 0x97, 0xfd, 0x16, 0xf9, 0x6a, 0xb5, 0x20,
	0x32, 0x7a, 0xdc, 0xde, 0xc6, 0xeb, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1e, 0x8b, 0xd6, 0x98,
	0x70, 0x02, 0x00, 0x00,
}

func (this *IBCParams) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*IBCParams)
	if !ok {
		that2, ok := that.(IBCParams)
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
	if this.UnbondingTime != that1.UnbondingTime {
		return false
	}
	if this.MaxValidators != that1.MaxValidators {
		return false
	}
	if this.MaxEntries != that1.MaxEntries {
		return false
	}
	if this.HistoricalEntries != that1.HistoricalEntries {
		return false
	}
	if this.BondDenom != that1.BondDenom {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (m *IBCParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IBCParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IBCParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.BondDenom) > 0 {
		i -= len(m.BondDenom)
		copy(dAtA[i:], m.BondDenom)
		i = encodeVarintStaking(dAtA, i, uint64(len(m.BondDenom)))
		i--
		dAtA[i] = 0x2a
	}
	if m.HistoricalEntries != 0 {
		i = encodeVarintStaking(dAtA, i, uint64(m.HistoricalEntries))
		i--
		dAtA[i] = 0x20
	}
	if m.MaxEntries != 0 {
		i = encodeVarintStaking(dAtA, i, uint64(m.MaxEntries))
		i--
		dAtA[i] = 0x18
	}
	if m.MaxValidators != 0 {
		i = encodeVarintStaking(dAtA, i, uint64(m.MaxValidators))
		i--
		dAtA[i] = 0x10
	}
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.UnbondingTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.UnbondingTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintStaking(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintStaking(dAtA []byte, offset int, v uint64) int {
	offset -= sovStaking(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *IBCParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.UnbondingTime)
	n += 1 + l + sovStaking(uint64(l))
	if m.MaxValidators != 0 {
		n += 1 + sovStaking(uint64(m.MaxValidators))
	}
	if m.MaxEntries != 0 {
		n += 1 + sovStaking(uint64(m.MaxEntries))
	}
	if m.HistoricalEntries != 0 {
		n += 1 + sovStaking(uint64(m.HistoricalEntries))
	}
	l = len(m.BondDenom)
	if l > 0 {
		n += 1 + l + sovStaking(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovStaking(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStaking(x uint64) (n int) {
	return sovStaking(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *IBCParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStaking
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
			return fmt.Errorf("proto: IBCParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IBCParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStaking
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
				return ErrInvalidLengthStaking
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.UnbondingTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxValidators", wireType)
			}
			m.MaxValidators = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxValidators |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxEntries", wireType)
			}
			m.MaxEntries = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxEntries |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HistoricalEntries", wireType)
			}
			m.HistoricalEntries = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStaking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HistoricalEntries |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BondDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStaking
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
				return ErrInvalidLengthStaking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStaking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BondDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStaking(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStaking
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipStaking(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStaking
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
					return 0, ErrIntOverflowStaking
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
					return 0, ErrIntOverflowStaking
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
				return 0, ErrInvalidLengthStaking
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStaking
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStaking
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStaking        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStaking          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStaking = fmt.Errorf("proto: unexpected end of group")
)

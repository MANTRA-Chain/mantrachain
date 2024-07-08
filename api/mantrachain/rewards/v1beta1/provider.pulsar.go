// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package rewardsv1beta1

import (
	fmt "fmt"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sort "sort"
	sync "sync"
)

var _ protoreflect.List = (*_Provider_2_list)(nil)

type _Provider_2_list struct {
	list *[]*ProviderPair
}

func (x *_Provider_2_list) Len() int {
	if x.list == nil {
		return 0
	}
	return len(*x.list)
}

func (x *_Provider_2_list) Get(i int) protoreflect.Value {
	return protoreflect.ValueOfMessage((*x.list)[i].ProtoReflect())
}

func (x *_Provider_2_list) Set(i int, value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*ProviderPair)
	(*x.list)[i] = concreteValue
}

func (x *_Provider_2_list) Append(value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*ProviderPair)
	*x.list = append(*x.list, concreteValue)
}

func (x *_Provider_2_list) AppendMutable() protoreflect.Value {
	v := new(ProviderPair)
	*x.list = append(*x.list, v)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_Provider_2_list) Truncate(n int) {
	for i := n; i < len(*x.list); i++ {
		(*x.list)[i] = nil
	}
	*x.list = (*x.list)[:n]
}

func (x *_Provider_2_list) NewElement() protoreflect.Value {
	v := new(ProviderPair)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_Provider_2_list) IsValid() bool {
	return x.list != nil
}

var _ protoreflect.Map = (*_Provider_3_map)(nil)

type _Provider_3_map struct {
	m *map[uint64]uint64
}

func (x *_Provider_3_map) Len() int {
	if x.m == nil {
		return 0
	}
	return len(*x.m)
}

func (x *_Provider_3_map) Range(f func(protoreflect.MapKey, protoreflect.Value) bool) {
	if x.m == nil {
		return
	}
	for k, v := range *x.m {
		mapKey := (protoreflect.MapKey)(protoreflect.ValueOfUint64(k))
		mapValue := protoreflect.ValueOfUint64(v)
		if !f(mapKey, mapValue) {
			break
		}
	}
}

func (x *_Provider_3_map) Has(key protoreflect.MapKey) bool {
	if x.m == nil {
		return false
	}
	keyUnwrapped := key.Uint()
	concreteValue := keyUnwrapped
	_, ok := (*x.m)[concreteValue]
	return ok
}

func (x *_Provider_3_map) Clear(key protoreflect.MapKey) {
	if x.m == nil {
		return
	}
	keyUnwrapped := key.Uint()
	concreteKey := keyUnwrapped
	delete(*x.m, concreteKey)
}

func (x *_Provider_3_map) Get(key protoreflect.MapKey) protoreflect.Value {
	if x.m == nil {
		return protoreflect.Value{}
	}
	keyUnwrapped := key.Uint()
	concreteKey := keyUnwrapped
	v, ok := (*x.m)[concreteKey]
	if !ok {
		return protoreflect.Value{}
	}
	return protoreflect.ValueOfUint64(v)
}

func (x *_Provider_3_map) Set(key protoreflect.MapKey, value protoreflect.Value) {
	if !key.IsValid() || !value.IsValid() {
		panic("invalid key or value provided")
	}
	keyUnwrapped := key.Uint()
	concreteKey := keyUnwrapped
	valueUnwrapped := value.Uint()
	concreteValue := valueUnwrapped
	(*x.m)[concreteKey] = concreteValue
}

func (x *_Provider_3_map) Mutable(key protoreflect.MapKey) protoreflect.Value {
	panic("should not call Mutable on protoreflect.Map whose value is not of type protoreflect.Message")
}

func (x *_Provider_3_map) NewValue() protoreflect.Value {
	v := uint64(0)
	return protoreflect.ValueOfUint64(v)
}

func (x *_Provider_3_map) IsValid() bool {
	return x.m != nil
}

var (
	md_Provider                protoreflect.MessageDescriptor
	fd_Provider_index          protoreflect.FieldDescriptor
	fd_Provider_pairs          protoreflect.FieldDescriptor
	fd_Provider_pair_id_to_idx protoreflect.FieldDescriptor
)

func init() {
	file_mantrachain_rewards_v1beta1_provider_proto_init()
	md_Provider = File_mantrachain_rewards_v1beta1_provider_proto.Messages().ByName("Provider")
	fd_Provider_index = md_Provider.Fields().ByName("index")
	fd_Provider_pairs = md_Provider.Fields().ByName("pairs")
	fd_Provider_pair_id_to_idx = md_Provider.Fields().ByName("pair_id_to_idx")
}

var _ protoreflect.Message = (*fastReflection_Provider)(nil)

type fastReflection_Provider Provider

func (x *Provider) ProtoReflect() protoreflect.Message {
	return (*fastReflection_Provider)(x)
}

func (x *Provider) slowProtoReflect() protoreflect.Message {
	mi := &file_mantrachain_rewards_v1beta1_provider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_Provider_messageType fastReflection_Provider_messageType
var _ protoreflect.MessageType = fastReflection_Provider_messageType{}

type fastReflection_Provider_messageType struct{}

func (x fastReflection_Provider_messageType) Zero() protoreflect.Message {
	return (*fastReflection_Provider)(nil)
}
func (x fastReflection_Provider_messageType) New() protoreflect.Message {
	return new(fastReflection_Provider)
}
func (x fastReflection_Provider_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_Provider
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_Provider) Descriptor() protoreflect.MessageDescriptor {
	return md_Provider
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_Provider) Type() protoreflect.MessageType {
	return _fastReflection_Provider_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_Provider) New() protoreflect.Message {
	return new(fastReflection_Provider)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_Provider) Interface() protoreflect.ProtoMessage {
	return (*Provider)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_Provider) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.Index != "" {
		value := protoreflect.ValueOfString(x.Index)
		if !f(fd_Provider_index, value) {
			return
		}
	}
	if len(x.Pairs) != 0 {
		value := protoreflect.ValueOfList(&_Provider_2_list{list: &x.Pairs})
		if !f(fd_Provider_pairs, value) {
			return
		}
	}
	if len(x.PairIdToIdx) != 0 {
		value := protoreflect.ValueOfMap(&_Provider_3_map{m: &x.PairIdToIdx})
		if !f(fd_Provider_pair_id_to_idx, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_Provider) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "mantrachain.rewards.v1beta1.Provider.index":
		return x.Index != ""
	case "mantrachain.rewards.v1beta1.Provider.pairs":
		return len(x.Pairs) != 0
	case "mantrachain.rewards.v1beta1.Provider.pair_id_to_idx":
		return len(x.PairIdToIdx) != 0
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: mantrachain.rewards.v1beta1.Provider"))
		}
		panic(fmt.Errorf("message mantrachain.rewards.v1beta1.Provider does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Provider) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "mantrachain.rewards.v1beta1.Provider.index":
		x.Index = ""
	case "mantrachain.rewards.v1beta1.Provider.pairs":
		x.Pairs = nil
	case "mantrachain.rewards.v1beta1.Provider.pair_id_to_idx":
		x.PairIdToIdx = nil
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: mantrachain.rewards.v1beta1.Provider"))
		}
		panic(fmt.Errorf("message mantrachain.rewards.v1beta1.Provider does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_Provider) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "mantrachain.rewards.v1beta1.Provider.index":
		value := x.Index
		return protoreflect.ValueOfString(value)
	case "mantrachain.rewards.v1beta1.Provider.pairs":
		if len(x.Pairs) == 0 {
			return protoreflect.ValueOfList(&_Provider_2_list{})
		}
		listValue := &_Provider_2_list{list: &x.Pairs}
		return protoreflect.ValueOfList(listValue)
	case "mantrachain.rewards.v1beta1.Provider.pair_id_to_idx":
		if len(x.PairIdToIdx) == 0 {
			return protoreflect.ValueOfMap(&_Provider_3_map{})
		}
		mapValue := &_Provider_3_map{m: &x.PairIdToIdx}
		return protoreflect.ValueOfMap(mapValue)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: mantrachain.rewards.v1beta1.Provider"))
		}
		panic(fmt.Errorf("message mantrachain.rewards.v1beta1.Provider does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Provider) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "mantrachain.rewards.v1beta1.Provider.index":
		x.Index = value.Interface().(string)
	case "mantrachain.rewards.v1beta1.Provider.pairs":
		lv := value.List()
		clv := lv.(*_Provider_2_list)
		x.Pairs = *clv.list
	case "mantrachain.rewards.v1beta1.Provider.pair_id_to_idx":
		mv := value.Map()
		cmv := mv.(*_Provider_3_map)
		x.PairIdToIdx = *cmv.m
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: mantrachain.rewards.v1beta1.Provider"))
		}
		panic(fmt.Errorf("message mantrachain.rewards.v1beta1.Provider does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Provider) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "mantrachain.rewards.v1beta1.Provider.pairs":
		if x.Pairs == nil {
			x.Pairs = []*ProviderPair{}
		}
		value := &_Provider_2_list{list: &x.Pairs}
		return protoreflect.ValueOfList(value)
	case "mantrachain.rewards.v1beta1.Provider.pair_id_to_idx":
		if x.PairIdToIdx == nil {
			x.PairIdToIdx = make(map[uint64]uint64)
		}
		value := &_Provider_3_map{m: &x.PairIdToIdx}
		return protoreflect.ValueOfMap(value)
	case "mantrachain.rewards.v1beta1.Provider.index":
		panic(fmt.Errorf("field index of message mantrachain.rewards.v1beta1.Provider is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: mantrachain.rewards.v1beta1.Provider"))
		}
		panic(fmt.Errorf("message mantrachain.rewards.v1beta1.Provider does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_Provider) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "mantrachain.rewards.v1beta1.Provider.index":
		return protoreflect.ValueOfString("")
	case "mantrachain.rewards.v1beta1.Provider.pairs":
		list := []*ProviderPair{}
		return protoreflect.ValueOfList(&_Provider_2_list{list: &list})
	case "mantrachain.rewards.v1beta1.Provider.pair_id_to_idx":
		m := make(map[uint64]uint64)
		return protoreflect.ValueOfMap(&_Provider_3_map{m: &m})
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: mantrachain.rewards.v1beta1.Provider"))
		}
		panic(fmt.Errorf("message mantrachain.rewards.v1beta1.Provider does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_Provider) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in mantrachain.rewards.v1beta1.Provider", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_Provider) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Provider) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_Provider) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_Provider) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*Provider)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		l = len(x.Index)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if len(x.Pairs) > 0 {
			for _, e := range x.Pairs {
				l = options.Size(e)
				n += 1 + l + runtime.Sov(uint64(l))
			}
		}
		if len(x.PairIdToIdx) > 0 {
			SiZeMaP := func(k uint64, v uint64) {
				mapEntrySize := 1 + runtime.Sov(uint64(k)) + 1 + runtime.Sov(uint64(v))
				n += mapEntrySize + 1 + runtime.Sov(uint64(mapEntrySize))
			}
			if options.Deterministic {
				sortme := make([]uint64, 0, len(x.PairIdToIdx))
				for k := range x.PairIdToIdx {
					sortme = append(sortme, k)
				}
				sort.Slice(sortme, func(i, j int) bool {
					return sortme[i] < sortme[j]
				})
				for _, k := range sortme {
					v := x.PairIdToIdx[k]
					SiZeMaP(k, v)
				}
			} else {
				for k, v := range x.PairIdToIdx {
					SiZeMaP(k, v)
				}
			}
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*Provider)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if len(x.PairIdToIdx) > 0 {
			MaRsHaLmAp := func(k uint64, v uint64) (protoiface.MarshalOutput, error) {
				baseI := i
				i = runtime.EncodeVarint(dAtA, i, uint64(v))
				i--
				dAtA[i] = 0x10
				i = runtime.EncodeVarint(dAtA, i, uint64(k))
				i--
				dAtA[i] = 0x8
				i = runtime.EncodeVarint(dAtA, i, uint64(baseI-i))
				i--
				dAtA[i] = 0x1a
				return protoiface.MarshalOutput{}, nil
			}
			if options.Deterministic {
				keysForPairIdToIdx := make([]uint64, 0, len(x.PairIdToIdx))
				for k := range x.PairIdToIdx {
					keysForPairIdToIdx = append(keysForPairIdToIdx, uint64(k))
				}
				sort.Slice(keysForPairIdToIdx, func(i, j int) bool {
					return keysForPairIdToIdx[i] < keysForPairIdToIdx[j]
				})
				for iNdEx := len(keysForPairIdToIdx) - 1; iNdEx >= 0; iNdEx-- {
					v := x.PairIdToIdx[uint64(keysForPairIdToIdx[iNdEx])]
					out, err := MaRsHaLmAp(keysForPairIdToIdx[iNdEx], v)
					if err != nil {
						return out, err
					}
				}
			} else {
				for k := range x.PairIdToIdx {
					v := x.PairIdToIdx[k]
					out, err := MaRsHaLmAp(k, v)
					if err != nil {
						return out, err
					}
				}
			}
		}
		if len(x.Pairs) > 0 {
			for iNdEx := len(x.Pairs) - 1; iNdEx >= 0; iNdEx-- {
				encoded, err := options.Marshal(x.Pairs[iNdEx])
				if err != nil {
					return protoiface.MarshalOutput{
						NoUnkeyedLiterals: input.NoUnkeyedLiterals,
						Buf:               input.Buf,
					}, err
				}
				i -= len(encoded)
				copy(dAtA[i:], encoded)
				i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
				i--
				dAtA[i] = 0x12
			}
		}
		if len(x.Index) > 0 {
			i -= len(x.Index)
			copy(dAtA[i:], x.Index)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.Index)))
			i--
			dAtA[i] = 0xa
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*Provider)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: Provider: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: Provider: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Index = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 2:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Pairs", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Pairs = append(x.Pairs, &ProviderPair{})
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.Pairs[len(x.Pairs)-1]); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field PairIdToIdx", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if x.PairIdToIdx == nil {
					x.PairIdToIdx = make(map[uint64]uint64)
				}
				var mapkey uint64
				var mapvalue uint64
				for iNdEx < postIndex {
					entryPreIndex := iNdEx
					var wire uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
						}
						if iNdEx >= l {
							return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
								return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
							}
							if iNdEx >= l {
								return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
								return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
							}
							if iNdEx >= l {
								return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
						skippy, err := runtime.Skip(dAtA[iNdEx:])
						if err != nil {
							return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
						}
						if (skippy < 0) || (iNdEx+skippy) < 0 {
							return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
						}
						if (iNdEx + skippy) > postIndex {
							return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
						}
						iNdEx += skippy
					}
				}
				x.PairIdToIdx[mapkey] = mapvalue
				iNdEx = postIndex
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: mantrachain/rewards/v1beta1/provider.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Provider
type Provider struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index       string            `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	Pairs       []*ProviderPair   `protobuf:"bytes,2,rep,name=pairs,proto3" json:"pairs,omitempty"`
	PairIdToIdx map[uint64]uint64 `protobuf:"bytes,3,rep,name=pair_id_to_idx,json=pairIdToIdx,proto3" json:"pair_id_to_idx,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *Provider) Reset() {
	*x = Provider{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mantrachain_rewards_v1beta1_provider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Provider) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Provider) ProtoMessage() {}

// Deprecated: Use Provider.ProtoReflect.Descriptor instead.
func (*Provider) Descriptor() ([]byte, []int) {
	return file_mantrachain_rewards_v1beta1_provider_proto_rawDescGZIP(), []int{0}
}

func (x *Provider) GetIndex() string {
	if x != nil {
		return x.Index
	}
	return ""
}

func (x *Provider) GetPairs() []*ProviderPair {
	if x != nil {
		return x.Pairs
	}
	return nil
}

func (x *Provider) GetPairIdToIdx() map[uint64]uint64 {
	if x != nil {
		return x.PairIdToIdx
	}
	return nil
}

var File_mantrachain_rewards_v1beta1_provider_proto protoreflect.FileDescriptor

var file_mantrachain_rewards_v1beta1_provider_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x6d, 0x61, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x72, 0x65,
	0x77, 0x61, 0x72, 0x64, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x6d, 0x61,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64,
	0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x28, 0x6d, 0x61, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x2f, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xfe, 0x01, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x3f, 0x0a, 0x05, 0x70, 0x61, 0x69, 0x72, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6d, 0x61, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x50, 0x61, 0x69, 0x72,
	0x52, 0x05, 0x70, 0x61, 0x69, 0x72, 0x73, 0x12, 0x5b, 0x0a, 0x0e, 0x70, 0x61, 0x69, 0x72, 0x5f,
	0x69, 0x64, 0x5f, 0x74, 0x6f, 0x5f, 0x69, 0x64, 0x78, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x36, 0x2e, 0x6d, 0x61, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65,
	0x77, 0x61, 0x72, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x50, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x50, 0x61, 0x69, 0x72, 0x49, 0x64, 0x54, 0x6f, 0x49,
	0x64, 0x78, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x70, 0x61, 0x69, 0x72, 0x49, 0x64, 0x54,
	0x6f, 0x49, 0x64, 0x78, 0x1a, 0x3e, 0x0a, 0x10, 0x50, 0x61, 0x69, 0x72, 0x49, 0x64, 0x54, 0x6f,
	0x49, 0x64, 0x78, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x42, 0xfb, 0x01, 0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x61, 0x6e,
	0x74, 0x72, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73,
	0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x42, 0x0d, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3b, 0x63, 0x6f, 0x73, 0x6d, 0x6f,
	0x73, 0x73, 0x64, 0x6b, 0x2e, 0x69, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x61, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x2f,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x4d, 0x52, 0x58, 0xaa, 0x02, 0x1b, 0x4d,
	0x61, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x52, 0x65, 0x77, 0x61, 0x72,
	0x64, 0x73, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xca, 0x02, 0x1b, 0x4d, 0x61, 0x6e,
	0x74, 0x72, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5c, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73,
	0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xe2, 0x02, 0x27, 0x4d, 0x61, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5c, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x5c, 0x56,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x1d, 0x4d, 0x61, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x3a, 0x3a, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mantrachain_rewards_v1beta1_provider_proto_rawDescOnce sync.Once
	file_mantrachain_rewards_v1beta1_provider_proto_rawDescData = file_mantrachain_rewards_v1beta1_provider_proto_rawDesc
)

func file_mantrachain_rewards_v1beta1_provider_proto_rawDescGZIP() []byte {
	file_mantrachain_rewards_v1beta1_provider_proto_rawDescOnce.Do(func() {
		file_mantrachain_rewards_v1beta1_provider_proto_rawDescData = protoimpl.X.CompressGZIP(file_mantrachain_rewards_v1beta1_provider_proto_rawDescData)
	})
	return file_mantrachain_rewards_v1beta1_provider_proto_rawDescData
}

var file_mantrachain_rewards_v1beta1_provider_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_mantrachain_rewards_v1beta1_provider_proto_goTypes = []interface{}{
	(*Provider)(nil),     // 0: mantrachain.rewards.v1beta1.Provider
	nil,                  // 1: mantrachain.rewards.v1beta1.Provider.PairIdToIdxEntry
	(*ProviderPair)(nil), // 2: mantrachain.rewards.v1beta1.ProviderPair
}
var file_mantrachain_rewards_v1beta1_provider_proto_depIdxs = []int32{
	2, // 0: mantrachain.rewards.v1beta1.Provider.pairs:type_name -> mantrachain.rewards.v1beta1.ProviderPair
	1, // 1: mantrachain.rewards.v1beta1.Provider.pair_id_to_idx:type_name -> mantrachain.rewards.v1beta1.Provider.PairIdToIdxEntry
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_mantrachain_rewards_v1beta1_provider_proto_init() }
func file_mantrachain_rewards_v1beta1_provider_proto_init() {
	if File_mantrachain_rewards_v1beta1_provider_proto != nil {
		return
	}
	file_mantrachain_rewards_v1beta1_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_mantrachain_rewards_v1beta1_provider_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Provider); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mantrachain_rewards_v1beta1_provider_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mantrachain_rewards_v1beta1_provider_proto_goTypes,
		DependencyIndexes: file_mantrachain_rewards_v1beta1_provider_proto_depIdxs,
		MessageInfos:      file_mantrachain_rewards_v1beta1_provider_proto_msgTypes,
	}.Build()
	File_mantrachain_rewards_v1beta1_provider_proto = out.File
	file_mantrachain_rewards_v1beta1_provider_proto_rawDesc = nil
	file_mantrachain_rewards_v1beta1_provider_proto_goTypes = nil
	file_mantrachain_rewards_v1beta1_provider_proto_depIdxs = nil
}
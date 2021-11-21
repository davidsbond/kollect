// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: kollect/event/v1/event.proto

// Package kollect.event.v1 contains the definition of the event envelope which wraps all events published by kollect
// components.

package event

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Envelope is used as a wrapper around all published events that contains additional metadata, such as a unique
// identifier, timestamps etc.
type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// id is a unique identifier for the event.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// timestamp represents the time at which an event was created.
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// applies_at indicates the time at which the event described in Payload happened. It does not always
	// match Timestamp.
	AppliesAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=applies_at,json=appliesAt,proto3" json:"applies_at,omitempty"`
	// payload contains the content of the event,
	Payload *anypb.Any `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kollect_event_v1_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_kollect_event_v1_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_kollect_event_v1_event_proto_rawDescGZIP(), []int{0}
}

func (x *Envelope) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Envelope) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Envelope) GetAppliesAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AppliesAt
	}
	return nil
}

func (x *Envelope) GetPayload() *anypb.Any {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_kollect_event_v1_event_proto protoreflect.FileDescriptor

var file_kollect_event_v1_event_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x6b, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f,
	0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10,
	0x6b, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbf, 0x01, 0x0a,
	0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x39, 0x0a, 0x0a, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x73, 0x5f, 0x61,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x65, 0x73, 0x41, 0x74, 0x12, 0x2e,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x3c,
	0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x76,
	0x69, 0x64, 0x73, 0x62, 0x6f, 0x6e, 0x64, 0x2f, 0x6b, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x2f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kollect_event_v1_event_proto_rawDescOnce sync.Once
	file_kollect_event_v1_event_proto_rawDescData = file_kollect_event_v1_event_proto_rawDesc
)

func file_kollect_event_v1_event_proto_rawDescGZIP() []byte {
	file_kollect_event_v1_event_proto_rawDescOnce.Do(func() {
		file_kollect_event_v1_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_kollect_event_v1_event_proto_rawDescData)
	})
	return file_kollect_event_v1_event_proto_rawDescData
}

var file_kollect_event_v1_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_kollect_event_v1_event_proto_goTypes = []interface{}{
	(*Envelope)(nil),              // 0: kollect.event.v1.Envelope
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*anypb.Any)(nil),             // 2: google.protobuf.Any
}
var file_kollect_event_v1_event_proto_depIdxs = []int32{
	1, // 0: kollect.event.v1.Envelope.timestamp:type_name -> google.protobuf.Timestamp
	1, // 1: kollect.event.v1.Envelope.applies_at:type_name -> google.protobuf.Timestamp
	2, // 2: kollect.event.v1.Envelope.payload:type_name -> google.protobuf.Any
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_kollect_event_v1_event_proto_init() }
func file_kollect_event_v1_event_proto_init() {
	if File_kollect_event_v1_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kollect_event_v1_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Envelope); i {
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
			RawDescriptor: file_kollect_event_v1_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kollect_event_v1_event_proto_goTypes,
		DependencyIndexes: file_kollect_event_v1_event_proto_depIdxs,
		MessageInfos:      file_kollect_event_v1_event_proto_msgTypes,
	}.Build()
	File_kollect_event_v1_event_proto = out.File
	file_kollect_event_v1_event_proto_rawDesc = nil
	file_kollect_event_v1_event_proto_goTypes = nil
	file_kollect_event_v1_event_proto_depIdxs = nil
}

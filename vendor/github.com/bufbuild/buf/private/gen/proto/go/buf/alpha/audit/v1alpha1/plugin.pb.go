// Copyright 2020-2022 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        (unknown)
// source: buf/alpha/audit/v1alpha1/plugin.proto

package auditv1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BufAlphaRegistryV1Alpha1PluginVisibility int32

const (
	BufAlphaRegistryV1Alpha1PluginVisibility_BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_UNSPECIFIED BufAlphaRegistryV1Alpha1PluginVisibility = 0
	BufAlphaRegistryV1Alpha1PluginVisibility_BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_PUBLIC      BufAlphaRegistryV1Alpha1PluginVisibility = 1
	BufAlphaRegistryV1Alpha1PluginVisibility_BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_PRIVATE     BufAlphaRegistryV1Alpha1PluginVisibility = 2
)

// Enum value maps for BufAlphaRegistryV1Alpha1PluginVisibility.
var (
	BufAlphaRegistryV1Alpha1PluginVisibility_name = map[int32]string{
		0: "BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_UNSPECIFIED",
		1: "BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_PUBLIC",
		2: "BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_PRIVATE",
	}
	BufAlphaRegistryV1Alpha1PluginVisibility_value = map[string]int32{
		"BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_UNSPECIFIED": 0,
		"BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_PUBLIC":      1,
		"BUF_ALPHA_REGISTRY_V1_ALPHA1_PLUGIN_VISIBILITY_PRIVATE":     2,
	}
)

func (x BufAlphaRegistryV1Alpha1PluginVisibility) Enum() *BufAlphaRegistryV1Alpha1PluginVisibility {
	p := new(BufAlphaRegistryV1Alpha1PluginVisibility)
	*p = x
	return p
}

func (x BufAlphaRegistryV1Alpha1PluginVisibility) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BufAlphaRegistryV1Alpha1PluginVisibility) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_audit_v1alpha1_plugin_proto_enumTypes[0].Descriptor()
}

func (BufAlphaRegistryV1Alpha1PluginVisibility) Type() protoreflect.EnumType {
	return &file_buf_alpha_audit_v1alpha1_plugin_proto_enumTypes[0]
}

func (x BufAlphaRegistryV1Alpha1PluginVisibility) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BufAlphaRegistryV1Alpha1PluginVisibility.Descriptor instead.
func (BufAlphaRegistryV1Alpha1PluginVisibility) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescGZIP(), []int{0}
}

type BufAlphaRegistryV1Alpha1PluginVersionMapping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PluginOwner string `protobuf:"bytes,1,opt,name=plugin_owner,json=pluginOwner,proto3" json:"plugin_owner,omitempty"`
	PluginName  string `protobuf:"bytes,2,opt,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	Version     string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Deleted     bool   `protobuf:"varint,4,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionMapping) Reset() {
	*x = BufAlphaRegistryV1Alpha1PluginVersionMapping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionMapping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BufAlphaRegistryV1Alpha1PluginVersionMapping) ProtoMessage() {}

func (x *BufAlphaRegistryV1Alpha1PluginVersionMapping) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BufAlphaRegistryV1Alpha1PluginVersionMapping.ProtoReflect.Descriptor instead.
func (*BufAlphaRegistryV1Alpha1PluginVersionMapping) Descriptor() ([]byte, []int) {
	return file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionMapping) GetPluginOwner() string {
	if x != nil {
		return x.PluginOwner
	}
	return ""
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionMapping) GetPluginName() string {
	if x != nil {
		return x.PluginName
	}
	return ""
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionMapping) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionMapping) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

type BufAlphaRegistryV1Alpha1PluginConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PluginOwner string   `protobuf:"bytes,1,opt,name=plugin_owner,json=pluginOwner,proto3" json:"plugin_owner,omitempty"`
	PluginName  string   `protobuf:"bytes,2,opt,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	Parameters  []string `protobuf:"bytes,3,rep,name=parameters,proto3" json:"parameters,omitempty"`
	Deleted     bool     `protobuf:"varint,4,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *BufAlphaRegistryV1Alpha1PluginConfig) Reset() {
	*x = BufAlphaRegistryV1Alpha1PluginConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BufAlphaRegistryV1Alpha1PluginConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BufAlphaRegistryV1Alpha1PluginConfig) ProtoMessage() {}

func (x *BufAlphaRegistryV1Alpha1PluginConfig) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BufAlphaRegistryV1Alpha1PluginConfig.ProtoReflect.Descriptor instead.
func (*BufAlphaRegistryV1Alpha1PluginConfig) Descriptor() ([]byte, []int) {
	return file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *BufAlphaRegistryV1Alpha1PluginConfig) GetPluginOwner() string {
	if x != nil {
		return x.PluginOwner
	}
	return ""
}

func (x *BufAlphaRegistryV1Alpha1PluginConfig) GetPluginName() string {
	if x != nil {
		return x.PluginName
	}
	return ""
}

func (x *BufAlphaRegistryV1Alpha1PluginConfig) GetParameters() []string {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *BufAlphaRegistryV1Alpha1PluginConfig) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

type BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary) Reset() {
	*x = BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary) ProtoMessage() {}

func (x *BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary.ProtoReflect.Descriptor instead.
func (*BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary) Descriptor() ([]byte, []int) {
	return file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescGZIP(), []int{2}
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_buf_alpha_audit_v1alpha1_plugin_proto protoreflect.FileDescriptor

var file_buf_alpha_audit_v1alpha1_plugin_proto_rawDesc = []byte{
	0x0a, 0x25, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x61, 0x75, 0x64, 0x69,
	0x74, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x22, 0xa6, 0x01, 0x0a, 0x2c, 0x42, 0x75, 0x66, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x56, 0x31, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x50, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x70, 0x70, 0x69,
	0x6e, 0x67, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x6f, 0x77, 0x6e,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x4f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0xa4, 0x01, 0x0a, 0x24, 0x42,
	0x75, 0x66, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x56,
	0x31, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5f, 0x6f, 0x77,
	0x6e, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x22, 0x63, 0x0a, 0x33, 0x42, 0x75, 0x66, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x56, 0x31, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x50, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x4c, 0x69, 0x62, 0x72, 0x61, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2a, 0xe1, 0x01, 0x0a, 0x28, 0x42, 0x75, 0x66, 0x41, 0x6c,
	0x70, 0x68, 0x61, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x56, 0x31, 0x41, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c,
	0x69, 0x74, 0x79, 0x12, 0x3e, 0x0a, 0x3a, 0x42, 0x55, 0x46, 0x5f, 0x41, 0x4c, 0x50, 0x48, 0x41,
	0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x52, 0x59, 0x5f, 0x56, 0x31, 0x5f, 0x41, 0x4c, 0x50,
	0x48, 0x41, 0x31, 0x5f, 0x50, 0x4c, 0x55, 0x47, 0x49, 0x4e, 0x5f, 0x56, 0x49, 0x53, 0x49, 0x42,
	0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x39, 0x0a, 0x35, 0x42, 0x55, 0x46, 0x5f, 0x41, 0x4c, 0x50, 0x48, 0x41,
	0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x52, 0x59, 0x5f, 0x56, 0x31, 0x5f, 0x41, 0x4c, 0x50,
	0x48, 0x41, 0x31, 0x5f, 0x50, 0x4c, 0x55, 0x47, 0x49, 0x4e, 0x5f, 0x56, 0x49, 0x53, 0x49, 0x42,
	0x49, 0x4c, 0x49, 0x54, 0x59, 0x5f, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x10, 0x01, 0x12, 0x3a,
	0x0a, 0x36, 0x42, 0x55, 0x46, 0x5f, 0x41, 0x4c, 0x50, 0x48, 0x41, 0x5f, 0x52, 0x45, 0x47, 0x49,
	0x53, 0x54, 0x52, 0x59, 0x5f, 0x56, 0x31, 0x5f, 0x41, 0x4c, 0x50, 0x48, 0x41, 0x31, 0x5f, 0x50,
	0x4c, 0x55, 0x47, 0x49, 0x4e, 0x5f, 0x56, 0x49, 0x53, 0x49, 0x42, 0x49, 0x4c, 0x49, 0x54, 0x59,
	0x5f, 0x50, 0x52, 0x49, 0x56, 0x41, 0x54, 0x45, 0x10, 0x02, 0x42, 0x83, 0x02, 0x0a, 0x1c, 0x63,
	0x6f, 0x6d, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x61, 0x75, 0x64,
	0x69, 0x74, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0b, 0x50, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x53, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f,
	0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x3b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xa2,
	0x02, 0x03, 0x42, 0x41, 0x41, 0xaa, 0x02, 0x18, 0x42, 0x75, 0x66, 0x2e, 0x41, 0x6c, 0x70, 0x68,
	0x61, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x74, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0xca, 0x02, 0x18, 0x42, 0x75, 0x66, 0x5c, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x5c, 0x41, 0x75, 0x64,
	0x69, 0x74, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x24, 0x42, 0x75,
	0x66, 0x5c, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x5c, 0x41, 0x75, 0x64, 0x69, 0x74, 0x5c, 0x56, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x1b, 0x42, 0x75, 0x66, 0x3a, 0x3a, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x3a,
	0x3a, 0x41, 0x75, 0x64, 0x69, 0x74, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescOnce sync.Once
	file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescData = file_buf_alpha_audit_v1alpha1_plugin_proto_rawDesc
)

func file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescGZIP() []byte {
	file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescOnce.Do(func() {
		file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescData)
	})
	return file_buf_alpha_audit_v1alpha1_plugin_proto_rawDescData
}

var file_buf_alpha_audit_v1alpha1_plugin_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_buf_alpha_audit_v1alpha1_plugin_proto_goTypes = []interface{}{
	(BufAlphaRegistryV1Alpha1PluginVisibility)(0),               // 0: buf.alpha.audit.v1alpha1.BufAlphaRegistryV1Alpha1PluginVisibility
	(*BufAlphaRegistryV1Alpha1PluginVersionMapping)(nil),        // 1: buf.alpha.audit.v1alpha1.BufAlphaRegistryV1Alpha1PluginVersionMapping
	(*BufAlphaRegistryV1Alpha1PluginConfig)(nil),                // 2: buf.alpha.audit.v1alpha1.BufAlphaRegistryV1Alpha1PluginConfig
	(*BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary)(nil), // 3: buf.alpha.audit.v1alpha1.BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary
}
var file_buf_alpha_audit_v1alpha1_plugin_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_buf_alpha_audit_v1alpha1_plugin_proto_init() }
func file_buf_alpha_audit_v1alpha1_plugin_proto_init() {
	if File_buf_alpha_audit_v1alpha1_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BufAlphaRegistryV1Alpha1PluginVersionMapping); i {
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
		file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BufAlphaRegistryV1Alpha1PluginConfig); i {
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
		file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BufAlphaRegistryV1Alpha1PluginVersionRuntimeLibrary); i {
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
			RawDescriptor: file_buf_alpha_audit_v1alpha1_plugin_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_buf_alpha_audit_v1alpha1_plugin_proto_goTypes,
		DependencyIndexes: file_buf_alpha_audit_v1alpha1_plugin_proto_depIdxs,
		EnumInfos:         file_buf_alpha_audit_v1alpha1_plugin_proto_enumTypes,
		MessageInfos:      file_buf_alpha_audit_v1alpha1_plugin_proto_msgTypes,
	}.Build()
	File_buf_alpha_audit_v1alpha1_plugin_proto = out.File
	file_buf_alpha_audit_v1alpha1_plugin_proto_rawDesc = nil
	file_buf_alpha_audit_v1alpha1_plugin_proto_goTypes = nil
	file_buf_alpha_audit_v1alpha1_plugin_proto_depIdxs = nil
}

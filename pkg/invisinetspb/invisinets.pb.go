//
//Copyright 2023 The Invisinets Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: invisinets.proto

package invisinetspb

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

type Direction int32

const (
	Direction_INBOUND  Direction = 0
	Direction_OUTBOUND Direction = 1
)

// Enum value maps for Direction.
var (
	Direction_name = map[int32]string{
		0: "INBOUND",
		1: "OUTBOUND",
	}
	Direction_value = map[string]int32{
		"INBOUND":  0,
		"OUTBOUND": 1,
	}
)

func (x Direction) Enum() *Direction {
	p := new(Direction)
	*p = x
	return p
}

func (x Direction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Direction) Descriptor() protoreflect.EnumDescriptor {
	return file_invisinets_proto_enumTypes[0].Descriptor()
}

func (Direction) Type() protoreflect.EnumType {
	return &file_invisinets_proto_enumTypes[0]
}

func (x Direction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Direction.Descriptor instead.
func (Direction) EnumDescriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{0}
}

// Provides the necessary URI/ID to find the Invisinets networks (eg, subscription + resource group in Azure or project in GCP)
type InvisinetsDeployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *InvisinetsDeployment) Reset() {
	*x = InvisinetsDeployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvisinetsDeployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvisinetsDeployment) ProtoMessage() {}

func (x *InvisinetsDeployment) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvisinetsDeployment.ProtoReflect.Descriptor instead.
func (*InvisinetsDeployment) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{0}
}

func (x *InvisinetsDeployment) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type AddressSpaceList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AddressSpaces []string `protobuf:"bytes,1,rep,name=address_spaces,json=addressSpaces,proto3" json:"address_spaces,omitempty"`
}

func (x *AddressSpaceList) Reset() {
	*x = AddressSpaceList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressSpaceList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressSpaceList) ProtoMessage() {}

func (x *AddressSpaceList) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressSpaceList.ProtoReflect.Descriptor instead.
func (*AddressSpaceList) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{1}
}

func (x *AddressSpaceList) GetAddressSpaces() []string {
	if x != nil {
		return x.AddressSpaces
	}
	return nil
}

type ResourceID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ResourceID) Reset() {
	*x = ResourceID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceID) ProtoMessage() {}

func (x *ResourceID) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceID.ProtoReflect.Descriptor instead.
func (*ResourceID) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{2}
}

func (x *ResourceID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ResourceDescription struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description  []byte `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	AddressSpace string `protobuf:"bytes,2,opt,name=address_space,json=addressSpace,proto3" json:"address_space,omitempty"`
}

func (x *ResourceDescription) Reset() {
	*x = ResourceDescription{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceDescription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceDescription) ProtoMessage() {}

func (x *ResourceDescription) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceDescription.ProtoReflect.Descriptor instead.
func (*ResourceDescription) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{3}
}

func (x *ResourceDescription) GetDescription() []byte {
	if x != nil {
		return x.Description
	}
	return nil
}

func (x *ResourceDescription) GetAddressSpace() string {
	if x != nil {
		return x.AddressSpace
	}
	return ""
}

type BasicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success         bool        `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message         string      `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	UpdatedResource *ResourceID `protobuf:"bytes,3,opt,name=updated_resource,json=updatedResource,proto3,oneof" json:"updated_resource,omitempty"`
}

func (x *BasicResponse) Reset() {
	*x = BasicResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicResponse) ProtoMessage() {}

func (x *BasicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicResponse.ProtoReflect.Descriptor instead.
func (*BasicResponse) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{4}
}

func (x *BasicResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *BasicResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *BasicResponse) GetUpdatedResource() *ResourceID {
	if x != nil {
		return x.UpdatedResource
	}
	return nil
}

type PermitListRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Tag       []string  `protobuf:"bytes,2,rep,name=tag,proto3" json:"tag,omitempty"`
	Direction Direction `protobuf:"varint,3,opt,name=direction,proto3,enum=invisinetspb.Direction" json:"direction,omitempty"`
	SrcPort   int32     `protobuf:"varint,4,opt,name=src_port,json=srcPort,proto3" json:"src_port,omitempty"`
	DstPort   int32     `protobuf:"varint,5,opt,name=dst_port,json=dstPort,proto3" json:"dst_port,omitempty"`
	Protocol  int32     `protobuf:"varint,6,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *PermitListRule) Reset() {
	*x = PermitListRule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PermitListRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermitListRule) ProtoMessage() {}

func (x *PermitListRule) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermitListRule.ProtoReflect.Descriptor instead.
func (*PermitListRule) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{5}
}

func (x *PermitListRule) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PermitListRule) GetTag() []string {
	if x != nil {
		return x.Tag
	}
	return nil
}

func (x *PermitListRule) GetDirection() Direction {
	if x != nil {
		return x.Direction
	}
	return Direction_INBOUND
}

func (x *PermitListRule) GetSrcPort() int32 {
	if x != nil {
		return x.SrcPort
	}
	return 0
}

func (x *PermitListRule) GetDstPort() int32 {
	if x != nil {
		return x.DstPort
	}
	return 0
}

func (x *PermitListRule) GetProtocol() int32 {
	if x != nil {
		return x.Protocol
	}
	return 0
}

type PermitList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AssociatedResource string            `protobuf:"bytes,1,opt,name=associated_resource,json=associatedResource,proto3" json:"associated_resource,omitempty"`
	Rules              []*PermitListRule `protobuf:"bytes,2,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *PermitList) Reset() {
	*x = PermitList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PermitList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermitList) ProtoMessage() {}

func (x *PermitList) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermitList.ProtoReflect.Descriptor instead.
func (*PermitList) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{6}
}

func (x *PermitList) GetAssociatedResource() string {
	if x != nil {
		return x.AssociatedResource
	}
	return ""
}

func (x *PermitList) GetRules() []*PermitListRule {
	if x != nil {
		return x.Rules
	}
	return nil
}

var File_invisinets_proto protoreflect.FileDescriptor

var file_invisinets_proto_rawDesc = []byte{
	0x0a, 0x10, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0c, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62,
	0x22, 0x26, 0x0a, 0x14, 0x49, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x44, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x39, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x53, 0x70, 0x61, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x53, 0x70, 0x61,
	0x63, 0x65, 0x73, 0x22, 0x1c, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49,
	0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x5c, 0x0a, 0x13, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x53, 0x70, 0x61, 0x63, 0x65, 0x22,
	0xa2, 0x01, 0x0a, 0x0d, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x48, 0x0a, 0x10, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x44, 0x48, 0x00, 0x52, 0x0f, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x88, 0x01, 0x01, 0x42,
	0x13, 0x0a, 0x11, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x22, 0xbb, 0x01, 0x0a, 0x0e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x35, 0x0a, 0x09, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x69,
	0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x44, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x19, 0x0a, 0x08, 0x73, 0x72, 0x63, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x07, 0x73, 0x72, 0x63, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x64,
	0x73, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x64,
	0x73, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x22, 0x71, 0x0a, 0x0a, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x2f, 0x0a, 0x13, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x61,
	0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x32, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05,
	0x72, 0x75, 0x6c, 0x65, 0x73, 0x2a, 0x26, 0x0a, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e, 0x42, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x00, 0x12,
	0x0c, 0x0a, 0x08, 0x4f, 0x55, 0x54, 0x42, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x32, 0xa7, 0x03,
	0x0a, 0x0b, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x5c, 0x0a,
	0x14, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x53,
	0x70, 0x61, 0x63, 0x65, 0x73, 0x12, 0x22, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65,
	0x74, 0x73, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x44,
	0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x1e, 0x2e, 0x69, 0x6e, 0x76, 0x69,
	0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x53, 0x70, 0x61, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x52, 0x0a, 0x0e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x21, 0x2e,
	0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x1a, 0x1b, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e,
	0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x45, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x18, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x44, 0x1a, 0x18, 0x2e, 0x69, 0x6e, 0x76,
	0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x50, 0x65, 0x72,
	0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x18, 0x2e, 0x69,
	0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x6d,
	0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e,
	0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x50, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50,
	0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x18,
	0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73,
	0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4e, 0x65, 0x74, 0x53, 0x79, 0x73, 0x2f, 0x69, 0x6e, 0x76,
	0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x69, 0x6e, 0x76, 0x69,
	0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_invisinets_proto_rawDescOnce sync.Once
	file_invisinets_proto_rawDescData = file_invisinets_proto_rawDesc
)

func file_invisinets_proto_rawDescGZIP() []byte {
	file_invisinets_proto_rawDescOnce.Do(func() {
		file_invisinets_proto_rawDescData = protoimpl.X.CompressGZIP(file_invisinets_proto_rawDescData)
	})
	return file_invisinets_proto_rawDescData
}

var file_invisinets_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_invisinets_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_invisinets_proto_goTypes = []interface{}{
	(Direction)(0),               // 0: invisinetspb.Direction
	(*InvisinetsDeployment)(nil), // 1: invisinetspb.InvisinetsDeployment
	(*AddressSpaceList)(nil),     // 2: invisinetspb.AddressSpaceList
	(*ResourceID)(nil),           // 3: invisinetspb.ResourceID
	(*ResourceDescription)(nil),  // 4: invisinetspb.ResourceDescription
	(*BasicResponse)(nil),        // 5: invisinetspb.BasicResponse
	(*PermitListRule)(nil),       // 6: invisinetspb.PermitListRule
	(*PermitList)(nil),           // 7: invisinetspb.PermitList
}
var file_invisinets_proto_depIdxs = []int32{
	3, // 0: invisinetspb.BasicResponse.updated_resource:type_name -> invisinetspb.ResourceID
	0, // 1: invisinetspb.PermitListRule.direction:type_name -> invisinetspb.Direction
	6, // 2: invisinetspb.PermitList.rules:type_name -> invisinetspb.PermitListRule
	1, // 3: invisinetspb.CloudPlugin.GetUsedAddressSpaces:input_type -> invisinetspb.InvisinetsDeployment
	4, // 4: invisinetspb.CloudPlugin.CreateResource:input_type -> invisinetspb.ResourceDescription
	3, // 5: invisinetspb.CloudPlugin.GetPermitList:input_type -> invisinetspb.ResourceID
	7, // 6: invisinetspb.CloudPlugin.AddPermitListRules:input_type -> invisinetspb.PermitList
	7, // 7: invisinetspb.CloudPlugin.DeletePermitListRules:input_type -> invisinetspb.PermitList
	2, // 8: invisinetspb.CloudPlugin.GetUsedAddressSpaces:output_type -> invisinetspb.AddressSpaceList
	5, // 9: invisinetspb.CloudPlugin.CreateResource:output_type -> invisinetspb.BasicResponse
	7, // 10: invisinetspb.CloudPlugin.GetPermitList:output_type -> invisinetspb.PermitList
	5, // 11: invisinetspb.CloudPlugin.AddPermitListRules:output_type -> invisinetspb.BasicResponse
	5, // 12: invisinetspb.CloudPlugin.DeletePermitListRules:output_type -> invisinetspb.BasicResponse
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_invisinets_proto_init() }
func file_invisinets_proto_init() {
	if File_invisinets_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_invisinets_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvisinetsDeployment); i {
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
		file_invisinets_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddressSpaceList); i {
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
		file_invisinets_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceID); i {
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
		file_invisinets_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceDescription); i {
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
		file_invisinets_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicResponse); i {
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
		file_invisinets_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PermitListRule); i {
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
		file_invisinets_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PermitList); i {
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
	file_invisinets_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_invisinets_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_invisinets_proto_goTypes,
		DependencyIndexes: file_invisinets_proto_depIdxs,
		EnumInfos:         file_invisinets_proto_enumTypes,
		MessageInfos:      file_invisinets_proto_msgTypes,
	}.Build()
	File_invisinets_proto = out.File
	file_invisinets_proto_rawDesc = nil
	file_invisinets_proto_goTypes = nil
	file_invisinets_proto_depIdxs = nil
}

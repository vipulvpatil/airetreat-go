// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: protos/server.proto

package protos

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

type TestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Test string `protobuf:"bytes,1,opt,name=test,proto3" json:"test,omitempty"`
}

func (x *TestRequest) Reset() {
	*x = TestRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRequest) ProtoMessage() {}

func (x *TestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRequest.ProtoReflect.Descriptor instead.
func (*TestRequest) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{0}
}

func (x *TestRequest) GetTest() string {
	if x != nil {
		return x.Test
	}
	return ""
}

type TestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Test string `protobuf:"bytes,1,opt,name=test,proto3" json:"test,omitempty"`
}

func (x *TestResponse) Reset() {
	*x = TestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResponse) ProtoMessage() {}

func (x *TestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResponse.ProtoReflect.Descriptor instead.
func (*TestResponse) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{1}
}

func (x *TestResponse) GetTest() string {
	if x != nil {
		return x.Test
	}
	return ""
}

type InstagramAccount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InstagramAccount) Reset() {
	*x = InstagramAccount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstagramAccount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstagramAccount) ProtoMessage() {}

func (x *InstagramAccount) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstagramAccount.ProtoReflect.Descriptor instead.
func (*InstagramAccount) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{2}
}

type AddInstagramAccountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InstagramUserId       string `protobuf:"bytes,1,opt,name=instagramUserId,proto3" json:"instagramUserId,omitempty"`
	ShortLivedAccessToken string `protobuf:"bytes,2,opt,name=shortLivedAccessToken,proto3" json:"shortLivedAccessToken,omitempty"`
}

func (x *AddInstagramAccountRequest) Reset() {
	*x = AddInstagramAccountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddInstagramAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddInstagramAccountRequest) ProtoMessage() {}

func (x *AddInstagramAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddInstagramAccountRequest.ProtoReflect.Descriptor instead.
func (*AddInstagramAccountRequest) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{3}
}

func (x *AddInstagramAccountRequest) GetInstagramUserId() string {
	if x != nil {
		return x.InstagramUserId
	}
	return ""
}

func (x *AddInstagramAccountRequest) GetShortLivedAccessToken() string {
	if x != nil {
		return x.ShortLivedAccessToken
	}
	return ""
}

type AddInstagramAccountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddInstagramAccountResponse) Reset() {
	*x = AddInstagramAccountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddInstagramAccountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddInstagramAccountResponse) ProtoMessage() {}

func (x *AddInstagramAccountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddInstagramAccountResponse.ProtoReflect.Descriptor instead.
func (*AddInstagramAccountResponse) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{4}
}

type DeleteInstagramAccountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteInstagramAccountRequest) Reset() {
	*x = DeleteInstagramAccountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteInstagramAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteInstagramAccountRequest) ProtoMessage() {}

func (x *DeleteInstagramAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteInstagramAccountRequest.ProtoReflect.Descriptor instead.
func (*DeleteInstagramAccountRequest) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteInstagramAccountRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteInstagramAccountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteInstagramAccountResponse) Reset() {
	*x = DeleteInstagramAccountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteInstagramAccountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteInstagramAccountResponse) ProtoMessage() {}

func (x *DeleteInstagramAccountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteInstagramAccountResponse.ProtoReflect.Descriptor instead.
func (*DeleteInstagramAccountResponse) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{6}
}

type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{7}
}

type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_server_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_server_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_protos_server_proto_rawDescGZIP(), []int{8}
}

var File_protos_server_proto protoreflect.FileDescriptor

var file_protos_server_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x22, 0x21, 0x0a,
	0x0b, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x73, 0x74,
	0x22, 0x22, 0x0a, 0x0c, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x65, 0x73, 0x74, 0x22, 0x12, 0x0a, 0x10, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61,
	0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x7c, 0x0a, 0x1a, 0x41, 0x64, 0x64, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x0f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x67,
	0x72, 0x61, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x34, 0x0a, 0x15, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x76, 0x65, 0x64, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x15, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x76, 0x65, 0x64, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x1d, 0x0a, 0x1b, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2f, 0x0a, 0x1d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x20, 0x0a, 0x1e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x0e, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0f, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x90, 0x02, 0x0a, 0x0c, 0x53, 0x6f,
	0x63, 0x69, 0x61, 0x6c, 0x4d, 0x69, 0x6e, 0x65, 0x47, 0x6f, 0x12, 0x33, 0x0a, 0x04, 0x54, 0x65,
	0x73, 0x74, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x54, 0x65, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x60, 0x0a, 0x13, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e,
	0x41, 0x64, 0x64, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x69, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x67, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x25, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x67, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x4c, 0x0a, 0x12,
	0x53, 0x6f, 0x63, 0x69, 0x61, 0x6c, 0x4d, 0x69, 0x6e, 0x65, 0x47, 0x6f, 0x48, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x12, 0x36, 0x0a, 0x05, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x14, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x69, 0x70, 0x75, 0x6c, 0x76, 0x70,
	0x61, 0x74, 0x69, 0x6c, 0x2f, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x6c, 0x6d, 0x69, 0x6e, 0x65, 0x2d,
	0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_protos_server_proto_rawDescOnce sync.Once
	file_protos_server_proto_rawDescData = file_protos_server_proto_rawDesc
)

func file_protos_server_proto_rawDescGZIP() []byte {
	file_protos_server_proto_rawDescOnce.Do(func() {
		file_protos_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_server_proto_rawDescData)
	})
	return file_protos_server_proto_rawDescData
}

var file_protos_server_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_protos_server_proto_goTypes = []interface{}{
	(*TestRequest)(nil),                    // 0: protos.TestRequest
	(*TestResponse)(nil),                   // 1: protos.TestResponse
	(*InstagramAccount)(nil),               // 2: protos.InstagramAccount
	(*AddInstagramAccountRequest)(nil),     // 3: protos.AddInstagramAccountRequest
	(*AddInstagramAccountResponse)(nil),    // 4: protos.AddInstagramAccountResponse
	(*DeleteInstagramAccountRequest)(nil),  // 5: protos.DeleteInstagramAccountRequest
	(*DeleteInstagramAccountResponse)(nil), // 6: protos.DeleteInstagramAccountResponse
	(*CheckRequest)(nil),                   // 7: protos.CheckRequest
	(*CheckResponse)(nil),                  // 8: protos.CheckResponse
}
var file_protos_server_proto_depIdxs = []int32{
	0, // 0: protos.AiRetreatGo.Test:input_type -> protos.TestRequest
	3, // 1: protos.AiRetreatGo.AddInstagramAccount:input_type -> protos.AddInstagramAccountRequest
	5, // 2: protos.AiRetreatGo.DeleteInstagramAccount:input_type -> protos.DeleteInstagramAccountRequest
	7, // 3: protos.AiRetreatGoHealth.Check:input_type -> protos.CheckRequest
	1, // 4: protos.AiRetreatGo.Test:output_type -> protos.TestResponse
	4, // 5: protos.AiRetreatGo.AddInstagramAccount:output_type -> protos.AddInstagramAccountResponse
	6, // 6: protos.AiRetreatGo.DeleteInstagramAccount:output_type -> protos.DeleteInstagramAccountResponse
	8, // 7: protos.AiRetreatGoHealth.Check:output_type -> protos.CheckResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_server_proto_init() }
func file_protos_server_proto_init() {
	if File_protos_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRequest); i {
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
		file_protos_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestResponse); i {
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
		file_protos_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstagramAccount); i {
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
		file_protos_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddInstagramAccountRequest); i {
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
		file_protos_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddInstagramAccountResponse); i {
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
		file_protos_server_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteInstagramAccountRequest); i {
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
		file_protos_server_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteInstagramAccountResponse); i {
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
		file_protos_server_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRequest); i {
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
		file_protos_server_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse); i {
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
			RawDescriptor: file_protos_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_protos_server_proto_goTypes,
		DependencyIndexes: file_protos_server_proto_depIdxs,
		MessageInfos:      file_protos_server_proto_msgTypes,
	}.Build()
	File_protos_server_proto = out.File
	file_protos_server_proto_rawDesc = nil
	file_protos_server_proto_goTypes = nil
	file_protos_server_proto_depIdxs = nil
}

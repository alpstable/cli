// Copyright 2023 The Gidari CLI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.6
// source: web.proto

package web

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

type WriteType int32

const (
	// CSV writes the response as a CSV file.
	WriteType_CSV WriteType = 0
	// MongoDB writes the response to a MongoDB database.
	WriteType_MONGO WriteType = 1
)

// Enum value maps for WriteType.
var (
	WriteType_name = map[int32]string{
		0: "CSV",
		1: "MONGO",
	}
	WriteType_value = map[string]int32{
		"CSV":   0,
		"MONGO": 1,
	}
)

func (x WriteType) Enum() *WriteType {
	p := new(WriteType)
	*p = x
	return p
}

func (x WriteType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WriteType) Descriptor() protoreflect.EnumDescriptor {
	return file_web_proto_enumTypes[0].Descriptor()
}

func (WriteType) Type() protoreflect.EnumType {
	return &file_web_proto_enumTypes[0]
}

func (x WriteType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WriteType.Descriptor instead.
func (WriteType) EnumDescriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{0}
}

type BasicAuth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *BasicAuth) Reset() {
	*x = BasicAuth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicAuth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicAuth) ProtoMessage() {}

func (x *BasicAuth) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicAuth.ProtoReflect.Descriptor instead.
func (*BasicAuth) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{0}
}

func (x *BasicAuth) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *BasicAuth) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type CoinbaseAuth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key        string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Secret     string `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	Passphrase string `protobuf:"bytes,3,opt,name=passphrase,proto3" json:"passphrase,omitempty"`
}

func (x *CoinbaseAuth) Reset() {
	*x = CoinbaseAuth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoinbaseAuth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoinbaseAuth) ProtoMessage() {}

func (x *CoinbaseAuth) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoinbaseAuth.ProtoReflect.Descriptor instead.
func (*CoinbaseAuth) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{1}
}

func (x *CoinbaseAuth) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *CoinbaseAuth) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

func (x *CoinbaseAuth) GetPassphrase() string {
	if x != nil {
		return x.Passphrase
	}
	return ""
}

type Auth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Auth:
	//
	//	*Auth_Basic
	//	*Auth_Coinbase
	Auth isAuth_Auth `protobuf_oneof:"auth"`
}

func (x *Auth) Reset() {
	*x = Auth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Auth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Auth) ProtoMessage() {}

func (x *Auth) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Auth.ProtoReflect.Descriptor instead.
func (*Auth) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{2}
}

func (m *Auth) GetAuth() isAuth_Auth {
	if m != nil {
		return m.Auth
	}
	return nil
}

func (x *Auth) GetBasic() *BasicAuth {
	if x, ok := x.GetAuth().(*Auth_Basic); ok {
		return x.Basic
	}
	return nil
}

func (x *Auth) GetCoinbase() *CoinbaseAuth {
	if x, ok := x.GetAuth().(*Auth_Coinbase); ok {
		return x.Coinbase
	}
	return nil
}

type isAuth_Auth interface {
	isAuth_Auth()
}

type Auth_Basic struct {
	Basic *BasicAuth `protobuf:"bytes,1,opt,name=basic,proto3,oneof"`
}

type Auth_Coinbase struct {
	Coinbase *CoinbaseAuth `protobuf:"bytes,2,opt,name=coinbase,proto3,oneof"`
}

func (*Auth_Basic) isAuth_Auth() {}

func (*Auth_Coinbase) isAuth_Auth() {}

type Writer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type       WriteType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.WriteType" json:"type,omitempty"`
	Dir        string    `protobuf:"bytes,2,opt,name=dir,proto3" json:"dir,omitempty"`
	Database   string    `protobuf:"bytes,3,opt,name=database,proto3" json:"database,omitempty"`
	ConnString string    `protobuf:"bytes,4,opt,name=conn_string,json=connString,proto3" json:"conn_string,omitempty"`
}

func (x *Writer) Reset() {
	*x = Writer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Writer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Writer) ProtoMessage() {}

func (x *Writer) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Writer.ProtoReflect.Descriptor instead.
func (*Writer) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{3}
}

func (x *Writer) GetType() WriteType {
	if x != nil {
		return x.Type
	}
	return WriteType_CSV
}

func (x *Writer) GetDir() string {
	if x != nil {
		return x.Dir
	}
	return ""
}

func (x *Writer) GetDatabase() string {
	if x != nil {
		return x.Database
	}
	return ""
}

func (x *Writer) GetConnString() string {
	if x != nil {
		return x.ConnString
	}
	return ""
}

type CSVConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File string `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *CSVConfig) Reset() {
	*x = CSVConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CSVConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CSVConfig) ProtoMessage() {}

func (x *CSVConfig) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CSVConfig.ProtoReflect.Descriptor instead.
func (*CSVConfig) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{4}
}

func (x *CSVConfig) GetFile() string {
	if x != nil {
		return x.File
	}
	return ""
}

type MongoConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Collection string   `protobuf:"bytes,1,opt,name=collection,proto3" json:"collection,omitempty"`
	Index      []string `protobuf:"bytes,2,rep,name=index,proto3" json:"index,omitempty"`
}

func (x *MongoConfig) Reset() {
	*x = MongoConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MongoConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MongoConfig) ProtoMessage() {}

func (x *MongoConfig) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MongoConfig.ProtoReflect.Descriptor instead.
func (*MongoConfig) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{5}
}

func (x *MongoConfig) GetCollection() string {
	if x != nil {
		return x.Collection
	}
	return ""
}

func (x *MongoConfig) GetIndex() []string {
	if x != nil {
		return x.Index
	}
	return nil
}

type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method string `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Url    string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	// QueryParams are the query parameters to use for the HTTP request.
	QueryParams map[string]string `protobuf:"bytes,3,rep,name=query_params,json=queryParams,proto3" json:"query_params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Writers are the writers to use for the transporting the response
	// to local storage.
	Writers []*Writer `protobuf:"bytes,4,rep,name=writers,proto3" json:"writers,omitempty"`
	// Auth is the request-specific auth, if blank the default global auth
	// defined on the web.Request will be used.
	Auth *Auth `protobuf:"bytes,5,opt,name=auth,proto3" json:"auth,omitempty"`
	// CSV is the CSV configuration for the request.
	Csv *CSVConfig `protobuf:"bytes,6,opt,name=csv,proto3" json:"csv,omitempty"`
	// Mongo is the MongoDB configuration for the request.
	Mongo *MongoConfig `protobuf:"bytes,7,opt,name=mongo,proto3" json:"mongo,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{6}
}

func (x *WriteRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *WriteRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *WriteRequest) GetQueryParams() map[string]string {
	if x != nil {
		return x.QueryParams
	}
	return nil
}

func (x *WriteRequest) GetWriters() []*Writer {
	if x != nil {
		return x.Writers
	}
	return nil
}

func (x *WriteRequest) GetAuth() *Auth {
	if x != nil {
		return x.Auth
	}
	return nil
}

func (x *WriteRequest) GetCsv() *CSVConfig {
	if x != nil {
		return x.Csv
	}
	return nil
}

func (x *WriteRequest) GetMongo() *MongoConfig {
	if x != nil {
		return x.Mongo
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{7}
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requests []*WriteRequest `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	// Writers are the writers to use for the transporting the response
	// to local storage. This is the default for all write requests and
	// can be overridden on a per-write request basis.
	Writers []*Writer `protobuf:"bytes,2,rep,name=writers,proto3" json:"writers,omitempty"`
	// Auth is the global auth for all requests, it can be overridden by the
	// individual requests.
	Auth *Auth `protobuf:"bytes,3,opt,name=auth,proto3" json:"auth,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_web_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_web_proto_rawDescGZIP(), []int{8}
}

func (x *Request) GetRequests() []*WriteRequest {
	if x != nil {
		return x.Requests
	}
	return nil
}

func (x *Request) GetWriters() []*Writer {
	if x != nil {
		return x.Writers
	}
	return nil
}

func (x *Request) GetAuth() *Auth {
	if x != nil {
		return x.Auth
	}
	return nil
}

var File_web_proto protoreflect.FileDescriptor

var file_web_proto_rawDesc = []byte{
	0x0a, 0x09, 0x77, 0x65, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x43, 0x0a, 0x09, 0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x12,
	0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x58, 0x0a, 0x0c, 0x43, 0x6f, 0x69, 0x6e, 0x62,
	0x61, 0x73, 0x65, 0x41, 0x75, 0x74, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x73, 0x73, 0x70, 0x68, 0x72, 0x61, 0x73, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x73, 0x73, 0x70, 0x68, 0x72, 0x61, 0x73,
	0x65, 0x22, 0x6b, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x28, 0x0a, 0x05, 0x62, 0x61, 0x73,
	0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x41, 0x75, 0x74, 0x68, 0x48, 0x00, 0x52, 0x05, 0x62, 0x61,
	0x73, 0x69, 0x63, 0x12, 0x31, 0x0a, 0x08, 0x63, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f,
	0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x41, 0x75, 0x74, 0x68, 0x48, 0x00, 0x52, 0x08, 0x63, 0x6f,
	0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x22, 0x7d,
	0x0a, 0x06, 0x57, 0x72, 0x69, 0x74, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57,
	0x72, 0x69, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x64, 0x69, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x69, 0x72,
	0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x1f, 0x0a,
	0x09, 0x43, 0x53, 0x56, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x69,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x43,
	0x0a, 0x0b, 0x4d, 0x6f, 0x6e, 0x67, 0x6f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1e, 0x0a,
	0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a,
	0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x22, 0xd9, 0x02, 0x0a, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x47,
	0x0a, 0x0c, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x72, 0x69,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x27, 0x0a, 0x07, 0x77, 0x72, 0x69, 0x74, 0x65,
	0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x72, 0x52, 0x07, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x73,
	0x12, 0x1f, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x04, 0x61, 0x75, 0x74,
	0x68, 0x12, 0x22, 0x0a, 0x03, 0x63, 0x73, 0x76, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x53, 0x56, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x03, 0x63, 0x73, 0x76, 0x12, 0x28, 0x0a, 0x05, 0x6d, 0x6f, 0x6e, 0x67, 0x6f, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x6f, 0x6e,
	0x67, 0x6f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x05, 0x6d, 0x6f, 0x6e, 0x67, 0x6f, 0x1a,
	0x3e, 0x0a, 0x10, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x0a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x84, 0x01, 0x0a, 0x07,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2f, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x08,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x12, 0x27, 0x0a, 0x07, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x72, 0x52, 0x07, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72,
	0x73, 0x12, 0x1f, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x04, 0x61, 0x75,
	0x74, 0x68, 0x2a, 0x1f, 0x0a, 0x09, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x07, 0x0a, 0x03, 0x43, 0x53, 0x56, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x4f, 0x4e, 0x47,
	0x4f, 0x10, 0x01, 0x32, 0x31, 0x0a, 0x03, 0x57, 0x65, 0x62, 0x12, 0x2a, 0x0a, 0x05, 0x57, 0x72,
	0x69, 0x74, 0x65, 0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x3b, 0x77, 0x65, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_web_proto_rawDescOnce sync.Once
	file_web_proto_rawDescData = file_web_proto_rawDesc
)

func file_web_proto_rawDescGZIP() []byte {
	file_web_proto_rawDescOnce.Do(func() {
		file_web_proto_rawDescData = protoimpl.X.CompressGZIP(file_web_proto_rawDescData)
	})
	return file_web_proto_rawDescData
}

var file_web_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_web_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_web_proto_goTypes = []interface{}{
	(WriteType)(0),       // 0: proto.WriteType
	(*BasicAuth)(nil),    // 1: proto.BasicAuth
	(*CoinbaseAuth)(nil), // 2: proto.CoinbaseAuth
	(*Auth)(nil),         // 3: proto.Auth
	(*Writer)(nil),       // 4: proto.Writer
	(*CSVConfig)(nil),    // 5: proto.CSVConfig
	(*MongoConfig)(nil),  // 6: proto.MongoConfig
	(*WriteRequest)(nil), // 7: proto.WriteRequest
	(*Response)(nil),     // 8: proto.Response
	(*Request)(nil),      // 9: proto.Request
	nil,                  // 10: proto.WriteRequest.QueryParamsEntry
}
var file_web_proto_depIdxs = []int32{
	1,  // 0: proto.Auth.basic:type_name -> proto.BasicAuth
	2,  // 1: proto.Auth.coinbase:type_name -> proto.CoinbaseAuth
	0,  // 2: proto.Writer.type:type_name -> proto.WriteType
	10, // 3: proto.WriteRequest.query_params:type_name -> proto.WriteRequest.QueryParamsEntry
	4,  // 4: proto.WriteRequest.writers:type_name -> proto.Writer
	3,  // 5: proto.WriteRequest.auth:type_name -> proto.Auth
	5,  // 6: proto.WriteRequest.csv:type_name -> proto.CSVConfig
	6,  // 7: proto.WriteRequest.mongo:type_name -> proto.MongoConfig
	7,  // 8: proto.Request.requests:type_name -> proto.WriteRequest
	4,  // 9: proto.Request.writers:type_name -> proto.Writer
	3,  // 10: proto.Request.auth:type_name -> proto.Auth
	9,  // 11: proto.Web.Write:input_type -> proto.Request
	8,  // 12: proto.Web.Write:output_type -> proto.Response
	12, // [12:13] is the sub-list for method output_type
	11, // [11:12] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_web_proto_init() }
func file_web_proto_init() {
	if File_web_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_web_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicAuth); i {
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
		file_web_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CoinbaseAuth); i {
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
		file_web_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Auth); i {
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
		file_web_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Writer); i {
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
		file_web_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CSVConfig); i {
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
		file_web_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MongoConfig); i {
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
		file_web_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteRequest); i {
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
		file_web_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_web_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
	file_web_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Auth_Basic)(nil),
		(*Auth_Coinbase)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_web_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_web_proto_goTypes,
		DependencyIndexes: file_web_proto_depIdxs,
		EnumInfos:         file_web_proto_enumTypes,
		MessageInfos:      file_web_proto_msgTypes,
	}.Build()
	File_web_proto = out.File
	file_web_proto_rawDesc = nil
	file_web_proto_goTypes = nil
	file_web_proto_depIdxs = nil
}

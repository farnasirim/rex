// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: rex.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ReadRequest_File int32

const (
	// cannot start enum values from 1.
	ReadRequest_STDOUT ReadRequest_File = 0
	ReadRequest_STDERR ReadRequest_File = 1
)

// Enum value maps for ReadRequest_File.
var (
	ReadRequest_File_name = map[int32]string{
		0: "STDOUT",
		1: "STDERR",
	}
	ReadRequest_File_value = map[string]int32{
		"STDOUT": 0,
		"STDERR": 1,
	}
)

func (x ReadRequest_File) Enum() *ReadRequest_File {
	p := new(ReadRequest_File)
	*p = x
	return p
}

func (x ReadRequest_File) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ReadRequest_File) Descriptor() protoreflect.EnumDescriptor {
	return file_rex_proto_enumTypes[0].Descriptor()
}

func (ReadRequest_File) Type() protoreflect.EnumType {
	return &file_rex_proto_enumTypes[0]
}

func (x ReadRequest_File) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReadRequest_File.Descriptor instead.
func (ReadRequest_File) EnumDescriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{8, 0}
}

// ExecRequest specifies what binary needs to be Exec'd and how.
type ExecRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// path is the absolute or relative (to server dir) path to an executable
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// args is a list of command line args that will be passed to the
	// executable upon execution.
	Args []string `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
}

func (x *ExecRequest) Reset() {
	*x = ExecRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecRequest) ProtoMessage() {}

func (x *ExecRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecRequest.ProtoReflect.Descriptor instead.
func (*ExecRequest) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{0}
}

func (x *ExecRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ExecRequest) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

// ExecResponse embodies the identifier of the newly created process if the
// call to Exec had been successful.
type ExecResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ProcessUUID is the UUID of the process that is created when Exec returns
	// without error.
	ProcessUUID string `protobuf:"bytes,1,opt,name=ProcessUUID,proto3" json:"ProcessUUID,omitempty"` // TODO: validate UUID here
}

func (x *ExecResponse) Reset() {
	*x = ExecResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecResponse) ProtoMessage() {}

func (x *ExecResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecResponse.ProtoReflect.Descriptor instead.
func (*ExecResponse) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{1}
}

func (x *ExecResponse) GetProcessUUID() string {
	if x != nil {
		return x.ProcessUUID
	}
	return ""
}

// ProcessInfo is the summarized information about a particular process
type ProcessInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessUUID string               `protobuf:"bytes,1,opt,name=processUUID,proto3" json:"processUUID,omitempty"`
	Pid         int32                `protobuf:"varint,2,opt,name=pid,proto3" json:"pid,omitempty"`
	ExitCode    int32                `protobuf:"varint,3,opt,name=exitCode,proto3" json:"exitCode,omitempty"`
	Running     bool                 `protobuf:"varint,4,opt,name=running,proto3" json:"running,omitempty"`
	Path        string               `protobuf:"bytes,5,opt,name=path,proto3" json:"path,omitempty"`
	Args        []string             `protobuf:"bytes,6,rep,name=args,proto3" json:"args,omitempty"`
	OwnerUUID   string               `protobuf:"bytes,7,opt,name=ownerUUID,proto3" json:"ownerUUID,omitempty"`
	Create      *timestamp.Timestamp `protobuf:"bytes,8,opt,name=create,proto3" json:"create,omitempty"`
	Exit        *timestamp.Timestamp `protobuf:"bytes,9,opt,name=exit,proto3" json:"exit,omitempty"`
}

func (x *ProcessInfo) Reset() {
	*x = ProcessInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessInfo) ProtoMessage() {}

func (x *ProcessInfo) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessInfo.ProtoReflect.Descriptor instead.
func (*ProcessInfo) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{2}
}

func (x *ProcessInfo) GetProcessUUID() string {
	if x != nil {
		return x.ProcessUUID
	}
	return ""
}

func (x *ProcessInfo) GetPid() int32 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *ProcessInfo) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

func (x *ProcessInfo) GetRunning() bool {
	if x != nil {
		return x.Running
	}
	return false
}

func (x *ProcessInfo) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ProcessInfo) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *ProcessInfo) GetOwnerUUID() string {
	if x != nil {
		return x.OwnerUUID
	}
	return ""
}

func (x *ProcessInfo) GetCreate() *timestamp.Timestamp {
	if x != nil {
		return x.Create
	}
	return nil
}

func (x *ProcessInfo) GetExit() *timestamp.Timestamp {
	if x != nil {
		return x.Exit
	}
	return nil
}

// ProcessInfoList embodies a list of ProcessInfo messages
type ProcessInfoList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Processes []*ProcessInfo `protobuf:"bytes,1,rep,name=processes,proto3" json:"processes,omitempty"`
}

func (x *ProcessInfoList) Reset() {
	*x = ProcessInfoList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessInfoList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessInfoList) ProtoMessage() {}

func (x *ProcessInfoList) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessInfoList.ProtoReflect.Descriptor instead.
func (*ProcessInfoList) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{3}
}

func (x *ProcessInfoList) GetProcesses() []*ProcessInfo {
	if x != nil {
		return x.Processes
	}
	return nil
}

type ListProcessInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListProcessInfoRequest) Reset() {
	*x = ListProcessInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListProcessInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProcessInfoRequest) ProtoMessage() {}

func (x *ListProcessInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListProcessInfoRequest.ProtoReflect.Descriptor instead.
func (*ListProcessInfoRequest) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{4}
}

type GetProcessInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessUUID string `protobuf:"bytes,1,opt,name=processUUID,proto3" json:"processUUID,omitempty"`
}

func (x *GetProcessInfoRequest) Reset() {
	*x = GetProcessInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProcessInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProcessInfoRequest) ProtoMessage() {}

func (x *GetProcessInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProcessInfoRequest.ProtoReflect.Descriptor instead.
func (*GetProcessInfoRequest) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{5}
}

func (x *GetProcessInfoRequest) GetProcessUUID() string {
	if x != nil {
		return x.ProcessUUID
	}
	return ""
}

type KillRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessUUID string `protobuf:"bytes,1,opt,name=processUUID,proto3" json:"processUUID,omitempty"`
	Signal      int32  `protobuf:"varint,2,opt,name=signal,proto3" json:"signal,omitempty"`
}

func (x *KillRequest) Reset() {
	*x = KillRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KillRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KillRequest) ProtoMessage() {}

func (x *KillRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KillRequest.ProtoReflect.Descriptor instead.
func (*KillRequest) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{6}
}

func (x *KillRequest) GetProcessUUID() string {
	if x != nil {
		return x.ProcessUUID
	}
	return ""
}

func (x *KillRequest) GetSignal() int32 {
	if x != nil {
		return x.Signal
	}
	return 0
}

type KillResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *KillResponse) Reset() {
	*x = KillResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KillResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KillResponse) ProtoMessage() {}

func (x *KillResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KillResponse.ProtoReflect.Descriptor instead.
func (*KillResponse) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{7}
}

type ReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessUUID string           `protobuf:"bytes,1,opt,name=processUUID,proto3" json:"processUUID,omitempty"`
	Target      ReadRequest_File `protobuf:"varint,2,opt,name=target,proto3,enum=ReadRequest_File" json:"target,omitempty"`
}

func (x *ReadRequest) Reset() {
	*x = ReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRequest) ProtoMessage() {}

func (x *ReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRequest.ProtoReflect.Descriptor instead.
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{8}
}

func (x *ReadRequest) GetProcessUUID() string {
	if x != nil {
		return x.ProcessUUID
	}
	return ""
}

func (x *ReadRequest) GetTarget() ReadRequest_File {
	if x != nil {
		return x.Target
	}
	return ReadRequest_STDOUT
}

type ReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ReadResponse) Reset() {
	*x = ReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rex_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadResponse) ProtoMessage() {}

func (x *ReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rex_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadResponse.ProtoReflect.Descriptor instead.
func (*ReadResponse) Descriptor() ([]byte, []int) {
	return file_rex_proto_rawDescGZIP(), []int{9}
}

func (x *ReadResponse) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

var File_rex_proto protoreflect.FileDescriptor

var file_rex_proto_rawDesc = []byte{
	0x0a, 0x09, 0x72, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x0b,
	0x45, 0x78, 0x65, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12,
	0x12, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x61,
	0x72, 0x67, 0x73, 0x22, 0x30, 0x0a, 0x0c, 0x45, 0x78, 0x65, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x55, 0x55,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x55, 0x55, 0x49, 0x44, 0x22, 0xa1, 0x02, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x55, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x55, 0x55, 0x49, 0x44, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78, 0x69,
	0x74, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x65, 0x78, 0x69,
	0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x77, 0x6e, 0x65, 0x72,
	0x55, 0x55, 0x49, 0x44, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x55, 0x55, 0x49, 0x44, 0x12, 0x32, 0x0a, 0x06, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x06, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x65, 0x78, 0x69,
	0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x04, 0x65, 0x78, 0x69, 0x74, 0x22, 0x3d, 0x0a, 0x0f, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x09,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0x18, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x39, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x55, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x55, 0x55, 0x49, 0x44, 0x22, 0x47, 0x0a,
	0x0b, 0x4b, 0x69, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x55, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x55, 0x55, 0x49, 0x44, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x22, 0x0e, 0x0a, 0x0c, 0x4b, 0x69, 0x6c, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x7a, 0x0a, 0x0b, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x55, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x55, 0x55, 0x49, 0x44, 0x12, 0x29, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x22, 0x1e, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54,
	0x44, 0x4f, 0x55, 0x54, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54, 0x44, 0x45, 0x52, 0x52,
	0x10, 0x01, 0x22, 0x28, 0x0a, 0x0c, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0xf4, 0x01, 0x0a,
	0x03, 0x52, 0x65, 0x78, 0x12, 0x25, 0x0a, 0x04, 0x45, 0x78, 0x65, 0x63, 0x12, 0x0c, 0x2e, 0x45,
	0x78, 0x65, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x45, 0x78, 0x65,
	0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0f, 0x4c,
	0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x17,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e,
	0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49,
	0x6e, 0x66, 0x6f, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x04, 0x4b, 0x69, 0x6c, 0x6c, 0x12, 0x0c, 0x2e,
	0x4b, 0x69, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x4b, 0x69,
	0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x04,
	0x52, 0x65, 0x61, 0x64, 0x12, 0x0c, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x66, 0x61, 0x72, 0x6e, 0x61, 0x73, 0x69, 0x72, 0x69, 0x6d, 0x2f, 0x72, 0x65, 0x78,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rex_proto_rawDescOnce sync.Once
	file_rex_proto_rawDescData = file_rex_proto_rawDesc
)

func file_rex_proto_rawDescGZIP() []byte {
	file_rex_proto_rawDescOnce.Do(func() {
		file_rex_proto_rawDescData = protoimpl.X.CompressGZIP(file_rex_proto_rawDescData)
	})
	return file_rex_proto_rawDescData
}

var file_rex_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rex_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_rex_proto_goTypes = []interface{}{
	(ReadRequest_File)(0),          // 0: ReadRequest.File
	(*ExecRequest)(nil),            // 1: ExecRequest
	(*ExecResponse)(nil),           // 2: ExecResponse
	(*ProcessInfo)(nil),            // 3: ProcessInfo
	(*ProcessInfoList)(nil),        // 4: ProcessInfoList
	(*ListProcessInfoRequest)(nil), // 5: ListProcessInfoRequest
	(*GetProcessInfoRequest)(nil),  // 6: GetProcessInfoRequest
	(*KillRequest)(nil),            // 7: KillRequest
	(*KillResponse)(nil),           // 8: KillResponse
	(*ReadRequest)(nil),            // 9: ReadRequest
	(*ReadResponse)(nil),           // 10: ReadResponse
	(*timestamp.Timestamp)(nil),    // 11: google.protobuf.Timestamp
}
var file_rex_proto_depIdxs = []int32{
	11, // 0: ProcessInfo.create:type_name -> google.protobuf.Timestamp
	11, // 1: ProcessInfo.exit:type_name -> google.protobuf.Timestamp
	3,  // 2: ProcessInfoList.processes:type_name -> ProcessInfo
	0,  // 3: ReadRequest.target:type_name -> ReadRequest.File
	1,  // 4: Rex.Exec:input_type -> ExecRequest
	5,  // 5: Rex.ListProcessInfo:input_type -> ListProcessInfoRequest
	6,  // 6: Rex.GetProcessInfo:input_type -> GetProcessInfoRequest
	7,  // 7: Rex.Kill:input_type -> KillRequest
	9,  // 8: Rex.Read:input_type -> ReadRequest
	2,  // 9: Rex.Exec:output_type -> ExecResponse
	4,  // 10: Rex.ListProcessInfo:output_type -> ProcessInfoList
	3,  // 11: Rex.GetProcessInfo:output_type -> ProcessInfo
	8,  // 12: Rex.Kill:output_type -> KillResponse
	10, // 13: Rex.Read:output_type -> ReadResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_rex_proto_init() }
func file_rex_proto_init() {
	if File_rex_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rex_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecRequest); i {
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
		file_rex_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecResponse); i {
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
		file_rex_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessInfo); i {
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
		file_rex_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessInfoList); i {
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
		file_rex_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListProcessInfoRequest); i {
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
		file_rex_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProcessInfoRequest); i {
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
		file_rex_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KillRequest); i {
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
		file_rex_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KillResponse); i {
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
		file_rex_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadRequest); i {
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
		file_rex_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadResponse); i {
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
			RawDescriptor: file_rex_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rex_proto_goTypes,
		DependencyIndexes: file_rex_proto_depIdxs,
		EnumInfos:         file_rex_proto_enumTypes,
		MessageInfos:      file_rex_proto_msgTypes,
	}.Build()
	File_rex_proto = out.File
	file_rex_proto_rawDesc = nil
	file_rex_proto_goTypes = nil
	file_rex_proto_depIdxs = nil
}

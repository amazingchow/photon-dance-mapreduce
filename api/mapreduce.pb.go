// Code generated by protoc-gen-go. DO NOT EDIT.
// source: amazingchow/mapreduce/pb/mapreduce.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type WorkerStatus int32

const (
	WorkerStatus_WORKER_STATUS_ALIVE     WorkerStatus = 0
	WorkerStatus_WORKER_STATUS_UNHEALTHY WorkerStatus = 1
)

var WorkerStatus_name = map[int32]string{
	0: "WORKER_STATUS_ALIVE",
	1: "WORKER_STATUS_UNHEALTHY",
}

var WorkerStatus_value = map[string]int32{
	"WORKER_STATUS_ALIVE":     0,
	"WORKER_STATUS_UNHEALTHY": 1,
}

func (x WorkerStatus) String() string {
	return proto.EnumName(WorkerStatus_name, int32(x))
}

func (WorkerStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{0}
}

type TaskType int32

const (
	TaskType_TASK_TYPE_MAP    TaskType = 0
	TaskType_TASK_TYPE_REDUCE TaskType = 1
)

var TaskType_name = map[int32]string{
	0: "TASK_TYPE_MAP",
	1: "TASK_TYPE_REDUCE",
}

var TaskType_value = map[string]int32{
	"TASK_TYPE_MAP":    0,
	"TASK_TYPE_REDUCE": 1,
}

func (x TaskType) String() string {
	return proto.EnumName(TaskType_name, int32(x))
}

func (TaskType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{1}
}

type TaskStatus int32

const (
	TaskStatus_TASK_STATUS_UNALLOTED TaskStatus = 0
	TaskStatus_TASK_STATUS_ALLOTED   TaskStatus = 1
	TaskStatus_TASK_STATUS_DONE      TaskStatus = 2
)

var TaskStatus_name = map[int32]string{
	0: "TASK_STATUS_UNALLOTED",
	1: "TASK_STATUS_ALLOTED",
	2: "TASK_STATUS_DONE",
}

var TaskStatus_value = map[string]int32{
	"TASK_STATUS_UNALLOTED": 0,
	"TASK_STATUS_ALLOTED":   1,
	"TASK_STATUS_DONE":      2,
}

func (x TaskStatus) String() string {
	return proto.EnumName(TaskStatus_name, int32(x))
}

func (TaskStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{2}
}

type IntercomType int32

const (
	IntercomType_INTERCOM_TYPE_ASK_TASK           IntercomType = 0
	IntercomType_INTERCOM_TYPE_FINISH_MAP_TASK    IntercomType = 1
	IntercomType_INTERCOM_TYPE_FINISH_REDUCE_TASK IntercomType = 2
	IntercomType_INTERCOM_TYPE_SEND_INTER_FILE    IntercomType = 3
)

var IntercomType_name = map[int32]string{
	0: "INTERCOM_TYPE_ASK_TASK",
	1: "INTERCOM_TYPE_FINISH_MAP_TASK",
	2: "INTERCOM_TYPE_FINISH_REDUCE_TASK",
	3: "INTERCOM_TYPE_SEND_INTER_FILE",
}

var IntercomType_value = map[string]int32{
	"INTERCOM_TYPE_ASK_TASK":           0,
	"INTERCOM_TYPE_FINISH_MAP_TASK":    1,
	"INTERCOM_TYPE_FINISH_REDUCE_TASK": 2,
	"INTERCOM_TYPE_SEND_INTER_FILE":    3,
}

func (x IntercomType) String() string {
	return proto.EnumName(IntercomType_name, int32(x))
}

func (IntercomType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{3}
}

type Task struct {
	Inputs               []string `protobuf:"bytes,1,rep,name=inputs,proto3" json:"inputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{0}
}

func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (m *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(m, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetInputs() []string {
	if m != nil {
		return m.Inputs
	}
	return nil
}

type Worker struct {
	Id                   string       `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Status               WorkerStatus `protobuf:"varint,5,opt,name=status,proto3,enum=amazingchow.mapreduce.WorkerStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Worker) Reset()         { *m = Worker{} }
func (m *Worker) String() string { return proto.CompactTextString(m) }
func (*Worker) ProtoMessage()    {}
func (*Worker) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{1}
}

func (m *Worker) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Worker.Unmarshal(m, b)
}
func (m *Worker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Worker.Marshal(b, m, deterministic)
}
func (m *Worker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Worker.Merge(m, src)
}
func (m *Worker) XXX_Size() int {
	return xxx_messageInfo_Worker.Size(m)
}
func (m *Worker) XXX_DiscardUnknown() {
	xxx_messageInfo_Worker.DiscardUnknown(m)
}

var xxx_messageInfo_Worker proto.InternalMessageInfo

func (m *Worker) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Worker) GetStatus() WorkerStatus {
	if m != nil {
		return m.Status
	}
	return WorkerStatus_WORKER_STATUS_ALIVE
}

// ---------- request + response ----------
type AddTaskRequest struct {
	Task                 *Task    `protobuf:"bytes,1,opt,name=task,proto3" json:"task,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddTaskRequest) Reset()         { *m = AddTaskRequest{} }
func (m *AddTaskRequest) String() string { return proto.CompactTextString(m) }
func (*AddTaskRequest) ProtoMessage()    {}
func (*AddTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{2}
}

func (m *AddTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddTaskRequest.Unmarshal(m, b)
}
func (m *AddTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddTaskRequest.Marshal(b, m, deterministic)
}
func (m *AddTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddTaskRequest.Merge(m, src)
}
func (m *AddTaskRequest) XXX_Size() int {
	return xxx_messageInfo_AddTaskRequest.Size(m)
}
func (m *AddTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddTaskRequest proto.InternalMessageInfo

func (m *AddTaskRequest) GetTask() *Task {
	if m != nil {
		return m.Task
	}
	return nil
}

type AddTaskResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddTaskResponse) Reset()         { *m = AddTaskResponse{} }
func (m *AddTaskResponse) String() string { return proto.CompactTextString(m) }
func (*AddTaskResponse) ProtoMessage()    {}
func (*AddTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{3}
}

func (m *AddTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddTaskResponse.Unmarshal(m, b)
}
func (m *AddTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddTaskResponse.Marshal(b, m, deterministic)
}
func (m *AddTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddTaskResponse.Merge(m, src)
}
func (m *AddTaskResponse) XXX_Size() int {
	return xxx_messageInfo_AddTaskResponse.Size(m)
}
func (m *AddTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddTaskResponse proto.InternalMessageInfo

type ListWorkersRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListWorkersRequest) Reset()         { *m = ListWorkersRequest{} }
func (m *ListWorkersRequest) String() string { return proto.CompactTextString(m) }
func (*ListWorkersRequest) ProtoMessage()    {}
func (*ListWorkersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{4}
}

func (m *ListWorkersRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListWorkersRequest.Unmarshal(m, b)
}
func (m *ListWorkersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListWorkersRequest.Marshal(b, m, deterministic)
}
func (m *ListWorkersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListWorkersRequest.Merge(m, src)
}
func (m *ListWorkersRequest) XXX_Size() int {
	return xxx_messageInfo_ListWorkersRequest.Size(m)
}
func (m *ListWorkersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListWorkersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListWorkersRequest proto.InternalMessageInfo

type ListWorkersResponse struct {
	Workers              []*Worker `protobuf:"bytes,1,rep,name=workers,proto3" json:"workers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ListWorkersResponse) Reset()         { *m = ListWorkersResponse{} }
func (m *ListWorkersResponse) String() string { return proto.CompactTextString(m) }
func (*ListWorkersResponse) ProtoMessage()    {}
func (*ListWorkersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{5}
}

func (m *ListWorkersResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListWorkersResponse.Unmarshal(m, b)
}
func (m *ListWorkersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListWorkersResponse.Marshal(b, m, deterministic)
}
func (m *ListWorkersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListWorkersResponse.Merge(m, src)
}
func (m *ListWorkersResponse) XXX_Size() int {
	return xxx_messageInfo_ListWorkersResponse.Size(m)
}
func (m *ListWorkersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListWorkersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListWorkersResponse proto.InternalMessageInfo

func (m *ListWorkersResponse) GetWorkers() []*Worker {
	if m != nil {
		return m.Workers
	}
	return nil
}

type IntercomRequest struct {
	MsgType              IntercomType `protobuf:"varint,1,opt,name=MsgType,proto3,enum=amazingchow.mapreduce.IntercomType" json:"MsgType,omitempty"`
	MsgContent           string       `protobuf:"bytes,2,opt,name=MsgContent,proto3" json:"MsgContent,omitempty"`
	Extra                string       `protobuf:"bytes,3,opt,name=Extra,proto3" json:"Extra,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *IntercomRequest) Reset()         { *m = IntercomRequest{} }
func (m *IntercomRequest) String() string { return proto.CompactTextString(m) }
func (*IntercomRequest) ProtoMessage()    {}
func (*IntercomRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{6}
}

func (m *IntercomRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IntercomRequest.Unmarshal(m, b)
}
func (m *IntercomRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IntercomRequest.Marshal(b, m, deterministic)
}
func (m *IntercomRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IntercomRequest.Merge(m, src)
}
func (m *IntercomRequest) XXX_Size() int {
	return xxx_messageInfo_IntercomRequest.Size(m)
}
func (m *IntercomRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IntercomRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IntercomRequest proto.InternalMessageInfo

func (m *IntercomRequest) GetMsgType() IntercomType {
	if m != nil {
		return m.MsgType
	}
	return IntercomType_INTERCOM_TYPE_ASK_TASK
}

func (m *IntercomRequest) GetMsgContent() string {
	if m != nil {
		return m.MsgContent
	}
	return ""
}

func (m *IntercomRequest) GetExtra() string {
	if m != nil {
		return m.Extra
	}
	return ""
}

type IntercomResponse struct {
	TaskType             TaskType `protobuf:"varint,1,opt,name=TaskType,proto3,enum=amazingchow.mapreduce.TaskType" json:"TaskType,omitempty"`
	File                 string   `protobuf:"bytes,2,opt,name=File,proto3" json:"File,omitempty"`
	NReduce              int32    `protobuf:"varint,3,opt,name=NReduce,proto3" json:"NReduce,omitempty"`
	MapTaskAllocated     int32    `protobuf:"varint,4,opt,name=MapTaskAllocated,proto3" json:"MapTaskAllocated,omitempty"`
	ReduceTaskAllocated  int32    `protobuf:"varint,5,opt,name=ReduceTaskAllocated,proto3" json:"ReduceTaskAllocated,omitempty"`
	ReduceFiles          []string `protobuf:"bytes,6,rep,name=ReduceFiles,proto3" json:"ReduceFiles,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IntercomResponse) Reset()         { *m = IntercomResponse{} }
func (m *IntercomResponse) String() string { return proto.CompactTextString(m) }
func (*IntercomResponse) ProtoMessage()    {}
func (*IntercomResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5cf3aee7a1e5ec0, []int{7}
}

func (m *IntercomResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IntercomResponse.Unmarshal(m, b)
}
func (m *IntercomResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IntercomResponse.Marshal(b, m, deterministic)
}
func (m *IntercomResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IntercomResponse.Merge(m, src)
}
func (m *IntercomResponse) XXX_Size() int {
	return xxx_messageInfo_IntercomResponse.Size(m)
}
func (m *IntercomResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IntercomResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IntercomResponse proto.InternalMessageInfo

func (m *IntercomResponse) GetTaskType() TaskType {
	if m != nil {
		return m.TaskType
	}
	return TaskType_TASK_TYPE_MAP
}

func (m *IntercomResponse) GetFile() string {
	if m != nil {
		return m.File
	}
	return ""
}

func (m *IntercomResponse) GetNReduce() int32 {
	if m != nil {
		return m.NReduce
	}
	return 0
}

func (m *IntercomResponse) GetMapTaskAllocated() int32 {
	if m != nil {
		return m.MapTaskAllocated
	}
	return 0
}

func (m *IntercomResponse) GetReduceTaskAllocated() int32 {
	if m != nil {
		return m.ReduceTaskAllocated
	}
	return 0
}

func (m *IntercomResponse) GetReduceFiles() []string {
	if m != nil {
		return m.ReduceFiles
	}
	return nil
}

func init() {
	proto.RegisterEnum("amazingchow.mapreduce.WorkerStatus", WorkerStatus_name, WorkerStatus_value)
	proto.RegisterEnum("amazingchow.mapreduce.TaskType", TaskType_name, TaskType_value)
	proto.RegisterEnum("amazingchow.mapreduce.TaskStatus", TaskStatus_name, TaskStatus_value)
	proto.RegisterEnum("amazingchow.mapreduce.IntercomType", IntercomType_name, IntercomType_value)
	proto.RegisterType((*Task)(nil), "amazingchow.mapreduce.Task")
	proto.RegisterType((*Worker)(nil), "amazingchow.mapreduce.Worker")
	proto.RegisterType((*AddTaskRequest)(nil), "amazingchow.mapreduce.AddTaskRequest")
	proto.RegisterType((*AddTaskResponse)(nil), "amazingchow.mapreduce.AddTaskResponse")
	proto.RegisterType((*ListWorkersRequest)(nil), "amazingchow.mapreduce.ListWorkersRequest")
	proto.RegisterType((*ListWorkersResponse)(nil), "amazingchow.mapreduce.ListWorkersResponse")
	proto.RegisterType((*IntercomRequest)(nil), "amazingchow.mapreduce.IntercomRequest")
	proto.RegisterType((*IntercomResponse)(nil), "amazingchow.mapreduce.IntercomResponse")
}

func init() {
	proto.RegisterFile("amazingchow/mapreduce/pb/mapreduce.proto", fileDescriptor_e5cf3aee7a1e5ec0)
}

var fileDescriptor_e5cf3aee7a1e5ec0 = []byte{
	// 712 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x4b, 0x6f, 0xda, 0x5a,
	0x10, 0xc6, 0x84, 0x47, 0x32, 0xe4, 0xe1, 0x0c, 0x79, 0x38, 0xe4, 0x26, 0x97, 0xeb, 0xdb, 0xa6,
	0x94, 0x05, 0xb4, 0x64, 0x51, 0xa9, 0x51, 0x17, 0x2e, 0x38, 0x0a, 0x0a, 0x8f, 0xc8, 0x36, 0x89,
	0x52, 0xa9, 0x42, 0x0e, 0x58, 0xd4, 0x0a, 0xd8, 0x2e, 0xe7, 0x90, 0x34, 0x5d, 0x76, 0xd1, 0x3f,
	0x50, 0xa9, 0x7f, 0xac, 0xfb, 0xae, 0xfa, 0x37, 0x2a, 0x55, 0x3e, 0xb6, 0xc1, 0x28, 0x90, 0x74,
	0xe7, 0x33, 0xf3, 0x7d, 0xf3, 0x7d, 0x33, 0x3e, 0x73, 0x20, 0xa7, 0x0f, 0xf4, 0xcf, 0xa6, 0xd5,
	0xeb, 0x7c, 0xb0, 0x6f, 0x8b, 0x03, 0xdd, 0x19, 0x1a, 0xdd, 0x51, 0xc7, 0x28, 0x3a, 0x57, 0x93,
	0x43, 0xc1, 0x19, 0xda, 0xd4, 0xc6, 0xcd, 0x10, 0xb2, 0x30, 0x4e, 0x66, 0xfe, 0xe9, 0xd9, 0x76,
	0xaf, 0x6f, 0x14, 0x75, 0xc7, 0x2c, 0xea, 0x96, 0x65, 0x53, 0x9d, 0x9a, 0xb6, 0x45, 0x3c, 0x92,
	0xb8, 0x0f, 0x31, 0x4d, 0x27, 0xd7, 0xb8, 0x05, 0x09, 0xd3, 0x72, 0x46, 0x94, 0x08, 0x5c, 0x76,
	0x21, 0xb7, 0xa4, 0xf8, 0x27, 0xb1, 0x05, 0x89, 0x0b, 0x7b, 0x78, 0x6d, 0x0c, 0x71, 0x15, 0xa2,
	0x66, 0x57, 0x88, 0x66, 0xb9, 0xdc, 0x92, 0x12, 0x35, 0xbb, 0x78, 0x04, 0x09, 0x42, 0x75, 0x3a,
	0x22, 0x42, 0x3c, 0xcb, 0xe5, 0x56, 0x4b, 0xff, 0x17, 0x66, 0xea, 0x17, 0x3c, 0xba, 0xca, 0xa0,
	0x8a, 0x4f, 0x11, 0x25, 0x58, 0x95, 0xba, 0x5d, 0x57, 0x59, 0x31, 0x3e, 0x8e, 0x0c, 0x42, 0xb1,
	0x08, 0x31, 0xaa, 0x93, 0x6b, 0x81, 0xcb, 0x72, 0xb9, 0x54, 0x69, 0x77, 0x4e, 0x31, 0xc6, 0x60,
	0x40, 0x71, 0x1d, 0xd6, 0xc6, 0x25, 0x88, 0x63, 0x5b, 0xc4, 0x10, 0x37, 0x00, 0x6b, 0x26, 0xa1,
	0x9e, 0x22, 0xf1, 0x2b, 0x8b, 0x0d, 0x48, 0x4f, 0x45, 0x3d, 0x30, 0xbe, 0x82, 0xe4, 0xad, 0x17,
	0x62, 0x2d, 0xa7, 0x4a, 0x7b, 0x0f, 0x36, 0xa0, 0x04, 0x68, 0xf1, 0x2b, 0x07, 0x6b, 0x55, 0x8b,
	0x1a, 0xc3, 0x8e, 0x3d, 0x08, 0xdc, 0xbf, 0x81, 0x64, 0x9d, 0xf4, 0xb4, 0x3b, 0xc7, 0x60, 0x0d,
	0xcc, 0x9f, 0x46, 0x40, 0x74, 0xa1, 0x4a, 0xc0, 0xc1, 0x7d, 0x80, 0x3a, 0xe9, 0x95, 0x6d, 0x8b,
	0x1a, 0x16, 0xf5, 0x67, 0x1c, 0x8a, 0xe0, 0x06, 0xc4, 0xe5, 0x4f, 0x74, 0xa8, 0x0b, 0x0b, 0x2c,
	0xe5, 0x1d, 0xc4, 0xdf, 0x1c, 0xf0, 0x13, 0x23, 0x7e, 0x5b, 0x47, 0xb0, 0xe8, 0xce, 0x24, 0x64,
	0xe5, 0xdf, 0x07, 0x66, 0xc9, 0x6c, 0x8c, 0x09, 0x88, 0x10, 0x3b, 0x36, 0xfb, 0x86, 0xef, 0x80,
	0x7d, 0xa3, 0x00, 0xc9, 0x86, 0xc2, 0x18, 0x4c, 0x3d, 0xae, 0x04, 0x47, 0xcc, 0x03, 0x5f, 0xd7,
	0x1d, 0x97, 0x2c, 0xf5, 0xfb, 0x76, 0x47, 0xa7, 0x46, 0x57, 0x88, 0x31, 0xc8, 0xbd, 0x38, 0xbe,
	0x80, 0xb4, 0xc7, 0x9a, 0x86, 0xc7, 0x19, 0x7c, 0x56, 0x0a, 0xb3, 0x90, 0xf2, 0xc2, 0xae, 0x0b,
	0x22, 0x24, 0xd8, 0xb5, 0x0c, 0x87, 0xf2, 0x15, 0x58, 0x0e, 0x5f, 0x2e, 0xdc, 0x86, 0xf4, 0x45,
	0x53, 0x39, 0x95, 0x95, 0xb6, 0xaa, 0x49, 0x5a, 0x4b, 0x6d, 0x4b, 0xb5, 0xea, 0xb9, 0xcc, 0x47,
	0x70, 0x17, 0xb6, 0xa7, 0x13, 0xad, 0xc6, 0x89, 0x2c, 0xd5, 0xb4, 0x93, 0x4b, 0x9e, 0xcb, 0x1f,
	0x4e, 0x06, 0x86, 0xeb, 0xb0, 0xa2, 0x49, 0xea, 0x69, 0x5b, 0xbb, 0x3c, 0x93, 0xdb, 0x75, 0xe9,
	0x8c, 0x8f, 0xe0, 0x06, 0xf0, 0x93, 0x90, 0x22, 0x57, 0x5a, 0x65, 0x99, 0xe7, 0xf2, 0xe7, 0x00,
	0x2e, 0xc9, 0x17, 0xde, 0x81, 0x4d, 0x86, 0x19, 0x57, 0x97, 0x6a, 0xb5, 0xa6, 0x26, 0x57, 0xf8,
	0x88, 0xeb, 0x29, 0x9c, 0x0a, 0x12, 0xdc, 0xb8, 0xae, 0x9f, 0xa8, 0x34, 0x1b, 0x32, 0x1f, 0xcd,
	0x7f, 0xe7, 0x60, 0x39, 0x7c, 0x45, 0x30, 0x03, 0x5b, 0xd5, 0x86, 0x26, 0x2b, 0xe5, 0x66, 0xdd,
	0xb3, 0xc0, 0xbc, 0x48, 0xea, 0x29, 0x1f, 0xc1, 0xff, 0x60, 0x6f, 0x3a, 0x77, 0x5c, 0x6d, 0x54,
	0xd5, 0x13, 0xd7, 0xb8, 0x07, 0xe1, 0xf0, 0x09, 0x64, 0x67, 0x42, 0xbc, 0x46, 0x3c, 0x54, 0xf4,
	0x7e, 0x21, 0x55, 0x6e, 0x54, 0xda, 0x2c, 0xd4, 0x3e, 0xae, 0xd6, 0x64, 0x7e, 0xa1, 0xf4, 0x33,
	0x0a, 0xe9, 0xba, 0xee, 0x78, 0xe3, 0x57, 0xce, 0xca, 0xaa, 0x31, 0xbc, 0x31, 0x3b, 0x06, 0x0e,
	0x20, 0xe9, 0x6f, 0x21, 0x3e, 0x9d, 0x73, 0xcf, 0xa6, 0x17, 0x3d, 0x73, 0xf0, 0x18, 0xcc, 0x5f,
	0xe6, 0xf4, 0x97, 0x1f, 0xbf, 0xbe, 0x45, 0x57, 0xc4, 0xc5, 0xe2, 0xcd, 0xcb, 0xa2, 0xbb, 0xf1,
	0xaf, 0xb9, 0x3c, 0xde, 0x41, 0x2a, 0xb4, 0xcb, 0xf8, 0x7c, 0x4e, 0xad, 0xfb, 0xaf, 0x40, 0x26,
	0xff, 0x37, 0xd0, 0x69, 0x69, 0x4c, 0xb9, 0xd2, 0xfe, 0xda, 0xe3, 0x7b, 0x58, 0x0c, 0xfe, 0x0c,
	0x1e, 0x3c, 0xb2, 0xdd, 0x81, 0xe8, 0xb3, 0x47, 0x71, 0xbe, 0x62, 0xe4, 0xed, 0xee, 0xbb, 0x9d,
	0xd9, 0x2f, 0xbd, 0xee, 0x98, 0x57, 0x09, 0xf6, 0x58, 0x1f, 0xfe, 0x09, 0x00, 0x00, 0xff, 0xff,
	0xa2, 0x39, 0x6f, 0x99, 0x0d, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MapReduceRPCServiceClient is the client API for MapReduceRPCService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MapReduceRPCServiceClient interface {
	AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskResponse, error)
	ListWorkers(ctx context.Context, in *ListWorkersRequest, opts ...grpc.CallOption) (*ListWorkersResponse, error)
	Intercom(ctx context.Context, in *IntercomRequest, opts ...grpc.CallOption) (*IntercomResponse, error)
}

type mapReduceRPCServiceClient struct {
	cc *grpc.ClientConn
}

func NewMapReduceRPCServiceClient(cc *grpc.ClientConn) MapReduceRPCServiceClient {
	return &mapReduceRPCServiceClient{cc}
}

func (c *mapReduceRPCServiceClient) AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskResponse, error) {
	out := new(AddTaskResponse)
	err := c.cc.Invoke(ctx, "/amazingchow.mapreduce.MapReduceRPCService/AddTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mapReduceRPCServiceClient) ListWorkers(ctx context.Context, in *ListWorkersRequest, opts ...grpc.CallOption) (*ListWorkersResponse, error) {
	out := new(ListWorkersResponse)
	err := c.cc.Invoke(ctx, "/amazingchow.mapreduce.MapReduceRPCService/ListWorkers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mapReduceRPCServiceClient) Intercom(ctx context.Context, in *IntercomRequest, opts ...grpc.CallOption) (*IntercomResponse, error) {
	out := new(IntercomResponse)
	err := c.cc.Invoke(ctx, "/amazingchow.mapreduce.MapReduceRPCService/Intercom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MapReduceRPCServiceServer is the server API for MapReduceRPCService service.
type MapReduceRPCServiceServer interface {
	AddTask(context.Context, *AddTaskRequest) (*AddTaskResponse, error)
	ListWorkers(context.Context, *ListWorkersRequest) (*ListWorkersResponse, error)
	Intercom(context.Context, *IntercomRequest) (*IntercomResponse, error)
}

// UnimplementedMapReduceRPCServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMapReduceRPCServiceServer struct {
}

func (*UnimplementedMapReduceRPCServiceServer) AddTask(ctx context.Context, req *AddTaskRequest) (*AddTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTask not implemented")
}
func (*UnimplementedMapReduceRPCServiceServer) ListWorkers(ctx context.Context, req *ListWorkersRequest) (*ListWorkersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWorkers not implemented")
}
func (*UnimplementedMapReduceRPCServiceServer) Intercom(ctx context.Context, req *IntercomRequest) (*IntercomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Intercom not implemented")
}

func RegisterMapReduceRPCServiceServer(s *grpc.Server, srv MapReduceRPCServiceServer) {
	s.RegisterService(&_MapReduceRPCService_serviceDesc, srv)
}

func _MapReduceRPCService_AddTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapReduceRPCServiceServer).AddTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/amazingchow.mapreduce.MapReduceRPCService/AddTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapReduceRPCServiceServer).AddTask(ctx, req.(*AddTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MapReduceRPCService_ListWorkers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWorkersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapReduceRPCServiceServer).ListWorkers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/amazingchow.mapreduce.MapReduceRPCService/ListWorkers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapReduceRPCServiceServer).ListWorkers(ctx, req.(*ListWorkersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MapReduceRPCService_Intercom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IntercomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapReduceRPCServiceServer).Intercom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/amazingchow.mapreduce.MapReduceRPCService/Intercom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapReduceRPCServiceServer).Intercom(ctx, req.(*IntercomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MapReduceRPCService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "amazingchow.mapreduce.MapReduceRPCService",
	HandlerType: (*MapReduceRPCServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddTask",
			Handler:    _MapReduceRPCService_AddTask_Handler,
		},
		{
			MethodName: "ListWorkers",
			Handler:    _MapReduceRPCService_ListWorkers_Handler,
		},
		{
			MethodName: "Intercom",
			Handler:    _MapReduceRPCService_Intercom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "amazingchow/mapreduce/pb/mapreduce.proto",
}

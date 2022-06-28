// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: tasks/tasks.proto

package challenges

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TasksGrpcClient is the client API for TasksGrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TasksGrpcClient interface {
	RetrieveAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*TaskList, error)
	FilterTasks(ctx context.Context, in *FilterRequest, opts ...grpc.CallOption) (*TaskList, error)
	RetrieveTaskByID(ctx context.Context, in *TaskID, opts ...grpc.CallOption) (*ExistingTask, error)
	Create(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*TaskID, error)
	Update(ctx context.Context, in *UpdateTaskRequest, opts ...grpc.CallOption) (*Empty, error)
}

type tasksGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewTasksGrpcClient(cc grpc.ClientConnInterface) TasksGrpcClient {
	return &tasksGrpcClient{cc}
}

func (c *tasksGrpcClient) RetrieveAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*TaskList, error) {
	out := new(TaskList)
	err := c.cc.Invoke(ctx, "/tasks.TasksGrpc/RetrieveAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksGrpcClient) FilterTasks(ctx context.Context, in *FilterRequest, opts ...grpc.CallOption) (*TaskList, error) {
	out := new(TaskList)
	err := c.cc.Invoke(ctx, "/tasks.TasksGrpc/FilterTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksGrpcClient) RetrieveTaskByID(ctx context.Context, in *TaskID, opts ...grpc.CallOption) (*ExistingTask, error) {
	out := new(ExistingTask)
	err := c.cc.Invoke(ctx, "/tasks.TasksGrpc/RetrieveTaskByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksGrpcClient) Create(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*TaskID, error) {
	out := new(TaskID)
	err := c.cc.Invoke(ctx, "/tasks.TasksGrpc/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tasksGrpcClient) Update(ctx context.Context, in *UpdateTaskRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/tasks.TasksGrpc/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TasksGrpcServer is the server API for TasksGrpc service.
// All implementations must embed UnimplementedTasksGrpcServer
// for forward compatibility
type TasksGrpcServer interface {
	RetrieveAll(context.Context, *Empty) (*TaskList, error)
	FilterTasks(context.Context, *FilterRequest) (*TaskList, error)
	RetrieveTaskByID(context.Context, *TaskID) (*ExistingTask, error)
	Create(context.Context, *CreateTaskRequest) (*TaskID, error)
	Update(context.Context, *UpdateTaskRequest) (*Empty, error)
	mustEmbedUnimplementedTasksGrpcServer()
}

// UnimplementedTasksGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedTasksGrpcServer struct {
}

func (UnimplementedTasksGrpcServer) RetrieveAll(context.Context, *Empty) (*TaskList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveAll not implemented")
}
func (UnimplementedTasksGrpcServer) FilterTasks(context.Context, *FilterRequest) (*TaskList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterTasks not implemented")
}
func (UnimplementedTasksGrpcServer) RetrieveTaskByID(context.Context, *TaskID) (*ExistingTask, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveTaskByID not implemented")
}
func (UnimplementedTasksGrpcServer) Create(context.Context, *CreateTaskRequest) (*TaskID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedTasksGrpcServer) Update(context.Context, *UpdateTaskRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedTasksGrpcServer) mustEmbedUnimplementedTasksGrpcServer() {}

// UnsafeTasksGrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TasksGrpcServer will
// result in compilation errors.
type UnsafeTasksGrpcServer interface {
	mustEmbedUnimplementedTasksGrpcServer()
}

func RegisterTasksGrpcServer(s grpc.ServiceRegistrar, srv TasksGrpcServer) {
	s.RegisterService(&TasksGrpc_ServiceDesc, srv)
}

func _TasksGrpc_RetrieveAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksGrpcServer).RetrieveAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.TasksGrpc/RetrieveAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksGrpcServer).RetrieveAll(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TasksGrpc_FilterTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksGrpcServer).FilterTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.TasksGrpc/FilterTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksGrpcServer).FilterTasks(ctx, req.(*FilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TasksGrpc_RetrieveTaskByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksGrpcServer).RetrieveTaskByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.TasksGrpc/RetrieveTaskByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksGrpcServer).RetrieveTaskByID(ctx, req.(*TaskID))
	}
	return interceptor(ctx, in, info, handler)
}

func _TasksGrpc_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksGrpcServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.TasksGrpc/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksGrpcServer).Create(ctx, req.(*CreateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TasksGrpc_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TasksGrpcServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tasks.TasksGrpc/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TasksGrpcServer).Update(ctx, req.(*UpdateTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TasksGrpc_ServiceDesc is the grpc.ServiceDesc for TasksGrpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TasksGrpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tasks.TasksGrpc",
	HandlerType: (*TasksGrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RetrieveAll",
			Handler:    _TasksGrpc_RetrieveAll_Handler,
		},
		{
			MethodName: "FilterTasks",
			Handler:    _TasksGrpc_FilterTasks_Handler,
		},
		{
			MethodName: "RetrieveTaskByID",
			Handler:    _TasksGrpc_RetrieveTaskByID_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _TasksGrpc_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _TasksGrpc_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tasks/tasks.proto",
}

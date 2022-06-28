package proto

import (
	"context"
	"example/challenges/internal"
	"example/challenges/internal/mux"
	"fmt"
	"net"
	"strings"

	pb "example/challenges/internal/proto/tasks"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TasksGrpcServer struct {
	pb.UnimplementedTasksGrpcServer
	table internal.TaskTable
}

func NewTasksGrpcServer(address string) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen to %v, %v", address, err)
	}

	table, err := mux.SelectDBImpl()

	if err != nil {
		return nil, lis, fmt.Errorf("failed to listen to select DB impl, %v", err)
	}

	srv := TasksGrpcServer{table: table}

	s := grpc.NewServer()
	pb.RegisterTasksGrpcServer(s, &srv)

	return s, lis, nil
}

func (s *TasksGrpcServer) RetrieveAll(ctx context.Context, in *pb.Empty) (*pb.TaskList, error) {
	tasks, err := s.table.RetrieveAll()

	if err != nil {
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return getProtoTaskList(tasks), nil
}

func (s *TasksGrpcServer) FilterTasks(ctx context.Context, in *pb.FilterRequest) (*pb.TaskList, error) {
	completed := in.GetCompleted()
	tasks, err := s.table.Filter(completed)

	if err != nil {
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return getProtoTaskList(tasks), nil
}

func (s *TasksGrpcServer) RetrieveTaskByID(ctx context.Context, in *pb.TaskID) (*pb.ExistingTask, error) {
	id := in.GetId()
	t, err := s.table.FindByID(id)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, status.Newf(codes.InvalidArgument, err.Error()).Err()
		}
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return &pb.ExistingTask{ID: t.ID, Name: t.Name, Completed: t.Completed}, nil
}

func (s *TasksGrpcServer) Create(ctx context.Context, in *pb.NewTask) (*pb.TaskID, error) {
	task := internal.Task{Name: in.GetName(), Completed: in.GetCompleted()}
	id, err := s.table.Create(task)

	if err != nil {
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return &pb.TaskID{Id: id}, nil
}

func (s *TasksGrpcServer) Update(ctx context.Context, in *pb.ExistingTask) (*pb.Empty, error) {
	task := internal.Task{ID: in.GetID(), Name: in.GetName(), Completed: in.GetCompleted()}
	_, err := s.table.Update(task)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, status.Newf(codes.InvalidArgument, err.Error()).Err()
		}
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return &pb.Empty{}, nil
}

func getProtoTaskList(tasks []internal.Task) *pb.TaskList {
	result := []*pb.ExistingTask{}

	for _, t := range tasks {
		result = append(result, &pb.ExistingTask{ID: t.ID, Name: t.Name, Completed: t.Completed})
	}

	return &pb.TaskList{Tasks: result}
}
